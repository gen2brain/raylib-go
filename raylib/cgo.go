package rl

/*
#include "external/glfw/src/context.c"
#include "external/glfw/src/init.c"
#include "external/glfw/src/input.c"
#include "external/glfw/src/monitor.c"
#include "external/glfw/src/vulkan.c"
#include "external/glfw/src/window.c"

// OSMesa for headless GLFW functions
#include "external/glfw/src/null_init.c"
#include "external/glfw/src/null_monitor.c"
#include "external/glfw/src/null_window.c"
#include "external/glfw/src/null_joystick.c"
#include "external/glfw/src/osmesa_context.c"

#include "external/glfw/src/posix_thread.c"
#include "external/glfw/src/posix_time.c"
#include "external/glfw/src/xkb_unicode.c"

#cgo CFLAGS: -Iexternal/glfw/include -DPLATFORM_DESKTOP -Wno-stringop-overflow  -D_GLFW_OSMESA
#cgo LDFLAGS: -lGL -lm -pthread -ldl -lrt -lOSMesa

#cgo opengl11 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo opengl21 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo opengl43 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo !opengl11,!opengl21,!opengl43 CFLAGS: -DGRAPHICS_API_OPENGL_33

#cgo CFLAGS: -std=gnu99 -Wno-missing-braces -Wno-unused-result -Wno-implicit-function-declaration

*/
import "C"
