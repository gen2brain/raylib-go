package rl

// GetCameraForward - Returns the cameras forward vector (normalized)
func GetCameraForward(camera *Camera) Vector3 {
	return Vector3Normalize(Vector3Subtract(camera.Target, camera.Position))
}

// GetCameraUp - Returns the cameras up vector (normalized)
// Note: The up vector might not be perpendicular to the forward vector
func GetCameraUp(camera *Camera) Vector3 {
	return Vector3Normalize(camera.Up)
}

// GetCameraRight - Returns the cameras right vector (normalized)
func GetCameraRight(camera *Camera) Vector3 {
	forward := GetCameraForward(camera)
	up := GetCameraUp(camera)

	return Vector3CrossProduct(forward, up)
}

// CameraMoveForward - Moves the camera in its forward direction
func CameraMoveForward(camera *Camera, distance float32, moveInWorldPlane uint8) {
	forward := GetCameraForward(camera)

	if moveInWorldPlane != 0 {
		// Project vector onto world plane
		forward.Y = float32(0)
		forward = Vector3Normalize(forward)
	}

	// Scale by distance
	forward = Vector3Scale(forward, distance)

	// Move position and target
	camera.Position = Vector3Add(camera.Position, forward)
	camera.Target = Vector3Add(camera.Target, forward)
}

// CameraMoveUp - Moves the camera in its up direction
func CameraMoveUp(camera *Camera, distance float32) {
	up := GetCameraUp(camera)

	// Scale by distance
	up = Vector3Scale(up, distance)

	// Move position and target
	camera.Position = Vector3Add(camera.Position, up)
	camera.Target = Vector3Add(camera.Target, up)
}

// CameraMoveRight - Moves the camera target in its current right direction
func CameraMoveRight(camera *Camera, distance float32, moveInWorldPlane uint8) {
	right := GetCameraRight(camera)

	if moveInWorldPlane != 0 {
		// Project vector onto world plane
		right.Y = float32(0)
		right = Vector3Normalize(right)
	}

	// Scale by distance
	right = Vector3Scale(right, distance)

	// Move position and target
	camera.Position = Vector3Add(camera.Position, right)
	camera.Target = Vector3Add(camera.Target, right)
}

// CameraMoveToTarget - Moves the camera position closer/farther to/from the camera target
func CameraMoveToTarget(camera *Camera, delta float32) {
	distance := Vector3Distance(camera.Position, camera.Target)

	// Apply delta
	distance = distance + delta

	// Distance must be greater than 0
	if distance <= float32(0) {
		distance = 0.001
	}

	// Set new distance by moving the position along the forward vector
	forward := GetCameraForward(camera)
	camera.Position = Vector3Add(camera.Target, Vector3Scale(forward, -distance))
}

// CameraYaw - Rotates the camera around its up vector
// Yaw is "looking left and right"
// If rotateAroundTarget is false, the camera rotates around its position
// Note: angle must be provided in radians
func CameraYaw(camera *Camera, angle float32, rotateAroundTarget uint8) {
	// Rotation axis
	var up = GetCameraUp(camera)

	// View vector
	var targetPosition = Vector3Subtract(camera.Target, camera.Position)

	// Rotate view vector around up axis
	targetPosition = Vector3RotateByAxisAngle(targetPosition, up, angle)

	if rotateAroundTarget != 0 {
		// Move position relative to target
		camera.Position = Vector3Subtract(camera.Target, targetPosition)
	} else {
		// Move target relative to position
		camera.Target = Vector3Add(camera.Position, targetPosition)
	}
}

