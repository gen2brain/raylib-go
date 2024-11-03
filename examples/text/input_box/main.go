/*******************************************************************************************
 *
 *   raylib [text] example - Input Box
 *
 *   Example originally created with raylib 1.7, last time updated with raylib 3.5
 *
 *   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
 *   BSD-like license that allows static linking with closed source software
 *
 *   Copyright (c) 2017-2024 Ramon Santamaria (@raysan5)
 *
 ********************************************************************************************/
package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth   = 800
	screenHeight  = 450
	maxInputChars = 9
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - input box")

	var name [maxInputChars]rune
	var letterCount, framesCounter int32
	var mouseOnText bool

	textBox := rl.Rectangle{X: screenWidth/2.0 - 100, Y: 180, Width: 225, Height: 50}

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		mouseOnText = rl.CheckCollisionPointRec(rl.GetMousePosition(), textBox)

		if mouseOnText {
			// Set the window's cursor to the I-Beam
			rl.SetMouseCursor(rl.MouseCursorIBeam)

			// Get char pressed (unicode character) on the queue
			key := rl.GetCharPressed()

			// Check if more characters have been pressed on the same frame
			for key > 0 {
				// NOTE: Only allow keys in range [32..125]
				if (key >= 32) && (key <= 125) && (letterCount < maxInputChars) {
					name[letterCount] = key
					letterCount++
				}

				key = rl.GetCharPressed() // Check next character in the queue
			}

			if rl.IsKeyPressed(rl.KeyBackspace) {
				letterCount--
				if letterCount < 0 {
					letterCount = 0
				}
				name[letterCount] = 0
			}
		} else {
			rl.SetMouseCursor(rl.MouseCursorDefault)
		}

		if mouseOnText {
			framesCounter++
		} else {
			framesCounter = 0
		}

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("PLACE MOUSE OVER INPUT BOX!", 240, 140, 20, rl.Gray)

		rl.DrawRectangleRec(textBox, rl.LightGray)
		col := rl.DarkGray
		if mouseOnText {
			col = rl.Red
		}
		rl.DrawRectangleLines(int32(textBox.X), int32(textBox.Y), int32(textBox.Width), int32(textBox.Height), col)

		rl.DrawText(getString(name), int32(textBox.X)+5, int32(textBox.Y)+8, 40, rl.Maroon)
		text := fmt.Sprintf("INPUT CHARS: %d/%d", letterCount, maxInputChars)
		rl.DrawText(text, 315, 250, 20, rl.DarkGray)

		if mouseOnText {
			if letterCount < maxInputChars {
				// Draw blinking underscore char
				if ((framesCounter / 20) % 2) == 0 {
					x := int32(textBox.X) + 8 + rl.MeasureText(getString(name), 40)
					rl.DrawText("_", x, int32(textBox.Y)+12, 40, rl.Maroon)
				}
			} else {
				rl.DrawText("Press BACKSPACE to delete chars...", 230, 300, 20, rl.Gray)
			}
		}

		rl.EndDrawing()
	}

	// De-Initialization
	rl.CloseWindow() // Close window and OpenGL context
}

func getString(r [maxInputChars]rune) string {
	var s string
	for _, char := range r {
		if char == 0 {
			return s
		}
		s += string(char)
	}
	return s
}
