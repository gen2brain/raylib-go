//go:build !cgo && (freebsd || linux)

package rl

import (
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
	"golang.org/x/sys/unix"
)

// bundle bundles the function pointer "sym" and the ffi call interface "cif"
type bundle struct {
	sym uintptr
	cif ffi.Cif
}

func newBundle(name string, rType *ffi.Type, aTypes ...*ffi.Type) *bundle {
	b := new(bundle)
	var err error

	if b.sym, err = purego.Dlsym(raylibDll, name); err != nil {
		panic(err)
	}

	nArgs := uint32(len(aTypes))

	if status := ffi.PrepCif(&b.cif, ffi.DefaultAbi, nArgs, rType, aTypes...); status != ffi.OK {
		panic(status)
	}

	return b
}

var (
	// raylibDll is the pointer to the shared library
	raylibDll uintptr
)

var (
	initWindow          *bundle
	closeWindow         *bundle
	setTraceLogCallback *bundle
)

func init() {
	raylibDll = loadLibrary()

	initWindow = newBundle("InitWindow", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer)
	closeWindow = newBundle("CloseWindow", &ffi.TypeVoid)
	setTraceLogCallback = newBundle("SetTraceLogCallback", &ffi.TypeVoid, &ffi.TypePointer)
}

// InitWindow - Initialize window and OpenGL context
func InitWindow(width int32, height int32, title string) {
	ctitle, err := unix.BytePtrFromString(title)
	if err != nil {
		panic(err)
	}
	ffi.Call(&initWindow.cif, initWindow.sym, nil, unsafe.Pointer(&width), unsafe.Pointer(&height), unsafe.Pointer(&ctitle))
}

// CloseWindow - Close window and unload OpenGL context
func CloseWindow() {
	ffi.Call(&closeWindow.cif, closeWindow.sym, nil)
}

// SetTraceLogCallback - Set custom trace log
func SetTraceLogCallback(fn TraceLogCallbackFun) {
	callback := traceLogCallbackWrapper(fn)
	ffi.Call(&setTraceLogCallback.cif, setTraceLogCallback.sym, nil, unsafe.Pointer(&callback))
}
