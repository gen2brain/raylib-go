package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - 3d camera free")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(10.0, 10.0, 10.0) // Camera position
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)      // Camera looking at point
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)          // Camera up vector (rotation towards target)
	camera.Fovy = 45.0                                    // Camera field-of-view Y

	cubePosition := rl.NewVector3(0.0, 0.0, 0.0)
	cubeScreenPosition := rl.Vector2{}

	rl.SetCameraMode(camera, rl.CameraFree) // Set a free camera mode

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera) // Update camera

		// Calculate cube screen space position (with a little offset to be in top)
		cubeScreenPosition = rl.GetWorldToScreen(rl.NewVector3(cubePosition.X, cubePosition.Y+2.5, cubePosition.Z), camera)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawCube(cubePosition, 2.0, 2.0, 2.0, rl.Red)
		rl.DrawCubeWires(cubePosition, 2.0, 2.0, 2.0, rl.Maroon)

		rl.DrawGrid(10, 1.0)

		rl.EndMode3D()

		rl.DrawText("Enemy: 100 / 100", int32(cubeScreenPosition.X)-rl.MeasureText("Enemy: 100 / 100", 20)/2, int32(cubeScreenPosition.Y), 20, rl.Black)
		rl.DrawText("Text is always on top of the cube", (screenWidth-rl.MeasureText("Text is always on top of the cube", 20))/2, 25, 20, rl.Gray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
