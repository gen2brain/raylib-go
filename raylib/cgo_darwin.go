//go:build darwin
// +build darwin

package rl

/*
#include "external/glfw/src/context.c"
#include "external/glfw/src/init.c"
#include "external/glfw/src/input.c"
#include "external/glfw/src/monitor.c"
#include "external/glfw/src/vulkan.c"
#include "external/glfw/src/window.c"

#include "external/glfw/src/cocoa_init.m"
#include "external/glfw/src/cocoa_joystick.m"
#include "external/glfw/src/cocoa_monitor.m"
#include "external/glfw/src/cocoa_time.c"
#include "external/glfw/src/cocoa_window.m"
#include "external/glfw/src/posix_thread.c"
#include "external/glfw/src/nsgl_context.m"
#include "external/glfw/src/egl_context.c"
#include "external/glfw/src/osmesa_context.c"

#cgo darwin LDFLAGS: -framework OpenGL -framework Cocoa -framework IOKit -framework CoreVideo -framework CoreFoundation
#cgo darwin CFLAGS: -x objective-c -Iexternal/glfw/include -D_GLFW_COCOA -D_GLFW_USE_CHDIR -D_GLFW_USE_MENUBAR -D_GLFW_USE_RETINA -Wno-deprecated-declarations -DPLATFORM_DESKTOP

#cgo darwin,opengl11 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo darwin,opengl21 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo darwin,opengl43 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo darwin,!opengl11,!opengl21,!opengl43 CFLAGS: -DGRAPHICS_API_OPENGL_33
*/
import "C"
