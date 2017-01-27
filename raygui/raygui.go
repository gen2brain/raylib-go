// raygui is simple and easy-to-use IMGUI (immediate mode GUI API) library.
package raygui

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/gen2brain/raylib-go/raylib"
)

// GUI property
type Property int32

// GUI properties enumeration
const (
	GlobalBaseColor Property = iota
	GlobalBorderColor
	GlobalTextColor
	GlobalTextFontsize
	GlobalBorderWidth
	GlobalBackgroundColor
	GlobalLinesColor
	LabelBorderWidth
	LabelTextColor
	LabelTextPadding
	ButtonBorderWidth
	ButtonTextPadding
	ButtonDefaultBorderColor
	ButtonDefaultInsideColor
	ButtonDefaultTextColor
	ButtonHoverBorderColor
	ButtonHoverInsideColor
	ButtonHoverTextColor
	ButtonPressedBorderColor
	ButtonPressedInsideColor
	ButtonPressedTextColor
	ToggleTextPadding
	ToggleBorderWidth
	ToggleDefaultBorderColor
	ToggleDefaultInsideColor
	ToggleDefaultTextColor
	ToggleHoverBorderColor
	ToggleHoverInsideColor
	ToggleHoverTextColor
	TogglePressedBorderColor
	TogglePressedInsideColor
	TogglePressedTextColor
	ToggleActiveBorderColor
	ToggleActiveInsideColor
	ToggleActiveTextColor
	TogglegroupPadding
	SliderBorderWidth
	SliderButtonBorderWidth
	SliderBorderColor
	SliderInsideColor
	SliderDefaultColor
	SliderHoverColor
	SliderActiveColor
	SliderbarBorderColor
	SliderbarInsideColor
	SliderbarDefaultColor
	SliderbarHoverColor
	SliderbarActiveColor
	SliderbarZeroLineColor
	ProgressbarBorderColor
	ProgressbarInsideColor
	ProgressbarProgressColor
	ProgressbarBorderWidth
	SpinnerLabelBorderColor
	SpinnerLabelInsideColor
	SpinnerDefaultButtonBorderColor
	SpinnerDefaultButtonInsideColor
	SpinnerDefaultSymbolColor
	SpinnerDefaultTextColor
	SpinnerHoverButtonBorderColor
	SpinnerHoverButtonInsideColor
	SpinnerHoverSymbolColor
	SpinnerHoverTextColor
	SpinnerPressedButtonBorderColor
	SpinnerPressedButtonInsideColor
	SpinnerPressedSymbolColor
	SpinnerPressedTextColor
	ComboboxPadding
	ComboboxButtonWidth
	ComboboxButtonHeight
	ComboboxBorderWidth
	ComboboxDefaultBorderColor
	ComboboxDefaultInsideColor
	ComboboxDefaultTextColor
	ComboboxDefaultListTextColor
	ComboboxHoverBorderColor
	ComboboxHoverInsideColor
	ComboboxHoverTextColor
	ComboboxHoverListTextColor
	ComboboxPressedBorderColor
	ComboboxPressedInsideColor
	ComboboxPressedTextColor
	ComboboxPressedListBorderColor
	ComboboxPressedListInsideColor
	ComboboxPressedListTextColor
	CheckboxDefaultBorderColor
	CheckboxDefaultInsideColor
	CheckboxHoverBorderColor
	CheckboxHoverInsideColor
	CheckboxClickBorderColor
	CheckboxClickInsideColor
	CheckboxStatusActiveColor
	CheckboxInsideWidth
	TextboxBorderWidth
	TextboxBorderColor
	TextboxInsideColor
	TextboxTextColor
	TextboxLineColor
	TextboxTextFontsize
)

// GUI elements states

const (
	ButtonDefault = iota
	ButtonHover
	ButtonPressed
	ButtonClicked
)

const (
	ToggleUnactive = iota
	ToggleHover
	TogglePressed
	ToggleActive
)

const (
	ComboboxUnactive = iota
	ComboboxHover
	ComboboxPressed
	ComboboxActive
)

const (
	SpinnerDefault = iota
	SpinnerHover
	SpinnerPressed
)

const (
	CheckboxStatus = iota
	CheckboxHover
	CheckboxPressed
)

const (
	SliderDefault = iota
	SliderHover
	SliderActive
)

