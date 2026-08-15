package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - ascii rendering")

	fudesumi := rl.LoadTexture("fudesumi.png")
	raysan := rl.LoadTexture("raysan.png")
	shader := rl.LoadShader("", "ascii.fs")

	// These locations are used to send data to the GPU
	resolutionLoc := rl.GetShaderLocation(shader, "resolution")
	fontSizeLoc := rl.GetShaderLocation(shader, "fontSize")

	// Set the character size for the ASCII effect (Fontsize should be 9 or more)
	fontSize := float32(9.0)

	// Send the updated values to the shader
	resolution := []float32{float32(screenWidth), float32(screenHeight)}
	rl.SetShaderValue(shader, resolutionLoc, resolution, rl.ShaderUniformVec2)

	circlePos := rl.NewVector2(40.0, float32(screenHeight)*0.5)
	circleSpeed := float32(1.0)

	// RenderTexture to apply postprocessing
	target := rl.LoadRenderTexture(screenWidth, screenHeight)

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		circlePos.X += circleSpeed
		if circlePos.X > 200.0 || circlePos.X < 40.0 {
			circleSpeed *= -1 // Revert speed
		}

		if rl.IsKeyPressed(rl.KeyLeft) && fontSize > 9.0 {
			fontSize -= 1.0 // Reduce fontSize
		}
		if rl.IsKeyPressed(rl.KeyRight) && fontSize < 15.0 {
			fontSize += 1.0 // Increase fontSize
		}

		// Set fontsize for the shader
		rl.SetShaderValue(shader, fontSizeLoc, []float32{fontSize}, rl.ShaderUniformFloat)

		// Draw
		rl.BeginTextureMode(target)

		rl.ClearBackground(rl.White)

		// Draw scene in our render texture
		rl.DrawTexture(fudesumi, 500, -30, rl.White)
		rl.DrawTextureV(raysan, circlePos, rl.White)

		rl.EndTextureMode()

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginShaderMode(shader)

		// Draw scene texture (rendered earlier) to screen with vertical flip
		srcRec := rl.NewRectangle(0, 0, float32(target.Texture.Width), -float32(target.Texture.Height))
		rl.DrawTextureRec(target.Texture, srcRec, rl.Vector2Zero(), rl.White)

		rl.EndShaderMode()

		rl.DrawRectangle(0, 0, screenWidth, 40, rl.Black)
		rl.DrawText(fmt.Sprintf("Ascii effect - FontSize:%2.0f - [Left] -1 [Right] +1 ", fontSize), 120, 10, 20, rl.LightGray)
		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	rl.UnloadTexture(fudesumi)
	rl.UnloadTexture(raysan)
	rl.UnloadShader(shader)
	rl.UnloadRenderTexture(target)

	rl.CloseWindow()
}
