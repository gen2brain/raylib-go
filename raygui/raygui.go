package raygui

/*
#cgo CFLAGS: -DRAYGUI_IMPLEMENTATION -Wno-unused-result
#include "raygui.h"
#include <stdlib.h>
*/
import "C"

import (
	"strings"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCROLLBAR_LEFT_SIDE  = 0
	SCROLLBAR_RIGHT_SIDE = 1
)

// GuiStyleProp - Style property
type GuiStyleProp struct {
	controlId     uint16
	propertyId    uint16
	propertyValue uint32
}

// STATE_NORMAL - Gui control state
const (
	STATE_NORMAL   int32 = 0
	STATE_FOCUSED        = 1
	STATE_PRESSED        = 2
	STATE_DISABLED       = 3
)

// GuiState .
type GuiState = int32

// TEXT_ALIGN_LEFT - Gui control text alignment
const (
	TEXT_ALIGN_LEFT   int32 = 0
	TEXT_ALIGN_CENTER       = 1
	TEXT_ALIGN_RIGHT        = 2
)

// GuiTextAlignment .
type GuiTextAlignment = int32

// DEFAULT - Gui controls
const (
	DEFAULT     int32 = 0
	LABEL             = 1
	BUTTON            = 2
	TOGGLE            = 3
	SLIDER            = 4
	PROGRESSBAR       = 5
	CHECKBOX          = 6
	COMBOBOX          = 7
	DROPDOWNBOX       = 8
	TEXTBOX           = 9
	VALUEBOX          = 10
	SPINNER           = 11
	LISTVIEW          = 12
	COLORPICKER       = 13
	SCROLLBAR         = 14
	STATUSBAR         = 15
)

// GuiControl .
type GuiControl = int32

// Gui base properties for every control
// NOTE: RAYGUI_MAX_PROPS_BASE properties (by default 16 properties)
const (
	BORDER_COLOR_NORMAL   int32 = 0
	BASE_COLOR_NORMAL           = 1
	TEXT_COLOR_NORMAL           = 2
	BORDER_COLOR_FOCUSED        = 3
	BASE_COLOR_FOCUSED          = 4
	TEXT_COLOR_FOCUSED          = 5
	BORDER_COLOR_PRESSED        = 6
	BASE_COLOR_PRESSED          = 7
	TEXT_COLOR_PRESSED          = 8
	BORDER_COLOR_DISABLED       = 9
	BASE_COLOR_DISABLED         = 10
	TEXT_COLOR_DISABLED         = 11
	BORDER_WIDTH                = 12
	TEXT_PADDING                = 13
	TEXT_ALIGNMENT              = 14
	RESERVED                    = 15
)

// GuiControlProperty .
type GuiControlProperty = int32

// DEFAULT extended properties
// NOTE: Those properties are common to all controls or global
const (
	TEXT_SIZE        int32 = 16
	TEXT_SPACING           = 17
	LINE_COLOR             = 18
	BACKGROUND_COLOR       = 19
)

// GuiDefaultProperty .
type GuiDefaultProperty = int32

// GROUP_PADDING .
const (
	GROUP_PADDING int32 = 16
)

// GuiToggleProperty .
type GuiToggleProperty = int32

const (
	// Slider size of internal bar
	SLIDER_WIDTH int32 = 16
	// Slider/SliderBar internal bar padding
	SLIDER_PADDING = 17
)

// GuiSliderProperty .
type GuiSliderProperty = int32

const (
	// ProgressBar internal padding
	PROGRESS_PADDING int32 = 16
)

// GuiProgressBarProperty .
type GuiProgressBarProperty = int32

const (
	ARROWS_SIZE           int32 = 16
	ARROWS_VISIBLE              = 17
	SCROLL_SLIDER_PADDING       = 18
	SCROLL_SLIDER_SIZE          = 19
	SCROLL_PADDING              = 20
	SCROLL_SPEED                = 21
)

// GuiScrollBarProperty .
type GuiScrollBarProperty = int32

const (
	CHECK_PADDING int32 = 16
)

// GuiCheckBoxProperty .
type GuiCheckBoxProperty = int32

const (
	// ComboBox right button width
	COMBO_BUTTON_WIDTH int32 = 16
	// ComboBox button separation
	COMBO_BUTTON_SPACING = 17
)

// GuiComboBoxProperty .
type GuiComboBoxProperty = int32