// Current GUI style (default light)
var style []int = []int{
	0xf5f5f5ff, // GLOBAL_BASE_COLOR
	0xf5f5f5ff, // GLOBAL_BORDER_COLOR
	0xf5f5f5ff, // GLOBAL_TEXT_COLOR
	10,         // GLOBAL_TEXT_FONTSIZE
	1,          // GLOBAL_BORDER_WIDTH
	0xf5f5f5ff, // BACKGROUND_COLOR
	0x90abb5ff, // LINES_COLOR
	1,          // LABEL_BORDER_WIDTH
	0x4d4d4dff, // LABEL_TEXT_COLOR
	20,         // LABEL_TEXT_PADDING
	2,          // BUTTON_BORDER_WIDTH
	20,         // BUTTON_TEXT_PADDING
	0x828282ff, // BUTTON_DEFAULT_BORDER_COLOR
	0xc8c8c8ff, // BUTTON_DEFAULT_INSIDE_COLOR
	0x4d4d4dff, // BUTTON_DEFAULT_TEXT_COLOR
	0xc8c8c8ff, // BUTTON_HOVER_BORDER_COLOR
	0xffffffff, // BUTTON_HOVER_INSIDE_COLOR
	0x868686ff, // BUTTON_HOVER_TEXT_COLOR
	0x7bb0d6ff, // BUTTON_PRESSED_BORDER_COLOR
	0xbcecffff, // BUTTON_PRESSED_INSIDE_COLOR
	0x5f9aa7ff, // BUTTON_PRESSED_TEXT_COLOR
	20,         // TOGGLE_TEXT_PADDING
	1,          // TOGGLE_BORDER_WIDTH
	0x828282ff, // TOGGLE_DEFAULT_BORDER_COLOR
	0xc8c8c8ff, // TOGGLE_DEFAULT_INSIDE_COLOR
	0x828282ff, // TOGGLE_DEFAULT_TEXT_COLOR
	0xc8c8c8ff, // TOGGLE_HOVER_BORDER_COLOR
	0xffffffff, // TOGGLE_HOVER_INSIDE_COLOR
	0x828282ff, // TOGGLE_HOVER_TEXT_COLOR
	0xbdd7eaff, // TOGGLE_PRESSED_BORDER_COLOR
	0xddf5ffff, // TOGGLE_PRESSED_INSIDE_COLOR
	0xafccd3ff, // TOGGLE_PRESSED_TEXT_COLOR
	0x7bb0d6ff, // TOGGLE_ACTIVE_BORDER_COLOR
	0xbcecffff, // TOGGLE_ACTIVE_INSIDE_COLOR
	0x5f9aa7ff, // TOGGLE_ACTIVE_TEXT_COLOR
	3,          // TOGGLEGROUP_PADDING
	1,          // SLIDER_BORDER_WIDTH
	1,          // SLIDER_BUTTON_BORDER_WIDTH
	0x828282ff, // SLIDER_BORDER_COLOR
	0xc8c8c8ff, // SLIDER_INSIDE_COLOR
	0xbcecffff, // SLIDER_DEFAULT_COLOR
	0xffffffff, // SLIDER_HOVER_COLOR
	0xddf5ffff, // SLIDER_ACTIVE_COLOR
	0x828282ff, // SLIDERBAR_BORDER_COLOR
	0xc8c8c8ff, // SLIDERBAR_INSIDE_COLOR
	0xbcecffff, // SLIDERBAR_DEFAULT_COLOR
	0xffffffff, // SLIDERBAR_HOVER_COLOR
	0xddf5ffff, // SLIDERBAR_ACTIVE_COLOR
	0x828282ff, // SLIDERBAR_ZERO_LINE_COLOR
	0x828282ff, // PROGRESSBAR_BORDER_COLOR
	0xc8c8c8ff, // PROGRESSBAR_INSIDE_COLOR
	0xbcecffff, // PROGRESSBAR_PROGRESS_COLOR
	2,          // PROGRESSBAR_BORDER_WIDTH
	0x828282ff, // SPINNER_LABEL_BORDER_COLOR
	0xc8c8c8ff, // SPINNER_LABEL_INSIDE_COLOR
	0x828282ff, // SPINNER_DEFAULT_BUTTON_BORDER_COLOR
	0xc8c8c8ff, // SPINNER_DEFAULT_BUTTON_INSIDE_COLOR
	0x000000ff, // SPINNER_DEFAULT_SYMBOL_COLOR
	0x000000ff, // SPINNER_DEFAULT_TEXT_COLOR
	0xc8c8c8ff, // SPINNER_HOVER_BUTTON_BORDER_COLOR
	0xffffffff, // SPINNER_HOVER_BUTTON_INSIDE_COLOR
	0x000000ff, // SPINNER_HOVER_SYMBOL_COLOR
	0x000000ff, // SPINNER_HOVER_TEXT_COLOR
	0x7bb0d6ff, // SPINNER_PRESSED_BUTTON_BORDER_COLOR
	0xbcecffff, // SPINNER_PRESSED_BUTTON_INSIDE_COLOR
	0x5f9aa7ff, // SPINNER_PRESSED_SYMBOL_COLOR
	0x000000ff, // SPINNER_PRESSED_TEXT_COLOR
	1,          // COMBOBOX_PADDING
	30,         // COMBOBOX_BUTTON_WIDTH
	20,         // COMBOBOX_BUTTON_HEIGHT
	1,          // COMBOBOX_BORDER_WIDTH
	0x828282ff, // COMBOBOX_DEFAULT_BORDER_COLOR
	0xc8c8c8ff, // COMBOBOX_DEFAULT_INSIDE_COLOR
	0x828282ff, // COMBOBOX_DEFAULT_TEXT_COLOR
	0x828282ff, // COMBOBOX_DEFAULT_LIST_TEXT_COLOR
	0xc8c8c8ff, // COMBOBOX_HOVER_BORDER_COLOR
	0xffffffff, // COMBOBOX_HOVER_INSIDE_COLOR
	0x828282ff, // COMBOBOX_HOVER_TEXT_COLOR
	0x828282ff, // COMBOBOX_HOVER_LIST_TEXT_COLOR
	0x7bb0d6ff, // COMBOBOX_PRESSED_BORDER_COLOR
	0xbcecffff, // COMBOBOX_PRESSED_INSIDE_COLOR
	0x5f9aa7ff, // COMBOBOX_PRESSED_TEXT_COLOR
	0x0078acff, // COMBOBOX_PRESSED_LIST_BORDER_COLOR
	0x66e7ffff, // COMBOBOX_PRESSED_LIST_INSIDE_COLOR
	0x0078acff, // COMBOBOX_PRESSED_LIST_TEXT_COLOR
	0x828282ff, // CHECKBOX_DEFAULT_BORDER_COLOR
	0xffffffff, // CHECKBOX_DEFAULT_INSIDE_COLOR
	0xc8c8c8ff, // CHECKBOX_HOVER_BORDER_COLOR
	0xffffffff, // CHECKBOX_HOVER_INSIDE_COLOR
	0x66e7ffff, // CHECKBOX_CLICK_BORDER_COLOR
	0xddf5ffff, // CHECKBOX_CLICK_INSIDE_COLOR
	0xbcecffff, // CHECKBOX_STATUS_ACTIVE_COLOR
	1,          // CHECKBOX_INSIDE_WIDTH
	1,          // TEXTBOX_BORDER_WIDTH
	0x828282ff, // TEXTBOX_BORDER_COLOR
	0xf5f5f5ff, // TEXTBOX_INSIDE_COLOR
	0x000000ff, // TEXTBOX_TEXT_COLOR
	0x000000ff, // TEXTBOX_LINE_COLOR
	10,         // TEXTBOX_TEXT_FONTSIZE
}

// GUI property names (to read/write style text files)
var propertyName []string = []string{
	"GLOBAL_BASE_COLOR",
	"GLOBAL_BORDER_COLOR",
	"GLOBAL_TEXT_COLOR",
	"GLOBAL_TEXT_FONTSIZE",
	"GLOBAL_BORDER_WIDTH",
	"BACKGROUND_COLOR",
	"LINES_COLOR",
	"LABEL_BORDER_WIDTH",
	"LABEL_TEXT_COLOR",
	"LABEL_TEXT_PADDING",
	"BUTTON_BORDER_WIDTH",
	"BUTTON_TEXT_PADDING",
	"BUTTON_DEFAULT_BORDER_COLOR",
	"BUTTON_DEFAULT_INSIDE_COLOR",
	"BUTTON_DEFAULT_TEXT_COLOR",
	"BUTTON_HOVER_BORDER_COLOR",
	"BUTTON_HOVER_INSIDE_COLOR",
	"BUTTON_HOVER_TEXT_COLOR",
	"BUTTON_PRESSED_BORDER_COLOR",
	"BUTTON_PRESSED_INSIDE_COLOR",
	"BUTTON_PRESSED_TEXT_COLOR",
	"TOGGLE_TEXT_PADDING",
	"TOGGLE_BORDER_WIDTH",
	"TOGGLE_DEFAULT_BORDER_COLOR",
	"TOGGLE_DEFAULT_INSIDE_COLOR",
	"TOGGLE_DEFAULT_TEXT_COLOR",
	"TOGGLE_HOVER_BORDER_COLOR",
	"TOGGLE_HOVER_INSIDE_COLOR",
	"TOGGLE_HOVER_TEXT_COLOR",
	"TOGGLE_PRESSED_BORDER_COLOR",
	"TOGGLE_PRESSED_INSIDE_COLOR",
	"TOGGLE_PRESSED_TEXT_COLOR",
	"TOGGLE_ACTIVE_BORDER_COLOR",
	"TOGGLE_ACTIVE_INSIDE_COLOR",
	"TOGGLE_ACTIVE_TEXT_COLOR",
	"TOGGLEGROUP_PADDING",
	"SLIDER_BORDER_WIDTH",
	"SLIDER_BUTTON_BORDER_WIDTH",
	"SLIDER_BORDER_COLOR",
	"SLIDER_INSIDE_COLOR",
	"SLIDER_DEFAULT_COLOR",
	"SLIDER_HOVER_COLOR",
	"SLIDER_ACTIVE_COLOR",
	"SLIDERBAR_BORDER_COLOR",
	"SLIDERBAR_INSIDE_COLOR",
	"SLIDERBAR_DEFAULT_COLOR",
	"SLIDERBAR_HOVER_COLOR",
	"SLIDERBAR_ACTIVE_COLOR",
	"SLIDERBAR_ZERO_LINE_COLOR",
	"PROGRESSBAR_BORDER_COLOR",
	"PROGRESSBAR_INSIDE_COLOR",
	"PROGRESSBAR_PROGRESS_COLOR",
	"PROGRESSBAR_BORDER_WIDTH",
	"SPINNER_LABEL_BORDER_COLOR",
	"SPINNER_LABEL_INSIDE_COLOR",
	"SPINNER_DEFAULT_BUTTON_BORDER_COLOR",
	"SPINNER_DEFAULT_BUTTON_INSIDE_COLOR",
	"SPINNER_DEFAULT_SYMBOL_COLOR",
	"SPINNER_DEFAULT_TEXT_COLOR",
	"SPINNER_HOVER_BUTTON_BORDER_COLOR",
	"SPINNER_HOVER_BUTTON_INSIDE_COLOR",
	"SPINNER_HOVER_SYMBOL_COLOR",
	"SPINNER_HOVER_TEXT_COLOR",
	"SPINNER_PRESSED_BUTTON_BORDER_COLOR",
	"SPINNER_PRESSED_BUTTON_INSIDE_COLOR",
	"SPINNER_PRESSED_SYMBOL_COLOR",
	"SPINNER_PRESSED_TEXT_COLOR",
	"COMBOBOX_PADDING",
	"COMBOBOX_BUTTON_WIDTH",
	"COMBOBOX_BUTTON_HEIGHT",
	"COMBOBOX_BORDER_WIDTH",
	"COMBOBOX_DEFAULT_BORDER_COLOR",
	"COMBOBOX_DEFAULT_INSIDE_COLOR",
	"COMBOBOX_DEFAULT_TEXT_COLOR",
	"COMBOBOX_DEFAULT_LIST_TEXT_COLOR",
	"COMBOBOX_HOVER_BORDER_COLOR",
	"COMBOBOX_HOVER_INSIDE_COLOR",
	"COMBOBOX_HOVER_TEXT_COLOR",
	"COMBOBOX_HOVER_LIST_TEXT_COLOR",
	"COMBOBOX_PRESSED_BORDER_COLOR",
	"COMBOBOX_PRESSED_INSIDE_COLOR",
	"COMBOBOX_PRESSED_TEXT_COLOR",
	"COMBOBOX_PRESSED_LIST_BORDER_COLOR",
	"COMBOBOX_PRESSED_LIST_INSIDE_COLOR",
	"COMBOBOX_PRESSED_LIST_TEXT_COLOR",
	"CHECKBOX_DEFAULT_BORDER_COLOR",
	"CHECKBOX_DEFAULT_INSIDE_COLOR",
	"CHECKBOX_HOVER_BORDER_COLOR",
	"CHECKBOX_HOVER_INSIDE_COLOR",
	"CHECKBOX_CLICK_BORDER_COLOR",
	"CHECKBOX_CLICK_INSIDE_COLOR",
	"CHECKBOX_STATUS_ACTIVE_COLOR",
	"CHECKBOX_INSIDE_WIDTH",
	"TEXTBOX_BORDER_WIDTH",
	"TEXTBOX_BORDER_COLOR",
	"TEXTBOX_INSIDE_COLOR",
	"TEXTBOX_TEXT_COLOR",
	"TEXTBOX_LINE_COLOR",
	"TEXTBOX_TEXT_FONTSIZE",
}

