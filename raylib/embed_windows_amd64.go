//go:build !cgo && !raylib_no_embed
// +build !cgo,!raylib_no_embed

package rl

import _ "embed"

//go:embed libs/windows_amd64/raylib.dll.gz
var embeddedLib []byte

const embeddedLibName = "raylib.dll"
