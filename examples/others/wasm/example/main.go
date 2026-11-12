package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(800, 450, "raylib [others] example - wasm")

	texture := rl.LoadTexture("resources/raylib_logo.png")

	ball := rl.NewVector2(400, 225)
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		speed := 240 * rl.GetFrameTime()
		if rl.IsKeyDown(rl.KeyRight) {
			ball.X += speed
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			ball.X -= speed
		}
		if rl.IsKeyDown(rl.KeyUp) {
			ball.Y -= speed
		}
		if rl.IsKeyDown(rl.KeyDown) {
			ball.Y += speed
		}
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			ball = rl.GetMousePosition()
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawTexture(texture, 800/2-texture.Width/2, 450/2-texture.Height/2, rl.White)
		rl.DrawCircleV(ball, 30, rl.Maroon)
		rl.DrawText("move the ball with arrow keys or the mouse", 10, 10, 20, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("%d FPS", rl.GetFPS()), 10, 420, 20, rl.Lime)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)
	rl.CloseWindow()
}
