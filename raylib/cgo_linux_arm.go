// +build linux,arm,!android

package raylib

/*
#cgo linux,arm LDFLAGS: -lGLESv2 -lEGL -lpthread -lrt -lm -lbcm_host -lvcos -lvchiq_arm -lopenal
#cgo linux,arm CFLAGS: -DPLATFORM_RPI -DGRAPHICS_API_OPENGL_ES2
*/
import "C"
