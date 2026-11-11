//go:build !cgo && !raylib_no_embed
// +build !cgo,!raylib_no_embed

package rl

import _ "embed"

//go:embed libs/darwin_arm64/libraylib.dylib.gz
var embeddedLib []byte

const embeddedLibName = "libraylib.dylib"
