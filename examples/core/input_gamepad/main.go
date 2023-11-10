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

	var gamepad int32 = 0 // which gamepad to display

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		if rl.IsGamepadAvailable(gamepad) {
			rl.DrawText(fmt.Sprintf("GP1: %s", rl.GetGamepadName(gamepad)), 10, 10, 10, rl.Black)

			if rl.GetGamepadName(gamepad) == xbox360NameID {
				rl.DrawTexture(texXboxPad, 0, 0, rl.DarkGray)

				// Draw buttons: xbox home
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonMiddle) {
					rl.DrawCircle(394, 89, 19, rl.Red)
				}

				// Draw buttons: basic
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonMiddleRight) {
					rl.DrawCircle(436, 150, 9, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonMiddleLeft) {
					rl.DrawCircle(352, 150, 9, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceLeft) {
					rl.DrawCircle(501, 151, 15, rl.Blue)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceDown) {
					rl.DrawCircle(536, 187, 15, rl.Lime)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceRight) {
					rl.DrawCircle(572, 151, 15, rl.Maroon)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceUp) {
					rl.DrawCircle(536, 115, 15, rl.Gold)
				}

				// Draw buttons: d-pad
				rl.DrawRectangle(317, 202, 19, 71, rl.Black)
				rl.DrawRectangle(293, 228, 69, 19, rl.Black)
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceUp) {
					rl.DrawRectangle(317, 202, 19, 26, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceDown) {
					rl.DrawRectangle(317, 202+45, 19, 26, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceLeft) {
					rl.DrawRectangle(292, 228, 25, 19, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceRight) {
					rl.DrawRectangle(292+44, 228, 26, 19, rl.Red)
				}

				// Draw buttons: left-right back
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftTrigger1) {
					rl.DrawCircle(259, 61, 20, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightTrigger1) {
					rl.DrawCircle(536, 61, 20, rl.Red)
				}

				// Draw axis: left joystick
				rl.DrawCircle(259, 152, 39, rl.Black)
				rl.DrawCircle(259, 152, 34, rl.LightGray)
				rl.DrawCircle(int32(259+(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftX)*20)),
					int32(152-(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftY)*20)), 25, rl.Black)

				// Draw axis: right joystick
				rl.DrawCircle(461, 237, 38, rl.Black)
				rl.DrawCircle(461, 237, 33, rl.LightGray)
				rl.DrawCircle(int32(461+(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightX)*20)),
					int32(237-(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightY)*20)), 25, rl.Black)

				// Draw axis: left-right triggers
				rl.DrawRectangle(170, 30, 15, 70, rl.Gray)
				rl.DrawRectangle(604, 30, 15, 70, rl.Gray)
				rl.DrawRectangle(170, 30, 15, int32(((1.0+rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftTrigger))/2.0)*70), rl.Red)
				rl.DrawRectangle(604, 30, 15, int32(((1.0+rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightTrigger))/2.0)*70), rl.Red)

			} else if rl.GetGamepadName(gamepad) == ps3NameID {
				rl.DrawTexture(texPs3Pad, 0, 0, rl.DarkGray)

				// Draw buttons: ps
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonMiddle) {
					rl.DrawCircle(396, 222, 13, rl.Red)
				}

				// Draw buttons: basic
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonMiddleLeft) {
					rl.DrawRectangle(328, 170, 32, 13, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonMiddleRight) {
					rl.DrawTriangle(rl.NewVector2(436, 168), rl.NewVector2(436, 185), rl.NewVector2(464, 177), rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceUp) {
					rl.DrawCircle(557, 144, 13, rl.Lime)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceRight) {
					rl.DrawCircle(586, 173, 13, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceDown) {
					rl.DrawCircle(557, 203, 13, rl.Violet)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceLeft) {
					rl.DrawCircle(527, 173, 13, rl.Pink)
				}

				// Draw buttons: d-pad
				rl.DrawRectangle(225, 132, 24, 84, rl.Black)
				rl.DrawRectangle(195, 161, 84, 25, rl.Black)
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceUp) {
					rl.DrawRectangle(225, 132, 24, 29, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceDown) {
					rl.DrawRectangle(225, 132+54, 24, 30, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceLeft) {
					rl.DrawRectangle(195, 161, 30, 25, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceRight) {
					rl.DrawRectangle(195+54, 161, 30, 25, rl.Red)
				}

				// Draw buttons: left-right back buttons
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftTrigger1) {
					rl.DrawCircle(239, 82, 20, rl.Red)
				}
				if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightTrigger1) {
					rl.DrawCircle(557, 82, 20, rl.Red)
				}

				// Draw axis: left joystick
				rl.DrawCircle(319, 255, 35, rl.Black)
				rl.DrawCircle(319, 255, 31, rl.LightGray)
				rl.DrawCircle(int32(319+(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftX)*20)),
					int32(255+(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftY)*20)), 25, rl.Black)

				// Draw axis: right joystick
				rl.DrawCircle(475, 255, 35, rl.Black)
				rl.DrawCircle(475, 255, 31, rl.LightGray)
				rl.DrawCircle(int32(475+(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightX)*20)),
					int32(255+(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightY)*20)), 25, rl.Black)

				// Draw axis: left-right triggers
				rl.DrawRectangle(169, 48, 15, 70, rl.Gray)
				rl.DrawRectangle(611, 48, 15, 70, rl.Gray)
				rl.DrawRectangle(169, 48, 15, int32(((1.0-rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftTrigger))/2.0)*70), rl.Red)
				rl.DrawRectangle(611, 48, 15, int32(((1.0-rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightTrigger))/2.0)*70), rl.Red)
			} else {
				rl.DrawText("- GENERIC GAMEPAD -", 280, 180, 20, rl.Gray)

				// TODO: Draw generic gamepad
			}

			rl.DrawText(fmt.Sprintf("DETECTED AXIS [%d]:", rl.GetGamepadAxisCount(gamepad)), 10, 50, 10, rl.Maroon)

			for i := int32(0); i < rl.GetGamepadAxisCount(gamepad); i++ {
				rl.DrawText(fmt.Sprintf("AXIS %d: %.02f", i, rl.GetGamepadAxisMovement(gamepad, i)), 20, 70+20*i, 10, rl.DarkGray)
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
