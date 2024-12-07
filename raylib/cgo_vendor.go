//go:build required
// +build required

package rl

import (
	_ "github.com/gen2brain/raylib-go/raylib/external"
	_ "github.com/gen2brain/raylib-go/raylib/external/android/native_app_glue"
	_ "github.com/gen2brain/raylib-go/raylib/external/glfw/include/GLFW"
	_ "github.com/gen2brain/raylib-go/raylib/external/glfw/src"
	_ "github.com/gen2brain/raylib-go/raylib/platforms"
)
