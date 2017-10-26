// +build linux,!arm,!arm64

package raylib

/*
#include "external/glfw/src/context.c"
#include "external/glfw/src/init.c"
#include "external/glfw/src/input.c"
#include "external/glfw/src/monitor.c"
#include "external/glfw/src/vulkan.c"
#include "external/glfw/src/window.c"

#include "external/glfw/src/x11_init.c"
#include "external/glfw/src/x11_monitor.c"
#include "external/glfw/src/x11_window.c"
#include "external/glfw/src/glx_context.c"
#include "external/glfw/src/linux_joystick.c"
#include "external/glfw/src/posix_time.c"
#include "external/glfw/src/posix_tls.c"
#include "external/glfw/src/xkb_unicode.c"
#include "external/glfw/src/egl_context.c"

#cgo linux LDFLAGS: -lGL -lm -pthread -ldl -lrt -lX11 -lXrandr -lXinerama -lXi -lXxf86vm -lXcursor
#cgo linux CFLAGS: -D_GLFW_X11 -Iexternal/glfw/include -DPLATFORM_DESKTOP

#cgo linux,!noaudio LDFLAGS: -lopenal

#cgo linux,!static CFLAGS: -DSHARED_OPENAL
#cgo linux,static CFLAGS: -DAL_LIBTYPE_STATIC


#cgo linux,opengl11 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo linux,opengl21 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo linux,!opengl11,!opengl21 CFLAGS: -DGRAPHICS_API_OPENGL_33
*/
import "C"
