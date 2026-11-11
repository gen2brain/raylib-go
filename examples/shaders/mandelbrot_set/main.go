package main

import (
	"math"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// A few good interesting places
var pointsOfInterest = [6][3]float32{
	{-1.76826775, -0.00422996283, 28435.9238},
	{0.322004497, -0.0357099883, 56499.7266},
	{-0.748880744, -0.0562955774, 9237.59082},
	{-1.78385007, -0.0156200649, 14599.5283},
	{-0.0985441282, -0.924688697, 26259.8535},
	{0.317785531, -0.0322612226, 29297.9258},
}

const (
	screenWidth    = 800
	screenHeight   = 450
	zoomSpeed      = 1.01
	offsetSpeedMul = 2.0
	startingZoom   = 0.6
)

var startingOffset = [2]float32{-0.5, 0.0}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - mandelbrot set")
	shader := rl.LoadShader("", "mandelbrot_set.fs")
	target := rl.LoadRenderTexture(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()))

	// Offset and zoom to draw the mandelbrot set at
	offset := []float32{startingOffset[0], startingOffset[1]}
	zoom := float32(startingZoom)
	maxIterations := int32(333)
	maxIterationsMultiplier := float32(166.5)

	// Get variable (uniform) locations on the shader
	zoomLoc := rl.GetShaderLocation(shader, "zoom")
	offsetLoc := rl.GetShaderLocation(shader, "offset")
	maxIterationsLoc := rl.GetShaderLocation(shader, "maxIterations")

	// Upload initial shader uniform values
	rl.SetShaderValue(shader, zoomLoc, []float32{zoom}, rl.ShaderUniformFloat)
	rl.SetShaderValue(shader, offsetLoc, offset, rl.ShaderUniformVec2)
	rl.SetShaderValue(shader, maxIterationsLoc, unsafe.Slice((*float32)(unsafe.Pointer(&maxIterations)), 1), rl.ShaderUniformInt)

	showControls := true

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		updateShader := false

		// Press [1 - 6] to reset to a point of interest
		if rl.IsKeyPressed(rl.KeyOne) ||
			rl.IsKeyPressed(rl.KeyTwo) ||
			rl.IsKeyPressed(rl.KeyThree) ||
			rl.IsKeyPressed(rl.KeyFour) ||
			rl.IsKeyPressed(rl.KeyFive) ||
			rl.IsKeyPressed(rl.KeySix) {
			interestIndex := 0
			if rl.IsKeyPressed(rl.KeyOne) {
				interestIndex = 0
			} else if rl.IsKeyPressed(rl.KeyTwo) {
				interestIndex = 1
			} else if rl.IsKeyPressed(rl.KeyThree) {
				interestIndex = 2
			} else if rl.IsKeyPressed(rl.KeyFour) {
				interestIndex = 3
			} else if rl.IsKeyPressed(rl.KeyFive) {
				interestIndex = 4
			} else if rl.IsKeyPressed(rl.KeySix) {
				interestIndex = 5
			}

			offset[0] = pointsOfInterest[interestIndex][0]
			offset[1] = pointsOfInterest[interestIndex][1]
			zoom = pointsOfInterest[interestIndex][2]
			updateShader = true
		}

		// If "R" is pressed, reset zoom and offset
		if rl.IsKeyPressed(rl.KeyR) {
			offset[0] = startingOffset[0]
			offset[1] = startingOffset[1]
			zoom = startingZoom
			updateShader = true
		}

		if rl.IsKeyPressed(rl.KeyF1) {
			showControls = !showControls // Toggle whether or not to show controls
		}

		// Change number of max iterations with UP and DOWN keys
		if rl.IsKeyPressed(rl.KeyUp) {
			maxIterationsMultiplier *= 1.4
			updateShader = true
		} else if rl.IsKeyPressed(rl.KeyDown) {
			maxIterationsMultiplier /= 1.4
			updateShader = true
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

			updateShader = true
		}

		// Update shader uniform values when parameters change
		if updateShader {
			inner := float64(1.0 - float32(math.Sqrt(float64(37.5*zoom))))
			absInner := math.Abs(inner)
			sqrtInner := math.Sqrt(absInner)
			iter := float32(math.Sqrt(2.0*sqrtInner)) * maxIterationsMultiplier
			maxIterations = int32(iter)

			rl.SetShaderValue(shader, zoomLoc, []float32{zoom}, rl.ShaderUniformFloat)
			rl.SetShaderValue(shader, offsetLoc, offset, rl.ShaderUniformVec2)
			rl.SetShaderValue(shader, maxIterationsLoc, unsafe.Slice((*float32)(unsafe.Pointer(&maxIterations)), 1), rl.ShaderUniformInt)
		}

		// Draw
		// Using a render texture to draw Mandelbrot set
		rl.BeginTextureMode(target)

		rl.ClearBackground(rl.Black)

		// Draw a rectangle in shader mode to be used as shader canvas
		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Black)

		rl.EndTextureMode()

		BeginDrawing := rl.BeginDrawing
		BeginDrawing()

		rl.ClearBackground(rl.Black)

		// Draw the saved texture and rendered mandelbrot set with shader
		rl.BeginShaderMode(shader)

		rl.DrawTextureEx(target.Texture, rl.Vector2Zero(), 0.0, 1.0, rl.White)

		rl.EndShaderMode()

		if showControls {
			rl.DrawText("Press Mouse buttons right/left to zoom in/out and move", 10, 15, 10, rl.RayWhite)
			rl.DrawText("Press F1 to toggle these controls", 10, 30, 10, rl.RayWhite)
			rl.DrawText("Press [1 - 6] to change point of interest", 10, 45, 10, rl.RayWhite)
			rl.DrawText("Press UP | DOWN to change number of iterations", 10, 60, 10, rl.RayWhite)
			rl.DrawText("Press R to recenter the camera", 10, 75, 10, rl.RayWhite)
		}

		rl.EndDrawing()

	}
	rl.UnloadRenderTexture(target)
	rl.UnloadShader(shader)

	rl.CloseWindow()
}
