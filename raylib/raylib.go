// Package raylib - Go bindings for raylib, a simple and easy-to-use library to learn videogames programming
package raylib

import (
	"runtime"
)

func init() {
	runtime.LockOSThread()
}
