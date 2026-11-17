package cms

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func testChain(t *testing.T, key crypto.Signer) []*x509.Certificate {
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(5), Subject: pkix.Name{CommonName: "Test"},
		NotBefore: time.Now(), NotAfter: time.Now().Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, key.Public(), key)
	if err != nil {
		t.Fatal(err)
	}
	cert, _ := x509.ParseCertificate(der)
	return []*x509.Certificate{cert}
}

func opensslVerify(t *testing.T, sig, content []byte, cert *x509.Certificate) {
	if _, err := exec.LookPath("openssl"); err != nil {
		t.Skip("openssl not found")
	}
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "sig.der"), sig, 0o644)
	os.WriteFile(filepath.Join(dir, "cert.pem"), pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}), 0o644)
	args := []string{"cms", "-verify", "-inform", "DER", "-in", filepath.Join(dir, "sig.der"), "-CAfile", filepath.Join(dir, "cert.pem"), "-out", os.DevNull}
	if content != nil {
		os.WriteFile(filepath.Join(dir, "content"), content, 0o644)
		args = append(args, "-content", filepath.Join(dir, "content"))
	}
	if out, err := exec.Command("openssl", args...).CombinedOutput(); err != nil {
		t.Fatalf("openssl %v: %v\n%s", args, err, out)
	}
}

func TestSignDetachedWithAttributes(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	chain := testChain(t, key)
	content := []byte("detached content")
	sum := sha256.Sum256(content)
	sig, err := Sign(Options{
		Key: key, Chain: chain,
		SignedAttrs: []Attribute{NewAttribute(OIDContentType, OIDData), NewAttribute(OIDMessageDigest, sum[:])},
	})
	if err != nil {
		t.Fatal(err)
	}
	opensslVerify(t, sig, content, chain[0])
}

func TestSignDetachedWithoutAttributes(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	chain := testChain(t, key)
	content := []byte("signed directly")
	sig, err := Sign(Options{Key: key, Chain: chain, Signed: content})
	if err != nil {
		t.Fatal(err)
	}
	opensslVerify(t, sig, content, chain[0])
}

func TestSignEmbedded(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	chain := testChain(t, key)
	content := Marshal([]byte("embedded octets"))
	sum := sha256.Sum256([]byte("embedded octets"))
	sig, err := Sign(Options{
		Key: key, Chain: chain, Content: content,
		SignedAttrs: []Attribute{NewAttribute(OIDContentType, OIDData), NewAttribute(OIDMessageDigest, sum[:])},
	})
	if err != nil {
		t.Fatal(err)
	}
	opensslVerify(t, sig, nil, chain[0])
}
