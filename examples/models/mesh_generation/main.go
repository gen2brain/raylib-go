package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	numModels := 8

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - mesh generation")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(10.0, 5.0, 10.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	checked := rl.GenImageChecked(2, 2, 1, 1, rl.Black, rl.Red)
	texture := rl.LoadTextureFromImage(checked)
	rl.UnloadImage(checked)

	models := make([]rl.Model, numModels)

	models[0] = rl.LoadModelFromMesh(rl.GenMeshPlane(2, 2, 4, 3))
	models[1] = rl.LoadModelFromMesh(rl.GenMeshCube(2, 1, 2))
	models[2] = rl.LoadModelFromMesh(rl.GenMeshSphere(2, 32, 32))
	models[3] = rl.LoadModelFromMesh(rl.GenMeshHemiSphere(2, 16, 16))
	models[4] = rl.LoadModelFromMesh(rl.GenMeshCylinder(1, 2, 16))
	models[5] = rl.LoadModelFromMesh(rl.GenMeshTorus(0.25, 4, 16, 32))
	models[6] = rl.LoadModelFromMesh(rl.GenMeshKnot(1, 2, 16, 128))
	models[7] = rl.LoadModelFromMesh(rl.GenMeshPoly(5, 2))

	for i := 0; i < numModels; i++ {
		rl.SetMaterialTexture(models[i].Materials, rl.MapDiffuse, texture)
	}

	position := rl.Vector3Zero()

	currentModel := 0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rl.UpdateCamera(&camera, rl.CameraOrbital)

		if rl.IsKeyPressed(rl.KeyUp) {
			currentModel++
			if currentModel >= numModels {
				currentModel = 0
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(models[currentModel], position, 1, rl.White)
		rl.DrawGrid(10, 1)

		rl.EndMode3D()

		rl.DrawRectangle(10, 10, 310, 30, rl.Fade(rl.SkyBlue, 0.5))
		rl.DrawRectangleLines(10, 10, 310, 30, rl.Fade(rl.DarkBlue, 0.5))
		rl.DrawText("UP ARROW KEY TO CHANGE MODELS", 20, 20, 10, rl.Blue)

		txt := "PLANE"
		switch currentModel {
		case 1:
			txt = "CUBE"
		case 2:
			txt = "SPHERE"
		case 3:
			txt = "HEMISPHERE"
		case 4:
			txt = "CYLINDER"
		case 5:
			txt = "TORUS"
		case 6:
			txt = "KNOT"
		case 7:
			txt = "POLY"
		}
		txtlen := rl.MeasureText(txt, 20)
		rl.DrawText(txt, screenWidth/2-txtlen/2, 10, 20, rl.DarkBlue)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)
	for i := 0; i < numModels; i++ {
		rl.UnloadModel(models[i])
	}

	rl.CloseWindow()
}
