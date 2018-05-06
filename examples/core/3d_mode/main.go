package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	raylib.InitWindow(800, 450, "raylib [core] example - 3d mode")

	camera := raylib.Camera3D{}
	camera.Position = raylib.NewVector3(0.0, 10.0, 10.0)
	camera.Target = raylib.NewVector3(0.0, 0.0, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Type = raylib.CameraPerspective

	cubePosition := raylib.NewVector3(0.0, 0.0, 0.0)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode3D(camera)

		raylib.DrawCube(cubePosition, 2.0, 2.0, 2.0, raylib.Red)
		raylib.DrawCubeWires(cubePosition, 2.0, 2.0, 2.0, raylib.Maroon)

		raylib.DrawGrid(10, 1.0)

		raylib.EndMode3D()

		raylib.DrawText("Welcome to the third dimension!", 10, 40, 20, raylib.DarkGray)

		raylib.DrawFPS(10, 10)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
