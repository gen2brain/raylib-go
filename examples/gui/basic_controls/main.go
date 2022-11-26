package main

import (
	"fmt"
	"strings"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagVsyncHint)

	rl.InitWindow(screenWidth, screenHeight, "raylib [gui] example - basic controls")

	buttonToggle := true
	buttonClicked := false
	checkboxChecked := false

	spinnerValue := int32(5)
	sliderValue := float32(10)
	sliderBarValue := float32(70)
	progressValue := float32(0.5)

	comboActive := int32(0)
	comboLastActive := int32(0)
	toggleActive := int32(0)

	toggleText := "Item 1;Item 2;Item 3"
	comboText := []string{
		"ashes", "bluish", "candy", "cherry", "cyber", "dark",
		"default", "enefete", "jungle", "lavanda", "sunny", "terminal",
	}
	comboList := strings.Join(comboText, ";")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if buttonClicked {
			progressValue += 0.1
			if progressValue >= 1.1 {
				progressValue = 0.0
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Beige)

		gui.Label(rl.NewRectangle(50, 50, 80, 20), "Label")

		buttonClicked = gui.Button(rl.NewRectangle(50, 70, 80, 40), "Button")

		checkboxChecked = gui.CheckBox(rl.NewRectangle(50, 140, 20, 20), "Checkbox", checkboxChecked)

		gui.ProgressBar(rl.NewRectangle(50, 210, 200, 20),
			fmt.Sprintf("%.1f", progressValue),
			"ProgressBar",
			progressValue, 0, 1)

		sliderValue = gui.Slider(rl.NewRectangle(50, 280, 200, 20), "Slider", "", sliderValue, 0, 100)
		gui.Label(rl.NewRectangle(200+50+5, 280, 20, 20), fmt.Sprintf("%.0f", sliderValue))

		buttonToggle = gui.Toggle(rl.NewRectangle(50, 350, 100, 40), "ButtonToggle", buttonToggle)

		toggleActive = gui.ToggleGroup(rl.NewRectangle(500, 70, 60, 30), toggleText, toggleActive)

		sliderBarValue = gui.SliderBar(rl.NewRectangle(500, 140, 200, 20), "SliderBar", "", sliderBarValue, 0, 100)
		gui.Label(rl.NewRectangle(500+200+5, 140, 20, 20), fmt.Sprintf("%.0f", sliderBarValue))

		_ = gui.Spinner(rl.NewRectangle(500, 210, 200, 20), "Spinner", &spinnerValue, 0, 100, true)

		gui.Label(rl.NewRectangle(500, 260, 200, 20), "ComboBox")
		comboActive = gui.ComboBox(rl.NewRectangle(500, 280, 200, 20), comboList, comboActive)

		if comboLastActive != comboActive {
			ch := comboText[comboActive] // choosed name
			filename := fmt.Sprintf("styles/%s/%s.rgs", ch, ch)
			gui.LoadStyle(filename)
			comboLastActive = comboActive
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
