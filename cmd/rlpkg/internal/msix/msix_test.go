package msix

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/xml"
	"io"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/authenticode"
)

const testManifest = `<?xml version="1.0" encoding="utf-8"?>
<Package xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10">
  <Identity Name="com.example.demo" Publisher="CN=Demo, O=&#34;Demo, Inc.&#34;, C=RS" Version="1.0.0.1" ProcessorArchitecture="x64"/>
</Package>
`

func testSigner(t *testing.T, subject pkix.Name) *authenticode.Signer {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(3), Subject: subject, NotBefore: time.Now(), NotAfter: time.Now().Add(time.Hour),
		KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning}}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return &authenticode.Signer{Key: key, Chain: []*x509.Certificate{cert}}
}

func testFiles() []File {
	exe := make([]byte, 3*blockSize+1234)
	for i := range exe {
		exe[i] = byte(i * i >> 3)
	}
	return []File{{"demo app.exe", exe}, {"Assets/StoreLogo.png", []byte("\x89PNG not really")}}
}

func readAll(t *testing.T, data []byte) map[string][]byte {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{}
	for _, f := range zr.File {
		r, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("%s: %v", f.Name, err)
		}
		files[f.Name] = body
	}
	return files
}

func TestWrite(t *testing.T) {
	data, err := Write([]byte(testManifest), testFiles(), nil)
	if err != nil {
		t.Fatal(err)
	}
	files := readAll(t, data)
	for _, name := range []string{"demo%20app.exe", "Assets/StoreLogo.png", manifestName, blockMapName, contentTypesName} {
		if files[name] == nil {
			t.Errorf("missing %s", name)
		}
	}
	if files[signatureName] != nil || bytes.Contains(files[contentTypesName], []byte(signatureName)) {
		t.Error("unsigned package mentions a signature")
	}

	var blockMap struct {
		File []struct {
			Name  string `xml:"Name,attr"`
			Size  int    `xml:"Size,attr"`
			Block []struct {
				Hash string `xml:"Hash,attr"`
				Size int    `xml:"Size,attr"`
			}
		}
	}
	if err := xml.Unmarshal(files[blockMapName], &blockMap); err != nil {
		t.Fatal(err)
	}
	if len(blockMap.File) != 3 || blockMap.File[0].Name != "demo app.exe" || blockMap.File[1].Name != `Assets\StoreLogo.png` {
		t.Fatalf("block map files: %+v", blockMap.File)
	}
	exe := testFiles()[0].Data
	if got := blockMap.File[0]; got.Size != len(exe) || len(got.Block) != 4 {
		t.Fatalf("exe: size %d, %d blocks", got.Size, len(got.Block))
	}
	for i, b := range blockMap.File[0].Block {
		sum := sha256.Sum256(exe[i*blockSize : min((i+1)*blockSize, len(exe))])
		if b.Hash != base64.StdEncoding.EncodeToString(sum[:]) || b.Size == 0 {
			t.Errorf("block %d: hash or size wrong", i)
		}
	}
	if blockMap.File[1].Block[0].Size != 0 {
		t.Error("stored file has a compressed block size")
	}

	if _, err := Write([]byte(testManifest), []File{{"appxblockmap.xml", nil}}, nil); err == nil {
		t.Error("reserved name accepted")
	}
}

func TestSign(t *testing.T) {
	signer := testSigner(t, pkix.Name{Country: []string{"RS"}, Organization: []string{"Demo, Inc."}, CommonName: "Demo"})
	unsigned, err := Write([]byte(testManifest), testFiles(), nil)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := Write([]byte(testManifest), testFiles(), signer)
	if err != nil {
		t.Fatal(err)
	}
	resigned, err := Sign(signed, signer)
	if err != nil {
		t.Fatal(err)
	}
	later, err := Sign(unsigned, signer)
	if err != nil {
		t.Fatal(err)
	}
	for _, data := range [][]byte{signed, resigned, later} {
		a, err := parse(data)
		if err != nil {
			t.Fatal(err)
		}
		names := make([]string, len(a.entries))
		for i, e := range a.entries {
			names[i] = e.name
		}
		if n := len(names); names[n-1] != signatureName || names[n-2] != contentTypesName || slices.Index(names, signatureName) != n-1 {
			t.Fatalf("entries: %v", names)
		}
		files := readAll(t, data)
		p7x := files[signatureName]
		if !bytes.HasPrefix(p7x, []byte("PKCX")) {
			t.Fatal("signature has no PKCX header")
		}
		sig := a.entries[len(a.entries)-1]
		for _, part := range []struct {
			tag  string
			data []byte
		}{
			{"AXPC", data[:sig.offset]},
			{"AXCD", directory(a.entries[:len(a.entries)-1], sig.offset)},
			{"AXCT", files[contentTypesName]},
			{"AXBM", files[blockMapName]},
		} {
			sum := sha256.Sum256(part.data)
			if !bytes.Contains(p7x, append([]byte(part.tag), sum[:]...)) {
				t.Errorf("%s digest not in the signature", part.tag)
			}
		}
	}

	other := testSigner(t, pkix.Name{CommonName: "Someone Else"})
	if _, err := Sign(unsigned, other); err == nil {
		t.Error("publisher mismatch accepted")
	}
	if _, err := Sign([]byte("not a package"), signer); err == nil {
		t.Error("garbage accepted")
	}
}

func TestCanonical(t *testing.T) {
	for _, tc := range []struct{ a, b string }{
		{`CN=Demo,O="Demo, Inc.",C=RS`, `CN=Demo, O="Demo, Inc.", C=RS`},
		{`CN=A +  OU=B`, `CN=A +OU=B`},
	} {
		if canonical(tc.a) != canonical(tc.b) {
			t.Errorf("%s and %s differ", tc.a, tc.b)
		}
	}
	if canonical(`O="A,  B"`) == canonical(`O="A, B"`) {
		t.Error("spaces inside quotes were dropped")
	}
	if got := CommonName("Jane, Inc."); got != `CN="Jane, Inc."` {
		t.Errorf("CommonName: %s", got)
	}
	if ValidName("ab") || ValidName("com.example_app") || !ValidName("com.example.my-app") {
		t.Error("ValidName")
	}
}

func TestPublisher(t *testing.T) {
	signer := testSigner(t, pkix.Name{Country: []string{"RS"}, Province: []string{"Belgrade"}, Organization: []string{`Demo "Apps", Inc.`}, CommonName: "Demo"})
	got, err := Publisher(signer.Chain[0])
	if err != nil {
		t.Fatal(err)
	}
	if want := `CN=Demo, O="Demo ""Apps"", Inc.", S=Belgrade, C=RS`; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
