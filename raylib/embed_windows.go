//go:build !cgo && !raylib_no_embed && windows && (amd64 || arm64)
// +build !cgo
// +build !raylib_no_embed
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
