package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	raylib.InitWindow(800, 450, "raylib [core] example - 3d camera free")

	camera := raylib.Camera3D{}
	camera.Position = raylib.NewVector3(10.0, 10.0, 10.0)
	camera.Target = raylib.NewVector3(0.0, 0.0, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Type = raylib.CameraPerspective

	cubePosition := raylib.NewVector3(0.0, 0.0, 0.0)

	raylib.SetCameraMode(camera, raylib.CameraFree) // Set a free camera mode

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera) // Update camera

		if raylib.IsKeyDown(raylib.KeyZ) {
			camera.Target = raylib.NewVector3(0.0, 0.0, 0.0)
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode3D(camera)

		raylib.DrawCube(cubePosition, 2.0, 2.0, 2.0, raylib.Red)
		raylib.DrawCubeWires(cubePosition, 2.0, 2.0, 2.0, raylib.Maroon)

		raylib.DrawGrid(10, 1.0)

		raylib.EndMode3D()

		raylib.DrawRectangle(10, 10, 320, 133, raylib.Fade(raylib.SkyBlue, 0.5))
		raylib.DrawRectangleLines(10, 10, 320, 133, raylib.Blue)

		raylib.DrawText("Free camera default controls:", 20, 20, 10, raylib.Black)
		raylib.DrawText("- Mouse Wheel to Zoom in-out", 40, 40, 10, raylib.DarkGray)
		raylib.DrawText("- Mouse Wheel Pressed to Pan", 40, 60, 10, raylib.DarkGray)
		raylib.DrawText("- Alt + Mouse Wheel Pressed to Rotate", 40, 80, 10, raylib.DarkGray)
		raylib.DrawText("- Alt + Ctrl + Mouse Wheel Pressed for Smooth Zoom", 40, 100, 10, raylib.DarkGray)
		raylib.DrawText("- Z to zoom to (0, 0, 0)", 40, 120, 10, raylib.DarkGray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
