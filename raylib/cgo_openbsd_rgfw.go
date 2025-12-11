//go:build openbsd && rgfw && !linux && !sdl && !sdl3 && !drm && !android
// +build openbsd,rgfw,!linux,!sdl,!sdl3,!drm,!android

package rl

/*
#cgo openbsd CFLAGS: -I. -I/usr/X11R6/include -DPLATFORM_DESKTOP_RGFW
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib -lX11 -lXrandr -lXinerama -lXi -lXxf86vm -lXcursor -lm -lpthread -ldl -lrt

#cgo openbsd,!es2,!es3 LDFLAGS: -lGL

#cgo openbsd,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo openbsd,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo openbsd,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo openbsd,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo openbsd,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo openbsd,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3
*/
import "C"
