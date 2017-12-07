// +build android,!js

package raylib

/*
#include "raylib.h"
#include <android/asset_manager.h>
#include <android_native_app_glue.h>

extern void android_main(struct android_app *app);
extern void android_init(void *state);

AAssetManager* asset_manager;
const char* internal_storage_path;
*/
import "C"

import (
	"errors"
	"fmt"
	"io"
	"unsafe"
)

var callbackHolder func(unsafe.Pointer)

// InitWindow - Initialize Window and OpenGL Graphics
func InitWindow(width int32, height int32, t interface{}) {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)

	app, ok := t.(unsafe.Pointer)
	if ok {
		C.InitWindow(cwidth, cheight, app)
		C.android_init(app)
	}
}

// SetCallbackFunc - Sets callback function
func SetCallbackFunc(callback func(unsafe.Pointer)) {
	callbackHolder = callback
}

//export androidMain
func androidMain(app *C.struct_android_app) {
	if callbackHolder != nil {
		callbackHolder(unsafe.Pointer(app))
	}
}

// ShowCursor - Shows cursor
func ShowCursor() {
	return
}

// HideCursor - Hides cursor
func HideCursor() {
	return
}

// IsCursorHidden - Returns true if cursor is not visible
func IsCursorHidden() bool {
	return false
}

// EnableCursor - Enables cursor
func EnableCursor() {
	return
}

// DisableCursor - Disables cursor
func DisableCursor() {
	return
}

// IsFileDropped - Check if a file have been dropped into window
func IsFileDropped() bool {
	return false
}

// GetDroppedFiles - Retrieve dropped files into window
func GetDroppedFiles(count *int32) (files []string) {
	return
}

// ClearDroppedFiles - Clear dropped files paths buffer
func ClearDroppedFiles() {
	return
}

// OpenAsset - Open asset
func OpenAsset(name string) (Asset, error) {
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

func (a *asset) Seek(offset int64, whence int) (int64, error) {
	off := C.AAsset_seek(a.ptr, C.off_t(offset), C.int(whence))
	if off == -1 {
		return 0, errors.New(fmt.Sprintf("bad result for offset=%d, whence=%d", offset, whence))
	}
	return int64(off), nil
}

func (a *asset) Close() error {
	C.AAsset_close(a.ptr)
	return nil
}
