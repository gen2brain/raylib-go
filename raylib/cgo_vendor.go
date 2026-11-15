//go:build required
// +build required

package rl

import (
	_ "github.com/gen2brain/raylib-go/raylib/external"
	_ "github.com/gen2brain/raylib-go/raylib/external/android/native_app_glue"
	_ "github.com/gen2brain/raylib-go/raylib/external/glfw/include/GLFW"
	_ "github.com/gen2brain/raylib-go/raylib/external/glfw/src"
	_ "github.com/gen2brain/raylib-go/raylib/external/ios"
	_ "github.com/gen2brain/raylib-go/raylib/external/ios/GLES2"
	_ "github.com/gen2brain/raylib-go/raylib/external/ios/GLES3"
	_ "github.com/gen2brain/raylib-go/raylib/external/wayland"
	_ "github.com/gen2brain/raylib-go/raylib/platforms"
)
