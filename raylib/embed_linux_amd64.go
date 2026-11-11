//go:build !cgo && !raylib_no_embed
// +build !cgo,!raylib_no_embed

package rl

import _ "embed"

//go:embed libs/linux_amd64/libraylib.so.gz
var embeddedLib []byte

const embeddedLibName = "libraylib.so"
