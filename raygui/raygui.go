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

// Constrain rectangle will ensure that if width/height are below given minimums, they will
// be set to an ideal minimum.
func ConstrainRectangle(bounds *rl.Rectangle, minWidth, idealWidth, minHeight, idealHeight int32) {
	if int32(bounds.Width) < minWidth {
		bounds.Width = float32(idealWidth)
	}
	if int32(bounds.Height) < minHeight {
		bounds.Height = float32(idealHeight)
	}
}

// InsetRectangle returns the dimensions of a rectangle inset by a margin within an outer rectangle.
func InsetRectangle(outer rl.RectangleInt32, inset int32) rl.RectangleInt32 {
	return rl.RectangleInt32{
		X: outer.X + inset, Y: outer.Y + inset,
		Width: outer.Width - 2*inset, Height: outer.Height - 2*inset,
	}
}

// DrawInsetRectangle is a helper to draw a box inset by a margin of an outer container.
func DrawInsetRectangle(outer rl.RectangleInt32, inset int32, color rl.Color) {
	inside := InsetRectangle(outer, inset)
	rl.DrawRectangle(inside.X, inside.Y, inside.Width, inside.Height, color)
}

// DrawBorderedRectangle is a helper to draw a box with a border around it.
func DrawBorderedRectangle(bounds rl.RectangleInt32, borderWidth int32, borderColor, insideColor rl.Color) {
	inside := InsetRectangle(bounds, borderWidth)
	rl.DrawRectangle(bounds.X, bounds.Y, bounds.Width, bounds.Height, borderColor)
	rl.DrawRectangle(inside.X, inside.Y, inside.Width, inside.Height, insideColor)
}
