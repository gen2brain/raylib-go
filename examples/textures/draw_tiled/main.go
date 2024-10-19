/*******************************************************************************************
 *
 *   raylib [textures] example - Draw part of the texture tiled
 *
 *   Example originally created with raylib 3.0, last time updated with raylib 4.2
 *
 *   Example contributed by Vlad Adrian (@demizdor) and reviewed by Ramon Santamaria (@raysan5)
 *
 *   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
 *   BSD-like license that allows static linking with closed source software
 *
 *   Copyright (c) 2020-2024 Vlad Adrian (@demizdor) and Ramon Santamaria (@raysan5)
 *
 ********************************************************************************************/
package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	optWidth     = 220 // Max width for the options container
	marginSize   = 8   // Size for the margins
	colorSize    = 16  // Size of the color select buttons
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable) // Make the window resizable
	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - Draw part of a texture tiled")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	texPattern := rl.LoadTexture("patterns.png")
	rl.SetTextureFilter(texPattern, rl.FilterTrilinear) // Makes the texture smoother when upscaled

	// Coordinates for all patterns inside the texture
	recPattern := []rl.Rectangle{
		{3, 3, 66, 66},
		{75, 3, 100, 100},
		{3, 75, 66, 66},
		{7, 156, 50, 50},
		{85, 106, 90, 45},
		{75, 154, 100, 60},
	}

	// Setup colors
	colors := []rl.Color{rl.Black, rl.Maroon, rl.Orange, rl.Blue, rl.Purple,
		rl.Beige, rl.Lime, rl.Red, rl.DarkGray, rl.SkyBlue}
	var maxColors = len(colors)
	colorRec := make([]rl.Rectangle, maxColors)

	// Calculate rectangle for each color
	var x, y float32
	for i := 0; i < maxColors; i++ {
		colorRec[i].X = 2.0 + marginSize + x
		colorRec[i].Y = 22.0 + 256.0 + marginSize + y
		colorRec[i].Width = colorSize * 2.0
		colorRec[i].Height = colorSize

		if i == (maxColors/2 - 1) {
			x = 0
			y += colorSize + marginSize
		} else {
			x += colorSize*2 + marginSize
		}
	}

	activePattern := 0
	activeCol := 0
	scale := float32(1.0)
	rotation := float32(0.0)

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Handle mouse
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			mouse := rl.GetMousePosition()

			// Check which pattern was clicked and set it as the active pattern
			for i := 0; i < len(recPattern); i++ {
				r := rl.Rectangle{
					X:      2 + marginSize + recPattern[i].X,
					Y:      40 + marginSize + recPattern[i].Y,
					Width:  recPattern[i].Width,
					Height: recPattern[i].Height,
				}
				if rl.CheckCollisionPointRec(mouse, r) {
					activePattern = i
					break
				}
			}

			// Check to see which color was clicked and set it as the active color
			for i := 0; i < maxColors; i++ {
				if rl.CheckCollisionPointRec(mouse, colorRec[i]) {
					activeCol = i
					break
				}
			}
		}

		// Handle keys

		// Change scale
		if rl.IsKeyPressed(rl.KeyUp) {
			scale += 0.25
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			scale -= 0.25
		}
		scale = clamp(scale, 0.25, 10)

		// Change rotation
		if rl.IsKeyPressed(rl.KeyLeft) {
			rotation -= 25.0
		}
		if rl.IsKeyPressed(rl.KeyRight) {
			rotation += 25.0
		}

		// Reset
		if rl.IsKeyPressed(rl.KeySpace) {
			rotation = 0.0
			scale = 1.0
		}

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw the tiled area
		src := rl.Rectangle{
			X:      optWidth + marginSize,
			Y:      marginSize,
			Width:  float32(rl.GetScreenWidth()) - optWidth - 2.0*marginSize,
			Height: float32(rl.GetScreenHeight()) - 2.0*marginSize,
		}
		DrawTextureTiled(texPattern, recPattern[activePattern], src, rl.Vector2{}, rotation, scale, colors[activeCol])

		// Draw options
		rl.DrawRectangle(marginSize, marginSize, optWidth-marginSize, int32(rl.GetScreenHeight())-2*marginSize,
			rl.ColorAlpha(rl.LightGray, 0.5))

		rl.DrawText("Select Pattern", 2+marginSize, 30+marginSize, 10, rl.Black)
		rl.DrawTexture(texPattern, 2+marginSize, 40+marginSize, rl.Black)
		rl.DrawRectangle(int32(2+marginSize+recPattern[activePattern].X),
			int32(40+marginSize+recPattern[activePattern].Y),
			int32(recPattern[activePattern].Width),
			int32(recPattern[activePattern].Height), rl.ColorAlpha(rl.DarkBlue, 0.3))

		rl.DrawText("Select Color", 2+marginSize, 10+256+marginSize, 10, rl.Black)
		for i := 0; i < maxColors; i++ {
			rl.DrawRectangleRec(colorRec[i], colors[i])
			if activeCol == i {
				rl.DrawRectangleLinesEx(colorRec[i], 3, rl.ColorAlpha(rl.White, 0.5))
			}
		}

		rl.DrawText("Scale (UP/DOWN to change)", 2+marginSize, 80+256+marginSize, 10, rl.Black)
		rl.DrawText(fmt.Sprintf("%.2fx", scale), 2+marginSize, 92+256+marginSize, 20, rl.Black)

		rl.DrawText("Rotation (LEFT/RIGHT to change)", 2+marginSize, 122+256+marginSize, 10, rl.Black)
		rl.DrawText(fmt.Sprintf("%.0f degrees", rotation), 2+marginSize, 134+256+marginSize, 20, rl.Black)

		rl.DrawText("Press [SPACE] to reset", 2+marginSize, 164+256+marginSize, 10, rl.DarkBlue)

		// Draw FPS
		rl.DrawText(fmt.Sprintf("%d FPS", rl.GetFPS()), 2+marginSize, 2+marginSize, 20, rl.Black)
		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadTexture(texPattern) // Unload texture

	rl.CloseWindow() // Close window and OpenGL context
}

