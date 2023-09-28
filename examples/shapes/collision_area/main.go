package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	screenWidth  = int32(1280)
	screenHeight = int32(720)

	pause     = false
	collision = false

	boxCollision     = rl.Rectangle{}
	screenUpperLimit = float32(40)
	boxAspeedX       = float32(4)
)

func main() {

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - collision area")

	rl.SetTargetFPS(60)

	boxA := rl.NewRectangle(10, float32(rl.GetScreenHeight()/2)-50, 200, 100)
	boxB := rl.NewRectangle(float32(rl.GetScreenWidth()/2)-30, float32(rl.GetScreenHeight()/2)-30, 60, 60)

	for !rl.WindowShouldClose() {

		mousePos := rl.GetMousePosition()

		if !pause {
			boxA.X += boxAspeedX
		}
		if boxA.X+boxA.Width >= float32(rl.GetScreenWidth()) || boxA.X <= 0 {
			boxAspeedX *= -1
		}

		boxB.X = mousePos.X - boxB.Width/2
		boxB.Y = mousePos.Y - boxB.Height/2

		if boxB.X+boxB.Width >= float32(rl.GetScreenWidth()) {
			boxB.X = float32(rl.GetScreenWidth()) - boxB.Width
		} else if boxB.X <= 0 {
			boxB.X = 0
		}

		if boxB.Y+boxB.Height >= float32(rl.GetScreenHeight()) {
			boxB.Y = float32(rl.GetScreenHeight()) - boxB.Height
		} else if boxB.X <= screenUpperLimit {
			boxB.Y = screenUpperLimit
		}

		collision := rl.CheckCollisionRecs(boxA, boxB)

		if collision {
			boxCollision = rl.GetCollisionRec(boxA, boxB)
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			pause = !pause
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		if collision {
			rl.DrawRectangle(0, 0, screenWidth, int32(screenUpperLimit), rl.Red)
			rl.DrawRectangleRec(boxCollision, rl.Lime)
			rl.DrawText("COLLISION", int32(rl.GetScreenWidth()/2)-(rl.MeasureText("COLLISION", 20)/2), int32(screenUpperLimit/2)-10, 20, rl.Black)
			rl.DrawText("Collision Area: "+fmt.Sprint(boxCollision.Width*boxCollision.Height), int32(rl.GetScreenWidth()/2)-100, int32(screenUpperLimit+10), 20, rl.Black)
		} else {
			rl.DrawRectangle(0, 0, screenWidth, int32(screenUpperLimit), rl.Black)
		}

		rl.DrawRectangleRec(boxA, rl.Orange)
		rl.DrawRectangleRec(boxB, rl.Blue)

		rl.DrawText("Press SPACE to PAUSE/RESUME", 20, int32(rl.GetScreenHeight())-35, 20, rl.Black)

		rl.EndDrawing()

	}

	rl.CloseWindow()
}
