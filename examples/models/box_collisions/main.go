package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [models] example - box collisions")

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(0.0, 10.0, 10.0)
	camera.Target = raylib.NewVector3(0.0, 0.0, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Type = raylib.CameraPerspective

	playerPosition := raylib.NewVector3(0.0, 1.0, 2.0)
	playerSize := raylib.NewVector3(1.0, 2.0, 1.0)
	playerColor := raylib.Green

	enemyBoxPos := raylib.NewVector3(-4.0, 1.0, 0.0)
	enemyBoxSize := raylib.NewVector3(2.0, 2.0, 2.0)

	enemySpherePos := raylib.NewVector3(4.0, 0.0, 0.0)
	enemySphereSize := float32(1.5)

	collision := false

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update

		// Move player
		if raylib.IsKeyDown(raylib.KeyRight) {
			playerPosition.X += 0.2
		} else if raylib.IsKeyDown(raylib.KeyLeft) {
			playerPosition.X -= 0.2
		} else if raylib.IsKeyDown(raylib.KeyDown) {
			playerPosition.Z += 0.2
		} else if raylib.IsKeyDown(raylib.KeyUp) {
			playerPosition.Z -= 0.2
		}

		collision = false

		// Check collisions player vs enemy-box
		if raylib.CheckCollisionBoxes(
			raylib.NewBoundingBox(
				raylib.NewVector3(playerPosition.X-playerSize.X/2, playerPosition.Y-playerSize.Y/2, playerPosition.Z-playerSize.Z/2),
				raylib.NewVector3(playerPosition.X+playerSize.X/2, playerPosition.Y+playerSize.Y/2, playerPosition.Z+playerSize.Z/2)),
			raylib.NewBoundingBox(
				raylib.NewVector3(enemyBoxPos.X-enemyBoxSize.X/2, enemyBoxPos.Y-enemyBoxSize.Y/2, enemyBoxPos.Z-enemyBoxSize.Z/2),
				raylib.NewVector3(enemyBoxPos.X+enemyBoxSize.X/2, enemyBoxPos.Y+enemyBoxSize.Y/2, enemyBoxPos.Z+enemyBoxSize.Z/2)),
		) {
			collision = true
		}

		// Check collisions player vs enemy-sphere
		if raylib.CheckCollisionBoxSphere(
			raylib.NewBoundingBox(
				raylib.NewVector3(playerPosition.X-playerSize.X/2, playerPosition.Y-playerSize.Y/2, playerPosition.Z-playerSize.Z/2),
				raylib.NewVector3(playerPosition.X+playerSize.X/2, playerPosition.Y+playerSize.Y/2, playerPosition.Z+playerSize.Z/2)),
			enemySpherePos,
			enemySphereSize,
		) {
			collision = true
		}

		if collision {
			playerColor = raylib.Red
		} else {
			playerColor = raylib.Green
		}

		// Draw

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode3D(camera)

		// Draw enemy-box
		raylib.DrawCube(enemyBoxPos, enemyBoxSize.X, enemyBoxSize.Y, enemyBoxSize.Z, raylib.Gray)
		raylib.DrawCubeWires(enemyBoxPos, enemyBoxSize.X, enemyBoxSize.Y, enemyBoxSize.Z, raylib.DarkGray)

		// Draw enemy-sphere
		raylib.DrawSphere(enemySpherePos, enemySphereSize, raylib.Gray)
		raylib.DrawSphereWires(enemySpherePos, enemySphereSize, 16, 16, raylib.DarkGray)

		// Draw player
		raylib.DrawCubeV(playerPosition, playerSize, playerColor)

		raylib.DrawGrid(10, 1.0) // Draw a grid

		raylib.EndMode3D()

		raylib.DrawText("Move player with cursors to collide", 220, 40, 20, raylib.Gray)

		raylib.DrawFPS(10, 10)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
