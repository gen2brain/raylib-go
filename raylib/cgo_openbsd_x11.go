//go:build openbsd && !glfw && !linux && !rgfw && !drm && !sdl && !sdl3 && !android && !memory && !nocgo
// +build openbsd,!glfw,!linux,!rgfw,!drm,!sdl,!sdl3,!android,!memory,!nocgo

package rl

/*
#cgo openbsd CFLAGS: -I. -I/usr/X11R6/include -I/usr/local/include -DPLATFORM_DESKTOP_X11
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib -L/usr/local/lib -lX11 -lXext -lXrandr -lXi -lXcursor -lxkbcommon -lm -pthread
#cgo openbsd,!software LDFLAGS: -lEGL
#cgo openbsd,es2,!software openbsd,es3,!software LDFLAGS: -lGLESv2
#cgo openbsd,opengl11,!es2,!es3,!software LDFLAGS: -lGL

#cgo openbsd,opengl11,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo openbsd,opengl21,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo openbsd,opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo openbsd,!opengl11,!opengl21,!opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo openbsd,es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo openbsd,es3,!es2,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES3
#cgo openbsd,software CFLAGS: -DGRAPHICS_API_OPENGL_SOFTWARE
*/
import "C"
