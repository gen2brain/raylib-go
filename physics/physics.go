// Package physics - 2D Physics library for videogames
//
// A port of Victor Fisac's physac engine (https://github.com/raysan5/raylib/blob/master/src/physac.h)
package physics

import (
	"fmt"
	"math"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/raymath"
)

// Physics shape type
type ShapeType int

// Physics shape types
const (
	// Circle type
	Circle ShapeType = 1
	// Polygon type
	Polygon ShapeType = 2
)

// PolygonData type
type PolygonData struct {
	// Current used vertex and normals count
	VertexCount int
	// Polygon vertex positions vectors
	Vertices [maxVertices]raylib.Vector2
	// Polygon vertex normals vectors
	Normals [maxVertices]raylib.Vector2
	// Vertices transform matrix 2x2
	Transform raylib.Mat2
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
	VertexData *PolygonData
}

// Body type
type Body struct {
	// Enabled dynamics state (collisions are calculated anyway)
	Enabled bool
	// Physics body shape pivot
	Position raylib.Vector2
	// Current linear velocity applied to position
	Velocity raylib.Vector2
	// Current linear force (reset to 0 every step)
	Force raylib.Vector2
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
	Shape *Shape
}

// Manifold type
type Manifold struct {
	// Manifold first physics body reference
	BodyA *Body
	// Manifold second physics body reference
	BodyB *Body
	// Depth of penetration from collision
	Penetration float32
	// Normal direction vector from 'a' to 'b'
	Normal raylib.Vector2
	// Points of contact during collision
	Contacts [2]raylib.Vector2
	// Current collision number of contacts
	ContactsCount int
	// Mixed restitution during collision
	Restitution float32
	// Mixed dynamic friction during collision
	DynamicFriction float32
	// Mixed static friction during collision
	StaticFriction float32
}

// Defines
const (
	maxVertices    = 24
	circleVertices = 24

	collisionIterations   = 100
	penetrationAllowance  = 0.05
	penetrationCorrection = 0.4

	fltMax  = 3.402823466 + 38
	epsilon = 0.000001
)

var (
	// Physics world gravity force
	gravityForce raylib.Vector2

	// Physics bodies pointers array
	bodies []*Body

	// Physics world current bodies counter
	bodiesCount int

	// Physics manifolds pointers array
	manifolds []*Manifold

	// Physics world current manifolds counter
	manifoldsCount int

	// Delta time used for physics steps
	deltaTime float32
)

// Init - initializes physics values
func Init() {
	gravityForce = raylib.NewVector2(0, 9.81/1000)

	bodies = make([]*Body, 0)

	manifolds = make([]*Manifold, 0)
	manifoldsCount = 0
}

// SetGravity - Sets physics global gravity force
func SetGravity(x, y float32) {
	gravityForce.X = x
	gravityForce.Y = y
}

// CreateBodyCircle - Creates a new circle physics body with generic parameters
func CreateBodyCircle(pos raylib.Vector2, radius, density float32) *Body {
	bodyCircle := CreateBodyPolygon(pos, radius, circleVertices, density)
	//bodyCircle.Shape.Type = Circle
	return bodyCircle
}

// CreateBodyRectangle - Creates a new rectangle physics body with generic parameters
func CreateBodyRectangle(pos raylib.Vector2, width, height, density float32) *Body {
	newBody := &Body{}
	// Initialize new body with generic values
	newBody.Enabled = true
	newBody.Position = pos
	newBody.Velocity = raylib.Vector2{}
	newBody.Force = raylib.Vector2{}
	newBody.AngularVelocity = 0
	newBody.Torque = 0
	newBody.Orient = 0

	newBody.Shape = &Shape{}
	newBody.Shape.Type = Polygon
	newBody.Shape.Body = newBody
	newBody.Shape.VertexData = createRectanglePolygon(pos, raylib.NewVector2(width, height))

	// Calculate centroid and moment of inertia
	center := raylib.Vector2{}
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

		D := raymath.Vector2CrossProduct(p1, p2)
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

	if newBody.Mass != 0.0 {
		newBody.InverseMass = 1.0 / newBody.Mass
	}

	if newBody.Inertia != 0.0 {
		newBody.InverseInertia = 1.0 / newBody.Inertia
	}

	// Add new body to bodies pointers array and update bodies count
	bodies = append(bodies, newBody)
	bodiesCount++

	return newBody
}

// CreateBodyPolygon - Creates a new polygon physics body with generic parameters
func CreateBodyPolygon(pos raylib.Vector2, radius float32, sides int32, density float32) *Body {
	newBody := &Body{}
	// Initialize new body with generic values
	newBody.Enabled = true
	newBody.Position = pos
	newBody.Velocity = raylib.Vector2{}
	newBody.Force = raylib.Vector2{}
	newBody.AngularVelocity = 0
	newBody.Torque = 0
	newBody.Orient = 0

	newBody.Shape = &Shape{}
	newBody.Shape.Type = Polygon
	newBody.Shape.Body = newBody
	newBody.Shape.VertexData = createRandomPolygon(radius, sides)

	// Calculate centroid and moment of inertia
	center := raylib.Vector2{}
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

		cross := raymath.Vector2CrossProduct(position1, position2)
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

	if newBody.Mass != 0.0 {
		newBody.InverseMass = 1.0 / newBody.Mass
	}

	if newBody.Inertia != 0.0 {
		newBody.InverseInertia = 1.0 / newBody.Inertia
	}

	// Add new body to bodies pointers array and update bodies count
	bodies = append(bodies, newBody)
	bodiesCount++

	return newBody
}

