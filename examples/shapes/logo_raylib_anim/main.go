package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - raylib logo animation")

	logoPositionX := screenWidth/2 - 128
	logoPositionY := screenHeight/2 - 128

	framesCounter := 0
	lettersCount := int32(0)

	topSideRecWidth := int32(16)
	leftSideRecHeight := int32(16)

	bottomSideRecWidth := int32(16)
	rightSideRecHeight := int32(16)

	state := 0            // Tracking animation states (State Machine)
	alpha := float32(1.0) // Useful for fading

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
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

			if lettersCount >= 6 { // When all letters have appeared, just fade out everything
				alpha -= 0.02

				if alpha <= 0.0 {
					alpha = 0.0
					state = 4
				}
			}
		} else if state == 4 { // State 4: Reset and Replay
			if raylib.IsKeyPressed(raylib.KeyR) {
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

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RayWhite)

		if state == 0 {
			if (framesCounter/15)%2 == 0 {
				raylib.DrawRectangle(logoPositionX, logoPositionY, 16, 16, raylib.Black)
			}
		} else if state == 1 {
			raylib.DrawRectangle(logoPositionX, logoPositionY, topSideRecWidth, 16, raylib.Black)
			raylib.DrawRectangle(logoPositionX, logoPositionY, 16, leftSideRecHeight, raylib.Black)
		} else if state == 2 {
			raylib.DrawRectangle(logoPositionX, logoPositionY, topSideRecWidth, 16, raylib.Black)
			raylib.DrawRectangle(logoPositionX, logoPositionY, 16, leftSideRecHeight, raylib.Black)

			raylib.DrawRectangle(logoPositionX+240, logoPositionY, 16, rightSideRecHeight, raylib.Black)
			raylib.DrawRectangle(logoPositionX, logoPositionY+240, bottomSideRecWidth, 16, raylib.Black)
		} else if state == 3 {
			raylib.DrawRectangle(logoPositionX, logoPositionY, topSideRecWidth, 16, raylib.Fade(raylib.Black, alpha))
			raylib.DrawRectangle(logoPositionX, logoPositionY+16, 16, leftSideRecHeight-32, raylib.Fade(raylib.Black, alpha))

			raylib.DrawRectangle(logoPositionX+240, logoPositionY+16, 16, rightSideRecHeight-32, raylib.Fade(raylib.Black, alpha))
			raylib.DrawRectangle(logoPositionX, logoPositionY+240, bottomSideRecWidth, 16, raylib.Fade(raylib.Black, alpha))

			raylib.DrawRectangle(screenWidth/2-112, screenHeight/2-112, 224, 224, raylib.Fade(raylib.RayWhite, alpha))

			text := "raylib"
			length := int32(len(text))
			if lettersCount > length {
				lettersCount = length
			}

			raylib.DrawText(text[0:lettersCount], screenWidth/2-44, screenHeight/2+48, 50, raylib.Fade(raylib.Black, alpha))
		} else if state == 4 {
			raylib.DrawText("[R] REPLAY", 340, 200, 20, raylib.Gray)
		}

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
