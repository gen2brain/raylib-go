//go:build ios && !nocgo
// +build ios,!nocgo

package rl

/*
#cgo ios LDFLAGS: -framework UIKit -framework Foundation -framework CoreFoundation -framework CoreGraphics -framework QuartzCore -framework OpenGLES
#cgo ios LDFLAGS: -framework AVFoundation -framework AudioToolbox -framework CoreAudio -framework GameController -framework CoreHaptics
#cgo ios CFLAGS: -x objective-c -std=gnu11 -DPLATFORM_UIKIT -DGLES_SILENCE_DEPRECATION -DFBO_READ_TEXTURE_AS_RGBA -Iexternal/ios -Wno-deprecated-declarations -Wno-implicit-const-int-float-conversion

#cgo ios,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo ios,es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES3
#cgo ios,software CFLAGS: -DGRAPHICS_API_OPENGL_SOFTWARE
*/
import "C"
