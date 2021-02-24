package raygui

import "github.com/gen2brain/raylib-go/raylib"

// buttonColoring describes the per-state properties for a Button control.
type buttonColoring struct {
	Border, Inside, Text Property
}

// buttonColors lists the styling for each supported state.
var buttonColors = map[ControlState]buttonColoring{
	Normal:  {ButtonDefaultBorderColor, ButtonDefaultInsideColor, ButtonDefaultTextColor},
	Clicked: {ButtonDefaultBorderColor, ButtonDefaultInsideColor, ButtonDefaultTextColor},
	Focused: {ButtonHoverBorderColor, ButtonHoverInsideColor, ButtonHoverTextColor},
	Pressed: {ButtonPressedBorderColor, ButtonPressedInsideColor, ButtonPressedTextColor},
}

// Button - Button element, returns true when clicked
func Button(bounds rl.Rectangle, text string) bool {
	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(text, textHeight)

	ConstrainRectangle(&bounds, textWidth, textWidth+GetStyle32(ButtonTextPadding), textHeight, textHeight+GetStyle32(ButtonTextPadding)/2)

	// Determine what state we're in and whether its valid.
	state := GetInteractionState(bounds)
	colors, exist := buttonColors[state]
	if !exist {
		return false
	}

	// Draw control
	b := bounds.ToInt32()
	DrawBorderedRectangle(b, GetStyle32(ButtonBorderWidth), GetStyleColor(colors.Border), GetStyleColor(colors.Inside))
	rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, GetStyleColor(colors.Text))

	return state == Clicked
}
