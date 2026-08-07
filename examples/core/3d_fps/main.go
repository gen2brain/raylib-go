package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Movement constants
const (
	gravity      float32 = 32.0
	maxSpeed     float32 = 20.0
	crouchSpeed  float32 = 5.0
	jumpForce    float32 = 12.0
	maxAccel     float32 = 150.0
	friction     float32 = 0.86
	airDrag      float32 = 0.98
	control      float32 = 15.0
	crouchHeight float32 = 0.0
	standHeight  float32 = 1.0
	bottomHeight float32 = 0.5
)

// Body structure
type Body struct {
	Position   rl.Vector3
	Velocity   rl.Vector3
	Dir        rl.Vector3
	IsGrounded bool
}

// Global Variables
var (
	sensitivity  = rl.NewVector2(0.001, 0.001)
	player       = Body{}
	lookRotation = rl.NewVector2(0, 0)
	headTimer    = float32(0.0)
	walkLerp     = float32(0.0)
	headLerp     = standHeight
	lean         = rl.NewVector2(0, 0)
)

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - 3d camera fps")
	defer rl.CloseWindow()

	// Initialize camera variables
	camera := rl.Camera3D{}
	camera.Fovy = 60.0
	camera.Projection = rl.CameraPerspective
	camera.Position = rl.NewVector3(
		player.Position.X,
		player.Position.Y+(bottomHeight+headLerp),
		player.Position.Z,
	)

	updateCameraFPS(&camera) // Update camera parameters

	rl.DisableCursor() // Limit cursor to relative movement inside the window
	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		delta := rl.GetFrameTime()

		// 1. Mouse look
		mouseDelta := rl.GetMouseDelta()
		lookRotation.X -= mouseDelta.X * sensitivity.X
		lookRotation.Y += mouseDelta.Y * sensitivity.Y

		// 2. Arrow keys look
		const keyLookSpeed float32 = 2.0
		if rl.IsKeyDown(rl.KeyRight) {
			lookRotation.X -= keyLookSpeed * delta
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			lookRotation.X += keyLookSpeed * delta
		}
		if rl.IsKeyDown(rl.KeyUp) {
			lookRotation.Y -= keyLookSpeed * delta
		}
		if rl.IsKeyDown(rl.KeyDown) {
			lookRotation.Y += keyLookSpeed * delta
		}

		sideway := b2f(rl.IsKeyDown(rl.KeyD)) - b2f(rl.IsKeyDown(rl.KeyA))
		forward := b2f(rl.IsKeyDown(rl.KeyW)) - b2f(rl.IsKeyDown(rl.KeyS))
		crouching := rl.IsKeyDown(rl.KeyLeftControl)

		updateBody(&player, lookRotation.X, sideway, forward, rl.IsKeyPressed(rl.KeySpace), crouching)

		targetHeight := standHeight
		if crouching {
			targetHeight = crouchHeight
		}
		headLerp = rl.Lerp(headLerp, targetHeight, 20.0*delta)

		camera.Position = rl.NewVector3(
			player.Position.X,
			player.Position.Y+(bottomHeight+headLerp),
			player.Position.Z,
		)

		if player.IsGrounded && (forward != 0 || sideway != 0) {
			headTimer += delta * 3.0
			walkLerp = rl.Lerp(walkLerp, 1.0, 10.0*delta)
			camera.Fovy = rl.Lerp(camera.Fovy, 55.0, 5.0*delta)
		} else {
			walkLerp = rl.Lerp(walkLerp, 0.0, 10.0*delta)
			camera.Fovy = rl.Lerp(camera.Fovy, 60.0, 5.0*delta)
		}

		lean.X = rl.Lerp(lean.X, sideway*0.02, 10.0*delta)
		lean.Y = rl.Lerp(lean.Y, forward*0.015, 10.0*delta)

		updateCameraFPS(&camera)

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)
		drawLevel()
		rl.EndMode3D()

		// Draw info box
		rl.DrawRectangle(5, 5, 330, 75, rl.Fade(rl.SkyBlue, 0.5))
		rl.DrawRectangleLines(5, 5, 330, 75, rl.Blue)

		rl.DrawText("Camera controls:", 15, 15, 10, rl.Black)
		rl.DrawText("- Move keys: W, A, S, D, Space, Left-Ctrl", 15, 30, 10, rl.Black)
		rl.DrawText("- Look around: arrow keys or mouse", 15, 45, 10, rl.Black)

		velLen := rl.Vector2Length(rl.NewVector2(player.Velocity.X, player.Velocity.Z))
		rl.DrawText(fmt.Sprintf("- Velocity Len: (%06.3f)", velLen), 15, 60, 10, rl.Black)

		rl.EndDrawing()
	}
}

// Helper function to convert bool to float32 (1.0 or 0.0)
func b2f(b bool) float32 {
	if b {
		return 1.0
	}
	return 0.0
}

