/*******************************************************************************************
*
*   raylib [core] example - window flags
*
*   Example originally created with raylib 3.5, last time updated with raylib 3.5
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2020-2024 Ramon Santamaria (@raysan5)
*
********************************************************************************************/

package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

// Possible window flags
/*
   FLAG_VSYNC_HINT
   FLAG_FULLSCREEN_MODE    -> not working properly -> wrong scaling!
   FLAG_WINDOW_RESIZABLE
   FLAG_WINDOW_UNDECORATED
   FLAG_WINDOW_TRANSPARENT
   FLAG_WINDOW_HIDDEN
   FLAG_WINDOW_MINIMIZED   -> Not supported on window creation
   FLAG_WINDOW_MAXIMIZED   -> Not supported on window creation
   FLAG_WINDOW_UNFOCUSED
   FLAG_WINDOW_TOPMOST
   FLAG_WINDOW_HIGHDPI     -> errors after minimize-resize, fb size is recalculated
   FLAG_WINDOW_ALWAYS_RUN
   FLAG_MSAA_4X_HINT
*/

func main() {
	// Set configuration flags for window creation
	//SetConfigFlags(FLAG_VSYNC_HINT | FLAG_MSAA_4X_HINT | FLAG_WINDOW_HIGHDPI);
	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - Window Flags")

	ballPosition := rl.Vector2{X: float32(rl.GetScreenWidth()) / 2.0, Y: float32(rl.GetScreenHeight()) / 2.0}
	ballSpeed := rl.Vector2{X: 5.0, Y: 4.0}
	ballRadius := float32(20.0)
	framesCounter := 0

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		if rl.IsKeyPressed(rl.KeyF) {
			rl.ToggleFullscreen() // modifies window size when scaling!
		}

		if rl.IsKeyPressed(rl.KeyR) {
			if rl.IsWindowState(rl.FlagWindowResizable) {
				rl.ClearWindowState(rl.FlagWindowResizable)
			} else {
				rl.SetWindowState(rl.FlagWindowResizable)
			}
		}

		if rl.IsKeyPressed(rl.KeyD) {
			if rl.IsWindowState(rl.FlagWindowUndecorated) {
				rl.ClearWindowState(rl.FlagWindowUndecorated)
			} else {
				rl.SetWindowState(rl.FlagWindowUndecorated)
			}
		}

		if rl.IsKeyPressed(rl.KeyH) {
			if !rl.IsWindowState(rl.FlagWindowHidden) {
				rl.SetWindowState(rl.FlagWindowHidden)
			}

			framesCounter = 0
		}

		if rl.IsWindowState(rl.FlagWindowHidden) {
			framesCounter++
			if framesCounter >= 240 {
				rl.ClearWindowState(rl.FlagWindowHidden) // Show window after 3 seconds
			}
		}

		if rl.IsKeyPressed(rl.KeyN) {
			if !rl.IsWindowState(rl.FlagWindowMinimized) {
				rl.MinimizeWindow()
			}

			framesCounter = 0
		}

		if rl.IsWindowState(rl.FlagWindowMinimized) {
			framesCounter++
			if framesCounter >= 240 {
				rl.RestoreWindow() // Restore window after 3 seconds
			}
		}

		if rl.IsKeyPressed(rl.KeyM) {
			// NOTE: Requires FLAG_WINDOW_RESIZABLE enabled!
			if rl.IsWindowState(rl.FlagWindowMaximized) {
				rl.RestoreWindow()
			} else {
				rl.MaximizeWindow()
			}
		}

		if rl.IsKeyPressed(rl.KeyU) {
			if rl.IsWindowState(rl.FlagWindowUnfocused) {
				rl.ClearWindowState(rl.FlagWindowUnfocused)
			} else {
				rl.SetWindowState(rl.FlagWindowUnfocused)
			}
		}

		if rl.IsKeyPressed(rl.KeyT) {
			if rl.IsWindowState(rl.FlagWindowTopmost) {
				rl.ClearWindowState(rl.FlagWindowTopmost)
			} else {
				rl.SetWindowState(rl.FlagWindowTopmost)
			}
		}

		if rl.IsKeyPressed(rl.KeyA) {
			if rl.IsWindowState(rl.FlagWindowAlwaysRun) {
				rl.ClearWindowState(rl.FlagWindowAlwaysRun)
			} else {
				rl.SetWindowState(rl.FlagWindowAlwaysRun)
			}
		}

		if rl.IsKeyPressed(rl.KeyV) {
			if rl.IsWindowState(rl.FlagVsyncHint) {
				rl.ClearWindowState(rl.FlagVsyncHint)
			} else {
				rl.SetWindowState(rl.FlagVsyncHint)
			}
		}

		// Bouncing ball logic
		ballPosition.X += ballSpeed.X
		ballPosition.Y += ballSpeed.Y
		if (ballPosition.X >= (float32(rl.GetScreenWidth()) - ballRadius)) || (ballPosition.X <= ballRadius) {
			ballSpeed.X *= -1.0
		}
		if (ballPosition.Y >= (float32(rl.GetScreenHeight()) - ballRadius)) || (ballPosition.Y <= ballRadius) {
			ballSpeed.Y *= -1.0
		}

		rl.BeginDrawing()
		if rl.IsWindowState(rl.FlagWindowTransparent) {
			rl.ClearBackground(rl.Blank)
		} else {
			rl.ClearBackground(rl.White)
		}

		rl.DrawCircleV(ballPosition, ballRadius, rl.Maroon)
		rl.DrawRectangleLinesEx(
			rl.Rectangle{
				Width:  float32(rl.GetScreenWidth()),
				Height: float32(rl.GetScreenHeight()),
			},
			4, rl.White,
		)

		rl.DrawCircleV(rl.GetMousePosition(), 10, rl.DarkBlue)

		rl.DrawFPS(10, 10)

		rl.DrawText(
			fmt.Sprintf("Screen Size: [%d, %d]", rl.GetScreenWidth(), rl.GetScreenHeight()),
			10,
			40,
			10,
			rl.Green,
		)

		// Draw window state info
		rl.DrawText("Following flags can be set after window creation:", 10, 60, 10, rl.Gray)
		if rl.IsWindowState(rl.FlagFullscreenMode) {
			rl.DrawText("[F] FLAG_FULLSCREEN_MODE: on", 10, 80, 10, rl.Lime)
		} else {
			rl.DrawText("[F] FLAG_FULLSCREEN_MODE: off", 10, 80, 10, rl.Maroon)
		}
		if rl.IsWindowState(rl.FlagWindowResizable) {
			rl.DrawText("[R] FLAG_WINDOW_RESIZABLE: on", 10, 100, 10, rl.Lime)
		} else {
			rl.DrawText("[R] FLAG_WINDOW_RESIZABLE: off", 10, 100, 10, rl.Maroon)
		}
		if rl.IsWindowState(rl.FlagWindowUndecorated) {
			rl.DrawText("[D] FLAG_WINDOW_UNDECORATED: on", 10, 120, 10, rl.Lime)
		} else {
			rl.DrawText("[D] FLAG_WINDOW_UNDECORATED: off", 10, 120, 10, rl.Maroon)
		}
		if rl.IsWindowState(rl.FlagWindowHidden) {
			rl.DrawText("[H] FLAG_WINDOW_HIDDEN: on", 10, 140, 10, rl.Lime)
		} else {
			rl.DrawText("[H] FLAG_WINDOW_HIDDEN: off", 10, 140, 10, rl.Maroon)
		}
		if rl.IsWindowState(rl.FlagWindowMinimized) {
			rl.DrawText("[N] FLAG_WINDOW_MINIMIZED: on", 10, 160, 10, rl.Lime)
		} else {
			rl.DrawText("[N] FLAG_WINDOW_MINIMIZED: off", 10, 160, 10, rl.Maroon)
		}
		if rl.IsWindowState(rl.FlagWindowMaximized) {
			rl.DrawText("[M] FLAG_WINDOW_MAXIMIZED: on", 10, 180, 10, rl.Lime)
		} else {
			rl.DrawText("[M] FLAG_WINDOW_MAXIMIZED: off", 10, 180, 10, rl.Maroon)
		}
		if rl.IsWindowState(rl.FlagWindowUnfocused) {
			rl.DrawText("[G] FLAG_WINDOW_UNFOCUSED: on", 10, 200, 10, rl.Lime)
		} else {
			rl.DrawText("[U] FLAG_WINDOW_UNFOCUSED: off", 10, 200, 10, rl.Maroon)
		}
		if rl.IsWindowState(rl.FlagWindowTopmost) {
			rl.DrawText("[T] FLAG_WINDOW_TOPMOST: on", 10, 220, 10, rl.Lime)
		} else {
			rl.DrawText("[T] FLAG_WINDOW_TOPMOST: off", 10, 220, 10, rl.Maroon)
		}
		if rl.IsWindowState(rl.FlagWindowAlwaysRun) {
			rl.DrawText("[A] FLAG_WINDOW_ALWAYS_RUN: on", 10, 240, 10, rl.Lime)
		} else {
			rl.DrawText("[A] FLAG_WINDOW_ALWAYS_RUN: off", 10, 240, 10, rl.Maroon)
		}
		if rl.IsWindowState(rl.FlagVsyncHint) {
			rl.DrawText("[V] FLAG_VSYNC_HINT: on", 10, 260, 10,
				rl.Lime)
		} else {
			rl.DrawText("[V] FLAG_VSYNC_HINT: off", 10, 260, 10, rl.Maroon)
		}

		rl.DrawText("Following flags can only be set before window creation:", 10, 300, 10, rl.Gray)
		if rl.IsWindowState(rl.FlagWindowHighdpi) {
			rl.DrawText("FLAG_WINDOW_HIGHDPI: on", 10, 320, 10, rl.Lime)
		} else {
			rl.DrawText("FLAG_WINDOW_HIGHDPI: off", 10, 320, 10, rl.Maroon)
		}
		if rl.IsWindowState(rl.FlagWindowTransparent) {
			rl.DrawText("FLAG_WINDOW_TRANSPARENT: on", 10, 340, 10, rl.Lime)
		} else {
			rl.DrawText("FLAG_WINDOW_TRANSPARENT: off", 10, 340, 10, rl.Maroon)
		}
		if rl.IsWindowState(rl.FlagMsaa4xHint) {
			rl.DrawText("FLAG_MSAA_4X_HINT: on", 10, 360, 10, rl.Lime)
		} else {
			rl.DrawText("FLAG_MSAA_4X_HINT: off", 10, 360, 10, rl.Maroon)
		}

		rl.EndDrawing()
	}
	rl.CloseWindow()
}
