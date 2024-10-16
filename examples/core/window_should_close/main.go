package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	screenWidth  = 850
	screenHeight = 480
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - window should close")

	rl.SetExitKey(rl.KeyNull) // Disable KEY_ESCAPE to close window, X-button still works

	exitWindowRequested := false // Flag to request window to exit
	exitWindow := false          // Flag to set window to exit

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	for !exitWindow {
		if rl.WindowShouldClose() || rl.IsMouseButtonPressed(rl.KeyEscape) {
			exitWindowRequested = true
		}

		if exitWindowRequested {
			if rl.IsKeyPressed(rl.KeyY) {
				exitWindow = true
			} else if rl.IsKeyPressed(rl.KeyN) {
				exitWindowRequested = false
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		if exitWindowRequested {
			rl.DrawRectangle(0, 100, screenWidth, 200, rl.Black)
			rl.DrawText("Are you sure you want to exit the program? [Y/N]", 40, 180, 30, rl.White)
		} else {
			rl.DrawText("Try to close the window to get confirmation message!", 120, 200, 20, rl.LightGray)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow() // Close window and OpenGL context
}
