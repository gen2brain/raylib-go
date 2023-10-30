package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(800)

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - model animation")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(10.0, 15.0, 10.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 75.0
	camera.Projection = rl.CameraPerspective

	model := rl.LoadModel("guy.iqm")
	texture := rl.LoadTexture("guytex.png")
	rl.SetMaterialTexture(model.Materials, rl.MapDiffuse, texture)

	position := rl.NewVector3(0, 0, 0)

	anims := rl.LoadModelAnimations("guyanim.iqm")
	animFrameCount := 0

	rl.DisableCursor()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rl.UpdateCamera(&camera, rl.CameraOrbital)

		if rl.IsKeyDown(rl.KeySpace) {
			animFrameCount++
			animCurrent := anims[0]
			animFrameNum := animCurrent.FrameCount

			rl.UpdateModelAnimation(model, anims[0], int32(animFrameCount))
			if animFrameCount >= int(animFrameNum) {
				animFrameCount = 0
			}

		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModelEx(model, position, rl.NewVector3(1, 0, 0), -90, rl.NewVector3(1, 1, 1), rl.White)

		rl.DrawGrid(10, 1)

		rl.EndMode3D()

		rl.DrawText("PRESS SPACE to PLAY MODEL ANIMATION", 10, 10, 20, rl.Black)
		rl.DrawText("(c) Guy IQM 3D model by @culacant", 10, 30, 10, rl.Black)

		rl.EndDrawing()
	}

	rl.UnloadModel(model)
	rl.UnloadTexture(texture)

	rl.CloseWindow()
}
