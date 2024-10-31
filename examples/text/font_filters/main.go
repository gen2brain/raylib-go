/*******************************************************************************************
*
*   raylib [text] example - Font filters
*
*   NOTE: After font loading, font texture atlas filter could be configured for a softer
*   display of the font when scaling it to different sizes, that way, it's not required
*   to generate multiple fonts at multiple sizes (as long as the scaling is not very different)
*
*   Example originally created with raylib 1.3, last time updated with raylib 4.2
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2015-2024 Ramon Santamaria (@raysan5)
*
********************************************************************************************/
package main

import (
	"fmt"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - font filters")
	msg := "Loaded Font"

	// NOTE: Textures/Fonts MUST be loaded after Window initialization (OpenGL context is required)

	// TTF Font loading with custom generation parameters
	font := rl.LoadFontEx("KAISG.ttf", 96, nil, 0)

	// Generate mipmap levels to use trilinear filtering
	// NOTE: On 2D drawing it won't be noticeable, it looks like FILTER_BILINEAR
	rl.GenTextureMipmaps(&font.Texture)

	fontSize := float32(font.BaseSize)
	fontPosition := rl.Vector2{X: 40.0, Y: screenHeight/2.0 - 80.0}
	textSize := rl.Vector2{}

	// Setup texture scaling filter
	rl.SetTextureFilter(font.Texture, rl.FilterPoint)
	currentFontFilter := 0 // TEXTURE_FILTER_POINT

	rl.SetTargetFPS(60)           // Set our game to run at 60 frames-per-second
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		fontSize += rl.GetMouseWheelMove() * 4.0

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

		textSize = rl.MeasureTextEx(font, msg, fontSize, 0)

		if rl.IsKeyDown(rl.KeyLeft) {
			fontPosition.X -= 10
		} else if rl.IsKeyDown(rl.KeyRight) {
			fontPosition.X += 10
		}

		// Load a dropped TTF file dynamically (at current fontSize)
		if rl.IsFileDropped() {
			droppedFiles := rl.LoadDroppedFiles()

			// NOTE: We only support first ttf file dropped
			if filepath.Ext(droppedFiles[0]) == ".ttf" {
				rl.UnloadFont(font)
				font = rl.LoadFontEx(droppedFiles[0], int32(fontSize), nil, 0)
			}

			rl.UnloadDroppedFiles() // Unload file paths from memory
		}

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("Use mouse wheel to change font size", 20, 20, 10, rl.Gray)
		rl.DrawText("Use KEY_RIGHT and KEY_LEFT to move text", 20, 40, 10, rl.Gray)
		rl.DrawText("Use 1, 2, 3 to change texture filter", 20, 60, 10, rl.Gray)
		rl.DrawText("Drop a new TTF font for dynamic loading", 20, 80, 10, rl.DarkGray)
		rl.DrawTextEx(font, msg, fontPosition, fontSize, 0, rl.Black)

		// TODO: It seems texSize measurement is not accurate due to chars offsets...
		//rl.DrawRectangleLines(fontPosition.X, fontPosition.Y, textSize.X, textSize.Y, rl.Red);

		rl.DrawRectangle(0, screenHeight-80, screenWidth, 80, rl.LightGray)
		text := fmt.Sprintf("Font size: %02.02f", fontSize)
		rl.DrawText(text, 20, screenHeight-50, 10, rl.DarkGray)
		text = fmt.Sprintf("Text size: [%02.02f, %02.02f]", textSize.X, textSize.Y)
		rl.DrawText(text, 20, screenHeight-30, 10, rl.DarkGray)
		rl.DrawText("CURRENT TEXTURE FILTER:", 250, 400, 20, rl.Gray)

		switch currentFontFilter {
		case 0:
			rl.DrawText("POINT", 570, 400, 20, rl.Black)
		case 1:
			rl.DrawText("BILINEAR", 570, 400, 20, rl.Black)
		case 2:
			rl.DrawText("TRILINEAR", 570, 400, 20, rl.Black)
		}

		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadFont(font) // Font unloading
	rl.CloseWindow()    // Close window and OpenGL context
}
