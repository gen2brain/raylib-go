package raygui

import "github.com/gen2brain/raylib-go/raylib"

// togglebuttonColoring describes the per-state colors for a ToggleBox control.
type togglebuttonColoring struct {
	Border, Inside, Text Property
}

// togglebuttonColors lists the styling for each supported state.
var togglebuttonColors = map[ControlState]togglebuttonColoring{
	Normal: {ToggleDefaultBorderColor, ToggleDefaultInsideColor, ToggleDefaultTextColor},
	// Hijacking 'Clicked' for the 'active' state.
	Clicked: {ToggleActiveBorderColor, ToggleActiveInsideColor, ToggleDefaultTextColor},
	Pressed: {TogglePressedBorderColor, TogglePressedInsideColor, TogglePressedTextColor},
	Focused: {ToggleHoverBorderColor, ToggleHoverInsideColor, ToggleHoverTextColor},
}

// ToggleButton - Toggle Button element, returns true when active
func ToggleButton(bounds rl.Rectangle, text string, active bool) bool {
	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(text, textHeight)

	ConstrainRectangle(&bounds, textWidth, textWidth+GetStyle32(ToggleTextPadding), textHeight, textHeight+GetStyle32(ToggleTextPadding))

	state := GetInteractionState(bounds)
	if state == Clicked {
		active = !active
		state = Normal
	}

	// Hijack 'Clicked' as the 'active' state
	if state == Normal && active {
		state = Clicked
	}

	colors, exists := togglebuttonColors[state]
	if !exists {
		return active
	}

	// Draw control
	b := bounds.ToInt32()
	DrawBorderedRectangle(b, GetStyle32(ToggleBorderWidth), GetStyleColor(colors.Border), GetStyleColor(colors.Inside))
	rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, GetStyleColor(colors.Text))

	return active
}
