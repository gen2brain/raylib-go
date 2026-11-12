//go:build darwin && (sdl || sdl3) && !rgfw && !memory && !nocgo
// +build darwin
// +build sdl sdl3
// +build !rgfw
// +build !memory
// +build !nocgo

package rl

/*
#cgo darwin LDFLAGS: -framework Cocoa -framework IOKit -framework CoreVideo -framework CoreFoundation
#cgo darwin CFLAGS: -Wno-deprecated-declarations -Wno-implicit-const-int-float-conversion
#cgo darwin,sdl CFLAGS: -DPLATFORM_DESKTOP_SDL
#cgo darwin,sdl3 CFLAGS: -DPLATFORM_DESKTOP_SDL -DPLATFORM_DESKTOP_SDL3 -DUSING_SDL3_PROJECT
#cgo darwin,sdl pkg-config: sdl2
#cgo darwin,sdl3 pkg-config: sdl3

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