const (
	// DropdownBox arrow separation from border and items
	ARROW_PADDING int32 = 16
	// DropdownBox items separation
	DROPDOWN_ITEMS_SPACING = 17
)

// GuiDropdownBoxProperty .
type GuiDropdownBoxProperty = int32

const (
	// TextBox/TextBoxMulti/ValueBox/Spinner inner text padding
	TEXT_INNER_PADDING int32 = 16
	// TextBoxMulti lines separation
	TEXT_LINES_SPACING = 17
)

// GuiTextBoxProperty .
type GuiTextBoxProperty = int32

const (
	// Spinner left/right buttons width
	SPIN_BUTTON_WIDTH int32 = 16
	// Spinner buttons separation
	SPIN_BUTTON_SPACING = 17
)

// GuiSpinnerProperty .
type GuiSpinnerProperty = int32

const (
	// ListView items height
	LIST_ITEMS_HEIGHT int32 = 16
	// ListView items separation
	LIST_ITEMS_SPACING = 17
	// ListView scrollbar size (usually width)
	SCROLLBAR_WIDTH = 18
	// ListView scrollbar side (0-left, 1-right)
	SCROLLBAR_SIDE = 19
)

// GuiListViewProperty .
type GuiListViewProperty = int32

const (
	COLOR_SELECTOR_SIZE int32 = 16
	// rl.ColorPicker right hue bar width
	HUEBAR_WIDTH = 17
	// rl.ColorPicker right hue bar separation from panel
	HUEBAR_PADDING = 18
	// rl.ColorPicker right hue bar selector height
	HUEBAR_SELECTOR_HEIGHT = 19
	// rl.ColorPicker right hue bar selector overflow
	HUEBAR_SELECTOR_OVERFLOW = 20
)

// GuiColorPickerProperty .
type GuiColorPickerProperty = int32

// GuiEnable - Enable gui controls (global state)
func Enable() {
	C.GuiEnable()
}

// GuiDisable - Disable gui controls (global state)
func Disable() {
	C.GuiDisable()
}

// GuiLock - Lock gui controls (global state)
func Lock() {
	C.GuiLock()
}

// GuiUnlock - Unlock gui controls (global state)
func Unlock() {
	C.GuiUnlock()
}

// GuiIsLocked - Check if gui is locked (global state)
func IsLocked() bool {
	return bool(C.GuiIsLocked())
}

// GuiFade - Set gui controls alpha (global state), alpha goes from 0.0f to 1.0f
func Fade(alpha float32) {
	calpha := C.float(alpha)
	C.GuiFade(calpha)
}

// GuiSetState - Set gui state (global state)
func SetState(state int32) {
	cstate := C.int(state)
	C.GuiSetState(cstate)
}

// GuiGetState - Get gui state (global state)
func GetState() int32 {
	return int32(C.GuiGetState())
}

// GuiSetStyle .
func SetStyle(control int32, property int32, value int64) {
	ccontrol := C.int(control)
	cproperty := C.int(property)
	cvalue := C.int(value)
	C.GuiSetStyle(ccontrol, cproperty, cvalue)
}

// GuiGetStyle - Get one style property
func GetStyle(control int32, property int32) int64 {
	ccontrol := C.int(control)
	cproperty := C.int(property)
	return int64(C.GuiGetStyle(ccontrol, cproperty))
}

// GuiWindowBox - Window Box control, shows a window that can be closed
func WindowBox(bounds rl.Rectangle, title string) bool {
	var cbounds C.struct_Rectangle
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	return bool(C.GuiWindowBox(cbounds, ctitle))
}

// GuiGroupBox - Group Box control with text name
func GroupBox(bounds rl.Rectangle, text string) {
	var cbounds C.struct_Rectangle
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.GuiGroupBox(cbounds, ctext)
}

// GuiLine - Line separator control, could contain text
func Line(bounds rl.Rectangle, text string) {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.GuiLine(cbounds, ctext)
}

// GuiPanel - Panel control, useful to group controls
func Panel(bounds rl.Rectangle, text string) {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.GuiPanel(cbounds, ctext)
}

