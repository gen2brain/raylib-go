package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	StorageScore   = 0
	StorageHiscore = 1
)

func main() {
	raylib.InitWindow(800, 450, "raylib [core] example - storage save/load values")

	score := int32(0)
	hiscore := int32(0)

	framesCounter := 0

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if raylib.IsKeyPressed(raylib.KeyR) {
			score = raylib.GetRandomValue(1000, 2000)
			hiscore = raylib.GetRandomValue(2000, 4000)
		}

		if raylib.IsKeyPressed(raylib.KeyEnter) {
			raylib.StorageSaveValue(StorageScore, score)
			raylib.StorageSaveValue(StorageHiscore, hiscore)
		} else if raylib.IsKeyPressed(raylib.KeySpace) {
			// NOTE: If requested position could not be found, value 0 is returned
			score = raylib.StorageLoadValue(StorageScore)
			hiscore = raylib.StorageLoadValue(StorageHiscore)
		}

		framesCounter++

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText(fmt.Sprintf("SCORE: %d", score), 280, 130, 40, raylib.Maroon)
		raylib.DrawText(fmt.Sprintf("HI-SCORE: %d", hiscore), 210, 200, 50, raylib.Black)

		raylib.DrawText(fmt.Sprintf("frames: %d", framesCounter), 10, 10, 20, raylib.Lime)

		raylib.DrawText("Press R to generate random numbers", 220, 40, 20, raylib.LightGray)
		raylib.DrawText("Press ENTER to SAVE values", 250, 310, 20, raylib.LightGray)
		raylib.DrawText("Press SPACE to LOAD values", 252, 350, 20, raylib.LightGray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
