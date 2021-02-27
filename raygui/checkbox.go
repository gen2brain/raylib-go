package raygui

import "github.com/gen2brain/raylib-go/raylib"

// checkboxColoring describes the per-state properties for a CheckBox control.
type checkboxColoring struct {
	Border, Inside Property
}

// checkboxColors lists the styling for each supported state.
var checkboxColors = map[ControlState]checkboxColoring{
	Normal:  {CheckboxDefaultBorderColor, CheckboxDefaultInsideColor},
	Clicked: {CheckboxDefaultBorderColor, CheckboxDefaultInsideColor},
	Pressed: {CheckboxClickBorderColor, CheckboxClickInsideColor},
	Focused: {CheckboxHoverBorderColor, CheckboxHoverInsideColor},
}

// CheckBox - Check Box element, returns true when active
func CheckBox(bounds rl.Rectangle, checked bool) bool {
	state := GetInteractionState(bounds)
	colors, exists := checkboxColors[state]
	if !exists {
		return checked
	}

	// Update control
	if state == Clicked {
		checked = !checked
	}

	// Render control
	box := bounds.ToInt32()
	DrawBorderedRectangle(box, GetStyle32(ToggleBorderWidth), GetStyleColor(colors.Border), GetStyleColor(colors.Inside))
	if checked {
		// Show the inner button.
		DrawInsetRectangle(box, GetStyle32(CheckboxInsideWidth), GetStyleColor(CheckboxDefaultActiveColor))
	}

	return checked
}
