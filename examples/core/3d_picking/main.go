package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [core] example - 3d picking")

	camera := raylib.Camera3D{}
	camera.Position = raylib.NewVector3(10.0, 10.0, 10.0)
	camera.Target = raylib.NewVector3(0.0, 0.0, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Type = raylib.CameraPerspective

	cubePosition := raylib.NewVector3(0.0, 1.0, 0.0)
	cubeSize := raylib.NewVector3(2.0, 2.0, 2.0)

	var ray raylib.Ray

	collision := false

	raylib.SetCameraMode(camera, raylib.CameraFree) // Set a free camera mode

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera) // Update camera

		if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
			// NOTE: This function is NOT WORKING properly!
			ray = raylib.GetMouseRay(raylib.GetMousePosition(), camera)

			// Check collision between ray and box
			min := raylib.NewVector3(cubePosition.X-cubeSize.X/2, cubePosition.Y-cubeSize.Y/2, cubePosition.Z-cubeSize.Z/2)
			max := raylib.NewVector3(cubePosition.X+cubeSize.X/2, cubePosition.Y+cubeSize.Y/2, cubePosition.Z+cubeSize.Z/2)
			collision = raylib.CheckCollisionRayBox(ray, raylib.NewBoundingBox(min, max))
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode3D(camera)

		if collision {
			raylib.DrawCube(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, raylib.Red)
			raylib.DrawCubeWires(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, raylib.Maroon)

			raylib.DrawCubeWires(cubePosition, cubeSize.X+0.2, cubeSize.Y+0.2, cubeSize.Z+0.2, raylib.Green)
		} else {
			raylib.DrawCube(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, raylib.Gray)
			raylib.DrawCubeWires(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, raylib.DarkGray)
		}

		raylib.DrawRay(ray, raylib.Maroon)

		raylib.DrawGrid(10, 1.0)

		raylib.EndMode3D()

		raylib.DrawText("Try selecting the box with mouse!", 240, 10, 20, raylib.DarkGray)

		if collision {
			raylib.DrawText("BOX SELECTED", (screenWidth-raylib.MeasureText("BOX SELECTED", 30))/2, int32(float32(screenHeight)*0.1), 30, raylib.Green)
		}

		raylib.DrawFPS(10, 10)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
