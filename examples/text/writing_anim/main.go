package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - text writing anim")

	message := "This sample illustrates a text writing\nanimation effect! Check it out! ;)"
	length := len(message)

	framesCounter := 0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update
		if rl.IsKeyDown(rl.KeySpace) {
			framesCounter += 8
		} else {
			framesCounter++
		}

		if rl.IsKeyPressed(rl.KeyEnter) {
			framesCounter = 0
		}

		if framesCounter/10 > length {
			framesCounter = length * 10
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText(message[0:framesCounter/10], 210, 160, 20, rl.Maroon)

		rl.DrawText("PRESS [ENTER] to RESTART!", 240, 260, 20, rl.LightGray)
		rl.DrawText("PRESS [SPACE] to SPEED UP!", 239, 300, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
