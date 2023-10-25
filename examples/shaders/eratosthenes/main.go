package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - Eratosthenes")

	shader := rl.LoadShader("", "eratosthenes.fs")
	target := rl.LoadRenderTexture(screenWidth, screenHeight)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rl.BeginTextureMode(target)
		rl.ClearBackground(rl.Black)
		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Black)
		rl.EndTextureMode()

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginShaderMode(shader)

		rl.DrawTextureRec(target.Texture, rl.NewRectangle(0, 0, float32(target.Texture.Width), float32(target.Texture.Height)), rl.Vector2Zero(), rl.White)

		rl.EndShaderMode()

		rl.EndDrawing()
	}

	rl.UnloadShader(shader)
	rl.UnloadRenderTexture(target)

	rl.CloseWindow()
}
