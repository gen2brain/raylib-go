/*******************************************************************************************
*
*   raylib [shaders] example - Vertex displacement
*
*   Example complexity rating: [★★★☆] 3/4
*
*   Example originally created with raylib 5.0, last time updated with raylib 5.5
*
*   Example originally contributed by Alex ZH (@ZzzhHe) and reviewed by Ramon Santamaria (@raysan5)
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2023-2025 Alex ZH (@ZzzhHe)
*
********************************************************************************************/

package main

import (
	"fmt"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// NOTE: Usage: `PLATFORM_DESKTOP=1 go run .`
var glslVersion int

func init() {
	if v, ok := os.LookupEnv("PLATFORM_DESKTOP"); ok && v == "1" {
		glslVersion = 330
	} else { // PLATFORM_ANDROID, PLATFORM_WEB
		glslVersion = 100
	}
}

// ------------------------------------------------------------------------------------
// Program main entry point
// ------------------------------------------------------------------------------------
func main() {
	// Initialization
	//--------------------------------------------------------------------------------------
	const screenWidth int32 = 800
	const screenHeight int32 = 450

	rl.SetConfigFlags(rl.FlagMsaa4xHint) // Enable Multi Sampling Anti Aliasing 4x (if available)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - vertex displacement")
	defer rl.CloseWindow() // Close window and OpenGL context

	// Set up camera
	camera := rl.Camera{
		Position:   rl.NewVector3(20.0, 5.0, -20.0),
		Target:     rl.Vector3Zero(),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       60.0,
		Projection: rl.CameraPerspective,
	}

	// Load vertex and fragment shaders
	shaderDir := fmt.Sprintf("glsl%d", glslVersion)

	shader := rl.LoadShader(
		filepath.Join(shaderDir, "vertex_displacement.vs"),
		filepath.Join(shaderDir, "vertex_displacement.fs"),
	)
	defer rl.UnloadShader(shader)

	timeLoc := rl.GetShaderLocation(shader, "time")

	// Load perlin noise texture
	perlinNoiseImage := rl.GenImagePerlinNoise(512, 512, 0, 0, 1.0)
	perlinNoiseMap := rl.LoadTextureFromImage(perlinNoiseImage)
	defer rl.UnloadTexture(perlinNoiseMap)
	rl.UnloadImage(perlinNoiseImage)

	// Set shader uniform location
	perlinNoiseMapLoc := rl.GetShaderLocation(shader, "perlinNoiseMap")
	rl.EnableShader(shader.ID)
	rl.ActiveTextureSlot(1)
	rl.EnableTexture(perlinNoiseMap.ID)
	rl.SetUniformSampler(perlinNoiseMapLoc, 1)

	// Create a plane mesh and model
	planeMesh := rl.GenMeshPlane(50, 50, 50, 50)
	planeModel := rl.LoadModelFromMesh(planeMesh)
	defer rl.UnloadModel(planeModel)

	// Set plane model material
	planeModel.Materials.Shader = shader

	timer := float32(0)

	rl.DisableCursor()

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second
	//--------------------------------------------------------------------------------------

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		//----------------------------------------------------------------------------------
		rl.UpdateCamera(&camera, rl.CameraFree) // Update camera with free camera mode

		deltaTime := rl.GetFrameTime()
		timer += deltaTime
		timeValue := []float32{timer}

		rl.SetShaderValue(shader, timeLoc, timeValue, rl.ShaderUniformFloat) // Send time value to shader
		//----------------------------------------------------------------------------------

		// Draw
		//----------------------------------------------------------------------------------
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.BeginShaderMode(shader)

		// Draw plane model
		rl.DrawModel(planeModel, rl.Vector3Zero(), 1.0, rl.White)

		rl.EndShaderMode()

		rl.EndMode3D()

		rl.DrawText("Vertex displacement", 10, 10, 20, rl.DarkGray)

		rl.DrawFPS(10, 40)

		rl.EndDrawing()
		//----------------------------------------------------------------------------------
	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	// NOTE: Unload all loaded resources at this point (that are not `defer`-ed)
	//--------------------------------------------------------------------------------------
}
