package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const numProcesses = 8

// Process
const (
	None = iota
	ColorGrayscale
	ColorTint
	ColorInvert
	ColorContrast
	ColorBrightness
	FlipVertical
	FlipHorizontal
)

var processText = []string{
	"NO PROCESSING",
	"COLOR GRAYSCALE",
	"COLOR TINT",
	"COLOR INVERT",
	"COLOR CONTRAST",
	"COLOR BRIGHTNESS",
	"FLIP VERTICAL",
	"FLIP HORIZONTAL",
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - image processing")

	image := rl.LoadImage("parrots.png")           // Loaded in CPU memory (RAM)
	rl.ImageFormat(image, rl.UncompressedR8g8b8a8) // Format image to RGBA 32bit (required for texture update)
	texture := rl.LoadTextureFromImage(image)      // Image converted to texture, GPU memory (VRAM)

	currentProcess := None
	textureReload := false

	selectRecs := make([]rl.Rectangle, numProcesses)

	for i := 0; i < numProcesses; i++ {
		selectRecs[i] = rl.NewRectangle(40, 50+32*float32(i), 150, 30)
	}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyDown) {
			currentProcess++
			if currentProcess > 7 {
				currentProcess = 0
			}
			textureReload = true
		} else if rl.IsKeyPressed(rl.KeyUp) {
			currentProcess--
			if currentProcess < 0 {
				currentProcess = 7
			}
			textureReload = true
		}

		if textureReload {
			rl.UnloadImage(image)               // Unload current image data
			image = rl.LoadImage("parrots.png") // Re-load image data

			// NOTE: Image processing is a costly CPU process to be done every frame,
			// If image processing is required in a frame-basis, it should be done
			// with a texture and by shaders
			switch currentProcess {
			case ColorGrayscale:
				rl.ImageColorGrayscale(image)
				break
			case ColorTint:
				rl.ImageColorTint(image, rl.Green)
				break
			case ColorInvert:
				rl.ImageColorInvert(image)
				break
			case ColorContrast:
				rl.ImageColorContrast(image, -40)
				break
			case ColorBrightness:
				rl.ImageColorBrightness(image, -80)
				break
			case FlipVertical:
				rl.ImageFlipVertical(image)
				break
			case FlipHorizontal:
				rl.ImageFlipHorizontal(image)
				break
			default:
				break
			}

			pixels := rl.LoadImageColors(image) // Get pixel data from image (RGBA 32bit)
			rl.UpdateTexture(texture, pixels)   // Update texture with new image data

			textureReload = false
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("IMAGE PROCESSING:", 40, 30, 10, rl.DarkGray)

		// Draw rectangles
		for i := 0; i < numProcesses; i++ {
			if i == currentProcess {
				rl.DrawRectangleRec(selectRecs[i], rl.SkyBlue)
				rl.DrawRectangleLines(int32(selectRecs[i].X), int32(selectRecs[i].Y), int32(selectRecs[i].Width), int32(selectRecs[i].Height), rl.Blue)
				rl.DrawText(processText[i], int32(selectRecs[i].X+selectRecs[i].Width/2)-rl.MeasureText(processText[i], 10)/2, int32(selectRecs[i].Y)+11, 10, rl.DarkBlue)
			} else {
				rl.DrawRectangleRec(selectRecs[i], rl.LightGray)
				rl.DrawRectangleLines(int32(selectRecs[i].X), int32(selectRecs[i].Y), int32(selectRecs[i].Width), int32(selectRecs[i].Height), rl.Gray)
				rl.DrawText(processText[i], int32(selectRecs[i].X+selectRecs[i].Width/2)-rl.MeasureText(processText[i], 10)/2, int32(selectRecs[i].Y)+11, 10, rl.DarkGray)
			}
		}

		rl.DrawTexture(texture, screenWidth-texture.Width-60, screenHeight/2-texture.Height/2, rl.White)
		rl.DrawRectangleLines(screenWidth-texture.Width-60, screenHeight/2-texture.Height/2, texture.Width, texture.Height, rl.Black)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)

	rl.CloseWindow()
}
