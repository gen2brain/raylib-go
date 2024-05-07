//go:build !cgo && (darwin || openbsd || freebsd || linux)

package rl

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/unix"
)

// loadLibrary loads the raylib dll and panics on error
func loadLibrary() uintptr {
	libname := "raylib.so"
	switch runtime.GOOS {
	case "darwin":
		libname = "raylib.dylib"
	}

	handle, err := purego.Dlopen(libname, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(fmt.Errorf("cannot load library %s: %w", libname, err))
	}

	proc, err := purego.Dlsym(handle, "raylib_version")
	if err != nil {
		panic(err)
	}

	version := unix.BytePtrToString(**(***byte)(unsafe.Pointer(&proc)))
	if version != requiredVersion {
		panic(fmt.Errorf("version %s of %s doesn't match the required version %s", version, libname, requiredVersion))
	}

	return uintptr(handle)
}

func traceLogCallbackWrapper(fn TraceLogCallbackFun) uintptr {
	return purego.NewCallback(func(logLevel int32, text *byte) uintptr {
		fn(int(logLevel), unix.BytePtrToString(text))
		return 0
	})
}
