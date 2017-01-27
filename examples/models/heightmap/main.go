package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [models] example - heightmap loading and drawing")

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(18.0, 16.0, 18.0)
	camera.Target = raylib.NewVector3(0.0, 0.0, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0

	image := raylib.LoadImage("heightmap.png")                        // Load heightmap image (RAM)
	texture := raylib.LoadTextureFromImage(image)                     // Convert image to texture (VRAM)
	hmap := raylib.LoadHeightmap(image, raylib.NewVector3(16, 8, 16)) // Load heightmap model with defined size
	hmap.Material.TexDiffuse = texture                                // Set map diffuse texture
	mapPosition := raylib.NewVector3(-8.0, 0.0, -8.0)                 // Set model position

	raylib.UnloadImage(image) // Unload cubesmap image from RAM, already uploaded to VRAM

	raylib.SetCameraMode(camera, raylib.CameraOrbital) // Set an orbital camera mode

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update

		raylib.UpdateCamera(&camera) // Update camera

		// Draw

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.Begin3dMode(camera)

		// NOTE: Model is scaled to 1/4 of its original size (128x128 units)
		raylib.DrawModel(hmap, mapPosition, 1.0, raylib.Red)

		raylib.DrawGrid(20, 1.0)

		raylib.End3dMode()

		raylib.DrawTexture(texture, screenWidth-texture.Width-20, 20, raylib.White)
		raylib.DrawRectangleLines(screenWidth-texture.Width-20, 20, texture.Width, texture.Height, raylib.Green)

		raylib.DrawFPS(10, 10)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(texture) // Unload map texture
	raylib.UnloadModel(hmap)      // Unload map model

	raylib.CloseWindow()
}
