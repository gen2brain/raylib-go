/*******************************************************************************************
*
*   raylib [models] example - Load models M3D
*
*   Example originally created with raylib 4.5, last time updated with raylib 4.5
*
*   Example contributed by bzt (@bztsrc) and reviewed by Ramon Santamaria (@raysan5)
*
*   NOTES:
*     - Model3D (M3D) fileformat specs: https://gitlab.com/bztsrc/model3d
*     - Bender M3D exported: https://gitlab.com/bztsrc/model3d/-/tree/master/blender
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2022-2024 bzt (@bztsrc)
*
********************************************************************************************/
package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

func main() {
	// Initialization
	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - M3D model loading")

	// Define the camera to look into our 3d world
	camera := rl.Camera{
		Position:   rl.NewVector3(1.5, 1.5, 1.5),
		Target:     rl.NewVector3(0.0, 0.4, 0.0),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}

	position := rl.NewVector3(0.0, 0.0, 0.0)

	modelFileName := "cesium_man.m3d"
	drawMesh := true
	drawSkeleton := true
	animPlaying := false // Store anim state, what to draw

	// Load model
	model := rl.LoadModel(modelFileName)

	// Load animations

	animFrameCounter := 0
	animID := 0
	anims := rl.LoadModelAnimations(modelFileName)
	animsCount := int32(len(anims))

	rl.DisableCursor()
	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)

		if animsCount > 0 {
			// Play animation when space bar is held down (or step one frame with N)
			if rl.IsKeyDown(rl.KeySpace) || rl.IsKeyPressed(rl.KeyN) {
				animFrameCounter++
				if animFrameCounter >= int(anims[animID].FrameCount) {
					animFrameCounter = 0
				}
				rl.UpdateModelAnimation(model, anims[animID], int32(animFrameCounter))
				animPlaying = true
			}

			// Select animation by pressing C
			if rl.IsKeyPressed(rl.KeyC) {
				animFrameCounter = 0
				animID++
				if animID >= int(animsCount) {
					animID = 0
				}
				rl.UpdateModelAnimation(model, anims[animID], 0)
				animPlaying = true
			}
		}

		// Toggle skeleton drawing
		if rl.IsKeyPressed(rl.KeyB) {
			drawSkeleton = !drawSkeleton
		}

		// Toggle mesh drawing
		if rl.IsKeyPressed(rl.KeyM) {
			drawMesh = !drawMesh
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		// Draw 3d model with texture
		if drawMesh {
			rl.DrawModel(model, position, 1.0, rl.White)
		}

		// Draw the animated skeleton
		if drawSkeleton {
			modelBones := model.GetBones()
			modelPoses := model.GetBindPose()
			anim := anims[animID]
			animBones := anim.GetBones()
			for bone := 0; bone < int(model.BoneCount)-1; bone++ {
				if !animPlaying || animsCount == 0 {
					// Display the bind-pose skeleton
					rl.DrawCube(modelPoses[bone].Translation, 0.04, 0.04, 0.04, rl.Red)
					if modelBones[bone].Parent >= 0 {
						rl.DrawLine3D(modelPoses[bone].Translation, modelPoses[modelBones[bone].Parent].Translation, rl.Red)
					}
				} else {
					// // Display the frame-pose skeleton
					pos := anim.GetFramePose(animFrameCounter, bone).Translation
					rl.DrawCube(pos, 0.05, 0.05, 0.05, rl.Red)
					if animBones[bone].Parent >= 0 {
						endPos := anim.GetFramePose(animFrameCounter, int(animBones[bone].Parent)).Translation
						rl.DrawLine3D(pos, endPos, rl.Red)
					}
				}
			}
		}

		rl.DrawGrid(10, 1.0)

		rl.EndMode3D()

		rl.DrawText("PRESS SPACE to PLAY MODEL ANIMATION", 10, screenHeight-80, 10, rl.Maroon)
		rl.DrawText("PRESS N to STEP ONE ANIMATION FRAME", 10, screenHeight-60, 10, rl.DarkGray)
		rl.DrawText("PRESS C to CYCLE THROUGH ANIMATIONS", 10, screenHeight-40, 10, rl.DarkGray)
		rl.DrawText("PRESS M to toggle MESH, B to toggle SKELETON DRAWING", 10, screenHeight-20, 10, rl.DarkGray)
		rl.DrawText("(c) CesiumMan model by KhronosGroup", screenWidth-210, screenHeight-20, 10, rl.Gray)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadModelAnimations(anims)
	rl.UnloadModel(model)
	rl.CloseWindow()
}
