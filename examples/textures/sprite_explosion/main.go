package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	numFramesPerLine, numLines = 5, 5
)

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - sprite explosion")

	rl.InitAudioDevice()
	fxBoom := rl.LoadSound("boom.wav")
	explosion := rl.LoadTexture("explosion.png")

	frameW := float32(explosion.Width / int32(numFramesPerLine))
	frameH := float32(explosion.Height / int32(numLines))

	currentFrame, currentLine := 0, 0

	frameRec := rl.NewRectangle(0, 0, frameW, frameH)
	position := rl.NewVector2(0, 0)

	active := false
	framesCount := 0

	rl.SetTargetFPS(120)

	for !rl.WindowShouldClose() {

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) && !active {
			position = rl.GetMousePosition()
			active = true
			position.X -= frameW / 2
			position.Y -= frameH / 2
			rl.PlaySound(fxBoom)
		}

		if active {
			framesCount++
			if framesCount > 2 {
				currentFrame++
				if currentFrame >= numFramesPerLine {
					currentFrame = 0
					currentLine++
					if currentLine >= numLines {
						currentLine = 0
						active = false
					}
				}

				framesCount = 0
			}
		}

		frameRec.X = frameW * float32(currentFrame)
		frameRec.Y = frameH * float32(currentLine)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("click left mouse on screen to explode", 19, 21, 20, rl.Red)
		rl.DrawText("click left mouse on screen to explode", 20, 20, 20, rl.Black)

		if active {
			rl.DrawTextureRec(explosion, frameRec, position, rl.White)
		}

		rl.EndDrawing()
	}

	rl.UnloadTexture(explosion)
	rl.UnloadSound(fxBoom)
	rl.CloseAudioDevice()

	rl.CloseWindow()
}