// CameraPitch - Rotates the camera around its right vector, pitch is "looking up and down"
//   - lockView prevents camera overrotation (aka "somersaults")
//   - rotateAroundTarget defines if rotation is around target or around its position
//   - rotateUp rotates the up direction as well (typically only useful in CAMERA_FREE)
//
// NOTE: angle must be provided in radians
func CameraPitch(camera *Camera, angle float32, lockView uint8, rotateAroundTarget uint8, rotateUp uint8) {
	// Up direction
	var up = GetCameraUp(camera)

	// View vector
	var targetPosition = Vector3Subtract(camera.Target, camera.Position)

	if lockView != 0 {
		// In these camera modes we clamp the Pitch angle
		// to allow only viewing straight up or down.

		// Clamp view up
		maxAngleUp := Vector3Angle(up, targetPosition)
		maxAngleUp = maxAngleUp - 0.001 // avoid numerical errors
		if angle > maxAngleUp {
			angle = maxAngleUp
		}

		// Clamp view down
		maxAngleDown := Vector3Angle(Vector3Negate(up), targetPosition)
		maxAngleDown = maxAngleDown * -1.0  // downwards angle is negative
		maxAngleDown = maxAngleDown + 0.001 // avoid numerical errors
		if angle < maxAngleDown {
			angle = maxAngleDown
		}
	}

	// Rotation axis
	var right = GetCameraRight(camera)

	// Rotate view vector around right axis
	targetPosition = Vector3RotateByAxisAngle(targetPosition, right, angle)

	if rotateAroundTarget != 0 {
		// Move position relative to target
		camera.Position = Vector3Subtract(camera.Target, targetPosition)
	} else {
		// Move target relative to position
		camera.Target = Vector3Add(camera.Position, targetPosition)
	}

	if rotateUp != 0 {
		// Rotate up direction around right axis
		camera.Up = Vector3RotateByAxisAngle(camera.Up, right, angle)
	}
}

// CameraRoll - Rotates the camera around its forward vector
// Roll is "turning your head sideways to the left or right"
// Note: angle must be provided in radians
func CameraRoll(camera *Camera, angle float32) {
	// Rotation axis
	var forward = GetCameraForward(camera)

	// Rotate up direction around forward axis
	camera.Up = Vector3RotateByAxisAngle(camera.Up, forward, angle)
}

// GetCameraViewMatrix - Returns the camera view matrix
func GetCameraViewMatrix(camera *Camera) Matrix {
	return MatrixLookAt(camera.Position, camera.Target, camera.Up)
}

// GetCameraProjectionMatrix - Returns the camera projection matrix
func GetCameraProjectionMatrix(camera *Camera, aspect float32) Matrix {
	if camera.Projection == CameraPerspective {
		return MatrixPerspective(camera.Fovy*(Pi/180.0), aspect, 0.01, 1000.0)
	} else if camera.Projection == CameraOrthographic {
		top := camera.Fovy / 2.0
		right := top * aspect

		return MatrixOrtho(-right, right, -top, top, 0.01, 1000.0)
	}

	return MatrixIdentity()
}

