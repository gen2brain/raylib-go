package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - directional billboard")
	defer rl.CloseWindow()

	// Set up the camera
	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(2.0, 1.0, 2.0) // Starting position
	camera.Target = rl.NewVector3(0.0, 0.5, 0.0)   // Target position
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)       // Up vector
	camera.Fovy = 45.0                             // FOV
	camera.Projection = rl.CameraPerspective       // Projection type

	// Load billboard texture directly from root folder
	skillbot := rl.LoadTexture("skillbot.png")
	// Timer to update animation
	animTimer := float32(0.0)
	// Animation frame
	anim := uint32(0)

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		rl.UpdateCamera(&camera, rl.CameraOrbital)

		// Update timer with delta time
		animTimer += rl.GetFrameTime()

		// Update frame index after a certain amount of time (half a second)
		if animTimer > 0.5 {
			animTimer = 0.0
			anim++
		}

		// Reset frame index to zero on overflow
		if anim >= 4 {
			anim = 0
		}

		// Find current direction frame based on camera position relative to billboard object
		angle := rl.Vector2Angle(rl.NewVector2(2.0, 0.0), rl.NewVector2(camera.Position.X, camera.Position.Z))
		dir := float32(math.Floor(float64(((angle / float32(math.Pi)) * 4.0) + 0.25)))

		// Correct frame index if angle is negative
		if dir < 0.0 {
			dir = 8.0 - float32(math.Abs(float64(int32(dir))))
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawGrid(10, 1.0)

		// Draw billboard pointing straight up to sky, rotated relative to camera
		sourceRec := rl.NewRectangle(0.0+(float32(anim)*24.0), 0.0+(dir*24.0), 24.0, 24.0)
		rl.DrawBillboardPro(
			camera,
			skillbot,
			sourceRec,
			rl.Vector3Zero(),
			rl.NewVector3(0.0, 1.0, 0.0),
			rl.Vector2One(),
			rl.NewVector2(0.5, 0.0),
			0.0,
			rl.White,
		)

		rl.EndMode3D()

		// Render various variables for reference
		rl.DrawText(fmt.Sprintf("animation: %d", anim), 10, 10, 20, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("direction frame: %.0f", dir), 10, 40, 20, rl.DarkGray)

		rl.EndDrawing()
	}
	rl.UnloadTexture(skillbot)
	rl.CloseWindow()
}
