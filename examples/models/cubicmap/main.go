package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [models] example - cubesmap loading and drawing")

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(16.0, 14.0, 16.0)
	camera.Target = raylib.NewVector3(0.0, 0.0, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Type = raylib.CameraPerspective

	image := raylib.LoadImage("cubicmap.png")      // Load cubicmap image (RAM)
	cubicmap := raylib.LoadTextureFromImage(image) // Convert image to texture to display (VRAM)

	mesh := raylib.GenMeshCubicmap(*image, raylib.NewVector3(1.0, 1.0, 1.0))
	model := raylib.LoadModelFromMesh(mesh)

	// NOTE: By default each cube is mapped to one part of texture atlas
	texture := raylib.LoadTexture("cubicmap_atlas.png")      // Load map texture
	model.Material.Maps[raylib.MapDiffuse].Texture = texture // Set map diffuse texture

	mapPosition := raylib.NewVector3(-16.0, 0.0, -8.0) // Set model position

	raylib.UnloadImage(image) // Unload cubicmap image from RAM, already uploaded to VRAM

	raylib.SetCameraMode(camera, raylib.CameraOrbital) // Set an orbital camera mode

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update

		raylib.UpdateCamera(&camera) // Update camera

		// Draw

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode3D(camera)

		raylib.DrawModel(model, mapPosition, 1.0, raylib.White)

		raylib.EndMode3D()

		raylib.DrawTextureEx(cubicmap, raylib.NewVector2(float32(screenWidth-cubicmap.Width*4-20), 20), 0.0, 4.0, raylib.White)
		raylib.DrawRectangleLines(screenWidth-cubicmap.Width*4-20, 20, cubicmap.Width*4, cubicmap.Height*4, raylib.Green)

		raylib.DrawText("cubicmap image used to", 658, 90, 10, raylib.Gray)
		raylib.DrawText("generate map 3d model", 658, 104, 10, raylib.Gray)

		raylib.DrawFPS(10, 10)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(cubicmap) // Unload cubicmap texture
	raylib.UnloadTexture(texture)  // Unload map texture
	raylib.UnloadModel(model)      // Unload map model

	raylib.CloseWindow()
}
