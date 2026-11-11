//go:build (!cgo || nocgo) && !extlib
// +build !cgo nocgo
// +build !extlib

package rl

import _ "embed"

//go:embed libs/windows_arm64/raylib.dll.gz
var embeddedLib []byte

const embeddedLibName = "raylib.dll"