// Get background color
func BackgroundColor() raylib.Color {
	return raylib.GetColor(int32(style[GlobalBackgroundColor]))
}

// Get lines color
func LinesColor() raylib.Color {
	return raylib.GetColor(int32(style[GlobalLinesColor]))
}

// Label element, show text
func Label(bounds raylib.Rectangle, text string) {
	LabelEx(bounds, text, raylib.GetColor(int32(style[LabelTextColor])), raylib.NewColor(0, 0, 0, 0), raylib.NewColor(0, 0, 0, 0))
}

// Label element extended, configurable colors
func LabelEx(bounds raylib.Rectangle, text string, textColor, border, inner raylib.Color) {
	// Update control
	textWidth := raylib.MeasureText(text, int32(style[GlobalTextFontsize]))
	textHeight := int32(style[GlobalTextFontsize])

	if bounds.Width < textWidth {
		bounds.Width = textWidth + int32(style[LabelTextPadding])
	}
	if bounds.Height < textHeight {
		bounds.Height = textHeight + int32(style[LabelTextPadding])/2
	}

	// Draw control
	raylib.DrawRectangleRec(bounds, border)
	raylib.DrawRectangle(bounds.X+int32(style[LabelBorderWidth]), bounds.Y+int32(style[LabelBorderWidth]), bounds.Width-(2*int32(style[LabelBorderWidth])), bounds.Height-(2*int32(style[LabelBorderWidth])), inner)
	raylib.DrawText(text, bounds.X+((bounds.Width/2)-(textWidth/2)), bounds.Y+((bounds.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), textColor)
}

// Button element, returns true when clicked
func Button(bounds raylib.Rectangle, text string) bool {
	buttonState := ButtonDefault
	mousePoint := raylib.GetMousePosition()

	textWidth := raylib.MeasureText(text, int32(style[GlobalTextFontsize]))
	textHeight := int32(style[GlobalTextFontsize])

	// Update control
	if bounds.Width < textWidth {
		bounds.Width = textWidth + int32(style[ButtonTextPadding])
	}

	if bounds.Height < textHeight {
		bounds.Height = textHeight + int32(style[ButtonTextPadding])/2
	}

	if raylib.CheckCollisionPointRec(mousePoint, bounds) {
		if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			buttonState = ButtonPressed
		} else if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) {
			buttonState = ButtonClicked
		} else {
			buttonState = ButtonHover
		}
	}

	// Draw control
	switch buttonState {
	case ButtonDefault:
		raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[ButtonDefaultBorderColor])))
		raylib.DrawRectangle(bounds.X+int32(style[ButtonBorderWidth]), bounds.Y+int32(style[ButtonBorderWidth]), bounds.Width-(2*int32(style[ButtonBorderWidth])), bounds.Height-(2*int32(style[ButtonBorderWidth])), raylib.GetColor(int32(style[ButtonDefaultInsideColor])))
		raylib.DrawText(text, bounds.X+((bounds.Width/2)-(raylib.MeasureText(text, int32(style[GlobalTextFontsize]))/2)), bounds.Y+((bounds.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ButtonDefaultTextColor])))
		break

	case ButtonHover:
		raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[ButtonHoverBorderColor])))
		raylib.DrawRectangle(bounds.X+int32(style[ButtonBorderWidth]), bounds.Y+int32(style[ButtonBorderWidth]), bounds.Width-(2*int32(style[ButtonBorderWidth])), bounds.Height-(2*int32(style[ButtonBorderWidth])), raylib.GetColor(int32(style[ButtonHoverInsideColor])))
		raylib.DrawText(text, bounds.X+((bounds.Width/2)-(raylib.MeasureText(text, int32(style[GlobalTextFontsize]))/2)), bounds.Y+((bounds.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ButtonHoverTextColor])))
		break

	case ButtonPressed:
		raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[ButtonPressedBorderColor])))
		raylib.DrawRectangle(bounds.X+int32(style[ButtonBorderWidth]), bounds.Y+int32(style[ButtonBorderWidth]), bounds.Width-(2*int32(style[ButtonBorderWidth])), bounds.Height-(2*int32(style[ButtonBorderWidth])), raylib.GetColor(int32(style[ButtonPressedInsideColor])))
		raylib.DrawText(text, bounds.X+((bounds.Width/2)-(raylib.MeasureText(text, int32(style[GlobalTextFontsize]))/2)), bounds.Y+((bounds.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ButtonPressedTextColor])))
		break

	case ButtonClicked:
		raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[ButtonPressedBorderColor])))
		raylib.DrawRectangle(bounds.X+int32(style[ButtonBorderWidth]), bounds.Y+int32(style[ButtonBorderWidth]), bounds.Width-(2*int32(style[ButtonBorderWidth])), bounds.Height-(2*int32(style[ButtonBorderWidth])), raylib.GetColor(int32(style[ButtonPressedInsideColor])))
		break

	default:
		break
	}

	if buttonState == ButtonClicked {
		return true
	} else {
		return false
	}
}

