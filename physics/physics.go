// Package physics - 2D Physics library for videogames
//
// A port of Victor Fisac's physac engine (https://github.com/raysan5/raylib/blob/master/src/physac.h)
package physics

import (
	"math"

	"github.com/gen2brain/raylib-go/raylib"
)

// ShapeType type
type ShapeType int

// Physics shape types
const (
	// Circle type
	CircleShape ShapeType = iota
	// Polygon type
	PolygonShape
)

// Polygon type
type Polygon struct {
	// Current used vertex and normals count
	VertexCount int
	// Polygon vertex positions vectors
	Vertices [maxVertices]rl.Vector2
	// Polygon vertex normals vectors
	Normals [maxVertices]rl.Vector2
	// Vertices transform matrix 2x2
	Transform rl.Mat2
}

// Shape type
type Shape struct {
	// Physics shape type (circle or polygon)
	Type ShapeType
	// Shape physics body reference
	Body *Body
	// Circle shape radius (used for circle shapes)
	Radius float32
	// Polygon shape vertices position and normals data (just used for polygon shapes)
	VertexData Polygon
}

// Body type
type Body struct {
	// Enabled dynamics state (collisions are calculated anyway)
	Enabled bool
	// Physics body shape pivot
	Position rl.Vector2
	// Current linear velocity applied to position
	Velocity rl.Vector2
	// Current linear force (reset to 0 every step)
	Force rl.Vector2
	// Current angular velocity applied to orient
	AngularVelocity float32
	// Current angular force (reset to 0 every step)
	Torque float32
	// Rotation in radians
	Orient float32
	// Moment of inertia
	Inertia float32
	// Inverse value of inertia
	InverseInertia float32
	// Physics body mass
	Mass float32
	// Inverse value of mass
	InverseMass float32
	// Friction when the body has not movement (0 to 1)
	StaticFriction float32
	// Friction when the body has movement (0 to 1)
	DynamicFriction float32
	// Restitution coefficient of the body (0 to 1)
	Restitution float32
	// Apply gravity force to dynamics
	UseGravity bool
	// Physics grounded on other body state
	IsGrounded bool
	// Physics rotation constraint
	FreezeOrient bool
	// Physics body shape information (type, radius, vertices, normals)
	Shape Shape
}

// manifold type
type manifold struct {
	// Manifold first physics body reference
	BodyA *Body
	// Manifold second physics body reference
	BodyB *Body
	// Depth of penetration from collision
	Penetration float32
	// Normal direction vector from 'a' to 'b'
	Normal rl.Vector2
	// Points of contact during collision
	Contacts [2]rl.Vector2
	// Current collision number of contacts
	ContactsCount int
	// Mixed restitution during collision
	Restitution float32
	// Mixed dynamic friction during collision
	DynamicFriction float32
	// Mixed static friction during collision
	StaticFriction float32
}

// Constants
const (
	maxBodies    = 64
	maxManifolds = 4096

	maxVertices    = 24
	circleVertices = 24

	collisionIterations   = 100
	penetrationAllowance  = 0.05
	penetrationCorrection = 0.4

	fltMax  = 3.402823466 + 38
	epsilon = 0.000001
)

// Globals
var (
	// Physics bodies pointers
	bodies []*Body

	// Physics manifolds pointers
	manifolds []*manifold

	// Physics world gravity force
	gravityForce rl.Vector2

	// Delta time used for physics steps, in milliseconds
	deltaTime float32
)

// Init - initializes physics values
func Init() {
	deltaTime = 1.0 / 60.0 / 10.0 * 1000
	gravityForce = rl.NewVector2(0, 9.81)

	bodies = make([]*Body, 0, maxBodies)
	manifolds = make([]*manifold, 0, maxManifolds)
}

// Sets physics fixed time step in milliseconds. 1.666666 by default
func SetPhysicsTimeStep(delta float32) {
	deltaTime = delta
}

// SetGravity - Sets physics global gravity force
func SetGravity(x, y float32) {
	gravityForce.X = x
	gravityForce.Y = y
}

// NewBodyCircle - Creates a new circle physics body with generic parameters
func NewBodyCircle(pos rl.Vector2, radius, density float32) *Body {
	return NewBodyPolygon(pos, radius, circleVertices, density)
}

// NewBodyRectangle - Creates a new rectangle physics body with generic parameters
func NewBodyRectangle(pos rl.Vector2, width, height, density float32) *Body {
	newBody := &Body{}

	// Initialize new body with generic values
	newBody.Enabled = true
	newBody.Position = pos
	newBody.Velocity = rl.Vector2{}
	newBody.Force = rl.Vector2{}
	newBody.AngularVelocity = 0
	newBody.Torque = 0
	newBody.Orient = 0

	newBody.Shape = Shape{}
	newBody.Shape.Type = PolygonShape
	newBody.Shape.Body = newBody
	newBody.Shape.VertexData = newRectanglePolygon(pos, rl.NewVector2(width, height))

	// Calculate centroid and moment of inertia
	center := rl.Vector2{}
	area := float32(0.0)
	inertia := float32(0.0)
	k := float32(1.0) / 3.0

	for i := 0; i < newBody.Shape.VertexData.VertexCount; i++ {
		// Triangle vertices, third vertex implied as (0, 0)
		p1 := newBody.Shape.VertexData.Vertices[i]
		nextIndex := 0
		if i+1 < newBody.Shape.VertexData.VertexCount {
			nextIndex = i + 1
		}
		p2 := newBody.Shape.VertexData.Vertices[nextIndex]

		D := rl.Vector2CrossProduct(p1, p2)
		triangleArea := D / 2

		area += triangleArea

		// Use area to weight the centroid average, not just vertex position
		center.X += triangleArea * k * (p1.X + p2.X)
		center.Y += triangleArea * k * (p1.Y + p2.Y)

		intx2 := p1.X*p1.X + p2.X*p1.X + p2.X*p2.X
		inty2 := p1.Y*p1.Y + p2.Y*p1.Y + p2.Y*p2.Y
		inertia += (0.25 * k * D) * (intx2 + inty2)
	}

	center.X *= 1.0 / area
	center.Y *= 1.0 / area

	// Translate vertices to centroid (make the centroid (0, 0) for the polygon in model space)
	// NOTE: this is not really necessary
	for i := 0; i < newBody.Shape.VertexData.VertexCount; i++ {
		newBody.Shape.VertexData.Vertices[i].X -= center.X
		newBody.Shape.VertexData.Vertices[i].Y -= center.Y
	}

	newBody.Mass = density * area
	newBody.Inertia = density * inertia
	newBody.StaticFriction = 0.4
	newBody.DynamicFriction = 0.2
	newBody.Restitution = 0
	newBody.UseGravity = true
	newBody.IsGrounded = false
	newBody.FreezeOrient = false

	if newBody.Mass != 0 {
		newBody.InverseMass = 1.0 / newBody.Mass
	}

	if newBody.Inertia != 0 {
		newBody.InverseInertia = 1.0 / newBody.Inertia
	}

	// Add new body to bodies pointers
	bodies = append(bodies, newBody)

	return newBody
}

