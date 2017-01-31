// +build windows,!static

package raylib

/*
#cgo windows LDFLAGS: -lglfw3 -lopengl32 -lgdi32 -lOpenAL32 -lwinmm -lole32
#cgo windows CFLAGS: -DPLATFORM_DESKTOP -DGRAPHICS_API_OPENGL_33 -DSHARED_OPENAL
*/
import "C"
