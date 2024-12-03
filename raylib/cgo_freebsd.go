//go:build freebsd && !linux && !rgfw && !drm && !sdl && !sdl3 && !android
// +build freebsd,!linux,!rgfw,!drm,!sdl,!sdl3,!android

package rl

/*
#include "external/glfw/src/context.c"
#include "external/glfw/src/init.c"
#include "external/glfw/src/input.c"
#include "external/glfw/src/monitor.c"
#include "external/glfw/src/platform.c"
#include "external/glfw/src/vulkan.c"
#include "external/glfw/src/window.c"

#include "external/glfw/src/x11_init.c"
#include "external/glfw/src/x11_monitor.c"
#include "external/glfw/src/x11_window.c"
#include "external/glfw/src/glx_context.c"

#include "external/glfw/src/null_joystick.c"
#include "external/glfw/src/posix_module.c"
#include "external/glfw/src/posix_poll.c"
#include "external/glfw/src/posix_thread.c"
#include "external/glfw/src/posix_time.c"
#include "external/glfw/src/xkb_unicode.c"
#include "external/glfw/src/egl_context.c"
#include "external/glfw/src/osmesa_context.c"

GLFWbool _glfwConnectNull(int platformID, _GLFWplatform* platform) {
	return GLFW_TRUE;
}

#cgo freebsd CFLAGS: -I. -I/usr/local/include -Iexternal/glfw/include -DPLATFORM_DESKTOP -D_GLFW_X11
#cgo freebsd LDFLAGS: -L/usr/local/lib -lm -pthread -ldl -lrt -lX11

#cgo freebsd,!es2,!es3 LDFLAGS: -lGL

#cgo freebsd,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo freebsd,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo freebsd,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo freebsd,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo freebsd,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo freebsd,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3
*/
import "C"
