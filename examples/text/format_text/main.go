package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [text] example - text formatting")

	score := 100020
	hiscore := 200450
	lives := 5

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText(fmt.Sprintf("Score: %08d", score), 200, 80, 20, raylib.Red)

		raylib.DrawText(fmt.Sprintf("HiScore: %08d", hiscore), 200, 120, 20, raylib.Green)

		raylib.DrawText(fmt.Sprintf("Lives: %02d", lives), 200, 160, 40, raylib.Blue)

		raylib.DrawText(fmt.Sprintf("Elapsed Time: %02.02f ms", raylib.GetFrameTime()*1000), 200, 220, 20, raylib.Black)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
