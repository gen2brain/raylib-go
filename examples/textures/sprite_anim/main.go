/*******************************************************************************************
*
*   raylib [textures] example - Sprite animation
*
*   Example originally created with raylib 1.3, last time updated with raylib 1.3
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2014-2024 Ramon Santamaria (@raysan5)
*
********************************************************************************************/
package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth    = 800
	screenHeight   = 450
	maxFramesSpeed = 15
	minFramesSpeed = 1
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [texture] example - sprite anim")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	scarfy := rl.LoadTexture("scarfy.png") // Texture loading

	position := rl.Vector2{X: 350.0, Y: 280.0}
	frameRec := rl.Rectangle{Width: float32(scarfy.Width / 6), Height: float32(scarfy.Height)}
	var currentFrame, framesCounter, framesSpeed int32 = 0, 0, 8

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		framesCounter++

		if framesCounter >= (60 / framesSpeed) {
			framesCounter = 0
			currentFrame++

			if currentFrame > 5 {
				currentFrame = 0
			}

			frameRec.X = float32(currentFrame*scarfy.Width) / 6
		}

		// Control frames speed
		if rl.IsKeyPressed(rl.KeyRight) {
			framesSpeed++
		} else if rl.IsKeyPressed(rl.KeyLeft) {
			framesSpeed--
		}

		// Make sure that framesSpeed is between minFramesSpeed and maxFramesSpeed
		framesSpeed = clamp(framesSpeed, minFramesSpeed, maxFramesSpeed)

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(scarfy, 15, 40, rl.White)
		rl.DrawRectangleLines(15, 40, scarfy.Width, scarfy.Height, rl.Lime)
		rl.DrawRectangleLines(15+int32(frameRec.X), 40+int32(frameRec.Y), int32(frameRec.Width), int32(frameRec.Height), rl.Red)

		rl.DrawText("FRAME SPEED: ", 165, 210, 10, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("%02d FPS", framesSpeed), 575, 210, 10, rl.DarkGray)
		rl.DrawText("PRESS RIGHT/LEFT KEYS to CHANGE SPEED!", 290, 240, 10, rl.DarkGray)

		for i := int32(0); i < maxFramesSpeed; i++ {
			if i < framesSpeed {
				rl.DrawRectangle(250+21*i, 205, 20, 20, rl.Red)
			}
			rl.DrawRectangleLines(250+21*i, 205, 20, 20, rl.Maroon)
		}

		rl.DrawTextureRec(scarfy, frameRec, position, rl.White) // Draw part of the texture

		rl.DrawText("(c) Scarfy sprite by Eiden Marsal", screenWidth-200, screenHeight-20, 10, rl.Gray)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadTexture(scarfy) // Texture unloading
	rl.CloseWindow()         // Close window and OpenGL context
}

func clamp(value, minValue, maxValue int32) int32 {
	return max(minValue, min(value, maxValue))
}
