//go:build darwin && sdl2 && !rgfw
// +build darwin,sdl2,!rgfw

package rl

/*
#cgo darwin LDFLAGS: -framework Cocoa -framework IOKit -framework CoreVideo -framework CoreFoundation
#cgo darwin CFLAGS: -Wno-deprecated-declarations -Wno-implicit-const-int-float-conversion
#cgo darwin,sdl2 CFLAGS: -DPLATFORM_DESKTOP_SDL -DUSING_SDL2_PROJECT
#cgo darwin,sdl2 pkg-config: sdl2

#cgo darwin,!es2,!es3 LDFLAGS: -framework OpenGL

#cgo darwin,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo darwin,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo darwin,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo darwin,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo darwin,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo darwin,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3
*/
import "C"
