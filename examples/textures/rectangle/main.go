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

	raylib.InitWindow(screenWidth, screenHeight, "raylib [textures] example - texture loading and drawing")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	scarfy := raylib.LoadTexture("scarfy.png") // Texture loading

	position := raylib.NewVector2(350.0, 280.0)
	frameRec := raylib.NewRectangle(0, 0, float32(scarfy.Width/6), float32(scarfy.Height))
	currentFrame := float32(0)

	framesCounter := 0
	framesSpeed := 8 // Number of spritesheet frames shown by second

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		framesCounter++

		if framesCounter >= (60 / framesSpeed) {
			framesCounter = 0
			currentFrame++

			if currentFrame > 5 {
				currentFrame = 0
			}

			frameRec.X = currentFrame * float32(scarfy.Width) / 6
		}

		if raylib.IsKeyPressed(raylib.KeyRight) {
			framesSpeed++
		} else if raylib.IsKeyPressed(raylib.KeyLeft) {
			framesSpeed--
		}

		if framesSpeed > maxFrameSpeed {
			framesSpeed = maxFrameSpeed
		} else if framesSpeed < minFrameSpeed {
			framesSpeed = minFrameSpeed
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawTexture(scarfy, 15, 40, raylib.White)
		raylib.DrawRectangleLines(15, 40, scarfy.Width, scarfy.Height, raylib.Lime)
		raylib.DrawRectangleLines(15+int32(frameRec.X), 40+int32(frameRec.Y), int32(frameRec.Width), int32(frameRec.Height), raylib.Red)

		raylib.DrawText("FRAME SPEED: ", 165, 210, 10, raylib.DarkGray)
		raylib.DrawText(fmt.Sprintf("%02d FPS", framesSpeed), 575, 210, 10, raylib.DarkGray)
		raylib.DrawText("PRESS RIGHT/LEFT KEYS to CHANGE SPEED!", 290, 240, 10, raylib.DarkGray)

		for i := 0; i < maxFrameSpeed; i++ {
			if i < framesSpeed {
				raylib.DrawRectangle(int32(250+21*i), 205, 20, 20, raylib.Red)
			}
			raylib.DrawRectangleLines(int32(250+21*i), 205, 20, 20, raylib.Maroon)
		}

		raylib.DrawTextureRec(scarfy, frameRec, position, raylib.White) // Draw part of the texture

		raylib.DrawText("(c) Scarfy sprite by Eiden Marsal", screenWidth-200, screenHeight-20, 10, raylib.Gray)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(scarfy)

	raylib.CloseWindow()
}
