// Package raygui - Simple and easy-to-use immediate-mode gui library, a Go port of raygui (https://github.com/raysan5/raygui).
package raygui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type (
	ControlID     uint16
	PropertyID    uint16
	PropertyValue uint32

	IconID uint32
)

const (
	borderElement PropertyID = iota
	baseElement
	textElement
	otherElement
)

var (
	guiState  = STATE_NORMAL
	guiFont   rl.Font
	guiLocked bool
	guiAlpha  float32 = 1.0

	guiIconScale int32 = 1
	guiIconsPtr        = guiIcons[:]

	guiTooltip    bool
	guiTooltipPtr string

	guiControlExclusiveMode bool
	guiControlExclusiveRec  rl.Rectangle

	textBoxCursorIndex int32
	autoCursorCounter  int32

	guiStyle       [RAYGUI_MAX_CONTROLS * (RAYGUI_MAX_PROPS_BASE + RAYGUI_MAX_PROPS_EXTENDED)]uint32
	guiStyleLoaded bool
)

func guiButtonDown() bool {
	return rl.IsMouseButtonDown(rl.MouseLeftButton) || rl.IsGamepadButtonDown(0, rl.GamepadButtonRightFaceDown)
}

func guiButtonPressed() bool {
	return rl.IsMouseButtonPressed(rl.MouseLeftButton) || rl.IsGamepadButtonPressed(0, rl.GamepadButtonRightFaceDown)
}

func guiButtonPressedMid() bool {
	return rl.IsMouseButtonPressed(rl.MouseMiddleButton) || rl.IsGamepadButtonPressed(0, rl.GamepadButtonRightFaceUp)
}

func guiButtonReleased() bool {
	return rl.IsMouseButtonReleased(rl.MouseLeftButton) || rl.IsGamepadButtonReleased(0, rl.GamepadButtonRightFaceDown)
}

func guiScrollDelta() float32 {
	return rl.GetMouseWheelMove() + (rl.GetGamepadAxisMovement(0, rl.GamepadAxisRightTrigger) + 1) - (rl.GetGamepadAxisMovement(0, rl.GamepadAxisLeftTrigger) + 1)
}

func guiPointerPosition() rl.Vector2 {
	return rl.GetMousePosition()
}

func guiKeyDown(key int32) bool {
	return rl.IsKeyDown(key)
}

func guiKeyPressed(key int32) bool {
	return rl.IsKeyPressed(key)
}

func guiInputKey() int32 {
	return rl.GetCharPressed()
}

func checkBoundsID(src, dst rl.Rectangle) bool {
	return int32(src.X) == int32(dst.X) && int32(src.Y) == int32(dst.Y) && int32(src.Width) == int32(dst.Width) && int32(src.Height) == int32(dst.Height)
}

// Enable - Enable gui global state
func Enable() {
	if guiState == STATE_DISABLED {
		guiState = STATE_NORMAL
	}
}

// Disable - Disable gui global state
func Disable() {
	if guiState == STATE_NORMAL {
		guiState = STATE_DISABLED
	}
}

// Lock - Lock gui global state
func Lock() {
	guiLocked = true
}

// Unlock - Unlock gui global state
func Unlock() {
	guiLocked = false
}

// IsLocked - Check if gui is locked (global state)
func IsLocked() bool {
	return guiLocked
}

// SetAlpha - Set gui controls alpha global state
func SetAlpha(alpha float32) {
	if alpha < 0.0 {
		alpha = 0.0
	} else if alpha > 1.0 {
		alpha = 1.0
	}

	guiAlpha = alpha
}

// SetState - Set gui state (global state)
func SetState(state PropertyValue) {
	guiState = state
}

// GetState - Get gui state (global state)
func GetState() PropertyValue {
	return guiState
}

// SetFont - Set custom gui font
func SetFont(font rl.Font) {
	if font.Texture.ID > 0 {
		if !guiStyleLoaded {
			LoadStyleDefault()
		}

		guiFont = font
	}
}

// GetFont - Get custom gui font
func GetFont() rl.Font {
	return guiFont
}

// SetStyle - Set control style property value
func SetStyle(control ControlID, property PropertyID, value PropertyValue) {
	if !guiStyleLoaded {
		LoadStyleDefault()
	}
	guiStyle[int(control)*(RAYGUI_MAX_PROPS_BASE+RAYGUI_MAX_PROPS_EXTENDED)+int(property)] = uint32(value)

	if control == 0 && property < RAYGUI_MAX_PROPS_BASE {
		for i := 1; i < RAYGUI_MAX_CONTROLS; i++ {
			guiStyle[i*(RAYGUI_MAX_PROPS_BASE+RAYGUI_MAX_PROPS_EXTENDED)+int(property)] = uint32(value)
		}
	}
}

// GetStyle - Get control style property value
func GetStyle(control ControlID, property PropertyID) PropertyValue {
	if !guiStyleLoaded {
		LoadStyleDefault()
	}
	return PropertyValue(guiStyle[int(control)*(RAYGUI_MAX_PROPS_BASE+RAYGUI_MAX_PROPS_EXTENDED)+int(property)])
}

func getStyle(control ControlID, property PropertyID) int32 {
	return int32(GetStyle(control, property))
}

func setStyle(control ControlID, property PropertyID, value int32) {
	SetStyle(control, property, PropertyValue(uint32(value)))
}