// ScrollPanel control
func ScrollPanel(bounds rl.Rectangle, text string, content rl.Rectangle, scroll *rl.Vector2) rl.Rectangle {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	var ccontent C.struct_Rectangle
	ccontent.x = C.float(content.X)
	ccontent.y = C.float(content.Y)
	ccontent.width = C.float(content.Width)
	ccontent.height = C.float(content.Height)
	var cscroll C.struct_Vector2
	cscroll.x = C.float(scroll.X)
	cscroll.y = C.float(scroll.Y)
	defer func() {
		scroll.X = float32(cscroll.x)
		scroll.Y = float32(cscroll.y)
	}()

	var res C.struct_Rectangle
	res = C.GuiScrollPanel(cbounds, ctext, ccontent, &cscroll)

	var goRes rl.Rectangle
	goRes.X = float32(res.x)
	goRes.Y = float32(res.y)
	goRes.Width = float32(res.width)
	goRes.Height = float32(res.height)

	return goRes
}

// ScrollBar control (used by GuiScrollPanel())
func ScrollBar(bounds rl.Rectangle, value, minValue, maxValue int32) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)

	cvalue := C.int(value)
	cminValue := C.int(minValue)
	cmaxValue := C.int(maxValue)

	return int32(C.GuiScrollBar(cbounds, cvalue, cminValue, cmaxValue))
}

// Label control, shows text
func Label(bounds rl.Rectangle, text string) {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.GuiLabel(cbounds, ctext)
}

// Button control, returns true when clicked
func Button(bounds rl.Rectangle, text string) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	return bool(C.GuiButton(cbounds, ctext))
}

// LabelButton control, show true when clicked
func LabelButton(bounds rl.Rectangle, text string) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	return bool(C.GuiLabelButton(cbounds, ctext))
}

// Toggle control, returns true when active
func Toggle(bounds rl.Rectangle, text string, active bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cactive := C.bool(active)
	return bool(C.GuiToggle(cbounds, ctext, cactive))
}

// ToggleGroup control, returns active toggle index
func ToggleGroup(bounds rl.Rectangle, text string, active int32) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cactive := C.int(active)
	return int32(C.GuiToggleGroup(cbounds, ctext, cactive))
}

// CheckBox control, returns true when active
func CheckBox(bounds rl.Rectangle, text string, checked bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cchecked := C.bool(checked)
	return bool(C.GuiCheckBox(cbounds, ctext, cchecked))
}

// ComboBox control, returns selected item index
func ComboBox(bounds rl.Rectangle, text string, active int32) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cactive := C.int(active)
	return int32(C.GuiComboBox(cbounds, ctext, cactive))
}

// Spinner control, returns selected value
func Spinner(bounds rl.Rectangle, text string, value *int32, minValue, maxValue int, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	if value == nil {
		value = new(int32)
	}
	cvalue := C.int(*value)
	defer func() {
		*value = int32(cvalue)
	}()

	cminValue := C.int(minValue)
	cmaxValue := C.int(maxValue)
	ceditMode := C.bool(editMode)

	return bool(C.GuiSpinner(cbounds, ctext, &cvalue, cminValue, cmaxValue, ceditMode))
}

// Slider control
func Slider(bounds rl.Rectangle, textLeft string, textRight string, value float32, minValue float32, maxValue float32) float32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctextLeft := C.CString(textLeft)
	defer C.free(unsafe.Pointer(ctextLeft))
	ctextRight := C.CString(textRight)
	defer C.free(unsafe.Pointer(ctextRight))
	cvalue := C.float(value)
	cminValue := C.float(minValue)
	cmaxValue := C.float(maxValue)
	return float32(C.GuiSlider(cbounds, ctextLeft, ctextRight, cvalue, cminValue, cmaxValue))
}

// SliderBar control, returns selected value
func SliderBar(bounds rl.Rectangle, textLeft string, textRight string, value float32, minValue float32, maxValue float32) float32 {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	ctextLeft := C.CString(textLeft)
	defer C.free(unsafe.Pointer(ctextLeft))
	ctextRight := C.CString(textRight)
	defer C.free(unsafe.Pointer(ctextRight))
	cvalue := C.float(value)
	cminValue := C.float(minValue)
	cmaxValue := C.float(maxValue)
	return float32(C.GuiSliderBar(cbounds, ctextLeft, ctextRight, cvalue, cminValue, cmaxValue))
}