// Update body considering current world state
func updateBody(body *Body, rot, side, forward float32, jumpPressed, crouchHold bool) {
	input := rl.NewVector2(side, -forward)

	delta := rl.GetFrameTime()

	if !body.IsGrounded {
		body.Velocity.Y -= gravity * delta
	}

	if body.IsGrounded && jumpPressed {
		body.Velocity.Y = jumpForce
		body.IsGrounded = false
	}

	rot64 := float64(rot)
	front := rl.NewVector3(float32(math.Sin(rot64)), 0.0, float32(math.Cos(rot64)))
	right := rl.NewVector3(float32(math.Cos(-rot64)), 0.0, float32(math.Sin(-rot64)))

	desiredDir := rl.NewVector3(input.X*right.X+input.Y*front.X, 0.0, input.X*right.Z+input.Y*front.Z)
	body.Dir = rl.Vector3Lerp(body.Dir, desiredDir, control*delta)

	decel := airDrag
	if body.IsGrounded {
		decel = friction
	}

	hvel := rl.NewVector3(body.Velocity.X*decel, 0.0, body.Velocity.Z*decel)

	hvelLength := rl.Vector3Length(hvel)
	if hvelLength < (maxSpeed * 0.01) {
		hvel = rl.NewVector3(0, 0, 0)
	}

	// This is what creates strafing
	speed := rl.Vector3DotProduct(hvel, body.Dir)

	maxSpd := maxSpeed
	if crouchHold {
		maxSpd = crouchSpeed
	}

	accel := rl.Clamp(maxSpd-speed, 0.0, maxAccel*delta)
	hvel.X += body.Dir.X * accel
	hvel.Z += body.Dir.Z * accel

	body.Velocity.X = hvel.X
	body.Velocity.Z = hvel.Z

	body.Position.X += body.Velocity.X * delta
	body.Position.Y += body.Velocity.Y * delta
	body.Position.Z += body.Velocity.Z * delta

	// Collision system against the floor
	if body.Position.Y <= 0.0 {
		body.Position.Y = 0.0
		body.Velocity.Y = 0.0
		body.IsGrounded = true
	}
}

// Update camera for FPS behaviour
func updateCameraFPS(camera *rl.Camera3D) {
	up := rl.NewVector3(0.0, 1.0, 0.0)
	targetOffset := rl.NewVector3(0.0, 0.0, -1.0)

	// Left and right
	yaw := rl.Vector3RotateByAxisAngle(targetOffset, up, lookRotation.X)

	// Clamp view up
	maxAngleUp := rl.Vector3Angle(up, yaw)
	maxAngleUp -= 0.001
	if -lookRotation.Y > maxAngleUp {
		lookRotation.Y = -maxAngleUp
	}

	// Clamp view down
	maxAngleDown := rl.Vector3Angle(rl.Vector3Negate(up), yaw)
	maxAngleDown *= -1.0
	maxAngleDown += 0.001
	if -lookRotation.Y < maxAngleDown {
		lookRotation.Y = -maxAngleDown
	}

	// Up and down
	right := rl.Vector3Normalize(rl.Vector3CrossProduct(yaw, up))

	// Rotate view vector around right axis
	pitchAngle := -lookRotation.Y - lean.Y
	minPitch := float32(-math.Pi/2.0 + 0.0001)
	maxPitch := float32(math.Pi/2.0 - 0.0001)
	pitchAngle = rl.Clamp(pitchAngle, minPitch, maxPitch)

	pitch := rl.Vector3RotateByAxisAngle(yaw, right, pitchAngle)

	// Head animation
	headTimerPi := float64(headTimer * math.Pi)
	headSin := float32(math.Sin(headTimerPi))
	headCos := float32(math.Cos(headTimerPi))

	const stepRotation float32 = 0.01
	camera.Up = rl.Vector3RotateByAxisAngle(up, pitch, headSin*stepRotation+lean.X)

	// Camera BOB
	const bobSide float32 = 0.1
	const bobUp float32 = 0.15
	bobbing := rl.Vector3Scale(right, headSin*bobSide)
	bobbing.Y = float32(math.Abs(float64(headCos * bobUp)))

	camera.Position = rl.Vector3Add(camera.Position, rl.Vector3Scale(bobbing, walkLerp))
	camera.Target = rl.Vector3Add(camera.Position, pitch)
}

// Draw game level
func drawLevel() {
	const floorExtent = 25
	const tileSize float32 = 5.0
	tileColor1 := rl.NewColor(150, 200, 200, 255)

	// Floor tiles
	for y := -floorExtent; y < floorExtent; y++ {
		for x := -floorExtent; x < floorExtent; x++ {
			if (y&1 != 0) && (x&1 != 0) {
				rl.DrawPlane(
					rl.NewVector3(float32(x)*tileSize, 0.0, float32(y)*tileSize),
					rl.NewVector2(tileSize, tileSize),
					tileColor1,
				)
			} else if (y&1 == 0) && (x&1 == 0) {
				rl.DrawPlane(
					rl.NewVector3(float32(x)*tileSize, 0.0, float32(y)*tileSize),
					rl.NewVector2(tileSize, tileSize),
					rl.LightGray,
				)
			}
		}
	}

	towerSize := rl.NewVector3(16.0, 32.0, 16.0)
	towerColor := rl.NewColor(150, 200, 200, 255)

	towerPos := rl.NewVector3(16.0, 16.0, 16.0)
	rl.DrawCubeV(towerPos, towerSize, towerColor)
	rl.DrawCubeWiresV(towerPos, towerSize, rl.DarkBlue)

	towerPos.X *= -1
	rl.DrawCubeV(towerPos, towerSize, towerColor)
	rl.DrawCubeWiresV(towerPos, towerSize, rl.DarkBlue)

	towerPos.Z *= -1
	rl.DrawCubeV(towerPos, towerSize, towerColor)
	rl.DrawCubeWiresV(towerPos, towerSize, rl.DarkBlue)

	towerPos.X *= -1
	rl.DrawCubeV(towerPos, towerSize, towerColor)
	rl.DrawCubeWiresV(towerPos, towerSize, rl.DarkBlue)

	// Red sun
	rl.DrawSphere(rl.NewVector3(300.0, 300.0, 0.0), 100.0, rl.NewColor(255, 0, 0, 255))
}
