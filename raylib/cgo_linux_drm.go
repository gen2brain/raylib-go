//go:build linux && drm && !rgfw && !sdl && !sdl3 && !android && !memory && !nocgo
// +build linux,drm,!rgfw,!sdl,!sdl3,!android,!memory,!nocgo

package rl

/*
#cgo linux,drm LDFLAGS: -ldrm -lpthread -lrt -lm -ldl
#cgo linux,drm,!software LDFLAGS: -lGLESv2 -lEGL -lgbm
#cgo linux,drm CFLAGS: -DPLATFORM_DRM -I/usr/include/libdrm -include stdint.h -Duint=uint32_t
#cgo linux,drm,!software CFLAGS: -DGRAPHICS_API_OPENGL_ES2 -DEGL_NO_X11
#cgo linux,drm,software CFLAGS: -DGRAPHICS_API_OPENGL_SOFTWARE
*/
import "C"
