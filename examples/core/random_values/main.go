package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	raylib.InitWindow(800, 450, "raylib [core] example - generate random values")

	framesCounter := 0                        // Variable used to count frames
	randValue := raylib.GetRandomValue(-8, 5) // Get a random integer number between -8 and 5 (both included)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		framesCounter++

		// Every two seconds (120 frames) a new random value is generated
		if ((framesCounter / 120) % 2) == 1 {
			randValue = raylib.GetRandomValue(-8, 5)
			framesCounter = 0
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("Every 2 seconds a new random value is generated:", 130, 100, 20, raylib.Maroon)

		raylib.DrawText(fmt.Sprintf("%d", randValue), 360, 180, 80, raylib.LightGray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
