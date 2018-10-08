package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib [shapes] example - raylib color palette")
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("raylib color palette", 28, 42, 20, rl.Black)

		rl.DrawRectangle(26, 80, 100, 100, rl.DarkGray)
		rl.DrawRectangle(26, 188, 100, 100, rl.Gray)
		rl.DrawRectangle(26, 296, 100, 100, rl.LightGray)
		rl.DrawRectangle(134, 80, 100, 100, rl.Maroon)
		rl.DrawRectangle(134, 188, 100, 100, rl.Red)
		rl.DrawRectangle(134, 296, 100, 100, rl.Pink)
		rl.DrawRectangle(242, 80, 100, 100, rl.Orange)
		rl.DrawRectangle(242, 188, 100, 100, rl.Gold)
		rl.DrawRectangle(242, 296, 100, 100, rl.Yellow)
		rl.DrawRectangle(350, 80, 100, 100, rl.DarkGreen)
		rl.DrawRectangle(350, 188, 100, 100, rl.Lime)
		rl.DrawRectangle(350, 296, 100, 100, rl.Green)
		rl.DrawRectangle(458, 80, 100, 100, rl.DarkBlue)
		rl.DrawRectangle(458, 188, 100, 100, rl.Blue)
		rl.DrawRectangle(458, 296, 100, 100, rl.SkyBlue)
		rl.DrawRectangle(566, 80, 100, 100, rl.DarkPurple)
		rl.DrawRectangle(566, 188, 100, 100, rl.Violet)
		rl.DrawRectangle(566, 296, 100, 100, rl.Purple)
		rl.DrawRectangle(674, 80, 100, 100, rl.DarkBrown)
		rl.DrawRectangle(674, 188, 100, 100, rl.Brown)
		rl.DrawRectangle(674, 296, 100, 100, rl.Beige)

		rl.DrawText("DARKGRAY", 65, 166, 10, rl.Black)
		rl.DrawText("GRAY", 93, 274, 10, rl.Black)
		rl.DrawText("LIGHTGRAY", 61, 382, 10, rl.Black)
		rl.DrawText("MAROON", 186, 166, 10, rl.Black)
		rl.DrawText("RED", 208, 274, 10, rl.Black)
		rl.DrawText("PINK", 204, 382, 10, rl.Black)
		rl.DrawText("ORANGE", 295, 166, 10, rl.Black)
		rl.DrawText("GOLD", 310, 274, 10, rl.Black)
		rl.DrawText("YELLOW", 300, 382, 10, rl.Black)
		rl.DrawText("DARKGREEN", 382, 166, 10, rl.Black)
		rl.DrawText("LIME", 420, 274, 10, rl.Black)
		rl.DrawText("GREEN", 410, 382, 10, rl.Black)
		rl.DrawText("DARKBLUE", 498, 166, 10, rl.Black)
		rl.DrawText("BLUE", 526, 274, 10, rl.Black)
		rl.DrawText("SKYBLUE", 505, 382, 10, rl.Black)
		rl.DrawText("DARKPURPLE", 592, 166, 10, rl.Black)
		rl.DrawText("VIOLET", 621, 274, 10, rl.Black)
		rl.DrawText("PURPLE", 620, 382, 10, rl.Black)
		rl.DrawText("DARKBROWN", 705, 166, 10, rl.Black)
		rl.DrawText("BROWN", 733, 274, 10, rl.Black)
		rl.DrawText("BEIGE", 737, 382, 10, rl.Black)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
