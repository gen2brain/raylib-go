//go:build freebsd && !linux && !drm && sdl && !android
// +build freebsd,!linux,!drm,sdl,!android

package rl

/*
#cgo freebsd CFLAGS: -I. -I/usr/local/include -DPLATFORM_DESKTOP_SDL
#cgo freebsd LDFLAGS: -L/usr/local/lib

#cgo freebsd LDFLAGS: -lm -pthread -ldl -lrt -lX11
#cgo freebsd pkg-config: sdl2

#cgo freebsd,!angle LDFLAGS: -lGL

#cgo freebsd,opengl11,!angle CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo freebsd,opengl21,!angle CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo freebsd,opengl43,!angle CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo freebsd,!opengl11,!opengl21,!opengl43,!angle CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo freebsd,angle CFLAGS: -DGRAPHICS_API_OPENGL_ES2
*/
import "C"
