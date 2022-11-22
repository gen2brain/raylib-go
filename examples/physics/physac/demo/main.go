package main

import (
	"github.com/gen2brain/raylib-go/physics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(screenWidth, screenHeight, "Physac [raylib] - physics demo")

	// Physac logo drawing position
	logoX := screenWidth - rl.MeasureText("Physac", 30) - 10
	logoY := int32(15)

	// Initialize physics and default physics bodies
	physics.Init()

	// Create floor rectangle physics body
	floor := physics.NewBodyRectangle(rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)), 500, 100, 10)
	floor.Enabled = false // Disable body state to convert it to static (no dynamics, but collisions)

	// Create obstacle circle physics body
	circle := physics.NewBodyCircle(rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2), 45, 10)
	circle.Enabled = false // Disable body state to convert it to static (no dynamics, but collisions)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		// Update created physics objects
		physics.Update()

		if rl.IsKeyPressed(rl.KeyR) { // Reset physics input
			physics.Reset()

			floor = physics.NewBodyRectangle(rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)), 500, 100, 10)
			floor.Enabled = false

			circle = physics.NewBodyCircle(rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2), 45, 10)
			circle.Enabled = false
		}

		// Physics body creation inputs
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			physics.NewBodyPolygon(rl.GetMousePosition(), float32(rl.GetRandomValue(20, 80)), int(rl.GetRandomValue(3, 8)), 10)
		} else if rl.IsMouseButtonPressed(rl.MouseRightButton) {
			physics.NewBodyCircle(rl.GetMousePosition(), float32(rl.GetRandomValue(10, 45)), 10)
		}

		// Destroy falling physics bodies
		for i := 0; i < physics.GetBodiesCount(); i++ {
			body := physics.GetBody(i)
			if body.Position.Y > float32(screenHeight)*2 {
				body.Destroy()
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		rl.DrawFPS(screenWidth-90, screenHeight-30)

		// Draw created physics bodies
		for i := 0; i < physics.GetBodiesCount(); i++ {
			body := physics.GetBody(i)

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

		rl.DrawText("Left mouse button to create a polygon", 10, 10, 10, rl.White)
		rl.DrawText("Right mouse button to create a circle", 10, 25, 10, rl.White)
		rl.DrawText("Press 'R' to reset example", 10, 40, 10, rl.White)

		rl.DrawText("Physac", logoX, logoY, 30, rl.White)
		rl.DrawText("Powered by", logoX+50, logoY-7, 10, rl.White)

		rl.EndDrawing()
	}

	physics.Close() // Unitialize physics

	rl.CloseWindow()
}
