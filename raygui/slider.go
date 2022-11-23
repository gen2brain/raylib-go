package raygui

import rl "github.com/gen2brain/raylib-go/raylib"

// Slider - Slider element, returns selected value
func Slider(bounds rl.Rectangle, value, minValue, maxValue float32) float32 {
	b := bounds.ToInt32()
	sliderPos := float32(0)
	state := Normal

	buttonTravelDistance := float32(0)
	mousePoint := rl.GetMousePosition()

	// Update control
	if value < minValue {
		value = minValue
	} else if value >= maxValue {
		value = maxValue
	}

	sliderPos = (value - minValue) / (maxValue - minValue)

	sliderButton := rl.RectangleInt32{}
	sliderButton.Width = (b.Width-(2*int32(style[SliderButtonBorderWidth])))/10 - 8
	sliderButton.Height = b.Height - (2 * int32(style[SliderBorderWidth]+2*style[SliderButtonBorderWidth]))

	sliderButtonMinPos := b.X + int32(style[SliderBorderWidth]) + int32(style[SliderButtonBorderWidth])
	sliderButtonMaxPos := b.X + b.Width - (int32(style[SliderBorderWidth]) + int32(style[SliderButtonBorderWidth]) + sliderButton.Width)

	buttonTravelDistance = float32(sliderButtonMaxPos - sliderButtonMinPos)

	sliderButton.X = b.X + int32(style[SliderBorderWidth]) + int32(style[SliderButtonBorderWidth]) + int32(sliderPos*buttonTravelDistance)
	sliderButton.Y = b.Y + int32(style[SliderBorderWidth]) + int32(style[SliderButtonBorderWidth])

	if rl.CheckCollisionPointRec(mousePoint, bounds) {
		state = Focused

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			state = Pressed
		}

		if state == Pressed && rl.IsMouseButtonDown(rl.MouseLeftButton) {
			sliderButton.X = int32(mousePoint.X) - sliderButton.Width/2

			if sliderButton.X <= sliderButtonMinPos {
				sliderButton.X = sliderButtonMinPos
			} else if sliderButton.X >= sliderButtonMaxPos {
				sliderButton.X = sliderButtonMaxPos
			}

			sliderPos = float32(sliderButton.X-sliderButtonMinPos) / buttonTravelDistance
		}
	} else {
		state = Normal
	}

	// Draw control
	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(uint(style[SliderBorderColor])))
	rl.DrawRectangle(b.X+int32(style[SliderBorderWidth]), b.Y+int32(style[SliderBorderWidth]), b.Width-(2*int32(style[SliderBorderWidth])), b.Height-(2*int32(style[SliderBorderWidth])), rl.GetColor(uint(style[SliderInsideColor])))

	switch state {
	case Normal:
		rl.DrawRectangle(sliderButton.X, sliderButton.Y, sliderButton.Width, sliderButton.Height, rl.GetColor(uint(style[SliderDefaultColor])))
		break
	case Focused:
		rl.DrawRectangle(sliderButton.X, sliderButton.Y, sliderButton.Width, sliderButton.Height, rl.GetColor(uint(style[SliderHoverColor])))
		break
	case Pressed:
		rl.DrawRectangle(sliderButton.X, sliderButton.Y, sliderButton.Width, sliderButton.Height, rl.GetColor(uint(style[SliderActiveColor])))
		break
	default:
		break
	}

	return minValue + (maxValue-minValue)*sliderPos
}
