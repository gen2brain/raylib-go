package raygui

import rl "github.com/gen2brain/raylib-go/raylib"

// SliderBar - Slider Bar element, returns selected value
func SliderBar(bounds rl.Rectangle, value, minValue, maxValue float32) float32 {
	b := bounds.ToInt32()
	state := Normal

	mousePoint := rl.GetMousePosition()

	fixedValue := float32(0)
	fixedMinValue := float32(0)

	fixedValue = value - minValue
	maxValue = maxValue - minValue
	fixedMinValue = 0

	// Update control
	if fixedValue <= fixedMinValue {
		fixedValue = fixedMinValue
	} else if fixedValue >= maxValue {
		fixedValue = maxValue
	}

	sliderBar := rl.RectangleInt32{}

	sliderBar.X = b.X + int32(style[SliderBorderWidth])
	sliderBar.Y = b.Y + int32(style[SliderBorderWidth])
	sliderBar.Width = int32((fixedValue * (float32(b.Width) - 2*float32(style[SliderBorderWidth]))) / (maxValue - fixedMinValue))
	sliderBar.Height = b.Height - 2*int32(style[SliderBorderWidth])

	if rl.CheckCollisionPointRec(mousePoint, bounds) {
		state = Focused

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			state = Pressed

			sliderBar.Width = (int32(mousePoint.X) - b.X - int32(style[SliderBorderWidth]))

			if int32(mousePoint.X) <= (b.X + int32(style[SliderBorderWidth])) {
				sliderBar.Width = 0
			} else if int32(mousePoint.X) >= (b.X + b.Width - int32(style[SliderBorderWidth])) {
				sliderBar.Width = b.Width - 2*int32(style[SliderBorderWidth])
			}
		}
	} else {
		state = Normal
	}

	fixedValue = (float32(sliderBar.Width) * (maxValue - fixedMinValue)) / (float32(b.Width) - 2*float32(style[SliderBorderWidth]))

	// Draw control
	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(uint(style[SliderbarBorderColor])))
	rl.DrawRectangle(b.X+int32(style[SliderBorderWidth]), b.Y+int32(style[SliderBorderWidth]), b.Width-(2*int32(style[SliderBorderWidth])), b.Height-(2*int32(style[SliderBorderWidth])), rl.GetColor(uint(style[SliderbarInsideColor])))

	switch state {
	case Normal:
		rl.DrawRectangle(sliderBar.X, sliderBar.Y, sliderBar.Width, sliderBar.Height, rl.GetColor(uint(style[SliderbarDefaultColor])))
		break
	case Focused:
		rl.DrawRectangle(sliderBar.X, sliderBar.Y, sliderBar.Width, sliderBar.Height, rl.GetColor(uint(style[SliderbarHoverColor])))
		break
	case Pressed:
		rl.DrawRectangle(sliderBar.X, sliderBar.Y, sliderBar.Width, sliderBar.Height, rl.GetColor(uint(style[SliderbarActiveColor])))
		break
	default:
		break
	}

	if minValue < 0 && maxValue > 0 {
		rl.DrawRectangle((b.X+int32(style[SliderBorderWidth]))-int32(minValue*(float32(b.Width-(int32(style[SliderBorderWidth])*2))/maxValue)), sliderBar.Y, 1, sliderBar.Height, rl.GetColor(uint(style[SliderbarZeroLineColor])))
	}

	return fixedValue + minValue
}
