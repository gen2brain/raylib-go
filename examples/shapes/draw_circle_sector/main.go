package main

import (
	"math"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - draw circle sector")

	cntr := rl.NewVector2(float32(rl.GetScreenWidth()-300)/2, float32(rl.GetScreenHeight())/2)

	outerRad, startAngle, endAngle, segments := float32(180), float32(0), float32(180), float32(10)
	minSegments := 4

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawLine(500, 0, 500, int32(rl.GetScreenHeight()), rl.Fade(rl.LightGray, 0.5))
		rl.DrawRectangle(500, 0, int32(rl.GetScreenWidth()-500), int32(rl.GetScreenHeight()), rl.Fade(rl.LightGray, 0.3))

		rl.DrawCircleSector(cntr, outerRad, startAngle, endAngle, int32(segments), rl.Fade(rl.Maroon, 0.5))
		rl.DrawCircleSectorLines(cntr, outerRad, startAngle, endAngle, int32(segments), rl.Fade(rl.Maroon, 0.8))

		startAngle = rg.SliderBar(rl.NewRectangle(600, 40, 120, 20), "Start Angle", "", startAngle, 0, 720)
		endAngle = rg.SliderBar(rl.NewRectangle(600, 70, 120, 20), "End Angle", "", endAngle, 0, 720)
		outerRad = rg.SliderBar(rl.NewRectangle(600, 140, 120, 20), "Radius", "", outerRad, 0, 200)
		segments = rg.SliderBar(rl.NewRectangle(600, 170, 120, 20), "Segments", "", segments, 0, 100)

		minSegments = int(math.Ceil(float64(endAngle-startAngle) / 90))

		if segments >= float32(minSegments) {
			rl.DrawText("MODE: MANUAL", 600, 200, 10, rl.Maroon)
		} else {
			rl.DrawText("MODE: AUTO", 600, 200, 10, rl.DarkGray)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
