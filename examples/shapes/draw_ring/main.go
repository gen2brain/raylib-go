/*******************************************************************************************
*
*   raylib [shapes] example - draw ring (with gui options)
*
*   Example originally created with raylib 2.5, last time updated with raylib 2.5
*
*   Example contributed by Vlad Adrian (@demizdor) and reviewed by Ramon Santamaria (@raysan5)
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2018-2024 Vlad Adrian (@demizdor) and Ramon Santamaria (@raysan5)
*
********************************************************************************************/
package main

import (
	"fmt"
	"math"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - draw ring")

	center := rl.Vector2{X: (float32(rl.GetScreenWidth()) - 300) / 2.0, Y: float32(rl.GetScreenHeight()) / 2.0}

	var innerRadius, outerRadius float32 = 80.0, 190.0
	var startAngle, endAngle float32 = 0.0, 360
	var segments float32 = 0.0

	drawRing, drawRingLines, drawSectorLines := true, false, false

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		// NOTE: All variables update happens inside GUI control functions

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawLine(500, 0, 500, int32(rl.GetScreenHeight()), rl.Fade(rl.LightGray, 0.6))
		rl.DrawRectangle(500, 0, int32(rl.GetScreenWidth())-500, int32(rl.GetScreenHeight()),
			rl.Fade(rl.LightGray, 0.3))

		if drawRing {
			rl.DrawRing(center, innerRadius, outerRadius, startAngle, endAngle, int32(
				segments), rl.Fade(rl.Maroon, 0.3))
		}
		if drawRingLines {
			rl.DrawRingLines(center, innerRadius, outerRadius, startAngle, endAngle, int32(
				segments), rl.Fade(rl.Black, 0.4))
		}
		if drawSectorLines {
			rl.DrawCircleSectorLines(center, outerRadius, startAngle, endAngle, int32(
				segments), rl.Fade(rl.Black, 0.4))
		}

		// Draw GUI controls
		//------------------------------------------------------------------------------
		startAngle = gui.Slider(rl.Rectangle{
			X: 600, Y: 40, Width: 120, Height: 20,
		}, "StartAngle", fmt.Sprintf("%.2f", startAngle), startAngle, -450, 450)
		endAngle = gui.Slider(rl.Rectangle{
			X: 600, Y: 70, Width: 120, Height: 20,
		}, "EndAngle", fmt.Sprintf("%.2f", endAngle), endAngle, -450, 450)

		innerRadius = gui.Slider(rl.Rectangle{
			X: 600, Y: 140, Width: 120, Height: 20,
		}, "InnerRadius", fmt.Sprintf("%.2f", innerRadius), innerRadius, 0, 100)
		outerRadius = gui.Slider(rl.Rectangle{
			X: 600, Y: 170, Width: 120, Height: 20,
		}, "OuterRadius", fmt.Sprintf("%.2f", outerRadius), outerRadius, 0, 200)

		segments = gui.Slider(rl.Rectangle{
			X: 600, Y: 240, Width: 120, Height: 20,
		}, "Segments", fmt.Sprintf("%.2f", segments), segments, 0, 100)

		drawRing = gui.CheckBox(rl.Rectangle{
			X: 600, Y: 320, Width: 20, Height: 20,
		}, "Draw Ring", drawRing)
		drawRingLines = gui.CheckBox(rl.Rectangle{
			X: 600, Y: 350, Width: 20, Height: 20,
		}, "Draw Ring Lines", drawRingLines)
		drawSectorLines = gui.CheckBox(rl.Rectangle{
			X: 600, Y: 380, Width: 20, Height: 20,
		}, "Draw Sector Lines", drawSectorLines)

		minSegments := ceil((endAngle - startAngle) / 90)
		var color = rl.DarkGray
		if segments >= minSegments {
			color = rl.Maroon
		}
		var mode = "AUTO"
		if segments >= minSegments {
			mode = "MANUAL"
		}
		rl.DrawText(fmt.Sprintf("MODE: %s", mode), 600, 270, 10, color)
		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.CloseWindow() // Close window and OpenGL context
}

func ceil(value float32) float32 {
	return float32(math.Ceil(float64(value)))
}
