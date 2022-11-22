package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - generate random values")

	framesCounter := 0                    // Variable used to count frames
	randValue := rl.GetRandomValue(-8, 5) // Get a random integer number between -8 and 5 (both included)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		framesCounter++

		// Every two seconds (120 frames) a new random value is generated
		if ((framesCounter / 120) % 2) == 1 {
			randValue = rl.GetRandomValue(-8, 5)
			framesCounter = 0
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("Every 2 seconds a new random value is generated:", 130, 100, 20, rl.Maroon)

		rl.DrawText(fmt.Sprintf("%d", randValue), 360, 180, 80, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
