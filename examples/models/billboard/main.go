package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - drawing billboards")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(5.0, 4.0, 5.0)
	camera.Target = rl.NewVector3(0.0, 2.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	bill := rl.LoadTexture("billboard.png")      // Our texture billboard
	billPosition := rl.NewVector3(0.0, 2.0, 0.0) // Position where draw billboard

	rl.SetCameraMode(camera, rl.CameraOrbital) // Set an orbital camera mode

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera) // Update camera

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawGrid(10, 1.0) // Draw a grid

		rl.DrawBillboard(camera, bill, billPosition, 2.0, rl.White)

		rl.EndMode3D()

		rl.EndDrawing()
	}

	rl.UnloadTexture(bill) // Unload texture

	rl.CloseWindow()
}