// Toggle Button element, returns true when active
func ToggleButton(bounds raylib.Rectangle, text string, toggle bool) bool {
	toggleState := ToggleUnactive
	toggleButton := bounds
	mousePoint := raylib.GetMousePosition()

	textWidth := raylib.MeasureText(text, int32(style[GlobalTextFontsize]))
	textHeight := int32(style[GlobalTextFontsize])

	// Update control
	if toggleButton.Width < textWidth {
		toggleButton.Width = textWidth + int32(style[ToggleTextPadding])
	}
	if toggleButton.Height < textHeight {
		toggleButton.Height = textHeight + int32(style[ToggleTextPadding])/2
	}

	if toggle {
		toggleState = ToggleActive
	} else {
		toggleState = ToggleUnactive
	}

	if raylib.CheckCollisionPointRec(mousePoint, toggleButton) {
		if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			toggleState = TogglePressed
		} else if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) {
			if toggle {
				toggle = false
				toggleState = ToggleUnactive
			} else {
				toggle = true
				toggleState = ToggleActive
			}
		} else {
			toggleState = ToggleHover
		}
	}

	// Draw control
	switch toggleState {
	case ToggleUnactive:
		raylib.DrawRectangleRec(toggleButton, raylib.GetColor(int32(style[ToggleDefaultBorderColor])))
		raylib.DrawRectangle(toggleButton.X+int32(style[ToggleBorderWidth]), toggleButton.Y+int32(style[ToggleBorderWidth]), toggleButton.Width-(2*int32(style[ToggleBorderWidth])), toggleButton.Height-(2*int32(style[ToggleBorderWidth])), raylib.GetColor(int32(style[ToggleDefaultInsideColor])))
		raylib.DrawText(text, toggleButton.X+((toggleButton.Width/2)-(raylib.MeasureText(text, int32(style[GlobalTextFontsize]))/2)), toggleButton.Y+((toggleButton.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ToggleDefaultTextColor])))
		break
	case ToggleHover:
		raylib.DrawRectangleRec(toggleButton, raylib.GetColor(int32(style[ToggleHoverBorderColor])))
		raylib.DrawRectangle(toggleButton.X+int32(style[ToggleBorderWidth]), toggleButton.Y+int32(style[ToggleBorderWidth]), toggleButton.Width-(2*int32(style[ToggleBorderWidth])), toggleButton.Height-(2*int32(style[ToggleBorderWidth])), raylib.GetColor(int32(style[ToggleHoverInsideColor])))
		raylib.DrawText(text, toggleButton.X+((toggleButton.Width/2)-(raylib.MeasureText(text, int32(style[GlobalTextFontsize]))/2)), toggleButton.Y+((toggleButton.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ToggleHoverTextColor])))
		break
	case TogglePressed:
		raylib.DrawRectangleRec(toggleButton, raylib.GetColor(int32(style[TogglePressedBorderColor])))
		raylib.DrawRectangle(toggleButton.X+int32(style[ToggleBorderWidth]), toggleButton.Y+int32(style[ToggleBorderWidth]), toggleButton.Width-(2*int32(style[ToggleBorderWidth])), toggleButton.Height-(2*int32(style[ToggleBorderWidth])), raylib.GetColor(int32(style[TogglePressedInsideColor])))
		raylib.DrawText(text, toggleButton.X+((toggleButton.Width/2)-(raylib.MeasureText(text, int32(style[GlobalTextFontsize]))/2)), toggleButton.Y+((toggleButton.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[TogglePressedTextColor])))
		break
	case ToggleActive:
		raylib.DrawRectangleRec(toggleButton, raylib.GetColor(int32(style[ToggleActiveBorderColor])))
		raylib.DrawRectangle(toggleButton.X+int32(style[ToggleBorderWidth]), toggleButton.Y+int32(style[ToggleBorderWidth]), toggleButton.Width-(2*int32(style[ToggleBorderWidth])), toggleButton.Height-(2*int32(style[ToggleBorderWidth])), raylib.GetColor(int32(style[ToggleActiveInsideColor])))
		raylib.DrawText(text, toggleButton.X+((toggleButton.Width/2)-(raylib.MeasureText(text, int32(style[GlobalTextFontsize]))/2)), toggleButton.Y+((toggleButton.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ToggleActiveTextColor])))
		break
	default:
		break
	}

	return toggle
}

// Toggle Group element, returns toggled button index
func ToggleGroup(bounds raylib.Rectangle, toggleText []string, toggleActive int) int {
	for i := 0; i < len(toggleText); i++ {
		if i == toggleActive {
			ToggleButton(raylib.NewRectangle(bounds.X+int32(i)*(bounds.Width+int32(style[TogglegroupPadding])), bounds.Y, bounds.Width, bounds.Height), toggleText[i], true)
		} else if ToggleButton(raylib.NewRectangle(bounds.X+int32(i)*(bounds.Width+int32(style[TogglegroupPadding])), bounds.Y, bounds.Width, bounds.Height), toggleText[i], false) {
			toggleActive = i
		}
	}

	return toggleActive
}

