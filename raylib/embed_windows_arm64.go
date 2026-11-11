//go:build !cgo && !extlib
// +build !cgo,!extlib

package rl

import _ "embed"

//go:embed libs/windows_arm64/raylib.dll.gz
var embeddedLib []byte

const embeddedLibName = "raylib.dll"
