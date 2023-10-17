package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const screenW = int32(1280)
const screenH = int32(720)
const virtualScreenW = screenW / 5
const virtualScreenH = screenH / 5
const virtualRatio = float32(screenW) / float32(virtualScreenW)

func main() {

	rl.InitWindow(screenW, screenH, "raylib [core] example - smooth pixel-perfect camera")

	worldSpaceCam := rl.Camera2D{}
	worldSpaceCam.Zoom = 1

	screenSpaceCam := rl.Camera2D{}
	screenSpaceCam.Zoom = 1

	target := rl.LoadRenderTexture(virtualScreenW, virtualScreenH)

	rec1 := rl.NewRectangle(120, 60, 80, 40)
	rec2 := rl.NewRectangle(130, 70, 90, 30)
	rec3 := rl.NewRectangle(140, 80, 65, 45)

	sourceRec := rl.NewRectangle(0, 0, float32(target.Texture.Width), -float32(target.Texture.Height))
	destRec := rl.NewRectangle(-virtualRatio, -virtualRatio, float32(screenW)+(virtualRatio*2), float32(screenH)+(virtualRatio*2))

	origin := rl.NewVector2(0, 0)

	rotation, camX, camY := float32(0), float32(0), float32(0)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rotation += 60 * rl.GetFrameTime()

		camX = float32(math.Sin(rl.GetTime()) * 50)
		camY = float32(math.Cos(rl.GetTime()) * 30)

		screenSpaceCam.Target = rl.NewVector2(camX, camY)

		worldSpaceCam.Target.X = screenSpaceCam.Target.X
		screenSpaceCam.Target.X -= worldSpaceCam.Target.X
		screenSpaceCam.Target.X *= virtualRatio

		worldSpaceCam.Target.Y = screenSpaceCam.Target.Y
		screenSpaceCam.Target.Y -= worldSpaceCam.Target.Y
		screenSpaceCam.Target.Y *= virtualRatio

		rl.BeginTextureMode(target)
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode2D(worldSpaceCam)

		rl.DrawRectanglePro(rec1, origin, rotation, rl.Black)
		rl.DrawRectanglePro(rec2, origin, -rotation, rl.Red)
		rl.DrawRectanglePro(rec3, origin, rotation+45, rl.Blue)

		rl.EndMode2D()
		rl.EndTextureMode()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Red)
		rl.BeginMode2D(screenSpaceCam)
		rl.DrawTexturePro(target.Texture, sourceRec, destRec, origin, 0, rl.White)
		rl.EndMode2D()

		rl.DrawText("screen res "+fmt.Sprint(screenW)+"x"+fmt.Sprint(screenH), 10, 10, 20, rl.Black)
		rl.DrawText("world res "+fmt.Sprint(virtualScreenW)+"x"+fmt.Sprint(virtualScreenH), 10, 30, 20, rl.Black)
		rl.DrawFPS(screenW-100, 10)

		rl.EndDrawing()
	}

	rl.UnloadRenderTexture(target)

	rl.CloseWindow()
}
