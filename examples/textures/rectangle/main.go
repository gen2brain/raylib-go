package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	maxFrameSpeed = 15
	minFrameSpeed = 1
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - texture loading and drawing")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	scarfy := rl.LoadTexture("scarfy.png") // Texture loading

	position := rl.NewVector2(350.0, 280.0)
	frameRec := rl.NewRectangle(0, 0, float32(scarfy.Width/6), float32(scarfy.Height))
	currentFrame := float32(0)

	framesCounter := 0
	framesSpeed := 8 // Number of spritesheet frames shown by second

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		framesCounter++

		if framesCounter >= (60 / framesSpeed) {
			framesCounter = 0
			currentFrame++

			if currentFrame > 5 {
				currentFrame = 0
			}

			frameRec.X = currentFrame * float32(scarfy.Width) / 6
		}

		if rl.IsKeyPressed(rl.KeyRight) {
			framesSpeed++
		} else if rl.IsKeyPressed(rl.KeyLeft) {
			framesSpeed--
		}

		if framesSpeed > maxFrameSpeed {
			framesSpeed = maxFrameSpeed
		} else if framesSpeed < minFrameSpeed {
			framesSpeed = minFrameSpeed
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(scarfy, 15, 40, rl.White)
		rl.DrawRectangleLines(15, 40, scarfy.Width, scarfy.Height, rl.Lime)
		rl.DrawRectangleLines(15+int32(frameRec.X), 40+int32(frameRec.Y), int32(frameRec.Width), int32(frameRec.Height), rl.Red)

		rl.DrawText("FRAME SPEED: ", 165, 210, 10, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("%02d FPS", framesSpeed), 575, 210, 10, rl.DarkGray)
		rl.DrawText("PRESS RIGHT/LEFT KEYS to CHANGE SPEED!", 290, 240, 10, rl.DarkGray)

		for i := 0; i < maxFrameSpeed; i++ {
			if i < framesSpeed {
				rl.DrawRectangle(int32(250+21*i), 205, 20, 20, rl.Red)
			}
			rl.DrawRectangleLines(int32(250+21*i), 205, 20, 20, rl.Maroon)
		}

		rl.DrawTextureRec(scarfy, frameRec, position, rl.White) // Draw part of the texture

		rl.DrawText("(c) Scarfy sprite by Eiden Marsal", screenWidth-200, screenHeight-20, 10, rl.Gray)

		rl.EndDrawing()
	}

	rl.UnloadTexture(scarfy)

	rl.CloseWindow()
}
