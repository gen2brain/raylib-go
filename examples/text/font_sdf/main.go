package main

import (
	_ "embed"
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed AnonymousPro-Bold.ttf
var fileData []byte

func main() {
	// Initialization
	const screenWidth, screenHeight = 800, 450
	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - SDF fonts")
	defer rl.CloseWindow() // Close window and OpenGL context

	// NOTE: Textures/Fonts MUST be loaded after Window initialization (OpenGL context is required)

	const msg = "Signed Distance Fields"

	// Default font generation from TTF font
	fontDefault := rl.Font{BaseSize: 16, CharsCount: 95}
	defer rl.UnloadFont(fontDefault) // Default font unloading

	// Loading font data from memory data
	// Parameters > font size: 16, no glyphs array provided (nil), glyphs count: 95 (autogenerate chars array)
	glyphs := rl.LoadFontData(fileData, 16, nil, 95, rl.FontDefault)
	fontDefault.Chars = &glyphs[0]

	// Parameters >  font size: 16, glyphs padding in image: 4 px, pack method: 0 (default)
	atlas := rl.GenImageFontAtlas(unsafe.Slice(fontDefault.Chars, fontDefault.CharsCount), unsafe.Slice(&fontDefault.Recs, fontDefault.CharsCount), 16, 4, 0)
	fontDefault.Texture = rl.LoadTextureFromImage(&atlas)
	rl.UnloadImage(&atlas)

	// SDF font generation from TTF font
	fontSDF := rl.Font{BaseSize: 16, CharsCount: 95}
	defer rl.UnloadFont(fontSDF) // SDF font unloading

	// Parameters > font size: 16, no glyphs array provided (nil), glyphs count: 0 (defaults to 95)
	glyphsSDF := rl.LoadFontData(fileData, 16, nil, 0, rl.FontSdf)
	fontSDF.Chars = &glyphsSDF[0]
	// Parameters > font size: 16, glyphs padding in image: 0 px, pack method: 1 (Skyline algorithm)
	atlas = rl.GenImageFontAtlas(unsafe.Slice(fontSDF.Chars, fontSDF.CharsCount), unsafe.Slice(&fontSDF.Recs, fontSDF.CharsCount), 16, 0, 1)
	fontSDF.Texture = rl.LoadTextureFromImage(&atlas)
	rl.UnloadImage(&atlas)

	// Load SDF required shader (we use default vertex shader)
	shader := rl.LoadShader("", "sdf.fs")
	defer rl.UnloadShader(shader)                           // Unload SDF shader
	rl.SetTextureFilter(fontSDF.Texture, rl.FilterBilinear) // Required for SDF font

	fontPosition := rl.NewVector2(40, screenHeight/2.0-50)
	textSize := rl.Vector2Zero()
	fontSize := float32(16)
	currentFont := 0 // 0 - fontDefault, 1 - fontSDF

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		fontSize += rl.GetMouseWheelMove() * 8.0

		if fontSize < 6 {
			fontSize = 6
		}

		if rl.IsKeyDown(rl.KeySpace) {
			currentFont = 1
			textSize = rl.MeasureTextEx(fontSDF, msg, fontSize, 0)
		} else {
			currentFont = 0
			textSize = rl.MeasureTextEx(fontDefault, msg, fontSize, 0)
		}

		fontPosition.X = float32(rl.GetScreenWidth()/2) - textSize.X/2
		fontPosition.Y = float32(rl.GetScreenHeight()/2) - textSize.Y/2 + 80

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		if currentFont == 1 {
			// NOTE: SDF fonts require a custom SDf shader to compute fragment color
			rl.BeginShaderMode(shader) // Activate SDF font shader
			rl.DrawTextEx(fontSDF, msg, fontPosition, fontSize, 0, rl.Black)
			rl.EndShaderMode() // Activate our default shader for next drawings
			rl.DrawTexture(fontSDF.Texture, 10, 10, rl.Black)
			rl.DrawText("SDF!", 320, 20, 80, rl.Red)
		} else {
			rl.DrawTextEx(fontDefault, msg, fontPosition, fontSize, 0, rl.Black)
			rl.DrawTexture(fontDefault.Texture, 10, 10, rl.Black)
			rl.DrawText("default font", 315, 40, 30, rl.Gray)
		}

		rl.DrawText("FONT SIZE: 16.0", int32(rl.GetScreenWidth()-240), 20, 20, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("RENDER SIZE: %02.02f", fontSize), int32(rl.GetScreenWidth()-240), 50, 20, rl.DarkGray)
		rl.DrawText("Use MOUSE WHEEL to SCALE TEXT!", int32(rl.GetScreenWidth()-240), 90, 10, rl.DarkGray)

		rl.DrawText("HOLD SPACE to USE SDF FONT VERSION!", 340, int32(rl.GetScreenHeight()-30), 20, rl.Maroon)

		rl.EndDrawing()

	}
}
