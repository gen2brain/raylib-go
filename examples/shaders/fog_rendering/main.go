package main

import (
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const screenWidth = 800
	const screenHeight = 450

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - fog rendering")
	defer rl.CloseWindow()

	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(2.0, 2.0, 6.0)
	camera.Target = rl.NewVector3(0.0, 0.5, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	modelA := rl.LoadModelFromMesh(rl.GenMeshTorus(0.4, 1.0, 16, 32))
	modelB := rl.LoadModelFromMesh(rl.GenMeshCube(1.0, 1.0, 1.0))
	modelC := rl.LoadModelFromMesh(rl.GenMeshSphere(0.5, 32, 32))
	texture := rl.LoadTexture("texel_checker.png")

	materialsA := unsafe.Slice(modelA.Materials, modelA.MaterialCount)
	materialsB := unsafe.Slice(modelB.Materials, modelB.MaterialCount)
	materialsC := unsafe.Slice(modelC.Materials, modelC.MaterialCount)

	rl.SetMaterialTexture(&materialsA[0], rl.MapDiffuse, texture)
	rl.SetMaterialTexture(&materialsB[0], rl.MapDiffuse, texture)
	rl.SetMaterialTexture(&materialsC[0], rl.MapDiffuse, texture)

	shader := rl.LoadShader("lighting.vs", "fog.fs")

	shaderLocs := unsafe.Slice(shader.Locs, 32)
	shaderLocs[rl.ShaderLocMatrixModel] = rl.GetShaderLocation(shader, "matModel")
	shaderLocs[rl.ShaderLocVectorView] = rl.GetShaderLocation(shader, "viewPos")

	ambient := []float32{0.2, 0.2, 0.2, 1.0}
	ambientLoc := rl.GetShaderLocation(shader, "ambient")
	rl.SetShaderValue(shader, ambientLoc, ambient, rl.ShaderUniformVec4)

	fogColorNorm := rl.ColorNormalize(rl.Gray)
	fogColor := []float32{fogColorNorm.X, fogColorNorm.Y, fogColorNorm.Z, fogColorNorm.W}
	fogColorLoc := rl.GetShaderLocation(shader, "fogColor")
	rl.SetShaderValue(shader, fogColorLoc, fogColor, rl.ShaderUniformVec4)

	fogDensity := float32(0.15)
	fogDensityLoc := rl.GetShaderLocation(shader, "fogDensity")
	rl.SetShaderValue(shader, fogDensityLoc, []float32{fogDensity}, rl.ShaderUniformFloat)

	materialsA[0].Shader = shader
	materialsB[0].Shader = shader
	materialsC[0].Shader = shader

	CreateLight(LightTypePoint, rl.NewVector3(0, 2, 6), rl.Vector3Zero(), rl.White, shader)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rl.UpdateCamera(&camera, rl.CameraOrbital)

		if rl.IsKeyDown(rl.KeyUp) {
			fogDensity += 0.001
			if fogDensity > 1.0 {
				fogDensity = 1.0
			}
		}

		if rl.IsKeyDown(rl.KeyDown) {
			fogDensity -= 0.001
			if fogDensity < 0.0 {
				fogDensity = 0.0
			}
		}

		rl.SetShaderValue(shader, fogDensityLoc, []float32{fogDensity}, rl.ShaderUniformFloat)

		modelA.Transform = rl.MatrixMultiply(modelA.Transform, rl.MatrixRotateX(-0.025))
		modelA.Transform = rl.MatrixMultiply(modelA.Transform, rl.MatrixRotateZ(0.012))

		rl.SetShaderValue(shader, shaderLocs[rl.ShaderLocVectorView], []float32{camera.Position.X, camera.Position.Y, camera.Position.Z}, rl.ShaderUniformVec3)

		rl.BeginDrawing()

		rl.ClearBackground(rl.Gray)

		rl.BeginMode3D(camera)

		// Draw models
		rl.DrawModel(modelA, rl.Vector3Zero(), 1.0, rl.White)
		rl.DrawModel(modelB, rl.NewVector3(-2.6, 0, 0), 1.0, rl.White)
		rl.DrawModel(modelC, rl.NewVector3(2.6, 0, 0), 1.0, rl.White)

		for i := -20; i < 20; i += 2 {
			rl.DrawModel(modelA, rl.NewVector3(float32(i), 0, 2), 1.0, rl.White)
		}

		rl.EndMode3D()

		rl.DrawText(fmt.Sprintf("Use KEY_UP/KEY_DOWN to change fog density [%.2f]", fogDensity), 10, 10, 20, rl.RayWhite)

		rl.EndDrawing()
	}

	rl.UnloadModel(modelA)
	rl.UnloadModel(modelB)
	rl.UnloadModel(modelC)
	rl.UnloadTexture(texture)
	rl.UnloadShader(shader)
}
