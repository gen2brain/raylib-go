package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	xbox360NameID = "Xbox 360 Controller"
	ps3NameID     = "PLAYSTATION(R)3 Controller"
)

func main() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint) // Set MSAA 4X hint before windows creation

	rl.InitWindow(800, 450, "raylib [core] example - gamepad input")

	texPs3Pad := rl.LoadTexture("ps3.png")
	texXboxPad := rl.LoadTexture("xbox.png")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		if rl.IsGamepadAvailable(rl.GamepadPlayer1) {
			rl.DrawText(fmt.Sprintf("GP1: %s", rl.GetGamepadName(rl.GamepadPlayer1)), 10, 10, 10, rl.Black)

			if rl.GetGamepadName(rl.GamepadPlayer1) == xbox360NameID {
				rl.DrawTexture(texXboxPad, 0, 0, rl.DarkGray)

				// Draw buttons: xbox home
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonHome) {
					rl.DrawCircle(394, 89, 19, rl.Red)
				}

				// Draw buttons: basic
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonStart) {
					rl.DrawCircle(436, 150, 9, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonSelect) {
					rl.DrawCircle(352, 150, 9, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonX) {
					rl.DrawCircle(501, 151, 15, rl.Blue)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonA) {
					rl.DrawCircle(536, 187, 15, rl.Lime)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonB) {
					rl.DrawCircle(572, 151, 15, rl.Maroon)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonY) {
					rl.DrawCircle(536, 115, 15, rl.Gold)
				}

				// Draw buttons: d-pad
				rl.DrawRectangle(317, 202, 19, 71, rl.Black)
				rl.DrawRectangle(293, 228, 69, 19, rl.Black)
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonUp) {
					rl.DrawRectangle(317, 202, 19, 26, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonDown) {
					rl.DrawRectangle(317, 202+45, 19, 26, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonLeft) {
					rl.DrawRectangle(292, 228, 25, 19, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonRight) {
					rl.DrawRectangle(292+44, 228, 26, 19, rl.Red)
				}

				// Draw buttons: left-right back
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonLb) {
					rl.DrawCircle(259, 61, 20, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadXboxButtonRb) {
					rl.DrawCircle(536, 61, 20, rl.Red)
				}

				// Draw axis: left joystick
				rl.DrawCircle(259, 152, 39, rl.Black)
				rl.DrawCircle(259, 152, 34, rl.LightGray)
				rl.DrawCircle(int32(259+(rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadXboxAxisLeftX)*20)),
					int32(152-(rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadXboxAxisLeftY)*20)), 25, rl.Black)

				// Draw axis: right joystick
				rl.DrawCircle(461, 237, 38, rl.Black)
				rl.DrawCircle(461, 237, 33, rl.LightGray)
				rl.DrawCircle(int32(461+(rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadXboxAxisRightX)*20)),
					int32(237-(rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadXboxAxisRightY)*20)), 25, rl.Black)

				// Draw axis: left-right triggers
				rl.DrawRectangle(170, 30, 15, 70, rl.Gray)
				rl.DrawRectangle(604, 30, 15, 70, rl.Gray)
				rl.DrawRectangle(170, 30, 15, int32(((1.0+rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadXboxAxisLt))/2.0)*70), rl.Red)
				rl.DrawRectangle(604, 30, 15, int32(((1.0+rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadXboxAxisRt))/2.0)*70), rl.Red)

			} else if rl.GetGamepadName(rl.GamepadPlayer1) == ps3NameID {
				rl.DrawTexture(texPs3Pad, 0, 0, rl.DarkGray)

				// Draw buttons: ps
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonPs) {
					rl.DrawCircle(396, 222, 13, rl.Red)
				}

				// Draw buttons: basic
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonSelect) {
					rl.DrawRectangle(328, 170, 32, 13, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonStart) {
					rl.DrawTriangle(rl.NewVector2(436, 168), rl.NewVector2(436, 185), rl.NewVector2(464, 177), rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonTriangle) {
					rl.DrawCircle(557, 144, 13, rl.Lime)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonCircle) {
					rl.DrawCircle(586, 173, 13, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonCross) {
					rl.DrawCircle(557, 203, 13, rl.Violet)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonSquare) {
					rl.DrawCircle(527, 173, 13, rl.Pink)
				}

				// Draw buttons: d-pad
				rl.DrawRectangle(225, 132, 24, 84, rl.Black)
				rl.DrawRectangle(195, 161, 84, 25, rl.Black)
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonUp) {
					rl.DrawRectangle(225, 132, 24, 29, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonDown) {
					rl.DrawRectangle(225, 132+54, 24, 30, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonLeft) {
					rl.DrawRectangle(195, 161, 30, 25, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonRight) {
					rl.DrawRectangle(195+54, 161, 30, 25, rl.Red)
				}

				// Draw buttons: left-right back buttons
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonL1) {
					rl.DrawCircle(239, 82, 20, rl.Red)
				}
				if rl.IsGamepadButtonDown(rl.GamepadPlayer1, rl.GamepadPs3ButtonR1) {
					rl.DrawCircle(557, 82, 20, rl.Red)
				}

				// Draw axis: left joystick
				rl.DrawCircle(319, 255, 35, rl.Black)
				rl.DrawCircle(319, 255, 31, rl.LightGray)
				rl.DrawCircle(int32(319+(rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadPs3AxisLeftX)*20)),
					int32(255+(rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadPs3AxisLeftY)*20)), 25, rl.Black)

				// Draw axis: right joystick
				rl.DrawCircle(475, 255, 35, rl.Black)
				rl.DrawCircle(475, 255, 31, rl.LightGray)
				rl.DrawCircle(int32(475+(rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadPs3AxisRightX)*20)),
					int32(255+(rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadPs3AxisRightY)*20)), 25, rl.Black)

				// Draw axis: left-right triggers
				rl.DrawRectangle(169, 48, 15, 70, rl.Gray)
				rl.DrawRectangle(611, 48, 15, 70, rl.Gray)
				rl.DrawRectangle(169, 48, 15, int32(((1.0-rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadPs3AxisL2))/2.0)*70), rl.Red)
				rl.DrawRectangle(611, 48, 15, int32(((1.0-rl.GetGamepadAxisMovement(rl.GamepadPlayer1, rl.GamepadPs3AxisR2))/2.0)*70), rl.Red)
			} else {
				rl.DrawText("- GENERIC GAMEPAD -", 280, 180, 20, rl.Gray)

				// TODO: Draw generic gamepad
			}

			rl.DrawText(fmt.Sprintf("DETECTED AXIS [%d]:", rl.GetGamepadAxisCount(rl.GamepadPlayer1)), 10, 50, 10, rl.Maroon)

			for i := int32(0); i < rl.GetGamepadAxisCount(rl.GamepadPlayer1); i++ {
				rl.DrawText(fmt.Sprintf("AXIS %d: %.02f", i, rl.GetGamepadAxisMovement(rl.GamepadPlayer1, i)), 20, 70+20*i, 10, rl.DarkGray)
			}

			if rl.GetGamepadButtonPressed() != -1 {
				rl.DrawText(fmt.Sprintf("DETECTED BUTTON: %d", rl.GetGamepadButtonPressed()), 10, 430, 10, rl.Red)
			} else {
				rl.DrawText("DETECTED BUTTON: NONE", 10, 430, 10, rl.Gray)
			}
		} else {
			rl.DrawText("GP1: NOT DETECTED", 10, 10, 10, rl.Gray)

			rl.DrawTexture(texXboxPad, 0, 0, rl.LightGray)
		}

		rl.EndDrawing()
	}

	rl.UnloadTexture(texPs3Pad)
	rl.UnloadTexture(texXboxPad)

	rl.CloseWindow()
}
