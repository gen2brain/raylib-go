package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	rl.SetConfigFlags(rl.FlagMsaa4xHint)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - bouncing ball")

	ballPos := rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2)
	ballSpeed := rl.NewVector2(5, 4)
	ballRadius := 20

	pause := false
	frames := 0

	rl.SetTargetFPS(60)

	rl.SetMousePosition(0, 0)

	for !rl.WindowShouldClose() {

		if rl.IsKeyPressed(rl.KeySpace) {
			pause = !pause
		}

		if !pause {
			ballPos.X += ballSpeed.X
			ballPos.Y += ballSpeed.Y
			if ballPos.X >= float32(screenWidth)-float32(ballRadius) || ballPos.X <= float32(ballRadius) {
				ballSpeed.X *= -1
			}
			if ballPos.Y >= float32(screenHeight)-float32(ballRadius) || ballPos.Y <= float32(ballRadius) {
				ballSpeed.Y *= -1
			}
		} else {
			frames++
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("SPACE key to pause", 10, 10, 20, rl.Black)

		rl.DrawCircleV(ballPos, float32(ballRadius), rl.Red)

		if pause && (frames/30)%2 == 0 {
			rl.DrawText("PAUSED", 10, screenHeight-40, 30, rl.Black)
		}

		rl.DrawFPS(screenWidth-100, 10)

		rl.EndDrawing()

	}

	rl.CloseWindow()
}
