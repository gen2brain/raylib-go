// Package physics - 2D Physics library for videogames
//
// A port of Victor Fisac's physac engine (https://github.com/raysan5/raylib/blob/master/src/physac.h)
package physics

import (
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
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
	Positions [maxVertices]rl.Vector2
	// Polygon vertex normals vectors
	Normals [maxVertices]rl.Vector2
}

// Shape type
type Shape struct {
	// Physics shape type (circle or polygon)
	Type ShapeType
	// Shape physics body reference
	Body *Body
	// Circle shape radius (used for circle shapes)
	Radius float32
	// Vertices transform matrix 2x2
	Transform rl.Mat2
	// Polygon shape vertices position and normals data (just used for polygon shapes)
	VertexData Polygon
}

// Body type
type Body struct {
	// Reference unique identifier
	ID int
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

// Manifold type
type Manifold struct {
	// Reference unique identifier
	ID int
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
	maxBodies      = 64
	maxManifolds   = 4096
	maxVertices    = 24
	circleVertices = 24

	collisionIterations   = 100
	penetrationAllowance  = 0.05
	penetrationCorrection = 0.4

	degToRad = math.Pi / 180
	epsilon  = 0.000001
	physacK  = 1.0 / 3.0
)

// Globals
var (
	// Offset time for MONOTONIC clock
	baseTime = time.Now()

	// Start time in milliseconds
	startTime float32

	// Delta time used for physics steps, in milliseconds
	deltaTime float32 = 1.0 / 60.0 / 10.0 * 1000

	// Current time in milliseconds
	currentTime float32

	// Hi-res clock frequency
	frequency uint64 = 0

	// Physics time step delta time accumulator
	accumulator float32

	// Physics world gravity force
	gravityForce = rl.NewVector2(0, 9.81)

	// Physics bodies pointers array
	bodies [64]*Body

	// Physics world current bodies counter
	bodiesCount int

	// Physics bodies pointers array
	manifolds [4096]*Manifold

	// Physics world current manifolds counter
	manifoldsCount int
)

// Init - Initializes physics values, pointers and creates physics loop thread
func Init() {
	initTimer()
	accumulator = 0
}

// SetGravity - Sets physics global gravity force
func SetGravity(x, y float32) {
	gravityForce.X = x
	gravityForce.Y = y
}

// NewBodyCircle - Creates a new circle physics body with generic parameters
func NewBodyCircle(pos rl.Vector2, radius, density float32) *Body {
	newID := findAvailableBodyIndex()
	if newID < 0 {
		return nil
	}

	// Initialize new body with generic values
	newBody := &Body{
		ID:              newID,
		Enabled:         true,
		Position:        pos,
		Velocity:        rl.Vector2{},
		Force:           rl.Vector2{},
		AngularVelocity: 0.0,
		Torque:          0.0,
		Orient:          0.0,
		Shape: Shape{
			Type:       CircleShape,
			Radius:     radius,
			Transform:  rl.Mat2Radians(0.0),
			VertexData: Polygon{},
		},
		StaticFriction:  0.4,
		DynamicFriction: 0.2,
		Restitution:     0.0,
		UseGravity:      true,
		IsGrounded:      false,
		FreezeOrient:    false,
	}

	newBody.Shape.Body = newBody

	newBody.Mass = math.Pi * radius * radius * density
	newBody.InverseMass = safeDiv(1.0, newBody.Mass)
	newBody.Inertia = newBody.Mass * radius * radius
	newBody.InverseInertia = safeDiv(1.0, newBody.Inertia)

	// Add new body to bodies pointers array and update bodies count
	bodies[bodiesCount] = newBody
	bodiesCount++
	return newBody
}

// NewBodyRectangle - Creates a new rectangle physics body with generic parameters
func NewBodyRectangle(pos rl.Vector2, width, height, density float32) *Body {
	newID := findAvailableBodyIndex()
	if newID < 0 {
		return nil
	}

	// Initialize new body with generic values
	newBody := &Body{
		ID:              newID,
		Enabled:         true,
		Position:        pos,
		Velocity:        rl.Vector2{},
		Force:           rl.Vector2{},
		AngularVelocity: 0.0,
		Torque:          0.0,
		Orient:          0.0,
		Shape: Shape{
			Type:       PolygonShape,
			Transform:  rl.Mat2Radians(0.0),
			VertexData: createRectanglePolygon(pos, rl.NewVector2(width, height)),
		},
		StaticFriction:  0.4,
		DynamicFriction: 0.2,
		Restitution:     0.0,
		UseGravity:      true,
		IsGrounded:      false,
		FreezeOrient:    false,
	}

	// Calculate centroid and moment of inertia
	newBody.Shape.Body = newBody
	var center rl.Vector2
	area := float32(0.0)
	inertia := float32(0.0)

	for i := 0; i < newBody.Shape.VertexData.VertexCount; i++ {
		// Triangle vertices, third vertex implied as (0, 0)
		nextIndex := getNextIndex(i, newBody.Shape.VertexData.VertexCount)
		p1 := newBody.Shape.VertexData.Positions[i]
		p2 := newBody.Shape.VertexData.Positions[nextIndex]

		D := rl.Vector2CrossProduct(p1, p2)
		triangleArea := D / 2

		area += triangleArea

		// Use area to weight the centroid average, not just vertex position
		center.X += triangleArea * physacK * (p1.X + p2.X)
		center.Y += triangleArea * physacK * (p1.Y + p2.Y)

		var intx2 float32 = p1.X*p1.X + p2.X*p1.X + p2.X*p2.X
		var inty2 float32 = p1.Y*p1.Y + p2.Y*p1.Y + p2.Y*p2.Y
		inertia += (0.25 * physacK * D) * (intx2 + inty2)
	}

	center.X *= 1.0 / area
	center.Y *= 1.0 / area

	// Translate vertices to centroid (make the centroid (0, 0) for the polygon in model space)
	// Note: this is not really necessary
	for i := 0; i < newBody.Shape.VertexData.VertexCount; i++ {
		newBody.Shape.VertexData.Positions[i].X -= center.X
		newBody.Shape.VertexData.Positions[i].Y -= center.Y
	}

	newBody.Mass = density * area
	newBody.InverseMass = safeDiv(1.0, newBody.Mass)
	newBody.Inertia = density * inertia
	newBody.InverseInertia = safeDiv(1.0, newBody.Inertia)

	// Add new body to bodies pointers array and update bodies count
	bodies[bodiesCount] = newBody
	bodiesCount++
	return newBody
}

// NewBodyPolygon - Creates a new polygon physics body with generic parameters
func NewBodyPolygon(pos rl.Vector2, radius float32, sides int, density float32) *Body {
	newID := findAvailableBodyIndex()
	if newID < 0 {
		return nil
	}

	// Initialize new body with generic values
	newBody := &Body{
		ID:              newID,
		Enabled:         true,
		Position:        pos,
		Velocity:        rl.Vector2{},
		Force:           rl.Vector2{},
		AngularVelocity: 0.0,
		Torque:          0.0,
		Orient:          0.0,
		Shape: Shape{
			Type:       PolygonShape,
			Transform:  rl.Mat2Radians(0),
			VertexData: createRandomPolygon(radius, sides),
		},
		StaticFriction:  0.4,
		DynamicFriction: 0.2,
		Restitution:     0.0,
		UseGravity:      true,
		IsGrounded:      false,
		FreezeOrient:    false,
	}

	newBody.Shape.Body = newBody

	// Calculate centroid and moment of inertia
	var center rl.Vector2
	area := float32(0.0)
	inertia := float32(0.0)

	for i := 0; i < newBody.Shape.VertexData.VertexCount; i++ {
		// Triangle vertices, third vertex implied as (0, 0)
		nextIndex := getNextIndex(i, newBody.Shape.VertexData.VertexCount)
		position1 := newBody.Shape.VertexData.Positions[i]
		position2 := newBody.Shape.VertexData.Positions[nextIndex]

		cross := rl.Vector2CrossProduct(position1, position2)
		triangleArea := cross / 2

		area += triangleArea

		// Use area to weight the centroid average, not just vertex position
		center.X += triangleArea * physacK * (position1.X + position2.X)
		center.Y += triangleArea * physacK * (position1.Y + position2.Y)

		intx2 := position1.X*position1.X + position2.X*position1.X + position2.X*position2.X
		inty2 := position1.Y*position1.Y + position2.Y*position1.Y + position2.Y*position2.Y
		inertia += (0.25 * physacK * cross) * (intx2 + inty2)
	}

	center.X *= 1.0 / area
	center.Y *= 1.0 / area

	// Translate vertices to centroid (make the centroid (0, 0) for the polygon in model space)
	// Note: this is not really necessary
	for i := 0; i < newBody.Shape.VertexData.VertexCount; i++ {
		newBody.Shape.VertexData.Positions[i].X -= center.X
		newBody.Shape.VertexData.Positions[i].Y -= center.Y
	}

	newBody.Mass = density * area
	newBody.InverseMass = safeDiv(1.0, newBody.Mass)
	newBody.Inertia = density * inertia
	newBody.InverseInertia = safeDiv(1.0, newBody.Inertia)

	// Add new body to bodies pointers array and update bodies count
	bodies[bodiesCount] = newBody
	bodiesCount++
	return newBody
}

// Reset - Destroys created physics bodies and manifolds
func Reset() {
	Close()
}

// AddForce - Adds a force to a physics body
func AddForce(body *Body, force rl.Vector2) {
	if body != nil {
		body.Force = rl.Vector2Add(body.Force, force)
	}
}

// AddTorque - Adds an angular force to a physics body
func AddTorque(body *Body, amount float32) {
	if body != nil {
		body.Torque += amount
	}
}

// Shatter - Shatters a polygon shape physics body to little physics bodies with explosion force
func Shatter(body *Body, position rl.Vector2, force float32) {
	if body == nil || body.Shape.Type != PolygonShape {
		return
	}

	vertexData := body.Shape.VertexData
	collision := false

	for i := 0; i < vertexData.VertexCount; i++ {
		nextIndex := getNextIndex(i, vertexData.VertexCount)
		pA := body.Position
		pB := rl.Mat2MultiplyVector2(body.Shape.Transform,
			rl.Vector2Add(body.Position, vertexData.Positions[i]))
		pC := rl.Mat2MultiplyVector2(body.Shape.Transform,
			rl.Vector2Add(body.Position, vertexData.Positions[nextIndex]))

		// Check collision between each triangle
		alpha := ((pB.Y-pC.Y)*(position.X-pC.X) + (pC.X-pB.X)*(position.Y-pC.Y)) /
			((pB.Y-pC.Y)*(pA.X-pC.X) + (pC.X-pB.X)*(pA.Y-pC.Y))
		beta := ((pC.Y-pA.Y)*(position.X-pC.X) + (pA.X-pC.X)*(position.Y-pC.Y)) /
			((pB.Y-pC.Y)*(pA.X-pC.X) + (pC.X-pB.X)*(pA.Y-pC.Y))
		gamma := 1.0 - alpha - beta

		if alpha > 0 && beta > 0 && gamma > 0 {
			collision = true
			break
		}
	}

	if !collision {
		return
	}

	count := vertexData.VertexCount
	bodyPos := body.Position
	trans := body.Shape.Transform

	vertices := make([]rl.Vector2, count)
	for i := 0; i < count; i++ {
		vertices[i] = vertexData.Positions[i]
	}

	// Destroy shattered physics body
	body.Destroy()
	for i := 0; i < count; i++ {
		nextIndex := getNextIndex(i, count)
		center := triangleBarycenter(vertices[i], vertices[nextIndex], rl.Vector2{})
		center = rl.Vector2Add(bodyPos, center)
		offset := rl.Vector2Subtract(center, bodyPos)

		var newBody *Body = NewBodyPolygon(center, 10, 3, 10)
		var newData Polygon = Polygon{}
		newData.VertexCount = 3
		newData.Positions[0] = rl.Vector2Subtract(vertices[i], offset)
		newData.Positions[1] = rl.Vector2Subtract(vertices[nextIndex], offset)
		newData.Positions[2] = rl.Vector2Subtract(position, center)

		// Separate vertices to avoid unnecessary physics collisions
		newData.Positions[0].X *= 0.95
		newData.Positions[0].Y *= 0.95
		newData.Positions[1].X *= 0.95
		newData.Positions[1].Y *= 0.95
		newData.Positions[2].X *= 0.95
		newData.Positions[2].Y *= 0.95

		// Calculate polygon faces normals
		for j := 0; j < newData.VertexCount; j++ {
			nextVertex := getNextIndex(j, newData.VertexCount)
			face := rl.Vector2Subtract(newData.Positions[nextVertex], newData.Positions[j])

			newData.Normals[j] = rl.NewVector2(face.Y, -face.X)
			normalize(&newData.Normals[j])
		}

		// Apply computed vertex data to new physics body shape
		newBody.Shape.VertexData = newData
		newBody.Shape.Transform = trans

		// Calculate centroid and moment of inertia
		center = rl.Vector2{}
		area := float32(0.0)
		inertia := float32(0.0)

		for j := 0; j < newBody.Shape.VertexData.VertexCount; j++ {
			// Triangle vertices, third vertex implied as (0, 0)
			nextVertex := getNextIndex(j, newBody.Shape.VertexData.VertexCount)
			p1 := newBody.Shape.VertexData.Positions[j]
			p2 := newBody.Shape.VertexData.Positions[nextVertex]

			D := rl.Vector2CrossProduct(p1, p2)
			triangleArea := D / 2

			area += triangleArea

			// Use area to weight the centroid average, not just vertex position
			center.X += triangleArea * physacK * (p1.X + p2.X)
			center.Y += triangleArea * physacK * (p1.Y + p2.Y)

			intx2 := p1.X*p1.X + p2.X*p1.X + p2.X*p2.X
			inty2 := p1.Y*p1.Y + p2.Y*p1.Y + p2.Y*p2.Y
			inertia += (0.25*physacK*D)*intx2 + inty2
		}

		center.X *= 1.0 / area
		center.Y *= 1.0 / area

		newBody.Mass = area
		newBody.InverseMass = safeDiv(1.0, newBody.Mass)
		newBody.Inertia = inertia
		newBody.InverseInertia = safeDiv(1.0, newBody.Inertia)

		// Calculate explosion force direction
		pointA := newBody.Position
		pointB := rl.Vector2Subtract(newData.Positions[1], newData.Positions[0])
		pointB.X /= 2.0
		pointB.Y /= 2.0
		forceDirection := rl.Vector2Subtract(
			rl.Vector2Add(pointA, rl.Vector2Add(newData.Positions[0], pointB)),
			newBody.Position,
		)
		normalize(&forceDirection)
		forceDirection.X *= force
		forceDirection.Y *= force

		// Apply force to new physics body
		AddForce(newBody, forceDirection)
	}
}

// GetBodies - Returns the slice of created physics bodies
func GetBodies() []*Body {
	return bodies[:]
}

// GetBodiesCount - Returns the current amount of created physics bodies
func GetBodiesCount() int {
	return bodiesCount
}

// GetBody - Returns a physics body of the bodies pool at a specific index
func GetBody(index int) *Body {
	return bodies[index]
}

// GetShapeType - Returns the physics body shape type (PHYSICS_CIRCLE or PHYSICS_POLYGON)
func GetShapeType(index int) ShapeType {
	result := ShapeType(-1)
	if index < bodiesCount {
		if bodies[index] != nil {
			result = bodies[index].Shape.Type
		}
	}
	return result
}

// GetShapeVerticesCount - Returns the amount of vertices of a physics body shape
func GetShapeVerticesCount(index int) int {
	var result int = 0
	if index < bodiesCount {
		if bodies[index] != nil {
			switch bodies[index].Shape.Type {
			case CircleShape:
				result = circleVertices
			case PolygonShape:
				result = bodies[index].Shape.VertexData.VertexCount
			default:
			}
		}
	}
	return result
}

// GetShapeVertex - Returns transformed position of a body shape (body position + vertex transformed position)
func (b *Body) GetShapeVertex(vertex int) rl.Vector2 {
	var position rl.Vector2

	switch b.Shape.Type {
	case CircleShape:
		angle := 360.0 / circleVertices * float64(vertex) * (degToRad)
		position.X = b.Position.X + float32(math.Cos(angle))*b.Shape.Radius
		position.Y = b.Position.Y + float32(math.Sin(angle))*b.Shape.Radius
	case PolygonShape:
		position = rl.Vector2Add(
			b.Position,
			rl.Mat2MultiplyVector2(b.Shape.Transform, b.Shape.VertexData.Positions[vertex]),
		)
	}
	return position
}

// SetBodyRotation - Sets physics body shape transform based on radians parameter
func (b *Body) SetRotation(radians float32) {
	b.Orient = radians
	if b.Shape.Type == PolygonShape {
		b.Shape.Transform = rl.Mat2Radians(radians)
	}
}

// Destroy - Unitializes and destroy a physics body
func (b *Body) Destroy() {
	id := b.ID
	index := -1
	for i := 0; i < bodiesCount; i++ {
		if bodies[i].ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}

	bodies[index] = nil

	// Reorder physics bodies pointers array and its catched index
	for i := index; i+1 < bodiesCount; i++ {
		bodies[i] = bodies[i+1]
	}

	// Update physics bodies count
	bodiesCount--
}

// Close - Unitializes physics pointers
func Close() {
	// Unitialize physics manifolds dynamic memory allocations
	for i := manifoldsCount - 1; i >= 0; i-- {
		destroyManifold(manifolds[i])
	}

	// Unitialize physics bodies dynamic memory allocations
	for i := bodiesCount - 1; i >= 0; i-- {
		bodies[i].Destroy()
	}
}

// findAvailableBodyIndex - Finds a valid index for a new physics body initialization
func findAvailableBodyIndex() int {
	index := -1
	for i := 0; i < maxBodies; i++ {
		currentID := i

		// Check if current id already exist in other physics body
		for k := 0; k < bodiesCount; k++ {
			if bodies[k].ID == currentID {
				currentID++
				break
			}
		}

		// If it is not used, use it as new physics body id
		if currentID == i {
			index = i
			break
		}
	}
	return index
}

// createRandomPolygon - Creates a random polygon shape with max vertex distance from polygon pivot
func createRandomPolygon(radius float32, sides int) Polygon {
	var data Polygon = Polygon{}
	data.VertexCount = sides

	// Calculate polygon vertices positions
	for i := 0; i < data.VertexCount; i++ {
		data.Positions[i].X = float32(math.Cos(360.0/float64(sides)*float64(i)*degToRad)) * radius
		data.Positions[i].Y = float32(math.Sin(360.0/float64(sides)*float64(i)*degToRad)) * radius
	}

	// Calculate polygon faces normals
	for i := 0; i < data.VertexCount; i++ {
		nextIndex := getNextIndex(i, sides)
		face := rl.Vector2Subtract(data.Positions[nextIndex], data.Positions[i])

		data.Normals[i] = rl.NewVector2(face.Y, -face.X)
		normalize(&data.Normals[i])
	}

	return data
}

// createRectanglePolygon - Creates a rectangle polygon shape based on a min and max positions
func createRectanglePolygon(pos rl.Vector2, size rl.Vector2) Polygon {
	var data Polygon = Polygon{}
	data.VertexCount = 4

	// Calculate polygon vertices positions
	data.Positions[0] = rl.NewVector2(pos.X+size.X/2, pos.Y-size.Y/2)
	data.Positions[1] = rl.NewVector2(pos.X+size.X/2, pos.Y+size.Y/2)
	data.Positions[2] = rl.NewVector2(pos.X-size.X/2, pos.Y+size.Y/2)
	data.Positions[3] = rl.NewVector2(pos.X-size.X/2, pos.Y-size.Y/2)

	// Calculate polygon faces normals
	for i := 0; i < data.VertexCount; i++ {
		nextIndex := getNextIndex(i, data.VertexCount)
		face := rl.Vector2Subtract(data.Positions[nextIndex], data.Positions[i])

		data.Normals[i] = rl.NewVector2(face.Y, -face.X)
		normalize(&data.Normals[i])
	}

	return data
}

// step - Does physics steps calculations (dynamics, collisions and position corrections)
func step() {
	// Clear previous generated collisions information
	for i := manifoldsCount - 1; i >= 0; i-- {
		if manifold := manifolds[i]; manifold != nil {
			destroyManifold(manifold)
		}
	}

	// Reset physics bodies grounded state
	for i := 0; i < bodiesCount; i++ {
		bodies[i].IsGrounded = false
	}

	// Generate new collision information
	for i := 0; i < bodiesCount; i++ {
		bodyA := bodies[i]
		if bodyA == nil {
			continue
		}

		for j := i + 1; j < bodiesCount; j++ {
			var bodyB *Body = bodies[j]
			if bodyB == nil || bodyA.InverseMass == 0 && bodyB.InverseMass == 0 {
				continue
			}

			manifold := createManifold(bodyA, bodyB)
			solveManifold(manifold)

			if manifold.ContactsCount > 0 {
				// Create a new manifold with same information as previously solved manifold and add it to the manifolds pool last slot
				newManifold := createManifold(bodyA, bodyB)
				newManifold.Penetration = manifold.Penetration
				newManifold.Normal = manifold.Normal
				newManifold.Contacts[0] = manifold.Contacts[0]
				newManifold.Contacts[1] = manifold.Contacts[1]
				newManifold.ContactsCount = manifold.ContactsCount
				newManifold.Restitution = manifold.Restitution
				newManifold.DynamicFriction = manifold.DynamicFriction
				newManifold.StaticFriction = manifold.StaticFriction
			}
		}
	}

	// Integrate forces to physics bodies
	for i := 0; i < bodiesCount; i++ {
		if body := bodies[i]; body != nil {
			integrateForces(body)
		}
	}

	// Initialize physics manifolds to solve collisions
	for i := 0; i < manifoldsCount; i++ {
		if manifold := manifolds[i]; manifold != nil {
			initializeManifolds(manifold)
		}
	}

	// Integrate physics collisions impulses to solve collisions
	for i := 0; i < collisionIterations; i++ {
		for j := 0; j < manifoldsCount; j++ {
			if manifold := manifolds[i]; manifold != nil {
				integrateImpulses(manifold)
			}
		}
	}

	// Integrate velocity to physics bodies
	for i := 0; i < bodiesCount; i++ {
		if body := bodies[i]; body != nil {
			integrateVelocity(body)
		}
	}

	// Correct physics bodies positions based on manifolds collision information
	for i := 0; i < manifoldsCount; i++ {
		if manifold := manifolds[i]; manifold != nil {
			correctPositions(manifold)
		}
	}

	// Clear physics bodies forces
	for i := 0; i < bodiesCount; i++ {
		if body := bodies[i]; body != nil {
			body.Force = rl.Vector2{}
			body.Torque = 0
		}
	}
}

// Update - Runs physics step
func Update() {
	// Calculate current time
	currentTime = getCurrentTime()

	// Calculate current delta time
	var delta float32 = currentTime - startTime

	// Store the time elapsed since the last frame began
	accumulator += delta

	// Fixed time stepping loop
	for accumulator >= deltaTime {
		step()
		accumulator -= deltaTime
	}

	// Record the starting of this frame
	startTime = currentTime
}

// SetTimeStep - Sets physics fixed time step in milliseconds. 1.666666 by default
func SetTimeStep(delta float32) {
	deltaTime = delta
}

// findAvailableManifoldIndex - Finds a valid index for a new manifold initialization
func findAvailableManifoldIndex() int {
	index := -1
	for i := 0; i < maxManifolds; i++ {
		var currentId int = i

		// Check if current id already exist in other physics body
		for k := 0; k < manifoldsCount; k++ {
			if manifolds[k].ID == currentId {
				currentId++
				break
			}
		}

		// If it is not used, use it as new physics body id
		if currentId == i {
			index = i
			break
		}
	}
	return index
}

// createManifold - Creates a new physics manifold to solve collision
func createManifold(a *Body, b *Body) *Manifold {
	newID := findAvailableManifoldIndex()
	if newID < 0 {
		return nil
	}

	// Initialize new manifold with generic values
	newManifold := &Manifold{
		ID:              newID,
		BodyA:           a,
		BodyB:           b,
		Penetration:     0,
		Normal:          rl.Vector2{},
		Contacts:        [2]rl.Vector2{},
		ContactsCount:   0,
		Restitution:     0.0,
		DynamicFriction: 0.0,
		StaticFriction:  0.0,
	}

	// Add new contact to conctas pointers array and update contacts count
	manifolds[manifoldsCount] = newManifold
	manifoldsCount++

	return newManifold
}

// destroyManifold - Unitializes and destroys a physics manifold
func destroyManifold(manifold *Manifold) {
	if manifold == nil {
		return
	}

	id := manifold.ID
	index := -1
	for i := 0; i < manifoldsCount; i++ {
		if manifolds[i].ID == id {
			index = i
			break
		}
	}
	if index < 0 {
		return
	}

	manifolds[index] = nil

	// Reorder physics manifolds pointers array and its catched index
	for i := index; i < manifoldsCount; i++ {
		if (i + 1) < manifoldsCount {
			manifolds[i] = manifolds[i+1]
		}
	}

	// Update physics manifolds count
	manifoldsCount--
}

// solveManifold - Solves a created physics manifold between two physics bodies
func solveManifold(manifold *Manifold) {
	switch manifold.BodyA.Shape.Type {
	case CircleShape:
		switch manifold.BodyB.Shape.Type {
		case CircleShape:
			solveCircleToCircle(manifold)
		case PolygonShape:
			solveCircleToPolygon(manifold)
		}
	case PolygonShape:
		switch manifold.BodyB.Shape.Type {
		case CircleShape:
			solvePolygonToCircle(manifold)
		case PolygonShape:
			solvePolygonToPolygon(manifold)
		}
	}

	// Update physics body grounded state if normal direction is down and grounded state
	// is not set yet in previous manifolds
	if !manifold.BodyB.IsGrounded {
		manifold.BodyB.IsGrounded = manifold.Normal.Y < 0
	}
}

// solveCircleToCircle - Solves collision between two circle shape physics bodies
func solveCircleToCircle(manifold *Manifold) {
	bodyA, bodyB := manifold.BodyA, manifold.BodyB
	if bodyA == nil || bodyB == nil {
		return
	}

	// Calculate translational vector, which is normal
	var normal rl.Vector2 = rl.Vector2Subtract(bodyB.Position, bodyA.Position)

	distSqr := rl.Vector2LenSqr(normal)
	radius := bodyA.Shape.Radius + bodyB.Shape.Radius

	// Check if circles are not in contact
	if distSqr >= radius*radius {
		manifold.ContactsCount = 0
		return
	}

	distance := float32(math.Sqrt(float64(distSqr)))
	manifold.ContactsCount = 1
	if distance == 0 {
		manifold.Penetration = bodyA.Shape.Radius
		manifold.Normal = rl.NewVector2(1, 0)
		manifold.Contacts[0] = bodyA.Position
	} else {
		manifold.Penetration = radius - distance
		// Faster than using normalize() due to sqrt is already performed
		manifold.Normal = rl.NewVector2(
			normal.X/distance,
			normal.Y/distance,
		)
		manifold.Contacts[0] = rl.NewVector2(
			manifold.Normal.X*bodyA.Shape.Radius+bodyA.Position.X,
			manifold.Normal.Y*bodyA.Shape.Radius+bodyA.Position.Y,
		)
	}

	// Update physics body grounded state if normal direction is down
	if !bodyA.IsGrounded {
		bodyA.IsGrounded = manifold.Normal.Y < 0
	}
}

// solveCircleToPolygon - Solves collision between a circle to a polygon shape physics bodies
func solveCircleToPolygon(manifold *Manifold) {
	bodyA, bodyB := manifold.BodyA, manifold.BodyB
	if bodyA == nil || bodyB == nil {
		return
	}
	solveDifferentShapes(manifold, bodyA, bodyB)
}

// solvePolygonToCircle - Solves collision between a polygon to a circle shape physics bodies
func solvePolygonToCircle(manifold *Manifold) {
	bodyA, bodyB := manifold.BodyA, manifold.BodyB
	if bodyA == nil || bodyB == nil {
		return
	}
	solveDifferentShapes(manifold, bodyB, bodyA)
	manifold.Normal.X *= -1.0
	manifold.Normal.Y *= -1.0
}

// solveDifferentShapes - Solves collision between two different types of shapes
func solveDifferentShapes(manifold *Manifold, bodyA *Body, bodyB *Body) {
	manifold.ContactsCount = 0

	// Transform circle center to polygon transform space
	center := rl.Mat2MultiplyVector2(
		rl.Mat2Transpose(bodyB.Shape.Transform),
		rl.Vector2Subtract(bodyA.Position, bodyB.Position),
	)

	// Find edge with minimum penetration
	// It is the same concept as using support points in SolvePolygonToPolygon
	separation := float32(-math.MaxFloat32)
	faceNormal := 0
	vertexData := bodyB.Shape.VertexData

	for i := 0; i < vertexData.VertexCount; i++ {
		currentSeparation := rl.Vector2DotProduct(
			vertexData.Normals[i],
			rl.Vector2Subtract(center, vertexData.Positions[i]),
		)

		if currentSeparation > bodyA.Shape.Radius {
			return
		}

		if currentSeparation > separation {
			separation = currentSeparation
			faceNormal = i
		}
	}

	// Grab face's vertices
	nextIndex := getNextIndex(faceNormal, vertexData.VertexCount)
	v1 := vertexData.Positions[faceNormal]
	v2 := vertexData.Positions[nextIndex]

	// Check to see if center is within polygon
	if separation < epsilon {
		manifold.ContactsCount = 1
		var normal rl.Vector2 = rl.Mat2MultiplyVector2(bodyB.Shape.Transform, vertexData.Normals[faceNormal])
		manifold.Normal = rl.NewVector2(-normal.X, -normal.Y)
		manifold.Contacts[0] = rl.NewVector2(
			manifold.Normal.X*bodyA.Shape.Radius+bodyA.Position.X,
			manifold.Normal.Y*bodyA.Shape.Radius+bodyA.Position.Y,
		)
		manifold.Penetration = bodyA.Shape.Radius
		return
	}

	// Determine which voronoi region of the edge center of circle lies within
	dot1 := rl.Vector2DotProduct(rl.Vector2Subtract(center, v1), rl.Vector2Subtract(v2, v1))
	dot2 := rl.Vector2DotProduct(rl.Vector2Subtract(center, v2), rl.Vector2Subtract(v1, v2))
	manifold.Penetration = bodyA.Shape.Radius - separation

	switch {
	case dot1 <= 0: // Closest to v1
		if rl.Vector2Distance(center, v1) > bodyA.Shape.Radius*bodyA.Shape.Radius {
			return
		}

		manifold.ContactsCount = 1
		var normal rl.Vector2 = rl.Vector2Subtract(v1, center)
		normal = rl.Mat2MultiplyVector2(bodyB.Shape.Transform, normal)
		normalize(&normal)
		manifold.Normal = normal
		v1 = rl.Mat2MultiplyVector2(bodyB.Shape.Transform, v1)
		v1 = rl.Vector2Add(v1, bodyB.Position)
		manifold.Contacts[0] = v1

	case dot2 <= 0: // Closest to v2
		if rl.Vector2Distance(center, v2) > bodyA.Shape.Radius*bodyA.Shape.Radius {
			return
		}

		manifold.ContactsCount = 1
		var normal rl.Vector2 = rl.Vector2Subtract(v2, center)
		v2 = rl.Mat2MultiplyVector2(bodyB.Shape.Transform, v2)
		v2 = rl.Vector2Add(v2, bodyB.Position)
		manifold.Contacts[0] = v2
		normal = rl.Mat2MultiplyVector2(bodyB.Shape.Transform, normal)
		normalize(&normal)
		manifold.Normal = normal

	default: // Closest to face
		var normal rl.Vector2 = vertexData.Normals[faceNormal]

		if rl.Vector2DotProduct(rl.Vector2Subtract(center, v1), normal) > bodyA.Shape.Radius {
			return
		}

		normal = rl.Mat2MultiplyVector2(bodyB.Shape.Transform, normal)
		manifold.Normal = rl.NewVector2(-normal.X, -normal.Y)
		manifold.Contacts[0] = rl.NewVector2(
			manifold.Normal.X*bodyA.Shape.Radius+bodyA.Position.X,
			manifold.Normal.Y*bodyA.Shape.Radius+bodyA.Position.Y,
		)
		manifold.ContactsCount = 1
	}
}

// solvePolygonToPolygon - Solves collision between two polygons shape physics bodies
func solvePolygonToPolygon(manifold *Manifold) {
	bodyA, bodyB := manifold.BodyA, manifold.BodyB
	if bodyA == nil || bodyB == nil {
		return
	}

	shapeA, shapeB := bodyA.Shape, bodyB.Shape
	manifold.ContactsCount = 0

	// Check for separating axis with A shape's face planes
	faceA, penetrationA := findAxisLeastPenetration(shapeA, shapeB)
	if penetrationA >= 0 {
		return
	}

	// Check for separating axis with B shape's face planes
	faceB, penetrationB := findAxisLeastPenetration(shapeB, shapeA)
	if penetrationB >= 0 {
		return
	}

	referenceIndex := 0
	flip := false // Always point from A shape to B shape

	var refPoly Shape // Reference
	var incPoly Shape // Incident

	// Determine which shape contains reference face
	if biasGreaterThan(penetrationA, penetrationB) {
		refPoly = shapeA
		incPoly = shapeB
		referenceIndex = faceA
	} else {
		refPoly = shapeB
		incPoly = shapeA
		referenceIndex = faceB
		flip = true
	}

	// World space incident face
	var incidentFace [2]rl.Vector2
	findIncidentFace(&incidentFace[0], &incidentFace[1], refPoly, incPoly, referenceIndex)

	// Setup reference face vertices
	refData := refPoly.VertexData
	v1 := refData.Positions[referenceIndex]
	referenceIndex = getNextIndex(referenceIndex, refData.VertexCount)
	v2 := refData.Positions[referenceIndex]

	// Transform vertices to world space
	v1 = rl.Mat2MultiplyVector2(refPoly.Transform, v1)
	v1 = rl.Vector2Add(v1, refPoly.Body.Position)
	v2 = rl.Mat2MultiplyVector2(refPoly.Transform, v2)
	v2 = rl.Vector2Add(v2, refPoly.Body.Position)

	// Calculate reference face side normal in world space
	sidePlaneNormal := rl.Vector2Subtract(v2, v1)
	normalize(&sidePlaneNormal)

	// Orthogonalize
	refFaceNormal := rl.NewVector2(sidePlaneNormal.Y, -sidePlaneNormal.X)
	refC := rl.Vector2DotProduct(refFaceNormal, v1)
	negSide := rl.Vector2DotProduct(sidePlaneNormal, v1) * float32(-1)
	posSide := rl.Vector2DotProduct(sidePlaneNormal, v2)

	// Clip incident face to reference face side planes (due to floating point error, possible to not have required points
	if clip(rl.NewVector2(-sidePlaneNormal.X, -sidePlaneNormal.Y), negSide, &incidentFace[0], &incidentFace[1]) < 2 {
		return
	}

	if clip(sidePlaneNormal, posSide, &incidentFace[0], &incidentFace[1]) < 2 {
		return
	}

	// Flip normal if required
	if flip {
		manifold.Normal = rl.NewVector2(-refFaceNormal.X, -refFaceNormal.Y)
	} else {
		manifold.Normal = refFaceNormal
	}

	// Keep points behind reference face
	currentPoint := 0 // Clipped points behind reference face
	separation := rl.Vector2DotProduct(refFaceNormal, incidentFace[0]) - refC

	if separation <= 0 {
		manifold.Contacts[currentPoint] = incidentFace[0]
		manifold.Penetration = -separation
		currentPoint++
	} else {
		manifold.Penetration = 0
	}

	separation = rl.Vector2DotProduct(refFaceNormal, incidentFace[1]) - refC

	if separation <= 0 {
		manifold.Contacts[currentPoint] = incidentFace[1]
		manifold.Penetration += -separation
		currentPoint++

		// Calculate total penetration average
		manifold.Penetration /= float32(currentPoint)
	}

	manifold.ContactsCount = currentPoint
}

// integrateForces - Integrates physics forces into velocity
func integrateForces(body *Body) {
	if body == nil || body.InverseMass == 0 || !body.Enabled {
		return
	}

	body.Velocity.X += body.Force.X * body.InverseMass * (deltaTime / 2.0)
	body.Velocity.Y += body.Force.Y * body.InverseMass * (deltaTime / 2.0)

	if body.UseGravity {
		body.Velocity.X += gravityForce.X * (deltaTime / 1000 / 2.0)
		body.Velocity.Y += gravityForce.Y * (deltaTime / 1000 / 2.0)
	}

	if !body.FreezeOrient {
		body.AngularVelocity += body.Torque * body.InverseInertia * (deltaTime / 2.0)
	}
}

// initializeManifolds - Initializes physics manifolds to solve collisions
func initializeManifolds(manifold *Manifold) {
	bodyA, bodyB := manifold.BodyA, manifold.BodyB

	if bodyA == nil || bodyB == nil {
		return
	}

	// // Calculate average restitution, static and dynamic friction
	manifold.Restitution = float32(math.Sqrt(float64(bodyA.Restitution * bodyB.Restitution)))
	manifold.StaticFriction = float32(math.Sqrt(float64(bodyA.StaticFriction * bodyB.StaticFriction)))
	manifold.DynamicFriction = float32(math.Sqrt(float64(bodyA.DynamicFriction * bodyB.DynamicFriction)))

	for i := 0; i < manifold.ContactsCount; i++ {
		// Caculate radius from center of mass to contact
		radiusA := rl.Vector2Subtract(manifold.Contacts[i], bodyA.Position)
		radiusB := rl.Vector2Subtract(manifold.Contacts[i], bodyB.Position)

		crossA := rl.Vector2Cross(bodyA.AngularVelocity, radiusA)
		crossB := rl.Vector2Cross(bodyB.AngularVelocity, radiusB)
		radiusV := rl.NewVector2(
			bodyB.Velocity.X+crossB.X-bodyA.Velocity.X-crossA.X,
			bodyB.Velocity.Y+crossB.Y-bodyA.Velocity.Y-crossA.Y,
		)

		// Determine if we should perform a resting collision or not;
		// The idea is if the only thing moving this object is gravity, then the collision should be
		// performed without any restitution
		rad := rl.NewVector2(gravityForce.X*deltaTime/1000, gravityForce.Y*deltaTime/1000)
		if rl.Vector2LenSqr(radiusV) < (rl.Vector2LenSqr(rad) + epsilon) {
			manifold.Restitution = 0
		}
	}
}

// integrateImpulses - Integrates physics collisions impulses to solve collisions
func integrateImpulses(manifold *Manifold) {
	bodyA, bodyB := manifold.BodyA, manifold.BodyB

	if bodyA == nil || bodyB == nil {
		return
	}

	// Early out and positional correct if both objects have infinite mass
	if math.Abs(float64(bodyA.InverseMass+bodyB.InverseMass)) <= epsilon {
		bodyA.Velocity = rl.Vector2{}
		bodyB.Velocity = rl.Vector2{}
		return
	}

	for i := 0; i < manifold.ContactsCount; i++ {
		// Calculate radius from center of mass to contact
		radiusA := rl.Vector2Subtract(manifold.Contacts[i], bodyA.Position)
		radiusB := rl.Vector2Subtract(manifold.Contacts[i], bodyB.Position)

		// Calculate relative velocity
		radiusV := rl.NewVector2(
			bodyB.Velocity.X+rl.Vector2Cross(bodyB.AngularVelocity, radiusB).X-
				bodyA.Velocity.X-rl.Vector2Cross(bodyA.AngularVelocity, radiusA).X,
			bodyB.Velocity.Y+rl.Vector2Cross(bodyB.AngularVelocity, radiusB).Y-
				bodyA.Velocity.Y-rl.Vector2Cross(bodyA.AngularVelocity, radiusA).Y,
		)

		// Relative velocity along the normal
		contactVelocity := rl.Vector2DotProduct(radiusV, manifold.Normal)

		// Do not resolve if velocities are separating
		if contactVelocity > 0 {
			return
		}

		raCrossN := rl.Vector2CrossProduct(radiusA, manifold.Normal)
		rbCrossN := rl.Vector2CrossProduct(radiusB, manifold.Normal)

		inverseMassSum := bodyA.InverseMass + bodyB.InverseMass +
			(raCrossN*raCrossN)*bodyA.InverseInertia + (rbCrossN*rbCrossN)*bodyB.InverseInertia

		// Calculate impulse scalar value
		impulse := -(manifold.Restitution + 1.0) * contactVelocity
		impulse /= inverseMassSum
		impulse /= float32(manifold.ContactsCount)

		// Apply impulse to each physics body
		impulseV := rl.NewVector2(manifold.Normal.X*impulse, manifold.Normal.Y*impulse)

		if bodyA.Enabled {
			bodyA.Velocity.X += bodyA.InverseMass * (-impulseV.X)
			bodyA.Velocity.Y += bodyA.InverseMass * (-impulseV.Y)

			if !bodyA.FreezeOrient {
				bodyA.AngularVelocity += bodyA.InverseInertia *
					rl.Vector2CrossProduct(radiusA, rl.NewVector2(-impulseV.X, -impulseV.Y))
			}
		}

		if bodyB.Enabled {
			bodyB.Velocity.X += bodyB.InverseMass * impulseV.X
			bodyB.Velocity.Y += bodyB.InverseMass * impulseV.Y

			if !bodyB.FreezeOrient {
				bodyB.AngularVelocity += bodyB.InverseInertia * rl.Vector2CrossProduct(radiusB, impulseV)
			}
		}

		// Apply friction impulse to each physics body
		radiusV.X = 0 +
			bodyB.Velocity.X + rl.Vector2Cross(bodyB.AngularVelocity, radiusB).X -
			bodyA.Velocity.X - rl.Vector2Cross(bodyA.AngularVelocity, radiusA).X
		radiusV.Y = 0 +
			bodyB.Velocity.Y + rl.Vector2Cross(bodyB.AngularVelocity, radiusB).Y -
			bodyA.Velocity.Y - rl.Vector2Cross(bodyA.AngularVelocity, radiusA).Y

		tangent := rl.NewVector2(
			radiusV.X-manifold.Normal.X*rl.Vector2DotProduct(radiusV, manifold.Normal),
			radiusV.Y-manifold.Normal.Y*rl.Vector2DotProduct(radiusV, manifold.Normal),
		)
		normalize(&tangent)

		// Calculate impulse tangent magnitude
		impulseTangent := -rl.Vector2DotProduct(radiusV, tangent)
		impulseTangent /= inverseMassSum
		impulseTangent /= float32(manifold.ContactsCount)

		absImpulseTangent := float32(math.Abs(float64(impulseTangent)))

		// Don't apply tiny friction impulses
		if absImpulseTangent <= epsilon {
			return
		}

		// Apply coulumb's law
		var tangentImpulse rl.Vector2
		if absImpulseTangent < impulse*manifold.StaticFriction {
			tangentImpulse = rl.NewVector2(tangent.X*impulseTangent, tangent.Y*impulseTangent)
		} else {
			tangentImpulse = rl.NewVector2(
				tangent.X*(-impulse)*manifold.DynamicFriction,
				tangent.Y*(-impulse)*manifold.DynamicFriction,
			)
		}

		// Apply friction impulse
		if bodyA.Enabled {
			bodyA.Velocity.X += bodyA.InverseMass * (-tangentImpulse.X)
			bodyA.Velocity.Y += bodyA.InverseMass * (-tangentImpulse.Y)

			if !bodyA.FreezeOrient {
				bodyA.AngularVelocity += bodyA.InverseInertia *
					rl.Vector2CrossProduct(radiusA, rl.NewVector2(-tangentImpulse.X, -tangentImpulse.Y))
			}
		}

		if bodyB.Enabled {
			bodyB.Velocity.X += bodyB.InverseMass * tangentImpulse.X
			bodyB.Velocity.Y += bodyB.InverseMass * tangentImpulse.Y

			if !bodyB.FreezeOrient {
				bodyB.AngularVelocity += bodyB.InverseInertia * rl.Vector2CrossProduct(radiusB, tangentImpulse)
			}
		}
	}
}

// integrateVelocity - Integrates physics velocity into position and forces
func integrateVelocity(body *Body) {
	if body == nil || !body.Enabled {
		return
	}

	body.Position.X += body.Velocity.X * deltaTime
	body.Position.Y += body.Velocity.Y * deltaTime

	if !body.FreezeOrient {
		body.Orient += body.AngularVelocity * deltaTime
	}

	rl.Mat2Set(&body.Shape.Transform, body.Orient)

	integrateForces(body)
}

// correctPositions - Corrects physics bodies positions based on manifolds collision information
func correctPositions(manifold *Manifold) {
	bodyA, bodyB := manifold.BodyA, manifold.BodyB
	if bodyA == nil || bodyB == nil {
		return
	}

	corrCoeff := float32(math.Max(float64(manifold.Penetration-penetrationAllowance), 0)) /
		(bodyA.InverseMass + bodyB.InverseMass) * penetrationCorrection
	correction := rl.NewVector2(corrCoeff*manifold.Normal.X, corrCoeff*manifold.Normal.Y)

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
	bestProjection := float32(-math.MaxFloat32)
	bestVertex := rl.Vector2{}

	for i := 0; i < shape.VertexData.VertexCount; i++ {
		vertex := shape.VertexData.Positions[i]
		projection := rl.Vector2DotProduct(vertex, dir)

		if projection > bestProjection {
			bestVertex = vertex
			bestProjection = projection
		}
	}

	return bestVertex
}

// findAxisLeastPenetration - Finds polygon shapes axis least penetration
func findAxisLeastPenetration(shapeA Shape, shapeB Shape) (int, float32) {
	bestIndex := 0
	bestDistance := float32(-math.MaxFloat32)

	dataA := shapeA.VertexData

	for i := 0; i < dataA.VertexCount; i++ {
		// Retrieve a face normal from A shape
		normal := dataA.Normals[i]
		transNormal := rl.Mat2MultiplyVector2(shapeA.Transform, normal)

		// Transform face normal into B shape's model space
		buT := rl.Mat2Transpose(shapeB.Transform)
		normal = rl.Mat2MultiplyVector2(buT, transNormal)

		// Retrieve vertex on face from A shape, transform into B shape's model space
		support := getSupport(shapeB, rl.NewVector2(-normal.X, -normal.Y))

		// Retrieve vertex on face from A shape, transform into B shape's model space
		vertex := dataA.Positions[i]
		vertex = rl.Mat2MultiplyVector2(shapeA.Transform, vertex)
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

	return bestIndex, float32(bestDistance)
}

// findIncidentFace - Finds two polygon shapes incident face
func findIncidentFace(v0 *rl.Vector2, v1 *rl.Vector2, ref Shape, inc Shape, index int) {
	refData := ref.VertexData
	incData := inc.VertexData

	refNormal := refData.Normals[index]

	// Calculate normal in incident's frame of reference
	refNormal = rl.Mat2MultiplyVector2(ref.Transform, refNormal)                   // To world space
	refNormal = rl.Mat2MultiplyVector2(rl.Mat2Transpose(inc.Transform), refNormal) // To incident's model space

	// Find most anti-normal face on polygon
	incidentFace := 0
	minDot := float32(math.MaxFloat32)
	for i := 0; i < incData.VertexCount; i++ {
		dot := rl.Vector2DotProduct(refNormal, incData.Normals[i])
		if dot < minDot {
			minDot = dot
			incidentFace = i
		}
	}

	// Assign face vertices for incident face
	*v0 = rl.Mat2MultiplyVector2(inc.Transform, incData.Positions[incidentFace])
	*v0 = rl.Vector2Add(*v0, inc.Body.Position)
	incidentFace = getNextIndex(incidentFace, incData.VertexCount)
	*v1 = rl.Mat2MultiplyVector2(inc.Transform, incData.Positions[incidentFace])
	*v1 = rl.Vector2Add(*v1, inc.Body.Position)
}

// clip - Calculates clipping based on a normal and two faces
func clip(normal rl.Vector2, clip float32, faceA *rl.Vector2, faceB *rl.Vector2) int {
	sp := 0
	out := [2]rl.Vector2{*faceA, *faceB}

	// Retrieve distances from each endpoint to the line
	distanceA := rl.Vector2DotProduct(normal, *faceA) - clip
	distanceB := rl.Vector2DotProduct(normal, *faceB) - clip

	// If negative (behind plane)
	if distanceA <= 0 {
		out[sp] = *faceA
		sp++
	}

	if distanceB <= 0 {
		out[sp] = *faceB
		sp++
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

// biasGreaterThan - Checks if values are between bias range
func biasGreaterThan(valueA float32, valueB float32) bool {
	return valueA >= (valueB*0.95 + valueA*0.01)
}

// triangleBarycenter - Returns the barycenter of a triangle given by 3 points
func triangleBarycenter(v1 rl.Vector2, v2 rl.Vector2, v3 rl.Vector2) rl.Vector2 {
	return rl.NewVector2(
		(v1.X+v2.X+v3.X)/3,
		(v1.Y+v2.Y+v3.Y)/3,
	)
}

// initTimer - Initializes hi-resolution MONOTONIC timer
func initTimer() {
	rand.Seed(getTimeCount())
	frequency = 1000000000
	startTime = getCurrentTime() // Get current time
}

// getTimeCount - Gets hi-res MONOTONIC time measure in nanoseconds
func getTimeCount() int64 {
	return time.Since(baseTime).Nanoseconds()
}

// getCurrentTime - Gets current time measure in milliseconds
func getCurrentTime() float32 {
	return float32(getTimeCount()) / float32(frequency) * 1000
}

// normalize - Returns the normalized values of a vector
func normalize(vector *rl.Vector2) {
	aux := *vector
	length := float32(math.Sqrt(float64(aux.X*aux.X + aux.Y*aux.Y)))
	if length == 0 {
		length = 1.0
	}
	ilength := 1.0 / length
	vector.X *= ilength
	vector.Y *= ilength
}

func getNextIndex(i, count int) int {
	if i+1 < count {
		return i + 1
	}
	return 0
}

func safeDiv(a, b float32) float32 {
	if b == 0 {
		return 0
	}
	return a / b
}
