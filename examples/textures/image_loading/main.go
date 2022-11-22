package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - image loading")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)

	image := rl.LoadImage("raylib_logo.png")  // Loaded in CPU memory (RAM)
	texture := rl.LoadTextureFromImage(image) // Image converted to texture, GPU memory (VRAM)

	rl.UnloadImage(image) // Once image has been converted to texture and uploaded to VRAM, it can be unloaded from RAM

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2, rl.White)

		rl.DrawText("this IS a texture loaded from an image!", 300, 370, 10, rl.Gray)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)

	rl.CloseWindow()
}
