//go:build openbsd && !linux && !drm && sdl && !android
// +build openbsd,!linux,!drm,sdl,!android

package rl

/*
#cgo openbsd CFLAGS: -I. -I/usr/X11R6/include -DPLATFORM_DESKTOP_SDL
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib

#cgo openbsd LDFLAGS: -lm -pthread -lX11
#cgo openbsd pkg-config: sdl2

#cgo openbsd,!angle LDFLAGS: -lGL

#cgo openbsd,opengl11,!angle CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo openbsd,opengl21,!angle CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo openbsd,opengl43,!angle CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo openbsd,!opengl11,!opengl21,!opengl43,!angle CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo openbsd,angle CFLAGS: -DGRAPHICS_API_OPENGL_ES2
*/
import "C"
