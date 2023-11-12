package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	screenW = int32(800)
	screenH = int32(450)
)

func main() {

	rl.InitWindow(screenW, screenH, "raylib [core] example - scissor test")

	scissorArea := rl.NewRectangle(0, 0, 300, 300)
	scissorMode := true

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		if rl.IsKeyPressed(rl.KeyS) {
			scissorMode = !scissorMode
		}

		scissorArea.X = float32(rl.GetMouseX())
		scissorArea.Y = float32(rl.GetMouseY())

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		if scissorMode {
			rl.BeginScissorMode(scissorArea.ToInt32().X, scissorArea.ToInt32().Y, scissorArea.ToInt32().Width, scissorArea.ToInt32().Height)
		}

		rl.DrawRectangle(0, 0, screenW, screenH, rl.Red)
		rl.DrawText("MOVE MOUSE TO REVEAL TEXT", 190, 200, 20, rl.Black)

		if scissorMode {
			rl.EndScissorMode()
		}

		rl.DrawRectangleLinesEx(scissorArea, 1, rl.Black)
		rl.DrawText("S KEY TO TOGGLE MODE", 10, 10, 20, rl.Black)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
