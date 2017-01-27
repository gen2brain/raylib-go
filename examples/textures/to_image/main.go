package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [textures] example - texture to image")

	image := raylib.LoadImage("raylib_logo.png")  // Load image data into CPU memory (RAM)
	texture := raylib.LoadTextureFromImage(image) // Image converted to texture, GPU memory (RAM -> VRAM)
	raylib.UnloadImage(image)                     // Unload image data from CPU memory (RAM)

	image = raylib.GetTextureData(texture) // Retrieve image data from GPU memory (VRAM -> RAM)
	raylib.UnloadTexture(texture)          // Unload texture from GPU memory (VRAM)

	texture = raylib.LoadTextureFromImage(image) // Recreate texture from retrieved image data (RAM -> VRAM)
	raylib.UnloadImage(image)                    // Unload retrieved image data from CPU memory (RAM)

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
