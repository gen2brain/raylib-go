//go:build !cgo && (freebsd || linux)

package rl

import (
	"image/color"
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

func (b *bundle) call(rValue unsafe.Pointer, aValues ...unsafe.Pointer) {
	ffi.Call(&b.cif, b.sym, rValue, aValues...)
}

var (
	// raylibDll is the pointer to the shared library
	raylibDll uintptr
)

var (
	initWindow          *bundle
	closeWindow         *bundle
	windowShouldClose   *bundle
	clearBackground     *bundle
	beginDrawing        *bundle
	endDrawing          *bundle
	setTraceLogCallback *bundle
)

func init() {
	raylibDll = loadLibrary()

	initWindow = newBundle("InitWindow", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer)
	closeWindow = newBundle("CloseWindow", &ffi.TypeVoid)
	windowShouldClose = newBundle("WindowShouldClose", &ffi.TypeUint32)
	clearBackground = newBundle("ClearBackground", &ffi.TypeVoid, &typeColor)
	beginDrawing = newBundle("BeginDrawing", &ffi.TypeVoid)
	endDrawing = newBundle("EndDrawing", &ffi.TypeVoid)
	setTraceLogCallback = newBundle("SetTraceLogCallback", &ffi.TypeVoid, &ffi.TypePointer)
}

// InitWindow - Initialize window and OpenGL context
func InitWindow(width int32, height int32, title string) {
	ctitle, err := unix.BytePtrFromString(title)
	if err != nil {
		panic(err)
	}
	initWindow.call(nil, unsafe.Pointer(&width), unsafe.Pointer(&height), unsafe.Pointer(&ctitle))
}

// CloseWindow - Close window and unload OpenGL context
func CloseWindow() {
	closeWindow.call(nil)
}

// WindowShouldClose - Check if application should close (KEY_ESCAPE pressed or windows close icon clicked)
func WindowShouldClose() bool {
	var close uint32
	windowShouldClose.call(unsafe.Pointer(&close))
	return close != 0
}

// ClearBackground - Set background color (framebuffer clear color)
func ClearBackground(col color.RGBA) {
	clearBackground.call(nil, unsafe.Pointer(&col))
}

// BeginDrawing - Setup canvas (framebuffer) to start drawing
func BeginDrawing() {
	beginDrawing.call(nil)
}

// EndDrawing - End canvas drawing and swap buffers (double buffering)
func EndDrawing() {
	endDrawing.call(nil)
}

// SetTraceLogCallback - Set custom trace log
func SetTraceLogCallback(fn TraceLogCallbackFun) {
	callback := traceLogCallbackWrapper(fn)
	setTraceLogCallback.call(nil, unsafe.Pointer(&callback))
}
