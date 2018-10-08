// +build linux,arm,!android

package rl

/*
#cgo linux,arm LDFLAGS: -L/opt/vc/lib -L/opt/vc/lib64 -lbrcmGLESv2 -lbrcmEGL -lpthread -lrt -lm -lbcm_host -lvcos -lvchiq_arm -ldl
#cgo linux,arm CFLAGS: -DPLATFORM_RPI -DGRAPHICS_API_OPENGL_ES2 -Iexternal -I/opt/vc/include -I/opt/vc/include/interface/vcos -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
*/
import "C"
