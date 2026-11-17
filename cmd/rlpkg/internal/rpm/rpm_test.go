package rpm

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/pgp"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/pkgtree"
)

func TestWrite(t *testing.T) {
	data, err := Write(Package{
		Name: "demo", Version: "1.0", Release: "1", Arch: "x86_64", Summary: "Demo app", Description: "Demo app", License: "MIT", Vendor: "Demo",
		Files: []pkgtree.File{
			{Path: "/usr/bin/demo", Mode: 0o755, Data: []byte("#!/bin/sh\necho demo\n")},
			{Path: "/usr/share/applications/demo.desktop", Mode: 0o644, Data: []byte("[Desktop Entry]\n")},
			{Path: "/usr/bin/demo-link", Mode: 0o777, Link: "demo"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte{0xed, 0xab, 0xee, 0xdb}) {
		t.Fatal("bad lead")
	}
	if _, err := exec.LookPath("rpm"); err != nil {
		t.Skip("rpm not found")
	}
	path := filepath.Join(t.TempDir(), "demo.rpm")
	os.WriteFile(path, data, 0o644)
	out, err := exec.Command("rpm", "-qpi", path).CombinedOutput()
	if err != nil {
		t.Fatalf("rpm -qpi: %v\n%s", err, out)
	}
	for _, want := range []string{"Name        : demo", "Version     : 1.0", "Architecture: x86_64", "License     : MIT"} {
		if !bytes.Contains(out, []byte(want)) {
			t.Errorf("missing %q in\n%s", want, out)
		}
	}
	out, _ = exec.Command("rpm", "-K", "-v", path).CombinedOutput()
	if !bytes.Contains(out, []byte("Header SHA256 digest: OK")) || !bytes.Contains(out, []byte("Payload SHA256 digest: OK")) {
		t.Errorf("rpm -K:\n%s", out)
	}
	out, err = exec.Command("rpm", "-qpl", path).CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("/usr/bin/demo-link")) {
		t.Errorf("rpm -qpl: %v\n%s", err, out)
	}
}

func testSigner(t *testing.T, algo packet.PublicKeyAlgorithm) (pgp.Signer, *openpgp.Entity) {
	t.Helper()
	entity, err := openpgp.NewEntity("Demo", "", "demo@example.com", &packet.Config{Algorithm: algo, RSABits: 2048})
	if err != nil {
		t.Fatal(err)
	}
	return func(data []byte) ([]byte, error) {
		var sig bytes.Buffer
		err := openpgp.DetachSign(&sig, entity, bytes.NewReader(data), nil)
		return sig.Bytes(), err
	}, entity
}

func checkSignatures(t *testing.T, data []byte, entity *openpgp.Entity, hdrTag, allTag int32) {
	t.Helper()
	entries, sigLen, err := parseHeader(data[96:])
	if err != nil {
		t.Fatal(err)
	}
	rest := 96 + (sigLen+7)/8*8
	_, hdrLen, err := parseHeader(data[rest:])
	if err != nil {
		t.Fatal(err)
	}
	signed := map[int32][]byte{hdrTag: data[rest : rest+hdrLen], allTag: data[rest:]}
	for _, e := range entries {
		if !slices.Contains([]int32{tagSigPGP, tagSigGPG, tagSigDSA, tagSigRSA}, e.tag) {
			continue
		}
		body, ok := signed[e.tag]
		if !ok {
			t.Fatalf("unexpected signature tag %d", e.tag)
		}
		if _, err := openpgp.CheckDetachedSignature(openpgp.EntityList{entity}, bytes.NewReader(body), bytes.NewReader(e.data), nil); err != nil {
			t.Fatalf("tag %d: %v", e.tag, err)
		}
		delete(signed, e.tag)
	}
	if len(signed) != 0 {
		t.Fatalf("missing signature tags, left %d", len(signed))
	}
}

func TestSign(t *testing.T) {
	sign, entity := testSigner(t, packet.PubKeyAlgoRSA)
	data, err := Write(Package{
		Name: "demo", Version: "1.0", Release: "1", Arch: "x86_64", Summary: "Demo app", Description: "Demo app", License: "MIT", Vendor: "Demo",
		Files: []pkgtree.File{{Path: "/usr/bin/demo", Mode: 0o755, Data: []byte("#!/bin/sh\n")}},
		Sign:  sign,
	})
	if err != nil {
		t.Fatal(err)
	}
	checkSignatures(t, data, entity, tagSigRSA, tagSigPGP)

	unsigned, err := Sign(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	sign, entity = testSigner(t, packet.PubKeyAlgoEdDSA)
	resigned, err := Sign(unsigned, sign)
	if err != nil {
		t.Fatal(err)
	}
	checkSignatures(t, resigned, entity, tagSigDSA, tagSigGPG)

	if _, err := Sign([]byte("not a package"), sign); err == nil {
		t.Error("garbage accepted")
	}
	lead := append([]byte{0xed, 0xab, 0xee, 0xdb}, make([]byte, 92)...)
	v6 := assemble(lead, header(tagHeaderSignatures, nil), header(tagHeaderImmutable, []entry{int32s(tagRPMFormat, []int32{6})}), nil)
	if _, err := Sign(v6, sign); err == nil {
		t.Error("v6 package accepted")
	}
	v4 := assemble(lead, header(tagHeaderSignatures, nil), header(tagHeaderImmutable, []entry{int32s(tagRPMFormat, []int32{4})}), nil)
	if _, err := Sign(v4, sign); err != nil {
		t.Errorf("v4 package: %v", err)
	}
}
