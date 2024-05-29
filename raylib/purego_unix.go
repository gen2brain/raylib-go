//go:build !cgo && (freebsd || linux)

package rl

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/unix"
)

const (
	libname         = "libraylib.so"
	requiredVersion = "5.0"
)

var vsprintf uintptr

func init() {
	var filename string
	switch runtime.GOOS {
	case "linux":
		filename = "libc.so.6"
	case "freebsd":
		filename = "libc.so.7"
	}

	handle, err := purego.Dlopen(filename, purego.RTLD_LAZY)
	if err == nil {
		vsprintf, _ = purego.Dlsym(handle, "vsprintf")
	}
}

// loadLibrary loads the raylib dll and panics on error
func loadLibrary() uintptr {
	handle, err := purego.Dlopen(libname, purego.RTLD_LAZY)
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

	//DEBUG
	fmt.Println(version)

	return uintptr(handle)
}

func traceLogCallbackWrapper(fn TraceLogCallbackFun) uintptr {
	return purego.NewCallback(func(logLevel int32, text *byte, args unsafe.Pointer) uintptr {
		if vsprintf != 0 {
			var buffer [1024]byte // Max size is 1024 (see https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-wvsprintfa)
			_, _, errno := purego.SyscallN(vsprintf, uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(text)), uintptr(args))
			if errno == 0 {
				text = &buffer[0]
			}
		}
		fn(int(logLevel), unix.BytePtrToString(text))
		return 0
	})
}
