package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - mouse input")
	rl.SetTargetFPS(60)

	ballColor := rl.DarkBlue

	for !rl.WindowShouldClose() {
		ballPosition := rl.GetMousePosition()

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			ballColor = rl.Maroon
		} else if rl.IsMouseButtonPressed(rl.MouseMiddleButton) {
			ballColor = rl.Lime
		} else if rl.IsMouseButtonPressed(rl.MouseRightButton) {
			ballColor = rl.DarkBlue
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawCircleV(ballPosition, 40, ballColor)

		rl.DrawText("move ball with mouse and click mouse button to change color", 10, 10, 20, rl.DarkGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
