package rl

/*
#include "raylib.h"
#include "raygui.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func GuiLoadStyle(fileName string) {
	cfilename := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfilename))
	C.GuiLoadStyle(cfilename)
}

func GuiGetStyle(control, property int) int {
	ccontrol := (C.int)(control)
	cproperty := (C.int)(property)
	return int(C.GuiGetStyle(ccontrol, cproperty))
}

func GuiSetStyle(control, property, value int) {
	ccontrol := (C.int)(control)
	cproperty := (C.int)(property)
	cvalue := (C.int)(value)
	C.GuiSetStyle(ccontrol, cproperty, cvalue)
}

func GuiLabel(bounds Rectangle, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.GuiLabel(*bounds.cptr(), ctext)
}
