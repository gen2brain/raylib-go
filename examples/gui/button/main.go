package main

import (
	"fmt"

	rl  "github.com/Konstantin8105/raylib"
	gui "github.com/Konstantin8105/raygui3.5"
)

func main() {
	rl.InitWindow(800, 450, "raylib [physics] example - box2d")

	rl.SetTargetFPS(60)

	var button bool

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		button = gui.Button(rl.NewRectangle(50, 150, 100, 40), "Click")
		if button {
			fmt.Println("Clicked on button")
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}