// NewBodyPolygon - Creates a new polygon physics body with generic parameters
func NewBodyPolygon(pos rl.Vector2, radius float32, sides int, density float32) *Body {
	newBody := &Body{}

	// Initialize new body with generic values
	newBody.Enabled = true
	newBody.Position = pos
	newBody.Velocity = rl.Vector2{}
	newBody.Force = rl.Vector2{}
	newBody.AngularVelocity = 0
	newBody.Torque = 0
	newBody.Orient = 0

	newBody.Shape = Shape{}
	newBody.Shape.Type = PolygonShape
	newBody.Shape.Body = newBody
	newBody.Shape.VertexData = newRandomPolygon(radius, sides)

	// Calculate centroid and moment of inertia
	center := rl.Vector2{}
	area := float32(0.0)
	inertia := float32(0.0)
	alpha := float32(1.0) / 3.0

	for i := 0; i < newBody.Shape.VertexData.VertexCount; i++ {
		// Triangle vertices, third vertex implied as (0, 0)
		position1 := newBody.Shape.VertexData.Vertices[i]
		nextIndex := 0
		if i+1 < newBody.Shape.VertexData.VertexCount {
			nextIndex = i + 1
		}
		position2 := newBody.Shape.VertexData.Vertices[nextIndex]

		cross := rl.Vector2CrossProduct(position1, position2)
		triangleArea := cross / 2

		area += triangleArea

		// Use area to weight the centroid average, not just vertex position
		center.X += triangleArea * alpha * (position1.X + position2.X)
		center.Y += triangleArea * alpha * (position1.Y + position2.Y)

		intx2 := position1.X*position1.X + position2.X*position1.X + position2.X*position2.X
		inty2 := position1.Y*position1.Y + position2.Y*position1.Y + position2.Y*position2.Y
		inertia += (0.25 * alpha * cross) * (intx2 + inty2)
	}

	center.X *= 1.0 / area
	center.Y *= 1.0 / area

	// Translate vertices to centroid (make the centroid (0, 0) for the polygon in model space)
	// Note: this is not really necessary
	for i := 0; i < newBody.Shape.VertexData.VertexCount; i++ {
		newBody.Shape.VertexData.Vertices[i].X -= center.X
		newBody.Shape.VertexData.Vertices[i].Y -= center.Y
	}

	newBody.Mass = density * area
	newBody.Inertia = density * inertia
	newBody.StaticFriction = 0.4
	newBody.DynamicFriction = 0.2
	newBody.Restitution = 0
	newBody.UseGravity = true
	newBody.IsGrounded = false
	newBody.FreezeOrient = false

	if newBody.Mass != 0 {
		newBody.InverseMass = 1.0 / newBody.Mass
	}

	if newBody.Inertia != 0 {
		newBody.InverseInertia = 1.0 / newBody.Inertia
	}

	// Add new body to bodies pointers
	bodies = append(bodies, newBody)

	return newBody
}

// GetBodies - Returns the slice of created physics bodies
func GetBodies() []*Body {
	return bodies
}

// GetBodiesCount - Returns the current amount of created physics bodies
func GetBodiesCount() int {
	return len(bodies)
}

// GetBody - Returns a physics body of the bodies pool at a specific index
func GetBody(index int) *Body {
	var body *Body

	if index < len(bodies) {
		body = bodies[index]
	} else {
		rl.TraceLog(rl.LogDebug, "[PHYSAC] physics body index is out of bounds")
	}

	return body
}

// GetShapeType - Returns the physics body shape type (Circle or Polygon)
func GetShapeType(index int) ShapeType {
	var result ShapeType

	if index < len(bodies) {
		result = bodies[index].Shape.Type
	} else {
		rl.TraceLog(rl.LogDebug, "[PHYSAC] physics body index is out of bounds")
	}

	return result
}

// GetShapeVerticesCount - Returns the amount of vertices of a physics body shape
func GetShapeVerticesCount(index int) int {
	result := 0

	if index < len(bodies) {
		switch bodies[index].Shape.Type {
		case CircleShape:
			result = circleVertices
			break
		case PolygonShape:
			result = bodies[index].Shape.VertexData.VertexCount
			break
		}
	} else {
		rl.TraceLog(rl.LogDebug, "[PHYSAC] physics body index is out of bounds")
	}

	return result
}

// DestroyBody - Unitializes and destroys a physics body
func DestroyBody(body *Body) bool {
	for index, b := range bodies {
		if b == body {
			// Free body allocated memory
			bodies = append(bodies[:index], bodies[index+1:]...)
			return true
		}
	}

	return false
}

