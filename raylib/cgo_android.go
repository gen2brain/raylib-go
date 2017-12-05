// +build android

package raylib

/*
#cgo android LDFLAGS: -llog -landroid -lEGL -lGLESv2 -lOpenSLES -lm -landroid_native_app_glue -u ANativeActivity_onCreate
#cgo android CFLAGS: -DPLATFORM_ANDROID -DGRAPHICS_API_OPENGL_ES2 -Iexternal
*/
import "C"
