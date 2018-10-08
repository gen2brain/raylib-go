package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - geometric shapes")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(0.0, 10.0, 10.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawCube(rl.NewVector3(-4.0, 0.0, 2.0), 2.0, 5.0, 2.0, rl.Red)
		rl.DrawCubeWires(rl.NewVector3(-4.0, 0.0, 2.0), 2.0, 5.0, 2.0, rl.Gold)
		rl.DrawCubeWires(rl.NewVector3(-4.0, 0.0, -2.0), 3.0, 6.0, 2.0, rl.Maroon)

		rl.DrawSphere(rl.NewVector3(-1.0, 0.0, -2.0), 1.0, rl.Green)
		rl.DrawSphereWires(rl.NewVector3(1.0, 0.0, 2.0), 2.0, 16, 16, rl.Lime)

		rl.DrawCylinder(rl.NewVector3(4.0, 0.0, -2.0), 1.0, 2.0, 3.0, 4, rl.SkyBlue)
		rl.DrawCylinderWires(rl.NewVector3(4.0, 0.0, -2.0), 1.0, 2.0, 3.0, 4, rl.DarkBlue)
		rl.DrawCylinderWires(rl.NewVector3(4.5, -1.0, 2.0), 1.0, 1.0, 2.0, 6, rl.Brown)

		rl.DrawCylinder(rl.NewVector3(1.0, 0.0, -4.0), 0.0, 1.5, 3.0, 8, rl.Gold)
		rl.DrawCylinderWires(rl.NewVector3(1.0, 0.0, -4.0), 0.0, 1.5, 3.0, 8, rl.Pink)

		rl.DrawGrid(10, 1.0) // Draw a grid

		rl.EndMode3D()

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
