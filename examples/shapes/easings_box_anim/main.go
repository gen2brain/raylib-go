package main

import (
	ez "github.com/gen2brain/raylib-go/easings"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	fps = int32(60)
)

const (
	screenWidth  = int32(800)
	screenHeight = int32(450)
)

func main() {

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - easings box anim")

	rec := rl.NewRectangle(float32(screenWidth/2), -100, 100, 100)
	rotation := float32(0)
	alpha := float32(1)

	state := 0
	frames := 0

	rl.SetTargetFPS(fps)

	for !rl.WindowShouldClose() {

		switch state {
		case 0:
			frames++

			rec.Y = ez.ElasticOut(float32(frames), -100, float32(screenHeight/2)+100, 120)

			if frames >= 120 {
				frames = 0
				state = 1
			}

		case 1:
			frames++
			rec.Height = ez.BounceOut(float32(frames), 100, -90, 120)
			rec.Width = ez.BounceOut(float32(frames), 100, float32(screenWidth), 120)

			if frames >= 120 {
				frames++
				state = 2
			}
		case 2:
			frames++
			rotation = ez.QuadOut(float32(frames), 0, 270, 240)
			if frames >= 240 {
				frames = 0
				state = 3
			}
		case 3:
			frames++
			rec.Height = ez.CircOut(float32(frames), 10, float32(screenWidth), 120)
			if frames >= 120 {
				frames = 0
				state = 4
			}

		case 4:
			frames++
			alpha = ez.SineOut(float32(frames), 1, -1, 160)
			if frames >= 160 {
				frames = 0
				state = 5
			}
		default:
			break
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			rec = rl.NewRectangle(float32(screenWidth/2), -100, 100, 100)
			rotation = 0
			alpha = 1
			state = 0
			frames = 0
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawRectanglePro(rec, rl.NewVector2(rec.Width/2, rec.Height/2), rotation, rl.Fade(rl.Black, alpha))

		if state == 5 {
			txtlen := rl.MeasureText("SPACE to replay", 20)
			rl.DrawText("SPACE to replay", (screenWidth/2)-txtlen/2, 200, 20, rl.Black)
		}
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
