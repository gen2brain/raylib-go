package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.SetConfigFlags(raylib.FlagMsaa4xHint) // Enable Multi Sampling Anti Aliasing 4x (if available)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - model shader")

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(3.0, 3.0, 3.0)
	camera.Target = raylib.NewVector3(0.0, 1.5, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0

	dwarf := raylib.LoadModel("dwarf.obj")                                 // Load OBJ model
	texture := raylib.LoadTexture("dwarf_diffuse.png")                     // Load model texture
	shader := raylib.LoadShader("glsl330/base.vs", "glsl330/grayscale.fs") // Load model shader

	dwarf.Material.Shader = shader                           // Set shader effect to 3d model
	dwarf.Material.Maps[raylib.MapDiffuse].Texture = texture // Set dwarf model diffuse texture

	position := raylib.NewVector3(0.0, 0.0, 0.0) // Set model position

	raylib.SetCameraMode(camera, raylib.CameraFree) // Set free camera mode

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera) // Update camera

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode3D(camera)

		raylib.DrawModel(dwarf, position, 2.0, raylib.White) // Draw 3d model with texture

		raylib.DrawGrid(10, 1.0) // Draw a grid

		raylib.EndMode3D()

		raylib.DrawText("(c) Dwarf 3D model by David Moreno", screenWidth-200, screenHeight-20, 10, raylib.Gray)

		raylib.DrawText(fmt.Sprintf("Camera position: (%.2f, %.2f, %.2f)", camera.Position.X, camera.Position.Y, camera.Position.Z), 600, 20, 10, raylib.Black)
		raylib.DrawText(fmt.Sprintf("Camera target: (%.2f, %.2f, %.2f)", camera.Target.X, camera.Target.Y, camera.Target.Z), 600, 40, 10, raylib.Gray)

		raylib.DrawFPS(10, 10)

		raylib.EndDrawing()
	}

	raylib.UnloadShader(shader)   // Unload shader
	raylib.UnloadTexture(texture) // Unload texture
	raylib.UnloadModel(dwarf)     // Unload model

	raylib.CloseWindow()
}
