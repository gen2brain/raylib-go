package main

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	boneSockets     = 3
	boneSocketHat   = 0
	boneSocketHandR = 1
	boneSocketHandL = 2
)

// Helper function to extract bone name string from BoneInfo
func getBoneName(bone rl.BoneInfo) string {
	var bytes []byte
	for _, b := range bone.Name {
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

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - bone socket")
	defer rl.CloseWindow()

	// Define the camera to look into our 3D world
	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(5.0, 5.0, 5.0) // Camera position
	camera.Target = rl.NewVector3(0.0, 2.0, 0.0)   // Camera looking at point
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)       // Camera up vector
	camera.Fovy = 45.0                             // Camera field-of-view Y
	camera.Projection = rl.CameraPerspective       // Camera projection type

	// Load gltf character model directly from root folder
	characterModel := rl.LoadModel("greenman.glb")

	// Load equipment models directly from root folder
	equipModel := [boneSockets]rl.Model{
		rl.LoadModel("greenman_hat.glb"),
		rl.LoadModel("greenman_sword.glb"),
		rl.LoadModel("greenman_shield.glb"),
	}

	showEquip := [3]bool{true, true, true} // Toggle on/off equip

	// Load gltf model animations
	modelAnimations := rl.LoadModelAnimations("greenman.glb")
	animsCount := len(modelAnimations)

	animIndex := 0
	animCurrentFrame := int32(0)

	// Indices of bones for sockets
	boneSocketIndex := [boneSockets]int32{-1, -1, -1}

	// Search bones for sockets via characterModel.Skeleton
	bones := unsafe.Slice(characterModel.Skeleton.Bones, characterModel.Skeleton.BoneCount)
	for i := int32(0); i < characterModel.Skeleton.BoneCount; i++ {
		name := getBoneName(bones[i])
		if name == "socket_hat" {
			boneSocketIndex[boneSocketHat] = i
			continue
		}
		if name == "socket_hand_R" {
			boneSocketIndex[boneSocketHandR] = i
			continue
		}
		if name == "socket_hand_L" {
			boneSocketIndex[boneSocketHandL] = i
			continue
		}
	}

	position := rl.NewVector3(0.0, 0.0, 0.0) // Set model position
	angle := 0                               // Set angle for character rotation

	rl.DisableCursor() // Limit cursor to relative movement inside the window
	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		rl.UpdateCamera(&camera, rl.CameraThirdPerson)

		// Rotate character
		if rl.IsKeyDown(rl.KeyF) {
			angle = (angle + 1) % 360
		} else if rl.IsKeyDown(rl.KeyH) {
			angle = (360 + angle - 1) % 360
		}

		// Select current animation
		if animsCount > 0 {
			if rl.IsKeyPressed(rl.KeyT) {
				animIndex = (animIndex + 1) % animsCount
			} else if rl.IsKeyPressed(rl.KeyG) {
				animIndex = (animIndex + animsCount - 1) % animsCount
			}
		}

		// Toggle shown equipment
		if rl.IsKeyPressed(rl.KeyOne) {
			showEquip[boneSocketHat] = !showEquip[boneSocketHat]
		}
		if rl.IsKeyPressed(rl.KeyTwo) {
			showEquip[boneSocketHandR] = !showEquip[boneSocketHandR]
		}
		if rl.IsKeyPressed(rl.KeyThree) {
			showEquip[boneSocketHandL] = !showEquip[boneSocketHandL]
		}

		if animsCount > 0 {
			anim := modelAnimations[animIndex]
			animCurrentFrame = (animCurrentFrame + 1) % anim.KeyframeCount
			rl.UpdateModelAnimation(characterModel, anim, float32(animCurrentFrame))
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		charMeshes := unsafe.Slice(characterModel.Meshes, characterModel.MeshCount)
		charMaterials := unsafe.Slice(characterModel.Materials, characterModel.MaterialCount)

		// Draw character
		characterRotate := rl.QuaternionFromAxisAngle(rl.NewVector3(0.0, 1.0, 0.0), float32(angle)*rl.Deg2rad)
		characterModel.Transform = rl.MatrixMultiply(rl.QuaternionToMatrix(characterRotate), rl.MatrixTranslate(position.X, position.Y, position.Z))

		if animsCount > 0 {
			anim := modelAnimations[animIndex]
			rl.UpdateModelAnimation(characterModel, anim, float32(animCurrentFrame))
			rl.DrawMesh(charMeshes[0], charMaterials[1], characterModel.Transform)

			// Draw equipment (hat, sword, shield)
			keyframePosesList := unsafe.Slice(anim.KeyframePoses, anim.KeyframeCount)
			keyframePoses := unsafe.Slice(keyframePosesList[animCurrentFrame], characterModel.Skeleton.BoneCount)
			bindPoses := unsafe.Slice(characterModel.Skeleton.BindPose, characterModel.Skeleton.BoneCount)

			for i := 0; i < boneSockets; i++ {
				if !showEquip[i] || boneSocketIndex[i] == -1 {
					continue
				}

				socketIdx := boneSocketIndex[i]
				transform := keyframePoses[socketIdx]
				inRotation := bindPoses[socketIdx].Rotation
				outRotation := transform.Rotation

				// Calculate socket rotation
				rotate := rl.QuaternionMultiply(outRotation, rl.QuaternionInvert(inRotation))
				matrixTransform := rl.QuaternionToMatrix(rotate)
				matrixTransform = rl.MatrixMultiply(matrixTransform, rl.MatrixTranslate(transform.Translation.X, transform.Translation.Y, transform.Translation.Z))
				matrixTransform = rl.MatrixMultiply(matrixTransform, characterModel.Transform)

				equipMeshes := unsafe.Slice(equipModel[i].Meshes, equipModel[i].MeshCount)
				equipMaterials := unsafe.Slice(equipModel[i].Materials, equipModel[i].MaterialCount)

				// Draw mesh at socket position with socket angle rotation
				rl.DrawMesh(equipMeshes[0], equipMaterials[1], matrixTransform)
			}
		} else {
			rl.DrawMesh(charMeshes[0], charMaterials[1], characterModel.Transform)
		}

		rl.DrawGrid(10, 1.0)

		rl.EndMode3D()

		rl.DrawText("Use the T/G to switch animation", 10, 10, 20, rl.Gray)
		rl.DrawText("Use the F/H to rotate character left/right", 10, 35, 20, rl.Gray)
		rl.DrawText("Use the 1,2,3 to toggle shown of hat, sword and shield", 10, 60, 20, rl.Gray)

		rl.EndDrawing()
	}
	rl.UnloadModel(characterModel)
	rl.UnloadModelAnimations(modelAnimations)
	for i := 0; i < boneSockets; i++ {
		rl.UnloadModel(equipModel[i])
	}
	rl.CloseWindow()
}
