package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialization
	//--------------------------------------------------------------------------------------
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - raylib logo animation")

	var logoPositionX int32 = screenWidth/2 - 128
	var logoPositionY int32 = screenHeight/2 - 128

	framesCounter := 0
	lettersCount := 0

	var topSideRecWidth int32 = 16
	var leftSideRecHeight int32 = 16

	var bottomSideRecWidth int32 = 16
	var rightSideRecHeight int32 = 16

	state := 0            // Tracking animation states (State Machine)
	alpha := float32(1.0) // Useful for fading

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second
	//--------------------------------------------------------------------------------------

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		//----------------------------------------------------------------------------------
		if state == 0 { // State 0: Small box blinking
			framesCounter++

			if framesCounter == 120 {
				state = 1
				framesCounter = 0 // Reset counter... will be used later...
			}
		} else if state == 1 { // State 1: Top and left bars growing
			topSideRecWidth += 4
			leftSideRecHeight += 4

			if topSideRecWidth == 256 {
				state = 2
			}
		} else if state == 2 { // State 2: Bottom and right bars growing
			bottomSideRecWidth += 4
			rightSideRecHeight += 4

			if bottomSideRecWidth == 256 {
				state = 3
			}
		} else if state == 3 { // State 3: Letters appearing (one by one)
			framesCounter++

			if framesCounter%12 == 0 { // Every 12 frames, one more letter!
				lettersCount++
				framesCounter = 0
			}

			if lettersCount >= 10 { // When all letters have appeared, just fade out everything
				alpha -= 0.02

				if alpha <= 0.0 {
					alpha = 0.0
					state = 4
				}
			}
		} else if state == 4 { // State 4: Reset and Replay
			if rl.IsKeyPressed(rl.KeyR) {
				framesCounter = 0
				lettersCount = 0

				topSideRecWidth = 16
				leftSideRecHeight = 16

				bottomSideRecWidth = 16
				rightSideRecHeight = 16

				alpha = 1.0
				state = 0 // Return to State 0
			}
		}
		//----------------------------------------------------------------------------------

		// Draw
		//----------------------------------------------------------------------------------
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		if state == 0 {
			if (framesCounter/15)%2 == 1 {
				rl.DrawRectangle(logoPositionX, logoPositionY, 16, 16, rl.Black)
			}
		} else if state == 1 {
			rl.DrawRectangle(logoPositionX, logoPositionY, topSideRecWidth, 16, rl.Black)
			rl.DrawRectangle(logoPositionX, logoPositionY, 16, leftSideRecHeight, rl.Black)
		} else if state == 2 {
			rl.DrawRectangle(logoPositionX, logoPositionY, topSideRecWidth, 16, rl.Black)
			rl.DrawRectangle(logoPositionX, logoPositionY, 16, leftSideRecHeight, rl.Black)

			rl.DrawRectangle(logoPositionX+240, logoPositionY, 16, rightSideRecHeight, rl.Black)
			rl.DrawRectangle(logoPositionX, logoPositionY+240, bottomSideRecWidth, 16, rl.Black)
		} else if state == 3 {
			rl.DrawRectangle(logoPositionX, logoPositionY, topSideRecWidth, 16, rl.Fade(rl.Black, alpha))
			rl.DrawRectangle(logoPositionX, logoPositionY+16, 16, leftSideRecHeight-32, rl.Fade(rl.Black, alpha))

			rl.DrawRectangle(logoPositionX+240, logoPositionY+16, 16, rightSideRecHeight-32, rl.Fade(rl.Black, alpha))
			rl.DrawRectangle(logoPositionX, logoPositionY+240, bottomSideRecWidth, 16, rl.Fade(rl.Black, alpha))

			rl.DrawRectangle(screenWidth/2-112, screenHeight/2-112, 224, 224, rl.Fade(rl.RayWhite, alpha))

			rl.DrawText("raylib", screenWidth/2-44, screenHeight/2+48, 50, rl.Fade(rl.Black, alpha))
		} else if state == 4 {
			rl.DrawText("[R] REPLAY", 340, 200, 20, rl.Gray)
		}

		rl.EndDrawing()
		//----------------------------------------------------------------------------------
	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	rl.CloseWindow() // Close window and OpenGL context
	//--------------------------------------------------------------------------------------
}