// ProgressBar control, shows current progress value
func ProgressBar(bounds rl.Rectangle, textLeft string, textRight string, value float32, minValue float32, maxValue float32) float32 {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	ctextLeft := C.CString(textLeft)
	defer C.free(unsafe.Pointer(ctextLeft))
	ctextRight := C.CString(textRight)
	defer C.free(unsafe.Pointer(ctextRight))
	cvalue := C.float(value)
	cminValue := C.float(minValue)
	cmaxValue := C.float(maxValue)
	return float32(C.GuiProgressBar(cbounds, ctextLeft, ctextRight, cvalue, cminValue, cmaxValue))
}

// StatusBar control, shows info text
func StatusBar(bounds rl.Rectangle, text string) {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.GuiStatusBar(cbounds, ctext)
}

// DummyRec control for placeholders
func DummyRec(bounds rl.Rectangle, text string) {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.GuiDummyRec(cbounds, ctext)
}

// Grid control, returns mouse cell position
func Grid(bounds rl.Rectangle, text string, spacing float32, subdivs int32) rl.Vector2 {
	var cbounds C.struct_Rectangle
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cspacing := C.float(spacing)
	csubdivs := C.int(subdivs)
	cResult := C.GuiGrid(cbounds, ctext, cspacing, csubdivs)
	var goRes rl.Vector2
	goRes.Y = float32(cResult.y)
	goRes.X = float32(cResult.x)
	return goRes
}

// ListView control, returns selected list item index
func ListView(bounds rl.Rectangle, text string, scrollIndex *int32, active int32) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	if scrollIndex == nil {
		scrollIndex = new(int32)
	}
	cscrollIndex := C.int(*scrollIndex)
	defer func() {
		*scrollIndex = int32(cscrollIndex)
	}()

	cactive := C.int(active)

	return int32(C.GuiListView(cbounds, ctext, &cscrollIndex, cactive))
}

// MessageBox control, displays a message
func MessageBox(bounds rl.Rectangle, title string, message string, buttons string) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	cmessage := C.CString(message)
	defer C.free(unsafe.Pointer(cmessage))
	cbuttons := C.CString(buttons)
	defer C.free(unsafe.Pointer(cbuttons))
	return int32(C.GuiMessageBox(cbounds, ctitle, cmessage, cbuttons))
}

// ColorPicker control (multiple color controls)
func ColorPicker(bounds rl.Rectangle, text string, color rl.Color) rl.Color {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	var ccolor C.struct_Color
	ccolor.r = C.uchar(color.R)
	ccolor.g = C.uchar(color.G)
	ccolor.b = C.uchar(color.B)
	ccolor.a = C.uchar(color.A)
	cResult := C.GuiColorPicker(cbounds, ctext, ccolor)
	var goRes rl.Color
	goRes.A = byte(cResult.a)
	goRes.R = byte(cResult.r)
	goRes.G = byte(cResult.g)
	goRes.B = byte(cResult.b)
	return goRes
}

// ColorPanel control
func ColorPanel(bounds rl.Rectangle, text string, color rl.Color) rl.Color {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	var ccolor C.struct_Color
	ccolor.b = C.uchar(color.B)
	ccolor.a = C.uchar(color.A)
	ccolor.r = C.uchar(color.R)
	ccolor.g = C.uchar(color.G)
	cResult := C.GuiColorPanel(cbounds, ctext, ccolor)
	var goRes rl.Color
	goRes.A = byte(cResult.a)
	goRes.R = byte(cResult.r)
	goRes.G = byte(cResult.g)
	goRes.B = byte(cResult.b)
	return goRes
}

// ColorBarAlpha control
func ColorBarAlpha(bounds rl.Rectangle, text string, alpha float32) float32 {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	calpha := C.float(alpha)
	return float32(C.GuiColorBarAlpha(cbounds, ctext, calpha))
}

// ColorBarHue control
func ColorBarHue(bounds rl.Rectangle, text string, value float32) float32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cvalue := C.float(value)
	return float32(C.GuiColorBarHue(cbounds, ctext, cvalue))
}

// DropdownBox control
func DropdownBox(bounds rl.Rectangle, text string, active *int32, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	if active == nil {
		active = new(int32)
	}
	cactive := C.int(*active)
	defer func() {
		*active = int32(cactive)
	}()

	ceditMode := C.bool(editMode)

	return bool(C.GuiDropdownBox(cbounds, ctext, &cactive, ceditMode))
}

