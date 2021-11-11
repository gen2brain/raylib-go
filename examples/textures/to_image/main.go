package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - texture to image")

	image := rl.LoadImage("raylib_logo.png")  // Load image data into CPU memory (RAM)
	texture := rl.LoadTextureFromImage(image) // Image converted to texture, GPU memory (RAM -> VRAM)
	rl.UnloadImage(image)                     // Unload image data from CPU memory (RAM)

	image = rl.LoadImageFromTexture(texture) // Retrieve image data from GPU memory (VRAM -> RAM)
	rl.UnloadTexture(texture)                // Unload texture from GPU memory (VRAM)

	texture = rl.LoadTextureFromImage(image) // Recreate texture from retrieved image data (RAM -> VRAM)
	rl.UnloadImage(image)                    // Unload retrieved image data from CPU memory (RAM)

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
