package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const numTextures = 7

func main() {
	screenWidth := 800
	screenHeight := 450

	raylib.InitWindow(int32(screenWidth), int32(screenHeight), "raylib [textures] example - procedural images generation")

	verticalGradient := raylib.GenImageGradientV(screenWidth, screenHeight, raylib.Red, raylib.Blue)
	horizontalGradient := raylib.GenImageGradientH(screenWidth, screenHeight, raylib.Red, raylib.Blue)
	radialGradient := raylib.GenImageGradientRadial(screenWidth, screenHeight, 0, raylib.White, raylib.Black)
	checked := raylib.GenImageChecked(screenWidth, screenHeight, 32, 32, raylib.Red, raylib.Blue)
	whiteNoise := raylib.GenImageWhiteNoise(screenWidth, screenHeight, 0.5)
	perlinNoise := raylib.GenImagePerlinNoise(screenWidth, screenHeight, 50, 50, 4.0)
	cellular := raylib.GenImageCellular(screenWidth, screenHeight, 32)

	textures := make([]raylib.Texture2D, numTextures)
	textures[0] = raylib.LoadTextureFromImage(verticalGradient)
	textures[1] = raylib.LoadTextureFromImage(horizontalGradient)
	textures[2] = raylib.LoadTextureFromImage(radialGradient)
	textures[3] = raylib.LoadTextureFromImage(checked)
	textures[4] = raylib.LoadTextureFromImage(whiteNoise)
	textures[5] = raylib.LoadTextureFromImage(perlinNoise)
	textures[6] = raylib.LoadTextureFromImage(cellular)

	// Unload image data (CPU RAM)
	raylib.UnloadImage(verticalGradient)
	raylib.UnloadImage(horizontalGradient)
	raylib.UnloadImage(radialGradient)
	raylib.UnloadImage(checked)
	raylib.UnloadImage(whiteNoise)
	raylib.UnloadImage(perlinNoise)
	raylib.UnloadImage(cellular)

	currentTexture := 0

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
			currentTexture = (currentTexture + 1) % numTextures // Cycle between the textures
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawTexture(textures[currentTexture], 0, 0, raylib.White)

		raylib.DrawRectangle(30, 400, 325, 30, raylib.Fade(raylib.SkyBlue, 0.5))
		raylib.DrawRectangleLines(30, 400, 325, 30, raylib.Fade(raylib.White, 0.5))
		raylib.DrawText("MOUSE LEFT BUTTON to CYCLE PROCEDURAL TEXTURES", 40, 410, 10, raylib.White)

		switch currentTexture {
		case 0:
			raylib.DrawText("VERTICAL GRADIENT", 560, 10, 20, raylib.RayWhite)
			break
		case 1:
			raylib.DrawText("HORIZONTAL GRADIENT", 540, 10, 20, raylib.RayWhite)
			break
		case 2:
			raylib.DrawText("RADIAL GRADIENT", 580, 10, 20, raylib.LightGray)
			break
		case 3:
			raylib.DrawText("CHECKED", 680, 10, 20, raylib.RayWhite)
			break
		case 4:
			raylib.DrawText("WHITE NOISE", 640, 10, 20, raylib.Red)
			break
		case 5:
			raylib.DrawText("PERLIN NOISE", 630, 10, 20, raylib.RayWhite)
			break
		case 6:
			raylib.DrawText("CELLULAR", 670, 10, 20, raylib.RayWhite)
			break
		default:
			break
		}

		raylib.EndDrawing()
	}

	for _, t := range textures {
		raylib.UnloadTexture(t)
	}

	raylib.CloseWindow()
}
