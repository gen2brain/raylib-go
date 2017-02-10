package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [text] example - font selector")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	fonts := make([]raylib.SpriteFont, 8) // SpriteFont array

	fonts[0] = raylib.LoadSpriteFont("fonts/alagard.rbmf")       // SpriteFont loading
	fonts[1] = raylib.LoadSpriteFont("fonts/pixelplay.rbmf")     // SpriteFont loading
	fonts[2] = raylib.LoadSpriteFont("fonts/mecha.rbmf")         // SpriteFont loading
	fonts[3] = raylib.LoadSpriteFont("fonts/setback.rbmf")       // SpriteFont loading
	fonts[4] = raylib.LoadSpriteFont("fonts/romulus.rbmf")       // SpriteFont loading
	fonts[5] = raylib.LoadSpriteFont("fonts/pixantiqua.rbmf")    // SpriteFont loading
	fonts[6] = raylib.LoadSpriteFont("fonts/alpha_beta.rbmf")    // SpriteFont loading
	fonts[7] = raylib.LoadSpriteFont("fonts/jupiter_crash.rbmf") // SpriteFont loading

	currentFont := 0 // Selected font

	colors := [8]raylib.Color{raylib.Maroon, raylib.Orange, raylib.DarkGreen, raylib.DarkBlue, raylib.DarkPurple, raylib.Lime, raylib.Gold, raylib.Red}

	fontNames := [8]string{"[0] Alagard", "[1] PixelPlay", "[2] MECHA", "[3] Setback", "[4] Romulus", "[5] PixAntiqua", "[6] Alpha Beta", "[7] Jupiter Crash"}

	text := "THIS is THE FONT you SELECTED!" // Main text

	textSize := raylib.MeasureTextEx(fonts[currentFont], text, float32(fonts[currentFont].BaseSize)*3, 1)

	mousePoint := raylib.Vector2{}

	btnNextOutColor := raylib.DarkBlue // Button color (outside line)
	btnNextInColor := raylib.SkyBlue   // Button color (inside)

	framesCounter := 0 // Useful to count frames button is 'active' = clicked

	positionY := int32(180) // Text selector and button Y position

	btnNextRec := raylib.NewRectangle(673, positionY, 109, 44) // Button rectangle (useful for collision)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update

		// Keyboard-based font selection (easy)
		if raylib.IsKeyPressed(raylib.KeyRight) {
			if currentFont < 7 {
				currentFont++
			}
		}

		if raylib.IsKeyPressed(raylib.KeyLeft) {
			if currentFont > 0 {
				currentFont--
			}
		}

		if raylib.IsKeyPressed('0') {
			currentFont = 0
		} else if raylib.IsKeyPressed('1') {
			currentFont = 1
		} else if raylib.IsKeyPressed('2') {
			currentFont = 2
		} else if raylib.IsKeyPressed('3') {
			currentFont = 3
		} else if raylib.IsKeyPressed('4') {
			currentFont = 4
		} else if raylib.IsKeyPressed('5') {
			currentFont = 5
		} else if raylib.IsKeyPressed('6') {
			currentFont = 6
		} else if raylib.IsKeyPressed('7') {
			currentFont = 7
		}

		// Mouse-based font selection (NEXT button logic)
		mousePoint = raylib.GetMousePosition()

		if raylib.CheckCollisionPointRec(mousePoint, btnNextRec) {
			// Mouse hover button logic
			if framesCounter == 0 {
				btnNextOutColor = raylib.DarkPurple
				btnNextInColor = raylib.Purple
			}

			if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
				framesCounter = 20 // Frames button is 'active'
				btnNextOutColor = raylib.Maroon
				btnNextInColor = raylib.Red
			}
		} else {
			// Mouse not hover button
			btnNextOutColor = raylib.DarkBlue
			btnNextInColor = raylib.SkyBlue
		}

		if framesCounter > 0 {
			framesCounter--
		}

		if framesCounter == 1 { // We change font on frame 1
			currentFont++
			if currentFont > 7 {
				currentFont = 0
			}
		}

		// Text measurement for better positioning on screen
		textSize = raylib.MeasureTextEx(fonts[currentFont], text, float32(fonts[currentFont].BaseSize)*3, 1)

		// Draw
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("font selector - use arroys, button or numbers", 160, 80, 20, raylib.DarkGray)
		raylib.DrawLine(120, 120, 680, 120, raylib.DarkGray)

		raylib.DrawRectangle(18, positionY, 644, 44, raylib.DarkGray)
		raylib.DrawRectangle(20, positionY+2, 640, 40, raylib.LightGray)
		raylib.DrawText(fontNames[currentFont], 30, positionY+13, 20, raylib.Black)
		raylib.DrawText("< >", 610, positionY+8, 30, raylib.Black)

		raylib.DrawRectangleRec(btnNextRec, btnNextOutColor)
		raylib.DrawRectangle(675, positionY+2, 105, 40, btnNextInColor)
		raylib.DrawText("NEXT", 700, positionY+13, 20, btnNextOutColor)

		raylib.DrawTextEx(fonts[currentFont], text, raylib.NewVector2(float32(screenWidth)/2-textSize.X/2, 260+(70-textSize.Y)/2), float32(fonts[currentFont].BaseSize*3), 1, colors[currentFont])

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
