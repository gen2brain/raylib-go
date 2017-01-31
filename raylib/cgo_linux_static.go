// +build linux,static,!arm

package raylib

/*
#cgo linux LDFLAGS: -lglfw3 -lGL -lopenal -lm -pthread -ldl -lX11 -lXrandr -lXinerama -lXi -lXxf86vm -lXcursor
#cgo linux CFLAGS: -DPLATFORM_DESKTOP -DGRAPHICS_API_OPENGL_33 -DAL_LIBTYPE_STATIC
*/
import "C"
