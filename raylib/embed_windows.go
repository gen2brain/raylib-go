//go:build (!cgo || nocgo) && !extlib && windows && (amd64 || arm64)
// +build !cgo nocgo
// +build !extlib
// +build windows
// +build amd64 arm64

package rl

func afterLoad() {}

// unloadLib frees the embedded DLL and removes its temporary directory, a loaded DLL cannot be deleted
func unloadLib() {
	if extractedDir == "" {
		return
	}

	dll.Close()
	removeExtracted()
}
