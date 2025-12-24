//go:build openbsd && !linux && (sdl || sdl3) && !rgfw && !drm && !android
// +build openbsd
// +build !linux
// +build sdl sdl3
// +build !rgfw
// +build !drm
// +build !android

package rl

/*
#cgo openbsd CFLAGS: -I. -I/usr/X11R6/include
#cgo openbsd,sdl CFLAGS: -DPLATFORM_DESKTOP_SDL
#cgo openbsd,sdl3 CFLAGS: -DPLATFORM_DESKTOP_SDL -DPLATFORM_DESKTOP_SDL3
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib

#cgo openbsd,sdl pkg-config: sdl2
#cgo openbsd,sdl3 pkg-config: sdl3

#cgo openbsd,!es2,!es3 LDFLAGS: -lGL

#cgo openbsd,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo openbsd,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo openbsd,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo openbsd,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo openbsd,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo openbsd,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3
*/
import "C"
