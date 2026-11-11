//go:build !cgo && !extlib && (linux || darwin) && (amd64 || arm64)
// +build !cgo
// +build !extlib
// +build linux darwin
// +build amd64 arm64

package rl

// afterLoad removes the unpacked library, the loaded mapping stays valid
func afterLoad() {
	removeExtracted()
}

func unloadLib() {}
