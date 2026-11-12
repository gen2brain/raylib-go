//go:build android && !termux && !nocgo
// +build android,!termux,!nocgo

package rl

/*
#include "external/android/native_app_glue/android_native_app_glue.c"

#cgo android LDFLAGS: -llog -landroid -lEGL -lGLESv2 -lOpenSLES -lm -Wl,--wrap=fopen
#cgo android CFLAGS: -DPLATFORM_ANDROID -Iexternal/android/native_app_glue -Wno-implicit-const-int-float-conversion
#cgo android,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo android,es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES3
#cgo android,es3 LDFLAGS: -lGLESv3

#cgo android,arm CFLAGS: -march=armv7-a -mfloat-abi=softfp -mfpu=vfpv3-d16
*/
import "C"
