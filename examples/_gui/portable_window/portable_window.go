package main

import (
	"fmt"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

/*******************************************************************************************
*
*   raygui - portable window
*
*   DEPENDENCIES:
*       raylib 4.0  - Windowing/input management and drawing.
*       raygui 3.0  - Immediate-mode GUI controls.
*
*   COMPILATION (Windows - MinGW):
*       gcc -o $(NAME_PART).exe $(FILE_NAME) -I../../src -lraylib -lopengl32 -lgdi32 -std=c99
*
*   LICENSE: zlib/libpng
*
*   Copyright (c) 2016-2022 Ramon Santamaria (@raysan5)
*
**********************************************************************************************/

// ------------------------------------------------------------------------------------
// Program main entry point
// ------------------------------------------------------------------------------------
func main() {
	// Initialization
	//---------------------------------------------------------------------------------------
	const (
		screenWidth  = 800
		screenHeight = 600
	)

	rl.SetConfigFlags(rl.FlagWindowUndecorated)
	rl.InitWindow(screenWidth, screenHeight, "raygui - portable window")

	// General variables
	var (
		mousePosition  = rl.Vector2{0, 0}
		windowPosition = rl.Vector2{500, 200}
		panOffset      = mousePosition
		dragWindow     = false
	)

	rl.SetWindowPosition(int(windowPosition.X), int(windowPosition.Y))

	exitWindow := false

	rl.SetTargetFPS(60)
	//--------------------------------------------------------------------------------------

	// Main game loop
	for !exitWindow && !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		//----------------------------------------------------------------------------------
		mousePosition = rl.GetMousePosition()

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if (rl.CheckCollisionPointRec(mousePosition, rl.Rectangle{0, 0, screenWidth, 20})) {
				dragWindow = true
				panOffset = mousePosition
			}
		}

		if dragWindow {
			windowPosition.X += (mousePosition.X - panOffset.X)
			windowPosition.Y += (mousePosition.Y - panOffset.Y)

			if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				dragWindow = false
			}

			rl.SetWindowPosition(int(windowPosition.X), int(windowPosition.Y))
		}
		//----------------------------------------------------------------------------------

		// Draw
		//----------------------------------------------------------------------------------
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		exitWindow = gui.WindowBox(rl.Rectangle{0, 0, screenWidth, screenHeight}, "#198# PORTABLE WINDOW")

		rl.DrawText(fmt.Sprintf("Mouse Position: [ %.0f, %.0f ]", mousePosition.X, mousePosition.Y), 10, 40, 10, rl.DarkGray)

		rl.EndDrawing()
		//----------------------------------------------------------------------------------
	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	rl.CloseWindow() // Close window and OpenGL context
	//--------------------------------------------------------------------------------------
}
