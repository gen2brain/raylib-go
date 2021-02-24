// Package raygui - Simple and easy-to-use IMGUI (immediate mode GUI API) library

package raygui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// GUI controls states
type ControlState int

const (
	Disabled ControlState = iota
	// Normal is the default state for rendering GUI elements.
	Normal
	// Focused indicates the mouse is hovering over the GUI element.
	Focused
	// Pressed indicates the mouse is hovering over the GUI element and LMB is pressed down.
	Pressed
	// Clicked indicates the mouse is hovering over the GUI element and LMB has just been released.
	Clicked
)

// IsColliding will return true if 'point' is within any of the given rectangles.
func IsInAny(point rl.Vector2, rectangles ...rl.Rectangle) bool {
	for _, rect := range rectangles {
		if rl.CheckCollisionPointRec(point, rect) {
			return true
		}
	}
	return false
}

// GetInteractionState determines the current state of a control based on mouse position and
// button states.
func GetInteractionState(rectangles ...rl.Rectangle) ControlState {
	switch {
	case !IsInAny(rl.GetMousePosition(), rectangles...):
		return Normal
	case rl.IsMouseButtonDown(rl.MouseLeftButton):
		return Pressed
	case rl.IsMouseButtonReleased(rl.MouseLeftButton) || rl.IsMouseButtonPressed(rl.MouseLeftButton):
		return Clicked
	default:
		return Focused
	}
}
