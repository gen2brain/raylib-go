package raygui

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
)

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

	activeText := comboText[active]

	// style sizing.
	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(activeText, textHeight)
	borderWidth := int32(style[ComboboxBorderWidth])
	textPadding := int32(style[ToggleTextPadding])

	b := bounds.ToInt32()
	if b.Width < textWidth {
		b.Width = textWidth + textPadding
		bounds.Width = float32(b.Width)
	}
	if b.Height < textHeight {
		b.Height = textHeight + textPadding
		bounds.Height = float32(b.Height)
	}

	// Identify what the counter is going to look like with max digits so we don't resize it.
	clickWidth := rl.MeasureText(fmt.Sprintf("%d/%d", comboCount, comboCount), b.Height)

	click := rl.NewRectangle(bounds.X+bounds.Width+float32(style[ComboboxPadding]), bounds.Y, float32(clickWidth), float32(b.Height))
	c := click.ToInt32()
	mousePoint := rl.GetMousePosition()
	state := Normal
	if rl.CheckCollisionPointRec(mousePoint, bounds) || rl.CheckCollisionPointRec(mousePoint, click) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			state = Pressed
		} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) || rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			state = Pressed
		} else {
			state = Focused
		}
	}

	// Draw control
	var borderColor, insideColor, listColor, textColor rl.Color

	switch state {
	case Normal:
		borderColor = rl.GetColor(int32(style[ComboboxDefaultBorderColor]))
		insideColor = rl.GetColor(int32(style[ComboboxDefaultInsideColor]))
		listColor = rl.GetColor(int32(style[ComboboxDefaultListTextColor]))
		textColor = rl.GetColor(int32(style[ComboboxDefaultTextColor]))

	case Focused:
		borderColor = rl.GetColor(int32(style[ComboboxHoverBorderColor]))
		insideColor = rl.GetColor(int32(style[ComboboxHoverInsideColor]))
		listColor = rl.GetColor(int32(style[ComboboxHoverListTextColor]))
		textColor = rl.GetColor(int32(style[ComboboxHoverTextColor]))

	case Pressed:
		borderColor = rl.GetColor(int32(style[ComboboxPressedBorderColor]))
		insideColor = rl.GetColor(int32(style[ComboboxPressedInsideColor]))
		listColor = rl.GetColor(int32(style[ComboboxPressedListTextColor]))
		textColor = rl.GetColor(int32(style[ComboboxPressedTextColor]))

	default:
		rl.TraceLog(rl.LogWarning, "ComboBox in unrecognized state %d", state)
		return -1
	}

	// Render the box itself
	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, borderColor)
	rl.DrawRectangle(b.X+borderWidth, b.Y+borderWidth, b.Width-(2*borderWidth), b.Height-(2*borderWidth), insideColor)
	rl.DrawText(activeText, b.X+((b.Width/2)-(rl.MeasureText(activeText, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, textColor)

	// Render the accompanying "clicks" box showing the element counter.
	rl.DrawRectangle(c.X, c.Y, c.Width, c.Height, borderColor)
	rl.DrawRectangle(c.X+borderWidth, c.Y+borderWidth, c.Width-(2*borderWidth), c.Height-(2*borderWidth), insideColor)
	companionText := fmt.Sprintf("%d/%d", active+1, comboCount)
	rl.DrawText(companionText, c.X+((c.Width/2)-(rl.MeasureText(companionText, textHeight)/2)), c.Y+((c.Height/2)-(textHeight/2)), textHeight, listColor)

	if rl.CheckCollisionPointRec(mousePoint, bounds) || rl.CheckCollisionPointRec(mousePoint, click) {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			active++
			if active >= comboCount {
				active = 0
			}
		}
	}

	return active
}