// ValueBox control, updates input text with numbers
func ValueBox(bounds rl.Rectangle, text string, value *int32, minValue, maxValue int, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	if value == nil {
		value = new(int32)
	}
	cvalue := C.int(*value)
	defer func() {
		*value = int32(cvalue)
	}()

	cminValue := C.int(minValue)
	cmaxValue := C.int(maxValue)
	ceditMode := C.bool(editMode)

	return bool(C.GuiValueBox(cbounds, ctext, &cvalue, cminValue, cmaxValue, ceditMode))
}

// TextBox control, updates input text
func TextBox(bounds rl.Rectangle, text *string, textSize int, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)

	bs := []byte(*text)
	if len(bs) == 0 {
		bs = []byte{byte(0)}
	}
	if 0 < len(bs) && bs[len(bs)-1] != byte(0) { // minimalize allocation
		bs = append(bs, byte(0)) // for next input symbols
	}
	ctext := (*C.char)(unsafe.Pointer(&bs[0]))
	defer func() {
		*text = strings.TrimSpace(strings.Trim(string(bs), "\x00"))
		// no need : C.free(unsafe.Pointer(ctext))
	}()

	ctextSize := C.int(textSize)
	ceditMode := C.bool(editMode)

	return bool(C.GuiTextBox(cbounds, ctext, ctextSize, ceditMode))
}

// LoadStyle file over global style variable (.rgs)
func LoadStyle(fileName string) {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	C.GuiLoadStyle(cfileName)
}

// LoadStyleDefault over global style
func LoadStyleDefault() {
	C.GuiLoadStyleDefault()
}

// IconText gets text with icon id prepended (if supported)
func IconText(iconId int32, text string) string {
	ciconId := C.int(iconId)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	return C.GoString(C.GuiIconText(ciconId, ctext))
}

