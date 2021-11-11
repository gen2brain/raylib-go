// GUI element appearance can be dynamically configured through Property values, the set of which
// forms a theme called the Style.
package raygui

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
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

	// Add new properties above.
	NumProperties
)

// GUI property names (to read/write style text files)
var propertyName = [NumProperties]string{
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

// Current GUI style (default light).
var style = [NumProperties]int64{
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

// LoadGuiStyle will load a GUI style from a file. See SaveGuiStyle.
func LoadGuiStyle(fileName string) {
	file, err := rl.OpenAsset(fileName)
	if err != nil {
		rl.TraceLog(rl.LogWarning, "[%s] GUI style file could not be opened: %w", fileName, err)
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

// SaveGuiStyle will write the current GUI style to a file in a format suitable for loading via LoadGuiStyle.
func SaveGuiStyle(fileName string) {
	var styleFile string
	for i := 0; i < len(propertyName); i++ {
		styleFile += fmt.Sprintf("%-40s0x%x\n", propertyName[i], GetStyleProperty(Property(i)))
	}

	if err := ioutil.WriteFile(fileName, []byte(styleFile), 0644); err != nil {
		rl.TraceLog(rl.LogWarning, "[%s] GUI style file could not be written: %w", fileName, err)
	}
}

// SetStyleProperty - Set one style property
func SetStyleProperty(guiProperty Property, value int64) {
	style[guiProperty] = value
}

// SetStyleColor - Set one style property to a color value
func SetStyleColor(guiProperty Property, value rl.Color) {
	style[guiProperty] = int64(rl.ColorToInt(value))
}

// GetStyleProperty - Get one style property
func GetStyleProperty(guiProperty Property) int64 {
	return style[int(guiProperty)]
}

// BackgroundColor will return the current background color
func BackgroundColor() rl.Color {
	return rl.GetColor(uint(style[GlobalBackgroundColor]))
}

// LinesColor will return the current color for lines
func LinesColor() rl.Color {
	return rl.GetColor(uint(style[GlobalLinesColor]))
}

// TextColor will return the current color for normal state
func TextColor() rl.Color {
	return rl.GetColor(uint(style[GlobalTextColor]))
}

// GetStyle32 will return the int32 for a given property of the current style
func GetStyle32(property Property) int32 {
	return int32(style[property])
}

// GetPropColor will return the Color value for a given property of the current style
func GetStyleColor(property Property) rl.Color {
	return rl.GetColor(uint(style[property]))
}