// Update - Physics steps calculations (dynamics, collisions and position corrections)
func Update() {
	deltaTime = rl.GetFrameTime() * 1000

	// Clear previous generated collisions information
	for _, m := range manifolds {
		destroyManifold(m)
	}

	// Reset physics bodies grounded state
	for _, b := range bodies {
		b.IsGrounded = false
	}

	// Generate new collision information
	bodiesCount := len(bodies)
	for i := 0; i < bodiesCount; i++ {

		bodyA := bodies[i]

		for j := i + 1; j < bodiesCount; j++ {

			bodyB := bodies[j]

			if (bodyA.InverseMass == 0) && (bodyB.InverseMass == 0) {
				continue
			}

			var m *manifold
			if bodyA.Shape.Type == PolygonShape && bodyB.Shape.Type == CircleShape {
				m = newManifold(bodyB, bodyA)
			} else {
				m = newManifold(bodyA, bodyB)
			}

			m.solveManifold()

			if m.ContactsCount > 0 {
				// Create a new manifold with same information as previously solved manifold and add it to the manifolds pool last slot
				newManifold := newManifold(bodyA, bodyB)
				newManifold.Penetration = m.Penetration
				newManifold.Normal = m.Normal
				newManifold.Contacts[0] = m.Contacts[0]
				newManifold.Contacts[1] = m.Contacts[1]
				newManifold.ContactsCount = m.ContactsCount
				newManifold.Restitution = m.Restitution
				newManifold.DynamicFriction = m.DynamicFriction
				newManifold.StaticFriction = m.StaticFriction
			}
		}
	}

	// Integrate forces to physics bodies
	for _, b := range bodies {
		b.integrateForces()
	}

	// Initialize physics manifolds to solve collisions
	for _, m := range manifolds {
		m.initializeManifolds()
	}

	// Integrate physics collisions impulses to solve collisions
	manifoldsCount := len(manifolds)
	for i := 0; i < collisionIterations; i++ {
		for j := 0; j < manifoldsCount; j++ {
			if i < manifoldsCount {
				manifolds[i].integrateImpulses()
			}
		}
	}

	// Integrate velocity to physics bodies
	for _, b := range bodies {
		b.integrateVelocity()
	}

	// Correct physics bodies positions based on manifolds collision information
	for _, m := range manifolds {
		m.correctPositions()
	}

	// Clear physics bodies forces
	for _, b := range bodies {
		b.Force = rl.Vector2{}
		b.Torque = 0
	}
}

// Reset - Destroys created physics bodies and manifolds
func Reset() {
	bodies = make([]*Body, 0, maxBodies)
	manifolds = make([]*manifold, 0, maxManifolds)
}

// Close - Unitializes physics pointers
func Close() {
	// Unitialize physics manifolds dynamic memory allocations
	for _, m := range manifolds {
		destroyManifold(m)
	}

	// Unitialize physics bodies dynamic memory allocations
	for _, b := range bodies {
		DestroyBody(b)
	}
}

// AddForce - Adds a force to a physics body
func (b *Body) AddForce(force rl.Vector2) {
	b.Force = rl.Vector2Add(b.Force, force)
}

// AddTorque - Adds an angular force to a physics body
func (b *Body) AddTorque(amount float32) {
	b.Torque += amount
}

// Shatter - Shatters a polygon shape physics body to little physics bodies with explosion force
func (b *Body) Shatter(position rl.Vector2, force float32) {
	if b.Shape.Type != PolygonShape {
		return
	}

	vertexData := b.Shape.VertexData
	collision := false

	for i := 0; i < vertexData.VertexCount; i++ {
		positionA := b.Position
		positionB := rl.Mat2MultiplyVector2(vertexData.Transform, rl.Vector2Add(b.Position, vertexData.Vertices[i]))

		nextIndex := 0
		if i+1 < vertexData.VertexCount {
			nextIndex = i + 1
		}

		positionC := rl.Mat2MultiplyVector2(vertexData.Transform, rl.Vector2Add(b.Position, vertexData.Vertices[nextIndex]))

		// Check collision between each triangle
		alpha := ((positionB.Y-positionC.Y)*(position.X-positionC.X) + (positionC.X-positionB.X)*(position.Y-positionC.Y)) /
			((positionB.Y-positionC.Y)*(positionA.X-positionC.X) + (positionC.X-positionB.X)*(positionA.Y-positionC.Y))

		beta := ((positionC.Y-positionA.Y)*(position.X-positionC.X) + (positionA.X-positionC.X)*(position.Y-positionC.Y)) /
			((positionB.Y-positionC.Y)*(positionA.X-positionC.X) + (positionC.X-positionB.X)*(positionA.Y-positionC.Y))

		gamma := 1.0 - alpha - beta

		if alpha > 0 && beta > 0 && gamma > 0 {
			collision = true
			break
		}
	}

	if collision {
		count := vertexData.VertexCount
		bodyPos := b.Position

		vertices := make([]rl.Vector2, count)
		trans := vertexData.Transform
		for i := 0; i < count; i++ {
			vertices[i] = vertexData.Vertices[i]
		}

		// Destroy shattered physics body
		DestroyBody(b)

		for i := 0; i < count; i++ {
			nextIndex := 0
			if i+1 < count {
				nextIndex = i + 1
			}

			center := triangleBarycenter(vertices[i], vertices[nextIndex], rl.NewVector2(0, 0))
			center = rl.Vector2Add(bodyPos, center)
			offset := rl.Vector2Subtract(center, bodyPos)

			newBody := NewBodyPolygon(center, 10, 3, 10) // Create polygon physics body with relevant values

			newData := Polygon{}
			newData.VertexCount = 3
			newData.Transform = trans

			newData.Vertices[0] = rl.Vector2Subtract(vertices[i], offset)
			newData.Vertices[1] = rl.Vector2Subtract(vertices[nextIndex], offset)
			newData.Vertices[2] = rl.Vector2Subtract(position, center)

			// Separate vertices to avoid unnecessary physics collisions
			newData.Vertices[0].X *= 0.95
			newData.Vertices[0].Y *= 0.95
			newData.Vertices[1].X *= 0.95
			newData.Vertices[1].Y *= 0.95
			newData.Vertices[2].X *= 0.95
			newData.Vertices[2].Y *= 0.95

			// Calculate polygon faces normals
			for j := 0; j < newData.VertexCount; j++ {
				nextVertex := 0
				if j+1 < newData.VertexCount {
					nextVertex = j + 1
				}

				face := rl.Vector2Subtract(newData.Vertices[nextVertex], newData.Vertices[j])

				newData.Normals[j] = rl.NewVector2(face.Y, -face.X)
				normalize(&newData.Normals[j])
			}

			// Apply computed vertex data to new physics body shape
			newBody.Shape.VertexData = newData

			// Calculate centroid and moment of inertia
			center = rl.NewVector2(0, 0)
			area := float32(0.0)
			inertia := float32(0.0)
			k := float32(1.0) / 3.0

			for j := 0; j < newBody.Shape.VertexData.VertexCount; j++ {
				// Triangle vertices, third vertex implied as (0, 0)
				p1 := newBody.Shape.VertexData.Vertices[j]
				nextVertex := 0
				if j+1 < newBody.Shape.VertexData.VertexCount {
					nextVertex = j + 1
				}
				p2 := newBody.Shape.VertexData.Vertices[nextVertex]

				D := rl.Vector2CrossProduct(p1, p2)
				triangleArea := D / 2

				area += triangleArea

				// Use area to weight the centroid average, not just vertex position
				center.X += triangleArea * k * (p1.X + p2.X)
				center.Y += triangleArea * k * (p1.Y + p2.Y)

				intx2 := p1.X*p1.X + p2.X*p1.X + p2.X*p2.X
				inty2 := p1.Y*p1.Y + p2.Y*p1.Y + p2.Y*p2.Y
				inertia += (0.25 * k * D) * (intx2 + inty2)
			}

			center.X *= 1.0 / area
			center.Y *= 1.0 / area

			newBody.Mass = area
			newBody.Inertia = inertia

			if newBody.Mass != 0 {
				newBody.InverseMass = 1.0 / newBody.Mass
			}

			if newBody.Inertia != 0 {
				newBody.InverseInertia = 1.0 / newBody.Inertia
			}

			// Calculate explosion force direction
			pointA := newBody.Position
			pointB := rl.Vector2Subtract(newData.Vertices[1], newData.Vertices[0])
			pointB.X /= 2
			pointB.Y /= 2
			forceDirection := rl.Vector2Subtract(rl.Vector2Add(pointA, rl.Vector2Add(newData.Vertices[0], pointB)), newBody.Position)
			normalize(&forceDirection)
			forceDirection.X *= force
			forceDirection.Y *= force

			// Apply force to new physics body
			newBody.AddForce(forceDirection)
		}
	}
}

