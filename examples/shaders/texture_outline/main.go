package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - texture outline")

	texture := rl.LoadTexture("gopher.png")

	shader := rl.LoadShader("", "outline.fs")

	cnt := rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2))

	outlineSize := []float32{2}
	outlineColor := []float32{1, 0, 0, 1}
	textureSize := []float32{float32(texture.Width), float32(texture.Height)}

	outlineSizeLoc := rl.GetShaderLocation(shader, "outlineSize")
	outlineColorLoc := rl.GetShaderLocation(shader, "outlineColor")
	textureSizeLoc := rl.GetShaderLocation(shader, "textureSize")

	rl.SetShaderValue(shader, outlineSizeLoc, outlineSize, rl.ShaderUniformFloat)
	rl.SetShaderValue(shader, outlineColorLoc, outlineColor, rl.ShaderUniformVec4)
	rl.SetShaderValue(shader, textureSizeLoc, textureSize, rl.ShaderUniformVec2)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		if rl.IsKeyPressed(rl.KeyUp) {
			outlineSize[0]++
		} else if rl.IsKeyPressed(rl.KeyDown) {
			outlineSize[0]--
			if outlineSize[0] < 1 {
				outlineSize[0] = 1
			}
		}

		rl.SetShaderValue(shader, outlineSizeLoc, outlineSize, rl.ShaderUniformFloat)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginShaderMode(shader)

		rl.DrawTexture(texture, int32(cnt.X)-texture.Width/2, int32(cnt.Y)-texture.Height/2, rl.White)

		rl.EndShaderMode()

		rl.DrawText("UP/DOWN ARROW KEYS INCREASE OUTLINE THICKNESS", 10, 10, 10, rl.Black)

		rl.EndDrawing()
	}

	rl.UnloadShader(shader)
	rl.UnloadTexture(texture)

	rl.CloseWindow()
}
