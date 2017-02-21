// Package physics - 2D Physics library for videogames
//
// A port of Victor Fisac's physac engine (https://github.com/raysan5/raylib/blob/master/src/physac.h)
package physics

import (
	"github.com/gen2brain/raylib-go/raylib"
)

// Defines
const (
	PhysacMaxBodies      = 64
	PhysacMaxManifolds   = 4096
	PhysacMaxVertices    = 24
	PhysacCircleVertices = 24

	PhysacDesiredDeltatime      = 1.0 / 60.0
	PhysacMaxTimestep           = 0.02
	PhysacCollisionIterations   = 100
	PhysacPenetrationAllowance  = 0.05
	PhysacPenetrationCorrection = 0.4
)

// Physics shape type
const (
	PhysicsCircle = iota
	PhysicsPolygon
)

// Mat2 type (used for polygon shape rotation matrix)
type Mat2 struct {
	M00 float32
	M01 float32
	M10 float32
	M11 float32
}

// PolygonData type
type PolygonData struct {
	// Current used vertex and normals count
	VertexCount uint32
	// Polygon vertex positions vectors
	Vertices [24]raylib.Vector2
	// Polygon vertex normals vectors
	Normals [24]raylib.Vector2
	// Vertices transform matrix 2x2
	Transform Mat2
}

// PhysicsShape type
type PhysicsShape struct {
	// Physics shape type (circle or polygon)
	Type uint32
	// Padding
	_ [4]byte
	// Shape physics body reference
	Body *PhysicsBodyData
	// Circle shape radius (used for circle shapes)
	Radius float32
	// Polygon shape vertices position and normals data (just used for polygon shapes)
	VertexData PolygonData
}

// PhysicsBodyData type
type PhysicsBodyData struct {
	// Reference unique identifier
	Id uint32
	// Enabled dynamics state (collisions are calculated anyway)
	Enabled uint32
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
	UseGravity uint32
	// Physics grounded on other body state
	IsGrounded uint32
	// Physics rotation constraint
	FreezeOrient uint32
	// Padding
	_ [4]byte
	// Physics body shape information (type, radius, vertices, normals)
	Shape PhysicsShape
}

// PhysicsManifoldData type
type PhysicsManifoldData struct {
	// Reference unique identifier
	Id uint32
	// Paddin
	_ [4]byte
	// Manifold first physics body reference
	BodyA *PhysicsBodyData
	// Manifold second physics body reference
	BodyB *PhysicsBodyData
	// Depth of penetration from collision
	Penetration float32
	// Normal direction vector from 'a' to 'b'
	Normal raylib.Vector2
	// Points of contact during collision
	Contacts [2]raylib.Vector2
	// Current collision number of contacts
	ContactsCount uint32
	// Mixed restitution during collision
	Restitution float32
	// Mixed dynamic friction during collision
	DynamicFriction float32
	// Mixed static friction during collision
	StaticFriction float32
	// Padding
	_ [4]byte
}
