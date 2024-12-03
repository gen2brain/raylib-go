//go:build darwin && !rgfw && !sdl && !sdl3
// +build darwin,!rgfw,!sdl,!sdl3

package rl

/*
#include "external/glfw/src/context.c"
#include "external/glfw/src/init.c"
#include "external/glfw/src/input.c"
#include "external/glfw/src/monitor.c"
#include "external/glfw/src/platform.c"
#include "external/glfw/src/vulkan.c"
#include "external/glfw/src/window.c"

#include "external/glfw/src/cocoa_init.m"
#include "external/glfw/src/cocoa_joystick.m"
#include "external/glfw/src/cocoa_monitor.m"
#include "external/glfw/src/cocoa_time.c"
#include "external/glfw/src/cocoa_window.m"
#include "external/glfw/src/posix_module.c"
#include "external/glfw/src/posix_thread.c"
#include "external/glfw/src/nsgl_context.m"
#include "external/glfw/src/egl_context.c"
#include "external/glfw/src/osmesa_context.c"

GLFWbool _glfwConnectNull(int platformID, _GLFWplatform* platform) {
	return GLFW_TRUE;
}

#cgo darwin LDFLAGS: -framework Cocoa -framework IOKit -framework CoreVideo -framework CoreFoundation
#cgo darwin CFLAGS: -x objective-c -Iexternal/glfw/include -D_GLFW_COCOA -D_GLFW_USE_CHDIR -D_GLFW_USE_MENUBAR -D_GLFW_USE_RETINA -Wno-deprecated-declarations -Wno-implicit-const-int-float-conversion -DPLATFORM_DESKTOP

#cgo darwin,!es2,!es3 LDFLAGS: -framework OpenGL

#cgo darwin,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo darwin,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo darwin,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo darwin,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo darwin,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo darwin,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3
*/
import "C"
