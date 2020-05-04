package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	storageScore   = 0
	storageHiscore = 1
)

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - storage save/load values")

	score := int32(0)
	hiscore := int32(0)

	framesCounter := 0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyR) {
			score = rl.GetRandomValue(1000, 2000)
			hiscore = rl.GetRandomValue(2000, 4000)
		}

		if rl.IsKeyPressed(rl.KeyEnter) {
			rl.SaveStorageValue(storageScore, score)
			rl.SaveStorageValue(storageHiscore, hiscore)
		} else if rl.IsKeyPressed(rl.KeySpace) {
			// NOTE: If requested position could not be found, value 0 is returned
			score = rl.LoadStorageValue(storageScore)
			hiscore = rl.LoadStorageValue(storageHiscore)
		}

		framesCounter++

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText(fmt.Sprintf("SCORE: %d", score), 280, 130, 40, rl.Maroon)
		rl.DrawText(fmt.Sprintf("HI-SCORE: %d", hiscore), 210, 200, 50, rl.Black)

		rl.DrawText(fmt.Sprintf("frames: %d", framesCounter), 10, 10, 20, rl.Lime)

		rl.DrawText("Press R to generate random numbers", 220, 40, 20, rl.LightGray)
		rl.DrawText("Press ENTER to SAVE values", 250, 310, 20, rl.LightGray)
		rl.DrawText("Press SPACE to LOAD values", 252, 350, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
