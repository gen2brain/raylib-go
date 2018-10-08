package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const maxCircles = 64

type circleWave struct {
	Position rl.Vector2
	Radius   float32
	Alpha    float32
	Speed    float32
	Color    rl.Color
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagMsaa4xHint) // NOTE: Try to enable MSAA 4X

	rl.InitWindow(screenWidth, screenHeight, "raylib [audio] example - module playing (streaming)")

	rl.InitAudioDevice()

	colors := []rl.Color{
		rl.Orange, rl.Red, rl.Gold, rl.Lime, rl.Blue,
		rl.Violet, rl.Brown, rl.LightGray, rl.Pink,
		rl.Yellow, rl.Green, rl.SkyBlue, rl.Purple, rl.Beige,
	}

	circles := make([]circleWave, maxCircles)

	for i := maxCircles - 1; i >= 0; i-- {
		c := circleWave{}

		c.Alpha = 0
		c.Radius = float32(rl.GetRandomValue(10, 40))

		x := rl.GetRandomValue(int32(c.Radius), screenWidth-int32(c.Radius))
		y := rl.GetRandomValue(int32(c.Radius), screenHeight-int32(c.Radius))
		c.Position = rl.NewVector2(float32(x), float32(y))

		c.Speed = float32(rl.GetRandomValue(1, 100)) / 20000.0
		c.Color = colors[rl.GetRandomValue(0, int32(len(colors)-1))]

		circles[i] = c
	}

	xm := rl.LoadMusicStream("mini1111.xm")
	rl.PlayMusicStream(xm)

	pause := false

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(xm) // Update music buffer with new stream data

		// Restart music playing (stop and play)
		if rl.IsKeyPressed(rl.KeySpace) {
			rl.StopMusicStream(xm)
			rl.PlayMusicStream(xm)
		}

		// Pause/Resume music playing
		if rl.IsKeyPressed(rl.KeyP) {
			pause = !pause

			if pause {
				rl.PauseMusicStream(xm)
			} else {
				rl.ResumeMusicStream(xm)
			}
		}

		// Get timePlayed scaled to bar dimensions
		timePlayed := int32(rl.GetMusicTimePlayed(xm)/rl.GetMusicTimeLength(xm)*float32(screenWidth-40)) * 2

		// Color circles animation
		for i := maxCircles - 1; (i >= 0) && !pause; i-- {
			circles[i].Alpha += circles[i].Speed
			circles[i].Radius += circles[i].Speed * 10.0

			if circles[i].Alpha > 1.0 {
				circles[i].Speed *= -1
			}

			if circles[i].Alpha <= 0.0 {
				circles[i].Alpha = 0.0
				circles[i].Radius = float32(rl.GetRandomValue(10, 40))
				circles[i].Position.X = float32(rl.GetRandomValue(int32(circles[i].Radius), screenWidth-int32(circles[i].Radius)))
				circles[i].Position.Y = float32(rl.GetRandomValue(int32(circles[i].Radius), screenHeight-int32(circles[i].Radius)))
				circles[i].Color = colors[rl.GetRandomValue(0, 13)]
				circles[i].Speed = float32(rl.GetRandomValue(1, 100)) / 20000.0
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		for i := maxCircles - 1; i >= 0; i-- {
			rl.DrawCircleV(circles[i].Position, float32(circles[i].Radius), rl.Fade(circles[i].Color, circles[i].Alpha))
		}

		// Draw time bar
		rl.DrawRectangle(20, screenHeight-20-12, screenWidth-40, 12, rl.LightGray)
		rl.DrawRectangle(20, screenHeight-20-12, timePlayed, 12, rl.Maroon)
		rl.DrawRectangleLines(20, screenHeight-20-12, screenWidth-40, 12, rl.Gray)

		rl.EndDrawing()
	}

	rl.UnloadMusicStream(xm)

	rl.CloseAudioDevice()

	rl.CloseWindow()
}
