package main

import (
	"fmt"
	"path/filepath"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	skyboxFilename := "skybox.png"

	rl.InitWindow(800, 450, "raylib [models] example - skybox loading and drawing")

	camera := rl.NewCamera3D(
		rl.NewVector3(1.0, 1.0, 1.0),
		rl.NewVector3(4.0, 1.0, 4.0),
		rl.NewVector3(0.0, 1.0, 0.0),
		45.0,
		rl.CameraPerspective,
	)

	// load skybox shader and set required locations
	skyboxShader := rl.LoadShader("skybox.vs", "skybox.fs")

	setShaderIntValue(skyboxShader, "environmentMap", rl.MapCubemap)

	// load skybox model
	cube := rl.GenMeshCube(1.0, 1.0, 1.0)
	skybox := rl.LoadModelFromMesh(cube)

	skybox.Materials.Shader = skyboxShader

	// load cubemap texture
	skyboxImg := rl.LoadImage(skyboxFilename)

	skyboxTexture := rl.LoadTextureCubemap(skyboxImg, rl.CubemapLayoutAutoDetect)

	rl.UnloadImage(skyboxImg)

	rl.SetMaterialTexture(skybox.Materials, rl.MapCubemap, skyboxTexture)

	// limit cursor to relative movement inside the window
	rl.DisableCursor()

	// set our game to run at 60 frames-per-second
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)

		// load new cubemap texture on drag&drop
		if rl.IsFileDropped() {
			droppedFiles := rl.LoadDroppedFiles()

			// only support one file dropped
			if len(droppedFiles) == 1 {
				switch filepath.Ext(droppedFiles[0]) {
				case ".png", ".jpg", ".bmp", ".tga":
					skyboxFilename = droppedFiles[0]

					rl.UnloadTexture(skyboxTexture)

					img := rl.LoadImage(skyboxFilename)

					skyboxTexture = rl.LoadTextureCubemap(img, rl.CubemapLayoutAutoDetect)

					rl.UnloadImage(img)

					rl.SetMaterialTexture(skybox.Materials, rl.MapCubemap, skyboxTexture)
				}
			}

			rl.UnloadDroppedFiles()
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.White)

		rl.BeginMode3D(camera)

		// we are inside the cube, we need to disable backface culling
		rl.DisableBackfaceCulling()
		rl.DisableDepthMask()

		rl.DrawModel(skybox, rl.NewVector3(0, 0, 0), 1.0, rl.White)

		// restore depth and backface culling
		rl.EnableBackfaceCulling()
		rl.EnableDepthMask()

		rl.DrawGrid(10, 1.0)

		rl.EndMode3D()

		rl.DrawText(
			fmt.Sprintf("File: %s", skyboxFilename),
			10,
			int32(rl.GetScreenHeight()-20),
			10,
			rl.Black,
		)

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	rl.UnloadModel(skybox)
	rl.UnloadTexture(skyboxTexture)
	rl.UnloadShader(skyboxShader)

	rl.CloseWindow()
}

func setShaderIntValue(shader rl.Shader, name string, value int32) {
	rl.SetShaderValue(
		shader,
		rl.GetShaderLocation(shader, name),
		unsafe.Slice((*float32)(unsafe.Pointer(&value)), 4),
		rl.ShaderUniformInt,
	)
}
