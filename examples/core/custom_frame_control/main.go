package main

/*******************************************************************************************
*
*   raylib [core] example - custom frame control
*
*   NOTE: WARNING: This is an example for advanced users willing to have full control over
*   the frame processes. By default, EndDrawing() calls the following processes:
*       1. Draw remaining batch data: rlDrawRenderBatchActive()
*       2. SwapScreenBuffer()
*       3. Frame time control: WaitTime()
*       4. PollInputEvents()
*
*   To avoid steps 2, 3 and 4, flag SUPPORT_CUSTOM_FRAME_CONTROL can be enabled in
*   config.h (it requires recompiling raylib). This way those steps are up to the user.
*
*   Note that enabling this flag invalidates some functions:
*       - GetFrameTime()
*       - SetTargetFPS()
*       - GetFPS()
*
*   Example originally created with raylib 4.0, last time updated with raylib 4.0
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2021-2024 Ramon Santamaria (@raysan5)
*
********************************************************************************************/

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - Custom Frame Control")

	previousTime := rl.GetTime()
	currentTime := 0.0
	updateDrawTime := 0.0
	waitTime := 0.0
	deltaTime := 0.0

	timeCounter := 0.0
	position := 0.0
	pause := false

	targetFPS := 60

	for !rl.WindowShouldClose() {
		rl.PollInputEvents() // Poll input events (SUPPORT_CUSTOM_FRAME_CONTROL)

		if rl.IsKeyPressed(rl.KeySpace) {
			pause = !pause
		}
		if rl.IsKeyPressed(rl.KeyUp) {
			targetFPS += 20
		} else if rl.IsKeyPressed(rl.KeyDown) {
			targetFPS -= 20
		}

		if targetFPS < 0 {
			targetFPS = 0
		}

		if !pause {
			position += 200 * deltaTime // We move at 200 pixels per second
			if position >= float64(rl.GetScreenWidth()) {
				position = 0
			}
			timeCounter += deltaTime // We count time (seconds)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		for i := 0; i < rl.GetScreenWidth()/200; i++ {
			rl.DrawRectangle(200*int32(i), 0, 1, int32(rl.GetScreenHeight()), rl.SkyBlue)
		}

		rl.DrawCircle(int32(position), int32(rl.GetScreenHeight()/2-25), 50, rl.Red)

		msg := fmt.Sprintf("%03.0f ms", timeCounter*1000)
		rl.DrawText(msg, int32(position-40), int32(rl.GetScreenHeight()/2-100), 20, rl.Maroon)
		msg = fmt.Sprintf("PosX: %03.0f", position)
		rl.DrawText(msg, int32(position-50), int32(rl.GetScreenHeight()/2+40), 20, rl.Black)

		msg = "Circle is moving at a constant 200 pixels/sec,\nindependently of the frame rate."
		rl.DrawText(msg, 10, 10, 20, rl.DarkGray)
		msg = "PRESS SPACE to PAUSE MOVEMENT"
		rl.DrawText(msg, 10, int32(rl.GetScreenHeight()-60), 20, rl.Gray)
		msg = "PRESS UP | DOWN to CHANGE TARGET FPS"
		rl.DrawText(msg, 10, int32(rl.GetScreenHeight()-30), 20, rl.Gray)
		msg = fmt.Sprintf("TARGET FPS: %d", targetFPS)
		rl.DrawText(msg, int32(rl.GetScreenWidth()-220), 10, 20, rl.Lime)
		msg = fmt.Sprintf("CURRENT FPS: %d", int(1/deltaTime))
		rl.DrawText(msg, int32(rl.GetScreenWidth()-220), 40, 20, rl.Lime)

		rl.EndDrawing()

		// NOTE: In case raylib is configured to SUPPORT_CUSTOM_FRAME_CONTROL,
		// Events polling, screen buffer swap and frame time control must be managed by the user

		rl.SwapScreenBuffer() // We want a fixed frame rate

		currentTime = rl.GetTime()
		updateDrawTime = currentTime - previousTime

		if targetFPS > 0 { // We want a fixed frame rate
			waitTime = (1 / float64(targetFPS)) - updateDrawTime
			if waitTime > 0 {
				rl.WaitTime(waitTime)
				currentTime = rl.GetTime()
				deltaTime = currentTime - previousTime
			}
		} else {
			deltaTime = updateDrawTime
		}

		previousTime = currentTime
	}

	rl.CloseWindow()
}
