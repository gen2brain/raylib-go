package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - drop files")

	rl.SetTargetFPS(60)

	var count int
	var droppedFiles []string

	for !rl.WindowShouldClose() {
		if rl.IsFileDropped() {
			droppedFiles = rl.LoadDroppedFiles()
			count = len(droppedFiles)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		if count == 0 {
			rl.DrawText("Drop your files to this window!", 100, 40, 20, rl.DarkGray)
		} else {
			rl.DrawText("Dropped files:", 100, 40, 20, rl.DarkGray)

			for i := 0; i < count; i++ {
				if i%2 == 0 {
					rl.DrawRectangle(0, int32(85+40*i), screenWidth, 40, rl.Fade(rl.LightGray, 0.5))
				} else {
					rl.DrawRectangle(0, int32(85+40*i), screenWidth, 40, rl.Fade(rl.LightGray, 0.3))
				}

				rl.DrawText(droppedFiles[i], 120, int32(100+i*40), 10, rl.Gray)
			}

			rl.DrawText("Drop new files...", 100, int32(150+count*40), 20, rl.DarkGray)
		}

		rl.EndDrawing()
	}

	rl.UnloadDroppedFiles()

	rl.CloseWindow()
}