// DrawTextureTiled draws a part of a texture (defined by a rectangle) with rotation and scale tiled into dest.
func DrawTextureTiled(texture rl.Texture2D, source, dest rl.Rectangle, origin rl.Vector2, rotation, scale float32,
	tint rl.Color) {

	if (texture.ID <= 0) || (scale <= 0.0) { // Want see an infinite loop?!...just delete this line!
		return
	}
	if (source.Width == 0) || (source.Height == 0) {
		return
	}

	tileWidth := source.Width * scale
	tileHeight := source.Height * scale
	if (dest.Width < tileWidth) && (dest.Height < tileHeight) {
		// Can fit only one tile
		src := rl.Rectangle{
			X:      source.X,
			Y:      source.Y,
			Width:  dest.Width / tileWidth * source.Width,
			Height: dest.Height / tileHeight * source.Height,
		}
		dst := rl.Rectangle{X: dest.X, Y: dest.Y, Width: dest.Width, Height: dest.Height}
		rl.DrawTexturePro(texture, src, dst, origin, rotation, tint)
	} else if dest.Width <= tileWidth {
		// Tiled vertically (one column)
		var dy float32
		for ; dy+tileHeight < dest.Height; dy += tileHeight {
			src := rl.Rectangle{
				X:      source.X,
				Y:      source.Y,
				Width:  dest.Width / tileWidth * source.Width,
				Height: source.Height,
			}
			dst := rl.Rectangle{X: dest.X, Y: dest.Y + dy, Width: dest.Width, Height: tileHeight}
			rl.DrawTexturePro(texture, src, dst, origin, rotation, tint)
		}

		// Fit last tile
		if dy < dest.Height {
			src := rl.Rectangle{X: source.X, Y: source.Y,
				Width:  (dest.Width / tileWidth) * source.Width,
				Height: ((dest.Height - dy) / tileHeight) * source.Height,
			}
			dst := rl.Rectangle{X: dest.X, Y: dest.Y + dy, Width: dest.Width, Height: dest.Height - dy}
			rl.DrawTexturePro(texture, src, dst, origin, rotation, tint)
		}
	} else if dest.Height <= tileHeight {
		// Tiled horizontally (one row)
		var dx float32
		for ; dx+tileWidth < dest.Width; dx += tileWidth {
			src := rl.Rectangle{
				X: source.X, Y: source.Y, Width: source.Width,
				Height: (dest.Height / tileHeight) * source.Height,
			}
			dst := rl.Rectangle{X: dest.X + dx, Y: dest.Y, Width: tileWidth, Height: dest.Height}
			rl.DrawTexturePro(texture, src, dst, origin, rotation, tint)
		}

		// Fit last tile
		if dx < dest.Width {
			src := rl.Rectangle{
				X: source.X, Y: source.Y, Width: ((dest.Width - dx) / tileWidth) * source.Width,
				Height: (dest.Height / tileHeight) * source.Height,
			}
			dst := rl.Rectangle{X: dest.X + dx, Y: dest.Y, Width: dest.Width - dx, Height: dest.Height}
			rl.DrawTexturePro(texture, src,
				dst, origin, rotation, tint)
		}
	} else {
		// Tiled both horizontally and vertically (rows and columns)
		var dx float32
		for ; dx+tileWidth < dest.Width; dx += tileWidth {
			var dy float32
			for ; dy+tileHeight < dest.Height; dy += tileHeight {
				dst := rl.Rectangle{X: dest.X + dx, Y: dest.Y + dy, Width: tileWidth, Height: tileHeight}
				rl.DrawTexturePro(texture, source, dst, origin, rotation, tint)
			}

			if dy < dest.Height {
				src := rl.Rectangle{
					X: source.X, Y: source.Y,
					Width: source.Width, Height: ((dest.Height - dy) / tileHeight) * source.Height,
				}
				dst := rl.Rectangle{
					X: dest.X + dx, Y: dest.Y + dy,
					Width: tileWidth, Height: dest.Height - dy,
				}
				rl.DrawTexturePro(texture, src, dst, origin, rotation, tint)
			}
		}

		// Fit last column of tiles
		if dx < dest.Width {
			var dy float32
			for ; dy+tileHeight < dest.Height; dy += tileHeight {
				src := rl.Rectangle{
					X: source.X, Y: source.Y,
					Width: ((dest.Width - dx) / tileWidth) * source.Width, Height: source.Height,
				}
				dst := rl.Rectangle{X: dest.X + dx, Y: dest.Y + dy, Width: dest.Width - dx, Height: tileHeight}
				rl.DrawTexturePro(texture, src, dst, origin, rotation, tint)
			}

			// Draw final tile in the bottom right corner
			if dy < dest.Height {
				src := rl.Rectangle{
					X: source.X, Y: source.Y,
					Width:  ((dest.Width - dx) / tileWidth) * source.Width,
					Height: ((dest.Height - dy) / tileHeight) * source.Height,
				}
				dst := rl.Rectangle{X: dest.X + dx, Y: dest.Y + dy, Width: dest.Width - dx, Height: dest.Height - dy}
				rl.DrawTexturePro(texture, src, dst, origin, rotation, tint)
			}
		}
	}
}

func clamp(value, minValue, maxValue float32) float32 {
	return min(maxValue, max(value, minValue))
}
