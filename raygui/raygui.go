package raygui

/*
#define RAYGUI_IMPLEMENTATION
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

// Gui control state
const (
	STATE_NORMAL   int32 = 0
	STATE_FOCUSED        = 1
	STATE_PRESSED        = 2
	STATE_DISABLED       = 3
)

// GuiState .
type GuiState = int32

// Gui control text alignment
const (
	TEXT_ALIGN_LEFT   int32 = 0
	TEXT_ALIGN_CENTER       = 1
	TEXT_ALIGN_RIGHT        = 2
)

// GuiTextAlignment .
type GuiTextAlignment = int32

// Gui control text alignment vertical
const (
	TEXT_ALIGN_TOP    int32 = 0
	TEXT_ALIGN_MIDDLE       = 1
	TEXT_ALIGN_BOTTOM       = 2
)

// GuiTextWrapMode .
type GuiTextWrapMode = int32

// Gui control text wrap mode
// NOTE: Useful for multiline text
const (
	TEXT_WRAP_NONE int32 = 0
	TEXT_WRAP_CHAR       = 1
	TEXT_WRAP_WORD       = 2
)

// GuiTextAlignmentVertical .
type GuiTextAlignmentVertical = int32

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
)

// GuiControlProperty .
type GuiControlProperty = int32

// DEFAULT extended properties
// NOTE: Those properties are common to all controls or global
const (
	TEXT_SIZE               int32 = 16
	TEXT_SPACING                  = 17
	LINE_COLOR                    = 18
	BACKGROUND_COLOR              = 19
	TEXT_LINE_SPACING             = 20
	TEXT_ALIGNMENT_VERTICAL       = 21
	TEXT_WRAP_MODE                = 22
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
func Fade(color rl.Color, alpha float32) {
	ccolor := (*C.Color)(unsafe.Pointer(&color))
	calpha := C.float(alpha)
	C.GuiFade(*ccolor, calpha)
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
	return C.GuiWindowBox(cbounds, ctitle) != 0
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

// Panel - Panel control, useful to group controls
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

// ScrollPanel control - Scroll Panel control
func ScrollPanel(bounds rl.Rectangle, text string, content rl.Rectangle, scroll *rl.Vector2, view *rl.Rectangle) int32 {
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
	var cview C.struct_Rectangle
	cview.x = C.float(view.X)
	cview.y = C.float(view.Y)
	cview.width = C.float(view.Width)
	cview.height = C.float(view.Height)

	res := C.GuiScrollPanel(cbounds, ctext, ccontent, &cscroll, &cview)

	return int32(res)
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
	return C.GuiButton(cbounds, ctext) != 0
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
	return C.GuiLabelButton(cbounds, ctext) != 0
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
	C.GuiToggle(cbounds, ctext, &cactive)
	return bool(cactive)
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
	C.GuiToggleGroup(cbounds, ctext, &cactive)
	return int32(cactive)
}

// ToggleSlider control, returns true when clicked
func ToggleSlider(bounds rl.Rectangle, text string, active int32) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cactive := C.int(active)
	C.GuiToggleSlider(cbounds, ctext, &cactive)
	return int32(cactive)
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
	C.GuiCheckBox(cbounds, ctext, &cchecked)
	return bool(cchecked)
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
	C.GuiComboBox(cbounds, ctext, &cactive)
	return int32(cactive)
}

// Spinner control, returns selected value
func Spinner(bounds rl.Rectangle, text string, value *int32, minValue, maxValue int, editMode bool) int32 {
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

	C.GuiSpinner(cbounds, ctext, &cvalue, cminValue, cmaxValue, ceditMode)
	return int32(cvalue)
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
	C.GuiSlider(cbounds, ctextLeft, ctextRight, &cvalue, cminValue, cmaxValue)
	return float32(cvalue)
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
	C.GuiSliderBar(cbounds, ctextLeft, ctextRight, &cvalue, cminValue, cmaxValue)
	return float32(cvalue)
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
	C.GuiProgressBar(cbounds, ctextLeft, ctextRight, &cvalue, cminValue, cmaxValue)
	return float32(cvalue)
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
func Grid(bounds rl.Rectangle, text string, spacing float32, subdivs int32, mouseCell *rl.Vector2) int32 {
	var cbounds C.struct_Rectangle
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cspacing := C.float(spacing)
	csubdivs := C.int(subdivs)
	var cmouseCell C.struct_Vector2
	cmouseCell.x = C.float(mouseCell.X)
	cmouseCell.y = C.float(mouseCell.Y)
	res := C.GuiGrid(cbounds, ctext, cspacing, csubdivs, &cmouseCell)
	return int32(res)
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

	C.GuiListView(cbounds, ctext, &cscrollIndex, &cactive)
	return int32(cactive)
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
	C.GuiColorPicker(cbounds, ctext, &ccolor)
	var goRes rl.Color
	goRes.A = byte(ccolor.a)
	goRes.R = byte(ccolor.r)
	goRes.G = byte(ccolor.g)
	goRes.B = byte(ccolor.b)
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
	C.GuiColorPanel(cbounds, ctext, &ccolor)
	var goRes rl.Color
	goRes.A = byte(ccolor.a)
	goRes.R = byte(ccolor.r)
	goRes.G = byte(ccolor.g)
	goRes.B = byte(ccolor.b)
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
	C.GuiColorBarAlpha(cbounds, ctext, &calpha)
	return float32(calpha)
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
	C.GuiColorBarHue(cbounds, ctext, &cvalue)
	return float32(cvalue)
}

// ColorPickerHSV - Color Picker control that avoids conversion to RGB on each call (multiple color controls)
func ColorPickerHSV(bounds rl.Rectangle, text string, colorHSV *rl.Vector3) int32 {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)

	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	var ccolorHSV C.struct_Vector3
	ccolorHSV.x = C.float(colorHSV.X)
	ccolorHSV.y = C.float(colorHSV.Y)
	ccolorHSV.z = C.float(colorHSV.Z)
	defer func() {
		colorHSV.X = float32(ccolorHSV.x)
		colorHSV.Y = float32(ccolorHSV.y)
		colorHSV.Z = float32(ccolorHSV.z)
	}()

	return int32(C.GuiColorPickerHSV(cbounds, ctext, &ccolorHSV))
}

// ColorPanelHSV - Color Panel control that returns HSV color value, used by GuiColorPickerHSV()
func ColorPanelHSV(bounds rl.Rectangle, text string, colorHSV *rl.Vector3) int32 {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)

	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	var ccolorHSV C.struct_Vector3
	ccolorHSV.x = C.float(colorHSV.X)
	ccolorHSV.y = C.float(colorHSV.Y)
	ccolorHSV.z = C.float(colorHSV.Z)
	defer func() {
		colorHSV.X = float32(ccolorHSV.x)
		colorHSV.Y = float32(ccolorHSV.y)
		colorHSV.Z = float32(ccolorHSV.z)
	}()

	return int32(C.GuiColorPanelHSV(cbounds, ctext, &ccolorHSV))
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

	return C.GuiDropdownBox(cbounds, ctext, &cactive, ceditMode) != 0
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

	return C.GuiValueBox(cbounds, ctext, &cvalue, cminValue, cmaxValue, ceditMode) != 0
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

	return C.GuiTextBox(cbounds, ctext, ctextSize, ceditMode) != 0
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

// LoadStyleFromMemory - Load style from memory (binary only)
func LoadStyleFromMemory(data []byte) {
	C.GuiLoadStyleFromMemory((*C.uchar)(unsafe.Pointer(&data[0])), C.int(len(data)))
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
	ICON_NONE int32 = iota
	ICON_FOLDER_FILE_OPEN
	ICON_FILE_SAVE_CLASSIC
	ICON_FOLDER_OPEN
	ICON_FOLDER_SAVE
	ICON_FILE_OPEN
	ICON_FILE_SAVE
	ICON_FILE_EXPORT
	ICON_FILE_ADD
	ICON_FILE_DELETE
	ICON_FILETYPE_TEXT
	ICON_FILETYPE_AUDIO
	ICON_FILETYPE_IMAGE
	ICON_FILETYPE_PLAY
	ICON_FILETYPE_VIDEO
	ICON_FILETYPE_INFO
	ICON_FILE_COPY
	ICON_FILE_CUT
	ICON_FILE_PASTE
	ICON_CURSOR_HAND
	ICON_CURSOR_POINTER
	ICON_CURSOR_CLASSIC
	ICON_PENCIL
	ICON_PENCIL_BIG
	ICON_BRUSH_CLASSIC
	ICON_BRUSH_PAINTER
	ICON_WATER_DROP
	ICON_COLOR_PICKER
	ICON_RUBBER
	ICON_COLOR_BUCKET
	ICON_TEXT_T
	ICON_TEXT_A
	ICON_SCALE
	ICON_RESIZE
	ICON_FILTER_POINT
	ICON_FILTER_BILINEAR
	ICON_CROP
	ICON_CROP_ALPHA
	ICON_SQUARE_TOGGLE
	ICON_SYMMETRY
	ICON_SYMMETRY_HORIZONTAL
	ICON_SYMMETRY_VERTICAL
	ICON_LENS
	ICON_LENS_BIG
	ICON_EYE_ON
	ICON_EYE_OFF
	ICON_FILTER_TOP
	ICON_FILTER
	ICON_TARGET_POINT
	ICON_TARGET_SMALL
	ICON_TARGET_BIG
	ICON_TARGET_MOVE
	ICON_CURSOR_MOVE
	ICON_CURSOR_SCALE
	ICON_CURSOR_SCALE_RIGHT
	ICON_CURSOR_SCALE_LEFT
	ICON_UNDO
	ICON_REDO
	ICON_REREDO
	ICON_MUTATE
	ICON_ROTATE
	ICON_REPEAT
	ICON_SHUFFLE
	ICON_EMPTYBOX
	ICON_TARGET
	ICON_TARGET_SMALL_FILL
	ICON_TARGET_BIG_FILL
	ICON_TARGET_MOVE_FILL
	ICON_CURSOR_MOVE_FILL
	ICON_CURSOR_SCALE_FILL
	ICON_CURSOR_SCALE_RIGHT_FILL
	ICON_CURSOR_SCALE_LEFT_FILL
	ICON_UNDO_FILL
	ICON_REDO_FILL
	ICON_REREDO_FILL
	ICON_MUTATE_FILL
	ICON_ROTATE_FILL
	ICON_REPEAT_FILL
	ICON_SHUFFLE_FILL
	ICON_EMPTYBOX_SMALL
	ICON_BOX
	ICON_BOX_TOP
	ICON_BOX_TOP_RIGHT
	ICON_BOX_RIGHT
	ICON_BOX_BOTTOM_RIGHT
	ICON_BOX_BOTTOM
	ICON_BOX_BOTTOM_LEFT
	ICON_BOX_LEFT
	ICON_BOX_TOP_LEFT
	ICON_BOX_CENTER
	ICON_BOX_CIRCLE_MASK
	ICON_POT
	ICON_ALPHA_MULTIPLY
	ICON_ALPHA_CLEAR
	ICON_DITHERING
	ICON_MIPMAPS
	ICON_BOX_GRID
	ICON_GRID
	ICON_BOX_CORNERS_SMALL
	ICON_BOX_CORNERS_BIG
	ICON_FOUR_BOXES
	ICON_GRID_FILL
	ICON_BOX_MULTISIZE
	ICON_ZOOM_SMALL
	ICON_ZOOM_MEDIUM
	ICON_ZOOM_BIG
	ICON_ZOOM_ALL
	ICON_ZOOM_CENTER
	ICON_BOX_DOTS_SMALL
	ICON_BOX_DOTS_BIG
	ICON_BOX_CONCENTRIC
	ICON_BOX_GRID_BIG
	ICON_OK_TICK
	ICON_CROSS
	ICON_ARROW_LEFT
	ICON_ARROW_RIGHT
	ICON_ARROW_DOWN
	ICON_ARROW_UP
	ICON_ARROW_LEFT_FILL
	ICON_ARROW_RIGHT_FILL
	ICON_ARROW_DOWN_FILL
	ICON_ARROW_UP_FILL
	ICON_AUDIO
	ICON_FX
	ICON_WAVE
	ICON_WAVE_SINUS
	ICON_WAVE_SQUARE
	ICON_WAVE_TRIANGULAR
	ICON_CROSS_SMALL
	ICON_PLAYER_PREVIOUS
	ICON_PLAYER_PLAY_BACK
	ICON_PLAYER_PLAY
	ICON_PLAYER_PAUSE
	ICON_PLAYER_STOP
	ICON_PLAYER_NEXT
	ICON_PLAYER_RECORD
	ICON_MAGNET
	ICON_LOCK_CLOSE
	ICON_LOCK_OPEN
	ICON_CLOCK
	ICON_TOOLS
	ICON_GEAR
	ICON_GEAR_BIG
	ICON_BIN
	ICON_HAND_POINTER
	ICON_LASER
	ICON_COIN
	ICON_EXPLOSION
	ICON_1UP
	ICON_PLAYER
	ICON_PLAYER_JUMP
	ICON_KEY
	ICON_DEMON
	ICON_TEXT_POPUP
	ICON_GEAR_EX
	ICON_CRACK
	ICON_CRACK_POINTS
	ICON_STAR
	ICON_DOOR
	ICON_EXIT
	ICON_MODE_2D
	ICON_MODE_3D
	ICON_CUBE
	ICON_CUBE_FACE_TOP
	ICON_CUBE_FACE_LEFT
	ICON_CUBE_FACE_FRONT
	ICON_CUBE_FACE_BOTTOM
	ICON_CUBE_FACE_RIGHT
	ICON_CUBE_FACE_BACK
	ICON_CAMERA
	ICON_SPECIAL
	ICON_LINK_NET
	ICON_LINK_BOXES
	ICON_LINK_MULTI
	ICON_LINK
	ICON_LINK_BROKE
	ICON_TEXT_NOTES
	ICON_NOTEBOOK
	ICON_SUITCASE
	ICON_SUITCASE_ZIP
	ICON_MAILBOX
	ICON_MONITOR
	ICON_PRINTER
	ICON_PHOTO_CAMERA
	ICON_PHOTO_CAMERA_FLASH
	ICON_HOUSE
	ICON_HEART
	ICON_CORNER
	ICON_VERTICAL_BARS
	ICON_VERTICAL_BARS_FILL
	ICON_LIFE_BARS
	ICON_INFO
	ICON_CROSSLINE
	ICON_HELP
	ICON_FILETYPE_ALPHA
	ICON_FILETYPE_HOME
	ICON_LAYERS_VISIBLE
	ICON_LAYERS
	ICON_WINDOW
	ICON_HIDPI
	ICON_FILETYPE_BINARY
	ICON_HEX
	ICON_SHIELD
	ICON_FILE_NEW
	ICON_FOLDER_ADD
	ICON_ALARM
	ICON_CPU
	ICON_ROM
	ICON_STEP_OVER
	ICON_STEP_INTO
	ICON_STEP_OUT
	ICON_RESTART
	ICON_BREAKPOINT_ON
	ICON_BREAKPOINT_OFF
	ICON_BURGER_MENU
	ICON_CASE_SENSITIVE
	ICON_REG_EXP
	ICON_FOLDER
	ICON_FILE
	ICON_SAND_TIMER
	ICON_220
	ICON_221
	ICON_222
	ICON_223
	ICON_224
	ICON_225
	ICON_226
	ICON_227
	ICON_228
	ICON_229
	ICON_230
	ICON_231
	ICON_232
	ICON_233
	ICON_234
	ICON_235
	ICON_236
	ICON_237
	ICON_238
	ICON_239
	ICON_240
	ICON_241
	ICON_242
	ICON_243
	ICON_244
	ICON_245
	ICON_246
	ICON_247
	ICON_248
	ICON_249
	ICON_250
	ICON_251
	ICON_252
	ICON_253
	ICON_254
	ICON_255
)

// TextInputBox control, ask for text
func TextInputBox(bounds rl.Rectangle, title, message, buttons string, text *string, textMaxSize int32, secretViewActive *bool) int32 {
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

	csecretViewActive := C.bool(*secretViewActive)
	defer func() {
		*secretViewActive = bool(csecretViewActive)
	}()

	return int32(C.GuiTextInputBox(cbounds, ctitle, cmessage, cbuttons, ctext, ctextMaxSize, &csecretViewActive))
}

// ListViewEx control with extended parameters
func ListViewEx(bounds rl.Rectangle, text []string, focus, scrollIndex *int32, active int32) int32 {
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

	C.GuiListViewEx(cbounds, (**C.char)(ctext.Pointer), count, &cfocus, &cscrollIndex, &cactive)
	return int32(cactive)
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
