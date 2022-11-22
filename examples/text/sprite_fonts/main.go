package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - sprite fonts usage")

	msg1 := "THIS IS A custom SPRITE FONT..."
	msg2 := "...and this is ANOTHER CUSTOM font..."
	msg3 := "...and a THIRD one! GREAT! :D"

	// NOTE: Textures/Fonts MUST be loaded after Window initialization (OpenGL context is required)
	font1 := rl.LoadFont("fonts/custom_mecha.png")         // Font loading
	font2 := rl.LoadFont("fonts/custom_alagard.png")       // Font loading
	font3 := rl.LoadFont("fonts/custom_jupiter_crash.png") // Font loading

	var fontPosition1, fontPosition2, fontPosition3 rl.Vector2

	fontPosition1.X = float32(screenWidth)/2 - rl.MeasureTextEx(font1, msg1, float32(font1.BaseSize), -3).X/2
	fontPosition1.Y = float32(screenHeight)/2 - float32(font1.BaseSize)/2 - 80

	fontPosition2.X = float32(screenWidth)/2 - rl.MeasureTextEx(font2, msg2, float32(font2.BaseSize), -2).X/2
	fontPosition2.Y = float32(screenHeight)/2 - float32(font2.BaseSize)/2 - 10

	fontPosition3.X = float32(screenWidth)/2 - rl.MeasureTextEx(font3, msg3, float32(font3.BaseSize), 2).X/2
	fontPosition3.Y = float32(screenHeight)/2 - float32(font3.BaseSize)/2 + 50

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawTextEx(font1, msg1, fontPosition1, float32(font1.BaseSize), -3, rl.White)
		rl.DrawTextEx(font2, msg2, fontPosition2, float32(font2.BaseSize), -2, rl.White)
		rl.DrawTextEx(font3, msg3, fontPosition3, float32(font3.BaseSize), 2, rl.White)

		rl.EndDrawing()
	}

	rl.UnloadFont(font1) // Font unloading
	rl.UnloadFont(font2) // Font unloading
	rl.UnloadFont(font3) // Font unloading

	rl.CloseWindow()
}
