package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagMsaa4xHint) // Enable Multi Sampling Anti Aliasing 4x (if available)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - custom uniform variable")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(8.0, 8.0, 8.0)
	camera.Target = rl.NewVector3(0.0, 1.5, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	obj := rl.LoadModel("barracks.obj")               // Load OBJ model
	texture := rl.LoadTexture("barracks_diffuse.png") // Load model texture

	rl.SetMaterialTexture(obj.Materials, rl.MapDiffuse, texture) // Set obj model diffuse texture

	position := rl.NewVector3(0.0, 0.0, 0.0) // Set model position

	shader := rl.LoadShader("", "glsl330/swirl.fs")

	// Get variable (uniform) location on the shader to connect with the program
	// NOTE: If uniform variable could not be found in the shader, function returns -1
	swirlCenterLoc := rl.GetShaderLocation(shader, "center")

	swirlCenter := make([]float32, 2)
	swirlCenter[0] = float32(screenWidth) / 2
	swirlCenter[1] = float32(screenHeight) / 2

	// Create a RenderTexture2D to be used for render to texture
	target := rl.LoadRenderTexture(screenWidth, screenHeight)

	// Setup orbital camera
	rl.SetCameraMode(camera, rl.CameraOrbital) // Set an orbital camera mode

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update
		//----------------------------------------------------------------------------------
		mousePosition := rl.GetMousePosition()

		swirlCenter[0] = mousePosition.X
		swirlCenter[1] = float32(screenHeight) - mousePosition.Y

		// Send new value to the shader to be used on drawing
		rl.SetShaderValue(shader, swirlCenterLoc, swirlCenter, rl.ShaderUniformVec2)

		rl.UpdateCamera(&camera) // Update camera

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginTextureMode(target) // Enable drawing to texture

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(obj, position, 0.5, rl.White) // Draw 3d model with texture

		rl.DrawGrid(10, 1.0) // Draw a grid

		rl.EndMode3D()

		rl.DrawText("TEXT DRAWN IN RENDER TEXTURE", 200, 10, 30, rl.Red)

		rl.EndTextureMode() // End drawing to texture (now we have a texture available for next passes)

		rl.BeginShaderMode(shader)

		// NOTE: Render texture must be y-flipped due to default OpenGL coordinates (left-bottom)
		rl.DrawTextureRec(target.Texture, rl.NewRectangle(0, 0, float32(target.Texture.Width), float32(-target.Texture.Height)), rl.NewVector2(0, 0), rl.White)

		rl.EndShaderMode()

		rl.DrawText("(c) Barracks 3D model by Alberto Cano", screenWidth-200, screenHeight-20, 10, rl.Gray)

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	rl.UnloadShader(shader)        // Unload shader
	rl.UnloadTexture(texture)      // Unload texture
	rl.UnloadModel(obj)            // Unload model
	rl.UnloadRenderTexture(target) // Unload render texture

	rl.CloseWindow()
}
