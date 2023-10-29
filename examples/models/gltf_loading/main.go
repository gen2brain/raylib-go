package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - gltf loading")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(5.0, 5.0, 5.0)
	camera.Target = rl.NewVector3(0.0, 2.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	model := rl.LoadModel("robot.glb")

	animIndex := 0
	animCurrentFrame := 0

	modelAnims := rl.LoadModelAnimations("robot.glb")

	position := rl.NewVector3(0, 0, 0)
	rl.DisableCursor()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rl.UpdateCamera(&camera, rl.CameraOrbital)

		if rl.IsKeyPressed(rl.KeyUp) {
			animIndex++
			if animIndex >= len(modelAnims) {
				animIndex = 0
			}
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			animIndex--
			if animIndex < 0 {
				animIndex = len(modelAnims) - 1
			}
		}

		animPlaying := modelAnims[animIndex]
		animCurrentFrame = (animCurrentFrame + 1) % int(animPlaying.FrameCount)
		rl.UpdateModelAnimation(model, animPlaying, int32(animCurrentFrame))

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(model, position, 0.8, rl.White)

		rl.EndMode3D()

		rl.DrawText("current animation number: "+fmt.Sprint(animIndex), 10, 10, 10, rl.Black)
		rl.DrawText("UP/DOWN ARROW KEYS CHANGE ANIMATION", 10, 30, 10, rl.Black)

		rl.EndDrawing()
	}

	rl.UnloadModel(model)

	rl.CloseWindow()
}
