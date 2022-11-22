package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagMsaa4xHint) // Enable Multi Sampling Anti Aliasing 4x (if available)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - model shader")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(3.0, 3.0, 3.0)
	camera.Target = rl.NewVector3(0.0, 1.5, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0

	dwarf := rl.LoadModel("dwarf.obj")                                 // Load OBJ model
	texture := rl.LoadTexture("dwarf_diffuse.png")                     // Load model texture
	shader := rl.LoadShader("glsl330/base.vs", "glsl330/grayscale.fs") // Load model shader

	rl.SetMaterialTexture(dwarf.Materials, rl.MapDiffuse, texture)
	dwarf.Materials.Shader = shader // Set shader effect to 3d model

	position := rl.NewVector3(0.0, 0.0, 0.0) // Set model position

	rl.SetCameraMode(camera, rl.CameraFree) // Set free camera mode

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera) // Update camera

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(dwarf, position, 2.0, rl.White) // Draw 3d model with texture

		rl.DrawGrid(10, 1.0) // Draw a grid

		rl.EndMode3D()

		rl.DrawText("(c) Dwarf 3D model by David Moreno", screenWidth-200, screenHeight-20, 10, rl.Gray)

		rl.DrawText(fmt.Sprintf("Camera position: (%.2f, %.2f, %.2f)", camera.Position.X, camera.Position.Y, camera.Position.Z), 600, 20, 10, rl.Black)
		rl.DrawText(fmt.Sprintf("Camera target: (%.2f, %.2f, %.2f)", camera.Target.X, camera.Target.Y, camera.Target.Z), 600, 40, 10, rl.Gray)

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	rl.UnloadShader(shader)   // Unload shader
	rl.UnloadTexture(texture) // Unload texture
	rl.UnloadModel(dwarf)     // Unload model

	rl.CloseWindow()
}
