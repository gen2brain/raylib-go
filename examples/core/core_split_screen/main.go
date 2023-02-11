package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

/*******************************************************************************************
*
*   raylib [core] example - split screen
*
*   Example originally created with raylib 3.7, last time updated with raylib 4.0
*
*   Example contributed by Jeffery Myers (@JeffM2501) and reviewed by Ramon Santamaria (@raysan5)
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2021-2022 Jeffery Myers (@JeffM2501)
*
********************************************************************************************/

var cameraPlayer1 rl.Camera
var cameraPlayer2 rl.Camera

// Scene drawing
func DrawScene() {
	count := float32(5.0)
	spacing := float32(4.0)

	// Grid of cube trees on a plane to make a "world"
	rl.DrawPlane(rl.Vector3{0, 0, 0}, rl.Vector2{50, 50}, rl.Beige) // Simple world plane

	for x := -float32(count * spacing); x <= count*spacing; x += spacing {
		for z := -float32(count * spacing); z <= count*spacing; z += spacing {
			rl.DrawCube(rl.Vector3{x, 1.5, z}, 1, 1, 1, rl.Lime)
			rl.DrawCube(rl.Vector3{x, 0.5, z}, 0.25, 1, 0.25, rl.Brown)
		}
	}

	// Draw a cube at each player's position
	rl.DrawCube(cameraPlayer1.Position, 1, 1, 1, rl.Red)
	rl.DrawCube(cameraPlayer2.Position, 1, 1, 1, rl.Blue)
}

// ------------------------------------------------------------------------------------
// Program main entry point
// ------------------------------------------------------------------------------------
func main() {
	// Initialization
	//--------------------------------------------------------------------------------------
	const (
		screenWidth  = 800
		screenHeight = 450
	)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - split screen")

	// Setup player 1 camera and screen
	cameraPlayer1.Fovy = 45.0
	cameraPlayer1.Up.Y = 1.0
	cameraPlayer1.Target.Y = 1.0
	cameraPlayer1.Position.Z = -3.0
	cameraPlayer1.Position.Y = 1.0

	screenPlayer1 := rl.LoadRenderTexture(screenWidth/2, screenHeight)

	// Setup player two camera and screen
	cameraPlayer2.Fovy = 45.0
	cameraPlayer2.Up.Y = 1.0
	cameraPlayer2.Target.Y = 3.0
	cameraPlayer2.Position.X = -3.0
	cameraPlayer2.Position.Y = 3.0

	screenPlayer2 := rl.LoadRenderTexture(screenWidth/2, screenHeight)

	// Build a flipped rectangle the size of the split view to use for drawing later
	splitScreenRect := rl.Rectangle{0.0, 0.0, float32(screenPlayer1.Texture.Width), float32(-screenPlayer1.Texture.Height)}

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second
	//--------------------------------------------------------------------------------------

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		//----------------------------------------------------------------------------------
		// If anyone moves this frame, how far will they move based on the time since the last frame
		// this moves thigns at 10 world units per second, regardless of the actual FPS
		offsetThisFrame := 10.0 * rl.GetFrameTime()

		// Move Player1 forward and backwards (no turning)
		if rl.IsKeyDown(rl.KeyW) {
			cameraPlayer1.Position.Z += offsetThisFrame
			cameraPlayer1.Target.Z += offsetThisFrame
		} else if rl.IsKeyDown(rl.KeyS) {
			cameraPlayer1.Position.Z -= offsetThisFrame
			cameraPlayer1.Target.Z -= offsetThisFrame
		}

		// Move Player2 forward and backwards (no turning)
		if rl.IsKeyDown(rl.KeyUp) {
			cameraPlayer2.Position.X += offsetThisFrame
			cameraPlayer2.Target.X += offsetThisFrame
		} else if rl.IsKeyDown(rl.KeyDown) {
			cameraPlayer2.Position.X -= offsetThisFrame
			cameraPlayer2.Target.X -= offsetThisFrame
		}
		//----------------------------------------------------------------------------------

		// Draw
		//----------------------------------------------------------------------------------
		// Draw Player1 view to the render texture
		rl.BeginTextureMode(screenPlayer1)
		rl.ClearBackground(rl.SkyBlue)
		rl.BeginMode3D(cameraPlayer1)
		DrawScene()
		rl.EndMode3D()
		rl.DrawText("PLAYER1 W/S to move", 10, 10, 20, rl.Red)
		rl.EndTextureMode()

		// Draw Player2 view to the render texture
		rl.BeginTextureMode(screenPlayer2)
		rl.ClearBackground(rl.SkyBlue)
		rl.BeginMode3D(cameraPlayer2)
		DrawScene()
		rl.EndMode3D()
		rl.DrawText("PLAYER2 UP/DOWN to move", 10, 10, 20, rl.Blue)
		rl.EndTextureMode()

		// Draw both views render textures to the screen side by side
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawTextureRec(screenPlayer1.Texture, splitScreenRect, rl.Vector2{0, 0}, rl.White)
		rl.DrawTextureRec(screenPlayer2.Texture, splitScreenRect, rl.Vector2{screenWidth / 2.0, 0}, rl.White)
		rl.EndDrawing()
	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	rl.UnloadRenderTexture(screenPlayer1) // Unload render texture
	rl.UnloadRenderTexture(screenPlayer2) // Unload render texture

	rl.CloseWindow() // Close window and OpenGL context
	//--------------------------------------------------------------------------------------
}
