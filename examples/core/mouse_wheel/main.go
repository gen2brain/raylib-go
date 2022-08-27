package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - mouse wheel")

	boxPositionY := screenHeight/2 - 40
	scrollSpeed := int32(4) // Scrolling speed in pixels

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		boxPositionY -= int32(rl.GetMouseWheelMove()) * scrollSpeed

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawRectangle(screenWidth/2-40, boxPositionY, 80, 80, rl.Maroon)

		rl.DrawText("Use mouse wheel to move the square up and down!", 10, 10, 20, rl.Gray)
		rl.DrawText(fmt.Sprintf("Box position Y: %d", boxPositionY), 10, 40, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
