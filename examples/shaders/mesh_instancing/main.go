package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_INSTANCES = 100000

func main() {
	var (
		screenWidth         = int32(800)   // Framebuffer width
		screenHeight        = int32(450)   // Framebuffer height
		fps                 = 60           // Frames per second
		speed               = 30           // Speed of jump animation
		groups              = 2            // Count of separate groups jumping around
		amp                 = float32(10)  // Maximum amplitude of jump
		variance            = float32(0.8) // Global variance in jump height
		loop                = float32(0.0) // Individual cube's computed loop timer
		textPositionY int32 = 300
		framesCounter       = 0
	)

	rl.SetConfigFlags(rl.FlagMsaa4xHint) // Enable Multi Sampling Anti Aliasing 4x (if available)
	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - mesh instancing")

	// Define the camera to look into our 3d world
	camera := rl.Camera{
		Position:   rl.NewVector3(-125.0, 125.0, -125.0),
		Target:     rl.NewVector3(0.0, 0.0, 0.0),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}

	cube := rl.GenMeshCube(1.0, 1.0, 1.0)

	rotations := make([]rl.Matrix, MAX_INSTANCES)    // Rotation state of instances
	rotationsInc := make([]rl.Matrix, MAX_INSTANCES) // Per-frame rotation animation of instances
	translations := make([]rl.Matrix, MAX_INSTANCES) // Locations of instances

	// Scatter random cubes around
	for i := 0; i < MAX_INSTANCES; i++ {
		x := float32(rl.GetRandomValue(-50, 50))
		y := float32(rl.GetRandomValue(-50, 50))
		z := float32(rl.GetRandomValue(-50, 50))
		translations[i] = rl.MatrixTranslate(x, y, z)

		x = float32(rl.GetRandomValue(0, 360))
		y = float32(rl.GetRandomValue(0, 360))
		z = float32(rl.GetRandomValue(0, 360))
		axis := rl.Vector3Normalize(rl.NewVector3(x, y, z))
		angle := float32(rl.GetRandomValue(0, 10)) * rl.Deg2rad

		rotationsInc[i] = rl.MatrixRotate(axis, angle)
		rotations[i] = rl.MatrixIdentity()
	}

	transforms := make([]rl.Matrix, MAX_INSTANCES)

	shader := rl.LoadShader("glsl330/base_lighting_instanced.vs", "glsl330/lighting.fs")
	shader.UpdateLocation(rl.LocMatrixMvp, rl.GetShaderLocation(shader, "mvp"))
	shader.UpdateLocation(rl.LocVectorView, rl.GetShaderLocation(shader, "viewPos"))
	shader.UpdateLocation(rl.LocMatrixModel, rl.GetShaderLocationAttrib(shader, "instanceTransform"))

	// ambient light level
	ambientLoc := rl.GetShaderLocation(shader, "ambient")
	rl.SetShaderValue(shader, ambientLoc, []float32{0.2, 0.2, 0.2, 1.0}, rl.ShaderUniformVec4)
	NewLight(LightTypeDirectional, rl.NewVector3(50.0, 50.0, 0.0), rl.Vector3Zero(), rl.White, shader)

	material := rl.LoadMaterialDefault()
	material.Shader = shader
	mmap := material.GetMap(rl.MapDiffuse)
	mmap.Color = rl.Red

	rl.SetCameraMode(camera, rl.CameraOrbital)

	rl.SetTargetFPS(int32(fps))
	for !rl.WindowShouldClose() {
		// Update
		//----------------------------------------------------------------------------------

		textPositionY = 300
		framesCounter++

		if rl.IsKeyDown(rl.KeyUp) {
			amp += 0.5
		}
		if rl.IsKeyDown(rl.KeyDown) {
			if amp <= 1 {
				amp = 1
			} else {
				amp -= 1
			}
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			if variance < 0 {
				variance = 0
			} else {
				variance -= 0.01
			}
		}
		if rl.IsKeyDown(rl.KeyRight) {
			if variance > 1 {
				variance = 1
			} else {
				variance += 0.01
			}
		}
		if rl.IsKeyDown(rl.KeyOne) {
			groups = 1
		}
		if rl.IsKeyDown(rl.KeyTwo) {
			groups = 2
		}
		if rl.IsKeyDown(rl.KeyThree) {
			groups = 3
		}
		if rl.IsKeyDown(rl.KeyFour) {
			groups = 4
		}
		if rl.IsKeyDown(rl.KeyFive) {
			groups = 5
		}
		if rl.IsKeyDown(rl.KeySix) {
			groups = 6
		}
		if rl.IsKeyDown(rl.KeySeven) {
			groups = 7
		}
		if rl.IsKeyDown(rl.KeyEight) {
			groups = 8
		}
		if rl.IsKeyDown(rl.KeyNine) {
			groups = 9
		}
		if rl.IsKeyDown(rl.KeyW) {
			groups = 7
			amp = 25
			speed = 18
			variance = float32(0.70)
		}
		if rl.IsKeyDown(rl.KeyEqual) {
			if float32(speed) <= float32(fps)*0.25 {
				speed = int(float32(fps) * 0.25)
			} else {
				speed = int(float32(speed) * 0.95)
			}
		}
		if rl.IsKeyDown(rl.KeyKpAdd) {
			if float32(speed) <= float32(fps)*0.25 {
				speed = int(float32(fps) * 0.25)
			} else {
				speed = int(float32(speed) * 0.95)
			}
		}
		if rl.IsKeyDown(rl.KeyMinus) {
			speed = int(math.Max(float64(speed)*1.02, float64(speed)+1))
		}
		if rl.IsKeyDown(rl.KeyKpSubtract) {
			speed = int(math.Max(float64(speed)*1.02, float64(speed)+1))
		}

		// Update the light shader with the camera view position
		rl.SetShaderValue(shader, shader.GetLocation(rl.LocVectorView),
			[]float32{camera.Position.X, camera.Position.Y, camera.Position.Z}, rl.ShaderUniformVec3)

		// Apply per-instance transformations
		for i := 0; i < MAX_INSTANCES; i++ {
			rotations[i] = rl.MatrixMultiply(rotations[i], rotationsInc[i])
			transforms[i] = rl.MatrixMultiply(rotations[i], translations[i])

			// Get the animation cycle's framesCounter for this instance
			loop = float32((framesCounter+int(float32(i%groups)/float32(groups)*float32(speed)))%speed) / float32(speed)

			// Calculate the y according to loop cycle
			y := float32(math.Sin(float64(loop)*rl.Pi*2)) * amp *
				((1 - variance) + (float32(variance) * float32(i%(groups*10)) / float32(groups*10)))

			// Clamp to floor
			if y < 0 {
				y = 0
			}

			transforms[i] = rl.MatrixMultiply(transforms[i], rl.MatrixTranslate(0.0, y, 0.0))
		}

		rl.UpdateCamera(&camera) // Update camera
		//----------------------------------------------------------------------------------

		// Draw
		//----------------------------------------------------------------------------------
		rl.BeginDrawing()
		{
			rl.ClearBackground(rl.RayWhite)

			rl.BeginMode3D(camera)
			//rl.DrawMesh(cube, material, rl.MatrixIdentity())
			rl.DrawMeshInstanced(cube, material, transforms, MAX_INSTANCES)
			rl.EndMode3D()

			rl.DrawText("A CUBE OF DANCING CUBES!", 490, 10, 20, rl.Maroon)
			rl.DrawText("PRESS KEYS:", 10, textPositionY, 20, rl.Black)

			textPositionY += 25
			rl.DrawText("1 - 9", 10, textPositionY, 10, rl.Black)
			rl.DrawText(": Number of groups", 50, textPositionY, 10, rl.Black)
			rl.DrawText(fmt.Sprintf(": %d", groups), 160, textPositionY, 10, rl.Black)

			textPositionY += 15
			rl.DrawText("UP", 10, textPositionY, 10, rl.Black)
			rl.DrawText(": increase amplitude", 50, textPositionY, 10, rl.Black)
			rl.DrawText(fmt.Sprintf(": %.2f", amp), 160, textPositionY, 10, rl.Black)

			textPositionY += 15
			rl.DrawText("DOWN", 10, textPositionY, 10, rl.Black)
			rl.DrawText(": decrease amplitude", 50, textPositionY, 10, rl.Black)

			textPositionY += 15
			rl.DrawText("LEFT", 10, textPositionY, 10, rl.Black)
			rl.DrawText(": decrease variance", 50, textPositionY, 10, rl.Black)
			rl.DrawText(fmt.Sprintf(": %.2f", variance), 160, textPositionY, 10, rl.Black)

			textPositionY += 15
			rl.DrawText("RIGHT", 10, textPositionY, 10, rl.Black)
			rl.DrawText(": increase variance", 50, textPositionY, 10, rl.Black)

			textPositionY += 15
			rl.DrawText("+/=", 10, textPositionY, 10, rl.Black)
			rl.DrawText(": increase speed", 50, textPositionY, 10, rl.Black)
			rl.DrawText(fmt.Sprintf(": %d = %f loops/sec", speed, float32(fps)/float32(speed)), 160, textPositionY, 10, rl.Black)

			textPositionY += 15
			rl.DrawText("-", 10, textPositionY, 10, rl.Black)
			rl.DrawText(": decrease speed", 50, textPositionY, 10, rl.Black)

			textPositionY += 15
			rl.DrawText("W", 10, textPositionY, 10, rl.Black)
			rl.DrawText(": Wild setup!", 50, textPositionY, 10, rl.Black)

			rl.DrawFPS(10, 10)
		}
		rl.EndDrawing()
		//----------------------------------------------------------------------------------
	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	rl.CloseWindow() // Close window and OpenGL context
	//--------------------------------------------------------------------------------------
}
