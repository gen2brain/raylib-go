// +build js

package raylib

import (
	"os"
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
	js.Global.Set("_go_update_function", fn)
	js.Global.Get("Module").Call("_emscripten_set_main_loop_go", fps, infinite)
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
func OpenAsset(name string) (Asset, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}
