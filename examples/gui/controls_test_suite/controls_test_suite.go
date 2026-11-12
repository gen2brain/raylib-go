package main

import (
	"fmt"
	"path/filepath"
	"strings"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

/*******************************************************************************************
*
*   raygui - controls test suite
*
*   DEPENDENCIES:
*       raylib 6.2 - Windowing/input management and drawing.
*       raygui 5.1-dev - Immediate-mode GUI controls.
*
*   LICENSE: zlib/libpng
*
*   Copyright (c) 2016-2026 Ramon Santamaria (@raysan5)
*
**********************************************************************************************/

var styles = []string{"", "jungle", "lavanda", "dark", "bluish", "cyber", "terminal", "candy", "cherry", "ashes", "enefete", "sunny", "amber", "genesis", "brick", "pocket", "turbo", "wisteria", "advance", "rltech"}

func main() {
	const screenWidth = 960
	const screenHeight = 580

	rl.InitWindow(screenWidth, screenHeight, "raygui - controls test suite")
	rl.SetExitKey(0)

	var tabActive int32

	var dropdownBox000Active int32
	dropDown000EditMode := false

	var dropdownBox001Active int32
	dropDown001EditMode := false

	var spinner001Value int32
	spinnerEditMode := false

	var valueBox002Value int32
	valueBoxEditMode := false

	textBoxText := "Text box"
	textBoxEditMode := false

	textBoxMultiText := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.\n\nDuis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.\n\nThisisastringlongerthanexpectedwithoutspacestotestcharbreaksforthosecases,checkingifworkingasexpected.\n\nExcepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	textBoxMultiEditMode := false

	var listViewScrollIndex int32
	listViewActive := int32(-1)

	var listViewExScrollIndex int32
	listViewExActive := int32(2)
	listViewExFocus := int32(-1)
	listViewExList := []string{"This", "is", "a", "list view", "with", "disable", "elements", "amazing!"}

	colorPickerValue := rl.Red

	sliderValue := float32(50.0)
	sliderBarValue := float32(60)
	progressValue := float32(0.1)

	forceSquaredChecked := false

	alphaValue := float32(0.5)

	var visualStyleActive int32
	var prevVisualStyleActive int32

	var toggleGroupActive int32
	var toggleSliderActive int32

	var viewScroll rl.Vector2

	exitWindow := false
	showMessageBox := false

	textInput := ""
	showTextInputBox := false

	alpha := float32(1.0)

	rl.SetTargetFPS(60)

	for !exitWindow {
		exitWindow = rl.WindowShouldClose()

		if rl.IsKeyPressed(rl.KeyEscape) {
			showMessageBox = !showMessageBox
		}

		if rl.IsKeyDown(rl.KeyLeftControl) && rl.IsKeyPressed(rl.KeyS) {
			showTextInputBox = true
		}

		if rl.IsFileDropped() {
			droppedFiles := rl.LoadDroppedFiles()

			if len(droppedFiles) > 0 && strings.EqualFold(filepath.Ext(droppedFiles[0]), ".rgs") {
				gui.LoadStyle(droppedFiles[0])
			}

			rl.UnloadDroppedFiles()
		}

		if alpha < 0.0 {
			alpha = 0.0
		}
		if rl.IsKeyPressed(rl.KeySpace) {
			alpha = 1.0
		}

		gui.SetAlpha(alpha)

		if rl.IsKeyPressed(rl.KeyLeft) {
			progressValue -= 0.1
		} else if rl.IsKeyPressed(rl.KeyRight) {
			progressValue += 0.1
		}
		if progressValue > 1.0 {
			progressValue = 1.0
		} else if progressValue < 0.0 {
			progressValue = 0.0
		}

		if visualStyleActive != prevVisualStyleActive {
			gui.LoadStyleDefault()

			if visualStyleActive > 0 && int(visualStyleActive) < len(styles) {
				gui.LoadStyle(filepath.Join("styles", "style_"+styles[visualStyleActive]+".rgs"))
			}

			gui.SetStyle(gui.LABEL, gui.TEXT_ALIGNMENT, gui.TEXT_ALIGN_LEFT)

			prevVisualStyleActive = visualStyleActive
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

		if dropDown000EditMode || dropDown001EditMode {
			gui.Lock()
		}
		if showTextInputBox {
			gui.Lock()
		}

		gui.SetStyle(gui.TABBAR, gui.TAB_ITEMS_WIDTH, 140)
		gui.TabBar(rl.Rectangle{X: 0, Y: 8, Width: float32(rl.GetScreenWidth()), Height: 26}, "#176#TAB_01;#177#TAB_02;#178#TAB_03;#179#TAB_04;#180#TAB_05;#181#TAB_06", nil, &tabActive)

		gui.CheckBox(rl.Rectangle{X: 25, Y: 108 + 20, Width: 15, Height: 15}, "FORCE CHECK!", &forceSquaredChecked)

		gui.SetStyle(gui.TEXTBOX, gui.TEXT_ALIGNMENT, gui.TEXT_ALIGN_CENTER)
		if gui.Spinner(rl.Rectangle{X: 25, Y: 135 + 20, Width: 125, Height: 30}, "", &spinner001Value, 0, 100, spinnerEditMode) != 0 {
			spinnerEditMode = !spinnerEditMode
		}
		if gui.ValueBox(rl.Rectangle{X: 25, Y: 175 + 20, Width: 125, Height: 30}, "", &valueBox002Value, 0, 100, valueBoxEditMode) != 0 {
			valueBoxEditMode = !valueBoxEditMode
		}
		gui.SetStyle(gui.TEXTBOX, gui.TEXT_ALIGNMENT, gui.TEXT_ALIGN_LEFT)
		if gui.TextBox(rl.Rectangle{X: 25, Y: 215 + 20, Width: 125, Height: 30}, &textBoxText, 64, textBoxEditMode) != 0 {
			textBoxEditMode = !textBoxEditMode
		}

		gui.SetStyle(gui.BUTTON, gui.TEXT_ALIGNMENT, gui.TEXT_ALIGN_CENTER)

		if gui.Button(rl.Rectangle{X: 25, Y: 255 + 20, Width: 125, Height: 30}, gui.IconText(gui.ICON_FILE_SAVE, "Save File")) != 0 {
			showTextInputBox = true
		}

		gui.GroupBox(rl.Rectangle{X: 25, Y: 310 + 20, Width: 125, Height: 150}, "STATES")
		gui.SetState(gui.STATE_NORMAL)
		gui.Button(rl.Rectangle{X: 30, Y: 320 + 20, Width: 115, Height: 30}, "NORMAL")
		gui.SetState(gui.STATE_FOCUSED)
		gui.Button(rl.Rectangle{X: 30, Y: 355 + 20, Width: 115, Height: 30}, "FOCUSED")
		gui.SetState(gui.STATE_PRESSED)
		gui.Button(rl.Rectangle{X: 30, Y: 390 + 20, Width: 115, Height: 30}, "#15#PRESSED")
		gui.SetState(gui.STATE_DISABLED)
		gui.Button(rl.Rectangle{X: 30, Y: 425 + 20, Width: 115, Height: 30}, "DISABLED")
		gui.SetState(gui.STATE_NORMAL)

		gui.SetStyle(gui.COMBOBOX, gui.COMBO_BUTTON_WIDTH, 40)
		gui.ComboBox(rl.Rectangle{X: 25, Y: 480 + 20, Width: 125, Height: 30},
			"default;Jungle;Lavanda;Dark;Bluish;Cyber;Terminal;Candy;Cherry;Ashes;Enefete;Sunny;Amber;Genesis;Brick;Pocket;Turbo;Wisteria;Advance;RLTech", &visualStyleActive)

		if dropDown000EditMode || dropDown001EditMode {
			gui.Unlock()
		}
		if showTextInputBox {
			gui.Lock()
		}

		gui.SetStyle(gui.DROPDOWNBOX, gui.TEXT_PADDING, 4)
		gui.SetStyle(gui.DROPDOWNBOX, gui.TEXT_ALIGNMENT, gui.TEXT_ALIGN_LEFT)
		if gui.DropdownBox(rl.Rectangle{X: 25, Y: 65 + 20, Width: 125, Height: 30}, "#01#ONE;#02#TWO;#03#THREE;#04#FOUR", &dropdownBox001Active, dropDown001EditMode) != 0 {
			dropDown001EditMode = !dropDown001EditMode
		}
		gui.SetStyle(gui.DROPDOWNBOX, gui.TEXT_ALIGNMENT, gui.TEXT_ALIGN_CENTER)
		gui.SetStyle(gui.DROPDOWNBOX, gui.TEXT_PADDING, 0)

		if gui.DropdownBox(rl.Rectangle{X: 25, Y: 25 + 20, Width: 125, Height: 30}, "#27#ONE;#28#TWO;#29#THREE", &dropdownBox000Active, dropDown000EditMode) != 0 {
			dropDown000EditMode = !dropDown000EditMode
		}

		gui.ListView(rl.Rectangle{X: 165, Y: 25 + 20, Width: 140, Height: 124}, "Charmander;Bulbasaur;#18#Squirtle;Pikachu;Eevee;Pidgey", &listViewScrollIndex, &listViewActive)
		gui.ListViewEx(rl.Rectangle{X: 165, Y: 162 + 20, Width: 140, Height: 184}, listViewExList, &listViewExScrollIndex, &listViewExActive, &listViewExFocus)
		gui.SetStyle(gui.LISTVIEW, gui.LIST_ITEMS_BORDER_NORMAL, 0)

		gui.ToggleGroup(rl.Rectangle{X: 165, Y: 360 + 20, Width: 140, Height: 24}, "#1#ONE\n#3#TWO\n#8#THREE\n#23#", &toggleGroupActive)
		gui.SetStyle(gui.SLIDER, gui.SLIDER_PADDING, 2)
		gui.ToggleSlider(rl.Rectangle{X: 165, Y: 480 + 20, Width: 140, Height: 30}, "ON;OFF", &toggleSliderActive)
		gui.SetStyle(gui.SLIDER, gui.SLIDER_PADDING, 0)

		gui.Panel(rl.Rectangle{X: 320, Y: 25 + 20, Width: 225, Height: 140}, "Panel Info")
		gui.ColorPicker(rl.Rectangle{X: 320, Y: 185 + 20, Width: 196, Height: 192}, "", &colorPickerValue)

		gui.Slider(rl.Rectangle{X: 355, Y: 400 + 20, Width: 165, Height: 20}, "TEST", fmt.Sprintf("%2.2f", sliderValue), &sliderValue, -50, 100)
		gui.SliderBar(rl.Rectangle{X: 320, Y: 430 + 20, Width: 200, Height: 20}, "", fmt.Sprintf("%d", int32(sliderBarValue)), &sliderBarValue, 0, 100)

		gui.ProgressBar(rl.Rectangle{X: 320, Y: 460 + 20, Width: 200, Height: 20}, "", fmt.Sprintf("%d%%", int32(progressValue*100)), &progressValue, 0.0, 1.0)
		gui.Enable()

		var view rl.Rectangle
		gui.ScrollPanel(rl.Rectangle{X: 560, Y: 25 + 20, Width: 102, Height: 354}, "", rl.Rectangle{X: 560, Y: 25, Width: 300, Height: 1200}, &viewScroll, &view)

		var mouseCell rl.Vector2
		gui.Grid(rl.Rectangle{X: 560, Y: 25 + 180 + 195 + 20, Width: 100, Height: 120}, "", 20, 3, &mouseCell)

		gui.ColorBarAlpha(rl.Rectangle{X: 320, Y: 490 + 20, Width: 200, Height: 30}, "", &alphaValue)

		gui.SetStyle(gui.DEFAULT, gui.TEXT_ALIGNMENT_VERTICAL, gui.TEXT_ALIGN_TOP)
		gui.SetStyle(gui.DEFAULT, gui.TEXT_WRAP_MODE, gui.TEXT_WRAP_WORD)
		if gui.TextBox(rl.Rectangle{X: 678, Y: 25 + 20, Width: 258, Height: 492}, &textBoxMultiText, 1024, textBoxMultiEditMode) != 0 {
			textBoxMultiEditMode = !textBoxMultiEditMode
		}
		gui.SetStyle(gui.DEFAULT, gui.TEXT_WRAP_MODE, gui.TEXT_WRAP_NONE)
		gui.SetStyle(gui.DEFAULT, gui.TEXT_ALIGNMENT_VERTICAL, gui.TEXT_ALIGN_MIDDLE)

		gui.SetStyle(gui.DEFAULT, gui.TEXT_ALIGNMENT, gui.TEXT_ALIGN_LEFT)
		gui.StatusBar(rl.Rectangle{X: 0, Y: float32(rl.GetScreenHeight()) - 20, Width: float32(rl.GetScreenWidth()), Height: 20}, "This is a status bar")
		gui.SetStyle(gui.DEFAULT, gui.TEXT_ALIGNMENT, gui.TEXT_ALIGN_CENTER)

		if showMessageBox {
			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.RayWhite, 0.8))
			btnActive := int32(-1)
			gui.MessageBox(rl.Rectangle{X: float32(rl.GetScreenWidth())/2 - 125, Y: float32(rl.GetScreenHeight())/2 - 50, Width: 250, Height: 100},
				gui.IconText(gui.ICON_EXIT, "Close Window"), "Do you really want to exit?", "Yes;No", &btnActive)

			if btnActive == 0 || btnActive == 2 {
				showMessageBox = false
			} else if btnActive == 1 {
				exitWindow = true
			}
		}

		if showTextInputBox {
			gui.Unlock()

			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.RayWhite, 0.8))
			btnActive := int32(-1)
			gui.TextInputBox(rl.Rectangle{X: float32(rl.GetScreenWidth())/2 - 120, Y: float32(rl.GetScreenHeight())/2 - 60, Width: 240, Height: 140},
				gui.IconText(gui.ICON_FILE_SAVE, "Save file as..."), "Introduce output file name:", &textInput, 255, "Ok;Cancel", &btnActive, nil)

			if btnActive == 0 || btnActive == 1 || btnActive == 2 {
				showTextInputBox = false
				textInput = ""
			}
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
