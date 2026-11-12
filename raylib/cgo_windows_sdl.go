//go:build windows && (sdl || sdl3) && !rgfw && !win32 && !memory && !nocgo
// +build windows
// +build sdl sdl3
// +build !rgfw
// +build !win32
// +build !memory
// +build !nocgo

package rl

/*
#cgo windows LDFLAGS: -lgdi32 -lwinmm -lole32
#cgo windows,sdl LDFLAGS: -lSDL2
#cgo windows,sdl3 LDFLAGS: -lSDL3
#cgo windows CFLAGS: -Iexternal -Wno-stringop-overflow
#cgo windows,sdl CFLAGS: -DPLATFORM_DESKTOP_SDL
#cgo windows,sdl3 CFLAGS: -DPLATFORM_DESKTOP_SDL -DPLATFORM_DESKTOP_SDL3 -DUSING_SDL3_PROJECT

#cgo windows,!es2,!es3,!software LDFLAGS: -lopengl32

#cgo windows,opengl11,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo windows,opengl21,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo windows,opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo windows,!opengl11,!opengl21,!opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo windows,es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo windows,es3,!es2,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES3
#cgo windows,software CFLAGS: -DGRAPHICS_API_OPENGL_SOFTWARE
*/
import "C"
