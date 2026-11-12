//go:build freebsd && rgfw && !linux && !drm && !sdl && !sdl3 && !android && !memory && !nocgo
// +build freebsd,rgfw,!linux,!drm,!sdl,!sdl3,!android,!memory,!nocgo

package rl

/*
#cgo freebsd CFLAGS: -I. -I/usr/local/include -DPLATFORM_DESKTOP_RGFW
#cgo freebsd LDFLAGS: -L/usr/local/lib -lX11 -lXrandr -lXinerama -lXi -lXxf86vm -lXcursor -lm -lpthread -ldl -lrt

#cgo freebsd,!es2,!es3 LDFLAGS: -lGL

#cgo freebsd,opengl11,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo freebsd,opengl21,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo freebsd,opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo freebsd,!opengl11,!opengl21,!opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo freebsd,es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo freebsd,es3,!es2,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES3
#cgo freebsd,software CFLAGS: -DGRAPHICS_API_OPENGL_SOFTWARE
*/
import "C"
