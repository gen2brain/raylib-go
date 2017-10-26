package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.SetConfigFlags(raylib.FlagMsaa4xHint | raylib.FlagVsyncHint) // Enable Multi Sampling Anti Aliasing 4x (if available)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - postprocessing shader")

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(3.0, 3.0, 3.0)
	camera.Target = raylib.NewVector3(0.0, 1.5, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0

	dwarf := raylib.LoadModel("dwarf.obj")             // Load OBJ model
	texture := raylib.LoadTexture("dwarf_diffuse.png") // Load model texture

	dwarf.Material.Maps[raylib.MapDiffuse].Texture = texture // Set dwarf model diffuse texture

	position := raylib.NewVector3(0.0, 0.0, 0.0) // Set model position

	shader := raylib.LoadShader("glsl330/base.vs", "glsl330/bloom.fs") // Load postpro shader

	// Create a RenderTexture2D to be used for render to texture
	target := raylib.LoadRenderTexture(screenWidth, screenHeight)

	raylib.SetCameraMode(camera, raylib.CameraOrbital) // Set free camera mode

	//raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera) // Update camera

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginTextureMode(target) // Enable drawing to texture

		raylib.Begin3dMode(camera)

		raylib.DrawModel(dwarf, position, 2.0, raylib.White) // Draw 3d model with texture

		raylib.DrawGrid(10, 1.0) // Draw a grid

		raylib.End3dMode()

		raylib.DrawText("HELLO POSTPROCESSING!", 70, 190, 50, raylib.Red)

		raylib.EndTextureMode() // End drawing to texture (now we have a texture available for next passes)

		raylib.BeginShaderMode(shader)

		// NOTE: Render texture must be y-flipped due to default OpenGL coordinates (left-bottom)
		raylib.DrawTextureRec(target.Texture, raylib.NewRectangle(0, 0, target.Texture.Width, -target.Texture.Height), raylib.NewVector2(0, 0), raylib.White)

		raylib.EndShaderMode()

		raylib.DrawText("(c) Dwarf 3D model by David Moreno", screenWidth-200, screenHeight-20, 10, raylib.Gray)

		raylib.DrawFPS(10, 10)

		raylib.EndDrawing()
	}

	raylib.UnloadShader(shader)        // Unload shader
	raylib.UnloadTexture(texture)      // Unload texture
	raylib.UnloadModel(dwarf)          // Unload model
	raylib.UnloadRenderTexture(target) // Unload render texture

	raylib.CloseWindow()
}
