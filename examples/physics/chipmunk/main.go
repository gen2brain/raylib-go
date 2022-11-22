package main

import (
	"fmt"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jakecoffman/cp"
)

var grabbableMaskBit uint = 1 << 31
var grabFilter = cp.ShapeFilter{
	cp.NO_GROUP, grabbableMaskBit, grabbableMaskBit,
}

func randUnitCircle() cp.Vector {
	v := cp.Vector{X: rand.Float64()*2.0 - 1.0, Y: rand.Float64()*2.0 - 1.0}
	if v.LengthSq() < 1.0 {
		return v
	}
	return randUnitCircle()
}

var simpleTerrainVerts = []cp.Vector{
	{350.00, 425.07}, {336.00, 436.55}, {272.00, 435.39}, {258.00, 427.63}, {225.28, 420.00}, {202.82, 396.00},
	{191.81, 388.00}, {189.00, 381.89}, {173.00, 380.39}, {162.59, 368.00}, {150.47, 319.00}, {128.00, 311.55},
	{119.14, 286.00}, {126.84, 263.00}, {120.56, 227.00}, {141.14, 178.00}, {137.52, 162.00}, {146.51, 142.00},
	{156.23, 136.00}, {158.00, 118.27}, {170.00, 100.77}, {208.43, 84.00}, {224.00, 69.65}, {249.30, 68.00},
	{257.00, 54.77}, {363.00, 45.94}, {374.15, 54.00}, {386.00, 69.60}, {413.00, 70.73}, {456.00, 84.89},
	{468.09, 99.00}, {467.09, 123.00}, {464.92, 135.00}, {469.00, 141.03}, {497.00, 148.67}, {513.85, 180.00},
	{509.56, 223.00}, {523.51, 247.00}, {523.00, 277.00}, {497.79, 311.00}, {478.67, 348.00}, {467.90, 360.00},
	{456.76, 382.00}, {432.95, 389.00}, {417.00, 411.32}, {373.00, 433.19}, {361.00, 430.02}, {350.00, 425.07},
}

// creates a circle with random placement
func addCircle(space *cp.Space, radius float64) {
	mass := radius * radius / 25.0
	body := space.AddBody(cp.NewBody(mass, cp.MomentForCircle(mass, 0, radius, cp.Vector{})))
	body.SetPosition(randUnitCircle().Mult(180))

	shape := space.AddShape(cp.NewCircle(body, radius, cp.Vector{}))
	shape.SetElasticity(0)
	shape.SetFriction(0.9)
}

// creates a simple terrain to contain bodies
func simpleTerrain() *cp.Space {
	space := cp.NewSpace()
	space.Iterations = 10
	space.SetGravity(cp.Vector{0, -100})
	space.SetCollisionSlop(0.5)

	offset := cp.Vector{X: -320, Y: -240}
	for i := 0; i < len(simpleTerrainVerts)-1; i++ {
		a := simpleTerrainVerts[i]
		b := simpleTerrainVerts[i+1]
		space.AddShape(cp.NewSegment(space.StaticBody, a.Add(offset), b.Add(offset), 0))
	}

	return space
}

func main() {
	const width, height = 800, 450
	const physicsTickrate = 1.0 / 60.0

	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(width, height, "raylib [physics] example - chipmunk")

	offset := rl.Vector2{X: width / 2, Y: height / 2}
	// since the example ported from elsewhere, flip the camera 180 and offset to center it
	camera := rl.NewCamera2D(offset, rl.Vector2{}, 180, 1)

	space := simpleTerrain()
	for i := 0; i < 1000; i++ {
		addCircle(space, 5)
	}
	mouseBody := cp.NewKinematicBody()
	var mouse cp.Vector
	var mouseJoint *cp.Constraint

	var accumulator, dt float32
	lastTime := rl.GetTime()
	for !rl.WindowShouldClose() {
		// calculate dt
		now := rl.GetTime()
		dt = float32(now - lastTime)
		lastTime = now

		// update the mouse position
		mousePos := rl.GetMousePosition()
		// alter the mouse coordinates based on the camera position, rotation
		mouse.X = float64(mousePos.X-camera.Offset.X) * -1
		mouse.Y = float64(mousePos.Y-camera.Offset.Y) * -1
		// smooth mouse movements to new position
		newPoint := mouseBody.Position().Lerp(mouse, 0.25)
		mouseBody.SetVelocityVector(newPoint.Sub(mouseBody.Position()).Mult(60.0))
		mouseBody.SetPosition(newPoint)

		// handle grabbing
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			result := space.PointQueryNearest(mouse, 5, grabFilter)
			if result.Shape != nil && result.Shape.Body().Mass() < cp.INFINITY {
				var nearest cp.Vector
				if result.Distance > 0 {
					nearest = result.Point
				} else {
					nearest = mouse
				}

				// create a new constraint where the mouse is to draw the body towards the mouse
				body := result.Shape.Body()
				mouseJoint = cp.NewPivotJoint2(mouseBody, body, cp.Vector{}, body.WorldToLocal(nearest))
				mouseJoint.SetMaxForce(50000)
				mouseJoint.SetErrorBias(math.Pow(1.0-0.15, 60.0))
				space.AddConstraint(mouseJoint)
			}
		} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) && mouseJoint != nil {
			space.RemoveConstraint(mouseJoint)
			mouseJoint = nil
		}

		// perform a fixed rate physics tick
		accumulator += dt
		for accumulator >= physicsTickrate {
			space.Step(physicsTickrate)
			accumulator -= physicsTickrate
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode2D(camera)

		// this is a generic way to iterate over the shapes in a space,
		// to avoid the type switch just keep a pointer to the shapes when they've been created
		space.EachShape(func(s *cp.Shape) {
			switch s.Class.(type) {
			case *cp.Segment:
				segment := s.Class.(*cp.Segment)
				a := segment.A()
				b := segment.B()
				rl.DrawLineV(v(a), v(b), rl.Black)
			case *cp.Circle:
				circle := s.Class.(*cp.Circle)
				pos := circle.Body().Position()
				rl.DrawCircleV(v(pos), float32(circle.Radius()), rl.Red)
			default:
				fmt.Println("unexpected shape", s.Class)
			}
		})

		rl.EndMode2D()
		rl.DrawFPS(0, 0)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func v(v cp.Vector) rl.Vector2 {
	return rl.Vector2{X: float32(v.X), Y: float32(v.Y)}
}