// Icons enumeration
const (
	ICON_NONE                    int32 = 0
	ICON_FOLDER_FILE_OPEN              = 1
	ICON_FILE_SAVE_CLASSIC             = 2
	ICON_FOLDER_OPEN                   = 3
	ICON_FOLDER_SAVE                   = 4
	ICON_FILE_OPEN                     = 5
	ICON_FILE_SAVE                     = 6
	ICON_FILE_EXPORT                   = 7
	ICON_FILE_ADD                      = 8
	ICON_FILE_DELETE                   = 9
	ICON_FILETYPE_TEXT                 = 10
	ICON_FILETYPE_AUDIO                = 11
	ICON_FILETYPE_IMAGE                = 12
	ICON_FILETYPE_PLAY                 = 13
	ICON_FILETYPE_VIDEO                = 14
	ICON_FILETYPE_INFO                 = 15
	ICON_FILE_COPY                     = 16
	ICON_FILE_CUT                      = 17
	ICON_FILE_PASTE                    = 18
	ICON_CURSOR_HAND                   = 19
	ICON_CURSOR_POINTER                = 20
	ICON_CURSOR_CLASSIC                = 21
	ICON_PENCIL                        = 22
	ICON_PENCIL_BIG                    = 23
	ICON_BRUSH_CLASSIC                 = 24
	ICON_BRUSH_PAINTER                 = 25
	ICON_WATER_DROP                    = 26
	ICON_COLOR_PICKER                  = 27
	ICON_RUBBER                        = 28
	ICON_COLOR_BUCKET                  = 29
	ICON_TEXT_T                        = 30
	ICON_TEXT_A                        = 31
	ICON_SCALE                         = 32
	ICON_RESIZE                        = 33
	ICON_FILTER_POINT                  = 34
	ICON_FILTER_BILINEAR               = 35
	ICON_CROP                          = 36
	ICON_CROP_ALPHA                    = 37
	ICON_SQUARE_TOGGLE                 = 38
	ICON_SYMMETRY                      = 39
	ICON_SYMMETRY_HORIZONTAL           = 40
	ICON_SYMMETRY_VERTICAL             = 41
	ICON_LENS                          = 42
	ICON_LENS_BIG                      = 43
	ICON_EYE_ON                        = 44
	ICON_EYE_OFF                       = 45
	ICON_FILTER_TOP                    = 46
	ICON_FILTER                        = 47
	ICON_TARGET_POINT                  = 48
	ICON_TARGET_SMALL                  = 49
	ICON_TARGET_BIG                    = 50
	ICON_TARGET_MOVE                   = 51
	ICON_CURSOR_MOVE                   = 52
	ICON_CURSOR_SCALE                  = 53
	ICON_CURSOR_SCALE_RIGHT            = 54
	ICON_CURSOR_SCALE_LEFT             = 55
	ICON_UNDO                          = 56
	ICON_REDO                          = 57
	ICON_REREDO                        = 58
	ICON_MUTATE                        = 59
	ICON_ROTATE                        = 60
	ICON_REPEAT                        = 61
	ICON_SHUFFLE                       = 62
	ICON_EMPTYBOX                      = 63
	ICON_TARGET                        = 64
	ICON_TARGET_SMALL_FILL             = 65
	ICON_TARGET_BIG_FILL               = 66
	ICON_TARGET_MOVE_FILL              = 67
	ICON_CURSOR_MOVE_FILL              = 68
	ICON_CURSOR_SCALE_FILL             = 69
	ICON_CURSOR_SCALE_RIGHT_FILL       = 70
	ICON_CURSOR_SCALE_LEFT_FILL        = 71
	ICON_UNDO_FILL                     = 72
	ICON_REDO_FILL                     = 73
	ICON_REREDO_FILL                   = 74
	ICON_MUTATE_FILL                   = 75
	ICON_ROTATE_FILL                   = 76
	ICON_REPEAT_FILL                   = 77
	ICON_SHUFFLE_FILL                  = 78
	ICON_EMPTYBOX_SMALL                = 79
	ICON_BOX                           = 80
	ICON_BOX_TOP                       = 81
	ICON_BOX_TOP_RIGHT                 = 82
	ICON_BOX_RIGHT                     = 83
	ICON_BOX_BOTTOM_RIGHT              = 84
	ICON_BOX_BOTTOM                    = 85
	ICON_BOX_BOTTOM_LEFT               = 86
	ICON_BOX_LEFT                      = 87
	ICON_BOX_TOP_LEFT                  = 88
	ICON_BOX_CENTER                    = 89
	ICON_BOX_CIRCLE_MASK               = 90
	ICON_POT                           = 91
	ICON_ALPHA_MULTIPLY                = 92
	ICON_ALPHA_CLEAR                   = 93
	ICON_DITHERING                     = 94
	ICON_MIPMAPS                       = 95
	ICON_BOX_GRID                      = 96
	ICON_GRID                          = 97
	ICON_BOX_CORNERS_SMALL             = 98
	ICON_BOX_CORNERS_BIG               = 99
	ICON_FOUR_BOXES                    = 100
	ICON_GRID_FILL                     = 101
	ICON_BOX_MULTISIZE                 = 102
	ICON_ZOOM_SMALL                    = 103
	ICON_ZOOM_MEDIUM                   = 104
	ICON_ZOOM_BIG                      = 105
	ICON_ZOOM_ALL                      = 106
	ICON_ZOOM_CENTER                   = 107
	ICON_BOX_DOTS_SMALL                = 108
	ICON_BOX_DOTS_BIG                  = 109
	ICON_BOX_CONCENTRIC                = 110
	ICON_BOX_GRID_BIG                  = 111
	ICON_OK_TICK                       = 112
	ICON_CROSS                         = 113
	ICON_ARROW_LEFT                    = 114
	ICON_ARROW_RIGHT                   = 115
	ICON_ARROW_DOWN                    = 116
	ICON_ARROW_UP                      = 117
	ICON_ARROW_LEFT_FILL               = 118
	ICON_ARROW_RIGHT_FILL              = 119
	ICON_ARROW_DOWN_FILL               = 120
	ICON_ARROW_UP_FILL                 = 121
	ICON_AUDIO                         = 122
	ICON_FX                            = 123
	ICON_WAVE                          = 124
	ICON_WAVE_SINUS                    = 125
	ICON_WAVE_SQUARE                   = 126
	ICON_WAVE_TRIANGULAR               = 127
	ICON_CROSS_SMALL                   = 128
	ICON_PLAYER_PREVIOUS               = 129
	ICON_PLAYER_PLAY_BACK              = 130
	ICON_PLAYER_PLAY                   = 131
	ICON_PLAYER_PAUSE                  = 132
	ICON_PLAYER_STOP                   = 133
	ICON_PLAYER_NEXT                   = 134
	ICON_PLAYER_RECORD                 = 135
	ICON_MAGNET                        = 136
	ICON_LOCK_CLOSE                    = 137
	ICON_LOCK_OPEN                     = 138
	ICON_CLOCK                         = 139
	ICON_TOOLS                         = 140
	ICON_GEAR                          = 141
	ICON_GEAR_BIG                      = 142
	ICON_BIN                           = 143
	ICON_HAND_POINTER                  = 144
	ICON_LASER                         = 145
	ICON_COIN                          = 146
	ICON_EXPLOSION                     = 147
	ICON_1UP                           = 148
	ICON_PLAYER                        = 149
	ICON_PLAYER_JUMP                   = 150
	ICON_KEY                           = 151
	ICON_DEMON                         = 152
	ICON_TEXT_POPUP                    = 153
	ICON_GEAR_EX                       = 154
	ICON_CRACK                         = 155
	ICON_CRACK_POINTS                  = 156
	ICON_STAR                          = 157
	ICON_DOOR                          = 158
	ICON_EXIT                          = 159
	ICON_MODE_2D                       = 160
	ICON_MODE_3D                       = 161
	ICON_CUBE                          = 162
	ICON_CUBE_FACE_TOP                 = 163
	ICON_CUBE_FACE_LEFT                = 164
	ICON_CUBE_FACE_FRONT               = 165
	ICON_CUBE_FACE_BOTTOM              = 166
	ICON_CUBE_FACE_RIGHT               = 167
	ICON_CUBE_FACE_BACK                = 168
	ICON_CAMERA                        = 169
	ICON_SPECIAL                       = 170
	ICON_LINK_NET                      = 171
	ICON_LINK_BOXES                    = 172
	ICON_LINK_MULTI                    = 173
	ICON_LINK                          = 174
	ICON_LINK_BROKE                    = 175
	ICON_TEXT_NOTES                    = 176
	ICON_NOTEBOOK                      = 177
	ICON_SUITCASE                      = 178
	ICON_SUITCASE_ZIP                  = 179
	ICON_MAILBOX                       = 180
	ICON_MONITOR                       = 181
	ICON_PRINTER                       = 182
	ICON_PHOTO_CAMERA                  = 183
	ICON_PHOTO_CAMERA_FLASH            = 184
	ICON_HOUSE                         = 185
	ICON_HEART                         = 186
	ICON_CORNER                        = 187
	ICON_VERTICAL_BARS                 = 188
	ICON_VERTICAL_BARS_FILL            = 189
	ICON_LIFE_BARS                     = 190
	ICON_INFO                          = 191
	ICON_CROSSLINE                     = 192
	ICON_HELP                          = 193
	ICON_FILETYPE_ALPHA                = 194
	ICON_FILETYPE_HOME                 = 195
	ICON_LAYERS_VISIBLE                = 196
	ICON_LAYERS                        = 197
	ICON_WINDOW                        = 198
	ICON_HIDPI                         = 199
	ICON_FILETYPE_BINARY               = 200
	ICON_HEX                           = 201
	ICON_SHIELD                        = 202
	ICON_FILE_NEW                      = 203
	ICON_FOLDER_ADD                    = 204
	ICON_ALARM                         = 205
	ICON_CPU                           = 206
	ICON_ROM                           = 207
	ICON_STEP_OVER                     = 208
	ICON_STEP_INTO                     = 209
	ICON_STEP_OUT                      = 210
	ICON_RESTART                       = 211
	ICON_BREAKPOINT_ON                 = 212
	ICON_BREAKPOINT_OFF                = 213
	ICON_BURGER_MENU                   = 214
	ICON_CASE_SENSITIVE                = 215
	ICON_REG_EXP                       = 216
	ICON_FOLDER                        = 217
	ICON_FILE                          = 218
	ICON_219                           = 219
	ICON_220                           = 220
	ICON_221                           = 221
	ICON_222                           = 222
	ICON_223                           = 223
	ICON_224                           = 224
	ICON_225                           = 225
	ICON_226                           = 226
	ICON_227                           = 227
	ICON_228                           = 228
	ICON_229                           = 229
	ICON_230                           = 230
	ICON_231                           = 231
	ICON_232                           = 232
	ICON_233                           = 233
	ICON_234                           = 234
	ICON_235                           = 235
	ICON_236                           = 236
	ICON_237                           = 237
	ICON_238                           = 238
	ICON_239                           = 239
	ICON_240                           = 240
	ICON_241                           = 241
	ICON_242                           = 242
	ICON_243                           = 243
	ICON_244                           = 244
	ICON_245                           = 245
	ICON_246                           = 246
	ICON_247                           = 247
	ICON_248                           = 248
	ICON_249                           = 249
	ICON_250                           = 250
	ICON_251                           = 251
	ICON_252                           = 252
	ICON_253                           = 253
	ICON_254                           = 254
	ICON_255                           = 255
)

