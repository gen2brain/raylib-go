//go:build darwin && sdl
// +build darwin,sdl

package rl

/*
#cgo darwin LDFLAGS: -framework Cocoa -framework IOKit -framework CoreVideo -framework CoreFoundation
#cgo darwin CFLAGS: -x objective-c -Wno-deprecated-declarations -Wno-implicit-const-int-float-conversion -DPLATFORM_DESKTOP_SDL
#cgo darwin pkg-config: sdl2

#cgo darwin,!es2 LDFLAGS: -framework OpenGL

#cgo darwin,opengl11,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo darwin,opengl21,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo darwin,opengl43,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo darwin,!opengl11,!opengl21,!opengl43,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo darwin,es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
*/
import "C"
