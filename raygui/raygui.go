// Package raygui - Simple and easy-to-use IMGUI (immediate mode GUI API) library

package raygui

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
)
