package main

import (
	"github.com/gen2brain/raylib-go/physics"
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(screenWidth, screenHeight, "Physac [raylib] - physics friction")

	// Physac logo drawing position
	logoX := screenWidth - rl.MeasureText("Physac", 30) - 10
	logoY := int32(15)

	// Initialize physics and default physics bodies
	physics.Init()

	// Create floor rectangle physics body
	floor := physics.NewBodyRectangle(rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)), float32(screenHeight), 100, 10)
	floor.Enabled = false // Disable body state to convert it to static (no dynamics, but collisions)
	wall := physics.NewBodyRectangle(rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)*0.8), 10, 80, 10)
	wall.Enabled = false // Disable body state to convert it to static (no dynamics, but collisions)

	// Create left ramp physics body
	rectLeft := physics.NewBodyRectangle(rl.NewVector2(25, float32(screenHeight)-5), 250, 250, 10)
	rectLeft.Enabled = false // Disable body state to convert it to static (no dynamics, but collisions)
	rectLeft.SetRotation(30 * rl.Deg2rad)

	// Create right ramp  physics body
	rectRight := physics.NewBodyRectangle(rl.NewVector2(float32(screenWidth)-25, float32(screenHeight)-5), 250, 250, 10)
	rectRight.Enabled = false // Disable body state to convert it to static (no dynamics, but collisions)
	rectRight.SetRotation(330 * rl.Deg2rad)

	// Create dynamic physics bodies
	bodyA := physics.NewBodyRectangle(rl.NewVector2(35, float32(screenHeight)*0.6), 40, 40, 10)
	bodyA.StaticFriction = 0.1
	bodyA.DynamicFriction = 0.1
	bodyA.SetRotation(30 * rl.Deg2rad)

	bodyB := physics.NewBodyRectangle(rl.NewVector2(float32(screenWidth)-35, float32(screenHeight)*0.6), 40, 40, 10)
	bodyB.StaticFriction = 1
	bodyB.DynamicFriction = 1
	bodyB.SetRotation(330 * rl.Deg2rad)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Physics steps calculations
		physics.Update()

		if rl.IsKeyPressed(rl.KeyR) { // Reset physics input
			// Reset dynamic physics bodies position, velocity and rotation
			bodyA.Position = rl.NewVector2(35, float32(screenHeight)*0.6)
			bodyA.Velocity = rl.NewVector2(0, 0)
			bodyA.AngularVelocity = 0
			bodyA.SetRotation(30 * rl.Deg2rad)

			bodyB.Position = rl.NewVector2(float32(screenWidth)-35, float32(screenHeight)*0.6)
			bodyB.Velocity = rl.NewVector2(0, 0)
			bodyB.AngularVelocity = 0
			bodyB.SetRotation(330 * rl.Deg2rad)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		rl.DrawFPS(screenWidth-90, screenHeight-30)

		// Draw created physics bodies
		bodiesCount := physics.GetBodiesCount()
		for i := 0; i < bodiesCount; i++ {
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

		rl.DrawRectangle(0, screenHeight-49, screenWidth, 49, rl.Black)

		rl.DrawText("Friction amount", (screenWidth-rl.MeasureText("Friction amount", 30))/2, 75, 30, rl.White)
		rl.DrawText("0.1", int32(bodyA.Position.X)-rl.MeasureText("0.1", 20)/2, int32(bodyA.Position.Y)-7, 20, rl.White)
		rl.DrawText("1", int32(bodyB.Position.X)-rl.MeasureText("1", 20)/2, int32(bodyB.Position.Y)-7, 20, rl.White)

		rl.DrawText("Press 'R' to reset example", 10, 10, 10, rl.White)

		rl.DrawText("Physac", logoX, logoY, 30, rl.White)
		rl.DrawText("Powered by", logoX+50, logoY-7, 10, rl.White)

		rl.EndDrawing()
	}

	physics.Close() // Unitialize physics

	rl.CloseWindow()
}
