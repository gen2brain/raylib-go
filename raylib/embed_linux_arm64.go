//go:build (!cgo || nocgo) && !extlib
// +build !cgo nocgo
// +build !extlib

package rl

import _ "embed"

//go:embed libs/linux_arm64/libraylib.so.gz
var embeddedLib []byte

const embeddedLibName = "libraylib.so"
