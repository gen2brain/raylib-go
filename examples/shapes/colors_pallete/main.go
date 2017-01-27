package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	raylib.InitWindow(800, 450, "raylib [shapes] example - raylib color palette")
	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("raylib color palette", 28, 42, 20, raylib.Black)

		raylib.DrawRectangle(26, 80, 100, 100, raylib.DarkGray)
		raylib.DrawRectangle(26, 188, 100, 100, raylib.Gray)
		raylib.DrawRectangle(26, 296, 100, 100, raylib.LightGray)
		raylib.DrawRectangle(134, 80, 100, 100, raylib.Maroon)
		raylib.DrawRectangle(134, 188, 100, 100, raylib.Red)
		raylib.DrawRectangle(134, 296, 100, 100, raylib.Pink)
		raylib.DrawRectangle(242, 80, 100, 100, raylib.Orange)
		raylib.DrawRectangle(242, 188, 100, 100, raylib.Gold)
		raylib.DrawRectangle(242, 296, 100, 100, raylib.Yellow)
		raylib.DrawRectangle(350, 80, 100, 100, raylib.DarkGreen)
		raylib.DrawRectangle(350, 188, 100, 100, raylib.Lime)
		raylib.DrawRectangle(350, 296, 100, 100, raylib.Green)
		raylib.DrawRectangle(458, 80, 100, 100, raylib.DarkBlue)
		raylib.DrawRectangle(458, 188, 100, 100, raylib.Blue)
		raylib.DrawRectangle(458, 296, 100, 100, raylib.SkyBlue)
		raylib.DrawRectangle(566, 80, 100, 100, raylib.DarkPurple)
		raylib.DrawRectangle(566, 188, 100, 100, raylib.Violet)
		raylib.DrawRectangle(566, 296, 100, 100, raylib.Purple)
		raylib.DrawRectangle(674, 80, 100, 100, raylib.DarkBrown)
		raylib.DrawRectangle(674, 188, 100, 100, raylib.Brown)
		raylib.DrawRectangle(674, 296, 100, 100, raylib.Beige)

		raylib.DrawText("DARKGRAY", 65, 166, 10, raylib.Black)
		raylib.DrawText("GRAY", 93, 274, 10, raylib.Black)
		raylib.DrawText("LIGHTGRAY", 61, 382, 10, raylib.Black)
		raylib.DrawText("MAROON", 186, 166, 10, raylib.Black)
		raylib.DrawText("RED", 208, 274, 10, raylib.Black)
		raylib.DrawText("PINK", 204, 382, 10, raylib.Black)
		raylib.DrawText("ORANGE", 295, 166, 10, raylib.Black)
		raylib.DrawText("GOLD", 310, 274, 10, raylib.Black)
		raylib.DrawText("YELLOW", 300, 382, 10, raylib.Black)
		raylib.DrawText("DARKGREEN", 382, 166, 10, raylib.Black)
		raylib.DrawText("LIME", 420, 274, 10, raylib.Black)
		raylib.DrawText("GREEN", 410, 382, 10, raylib.Black)
		raylib.DrawText("DARKBLUE", 498, 166, 10, raylib.Black)
		raylib.DrawText("BLUE", 526, 274, 10, raylib.Black)
		raylib.DrawText("SKYBLUE", 505, 382, 10, raylib.Black)
		raylib.DrawText("DARKPURPLE", 592, 166, 10, raylib.Black)
		raylib.DrawText("VIOLET", 621, 274, 10, raylib.Black)
		raylib.DrawText("PURPLE", 620, 382, 10, raylib.Black)
		raylib.DrawText("DARKBROWN", 705, 166, 10, raylib.Black)
		raylib.DrawText("BROWN", 733, 274, 10, raylib.Black)
		raylib.DrawText("BEIGE", 737, 382, 10, raylib.Black)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
