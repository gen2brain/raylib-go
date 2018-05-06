package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [text] example - ttf loading")

	msg := "TTF Font"

	// NOTE: Textures/Fonts MUST be loaded after Window initialization (OpenGL context is required)

	fontChars := int32(0)

	// TTF Font loading with custom generation parameters
	font := raylib.LoadFontEx("fonts/KAISG.ttf", 96, 0, &fontChars)

	// Generate mipmap levels to use trilinear filtering
	// NOTE: On 2D drawing it won't be noticeable, it looks like FILTER_BILINEAR
	raylib.GenTextureMipmaps(&font.Texture)

	fontSize := font.BaseSize
	fontPosition := raylib.NewVector2(40, float32(screenHeight)/2+50)
	textSize := raylib.Vector2{}

	raylib.SetTextureFilter(font.Texture, raylib.FilterPoint)
	currentFontFilter := 0 // FilterPoint

	count := int32(0)
	droppedFiles := make([]string, 0)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update
		//----------------------------------------------------------------------------------
		fontSize += raylib.GetMouseWheelMove() * 4.0

		// Choose font texture filter method
		if raylib.IsKeyPressed(raylib.KeyOne) {
			raylib.SetTextureFilter(font.Texture, raylib.FilterPoint)
			currentFontFilter = 0
		} else if raylib.IsKeyPressed(raylib.KeyTwo) {
			raylib.SetTextureFilter(font.Texture, raylib.FilterBilinear)
			currentFontFilter = 1
		} else if raylib.IsKeyPressed(raylib.KeyThree) {
			// NOTE: Trilinear filter won't be noticed on 2D drawing
			raylib.SetTextureFilter(font.Texture, raylib.FilterTrilinear)
			currentFontFilter = 2
		}

		textSize = raylib.MeasureTextEx(font, msg, float32(fontSize), 0)

		if raylib.IsKeyDown(raylib.KeyLeft) {
			fontPosition.X -= 10
		} else if raylib.IsKeyDown(raylib.KeyRight) {
			fontPosition.X += 10
		}

		// Load a dropped TTF file dynamically (at current fontSize)
		if raylib.IsFileDropped() {
			droppedFiles = raylib.GetDroppedFiles(&count)

			if count == 1 { // Only support one ttf file dropped
				raylib.UnloadFont(font)
				font = raylib.LoadFontEx(droppedFiles[0], fontSize, 0, &fontChars)
				raylib.ClearDroppedFiles()
			}
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("Use mouse wheel to change font size", 20, 20, 10, raylib.Gray)
		raylib.DrawText("Use KEY_RIGHT and KEY_LEFT to move text", 20, 40, 10, raylib.Gray)
		raylib.DrawText("Use 1, 2, 3 to change texture filter", 20, 60, 10, raylib.Gray)
		raylib.DrawText("Drop a new TTF font for dynamic loading", 20, 80, 10, raylib.DarkGray)

		raylib.DrawTextEx(font, msg, fontPosition, float32(fontSize), 0, raylib.Black)

		// TODO: It seems texSize measurement is not accurate due to chars offsets...
		//raylib.DrawRectangleLines(int32(fontPosition.X), int32(fontPosition.Y), int32(textSize.X), int32(textSize.Y), raylib.Red)

		raylib.DrawRectangle(0, screenHeight-80, screenWidth, 80, raylib.LightGray)
		raylib.DrawText(fmt.Sprintf("Font size: %02.02f", float32(fontSize)), 20, screenHeight-50, 10, raylib.DarkGray)
		raylib.DrawText(fmt.Sprintf("Text size: [%02.02f, %02.02f]", textSize.X, textSize.Y), 20, screenHeight-30, 10, raylib.DarkGray)
		raylib.DrawText("CURRENT TEXTURE FILTER:", 250, 400, 20, raylib.Gray)

		if currentFontFilter == 0 {
			raylib.DrawText("POINT", 570, 400, 20, raylib.Black)
		} else if currentFontFilter == 1 {
			raylib.DrawText("BILINEAR", 570, 400, 20, raylib.Black)
		} else if currentFontFilter == 2 {
			raylib.DrawText("TRILINEAR", 570, 400, 20, raylib.Black)
		}

		raylib.EndDrawing()
	}

	raylib.UnloadFont(font) // Font unloading

	raylib.ClearDroppedFiles() // Clear internal buffers

	raylib.CloseWindow()
}
