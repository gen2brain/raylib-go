package pgp

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

func writeKey(t *testing.T, entity *openpgp.Entity, armored bool) string {
	t.Helper()
	var buf bytes.Buffer
	if armored {
		w, err := armor.Encode(&buf, openpgp.PrivateKeyType, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := entity.SerializePrivateWithoutSigning(w, nil); err != nil {
			t.Fatal(err)
		}
		w.Close()
	} else if err := entity.SerializePrivateWithoutSigning(&buf, nil); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "key.asc")
	if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoad(t *testing.T) {
	for _, tc := range []struct {
		name    string
		algo    packet.PublicKeyAlgorithm
		armored bool
		pass    string
		rsa     bool
	}{
		{"rsa armored locked", packet.PubKeyAlgoRSA, true, "secret", true},
		{"eddsa binary", packet.PubKeyAlgoEdDSA, false, "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			entity, err := openpgp.NewEntity("Demo", "", "demo@example.com", &packet.Config{Algorithm: tc.algo, RSABits: 2048})
			if err != nil {
				t.Fatal(err)
			}
			if tc.pass != "" {
				if err := entity.EncryptPrivateKeys([]byte(tc.pass), nil); err != nil {
					t.Fatal(err)
				}
			}
			path := writeKey(t, entity, tc.armored)

			if tc.pass != "" {
				t.Setenv(PassphraseEnv, "wrong")
				if _, err := Load(path); err == nil {
					t.Fatal("wrong passphrase accepted")
				}
			}
			t.Setenv(PassphraseEnv, tc.pass)
			sign, err := Load(path)
			if err != nil {
				t.Fatal(err)
			}
			data := []byte("signed data")
			sig, err := sign(data)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := openpgp.CheckDetachedSignature(openpgp.EntityList{entity}, bytes.NewReader(data), bytes.NewReader(sig), nil); err != nil {
				t.Fatalf("binary signature: %v", err)
			}
			if rsa, err := IsRSA(sig); err != nil || rsa != tc.rsa {
				t.Errorf("IsRSA = %v, %v, want %v", rsa, err, tc.rsa)
			}
			armored, err := Armor(sig)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := openpgp.CheckArmoredDetachedSignature(openpgp.EntityList{entity}, bytes.NewReader(data), bytes.NewReader(armored), nil); err != nil {
				t.Fatalf("armored signature: %v", err)
			}
		})
	}
}

func TestLoadPublicKey(t *testing.T) {
	entity, err := openpgp.NewEntity("Demo", "", "demo@example.com", &packet.Config{Algorithm: packet.PubKeyAlgoEdDSA})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := entity.Serialize(&buf); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "pub.gpg")
	os.WriteFile(path, buf.Bytes(), 0o600)
	if _, err := Load(path); err == nil {
		t.Fatal("public key accepted")
	}
}
