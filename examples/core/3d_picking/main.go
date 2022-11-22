package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - 3d picking")

	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(10.0, 10.0, 10.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	cubePosition := rl.NewVector3(0.0, 1.0, 0.0)
	cubeSize := rl.NewVector3(2.0, 2.0, 2.0)

	var ray rl.Ray
	var collision rl.RayCollision

	rl.SetCameraMode(camera, rl.CameraFree) // Set a free camera mode

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera) // Update camera

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if !collision.Hit {
				ray = rl.GetMouseRay(rl.GetMousePosition(), camera)

				// Check collision between ray and box
				min := rl.NewVector3(cubePosition.X-cubeSize.X/2, cubePosition.Y-cubeSize.Y/2, cubePosition.Z-cubeSize.Z/2)
				max := rl.NewVector3(cubePosition.X+cubeSize.X/2, cubePosition.Y+cubeSize.Y/2, cubePosition.Z+cubeSize.Z/2)
				collision = rl.GetRayCollisionBox(ray, rl.NewBoundingBox(min, max))
			} else {
				collision.Hit = false
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		if collision.Hit {
			rl.DrawCube(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, rl.Red)
			rl.DrawCubeWires(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, rl.Maroon)

			rl.DrawCubeWires(cubePosition, cubeSize.X+0.2, cubeSize.Y+0.2, cubeSize.Z+0.2, rl.Green)
		} else {
			rl.DrawCube(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, rl.Gray)
			rl.DrawCubeWires(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, rl.DarkGray)
		}

		rl.DrawRay(ray, rl.Maroon)

		rl.DrawGrid(10, 1.0)

		rl.EndMode3D()

		rl.DrawText("Try selecting the box with mouse!", 240, 10, 20, rl.DarkGray)

		if collision.Hit {
			rl.DrawText("BOX SELECTED", (screenWidth-rl.MeasureText("BOX SELECTED", 30))/2, int32(float32(screenHeight)*0.1), 30, rl.Green)
		}

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
