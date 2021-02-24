package raygui

import "github.com/gen2brain/raylib-go/raylib"

// Button - Button element, returns true when clicked
func Button(bounds rl.Rectangle, text string) bool {
	b := bounds.ToInt32()
	state := Normal
	mousePoint := rl.GetMousePosition()
	clicked := false

	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(text, textHeight)

	// Update control
	if b.Width < textWidth {
		b.Width = textWidth + int32(style[ButtonTextPadding])
	}

	if b.Height < textHeight {
		b.Height = textHeight + int32(style[ButtonTextPadding])/2
	}

	if rl.CheckCollisionPointRec(mousePoint, bounds) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			state = Pressed
		} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) || rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			clicked = true
		} else {
			state = Focused
		}
	}

	// Draw control
	switch state {
	case Normal:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ButtonDefaultBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ButtonBorderWidth]), b.Y+int32(style[ButtonBorderWidth]), b.Width-(2*int32(style[ButtonBorderWidth])), b.Height-(2*int32(style[ButtonBorderWidth])), rl.GetColor(int32(style[ButtonDefaultInsideColor])))
		rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ButtonDefaultTextColor])))
		break

	case Focused:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ButtonHoverBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ButtonBorderWidth]), b.Y+int32(style[ButtonBorderWidth]), b.Width-(2*int32(style[ButtonBorderWidth])), b.Height-(2*int32(style[ButtonBorderWidth])), rl.GetColor(int32(style[ButtonHoverInsideColor])))
		rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ButtonHoverTextColor])))
		break

	case Pressed:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ButtonPressedBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ButtonBorderWidth]), b.Y+int32(style[ButtonBorderWidth]), b.Width-(2*int32(style[ButtonBorderWidth])), b.Height-(2*int32(style[ButtonBorderWidth])), rl.GetColor(int32(style[ButtonPressedInsideColor])))
		rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ButtonPressedTextColor])))
		break

	default:
		break
	}

	if clicked {
		return true
	}

	return false
}
