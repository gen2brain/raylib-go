package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [text] example - bmfont and ttf sprite fonts loading")

	msgBm := "THIS IS AN AngelCode SPRITE FONT"
	msgTtf := "THIS SPRITE FONT has been GENERATED from a TTF"

	// NOTE: Textures/Fonts MUST be loaded after Window initialization (OpenGL context is required)
	fontBm := raylib.LoadFont("fonts/bmfont.fnt")      // BMFont (AngelCode)
	fontTtf := raylib.LoadFont("fonts/pixantiqua.ttf") // TTF font

	fontPosition := raylib.Vector2{}

	fontPosition.X = float32(screenWidth)/2 - raylib.MeasureTextEx(fontBm, msgBm, float32(fontBm.BaseSize), 0).X/2
	fontPosition.Y = float32(screenHeight)/2 - float32(fontBm.BaseSize)/2 - 80

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawTextEx(fontBm, msgBm, fontPosition, float32(fontBm.BaseSize), 0, raylib.Maroon)
		raylib.DrawTextEx(fontTtf, msgTtf, raylib.NewVector2(75.0, 240.0), float32(fontTtf.BaseSize)*0.8, 2, raylib.Lime)

		raylib.EndDrawing()
	}

	raylib.UnloadFont(fontBm)  // AngelCode Font unloading
	raylib.UnloadFont(fontTtf) // TTF Font unloading

	raylib.CloseWindow()
}
