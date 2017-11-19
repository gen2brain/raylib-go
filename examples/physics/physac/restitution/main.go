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
	raylib.InitWindow(int32(screenWidth), int32(screenHeight), "Physac [raylib] - physics restitution")

	// Physac logo drawing position
	logoX := int32(screenWidth) - raylib.MeasureText("Physac", 30) - 10
	logoY := int32(15)

	// Initialize physics and default physics bodies
	physics.Init()

	// Create floor rectangle physics body
	floor := physics.NewBodyRectangle(raylib.NewVector2(screenWidth/2, screenHeight), screenWidth, 100, 10)
	floor.Enabled = false // Disable body state to convert it to static (no dynamics, but collisions)
	floor.Restitution = 1

	// Create circles physics body
	circleA := physics.NewBodyCircle(raylib.NewVector2(screenWidth*0.25, screenHeight/2), 30, 10)
	circleA.Restitution = 0
	circleB := physics.NewBodyCircle(raylib.NewVector2(screenWidth*0.5, screenHeight/2), 30, 10)
	circleB.Restitution = 0.5
	circleC := physics.NewBodyCircle(raylib.NewVector2(screenWidth*0.75, screenHeight/2), 30, 10)
	circleC.Restitution = 1

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update created physics objects
		physics.Update()

		if raylib.IsKeyPressed(raylib.KeyR) { // Reset physics input
			// Reset circles physics bodies position and velocity
			circleA.Position = raylib.NewVector2(screenWidth*0.25, screenHeight/2)
			circleA.Velocity = raylib.NewVector2(0, 0)
			circleB.Position = raylib.NewVector2(screenWidth*0.5, screenHeight/2)
			circleB.Velocity = raylib.NewVector2(0, 0)
			circleC.Position = raylib.NewVector2(screenWidth*0.75, screenHeight/2)
			circleC.Velocity = raylib.NewVector2(0, 0)
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

		raylib.DrawText("Restitution amount", (int32(screenWidth)-raylib.MeasureText("Restitution amount", 30))/2, 75, 30, raylib.White)
		raylib.DrawText("0", int32(circleA.Position.X)-raylib.MeasureText("0", 20)/2, int32(circleA.Position.Y)-7, 20, raylib.White)
		raylib.DrawText("0.5", int32(circleB.Position.X)-raylib.MeasureText("0.5", 20)/2, int32(circleB.Position.Y)-7, 20, raylib.White)
		raylib.DrawText("1", int32(circleC.Position.X)-raylib.MeasureText("1", 20)/2, int32(circleC.Position.Y)-7, 20, raylib.White)

		raylib.DrawText("Press 'R' to reset example", 10, 10, 10, raylib.White)

		raylib.DrawText("Physac", logoX, logoY, 30, raylib.White)
		raylib.DrawText("Powered by", logoX+50, logoY-7, 10, raylib.White)

		raylib.EndDrawing()
	}

	physics.Close() // Unitialize physics

	raylib.CloseWindow()
}
