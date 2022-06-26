package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib [raygui] example - basic controls")
	rl.GuiLoadStyle("./styles/cherry.txt.rgs")
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.GetColor(uint(rl.GuiGetStyle(0, 19))))

		rl.GuiLabel(rl.NewRectangle(30, 30, 600, 10), "HEY HEY HEY HEY!!!!")

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
