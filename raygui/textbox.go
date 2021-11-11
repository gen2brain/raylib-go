package raygui

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var backspaceHeld = false
var nextBackspace = time.Now()

// BackspaceRepeatDelay controls the time backspace must be held down before it will repeat.
var BackspaceRepeatDelay = 300 * time.Millisecond

// BackspaceRepeatInterval controls how frequently backspace registers after the initial delay.
var BackspaceRepeatInterval = 60 * time.Millisecond

// TextBox - Text Box element, updates input text
func TextBox(bounds rl.Rectangle, text string) string {
	b := bounds.ToInt32()

	letter := int32(-1)

	// Update control
	state := GetInteractionState(bounds)
	borderColor := TextboxBorderColor
	if state == Pressed || state == Focused {
		borderColor = ToggleActiveBorderColor

		framesCounter2++
		letter = rl.GetKeyPressed()
		if letter != -1 {
			if letter >= 32 && letter < 127 {
				text = fmt.Sprintf("%s%c", text, letter)
			}
		}

		backspacing := rl.IsKeyPressed(rl.KeyBackspace)
		if backspacing {
			nextBackspace = time.Now().Add(BackspaceRepeatDelay)
		} else if rl.IsKeyDown(rl.KeyBackspace) {
			backspacing = time.Since(nextBackspace) >= 0
			if backspacing {
				nextBackspace = time.Now().Add(BackspaceRepeatInterval)
			}
		}
		if backspacing && len(text) > 0 {
			text = text[:len(text)-1]
		}
	}

	DrawBorderedRectangle(b, GetStyle32(TextboxBorderWidth), GetStyleColor(borderColor), GetStyleColor(TextboxInsideColor))
	rl.DrawText(text, b.X+2, b.Y+int32(style[TextboxBorderWidth])+b.Height/2-int32(style[TextboxTextFontsize])/2, int32(style[TextboxTextFontsize]), GetStyleColor(TextboxTextColor))

	if state == Focused || state == Pressed {
		// Draw a cursor, when focused.
		if (framesCounter2/20)%2 == 0 {
			rl.DrawRectangle(b.X+4+rl.MeasureText(text, int32(style[GlobalTextFontsize])), b.Y+2, 1, b.Height-4, rl.GetColor(uint(style[TextboxLineColor])))
		}
	}

	return text
}
