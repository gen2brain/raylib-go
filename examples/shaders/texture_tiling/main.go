package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - texture tiling")

	cam := rl.Camera3D{}
	cam.Position = rl.NewVector3(4, 4, 4)
	cam.Target = rl.NewVector3(0, 0.5, 0)
	cam.Up = rl.NewVector3(0, 1, 0)
	cam.Fovy = 45
	cam.Projection = rl.CameraPerspective

	cube := rl.GenMeshCube(1, 1, 1)
	model := rl.LoadModelFromMesh(cube)

	texture := rl.LoadTexture("cubicmap_atlas.png")
	rl.SetMaterialTexture(model.Materials, rl.MapDiffuse, texture)

	tiling := []float32{3, 3}
	shader := rl.LoadShader("", "tiling.fs")
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "tiling"), tiling, rl.ShaderUniformVec2)
	model.Materials.Shader = shader

	rl.DisableCursor()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rl.UpdateCamera(&cam, rl.CameraOrbital)

		if rl.IsKeyPressed(rl.KeyZ) {
			cam.Target = rl.NewVector3(0, 0.5, 0)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(cam)

		rl.BeginShaderMode(shader)
		rl.DrawModel(model, rl.Vector3Zero(), 2, rl.White)
		rl.EndShaderMode()

		rl.DrawGrid(10, 1)

		rl.EndMode3D()

		rl.EndDrawing()
	}

	rl.UnloadShader(shader)
	rl.UnloadModel(model)
	rl.UnloadTexture(texture)

	rl.CloseWindow()
}
