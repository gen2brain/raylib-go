package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"os"
	"unsafe"
)

func InitRaylib(width int32, height int32) {
	initWindow(width, height, "raylib-go-headless")
}

// CloseRaylib - Close headless window and terminate graphics context
func CloseRaylib() {
	C.CloseWindow()
}

// RaylibShouldClose - Check if KEY_ESCAPE pressed ... should check for interrupts?
func RaylibShouldClose() bool {
	ret := C.WindowShouldClose()
	v := bool(ret)
	return v
}

// initWindow - Initialize Raylib and OpenGL Graphics
func initWindow(width int32, height int32, title string) {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)

	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	C.InitWindow(cwidth, cheight, ctitle)
	C.HideCursor()
}

// SetCallbackFunc - Sets callback function
func SetCallbackFunc(func()) {
	return
}

// OpenAsset - Open asset
func OpenAsset(name string) (Asset, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}
