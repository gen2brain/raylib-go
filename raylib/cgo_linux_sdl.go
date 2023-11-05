//go:build linux && !drm && sdl && !android
// +build linux,!drm,sdl,!android

package rl

/*
#cgo linux,!angle LDFLAGS: -lm
#cgo linux CFLAGS: -DPLATFORM_DESKTOP_SDL -Wno-stringop-overflow
#cgo linux pkg-config: sdl2

#cgo linux,!angle LDFLAGS: -lGL

#cgo linux,opengl11,!angle CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo linux,opengl21,!angle CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo linux,opengl43,!angle CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo linux,!opengl11,!opengl21,!opengl43,!angle CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo linux,angle CFLAGS: -DGRAPHICS_API_OPENGL_ES2
*/
import "C"
