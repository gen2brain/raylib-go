//go:build linux && wayland && !rgfw && !drm && !sdl && !sdl3 && !android && !memory && !nocgo
// +build linux
// +build wayland
// +build !rgfw
// +build !drm
// +build !sdl
// +build !sdl3
// +build !android
// +build !memory
// +build !nocgo

package rl

/*
#cgo linux CFLAGS: -DPLATFORM_DESKTOP_WAYLAND
#cgo linux LDFLAGS: -lwayland-client -lwayland-cursor -lxkbcommon -lm -ldl -lrt
#cgo linux,!software LDFLAGS: -lwayland-egl -lEGL
#cgo linux,es2,!software linux,es3,!software LDFLAGS: -lGLESv2
#cgo linux,opengl11,!es2,!es3,!software LDFLAGS: -lGL

#cgo linux,opengl11,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo linux,opengl21,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo linux,opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo linux,!opengl11,!opengl21,!opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo linux,es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo linux,es3,!es2,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES3
#cgo linux,software CFLAGS: -DGRAPHICS_API_OPENGL_SOFTWARE
*/
import "C"
