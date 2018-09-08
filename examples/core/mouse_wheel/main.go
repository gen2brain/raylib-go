package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [core] example - mouse wheel")

	boxPositionY := screenHeight/2 - 40
	scrollSpeed := int32(4) // Scrolling speed in pixels

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		boxPositionY -= raylib.GetMouseWheelMove() * scrollSpeed

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawRectangle(screenWidth/2-40, boxPositionY, 80, 80, raylib.Maroon)

		raylib.DrawText("Use mouse wheel to move the square up and down!", 10, 10, 20, raylib.Gray)
		raylib.DrawText(fmt.Sprintf("Box position Y: %d", boxPositionY), 10, 40, 20, raylib.LightGray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
