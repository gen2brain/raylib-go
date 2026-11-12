//go:build linux && rgfw && !drm && !sdl && !sdl3 && (!android || termux) && !memory && !nocgo
// +build linux
// +build rgfw
// +build !drm
// +build !sdl
// +build !sdl3
// +build !android termux
// +build !memory
// +build !nocgo

package rl

/*
#cgo linux,!es2 LDFLAGS: -lm
#cgo linux CFLAGS: -DPLATFORM_DESKTOP_RGFW -Wno-builtin-declaration-mismatch -Wno-discarded-qualifiers -Wno-int-conversion
#cgo linux LDFLAGS: -lX11 -lXrandr -lXinerama -lXi -lXxf86vm -lXcursor -lm -lpthread -ldl -lrt

#cgo linux LDFLAGS: -lGL
#cgo linux,es2,!software linux,es3,!software LDFLAGS: -lGLESv2

#cgo linux,opengl11,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo linux,opengl21,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo linux,opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo linux,!opengl11,!opengl21,!opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo linux,es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo linux,es3,!es2,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES3
#cgo linux,software CFLAGS: -DGRAPHICS_API_OPENGL_SOFTWARE
*/
import "C"
