/*******************************************************************************************
*
*   raylib [models] example - Plane rotations (yaw, pitch, roll)
*
*   Example originally created with raylib 1.8, last time updated with raylib 4.0
*
*   Example contributed by Berni (@Berni8k) and reviewed by Ramon Santamaria (@raysan5)
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2017-2024 Berni (@Berni8k) and Ramon Santamaria (@raysan5)
*
********************************************************************************************/

package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	screenWidth  = 800
	screenHeight = 450
)

func main() {
	//SetConfigFlags(FLAG_MSAA_4X_HINT | FLAG_WINDOW_HIGHDPI)
	title := "raylib [models] example - plane rotations (yaw, pitch, roll)"
	rl.InitWindow(screenWidth, screenHeight, title)

	camera := rl.Camera{
		Position: rl.Vector3{
			Y: 50.0,
			Z: -120.0,
		}, // Camera position perspective
		Target:     rl.Vector3{},         // Camera looking at point
		Up:         rl.Vector3{Y: 1.0},   // Camera up vector (rotation towards target)
		Fovy:       30.0,                 // Camera field-of-view Y
		Projection: rl.CameraPerspective, // Camera type
	}
	model := rl.LoadModel("plane.obj")                             // Load model
	texture := rl.LoadTexture("plane_diffuse.png")                 // Load model texture
	rl.SetMaterialTexture(model.Materials, rl.MapDiffuse, texture) // Set map diffuse texture

	var pitch, roll, yaw float32

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update

		// Plane pitch (x-axis) controls
		pitch = controlPlane(pitch, 0.6, rl.KeyDown, rl.KeyUp)
		// Plane yaw (y-axis) controls
		roll = controlPlane(roll, 1.0, rl.KeyLeft, rl.KeyRight)
		// Plane roll (z-axis) controls
		yaw = controlPlane(yaw, 1.0, rl.KeyA, rl.KeyS)

		// Transformation matrix for rotations
		rotationV := rl.Vector3{
			X: rl.Deg2rad * pitch,
			Y: rl.Deg2rad * yaw,
			Z: rl.Deg2rad * roll,
		}
		model.Transform = rl.MatrixRotateXYZ(rotationV)

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		// Draw 3D model (recommended to draw 3D always before 2D)
		rl.BeginMode3D(camera)

		// Draw 3d model with texture
		rl.DrawModel(model, rl.Vector3{Y: -8.0}, 1.0, rl.White)
		rl.DrawGrid(10, 10.0)

		rl.EndMode3D()

		// Draw controls info
		rl.DrawRectangle(30, 370, 260, 70, rl.Fade(rl.Green, 0.5))
		rl.DrawRectangleLines(30, 370, 260, 70, rl.Fade(rl.DarkGreen, 0.5))
		rl.DrawText("Pitch controlled with: KEY_UP / KEY_DOWN", 40, 380, 10, rl.DarkGray)
		rl.DrawText("Roll controlled with: KEY_LEFT / KEY_RIGHT", 40, 400, 10, rl.DarkGray)
		rl.DrawText("Yaw controlled with: KEY_A / KEY_S", 40, 420, 10, rl.DarkGray)

		rl.DrawText("(c) WWI Plane Model created by GiaHanLam", screenWidth-240, screenHeight-20, 10, rl.DarkGray)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadModel(model)     // Unload model data
	rl.UnloadTexture(texture) // Unload texture data

	rl.CloseWindow() // Close window and OpenGL context
}

func controlPlane(ctrl, value float32, key1, key2 int32) float32 {
	if rl.IsKeyDown(key1) {
		ctrl -= value
	} else if rl.IsKeyDown(key2) {
		ctrl += value
	} else {
		if ctrl > 0.0 {
			ctrl -= value / 2
		} else if ctrl < 0.0 {
			ctrl += value / 2
		}
	}
	return ctrl
}
