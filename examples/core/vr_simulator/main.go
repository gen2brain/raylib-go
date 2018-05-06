package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1080)
	screenHeight := int32(600)

	// NOTE: screenWidth/screenHeight should match VR device aspect ratio

	raylib.InitWindow(screenWidth, screenHeight, "raylib [core] example - vr simulator")

	// NOTE: default device (simulator)
	raylib.InitVrSimulator(raylib.GetVrDeviceInfo(raylib.HmdOculusRiftCv1)) // Init VR device (Oculus Rift CV1)

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(5.0, 2.0, 5.0) // Camera position
	camera.Target = raylib.NewVector3(0.0, 2.0, 0.0)   // Camera looking at point
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)       // Camera up vector (rotation towards target)
	camera.Fovy = 60.0                                 // Camera field-of-view Y

	cubePosition := raylib.NewVector3(0.0, 0.0, 0.0)

	raylib.SetCameraMode(camera, raylib.CameraFirstPerson) // Set first person camera mode

	raylib.SetTargetFPS(90)

	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera) // Update camera (simulator mode)

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginVrDrawing()

		raylib.BeginMode3D(camera)

		raylib.DrawCube(cubePosition, 2.0, 2.0, 2.0, raylib.Red)
		raylib.DrawCubeWires(cubePosition, 2.0, 2.0, 2.0, raylib.Maroon)

		raylib.DrawGrid(40, 1.0)

		raylib.EndMode3D()

		raylib.EndVrDrawing()

		raylib.DrawFPS(10, 10)

		raylib.EndDrawing()
	}

	raylib.CloseVrSimulator() // Close VR simulator

	raylib.CloseWindow()
}
