package raygui

import rl "github.com/gen2brain/raylib-go/raylib"

// Label - Label element, show text
func Label(bounds rl.Rectangle, text string) {
	LabelEx(bounds, text, rl.GetColor(uint(style[LabelTextColor])), rl.NewColor(0, 0, 0, 0), rl.NewColor(0, 0, 0, 0))
}

// LabelEx - Label element extended, configurable colors
func LabelEx(bounds rl.Rectangle, text string, textColor, border, inner rl.Color) {
	textHeight := GetStyle32(GlobalTextFontsize)
	textWidth := rl.MeasureText(text, textHeight)

	ConstrainRectangle(&bounds, textWidth, textWidth+GetStyle32(LabelTextPadding), textHeight, textHeight+GetStyle32(LabelTextPadding)/2)

	// Draw control
	b := bounds.ToInt32()
	DrawBorderedRectangle(b, GetStyle32(LabelBorderWidth), border, inner)
	rl.DrawText(text, b.X+((b.Width/2)-(textWidth/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, textColor)
}
