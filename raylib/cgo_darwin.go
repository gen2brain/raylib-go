// +build darwin

package raylib

/*
#cgo darwin LDFLAGS: -lglfw -framework OpenGL -framework OpenAL -framework Cocoa
#cgo darwin CFLAGS: -DPLATFORM_DESKTOP -DGRAPHICS_API_OPENGL_33
*/
import "C"
