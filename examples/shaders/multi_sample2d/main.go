package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagMsaa4xHint) // Enable Multi Sampling Anti Aliasing 4x (if available)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - multiple sample2D")

	imRed := rl.GenImageColor(int(screenWidth), int(screenHeight), rl.NewColor(uint8(255), uint8(0), uint8(0), uint8(255)))
	texRed := rl.LoadTextureFromImage(imRed)
	rl.UnloadImage(imRed)

	imBlue := rl.GenImageColor(int(screenWidth), int(screenHeight), rl.NewColor(uint8(0), uint8(0), uint8(255), uint8(255)))
	texBlue := rl.LoadTextureFromImage(imBlue)
	rl.UnloadImage(imBlue)

	shader := rl.LoadShader("", "color_mix.fs")

	texBlueLoc := rl.GetShaderLocation(shader, "texture1")
	dividerLoc := rl.GetShaderLocation(shader, "divider")
	dividerValue := []float32{0.5}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		if rl.IsKeyDown(rl.KeyRight) {
			dividerValue[0] += 0.01
		} else if rl.IsKeyDown(rl.KeyLeft) {
			dividerValue[0] -= 0.01
		}

		if dividerValue[0] < 0 {
			dividerValue[0] = 0
		} else if dividerValue[0] > 1 {
			dividerValue[0] = 1
		}

		rl.SetShaderValue(shader, dividerLoc, dividerValue, rl.ShaderUniformFloat)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginShaderMode(shader)

		rl.SetShaderValueTexture(shader, texBlueLoc, texBlue)
		rl.DrawTexture(texRed, 0, 0, rl.White)

		rl.EndShaderMode()

		rl.DrawText("USE LEFT/RIGHT ARROW KEYS TO MIX COLORS", 15, 17, 20, rl.Black)
		rl.DrawText("USE LEFT/RIGHT ARROW KEYS TO MIX COLORS", 16, 16, 20, rl.White)

		rl.EndDrawing()
	}

	rl.UnloadShader(shader)   // Unload shader
	rl.UnloadTexture(texRed)  // Unload texture
	rl.UnloadTexture(texBlue) // Unload texture

	rl.CloseWindow()
}
