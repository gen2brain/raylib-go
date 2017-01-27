package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - raylib logo using shapes")

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawRectangle(screenWidth/2-128, screenHeight/2-128, 256, 256, raylib.Black)
		raylib.DrawRectangle(screenWidth/2-112, screenHeight/2-112, 224, 224, raylib.RayWhite)
		raylib.DrawText("raylib", screenWidth/2-44, screenHeight/2+48, 50, raylib.Black)

		raylib.DrawText("this is NOT a texture!", 350, 370, 10, raylib.Gray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
