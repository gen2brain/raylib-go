package raygui

import "github.com/gen2brain/raylib-go/raylib"

// ToggleGroup - Toggle Group element, returns toggled button index
func ToggleGroup(bounds rl.Rectangle, toggleText []string, active int) int {
	padding := float32(style[TogglegroupPadding])
	for i := 0; i < len(toggleText); i++ {
		if i == active {
			ToggleButton(rl.NewRectangle(bounds.X+float32(i)*(bounds.Width+padding), bounds.Y, bounds.Width, bounds.Height), toggleText[i], true)
		} else if ToggleButton(rl.NewRectangle(bounds.X+float32(i)*(bounds.Width+padding), bounds.Y, bounds.Width, bounds.Height), toggleText[i], false) {
			active = i
		}
	}

	return active
}
