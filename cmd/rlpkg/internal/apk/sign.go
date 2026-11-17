package apk

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/cms"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/keys"
)

func DebugIdentity() (*keys.Identity, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "rlpkg", "debug.pem")
	if _, err := os.Stat(path); err == nil {
		return keys.Load(path)
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 64))
	if err != nil {
		return nil, err
	}
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "Android Debug", Organization: []string{"Android"}, Country: []string{"US"}},
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().AddDate(30, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	pem.Encode(&buf, &pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
	pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: der})
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr, "rlpkg: created debug signing key", path)
	return keys.Load(path)
}

func sign(id *keys.Identity, data []byte) ([]byte, error) {
	sum := sha256.Sum256(data)
	return id.Key.Sign(rand.Reader, sum[:], crypto.SHA256)
}

func isRSA(id *keys.Identity) bool {
	_, ok := id.Key.(*rsa.PrivateKey)
	return ok
}

func SignV1(entries []*Entry, id *keys.Identity) ([]*Entry, error) {
	return signJAR(entries, id, "X-Android-APK-Signed: 2\r\n")
}

// SignBundle JAR-signs an Android App Bundle, the only scheme bundletool and Play accept for it.
func SignBundle(entries []*Entry, id *keys.Identity) ([]*Entry, error) {
	return signJAR(entries, id, "")
}

func signJAR(entries []*Entry, id *keys.Identity, sfExtra string) ([]*Entry, error) {
	var manifest bytes.Buffer
	manifest.WriteString("Manifest-Version: 1.0\r\nCreated-By: rlpkg\r\n\r\n")
	var sf bytes.Buffer
	sf.WriteString("Signature-Version: 1.0\r\nCreated-By: rlpkg\r\n" + sfExtra)
	var sections bytes.Buffer
	for _, e := range entries {
		content, err := e.Content()
		if err != nil {
			return nil, err
		}
		sum := sha256.Sum256(content)
		section := "Name: " + e.Name + "\r\nSHA-256-Digest: " + base64.StdEncoding.EncodeToString(sum[:]) + "\r\n\r\n"
		manifest.WriteString(section)
		ssum := sha256.Sum256([]byte(section))
		sections.WriteString("Name: " + e.Name + "\r\nSHA-256-Digest: " + base64.StdEncoding.EncodeToString(ssum[:]) + "\r\n\r\n")
	}
	msum := sha256.Sum256(manifest.Bytes())
	sf.WriteString("SHA-256-Digest-Manifest: " + base64.StdEncoding.EncodeToString(msum[:]) + "\r\n\r\n")
	sf.Write(sections.Bytes())

	p7, err := cms.Sign(cms.Options{Key: id.Key, Chain: []*x509.Certificate{id.Cert}, Signed: sf.Bytes()})
	if err != nil {
		return nil, err
	}

	sigName := "META-INF/CERT.EC"
	if isRSA(id) {
		sigName = "META-INF/CERT.RSA"
	}
	var out []*Entry
	for _, f := range []struct {
		name string
		data []byte
	}{{"META-INF/MANIFEST.MF", manifest.Bytes()}, {"META-INF/CERT.SF", sf.Bytes()}, {sigName, p7}} {
		e, err := NewEntry(f.name, f.data, 8, 0)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return append(out, entries...), nil
}

const (
	sigV2BlockID    = 0x7109871a
	sigRSAPKCS1SHA2 = 0x0103
	sigECDSASHA2    = 0x0201
)

func SignV2(apk []byte, id *keys.Identity) ([]byte, error) {
	if len(apk) < 22 {
		return nil, errors.New("apk too short")
	}
	eocdOff := len(apk) - 22
	if binary.LittleEndian.Uint32(apk[eocdOff:]) != 0x06054b50 {
		return nil, errors.New("apk has a zip comment")
	}
	cdOff := int(binary.LittleEndian.Uint32(apk[eocdOff+16:]))
	cdSize := int(binary.LittleEndian.Uint32(apk[eocdOff+12:]))
	if cdOff+cdSize != eocdOff {
		return nil, errors.New("unexpected central directory layout")
	}

	eocd := append([]byte(nil), apk[eocdOff:]...)
	binary.LittleEndian.PutUint32(eocd[16:], uint32(cdOff))
	digest := v2Digest(apk[:cdOff], apk[cdOff:eocdOff], eocd)

	algo := uint32(sigECDSASHA2)
	if isRSA(id) {
		algo = sigRSAPKCS1SHA2
	}

	var signed bytes.Buffer
	signed.Write(v2Prefixed(v2Prefixed(binary.LittleEndian.AppendUint32(nil, algo), v2Prefixed(digest))))
	signed.Write(v2Prefixed(v2Prefixed(id.Cert.Raw)))
	signed.Write(v2Prefixed(nil))

	sig, err := sign(id, signed.Bytes())
	if err != nil {
		return nil, err
	}
	pub, err := x509.MarshalPKIXPublicKey(id.Key.Public())
	if err != nil {
		return nil, err
	}

	var signer bytes.Buffer
	signer.Write(v2Prefixed(signed.Bytes()))
	signer.Write(v2Prefixed(v2Prefixed(binary.LittleEndian.AppendUint32(nil, algo), v2Prefixed(sig))))
	signer.Write(v2Prefixed(pub))
	value := v2Prefixed(v2Prefixed(signer.Bytes()))

	pair := binary.LittleEndian.AppendUint64(nil, uint64(4+len(value)))
	pair = binary.LittleEndian.AppendUint32(pair, sigV2BlockID)
	pair = append(pair, value...)
	size := uint64(len(pair) + 8 + 16)
	block := binary.LittleEndian.AppendUint64(nil, size)
	block = append(block, pair...)
	block = binary.LittleEndian.AppendUint64(block, size)
	block = append(block, "APK Sig Block 42"...)

	out := make([]byte, 0, len(apk)+len(block))
	out = append(out, apk[:cdOff]...)
	out = append(out, block...)
	out = append(out, apk[cdOff:eocdOff]...)
	binary.LittleEndian.PutUint32(eocd[16:], uint32(cdOff+len(block)))
	return append(out, eocd...), nil
}

func v2Prefixed(parts ...[]byte) []byte {
	n := 0
	for _, p := range parts {
		n += len(p)
	}
	out := binary.LittleEndian.AppendUint32(make([]byte, 0, 4+n), uint32(n))
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func v2Digest(sections ...[]byte) []byte {
	var chunks [][]byte
	for _, sec := range sections {
		for len(sec) > 0 {
			n := min(len(sec), 1<<20)
			h := sha256.New()
			h.Write([]byte{0xa5})
			h.Write(binary.LittleEndian.AppendUint32(nil, uint32(n)))
			h.Write(sec[:n])
			chunks = append(chunks, h.Sum(nil))
			sec = sec[n:]
		}
	}
	h := sha256.New()
	h.Write([]byte{0x5a})
	h.Write(binary.LittleEndian.AppendUint32(nil, uint32(len(chunks))))
	for _, c := range chunks {
		h.Write(c)
	}
	return h.Sum(nil)
}
