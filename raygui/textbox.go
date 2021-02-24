package raygui

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
)

// TextBox - Text Box element, updates input text
func TextBox(bounds rl.Rectangle, text string) string {
	b := bounds.ToInt32()
	state := Normal

	mousePoint := rl.GetMousePosition()
	letter := int32(-1)

	// Update control
	if rl.CheckCollisionPointRec(mousePoint, bounds) {
		state = Focused // NOTE: PRESSED state is not used on this control

		framesCounter2++

		letter = rl.GetKeyPressed()
		if letter != -1 {
			if letter >= 32 && letter < 127 {
				text = fmt.Sprintf("%s%c", text, letter)
			}
		}

		if rl.IsKeyPressed(rl.KeyBackspace) {
			if len(text) > 0 {
				text = text[:len(text)-1]
			}
		}
	}

	// Draw control
	switch state {
	case Normal:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[TextboxBorderColor])))
		rl.DrawRectangle(b.X+int32(style[TextboxBorderWidth]), b.Y+int32(style[TextboxBorderWidth]), b.Width-(int32(style[TextboxBorderWidth])*2), b.Height-(int32(style[TextboxBorderWidth])*2), rl.GetColor(int32(style[TextboxInsideColor])))
		rl.DrawText(text, b.X+2, b.Y+int32(style[TextboxBorderWidth])+b.Height/2-int32(style[TextboxTextFontsize])/2, int32(style[TextboxTextFontsize]), rl.GetColor(int32(style[TextboxTextColor])))
		break
	case Focused:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ToggleActiveBorderColor])))
		rl.DrawRectangle(b.X+int32(style[TextboxBorderWidth]), b.Y+int32(style[TextboxBorderWidth]), b.Width-(int32(style[TextboxBorderWidth])*2), b.Height-(int32(style[TextboxBorderWidth])*2), rl.GetColor(int32(style[TextboxInsideColor])))
		rl.DrawText(text, b.X+2, b.Y+int32(style[TextboxBorderWidth])+b.Height/2-int32(style[TextboxTextFontsize])/2, int32(style[TextboxTextFontsize]), rl.GetColor(int32(style[TextboxTextColor])))

		if (framesCounter2/20)%2 == 0 && rl.CheckCollisionPointRec(mousePoint, bounds) {
			rl.DrawRectangle(b.X+4+rl.MeasureText(text, int32(style[GlobalTextFontsize])), b.Y+2, 1, b.Height-4, rl.GetColor(int32(style[TextboxLineColor])))
		}
		break
	case Pressed:
		break
	default:
		break
	}

	return text
}
