//go:build !cgo && windows
// +build !cgo,windows

package rl

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/windows"
)

const (
	libname         = "raylib.dll"
	requiredVersion = "5.0"
)

// loadLibrary loads the raylib dll and panics on error
func loadLibrary() uintptr {
	handle, err := windows.LoadLibrary(libname)
	if err != nil {
		panic(fmt.Errorf("cannot load library %s: %w", libname, err))
	}

	proc, err := windows.GetProcAddress(handle, "raylib_version")
	if err != nil {
		panic(err)
	}

	version := windows.BytePtrToString(**(***byte)(unsafe.Pointer(&proc)))
	if version != requiredVersion {
		panic(fmt.Errorf("version %s of %s doesn't match the required version %s", version, libname, requiredVersion))
	}

	return uintptr(handle)
}

func traceLogCallbackWrapper(fn TraceLogCallbackFun) uintptr {
	return purego.NewCallback(func(logLevel int32, text *byte) uintptr {
		fn(int(logLevel), windows.BytePtrToString(text))
		return 0
	})
}