// GetShapeVertex - Returns transformed position of a body shape (body position + vertex transformed position)
func (b *Body) GetShapeVertex(vertex int) rl.Vector2 {
	position := rl.Vector2{}

	switch b.Shape.Type {
	case CircleShape:
		position.X = b.Position.X + float32(math.Cos(360/float64(circleVertices)*float64(vertex)*rl.Deg2rad))*b.Shape.Radius
		position.Y = b.Position.Y + float32(math.Sin(360/float64(circleVertices)*float64(vertex)*rl.Deg2rad))*b.Shape.Radius
		break
	case PolygonShape:
		position = rl.Vector2Add(b.Position, rl.Mat2MultiplyVector2(b.Shape.VertexData.Transform, b.Shape.VertexData.Vertices[vertex]))
		break
	}

	return position
}

// SetRotation - Sets physics body shape transform based on radians parameter
func (b *Body) SetRotation(radians float32) {
	b.Orient = radians

	if b.Shape.Type == PolygonShape {
		b.Shape.VertexData.Transform = rl.Mat2Radians(radians)
	}
}

// integrateVelocity - Integrates physics velocity into position and forces
func (b *Body) integrateVelocity() {
	if !b.Enabled {
		return
	}

	b.Position.X += b.Velocity.X * deltaTime
	b.Position.Y += b.Velocity.Y * deltaTime

	if !b.FreezeOrient {
		b.Orient += b.AngularVelocity * deltaTime
	}

	rl.Mat2Set(&b.Shape.VertexData.Transform, b.Orient)

	b.integrateForces()
}

// integrateForces -  Integrates physics forces into velocity
func (b *Body) integrateForces() {
	if b.InverseMass == 0 || !b.Enabled {
		return
	}

	b.Velocity.X += (b.Force.X * b.InverseMass) * (deltaTime / 2)
	b.Velocity.Y += (b.Force.Y * b.InverseMass) * (deltaTime / 2)

	if b.UseGravity {
		b.Velocity.X += gravityForce.X * (deltaTime / 1000 / 2)
		b.Velocity.Y += gravityForce.Y * (deltaTime / 1000 / 2)
	}

	if !b.FreezeOrient {
		b.AngularVelocity += b.Torque * b.InverseInertia * (deltaTime / 2)
	}
}

// newRandomPolygon - Creates a random polygon shape with max vertex distance from polygon pivot
func newRandomPolygon(radius float32, sides int) Polygon {
	data := Polygon{}
	data.VertexCount = sides

	orient := rl.GetRandomValue(0, 360)
	data.Transform = rl.Mat2Radians(float32(orient) * rl.Deg2rad)

	// Calculate polygon vertices positions
	for i := 0; i < data.VertexCount; i++ {
		data.Vertices[i].X = float32(math.Cos(360/float64(sides)*float64(i)*rl.Deg2rad)) * radius
		data.Vertices[i].Y = float32(math.Sin(360/float64(sides)*float64(i)*rl.Deg2rad)) * radius
	}

	// Calculate polygon faces normals
	for i := 0; i < data.VertexCount; i++ {
		nextIndex := 0
		if i+1 < sides {
			nextIndex = i + 1
		}

		face := rl.Vector2Subtract(data.Vertices[nextIndex], data.Vertices[i])

		data.Normals[i] = rl.NewVector2(face.Y, -face.X)
		normalize(&data.Normals[i])
	}

	return data
}

// newRectanglePolygon - Creates a rectangle polygon shape based on a min and max positions
func newRectanglePolygon(pos, size rl.Vector2) Polygon {
	data := Polygon{}

	data.VertexCount = 4
	data.Transform = rl.Mat2Radians(0)

	// Calculate polygon vertices positions
	data.Vertices[0] = rl.NewVector2(pos.X+size.X/2, pos.Y-size.Y/2)
	data.Vertices[1] = rl.NewVector2(pos.X+size.X/2, pos.Y+size.Y/2)
	data.Vertices[2] = rl.NewVector2(pos.X-size.X/2, pos.Y+size.Y/2)
	data.Vertices[3] = rl.NewVector2(pos.X-size.X/2, pos.Y-size.Y/2)

	// Calculate polygon faces normals
	for i := 0; i < data.VertexCount; i++ {
		nextIndex := 0
		if i+1 < data.VertexCount {
			nextIndex = i + 1
		}
		face := rl.Vector2Subtract(data.Vertices[nextIndex], data.Vertices[i])

		data.Normals[i] = rl.NewVector2(face.Y, -face.X)
		normalize(&data.Normals[i])
	}

	return data
}

