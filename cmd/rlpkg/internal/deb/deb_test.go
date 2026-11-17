package deb

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/pgp"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/pkgtree"
)

func TestWrite(t *testing.T) {
	data, err := Write(Package{
		Name: "demo", Version: "1.0-1", Arch: "amd64", Maintainer: "Demo <demo@example.com>", Description: "Demo app", Section: "misc",
		Files: []pkgtree.File{
			{Path: "/usr/bin/demo", Mode: 0o755, Data: []byte("#!/bin/sh\necho demo\n")},
			{Path: "/usr/share/applications/demo.desktop", Mode: 0o644, Data: []byte("[Desktop Entry]\n")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte("!<arch>\ndebian-binary")) {
		t.Fatal("not an ar archive")
	}
	if _, err := exec.LookPath("dpkg-deb"); err != nil {
		t.Skip("dpkg-deb not found")
	}
	path := filepath.Join(t.TempDir(), "demo.deb")
	os.WriteFile(path, data, 0o644)
	out, err := exec.Command("dpkg-deb", "-I", path).CombinedOutput()
	if err != nil {
		t.Fatalf("dpkg-deb -I: %v\n%s", err, out)
	}
	for _, want := range []string{"Package: demo", "Version: 1.0-1", "Architecture: amd64"} {
		if !bytes.Contains(out, []byte(want)) {
			t.Errorf("missing %q in\n%s", want, out)
		}
	}
	out, err = exec.Command("dpkg-deb", "-c", path).CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("./usr/bin/demo")) || !bytes.Contains(out, []byte("root/root")) {
		t.Errorf("dpkg-deb -c: %v\n%s", err, out)
	}
}

func testSigner(t *testing.T) (pgp.Signer, *openpgp.Entity) {
	t.Helper()
	entity, err := openpgp.NewEntity("Demo", "", "demo@example.com", &packet.Config{Algorithm: packet.PubKeyAlgoEdDSA})
	if err != nil {
		t.Fatal(err)
	}
	return func(data []byte) ([]byte, error) {
		var sig bytes.Buffer
		err := openpgp.DetachSign(&sig, entity, bytes.NewReader(data), nil)
		return sig.Bytes(), err
	}, entity
}

func checkOrigin(t *testing.T, data []byte, entity *openpgp.Entity) {
	t.Helper()
	var names []string
	var signed, sig []byte
	for rest := data[len(arMagic):]; len(rest) > 0; {
		name := strings.TrimSpace(string(rest[:16]))
		size, err := strconv.Atoi(strings.TrimSpace(string(rest[48:58])))
		if err != nil {
			t.Fatal(err)
		}
		names = append(names, name)
		if name == originName {
			sig = rest[60 : 60+size]
		} else {
			signed = append(signed, rest[60:60+size]...)
		}
		rest = rest[min(60+size+size%2, len(rest)):]
	}
	if got := strings.Join(names, " "); got != "debian-binary control.tar.gz data.tar.gz _gpgorigin" {
		t.Fatalf("members: %s", got)
	}
	if _, err := openpgp.CheckDetachedSignature(openpgp.EntityList{entity}, bytes.NewReader(signed), bytes.NewReader(sig), nil); err != nil {
		t.Fatal(err)
	}
}

func TestSign(t *testing.T) {
	sign, entity := testSigner(t)
	data, err := Write(Package{
		Name: "demo", Version: "1.0-1", Arch: "amd64", Maintainer: "Demo <demo@example.com>", Description: "Demo app", Section: "misc",
		Files: []pkgtree.File{{Path: "/usr/bin/demo", Mode: 0o755, Data: []byte("#!/bin/sh\n")}},
		Sign:  sign,
	})
	if err != nil {
		t.Fatal(err)
	}
	checkOrigin(t, data, entity)

	sign, entity = testSigner(t)
	if data, err = Sign(data, sign); err != nil {
		t.Fatal(err)
	}
	checkOrigin(t, data, entity)

	if _, err := Sign([]byte("not a package"), sign); err == nil {
		t.Error("garbage accepted")
	}
}
