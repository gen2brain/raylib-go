package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - texture loading and drawing")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	texture := rl.LoadTexture("raylib_logo.png")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2, rl.White)
		rl.DrawText("this IS a texture!", 360, 370, 10, rl.Gray)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)

	rl.CloseWindow()
}
