package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	hmd := rl.GetVrDeviceInfo(rl.HmdOculusRiftCv1) // Oculus Rift CV1
	rl.InitWindow(int32(hmd.HScreenSize), int32(hmd.VScreenSize), "raylib [core] example - vr simulator")

	// NOTE: default device (simulator)
	rl.InitVrSimulator(hmd) // Init VR device

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(5.0, 2.0, 5.0) // Camera position
	camera.Target = rl.NewVector3(0.0, 2.0, 0.0)   // Camera looking at point
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)       // Camera up vector (rotation towards target)
	camera.Fovy = 60.0                                 // Camera field-of-view Y

	cubePosition := rl.NewVector3(0.0, 0.0, 0.0)

	rl.SetCameraMode(camera, rl.CameraFirstPerson) // Set first person camera mode

	rl.SetTargetFPS(90)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera) // Update camera (simulator mode)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginVrDrawing()

		rl.BeginMode3D(camera)

		rl.DrawCube(cubePosition, 2.0, 2.0, 2.0, rl.Red)
		rl.DrawCubeWires(cubePosition, 2.0, 2.0, 2.0, rl.Maroon)

		rl.DrawGrid(40, 1.0)

		rl.EndMode3D()

		rl.EndVrDrawing()

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	rl.CloseVrSimulator() // Close VR simulator

	rl.CloseWindow()
}
