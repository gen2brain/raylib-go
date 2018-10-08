package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - basic shapes drawing")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("some basic shapes available on raylib", 20, 20, 20, rl.DarkGray)

		rl.DrawLine(18, 42, screenWidth-18, 42, rl.Black)

		rl.DrawCircle(screenWidth/4, 120, 35, rl.DarkBlue)
		rl.DrawCircleGradient(screenWidth/4, 220, 60, rl.Green, rl.SkyBlue)
		rl.DrawCircleLines(screenWidth/4, 340, 80, rl.DarkBlue)

		rl.DrawRectangle(screenWidth/4*2-60, 100, 120, 60, rl.Red)
		rl.DrawRectangleGradientH(screenWidth/4*2-90, 170, 180, 130, rl.Maroon, rl.Gold)
		rl.DrawRectangleLines(screenWidth/4*2-40, 320, 80, 60, rl.Orange)

		rl.DrawTriangle(rl.NewVector2(float32(screenWidth)/4*3, 80),
			rl.NewVector2(float32(screenWidth)/4*3-60, 150),
			rl.NewVector2(float32(screenWidth)/4*3+60, 150), rl.Violet)

		rl.DrawTriangleLines(rl.NewVector2(float32(screenWidth)/4*3, 160),
			rl.NewVector2(float32(screenWidth)/4*3-20, 230),
			rl.NewVector2(float32(screenWidth)/4*3+20, 230), rl.DarkBlue)

		rl.DrawPoly(rl.NewVector2(float32(screenWidth)/4*3, 320), 6, 80, 0, rl.Brown)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
