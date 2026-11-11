//go:build windows && win32 && !rgfw && !sdl && !sdl3 && !nocgo
// +build windows,win32,!rgfw,!sdl,!sdl3,!nocgo

package rl

/*
#cgo windows LDFLAGS: -lgdi32 -lwinmm -lshcore
#cgo windows CFLAGS: -Iexternal -DPLATFORM_DESKTOP_WIN32 -D_WIN32_WINNT=0x0A00 -DWINVER=0x0A00 -Wno-stringop-overflow

#cgo windows,!es2,!es3 LDFLAGS: -lopengl32

#cgo windows,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo windows,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo windows,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo windows,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo windows,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo windows,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3
*/
import "C"
