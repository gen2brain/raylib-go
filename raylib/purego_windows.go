//go:build !cgo && windows
// +build !cgo,windows

package rl

import (
	"fmt"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/windows"
)

const (
	libname = "raylib.dll"
)

// loadLibrary loads the raylib dll and panics on error
func loadLibrary() uintptr {
	if handle, err := windows.LoadLibrary(libname); err != nil {
		panic(fmt.Errorf("cannot load library %s: %w", libname, err))
	} else {
		return uintptr(handle)
	}
}

func traceLogCallbackWrapper(fn TraceLogCallbackFun) uintptr {
	return purego.NewCallback(func(logLevel int32, text *byte) uintptr {
		fn(int(logLevel), windows.BytePtrToString(text))
		return 0
	})
}
