//go:build !android

package rl

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestAsset(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "img"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "img", "a.png"), []byte("os"), 0o644); err != nil {
		t.Fatal(err)
	}

	mapfs := fstest.MapFS{"assets/img/a.png": {Data: []byte("fs")}}

	for _, tc := range []struct {
		name  string
		asset *Asset
		want  string
	}{
		{"os", NewAsset(dir), "os"},
		{"fs", NewAssetFromFS(mapfs, "assets"), "fs"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.asset.ReadFile("img/a.png")
			if err != nil || string(data) != tc.want {
				t.Fatalf("ReadFile = %q, %v; want %q", data, err, tc.want)
			}

			f, err := tc.asset.Open("img/a.png")
			if err != nil {
				t.Fatal(err)
			}
			data, err = io.ReadAll(f)
			f.Close()
			if err != nil || string(data) != tc.want {
				t.Fatalf("Open = %q, %v; want %q", data, err, tc.want)
			}

			entries, err := tc.asset.ReadDir("img")
			if err != nil || len(entries) != 1 || entries[0].Name() != "a.png" {
				t.Fatalf("ReadDir = %v, %v; want [a.png]", entries, err)
			}
		})
	}
}