// AddForce - Adds a force to a physics body
func AddForce(body *Body, force raylib.Vector2) {
	body.Force = raymath.Vector2Add(body.Force, force)
}

// AddTorque - Adds an angular force to a physics body
func AddTorque(body *Body, amount float32) {
	body.Torque += amount
}

// Shatter - Shatters a polygon shape physics body to little physics bodies with explosion force
func Shatter(body *Body, position raylib.Vector2, force float32) {
}

// GetBodiesCount - Returns the current amount of created physics bodies
func GetBodiesCount() int {
	return bodiesCount
}

// GetBody - Returns a physics body of the bodies pool at a specific index
func GetBody(index int) *Body {
	body := &Body{}

	if index < bodiesCount {
		body = bodies[index]
	} else {
		raylib.TraceLog(raylib.LogDebug, "[PHYSAC] physics body index is out of bounds")
	}

	return body
}

// GetShapeType - Returns the physics body shape type (Circle or Polygon)
func GetShapeType(index int) ShapeType {
	var result ShapeType

	if index < bodiesCount {
		result = bodies[index].Shape.Type
	} else {
		raylib.TraceLog(raylib.LogDebug, "[PHYSAC] physics body index is out of bounds")
	}

	return result
}

// GetShapeVerticesCount - Returns the amount of vertices of a physics body shape
func GetShapeVerticesCount(index int) int {
	result := 0

	if index < bodiesCount {
		switch bodies[index].Shape.Type {
		case Circle:
			result = circleVertices
			break
		case Polygon:
			result = bodies[index].Shape.VertexData.VertexCount
			break
		}
	} else {
		raylib.TraceLog(raylib.LogDebug, "[PHYSAC] physics body index is out of bounds")
	}

	return result
}

// GetShapeVertex - Returns transformed position of a body shape (body position + vertex transformed position)
func GetShapeVertex(body *Body, vertex int) raylib.Vector2 {
	position := raylib.Vector2{}

	switch body.Shape.Type {
	case Circle:
		position.X = body.Position.X + float32(math.Cos(360/float64(circleVertices)*float64(vertex)*raylib.Deg2rad))*body.Shape.Radius
		position.Y = body.Position.Y + float32(math.Sin(360/float64(circleVertices)*float64(vertex)*raylib.Deg2rad))*body.Shape.Radius
		break
	case Polygon:
		vertexData := body.Shape.VertexData
		position = raymath.Vector2Add(body.Position, raymath.Mat2MultiplyVector2(vertexData.Transform, vertexData.Vertices[vertex]))
		break
	default:
		break
	}

	return position
}

// SetBodyRotation - Sets physics body shape transform based on radians parameter
func SetBodyRotation(body *Body, radians float32) {
	body.Orient = radians

	if body.Shape.Type == Polygon {
		body.Shape.VertexData.Transform = raymath.Mat2Radians(radians)
	}
}

// DestroyBody - Unitializes and destroy a physics body
func DestroyBody(index int) {
	// Free body allocated memory
	copy(bodies[index:], bodies[index+1:])
	bodies[len(bodies)-1] = &Body{}
	bodies = bodies[:len(bodies)-1]

	// Update physics bodies count
	bodiesCount--
}

