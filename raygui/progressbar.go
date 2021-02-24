package raygui

import "github.com/gen2brain/raylib-go/raylib"

// ProgressBar - Progress Bar element, shows current progress value
func ProgressBar(bounds rl.Rectangle, value float32) {
	b := bounds.ToInt32()
	if value > 1.0 {
		value = 1.0
	} else if value < 0.0 {
		value = 0.0
	}

	progressBar := rl.RectangleInt32{b.X + int32(style[ProgressbarBorderWidth]), b.Y + int32(style[ProgressbarBorderWidth]), b.Width - (int32(style[ProgressbarBorderWidth]) * 2), b.Height - (int32(style[ProgressbarBorderWidth]) * 2)}
	progressValue := rl.RectangleInt32{b.X + int32(style[ProgressbarBorderWidth]), b.Y + int32(style[ProgressbarBorderWidth]), int32(value * float32(b.Width-int32(style[ProgressbarBorderWidth])*2)), b.Height - (int32(style[ProgressbarBorderWidth]) * 2)}

	// Draw control
	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ProgressbarBorderColor])))
	rl.DrawRectangle(progressBar.X, progressBar.Y, progressBar.Width, progressBar.Height, rl.GetColor(int32(style[ProgressbarInsideColor])))
	rl.DrawRectangle(progressValue.X, progressValue.Y, progressValue.Width, progressValue.Height, rl.GetColor(int32(style[ProgressbarProgressColor])))
}
