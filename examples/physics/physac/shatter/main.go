package main

import (
	"github.com/gen2brain/raylib-go/physics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	velocity = 0.5
)

func main() {
	screenWidth := float32(800)
	screenHeight := float32(450)

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "Physac [raylib] - body shatter")

	// Physac logo drawing position
	logoX := int32(screenWidth) - rl.MeasureText("Physac", 30) - 10
	logoY := int32(15)

	// Initialize physics and default physics bodies
	physics.Init()
	physics.SetGravity(0, 0)

	// Create random polygon physics body to shatter
	physics.NewBodyPolygon(rl.NewVector2(screenWidth/2, screenHeight/2), float32(rl.GetRandomValue(80, 200)), int(rl.GetRandomValue(3, 8)), 10)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update created physics objects
		physics.Update()

		if rl.IsKeyPressed(rl.KeyR) { // Reset physics input
			physics.Reset()

			// Create random polygon physics body to shatter
			physics.NewBodyPolygon(rl.NewVector2(screenWidth/2, screenHeight/2), float32(rl.GetRandomValue(80, 200)), int(rl.GetRandomValue(3, 8)), 10)
		}

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			//for _, b := range physics.GetBodies() {
			for i := 0; i < physics.GetBodiesCount(); i++ {
				body := physics.GetBody(i)
				if body == nil {
					continue
				}

				physics.Shatter(body, rl.GetMousePosition(), 10/body.InverseMass)
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

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

				rl.DrawLineV(vertexA, vertexB, rl.Green) // Draw a line between two vertex positions
			}
		}

		rl.DrawText("Left mouse button in polygon area to shatter body\nPress 'R' to reset example", 10, 10, 10, rl.White)

		rl.DrawText("Physac", logoX, logoY, 30, rl.White)
		rl.DrawText("Powered by", logoX+50, logoY-7, 10, rl.White)

		rl.EndDrawing()
	}

	physics.Close() // Unitialize physics

	rl.CloseWindow()
}
