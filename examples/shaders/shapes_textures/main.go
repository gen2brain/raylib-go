package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - shapes and texture shaders")

	fudesumi := raylib.LoadTexture("fudesumi.png")

	// NOTE: Using GLSL 330 shader version, on OpenGL ES 2.0 use GLSL 100 shader version
	shader := raylib.LoadShader("glsl330/base.vs", "glsl330/grayscale.fs")

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		// Start drawing with default shader

		raylib.DrawText("USING DEFAULT SHADER", 20, 40, 10, raylib.Red)

		raylib.DrawCircle(80, 120, 35, raylib.DarkBlue)
		raylib.DrawCircleGradient(80, 220, 60, raylib.Green, raylib.SkyBlue)
		raylib.DrawCircleLines(80, 340, 80, raylib.DarkBlue)

		// Activate our custom shader to be applied on next shapes/textures drawings
		raylib.BeginShaderMode(shader)

		raylib.DrawText("USING CUSTOM SHADER", 190, 40, 10, raylib.Red)

		raylib.DrawRectangle(250-60, 90, 120, 60, raylib.Red)
		raylib.DrawRectangleGradientH(250-90, 170, 180, 130, raylib.Maroon, raylib.Gold)
		raylib.DrawRectangleLines(250-40, 320, 80, 60, raylib.Orange)

		// Activate our default shader for next drawings
		raylib.EndShaderMode()

		raylib.DrawText("USING DEFAULT SHADER", 370, 40, 10, raylib.Red)

		raylib.DrawTriangle(raylib.NewVector2(430, 80),
			raylib.NewVector2(430-60, 150),
			raylib.NewVector2(430+60, 150), raylib.Violet)

		raylib.DrawTriangleLines(raylib.NewVector2(430, 160),
			raylib.NewVector2(430-20, 230),
			raylib.NewVector2(430+20, 230), raylib.DarkBlue)

		raylib.DrawPoly(raylib.NewVector2(430, 320), 6, 80, 0, raylib.Brown)

		// Activate our custom shader to be applied on next shapes/textures drawings
		raylib.BeginShaderMode(shader)

		raylib.DrawTexture(fudesumi, 500, -30, raylib.White) // Using custom shader

		// Activate our default shader for next drawings
		raylib.EndShaderMode()

		raylib.DrawText("(c) Fudesumi sprite by Eiden Marsal", 380, screenHeight-20, 10, raylib.Gray)

		raylib.EndDrawing()
	}

	raylib.UnloadShader(shader)    // Unload shader
	raylib.UnloadTexture(fudesumi) // Unload texture

	raylib.CloseWindow() // Close window and OpenGL context
}
