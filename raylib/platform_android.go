// +build android

package raylib

/*
#include "raylib.h"
#include <android/asset_manager.h>
#include <android_native_app_glue.h>

extern void android_main(struct android_app *app);

AAssetManager* asset_manager;
extern void init_asset_manager(void *state);
*/
import "C"

import (
	"errors"
	"io"
	"unsafe"
)

var callbackHolder func(unsafe.Pointer)

// Initialize Window and OpenGL Graphics
func InitWindow(width int32, height int32, t interface{}) {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)

	app, ok := t.(unsafe.Pointer)
	if ok {
		C.InitWindow(cwidth, cheight, app)
		C.init_asset_manager(app)
	}
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

// Open asset
func OpenAsset(name string) (io.ReadCloser, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	a := &asset{C.AAssetManager_open(C.asset_manager, cname, C.AASSET_MODE_UNKNOWN)}

	if a.ptr == nil {
		return nil, errors.New("asset file could not be opened")
	}

	return a, nil
}

type asset struct {
	ptr *C.AAsset
}

func (a *asset) Read(p []byte) (n int, err error) {
	n = int(C.AAsset_read(a.ptr, unsafe.Pointer(&p[0]), C.size_t(len(p))))
	if n == 0 && len(p) > 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (a *asset) Close() error {
	C.AAsset_close(a.ptr)
	return nil
}
