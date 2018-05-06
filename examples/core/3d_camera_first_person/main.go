package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const (
	maxColumns = 20
)

func main() {
	raylib.InitWindow(800, 450, "raylib [core] example - 3d camera first person")

	camera := raylib.Camera3D{}
	camera.Position = raylib.NewVector3(4.0, 2.0, 4.0)
	camera.Target = raylib.NewVector3(0.0, 1.8, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 60.0
	camera.Type = raylib.CameraPerspective

	// Generates some random columns
	heights := make([]float32, maxColumns)
	positions := make([]raylib.Vector3, maxColumns)
	colors := make([]raylib.Color, maxColumns)

	for i := 0; i < maxColumns; i++ {
		heights[i] = float32(raylib.GetRandomValue(1, 12))
		positions[i] = raylib.NewVector3(float32(raylib.GetRandomValue(-15, 15)), heights[i]/2, float32(raylib.GetRandomValue(-15, 15)))
		colors[i] = raylib.NewColor(uint8(raylib.GetRandomValue(20, 255)), uint8(raylib.GetRandomValue(10, 55)), 30, 255)
	}

	raylib.SetCameraMode(camera, raylib.CameraFirstPerson) // Set a first person camera mode

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera) // Update camera

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode3D(camera)

		raylib.DrawPlane(raylib.NewVector3(0.0, 0.0, 0.0), raylib.NewVector2(32.0, 32.0), raylib.LightGray) // Draw ground
		raylib.DrawCube(raylib.NewVector3(-16.0, 2.5, 0.0), 1.0, 5.0, 32.0, raylib.Blue)                    // Draw a blue wall
		raylib.DrawCube(raylib.NewVector3(16.0, 2.5, 0.0), 1.0, 5.0, 32.0, raylib.Lime)                     // Draw a green wall
		raylib.DrawCube(raylib.NewVector3(0.0, 2.5, 16.0), 32.0, 5.0, 1.0, raylib.Gold)                     // Draw a yellow wall

		// Draw some cubes around
		for i := 0; i < maxColumns; i++ {
			raylib.DrawCube(positions[i], 2.0, heights[i], 2.0, colors[i])
			raylib.DrawCubeWires(positions[i], 2.0, heights[i], 2.0, raylib.Maroon)
		}

		raylib.EndMode3D()

		raylib.DrawRectangle(10, 10, 220, 70, raylib.Fade(raylib.SkyBlue, 0.5))
		raylib.DrawRectangleLines(10, 10, 220, 70, raylib.Blue)

		raylib.DrawText("First person camera default controls:", 20, 20, 10, raylib.Black)
		raylib.DrawText("- Move with keys: W, A, S, D", 40, 40, 10, raylib.DarkGray)
		raylib.DrawText("- Mouse move to look around", 40, 60, 10, raylib.DarkGray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
