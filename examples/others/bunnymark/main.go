package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// This is the maximum amount of elements (quads) per batch
// NOTE: This value is defined in [rlgl] module and can be changed there
const maxBatchElements = 8192

// Bunny type
type Bunny struct {
	Position rl.Vector2
	Speed    rl.Vector2
	Color    rl.Color
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [others] example - bunnymark")

	texture := rl.LoadTexture("wabbit_alpha.png")

	bunnies := make([]*Bunny, 0)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			// Create more bunnies
			for i := 0; i < 100; i++ {
				b := &Bunny{}
				b.Position = rl.GetMousePosition()
				b.Speed.X = float32(rl.GetRandomValue(-250, 250)) / 60.0
				b.Speed.Y = float32(rl.GetRandomValue(-250, 250)) / 60.0
				b.Color = rl.NewColor(uint8(rl.GetRandomValue(50, 240)), uint8(rl.GetRandomValue(80, 240)), uint8(rl.GetRandomValue(100, 240)), 255)

				bunnies = append(bunnies, b)
			}
		}

		// Update bunnies
		for _, b := range bunnies {
			b.Position.X += b.Speed.X
			b.Position.Y += b.Speed.Y

			if ((b.Position.X + float32(texture.Width/2)) > float32(screenWidth)) || ((b.Position.X + float32(texture.Width/2)) < 0) {
				b.Speed.X *= -1
			}

			if ((b.Position.Y + float32(texture.Height/2)) > float32(screenHeight)) || ((b.Position.Y + float32(texture.Height/2-40)) < 0) {
				b.Speed.Y *= -1
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		for _, b := range bunnies {
			// NOTE: When internal batch buffer limit is reached (MAX_BATCH_ELEMENTS),
			// a draw call is launched and buffer starts being filled again;
			// before issuing a draw call, updated vertex data from internal CPU buffer is send to GPU...
			// Process of sending data is costly and it could happen that GPU data has not been completely
			// processed for drawing while new data is tried to be sent (updating current in-use buffers)
			// it could generates a stall and consequently a frame drop, limiting the number of drawn bunnies
			rl.DrawTexture(texture, int32(b.Position.X), int32(b.Position.Y), b.Color)
		}

		rl.DrawRectangle(0, 0, screenWidth, 40, rl.Black)
		rl.DrawText(fmt.Sprintf("bunnies: %d", len(bunnies)), 120, 10, 20, rl.Green)
		rl.DrawText(fmt.Sprintf("batched draw calls: %d", 1+len(bunnies)/maxBatchElements), 320, 10, 20, rl.Maroon)

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)
	rl.CloseWindow()
}