// newManifold - Creates a new physics manifold to solve collision
func newManifold(a, b *Body) *manifold {
	newManifold := &manifold{}

	// Initialize new manifold with generic values
	newManifold.BodyA = a
	newManifold.BodyB = b
	newManifold.Penetration = 0
	newManifold.Normal = rl.Vector2{}
	newManifold.Contacts[0] = rl.Vector2{}
	newManifold.Contacts[1] = rl.Vector2{}
	newManifold.ContactsCount = 0
	newManifold.Restitution = 0
	newManifold.DynamicFriction = 0
	newManifold.StaticFriction = 0

	// Add new manifold to manifolds pointers
	manifolds = append(manifolds, newManifold)

	return newManifold
}

// destroyManifold - Unitializes and destroys a physics manifold
func destroyManifold(manifold *manifold) bool {
	for index, m := range manifolds {
		if m == manifold {
			// Free manifold allocated memory
			manifolds = append(manifolds[:index], manifolds[index+1:]...)
			return true
		}
	}

	return false
}

// solveManifold - Solves a created physics manifold between two physics bodies
func (m *manifold) solveManifold() {
	switch m.BodyA.Shape.Type {
	case CircleShape:
		switch m.BodyB.Shape.Type {
		case CircleShape:
			m.solveCircleToCircle()
			break
		case PolygonShape:
			m.solveCircleToPolygon()
			break
		}
	case PolygonShape:
		switch m.BodyB.Shape.Type {
		case CircleShape:
			m.solvePolygonToCircle()
			break
		case PolygonShape:
			m.solvePolygonToPolygon()
			break
		}
	}

	// Update physics body grounded state if normal direction is down and grounded state is not set yet in previous manifolds
	if !m.BodyB.IsGrounded {
		m.BodyB.IsGrounded = (m.Normal.Y < 0)
	}
}

// solveCircleToCircle - Solves collision between two circle shape physics bodies
func (m *manifold) solveCircleToCircle() {
	bodyA := m.BodyA
	bodyB := m.BodyB

	// Calculate translational vector, which is normal
	normal := rl.Vector2Subtract(bodyB.Position, bodyA.Position)

	distSqr := rl.Vector2LenSqr(normal)
	radius := bodyA.Shape.Radius + bodyB.Shape.Radius

	// Check if circles are not in contact
	if distSqr >= radius*radius {
		m.ContactsCount = 0
		return
	}

	distance := float32(math.Sqrt(float64(distSqr)))
	m.ContactsCount = 1

	if distance == 0 {
		m.Penetration = bodyA.Shape.Radius
		m.Normal = rl.NewVector2(1, 0)
		m.Contacts[0] = bodyA.Position
	} else {
		m.Penetration = radius - distance
		m.Normal = rl.NewVector2(normal.X/distance, normal.Y/distance) // Faster than using normalize() due to sqrt is already performed
		m.Contacts[0] = rl.NewVector2(m.Normal.X*bodyA.Shape.Radius+bodyA.Position.X, m.Normal.Y*bodyA.Shape.Radius+bodyA.Position.Y)
	}

	// Update physics body grounded state if normal direction is down
	if !bodyA.IsGrounded {
		bodyA.IsGrounded = (m.Normal.Y < 0)
	}
}

// solveCircleToPolygon - Solves collision between a circle to a polygon shape physics bodies
func (m *manifold) solveCircleToPolygon() {
	m.ContactsCount = 0

	// Transform circle center to polygon transform space
	center := m.BodyA.Position
	center = rl.Mat2MultiplyVector2(rl.Mat2Transpose(m.BodyB.Shape.VertexData.Transform), rl.Vector2Subtract(center, m.BodyB.Position))

	// Find edge with minimum penetration
	// It is the same concept as using support points in solvePolygonToPolygon
	separation := float32(-fltMax)
	faceNormal := 0
	vertexData := m.BodyB.Shape.VertexData

	for i := 0; i < vertexData.VertexCount; i++ {
		currentSeparation := rl.Vector2DotProduct(vertexData.Normals[i], rl.Vector2Subtract(center, vertexData.Vertices[i]))

		if currentSeparation > m.BodyA.Shape.Radius {
			return
		}

		if currentSeparation > separation {
			separation = currentSeparation
			faceNormal = i
		}
	}

	// Grab face's vertices
	v1 := vertexData.Vertices[faceNormal]
	nextIndex := 0
	if faceNormal+1 < vertexData.VertexCount {
		nextIndex = faceNormal + 1
	}
	v2 := vertexData.Vertices[nextIndex]

	// Check to see if center is within polygon
	if separation < epsilon {
		m.ContactsCount = 1
		normal := rl.Mat2MultiplyVector2(vertexData.Transform, vertexData.Normals[faceNormal])
		m.Normal = rl.NewVector2(-normal.X, -normal.Y)
		m.Contacts[0] = rl.NewVector2(m.Normal.X*m.BodyA.Shape.Radius+m.BodyA.Position.X, m.Normal.Y*m.BodyA.Shape.Radius+m.BodyA.Position.Y)
		m.Penetration = m.BodyA.Shape.Radius
		return
	}

	// Determine which voronoi region of the edge center of circle lies within
	dot1 := rl.Vector2DotProduct(rl.Vector2Subtract(center, v1), rl.Vector2Subtract(v2, v1))
	dot2 := rl.Vector2DotProduct(rl.Vector2Subtract(center, v2), rl.Vector2Subtract(v1, v2))
	m.Penetration = m.BodyA.Shape.Radius - separation

	if dot1 <= 0 { // Closest to v1
		if rl.Vector2Distance(center, v1) > m.BodyA.Shape.Radius*m.BodyA.Shape.Radius {
			return
		}

		m.ContactsCount = 1
		normal := rl.Vector2Subtract(v1, center)
		normal = rl.Mat2MultiplyVector2(vertexData.Transform, normal)
		normalize(&normal)
		m.Normal = normal
		v1 = rl.Mat2MultiplyVector2(vertexData.Transform, v1)
		v1 = rl.Vector2Add(v1, m.BodyB.Position)
		m.Contacts[0] = v1
	} else if dot2 <= 0 { // Closest to v2
		if rl.Vector2Distance(center, v2) > m.BodyA.Shape.Radius*m.BodyA.Shape.Radius {
			return
		}

		m.ContactsCount = 1
		normal := rl.Vector2Subtract(v2, center)
		v2 = rl.Mat2MultiplyVector2(vertexData.Transform, v2)
		v2 = rl.Vector2Add(v2, m.BodyB.Position)
		m.Contacts[0] = v2
		normal = rl.Mat2MultiplyVector2(vertexData.Transform, normal)
		normalize(&normal)
		m.Normal = normal
	} else { // Closest to face
		normal := vertexData.Normals[faceNormal]

		if rl.Vector2DotProduct(rl.Vector2Subtract(center, v1), normal) > m.BodyA.Shape.Radius {
			return
		}

		normal = rl.Mat2MultiplyVector2(vertexData.Transform, normal)
		m.Normal = rl.NewVector2(-normal.X, -normal.Y)
		m.Contacts[0] = rl.NewVector2(m.Normal.X*m.BodyA.Shape.Radius+m.BodyA.Position.X, m.Normal.Y*m.BodyA.Shape.Radius+m.BodyA.Position.Y)
		m.ContactsCount = 1
	}
}

