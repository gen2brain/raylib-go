package macho

import (
	"bytes"
	"compress/gzip"
	"debug/macho"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func darwinDylib(t *testing.T, arch string) []byte {
	_, file, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "raylib", "libs", "darwin_"+arch, "libraylib.dylib.gz")
	f, err := os.Open(path)
	if err != nil {
		t.Skip(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestUniversal(t *testing.T) {
	amd64 := darwinDylib(t, "amd64")
	arm64 := darwinDylib(t, "arm64")

	fat, err := Universal([][]byte{amd64, arm64})
	if err != nil {
		t.Fatal(err)
	}
	ff, err := macho.NewFatFile(bytes.NewReader(fat))
	if err != nil {
		t.Fatal(err)
	}
	if len(ff.Arches) != 2 {
		t.Fatalf("arches = %d, want 2", len(ff.Arches))
	}
	for i, want := range [][]byte{amd64, arm64} {
		a := ff.Arches[i]
		if a.Offset%(1<<a.Align) != 0 {
			t.Errorf("arch %d offset %d not aligned to 2^%d", i, a.Offset, a.Align)
		}
		if !bytes.Equal(fat[a.Offset:a.Offset+a.Size], want) {
			t.Errorf("arch %d contents differ", i)
		}
	}

	single, err := Universal([][]byte{arm64})
	if err != nil || !bytes.Equal(single, arm64) {
		t.Errorf("single slice was modified")
	}
}

func TestMacOSMinVersion(t *testing.T) {
	v, err := MinVersion(darwinDylib(t, "arm64"))
	if err != nil {
		t.Fatal(err)
	}
	if v != "11.0" {
		t.Errorf("minimum version = %q, want 11.0", v)
	}
}
