//go:build windows && sdl
// +build windows,sdl

package rl

/*
#cgo windows LDFLAGS: -lgdi32 -lwinmm -lole32 -lSDL2
#cgo windows CFLAGS: -Iexternal -DPLATFORM_DESKTOP_SDL -Wno-stringop-overflow

#cgo windows,!angle LDFLAGS: -lopengl32

#cgo windows,opengl11,!angle CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo windows,opengl21,!angle CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo windows,opengl43,!angle CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo windows,!opengl11,!opengl21,!opengl43,!angle CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo windows,angle CFLAGS: -DGRAPHICS_API_OPENGL_ES2
*/
import "C"
