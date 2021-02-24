// Package raygui - Simple and easy-to-use IMGUI (immediate mode GUI API) library
package raygui

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gen2brain/raylib-go/raylib"
)

// Property - GUI property
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
	boundsWidth
	boundsHeight
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
	CheckboxDefaultActiveColor
	CheckboxInsideWidth
	TextboxBorderWidth
	TextboxBorderColor
	TextboxInsideColor
	TextboxTextColor
	TextboxLineColor
	TextboxTextFontsize
)

// GUI controls states
const (
	Disabled = iota
	Normal
	Focused
	Pressed
)

// Current GUI style (default light)
var style = []int64{
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
var propertyName = []string{
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

// For spinner
var (
	framesCounter  int
	framesCounter2 int
	valueSpeed     bool
)

// BackgroundColor - Get background color
func BackgroundColor() rl.Color {
	return rl.GetColor(int32(style[GlobalBackgroundColor]))
}

// LinesColor - Get lines color
func LinesColor() rl.Color {
	return rl.GetColor(int32(style[GlobalLinesColor]))
}

// TextColor - Get text color for normal state
func TextColor() rl.Color {
	return rl.GetColor(int32(style[GlobalTextColor]))
}

// Label - Label element, show text
func Label(bounds rl.Rectangle, text string) {
	LabelEx(bounds, text, rl.GetColor(int32(style[LabelTextColor])), rl.NewColor(0, 0, 0, 0), rl.NewColor(0, 0, 0, 0))
}

// LabelEx - Label element extended, configurable colors
func LabelEx(bounds rl.Rectangle, text string, textColor, border, inner rl.Color) {
	b := bounds.ToInt32()
	// Update control
	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(text, textHeight)

	if b.Width < textWidth {
		b.Width = textWidth + int32(style[LabelTextPadding])
	}
	if b.Height < textHeight {
		b.Height = textHeight + int32(style[LabelTextPadding])/2
	}

	// Draw control
	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, border)
	rl.DrawRectangle(b.X+int32(style[LabelBorderWidth]), b.Y+int32(style[LabelBorderWidth]), b.Width-(2*int32(style[LabelBorderWidth])), b.Height-(2*int32(style[LabelBorderWidth])), inner)
	rl.DrawText(text, b.X+((b.Width/2)-(textWidth/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, textColor)
}

// Button - Button element, returns true when clicked
func Button(bounds rl.Rectangle, text string) bool {
	b := bounds.ToInt32()
	state := Normal
	mousePoint := rl.GetMousePosition()
	clicked := false

	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(text, textHeight)

	// Update control
	if b.Width < textWidth {
		b.Width = textWidth + int32(style[ButtonTextPadding])
	}

	if b.Height < textHeight {
		b.Height = textHeight + int32(style[ButtonTextPadding])/2
	}

	if rl.CheckCollisionPointRec(mousePoint, bounds) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			state = Pressed
		} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) || rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			clicked = true
		} else {
			state = Focused
		}
	}

	// Draw control
	switch state {
	case Normal:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ButtonDefaultBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ButtonBorderWidth]), b.Y+int32(style[ButtonBorderWidth]), b.Width-(2*int32(style[ButtonBorderWidth])), b.Height-(2*int32(style[ButtonBorderWidth])), rl.GetColor(int32(style[ButtonDefaultInsideColor])))
		rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ButtonDefaultTextColor])))
		break

	case Focused:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ButtonHoverBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ButtonBorderWidth]), b.Y+int32(style[ButtonBorderWidth]), b.Width-(2*int32(style[ButtonBorderWidth])), b.Height-(2*int32(style[ButtonBorderWidth])), rl.GetColor(int32(style[ButtonHoverInsideColor])))
		rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ButtonHoverTextColor])))
		break

	case Pressed:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ButtonPressedBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ButtonBorderWidth]), b.Y+int32(style[ButtonBorderWidth]), b.Width-(2*int32(style[ButtonBorderWidth])), b.Height-(2*int32(style[ButtonBorderWidth])), rl.GetColor(int32(style[ButtonPressedInsideColor])))
		rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ButtonPressedTextColor])))
		break

	default:
		break
	}

	if clicked {
		return true
	}

	return false
}

