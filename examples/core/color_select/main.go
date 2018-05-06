package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	raylib.InitWindow(800, 450, "raylib [core] example - color selection (collision detection)")

	colors := [21]raylib.Color{
		raylib.DarkGray, raylib.Maroon, raylib.Orange, raylib.DarkGreen, raylib.DarkBlue, raylib.DarkPurple,
		raylib.DarkBrown, raylib.Gray, raylib.Red, raylib.Gold, raylib.Lime, raylib.Blue, raylib.Violet, raylib.Brown,
		raylib.LightGray, raylib.Pink, raylib.Yellow, raylib.Green, raylib.SkyBlue, raylib.Purple, raylib.Beige,
	}

	colorsRecs := make([]raylib.Rectangle, 21) // Rectangles array

	// Fills colorsRecs data (for every rectangle)
	for i := 0; i < 21; i++ {
		r := raylib.Rectangle{}
		r.X = float32(20 + 100*(i%7) + 10*(i%7))
		r.Y = float32(60 + 100*(i/7) + 10*(i/7))
		r.Width = 100
		r.Height = 100

		colorsRecs[i] = r
	}

	selected := make([]bool, 21) // Selected rectangles indicator

	mousePoint := raylib.Vector2{}

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		mousePoint = raylib.GetMousePosition()

		for i := 0; i < 21; i++ { // Iterate along all the rectangles
			if raylib.CheckCollisionPointRec(mousePoint, colorsRecs[i]) {
				colors[i].A = 120

				if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
					selected[i] = !selected[i]
				}
			} else {
				colors[i].A = 255
			}
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		for i := 0; i < 21; i++ { // Draw all rectangles
			raylib.DrawRectangleRec(colorsRecs[i], colors[i])

			// Draw four rectangles around selected rectangle
			if selected[i] {
				raylib.DrawRectangle(int32(colorsRecs[i].X), int32(colorsRecs[i].Y), 100, 10, raylib.RayWhite)    // Square top rectangle
				raylib.DrawRectangle(int32(colorsRecs[i].X), int32(colorsRecs[i].Y), 10, 100, raylib.RayWhite)    // Square left rectangle
				raylib.DrawRectangle(int32(colorsRecs[i].X+90), int32(colorsRecs[i].Y), 10, 100, raylib.RayWhite) // Square right rectangle
				raylib.DrawRectangle(int32(colorsRecs[i].X), int32(colorsRecs[i].Y)+90, 100, 10, raylib.RayWhite) // Square bottom rectangle
			}
		}

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
