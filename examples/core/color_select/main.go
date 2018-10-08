package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - color selection (collision detection)")

	colors := [21]rl.Color{
		rl.DarkGray, rl.Maroon, rl.Orange, rl.DarkGreen, rl.DarkBlue, rl.DarkPurple,
		rl.DarkBrown, rl.Gray, rl.Red, rl.Gold, rl.Lime, rl.Blue, rl.Violet, rl.Brown,
		rl.LightGray, rl.Pink, rl.Yellow, rl.Green, rl.SkyBlue, rl.Purple, rl.Beige,
	}

	colorsRecs := make([]rl.Rectangle, 21) // Rectangles array

	// Fills colorsRecs data (for every rectangle)
	for i := 0; i < 21; i++ {
		r := rl.Rectangle{}
		r.X = float32(20 + 100*(i%7) + 10*(i%7))
		r.Y = float32(60 + 100*(i/7) + 10*(i/7))
		r.Width = 100
		r.Height = 100

		colorsRecs[i] = r
	}

	selected := make([]bool, 21) // Selected rectangles indicator

	mousePoint := rl.Vector2{}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		mousePoint = rl.GetMousePosition()

		for i := 0; i < 21; i++ { // Iterate along all the rectangles
			if rl.CheckCollisionPointRec(mousePoint, colorsRecs[i]) {
				colors[i].A = 120

				if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
					selected[i] = !selected[i]
				}
			} else {
				colors[i].A = 255
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		for i := 0; i < 21; i++ { // Draw all rectangles
			rl.DrawRectangleRec(colorsRecs[i], colors[i])

			// Draw four rectangles around selected rectangle
			if selected[i] {
				rl.DrawRectangle(int32(colorsRecs[i].X), int32(colorsRecs[i].Y), 100, 10, rl.RayWhite)    // Square top rectangle
				rl.DrawRectangle(int32(colorsRecs[i].X), int32(colorsRecs[i].Y), 10, 100, rl.RayWhite)    // Square left rectangle
				rl.DrawRectangle(int32(colorsRecs[i].X+90), int32(colorsRecs[i].Y), 10, 100, rl.RayWhite) // Square right rectangle
				rl.DrawRectangle(int32(colorsRecs[i].X), int32(colorsRecs[i].Y)+90, 100, 10, rl.RayWhite) // Square bottom rectangle
			}
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