// Combo Box element, returns selected item index
func ComboBox(bounds raylib.Rectangle, comboText []string, comboActive int) int {
	comboBoxState := ComboboxUnactive
	comboBoxButton := bounds
	click := raylib.NewRectangle(bounds.X+bounds.Width+int32(style[ComboboxPadding]), bounds.Y, int32(style[ComboboxButtonWidth]), int32(style[ComboboxButtonHeight]))
	mousePoint := raylib.GetMousePosition()

	textWidth := int32(0)
	textHeight := int32(style[GlobalTextFontsize])

	comboNum := len(comboText)

	for i := 0; i < comboNum; i++ {
		if i == comboActive {
			// Update control
			textWidth = raylib.MeasureText(comboText[i], int32(style[GlobalTextFontsize]))

			if comboBoxButton.Width < textWidth {
				comboBoxButton.Width = textWidth + int32(style[ToggleTextPadding])
			}
			if comboBoxButton.Height < textHeight {
				comboBoxButton.Height = textHeight + int32(style[ToggleTextPadding])/2
			}

			if raylib.CheckCollisionPointRec(mousePoint, comboBoxButton) || raylib.CheckCollisionPointRec(mousePoint, click) {
				if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
					comboBoxState = ComboboxPressed
				} else if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) {
					comboBoxState = ComboboxActive
				} else {
					comboBoxState = ComboboxHover
				}
			}

			// Draw control
			switch comboBoxState {
			case ComboboxUnactive:
				raylib.DrawRectangleRec(comboBoxButton, raylib.GetColor(int32(style[ComboboxDefaultBorderColor])))
				raylib.DrawRectangle(comboBoxButton.X+int32(style[ComboboxBorderWidth]), comboBoxButton.Y+int32(style[ComboboxBorderWidth]), comboBoxButton.Width-(2*int32(style[ComboboxBorderWidth])), comboBoxButton.Height-(2*int32(style[ComboboxBorderWidth])), raylib.GetColor(int32(style[ComboboxDefaultInsideColor])))

				raylib.DrawRectangleRec(click, raylib.GetColor(int32(style[ComboboxDefaultBorderColor])))
				raylib.DrawRectangle(click.X+int32(style[ComboboxBorderWidth]), click.Y+int32(style[ComboboxBorderWidth]), click.Width-(2*int32(style[ComboboxBorderWidth])), click.Height-(2*int32(style[ComboboxBorderWidth])), raylib.GetColor(int32(style[ComboboxDefaultInsideColor])))
				raylib.DrawText(fmt.Sprintf("%d/%d", comboActive+1, comboNum), click.X+((click.Width/2)-(raylib.MeasureText(fmt.Sprintf("%d/%d", comboActive+1, comboNum), int32(style[GlobalTextFontsize]))/2)), click.Y+((click.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ComboboxDefaultListTextColor])))
				raylib.DrawText(comboText[i], comboBoxButton.X+((comboBoxButton.Width/2)-(raylib.MeasureText(comboText[i], int32(style[GlobalTextFontsize]))/2)), comboBoxButton.Y+((comboBoxButton.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ComboboxDefaultTextColor])))
				break
			case ComboboxHover:
				raylib.DrawRectangleRec(comboBoxButton, raylib.GetColor(int32(style[ComboboxHoverBorderColor])))
				raylib.DrawRectangle(comboBoxButton.X+int32(style[ComboboxBorderWidth]), comboBoxButton.Y+int32(style[ComboboxBorderWidth]), comboBoxButton.Width-(2*int32(style[ComboboxBorderWidth])), comboBoxButton.Height-(2*int32(style[ComboboxBorderWidth])), raylib.GetColor(int32(style[ComboboxHoverInsideColor])))

				raylib.DrawRectangleRec(click, raylib.GetColor(int32(style[ComboboxHoverBorderColor])))
				raylib.DrawRectangle(click.X+int32(style[ComboboxBorderWidth]), click.Y+int32(style[ComboboxBorderWidth]), click.Width-(2*int32(style[ComboboxBorderWidth])), click.Height-(2*int32(style[ComboboxBorderWidth])), raylib.GetColor(int32(style[ComboboxHoverInsideColor])))
				raylib.DrawText(fmt.Sprintf("%d/%d", comboActive+1, comboNum), click.X+((click.Width/2)-(raylib.MeasureText(fmt.Sprintf("%d/%d", comboActive+1, comboNum), int32(style[GlobalTextFontsize]))/2)), click.Y+((click.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ComboboxHoverListTextColor])))
				raylib.DrawText(comboText[i], comboBoxButton.X+((comboBoxButton.Width/2)-(raylib.MeasureText(comboText[i], int32(style[GlobalTextFontsize]))/2)), comboBoxButton.Y+((comboBoxButton.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ComboboxHoverTextColor])))
				break
			case ComboboxPressed:
				raylib.DrawRectangleRec(comboBoxButton, raylib.GetColor(int32(style[ComboboxPressedBorderColor])))
				raylib.DrawRectangle(comboBoxButton.X+int32(style[ComboboxBorderWidth]), comboBoxButton.Y+int32(style[ComboboxBorderWidth]), comboBoxButton.Width-(2*int32(style[ComboboxBorderWidth])), comboBoxButton.Height-(2*int32(style[ComboboxBorderWidth])), raylib.GetColor(int32(style[ComboboxPressedInsideColor])))

				raylib.DrawRectangleRec(click, raylib.GetColor(int32(style[ComboboxPressedListBorderColor])))
				raylib.DrawRectangle(click.X+int32(style[ComboboxBorderWidth]), click.Y+int32(style[ComboboxBorderWidth]), click.Width-(2*int32(style[ComboboxBorderWidth])), click.Height-(2*int32(style[ComboboxBorderWidth])), raylib.GetColor(int32(style[ComboboxPressedListInsideColor])))
				raylib.DrawText(fmt.Sprintf("%d/%d", comboActive+1, comboNum), click.X+((click.Width/2)-(raylib.MeasureText(fmt.Sprintf("%d/%d", comboActive+1, comboNum), int32(style[GlobalTextFontsize]))/2)), click.Y+((click.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ComboboxPressedListTextColor])))
				raylib.DrawText(comboText[i], comboBoxButton.X+((comboBoxButton.Width/2)-(raylib.MeasureText(comboText[i], int32(style[GlobalTextFontsize]))/2)), comboBoxButton.Y+((comboBoxButton.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ComboboxPressedTextColor])))
				break
			case ComboboxActive:
				raylib.DrawRectangleRec(comboBoxButton, raylib.GetColor(int32(style[ComboboxPressedBorderColor])))
				raylib.DrawRectangle(comboBoxButton.X+int32(style[ComboboxBorderWidth]), comboBoxButton.Y+int32(style[ComboboxBorderWidth]), comboBoxButton.Width-(2*int32(style[ComboboxBorderWidth])), comboBoxButton.Height-(2*int32(style[ComboboxBorderWidth])), raylib.GetColor(int32(style[ComboboxPressedInsideColor])))

				raylib.DrawRectangleRec(click, raylib.GetColor(int32(style[ComboboxPressedListBorderColor])))
				raylib.DrawRectangle(click.X+int32(style[ComboboxBorderWidth]), click.Y+int32(style[ComboboxBorderWidth]), click.Width-(2*int32(style[ComboboxBorderWidth])), click.Height-(2*int32(style[ComboboxBorderWidth])), raylib.GetColor(int32(style[ComboboxPressedListInsideColor])))
				raylib.DrawText(fmt.Sprintf("%d/%d", comboActive+1, comboNum), click.X+((click.Width/2)-(raylib.MeasureText(fmt.Sprintf("%d/%d", comboActive+1, comboNum), int32(style[GlobalTextFontsize]))/2)), click.Y+((click.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ComboboxPressedListTextColor])))
				raylib.DrawText(comboText[i], comboBoxButton.X+((comboBoxButton.Width/2)-(raylib.MeasureText(comboText[i], int32(style[GlobalTextFontsize]))/2)), comboBoxButton.Y+((comboBoxButton.Height/2)-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[ComboboxPressedTextColor])))
				break
			default:
				break
			}

		}
	}

	if raylib.CheckCollisionPointRec(raylib.GetMousePosition(), bounds) || raylib.CheckCollisionPointRec(raylib.GetMousePosition(), click) {
		if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
			comboActive += 1
			if comboActive >= comboNum {
				comboActive = 0
			}
		}
	}

	return comboActive
}

// Check Box element, returns true when active
func CheckBox(bounds raylib.Rectangle, text string, checked bool) bool {
	checkBoxState := CheckboxStatus
	mousePoint := raylib.GetMousePosition()

	// Update control
	if raylib.CheckCollisionPointRec(mousePoint, bounds) {
		if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			checkBoxState = CheckboxPressed
		} else if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) {
			checkBoxState = CheckboxStatus
			checked = !checked
		} else {
			checkBoxState = CheckboxHover
		}
	}

	// Draw control
	switch checkBoxState {
	case CheckboxHover:
		raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[CheckboxHoverBorderColor])))
		raylib.DrawRectangle(bounds.X+int32(style[ToggleBorderWidth]), bounds.Y+int32(style[ToggleBorderWidth]), bounds.Width-(2*int32(style[ToggleBorderWidth])), bounds.Height-(2*int32(style[ToggleBorderWidth])), raylib.GetColor(int32(style[CheckboxHoverInsideColor])))
		break
	case CheckboxStatus:
		raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[CheckboxDefaultBorderColor])))
		raylib.DrawRectangle(bounds.X+int32(style[ToggleBorderWidth]), bounds.Y+int32(style[ToggleBorderWidth]), bounds.Width-(2*int32(style[ToggleBorderWidth])), bounds.Height-(2*int32(style[ToggleBorderWidth])), raylib.GetColor(int32(style[CheckboxDefaultInsideColor])))
		break
	case CheckboxPressed:
		raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[CheckboxClickBorderColor])))
		raylib.DrawRectangle(bounds.X+int32(style[ToggleBorderWidth]), bounds.Y+int32(style[ToggleBorderWidth]), bounds.Width-(2*int32(style[ToggleBorderWidth])), bounds.Height-(2*int32(style[ToggleBorderWidth])), raylib.GetColor(int32(style[CheckboxClickInsideColor])))
		break
	default:
		break
	}

	if text != "" {
		raylib.DrawText(text, bounds.X+bounds.Width+2, bounds.Y+((bounds.Height/2)-(int32(style[GlobalTextFontsize])/2)+1), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[LabelTextColor])))
	}

	if checked {
		raylib.DrawRectangle(bounds.X+int32(style[CheckboxInsideWidth]), bounds.Y+int32(style[CheckboxInsideWidth]), bounds.Width-(2*int32(style[CheckboxInsideWidth])), bounds.Height-(2*int32(style[CheckboxInsideWidth])), raylib.GetColor(int32(style[CheckboxStatusActiveColor])))
	}

	return checked
}

