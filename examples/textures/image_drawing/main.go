package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - texture loading and drawing")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	cat := rl.LoadImage("cat.png")                        // Load image in CPU memory (RAM)
	rl.ImageCrop(cat, rl.NewRectangle(100, 10, 280, 380)) // Crop an image piece
	rl.ImageFlipHorizontal(cat)                           // Flip cropped image horizontally
	rl.ImageResize(cat, 150, 200)                         // Resize flipped-cropped image

	parrots := rl.LoadImage("parrots.png") // Load image in CPU memory (RAM)

	// Draw one image over the other with a scaling of 1.5f
	rl.ImageDraw(parrots, cat, rl.NewRectangle(0, 0, float32(cat.Width), float32(cat.Height)), rl.NewRectangle(30, 40, float32(cat.Width)*1.5, float32(cat.Height)*1.5), rl.White)
	rl.ImageCrop(parrots, rl.NewRectangle(0, 50, float32(parrots.Width), float32(parrots.Height-100))) // Crop resulting image

	rl.UnloadImage(cat) // Unload image from RAM

	texture := rl.LoadTextureFromImage(parrots) // Image converted to texture, uploaded to GPU memory (VRAM)
	rl.UnloadImage(parrots)                     // Once image has been converted to texture and uploaded to VRAM, it can be unloaded from RAM

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2-40, rl.White)
		rl.DrawRectangleLines(screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2-40, texture.Width, texture.Height, rl.DarkGray)

		rl.DrawText("We are drawing only one texture from various images composed!", 240, 350, 10, rl.DarkGray)
		rl.DrawText("Source images have been cropped, scaled, flipped and copied one over the other.", 190, 370, 10, rl.DarkGray)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)

	rl.CloseWindow()
}
