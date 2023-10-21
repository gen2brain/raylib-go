package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - texture waves")

	texture := rl.LoadTexture("space.png")

	shader := rl.LoadShader("", "wave.fs")

	secondsLoc := rl.GetShaderLocation(shader, "seconds")
	freqXLoc := rl.GetShaderLocation(shader, "freqX")
	freqYLoc := rl.GetShaderLocation(shader, "freqY")
	ampXLoc := rl.GetShaderLocation(shader, "ampX")
	ampYLoc := rl.GetShaderLocation(shader, "ampY")
	speedXLoc := rl.GetShaderLocation(shader, "speedX")
	speedYLoc := rl.GetShaderLocation(shader, "speedY")

	freqX := []float32{25}
	freqY := []float32{25}
	ampX := []float32{5}
	ampY := []float32{5}
	speedX := []float32{8}
	speedY := []float32{8}

	screensize := []float32{float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())}

	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "size"), screensize, rl.ShaderUniformVec2)
	rl.SetShaderValue(shader, freqXLoc, freqX, rl.ShaderUniformFloat)
	rl.SetShaderValue(shader, freqYLoc, freqY, rl.ShaderUniformFloat)
	rl.SetShaderValue(shader, ampXLoc, ampX, rl.ShaderUniformFloat)
	rl.SetShaderValue(shader, ampYLoc, ampY, rl.ShaderUniformFloat)
	rl.SetShaderValue(shader, speedXLoc, speedX, rl.ShaderUniformFloat)
	rl.SetShaderValue(shader, speedYLoc, speedY, rl.ShaderUniformFloat)

	seconds := []float32{0}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		seconds[0] += rl.GetFrameTime()
		rl.SetShaderValue(shader, secondsLoc, seconds, rl.ShaderUniformFloat)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginShaderMode(shader)
		rl.DrawTexture(texture, 0, 0, rl.White)
		rl.DrawTexture(texture, texture.Width, 0, rl.White)
		rl.EndShaderMode()

		rl.EndDrawing()
	}

	rl.UnloadShader(shader)
	rl.UnloadTexture(texture)

	rl.CloseWindow()
}