// Step - Physics steps calculations (dynamics, collisions and position corrections)
func Step(delta float32) {
	deltaTime = delta

	// Clear previous generated collisions information
	for i := manifoldsCount - 1; i >= 0; i-- {
		destroyManifold(i)
	}

	// Reset physics bodies grounded state
	for i := 0; i < bodiesCount; i++ {
		bodies[i].IsGrounded = false
	}

	// Generate new collision information
	for i := 0; i < bodiesCount; i++ {
		bodyA := bodies[i]

		fmt.Printf("bodyA.Shape %+v\n", bodyA.Shape)

		for j := i + 1; j < bodiesCount; j++ {
			bodyB := bodies[j]

			if (bodyA.InverseMass == 0) && (bodyB.InverseMass == 0) {
				continue
			}

			fmt.Printf("bodyB.Shape %+v\n", bodyB.Shape)

			manifold := &Manifold{}
			if bodyA.Shape.Type == Polygon && bodyB.Shape.Type == Circle {
				manifold = createManifold(bodyB, bodyA)
			} else {
				manifold = createManifold(bodyA, bodyB)
			}

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
		integrateForces(bodies[i])
	}

	// Initialize physics manifolds to solve collisions
	for i := 0; i < manifoldsCount; i++ {
		initializeManifolds(manifolds[i])
	}

	// Integrate physics collisions impulses to solve collisions
	for i := 0; i < collisionIterations; i++ {
		for j := 0; j < manifoldsCount; j++ {
			integrateImpulses(manifolds[j])
		}
	}

	// Integrate velocity to physics bodies
	for i := 0; i < bodiesCount; i++ {
		integrateVelocity(bodies[i])
	}

	// Correct physics bodies positions based on manifolds collision information
	for i := 0; i < manifoldsCount; i++ {
		correctPositions(manifolds[i])
	}

	// Clear physics bodies forces
	for i := 0; i < bodiesCount; i++ {
		bodies[i].Force = raylib.Vector2{}
		bodies[i].Torque = 0
	}
}

// Reset - Destroys created physics bodies and manifolds and resets global values
func Reset() {
	bodies = make([]*Body, 0)
	bodiesCount = 0

	manifolds = make([]*Manifold, 0)
	manifoldsCount = 0

	raylib.TraceLog(raylib.LogDebug, "[PHYSAC] physics module reset successfully")
}

// Close - Unitializes physics pointers
func Close() {
	// Unitialize physics manifolds dynamic memory allocations
	for i := manifoldsCount - 1; i >= 0; i-- {
		destroyManifold(i)
	}

	// Unitialize physics bodies dynamic memory allocations
	for i := bodiesCount - 1; i >= 0; i-- {
		DestroyBody(i)
	}
}

// createRandomPolygon - Creates a random polygon shape with max vertex distance from polygon pivot
func createRandomPolygon(radius float32, sides int32) *PolygonData {
	data := &PolygonData{}
	data.VertexCount = int(sides)

	orient := raylib.GetRandomValue(0, 360)
	data.Transform = raymath.Mat2Radians(float32(orient) * raylib.Deg2rad)

	// Calculate polygon vertices positions
	for i := 0; i < data.VertexCount; i++ {
		data.Vertices[i].X = float32(math.Cos(360/float64(sides)*float64(i)*raylib.Deg2rad)) * radius
		data.Vertices[i].Y = float32(math.Sin(360/float64(sides)*float64(i)*raylib.Deg2rad)) * radius
	}

	// Calculate polygon faces normals
	for i := 0; i < data.VertexCount; i++ {
		nextIndex := 0
		if i+1 < int(sides) {
			nextIndex = i + 1
		}
		face := raymath.Vector2Subtract(data.Vertices[nextIndex], data.Vertices[i])

		data.Normals[i] = raylib.NewVector2(face.Y, -face.X)
		raymath.Vector2Normalize(&data.Normals[i])
	}

	return data
}

// createRectanglePolygon - Creates a rectangle polygon shape based on a min and max positions
func createRectanglePolygon(pos, size raylib.Vector2) *PolygonData {
	data := &PolygonData{}

	data.VertexCount = 4
	data.Transform = raymath.Mat2Radians(0)

	// Calculate polygon vertices positions
	data.Vertices[0] = raylib.NewVector2(pos.X+size.X/2, pos.Y-size.Y/2)
	data.Vertices[1] = raylib.NewVector2(pos.X+size.X/2, pos.Y+size.Y/2)
	data.Vertices[2] = raylib.NewVector2(pos.X-size.X/2, pos.Y+size.Y/2)
	data.Vertices[3] = raylib.NewVector2(pos.X-size.X/2, pos.Y-size.Y/2)

	// Calculate polygon faces normals
	for i := 0; i < data.VertexCount; i++ {
		nextIndex := 0
		if i+1 < data.VertexCount {
			nextIndex = i + 1
		}
		face := raymath.Vector2Subtract(data.Vertices[nextIndex], data.Vertices[i])

		data.Normals[i] = raylib.NewVector2(face.Y, -face.X)
		raymath.Vector2Normalize(&data.Normals[i])
	}

	return data
}

// createManifold - Creates a new physics manifold to solve collision
func createManifold(a, b *Body) *Manifold {
	newManifold := &Manifold{}

	// Initialize new manifold with generic values
	newManifold.BodyA = a
	newManifold.BodyB = b
	newManifold.Penetration = 0
	newManifold.Normal = raylib.Vector2{}
	newManifold.Contacts[0] = raylib.Vector2{}
	newManifold.Contacts[1] = raylib.Vector2{}
	newManifold.ContactsCount = 0
	newManifold.Restitution = 0
	newManifold.DynamicFriction = 0
	newManifold.StaticFriction = 0

	// Add new body to bodies pointers array and update bodies count
	manifolds = append(manifolds, newManifold)
	manifoldsCount++

	return newManifold
}

// destroyManifold - Unitializes and destroys a physics manifold
func destroyManifold(index int) {
	// Free manifold allocated memory
	copy(manifolds[index:], manifolds[index+1:])
	manifolds[len(manifolds)-1] = &Manifold{}
	manifolds = manifolds[:len(manifolds)-1]

	// Update physics manifolds count
	manifoldsCount--
}

// solveManifold - Solves a created physics manifold between two physics bodies
func solveManifold(manifold *Manifold) {
	fmt.Printf("%+v\n\n", manifold)
	switch manifold.BodyA.Shape.Type {
	case Circle:
		switch manifold.BodyB.Shape.Type {
		case Circle:
			solveCircleToCircle(manifold)
			break
		case Polygon:
			solveCircleToPolygon(manifold)
			break
		}
	case Polygon:
		switch manifold.BodyB.Shape.Type {
		case Circle:
			solvePolygonToCircle(manifold)
			break
		case Polygon:
			solvePolygonToPolygon(manifold)
			break
		}
	}
	fmt.Printf("%+v\n\n", manifold)
	fmt.Println()

	// Update physics body grounded state if normal direction is down and grounded state is not set yet in previous manifolds
	if !manifold.BodyB.IsGrounded {
		manifold.BodyB.IsGrounded = (manifold.Normal.Y < 0)
	}
}

// solveCircleToCircle - Solves collision between two circle shape physics bodies
func solveCircleToCircle(manifold *Manifold) {
	fmt.Println("solveCircleToCircle")
	bodyA := manifold.BodyA
	bodyB := manifold.BodyB

	// Calculate translational vector, which is normal
	normal := raymath.Vector2Subtract(bodyB.Position, bodyA.Position)

	distSqr := raymath.Vector2LenSqr(normal)
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
		manifold.Normal = raylib.NewVector2(1, 0)
		manifold.Contacts[0] = bodyA.Position
	} else {
		manifold.Penetration = radius - distance
		manifold.Normal = raylib.NewVector2(normal.X/distance, normal.Y/distance) // Faster than using MathNormalize() due to sqrt is already performed
		manifold.Contacts[0] = raylib.NewVector2(manifold.Normal.X*bodyA.Shape.Radius+bodyA.Position.X, manifold.Normal.Y*bodyA.Shape.Radius+bodyA.Position.Y)
	}

	// Update physics body grounded state if normal direction is down
	if !bodyA.IsGrounded {
		bodyA.IsGrounded = (manifold.Normal.Y < 0)
	}
}

