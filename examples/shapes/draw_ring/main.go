package main

import (
	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - draw ring")

	cntr := rl.NewVector2(float32(rl.GetScreenWidth()-300)/2, float32(rl.GetScreenHeight())/2)
	innerRad, outerRad := float32(80), float32(190)
	startAngle, endAngle, segments := float32(0), float32(360), float32(0)
	drawRing, drawRingLines, drawCircLines := true, false, false

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawLine(500, 0, 500, int32(rl.GetScreenHeight()), rl.Fade(rl.LightGray, 0.3))

		if drawRing {
			rl.DrawRing(cntr, innerRad, outerRad, startAngle, endAngle, int32(segments), rl.Fade(rl.Maroon, 0.7))
		}
		if drawRingLines {
			rl.DrawRingLines(cntr, innerRad, outerRad, startAngle, endAngle, int32(segments), rl.Fade(rl.Black, 0.7))
		}
		if drawCircLines {
			rl.DrawCircleSectorLines(cntr, outerRad, startAngle, endAngle, int32(segments), rl.Fade(rl.Black, 0.7))
		}

		startAngle = rg.SliderBar(rl.NewRectangle(600, 40, 120, 20), "Start Angle", "", startAngle, -450, 450)
		endAngle = rg.SliderBar(rl.NewRectangle(600, 70, 120, 20), "End Angle", "", endAngle, -450, 450)
		innerRad = rg.SliderBar(rl.NewRectangle(600, 140, 120, 20), "Inner Radius", "", innerRad, -450, 450)
		outerRad = rg.SliderBar(rl.NewRectangle(600, 170, 120, 20), "Outer Radius", "", outerRad, -450, 450)
		drawRing = rg.CheckBox(rl.NewRectangle(600, 320, 20, 20), "Draw Ring", drawRing)
		drawRingLines = rg.CheckBox(rl.NewRectangle(600, 350, 20, 20), "Draw Ring Lines", drawRingLines)
		drawCircLines = rg.CheckBox(rl.NewRectangle(600, 380, 20, 20), "Draw Ring", drawCircLines)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
