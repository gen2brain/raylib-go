package raygui

import rl "github.com/gen2brain/raylib-go/raylib"

// Label - Label element, show text
func Label(bounds rl.Rectangle, text string) {
	LabelEx(bounds, text, rl.GetColor(int32(style[LabelTextColor])), rl.NewColor(0, 0, 0, 0), rl.NewColor(0, 0, 0, 0))
}

// LabelEx - Label element extended, configurable colors
func LabelEx(bounds rl.Rectangle, text string, textColor, border, inner rl.Color) {
	b := bounds.ToInt32()
	// Update control
	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(text, textHeight)

	if b.Width < textWidth {
		b.Width = textWidth + int32(style[LabelTextPadding])
	}
	if b.Height < textHeight {
		b.Height = textHeight + int32(style[LabelTextPadding])/2
	}

	// Draw control
	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, border)
	rl.DrawRectangle(b.X+int32(style[LabelBorderWidth]), b.Y+int32(style[LabelBorderWidth]), b.Width-(2*int32(style[LabelBorderWidth])), b.Height-(2*int32(style[LabelBorderWidth])), inner)
	rl.DrawText(text, b.X+((b.Width/2)-(textWidth/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, textColor)
}
