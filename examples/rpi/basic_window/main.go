package main

import "github.com/gen2brain/raylib-go/raylib"

func main() {
	raylib.InitWindow(800, 450, "raylib [rpi] example - basic window")

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("Congrats! You created your first window!", 190, 200, 20, raylib.LightGray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
