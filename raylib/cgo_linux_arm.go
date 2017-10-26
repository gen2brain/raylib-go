// +build linux,arm,!android

package raylib

/*
#cgo linux,arm LDFLAGS: -lGLESv2 -lEGL -lpthread -lrt -lm -lbcm_host -lvcos -lvchiq_arm -lopenal
#cgo linux,arm CFLAGS: -DPLATFORM_RPI -DGRAPHICS_API_OPENGL_ES2 -I/opt/vc/include -I/opt/vc/include/interface/vcos -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
*/
import "C"