// Slider element, returns selected value
func Slider(bounds raylib.Rectangle, value, minValue, maxValue float32) float32 {
	sliderPos := float32(0)
	sliderState := SliderDefault

	buttonTravelDistance := float32(0)
	mousePoint := raylib.GetMousePosition()

	// Update control
	if value < minValue {
		value = minValue
	} else if value >= maxValue {
		value = maxValue
	}

	sliderPos = (value - minValue) / (maxValue - minValue)

	sliderButton := raylib.Rectangle{}
	sliderButton.Width = (bounds.Width-(2*int32(style[SliderButtonBorderWidth])))/10 - 8
	sliderButton.Height = bounds.Height - (2 * int32(style[SliderBorderWidth]+2*style[SliderButtonBorderWidth]))

	sliderButtonMinPos := bounds.X + int32(style[SliderBorderWidth]) + int32(style[SliderButtonBorderWidth])
	sliderButtonMaxPos := bounds.X + bounds.Width - (int32(style[SliderBorderWidth]) + int32(style[SliderButtonBorderWidth]) + sliderButton.Width)

	buttonTravelDistance = float32(sliderButtonMaxPos - sliderButtonMinPos)

	sliderButton.X = bounds.X + int32(style[SliderBorderWidth]) + int32(style[SliderButtonBorderWidth]) + int32(sliderPos*buttonTravelDistance)
	sliderButton.Y = bounds.Y + int32(style[SliderBorderWidth]) + int32(style[SliderButtonBorderWidth])

	if raylib.CheckCollisionPointRec(mousePoint, bounds) {
		sliderState = SliderHover

		if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			sliderState = SliderActive
		}

		if sliderState == SliderActive && raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			sliderButton.X = int32(mousePoint.X) - sliderButton.Width/2

			if sliderButton.X <= sliderButtonMinPos {
				sliderButton.X = sliderButtonMinPos
			} else if sliderButton.X >= sliderButtonMaxPos {
				sliderButton.X = sliderButtonMaxPos
			}

			sliderPos = float32(sliderButton.X-sliderButtonMinPos) / buttonTravelDistance
		}
	} else {
		sliderState = SliderDefault
	}

	// Draw control
	raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[SliderBorderColor])))
	raylib.DrawRectangle(bounds.X+int32(style[SliderBorderWidth]), bounds.Y+int32(style[SliderBorderWidth]), bounds.Width-(2*int32(style[SliderBorderWidth])), bounds.Height-(2*int32(style[SliderBorderWidth])), raylib.GetColor(int32(style[SliderInsideColor])))

	switch sliderState {
	case SliderDefault:
		raylib.DrawRectangleRec(sliderButton, raylib.GetColor(int32(style[SliderDefaultColor])))
		break
	case SliderHover:
		raylib.DrawRectangleRec(sliderButton, raylib.GetColor(int32(style[SliderHoverColor])))
		break
	case SliderActive:
		raylib.DrawRectangleRec(sliderButton, raylib.GetColor(int32(style[SliderActiveColor])))
		break
	default:
		break
	}

	return minValue + (maxValue-minValue)*sliderPos
}

// Slider Bar element, returns selected value
func SliderBar(bounds raylib.Rectangle, value, minValue, maxValue float32) float32 {
	sliderState := SliderDefault
	mousePoint := raylib.GetMousePosition()

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

	sliderBar := raylib.Rectangle{}

	sliderBar.X = bounds.X + int32(style[SliderBorderWidth])
	sliderBar.Y = bounds.Y + int32(style[SliderBorderWidth])
	sliderBar.Width = int32((fixedValue * (float32(bounds.Width) - 2*float32(style[SliderBorderWidth]))) / (maxValue - fixedMinValue))
	sliderBar.Height = bounds.Height - 2*int32(style[SliderBorderWidth])

	if raylib.CheckCollisionPointRec(mousePoint, bounds) {
		sliderState = SliderHover

		if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			sliderState = SliderActive

			sliderBar.Width = (int32(mousePoint.X) - bounds.X - int32(style[SliderBorderWidth]))

			if int32(mousePoint.X) <= (bounds.X + int32(style[SliderBorderWidth])) {
				sliderBar.Width = 0
			} else if int32(mousePoint.X) >= (bounds.X + bounds.Width - int32(style[SliderBorderWidth])) {
				sliderBar.Width = bounds.Width - 2*int32(style[SliderBorderWidth])
			}
		}
	} else {
		sliderState = SliderDefault
	}

	fixedValue = (float32(sliderBar.Width) * (maxValue - fixedMinValue)) / (float32(bounds.Width) - 2*float32(style[SliderBorderWidth]))

	// Draw control
	raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[SliderbarBorderColor])))
	raylib.DrawRectangle(bounds.X+int32(style[SliderBorderWidth]), bounds.Y+int32(style[SliderBorderWidth]), bounds.Width-(2*int32(style[SliderBorderWidth])), bounds.Height-(2*int32(style[SliderBorderWidth])), raylib.GetColor(int32(style[SliderbarInsideColor])))

	switch sliderState {
	case SliderDefault:
		raylib.DrawRectangleRec(sliderBar, raylib.GetColor(int32(style[SliderbarDefaultColor])))
		break
	case SliderHover:
		raylib.DrawRectangleRec(sliderBar, raylib.GetColor(int32(style[SliderbarHoverColor])))
		break
	case SliderActive:
		raylib.DrawRectangleRec(sliderBar, raylib.GetColor(int32(style[SliderbarActiveColor])))
		break
	default:
		break
	}

	if minValue < 0 && maxValue > 0 {
		raylib.DrawRectangle((bounds.X+int32(style[SliderBorderWidth]))-int32(minValue*(float32(bounds.Width-(int32(style[SliderBorderWidth])*2))/maxValue)), sliderBar.Y, 1, sliderBar.Height, raylib.GetColor(int32(style[SliderbarZeroLineColor])))
	}

	return fixedValue + minValue
}

// Progress Bar element, shows current progress value
func ProgressBar(bounds raylib.Rectangle, value float32) {
	if value > 1.0 {
		value = 1.0
	} else if value < 0.0 {
		value = 0.0
	}

	progressBar := raylib.NewRectangle(bounds.X+int32(style[ProgressbarBorderWidth]), bounds.Y+int32(style[ProgressbarBorderWidth]), bounds.Width-(int32(style[ProgressbarBorderWidth])*2), bounds.Height-(int32(style[ProgressbarBorderWidth])*2))
	progressValue := raylib.NewRectangle(bounds.X+int32(style[ProgressbarBorderWidth]), bounds.Y+int32(style[ProgressbarBorderWidth]), int32(value*float32(bounds.Width-int32(style[ProgressbarBorderWidth])*2)), bounds.Height-(int32(style[ProgressbarBorderWidth])*2))

	// Draw control
	raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[ProgressbarBorderColor])))
	raylib.DrawRectangleRec(progressBar, raylib.GetColor(int32(style[ProgressbarInsideColor])))
	raylib.DrawRectangleRec(progressValue, raylib.GetColor(int32(style[ProgressbarProgressColor])))
}

