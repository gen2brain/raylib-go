//go:build freebsd && !linux && sdl2 && !rgfw && !drm && !android
// +build freebsd,!linux,sdl2,!rgfw,!drm,!android

package rl

/*
#cgo freebsd CFLAGS: -I. -I/usr/local/include
#cgo freebsd,sdl2 CFLAGS: -DPLATFORM_DESKTOP_SDL -DUSING_SDL2_PROJECT
#cgo freebsd LDFLAGS: -L/usr/local/lib

#cgo freebsd,sdl2 pkg-config: sdl2

#cgo freebsd,!es2,!es3 LDFLAGS: -lGL

#cgo freebsd,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo freebsd,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo freebsd,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo freebsd,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo freebsd,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo freebsd,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3
*/
import "C"
