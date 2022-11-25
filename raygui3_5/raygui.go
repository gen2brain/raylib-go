package raygui3_5

/*
#cgo CFLAGS: -DRAYGUI_IMPLEMENTATION -I../raylib/
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

// GuiStyleProp - transpiled function from  C4GO/tests/raylib/raygui.h:332
//
//*
//*   raygui v3.5-dev - A simple and easy-to-use immediate-mode gui library
//*
//*   DESCRIPTION:
//*
//*   raygui is a tools-dev-focused immediate-mode-gui library based on raylib but also
//*   available as a standalone library, as long as input and drawing functions are provided.
//*
//*   Controls provided:
//*
//*   # Container/separators Controls
//*       - WindowBox     --> StatusBar, Panel
//*       - GroupBox      --> Line
//*       - Line
//*       - Panel         --> StatusBar
//*       - ScrollPanel   --> StatusBar
//*
//*   # Basic Controls
//*       - Label
//*       - Button
//*       - LabelButton   --> Label
//*       - Toggle
//*       - ToggleGroup   --> Toggle
//*       - CheckBox
//*       - ComboBox
//*       - DropdownBox
//*       - TextBox
//*       - TextBoxMulti
//*       - ValueBox      --> TextBox
//*       - Spinner       --> Button, ValueBox
//*       - Slider
//*       - SliderBar     --> Slider
//*       - ProgressBar
//*       - StatusBar
//*       - DummyRec
//*       - Grid
//*
//*   # Advance Controls
//*       - ListView
//*       - rl.ColorPicker   --> rl.ColorPanel, rl.ColorBarHue
//*       - MessageBox    --> Window, Label, Button
//*       - TextInputBox  --> Window, Label, TextBox, Button
//*
//*   It also provides a set of functions for styling the controls based on its properties (size, color).
//*
//*
//*   RAYGUI STYLE (guiStyle):
//*
//*   raygui uses a global data array for all gui style properties (allocated on data segment by default),
//*   when a new style is loaded, it is loaded over the global style... but a default gui style could always be
//*   recovered with GuiLoadStyleDefault() function, that overwrites the current style to the default one
//*
//*   The global style array size is fixed and depends on the number of controls and properties:
//*
//*       static unsigned int guiStyle[RAYGUI_MAX_CONTROLS*(RAYGUI_MAX_PROPS_BASE + RAYGUI_MAX_PROPS_EXTENDED)];
//*
//*   guiStyle size is by default: 16*(16 + 8) = 384*4 = 1536 bytes = 1.5 KB
//*
//*   Note that the first set of BASE properties (by default guiStyle[0..15]) belong to the generic style
//*   used for all controls, when any of those base values is set, it is automatically populated to all
//*   controls, so, specific control values overwriting generic style should be set after base values.
//*
//*   After the first BASE set we have the EXTENDED properties (by default guiStyle[16..23]), those
//*   properties are actually common to all controls and can not be overwritten individually (like BASE ones)
//*   Some of those properties are: TEXT_SIZE, TEXT_SPACING, LINE_COLOR, BACKGROUND_COLOR
//*
//*   Custom control properties can be defined using the EXTENDED properties for each independent control.
//*
//*   TOOL: rGuiStyler is a visual tool to customize raygui style.
//*
//*
//*   RAYGUI ICONS (guiIcons):
//*
//*   raygui could use a global array containing icons data (allocated on data segment by default),
//*   a custom icons set could be loaded over this array using GuiLoadIcons(), but loaded icons set
//*   must be same RAYGUI_ICON_SIZE and no more than RAYGUI_ICON_MAX_ICONS will be loaded
//*
//*   Every icon is codified in binary form, using 1 bit per pixel, so, every 16x16 icon
//*   requires 8 integers (16*16/32) to be stored in memory.
//*
//*   When the icon is draw, actually one quad per pixel is drawn if the bit for that pixel is set.
//*
//*   The global icons array size is fixed and depends on the number of icons and size:
//*
//*       static unsigned int guiIcons[RAYGUI_ICON_MAX_ICONS*RAYGUI_ICON_DATA_ELEMENTS];
//*
//*   guiIcons size is by default: 256*(16*16/32) = 2048*4 = 8192 bytes = 8 KB
//*
//*   TOOL: rGuiIcons is a visual tool to customize raygui icons and create new ones.
//*
//*
//*   CONFIGURATION:
//*
//*   #define RAYGUI_IMPLEMENTATION
//*       Generates the implementation of the library into the included file.
//*       If not defined, the library is in header only mode and can be included in other headers
//*       or source files without problems. But only ONE file should hold the implementation.
//*
//*   #define RAYGUI_STANDALONE
//*       Avoid raylib.h header inclusion in this file. Data types defined on raylib are defined
//*       internally in the library and input management and drawing functions must be provided by
//*       the user (check library implementation for further details).
//*
//*   #define RAYGUI_NO_ICONS
//*       Avoid including embedded ricons data (256 icons, 16x16 pixels, 1-bit per pixel, 2KB)
//*
//*   #define RAYGUI_CUSTOM_ICONS
//*       Includes custom ricons.h header defining a set of custom icons,
//*       this file can be generated using rGuiIcons tool
//*
//*
//*   VERSIONS HISTORY:
//*       3.5 (xx-xxx-2022) ADDED: Multiple new icons, useful for code editing tools
//*                         ADDED: GuiTabBar(), based on GuiToggle()
//*                         REMOVED: Unneeded icon editing functions
//*                         REDESIGNED: GuiDrawText() to divide drawing by lines
//*                         REMOVED: MeasureTextEx() dependency, logic directly implemented
//*                         REMOVED: DrawTextEx() dependency, logic directly implemented
//*                         ADDED: Helper functions to split text in separate lines
//*       3.2 (22-May-2022) RENAMED: Some enum values, for unification, avoiding prefixes
//*                         REMOVED: GuiScrollBar(), only internal
//*                         REDESIGNED: GuiPanel() to support text parameter
//*                         REDESIGNED: GuiScrollPanel() to support text parameter
//*                         REDESIGNED: GuiColorPicker() to support text parameter
//*                         REDESIGNED: GuiColorPanel() to support text parameter
//*                         REDESIGNED: GuiColorBarAlpha() to support text parameter
//*                         REDESIGNED: GuiColorBarHue() to support text parameter
//*                         REDESIGNED: GuiTextInputBox() to support password
//*       3.1 (12-Jan-2022) REVIEWED: Default style for consistency (aligned with rGuiLayout v2.5 tool)
//*                         REVIEWED: GuiLoadStyle() to support compressed font atlas image data and unload previous textures
//*                         REVIEWED: External icons usage logic
//*                         REVIEWED: GuiLine() for centered alignment when including text
//*                         RENAMED: Multiple controls properties definitions to prepend RAYGUI_
//*                         RENAMED: RICON_ references to RAYGUI_ICON_ for library consistency
//*                         Projects updated and multiple tweaks
//*       3.0 (04-Nov-2021) Integrated ricons data to avoid external file
//*                         REDESIGNED: GuiTextBoxMulti()
//*                         REMOVED: GuiImageButton*()
//*                         Multiple minor tweaks and bugs corrected
//*       2.9 (17-Mar-2021) REMOVED: Tooltip API
//*       2.8 (03-May-2020) Centralized rectangles drawing to GuiDrawRectangle()
//*       2.7 (20-Feb-2020) ADDED: Possible tooltips API
//*       2.6 (09-Sep-2019) ADDED: GuiTextInputBox()
//*                         REDESIGNED: GuiListView*(), GuiDropdownBox(), GuiSlider*(), GuiProgressBar(), GuiMessageBox()
//*                         REVIEWED: GuiTextBox(), GuiSpinner(), GuiValueBox(), GuiLoadStyle()
//*                         Replaced property INNER_PADDING by TEXT_PADDING, renamed some properties
//*                         ADDED: 8 new custom styles ready to use
//*                         Multiple minor tweaks and bugs corrected
//*       2.5 (28-May-2019) Implemented extended GuiTextBox(), GuiValueBox(), GuiSpinner()
//*       2.3 (29-Apr-2019) ADDED: rIcons auxiliar library and support for it, multiple controls reviewed
//*                         Refactor all controls drawing mechanism to use control state
//*       2.2 (05-Feb-2019) ADDED: GuiScrollBar(), GuiScrollPanel(), reviewed GuiListView(), removed Gui*Ex() controls
//*       2.1 (26-Dec-2018) REDESIGNED: GuiCheckBox(), GuiComboBox(), GuiDropdownBox(), GuiToggleGroup() > Use combined text string
//*                         REDESIGNED: Style system (breaking change)
//*       2.0 (08-Nov-2018) ADDED: Support controls guiLock and custom fonts
//*                         REVIEWED: GuiComboBox(), GuiListView()...
//*       1.9 (09-Oct-2018) REVIEWED: GuiGrid(), GuiTextBox(), GuiTextBoxMulti(), GuiValueBox()...
//*       1.8 (01-May-2018) Lot of rework and redesign to align with rGuiStyler and rGuiLayout
//*       1.5 (21-Jun-2017) Working in an improved styles system
//*       1.4 (15-Jun-2017) Rewritten all GUI functions (removed useless ones)
//*       1.3 (12-Jun-2017) Complete redesign of style system
//*       1.1 (01-Jun-2017) Complete review of the library
//*       1.0 (07-Jun-2016) Converted to header-only by Ramon Santamaria.
//*       0.9 (07-Mar-2016) Reviewed and tested by Albert Martos, Ian Eito, Sergio Martinez and Ramon Santamaria.
//*       0.8 (27-Aug-2015) Initial release. Implemented by Kevin Gato, Daniel Nicolás and Ramon Santamaria.
//*
//*
//*   CONTRIBUTORS:
//*
//*       Ramon Santamaria:   Supervision, review, redesign, update and maintenance
//*       Vlad Adrian:        Complete rewrite of GuiTextBox() to support extended features (2019)
//*       Sergio Martinez:    Review, testing (2015) and redesign of multiple controls (2018)
//*       Adria Arranz:       Testing and Implementation of additional controls (2018)
//*       Jordi Jorba:        Testing and Implementation of additional controls (2018)
//*       Albert Martos:      Review and testing of the library (2015)
//*       Ian Eito:           Review and testing of the library (2015)
//*       Kevin Gato:         Initial implementation of basic components (2014)
//*       Daniel Nicolas:     Initial implementation of basic components (2014)
//*
//*
//*   LICENSE: zlib/libpng
//*
//*   Copyright (c) 2014-2022 Ramon Santamaria (@raysan5)
//*
//*   This software is provided "as-is", without any express or implied warranty. In no event
//*   will the authors be held liable for any damages arising from the use of this software.
//*
//*   Permission is granted to anyone to use this software for any purpose, including commercial
//*   applications, and to alter it and redistribute it freely, subject to the following restrictions:
//*
//*     1. The origin of this software must not be misrepresented; you must not claim that you
//*     wrote the original software. If you use this software in a product, an acknowledgment
//*     in the product documentation would be appreciated but is not required.
//*
//*     2. Altered source versions must be plainly marked as such, and must not be misrepresented
//*     as being the original software.
//*
//*     3. This notice may not be removed or altered from any source distribution.
//*
//
// Function specifiers in case library is build/used as a shared library (Windows)
// NOTE: Microsoft specifiers to tell compiler that symbols are imported/exported from a .dll
// Function specifiers definition
//----------------------------------------------------------------------------------
// Defines and Macros
//----------------------------------------------------------------------------------
// Allow custom memory allocators
// Simple log system to avoid printf() calls if required
// NOTE: Avoiding those calls, also avoids const strings memory usage
//----------------------------------------------------------------------------------
// Types and Structures Definition
// NOTE: Some types are required for RAYGUI_STANDALONE usage
//----------------------------------------------------------------------------------
// Style property
type GuiStyleProp struct {
	controlId     uint16
	propertyId    uint16
	propertyValue uint32
}

// STATE_NORMAL - transpiled function from  C4GO/tests/raylib/raygui.h:339
// Gui control state
const (
	STATE_NORMAL   int32 = 0
	STATE_FOCUSED        = 1
	STATE_PRESSED        = 2
	STATE_DISABLED       = 3
)

// GuiState - transpiled function from  C4GO/tests/raylib/raygui.h:339
type GuiState = int32

// TEXT_ALIGN_LEFT - transpiled function from  C4GO/tests/raylib/raygui.h:347
// Gui control text alignment
const (
	TEXT_ALIGN_LEFT   int32 = 0
	TEXT_ALIGN_CENTER       = 1
	TEXT_ALIGN_RIGHT        = 2
)

// GuiTextAlignment - transpiled function from  C4GO/tests/raylib/raygui.h:347
type GuiTextAlignment = int32

// DEFAULT - transpiled function from  C4GO/tests/raylib/raygui.h:354
// Gui controls
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

// GuiControl - transpiled function from  C4GO/tests/raylib/raygui.h:354
type GuiControl = int32

// BORDER_COLOR_NORMAL - transpiled function from  C4GO/tests/raylib/raygui.h:377
// Default -> populates to all controls when set
// Basic controls
// Used also for: LABELBUTTON
// Used also for: TOGGLEGROUP
// Used also for: SLIDERBAR
// Used also for: TEXTBOXMULTI
// Uses: BUTTON, VALUEBOX
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

// GuiControlProperty - transpiled function from  C4GO/tests/raylib/raygui.h:377
type GuiControlProperty = int32

// TEXT_SIZE - transpiled function from  C4GO/tests/raylib/raygui.h:402
// Gui extended properties depend on control
// NOTE: RAYGUI_MAX_PROPS_EXTENDED properties (by default 8 properties)
//----------------------------------------------------------------------------------
// DEFAULT extended properties
// NOTE: Those properties are common to all controls or global
const (
	TEXT_SIZE        int32 = 16
	TEXT_SPACING           = 17
	LINE_COLOR             = 18
	BACKGROUND_COLOR       = 19
)

// GuiDefaultProperty - transpiled function from  C4GO/tests/raylib/raygui.h:402
type GuiDefaultProperty = int32

// GROUP_PADDING - transpiled function from  C4GO/tests/raylib/raygui.h:416
// Text size (glyphs max height)
// Text spacing between glyphs
// Line control color
// Background color
// Label
//typedef enum { } GuiLabelProperty;
// Button/Spinner
//typedef enum { } GuiButtonProperty;
// Toggle/ToggleGroup
const (
	GROUP_PADDING int32 = 16
)

// GuiToggleProperty - transpiled function from  C4GO/tests/raylib/raygui.h:416
type GuiToggleProperty = int32

// SLIDER_WIDTH - transpiled function from  C4GO/tests/raylib/raygui.h:421
// ToggleGroup separation between toggles
// Slider/SliderBar
const (
	SLIDER_WIDTH   int32 = 16
	SLIDER_PADDING       = 17
)

// GuiSliderProperty - transpiled function from  C4GO/tests/raylib/raygui.h:421
type GuiSliderProperty = int32

// PROGRESS_PADDING - transpiled function from  C4GO/tests/raylib/raygui.h:427
// Slider size of internal bar
// Slider/SliderBar internal bar padding
// ProgressBar
const (
	PROGRESS_PADDING int32 = 16
)

// GuiProgressBarProperty - transpiled function from  C4GO/tests/raylib/raygui.h:427
type GuiProgressBarProperty = int32

// ARROWS_SIZE - transpiled function from  C4GO/tests/raylib/raygui.h:432
// ProgressBar internal padding
// ScrollBar
const (
	ARROWS_SIZE           int32 = 16
	ARROWS_VISIBLE              = 17
	SCROLL_SLIDER_PADDING       = 18
	SCROLL_SLIDER_SIZE          = 19
	SCROLL_PADDING              = 20
	SCROLL_SPEED                = 21
)

// GuiScrollBarProperty - transpiled function from  C4GO/tests/raylib/raygui.h:432
type GuiScrollBarProperty = int32

// CHECK_PADDING - transpiled function from  C4GO/tests/raylib/raygui.h:442
// (SLIDERBAR, SLIDER_PADDING)
// CheckBox
const (
	CHECK_PADDING int32 = 16
)

// GuiCheckBoxProperty - transpiled function from  C4GO/tests/raylib/raygui.h:442
type GuiCheckBoxProperty = int32

// COMBO_BUTTON_WIDTH - transpiled function from  C4GO/tests/raylib/raygui.h:447
// CheckBox internal check padding
// ComboBox
const (
	COMBO_BUTTON_WIDTH   int32 = 16
	COMBO_BUTTON_SPACING       = 17
)

// GuiComboBoxProperty - transpiled function from  C4GO/tests/raylib/raygui.h:447
type GuiComboBoxProperty = int32

// ARROW_PADDING - transpiled function from  C4GO/tests/raylib/raygui.h:453
// ComboBox right button width
// ComboBox button separation
// DropdownBox
const (
	ARROW_PADDING          int32 = 16
	DROPDOWN_ITEMS_SPACING       = 17
)

// GuiDropdownBoxProperty - transpiled function from  C4GO/tests/raylib/raygui.h:453
type GuiDropdownBoxProperty = int32

// TEXT_INNER_PADDING - transpiled function from  C4GO/tests/raylib/raygui.h:459
// DropdownBox arrow separation from border and items
// DropdownBox items separation
// TextBox/TextBoxMulti/ValueBox/Spinner
const (
	TEXT_INNER_PADDING int32 = 16
	TEXT_LINES_SPACING       = 17
)

// GuiTextBoxProperty - transpiled function from  C4GO/tests/raylib/raygui.h:459
type GuiTextBoxProperty = int32

// SPIN_BUTTON_WIDTH - transpiled function from  C4GO/tests/raylib/raygui.h:465
// TextBox/TextBoxMulti/ValueBox/Spinner inner text padding
// TextBoxMulti lines separation
// Spinner
const (
	SPIN_BUTTON_WIDTH   int32 = 16
	SPIN_BUTTON_SPACING       = 17
)

// GuiSpinnerProperty - transpiled function from  C4GO/tests/raylib/raygui.h:465
type GuiSpinnerProperty = int32

// LIST_ITEMS_HEIGHT - transpiled function from  C4GO/tests/raylib/raygui.h:471
// Spinner left/right buttons width
// Spinner buttons separation
// ListView
const (
	LIST_ITEMS_HEIGHT  int32 = 16
	LIST_ITEMS_SPACING       = 17
	SCROLLBAR_WIDTH          = 18
	SCROLLBAR_SIDE           = 19
)

// GuiListViewProperty - transpiled function from  C4GO/tests/raylib/raygui.h:471
type GuiListViewProperty = int32

// COLOR_SELECTOR_SIZE - transpiled function from  C4GO/tests/raylib/raygui.h:479
// ListView items height
// ListView items separation
// ListView scrollbar size (usually width)
// ListView scrollbar side (0-left, 1-right)
// rl.ColorPicker
const (
	COLOR_SELECTOR_SIZE      int32 = 16
	HUEBAR_WIDTH                   = 17
	HUEBAR_PADDING                 = 18
	HUEBAR_SELECTOR_HEIGHT         = 19
	HUEBAR_SELECTOR_OVERFLOW       = 20
)

// GuiColorPickerProperty - transpiled function from  C4GO/tests/raylib/raygui.h:479
type GuiColorPickerProperty = int32

// GuiEnable - transpiled function from  C4GO/tests/raylib/raygui.h:504
// rl.ColorPicker right hue bar width
// rl.ColorPicker right hue bar separation from panel
// rl.ColorPicker right hue bar selector height
// rl.ColorPicker right hue bar selector overflow
//----------------------------------------------------------------------------------
// Global Variables Definition
//----------------------------------------------------------------------------------
// ...
//----------------------------------------------------------------------------------
// Module Functions Declaration
//----------------------------------------------------------------------------------
// Global gui state control functions
// Enable gui controls (global state)
func Enable() {
	C.GuiEnable()
}

// GuiDisable - transpiled function from  C4GO/tests/raylib/raygui.h:505
// Disable gui controls (global state)
func Disable() {
	C.GuiDisable()
}

// GuiLock - transpiled function from  C4GO/tests/raylib/raygui.h:506
// Lock gui controls (global state)
func Lock() {
	C.GuiLock()
}

// GuiUnlock - transpiled function from  C4GO/tests/raylib/raygui.h:507
// Unlock gui controls (global state)
func Unlock() {
	C.GuiUnlock()
}

// GuiIsLocked - transpiled function from  C4GO/tests/raylib/raygui.h:508
// Check if gui is locked (global state)
func IsLocked() bool {
	return bool(C.GuiIsLocked())
}

// GuiFade - transpiled function from  C4GO/tests/raylib/raygui.h:509
// Set gui controls alpha (global state), alpha goes from 0.0f to 1.0f
func Fade(alpha float32) {
	calpha := C.float(alpha)
	C.GuiFade(calpha)
}

// GuiSetState - transpiled function from  C4GO/tests/raylib/raygui.h:510
// Set gui state (global state)
func SetState(state int32) {
	cstate := C.int(state)
	C.GuiSetState(cstate)
}

// GuiGetState - transpiled function from  C4GO/tests/raylib/raygui.h:511
// Get gui state (global state)
func GetState() int32 {
	return int32(C.GuiGetState())
}

// GuiSetStyle - transpiled function from  C4GO/tests/raylib/raygui.h:518
func SetStyle(control int32, property int32, value int32) {
	ccontrol := C.int(control)
	cproperty := C.int(property)
	cvalue := C.int(value)
	C.GuiSetStyle(ccontrol, cproperty, cvalue)
}

// GuiGetStyle - transpiled function from  C4GO/tests/raylib/raygui.h:519
// Get one style property
func GetStyle(control int32, property int32) int32 {
	ccontrol := C.int(control)
	cproperty := C.int(property)
	return int32(C.GuiGetStyle(ccontrol, cproperty))
}

// GuiWindowBox - transpiled function from  C4GO/tests/raylib/raygui.h:522
// Container/separator controls, useful for controls organization
// Window Box control, shows a window that can be closed
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

// GuiGroupBox - transpiled function from  C4GO/tests/raylib/raygui.h:523
// Group Box control with text name
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

// GuiLine - transpiled function from  C4GO/tests/raylib/raygui.h:524
// Line separator control, could contain text
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

// GuiPanel - transpiled function from  C4GO/tests/raylib/raygui.h:525
// Panel control, useful to group controls
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

// Scroll Panel control
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

// Scroll bar control (used by GuiScrollPanel())
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

// GuiButton - transpiled function from  C4GO/tests/raylib/raygui.h:531
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

// GuiLabelButton - transpiled function from  C4GO/tests/raylib/raygui.h:532
// Label button control, show true when clicked
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

// Toggle Button control, returns true when active
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

// GuiToggleGroup - transpiled function from  C4GO/tests/raylib/raygui.h:534
// Toggle Group control, returns active toggle index
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

// Check Box control, returns true when active
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

// GuiComboBox - transpiled function from  C4GO/tests/raylib/raygui.h:536
// Combo Box control, returns selected item index
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

// GuiSlider - transpiled function from  C4GO/tests/raylib/raygui.h:542
// Spinner control, returns selected value
// Value Box control, updates input text with numbers
// Text Box control, updates input text
// Text Box control with multiple lines
// Slider control, returns selected value
// Spinner control, returns selected value
func Spinner(bounds rl.Rectangle, text string, value *int32, minValue, maxValue int, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	var cvalue C.int
	cvalue = C.int(*value)
	defer func() {
		*value = int32(cvalue)
	}()

	cminValue := C.int(minValue)
	cmaxValue := C.int(maxValue)
	ceditMode := C.bool(editMode)

	return bool(C.GuiSpinner(cbounds, ctext, &cvalue, cminValue, cmaxValue, ceditMode))
}

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

// GuiSliderBar - transpiled function from  C4GO/tests/raylib/raygui.h:543
// Slider Bar control, returns selected value
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

// GuiProgressBar - transpiled function from  C4GO/tests/raylib/raygui.h:544
// Progress Bar control, shows current progress value
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

// GuiStatusBar - transpiled function from  C4GO/tests/raylib/raygui.h:545
// Status Bar control, shows info text
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

// GuiDummyRec - transpiled function from  C4GO/tests/raylib/raygui.h:546
// Dummy control for placeholders
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

// GuiGrid - transpiled function from  C4GO/tests/raylib/raygui.h:547
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

// GuiMessageBox - transpiled function from  C4GO/tests/raylib/raygui.h:552
// Advance controls set

// List View control, returns selected list item index
// List View control
func ListView(bounds rl.Rectangle, text string, scrollIndex *int32, active int32) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	var cscrollIndex C.int
	cscrollIndex = C.int(*scrollIndex)
	defer func() {
		*scrollIndex = int32(cscrollIndex)
	}()

	cactive := C.int(active)

	return int32(C.GuiListView(cbounds, ctext, &cscrollIndex, cactive))
}

// Message Box control, displays a message
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

// GuiColorPicker - transpiled function from  C4GO/tests/raylib/raygui.h:554
// rl.Color Picker control (multiple color controls)
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

// GuiColorPanel - transpiled function from  C4GO/tests/raylib/raygui.h:555
// Color Panel control
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

// GuiColorBarAlpha - transpiled function from  C4GO/tests/raylib/raygui.h:556
// Color Bar Alpha control
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

// GuiColorBarHue - transpiled function from  C4GO/tests/raylib/raygui.h:557
// Color Bar Hue control
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

// Dropdown Box control
// NOTE: Returns mouse click
func DropdownBox(bounds rl.Rectangle, text string, active *int32, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	var cactive C.int
	cactive = C.int(*active)
	defer func() {
		*active = int32(cactive)
	}()

	ceditMode := C.bool(editMode)

	return bool(C.GuiDropdownBox(cbounds, ctext, &cactive, ceditMode))
}

// Value Box control, updates input text with numbers
// NOTE: Requires static variables: frameCounter
func ValueBox(bounds rl.Rectangle, text string, value *int32, minValue, maxValue int, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	var cvalue C.int
	cvalue = C.int(*value)
	defer func() {
		*value = int32(cvalue)
	}()

	cminValue := C.int(minValue)
	cmaxValue := C.int(maxValue)
	ceditMode := C.bool(editMode)

	return bool(C.GuiValueBox(cbounds, ctext, &cvalue, cminValue, cmaxValue, ceditMode))
}

// Text Box control, updates input text
// NOTE 2: Returns if KEY_ENTER pressed (useful for data validation)
func TextBox(bounds rl.Rectangle, text *string, textSize int, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)

	bs := []byte(*text)
	if 0 < len(bs) && bs[len(bs)-1] != byte(0) { // minimalize allocation
		bs = append(bs, byte(0)) // for next input symbols
	}
	ctext := (*C.char)(unsafe.Pointer(&bs[0]))
	defer func() {
		*text = strings.TrimSpace(string(bs))
		// no need : C.free(unsafe.Pointer(ctext))
	}()

	ctextSize := C.int(textSize)
	ceditMode := C.bool(editMode)

	return bool(C.GuiTextBox(cbounds, ctext, ctextSize, ceditMode))
}

// GuiLoadStyle - transpiled function from  C4GO/tests/raylib/raygui.h:560
// Styles loading functions
// Load style file over global style variable (.rgs)
func LoadStyle(fileName string) {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	C.GuiLoadStyle(cfileName)
}

// TODO
// GuiLoadStyleDefault - transpiled function from  C4GO/tests/raylib/raygui.h:561
// Load style default over global style
func LoadStyleDefault() {
	C.GuiLoadStyleDefault()
}

// GuiIconText - transpiled function from  C4GO/tests/raylib/raygui.h:564
// Icons functionality
// Get text with icon id prepended (if supported)
func IconText(iconId int32, text string) string {
	ciconId := C.int(iconId)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	return C.GoString(C.GuiIconText(ciconId, ctext))
}

// TODO
// GuiGetIcons - transpiled function from  C4GO/tests/raylib/raygui.h:567
// Get raygui icons data pointer
// func GetIcons() []uint32 {
// 	return C.GuiGetIcons()
// }

// Text Input Box control, ask for text
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
	if 0 < len(bs) && bs[len(bs)-1] != byte(0) { // minimalize allocation
		bs = append(bs, byte(0)) // for next input symbols
	}
	ctext := (*C.char)(unsafe.Pointer(&bs[0]))
	defer func() {
		*text = strings.TrimSpace(string(bs))
		// no need : C.free(unsafe.Pointer(ctext))
	}()

	ctextMaxSize := C.int(textMaxSize)

	var csecretViewActive C.int
	csecretViewActive = C.int(*secretViewActive)
	defer func() {
		*secretViewActive = int32(csecretViewActive)
	}()

	return int32(C.TextInputBox(cbounds, ctitle, cmessage, cbuttons, &ctext, textMaxSize, &csecretViewActive))
}

// List View with extended parameters
// Warning (*ast.FunctionDecl): {prefix: n:GuiListViewEx,t1:int (Rectangle, const char **, int, int *, int *, int),t2:}.  C4GO/tests/raylib/raygui.h:551 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `GuiListViewEx`. cannot parse C type: `const char **`

// Warning (*ast.FunctionDecl): {prefix: n:GuiTextBoxMulti,t1:_Bool (Rectangle, char *, int, _Bool),t2:}.  C4GO/tests/raylib/raygui.h:541 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `GuiTextBoxMulti`. cannot parse C type: `_Bool`

// Check Box control, returns true when active
// Warning (*ast.FunctionDecl): {prefix: n:GuiCheckBox,t1:_Bool (Rectangle, const char *, _Bool),t2:}.  C4GO/tests/raylib/raygui.h:535 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `GuiCheckBox`. cannot parse C type: `_Bool`

// Tab Bar control, returns TAB to be closed or -1
// Warning (*ast.FunctionDecl): {prefix: n:GuiTabBar,t1:int (Rectangle, const char **, int, int *),t2:}.  C4GO/tests/raylib/raygui.h:526 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `GuiTabBar`. cannot parse C type: `const char **`

// Warning (*ast.FunctionDecl): {prefix: n:SetAudioStreamCallback,t1:void (AudioStream, AudioCallback),t2:}.  C4GO/tests/raylib/raylib.h:1567 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `SetAudioStreamCallback`. field type is pointer: `rAudioBuffer *`

// Warning (*ast.FunctionDecl): {prefix: n:AttachAudioStreamProcessor,t1:void (AudioStream, AudioCallback),t2:}.  C4GO/tests/raylib/raylib.h:1569 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `AttachAudioStreamProcessor`. field type is pointer: `rAudioBuffer *`

// Warning (*ast.FunctionDecl): {prefix: n:DetachAudioStreamProcessor,t1:void (AudioStream, AudioCallback),t2:}.  C4GO/tests/raylib/raylib.h:1570 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `DetachAudioStreamProcessor`. field type is pointer: `rAudioBuffer *`

// Warning (*ast.FunctionDecl): {prefix: n:GuiSetFont,t1:void (Font),t2:}.  C4GO/tests/raylib/raygui.h:514 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `GuiSetFont`. field type is pointer: `Rectangle *`

// Warning (*ast.FunctionDecl): {prefix: n:GuiGetFont,t1:Font (void),t2:}.  C4GO/tests/raylib/raygui.h:515 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `GuiGetFont`. field type is pointer: `Rectangle *`
