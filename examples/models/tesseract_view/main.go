package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - tesseract view")
	defer rl.CloseWindow()

	// Define the camera to look into our 3D world
	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(4.0, 4.0, 4.0) // Camera position
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)   // Camera looking at point
	camera.Up = rl.NewVector3(0.0, 0.0, 1.0)       // Camera up vector
	camera.Fovy = 50.0                             // Camera field-of-view Y
	camera.Projection = rl.CameraPerspective       // Camera mode type

	// Find the coordinates by setting XYZW to +-1
	tesseract := [16]rl.Vector4{
		rl.NewVector4(1, 1, 1, 1), rl.NewVector4(1, 1, 1, -1),
		rl.NewVector4(1, 1, -1, 1), rl.NewVector4(1, 1, -1, -1),
		rl.NewVector4(1, -1, 1, 1), rl.NewVector4(1, -1, 1, -1),
		rl.NewVector4(1, -1, -1, 1), rl.NewVector4(1, -1, -1, -1),
		rl.NewVector4(-1, 1, 1, 1), rl.NewVector4(-1, 1, 1, -1),
		rl.NewVector4(-1, 1, -1, 1), rl.NewVector4(-1, 1, -1, -1),
		rl.NewVector4(-1, -1, 1, 1), rl.NewVector4(-1, -1, 1, -1),
		rl.NewVector4(-1, -1, -1, 1), rl.NewVector4(-1, -1, -1, -1),
	}

	var transformed [16]rl.Vector3
	var wValues [16]float32

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		rotation := rl.Deg2rad * 45.0 * float32(rl.GetTime())

		for i := 0; i < 16; i++ {
			p := tesseract[i]

			// Rotate the XW part of the vector
			rotXW := rl.Vector2Rotate(rl.NewVector2(p.X, p.W), rotation)
			p.X = rotXW.X
			p.W = rotXW.Y

			// Projection from XYZW to XYZ from perspective point (0, 0, 0, 3)
			c := float32(3.0) / (3.0 - p.W)
			p.X = c * p.X
			p.Y = c * p.Y
			p.Z = c * p.Z

			// Split XYZ coordinate and W values later for drawing
			transformed[i] = rl.NewVector3(p.X, p.Y, p.Z)
			wValues[i] = p.W
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)
		for i := 0; i < 16; i++ {
			// Draw spheres to indicate the W value
			rl.DrawSphere(transformed[i], float32(math.Abs(float64(wValues[i]*0.1))), rl.Red)

			for j := 0; j < 16; j++ {
				// Two lines are connected if they differ by 1 coordinate
				// This way we don't have to keep an edge list
				v1 := tesseract[i]
				v2 := tesseract[j]

				diff := 0
				if v1.X == v2.X {
					diff++
				}
				if v1.Y == v2.Y {
					diff++
				}
				if v1.Z == v2.Z {
					diff++
				}
				if v1.W == v2.W {
					diff++
				}

				// Draw only differing by 1 coordinate and lower index only (duplicate lines)
				if diff == 3 && i < j {
					rl.DrawLine3D(transformed[i], transformed[j], rl.Maroon)
				}
			}
		}
		rl.EndMode3D()

		rl.EndDrawing()
	}
	rl.CloseWindow()
}
