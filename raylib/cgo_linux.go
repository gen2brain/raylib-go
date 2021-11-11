//go:build linux && !drm && !rpi && !android
// +build linux,!drm,!rpi,!android

package rl

/*
#include "external/glfw/src/context.c"
#include "external/glfw/src/init.c"
#include "external/glfw/src/input.c"
#include "external/glfw/src/monitor.c"
#include "external/glfw/src/vulkan.c"
#include "external/glfw/src/window.c"

#ifdef _GLFW_WAYLAND
#include "external/glfw/src/wl_init.c"
#include "external/glfw/src/wl_monitor.c"
#include "external/glfw/src/wl_window.c"
#include "external/glfw/src/wayland-pointer-constraints-unstable-v1-client-protocol.c"
#include "external/glfw/src/wayland-relative-pointer-unstable-v1-client-protocol.c"
#include "external/glfw/src/wayland-idle-inhibit-unstable-v1-client-protocol.c"
#include "external/glfw/src/wayland-xdg-shell-client-protocol.c"
#include "external/glfw/src/wayland-xdg-decoration-client-protocol.c"
#include "external/glfw/src/wayland-viewporter-client-protocol.c"
#endif
#ifdef _GLFW_X11
#include "external/glfw/src/x11_init.c"
#include "external/glfw/src/x11_monitor.c"
#include "external/glfw/src/x11_window.c"
#include "external/glfw/src/glx_context.c"
#endif

#include "external/glfw/src/linux_joystick.c"
#include "external/glfw/src/posix_thread.c"
#include "external/glfw/src/posix_time.c"
#include "external/glfw/src/xkb_unicode.c"
#include "external/glfw/src/egl_context.c"
#include "external/glfw/src/osmesa_context.c"

#cgo linux CFLAGS: -Iexternal/glfw/include -DPLATFORM_DESKTOP -Wno-stringop-overflow

#cgo linux,!wayland LDFLAGS: -lGL -lm -pthread -ldl -lrt -lX11
#cgo linux,wayland LDFLAGS: -lGL -lm -pthread -ldl -lrt -lwayland-client -lwayland-cursor -lwayland-egl -lxkbcommon

#cgo linux,!wayland CFLAGS: -D_GLFW_X11
#cgo linux,wayland CFLAGS: -D_GLFW_WAYLAND

#cgo linux,opengl11 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo linux,opengl21 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo linux,opengl43 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo linux,!opengl11,!opengl21,!opengl43 CFLAGS: -DGRAPHICS_API_OPENGL_33
*/
import "C"
