package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [text] example - sprite fonts usage")

	msg1 := "THIS IS A custom SPRITE FONT..."
	msg2 := "...and this is ANOTHER CUSTOM font..."
	msg3 := "...and a THIRD one! GREAT! :D"

	// NOTE: Textures/Fonts MUST be loaded after Window initialization (OpenGL context is required)
	font1 := raylib.LoadFont("fonts/custom_mecha.png")         // Font loading
	font2 := raylib.LoadFont("fonts/custom_alagard.png")       // Font loading
	font3 := raylib.LoadFont("fonts/custom_jupiter_crash.png") // Font loading

	var fontPosition1, fontPosition2, fontPosition3 raylib.Vector2

	fontPosition1.X = float32(screenWidth)/2 - raylib.MeasureTextEx(font1, msg1, float32(font1.BaseSize), -3).X/2
	fontPosition1.Y = float32(screenHeight)/2 - float32(font1.BaseSize)/2 - 80

	fontPosition2.X = float32(screenWidth)/2 - raylib.MeasureTextEx(font2, msg2, float32(font2.BaseSize), -2).X/2
	fontPosition2.Y = float32(screenHeight)/2 - float32(font2.BaseSize)/2 - 10

	fontPosition3.X = float32(screenWidth)/2 - raylib.MeasureTextEx(font3, msg3, float32(font3.BaseSize), 2).X/2
	fontPosition3.Y = float32(screenHeight)/2 - float32(font3.BaseSize)/2 + 50

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawTextEx(font1, msg1, fontPosition1, float32(font1.BaseSize), -3, raylib.White)
		raylib.DrawTextEx(font2, msg2, fontPosition2, float32(font2.BaseSize), -2, raylib.White)
		raylib.DrawTextEx(font3, msg3, fontPosition3, float32(font3.BaseSize), 2, raylib.White)

		raylib.EndDrawing()
	}

	raylib.UnloadFont(font1) // Font unloading
	raylib.UnloadFont(font2) // Font unloading
	raylib.UnloadFont(font3) // Font unloading

	raylib.CloseWindow()
}
