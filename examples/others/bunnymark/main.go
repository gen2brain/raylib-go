package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

// Bunny type
type Bunny struct {
	Position rl.Vector2
	Speed    rl.Vector2
	Color    rl.Color
}

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(960)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - Bunnymark")

	texture := rl.LoadTexture("wabbit_alpha.png")

	bunnies := make([]*Bunny, 0)
	bunniesCount := 0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			// Create more bunnies
			for i := 0; i < 100; i++ {
				b := &Bunny{}
				b.Position = rl.GetMousePosition()
				b.Speed.X = float32(rl.GetRandomValue(250, 500)) / 60.0
				b.Speed.Y = float32(rl.GetRandomValue(250, 500)-500) / 60.0

				bunnies = append(bunnies, b)
				bunniesCount++
			}
		}

		// Update bunnies
		for _, b := range bunnies {
			b.Position.X += b.Speed.X
			b.Position.Y += b.Speed.Y

			if (b.Position.X > float32(screenWidth)) || (b.Position.X < 0) {
				b.Speed.X *= -1
			}

			if (b.Position.Y > float32(screenHeight)) || (b.Position.Y < 0) {
				b.Speed.Y *= -1
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		for _, b := range bunnies {
			// NOTE: When internal QUADS batch limit is reached, a draw call is launched and
			// batching buffer starts being filled again; before launching the draw call,
			// updated vertex data from internal buffer is send to GPU... it seems it generates
			// a stall and consequently a frame drop, limiting number of bunnies drawn at 60 fps
			rl.DrawTexture(texture, int32(b.Position.X), int32(b.Position.Y), rl.RayWhite)
		}

		rl.DrawRectangle(0, 0, screenWidth, 40, rl.LightGray)
		rl.DrawText("raylib bunnymark", 10, 10, 20, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("bunnies: %d", bunniesCount), 400, 10, 20, rl.Red)

		rl.DrawFPS(260, 10)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)
	rl.CloseWindow()
}
