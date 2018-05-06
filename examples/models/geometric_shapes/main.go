package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [models] example - geometric shapes")

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(0.0, 10.0, 10.0)
	camera.Target = raylib.NewVector3(0.0, 0.0, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode3D(camera)

		raylib.DrawCube(raylib.NewVector3(-4.0, 0.0, 2.0), 2.0, 5.0, 2.0, raylib.Red)
		raylib.DrawCubeWires(raylib.NewVector3(-4.0, 0.0, 2.0), 2.0, 5.0, 2.0, raylib.Gold)
		raylib.DrawCubeWires(raylib.NewVector3(-4.0, 0.0, -2.0), 3.0, 6.0, 2.0, raylib.Maroon)

		raylib.DrawSphere(raylib.NewVector3(-1.0, 0.0, -2.0), 1.0, raylib.Green)
		raylib.DrawSphereWires(raylib.NewVector3(1.0, 0.0, 2.0), 2.0, 16, 16, raylib.Lime)

		raylib.DrawCylinder(raylib.NewVector3(4.0, 0.0, -2.0), 1.0, 2.0, 3.0, 4, raylib.SkyBlue)
		raylib.DrawCylinderWires(raylib.NewVector3(4.0, 0.0, -2.0), 1.0, 2.0, 3.0, 4, raylib.DarkBlue)
		raylib.DrawCylinderWires(raylib.NewVector3(4.5, -1.0, 2.0), 1.0, 1.0, 2.0, 6, raylib.Brown)

		raylib.DrawCylinder(raylib.NewVector3(1.0, 0.0, -4.0), 0.0, 1.5, 3.0, 8, raylib.Gold)
		raylib.DrawCylinderWires(raylib.NewVector3(1.0, 0.0, -4.0), 0.0, 1.5, 3.0, 8, raylib.Pink)

		raylib.DrawGrid(10, 1.0) // Draw a grid

		raylib.EndMode3D()

		raylib.DrawFPS(10, 10)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