// Spinner element, returns selected value
func Spinner(bounds raylib.Rectangle, value, minValue, maxValue int) int {
	spinnerState := SpinnerDefault
	labelBoxBound := raylib.NewRectangle(bounds.X+bounds.Width/4+1, bounds.Y, bounds.Width/2, bounds.Height)
	leftButtonBound := raylib.NewRectangle(bounds.X, bounds.Y, bounds.Width/4, bounds.Height)
	rightButtonBound := raylib.NewRectangle(bounds.X+bounds.Width-bounds.Width/4+1, bounds.Y, bounds.Width/4, bounds.Height)
	mousePoint := raylib.GetMousePosition()

	textWidth := raylib.MeasureText(fmt.Sprintf("%d", value), int32(style[GlobalTextFontsize]))

	buttonSide := 0

	framesCounter := 0
	valueSpeed := false

	// Update control
	if raylib.CheckCollisionPointRec(mousePoint, leftButtonBound) || raylib.CheckCollisionPointRec(mousePoint, rightButtonBound) || raylib.CheckCollisionPointRec(mousePoint, labelBoxBound) {
		if raylib.IsKeyDown(raylib.KeyLeft) {
			spinnerState = SpinnerPressed
			buttonSide = 1

			if value > minValue {
				value -= 1
			}
		} else if raylib.IsKeyDown(raylib.KeyRight) {
			spinnerState = SpinnerPressed
			buttonSide = 2

			if value < maxValue {
				value += 1
			}
		}
	}

	if raylib.CheckCollisionPointRec(mousePoint, leftButtonBound) {
		buttonSide = 1
		spinnerState = SpinnerHover

		if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			if !valueSpeed {
				if value > minValue {
					value--
				}
				valueSpeed = true
			} else {
				framesCounter++
			}

			spinnerState = SpinnerPressed

			if value > minValue {
				if framesCounter >= 30 {
					value -= 1
				}
			}
		}
	} else if raylib.CheckCollisionPointRec(mousePoint, rightButtonBound) {
		buttonSide = 2
		spinnerState = SpinnerHover

		if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			if !valueSpeed {
				if value < maxValue {
					value++
				}
				valueSpeed = true
			} else {
				framesCounter++
			}

			spinnerState = SpinnerPressed

			if value < maxValue {
				if framesCounter >= 30 {
					value += 1
				}
			}
		}
	} else if !raylib.CheckCollisionPointRec(mousePoint, labelBoxBound) {
		buttonSide = 0
	}

	if raylib.IsMouseButtonUp(raylib.MouseLeftButton) {
		valueSpeed = false
		framesCounter = 0
	}

	// Draw control
	switch spinnerState {
	case SpinnerDefault:
		raylib.DrawRectangleRec(leftButtonBound, raylib.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
		raylib.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, raylib.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

		raylib.DrawRectangleRec(rightButtonBound, raylib.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
		raylib.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, raylib.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

		raylib.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(raylib.MeasureText("+", int32(style[GlobalTextFontsize])))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerDefaultSymbolColor])))
		raylib.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(raylib.MeasureText("-", int32(style[GlobalTextFontsize])))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerDefaultSymbolColor])))

		raylib.DrawRectangleRec(labelBoxBound, raylib.GetColor(int32(style[SpinnerLabelBorderColor])))
		raylib.DrawRectangle(labelBoxBound.X+1, labelBoxBound.Y+1, labelBoxBound.Width-2, labelBoxBound.Height-2, raylib.GetColor(int32(style[SpinnerLabelInsideColor])))

		raylib.DrawText(fmt.Sprintf("%d", value), labelBoxBound.X+(labelBoxBound.Width/2-textWidth/2), labelBoxBound.Y+(labelBoxBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerDefaultTextColor])))
		break
	case SpinnerHover:
		if buttonSide == 1 {
			raylib.DrawRectangleRec(leftButtonBound, raylib.GetColor(int32(style[SpinnerHoverButtonBorderColor])))
			raylib.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, raylib.GetColor(int32(style[SpinnerHoverButtonInsideColor])))

			raylib.DrawRectangleRec(rightButtonBound, raylib.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			raylib.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, raylib.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			raylib.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(raylib.MeasureText("+", int32(style[GlobalTextFontsize])))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerHoverSymbolColor])))
			raylib.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(raylib.MeasureText("-", int32(style[GlobalTextFontsize])))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerDefaultSymbolColor])))
		} else if buttonSide == 2 {
			raylib.DrawRectangleRec(leftButtonBound, raylib.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			raylib.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, raylib.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			raylib.DrawRectangleRec(rightButtonBound, raylib.GetColor(int32(style[SpinnerHoverButtonBorderColor])))
			raylib.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, raylib.GetColor(int32(style[SpinnerHoverButtonInsideColor])))

			raylib.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(raylib.MeasureText("+", int32(style[GlobalTextFontsize])))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerDefaultSymbolColor])))
			raylib.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(raylib.MeasureText("-", int32(style[GlobalTextFontsize])))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerHoverSymbolColor])))
		}

		raylib.DrawRectangleRec(labelBoxBound, raylib.GetColor(int32(style[SpinnerLabelBorderColor])))
		raylib.DrawRectangle(labelBoxBound.X+1, labelBoxBound.Y+1, labelBoxBound.Width-2, labelBoxBound.Height-2, raylib.GetColor(int32(style[SpinnerLabelInsideColor])))

		raylib.DrawText(fmt.Sprintf("%d", value), labelBoxBound.X+(labelBoxBound.Width/2-textWidth/2), labelBoxBound.Y+(labelBoxBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerHoverTextColor])))
		break
	case SpinnerPressed:
		if buttonSide == 1 {
			raylib.DrawRectangleRec(leftButtonBound, raylib.GetColor(int32(style[SpinnerPressedButtonBorderColor])))
			raylib.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, raylib.GetColor(int32(style[SpinnerPressedButtonInsideColor])))

			raylib.DrawRectangleRec(rightButtonBound, raylib.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			raylib.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, raylib.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			raylib.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(raylib.MeasureText("+", int32(style[GlobalTextFontsize])))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerPressedSymbolColor])))
			raylib.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(raylib.MeasureText("-", int32(style[GlobalTextFontsize])))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerDefaultSymbolColor])))
		} else if buttonSide == 2 {
			raylib.DrawRectangleRec(leftButtonBound, raylib.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			raylib.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, raylib.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			raylib.DrawRectangleRec(rightButtonBound, raylib.GetColor(int32(style[SpinnerPressedButtonBorderColor])))
			raylib.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, raylib.GetColor(int32(style[SpinnerPressedButtonInsideColor])))

			raylib.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(raylib.MeasureText("+", int32(style[GlobalTextFontsize])))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerDefaultSymbolColor])))
			raylib.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(raylib.MeasureText("-", int32(style[GlobalTextFontsize])))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerPressedSymbolColor])))
		}

		raylib.DrawRectangleRec(labelBoxBound, raylib.GetColor(int32(style[SpinnerLabelBorderColor])))
		raylib.DrawRectangle(labelBoxBound.X+1, labelBoxBound.Y+1, labelBoxBound.Width-2, labelBoxBound.Height-2, raylib.GetColor(int32(style[SpinnerLabelInsideColor])))

		raylib.DrawText(fmt.Sprintf("%d", value), labelBoxBound.X+(labelBoxBound.Width/2-textWidth/2), labelBoxBound.Y+(labelBoxBound.Height/2-(int32(style[GlobalTextFontsize])/2)), int32(style[GlobalTextFontsize]), raylib.GetColor(int32(style[SpinnerPressedTextColor])))
		break
	default:
		break
	}

	return value
}

// Text Box element, returns input text
func TextBox(bounds raylib.Rectangle, text string) string {
	keyBackspaceText := int32(259) // GLFW BACKSPACE: 3 + 256

	initPos := bounds.X + 4
	letter := int32(-1)
	framesCounter := 0

	mousePoint := raylib.GetMousePosition()

	// Update control
	framesCounter++

	letter = raylib.GetKeyPressed()

	if raylib.CheckCollisionPointRec(mousePoint, bounds) {
		if letter != -1 {
			if letter == keyBackspaceText {
				if len(text) > 0 {
					text = text[:len(text)-1]
				}
			} else {
				if letter >= 32 && letter < 127 {
					text = text + fmt.Sprintf("%c", letter)
				}
			}
		}
	}

	// Draw control
	if raylib.CheckCollisionPointRec(mousePoint, bounds) {
		raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[ToggleActiveBorderColor])))
	} else {
		raylib.DrawRectangleRec(bounds, raylib.GetColor(int32(style[TextboxBorderColor])))
	}

	raylib.DrawRectangle(bounds.X+int32(style[TextboxBorderWidth]), bounds.Y+int32(style[TextboxBorderWidth]), bounds.Width-(int32(style[TextboxBorderWidth])*2), bounds.Height-(int32(style[TextboxBorderWidth])*2), raylib.GetColor(int32(style[TextboxInsideColor])))

	for i := 0; i < len(text); i++ {
		raylib.DrawText(fmt.Sprintf("%c", text[i]), initPos, bounds.Y+int32(style[TextboxTextFontsize]), int32(style[TextboxTextFontsize]), raylib.GetColor(int32(style[TextboxTextColor])))
		initPos += (raylib.MeasureText(fmt.Sprintf("%c", text[i]), int32(style[GlobalTextFontsize])) + 2)
	}

	if (framesCounter/20)%2 == 0 && raylib.CheckCollisionPointRec(mousePoint, bounds) {
		raylib.DrawRectangle(initPos+2, bounds.Y+5, 1, 20, raylib.GetColor(int32(style[TextboxLineColor])))
	}

	return text
}

