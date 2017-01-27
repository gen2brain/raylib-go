// +build !android,arm

package raylib

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Initialize Window and OpenGL Graphics
func InitWindow(width int32, height int32, title string) {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.InitWindow(cwidth, cheight, ctitle)
}

// Shows cursor
func ShowCursor() {
	C.ShowCursor()
}

// Hides cursor
func HideCursor() {
	C.HideCursor()
}

// Returns true if cursor is not visible
func IsCursorHidden() bool {
	ret := C.IsCursorHidden()
	v := bool(int(ret) == 1)
	return v
}

// Enables cursor
func EnableCursor() {
	C.EnableCursor()
}

// Disables cursor
func DisableCursor() {
	C.DisableCursor()
}