// ToggleButton - Toggle Button element, returns true when active
func ToggleButton(bounds rl.Rectangle, text string, active bool) bool {
	b := bounds.ToInt32()
	state := Normal
	mousePoint := rl.GetMousePosition()

	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(text, textHeight)

	// Update control
	if b.Width < textWidth {
		b.Width = textWidth + int32(style[ToggleTextPadding])
	}
	if b.Height < textHeight {
		b.Height = textHeight + int32(style[ToggleTextPadding])/2
	}

	if rl.CheckCollisionPointRec(mousePoint, bounds) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			state = Pressed
		} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) || rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			state = Normal
			active = !active
		} else {
			state = Focused
		}
	}

	// Draw control
	switch state {
	case Normal:
		if active {
			rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ToggleActiveBorderColor])))
			rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[ToggleActiveInsideColor])))
			rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ToggleDefaultTextColor])))
		} else {
			rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ToggleDefaultBorderColor])))
			rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[ToggleDefaultInsideColor])))
			rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ToggleDefaultTextColor])))
		}
		break
	case Focused:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ToggleHoverBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[ToggleHoverInsideColor])))
		rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[ToggleHoverTextColor])))
		break
	case Pressed:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[TogglePressedBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[TogglePressedInsideColor])))
		rl.DrawText(text, b.X+((b.Width/2)-(rl.MeasureText(text, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, rl.GetColor(int32(style[TogglePressedTextColor])))
		break
	default:
		break
	}

	return active
}

// ToggleGroup - Toggle Group element, returns toggled button index
func ToggleGroup(bounds rl.Rectangle, toggleText []string, active int) int {
	for i := 0; i < len(toggleText); i++ {
		if i == active {
			ToggleButton(rl.NewRectangle(bounds.X+float32(i)*(bounds.Width+float32(style[TogglegroupPadding])), bounds.Y, bounds.Width, bounds.Height), toggleText[i], true)
		} else if ToggleButton(rl.NewRectangle(bounds.X+float32(i)*(bounds.Width+float32(style[TogglegroupPadding])), bounds.Y, bounds.Width, bounds.Height), toggleText[i], false) {
			active = i
		}
	}

	return active
}

// ComboBox draws a simplified version of a ComboBox allowing the user to select a string
// from a list accompanied by an N/M counter. The widget does not provide a drop-down/completion
// or any input support.
func ComboBox(bounds rl.Rectangle, comboText []string, active int) int {
	// Reject invalid selections and disable rendering.
	comboCount := len(comboText)
	if active < 0 || active >= comboCount {
		rl.TraceLog(rl.LogWarning, "ComboBox active expects 0 <= active <= %d", comboCount)
		return -1
	}

	activeText := comboText[active]

	// style sizing.
	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(activeText, textHeight)
	borderWidth := int32(style[ComboboxBorderWidth])
	textPadding := int32(style[ToggleTextPadding])

	b := bounds.ToInt32()
	if b.Width < textWidth {
		b.Width = textWidth + textPadding
		bounds.Width = float32(b.Width)
	}
	if b.Height < textHeight {
		b.Height = textHeight + textPadding
		bounds.Height = float32(b.Height)
	}

	// Identify what the counter is going to look like with max digits so we don't resize it.
	clickWidth := rl.MeasureText(fmt.Sprintf("%d/%d", comboCount, comboCount), b.Height)

	click := rl.NewRectangle(bounds.X+bounds.Width+float32(style[ComboboxPadding]), bounds.Y, float32(clickWidth), float32(b.Height))
	c := click.ToInt32()
	mousePoint := rl.GetMousePosition()
	state := Normal
	if rl.CheckCollisionPointRec(mousePoint, bounds) || rl.CheckCollisionPointRec(mousePoint, click) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			state = Pressed
		} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) || rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			state = Pressed
		} else {
			state = Focused
		}
	}

	// Draw control
	var borderColor, insideColor, listColor, textColor rl.Color

	switch state {
	case Normal:
		borderColor = rl.GetColor(int32(style[ComboboxDefaultBorderColor]))
		insideColor = rl.GetColor(int32(style[ComboboxDefaultInsideColor]))
		listColor = rl.GetColor(int32(style[ComboboxDefaultListTextColor]))
		textColor = rl.GetColor(int32(style[ComboboxDefaultTextColor]))

	case Focused:
		borderColor = rl.GetColor(int32(style[ComboboxHoverBorderColor]))
		insideColor = rl.GetColor(int32(style[ComboboxHoverInsideColor]))
		listColor = rl.GetColor(int32(style[ComboboxHoverListTextColor]))
		textColor = rl.GetColor(int32(style[ComboboxHoverTextColor]))

	case Pressed:
		borderColor = rl.GetColor(int32(style[ComboboxPressedBorderColor]))
		insideColor = rl.GetColor(int32(style[ComboboxPressedInsideColor]))
		listColor = rl.GetColor(int32(style[ComboboxPressedListTextColor]))
		textColor = rl.GetColor(int32(style[ComboboxPressedTextColor]))

	default:
		rl.TraceLog(rl.LogWarning, "ComboBox in unrecognized state %d", state)
		return -1
	}

	// Render the box itself
	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, borderColor)
	rl.DrawRectangle(b.X+borderWidth, b.Y+borderWidth, b.Width-(2*borderWidth), b.Height-(2*borderWidth), insideColor)
	rl.DrawText(activeText, b.X+((b.Width/2)-(rl.MeasureText(activeText, textHeight)/2)), b.Y+((b.Height/2)-(textHeight/2)), textHeight, textColor)

	// Render the accompanying "clicks" box showing the element counter.
	rl.DrawRectangle(c.X, c.Y, c.Width, c.Height, borderColor)
	rl.DrawRectangle(c.X+borderWidth, c.Y+borderWidth, c.Width-(2*borderWidth), c.Height-(2*borderWidth), insideColor)
	companionText := fmt.Sprintf("%d/%d", active+1, comboCount)
	rl.DrawText(companionText, c.X+((c.Width/2)-(rl.MeasureText(companionText, textHeight)/2)), c.Y+((c.Height/2)-(textHeight/2)), textHeight, listColor)

	if rl.CheckCollisionPointRec(mousePoint, bounds) || rl.CheckCollisionPointRec(mousePoint, click) {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			active++
			if active >= comboCount {
				active = 0
			}
		}
	}

	return active
}

