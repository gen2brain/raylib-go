package raygui3_5

/*
#cgo CFLAGS: -DRAYGUI_IMPLEMENTATION
#include "../raylib/raygui.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
	rl  "github.com/Konstantin8105/raylib-go/raylib"
)

// int GuiToggleGroup(Rectangle bounds, const char *text, int active)
func ToggleGroup(bounds rl.Rectangle, text string, active int) int {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cactive := C.int(active)
	res := C.GuiToggleGroup(cbounds, ctext, cactive)
	return int(res)
}

// bool GuiButton(Rectangle bounds, const char *text)
func Button(bounds rl.Rectangle, text string) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	res := C.GuiButton(cbounds, ctext)
	return bool(res)
}
