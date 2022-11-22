package rl

/*
#cgo CFLAGS: -DRAYGUI_IMPLEMENTATION
#include "raygui.h"
#include <stdlib.h>
*/
import "C"

import (
"unsafe"
"fmt"
)

// int GuiToggleGroup(Rectangle bounds, const char *text, int active)
func ToggleGroup(bounds Rectangle, text string, active int) int {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cactive := C.int(active)
	res := C.GuiToggleGroup(cbounds, ctext, cactive)
	fmt.Printf(">%v %v<", active, res)
	return int(res)
}


// bool GuiButton(Rectangle bounds, const char *text)
func Button(bounds Rectangle, text string) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	res := C.GuiButton(cbounds, ctext)
	fmt.Printf(">%v<", res)
	return bool(res)
}
