//go:build openbsd && !linux && !drm && !sdl && !android
// +build openbsd,!linux,!drm,!sdl,!android

package rl

/*
#include "external/glfw/src/context.c"
#include "external/glfw/src/init.c"
#include "external/glfw/src/input.c"
#include "external/glfw/src/monitor.c"
#include "external/glfw/src/platform.c"
#include "external/glfw/src/vulkan.c"
#include "external/glfw/src/window.c"

#ifdef _GLFW_WAYLAND
#include "external/glfw/src/wl_init.c"
#include "external/glfw/src/wl_monitor.c"
#include "external/glfw/src/wl_window.c"
#endif
#ifdef _GLFW_X11
#include "external/glfw/src/x11_init.c"
#include "external/glfw/src/x11_monitor.c"
#include "external/glfw/src/x11_window.c"
#include "external/glfw/src/glx_context.c"
#endif

#include "external/glfw/src/null_joystick.c"
#include "external/glfw/src/posix_module.c"
#include "external/glfw/src/posix_poll.c"
#include "external/glfw/src/posix_thread.c"
#include "external/glfw/src/posix_time.c"
#include "external/glfw/src/xkb_unicode.c"
#include "external/glfw/src/egl_context.c"
#include "external/glfw/src/osmesa_context.c"

#cgo openbsd CFLAGS: -I. -I/usr/X11R6/include -Iexternal/glfw/include -DPLATFORM_DESKTOP
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib

#cgo openbsd,!wayland LDFLAGS: -lm -pthread -lX11
#cgo openbsd,wayland LDFLAGS: -lm -pthread -lwayland-client -lwayland-cursor -lwayland-egl -lxkbcommon

#cgo openbsd,!es2,!es3 LDFLAGS: -lGL

#cgo openbsd,!wayland CFLAGS: -D_GLFW_X11
#cgo openbsd,wayland CFLAGS: -D_GLFW_WAYLAND

#cgo openbsd,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo openbsd,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo openbsd,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo openbsd,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo openbsd,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo openbsd,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3
*/
import "C"
