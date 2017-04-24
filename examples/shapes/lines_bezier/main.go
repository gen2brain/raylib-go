package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - cubic-bezier lines")

	start := raylib.NewVector2(0, 0)
	end := raylib.NewVector2(float32(screenWidth), float32(screenHeight))

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			start = raylib.GetMousePosition()
		} else if raylib.IsMouseButtonDown(raylib.MouseRightButton) {
			end = raylib.GetMousePosition()
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("USE MOUSE LEFT-RIGHT CLICK to DEFINE LINE START and END POINTS", 15, 20, 20, raylib.Gray)

		raylib.DrawLineBezier(start, end, 2.0, raylib.Red)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
