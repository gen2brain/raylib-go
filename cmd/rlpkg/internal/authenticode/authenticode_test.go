package authenticode

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/binary"
	"encoding/pem"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func testSigner(t *testing.T) *Signer {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(3), Subject: pkix.Name{CommonName: "Test"}, NotBefore: time.Now(), NotAfter: time.Now().Add(time.Hour),
		KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning}}
	der, _ := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	cert, _ := x509.ParseCertificate(der)
	return &Signer{Key: key, Chain: []*x509.Certificate{cert}}
}

func buildExe(t *testing.T) []byte {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module x\n\ngo 1.21\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644)
	cmd := exec.Command("go", "build", "-o", "x.exe", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOOS=windows", "GOARCH=amd64", "CGO_ENABLED=0")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("go build windows: %v\n%s", err, out)
	}
	data, _ := os.ReadFile(filepath.Join(dir, "x.exe"))
	return data
}

func TestSign(t *testing.T) {
	pe := buildExe(t)
	signed, err := Sign(pe, testSigner(t))
	if err != nil {
		t.Fatal(err)
	}
	l, err := layout(signed)
	if err != nil {
		t.Fatal(err)
	}
	if l.certOff == 0 || l.certOff+l.certSize != len(signed) || l.certOff%8 != 0 {
		t.Fatalf("certificate table off=%d size=%d len=%d", l.certOff, l.certSize, len(signed))
	}
	if binary.LittleEndian.Uint16(signed[l.certOff+4:]) != 0x200 || binary.LittleEndian.Uint16(signed[l.certOff+6:]) != 2 {
		t.Error("WIN_CERTIFICATE header wrong")
	}
	if got := binary.LittleEndian.Uint32(signed[l.checksumOff:]); got != checksum(signed, l.checksumOff) {
		t.Errorf("checksum %#x not consistent", got)
	}

	h1, _ := imageHash(signed, l)
	l0, _ := layout(pe)
	h0, _ := imageHash(pe, l0)
	if string(h0) != string(h1) {
		t.Error("image hash changed by signing")
	}

	resigned, err := Sign(signed, testSigner(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(resigned) < len(signed)-64 || len(resigned) > len(signed)+64 {
		t.Errorf("re-signing did not replace the table: %d vs %d", len(resigned), len(signed))
	}

	if _, err := exec.LookPath("openssl"); err != nil {
		t.Skip("openssl not found")
	}
	dir := t.TempDir()
	p7 := signed[l.certOff+8 : l.certOff+int(binary.LittleEndian.Uint32(signed[l.certOff:]))]
	os.WriteFile(filepath.Join(dir, "sig.der"), p7, 0o644)
	var leaf *x509.Certificate
	var sd struct {
		ContentType asn1.ObjectIdentifier
		Content     asn1.RawValue `asn1:"tag:0,explicit"`
	}
	if _, err := asn1.Unmarshal(p7, &sd); err != nil {
		t.Fatal(err)
	}
	var inner struct {
		Version          int
		DigestAlgorithms asn1.RawValue
		ContentInfo      asn1.RawValue
		Certificates     asn1.RawValue `asn1:"tag:0,optional"`
	}
	if _, err := asn1.Unmarshal(sd.Content.Bytes, &inner); err != nil {
		t.Fatal(err)
	}
	if leaf, err = x509.ParseCertificate(inner.Certificates.Bytes); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(dir, "cert.pem"), pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: leaf.Raw}), 0o644)
	out, err := exec.Command("openssl", "smime", "-verify", "-inform", "DER", "-in", filepath.Join(dir, "sig.der"),
		"-CAfile", filepath.Join(dir, "cert.pem"), "-purpose", "any", "-out", filepath.Join(dir, "content.der")).CombinedOutput()
	if err != nil {
		t.Fatalf("openssl smime -verify: %v\n%s", err, out)
	}
	content, _ := os.ReadFile(filepath.Join(dir, "content.der"))
	if !bytes.Contains(content, h1) {
		t.Error("verified content does not carry the image hash")
	}
}
