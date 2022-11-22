package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

var (
	maxGestureStrings int = 20
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - gestures detection")

	touchPosition := rl.NewVector2(0, 0)
	touchArea := rl.NewRectangle(220, 10, float32(screenWidth)-230, float32(screenHeight)-20)

	gestureStrings := make([]string, 0)

	currentGesture := rl.GestureNone
	lastGesture := rl.GestureNone

	//rl.SetGesturesEnabled(uint32(rl.GestureHold | rl.GestureDrag)) // Enable only some gestures to be detected

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		lastGesture = currentGesture
		currentGesture = rl.GetGestureDetected()
		touchPosition = rl.GetTouchPosition(0)

		if rl.CheckCollisionPointRec(touchPosition, touchArea) && currentGesture != rl.GestureNone {
			if currentGesture != lastGesture {
				switch currentGesture {
				case rl.GestureTap:
					gestureStrings = append(gestureStrings, "GESTURE TAP")
				case rl.GestureDoubletap:
					gestureStrings = append(gestureStrings, "GESTURE DOUBLETAP")
				case rl.GestureHold:
					gestureStrings = append(gestureStrings, "GESTURE HOLD")
				case rl.GestureDrag:
					gestureStrings = append(gestureStrings, "GESTURE DRAG")
				case rl.GestureSwipeRight:
					gestureStrings = append(gestureStrings, "GESTURE SWIPE RIGHT")
				case rl.GestureSwipeLeft:
					gestureStrings = append(gestureStrings, "GESTURE SWIPE LEFT")
				case rl.GestureSwipeUp:
					gestureStrings = append(gestureStrings, "GESTURE SWIPE UP")
				case rl.GestureSwipeDown:
					gestureStrings = append(gestureStrings, "GESTURE SWIPE DOWN")
				case rl.GesturePinchIn:
					gestureStrings = append(gestureStrings, "GESTURE PINCH IN")
				case rl.GesturePinchOut:
					gestureStrings = append(gestureStrings, "GESTURE PINCH OUT")
				}

				if len(gestureStrings) >= maxGestureStrings {
					gestureStrings = make([]string, 0)
				}
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawRectangleRec(touchArea, rl.Gray)
		rl.DrawRectangle(225, 15, screenWidth-240, screenHeight-30, rl.RayWhite)

		rl.DrawText("GESTURES TEST AREA", screenWidth-270, screenHeight-40, 20, rl.Fade(rl.Gray, 0.5))

		for i := 0; i < len(gestureStrings); i++ {
			if i%2 == 0 {
				rl.DrawRectangle(10, int32(30+20*i), 200, 20, rl.Fade(rl.LightGray, 0.5))
			} else {
				rl.DrawRectangle(10, int32(30+20*i), 200, 20, rl.Fade(rl.LightGray, 0.3))
			}

			if i < len(gestureStrings)-1 {
				rl.DrawText(gestureStrings[i], 35, int32(36+20*i), 10, rl.DarkGray)
			} else {
				rl.DrawText(gestureStrings[i], 35, int32(36+20*i), 10, rl.Maroon)
			}
		}

		rl.DrawRectangleLines(10, 29, 200, screenHeight-50, rl.Gray)
		rl.DrawText("DETECTED GESTURES", 50, 15, 10, rl.Gray)

		if currentGesture != rl.GestureNone {
			rl.DrawCircleV(touchPosition, 30, rl.Maroon)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