// UpdateCamera - Update camera position for selected mode
// Camera mode: CameraFree, CameraFirstPerson, CameraThirdPerson, CameraOrbital or Custom
func UpdateCamera(camera *Camera, mode CameraMode) {
	var mousePositionDelta = GetMouseDelta()

	moveInWorldPlaneBool := mode == CameraFirstPerson || mode == CameraThirdPerson
	var moveInWorldPlane uint8
	if moveInWorldPlaneBool {
		moveInWorldPlane = 1
	}

	rotateAroundTargetBool := mode == CameraThirdPerson || mode == CameraOrbital
	var rotateAroundTarget uint8
	if rotateAroundTargetBool {
		rotateAroundTarget = 1
	}

	lockViewBool := mode == CameraFirstPerson || mode == CameraThirdPerson || mode == CameraOrbital
	var lockView uint8
	if lockViewBool {
		lockView = 1
	}

	var rotateUp uint8

	if mode == CameraOrbital {
		// Orbital can just orbit
		var rotation = MatrixRotate(GetCameraUp(camera), 0.5*GetFrameTime())
		var view = Vector3Subtract(camera.Position, camera.Target)
		view = Vector3Transform(view, rotation)
		camera.Position = Vector3Add(camera.Target, view)
	} else {
		// Camera rotation
		if IsKeyDown(KeyDown) {
			CameraPitch(camera, -0.03, lockView, rotateAroundTarget, rotateUp)
		}
		if IsKeyDown(KeyUp) {
			CameraPitch(camera, 0.03, lockView, rotateAroundTarget, rotateUp)
		}
		if IsKeyDown(KeyRight) {
			CameraYaw(camera, -0.03, rotateAroundTarget)
		}
		if IsKeyDown(KeyLeft) {
			CameraYaw(camera, 0.03, rotateAroundTarget)
		}
		if IsKeyDown(KeyQ) {
			CameraRoll(camera, -0.03)
		}
		if IsKeyDown(KeyE) {
			CameraRoll(camera, 0.03)
		}

		// Camera movement
		if !(IsGamepadAvailable(0)) {
			// Camera pan (for CameraFree)
			if mode == CameraFree && IsMouseButtonDown(MouseMiddleButton) {
				var mouseDelta = GetMouseDelta()
				if mouseDelta.X > 0.0 {
					CameraMoveRight(camera, 0.2, moveInWorldPlane)
				}
				if mouseDelta.X < 0.0 {
					CameraMoveRight(camera, -0.2, moveInWorldPlane)
				}
				if mouseDelta.Y > 0.0 {
					CameraMoveUp(camera, -0.2)
				}
				if mouseDelta.Y < 0.0 {
					CameraMoveUp(camera, 0.2)
				}
			} else {
				// Mouse support
				CameraYaw(camera, -mousePositionDelta.X*0.003, rotateAroundTarget)
				CameraPitch(camera, -mousePositionDelta.Y*0.003, lockView, rotateAroundTarget, rotateUp)
			}

			// Keyboard support
			if IsKeyDown(KeyW) {
				CameraMoveForward(camera, 0.09, moveInWorldPlane)
			}
			if IsKeyDown(KeyA) {
				CameraMoveRight(camera, -0.09, moveInWorldPlane)
			}
			if IsKeyDown(KeyS) {
				CameraMoveForward(camera, -0.09, moveInWorldPlane)
			}
			if IsKeyDown(KeyD) {
				CameraMoveRight(camera, 0.09, moveInWorldPlane)
			}
		} else {
			// Gamepad controller support
			CameraYaw(camera, -(GetGamepadAxisMovement(0, GamepadAxisRightX)*float32(2))*0.003, rotateAroundTarget)
			CameraPitch(camera, -(GetGamepadAxisMovement(0, GamepadAxisRightY)*float32(2))*0.003, lockView, rotateAroundTarget, rotateUp)

			if GetGamepadAxisMovement(0, GamepadAxisLeftY) <= -0.25 {
				CameraMoveForward(camera, 0.09, moveInWorldPlane)
			}
			if GetGamepadAxisMovement(0, GamepadAxisLeftX) <= -0.25 {
				CameraMoveRight(camera, -0.09, moveInWorldPlane)
			}
			if GetGamepadAxisMovement(0, GamepadAxisLeftY) >= 0.25 {
				CameraMoveForward(camera, -0.09, moveInWorldPlane)
			}
			if GetGamepadAxisMovement(0, GamepadAxisLeftX) >= 0.25 {
				CameraMoveRight(camera, 0.09, moveInWorldPlane)
			}
		}

		if mode == CameraFree {
			if IsKeyDown(KeySpace) {
				CameraMoveUp(camera, 0.09)
			}
			if IsKeyDown(KeyLeftControl) {
				CameraMoveUp(camera, -0.09)
			}
		}
	}

	if mode == CameraThirdPerson || mode == CameraOrbital || mode == CameraFree {
		// Zoom target distance
		CameraMoveToTarget(camera, -GetMouseWheelMove())
		if IsKeyPressed(KeyKpSubtract) {
			CameraMoveToTarget(camera, 2.0)
		}
		if IsKeyPressed(KeyKpAdd) {
			CameraMoveToTarget(camera, -2.0)
		}
	}
}

// UpdateCameraPro - Update camera movement, movement/rotation values should be provided by user
func UpdateCameraPro(camera *Camera, movement Vector3, rotation Vector3, zoom float32) {
	// Required values
	// movement.X - Move forward/backward
	// movement.Y - Move right/left
	// movement.Z - Move up/down
	// rotation.X - yaw
	// rotation.Y - pitch
	// rotation.Z - roll
	// zoom - Move towards target

	lockView := uint8(1)
	rotateAroundTarget := uint8(0)
	rotateUp := uint8(0)
	moveInWorldPlane := uint8(1)

	// Camera rotation
	CameraPitch(camera, -rotation.Y*(Pi/180.0), lockView, rotateAroundTarget, rotateUp)
	CameraYaw(camera, -rotation.X*(Pi/180.0), rotateAroundTarget)
	CameraRoll(camera, rotation.Z*(Pi/180.0))

	// Camera movement
	CameraMoveForward(camera, movement.X, moveInWorldPlane)
	CameraMoveRight(camera, movement.Y, moveInWorldPlane)
	CameraMoveUp(camera, movement.Z)

	// Zoom target distance
	CameraMoveToTarget(camera, zoom)
}
