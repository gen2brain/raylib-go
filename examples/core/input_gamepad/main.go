package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	xbox360NameID = "Xbox 360 Controller"
	ps3NameID     = "PLAYSTATION(R)3 Controller"
)

func main() {
	raylib.SetConfigFlags(raylib.FlagMsaa4xHint) // Set MSAA 4X hint before windows creation

	raylib.InitWindow(800, 450, "raylib [core] example - gamepad input")

	texPs3Pad := raylib.LoadTexture("ps3.png")
	texXboxPad := raylib.LoadTexture("xbox.png")

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		if raylib.IsGamepadAvailable(raylib.GamepadPlayer1) {
			raylib.DrawText(fmt.Sprintf("GP1: %s", raylib.GetGamepadName(raylib.GamepadPlayer1)), 10, 10, 10, raylib.Black)

			if raylib.IsGamepadName(raylib.GamepadPlayer1, xbox360NameID) {
				raylib.DrawTexture(texXboxPad, 0, 0, raylib.DarkGray)

				// Draw buttons: xbox home
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonHome) {
					raylib.DrawCircle(394, 89, 19, raylib.Red)
				}

				// Draw buttons: basic
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonStart) {
					raylib.DrawCircle(436, 150, 9, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonSelect) {
					raylib.DrawCircle(352, 150, 9, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonX) {
					raylib.DrawCircle(501, 151, 15, raylib.Blue)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonA) {
					raylib.DrawCircle(536, 187, 15, raylib.Lime)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonB) {
					raylib.DrawCircle(572, 151, 15, raylib.Maroon)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonY) {
					raylib.DrawCircle(536, 115, 15, raylib.Gold)
				}

				// Draw buttons: d-pad
				raylib.DrawRectangle(317, 202, 19, 71, raylib.Black)
				raylib.DrawRectangle(293, 228, 69, 19, raylib.Black)
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonUp) {
					raylib.DrawRectangle(317, 202, 19, 26, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonDown) {
					raylib.DrawRectangle(317, 202+45, 19, 26, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonLeft) {
					raylib.DrawRectangle(292, 228, 25, 19, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonRight) {
					raylib.DrawRectangle(292+44, 228, 26, 19, raylib.Red)
				}

				// Draw buttons: left-right back
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonLb) {
					raylib.DrawCircle(259, 61, 20, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadXboxButtonRb) {
					raylib.DrawCircle(536, 61, 20, raylib.Red)
				}

				// Draw axis: left joystick
				raylib.DrawCircle(259, 152, 39, raylib.Black)
				raylib.DrawCircle(259, 152, 34, raylib.LightGray)
				raylib.DrawCircle(int32(259+(raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadXboxAxisLeftX)*20)),
					int32(152-(raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadXboxAxisLeftY)*20)), 25, raylib.Black)

				// Draw axis: right joystick
				raylib.DrawCircle(461, 237, 38, raylib.Black)
				raylib.DrawCircle(461, 237, 33, raylib.LightGray)
				raylib.DrawCircle(int32(461+(raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadXboxAxisRightX)*20)),
					int32(237-(raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadXboxAxisRightY)*20)), 25, raylib.Black)

				// Draw axis: left-right triggers
				raylib.DrawRectangle(170, 30, 15, 70, raylib.Gray)
				raylib.DrawRectangle(604, 30, 15, 70, raylib.Gray)
				raylib.DrawRectangle(170, 30, 15, int32(((1.0+raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadXboxAxisLt))/2.0)*70), raylib.Red)
				raylib.DrawRectangle(604, 30, 15, int32(((1.0+raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadXboxAxisRt))/2.0)*70), raylib.Red)

			} else if raylib.IsGamepadName(raylib.GamepadPlayer1, ps3NameID) {
				raylib.DrawTexture(texPs3Pad, 0, 0, raylib.DarkGray)

				// Draw buttons: ps
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonPs) {
					raylib.DrawCircle(396, 222, 13, raylib.Red)
				}

				// Draw buttons: basic
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonSelect) {
					raylib.DrawRectangle(328, 170, 32, 13, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonStart) {
					raylib.DrawTriangle(raylib.NewVector2(436, 168), raylib.NewVector2(436, 185), raylib.NewVector2(464, 177), raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonTriangle) {
					raylib.DrawCircle(557, 144, 13, raylib.Lime)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonCircle) {
					raylib.DrawCircle(586, 173, 13, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonCross) {
					raylib.DrawCircle(557, 203, 13, raylib.Violet)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonSquare) {
					raylib.DrawCircle(527, 173, 13, raylib.Pink)
				}

				// Draw buttons: d-pad
				raylib.DrawRectangle(225, 132, 24, 84, raylib.Black)
				raylib.DrawRectangle(195, 161, 84, 25, raylib.Black)
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonUp) {
					raylib.DrawRectangle(225, 132, 24, 29, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonDown) {
					raylib.DrawRectangle(225, 132+54, 24, 30, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonLeft) {
					raylib.DrawRectangle(195, 161, 30, 25, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonRight) {
					raylib.DrawRectangle(195+54, 161, 30, 25, raylib.Red)
				}

				// Draw buttons: left-right back buttons
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonL1) {
					raylib.DrawCircle(239, 82, 20, raylib.Red)
				}
				if raylib.IsGamepadButtonDown(raylib.GamepadPlayer1, raylib.GamepadPs3ButtonR1) {
					raylib.DrawCircle(557, 82, 20, raylib.Red)
				}

				// Draw axis: left joystick
				raylib.DrawCircle(319, 255, 35, raylib.Black)
				raylib.DrawCircle(319, 255, 31, raylib.LightGray)
				raylib.DrawCircle(int32(319+(raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadPs3AxisLeftX)*20)),
					int32(255+(raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadPs3AxisLeftY)*20)), 25, raylib.Black)

				// Draw axis: right joystick
				raylib.DrawCircle(475, 255, 35, raylib.Black)
				raylib.DrawCircle(475, 255, 31, raylib.LightGray)
				raylib.DrawCircle(int32(475+(raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadPs3AxisRightX)*20)),
					int32(255+(raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadPs3AxisRightY)*20)), 25, raylib.Black)

				// Draw axis: left-right triggers
				raylib.DrawRectangle(169, 48, 15, 70, raylib.Gray)
				raylib.DrawRectangle(611, 48, 15, 70, raylib.Gray)
				raylib.DrawRectangle(169, 48, 15, int32(((1.0-raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadPs3AxisL2))/2.0)*70), raylib.Red)
				raylib.DrawRectangle(611, 48, 15, int32(((1.0-raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, raylib.GamepadPs3AxisR2))/2.0)*70), raylib.Red)
			} else {
				raylib.DrawText("- GENERIC GAMEPAD -", 280, 180, 20, raylib.Gray)

				// TODO: Draw generic gamepad
			}

			raylib.DrawText(fmt.Sprintf("DETECTED AXIS [%d]:", raylib.GetGamepadAxisCount(raylib.GamepadPlayer1)), 10, 50, 10, raylib.Maroon)

			for i := int32(0); i < raylib.GetGamepadAxisCount(raylib.GamepadPlayer1); i++ {
				raylib.DrawText(fmt.Sprintf("AXIS %d: %.02f", i, raylib.GetGamepadAxisMovement(raylib.GamepadPlayer1, i)), 20, 70+20*i, 10, raylib.DarkGray)
			}

			if raylib.GetGamepadButtonPressed() != -1 {
				raylib.DrawText(fmt.Sprintf("DETECTED BUTTON: %d", raylib.GetGamepadButtonPressed()), 10, 430, 10, raylib.Red)
			} else {
				raylib.DrawText("DETECTED BUTTON: NONE", 10, 430, 10, raylib.Gray)
			}
		} else {
			raylib.DrawText("GP1: NOT DETECTED", 10, 10, 10, raylib.Gray)

			raylib.DrawTexture(texXboxPad, 0, 0, raylib.LightGray)
		}

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(texPs3Pad)
	raylib.UnloadTexture(texXboxPad)

	raylib.CloseWindow()
}
