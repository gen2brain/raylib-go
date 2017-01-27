package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [textures] example - image loading")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)

	image := raylib.LoadImage("raylib_logo.png")  // Loaded in CPU memory (RAM)
	texture := raylib.LoadTextureFromImage(image) // Image converted to texture, GPU memory (VRAM)

	raylib.UnloadImage(image) // Once image has been converted to texture and uploaded to VRAM, it can be unloaded from RAM

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2, raylib.White)

		raylib.DrawText("this IS a texture loaded from an image!", 300, 370, 10, raylib.Gray)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(texture)

	raylib.CloseWindow()
}
