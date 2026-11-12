//go:build linux && (sdl || sdl3) && !rgfw && !drm && (!android || termux) && !memory && !nocgo
// +build linux
// +build sdl sdl3
// +build !rgfw
// +build !drm
// +build !android termux
// +build !memory
// +build !nocgo

package rl

/*
#cgo linux,!es2 LDFLAGS: -lm
#cgo linux CFLAGS: -Wno-stringop-overflow
#cgo linux,sdl CFLAGS: -DPLATFORM_DESKTOP_SDL
#cgo linux,sdl3 CFLAGS: -DPLATFORM_DESKTOP_SDL -DPLATFORM_DESKTOP_SDL3 -DUSING_SDL3_PROJECT
#cgo linux,sdl pkg-config: sdl2
#cgo linux,sdl3 pkg-config: sdl3

#cgo linux,!es2,!es3,!software LDFLAGS: -lGL

#cgo linux,opengl11,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo linux,opengl21,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo linux,opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo linux,!opengl11,!opengl21,!opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo linux,es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo linux,es3,!es2,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES3
#cgo linux,software CFLAGS: -DGRAPHICS_API_OPENGL_SOFTWARE
*/
import "C"
