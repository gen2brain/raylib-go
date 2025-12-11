//go:build windows && sdl3 && !rgfw
// +build windows,sdl3,!rgfw

package rl

/*
#cgo windows LDFLAGS: -lgdi32 -lwinmm -lole32
#cgo windows,sdl3 LDFLAGS: -lSDL3
#cgo windows CFLAGS: -Iexternal -Wno-stringop-overflow
#cgo windows,sdl3 CFLAGS: -DPLATFORM_DESKTOP_SDL -DPLATFORM_DESKTOP_SDL3 -DUSING_SDL3_PACKAGE

#cgo windows,!es2,!es3 LDFLAGS: -lopengl32

#cgo windows,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo windows,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo windows,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo windows,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo windows,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo windows,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3
*/
import "C"
