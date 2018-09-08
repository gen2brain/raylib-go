package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [core] example - drop files")

	raylib.SetTargetFPS(60)

	count := int32(0)
	var droppedFiles []string

	for !raylib.WindowShouldClose() {
		if raylib.IsFileDropped() {
			droppedFiles = raylib.GetDroppedFiles(&count)
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)
		if count == 0 {
			raylib.DrawText("Drop your files to this window!", 100, 40, 20, raylib.DarkGray)
		} else {
			raylib.DrawText("Dropped files:", 100, 40, 20, raylib.DarkGray)

			for i := int32(0); i < count; i++ {
				if i%2 == 0 {
					raylib.DrawRectangle(0, int32(85+40*i), screenWidth, 40, raylib.Fade(raylib.LightGray, 0.5))
				} else {
					raylib.DrawRectangle(0, int32(85+40*i), screenWidth, 40, raylib.Fade(raylib.LightGray, 0.3))
				}

				raylib.DrawText(droppedFiles[i], 120, int32(100+i*40), 10, raylib.Gray)
			}

			raylib.DrawText("Drop new files...", 100, int32(150+count*40), 20, raylib.DarkGray)
		}

		raylib.EndDrawing()
	}

	raylib.ClearDroppedFiles()

	raylib.CloseWindow()
}
