package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - voxel loading")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(10.0, 10.0, 10.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	voxFiles := []string{"knight.vox", "sword.vox", "monument.vox"}

	var models []rl.Model

	for i := 0; i < len(voxFiles); i++ {
		models = append(models, rl.LoadModel(voxFiles[i]))
		bb := rl.GetModelBoundingBox(models[i])
		center := rl.Vector3Zero()
		center.X = bb.Min.X + ((bb.Max.X - bb.Min.X) / 2)
		center.Z = bb.Min.Z + ((bb.Max.Z - bb.Min.Z) / 2)

		matTranslate := rl.MatrixTranslate(-center.X, 0, -center.Z)
		models[i].Transform = matTranslate
	}

	currentModel := 0

	rl.DisableCursor()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rl.UpdateCamera(&camera, rl.CameraOrbital)

		if rl.IsKeyPressed(rl.KeyUp) {
			currentModel++
			if currentModel >= len(models) {
				currentModel = 0
			}

		}
		if rl.IsKeyPressed(rl.KeyDown) {
			currentModel--
			if currentModel < 0 {
				currentModel = len(models) - 1
			}

		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(models[currentModel], rl.Vector3Zero(), 1, rl.White)

		rl.EndMode3D()

		rl.DrawText("current voxel file: "+voxFiles[currentModel], 10, 10, 10, rl.Black)
		rl.DrawText("UP/DOWN ARROW KEYS CHANGE FILE", 10, 30, 10, rl.Black)
		rl.DrawText("MOUSE SCROLL OR KEYPAD + / - TO CHANGE ZOOM", 10, 50, 10, rl.Black)

		rl.EndDrawing()
	}

	for i := 0; i < len(models); i++ {
		rl.UnloadModel(models[i])
	}

	rl.CloseWindow()
}
