package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MAX_FRAME_SPEED = 15
	MIN_FRAME_SPEED = 1
)

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [texture] example - sprite anim")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	spriteFrames := float32(6)
	scarfy := rl.LoadTexture("scarfy.png") // Texture loading

	position := rl.Vector2{X: 350.0, Y: 280.0}
	frameRec := rl.Rectangle{X: 0.0, Y: 0.0, Width: float32(scarfy.Width) / spriteFrames, Height: float32(scarfy.Height)}
	currentFrame := 0

	framesCounter := 0
	framesSpeed := 8 // Number of spritesheet frames shown by second

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		framesCounter++

		if framesCounter >= (60 / framesSpeed) {
			framesCounter = 0
			currentFrame++

			if currentFrame > int(spriteFrames)-1 {
				currentFrame = 0
			}

			frameRec.X = float32(currentFrame) * float32(scarfy.Width) / spriteFrames // select which image to show
		}

		// Control frames speed
		if rl.IsKeyPressed(rl.KeyRight) {
			framesSpeed++
		} else if rl.IsKeyPressed(rl.KeyLeft) {
			framesSpeed--
		}

		if framesSpeed > MAX_FRAME_SPEED {
			framesSpeed = MAX_FRAME_SPEED
		} else if framesSpeed < MIN_FRAME_SPEED {
			framesSpeed = MIN_FRAME_SPEED
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(scarfy, 15, 40, rl.White)
		rl.DrawRectangleLines(15, 40, int32(scarfy.Width), int32(scarfy.Height), rl.Lime)
		rl.DrawRectangleLines(15+int32(frameRec.X), 40+int32(frameRec.Y), int32(frameRec.Width), int32(frameRec.Height), rl.Red)

		rl.DrawText("FRAME SPEED: ", 165, 210, 10, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("%02d FPS", framesSpeed), 575, 210, 10, rl.DarkGray)
		rl.DrawText("PRESS RIGHT/LEFT KEYS to CHANGE SPEED!", 290, 240, 10, rl.DarkGray)

		for i := 0; i < MAX_FRAME_SPEED; i++ {
			if i < framesSpeed {
				rl.DrawRectangle(250+21*int32(i), screenHeight/2.5+15, 20, 20, rl.Red)
			}
			rl.DrawRectangleLines(250+21*int32(i), screenHeight/2.5+15, 20, 20, rl.Maroon)
		}

		rl.DrawTextureRec(scarfy, frameRec, position, rl.White) // Draw part of the texture

		rl.DrawText("(c) Scarfy sprite by Eiden Marsal", screenWidth-200, screenHeight-20, 10, rl.Gray)

		rl.EndDrawing()
	}

	// De-Initialization
	defer rl.UnloadTexture(scarfy) // Texture unloading

	defer rl.CloseWindow() // Close window and OpenGL context
}
