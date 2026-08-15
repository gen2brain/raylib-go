package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// A few good julia sets
var pointsOfInterest = [6][2]float32{
	{-0.348827, 0.607167},
	{-0.786268, 0.169728},
	{-0.8, 0.156},
	{0.285, 0.0},
	{-0.835, -0.2321},
	{-0.70176, -0.3842},
}

const (
	screenWidth    = 800
	screenHeight   = 450
	zoomSpeed      = 1.01
	offsetSpeedMul = 2.0
	startingZoom   = 0.75
)

func main() {
	// Initialization
	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - julia set")

	shader := rl.LoadShader("", "julia_set.fs")
	target := rl.LoadRenderTexture(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()))
	cVal := []float32{pointsOfInterest[0][0], pointsOfInterest[0][1]}

	// Offset and zoom to draw the julia set at
	offset := []float32{0.0, 0.0}
	zoom := float32(startingZoom)

	// Get variable (uniform) locations on the shader
	cLoc := rl.GetShaderLocation(shader, "c")
	zoomLoc := rl.GetShaderLocation(shader, "zoom")
	offsetLoc := rl.GetShaderLocation(shader, "offset")

	// Upload initial shader uniform values
	rl.SetShaderValue(shader, cLoc, cVal, rl.ShaderUniformVec2)
	rl.SetShaderValue(shader, zoomLoc, []float32{zoom}, rl.ShaderUniformFloat)
	rl.SetShaderValue(shader, offsetLoc, offset, rl.ShaderUniformVec2)

	incrementSpeed := 0
	showControls := true

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		// Press [1 - 6] to reset c to a point of interest
		if rl.IsKeyPressed(rl.KeyOne) {
			cVal[0] = pointsOfInterest[0][0]
			cVal[1] = pointsOfInterest[0][1]
			rl.SetShaderValue(shader, cLoc, cVal, rl.ShaderUniformVec2)
		} else if rl.IsKeyPressed(rl.KeyTwo) {
			cVal[0] = pointsOfInterest[1][0]
			cVal[1] = pointsOfInterest[1][1]
			rl.SetShaderValue(shader, cLoc, cVal, rl.ShaderUniformVec2)
		} else if rl.IsKeyPressed(rl.KeyThree) {
			cVal[0] = pointsOfInterest[2][0]
			cVal[1] = pointsOfInterest[2][1]
			rl.SetShaderValue(shader, cLoc, cVal, rl.ShaderUniformVec2)
		} else if rl.IsKeyPressed(rl.KeyFour) {
			cVal[0] = pointsOfInterest[3][0]
			cVal[1] = pointsOfInterest[3][1]
			rl.SetShaderValue(shader, cLoc, cVal, rl.ShaderUniformVec2)
		} else if rl.IsKeyPressed(rl.KeyFive) {
			cVal[0] = pointsOfInterest[4][0]
			cVal[1] = pointsOfInterest[4][1]
			rl.SetShaderValue(shader, cLoc, cVal, rl.ShaderUniformVec2)
		} else if rl.IsKeyPressed(rl.KeySix) {
			cVal[0] = pointsOfInterest[5][0]
			cVal[1] = pointsOfInterest[5][1]
			rl.SetShaderValue(shader, cLoc, cVal, rl.ShaderUniformVec2)
		}

		// If "R" is pressed, reset zoom and offset
		if rl.IsKeyPressed(rl.KeyR) {
			zoom = startingZoom
			offset[0] = 0.0
			offset[1] = 0.0
			rl.SetShaderValue(shader, zoomLoc, []float32{zoom}, rl.ShaderUniformFloat)
			rl.SetShaderValue(shader, offsetLoc, offset, rl.ShaderUniformVec2)
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			incrementSpeed = 0 // Pause animation
		}

		if rl.IsKeyPressed(rl.KeyF1) {
			showControls = !showControls // Toggle whether or not to show controls
		}

		if rl.IsKeyPressed(rl.KeyRight) {
			incrementSpeed++
		} else if rl.IsKeyPressed(rl.KeyLeft) {
			incrementSpeed--
		}

		// If either left or right button is pressed, zoom in/out
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) || rl.IsMouseButtonDown(rl.MouseButtonRight) {
			// Change zoom. If Mouse left -> zoom in. Mouse right -> zoom out
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				zoom *= zoomSpeed
			} else {
				zoom *= 1.0 / zoomSpeed
			}

			mousePos := rl.GetMousePosition()
			offsetVelocityX := (mousePos.X/float32(screenWidth) - 0.5) * offsetSpeedMul / zoom
			offsetVelocityY := (mousePos.Y/float32(screenHeight) - 0.5) * offsetSpeedMul / zoom

			// Apply move velocity to camera
			dt := rl.GetFrameTime()
			offset[0] += dt * offsetVelocityX
			offset[1] += dt * offsetVelocityY

			// Update the shader uniform values
			rl.SetShaderValue(shader, zoomLoc, []float32{zoom}, rl.ShaderUniformFloat)
			rl.SetShaderValue(shader, offsetLoc, offset, rl.ShaderUniformVec2)
		}

		// Increment c value with time
		dc := rl.GetFrameTime() * float32(incrementSpeed) * 0.0005
		cVal[0] += dc
		cVal[1] += dc
		rl.SetShaderValue(shader, cLoc, cVal, rl.ShaderUniformVec2)

		// Draw
		// Using a render texture to draw Julia set
		rl.BeginTextureMode(target)

		rl.ClearBackground(rl.Black)

		// Draw a rectangle in shader mode to be used as shader canvas
		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Black)

		rl.EndTextureMode()

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		// Draw saved texture and rendered julia set with shader
		rl.BeginShaderMode(shader)

		rl.DrawTextureEx(target.Texture, rl.Vector2Zero(), 0.0, 1.0, rl.White)

		rl.EndShaderMode()

		if showControls {
			rl.DrawText("Press Mouse buttons right/left to zoom in/out and move", 10, 15, 10, rl.RayWhite)
			rl.DrawText("Press KEY_F1 to toggle these controls", 10, 30, 10, rl.RayWhite)
			rl.DrawText("Press KEYS [1 - 6] to change point of interest", 10, 45, 10, rl.RayWhite)
			rl.DrawText("Press KEY_LEFT | KEY_RIGHT to change speed", 10, 60, 10, rl.RayWhite)
			rl.DrawText("Press KEY_SPACE to stop movement animation", 10, 75, 10, rl.RayWhite)
			rl.DrawText("Press KEY_R to recenter the camera", 10, 90, 10, rl.RayWhite)
		}

		rl.EndDrawing()
	}
	rl.UnloadShader(shader)
	rl.UnloadRenderTexture(target)

	rl.CloseWindow()
}
