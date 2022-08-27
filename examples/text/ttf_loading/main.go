package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - ttf loading")

	msg := "TTF Font"

	// NOTE: Textures/Fonts MUST be loaded after Window initialization (OpenGL context is required)

	// TTF Font loading with custom generation parameters
	font := rl.LoadFontEx("fonts/KAISG.ttf", 96, nil)

	// Generate mipmap levels to use trilinear filtering
	// NOTE: On 2D drawing it won't be noticeable, it looks like FILTER_BILINEAR
	rl.GenTextureMipmaps(&font.Texture)

	fontSize := font.BaseSize
	fontPosition := rl.NewVector2(40, float32(screenHeight)/2+50)
	textSize := rl.Vector2{}

	rl.SetTextureFilter(font.Texture, rl.FilterPoint)
	currentFontFilter := 0 // FilterPoint

	count := 0
	droppedFiles := make([]string, 0)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update
		//----------------------------------------------------------------------------------
		fontSize += int32(rl.GetMouseWheelMove() * 4.0)

		// Choose font texture filter method
		if rl.IsKeyPressed(rl.KeyOne) {
			rl.SetTextureFilter(font.Texture, rl.FilterPoint)
			currentFontFilter = 0
		} else if rl.IsKeyPressed(rl.KeyTwo) {
			rl.SetTextureFilter(font.Texture, rl.FilterBilinear)
			currentFontFilter = 1
		} else if rl.IsKeyPressed(rl.KeyThree) {
			// NOTE: Trilinear filter won't be noticed on 2D drawing
			rl.SetTextureFilter(font.Texture, rl.FilterTrilinear)
			currentFontFilter = 2
		}

		textSize = rl.MeasureTextEx(font, msg, float32(fontSize), 0)

		if rl.IsKeyDown(rl.KeyLeft) {
			fontPosition.X -= 10
		} else if rl.IsKeyDown(rl.KeyRight) {
			fontPosition.X += 10
		}

		// Load a dropped TTF file dynamically (at current fontSize)
		if rl.IsFileDropped() {
			droppedFiles = rl.LoadDroppedFiles()
			count = len(droppedFiles)

			if count == 1 { // Only support one ttf file dropped
				rl.UnloadFont(font)
				font = rl.LoadFontEx(droppedFiles[0], fontSize, nil)
				rl.UnloadDroppedFiles()
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("Use mouse wheel to change font size", 20, 20, 10, rl.Gray)
		rl.DrawText("Use KEY_RIGHT and KEY_LEFT to move text", 20, 40, 10, rl.Gray)
		rl.DrawText("Use 1, 2, 3 to change texture filter", 20, 60, 10, rl.Gray)
		rl.DrawText("Drop a new TTF font for dynamic loading", 20, 80, 10, rl.DarkGray)

		rl.DrawTextEx(font, msg, fontPosition, float32(fontSize), 0, rl.Black)

		// TODO: It seems texSize measurement is not accurate due to chars offsets...
		//rl.DrawRectangleLines(int32(fontPosition.X), int32(fontPosition.Y), int32(textSize.X), int32(textSize.Y), rl.Red)

		rl.DrawRectangle(0, screenHeight-80, screenWidth, 80, rl.LightGray)
		rl.DrawText(fmt.Sprintf("Font size: %02.02f", float32(fontSize)), 20, screenHeight-50, 10, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("Text size: [%02.02f, %02.02f]", textSize.X, textSize.Y), 20, screenHeight-30, 10, rl.DarkGray)
		rl.DrawText("CURRENT TEXTURE FILTER:", 250, 400, 20, rl.Gray)

		if currentFontFilter == 0 {
			rl.DrawText("POINT", 570, 400, 20, rl.Black)
		} else if currentFontFilter == 1 {
			rl.DrawText("BILINEAR", 570, 400, 20, rl.Black)
		} else if currentFontFilter == 2 {
			rl.DrawText("TRILINEAR", 570, 400, 20, rl.Black)
		}

		rl.EndDrawing()
	}

	rl.UnloadFont(font) // Font unloading

	rl.UnloadDroppedFiles() // Clear internal buffers

	rl.CloseWindow()
}
