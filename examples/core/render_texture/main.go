package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - render texture")
	defer rl.CloseWindow()

	// Define a render texture to render
	renderTextureWidth := int32(300)
	renderTextureHeight := int32(300)
	target := rl.LoadRenderTexture(renderTextureWidth, renderTextureHeight)
	defer rl.UnloadRenderTexture(target)

	ballPosition := rl.NewVector2(float32(renderTextureWidth)/2.0, float32(renderTextureHeight)/2.0)
	ballSpeed := rl.NewVector2(5.0, 4.0)
	ballRadius := float32(20)

	rotation := float32(0.0)

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		// Ball movement logic
		ballPosition.X += ballSpeed.X
		ballPosition.Y += ballSpeed.Y

		// Check walls collision for bouncing
		if ballPosition.X >= (float32(renderTextureWidth)-ballRadius) || ballPosition.X <= ballRadius {
			ballSpeed.X *= -1.0
		}
		if ballPosition.Y >= (float32(renderTextureHeight)-ballRadius) || ballPosition.Y <= ballRadius {
			ballSpeed.Y *= -1.0
		}

		// Render texture rotation
		rotation += 0.5

		// Draw
		// Draw our scene to the render texture
		rl.BeginTextureMode(target)

		rl.ClearBackground(rl.SkyBlue)

		rl.DrawRectangle(0, 0, 20, 20, rl.Red)
		rl.DrawCircleV(ballPosition, ballRadius, rl.Maroon)

		rl.EndTextureMode()

		// Draw render texture to main framebuffer
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		// Draw our render texture with rotation applied
		// NOTE 1: We set the origin of the texture to the center of the render texture
		// NOTE 2: We flip vertically the texture setting negative source rectangle height
		sourceRec := rl.NewRectangle(0, 0, float32(target.Texture.Width), -float32(target.Texture.Height))
		destRec := rl.NewRectangle(float32(screenWidth)/2.0, float32(screenHeight)/2.0, float32(target.Texture.Width), float32(target.Texture.Height))
		origin := rl.NewVector2(float32(target.Texture.Width)/2.0, float32(target.Texture.Height)/2.0)

		rl.DrawTexturePro(target.Texture, sourceRec, destRec, origin, rotation, rl.White)

		rl.DrawText("DRAWING BOUNCING BALL INSIDE RENDER TEXTURE!", 10, screenHeight-40, 20, rl.Black)

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}
}
