package main

import (
	"path/filepath"
	"slices"
	"unsafe"

	"github.com/gen2brain/raylib-go/raylib"
)

var supportedFileTypes = []string{
	".obj",
	".gltf",
	".glb",
	".vox",
	".iqm",
	".m3d",
}

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

	meshes := unsafe.Slice(obj.Meshes, obj.MeshCount)
	bounds := rl.GetMeshBoundingBox(meshes[0]) // Set model bounds

	selected := false

	rl.DisableCursor()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)

		// Load new models/textures on drag&drop
		if rl.IsFileDropped() {
			droppedFiles := rl.LoadDroppedFiles()

			if len(droppedFiles) == 1 { // Only support one file dropped
				if slices.Contains(supportedFileTypes, filepath.Ext(droppedFiles[0])) { // Model file formats supported
					rl.UnloadModel(obj)                                          // Unload previous model
					obj = rl.LoadModel(droppedFiles[0])                          // Load new model
					rl.SetMaterialTexture(obj.Materials, rl.MapDiffuse, texture) // Set current map diffuse texture

					meshes = unsafe.Slice(obj.Meshes, obj.MeshCount)
					bounds = rl.GetMeshBoundingBox(meshes[0])

					// TODO: Move camera position from target enough distance to visualize model properly
				} else if filepath.Ext(droppedFiles[0]) == ".png" { // Texture file formats supported
					// Unload current model texture and load new one
					rl.UnloadTexture(texture)
					texture = rl.LoadTexture(droppedFiles[0])
					rl.SetMaterialTexture(obj.Materials, rl.MapDiffuse, texture) // Set current map diffuse texture
				}
			}

			rl.UnloadDroppedFiles() // Unload file paths from memory
		}

		// Select model on mouse click
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			// Check collision between ray and box
			if rl.GetRayCollisionBox(rl.GetMouseRay(rl.GetMousePosition(), camera), bounds).Hit {
				selected = !selected
			} else {
				selected = false
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode3D(camera)

		rl.DrawModel(obj, position, 1.0, rl.White) // Draw 3d model with texture
		rl.DrawGrid(20, 10.0)                      // Draw a grid

		if selected {
			rl.DrawBoundingBox(bounds, rl.Green) // Draw selection box
		}

		rl.EndMode3D()

		rl.DrawText("Drag & drop model to load mesh/texture", 10, screenHeight-20, 10, rl.Gray)
		if selected {
			rl.DrawText("Model selected!", screenWidth-110, 10, 10, rl.Green)
		}
		rl.DrawText("(c) Castle 3D model by Alberto Cano", screenWidth-200, screenHeight-20, 10, rl.Gray)

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
	rl.UnloadTexture(texture) // Unload texture
	rl.UnloadModel(obj)       // Unload model

	rl.CloseWindow()
}
