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
	raylib.InitWindow(int32(screenWidth), int32(screenHeight), "Physac [raylib] - body shatter")

	// Physac logo drawing position
	logoX := int32(screenWidth) - raylib.MeasureText("Physac", 30) - 10
	logoY := int32(15)

	// Initialize physics and default physics bodies
	physics.Init()
	physics.SetGravity(0, 0)

	// Create random polygon physics body to shatter
	physics.NewBodyPolygon(raylib.NewVector2(screenWidth/2, screenHeight/2), float32(raylib.GetRandomValue(80, 200)), int(raylib.GetRandomValue(3, 8)), 10)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update created physics objects
		physics.Update()

		if raylib.IsKeyPressed(raylib.KeyR) { // Reset physics input
			physics.Reset()

			// Create random polygon physics body to shatter
			physics.NewBodyPolygon(raylib.NewVector2(screenWidth/2, screenHeight/2), float32(raylib.GetRandomValue(80, 200)), int(raylib.GetRandomValue(3, 8)), 10)
		}

		if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
			for _, b := range physics.GetBodies() {
				b.Shatter(raylib.GetMousePosition(), 10/b.InverseMass)
			}
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.Black)

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

		raylib.DrawText("Left mouse button in polygon area to shatter body\nPress 'R' to reset example", 10, 10, 10, raylib.White)

		raylib.DrawText("Physac", logoX, logoY, 30, raylib.White)
		raylib.DrawText("Powered by", logoX+50, logoY-7, 10, raylib.White)

		raylib.EndDrawing()
	}

	physics.Close() // Unitialize physics

	raylib.CloseWindow()
}