// TextInputBox control, ask for text
func TextInputBox(bounds rl.Rectangle, title, message, buttons string, text *string, textMaxSize int32, secretViewActive *int32) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)

	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	cmessage := C.CString(message)
	defer C.free(unsafe.Pointer(cmessage))

	cbuttons := C.CString(buttons)
	defer C.free(unsafe.Pointer(cbuttons))

	bs := []byte(*text)
	if len(bs) == 0 {
		bs = []byte{byte(0)}
	}
	if 0 < len(bs) && bs[len(bs)-1] != byte(0) { // minimalize allocation
		bs = append(bs, byte(0)) // for next input symbols
	}
	ctext := (*C.char)(unsafe.Pointer(&bs[0]))
	defer func() {
		*text = strings.TrimSpace(strings.Trim(string(bs), "\x00"))
		// no need : C.free(unsafe.Pointer(ctext))
	}()

	ctextMaxSize := C.int(textMaxSize)

	if secretViewActive == nil {
		secretViewActive = new(int32)
	}
	csecretViewActive := C.int(*secretViewActive)
	defer func() {
		*secretViewActive = int32(csecretViewActive)
	}()

	return int32(C.GuiTextInputBox(cbounds, ctitle, cmessage, cbuttons, ctext, ctextMaxSize, &csecretViewActive))
}

