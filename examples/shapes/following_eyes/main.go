package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	rl.SetConfigFlags(rl.FlagMsaa4xHint)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - following eyes")

	scleraLpos := rl.NewVector2(float32(screenWidth/2)-100, float32(screenHeight/2))
	scleraRpos := rl.NewVector2(float32(screenWidth/2)+100, float32(screenHeight/2))
	scleraRad := 80

	irisLpos := rl.NewVector2(float32(screenWidth/2)-100, float32(screenHeight/2))
	irisRpos := rl.NewVector2(float32(screenWidth/2)+100, float32(screenHeight/2))
	irisRad := 24

	angle := float32(0)
	dx, dy, dxx, dyy := float32(0), float32(0), float32(0), float32(0)

	rl.SetTargetFPS(60)

	rl.SetMousePosition(0, 0)

	for !rl.WindowShouldClose() {

		irisLpos = rl.GetMousePosition()
		irisRpos = rl.GetMousePosition()

		if !rl.CheckCollisionPointCircle(irisLpos, scleraLpos, float32(scleraRad-20)) {

			dx = irisLpos.X - scleraLpos.X
			dy = irisLpos.Y - scleraLpos.Y

			angle = float32(math.Atan2(float64(dy), float64(dx)))

			dxx = (float32(scleraRad-irisRad) * float32(math.Cos(float64(angle))))
			dyy = (float32(scleraRad-irisRad) * float32(math.Sin(float64(angle))))

			irisLpos.X = scleraLpos.X + dxx
			irisLpos.Y = scleraLpos.Y + dyy

		}

		if !rl.CheckCollisionPointCircle(irisRpos, scleraRpos, float32(scleraRad)-20) {
			dx = irisRpos.X - scleraRpos.X
			dy = irisRpos.Y - scleraRpos.Y

			angle = float32(math.Atan2(float64(dy), float64(dx)))

			dxx = (float32(scleraRad-irisRad) * float32(math.Cos(float64(angle))))
			dyy = (float32(scleraRad-irisRad) * float32(math.Sin(float64(angle))))

			irisRpos.X = scleraRpos.X + dxx
			irisRpos.Y = scleraRpos.Y + dyy

		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawCircleV(scleraLpos, float32(scleraRad), rl.LightGray)
		rl.DrawCircleV(irisLpos, float32(irisRad), rl.Red)
		rl.DrawCircleV(irisLpos, 10, rl.Black)

		rl.DrawCircleV(scleraRpos, float32(scleraRad), rl.LightGray)
		rl.DrawCircleV(irisRpos, float32(irisRad), rl.Orange)
		rl.DrawCircleV(irisRpos, 10, rl.Black)

		rl.DrawFPS(screenWidth-100, 10)

		rl.EndDrawing()

	}

	rl.CloseWindow()
}
