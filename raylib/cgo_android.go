//go:build android
// +build android

package rl

/*
#include "external/android/native_app_glue/android_native_app_glue.c"

#cgo android LDFLAGS: -llog -landroid -lEGL -lGLESv2 -lOpenSLES -lm
#cgo android CFLAGS: -DPLATFORM_ANDROID -DGRAPHICS_API_OPENGL_ES2 -Iexternal/android/native_app_glue -Wno-implicit-const-int-float-conversion

#cgo android,arm CFLAGS: -march=armv7-a -mfloat-abi=softfp -mfpu=vfpv3-d16
*/
import "C"
