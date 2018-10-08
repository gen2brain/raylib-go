package main

import (
	"github.com/gen2brain/raylib-go/easings"
	"github.com/gen2brain/raylib-go/raygui"
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "raylib [easings] example - easings")

	currentTime := 0
	duration := float32(60)
	startPositionX := float32(screenWidth) / 4
	finalPositionX := startPositionX * 3
	currentPositionX := startPositionX

	ballPosition := rl.NewVector2(startPositionX, float32(screenHeight)/2)

	comboActive := 0
	comboLastActive := 0

	easingTypes := []string{"SineIn", "SineOut", "SineInOut", "BounceIn", "BounceOut", "BounceInOut", "BackIn", "BackOut", "BackInOut"}
	ease := easingTypes[comboActive]

	//rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if rl.IsKeyDown(rl.KeyR) {
			currentTime = 0
			currentPositionX = startPositionX
			ballPosition.X = currentPositionX
		}

		if comboLastActive != comboActive {
			currentTime = 0
			currentPositionX = startPositionX
			ballPosition.X = currentPositionX

			ease = easingTypes[comboActive]
			comboLastActive = comboActive
		}

		if currentPositionX < finalPositionX {
			switch ease {
			case "SineIn":
				currentPositionX = easings.SineIn(float32(currentTime), startPositionX, finalPositionX-startPositionX, duration)
			case "SineOut":
				currentPositionX = easings.SineOut(float32(currentTime), startPositionX, finalPositionX-startPositionX, duration)
			case "SineInOut":
				currentPositionX = easings.SineInOut(float32(currentTime), startPositionX, finalPositionX-startPositionX, duration)
			case "BounceIn":
				currentPositionX = easings.BounceIn(float32(currentTime), startPositionX, finalPositionX-startPositionX, duration)
			case "BounceOut":
				currentPositionX = easings.BounceOut(float32(currentTime), startPositionX, finalPositionX-startPositionX, duration)
			case "BounceInOut":
				currentPositionX = easings.BounceInOut(float32(currentTime), startPositionX, finalPositionX-startPositionX, duration)
			case "BackIn":
				currentPositionX = easings.BackIn(float32(currentTime), startPositionX, finalPositionX-startPositionX, duration)
			case "BackOut":
				currentPositionX = easings.BackOut(float32(currentTime), startPositionX, finalPositionX-startPositionX, duration)
			case "BackInOut":
				currentPositionX = easings.BackInOut(float32(currentTime), startPositionX, finalPositionX-startPositionX, duration)
			}

			ballPosition.X = currentPositionX
			currentTime++
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		raygui.Label(rl.NewRectangle(20, 20, 200, 20), "Easing Type:")
		comboActive = raygui.ComboBox(rl.NewRectangle(20, 40, 200, 20), easingTypes, comboActive)

		raygui.Label(rl.NewRectangle(20, 80, 200, 20), "Press R to reset")

		rl.DrawCircleV(ballPosition, 50, rl.Maroon)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
