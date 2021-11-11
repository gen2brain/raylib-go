package raygui

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

// comboboxColoring describes the per-state colors for a Combobox control.
type comboboxColoring struct {
	Border Property
	Inside Property
	List   Property
	Text   Property
}

// comboboxColors lists the styling for each supported state.
var comboboxColors = map[ControlState]comboboxColoring{
	Normal:  {ComboboxDefaultBorderColor, ComboboxDefaultInsideColor, ComboboxDefaultListTextColor, ComboboxDefaultTextColor},
	Clicked: {ComboboxDefaultBorderColor, ComboboxDefaultInsideColor, ComboboxDefaultListTextColor, ComboboxDefaultTextColor},
	Focused: {ComboboxHoverBorderColor, ComboboxHoverInsideColor, ComboboxHoverListTextColor, ComboboxHoverTextColor},
	Pressed: {ComboboxPressedBorderColor, ComboboxPressedInsideColor, ComboboxPressedListTextColor, ComboboxPressedTextColor},
}

// ComboBox draws a simplified version of a ComboBox allowing the user to select a string
// from a list accompanied by an N/M counter. The widget does not provide a drop-down/completion
// or any input support.
func ComboBox(bounds rl.Rectangle, comboText []string, active int) int {
	// Reject invalid selections and disable rendering.
	comboCount := len(comboText)
	if active < 0 || active >= comboCount {
		rl.TraceLog(rl.LogWarning, "ComboBox active expects 0 <= active <= %d", comboCount)
		return -1
	}

	// Calculate text dimensions.
	textHeight := GetStyle32(GlobalTextFontsize)
	activeText := comboText[active]
	textWidth := rl.MeasureText(activeText, textHeight)

	// Ensure box is large enough.
	ConstrainRectangle(&bounds, textWidth, textWidth+GetStyle32(ToggleTextPadding), textHeight, textHeight+GetStyle32(ToggleTextPadding))
	b := bounds.ToInt32()

	// Generate the worst-case sizing of the counter so we can avoid resizing it as the numbers go up/down.
	clickWidth := rl.MeasureText(fmt.Sprintf("%d/%d", comboCount, comboCount), b.Height)

	// Counter shows the index of the selection and the maximum number, e.g. "1/3".
	counter := rl.NewRectangle(bounds.X+bounds.Width+float32(style[ComboboxPadding]), bounds.Y, float32(clickWidth), float32(b.Height))
	c := counter.ToInt32()

	// Determine if the user is interacting with the control, and if so, which state it is in.
	state := GetInteractionState(bounds, counter)
	colors, exists := comboboxColors[state]
	if !exists {
		return active
	}

	// Update the control when the user releases the mouse over it.
	if state == Clicked {
		// increment but wrap to 0 on reaching end-of-list.
		active = (active + 1) % comboCount
	}

	// Render the box itself
	DrawBorderedRectangle(b, GetStyle32(ComboboxBorderWidth), GetStyleColor(colors.Border), GetStyleColor(colors.Inside))
	rl.DrawText(activeText, b.X+((b.Width/2)-(rl.MeasureText(activeText, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, GetStyleColor(colors.Text))

	// Render the accompanying "clicks" box showing the element counter.
	DrawBorderedRectangle(c, GetStyle32(ComboboxBorderWidth), GetStyleColor(colors.Border), GetStyleColor(colors.Inside))
	counterText := fmt.Sprintf("%d/%d", active+1, comboCount)
	rl.DrawText(counterText, c.X+((c.Width/2)-(rl.MeasureText(counterText, textHeight)/2)), c.Y+((c.Height/2)-(textHeight/2)), textHeight, GetStyleColor(colors.List))

	return active
}
