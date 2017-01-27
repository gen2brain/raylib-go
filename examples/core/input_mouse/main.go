package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	raylib.InitWindow(800, 450, "raylib [core] example - mouse input")
	raylib.SetTargetFPS(60)

	ballColor := raylib.DarkBlue

	for !raylib.WindowShouldClose() {
		ballPosition := raylib.GetMousePosition()

		if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
			ballColor = raylib.Maroon
		} else if raylib.IsMouseButtonPressed(raylib.MouseMiddleButton) {
			ballColor = raylib.Lime
		} else if raylib.IsMouseButtonPressed(raylib.MouseRightButton) {
			ballColor = raylib.DarkBlue
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)
		raylib.DrawCircleV(ballPosition, 40, ballColor)

		raylib.DrawText("move ball with mouse and click mouse button to change color", 10, 10, 20, raylib.DarkGray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
