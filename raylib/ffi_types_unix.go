//go:build !cgo && (freebsd || linux)

package rl

import "github.com/jupiterrider/ffi"

var (
	typeColor = ffi.Type{Type: ffi.Struct, Elements: &[]*ffi.Type{&ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, nil}[0]}
)