// CheckBox - Check Box element, returns true when active
func CheckBox(bounds rl.Rectangle, checked bool) bool {
	b := bounds.ToInt32()
	state := Normal
	mousePoint := rl.GetMousePosition()

	// Update control
	if rl.CheckCollisionPointRec(mousePoint, bounds) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			state = Pressed
		} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) || rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			state = Normal
			checked = !checked
		} else {
			state = Focused
		}
	}

	// Draw control
	switch state {
	case Normal:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[CheckboxDefaultBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[CheckboxDefaultInsideColor])))
		break
	case Focused:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[CheckboxHoverBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[CheckboxHoverInsideColor])))
		break
	case Pressed:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[CheckboxClickBorderColor])))
		rl.DrawRectangle(b.X+int32(style[ToggleBorderWidth]), b.Y+int32(style[ToggleBorderWidth]), b.Width-(2*int32(style[ToggleBorderWidth])), b.Height-(2*int32(style[ToggleBorderWidth])), rl.GetColor(int32(style[CheckboxClickInsideColor])))
		break
	default:
		break
	}

	if checked {
		rl.DrawRectangle(b.X+int32(style[CheckboxInsideWidth]), b.Y+int32(style[CheckboxInsideWidth]), b.Width-(2*int32(style[CheckboxInsideWidth])), b.Height-(2*int32(style[CheckboxInsideWidth])), rl.GetColor(int32(style[CheckboxDefaultActiveColor])))
	}

	return checked
}

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
	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[SliderBorderColor])))
	rl.DrawRectangle(b.X+int32(style[SliderBorderWidth]), b.Y+int32(style[SliderBorderWidth]), b.Width-(2*int32(style[SliderBorderWidth])), b.Height-(2*int32(style[SliderBorderWidth])), rl.GetColor(int32(style[SliderInsideColor])))

	switch state {
	case Normal:
		rl.DrawRectangle(sliderButton.X, sliderButton.Y, sliderButton.Width, sliderButton.Height, rl.GetColor(int32(style[SliderDefaultColor])))
		break
	case Focused:
		rl.DrawRectangle(sliderButton.X, sliderButton.Y, sliderButton.Width, sliderButton.Height, rl.GetColor(int32(style[SliderHoverColor])))
		break
	case Pressed:
		rl.DrawRectangle(sliderButton.X, sliderButton.Y, sliderButton.Width, sliderButton.Height, rl.GetColor(int32(style[SliderActiveColor])))
		break
	default:
		break
	}

	return minValue + (maxValue-minValue)*sliderPos
}

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
	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[SliderbarBorderColor])))
	rl.DrawRectangle(b.X+int32(style[SliderBorderWidth]), b.Y+int32(style[SliderBorderWidth]), b.Width-(2*int32(style[SliderBorderWidth])), b.Height-(2*int32(style[SliderBorderWidth])), rl.GetColor(int32(style[SliderbarInsideColor])))

	switch state {
	case Normal:
		rl.DrawRectangle(sliderBar.X, sliderBar.Y, sliderBar.Width, sliderBar.Height, rl.GetColor(int32(style[SliderbarDefaultColor])))
		break
	case Focused:
		rl.DrawRectangle(sliderBar.X, sliderBar.Y, sliderBar.Width, sliderBar.Height, rl.GetColor(int32(style[SliderbarHoverColor])))
		break
	case Pressed:
		rl.DrawRectangle(sliderBar.X, sliderBar.Y, sliderBar.Width, sliderBar.Height, rl.GetColor(int32(style[SliderbarActiveColor])))
		break
	default:
		break
	}

	if minValue < 0 && maxValue > 0 {
		rl.DrawRectangle((b.X+int32(style[SliderBorderWidth]))-int32(minValue*(float32(b.Width-(int32(style[SliderBorderWidth])*2))/maxValue)), sliderBar.Y, 1, sliderBar.Height, rl.GetColor(int32(style[SliderbarZeroLineColor])))
	}

	return fixedValue + minValue
}

