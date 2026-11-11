//go:build !cgo
// +build !cgo

package rl

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/gen2brain/raylib-go/raylib/internal/convert"
	"github.com/jupiterrider/ffi"
)

const (
	requiredVersion = "6.2"
)

var vsprintf ffi.Fun

func init() {
	var filename string
	funcname := "vsprintf"
	switch runtime.GOOS {
	case "linux":
		filename = "libc.so.6"
	case "freebsd":
		filename = "libc.so.7"
	case "windows":
		filename = "user32.dll"
		funcname = "wvsprintfA"
	case "darwin":
		filename = "libc.dylib"
	}

	if lib, err := ffi.Load(filename); err == nil {
		if fun, err := lib.Prep(funcname, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err == nil {
			vsprintf = fun
		}
	}
}

// loadLibrary loads the raylib shared library and panics on error
func loadLibrary() ffi.Lib {
	libname := extractLib()

	if len(libname) == 0 {
		switch runtime.GOOS {
		case "freebsd", "linux":
			libname = "libraylib.so.6.2.0"
		case "windows":
			libname = "raylib.dll"
		case "darwin":
			libname = "libraylib.6.2.0.dylib"
		}
	}

	lib, err := ffi.Load(libname)
	if err != nil {
		removeExtracted()
		panic(fmt.Errorf("cannot load library %s: %w", libname, err))
	}
	afterLoad()

	addr, err := lib.Get("raylib_version")
	if err != nil {
		panic(err)
	}

	version := convert.ToString(**(***byte)(unsafe.Pointer(&addr)))
	if version != requiredVersion {
		panic(fmt.Errorf("version %s of %s doesn't match the required version %s", version, libname, requiredVersion))
	}

	return lib
}

func traceLogCallbackWrapper(fn TraceLogCallbackFun) uintptr {
	return purego.NewCallback(func(logLevel int32, text *byte, args unsafe.Pointer) uintptr {
		if vsprintf.Addr != 0 {
			var buffer [1024]byte // Max size is 1024 (see https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-wvsprintfa)
			bufferPtr := &buffer[0]
			var ret int32
			vsprintf.Call(&ret, &bufferPtr, &text, &args)
			if ret > 0 {
				fn(int(logLevel), convert.ToString(bufferPtr))
				return 0
			}
		}
		fn(int(logLevel), convert.ToString(text))
		return 0
	})
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
