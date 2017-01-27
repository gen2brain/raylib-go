package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [text] example - text writing anim")

	message := "This sample illustrates a text writing\nanimation effect! Check it out! ;)"
	length := len(message)

	framesCounter := 0

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update
		if raylib.IsKeyDown(raylib.KeySpace) {
			framesCounter += 8
		} else {
			framesCounter++
		}

		if raylib.IsKeyPressed(raylib.KeyEnter) {
			framesCounter = 0
		}

		if framesCounter/10 > length {
			framesCounter = length * 10
		}

		// Draw
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText(message[0:framesCounter/10], 210, 160, 20, raylib.Maroon)

		raylib.DrawText("PRESS [ENTER] to RESTART!", 240, 260, 20, raylib.LightGray)
		raylib.DrawText("PRESS [SPACE] to SPEED UP!", 239, 300, 20, raylib.LightGray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
