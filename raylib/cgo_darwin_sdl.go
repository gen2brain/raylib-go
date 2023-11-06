//go:build darwin && sdl
// +build darwin,sdl

package rl

/*
#cgo darwin LDFLAGS: -framework Cocoa -framework IOKit -framework CoreVideo -framework CoreFoundation
#cgo darwin CFLAGS: -x objective-c -Wno-deprecated-declarations -Wno-implicit-const-int-float-conversion -DPLATFORM_DESKTOP_SDL
#cgo darwin pkg-config: sdl2

#cgo darwin,!angle LDFLAGS: -framework OpenGL

#cgo darwin,opengl11,!angle CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo darwin,opengl21,!angle CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo darwin,opengl43,!angle CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo darwin,!opengl11,!opengl21,!opengl43,!angle CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo darwin,angle CFLAGS: -DGRAPHICS_API_OPENGL_ES2
*/
import "C"
