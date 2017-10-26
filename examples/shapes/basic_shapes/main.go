package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - basic shapes drawing")

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("some basic shapes available on raylib", 20, 20, 20, raylib.DarkGray)

		raylib.DrawLine(18, 42, screenWidth-18, 42, raylib.Black)

		raylib.DrawCircle(screenWidth/4, 120, 35, raylib.DarkBlue)
		raylib.DrawCircleGradient(screenWidth/4, 220, 60, raylib.Green, raylib.SkyBlue)
		raylib.DrawCircleLines(screenWidth/4, 340, 80, raylib.DarkBlue)

		raylib.DrawRectangle(screenWidth/4*2-60, 100, 120, 60, raylib.Red)
		raylib.DrawRectangleGradientH(screenWidth/4*2-90, 170, 180, 130, raylib.Maroon, raylib.Gold)
		raylib.DrawRectangleLines(screenWidth/4*2-40, 320, 80, 60, raylib.Orange)

		raylib.DrawTriangle(raylib.NewVector2(float32(screenWidth)/4*3, 80),
			raylib.NewVector2(float32(screenWidth)/4*3-60, 150),
			raylib.NewVector2(float32(screenWidth)/4*3+60, 150), raylib.Violet)

		raylib.DrawTriangleLines(raylib.NewVector2(float32(screenWidth)/4*3, 160),
			raylib.NewVector2(float32(screenWidth)/4*3-20, 230),
			raylib.NewVector2(float32(screenWidth)/4*3+20, 230), raylib.DarkBlue)

		raylib.DrawPoly(raylib.NewVector2(float32(screenWidth)/4*3, 320), 6, 80, 0, raylib.Brown)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
