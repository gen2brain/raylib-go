package raygui

import "github.com/gen2brain/raylib-go/raylib"

// CheckBox - Check Box element, returns true when active
func CheckBox(bounds rl.Rectangle, checked bool) bool {
	b := bounds.ToInt32()
	state := Normal
	mousePoint := rl.GetMousePosition()

	// Update control
	if rl.CheckCollisionPointRec(mousePoint, bounds) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			state = Pressed
		} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) || rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			state = Normal
			checked = !checked
		} else {
			state = Focused
		}
	}

	// Draw control
	switch state {
	case Normal:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[CheckboxDefaultBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[CheckboxDefaultInsideColor])))
		break
	case Focused:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[CheckboxHoverBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[CheckboxHoverInsideColor])))
		break
	case Pressed:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[CheckboxClickBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[CheckboxClickInsideColor])))
		break
	default:
		break
	}

	if checked {
		rl.DrawRectangle(b.X+int32(style[CheckboxInsideWidth]), b.Y+int32(style[CheckboxInsideWidth]), b.Width-(2*int32(style[CheckboxInsideWidth])), b.Height-(2*int32(style[CheckboxInsideWidth])), rl.GetColor(int32(style[CheckboxDefaultActiveColor])))
	}

	return checked
}
