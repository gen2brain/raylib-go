package pgp

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

const PassphraseEnv = "RLPKG_GPG_PASSPHRASE"

type Signer func(data []byte) ([]byte, error)

func Load(path string) (Signer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	list, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(data))
	if err != nil {
		if list, err = openpgp.ReadKeyRing(bytes.NewReader(data)); err != nil {
			return nil, fmt.Errorf("%s: not an OpenPGP key: %w", path, err)
		}
	}
	if len(list) != 1 {
		return nil, fmt.Errorf("%s: needs exactly one key, has %d", path, len(list))
	}
	entity := list[0]
	if entity.PrivateKey == nil {
		return nil, fmt.Errorf("%s: no secret key, export it with gpg --export-secret-keys", path)
	}
	if err := entity.DecryptPrivateKeys([]byte(os.Getenv(PassphraseEnv))); err != nil {
		return nil, fmt.Errorf("%s: cannot unlock the key with the passphrase in %s: %w", path, PassphraseEnv, err)
	}
	return func(data []byte) ([]byte, error) {
		var sig bytes.Buffer
		err := openpgp.DetachSign(&sig, entity, bytes.NewReader(data), &packet.Config{DefaultHash: crypto.SHA256})
		return sig.Bytes(), err
	}, nil
}

func GPG(user string) Signer {
	return func(data []byte) ([]byte, error) {
		cmd := exec.Command("gpg", "--batch", "--local-user", user, "--digest-algo", "SHA256", "--detach-sign", "--output", "-")
		cmd.Stdin = bytes.NewReader(data)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		sig, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("gpg: %w: %s", err, bytes.TrimSpace(stderr.Bytes()))
		}
		return sig, nil
	}
}

func Armor(sig []byte) ([]byte, error) {
	var out bytes.Buffer
	w, err := armor.Encode(&out, "PGP SIGNATURE", nil)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(sig); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	out.WriteByte('\n')
	return out.Bytes(), nil
}

func IsRSA(sig []byte) (bool, error) {
	p, err := packet.Read(bytes.NewReader(sig))
	if err != nil {
		return false, err
	}
	s, ok := p.(*packet.Signature)
	if !ok {
		return false, errors.New("not an OpenPGP signature")
	}
	return s.PubKeyAlgo == packet.PubKeyAlgoRSA || s.PubKeyAlgo == packet.PubKeyAlgoRSASignOnly, nil
}
