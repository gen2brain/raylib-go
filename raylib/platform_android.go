// +build android

package raylib

/*
#include "raylib.h"
#include "android_native_app_glue.h"

extern void android_main(struct android_app *app);
*/
import "C"
import "unsafe"

var callbackHolder func(unsafe.Pointer)

// Initialize Window and OpenGL Graphics
func InitWindow(width int32, height int32, app unsafe.Pointer) {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	C.InitWindow(cwidth, cheight, app)
}

// Sets callback function
func SetCallbackFunc(callback func(unsafe.Pointer)) {
	callbackHolder = callback
}

//export androidMain
func androidMain(app *C.struct_android_app) {
	if callbackHolder != nil {
		callbackHolder(unsafe.Pointer(app))
	}
}
