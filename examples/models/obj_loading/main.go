package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - obj model loading")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(50.0, 50.0, 50.0)
	camera.Target = rl.NewVector3(0.0, 10.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	obj := rl.LoadModel("castle.obj")               // Load OBJ model
	texture := rl.LoadTexture("castle_diffuse.png") // Load model texture

	rl.SetMaterialTexture(obj.Materials, rl.MapDiffuse, texture) // Set map diffuse texture

	position := rl.NewVector3(0.0, 0.0, 0.0) // Set model position

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(obj, position, 1.0, rl.White) // Draw 3d model with texture

		rl.DrawGrid(20, 10.0) // Draw a grid

		rl.EndMode3D()

		rl.DrawText("(c) Castle 3D model by Alberto Cano", screenWidth-200, screenHeight-20, 10, rl.Gray)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture) // Unload texture
	rl.UnloadModel(obj)       // Unload model

	rl.CloseWindow()
}
