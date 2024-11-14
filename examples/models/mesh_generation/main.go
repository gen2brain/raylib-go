package main

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	numModels := 9

	rl.InitWindow(screenWidth, screenHeight, "raylib [models] example - mesh generation")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(10.0, 5.0, 10.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	checked := rl.GenImageChecked(2, 2, 1, 1, rl.Black, rl.Red)
	texture := rl.LoadTextureFromImage(checked)
	rl.UnloadImage(checked)

	models := make([]rl.Model, numModels)

	models[0] = rl.LoadModelFromMesh(rl.GenMeshPlane(2, 2, 4, 3))
	models[1] = rl.LoadModelFromMesh(rl.GenMeshCube(2, 1, 2))
	models[2] = rl.LoadModelFromMesh(rl.GenMeshSphere(2, 32, 32))
	models[3] = rl.LoadModelFromMesh(rl.GenMeshHemiSphere(2, 16, 16))
	models[4] = rl.LoadModelFromMesh(rl.GenMeshCylinder(1, 2, 16))
	models[5] = rl.LoadModelFromMesh(rl.GenMeshTorus(0.25, 4, 16, 32))
	models[6] = rl.LoadModelFromMesh(rl.GenMeshKnot(1, 2, 16, 128))
	models[7] = rl.LoadModelFromMesh(rl.GenMeshPoly(5, 2))
	models[8] = rl.LoadModelFromMesh(GenMeshCustom())

	for i := 0; i < numModels; i++ {
		rl.SetMaterialTexture(models[i].Materials, rl.MapDiffuse, texture)
	}

	position := rl.Vector3Zero()

	currentModel := 0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rl.UpdateCamera(&camera, rl.CameraOrbital)

		if rl.IsKeyPressed(rl.KeyUp) {
			currentModel = (currentModel + 1) % numModels // Cycle between the textures
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			// Adding numModels here is necessary to avoid a crash
			// where the golang % (modulus) operator, doesn't work as
			// one might expect for negative numbers.
			currentModel = (currentModel + numModels - 1) % numModels // Cycle between the textures
		}
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			currentModel = (currentModel + 1) % numModels // Cycle between the textures
		}
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(models[currentModel], position, 1, rl.White)
		rl.DrawGrid(10, 1)

		rl.EndMode3D()

		rl.DrawRectangle(10, 10, 310, 50, rl.Fade(rl.SkyBlue, 0.5))
		rl.DrawRectangleLines(10, 10, 310, 50, rl.Fade(rl.DarkBlue, 0.5))
		rl.DrawText("UP/DOWN ARROW KEY OR LEFT MOUSE", 20, 20, 10, rl.Blue)
		rl.DrawText("BUTTON TO CHANGE MODELS", 20, 40, 10, rl.Blue)

		txt := "PLANE"
		switch currentModel {
		case 1:
			txt = "CUBE"
		case 2:
			txt = "SPHERE"
		case 3:
			txt = "HEMISPHERE"
		case 4:
			txt = "CYLINDER"
		case 5:
			txt = "TORUS"
		case 6:
			txt = "KNOT"
		case 7:
			txt = "POLY"
		case 8:
			txt = "Custom (triangle)"
		}
		txtlen := rl.MeasureText(txt, 20)
		rl.DrawText(txt, screenWidth/2-txtlen/2, 10, 20, rl.DarkBlue)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)
	// Custom models meshes needs to be
	// cleared manually
	clearCustomMesh(models[8])
	for i := 0; i < numModels; i++ {
		rl.UnloadModel(models[i])
	}

	rl.CloseWindow()
}

func clearCustomMesh(model rl.Model) {
	// Vertices, Normals and Texcoords of your CUSTOM mesh are Go slices.
	// UnloadModel calls UnloadMesh for every mesh and UnloadMesh tries
	// to free your Go slices. This will panic because it cannot free
	// Go slices. Free() is a C function and it expects to free C memory
	// and not a Go slice. So clear the slices manually like this.
	model.Meshes.Vertices = nil
	model.Meshes.Normals = nil
	model.Meshes.Texcoords = nil
}

// GenMeshCustom generates a simple triangle mesh from code
func GenMeshCustom() rl.Mesh {
	mesh := rl.Mesh{
		TriangleCount: 1,
		VertexCount:   3,
	}

	var vertices, normals, texcoords []float32

	// 3 vertices
	vertices = addCoord(vertices, 0, 0, 0)
	vertices = addCoord(vertices, 1, 0, 2)
	vertices = addCoord(vertices, 2, 0, 0)
	mesh.Vertices = unsafe.SliceData(vertices)

	// 3 normals
	normals = addCoord(normals, 0, 1, 0)
	normals = addCoord(normals, 0, 1, 0)
	normals = addCoord(normals, 0, 1, 0)
	mesh.Normals = unsafe.SliceData(normals)

	// 3 texcoords
	texcoords = addCoord(texcoords, 0, 0)
	texcoords = addCoord(texcoords, 0.5, 1)
	texcoords = addCoord(texcoords, 1, 0)
	mesh.Texcoords = unsafe.SliceData(texcoords)

	// Upload mesh data from CPU (RAM) to GPU (VRAM) memory
	rl.UploadMesh(&mesh, false)

	return mesh
}

func addCoord(slice []float32, values ...float32) []float32 {
	for _, value := range values {
		slice = append(slice, value)
	}
	return slice
}