// solvePolygonToCircle - Solves collision between a polygon to a circle shape physics bodies
func (m *manifold) solvePolygonToCircle() {
	bodyA := m.BodyA
	bodyB := m.BodyB

	m.BodyA = bodyB
	m.BodyB = bodyA

	m.solveCircleToPolygon()

	m.Normal.X *= -1
	m.Normal.Y *= -1
}

// solvePolygonToPolygon - Solves collision between two polygons shape physics bodies
func (m *manifold) solvePolygonToPolygon() {
	bodyA := m.BodyA.Shape
	bodyB := m.BodyB.Shape
	m.ContactsCount = 0

	// Check for separating axis with A shape's face planes
	faceA, penetrationA := findAxisLeastPenetration(bodyA, bodyB)
	if penetrationA >= 0 {
		return
	}

	// Check for separating axis with B shape's face planes
	faceB, penetrationB := findAxisLeastPenetration(bodyB, bodyA)
	if penetrationB >= 0 {
		return
	}

	referenceIndex := 0
	flip := false // Always point from A shape to B shape

	refPoly := Shape{} // Reference
	incPoly := Shape{} // Incident

	// Determine which shape contains reference face
	if biasGreaterThan(penetrationA, penetrationB) {
		refPoly = bodyA
		incPoly = bodyB
		referenceIndex = faceA
	} else {
		refPoly = bodyB
		incPoly = bodyA
		referenceIndex = faceB
		flip = true
	}

	// World space incident face
	incidentFace0 := rl.Vector2{}
	incidentFace1 := rl.Vector2{}
	findIncidentFace(&incidentFace0, &incidentFace1, refPoly, incPoly, referenceIndex)

	// Setup reference face vertices
	refData := refPoly.VertexData
	v1 := refData.Vertices[referenceIndex]
	if referenceIndex+1 < refData.VertexCount {
		referenceIndex = referenceIndex + 1
	} else {
		referenceIndex = 0
	}
	v2 := refData.Vertices[referenceIndex]

	// Transform vertices to world space
	v1 = rl.Mat2MultiplyVector2(refData.Transform, v1)
	v1 = rl.Vector2Add(v1, refPoly.Body.Position)
	v2 = rl.Mat2MultiplyVector2(refData.Transform, v2)
	v2 = rl.Vector2Add(v2, refPoly.Body.Position)

	// Calculate reference face side normal in world space
	sidePlaneNormal := rl.Vector2Subtract(v2, v1)
	normalize(&sidePlaneNormal)

	// Orthogonalize
	refFaceNormal := rl.NewVector2(sidePlaneNormal.Y, -sidePlaneNormal.X)
	refC := rl.Vector2DotProduct(refFaceNormal, v1)
	negSide := rl.Vector2DotProduct(sidePlaneNormal, v1) * -1
	posSide := rl.Vector2DotProduct(sidePlaneNormal, v2)

	// clip incident face to reference face side planes (due to floating point error, possible to not have required points
	if clip(rl.NewVector2(-sidePlaneNormal.X, -sidePlaneNormal.Y), negSide, &incidentFace0, &incidentFace1) < 2 {
		return
	}
	if clip(sidePlaneNormal, posSide, &incidentFace0, &incidentFace1) < 2 {
		return
	}

	// Flip normal if required
	if flip {
		m.Normal = rl.NewVector2(-refFaceNormal.X, -refFaceNormal.Y)
	} else {
		m.Normal = refFaceNormal
	}

	// Keep points behind reference face
	currentPoint := 0 // clipped points behind reference face
	separation := rl.Vector2DotProduct(refFaceNormal, incidentFace0) - refC
	if separation <= 0 {
		m.Contacts[currentPoint] = incidentFace0
		m.Penetration = -separation
		currentPoint++
	} else {
		m.Penetration = 0
	}

	separation = rl.Vector2DotProduct(refFaceNormal, incidentFace1) - refC

	if separation <= 0 {
		m.Contacts[currentPoint] = incidentFace1
		m.Penetration += -separation
		currentPoint++

		// Calculate total penetration average
		m.Penetration /= float32(currentPoint)
	}

	m.ContactsCount = currentPoint
}

// initializeManifolds - Initializes physics manifolds to solve collisions
func (m *manifold) initializeManifolds() {
	bodyA := m.BodyA
	bodyB := m.BodyB

	// Calculate average restitution, static and dynamic friction
	m.Restitution = float32(math.Sqrt(float64(bodyA.Restitution * bodyB.Restitution)))
	m.StaticFriction = float32(math.Sqrt(float64(bodyA.StaticFriction * bodyB.StaticFriction)))
	m.DynamicFriction = float32(math.Sqrt(float64(bodyA.DynamicFriction * bodyB.DynamicFriction)))

	for i := 0; i < 2; i++ {
		// Caculate radius from center of mass to contact
		radiusA := rl.Vector2Subtract(m.Contacts[i], bodyA.Position)
		radiusB := rl.Vector2Subtract(m.Contacts[i], bodyB.Position)

		crossA := rl.Vector2Cross(bodyA.AngularVelocity, radiusA)
		crossB := rl.Vector2Cross(bodyB.AngularVelocity, radiusB)

		radiusV := rl.Vector2{}
		radiusV.X = bodyB.Velocity.X + crossB.X - bodyA.Velocity.X - crossA.X
		radiusV.Y = bodyB.Velocity.Y + crossB.Y - bodyA.Velocity.Y - crossA.Y

		// Determine if we should perform a resting collision or not;
		// The idea is if the only thing moving this object is gravity, then the collision should be performed without any restitution
		if rl.Vector2LenSqr(radiusV) < (rl.Vector2LenSqr(rl.NewVector2(gravityForce.X*deltaTime/1000, gravityForce.Y*deltaTime/1000)) + epsilon) {
			m.Restitution = 0
		}
	}
}

