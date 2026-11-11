//go:build (!cgo || nocgo) && !extlib && (windows || linux || darwin) && (amd64 || arm64)
// +build !cgo nocgo
// +build !extlib
// +build windows linux darwin
// +build amd64 arm64

package rl

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

var extractedDir string

// extractLib unpacks the embedded shared library into a new temporary directory and returns its path
func extractLib() string {
	dir, err := os.MkdirTemp("", "raylib-go-")
	if err != nil {
		panic(err)
	}
	extractedDir = dir

	zr, err := gzip.NewReader(bytes.NewReader(embeddedLib))
	if err != nil {
		removeExtracted()
		panic(err)
	}
	defer zr.Close()

	path := filepath.Join(dir, embeddedLibName)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o700)
	if err != nil {
		removeExtracted()
		panic(err)
	}

	_, err = io.Copy(f, zr)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		removeExtracted()
		panic(err)
	}

	return path
}

// removeExtracted deletes the temporary directory of the unpacked library
func removeExtracted() {
	if extractedDir != "" {
		os.RemoveAll(extractedDir)
		extractedDir = ""
	}
}
