package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	maxColumns   = 20
	screenWidth  = 800
	screenHeight = 450
	False        = 0
	True         = 1
)

var camModes = map[rl.CameraMode]string{
	rl.CameraFree:        "FREE",
	rl.CameraOrbital:     "ORBITAL",
	rl.CameraFirstPerson: "FIRST PERSON",
	rl.CameraThirdPerson: "THIRD PERSON",
}

var camProjections = map[rl.CameraProjection]string{
	rl.CameraPerspective:  "PERSPECTIVE",
	rl.CameraOrthographic: "ORTHOGRAPHIC",
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - 3d camera first person")

	// Define the camera to look into our 3d world (position, target, up vector)
	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(0.0, 2.0, 4.0) // Camera position
	camera.Target = rl.NewVector3(0.0, 2.0, 0.0)   // Camera looking at point
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)       // Camera up vector (rotation towards target)
	camera.Fovy = 60.0                             // Camera field-of-view Y
	camera.Projection = rl.CameraPerspective       // Camera projection type

	cameraMode := rl.CameraFirstPerson

	// Generates some random columns
	heights := make([]float32, maxColumns)
	positions := make([]rl.Vector3, maxColumns)
	colors := make([]rl.Color, maxColumns)

	for i := 0; i < maxColumns; i++ {
		heights[i] = rndF(1, 12)
		positions[i] = rl.NewVector3(rndF(-15, 15), heights[i]/2, rndF(-15, 15))
		colors[i] = rl.NewColor(rndU(20, 255), rndU(10, 55), 30, 255)
	}

	rl.DisableCursor()  // Limit cursor to relative movement inside the window
	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Switch camera mode
		if rl.IsKeyPressed(rl.KeyOne) {
			cameraMode = rl.CameraFree
			camera.Up = rl.NewVector3(0.0, 1.0, 0.0) // Reset roll
		}
		if rl.IsKeyPressed(rl.KeyTwo) {
			cameraMode = rl.CameraFirstPerson
			camera.Up = rl.NewVector3(0.0, 1.0, 0.0) // Reset roll
		}
		if rl.IsKeyPressed(rl.KeyThree) {
			cameraMode = rl.CameraThirdPerson
			camera.Up = rl.NewVector3(0.0, 1.0, 0.0) // Reset roll
		}
		if rl.IsKeyPressed(rl.KeyFour) {
			cameraMode = rl.CameraOrbital
			camera.Up = rl.NewVector3(0.0, 1.0, 0.0) // Reset roll
		}

		// Switch camera projection
		if rl.IsKeyPressed(rl.KeyP) {
			cameraMode = rl.CameraThirdPerson
			if camera.Projection == rl.CameraPerspective {
				// Create isometric view
				// Note: The target distance is related to the render distance in the orthographic projection
				camera.Position = rl.NewVector3(0.0, 2.0, -100.0)
				camera.Target = rl.NewVector3(0.0, 2.0, 0.0)
				camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
				camera.Projection = rl.CameraOrthographic
				camera.Fovy = 20.0 // near plane width in CAMERA_ORTHOGRAPHIC
				rl.CameraYaw(&camera, -135*rl.Deg2rad, True)
				rl.CameraPitch(&camera, -45*rl.Deg2rad, True, True, False)
			} else if camera.Projection == rl.CameraOrthographic {
				// Reset to default view
				camera.Position = rl.NewVector3(0.0, 2.0, 10.0)
				camera.Target = rl.NewVector3(0.0, 2.0, 0.0)
				camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
				camera.Projection = rl.CameraPerspective
				camera.Fovy = 60.0
			}
		}

		// Update camera computes movement internally depending on the camera mode.
		// Some default standard keyboard/mouse inputs are hardcoded to simplify use.
		// For advance camera controls, it's recommended to compute camera movement manually.
		rl.UpdateCamera(&camera, cameraMode) // Update camera

		/*
		   // Camera PRO usage example (EXPERIMENTAL)
		   // This new camera function allows custom movement/rotation values to be directly provided
		   // as input parameters, with this approach, rcamera module is internally independent of raylib inputs
		   UpdateCameraPro(&camera,
		       (Vector3){
		           (IsKeyDown(KEY_W) || IsKeyDown(KEY_UP))*0.1f -      // Move forward-backward
		           (IsKeyDown(KEY_S) || IsKeyDown(KEY_DOWN))*0.1f,
		           (IsKeyDown(KEY_D) || IsKeyDown(KEY_RIGHT))*0.1f -   // Move right-left
		           (IsKeyDown(KEY_A) || IsKeyDown(KEY_LEFT))*0.1f,
		           0.0f                                                // Move up-down
		       },
		       (Vector3){
		           GetMouseDelta().x*0.05f,                            // Rotation: yaw
		           GetMouseDelta().y*0.05f,                            // Rotation: pitch
		           0.0f                                                // Rotation: roll
		       },
		       GetMouseWheelMove()*2.0f);                              // Move to target (zoom)
		*/

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawPlane(rl.NewVector3(0.0, 0.0, 0.0), rl.NewVector2(32.0, 32.0), rl.LightGray) // Draw ground
		rl.DrawCube(rl.NewVector3(-16.0, 2.5, 0.0), 1.0, 5.0, 32.0, rl.Blue)                // Draw a blue wall
		rl.DrawCube(rl.NewVector3(16.0, 2.5, 0.0), 1.0, 5.0, 32.0, rl.Lime)                 // Draw a green wall
		rl.DrawCube(rl.NewVector3(0.0, 2.5, 16.0), 32.0, 5.0, 1.0, rl.Gold)                 // Draw a yellow wall

		// Draw some cubes around
		for i := 0; i < maxColumns; i++ {
			rl.DrawCube(positions[i], 2.0, heights[i], 2.0, colors[i])
			rl.DrawCubeWires(positions[i], 2.0, heights[i], 2.0, rl.Maroon)
		}

		// Draw player cube
		if cameraMode == rl.CameraThirdPerson {
			rl.DrawCube(camera.Target, 0.5, 0.5, 0.5, rl.Purple)
			rl.DrawCubeWires(camera.Target, 0.5, 0.5, 0.5, rl.DarkPurple)
		}

		rl.EndMode3D()

		// Draw info boxes
		rl.DrawRectangle(5, 5, 330, 100, rl.Fade(rl.SkyBlue, 0.5))
		rl.DrawRectangleLines(5, 5, 330, 100, rl.Blue)

		rl.DrawText("Camera controls:", 15, 15, 10, rl.Black)
		rl.DrawText("- Move keys: W, A, S, D, Space, Left-Ctrl", 15, 30, 10, rl.Black)
		rl.DrawText("- Look around: arrow keys or mouse", 15, 45, 10, rl.Black)
		rl.DrawText("- Camera mode keys: 1,2,3,4", 15, 60, 10, rl.Black)
		rl.DrawText("- Zoom keys: num-plus, num-minus or mouse scroll", 15, 75, 10, rl.Black)
		rl.DrawText("- Camera projection key: P", 15, 90, 10, rl.Black)

		rl.DrawRectangle(600, 5, 195, 100, rl.Fade(rl.SkyBlue, 0.5))
		rl.DrawRectangleLines(600, 5, 195, 100, rl.Blue)

		rl.DrawText("Camera status:", 610, 15, 10, rl.Black)
		rl.DrawText(s2T("Mode:", camModes[cameraMode]), 610, 30, 10, rl.Black)
		rl.DrawText(s2T("Projection:", camProjections[camera.Projection]), 610, 45, 10, rl.Black)
		rl.DrawText(s2T("Position:", v2T(camera.Position)), 610, 60, 10, rl.Black)
		rl.DrawText(s2T("Target:", v2T(camera.Target)), 610, 75, 10, rl.Black)
		rl.DrawText(s2T("Up:", v2T(camera.Up)), 610, 90, 10, rl.Black)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

// rndU generates a random uint8 value between min and max
func rndU(min, max int32) uint8 {
	return uint8(rl.GetRandomValue(min, max))
}

// rndF generates a random float32 value between min and max
func rndF(min, max int32) float32 {
	return float32(rl.GetRandomValue(min, max))
}

// s2t generates a Status item string from a name and value
func s2T(name, value string) string {
	return fmt.Sprintf(" - %s %s", name, value)
}

// v2T generates a string from a rl.Vector3
func v2T(v rl.Vector3) string {
	return fmt.Sprintf("%6.3f, %6.3f, %6.3f", v.X, v.Y, v.Z)
}