// Save GUI style file
func SaveGuiStyle(fileName string) {
	var styleFile string
	for i := 0; i < len(propertyName); i++ {
		styleFile += fmt.Sprintf("%-40s0x%x\n", propertyName[i], GetStyleProperty(Property(i)))
	}

	ioutil.WriteFile(fileName, []byte(styleFile), 0644)
}

// Load GUI style file
func LoadGuiStyle(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}

		id := fields[0]
		value := fields[1]

		for i := 0; i < len(propertyName); i++ {
			if id == propertyName[i] {
				if strings.HasPrefix(value, "0x") {
					value = value[2:]
				}

				v, err := strconv.ParseInt(value, 16, 64)
				if err == nil {
					style[i] = int(v)
				}
			}
		}
	}
}

// Set one style property
func SetStyleProperty(guiProperty Property, value int) {
	numColorSamples := 10

	if guiProperty == GlobalBaseColor {
		baseColor := raylib.GetColor(int32(value))
		fadeColor := make([]raylib.Color, numColorSamples)

		for i := 0; i < numColorSamples; i++ {
			fadeColor[i] = colorMultiply(baseColor, 1.0-float32(i)/float32(numColorSamples-1))
		}

		style[GlobalBaseColor] = value
		style[GlobalBackgroundColor] = int(raylib.GetHexValue(fadeColor[3]))
		style[ButtonDefaultInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[ButtonHoverInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[ButtonPressedInsideColor] = int(raylib.GetHexValue(fadeColor[5]))
		style[ToggleDefaultInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[ToggleHoverInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[TogglePressedInsideColor] = int(raylib.GetHexValue(fadeColor[5]))
		style[ToggleActiveInsideColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[SliderInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[SliderDefaultColor] = int(raylib.GetHexValue(fadeColor[6]))
		style[SliderHoverColor] = int(raylib.GetHexValue(fadeColor[7]))
		style[SliderActiveColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[SliderbarInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[SliderbarDefaultColor] = int(raylib.GetHexValue(fadeColor[6]))
		style[SliderbarHoverColor] = int(raylib.GetHexValue(fadeColor[7]))
		style[SliderbarActiveColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[SliderbarZeroLineColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[ProgressbarInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[ProgressbarProgressColor] = int(raylib.GetHexValue(fadeColor[6]))
		style[SpinnerLabelInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[SpinnerDefaultButtonInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[SpinnerHoverButtonInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[SpinnerPressedButtonInsideColor] = int(raylib.GetHexValue(fadeColor[5]))
		style[ComboboxDefaultInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[ComboboxHoverInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[ComboboxPressedInsideColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[ComboboxPressedListInsideColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[CheckboxDefaultInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[CheckboxClickInsideColor] = int(raylib.GetHexValue(fadeColor[6]))
		style[CheckboxStatusActiveColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[TextboxInsideColor] = int(raylib.GetHexValue(fadeColor[4]))
	} else if guiProperty == GlobalBorderColor {
		baseColor := raylib.GetColor(int32(value))
		fadeColor := make([]raylib.Color, numColorSamples)

		for i := 0; i < numColorSamples; i++ {
			fadeColor[i] = colorMultiply(baseColor, 1.0-float32(i)/float32(numColorSamples-1))
		}

		style[GlobalBorderColor] = value
		style[ButtonDefaultBorderColor] = int(raylib.GetHexValue(fadeColor[7]))
		style[ButtonHoverBorderColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[ButtonPressedBorderColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[ToggleDefaultBorderColor] = int(raylib.GetHexValue(fadeColor[7]))
		style[ToggleHoverBorderColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[TogglePressedBorderColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[ToggleActiveBorderColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[SliderBorderColor] = int(raylib.GetHexValue(fadeColor[7]))
		style[SliderbarBorderColor] = int(raylib.GetHexValue(fadeColor[7]))
		style[ProgressbarBorderColor] = int(raylib.GetHexValue(fadeColor[7]))
		style[SpinnerLabelBorderColor] = int(raylib.GetHexValue(fadeColor[7]))
		style[SpinnerDefaultButtonBorderColor] = int(raylib.GetHexValue(fadeColor[7]))
		style[SpinnerHoverButtonBorderColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[SpinnerPressedButtonBorderColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[ComboboxDefaultBorderColor] = int(raylib.GetHexValue(fadeColor[7]))
		style[ComboboxHoverBorderColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[ComboboxPressedBorderColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[ComboboxPressedListBorderColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[CheckboxDefaultBorderColor] = int(raylib.GetHexValue(fadeColor[7]))
		style[CheckboxHoverBorderColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[CheckboxClickBorderColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[TextboxBorderColor] = int(raylib.GetHexValue(fadeColor[7]))
	} else if guiProperty == GlobalTextColor {
		baseColor := raylib.GetColor(int32(value))
		fadeColor := make([]raylib.Color, numColorSamples)

		for i := 0; i < numColorSamples; i++ {
			fadeColor[i] = colorMultiply(baseColor, 1.0-float32(i)/float32(numColorSamples-1))
		}

		style[GlobalTextColor] = value
		style[LabelTextColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[ButtonDefaultTextColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[ButtonHoverTextColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[ButtonPressedTextColor] = int(raylib.GetHexValue(fadeColor[5]))
		style[ToggleDefaultTextColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[ToggleHoverTextColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[TogglePressedTextColor] = int(raylib.GetHexValue(fadeColor[5]))
		style[ToggleActiveTextColor] = int(raylib.GetHexValue(fadeColor[5]))
		style[SpinnerDefaultSymbolColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[SpinnerDefaultTextColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[SpinnerHoverSymbolColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[SpinnerHoverTextColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[SpinnerPressedSymbolColor] = int(raylib.GetHexValue(fadeColor[5]))
		style[SpinnerPressedTextColor] = int(raylib.GetHexValue(fadeColor[5]))
		style[ComboboxDefaultTextColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[ComboboxDefaultListTextColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[ComboboxHoverTextColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[ComboboxHoverListTextColor] = int(raylib.GetHexValue(fadeColor[8]))
		style[ComboboxPressedTextColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[ComboboxPressedListTextColor] = int(raylib.GetHexValue(fadeColor[4]))
		style[TextboxTextColor] = int(raylib.GetHexValue(fadeColor[9]))
		style[TextboxLineColor] = int(raylib.GetHexValue(fadeColor[6]))
	} else {
		style[guiProperty] = value
	}
}

// Get one style property
func GetStyleProperty(guiProperty Property) int {
	return style[int(guiProperty)]
}

func colorMultiply(baseColor raylib.Color, value float32) raylib.Color {
	multColor := baseColor

	if value > 1.0 {
		value = 1.0
	} else if value < 0.0 {
		value = 0.0
	}

	multColor.R += uint8((255 - float32(multColor.R)) * value)
	multColor.G += uint8((255 - float32(multColor.G)) * value)
	multColor.B += uint8((255 - float32(multColor.B)) * value)

	return multColor
}
