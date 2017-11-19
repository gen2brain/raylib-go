package main

import (
	"github.com/gen2brain/raylib-go/physics"
	"github.com/gen2brain/raylib-go/raylib"
)

const (
	velocity = 0.5
)

func main() {
	screenWidth := float32(800)
	screenHeight := float32(450)

	raylib.SetConfigFlags(raylib.FlagMsaa4xHint)
	raylib.InitWindow(int32(screenWidth), int32(screenHeight), "Physac [raylib] - physics movement")

	// Physac logo drawing position
	logoX := int32(screenWidth) - raylib.MeasureText("Physac", 30) - 10
	logoY := int32(15)

	// Initialize physics and default physics bodies
	physics.Init()

	// Create floor and walls rectangle physics body
	floor := physics.NewBodyRectangle(raylib.NewVector2(screenWidth/2, screenHeight), screenWidth, 100, 10)
	platformLeft := physics.NewBodyRectangle(raylib.NewVector2(screenWidth*0.25, screenHeight*0.6), screenWidth*0.25, 10, 10)
	platformRight := physics.NewBodyRectangle(raylib.NewVector2(screenWidth*0.75, screenHeight*0.6), screenWidth*0.25, 10, 10)
	wallLeft := physics.NewBodyRectangle(raylib.NewVector2(-5, screenHeight/2), 10, screenHeight, 10)
	wallRight := physics.NewBodyRectangle(raylib.NewVector2(screenWidth+5, screenHeight/2), 10, screenHeight, 10)

	// Disable dynamics to floor and walls physics bodies
	floor.Enabled = false
	platformLeft.Enabled = false
	platformRight.Enabled = false
	wallLeft.Enabled = false
	wallRight.Enabled = false

	// Create movement physics body
	body := physics.NewBodyRectangle(raylib.NewVector2(screenWidth/2, screenHeight/2), 50, 50, 1)
	body.FreezeOrient = true // Constrain body rotation to avoid little collision torque amounts

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update created physics objects
		physics.Update()

		if raylib.IsKeyPressed(raylib.KeyR) { // Reset physics input
			// Reset movement physics body position, velocity and rotation
			body.Position = raylib.NewVector2(screenWidth/2, screenHeight/2)
			body.Velocity = raylib.NewVector2(0, 0)
			body.SetRotation(0)
		}

		// Physics body creation inputs
		if raylib.IsKeyDown(raylib.KeyRight) {
			body.Velocity.X = velocity
		} else if raylib.IsKeyDown(raylib.KeyLeft) {
			body.Velocity.X = -velocity
		}

		if raylib.IsKeyDown(raylib.KeyUp) && body.IsGrounded {
			body.Velocity.Y = -velocity * 4
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.Black)

		raylib.DrawFPS(int32(screenWidth)-90, int32(screenHeight)-30)

		// Draw created physics bodies
		for i, body := range physics.GetBodies() {
			vertexCount := physics.GetShapeVerticesCount(i)
			for j := 0; j < vertexCount; j++ {
				// Get physics bodies shape vertices to draw lines
				// NOTE: GetShapeVertex() already calculates rotation transformations
				vertexA := body.GetShapeVertex(j)

				jj := 0
				if j+1 < vertexCount { // Get next vertex or first to close the shape
					jj = j + 1
				}

				vertexB := body.GetShapeVertex(jj)

				raylib.DrawLineV(vertexA, vertexB, raylib.Green) // Draw a line between two vertex positions
			}
		}

		raylib.DrawText("Use 'ARROWS' to move player", 10, 10, 10, raylib.White)
		raylib.DrawText("Press 'R' to reset example", 10, 30, 10, raylib.White)

		raylib.DrawText("Physac", logoX, logoY, 30, raylib.White)
		raylib.DrawText("Powered by", logoX+50, logoY-7, 10, raylib.White)

		raylib.EndDrawing()
	}

	physics.Close() // Unitialize physics

	raylib.CloseWindow()
}
