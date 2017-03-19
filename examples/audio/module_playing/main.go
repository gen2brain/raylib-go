package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const maxCircles = 64

type circleWave struct {
	Position raylib.Vector2
	Radius   float32
	Alpha    float32
	Speed    float32
	Color    raylib.Color
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.SetConfigFlags(raylib.FlagMsaa4xHint) // NOTE: Try to enable MSAA 4X

	raylib.InitWindow(screenWidth, screenHeight, "raylib [audio] example - module playing (streaming)")

	raylib.InitAudioDevice()

	colors := []raylib.Color{
		raylib.Orange, raylib.Red, raylib.Gold, raylib.Lime, raylib.Blue,
		raylib.Violet, raylib.Brown, raylib.LightGray, raylib.Pink,
		raylib.Yellow, raylib.Green, raylib.SkyBlue, raylib.Purple, raylib.Beige,
	}

	circles := make([]circleWave, maxCircles)

	for i := maxCircles - 1; i >= 0; i-- {
		c := circleWave{}

		c.Alpha = 0
		c.Radius = float32(raylib.GetRandomValue(10, 40))

		x := raylib.GetRandomValue(int32(c.Radius), screenWidth-int32(c.Radius))
		y := raylib.GetRandomValue(int32(c.Radius), screenHeight-int32(c.Radius))
		c.Position = raylib.NewVector2(float32(x), float32(y))

		c.Speed = float32(raylib.GetRandomValue(1, 100)) / 20000.0
		c.Color = colors[raylib.GetRandomValue(0, int32(len(colors)-1))]

		circles[i] = c
	}

	xm := raylib.LoadMusicStream("mini1111.xm")
	raylib.PlayMusicStream(xm)

	pause := false

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.UpdateMusicStream(xm) // Update music buffer with new stream data

		// Restart music playing (stop and play)
		if raylib.IsKeyPressed(raylib.KeySpace) {
			raylib.StopMusicStream(xm)
			raylib.PlayMusicStream(xm)
		}

		// Pause/Resume music playing
		if raylib.IsKeyPressed(raylib.KeyP) {
			pause = !pause

			if pause {
				raylib.PauseMusicStream(xm)
			} else {
				raylib.ResumeMusicStream(xm)
			}
		}

		// Get timePlayed scaled to bar dimensions
		timePlayed := int32(raylib.GetMusicTimePlayed(xm)/raylib.GetMusicTimeLength(xm)*float32(screenWidth-40)) * 2

		// Color circles animation
		for i := maxCircles - 1; (i >= 0) && !pause; i-- {
			circles[i].Alpha += circles[i].Speed
			circles[i].Radius += circles[i].Speed * 10.0

			if circles[i].Alpha > 1.0 {
				circles[i].Speed *= -1
			}

			if circles[i].Alpha <= 0.0 {
				circles[i].Alpha = 0.0
				circles[i].Radius = float32(raylib.GetRandomValue(10, 40))
				circles[i].Position.X = float32(raylib.GetRandomValue(int32(circles[i].Radius), screenWidth-int32(circles[i].Radius)))
				circles[i].Position.Y = float32(raylib.GetRandomValue(int32(circles[i].Radius), screenHeight-int32(circles[i].Radius)))
				circles[i].Color = colors[raylib.GetRandomValue(0, 13)]
				circles[i].Speed = float32(raylib.GetRandomValue(1, 100)) / 20000.0
			}
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		for i := maxCircles - 1; i >= 0; i-- {
			raylib.DrawCircleV(circles[i].Position, float32(circles[i].Radius), raylib.Fade(circles[i].Color, circles[i].Alpha))
		}

		// Draw time bar
		raylib.DrawRectangle(20, screenHeight-20-12, screenWidth-40, 12, raylib.LightGray)
		raylib.DrawRectangle(20, screenHeight-20-12, timePlayed, 12, raylib.Maroon)
		raylib.DrawRectangleLines(20, screenHeight-20-12, screenWidth-40, 12, raylib.Gray)

		raylib.EndDrawing()
	}

	raylib.UnloadMusicStream(xm)

	raylib.CloseAudioDevice()

	raylib.CloseWindow()
}
