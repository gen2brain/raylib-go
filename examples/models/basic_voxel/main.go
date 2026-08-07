package main

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const worldSize = 8 // Size of our voxel world (8x8x8 cubes)

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - basic voxel")
	defer rl.CloseWindow()

	rl.DisableCursor() // Lock mouse to window center

	// Define the camera to look into our 3D world (first person)
	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(-2.0, 0.0, -2.0) // Camera position at ground level
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)     // Camera looking at point
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)         // Camera up vector
	camera.Fovy = 45.0                               // Camera field-of-view Y
	camera.Projection = rl.CameraPerspective         // Camera projection type

	// Create a cube model
	cubeMesh := rl.GenMeshCube(1.0, 1.0, 1.0)
	cubeModel := rl.LoadModelFromMesh(cubeMesh)
	defer rl.UnloadModel(cubeModel)

	materials := unsafe.Slice(cubeModel.Materials, cubeModel.MaterialCount)
	maps := unsafe.Slice(materials[0].Maps, 12)
	maps[rl.MapDiffuse].Color = rl.Beige

	// Initialize voxel world - fill with voxels
	var voxels [worldSize][worldSize][worldSize]bool
	for x := 0; x < worldSize; x++ {
		for y := 0; y < worldSize; y++ {
			for z := 0; z < worldSize; z++ {
				voxels[x][y][z] = true
			}
		}
	}

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)

		// Handle voxel removal with mouse click
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			// Cast a ray from the screen center (where crosshair would be)
			screenCenter := rl.NewVector2(float32(rl.GetScreenWidth())/2.0, float32(rl.GetScreenHeight())/2.0)
			ray := rl.GetScreenToWorldRay(screenCenter, camera)

			// Check ray collision with all voxels
			closestDistance := float32(99999.0)
			closestVoxelPosition := rl.NewVector3(-1, -1, -1)
			voxelFound := false

			for x := 0; x < worldSize; x++ {
				for y := 0; y < worldSize; y++ {
					for z := 0; z < worldSize; z++ {
						if !voxels[x][y][z] {
							continue // Skip empty voxels
						}

						// Build a bounding box for this voxel
						position := rl.NewVector3(float32(x), float32(y), float32(z))
						box := rl.BoundingBox{
							Min: rl.NewVector3(position.X-0.5, position.Y-0.5, position.Z-0.5),
							Max: rl.NewVector3(position.X+0.5, position.Y+0.5, position.Z+0.5),
						}

						// Check ray-box collision
						collision := rl.GetRayCollisionBox(ray, box)
						if collision.Hit && collision.Distance < closestDistance {
							closestDistance = collision.Distance
							closestVoxelPosition = position
							voxelFound = true
						}
					}
				}
			}

			// Remove the closest voxel if one was hit
			if voxelFound {
				voxels[int(closestVoxelPosition.X)][int(closestVoxelPosition.Y)][int(closestVoxelPosition.Z)] = false
			}
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawGrid(10, 1.0)

		// Draw all voxels
		for x := 0; x < worldSize; x++ {
			for y := 0; y < worldSize; y++ {
				for z := 0; z < worldSize; z++ {
					if !voxels[x][y][z] {
						continue
					}

					position := rl.NewVector3(float32(x), float32(y), float32(z))
					rl.DrawModel(cubeModel, position, 1.0, rl.Beige)
					rl.DrawCubeWires(position, 1.0, 1.0, 1.0, rl.Black)
				}
			}
		}

		rl.EndMode3D()

		// Draw reference point for raycasting to delete blocks
		rl.DrawCircle(int32(rl.GetScreenWidth()/2), int32(rl.GetScreenHeight()/2), 4.0, rl.Red)

		rl.DrawText("Left-click a voxel to remove it!", 10, 10, 20, rl.DarkGray)
		rl.DrawText("WASD to move, mouse to look around", 10, 35, 10, rl.Gray)

		rl.EndDrawing()
	}
}
