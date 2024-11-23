//go:build !cgo && windows
// +build !cgo,windows

package rl

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/windows"
)

const (
	libname         = "raylib.dll"
	requiredVersion = "5.5"
)

var wvsprintfA uintptr

func init() {
	handle, err := windows.LoadLibrary("user32.dll")
	if err == nil {
		wvsprintfA, _ = windows.GetProcAddress(handle, "wvsprintfA")
	}
}

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
	return purego.NewCallback(func(logLevel int32, text *byte, args unsafe.Pointer) uintptr {
		if wvsprintfA != 0 {
			var buffer [1024]byte // Max size is 1024 (see https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-wvsprintfa)
			_, _, errno := syscall.SyscallN(wvsprintfA, uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(text)), uintptr(args))
			if errno == 0 {
				text = &buffer[0]
			}
		}
		fn(int(logLevel), windows.BytePtrToString(text))
		return 0
	})
}
