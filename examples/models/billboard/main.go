package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [models] example - drawing billboards")

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(5.0, 4.0, 5.0)
	camera.Target = raylib.NewVector3(0.0, 2.0, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Type = raylib.CameraPerspective

	bill := raylib.LoadTexture("billboard.png")      // Our texture billboard
	billPosition := raylib.NewVector3(0.0, 2.0, 0.0) // Position where draw billboard

	raylib.SetCameraMode(camera, raylib.CameraOrbital) // Set an orbital camera mode

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera) // Update camera

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode3D(camera)

		raylib.DrawBillboard(camera, bill, billPosition, 2.0, raylib.White)

		raylib.DrawGrid(10, 1.0) // Draw a grid

		raylib.EndMode3D()

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(bill) // Unload texture

	raylib.CloseWindow()
}
