package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [core] example - keyboard input")

	ballPosition := raylib.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if raylib.IsKeyDown(raylib.KeyRight) {
			ballPosition.X += 0.8
		}
		if raylib.IsKeyDown(raylib.KeyLeft) {
			ballPosition.X -= 0.8
		}
		if raylib.IsKeyDown(raylib.KeyUp) {
			ballPosition.Y -= 0.8
		}
		if raylib.IsKeyDown(raylib.KeyDown) {
			ballPosition.Y += 0.8
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("move the ball with arrow keys", 10, 10, 20, raylib.DarkGray)
		raylib.DrawCircleV(ballPosition, 50, raylib.Maroon)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
