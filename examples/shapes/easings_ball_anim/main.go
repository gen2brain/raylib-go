package main

import (
	ez "github.com/gen2brain/raylib-go/easings"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagMsaa4xHint)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - easings ball anim")

	ballPosX := -100
	ballRadius := 20
	ballAlpha := float32(0)

	state := 0
	frames := 0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		if state == 0 {

			frames++
			ballPosX = int(ez.ElasticOut(float32(frames), -100, float32(screenWidth/2)+100, 100))

			if frames >= 100 {
				frames = 0
				state = 1
			}

		} else if state == 1 {
			frames++
			ballRadius = int(ez.ElasticIn(float32(frames), 20, 500, 150))

			if frames >= 150 {
				frames = 0
				state = 2
			}
		} else if state == 2 {
			frames++
			ballAlpha = ez.CubicOut(float32(frames), 0, 1, 150)

			if frames >= 150 {
				frames = 0
				state = 3
			}
		} else if state == 3 {
			if rl.IsKeyPressed(rl.KeyEnter) {
				ballPosX = -100
				ballRadius = 20
				ballAlpha = 0
				state = 0
			}

		}

		if rl.IsKeyPressed(rl.KeyR) {
			frames = 0
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		if state >= 2 {
			rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.Green)
		}

		rl.DrawCircle(int32(ballPosX), 200, float32(ballRadius), rl.Fade(rl.Red, 1-ballAlpha))

		if state == 3 {
			textlen := rl.MeasureText("press ENTER to replay", 20)
			rl.DrawText("press ENTER to replay", (screenWidth/2)-textlen/2, 200, 20, rl.Black)
		}

		rl.EndDrawing()

	}

	rl.CloseWindow()
}
