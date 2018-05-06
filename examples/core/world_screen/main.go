package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [core] example - 3d camera free")

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(10.0, 10.0, 10.0) // Camera position
	camera.Target = raylib.NewVector3(0.0, 0.0, 0.0)      // Camera looking at point
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)          // Camera up vector (rotation towards target)
	camera.Fovy = 45.0                                    // Camera field-of-view Y

	cubePosition := raylib.NewVector3(0.0, 0.0, 0.0)
	cubeScreenPosition := raylib.Vector2{}

	raylib.SetCameraMode(camera, raylib.CameraFree) // Set a free camera mode

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera) // Update camera

		// Calculate cube screen space position (with a little offset to be in top)
		cubeScreenPosition = raylib.GetWorldToScreen(raylib.NewVector3(cubePosition.X, cubePosition.Y+2.5, cubePosition.Z), camera)

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode3D(camera)

		raylib.DrawCube(cubePosition, 2.0, 2.0, 2.0, raylib.Red)
		raylib.DrawCubeWires(cubePosition, 2.0, 2.0, 2.0, raylib.Maroon)

		raylib.DrawGrid(10, 1.0)

		raylib.EndMode3D()

		raylib.DrawText("Enemy: 100 / 100", int32(cubeScreenPosition.X)-raylib.MeasureText("Enemy: 100 / 100", 20)/2, int32(cubeScreenPosition.Y), 20, raylib.Black)
		raylib.DrawText("Text is always on top of the cube", (screenWidth-raylib.MeasureText("Text is always on top of the cube", 20))/2, 25, 20, raylib.Gray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
