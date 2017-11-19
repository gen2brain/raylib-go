package main

import (
	"github.com/gen2brain/raylib-go/physics"
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.SetConfigFlags(raylib.FlagMsaa4xHint)
	raylib.InitWindow(screenWidth, screenHeight, "Physac [raylib] - physics demo")

	// Physac logo drawing position
	logoX := screenWidth - raylib.MeasureText("Physac", 30) - 10
	logoY := int32(15)

	// Initialize physics and default physics bodies
	physics.Init()

	// Create floor rectangle physics body
	floor := physics.NewBodyRectangle(raylib.NewVector2(float32(screenWidth)/2, float32(screenHeight)), 500, 100, 10)
	floor.Enabled = false // Disable body state to convert it to static (no dynamics, but collisions)

	// Create obstacle circle physics body
	circle := physics.NewBodyCircle(raylib.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2), 45, 10)
	circle.Enabled = false // Disable body state to convert it to static (no dynamics, but collisions)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {

		// Update created physics objects
		physics.Update()

		if raylib.IsKeyPressed(raylib.KeyR) { // Reset physics input
			physics.Reset()

			floor = physics.NewBodyRectangle(raylib.NewVector2(float32(screenWidth)/2, float32(screenHeight)), 500, 100, 10)
			floor.Enabled = false

			circle = physics.NewBodyCircle(raylib.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2), 45, 10)
			circle.Enabled = false
		}

		// Physics body creation inputs
		if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
			physics.NewBodyPolygon(raylib.GetMousePosition(), float32(raylib.GetRandomValue(20, 80)), int(raylib.GetRandomValue(3, 8)), 10)
		} else if raylib.IsMouseButtonPressed(raylib.MouseRightButton) {
			physics.NewBodyCircle(raylib.GetMousePosition(), float32(raylib.GetRandomValue(10, 45)), 10)
		}

		// Destroy falling physics bodies
		for _, body := range physics.GetBodies() {
			if body.Position.Y > float32(screenHeight)*2 {
				physics.DestroyBody(body)
			}
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.Black)

		raylib.DrawFPS(screenWidth-90, screenHeight-30)

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

		raylib.DrawText("Left mouse button to create a polygon", 10, 10, 10, raylib.White)
		raylib.DrawText("Right mouse button to create a circle", 10, 25, 10, raylib.White)
		raylib.DrawText("Press 'R' to reset example", 10, 40, 10, raylib.White)

		raylib.DrawText("Physac", logoX, logoY, 30, raylib.White)
		raylib.DrawText("Powered by", logoX+50, logoY-7, 10, raylib.White)

		raylib.EndDrawing()
	}

	physics.Close() // Unitialize physics

	raylib.CloseWindow()
}
