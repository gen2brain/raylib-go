/*******************************************************************************************
*
*   raylib [models] example - rlgl module usage with push/pop matrix transformations
*
*   NOTE: This example uses [rlgl] module functionality (pseudo-OpenGL 1.1 style coding)
*
*   Example originally created with raylib 2.5, last time updated with raylib 4.0
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2018-2024 Ramon Santamaria (@raysan5)
*
********************************************************************************************/
package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450

	sunRadius   = 4.0
	earthRadius = 0.6
	moonRadius  = 0.16

	earthOrbitRadius = 8.0
	moonOrbitRadius  = 1.5

	rings, slices = 16, 16
)

func main() {
	// Initialization
	title := "raylib [models] example - rlgl module usage with push/pop matrix transformations"
	rl.InitWindow(screenWidth, screenHeight, title)

	// Define the camera to look into our 3d world
	camera := rl.Camera{
		Position: rl.Vector3{
			X: 16.0,
			Y: 16.0,
			Z: 16.0,
		}, // Camera position
		Target:     rl.Vector3{},         // Camera looking at point
		Up:         rl.Vector3{Y: 1.0},   // Camera up vector (rotation towards target)
		Fovy:       45.0,                 // Camera field-of-view Y
		Projection: rl.CameraPerspective, // Camera projection type
	}

	var rotationSpeed float32 = 0.2 // General system rotation speed
	var earthRotation float32       // Rotation of earth around itself (days) in degrees
	var earthOrbitRotation float32  // Rotation of earth around the Sun (years) in degrees
	var moonRotation float32        // Rotation of moon around itself
	var moonOrbitRotation float32   // Rotation of moon around earth in degrees

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		rl.UpdateCamera(&camera, rl.CameraOrbital)

		earthRotation += 5.0 * rotationSpeed
		earthOrbitRotation += 365 / 360.0 * (5.0 * rotationSpeed) * rotationSpeed
		moonRotation += 2.0 * rotationSpeed
		moonOrbitRotation += 8.0 * rotationSpeed

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode3D(camera)

		rl.PushMatrix()
		rl.Scalef(sunRadius, sunRadius, sunRadius) // Scale Sun
		DrawSphereBasic(rl.Gold)                   // Draw the Sun
		rl.PopMatrix()

		rl.PushMatrix()
		rl.Rotatef(earthOrbitRotation, 0.0, 1.0, 0.0) // Rotation for Earth orbit around Sun
		rl.Translatef(earthOrbitRadius, 0.0, 0.0)     // Translation for Earth orbit

		rl.PushMatrix()
		rl.Rotatef(earthRotation, 0.25, 1.0, 0.0)        // Rotation for Earth itself
		rl.Scalef(earthRadius, earthRadius, earthRadius) // Scale Earth

		DrawSphereBasic(rl.Blue) // Draw the Earth
		rl.PopMatrix()

		rl.Rotatef(moonOrbitRotation, 0.0, 1.0, 0.0)  // Rotation for Moon orbit around Earth
		rl.Translatef(moonOrbitRadius, 0.0, 0.0)      // Translation for Moon orbit
		rl.Rotatef(moonRotation, 0.0, 1.0, 0.0)       // Rotation for Moon itself
		rl.Scalef(moonRadius, moonRadius, moonRadius) // Scale Moon

		DrawSphereBasic(rl.LightGray) // Draw the Moon
		rl.PopMatrix()

		// Some reference elements (not affected by previous matrix transformations)
		rl.DrawCircle3D(rl.Vector3{}, earthOrbitRadius, rl.NewVector3(1, 0, 0), 90.0, rl.Fade(rl.Red, 0.5))
		rl.DrawGrid(20, 1.0)

		rl.EndMode3D()

		rl.DrawText("EARTH ORBITING AROUND THE SUN!", 400, 10, 20, rl.Maroon)
		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.CloseWindow() // Close window and OpenGL context
}

// DrawSphereBasic draws a sphere without any matrix transformation
// NOTE: Sphere is drawn in world position ( 0, 0, 0 ) with radius 1.0f
func DrawSphereBasic(color rl.Color) {
	// Make sure there is enough space in the internal render batch
	// buffer to store all required vertex, batch is reset if required
	rl.CheckRenderBatchLimit((rings + 2) * slices * 6)

	rl.Begin(rl.Triangles)
	rl.Color4ub(color.R, color.G, color.B, color.A)

	for ring := int32(0); ring < (rings + 2); ring++ {
		for slice := int32(0); slice < slices; slice++ {
			rl.Vertex3f(getCoords(ring, slice))
			rl.Vertex3f(getCoords(ring+1, slice+1))
			rl.Vertex3f(getCoords(ring+1, slice))
			rl.Vertex3f(getCoords(ring, slice))
			rl.Vertex3f(getCoords(ring, slice+1))
			rl.Vertex3f(getCoords(ring+1, slice+1))
		}
	}
	rl.End()
}

func getCoords(ring, slice int32) (x, y, z float32) {
	ringF := float64(ring)
	sliceF := float64(slice)

	// Calculate angels
	alpha := rl.Deg2rad * (270 + (180/(float64(rings)+1))*ringF)
	beta := rl.Deg2rad * (sliceF * 360 / float64(slices))

	// Calculate coords
	x = float32(math.Cos(alpha) * math.Sin(beta))
	y = float32(math.Sin(alpha))
	z = float32(math.Cos(alpha) * math.Cos(beta))

	return x, y, z
}
