//go:build linux && wayland && !x11 && !glfw && !glfw_x11 && !glfw_wayland && !rgfw && !drm && !sdl && !sdl3 && !android && !memory && !nocgo
// +build linux,wayland,!x11,!glfw,!glfw_x11,!glfw_wayland,!rgfw,!drm,!sdl,!sdl3,!android,!memory,!nocgo

package rl

/*
#cgo linux CFLAGS: -DPLATFORM_DESKTOP_WAYLAND
#cgo linux LDFLAGS: -lm -ldl -lrt
#cgo linux,es3,!es2,!software LDFLAGS: -lGLESv2
#cgo linux,opengl11,!es2,!es3,!software LDFLAGS: -lGL

#cgo linux,opengl11,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo linux,opengl21,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo linux,opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo linux,!arm64,!opengl11,!opengl21,!opengl43,!es2,!es3,!software linux,termux,!opengl11,!opengl21,!opengl43,!es2,!es3,!software linux,opengl33,!opengl11,!opengl21,!opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo linux,arm64,!termux,!opengl11,!opengl21,!opengl33,!opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo linux,es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo linux,es3,!es2,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES3
#cgo linux,software CFLAGS: -DGRAPHICS_API_OPENGL_SOFTWARE
*/
import "C"
