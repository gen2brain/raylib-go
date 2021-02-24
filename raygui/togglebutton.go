package raygui

import "github.com/gen2brain/raylib-go/raylib"

// ToggleButton - Toggle Button element, returns true when active
func ToggleButton(bounds rl.Rectangle, text string, active bool) bool {
	b := bounds.ToInt32()
	state := Normal
	mousePoint := rl.GetMousePosition()

	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(text, textHeight)

	// Update control
	if b.Width < textWidth {
		b.Width = textWidth + int32(style[ToggleTextPadding])
	}
	if b.Height < textHeight {
		b.Height = textHeight + int32(style[ToggleTextPadding])/2
	}

	if rl.CheckCollisionPointRec(mousePoint, bounds) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			state = Pressed
		} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) || rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			state = Normal
			active = !active
		} else {
			state = Focused
		}
	}

	// Draw control
	switch state {
	case Normal:
		if active {
			rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ToggleActiveBorderColor])))
			rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[ToggleActiveInsideColor])))
			rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ToggleDefaultTextColor])))
		} else {
			rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ToggleDefaultBorderColor])))
			rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[ToggleDefaultInsideColor])))
			rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ToggleDefaultTextColor])))
		}
		break
	case Focused:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ToggleHoverBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[ToggleHoverInsideColor])))
		rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ToggleHoverTextColor])))
		break
	case Pressed:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[TogglePressedBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[TogglePressedInsideColor])))
		rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[TogglePressedTextColor])))
		break
	default:
		break
	}

	return active
}
