package main

import (
	"fmt"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var exit_key = false

func main() {
	rl.InitWindow(800, 450, "raygui - button")

	rl.SetTargetFPS(60)

	var button bool

	for !(rl.WindowShouldClose() || exit_key) {
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
