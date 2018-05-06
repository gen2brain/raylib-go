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

	raylib.InitWindow(screenWidth, screenHeight, "raylib [core] example - gestures detection")

	touchPosition := raylib.NewVector2(0, 0)
	touchArea := raylib.NewRectangle(220, 10, float32(screenWidth)-230, float32(screenHeight)-20)

	gestureStrings := make([]string, 0)

	currentGesture := raylib.GestureNone
	lastGesture := raylib.GestureNone

	//raylib.SetGesturesEnabled(uint32(raylib.GestureHold | raylib.GestureDrag)) // Enable only some gestures to be detected

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		lastGesture = currentGesture
		currentGesture = raylib.GetGestureDetected()
		touchPosition = raylib.GetTouchPosition(0)

		if raylib.CheckCollisionPointRec(touchPosition, touchArea) && currentGesture != raylib.GestureNone {
			if currentGesture != lastGesture {
				switch currentGesture {
				case raylib.GestureTap:
					gestureStrings = append(gestureStrings, "GESTURE TAP")
				case raylib.GestureDoubletap:
					gestureStrings = append(gestureStrings, "GESTURE DOUBLETAP")
				case raylib.GestureHold:
					gestureStrings = append(gestureStrings, "GESTURE HOLD")
				case raylib.GestureDrag:
					gestureStrings = append(gestureStrings, "GESTURE DRAG")
				case raylib.GestureSwipeRight:
					gestureStrings = append(gestureStrings, "GESTURE SWIPE RIGHT")
				case raylib.GestureSwipeLeft:
					gestureStrings = append(gestureStrings, "GESTURE SWIPE LEFT")
				case raylib.GestureSwipeUp:
					gestureStrings = append(gestureStrings, "GESTURE SWIPE UP")
				case raylib.GestureSwipeDown:
					gestureStrings = append(gestureStrings, "GESTURE SWIPE DOWN")
				case raylib.GesturePinchIn:
					gestureStrings = append(gestureStrings, "GESTURE PINCH IN")
				case raylib.GesturePinchOut:
					gestureStrings = append(gestureStrings, "GESTURE PINCH OUT")
				}

				if len(gestureStrings) >= maxGestureStrings {
					gestureStrings = make([]string, 0)
				}
			}
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawRectangleRec(touchArea, raylib.Gray)
		raylib.DrawRectangle(225, 15, screenWidth-240, screenHeight-30, raylib.RayWhite)

		raylib.DrawText("GESTURES TEST AREA", screenWidth-270, screenHeight-40, 20, raylib.Fade(raylib.Gray, 0.5))

		for i := 0; i < len(gestureStrings); i++ {
			if i%2 == 0 {
				raylib.DrawRectangle(10, int32(30+20*i), 200, 20, raylib.Fade(raylib.LightGray, 0.5))
			} else {
				raylib.DrawRectangle(10, int32(30+20*i), 200, 20, raylib.Fade(raylib.LightGray, 0.3))
			}

			if i < len(gestureStrings)-1 {
				raylib.DrawText(gestureStrings[i], 35, int32(36+20*i), 10, raylib.DarkGray)
			} else {
				raylib.DrawText(gestureStrings[i], 35, int32(36+20*i), 10, raylib.Maroon)
			}
		}

		raylib.DrawRectangleLines(10, 29, 200, screenHeight-50, raylib.Gray)
		raylib.DrawText("DETECTED GESTURES", 50, 15, 10, raylib.Gray)

		if currentGesture != raylib.GestureNone {
			raylib.DrawCircleV(touchPosition, 30, raylib.Maroon)
		}

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
