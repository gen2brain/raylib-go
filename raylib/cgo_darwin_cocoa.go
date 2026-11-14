//go:build darwin && cocoa && !rgfw && !sdl && !sdl3 && !memory && !nocgo
// +build darwin,cocoa,!rgfw,!sdl,!sdl3,!memory,!nocgo

package rl

/*
#cgo darwin LDFLAGS: -framework Cocoa -framework IOKit -framework CoreFoundation -framework QuartzCore
#cgo darwin CFLAGS: -x objective-c -std=gnu11 -DPLATFORM_DESKTOP_COCOA -Wno-deprecated-declarations -Wno-implicit-const-int-float-conversion

#cgo darwin,!es2,!es3,!software LDFLAGS: -framework OpenGL

#cgo darwin,opengl11,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo darwin,opengl21,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo darwin,opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo darwin,!opengl11,!opengl21,!opengl43,!es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo darwin,es2,!es3,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo darwin,es3,!es2,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES3
#cgo darwin,software CFLAGS: -DGRAPHICS_API_OPENGL_SOFTWARE
*/
import "C"
