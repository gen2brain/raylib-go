package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagVsyncHint)

	rl.InitWindow(screenWidth, screenHeight, "raylib [gui] example - basic controls")

	buttonToggle := true
	buttonClicked := false
	checkboxChecked := false

	spinnerValue := 5
	sliderValue := float32(10)
	sliderBarValue := float32(70)
	progressValue := float32(0.5)

	comboActive := 0
	comboLastActive := 0
	toggleActive := 0

	toggleText := []string{"Item 1", "Item 2", "Item 3"}
	comboText := []string{"default_light", "default_dark", "hello_kitty", "monokai", "obsidian", "solarized_light", "solarized", "zahnrad"}

	var inputText string

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

		raygui.Label(rl.NewRectangle(50, 50, 80, 20), "Label")

		buttonClicked = raygui.Button(rl.NewRectangle(50, 70, 80, 40), "Button")

		raygui.Label(rl.NewRectangle(70, 140, 20, 20), "Checkbox")
		checkboxChecked = raygui.CheckBox(rl.NewRectangle(50, 140, 20, 20), checkboxChecked)

		raygui.Label(rl.NewRectangle(50, 190, 200, 20), "ProgressBar")
		raygui.ProgressBar(rl.NewRectangle(50, 210, 200, 20), progressValue)
		raygui.Label(rl.NewRectangle(200+50+5, 210, 20, 20), fmt.Sprintf("%.1f", progressValue))

		raygui.Label(rl.NewRectangle(50, 260, 200, 20), "Slider")
		sliderValue = raygui.Slider(rl.NewRectangle(50, 280, 200, 20), sliderValue, 0, 100)
		raygui.Label(rl.NewRectangle(200+50+5, 280, 20, 20), fmt.Sprintf("%.0f", sliderValue))

		buttonToggle = raygui.ToggleButton(rl.NewRectangle(50, 350, 100, 40), "ToggleButton", buttonToggle)

		raygui.Label(rl.NewRectangle(500, 50, 200, 20), "ToggleGroup")
		toggleActive = raygui.ToggleGroup(rl.NewRectangle(500, 70, 60, 30), toggleText, toggleActive)

		raygui.Label(rl.NewRectangle(500, 120, 200, 20), "SliderBar")
		sliderBarValue = raygui.SliderBar(rl.NewRectangle(500, 140, 200, 20), sliderBarValue, 0, 100)
		raygui.Label(rl.NewRectangle(500+200+5, 140, 20, 20), fmt.Sprintf("%.0f", sliderBarValue))

		raygui.Label(rl.NewRectangle(500, 190, 200, 20), "Spinner")
		spinnerValue = raygui.Spinner(rl.NewRectangle(500, 210, 200, 20), spinnerValue, 0, 100)

		raygui.Label(rl.NewRectangle(500, 260, 200, 20), "ComboBox")
		comboActive = raygui.ComboBox(rl.NewRectangle(500, 280, 200, 20), comboText, comboActive)

		if comboLastActive != comboActive {
			raygui.LoadGuiStyle(fmt.Sprintf("styles/%s.style", comboText[comboActive]))
			comboLastActive = comboActive
		}

		raygui.Label(rl.NewRectangle(500, 330, 200, 20), "TextBox")
		inputText = raygui.TextBox(rl.NewRectangle(500, 350, 200, 20), inputText)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