// ProgressBar - Progress Bar element, shows current progress value
func ProgressBar(bounds rl.Rectangle, value float32) {
	b := bounds.ToInt32()
	if value > 1.0 {
		value = 1.0
	} else if value < 0.0 {
		value = 0.0
	}

	progressBar := rl.RectangleInt32{b.X + int32(style[ProgressbarBorderWidth]), b.Y + int32(style[ProgressbarBorderWidth]), b.Width - (int32(style[ProgressbarBorderWidth]) * 2), b.Height - (int32(style[ProgressbarBorderWidth]) * 2)}
	progressValue := rl.RectangleInt32{b.X + int32(style[ProgressbarBorderWidth]), b.Y + int32(style[ProgressbarBorderWidth]), int32(value * float32(b.Width-int32(style[ProgressbarBorderWidth])*2)), b.Height - (int32(style[ProgressbarBorderWidth]) * 2)}

	// Draw control
	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ProgressbarBorderColor])))
	rl.DrawRectangle(progressBar.X, progressBar.Y, progressBar.Width, progressBar.Height, rl.GetColor(int32(style[ProgressbarInsideColor])))
	rl.DrawRectangle(progressValue.X, progressValue.Y, progressValue.Width, progressValue.Height, rl.GetColor(int32(style[ProgressbarProgressColor])))
}

