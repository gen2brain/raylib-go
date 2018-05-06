package main

import (
	"github.com/gen2brain/raylib-go/raylib"
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

	raylib.InitWindow(screenWidth, screenHeight, "raylib [textures] example - image processing")

	image := raylib.LoadImage("parrots.png")               // Loaded in CPU memory (RAM)
	raylib.ImageFormat(image, raylib.UncompressedR8g8b8a8) // Format image to RGBA 32bit (required for texture update)
	texture := raylib.LoadTextureFromImage(image)          // Image converted to texture, GPU memory (VRAM)

	currentProcess := None
	textureReload := false

	selectRecs := make([]raylib.Rectangle, numProcesses)

	for i := 0; i < numProcesses; i++ {
		selectRecs[i] = raylib.NewRectangle(40, 50+32*float32(i), 150, 30)
	}

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if raylib.IsKeyPressed(raylib.KeyDown) {
			currentProcess++
			if currentProcess > 7 {
				currentProcess = 0
			}
			textureReload = true
		} else if raylib.IsKeyPressed(raylib.KeyUp) {
			currentProcess--
			if currentProcess < 0 {
				currentProcess = 7
			}
			textureReload = true
		}

		if textureReload {
			raylib.UnloadImage(image)               // Unload current image data
			image = raylib.LoadImage("parrots.png") // Re-load image data

			// NOTE: Image processing is a costly CPU process to be done every frame,
			// If image processing is required in a frame-basis, it should be done
			// with a texture and by shaders
			switch currentProcess {
			case ColorGrayscale:
				raylib.ImageColorGrayscale(image)
				break
			case ColorTint:
				raylib.ImageColorTint(image, raylib.Green)
				break
			case ColorInvert:
				raylib.ImageColorInvert(image)
				break
			case ColorContrast:
				raylib.ImageColorContrast(image, -40)
				break
			case ColorBrightness:
				raylib.ImageColorBrightness(image, -80)
				break
			case FlipVertical:
				raylib.ImageFlipVertical(image)
				break
			case FlipHorizontal:
				raylib.ImageFlipHorizontal(image)
				break
			default:
				break
			}

			pixels := raylib.GetImageData(image)  // Get pixel data from image (RGBA 32bit)
			raylib.UpdateTexture(texture, pixels) // Update texture with new image data

			textureReload = false
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("IMAGE PROCESSING:", 40, 30, 10, raylib.DarkGray)

		// Draw rectangles
		for i := 0; i < numProcesses; i++ {
			if i == currentProcess {
				raylib.DrawRectangleRec(selectRecs[i], raylib.SkyBlue)
				raylib.DrawRectangleLines(int32(selectRecs[i].X), int32(selectRecs[i].Y), int32(selectRecs[i].Width), int32(selectRecs[i].Height), raylib.Blue)
				raylib.DrawText(processText[i], int32(selectRecs[i].X+selectRecs[i].Width/2)-raylib.MeasureText(processText[i], 10)/2, int32(selectRecs[i].Y)+11, 10, raylib.DarkBlue)
			} else {
				raylib.DrawRectangleRec(selectRecs[i], raylib.LightGray)
				raylib.DrawRectangleLines(int32(selectRecs[i].X), int32(selectRecs[i].Y), int32(selectRecs[i].Width), int32(selectRecs[i].Height), raylib.Gray)
				raylib.DrawText(processText[i], int32(selectRecs[i].X+selectRecs[i].Width/2)-raylib.MeasureText(processText[i], 10)/2, int32(selectRecs[i].Y)+11, 10, raylib.DarkGray)
			}
		}

		raylib.DrawTexture(texture, screenWidth-texture.Width-60, screenHeight/2-texture.Height/2, raylib.White)
		raylib.DrawRectangleLines(screenWidth-texture.Width-60, screenHeight/2-texture.Height/2, texture.Width, texture.Height, raylib.Black)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(texture)

	raylib.CloseWindow()
}