// integrateImpulses - Integrates physics collisions impulses to solve collisions
func (m *manifold) integrateImpulses() {
	bodyA := m.BodyA
	bodyB := m.BodyB

	// Early out and positional correct if both objects have infinite mass
	if math.Abs(float64(bodyA.InverseMass+bodyB.InverseMass)) <= epsilon {
		bodyA.Velocity = rl.Vector2{}
		bodyB.Velocity = rl.Vector2{}
		return
	}

	for i := 0; i < m.ContactsCount; i++ {
		// Calculate radius from center of mass to contact
		radiusA := rl.Vector2Subtract(m.Contacts[i], bodyA.Position)
		radiusB := rl.Vector2Subtract(m.Contacts[i], bodyB.Position)

		// Calculate relative velocity
		radiusV := rl.Vector2{}
		radiusV.X = bodyB.Velocity.X + rl.Vector2Cross(bodyB.AngularVelocity, radiusB).X - bodyA.Velocity.X - rl.Vector2Cross(bodyA.AngularVelocity, radiusA).X
		radiusV.Y = bodyB.Velocity.Y + rl.Vector2Cross(bodyB.AngularVelocity, radiusB).Y - bodyA.Velocity.Y - rl.Vector2Cross(bodyA.AngularVelocity, radiusA).Y

		// Relative velocity along the normal
		contactVelocity := rl.Vector2DotProduct(radiusV, m.Normal)

		// Do not resolve if velocities are separating
		if contactVelocity > 0 {
			return
		}

		raCrossN := rl.Vector2CrossProduct(radiusA, m.Normal)
		rbCrossN := rl.Vector2CrossProduct(radiusB, m.Normal)

		inverseMassSum := bodyA.InverseMass + bodyB.InverseMass + (raCrossN*raCrossN)*bodyA.InverseInertia + (rbCrossN*rbCrossN)*bodyB.InverseInertia

		// Calculate impulse scalar value
		impulse := -(1.0 + m.Restitution) * contactVelocity
		impulse /= inverseMassSum
		impulse /= float32(m.ContactsCount)

		// Apply impulse to each physics body
		impulseV := rl.NewVector2(m.Normal.X*impulse, m.Normal.Y*impulse)

		if bodyA.Enabled {
			bodyA.Velocity.X += bodyA.InverseMass * (-impulseV.X)
			bodyA.Velocity.Y += bodyA.InverseMass * (-impulseV.Y)
			if !bodyA.FreezeOrient {
				bodyA.AngularVelocity += bodyA.InverseInertia * rl.Vector2CrossProduct(radiusA, rl.NewVector2(-impulseV.X, -impulseV.Y))
			}
		}

		if bodyB.Enabled {
			bodyB.Velocity.X += bodyB.InverseMass * (impulseV.X)
			bodyB.Velocity.Y += bodyB.InverseMass * (impulseV.Y)
			if !bodyB.FreezeOrient {
				bodyB.AngularVelocity += bodyB.InverseInertia * rl.Vector2CrossProduct(radiusB, impulseV)
			}
		}

		// Apply friction impulse to each physics body
		radiusV.X = bodyB.Velocity.X + rl.Vector2Cross(bodyB.AngularVelocity, radiusB).X - bodyA.Velocity.X - rl.Vector2Cross(bodyA.AngularVelocity, radiusA).X
		radiusV.Y = bodyB.Velocity.Y + rl.Vector2Cross(bodyB.AngularVelocity, radiusB).Y - bodyA.Velocity.Y - rl.Vector2Cross(bodyA.AngularVelocity, radiusA).Y

		tangent := rl.NewVector2(radiusV.X-(m.Normal.X*rl.Vector2DotProduct(radiusV, m.Normal)), radiusV.Y-(m.Normal.Y*rl.Vector2DotProduct(radiusV, m.Normal)))
		normalize(&tangent)

		// Calculate impulse tangent magnitude
		impulseTangent := -(rl.Vector2DotProduct(radiusV, tangent))
		impulseTangent /= inverseMassSum
		impulseTangent /= float32(m.ContactsCount)

		absImpulseTangent := float32(math.Abs(float64(impulseTangent)))

		// Don't apply tiny friction impulses
		if absImpulseTangent <= epsilon {
			return
		}

		// Apply coulumb's law
		tangentImpulse := rl.Vector2{}
		if absImpulseTangent < impulse*m.StaticFriction {
			tangentImpulse = rl.NewVector2(tangent.X*impulseTangent, tangent.Y*impulseTangent)
		} else {
			tangentImpulse = rl.NewVector2(tangent.X*-impulse*m.DynamicFriction, tangent.Y*-impulse*m.DynamicFriction)
		}

		// Apply friction impulse
		if bodyA.Enabled {
			bodyA.Velocity.X += bodyA.InverseMass * (-tangentImpulse.X)
			bodyA.Velocity.Y += bodyA.InverseMass * (-tangentImpulse.Y)

			if !bodyA.FreezeOrient {
				bodyA.AngularVelocity += bodyA.InverseInertia * rl.Vector2CrossProduct(radiusA, rl.NewVector2(-tangentImpulse.X, -tangentImpulse.Y))
			}
		}

		if bodyB.Enabled {
			bodyB.Velocity.X += bodyB.InverseMass * (tangentImpulse.X)
			bodyB.Velocity.Y += bodyB.InverseMass * (tangentImpulse.Y)

			if !bodyB.FreezeOrient {
				bodyB.AngularVelocity += bodyB.InverseInertia * rl.Vector2CrossProduct(radiusB, tangentImpulse)
			}
		}
	}
}

