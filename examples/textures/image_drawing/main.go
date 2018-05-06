package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [textures] example - texture loading and drawing")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	cat := raylib.LoadImage("cat.png")                            // Load image in CPU memory (RAM)
	raylib.ImageCrop(cat, raylib.NewRectangle(100, 10, 280, 380)) // Crop an image piece
	raylib.ImageFlipHorizontal(cat)                               // Flip cropped image horizontally
	raylib.ImageResize(cat, 150, 200)                             // Resize flipped-cropped image

	parrots := raylib.LoadImage("parrots.png") // Load image in CPU memory (RAM)

	// Draw one image over the other with a scaling of 1.5f
	raylib.ImageDraw(parrots, cat, raylib.NewRectangle(0, 0, float32(cat.Width), float32(cat.Height)), raylib.NewRectangle(30, 40, float32(cat.Width)*1.5, float32(cat.Height)*1.5))
	raylib.ImageCrop(parrots, raylib.NewRectangle(0, 50, float32(parrots.Width), float32(parrots.Height-100))) // Crop resulting image

	raylib.UnloadImage(cat) // Unload image from RAM

	texture := raylib.LoadTextureFromImage(parrots) // Image converted to texture, uploaded to GPU memory (VRAM)
	raylib.UnloadImage(parrots)                     // Once image has been converted to texture and uploaded to VRAM, it can be unloaded from RAM

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2-40, raylib.White)
		raylib.DrawRectangleLines(screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2-40, texture.Width, texture.Height, raylib.DarkGray)

		raylib.DrawText("We are drawing only one texture from various images composed!", 240, 350, 10, raylib.DarkGray)
		raylib.DrawText("Source images have been cropped, scaled, flipped and copied one over the other.", 190, 370, 10, raylib.DarkGray)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(texture)

	raylib.CloseWindow()
}
