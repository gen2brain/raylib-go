package main

import (
	"fmt"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

/*******************************************************************************************
*
*   raygui - controls test suite
*
*   TEST CONTROLS:
*       - gui.DropdownBox()
*       - gui.CheckBox()
*       - gui.Spinner()
*       - gui.ValueBox()
*       - gui.TextBox()
*       - gui.Button()
*       - gui.ComboBox()
*       - gui.ListView()
*       - gui.ToggleGroup()
*       - gui.ColorPicker()
*       - gui.Slider()
*       - gui.SliderBar()
*       - gui.ProgressBar()
*       - gui.ColorBarAlpha()
*       - gui.ScrollPanel()
*
*
*   DEPENDENCIES:
*       raylib 4.0 - Windowing/input management and drawing.
*       raygui 3.2 - Immediate-mode GUI controls.
*
*   COMPILATION (Windows - MinGW):
*       gcc -o $(NAME_PART).exe $(FILE_NAME) -I../../src -lraylib -lopengl32 -lgdi32 -std=c99
*
*   LICENSE: zlib/libpng
*
*   Copyright (c) 2016-2022 Ramon Santamaria (@raysan5)
*
**********************************************************************************************/

//#define RAYGUI_CUSTOM_ICONS     // It requires providing gui_icons.h in the same directory
//#include "gui_icons.h"          // External icons data provided, it can be generated with rGuiIcons tool

// ------------------------------------------------------------------------------------
// Program main entry point
// ------------------------------------------------------------------------------------
func main() {
	// Initialization
	//---------------------------------------------------------------------------------------
	const (
		screenWidth  = 690
		screenHeight = 560
	)

	rl.InitWindow(screenWidth, screenHeight, "raygui - controls test suite")
	rl.SetExitKey(0)

	// GUI controls initialization
	//----------------------------------------------------------------------------------
	var (
		dropdownBox000Active int32 = 0
		dropDown000EditMode  bool  = false

		dropdownBox001Active int32 = 0
		dropDown001EditMode  bool  = false

		spinner001Value int32 = 0
		spinnerEditMode bool  = false

		valueBox002Value int32 = 0
		valueBoxEditMode bool  = false

		textBoxText          = "Text box"
		textBoxEditMode bool = false

		listViewScrollIndex int32 = 0
		listViewActive      int32 = -1

		listViewExScrollIndex int32 = 0
		listViewExActive      int32 = 2
		listViewExFocus       int32 = -1
		listViewExList              = []string{"This", "is", "a", "list view", "with", "disable", "elements", "amazing!"}

		colorPickerValue = rl.Red

		sliderValue    float32 = 50
		sliderBarValue float32 = 60
		progressValue  float32 = 0.4

		forceSquaredChecked bool = false

		alphaValue float32 = 0.5

		comboBoxActive int32 = 1

		toggleGroupActive int32 = 0

		viewScroll = rl.Vector2{0, 0}

		//----------------------------------------------------------------------------------

		// Custom GUI font loading
		//Font font = LoadFontEx("fonts/rainyhearts16.ttf", 12, 0, 0);
		//GuiSetFont(font);

		exitWindow     bool = false
		showMessageBox bool = false

		textInput        string
		showTextInputBox bool = false

		// TODO textInputFileName string
	)

	rl.SetTargetFPS(60)
	//--------------------------------------------------------------------------------------

	// Main game loop
	for !exitWindow { // Detect window close button or ESC key
		// Update
		//----------------------------------------------------------------------------------
		exitWindow = rl.WindowShouldClose()

		if rl.IsKeyPressed(rl.KeyEscape) {
			showMessageBox = !showMessageBox
		}

		if rl.IsKeyDown(rl.KeyLeftControl) && rl.IsKeyPressed(rl.KeyS) {
			showTextInputBox = true
		}

		// TODO if rl.IsFileDropped() {
		// TODO var droppedFiles gui.FilePathList = rl.LoadDroppedFiles()
		// TODO if (droppedFiles.count > 0) && rl.IsFileExtension(droppedFiles.paths[0], ".rgs") {
		// TODO 	gui.LoadStyle(droppedFiles.paths[0])
		// TODO }
		// TODO rl.UnloadDroppedFiles(droppedFiles) // Clear internal buffers
		// TODO }
		//----------------------------------------------------------------------------------

		// Draw
		//----------------------------------------------------------------------------------
		rl.BeginDrawing()

		rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

		// raygui: controls drawing
		//----------------------------------------------------------------------------------
		if dropDown000EditMode || dropDown001EditMode {
			gui.Lock()
		} else if !dropDown000EditMode && !dropDown001EditMode {
			gui.Unlock()
		}
		//GuiDisable();

		// First GUI column
		//GuiSetStyle(CHECKBOX, TEXT_ALIGNMENT, TEXT_ALIGN_LEFT);
		forceSquaredChecked = gui.CheckBox(rl.Rectangle{25, 108, 15, 15}, "FORCE CHECK!", forceSquaredChecked)

		gui.SetStyle(gui.TEXTBOX, gui.TEXT_ALIGNMENT, gui.TEXT_ALIGN_CENTER)
		//GuiSetStyle(VALUEBOX, TEXT_ALIGNMENT, TEXT_ALIGN_LEFT);
		gui.Spinner(rl.Rectangle{25, 135, 125, 30}, "", &spinner001Value, 0, 100, spinnerEditMode)

		if gui.ValueBox(rl.Rectangle{25, 175, 125, 30}, "", &valueBox002Value, 0, 100, valueBoxEditMode) {
			valueBoxEditMode = !valueBoxEditMode
		}
		gui.SetStyle(gui.TEXTBOX, gui.TEXT_ALIGNMENT, int64(gui.TEXT_ALIGN_LEFT))
		if gui.TextBox(rl.Rectangle{25, 215, 125, 30}, &textBoxText, 64, textBoxEditMode) {
			textBoxEditMode = !textBoxEditMode
		}

		gui.SetStyle(gui.BUTTON, gui.TEXT_ALIGNMENT, gui.TEXT_ALIGN_CENTER)

		if gui.Button(rl.Rectangle{25, 255, 125, 30}, gui.IconText(gui.ICON_FILE_SAVE, "Save File")) {
			showTextInputBox = true
		}

		gui.GroupBox(rl.Rectangle{25, 310, 125, 150}, "STATES")
		//GuiLock();
		gui.SetState(gui.STATE_NORMAL)
		if gui.Button(rl.Rectangle{30, 320, 115, 30}, "NORMAL") {
		}
		gui.SetState(gui.STATE_FOCUSED)
		if gui.Button(rl.Rectangle{30, 355, 115, 30}, "FOCUSED") {
		}
		gui.SetState(gui.STATE_PRESSED)
		if gui.Button(rl.Rectangle{30, 390, 115, 30}, "#15#PRESSED") {
		}
		gui.SetState(gui.STATE_DISABLED)
		if gui.Button(rl.Rectangle{30, 425, 115, 30}, "DISABLED") {
		}
		gui.SetState(gui.STATE_NORMAL)
		//GuiUnlock();

		comboBoxActive = gui.ComboBox(rl.Rectangle{25, 470, 125, 30}, "ONE;TWO;THREE;FOUR", comboBoxActive)

		// NOTE: gui.DropdownBox must draw after any other control that can be covered on unfolding
		gui.SetStyle(gui.DROPDOWNBOX, gui.TEXT_ALIGNMENT, int64(gui.TEXT_ALIGN_LEFT))
		if gui.DropdownBox(rl.Rectangle{25, 65, 125, 30}, "#01#ONE;#02#TWO;#03#THREE;#04#FOUR", &dropdownBox001Active, dropDown001EditMode) {
			dropDown001EditMode = !dropDown001EditMode
		}

		gui.SetStyle(gui.DROPDOWNBOX, gui.TEXT_ALIGNMENT, gui.TEXT_ALIGN_CENTER)
		if gui.DropdownBox(rl.Rectangle{25, 25, 125, 30}, "ONE;TWO;THREE", &dropdownBox000Active, dropDown000EditMode) {
			dropDown000EditMode = !dropDown000EditMode
		}

		// Second GUI column
		listViewActive = gui.ListView(rl.Rectangle{165, 25, 140, 140}, "Charmander;Bulbasaur;#18#Squirtel;Pikachu;Eevee;Pidgey", &listViewScrollIndex, listViewActive)
		listViewExActive = gui.ListViewEx(rl.Rectangle{165, 180, 140, 200}, listViewExList, &listViewExFocus, &listViewExScrollIndex, listViewExActive)

		toggleGroupActive = gui.ToggleGroup(rl.Rectangle{165, 400, 140, 25}, "#1#ONE\n#3#TWO\n#8#THREE\n#23#", toggleGroupActive)

		// Third GUI column
		gui.Panel(rl.NewRectangle(320, 25, 225, 140), "Panel Info")
		colorPickerValue = gui.ColorPicker(rl.Rectangle{320, 185, 196, 192}, "", colorPickerValue)

		sliderValue = gui.Slider(rl.Rectangle{355, 400, 165, 20}, "TEST",
			fmt.Sprintf("%2.2f", sliderValue), sliderValue, -50, 100)
		sliderBarValue = gui.SliderBar(rl.Rectangle{320, 430, 200, 20}, "",
			fmt.Sprintf("%2.2f", sliderBarValue), sliderBarValue, 0, 100)
		progressValue = gui.ProgressBar(rl.Rectangle{320, 460, 200, 20}, "", "", progressValue, 0, 1)

		// NOTE: View rectangle could be used to perform some scissor test
		var view rl.Rectangle
		gui.ScrollPanel(rl.Rectangle{560, 25, 102, 354}, "", rl.Rectangle{560, 25, 300, 1200}, &viewScroll, &view)

		var mouseCell rl.Vector2
		gui.Grid(rl.Rectangle{560, 25 + 180 + 195, 100, 120}, "", 20, 3, &mouseCell)

		alphaValue = gui.ColorBarAlpha(rl.Rectangle{320, 490, 200, 30}, "", alphaValue)

		gui.StatusBar(rl.Rectangle{0, float32(rl.GetScreenHeight()) - 20, float32(rl.GetScreenWidth()), 20}, "This is a status bar")

		if showMessageBox {
			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.RayWhite, 0.8))
			var result int32 = gui.MessageBox(rl.Rectangle{float32(rl.GetScreenWidth())/2 - 125, float32(rl.GetScreenHeight())/2 - 50, 250, 100}, gui.IconText(gui.ICON_EXIT, "Close Window"), "Do you really want to exit?", "Yes;No")

			if (result == 0) || (result == 2) {
				showMessageBox = false
			} else if result == 1 {
				exitWindow = true
			}
		}

		if showTextInputBox {
			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.RayWhite, 0.8))
			var secretViewActive bool
			var result int32 = gui.TextInputBox(
				rl.Rectangle{float32(rl.GetScreenWidth())/2 - 120, float32(rl.GetScreenHeight())/2 - 60, 240, 140},
				"Save",
				gui.IconText(gui.ICON_FILE_SAVE, "Save file as..."),
				"Ok;Cancel",
				&textInput, 255, &secretViewActive)

			if result == 1 {
				// TODO: Validate textInput value and save
				// strcpy(textInputFileName, textInput)
				// TODO textInputFileName = textInput
			}
			if (result == 0) || (result == 1) || (result == 2) {
				showTextInputBox = false
				//strcpy(textInput, "\0");
				textInput = ""
			}
		}
		//----------------------------------------------------------------------------------

		rl.EndDrawing()
		//----------------------------------------------------------------------------------
	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	rl.CloseWindow() // Close window and OpenGL context
	//--------------------------------------------------------------------------------------
}
