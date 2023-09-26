package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	rl.SetConfigFlags(rl.FlagMsaa4xHint) // ENABLE 4X MSAA IF AVAILABLE

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - raymarching")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(2.5, 2.5, 3.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.7)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 65.0

	shader := rl.LoadShader("", "raymarching.fs")

	viewEyeLoc := rl.GetShaderLocation(shader, "viewEye")
	viewCenterLoc := rl.GetShaderLocation(shader, "viewCenter")
	runtTimeLoc := rl.GetShaderLocation(shader, "runTime")
	resolutionLoc := rl.GetShaderLocation(shader, "resolution")

	resolution := []float32{float32(screenWidth), float32(screenHeight)}
	rl.SetShaderValue(shader, resolutionLoc, resolution, rl.ShaderUniformVec2)

	runTimer := float32(0)

	rl.DisableCursor()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraThirdPerson) // Update camera with free camera mode

		camPos := []float32{camera.Position.X, camera.Position.Y, camera.Position.Z}
		camTarget := []float32{camera.Target.X, camera.Target.Y, camera.Target.Z}

		deltaTime := rl.GetFrameTime()
		runTimer += deltaTime
		runTime := []float32{runTimer}

		rl.SetShaderValue(shader, viewEyeLoc, camPos, rl.ShaderUniformVec3)
		rl.SetShaderValue(shader, viewCenterLoc, camTarget, rl.ShaderUniformVec3)
		rl.SetShaderValue(shader, runtTimeLoc, runTime, rl.ShaderUniformFloat)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginShaderMode(shader)

		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.White)
		rl.EndShaderMode()

		rl.DrawText("(c) Raymarching shader by Iñigo Quilez. MIT License.", screenWidth-280, screenHeight-20, 10, rl.Black)

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	rl.UnloadShader(shader) // Unload shader

	rl.CloseWindow()
}