// Spinner - Spinner element, returns selected value
func Spinner(bounds rl.Rectangle, value, minValue, maxValue int) int {
	b := bounds.ToInt32()
	state := Normal

	mousePoint := rl.GetMousePosition()
	labelBoxBound := rl.RectangleInt32{b.X + b.Width/4 + 1, b.Y, b.Width / 2, b.Height}
	leftButtonBound := rl.RectangleInt32{b.X, b.Y, b.Width / 4, b.Height}
	rightButtonBound := rl.RectangleInt32{b.X + b.Width - b.Width/4 + 1, b.Y, b.Width / 4, b.Height}

	textHeight := int32(style[GlobalTextFontsize])
	textWidth := rl.MeasureText(fmt.Sprintf("%d", value), textHeight)

	buttonSide := 0

	// Update control
	if rl.CheckCollisionPointRec(mousePoint, leftButtonBound.ToFloat32()) || rl.CheckCollisionPointRec(mousePoint, rightButtonBound.ToFloat32()) || rl.CheckCollisionPointRec(mousePoint, labelBoxBound.ToFloat32()) {
		if rl.IsKeyDown(rl.KeyLeft) {
			state = Pressed
			buttonSide = 1

			if value > minValue {
				value--
			}
		} else if rl.IsKeyDown(rl.KeyRight) {
			state = Pressed
			buttonSide = 2

			if value < maxValue {
				value++
			}
		}
	}

	if rl.CheckCollisionPointRec(mousePoint, leftButtonBound.ToFloat32()) {
		buttonSide = 1
		state = Focused

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			if !valueSpeed {
				if value > minValue {
					value--
				}
				valueSpeed = true
			} else {
				framesCounter++
			}

			state = Pressed

			if value > minValue {
				if framesCounter >= 30 {
					value--
				}
			}
		}
	} else if rl.CheckCollisionPointRec(mousePoint, rightButtonBound.ToFloat32()) {
		buttonSide = 2
		state = Focused

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			if !valueSpeed {
				if value < maxValue {
					value++
				}
				valueSpeed = true
			} else {
				framesCounter++
			}

			state = Pressed

			if value < maxValue {
				if framesCounter >= 30 {
					value++
				}
			}
		}
	} else if !rl.CheckCollisionPointRec(mousePoint, labelBoxBound.ToFloat32()) {
		buttonSide = 0
	}

	if rl.IsMouseButtonUp(rl.MouseLeftButton) {
		valueSpeed = false
		framesCounter = 0
	}

	// Draw control
	switch state {
	case Normal:
		rl.DrawRectangle(leftButtonBound.X, leftButtonBound.Y, leftButtonBound.Width, leftButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
		rl.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

		rl.DrawRectangle(rightButtonBound.X, rightButtonBound.Y, rightButtonBound.Width, rightButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
		rl.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

		rl.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(rl.MeasureText("+", textHeight))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))
		rl.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(rl.MeasureText("-", textHeight))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))

		rl.DrawRectangle(labelBoxBound.X, labelBoxBound.Y, labelBoxBound.Width, labelBoxBound.Height, rl.GetColor(int32(style[SpinnerLabelBorderColor])))
		rl.DrawRectangle(labelBoxBound.X+1, labelBoxBound.Y+1, labelBoxBound.Width-2, labelBoxBound.Height-2, rl.GetColor(int32(style[SpinnerLabelInsideColor])))

		rl.DrawText(fmt.Sprintf("%d", value), labelBoxBound.X+(labelBoxBound.Width/2-textWidth/2), labelBoxBound.Y+(labelBoxBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultTextColor])))
		break
	case Focused:
		if buttonSide == 1 {
			rl.DrawRectangle(leftButtonBound.X, leftButtonBound.Y, leftButtonBound.Width, leftButtonBound.Height, rl.GetColor(int32(style[SpinnerHoverButtonBorderColor])))
			rl.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, rl.GetColor(int32(style[SpinnerHoverButtonInsideColor])))

			rl.DrawRectangle(rightButtonBound.X, rightButtonBound.Y, rightButtonBound.Width, rightButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			rl.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			rl.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(rl.MeasureText("+", textHeight))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerHoverSymbolColor])))
			rl.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(rl.MeasureText("-", textHeight))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))
		} else if buttonSide == 2 {
			rl.DrawRectangle(leftButtonBound.X, leftButtonBound.Y, leftButtonBound.Width, leftButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			rl.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			rl.DrawRectangle(rightButtonBound.X, rightButtonBound.Y, rightButtonBound.Width, rightButtonBound.Height, rl.GetColor(int32(style[SpinnerHoverButtonBorderColor])))
			rl.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, rl.GetColor(int32(style[SpinnerHoverButtonInsideColor])))

			rl.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(rl.MeasureText("+", textHeight))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))
			rl.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(rl.MeasureText("-", textHeight))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerHoverSymbolColor])))
		}

		rl.DrawRectangle(labelBoxBound.X, labelBoxBound.Y, labelBoxBound.Width, labelBoxBound.Height, rl.GetColor(int32(style[SpinnerLabelBorderColor])))
		rl.DrawRectangle(labelBoxBound.X+1, labelBoxBound.Y+1, labelBoxBound.Width-2, labelBoxBound.Height-2, rl.GetColor(int32(style[SpinnerLabelInsideColor])))

		rl.DrawText(fmt.Sprintf("%d", value), labelBoxBound.X+(labelBoxBound.Width/2-textWidth/2), labelBoxBound.Y+(labelBoxBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerHoverTextColor])))
		break
	case Pressed:
		if buttonSide == 1 {
			rl.DrawRectangle(leftButtonBound.X, leftButtonBound.Y, leftButtonBound.Width, leftButtonBound.Height, rl.GetColor(int32(style[SpinnerPressedButtonBorderColor])))
			rl.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, rl.GetColor(int32(style[SpinnerPressedButtonInsideColor])))

			rl.DrawRectangle(rightButtonBound.X, rightButtonBound.Y, rightButtonBound.Width, rightButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			rl.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			rl.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(rl.MeasureText("+", textHeight))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerPressedSymbolColor])))
			rl.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(rl.MeasureText("-", textHeight))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))
		} else if buttonSide == 2 {
			rl.DrawRectangle(leftButtonBound.X, leftButtonBound.Y, leftButtonBound.Width, leftButtonBound.Height, rl.GetColor(int32(style[SpinnerDefaultButtonBorderColor])))
			rl.DrawRectangle(leftButtonBound.X+2, leftButtonBound.Y+2, leftButtonBound.Width-4, leftButtonBound.Height-4, rl.GetColor(int32(style[SpinnerDefaultButtonInsideColor])))

			rl.DrawRectangle(rightButtonBound.X, rightButtonBound.Y, rightButtonBound.Width, rightButtonBound.Height, rl.GetColor(int32(style[SpinnerPressedButtonBorderColor])))
			rl.DrawRectangle(rightButtonBound.X+2, rightButtonBound.Y+2, rightButtonBound.Width-4, rightButtonBound.Height-4, rl.GetColor(int32(style[SpinnerPressedButtonInsideColor])))

			rl.DrawText("-", leftButtonBound.X+(leftButtonBound.Width/2-(rl.MeasureText("+", textHeight))/2), leftButtonBound.Y+(leftButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerDefaultSymbolColor])))
			rl.DrawText("+", rightButtonBound.X+(rightButtonBound.Width/2-(rl.MeasureText("-", textHeight))/2), rightButtonBound.Y+(rightButtonBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerPressedSymbolColor])))
		}

		rl.DrawRectangle(labelBoxBound.X, labelBoxBound.Y, labelBoxBound.Width, labelBoxBound.Height, rl.GetColor(int32(style[SpinnerLabelBorderColor])))
		rl.DrawRectangle(labelBoxBound.X+1, labelBoxBound.Y+1, labelBoxBound.Width-2, labelBoxBound.Height-2, rl.GetColor(int32(style[SpinnerLabelInsideColor])))

		rl.DrawText(fmt.Sprintf("%d", value), labelBoxBound.X+(labelBoxBound.Width/2-textWidth/2), labelBoxBound.Y+(labelBoxBound.Height/2-(textHeight/2)), textHeight, rl.GetColor(int32(style[SpinnerPressedTextColor])))
		break
	default:
		break
	}

	return value
}

