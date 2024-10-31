/*******************************************************************************************
*
*   raylib [text] example - Font loading
*
*   NOTE: raylib can load fonts from multiple input file formats:
*
*     - TTF/OTF > Sprite font atlas is generated on loading, user can configure
*                 some of the generation parameters (size, characters to include)
*     - BMFonts > Angel code font fileformat, sprite font image must be provided
*                 together with the .fnt file, font generation cna not be configured
*     - XNA Spritefont > Sprite font image, following XNA Spritefont conventions,
*                 Characters in image must follow some spacing and order rules
*
*   Example originally created with raylib 1.4, last time updated with raylib 3.0
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2016-2024 Ramon Santamaria (@raysan5)
*
********************************************************************************************/
package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	screenWidth  = 800
	screenHeight = 450
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - font loading")

	// Define characters to draw
	// NOTE: raylib supports UTF-8 encoding, following list is actually codified as UTF8 internally
	msg := "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHI\nJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmn\nopqrstuvwxyz" +
		"{|}~¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓ\nÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷\nøùúûüýþÿ"

	// NOTE: Textures/Fonts MUST be loaded after Window initialization (OpenGL context is required)

	// BMFont (AngelCode) : Font data and image atlas have been generated using external program
	fontBm := rl.LoadFont("pixantiqua.fnt")

	// TTF font : Font data and atlas are generated directly from TTF
	// NOTE: We define a font base size of 32 pixels tall and up-to 250 characters
	fontTtf := rl.LoadFontEx("pixantiqua.ttf", 32, nil, 250)

	rl.SetTextLineSpacing(16) // Set line spacing for multiline text (when line breaks are included '\n')

	useTtf := false

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		useTtf = rl.IsKeyDown(rl.KeySpace)

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("Hold SPACE to use TTF generated font", 20, 20, 20, rl.LightGray)

		if !useTtf {
			rl.DrawTextEx(fontBm, msg, rl.Vector2{X: 20.0, Y: 100.0}, float32(fontBm.BaseSize), 2, rl.Maroon)
			rl.DrawText("Using BMFont (Angelcode) imported", 20, int32(rl.GetScreenHeight())-30, 20, rl.Gray)
		} else {
			rl.DrawTextEx(fontTtf, msg, rl.Vector2{X: 20.0, Y: 100.0}, float32(fontTtf.BaseSize), 2, rl.Lime)
			rl.DrawText("Using TTF font generated", 20, int32(rl.GetScreenHeight())-30, 20, rl.Gray)
		}

		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadFont(fontBm)  // AngelCode Font unloading
	rl.UnloadFont(fontTtf) // TTF Font unloading
	rl.CloseWindow()       // Close window and OpenGL context
}
