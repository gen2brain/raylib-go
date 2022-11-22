package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - shapes and texture shaders")

	fudesumi := rl.LoadTexture("fudesumi.png")

	// NOTE: Using GLSL 330 shader version, on OpenGL ES 2.0 use GLSL 100 shader version
	shader := rl.LoadShader("glsl330/base.vs", "glsl330/grayscale.fs")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		// Start drawing with default shader

		rl.DrawText("USING DEFAULT SHADER", 20, 40, 10, rl.Red)

		rl.DrawCircle(80, 120, 35, rl.DarkBlue)
		rl.DrawCircleGradient(80, 220, 60, rl.Green, rl.SkyBlue)
		rl.DrawCircleLines(80, 340, 80, rl.DarkBlue)

		// Activate our custom shader to be applied on next shapes/textures drawings
		rl.BeginShaderMode(shader)

		rl.DrawText("USING CUSTOM SHADER", 190, 40, 10, rl.Red)

		rl.DrawRectangle(250-60, 90, 120, 60, rl.Red)
		rl.DrawRectangleGradientH(250-90, 170, 180, 130, rl.Maroon, rl.Gold)
		rl.DrawRectangleLines(250-40, 320, 80, 60, rl.Orange)

		// Activate our default shader for next drawings
		rl.EndShaderMode()

		rl.DrawText("USING DEFAULT SHADER", 370, 40, 10, rl.Red)

		rl.DrawTriangle(rl.NewVector2(430, 80),
			rl.NewVector2(430-60, 150),
			rl.NewVector2(430+60, 150), rl.Violet)

		rl.DrawTriangleLines(rl.NewVector2(430, 160),
			rl.NewVector2(430-20, 230),
			rl.NewVector2(430+20, 230), rl.DarkBlue)

		rl.DrawPoly(rl.NewVector2(430, 320), 6, 80, 0, rl.Brown)

		// Activate our custom shader to be applied on next shapes/textures drawings
		rl.BeginShaderMode(shader)

		rl.DrawTexture(fudesumi, 500, -30, rl.White) // Using custom shader

		// Activate our default shader for next drawings
		rl.EndShaderMode()

		rl.DrawText("(c) Fudesumi sprite by Eiden Marsal", 380, screenHeight-20, 10, rl.Gray)

		rl.EndDrawing()
	}

	rl.UnloadShader(shader)    // Unload shader
	rl.UnloadTexture(fudesumi) // Unload texture

	rl.CloseWindow() // Close window and OpenGL context
}
