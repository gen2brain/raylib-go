/*******************************************************************************************
*
*   raylib [text] example - Codepoints loading
*
*   Example originally created with raylib 4.2, last time updated with raylib 2.5
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2022-2024 Ramon Santamaria (@raysan5)
*
********************************************************************************************/
package main

import (
	"fmt"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - codepoints loading")

	// Text to be displayed, must be UTF-8 (save this code file as UTF-8)
	// NOTE: It can contain all the required text for the game,
	// this text will be scanned to get all the required codepoints
	text := "いろはにほへと　ちりぬるを\nわかよたれそ　つねならむ\nうゐのおくやま　けふこえて\nあさきゆめみし　ゑひもせす"

	// Get codepoints from text
	allCodepoints := []rune(text)

	// Removed duplicate codepoints to generate smaller font atlas
	slices.Sort(allCodepoints)
	codepoints := slices.Compact(allCodepoints)
	codepointsCount := len(codepoints)

	// Load font containing all the provided codepoint glyphs
	// A texture font atlas is automatically generated
	font := rl.LoadFontEx("DotGothic16-Regular.ttf", 36, codepoints, int32(codepointsCount))

	// Set bi-linear scale filter for better font scaling
	rl.SetTextureFilter(font.Texture, rl.FilterBilinear)
	// Set line spacing for multiline text (when line breaks are included '\n')
	rl.SetTextLineSpacing(20)

	showFontAtlas := false

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		if rl.IsKeyPressed(rl.KeySpace) {
			showFontAtlas = !showFontAtlas
		}

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawRectangle(0, 0, screenWidth, 70, rl.Black)
		msg := fmt.Sprintf("Total codepoints contained in provided text: %d", len(allCodepoints))
		rl.DrawText(msg, 10, 10, 20, rl.Green)
		msg = fmt.Sprintf("Total codepoints required for font atlas (duplicates excluded): %d", codepointsCount)
		rl.DrawText(msg, 10, 40, 20, rl.Green)

		if showFontAtlas {
			// Draw generated font texture atlas containing provided codepoints
			rl.DrawTexture(font.Texture, 150, 100, rl.Black)
			rl.DrawRectangleLines(150, 100, font.Texture.Width, font.Texture.Height, rl.Black)
		} else {
			// Draw provided text with loaded font, containing all required codepoint glyphs
			pos := rl.Vector2{
				X: 160,
				Y: 110,
			}
			rl.DrawTextEx(font, text, pos, 48, 5, rl.Black)
		}

		msg = "Press SPACE to toggle font atlas view!"
		rl.DrawText(msg, 10, screenHeight-30, 20, rl.Gray)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadFont(font) // Unload font
	rl.CloseWindow()    // Close window and OpenGL context
}
