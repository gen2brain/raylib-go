package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - cubesmap loading and drawing")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(16.0, 14.0, 16.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0

	image := rl.LoadImage("cubicmap.png")      // Load cubicmap image (RAM)
	cubicmap := rl.LoadTextureFromImage(image) // Convert image to texture to display (VRAM)

	mesh := rl.GenMeshCubicmap(*image, rl.NewVector3(1.0, 1.0, 1.0))
	model := rl.LoadModelFromMesh(mesh)

	// NOTE: By default each cube is mapped to one part of texture atlas
	texture := rl.LoadTexture("cubicmap_atlas.png")                // Load map texture
	rl.SetMaterialTexture(model.Materials, rl.MapDiffuse, texture) // Set map diffuse texture

	mapPosition := rl.NewVector3(-16.0, 0.0, -8.0) // Set model position

	rl.UnloadImage(image) // Unload cubicmap image from RAM, already uploaded to VRAM

	rl.SetCameraMode(camera, rl.CameraOrbital) // Set an orbital camera mode

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update

		rl.UpdateCamera(&camera) // Update camera

		// Draw

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(model, mapPosition, 1.0, rl.White)

		rl.EndMode3D()

		rl.DrawTextureEx(cubicmap, rl.NewVector2(float32(screenWidth-cubicmap.Width*4-20), 20), 0.0, 4.0, rl.White)
		rl.DrawRectangleLines(screenWidth-cubicmap.Width*4-20, 20, cubicmap.Width*4, cubicmap.Height*4, rl.Green)

		rl.DrawText("cubicmap image used to", 658, 90, 10, rl.Gray)
		rl.DrawText("generate map 3d model", 658, 104, 10, rl.Gray)

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	rl.UnloadTexture(cubicmap) // Unload cubicmap texture
	rl.UnloadTexture(texture)  // Unload map texture
	rl.UnloadModel(model)      // Unload map model

	rl.CloseWindow()
}
