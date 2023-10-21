package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - texture drawing")

	imBlank := rl.GenImageColor(1024, 1024, rl.Blank)
	texture := rl.LoadTextureFromImage(imBlank)
	rl.UnloadImage(imBlank)

	shader := rl.LoadShader("", "cubes_panning.fs")

	time := []float32{0}

	timeLoc := rl.GetShaderLocation(shader, "uTime")

	rl.SetShaderValue(shader, timeLoc, time, rl.ShaderUniformFloat)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		time = nil
		time = []float32{float32(rl.GetTime())}
		rl.SetShaderValue(shader, timeLoc, time, rl.ShaderUniformFloat)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginShaderMode(shader)
		rl.DrawTexture(texture, 0, 0, rl.White)
		rl.EndShaderMode()

		rl.DrawText("BACKGROUND is PAINTED and ANIMATED on SHADER!", 10, 10, 20, rl.Maroon)

		rl.EndDrawing()
	}

	rl.UnloadShader(shader) // Unload shader

	rl.CloseWindow()
}
