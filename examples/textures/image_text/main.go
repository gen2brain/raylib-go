package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - image text drawing")

	// TTF Font loading with custom generation parameters
	font := rl.LoadFontEx("fonts/KAISG.ttf", 64, nil)

	parrots := rl.LoadImage("parrots.png") // Load image in CPU memory (RAM)

	// Draw over image using custom font
	rl.ImageDrawTextEx(parrots, rl.NewVector2(20, 20), font, "[Parrots font drawing]", float32(font.BaseSize), 0, rl.White)

	texture := rl.LoadTextureFromImage(parrots) // Image converted to texture, uploaded to GPU memory (VRAM)

	rl.UnloadImage(parrots) // Once image has been converted to texture and uploaded to VRAM, it can be unloaded from RAM

	position := rl.NewVector2(float32(screenWidth)/2-float32(texture.Width)/2, float32(screenHeight)/2-float32(texture.Height)/2-20)

	showFont := false

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if rl.IsKeyDown(rl.KeySpace) {
			showFont = true
		} else {
			showFont = false
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		if !showFont {
			// Draw texture with text already drawn inside
			rl.DrawTextureV(texture, position, rl.White)

			// Draw text directly using sprite font
			rl.DrawTextEx(font, "[Parrots font drawing]", rl.NewVector2(position.X+20, position.Y+20+280), float32(font.BaseSize), 0, rl.White)
		} else {
			rl.DrawTexture(font.Texture, screenWidth/2-font.Texture.Width/2, 50, rl.Black)
		}

		rl.DrawText("PRESS SPACE to SEE USED SPRITEFONT ", 290, 420, 10, rl.DarkGray)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)
	rl.UnloadFont(font)

	rl.CloseWindow()
}
