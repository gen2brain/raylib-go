// +build linux,arm,!android

package raylib

/*
#cgo linux,arm LDFLAGS: -lbrcmGLESv2 -lbrcmEGL -lpthread -lrt -lm -lbcm_host -lvcos -lvchiq_arm
#cgo linux,arm CFLAGS: -DPLATFORM_RPI -DGRAPHICS_API_OPENGL_ES2 -I/opt/vc/include -I/opt/vc/include/interface/vcos -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads

#cgo linux,arm,!noaudio LDFLAGS: -lopenal
*/
import "C"
