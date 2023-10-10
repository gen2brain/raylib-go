package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	rl.SetConfigFlags(rl.FlagMsaa4xHint) //ENABLE 4X MSAA IF AVAILABLE

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - basic lighting")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(2.0, 4.0, 6.0)
	camera.Target = rl.NewVector3(0.0, 0.5, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	ground := rl.LoadModelFromMesh(rl.GenMeshPlane(10, 10, 3, 3))
	cube := rl.LoadModelFromMesh(rl.GenMeshCube(2, 4, 2))

	shader := rl.LoadShader("lighting.vs", "lighting.fs")

	*shader.Locs = rl.GetShaderLocation(shader, "viewPos")

	ambientLoc := rl.GetShaderLocation(shader, "ambient")
	shaderValue := []float32{0.1, 0.1, 0.1, 1.0}
	rl.SetShaderValue(shader, ambientLoc, shaderValue, rl.ShaderUniformVec4)

	ground.Materials.Shader = shader
	cube.Materials.Shader = shader

	lights := make([]Light, 4)
	lights[0] = NewLight(LightTypePoint, rl.NewVector3(-2, 1, -2), rl.NewVector3(0, 0, 0), rl.Yellow, shader)

	lights[1] = NewLight(LightTypePoint, rl.NewVector3(2, 1, 2), rl.NewVector3(0, 0, 0), rl.Red, shader)

	lights[2] = NewLight(LightTypePoint, rl.NewVector3(-2, 1, 2), rl.NewVector3(0, 0, 0), rl.Green, shader)

	lights[3] = NewLight(LightTypePoint, rl.NewVector3(2, 1, -2), rl.NewVector3(0, 0, 0), rl.Blue, shader)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraOrbital)

		cameraPos := []float32{camera.Position.X, camera.Position.Y, camera.Position.Z}
		rl.SetShaderValue(shader, *shader.Locs, cameraPos, rl.ShaderUniformVec3)

		if rl.IsKeyPressed(rl.KeyY) {
			lights[0].enabled *= -1
		}
		if rl.IsKeyPressed(rl.KeyR) {
			lights[1].enabled *= -1
		}
		if rl.IsKeyPressed(rl.KeyG) {
			lights[2].enabled *= -1
		}
		if rl.IsKeyPressed(rl.KeyB) {
			lights[3].enabled *= -1
		}

		for i := 0; i < len(lights); i++ {
			lights[i].UpdateValues()
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(ground, rl.NewVector3(0, 0, 0), 1, rl.White)
		rl.DrawModel(cube, rl.NewVector3(0, 0, 0), 1, rl.White)

		for i := 0; i < len(lights); i++ {
			if lights[i].enabled == 1 {
				rl.DrawSphereEx(lights[i].position, 0.2, 8, 8, lights[i].color)
			} else {
				rl.DrawSphereWires(lights[i].position, 0.2, 8, 8, rl.Fade(lights[i].color, 0.3))
			}
		}

		rl.DrawGrid(10, 1.0) // Draw a grid

		rl.EndMode3D()

		rl.DrawFPS(10, 10)

		rl.DrawText("KEYS [Y] [R] [G] [B] TURN LIGHTS ON/OFF", 10, 40, 20, rl.Black)

		rl.EndDrawing()
	}

	rl.UnloadShader(shader) // Unload shader
	rl.UnloadModel(cube)    // Unload model
	rl.UnloadModel(ground)  // Unload model

	rl.CloseWindow()
}
