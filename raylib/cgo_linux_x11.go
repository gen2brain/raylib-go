//go:build linux && ((x11 && !wayland) || termux) && !glfw && !glfw_x11 && !glfw_wayland && !rgfw && !drm && !sdl && !sdl3 && (!android || termux) && !memory && !nocgo
// +build linux
// +build x11,!wayland termux
// +build !glfw
// +build !glfw_x11
// +build !glfw_wayland
// +build !rgfw
// +build !drm
// +build !sdl
// +build !sdl3
// +build !android termux
// +build !memory
// +build !nocgo

package rl

/*
#cgo linux CFLAGS: -DPLATFORM_DESKTOP_X11
#cgo linux LDFLAGS: -lX11 -lXext -lXrandr -lXi -lXcursor -lxkbcommon -lm -ldl -lrt
#cgo linux,!software LDFLAGS: -lEGL
#cgo linux,termux,software LDFLAGS: -landroid-shmem
#cgo linux,es2,!software linux,es3,!software LDFLAGS: -lGLESv2
#cgo linux,arm64,!termux,!opengl11,!opengl21,!opengl33,!opengl43,!es2,!es3,!software LDFLAGS: -lGLESv2
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
