package raygui

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
)

// For spinner
var (
	framesCounter  int
	framesCounter2 int
	valueSpeed     bool
)

// Spinner - Spinner element, returns selected value
func Spinner(bounds rl.Rectangle, value, minValue, maxValue int) int {
	b := bounds.ToInt32()
	state := Normal

	mousePoint := rl.GetMousePosition()
	labelBoxBound := rl.RectangleInt32{b.X + b.Width/4 + 1, b.Y, b.Width / 2, b.Height}
	leftButtonBound := rl.RectangleInt32{b.X, b.Y, b.Width / 4, b.Height}
	rightButtonBound := rl.RectangleInt32{b.X + b.Width - b.Width/4 + 1, b.Y, b.Width / 4, b.Height}

	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(fmt.Sprintf("%d", value), textHeight)

	buttonSide := 0

	// Update control
	if rl.CheckCollisionPointRec(mousePoint, leftButtonBound.ToFloat32()) || rl.CheckCollisionPointRec(mousePoint, rightButtonBound.ToFloat32()) || rl.CheckCollisionPointRec(mousePoint, labelBoxBound.ToFloat32()) {
		if rl.IsKeyDown(rl.KeyLeft) {
			state = Pressed
			buttonSide = 1

			if value > minValue {
				value--
			}
		} else if rl.IsKeyDown(rl.KeyRight) {
			state = Pressed
			buttonSide = 2

			if value < maxValue {
				value++
			}
		}
	}

	if rl.CheckCollisionPointRec(mousePoint, leftButtonBound.ToFloat32()) {
		buttonSide = 1
		state = Focused

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			if !valueSpeed {
				if value > minValue {
					value--
				}
				valueSpeed = true
			} else {
				framesCounter++
			}

			state = Pressed

			if value > minValue {
				if framesCounter >= 30 {
					value--
				}
			}
		}
	} else if rl.CheckCollisionPointRec(mousePoint, rightButtonBound.ToFloat32()) {
		buttonSide = 2
		state = Focused

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			if !valueSpeed {
				if value < maxValue {
					value++
				}
				valueSpeed = true
			} else {
				framesCounter++
			}

			state = Pressed

			if value < maxValue {
				if framesCounter >= 30 {
					value++
				}
			}
		}
	} else if !rl.CheckCollisionPointRec(mousePoint, labelBoxBound.ToFloat32()) {
		buttonSide = 0
	}

	if rl.IsMouseButtonUp(rl.MouseLeftButton) {
		valueSpeed = false
		framesCounter = 0
	}

	// Draw control
	switch state {
	case Normal:
		rl.DrawRectangle(leftButtonBound.X, leftButtonBound.Y, leftButtonBound.Width, leftButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
		rl.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

		rl.DrawRectangle(rightButtonBound.X, rightButtonBound.Y, rightButtonBound.Width, rightButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
		rl.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

		rl.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(rl.MeasureText("+", textHeight))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))
		rl.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(rl.MeasureText("-", textHeight))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))

		rl.DrawRectangle(labelBoxBound.X, labelBoxBound.Y, labelBoxBound.Width, labelBoxBound.Height, rl.GetColor(int32(style[SpinnerLabelBorderColor])))
		rl.DrawRectangle(labelBoxBound.X+1, labelBoxBound.Y+1, labelBoxBound.Width-2, labelBoxBound.Height-2, rl.GetColor(int32(style[SpinnerLabelInsideColor])))

		rl.DrawText(fmt.Sprintf("%d", value), labelBoxBound.X+(labelBoxBound.Width/2-textWidth/2), labelBoxBound.Y+(labelBoxBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultTextColor])))
		break
	case Focused:
		if buttonSide == 1 {
			rl.DrawRectangle(leftButtonBound.X, leftButtonBound.Y, leftButtonBound.Width, leftButtonBound.Height, rl.GetColor(int32(style[SpinnerHoverButtonBorderColor])))
			rl.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, rl.GetColor(int32(style[SpinnerHoverButtonInsideColor])))

			rl.DrawRectangle(rightButtonBound.X, rightButtonBound.Y, rightButtonBound.Width, rightButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			rl.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			rl.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(rl.MeasureText("+", textHeight))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerHoverSymbolColor])))
			rl.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(rl.MeasureText("-", textHeight))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))
		} else if buttonSide == 2 {
			rl.DrawRectangle(leftButtonBound.X, leftButtonBound.Y, leftButtonBound.Width, leftButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			rl.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			rl.DrawRectangle(rightButtonBound.X, rightButtonBound.Y, rightButtonBound.Width, rightButtonBound.Height, rl.GetColor(int32(style[SpinnerHoverButtonBorderColor])))
			rl.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, rl.GetColor(int32(style[SpinnerHoverButtonInsideColor])))

			rl.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(rl.MeasureText("+", textHeight))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))
			rl.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(rl.MeasureText("-", textHeight))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerHoverSymbolColor])))
		}

		rl.DrawRectangle(labelBoxBound.X, labelBoxBound.Y, labelBoxBound.Width, labelBoxBound.Height, rl.GetColor(int32(style[SpinnerLabelBorderColor])))
		rl.DrawRectangle(labelBoxBound.X+1, labelBoxBound.Y+1, labelBoxBound.Width-2, labelBoxBound.Height-2, rl.GetColor(int32(style[SpinnerLabelInsideColor])))

		rl.DrawText(fmt.Sprintf("%d", value), labelBoxBound.X+(labelBoxBound.Width/2-textWidth/2), labelBoxBound.Y+(labelBoxBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerHoverTextColor])))
		break
	case Pressed:
		if buttonSide == 1 {
			rl.DrawRectangle(leftButtonBound.X, leftButtonBound.Y, leftButtonBound.Width, leftButtonBound.Height, rl.GetColor(int32(style[SpinnerPressedButtonBorderColor])))
			rl.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, rl.GetColor(int32(style[SpinnerPressedButtonInsideColor])))

			rl.DrawRectangle(rightButtonBound.X, rightButtonBound.Y, rightButtonBound.Width, rightButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			rl.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			rl.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(rl.MeasureText("+", textHeight))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerPressedSymbolColor])))
			rl.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(rl.MeasureText("-", textHeight))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))
		} else if buttonSide == 2 {
			rl.DrawRectangle(leftButtonBound.X, leftButtonBound.Y, leftButtonBound.Width, leftButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			rl.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			rl.DrawRectangle(rightButtonBound.X, rightButtonBound.Y, rightButtonBound.Width, rightButtonBound.Height, rl.GetColor(int32(style[SpinnerPressedButtonBorderColor])))
			rl.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, rl.GetColor(int32(style[SpinnerPressedButtonInsideColor])))

			rl.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(rl.MeasureText("+", textHeight))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))
			rl.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(rl.MeasureText("-", textHeight))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerPressedSymbolColor])))
		}

		rl.DrawRectangle(labelBoxBound.X, labelBoxBound.Y, labelBoxBound.Width, labelBoxBound.Height, rl.GetColor(int32(style[SpinnerLabelBorderColor])))
		rl.DrawRectangle(labelBoxBound.X+1, labelBoxBound.Y+1, labelBoxBound.Width-2, labelBoxBound.Height-2, rl.GetColor(int32(style[SpinnerLabelInsideColor])))

		rl.DrawText(fmt.Sprintf("%d", value), labelBoxBound.X+(labelBoxBound.Width/2-textWidth/2), labelBoxBound.Y+(labelBoxBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerPressedTextColor])))
		break
	default:
		break
	}

	return value
}
