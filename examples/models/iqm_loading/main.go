package main

import (
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Helper function to extract animation name string from ModelAnimation.Name
func getAnimName(anim rl.ModelAnimation) string {
	var bytes []byte
	for _, b := range anim.Name {
		if b == 0 {
			break
		}
		bytes = append(bytes, byte(b))
	}
	return string(bytes)
}

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - loading iqm")
	defer rl.CloseWindow()

	// Define the camera to look into our 3D world
	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(10.0, 10.0, 10.0) // Camera position
	camera.Target = rl.NewVector3(0.0, 4.0, 0.0)      // Camera looking at point
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)          // Camera up vector
	camera.Fovy = 45.0                                // Camera field-of-view Y
	camera.Projection = rl.CameraPerspective          // Camera mode type

	// Load the animated model mesh and basic data directly from root folder
	model := rl.LoadModel("guy.iqm")

	// Load model texture and set material map directly from root folder
	texture := rl.LoadTexture("guytex.png")

	materials := unsafe.Slice(model.Materials, model.MaterialCount)
	rl.SetMaterialTexture(&materials[0], rl.MapDiffuse, texture)

	position := rl.NewVector3(0.0, 0.0, 0.0) // Set model position

	// Load animation data directly from root folder
	anims := rl.LoadModelAnimations("guyanim.iqm")

	// Animation playing variables
	animIndex := 0
	animCurrentFrame := float32(0.0)

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		rl.UpdateCamera(&camera, rl.CameraOrbital)

		if rl.IsKeyPressed(rl.KeySpace) {
			animIndex++
			if animIndex > len(anims)-1 {
				animIndex = 0
			}
		}

		// Play animation
		if len(anims) > 0 {
			animCurrentFrame += 1.0
			rl.UpdateModelAnimation(model, anims[animIndex], animCurrentFrame)

			if animCurrentFrame >= float32(anims[animIndex].KeyframeCount) {
				animCurrentFrame = 0.0
			}
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModelEx(model, position, rl.NewVector3(1.0, 0.0, 0.0), -90.0, rl.NewVector3(1.0, 1.0, 1.0), rl.White)

		rl.DrawGrid(10, 1.0)

		rl.EndMode3D()

		animName := ""
		if len(anims) > 0 {
			animName = getAnimName(anims[animIndex])
		}

		rl.DrawText(fmt.Sprintf("Current animation: %s", animName), 10, 10, 20, rl.Maroon)
		rl.DrawText("Press space to switch animations", 10, screenHeight-30, 20, rl.Maroon)
		rl.DrawText("(c) Guy IQM 3D model by @culacant", screenWidth-200, screenHeight-20, 10, rl.Gray)

		rl.EndDrawing()
	}

	rl.UnloadModel(model)
	rl.UnloadTexture(texture)
	rl.UnloadModelAnimations(anims)

	rl.CloseWindow()
}
