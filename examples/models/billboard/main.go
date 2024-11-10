/*******************************************************************************************
*
*   raylib [models] example - Drawing billboards
*
*   Example originally created with raylib 1.3, last time updated with raylib 3.5
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2015-2024 Ramon Santamaria (@raysan5)
*
********************************************************************************************/package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - drawing billboards")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(5.0, 4.0, 5.0)
	camera.Target = rl.NewVector3(0.0, 2.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	bill := rl.LoadTexture("billboard.png")              // Our texture billboard
	billPositionStatic := rl.NewVector3(0.0, 2.0, 0.0)   // Position of static billboard
	billPositionRotating := rl.NewVector3(1.0, 2.0, 1.0) // Position of rotating billboard

	// Entire billboard texture, source is used to take a segment from a larger texture.
	source := rl.Rectangle{
		Width:  float32(bill.Width),
		Height: float32(bill.Height),
	}

	// NOTE: Billboard locked on axis-Y
	billUp := rl.Vector3{Y: 1.0}

	// Set the height of the rotating billboard to 1.0 with the aspect ratio fixed
	size := rl.Vector2{
		X: source.Width / source.Height,
		Y: 1.0,
	}

	// Rotate around origin
	// Here we choose to rotate around the image center
	origin := rl.Vector2Scale(size, 0.5)

	// Distance is needed for the correct billboard draw order
	// Larger distance (further away from the camera) should be drawn prior to smaller distance.
	var distanceStatic, distanceRotating, rotation float32

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update
		rl.UpdateCamera(&camera, rl.CameraOrbital) // Update camera with orbital camera mode

		rotation += 0.4
		distanceStatic = rl.Vector3Distance(camera.Position, billPositionStatic)
		distanceRotating = rl.Vector3Distance(camera.Position, billPositionRotating)

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode3D(camera)

		rl.DrawGrid(10, 1.0) // Draw a grid

		// Draw order matters!
		if distanceStatic > distanceRotating {
			rl.DrawBillboard(camera, bill, billPositionStatic, 2.0, rl.White)
			rl.DrawBillboardPro(camera, bill, source, billPositionRotating, billUp, size, origin, rotation, rl.White)
		} else {
			rl.DrawBillboardPro(camera, bill, source, billPositionRotating, billUp, size, origin, rotation, rl.White)
			rl.DrawBillboard(camera, bill, billPositionStatic, 2.0, rl.White)
		}

		rl.EndMode3D()
		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}

	rl.UnloadTexture(bill) // Unload texture

	rl.CloseWindow()
}
