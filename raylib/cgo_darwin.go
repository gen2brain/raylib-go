// +build darwin

package raylib

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
#include "external/glfw/src/cocoa_window.m"
#include "external/glfw/src/cocoa_time.c"
#include "external/glfw/src/posix_tls.c"
#include "external/glfw/src/nsgl_context.m"

#cgo darwin LDFLAGS: -framework OpenGL -framework OpenAL -framework Cocoa -framework IOKit -framework CoreVideo
#cgo darwin CFLAGS: -x objective-c -Iexternal/glfw/include -D_GLFW_COCOA -D_GLFW_USE_CHDIR -D_GLFW_USE_MENUBAR -D_GLFW_USE_RETINA -Wno-deprecated-declarations -DPLATFORM_DESKTOP -DGRAPHICS_API_OPENGL_33
*/
import "C"
