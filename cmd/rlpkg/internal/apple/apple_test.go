package apple

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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

func selfSigned(t *testing.T) *Signer {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(7), Subject: pkix.Name{CommonName: "Test", OrganizationalUnit: []string{"TEAM1"}},
		NotBefore: time.Now(), NotAfter: time.Now().Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature}
	der, _ := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	cert, _ := x509.ParseCertificate(der)
	return &Signer{Key: key, Chain: []*x509.Certificate{cert}}
}

func TestSignCMSVerifiesWithOpenSSL(t *testing.T) {
	s := selfSigned(t)
	cd := bytes.Repeat([]byte("cd"), 300)
	cms, err := SignCMS(cd, s.Key, s.Chain, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := exec.LookPath("openssl"); err != nil {
		t.Skip("openssl not found")
	}
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "cms.der"), cms, 0o644)
	os.WriteFile(filepath.Join(dir, "content"), cd, 0o644)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: s.Chain[0].Raw})
	os.WriteFile(filepath.Join(dir, "cert.pem"), certPEM, 0o644)
	out, err := exec.Command("openssl", "cms", "-verify", "-inform", "DER", "-in", filepath.Join(dir, "cms.der"),
		"-content", filepath.Join(dir, "content"), "-CAfile", filepath.Join(dir, "cert.pem"), "-out", os.DevNull).CombinedOutput()
	if err != nil {
		t.Fatalf("openssl cms -verify: %v\n%s", err, out)
	}
	out, _ = exec.Command("openssl", "cms", "-inform", "DER", "-in", filepath.Join(dir, "cms.der"), "-cmsout", "-print").CombinedOutput()
	for _, want := range []string{"1.2.840.113635.100.9.1", "1.2.840.113635.100.9.2", "messageDigest", "signingTime"} {
		if !bytes.Contains(out, []byte(want)) {
			t.Errorf("missing %s in CMS dump", want)
		}
	}
}

func TestPlistDER(t *testing.T) {
	v, err := ParsePlist([]byte(`<?xml version="1.0"?><plist version="1.0"><dict><key>b</key><true/><key>a</key><string>x</string><key>c</key><array><string>y</string></array></dict></plist>`))
	if err != nil {
		t.Fatal(err)
	}
	der, err := PlistDER(v)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x70, 0x1f, 0x02, 0x01, 0x01, 0xb0, 0x1a,
		0x30, 0x06, 0x0c, 0x01, 'a', 0x0c, 0x01, 'x',
		0x30, 0x06, 0x0c, 0x01, 'b', 0x01, 0x01, 0xff,
		0x30, 0x08, 0x0c, 0x01, 'c', 0x30, 0x03, 0x0c, 0x01, 'y'}
	if !bytes.Equal(der, want) {
		t.Errorf("DER = % x\nwant  % x", der, want)
	}
}

func TestCodeResourcesMacOS(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "Resources"), 0o755)
	os.MkdirAll(filepath.Join(root, "MacOS"), 0o755)
	os.WriteFile(filepath.Join(root, "Info.plist"), []byte("plist"), 0o644)
	os.WriteFile(filepath.Join(root, "Resources", "app.icns"), []byte("icns"), 0o644)
	os.WriteFile(filepath.Join(root, "MacOS", "app"), []byte("exe"), 0o755)
	os.MkdirAll(filepath.Join(root, "Frameworks"), 0o755)
	os.WriteFile(filepath.Join(root, "Frameworks", "libx.dylib"), []byte("dylib"), 0o755)
	data, err := CodeResources(root, true, "MacOS/app", map[string]Nested{"Frameworks/libx.dylib": {CDHash: bytes.Repeat([]byte{1}, 20), Requirement: "cdhash H\"01\""}})
	if err != nil {
		t.Fatal(err)
	}
	v, err := ParsePlist(data)
	if err != nil {
		t.Fatal(err)
	}
	m := v.(map[string]any)
	files2 := m["files2"].(map[string]any)
	if _, ok := files2["Resources/app.icns"]; !ok {
		t.Error("Resources/app.icns missing from files2")
	}
	if _, ok := files2["Info.plist"]; ok {
		t.Error("Info.plist must be omitted")
	}
	if _, ok := files2["MacOS/app"]; ok {
		t.Error("MacOS/ is a nested rule, the executable must not be sealed as a resource")
	}
	if n, ok := files2["Frameworks/libx.dylib"].(map[string]any); !ok || n["requirement"] != "cdhash H\"01\"" {
		t.Errorf("nested entry = %v", files2["Frameworks/libx.dylib"])
	}
	if _, ok := m["files"].(map[string]any)["Resources/app.icns"]; !ok {
		t.Error("files (v1) missing the resource")
	}
}
