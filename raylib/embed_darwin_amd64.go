//go:build !cgo && !extlib
// +build !cgo,!extlib

package rl

import _ "embed"

//go:embed libs/darwin_amd64/libraylib.dylib.gz
var embeddedLib []byte

const embeddedLibName = "libraylib.dylib"