// TextBoxMulti control with multiple lines
func TextBoxMulti(bounds rl.Rectangle, text *string, textSize int32, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)

	bs := []byte(*text)
	if len(bs) == 0 {
		bs = []byte{byte(0)}
	}
	if 0 < len(bs) && bs[len(bs)-1] != byte(0) { // minimalize allocation
		bs = append(bs, byte(0)) // for next input symbols
	}
	ctext := (*C.char)(unsafe.Pointer(&bs[0]))
	defer func() {
		*text = strings.TrimSpace(strings.Trim(string(bs), "\x00"))
		// no need : C.free(unsafe.Pointer(ctext))
	}()

	ctextSize := (C.int)(textSize)
	ceditMode := (C.bool)(editMode)

	return bool(C.GuiTextBoxMulti(cbounds, ctext, ctextSize, ceditMode))
}

// ListViewEx control with extended parameters
func ListViewEx(bounds rl.Rectangle, text []string, focus, scrollIndex *int32, active int32) int32 {
	// int GuiListViewEx(Rectangle bounds, const char **text, int count, int *focus, int *scrollIndex, int active)

	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)

	ctext := NewCStringArrayFromSlice(text)
	defer ctext.Free()

	count := C.int(len(text))

	if focus == nil {
		focus = new(int32)
	}
	cfocus := C.int(*focus)
	defer func() {
		*focus = int32(cfocus)
	}()

	if scrollIndex == nil {
		scrollIndex = new(int32)
	}
	cscrollIndex := C.int(*scrollIndex)
	defer func() {
		*scrollIndex = int32(cscrollIndex)
	}()

	cactive := C.int(active)

	return int32(C.GuiListViewEx(cbounds, (**C.char)(ctext.Pointer), count, &cfocus, &cscrollIndex, cactive))
}

// TabBar control
func TabBar(bounds rl.Rectangle, text []string, active *int32) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)

	ctext := NewCStringArrayFromSlice(text)
	defer ctext.Free()

	count := C.int(len(text))

	if active == nil {
		active = new(int32)
	}
	cactive := C.int(*active)
	defer func() {
		*active = int32(cactive)
	}()
	return int32(C.GuiTabBar(cbounds, (**C.char)(ctext.Pointer), count, &cactive))
}

// SetFont - set custom font (global state)
func SetFont(font rl.Font) {
	cfont := (*C.Font)(unsafe.Pointer(&font))
	C.GuiSetFont(*cfont)
}

// GetFont - get custom font (global state)
func GetFont() rl.Font {
	ret := C.GuiGetFont()
	ptr := unsafe.Pointer(&ret)
	return *(*rl.Font)(ptr)
}
