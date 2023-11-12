package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	screenW    = int32(800)
	screenH    = int32(440)
	playerSize = float32(40)
	cam1, cam2 rl.Camera2D
)

func main() {

	rl.InitWindow(screenW, screenH, "raylib [core] example - 2d camera split screen")

	player1 := rl.NewRectangle(200, 200, playerSize, playerSize)
	player2 := rl.NewRectangle(250, 200, playerSize, playerSize)

	cam1.Target = rl.NewVector2(player1.X, player1.Y)
	cam1.Offset = rl.NewVector2(200, 200)
	cam1.Rotation = 0
	cam1.Zoom = 1

	cam2 = cam1
	cam2.Target = rl.NewVector2(player2.X, player2.Y)

	screenCam1 := rl.LoadRenderTexture(screenW/2, screenH)
	screenCam2 := rl.LoadRenderTexture(screenW/2, screenH)

	splitScreenRec := rl.NewRectangle(0, 0, float32(screenCam1.Texture.Width), -float32(screenCam1.Texture.Height))

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		if rl.IsKeyDown(rl.KeyS) {
			player1.Y += 3
		} else if rl.IsKeyDown(rl.KeyW) {
			player1.Y -= 3
		}
		if rl.IsKeyDown(rl.KeyD) {
			player1.X += 3
		} else if rl.IsKeyDown(rl.KeyA) {
			player1.X -= 3
		}

		if rl.IsKeyDown(rl.KeyUp) {
			player2.Y -= 3
		} else if rl.IsKeyDown(rl.KeyDown) {
			player2.Y += 3
		}
		if rl.IsKeyDown(rl.KeyRight) {
			player2.X += 3
		} else if rl.IsKeyDown(rl.KeyLeft) {
			player2.X -= 3
		}

		cam1.Target = rl.NewVector2(player1.X, player1.Y)
		cam2.Target = rl.NewVector2(player2.X, player2.Y)

		rl.BeginTextureMode(screenCam1)
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode2D(cam1)

		for i := 0; i < int(screenW/int32(playerSize))+1; i++ {
			rl.DrawLineV(rl.NewVector2(playerSize*float32(i), 0), rl.NewVector2(playerSize*float32(i), float32(screenH)), rl.LightGray)
		}
		for i := 0; i < int(screenH/int32(playerSize))+1; i++ {
			rl.DrawLineV(rl.NewVector2(0, playerSize*float32(i)), rl.NewVector2(float32(screenW), playerSize*float32(i)), rl.LightGray)
		}
		for i := 0; i < int(screenW/int32(playerSize)); i++ {
			for j := 0; j < int(screenH/int32(playerSize)); j++ {
				rl.DrawText("["+fmt.Sprint(i)+","+fmt.Sprint(j)+"]", 10+int32(playerSize*float32(i)), 15+int32(playerSize*float32(j)), 10, rl.LightGray)
			}
		}

		rl.DrawRectangleRec(player1, rl.Red)
		rl.DrawRectangleRec(player2, rl.Blue)
		rl.EndMode2D()

		rl.DrawRectangle(0, 0, screenW/2, 30, rl.Fade(rl.RayWhite, 0.6))
		rl.DrawText("PLAYER 1 WASD KEYS", 10, 10, 10, rl.Maroon)
		rl.EndTextureMode()

		rl.BeginTextureMode(screenCam2)
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode2D(cam2)

		for i := 0; i < int(screenW/int32(playerSize))+1; i++ {
			rl.DrawLineV(rl.NewVector2(playerSize*float32(i), 0), rl.NewVector2(playerSize*float32(i), float32(screenH)), rl.LightGray)
		}
		for i := 0; i < int(screenH/int32(playerSize))+1; i++ {
			rl.DrawLineV(rl.NewVector2(0, playerSize*float32(i)), rl.NewVector2(float32(screenW), playerSize*float32(i)), rl.LightGray)
		}
		for i := 0; i < int(screenW/int32(playerSize)); i++ {
			for j := 0; j < int(screenH/int32(playerSize)); j++ {
				rl.DrawText("["+fmt.Sprint(i)+","+fmt.Sprint(j)+"]", 10+int32(playerSize*float32(i)), 15+int32(playerSize*float32(j)), 10, rl.LightGray)
			}
		}

		rl.DrawRectangleRec(player1, rl.Red)
		rl.DrawRectangleRec(player2, rl.Blue)
		rl.EndMode2D()

		rl.DrawRectangle(0, 0, screenW/2, 30, rl.Fade(rl.RayWhite, 0.6))
		rl.DrawText("PLAYER 2 ARROW KEYS", 10, 10, 10, rl.Maroon)
		rl.EndTextureMode()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.DrawTextureRec(screenCam1.Texture, splitScreenRec, rl.NewVector2(0, 0), rl.White)
		rl.DrawTextureRec(screenCam2.Texture, splitScreenRec, rl.NewVector2(float32(screenW/2), 0), rl.White)
		rl.DrawRectangle((screenW/2)-2, 0, 4, screenH, rl.LightGray)

		rl.EndDrawing()

	}

	rl.UnloadRenderTexture(screenCam1)
	rl.UnloadRenderTexture(screenCam2)

	rl.CloseWindow()
}
