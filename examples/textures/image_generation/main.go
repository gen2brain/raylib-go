package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const numTextures = 7

func main() {
	screenWidth := 800
	screenHeight := 450

	rl.InitWindow(int32(screenWidth), int32(screenHeight), "raylib [textures] example - procedural images generation")

	verticalGradient := rl.GenImageGradientV(screenWidth, screenHeight, rl.Red, rl.Blue)
	horizontalGradient := rl.GenImageGradientH(screenWidth, screenHeight, rl.Red, rl.Blue)
	radialGradient := rl.GenImageGradientRadial(screenWidth, screenHeight, 0, rl.White, rl.Black)
	checked := rl.GenImageChecked(screenWidth, screenHeight, 32, 32, rl.Red, rl.Blue)
	whiteNoise := rl.GenImageWhiteNoise(screenWidth, screenHeight, 0.5)
	cellular := rl.GenImageCellular(screenWidth, screenHeight, 32)

	textures := make([]rl.Texture2D, numTextures)
	textures[0] = rl.LoadTextureFromImage(verticalGradient)
	textures[1] = rl.LoadTextureFromImage(horizontalGradient)
	textures[2] = rl.LoadTextureFromImage(radialGradient)
	textures[3] = rl.LoadTextureFromImage(checked)
	textures[4] = rl.LoadTextureFromImage(whiteNoise)
	textures[5] = rl.LoadTextureFromImage(cellular)

	// Unload image data (CPU RAM)
	rl.UnloadImage(verticalGradient)
	rl.UnloadImage(horizontalGradient)
	rl.UnloadImage(radialGradient)
	rl.UnloadImage(checked)
	rl.UnloadImage(whiteNoise)
	rl.UnloadImage(cellular)

	currentTexture := 0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			currentTexture = (currentTexture + 1) % numTextures // Cycle between the textures
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(textures[currentTexture], 0, 0, rl.White)

		rl.DrawRectangle(30, 400, 325, 30, rl.Fade(rl.SkyBlue, 0.5))
		rl.DrawRectangleLines(30, 400, 325, 30, rl.Fade(rl.White, 0.5))
		rl.DrawText("MOUSE LEFT BUTTON to CYCLE PROCEDURAL TEXTURES", 40, 410, 10, rl.White)

		switch currentTexture {
		case 0:
			rl.DrawText("VERTICAL GRADIENT", 560, 10, 20, rl.RayWhite)
			break
		case 1:
			rl.DrawText("HORIZONTAL GRADIENT", 540, 10, 20, rl.RayWhite)
			break
		case 2:
			rl.DrawText("RADIAL GRADIENT", 580, 10, 20, rl.LightGray)
			break
		case 3:
			rl.DrawText("CHECKED", 680, 10, 20, rl.RayWhite)
			break
		case 4:
			rl.DrawText("WHITE NOISE", 640, 10, 20, rl.Red)
			break
		case 5:
			rl.DrawText("CELLULAR", 670, 10, 20, rl.RayWhite)
			break
		default:
			break
		}

		rl.EndDrawing()
	}

	for _, t := range textures {
		rl.UnloadTexture(t)
	}

	rl.CloseWindow()
}
