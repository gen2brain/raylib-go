// +build android

package raylib

/*
#cgo android LDFLAGS: -llog -landroid -lEGL -lGLESv2 -lOpenSLES -lopenal -lm -landroid_native_app_glue
#cgo android CFLAGS: -DPLATFORM_ANDROID -DGRAPHICS_API_OPENGL_ES2
*/
import "C"
