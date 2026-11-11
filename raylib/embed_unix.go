//go:build !cgo && !raylib_no_embed && (linux || darwin) && (amd64 || arm64)
// +build !cgo
// +build !raylib_no_embed
// +build linux darwin
// +build amd64 arm64

package rl

// afterLoad removes the unpacked library, the loaded mapping stays valid
func afterLoad() {
	removeExtracted()
}

func unloadLib() {}