// solveCircleToPolygon - Solves collision between a circle to a polygon shape physics bodies
func solveCircleToPolygon(manifold *Manifold) {
	fmt.Println("solveCircleToPolygon")
	bodyA := manifold.BodyA
	bodyB := manifold.BodyB

	manifold.ContactsCount = 0

	// Transform circle center to polygon transform space
	center := bodyA.Position
	center = raymath.Mat2MultiplyVector2(raymath.Mat2Transpose(bodyB.Shape.VertexData.Transform), raymath.Vector2Subtract(center, bodyB.Position))

	// Find edge with minimum penetration
	// It is the same concept as using support points in solvePolygonToPolygon
	separation := float32(-fltMax)
	faceNormal := 0
	vertexData := bodyB.Shape.VertexData

	for i := 0; i < vertexData.VertexCount; i++ {
		currentSeparation := raymath.Vector2DotProduct(vertexData.Normals[i], raymath.Vector2Subtract(center, vertexData.Vertices[i]))

		if currentSeparation > bodyA.Shape.Radius {
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
		manifold.ContactsCount = 1
		normal := raymath.Mat2MultiplyVector2(vertexData.Transform, vertexData.Normals[faceNormal])
		manifold.Normal = raylib.NewVector2(-normal.X, -normal.Y)
		manifold.Contacts[0] = raylib.NewVector2(manifold.Normal.X*bodyA.Shape.Radius+bodyA.Position.X, manifold.Normal.Y*bodyA.Shape.Radius+bodyA.Position.Y)
		manifold.Penetration = bodyA.Shape.Radius
		return
	}

	// Determine which voronoi region of the edge center of circle lies within
	dot1 := raymath.Vector2DotProduct(raymath.Vector2Subtract(center, v1), raymath.Vector2Subtract(v2, v1))
	dot2 := raymath.Vector2DotProduct(raymath.Vector2Subtract(center, v2), raymath.Vector2Subtract(v1, v2))
	manifold.Penetration = bodyA.Shape.Radius - separation

	if dot1 <= 0 { // Closest to v1
		if raymath.Vector2Distance(center, v1) > bodyA.Shape.Radius*bodyA.Shape.Radius {
			return
		}

		manifold.ContactsCount = 1
		normal := raymath.Vector2Subtract(v1, center)
		normal = raymath.Mat2MultiplyVector2(vertexData.Transform, normal)
		raymath.Vector2Normalize(&normal)
		manifold.Normal = normal
		v1 = raymath.Mat2MultiplyVector2(vertexData.Transform, v1)
		v1 = raymath.Vector2Add(v1, bodyB.Position)
		manifold.Contacts[0] = v1
	} else if dot2 <= 0 { // Closest to v2
		if raymath.Vector2Distance(center, v2) > bodyA.Shape.Radius*bodyA.Shape.Radius {
			return
		}

		manifold.ContactsCount = 1
		normal := raymath.Vector2Subtract(v2, center)
		v2 = raymath.Mat2MultiplyVector2(vertexData.Transform, v2)
		v2 = raymath.Vector2Add(v2, bodyB.Position)
		manifold.Contacts[0] = v2
		normal = raymath.Mat2MultiplyVector2(vertexData.Transform, normal)
		raymath.Vector2Normalize(&normal)
		manifold.Normal = normal
	} else { // Closest to face
		normal := vertexData.Normals[faceNormal]

		if raymath.Vector2DotProduct(raymath.Vector2Subtract(center, v1), normal) > bodyA.Shape.Radius {
			return
		}

		normal = raymath.Mat2MultiplyVector2(vertexData.Transform, normal)
		manifold.Normal = raylib.NewVector2(-normal.X, -normal.Y)
		manifold.Contacts[0] = raylib.NewVector2(manifold.Normal.X*bodyA.Shape.Radius+bodyA.Position.X, manifold.Normal.Y*bodyA.Shape.Radius+bodyA.Position.Y)
		manifold.ContactsCount = 1
	}
}

// solvePolygonToCircle - Solves collision between a polygon to a circle shape physics bodies
func solvePolygonToCircle(manifold *Manifold) {
	fmt.Println("solvePolygonToCircle")
	bodyA := manifold.BodyA
	bodyB := manifold.BodyB

	manifold.BodyA = bodyB
	manifold.BodyB = bodyA

	solveCircleToPolygon(manifold)

	manifold.Normal.X *= -1
	manifold.Normal.Y *= -1
}

// solvePolygonToPolygon - Solves collision between two polygons shape physics bodies
func solvePolygonToPolygon(manifold *Manifold) {
	//fmt.Println("solvePolygonToPolygon")
	shapeA := manifold.BodyA.Shape
	shapeB := manifold.BodyB.Shape

	manifold.ContactsCount = 0

	// Check for separating axis with A shape's face planes
	faceA := 0
	penetrationA := findAxisLeastPenetration(&faceA, shapeA, shapeB)
	if penetrationA >= 0 {
		return
	}

	// Check for separating axis with B shape's face planes
	faceB := 0
	penetrationB := findAxisLeastPenetration(&faceB, shapeB, shapeA)
	if penetrationB >= 0 {
		return
	}

	referenceIndex := 0
	flip := false // Always point from A shape to B shape

	refPoly := &Shape{} // Reference
	incPoly := &Shape{} // Incident

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
	//incidentFace := make([]raylib.Vector2, 2)
	incidentFace0 := raylib.Vector2{}
	incidentFace1 := raylib.Vector2{}
	findIncidentFace(&incidentFace0, &incidentFace1, refPoly, incPoly, referenceIndex)

	// Setup reference face vertices
	refData := refPoly.VertexData
	v1 := refData.Vertices[referenceIndex]
	if referenceIndex+1 < refData.VertexCount {
		referenceIndex = referenceIndex + 1
	}
	v2 := refData.Vertices[referenceIndex]

	// Transform vertices to world space
	v1 = raymath.Mat2MultiplyVector2(refData.Transform, v1)
	v1 = raymath.Vector2Add(v1, refPoly.Body.Position)
	v2 = raymath.Mat2MultiplyVector2(refData.Transform, v2)
	v2 = raymath.Vector2Add(v2, refPoly.Body.Position)

	// Calculate reference face side normal in world space
	sidePlaneNormal := raymath.Vector2Subtract(v2, v1)
	raymath.Vector2Normalize(&sidePlaneNormal)

	// Orthogonalize
	refFaceNormal := raylib.NewVector2(sidePlaneNormal.Y, -sidePlaneNormal.X)
	refC := raymath.Vector2DotProduct(refFaceNormal, v1)
	negSide := raymath.Vector2DotProduct(sidePlaneNormal, v1) * -1
	posSide := raymath.Vector2DotProduct(sidePlaneNormal, v2)

	// clip incident face to reference face side planes (due to floating point error, possible to not have required points
	if clip(raylib.NewVector2(-sidePlaneNormal.X, -sidePlaneNormal.Y), negSide, &incidentFace0, &incidentFace1) < 2 {
		return
	}
	if clip(sidePlaneNormal, posSide, &incidentFace0, &incidentFace1) < 2 {
		return
	}

	// Flip normal if required
	if flip {
		manifold.Normal = raylib.NewVector2(-refFaceNormal.X, -refFaceNormal.Y)
	} else {
		manifold.Normal = refFaceNormal
	}

	// Keep points behind reference face
	currentPoint := 0 // clipped points behind reference face
	separation := raymath.Vector2DotProduct(refFaceNormal, incidentFace0) - refC
	if separation <= 0 {
		manifold.Contacts[currentPoint] = incidentFace0
		manifold.Penetration = -separation
		currentPoint++
	} else {
		manifold.Penetration = 0
	}

	separation = raymath.Vector2DotProduct(refFaceNormal, incidentFace1) - refC

	if separation <= 0 {
		manifold.Contacts[currentPoint] = incidentFace1
		manifold.Penetration += -separation
		currentPoint++

		// Calculate total penetration average
		manifold.Penetration /= float32(currentPoint)
	}

	manifold.ContactsCount = currentPoint
}

// integrateForces -  Integrates physics forces into velocity
func integrateForces(body *Body) {
	if body.InverseMass == 0 || !body.Enabled {
		return
	}

	body.Velocity.X += (body.Force.X * body.InverseMass) * (deltaTime / 2)
	body.Velocity.Y += (body.Force.Y * body.InverseMass) * (deltaTime / 2)

	if body.UseGravity {
		body.Velocity.X += gravityForce.X * (deltaTime / 2)
		body.Velocity.Y += gravityForce.Y * (deltaTime / 2)
	}

	if !body.FreezeOrient {
		body.AngularVelocity += body.Torque * body.InverseInertia * (deltaTime / 2)
	}
}

// initializeManifolds - Initializes physics manifolds to solve collisions
func initializeManifolds(manifold *Manifold) {
	bodyA := manifold.BodyA
	bodyB := manifold.BodyB

	// Calculate average restitution, static and dynamic friction
	manifold.Restitution = float32(math.Sqrt(float64(bodyA.Restitution * bodyB.Restitution)))
	manifold.StaticFriction = float32(math.Sqrt(float64(bodyA.StaticFriction * bodyB.StaticFriction)))
	manifold.DynamicFriction = float32(math.Sqrt(float64(bodyA.DynamicFriction * bodyB.DynamicFriction)))

	for i := 0; i < 2; i++ {
		// Caculate radius from center of mass to contact
		radiusA := raymath.Vector2Subtract(manifold.Contacts[i], bodyA.Position)
		radiusB := raymath.Vector2Subtract(manifold.Contacts[i], bodyB.Position)

		crossA := raymath.Vector2Cross(bodyA.AngularVelocity, radiusA)
		crossB := raymath.Vector2Cross(bodyB.AngularVelocity, radiusB)

		radiusV := raylib.Vector2{}
		radiusV.X = bodyB.Velocity.X + crossB.X - bodyA.Velocity.X - crossA.X
		radiusV.Y = bodyB.Velocity.Y + crossB.Y - bodyA.Velocity.Y - crossA.Y

		// Determine if we should perform a resting collision or not;
		// The idea is if the only thing moving this object is gravity, then the collision should be performed without any restitution
		if raymath.Vector2LenSqr(radiusV) < (raymath.Vector2LenSqr(raylib.NewVector2(gravityForce.X*deltaTime, gravityForce.Y*deltaTime)) + epsilon) {
			manifold.Restitution = 0
		}
	}
}

// integrateImpulses - Integrates physics collisions impulses to solve collisions
func integrateImpulses(manifold *Manifold) {
	bodyA := manifold.BodyA
	bodyB := manifold.BodyB

	// Early out and positional correct if both objects have infinite mass
	if float32(math.Abs(float64(bodyA.InverseMass+bodyB.InverseMass))) <= epsilon {
		bodyA.Velocity = raylib.Vector2{}
		bodyB.Velocity = raylib.Vector2{}
		return
	}

	for i := 0; i < manifold.ContactsCount; i++ {
		// Calculate radius from center of mass to contact
		radiusA := raymath.Vector2Subtract(manifold.Contacts[i], bodyA.Position)
		radiusB := raymath.Vector2Subtract(manifold.Contacts[i], bodyB.Position)

		// Calculate relative velocity
		radiusV := raylib.Vector2{}
		radiusV.X = bodyB.Velocity.X + raymath.Vector2Cross(bodyB.AngularVelocity, radiusB).X - bodyA.Velocity.X - raymath.Vector2Cross(bodyA.AngularVelocity, radiusA).X
		radiusV.Y = bodyB.Velocity.Y + raymath.Vector2Cross(bodyB.AngularVelocity, radiusB).Y - bodyA.Velocity.Y - raymath.Vector2Cross(bodyA.AngularVelocity, radiusA).Y

		// Relative velocity along the normal
		contactVelocity := raymath.Vector2DotProduct(radiusV, manifold.Normal)

		// Do not resolve if velocities are separating
		if contactVelocity > 0 {
			return
		}

		raCrossN := raymath.Vector2CrossProduct(radiusA, manifold.Normal)
		rbCrossN := raymath.Vector2CrossProduct(radiusB, manifold.Normal)

		inverseMassSum := bodyA.InverseMass + bodyB.InverseMass + (raCrossN*raCrossN)*bodyA.InverseInertia + (rbCrossN*rbCrossN)*bodyB.InverseInertia

		// Calculate impulse scalar value
		impulse := -(1.0 + manifold.Restitution) * contactVelocity
		impulse /= inverseMassSum
		impulse /= float32(manifold.ContactsCount)

		// Apply impulse to each physics body
		impulseV := raylib.NewVector2(manifold.Normal.X*impulse, manifold.Normal.Y*impulse)

		if bodyA.Enabled {
			bodyA.Velocity.X += bodyA.InverseMass * (-impulseV.X)
			bodyA.Velocity.Y += bodyA.InverseMass * (-impulseV.Y)
			if !bodyA.FreezeOrient {
				bodyA.AngularVelocity += bodyA.InverseInertia * raymath.Vector2CrossProduct(radiusA, raylib.NewVector2(-impulseV.X, -impulseV.Y))
			}
		}

		if bodyB.Enabled {
			bodyB.Velocity.X += bodyB.InverseMass * (impulseV.X)
			bodyB.Velocity.Y += bodyB.InverseMass * (impulseV.Y)
			if !bodyB.FreezeOrient {
				bodyB.AngularVelocity += bodyB.InverseInertia * raymath.Vector2CrossProduct(radiusB, impulseV)
			}
		}

		// Apply friction impulse to each physics body
		radiusV.X = bodyB.Velocity.X + raymath.Vector2Cross(bodyB.AngularVelocity, radiusB).X - bodyA.Velocity.X - raymath.Vector2Cross(bodyA.AngularVelocity, radiusA).X
		radiusV.Y = bodyB.Velocity.Y + raymath.Vector2Cross(bodyB.AngularVelocity, radiusB).Y - bodyA.Velocity.Y - raymath.Vector2Cross(bodyA.AngularVelocity, radiusA).Y

		tangent := raylib.NewVector2(radiusV.X-(manifold.Normal.X*raymath.Vector2DotProduct(radiusV, manifold.Normal)), radiusV.Y-(manifold.Normal.Y*raymath.Vector2DotProduct(radiusV, manifold.Normal)))
		raymath.Vector2Normalize(&tangent)

		// Calculate impulse tangent magnitude
		impulseTangent := -raymath.Vector2DotProduct(radiusV, tangent)
		impulseTangent /= inverseMassSum
		impulseTangent /= float32(manifold.ContactsCount)

		absImpulseTangent := float32(math.Abs(float64(impulseTangent)))

		// Don't apply tiny friction impulses
		if absImpulseTangent <= epsilon {
			return
		}

		// Apply coulumb's law
		tangentImpulse := raylib.Vector2{}
		if absImpulseTangent < impulse*manifold.StaticFriction {
			tangentImpulse = raylib.NewVector2(tangent.X*impulseTangent, tangent.Y*impulseTangent)
		} else {
			tangentImpulse = raylib.NewVector2(tangent.X*-impulse*manifold.DynamicFriction, tangent.Y*-impulse*manifold.DynamicFriction)
		}

		// Apply friction impulse
		if bodyA.Enabled {
			bodyA.Velocity.X += bodyA.InverseMass * (-tangentImpulse.X)
			bodyA.Velocity.Y += bodyA.InverseMass * (-tangentImpulse.Y)

			if !bodyA.FreezeOrient {
				bodyA.AngularVelocity += bodyA.InverseInertia * raymath.Vector2CrossProduct(radiusA, raylib.NewVector2(-tangentImpulse.X, -tangentImpulse.Y))
			}
		}

		if bodyB.Enabled {
			bodyB.Velocity.X += bodyB.InverseMass * (tangentImpulse.X)
			bodyB.Velocity.Y += bodyB.InverseMass * (tangentImpulse.Y)

			if !bodyB.FreezeOrient {
				bodyB.AngularVelocity += bodyB.InverseInertia * raymath.Vector2CrossProduct(radiusB, tangentImpulse)
			}
		}
	}
}

// integrateVelocity - Integrates physics velocity into position and forces
func integrateVelocity(body *Body) {
	if !body.Enabled {
		return
	}

	body.Position.X += body.Velocity.X * deltaTime
	body.Position.Y += body.Velocity.Y * deltaTime

	if !body.FreezeOrient {
		body.Orient += body.AngularVelocity * deltaTime
	}

	raymath.Mat2Set(&body.Shape.VertexData.Transform, body.Orient)

	integrateForces(body)
}

// correctPositions - Corrects physics bodies positions based on manifolds collision information
func correctPositions(manifold *Manifold) {
	bodyA := manifold.BodyA
	bodyB := manifold.BodyB

	correction := raylib.Vector2{}
	correction.X = float32(math.Max(float64(manifold.Penetration-penetrationAllowance), 0)) / (bodyA.InverseMass + bodyB.InverseMass) * manifold.Normal.X * penetrationCorrection
	correction.Y = float32(math.Max(float64(manifold.Penetration-penetrationAllowance), 0)) / (bodyA.InverseMass + bodyB.InverseMass) * manifold.Normal.Y * penetrationCorrection

	if bodyA.Enabled {
		bodyA.Position.X -= correction.X * bodyA.InverseMass
		bodyA.Position.Y -= correction.Y * bodyA.InverseMass
	}

	if bodyB.Enabled {
		bodyB.Position.X += correction.X * bodyB.InverseMass
		bodyB.Position.Y += correction.Y * bodyB.InverseMass
	}
}

// Returns the extreme point along a direction within a polygon
func getSupport(shape *Shape, dir raylib.Vector2) raylib.Vector2 {
	bestProjection := float32(-fltMax)
	bestVertex := raylib.Vector2{}
	data := shape.VertexData

	for i := 0; i < data.VertexCount; i++ {
		vertex := data.Vertices[i]
		projection := raymath.Vector2DotProduct(vertex, dir)

		if projection > bestProjection {
			bestVertex = vertex
			bestProjection = projection
		}
	}

	return bestVertex
}

// findAxisLeastPenetration - Finds polygon shapes axis least penetration
func findAxisLeastPenetration(faceIndex *int, shapeA, shapeB *Shape) float32 {
	bestDistance := float32(-fltMax)
	bestIndex := 0

	dataA := shapeA.VertexData
	dataB := shapeB.VertexData

	for i := 0; i < dataA.VertexCount; i++ {
		// Retrieve a face normal from A shape
		normal := dataA.Normals[i]
		transNormal := raymath.Mat2MultiplyVector2(dataA.Transform, normal)

		// Transform face normal into B shape's model space
		buT := raymath.Mat2Transpose(dataB.Transform)
		normal = raymath.Mat2MultiplyVector2(buT, transNormal)

		// Retrieve support point from B shape along -n
		support := getSupport(shapeB, raylib.NewVector2(-normal.X, -normal.Y))

		// Retrieve vertex on face from A shape, transform into B shape's model space
		vertex := dataA.Vertices[i]
		vertex = raymath.Mat2MultiplyVector2(dataA.Transform, vertex)
		vertex = raymath.Vector2Add(vertex, shapeA.Body.Position)
		vertex = raymath.Vector2Subtract(vertex, shapeB.Body.Position)
		vertex = raymath.Mat2MultiplyVector2(buT, vertex)

		// Compute penetration distance in B shape's model space
		distance := raymath.Vector2DotProduct(normal, raymath.Vector2Subtract(support, vertex))

		// Store greatest distance
		if distance > bestDistance {
			bestDistance = distance
			bestIndex = i
		}
	}

	*faceIndex = bestIndex

	return bestDistance
}

// findIncidentFace - Finds two polygon shapes incident face
func findIncidentFace(v0, v1 *raylib.Vector2, ref, inc *Shape, index int) {
	refData := ref.VertexData
	incData := inc.VertexData

	referenceNormal := refData.Normals[index]

	// Calculate normal in incident's frame of reference
	referenceNormal = raymath.Mat2MultiplyVector2(refData.Transform, referenceNormal)                        // To world space
	referenceNormal = raymath.Mat2MultiplyVector2(raymath.Mat2Transpose(incData.Transform), referenceNormal) // To incident's model space

	// Find most anti-normal face on polygon
	incidentFace := 0
	minDot := float32(fltMax)

	for i := 0; i < incData.VertexCount; i++ {
		dot := raymath.Vector2DotProduct(referenceNormal, incData.Normals[i])

		if dot < minDot {
			minDot = dot
			incidentFace = i
		}
	}

	// Assign face vertices for incident face
	*v0 = raymath.Mat2MultiplyVector2(incData.Transform, incData.Vertices[incidentFace])
	*v0 = raymath.Vector2Add(*v0, inc.Body.Position)

	if incidentFace+1 < incData.VertexCount {
		incidentFace = incidentFace + 1
	}

	*v1 = raymath.Mat2MultiplyVector2(incData.Transform, incData.Vertices[incidentFace])
	*v1 = raymath.Vector2Add(*v1, inc.Body.Position)
}

// clip - Calculates clipping based on a normal and two faces
func clip(normal raylib.Vector2, clip float32, faceA, faceB *raylib.Vector2) int {
	sp := 0

	out := make([]raylib.Vector2, 2)
	out[0] = *faceA
	out[1] = *faceB

	// Retrieve distances from each endpoint to the line
	distanceA := raymath.Vector2DotProduct(normal, *faceA) - clip
	distanceB := raymath.Vector2DotProduct(normal, *faceB) - clip

	// If negative (behind plane)
	if distanceA <= 0 {
		spp := sp + 1
		out[spp] = *faceA
	}
	if distanceB <= 0 {
		spp := sp + 1
		out[spp] = *faceB
	}

	// If the points are on different sides of the plane
	if distanceA*distanceB < 0 {
		// Push intersection point
		alpha := distanceA / (distanceA - distanceB)
		out[sp] = *faceA
		delta := raymath.Vector2Subtract(*faceB, *faceA)
		delta.X *= alpha
		delta.Y *= alpha
		out[sp] = raymath.Vector2Add(out[sp], delta)
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
func triangleBarycenter(v1, v2, v3 raylib.Vector2) raylib.Vector2 {
	result := raylib.Vector2{}

	result.X = (v1.X + v2.X + v3.X) / 3
	result.Y = (v1.Y + v2.Y + v3.Y) / 3

	return result
}
