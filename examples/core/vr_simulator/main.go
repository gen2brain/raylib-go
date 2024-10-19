/*******************************************************************************************
*
*   raylib [core] example - VR Simulator (Oculus Rift CV1 parameters)
*
*   Example originally created with raylib 2.5, last time updated with raylib 4.0
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2017-2024 Ramon Santamaria (@raysan5)
*
********************************************************************************************/
package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	glslVersion  = 330 // Desktop
	// glslVersion = 100 // Android, web
)

func main() {
	// NOTE: screenWidth/screenHeight should match VR device aspect ratio
	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - vr simulator")

	// VR device parameters definition
	device := rl.VrDeviceInfo{
		// Oculus Rift CV1 parameters for simulator
		HResolution:            2160,     // Horizontal resolution in pixels
		VResolution:            1200,     // Vertical resolution in pixels
		HScreenSize:            0.133793, // Horizontal size in meters
		VScreenSize:            0.0669,   // Vertical size in meters
		EyeToScreenDistance:    0.041,    // Distance between eye and display in meters
		LensSeparationDistance: 0.07,     // Lens separation distance in meters
		InterpupillaryDistance: 0.07,     // IPD (distance between pupils) in meters

		// NOTE: CV1 uses fresnel-hybrid-asymmetric lenses with specific compute shaders
		// Following parameters are just an approximation to CV1 distortion stereo rendering

		// Lens distortion constant parameters
		LensDistortionValues: [4]float32{1.0, 0.22, 0.24, 0.0},
		// Chromatic aberration correction parameters
		ChromaAbCorrection: [4]float32{0.996, -0.004, 1.014, 0.0},
	}

	// Load VR stereo config for VR device parameters (Oculus Rift CV1 parameters)
	config := rl.LoadVrStereoConfig(device)

	// Distortion shader (uses device lens distortion and chroma)
	fileName := fmt.Sprintf("distortion%d.fs", glslVersion)
	distortion := rl.LoadShader("", fileName)

	// Update distortion shader with lens and distortion-scale parameters
	rl.SetShaderValue(distortion, rl.GetShaderLocation(distortion, "leftLensCenter"),
		config.LeftLensCenter[:], rl.ShaderUniformVec2)
	rl.SetShaderValue(distortion, rl.GetShaderLocation(distortion, "rightLensCenter"),
		config.RightLensCenter[:], rl.ShaderUniformVec2)
	rl.SetShaderValue(distortion, rl.GetShaderLocation(distortion, "leftScreenCenter"),
		config.LeftScreenCenter[:], rl.ShaderUniformVec2)
	rl.SetShaderValue(distortion, rl.GetShaderLocation(distortion, "rightScreenCenter"),
		config.RightScreenCenter[:], rl.ShaderUniformVec2)

	rl.SetShaderValue(distortion, rl.GetShaderLocation(distortion, "scale"),
		config.Scale[:], rl.ShaderUniformVec2)
	rl.SetShaderValue(distortion, rl.GetShaderLocation(distortion, "scaleIn"),
		config.ScaleIn[:], rl.ShaderUniformVec2)
	rl.SetShaderValue(distortion, rl.GetShaderLocation(distortion, "deviceWarpParam"),
		device.LensDistortionValues[:], rl.ShaderUniformVec4)
	rl.SetShaderValue(distortion, rl.GetShaderLocation(distortion, "chromaAbParam"),
		device.ChromaAbCorrection[:], rl.ShaderUniformVec4)

	// Initialize frame buffer for stereo rendering
	// NOTE: Screen size should match HMD aspect ratio
	target := rl.LoadRenderTexture(device.HResolution, device.VResolution)

	// The target's height is flipped (in the source Rectangle), due to OpenGL reasons
	sourceRec := rl.Rectangle{Width: float32(target.Texture.Width), Height: float32(-target.Texture.Height)}
	destRec := rl.Rectangle{Width: float32(rl.GetScreenWidth()), Height: float32(rl.GetScreenHeight())}

	// Define the camera to look into our 3d world

	camera := rl.Camera{
		Position:   rl.Vector3{X: 5, Y: 2, Z: 5},
		Target:     rl.Vector3{Y: 2},
		Up:         rl.Vector3{Y: 1},
		Fovy:       60.0,
		Projection: rl.CameraPerspective,
	}

	cubePosition := rl.Vector3{}

	rl.DisableCursor()  // Limit cursor to relative movement inside the window
	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)

		// Draw texture
		rl.BeginTextureMode(target)
		rl.ClearBackground(rl.RayWhite)
		rl.BeginVrStereoMode(config)
		rl.BeginMode3D(camera)

		rl.DrawCube(cubePosition, 2.0, 2.0, 2.0, rl.Red)
		rl.DrawCubeWires(cubePosition, 2.0, 2.0, 2.0, rl.Maroon)
		rl.DrawGrid(40, 1.0)

		rl.EndMode3D()
		rl.EndVrStereoMode()
		rl.EndTextureMode()

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.BeginShaderMode(distortion)
		rl.DrawTexturePro(target.Texture, sourceRec, destRec, rl.Vector2{}, 0.0, rl.White)
		rl.EndShaderMode()
		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadVrStereoConfig(config) // Unload stereo config

	rl.UnloadRenderTexture(target) // Unload stereo render fbo
	rl.UnloadShader(distortion)    // Unload distortion shader

	rl.CloseWindow() // Close window and OpenGL context
}