// correctPositions - Corrects physics bodies positions based on manifolds collision information
func (m *manifold) correctPositions() {
	bodyA := m.BodyA
	bodyB := m.BodyB

	correction := rl.Vector2{}
	correction.X = float32(math.Max(float64(m.Penetration-penetrationAllowance), 0)) / (bodyA.InverseMass + bodyB.InverseMass) * m.Normal.X * penetrationCorrection
	correction.Y = float32(math.Max(float64(m.Penetration-penetrationAllowance), 0)) / (bodyA.InverseMass + bodyB.InverseMass) * m.Normal.Y * penetrationCorrection

	if bodyA.Enabled {
		bodyA.Position.X -= correction.X * bodyA.InverseMass
		bodyA.Position.Y -= correction.Y * bodyA.InverseMass
	}

	if bodyB.Enabled {
		bodyB.Position.X += correction.X * bodyB.InverseMass
		bodyB.Position.Y += correction.Y * bodyB.InverseMass
	}
}

// getSupport - Returns the extreme point along a direction within a polygon
func getSupport(shape Shape, dir rl.Vector2) rl.Vector2 {
	bestProjection := float32(-fltMax)
	bestVertex := rl.Vector2{}

	for i := 0; i < shape.VertexData.VertexCount; i++ {
		vertex := shape.VertexData.Vertices[i]
		projection := rl.Vector2DotProduct(vertex, dir)

		if projection > bestProjection {
			bestVertex = vertex
			bestProjection = projection
		}
	}

	return bestVertex
}

// findAxisLeastPenetration - Finds polygon shapes axis least penetration
func findAxisLeastPenetration(shapeA, shapeB Shape) (int, float32) {
	bestIndex := 0
	bestDistance := float32(-fltMax)

	dataA := shapeA.VertexData
	dataB := shapeB.VertexData

	for i := 0; i < dataA.VertexCount; i++ {
		// Retrieve a face normal from A shape
		normal := dataA.Normals[i]
		transNormal := rl.Mat2MultiplyVector2(dataA.Transform, normal)

		// Transform face normal into B shape's model space
		buT := rl.Mat2Transpose(dataB.Transform)
		normal = rl.Mat2MultiplyVector2(buT, transNormal)

		// Retrieve support point from B shape along -n
		support := getSupport(shapeB, rl.NewVector2(-normal.X, -normal.Y))

		// Retrieve vertex on face from A shape, transform into B shape's model space
		vertex := dataA.Vertices[i]
		vertex = rl.Mat2MultiplyVector2(dataA.Transform, vertex)
		vertex = rl.Vector2Add(vertex, shapeA.Body.Position)
		vertex = rl.Vector2Subtract(vertex, shapeB.Body.Position)
		vertex = rl.Mat2MultiplyVector2(buT, vertex)

		// Compute penetration distance in B shape's model space
		distance := rl.Vector2DotProduct(normal, rl.Vector2Subtract(support, vertex))

		// Store greatest distance
		if distance > bestDistance {
			bestDistance = distance
			bestIndex = i
		}
	}

	return bestIndex, bestDistance
}

// findIncidentFace - Finds two polygon shapes incident face
func findIncidentFace(v0, v1 *rl.Vector2, ref, inc Shape, index int) {
	refData := ref.VertexData
	incData := inc.VertexData

	referenceNormal := refData.Normals[index]

	// Calculate normal in incident's frame of reference
	referenceNormal = rl.Mat2MultiplyVector2(refData.Transform, referenceNormal)                        // To world space
	referenceNormal = rl.Mat2MultiplyVector2(rl.Mat2Transpose(incData.Transform), referenceNormal) // To incident's model space

	// Find most anti-normal face on polygon
	incidentFace := 0
	minDot := float32(fltMax)

	for i := 0; i < incData.VertexCount; i++ {
		dot := rl.Vector2DotProduct(referenceNormal, incData.Normals[i])

		if dot < minDot {
			minDot = dot
			incidentFace = i
		}
	}

	// Assign face vertices for incident face
	*v0 = rl.Mat2MultiplyVector2(incData.Transform, incData.Vertices[incidentFace])
	*v0 = rl.Vector2Add(*v0, inc.Body.Position)

	if incidentFace+1 < incData.VertexCount {
		incidentFace = incidentFace + 1
	} else {
		incidentFace = 0
	}

	*v1 = rl.Mat2MultiplyVector2(incData.Transform, incData.Vertices[incidentFace])
	*v1 = rl.Vector2Add(*v1, inc.Body.Position)
}

// clip - Calculates clipping based on a normal and two faces
func clip(normal rl.Vector2, clip float32, faceA, faceB *rl.Vector2) int {
	sp := 0

	out := make([]rl.Vector2, 2)
	out[0] = *faceA
	out[1] = *faceB

	// Retrieve distances from each endpoint to the line
	distanceA := rl.Vector2DotProduct(normal, *faceA) - clip
	distanceB := rl.Vector2DotProduct(normal, *faceB) - clip

	// If negative (behind plane)
	if distanceA <= 0 {
		out[sp] = *faceA
		sp += 1
	}
	if distanceB <= 0 {
		out[sp] = *faceB
		sp += 1
	}

	// If the points are on different sides of the plane
	if distanceA*distanceB < 0 {
		// Push intersection point
		alpha := distanceA / (distanceA - distanceB)
		out[sp] = *faceA
		delta := rl.Vector2Subtract(*faceB, *faceA)
		delta.X *= alpha
		delta.Y *= alpha
		out[sp] = rl.Vector2Add(out[sp], delta)
		sp++
	}

	// Assign the new converted values
	*faceA = out[0]
	*faceB = out[1]

	return sp
}

// biasGreaterThan - Check if values are between bias range
func biasGreaterThan(valueA, valueB float32) bool {
	return valueA >= (valueB*0.95 + valueA*0.01)
}

// triangleBarycenter - Returns the barycenter of a triangle given by 3 points
func triangleBarycenter(v1, v2, v3 rl.Vector2) rl.Vector2 {
	result := rl.Vector2{}

	result.X = (v1.X + v2.X + v3.X) / 3
	result.Y = (v1.Y + v2.Y + v3.Y) / 3

	return result
}

// normalize - Normalize provided vector
func normalize(v *rl.Vector2) {
	var length, ilength float32

	aux := *v
	length = float32(math.Sqrt(float64(aux.X*aux.X + aux.Y*aux.Y)))

	if length == 0 {
		length = 1.0
	}

	ilength = 1.0 / length

	v.X *= ilength
	v.Y *= ilength
}
