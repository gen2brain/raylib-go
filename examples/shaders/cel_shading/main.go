package main

import (
	"fmt"
	"math"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	cullFaceFront = 0 // RL_CULL_FACE_FRONT
	cullFaceBack  = 1 // RL_CULL_FACE_BACK
)

func main() {
	const screenWidth = 800
	const screenHeight = 450

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - cel shading")
	defer rl.CloseWindow()

	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(9.0, 6.0, 9.0)
	camera.Target = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	model := rl.LoadModel("old_car_new.glb")
	defer rl.UnloadModel(model)

	materials := unsafe.Slice(model.Materials, model.MaterialCount)

	celShader := rl.LoadShader("cel.vs", "cel.fs")
	defer rl.UnloadShader(celShader)

	celShaderLocs := unsafe.Slice(celShader.Locs, 32)
	celShaderLocs[rl.ShaderLocVectorView] = rl.GetShaderLocation(celShader, "viewPos")

	defaultShader := materials[0].Shader
	materials[0].Shader = celShader

	numBands := float32(10.0)
	numBandsLoc := rl.GetShaderLocation(celShader, "numBands")
	rl.SetShaderValue(celShader, numBandsLoc, []float32{numBands}, rl.ShaderUniformFloat)

	outlineShader := rl.LoadShader("outline_hull.vs", "outline_hull.fs")
	defer rl.UnloadShader(outlineShader)

	outlineThicknessLoc := rl.GetShaderLocation(outlineShader, "outlineThickness")

	lights := make([]Light, maxLightsCount)
	lights[0] = CreateLight(LightTypeDirectional, rl.NewVector3(50.0, 50.0, 50.0), rl.Vector3Zero(), rl.White, celShader)

	celEnabled := true
	outlineEnabled := true

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraOrbital)

		cameraPos := []float32{camera.Position.X, camera.Position.Y, camera.Position.Z}
		rl.SetShaderValue(celShader, celShaderLocs[rl.ShaderLocVectorView], cameraPos, rl.ShaderUniformVec3)

		if rl.IsKeyPressed(rl.KeyZ) {
			celEnabled = !celEnabled
			if celEnabled {
				materials[0].Shader = celShader
			} else {
				materials[0].Shader = defaultShader
			}
		}

		if rl.IsKeyPressed(rl.KeyC) {
			outlineEnabled = !outlineEnabled
		}

		if rl.IsKeyPressed(rl.KeyE) || rl.IsKeyPressedRepeat(rl.KeyE) {
			numBands = rl.Clamp(numBands+1.0, 2.0, 20.0)
		}
		if rl.IsKeyPressed(rl.KeyQ) || rl.IsKeyPressedRepeat(rl.KeyQ) {
			numBands = rl.Clamp(numBands-1.0, 2.0, 20.0)
		}
		rl.SetShaderValue(celShader, numBandsLoc, []float32{numBands}, rl.ShaderUniformFloat)

		t := float32(rl.GetTime())
		lights[0].Position = rl.NewVector3(
			float32(math.Sin(float64(-t*0.3)))*5.0,
			5.0,
			float32(math.Cos(float64(-t*0.3)))*5.0,
		)

		for i := 0; i < maxLightsCount; i++ {
			if lights[i].Enabled == 1 {
				lights[i].UpdateValues()
			}
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		if outlineEnabled {

			thickness := float32(0.005)
			rl.SetShaderValue(outlineShader, outlineThicknessLoc, []float32{thickness}, rl.ShaderUniformFloat)

			rl.SetCullFace(cullFaceFront)

			materials[0].Shader = outlineShader

			rl.DrawModel(model, rl.Vector3Zero(), 0.75, rl.White)

			if celEnabled {
				materials[0].Shader = celShader
			} else {
				materials[0].Shader = defaultShader
			}

			rl.SetCullFace(cullFaceBack)
		}

		rl.DrawModel(model, rl.Vector3Zero(), 0.75, rl.White)
		rl.DrawSphereEx(lights[0].Position, 0.2, 50, 50, rl.Yellow)
		rl.DrawGrid(10, 10.0)

		rl.EndMode3D()

		rl.DrawFPS(10, 10)

		celStatus := "OFF"
		celColor := rl.DarkGray
		if celEnabled {
			celStatus = "ON"
			celColor = rl.DarkGreen
		}
		rl.DrawText(fmt.Sprintf("Cel: %s  [Z]", celStatus), 10, 65, 20, celColor)

		outlineStatus := "OFF"
		outlineColor := rl.DarkGray
		if outlineEnabled {
			outlineStatus = "ON"
			outlineColor = rl.DarkGreen
		}
		rl.DrawText(fmt.Sprintf("Outline: %s  [C]", outlineStatus), 10, 90, 20, outlineColor)

		rl.DrawText(fmt.Sprintf("Bands: %.0f  [Q/E]", numBands), 10, 115, 20, rl.DarkGray)

		rl.EndDrawing()
	}
}
