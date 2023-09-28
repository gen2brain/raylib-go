package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	mouseScaleMarkSize              = float32(12)
	rec                             = rl.NewRectangle(100, 100, 200, 80)
	mouseScaleReady, MouseScaleMode = false, false
)

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - rectangle scaling mouse")

	rl.SetTargetFPS(60)

	rl.SetMousePosition(0, 0)

	for !rl.WindowShouldClose() {

		mousePos := rl.GetMousePosition()

		if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(rec.X+rec.Width-mouseScaleMarkSize, rec.Y+rec.Height-mouseScaleMarkSize, mouseScaleMarkSize, mouseScaleMarkSize)) {
			mouseScaleReady = true
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				MouseScaleMode = true
			}
		} else {
			mouseScaleReady = false
		}

		if MouseScaleMode {

			mouseScaleReady = true
			rec.Width = mousePos.X - rec.X
			rec.Height = mousePos.Y - rec.Y

			// CHECK MIN MAX REC SIZES
			if rec.Width < mouseScaleMarkSize {
				rec.Width = rec.Width
			}
			if rec.Height < mouseScaleMarkSize {
				rec.Height = rec.Width
			}
			if rec.Width > (float32(rl.GetScreenWidth()) - rec.X) {
				rec.Width = float32(rl.GetScreenWidth()) - rec.X
			}
			if rec.Height > (float32(rl.GetScreenHeight()) - rec.Y) {
				rec.Height = float32(rl.GetScreenHeight()) - rec.Y
			}
			if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				MouseScaleMode = false
			}

		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("Scale rectangle dragging from bottom-right corner!", 10, 10, 20, rl.Black)

		rl.DrawRectangleRec(rec, rl.Fade(rl.Green, 0.5))

		if mouseScaleReady {
			rl.DrawRectangleLinesEx(rec, 1, rl.Red)
			rl.DrawTriangle(rl.NewVector2(rec.X+rec.Width-mouseScaleMarkSize, rec.Y+rec.Height), rl.NewVector2(rec.X+rec.Width, rec.Y+rec.Height), rl.NewVector2(rec.X+rec.Width, rec.Y+rec.Height-mouseScaleMarkSize), rl.Red)
		}

		rl.EndDrawing()

	}

	rl.CloseWindow()
}