// TextBox - Text Box element, updates input text
func TextBox(bounds rl.Rectangle, text string) string {
	b := bounds.ToInt32()
	state := Normal

	mousePoint := rl.GetMousePosition()
	letter := int32(-1)

	// Update control
	if rl.CheckCollisionPointRec(mousePoint, bounds) {
		state = Focused // NOTE: PRESSED state is not used on this control

		framesCounter2++

		letter = rl.GetKeyPressed()
		if letter != -1 {
			if letter >= 32 && letter < 127 {
				text = fmt.Sprintf("%s%c", text, letter)
			}
		}

		if rl.IsKeyPressed(rl.KeyBackspace) {
			if len(text) > 0 {
				text = text[:len(text)-1]
			}
		}
	}

	// Draw control
	switch state {
	case Normal:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[TextboxBorderColor])))
		rl.DrawRectangle(b.X+int32(style[TextboxBorderWidth]), b.Y+int32(style[TextboxBorderWidth]), b.Width-(int32(style[TextboxBorderWidth])*2), b.Height-(int32(style[TextboxBorderWidth])*2), rl.GetColor(int32(style[TextboxInsideColor])))
		rl.DrawText(text, b.X+2, b.Y+int32(style[TextboxBorderWidth])+b.Height/2-int32(style[TextboxTextFontsize])/2, int32(style[TextboxTextFontsize]), rl.GetColor(int32(style[TextboxTextColor])))
		break
	case Focused:
		rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, rl.GetColor(int32(style[ToggleActiveBorderColor])))
		rl.DrawRectangle(b.X+int32(style[TextboxBorderWidth]), b.Y+int32(style[TextboxBorderWidth]), b.Width-(int32(style[TextboxBorderWidth])*2), b.Height-(int32(style[TextboxBorderWidth])*2), rl.GetColor(int32(style[TextboxInsideColor])))
		rl.DrawText(text, b.X+2, b.Y+int32(style[TextboxBorderWidth])+b.Height/2-int32(style[TextboxTextFontsize])/2, int32(style[TextboxTextFontsize]), rl.GetColor(int32(style[TextboxTextColor])))

		if (framesCounter2/20)%2 == 0 && rl.CheckCollisionPointRec(mousePoint, bounds) {
			rl.DrawRectangle(b.X+4+rl.MeasureText(text, int32(style[GlobalTextFontsize])), b.Y+2, 1, b.Height-4, rl.GetColor(int32(style[TextboxLineColor])))
		}
		break
	case Pressed:
		break
	default:
		break
	}

	return text
}

// SaveGuiStyle - Save GUI style file
func SaveGuiStyle(fileName string) {
	var styleFile string
	for i := 0; i < len(propertyName); i++ {
		styleFile += fmt.Sprintf("%-40s0x%x\n", propertyName[i], GetStyleProperty(Property(i)))
	}

	ioutil.WriteFile(fileName, []byte(styleFile), 0644)
}

// LoadGuiStyle - Load GUI style file
func LoadGuiStyle(fileName string) {
	file, err := rl.OpenAsset(fileName)
	if err != nil {
		rl.TraceLog(rl.LogWarning, "[%s] GUI style file could not be opened", fileName)
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
					style[i] = int64(v)
				}
			}
		}
	}
}

// SetStyleProperty - Set one style property
func SetStyleProperty(guiProperty Property, value int64) {
	style[guiProperty] = value
}

// GetStyleProperty - Get one style property
func GetStyleProperty(guiProperty Property) int64 {
	return style[int(guiProperty)]
}
