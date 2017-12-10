// +build js

package raylib

import (
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

// InitWindow - Initialize Window and OpenGL Graphics
func InitWindow(width int32, height int32, t interface{}) {
	js.Global.Get("Module").Call("_InitWindow", width, height, t.(string))
}

// SetCallbackFunc - Sets callback function
func SetCallbackFunc(func(unsafe.Pointer)) {
}

// SetMainLoop - Sets main loop function
func SetMainLoop(fn func(), fps int, infinite bool) {
	js.Global.Get("Module").Call("_emscripten_set_main_loop", fn, fps, infinite)
}

// ShowCursor - Shows cursor
func ShowCursor() {
}

// HideCursor - Hides cursor
func HideCursor() {
}

// IsCursorHidden - Returns true if cursor is not visible
func IsCursorHidden() bool {
	return false
}

// EnableCursor - Enables cursor
func EnableCursor() {
}

// DisableCursor - Disables cursor
func DisableCursor() {
}

// IsFileDropped - Check if a file have been dropped into window
func IsFileDropped() bool {
	return false
}

// GetDroppedFiles - Retrieve dropped files into window
func GetDroppedFiles(count *int32) (f []string) {
	return
}

// ClearDroppedFiles - Clear dropped files paths buffer
func ClearDroppedFiles() {
}

// OpenAsset - Open asset
func OpenAsset(name string) (a Asset, err error) {
	defer func() {
		e := recover()

		if e == nil {
			return
		}

		if e, ok := e.(*js.Error); ok {
			err = e
		} else {
			panic(e)
		}
	}()

	ptr := js.Global.Get("FS").Call("open", name, "r")
	a = &asset{ptr, 0}

	return
}

type asset struct {
	ptr    *js.Object
	offset int64
}

func (a *asset) Read(p []byte) (n int, err error) {
	defer func() {
		e := recover()

		if e == nil {
			return
		}

		if e, ok := e.(*js.Error); ok {
			err = e
		} else {
			panic(e)
		}
	}()

	js.Global.Get("FS").Call("read", a.ptr, p, 0, cap(p), a.offset)
	n = len(p)
	return
}

func (a *asset) Seek(offset int64, whence int) (off int64, err error) {
	defer func() {
		e := recover()

		if e == nil {
			return
		}

		if e, ok := e.(*js.Error); ok {
			err = e
		} else {
			panic(e)
		}
	}()

	off = js.Global.Get("FS").Call("llseek", a.ptr, int(offset), int(whence)).Int64()
	a.offset = off
	return
}

func (a *asset) Close() (err error) {
	defer func() {
		e := recover()

		if e == nil {
			return
		}

		if e, ok := e.(*js.Error); ok {
			err = e
		} else {
			panic(e)
		}
	}()

	js.Global.Get("FS").Call("close", a.ptr)
	return nil
}
