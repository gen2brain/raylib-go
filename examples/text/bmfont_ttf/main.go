package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - bmfont and ttf sprite fonts loading")

	msgBm := "THIS IS AN AngelCode SPRITE FONT"
	msgTtf := "THIS SPRITE FONT has been GENERATED from a TTF"

	// NOTE: Textures/Fonts MUST be loaded after Window initialization (OpenGL context is required)
	fontBm := rl.LoadFont("fonts/bmfont.fnt")      // BMFont (AngelCode)
	fontTtf := rl.LoadFont("fonts/pixantiqua.ttf") // TTF font

	fontPosition := rl.Vector2{}

	fontPosition.X = float32(screenWidth)/2 - rl.MeasureTextEx(fontBm, msgBm, float32(fontBm.BaseSize), 0).X/2
	fontPosition.Y = float32(screenHeight)/2 - float32(fontBm.BaseSize)/2 - 80

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawTextEx(fontBm, msgBm, fontPosition, float32(fontBm.BaseSize), 0, rl.Maroon)
		rl.DrawTextEx(fontTtf, msgTtf, rl.NewVector2(75.0, 240.0), float32(fontTtf.BaseSize)*0.8, 2, rl.Lime)

		rl.EndDrawing()
	}

	rl.UnloadFont(fontBm)  // AngelCode Font unloading
	rl.UnloadFont(fontTtf) // TTF Font unloading

	rl.CloseWindow()
}
