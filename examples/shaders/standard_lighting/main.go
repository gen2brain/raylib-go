package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.SetConfigFlags(raylib.FlagMsaa4xHint | raylib.FlagVsyncHint) // Enable Multi Sampling Anti Aliasing 4x (if available)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - model shader")

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(4.0, 4.0, 4.0)
	camera.Target = raylib.NewVector3(0.0, 1.5, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0

	position := raylib.NewVector3(0.0, 0.0, 0.0) // Set model position

	dwarf := raylib.LoadModel("dwarf.obj") // Load OBJ model

	material := raylib.LoadStandardMaterial()

	material.TexDiffuse = raylib.LoadTexture("dwarf_diffuse.png")   // Load model diffuse texture
	material.TexNormal = raylib.LoadTexture("dwarf_normal.png")     // Load model normal texture
	material.TexSpecular = raylib.LoadTexture("dwarf_specular.png") // Load model specular texture
	material.ColDiffuse = raylib.White
	material.ColAmbient = raylib.NewColor(0, 0, 10, 255)
	material.ColSpecular = raylib.White
	material.Glossiness = 50.0

	dwarf.Material = material // Apply material to model

	spotLight := raylib.CreateLight(raylib.LightSpot, raylib.NewVector3(3.0, 5.0, 2.0), raylib.NewColor(255, 255, 255, 255))
	spotLight.Target = raylib.NewVector3(0.0, 0.0, 0.0)
	spotLight.Intensity = 2.0
	spotLight.Diffuse = raylib.NewColor(255, 100, 100, 255)
	spotLight.ConeAngle = 60.0

	dirLight := raylib.CreateLight(raylib.LightDirectional, raylib.NewVector3(0.0, -3.0, -3.0), raylib.NewColor(255, 255, 255, 255))
	dirLight.Target = raylib.NewVector3(1.0, -2.0, -2.0)
	dirLight.Intensity = 2.0
	dirLight.Diffuse = raylib.NewColor(100, 255, 100, 255)

	pointLight := raylib.CreateLight(raylib.LightPoint, raylib.NewVector3(0.0, 4.0, 5.0), raylib.NewColor(255, 255, 255, 255))
	pointLight.Intensity = 2.0
	pointLight.Diffuse = raylib.NewColor(100, 100, 255, 255)
	pointLight.Radius = 3.0

	// Setup orbital camera
	raylib.SetCameraMode(camera, raylib.CameraOrbital) // Set an orbital camera mode

	//raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera) // Update camera

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.Begin3dMode(camera)

		raylib.DrawModel(dwarf, position, 2.0, raylib.White) // Draw 3d model with texture

		raylib.DrawLight(spotLight)  // Draw spot light
		raylib.DrawLight(dirLight)   // Draw directional light
		raylib.DrawLight(pointLight) // Draw point light

		raylib.DrawGrid(10, 1.0) // Draw a grid

		raylib.End3dMode()

		raylib.DrawText("(c) Dwarf 3D model by David Moreno", screenWidth-200, screenHeight-20, 10, raylib.Gray)

		raylib.DrawFPS(10, 10)

		raylib.EndDrawing()
	}

	raylib.UnloadMaterial(material) // Unload material and assigned textures
	raylib.UnloadModel(dwarf)       // Unload model

	// Destroy all created lights
	raylib.DestroyLight(pointLight)
	raylib.DestroyLight(dirLight)
	raylib.DestroyLight(spotLight)

	raylib.CloseWindow()
}
