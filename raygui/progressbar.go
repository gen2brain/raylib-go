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

	borderWidth := GetStyle32(ProgressbarBorderWidth)
	progressBar := InsetRectangle(b, borderWidth)              // backing rectangle
	progressWidth := int32(value * float32(progressBar.Width)) // how much should be replaced with progress
	progressValue := rl.RectangleInt32{progressBar.X, progressBar.Y, progressWidth, progressBar.Height}

	// Draw control
	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, GetStyleColor(ProgressbarBorderColor))
	rl.DrawRectangle(progressBar.X, progressBar.Y, progressBar.Width, progressBar.Height, GetStyleColor(ProgressbarInsideColor))
	rl.DrawRectangle(progressValue.X, progressValue.Y, progressValue.Width, progressValue.Height, GetStyleColor(ProgressbarProgressColor))
}
