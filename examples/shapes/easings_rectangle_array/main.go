package main

import (
	ez "github.com/gen2brain/raylib-go/easings"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	playTime = float32(240)
	fps      = int32(60)
)

const (
	recsW, recsH = 50, 50
	screenWidth  = int32(800)
	screenHeight = int32(450)
	maxRecsX     = int(screenWidth) / recsW
	maxRecsY     = int(screenHeight) / recsH
)

func main() {

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - easings rectangle array")

	var recs [maxRecsX * maxRecsY]rl.Rectangle

	for y := 0; y < maxRecsY; y++ {

		for x := 0; x < maxRecsX; x++ {
			recs[y*maxRecsX+x].X = float32(recsW/2 + recsW*float32(x))
			recs[y*maxRecsX+x].Y = float32(recsH/2 + recsH*float32(y))
			recs[y*maxRecsX+x].Width = float32(recsW)
			recs[y*maxRecsX+x].Height = float32(recsH)
		}

	}

	rotation := float32(0)
	frameCount := float32(0)
	state := 0

	rl.SetTargetFPS(fps)

	for !rl.WindowShouldClose() {

		if state == 0 {
			frameCount++

			for i := 0; i < maxRecsX*maxRecsY; i++ {
				recs[i].Height = float32(ez.LinearIn(frameCount, recsH, -recsH, playTime))
				recs[i].Width = float32(ez.LinearIn(frameCount, recsW, -recsW, playTime))

				if recs[i].Height < 0 {
					recs[i].Height = 0
				}

				if recs[i].Width < 0 {
					recs[i].Width = 0
				}

				if recs[i].Height == 0 && recs[i].Width == 0 {
					state = 1
				}
				rotation = float32(ez.LinearIn(frameCount, 0, 360, playTime))
			}
		} else if state == 1 && rl.IsKeyPressed(rl.KeySpace) {
			frameCount = 0
			for i := 0; i < maxRecsX*maxRecsY; i++ {
				recs[i].Height = float32(recsH)
				recs[i].Width = float32(recsW)
			}

			state = 0
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		if state == 0 {
			for i := 0; i < maxRecsX*maxRecsY; i++ {
				rl.DrawRectanglePro(recs[i], rl.NewVector2(recs[i].Width/2, recs[i].Height/2), rotation, rl.Red)
			}
		} else if state == 1 {
			txtlen := rl.MeasureText("SPACE to replay", 20)
			rl.DrawText("SPACE to replay", (screenWidth/2)-txtlen/2, 200, 20, rl.Black)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
