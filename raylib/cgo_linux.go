//go:build linux && !rgfw && !drm && !sdl && !sdl3 && (!android || termux) && !memory && !nocgo
// +build linux
// +build !rgfw
// +build !drm
// +build !sdl
// +build !sdl3
// +build !android termux
// +build !memory
// +build !nocgo

package rl

/*
#if (_POSIX_C_SOURCE < 199309L)
#undef _POSIX_C_SOURCE
#define _POSIX_C_SOURCE 199309L
#endif
#if !defined(_GNU_SOURCE)
#define _GNU_SOURCE
#endif
#include <features.h>
#if defined(__GLIBC__)
#undef __GLIBC_USE_C2X_STRTOL
#define __GLIBC_USE_C2X_STRTOL 0
#undef __GLIBC_USE_C23_STRTOL
#define __GLIBC_USE_C23_STRTOL 0
#endif

#include "external/glfw/src/init.c"
#include "external/glfw/src/platform.c"
#include "external/glfw/src/context.c"
#include "external/glfw/src/monitor.c"
#include "external/glfw/src/window.c"
#include "external/glfw/src/input.c"
#include "external/glfw/src/vulkan.c"
#include "external/glfw/src/egl_context.c"
#include "external/glfw/src/osmesa_context.c"

GLFWbool _glfwConnectNull(int platformID, _GLFWplatform* platform) {
	return GLFW_FALSE;
}

#include "external/glfw/src/posix_module.c"
#include "external/glfw/src/posix_thread.c"
#include "external/glfw/src/posix_time.c"
#include "external/glfw/src/posix_poll.c"
#include "external/glfw/src/linux_joystick.c"
#include "external/glfw/src/xkb_unicode.c"

#if defined(_GLFW_WAYLAND)
#define createKeyTables createKeyTablesWayland
#define translateKey translateKeyWayland
#define acquireMonitor acquireMonitorWayland
#define releaseMonitor releaseMonitorWayland
#include "external/glfw/src/wl_init.c"
#include "external/glfw/src/wl_monitor.c"
#include "external/glfw/src/wl_window.c"
#undef createKeyTables
#undef translateKey
#undef acquireMonitor
#undef releaseMonitor
#endif

#if defined(_GLFW_X11)
#define createKeyTables createKeyTablesX11
#define translateKey translateKeyX11
#define acquireMonitor acquireMonitorX11
#define releaseMonitor releaseMonitorX11
#include "external/glfw/src/x11_init.c"
#include "external/glfw/src/x11_monitor.c"
#include "external/glfw/src/x11_window.c"
#include "external/glfw/src/glx_context.c"
#undef createKeyTables
#undef translateKey
#undef acquireMonitor
#undef releaseMonitor
#endif

#cgo linux CFLAGS: -Iexternal/glfw/include -DPLATFORM_DESKTOP_GLFW -Wno-stringop-overflow
#cgo linux LDFLAGS: -lm -pthread -ldl -lrt -lxkbcommon
#cgo linux,!glfw_x11,!termux LDFLAGS: -lwayland-client -lwayland-cursor -lwayland-egl
#cgo linux,!glfw_wayland linux,termux LDFLAGS: -lX11

#cgo linux,glfw_x11 linux,termux CFLAGS: -D_GLFW_X11
#cgo linux,glfw_wayland,!termux CFLAGS: -D_GLFW_WAYLAND
#cgo linux,!glfw_x11,!glfw_wayland,!termux CFLAGS: -D_GLFW_X11 -D_GLFW_WAYLAND

#cgo linux,!es2,!es3,!glfw_wayland linux,!es2,!es3,termux LDFLAGS: -lGL
#cgo linux,es3 LDFLAGS: -lGLESv2

#cgo linux,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo linux,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo linux,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo linux,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo linux,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo linux,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3
*/
import "C"
