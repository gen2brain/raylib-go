package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - raylib logo using shapes")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawRectangle(screenWidth/2-128, screenHeight/2-128, 256, 256, rl.Black)
		rl.DrawRectangle(screenWidth/2-112, screenHeight/2-112, 224, 224, rl.RayWhite)
		rl.DrawText("raylib", screenWidth/2-44, screenHeight/2+48, 50, rl.Black)

		rl.DrawText("this is NOT a texture!", 350, 370, 10, rl.Gray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
