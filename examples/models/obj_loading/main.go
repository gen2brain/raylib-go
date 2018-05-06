package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [models] example - obj model loading")

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(3.0, 3.0, 3.0)
	camera.Target = raylib.NewVector3(0.0, 1.5, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0

	dwarf := raylib.LoadModel("dwarf.obj")             // Load OBJ model
	texture := raylib.LoadTexture("dwarf_diffuse.png") // Load model texture

	dwarf.Material.Maps[raylib.MapDiffuse].Texture = texture // Set dwarf model diffuse texture

	position := raylib.NewVector3(0.0, 0.0, 0.0) // Set model position

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode3D(camera)

		raylib.DrawModel(dwarf, position, 2.0, raylib.White) // Draw 3d model with texture

		raylib.DrawGrid(10, 1.0) // Draw a grid

		raylib.DrawGizmo(position) // Draw gizmo

		raylib.EndMode3D()

		raylib.DrawText("(c) Dwarf 3D model by David Moreno", screenWidth-200, screenHeight-20, 10, raylib.Gray)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(texture) // Unload texture
	raylib.UnloadModel(dwarf)     // Unload model

	raylib.CloseWindow()
}
