package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [textures] example - image text drawing")

	// TTF Font loading with custom generation parameters
	var fontChars int32
	font := raylib.LoadFontEx("fonts/KAISG.ttf", 64, 0, &fontChars)

	parrots := raylib.LoadImage("parrots.png") // Load image in CPU memory (RAM)

	// Draw over image using custom font
	raylib.ImageDrawTextEx(parrots, raylib.NewVector2(20, 20), font, "[Parrots font drawing]", float32(font.BaseSize), 0, raylib.White)

	texture := raylib.LoadTextureFromImage(parrots) // Image converted to texture, uploaded to GPU memory (VRAM)

	raylib.UnloadImage(parrots) // Once image has been converted to texture and uploaded to VRAM, it can be unloaded from RAM

	position := raylib.NewVector2(float32(screenWidth)/2-float32(texture.Width)/2, float32(screenHeight)/2-float32(texture.Height)/2-20)

	showFont := false

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if raylib.IsKeyDown(raylib.KeySpace) {
			showFont = true
		} else {
			showFont = false
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		if !showFont {
			// Draw texture with text already drawn inside
			raylib.DrawTextureV(texture, position, raylib.White)

			// Draw text directly using sprite font
			raylib.DrawTextEx(font, "[Parrots font drawing]", raylib.NewVector2(position.X+20, position.Y+20+280), float32(font.BaseSize), 0, raylib.White)
		} else {
			raylib.DrawTexture(font.Texture, screenWidth/2-font.Texture.Width/2, 50, raylib.Black)
		}

		raylib.DrawText("PRESS SPACE to SEE USED SPRITEFONT ", 290, 420, 10, raylib.DarkGray)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(texture)
	raylib.UnloadFont(font)

	raylib.CloseWindow()
}
