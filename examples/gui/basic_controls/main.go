package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.SetConfigFlags(raylib.FlagVsyncHint)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [gui] example - basic controls")

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

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if buttonClicked {
			progressValue += 0.1
			if progressValue >= 1.1 {
				progressValue = 0.0
			}
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.Beige)

		raygui.Label(raylib.NewRectangle(50, 50, 80, 20), "Label")

		buttonClicked = raygui.Button(raylib.NewRectangle(50, 70, 80, 40), "Button")

		raygui.Label(raylib.NewRectangle(70, 140, 20, 20), "Checkbox")
		checkboxChecked = raygui.CheckBox(raylib.NewRectangle(50, 140, 20, 20), checkboxChecked)

		raygui.Label(raylib.NewRectangle(50, 190, 200, 20), "ProgressBar")
		raygui.ProgressBar(raylib.NewRectangle(50, 210, 200, 20), progressValue)
		raygui.Label(raylib.NewRectangle(200+50+5, 210, 20, 20), fmt.Sprintf("%.1f", progressValue))

		raygui.Label(raylib.NewRectangle(50, 260, 200, 20), "Slider")
		sliderValue = raygui.Slider(raylib.NewRectangle(50, 280, 200, 20), sliderValue, 0, 100)
		raygui.Label(raylib.NewRectangle(200+50+5, 280, 20, 20), fmt.Sprintf("%.0f", sliderValue))

		buttonToggle = raygui.ToggleButton(raylib.NewRectangle(50, 350, 100, 40), "ToggleButton", buttonToggle)

		raygui.Label(raylib.NewRectangle(500, 50, 200, 20), "ToggleGroup")
		toggleActive = raygui.ToggleGroup(raylib.NewRectangle(500, 70, 60, 30), toggleText, toggleActive)

		raygui.Label(raylib.NewRectangle(500, 120, 200, 20), "SliderBar")
		sliderBarValue = raygui.SliderBar(raylib.NewRectangle(500, 140, 200, 20), sliderBarValue, 0, 100)
		raygui.Label(raylib.NewRectangle(500+200+5, 140, 20, 20), fmt.Sprintf("%.0f", sliderBarValue))

		raygui.Label(raylib.NewRectangle(500, 190, 200, 20), "Spinner")
		spinnerValue = raygui.Spinner(raylib.NewRectangle(500, 210, 200, 20), spinnerValue, 0, 100)

		raygui.Label(raylib.NewRectangle(500, 260, 200, 20), "ComboBox")
		comboActive = raygui.ComboBox(raylib.NewRectangle(500, 280, 200, 20), comboText, comboActive)

		if comboLastActive != comboActive {
			raygui.LoadGuiStyle(fmt.Sprintf("styles/%s.style", comboText[comboActive]))
			comboLastActive = comboActive
		}

		raygui.Label(raylib.NewRectangle(500, 330, 200, 20), "TextBox")
		inputText = raygui.TextBox(raylib.NewRectangle(500, 350, 200, 20), inputText)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
