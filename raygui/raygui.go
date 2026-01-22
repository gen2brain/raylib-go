package raygui

/*
#define RAYGUI_IMPLEMENTATION
#include "raygui.h"
#include <stdlib.h>
*/
import "C"

import (
	"image/color"
	"strings"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type (
	ControlID     uint16
	PropertyID    uint16
	PropertyValue uint32

	IconID uint32
)

func (p PropertyID) IsExtended() bool {
	return p >= 16
}

func (v PropertyValue) AsColor() rl.Color {
	return rl.Color{uint8(v >> 24), uint8(v >> 16), uint8(v >> 8), uint8(v)}
}

func NewColorPropertyValue(color rl.Color) PropertyValue {
	return PropertyValue(uint32(color.R)<<24 + uint32(color.G)<<16 + uint32(color.B)<<8 + uint32(color.A))
}

// Gui control state
const (
	STATE_NORMAL PropertyValue = iota
	STATE_FOCUSED
	STATE_PRESSED
	STATE_DISABLED
)

// Gui control text alignment
const (
	TEXT_ALIGN_LEFT PropertyValue = iota
	TEXT_ALIGN_CENTER
	TEXT_ALIGN_RIGHT
)

// Gui control text alignment vertical
// NOTE: Text vertical position inside the text bounds
const (
	TEXT_ALIGN_TOP PropertyValue = iota
	TEXT_ALIGN_MIDDLE
	TEXT_ALIGN_BOTTOM
)

// Gui control text wrap mode
// NOTE: Useful for multiline text
const (
	TEXT_WRAP_NONE PropertyValue = iota
	TEXT_WRAP_CHAR
	TEXT_WRAP_WORD
)

const (
	SCROLLBAR_LEFT_SIDE PropertyValue = iota
	SCROLLBAR_RIGHT_SIDE
)

// Gui controls
const (
	// Default -> populates to all controls when set
	DEFAULT ControlID = iota

	// Basic controls
	LABEL // Used also for: LABELBUTTON
	BUTTON
	TOGGLE // Used also for: TOGGLEGROUP
	SLIDER // Used also for: SLIDERBAR, TOGGLESLIDER
	PROGRESSBAR
	CHECKBOX
	COMBOBOX
	DROPDOWNBOX
	TEXTBOX // Used also for: TEXTBOXMULTI
	VALUEBOX
	CONTROL11
	LISTVIEW
	COLORPICKER
	SCROLLBAR
	STATUSBAR
)

// ----------------------------------------------------------------------------------
// Module Types and Structures Definition
// ----------------------------------------------------------------------------------
// Gui control property style color element
const (
	BORDER PropertyID = iota
	BASE
	TEXT
	OTHER
)

// Gui base properties for every control
// NOTE: RAYGUI_MAX_PROPS_BASE properties (by default 16 properties)
const (
	BORDER_COLOR_NORMAL   PropertyID = iota // Control border color in STATE_NORMAL
	BASE_COLOR_NORMAL                       // Control base color in STATE_NORMAL
	TEXT_COLOR_NORMAL                       // Control text color in STATE_NORMAL
	BORDER_COLOR_FOCUSED                    // Control border color in STATE_FOCUSED
	BASE_COLOR_FOCUSED                      // Control base color in STATE_FOCUSED
	TEXT_COLOR_FOCUSED                      // Control text color in STATE_FOCUSED
	BORDER_COLOR_PRESSED                    // Control border color in STATE_PRESSED
	BASE_COLOR_PRESSED                      // Control base color in STATE_PRESSED
	TEXT_COLOR_PRESSED                      // Control text color in STATE_PRESSED
	BORDER_COLOR_DISABLED                   // Control border color in STATE_DISABLED
	BASE_COLOR_DISABLED                     // Control base color in STATE_DISABLED
	TEXT_COLOR_DISABLED                     // Control text color in STATE_DISABLED
	BORDER_WIDTH                            // Control border size, 0 for no border
	TEXT_PADDING                            // Control text padding, not considering border
	TEXT_ALIGNMENT                          // Control text horizontal alignment inside control text bound (after border and padding)
)

// Gui extended properties depend on control
// NOTE: RAYGUI_MAX_PROPS_EXTENDED properties (by default, max 8 properties)
// ----------------------------------------------------------------------------------
// DEFAULT extended properties
// NOTE: Those properties are common to all controls or global
// WARNING: We only have 8 slots for those properties by default!!! -> New global control: TEXT?
const (
	TEXT_SIZE               PropertyID = 16 + iota // Text size (glyphs max height)
	TEXT_SPACING                                   // Text spacing between glyphs
	LINE_COLOR                                     // Line control color
	BACKGROUND_COLOR                               // Background color
	TEXT_LINE_SPACING                              // Text spacing between lines
	TEXT_ALIGNMENT_VERTICAL                        // Text vertical alignment inside text bounds (after border and padding)
	TEXT_WRAP_MODE                                 // Text wrap-mode inside text bounds
)

// Toggle/ToggleGroup
const (
	GROUP_PADDING PropertyID = 16 + iota // ToggleGroup separation between toggles
)

// Slider/SliderBar
const (
	SLIDER_WIDTH   PropertyID = 16 + iota // Slider size of internal bar
	SLIDER_PADDING                        // Slider/SliderBar internal bar padding
)

// ProgressBar
const (
	PROGRESS_PADDING PropertyID = 16 + iota // ProgressBar internal padding
)

// ScrollBar
const (
	ARROWS_SIZE           PropertyID = 16 + iota // ScrollBar arrows size
	ARROWS_VISIBLE                               // ScrollBar arrows visible
	SCROLL_SLIDER_PADDING                        // ScrollBar slider internal padding
	SCROLL_SLIDER_SIZE                           // ScrollBar slider size
	SCROLL_PADDING                               // ScrollBar scroll padding from arrows
	SCROLL_SPEED                                 // ScrollBar scrolling speed
)

// CheckBox
const (
	CHECK_PADDING PropertyID = 16 + iota // CheckBox internal check padding
)

// ComboBox
const (
	COMBO_BUTTON_WIDTH   PropertyID = 16 + iota // ComboBox right button width
	COMBO_BUTTON_SPACING                        // ComboBox button separation
)

// DropdownBox
const (
	ARROW_PADDING          PropertyID = 16 + iota // DropdownBox arrow separation from border and items
	DROPDOWN_ITEMS_SPACING                        // DropdownBox items separation
	DROPDOWN_ARROW_HIDDEN                         // DropdownBox arrow hidden
	DROPDOWN_ROLL_UP                              // DropdownBox roll up flag (default rolls down)
)

// TextBox/TextBoxMulti/ValueBox/Spinner
const (
	TEXT_READONLY PropertyID = 16 + iota // TextBox in read-only mode: 0-text editable, 1-text no-editable
)

// ValueBox/Spinner
const (
	SPINNER_BUTTON_WIDTH   PropertyID = 16 // Spinner left/right buttons width
	SPINNER_BUTTON_SPACING                 // Spinner buttons separation
)

// ListView
const (
	LIST_ITEMS_HEIGHT        PropertyID = 16 + iota // ListView items height
	LIST_ITEMS_SPACING                              // ListView items separation
	SCROLLBAR_WIDTH                                 // ListView scrollbar size (usually width)
	SCROLLBAR_SIDE                                  // ListView scrollbar side (0-SCROLLBAR_LEFT_SIDE, 1-SCROLLBAR_RIGHT_SIDE)
	LIST_ITEMS_BORDER_NORMAL                        // ListView items border enabled in normal state
	LIST_ITEMS_BORDER_WIDTH                         // ListView items border width
)

// ColorPicker
const (
	COLOR_SELECTOR_SIZE      PropertyID = 16 + iota
	HUEBAR_WIDTH                        // ColorPicker right hue bar width
	HUEBAR_PADDING                      // ColorPicker right hue bar separation from panel
	HUEBAR_SELECTOR_HEIGHT              // ColorPicker right hue bar selector height
	HUEBAR_SELECTOR_OVERFLOW            // ColorPicker right hue bar selector overflow
)

// Icons enumeration
const (
	ICON_NONE                    IconID = 0
	ICON_FOLDER_FILE_OPEN        IconID = 1
	ICON_FILE_SAVE_CLASSIC       IconID = 2
	ICON_FOLDER_OPEN             IconID = 3
	ICON_FOLDER_SAVE             IconID = 4
	ICON_FILE_OPEN               IconID = 5
	ICON_FILE_SAVE               IconID = 6
	ICON_FILE_EXPORT             IconID = 7
	ICON_FILE_ADD                IconID = 8
	ICON_FILE_DELETE             IconID = 9
	ICON_FILETYPE_TEXT           IconID = 10
	ICON_FILETYPE_AUDIO          IconID = 11
	ICON_FILETYPE_IMAGE          IconID = 12
	ICON_FILETYPE_PLAY           IconID = 13
	ICON_FILETYPE_VIDEO          IconID = 14
	ICON_FILETYPE_INFO           IconID = 15
	ICON_FILE_COPY               IconID = 16
	ICON_FILE_CUT                IconID = 17
	ICON_FILE_PASTE              IconID = 18
	ICON_CURSOR_HAND             IconID = 19
	ICON_CURSOR_POINTER          IconID = 20
	ICON_CURSOR_CLASSIC          IconID = 21
	ICON_PENCIL                  IconID = 22
	ICON_PENCIL_BIG              IconID = 23
	ICON_BRUSH_CLASSIC           IconID = 24
	ICON_BRUSH_PAINTER           IconID = 25
	ICON_WATER_DROP              IconID = 26
	ICON_COLOR_PICKER            IconID = 27
	ICON_RUBBER                  IconID = 28
	ICON_COLOR_BUCKET            IconID = 29
	ICON_TEXT_T                  IconID = 30
	ICON_TEXT_A                  IconID = 31
	ICON_SCALE                   IconID = 32
	ICON_RESIZE                  IconID = 33
	ICON_FILTER_POINT            IconID = 34
	ICON_FILTER_BILINEAR         IconID = 35
	ICON_CROP                    IconID = 36
	ICON_CROP_ALPHA              IconID = 37
	ICON_SQUARE_TOGGLE           IconID = 38
	ICON_SYMMETRY                IconID = 39
	ICON_SYMMETRY_HORIZONTAL     IconID = 40
	ICON_SYMMETRY_VERTICAL       IconID = 41
	ICON_LENS                    IconID = 42
	ICON_LENS_BIG                IconID = 43
	ICON_EYE_ON                  IconID = 44
	ICON_EYE_OFF                 IconID = 45
	ICON_FILTER_TOP              IconID = 46
	ICON_FILTER                  IconID = 47
	ICON_TARGET_POINT            IconID = 48
	ICON_TARGET_SMALL            IconID = 49
	ICON_TARGET_BIG              IconID = 50
	ICON_TARGET_MOVE             IconID = 51
	ICON_CURSOR_MOVE             IconID = 52
	ICON_CURSOR_SCALE            IconID = 53
	ICON_CURSOR_SCALE_RIGHT      IconID = 54
	ICON_CURSOR_SCALE_LEFT       IconID = 55
	ICON_UNDO                    IconID = 56
	ICON_REDO                    IconID = 57
	ICON_REREDO                  IconID = 58
	ICON_MUTATE                  IconID = 59
	ICON_ROTATE                  IconID = 60
	ICON_REPEAT                  IconID = 61
	ICON_SHUFFLE                 IconID = 62
	ICON_EMPTYBOX                IconID = 63
	ICON_TARGET                  IconID = 64
	ICON_TARGET_SMALL_FILL       IconID = 65
	ICON_TARGET_BIG_FILL         IconID = 66
	ICON_TARGET_MOVE_FILL        IconID = 67
	ICON_CURSOR_MOVE_FILL        IconID = 68
	ICON_CURSOR_SCALE_FILL       IconID = 69
	ICON_CURSOR_SCALE_RIGHT_FILL IconID = 70
	ICON_CURSOR_SCALE_LEFT_FILL  IconID = 71
	ICON_UNDO_FILL               IconID = 72
	ICON_REDO_FILL               IconID = 73
	ICON_REREDO_FILL             IconID = 74
	ICON_MUTATE_FILL             IconID = 75
	ICON_ROTATE_FILL             IconID = 76
	ICON_REPEAT_FILL             IconID = 77
	ICON_SHUFFLE_FILL            IconID = 78
	ICON_EMPTYBOX_SMALL          IconID = 79
	ICON_BOX                     IconID = 80
	ICON_BOX_TOP                 IconID = 81
	ICON_BOX_TOP_RIGHT           IconID = 82
	ICON_BOX_RIGHT               IconID = 83
	ICON_BOX_BOTTOM_RIGHT        IconID = 84
	ICON_BOX_BOTTOM              IconID = 85
	ICON_BOX_BOTTOM_LEFT         IconID = 86
	ICON_BOX_LEFT                IconID = 87
	ICON_BOX_TOP_LEFT            IconID = 88
	ICON_BOX_CENTER              IconID = 89
	ICON_BOX_CIRCLE_MASK         IconID = 90
	ICON_POT                     IconID = 91
	ICON_ALPHA_MULTIPLY          IconID = 92
	ICON_ALPHA_CLEAR             IconID = 93
	ICON_DITHERING               IconID = 94
	ICON_MIPMAPS                 IconID = 95
	ICON_BOX_GRID                IconID = 96
	ICON_GRID                    IconID = 97
	ICON_BOX_CORNERS_SMALL       IconID = 98
	ICON_BOX_CORNERS_BIG         IconID = 99
	ICON_FOUR_BOXES              IconID = 100
	ICON_GRID_FILL               IconID = 101
	ICON_BOX_MULTISIZE           IconID = 102
	ICON_ZOOM_SMALL              IconID = 103
	ICON_ZOOM_MEDIUM             IconID = 104
	ICON_ZOOM_BIG                IconID = 105
	ICON_ZOOM_ALL                IconID = 106
	ICON_ZOOM_CENTER             IconID = 107
	ICON_BOX_DOTS_SMALL          IconID = 108
	ICON_BOX_DOTS_BIG            IconID = 109
	ICON_BOX_CONCENTRIC          IconID = 110
	ICON_BOX_GRID_BIG            IconID = 111
	ICON_OK_TICK                 IconID = 112
	ICON_CROSS                   IconID = 113
	ICON_ARROW_LEFT              IconID = 114
	ICON_ARROW_RIGHT             IconID = 115
	ICON_ARROW_DOWN              IconID = 116
	ICON_ARROW_UP                IconID = 117
	ICON_ARROW_LEFT_FILL         IconID = 118
	ICON_ARROW_RIGHT_FILL        IconID = 119
	ICON_ARROW_DOWN_FILL         IconID = 120
	ICON_ARROW_UP_FILL           IconID = 121
	ICON_AUDIO                   IconID = 122
	ICON_FX                      IconID = 123
	ICON_WAVE                    IconID = 124
	ICON_WAVE_SINUS              IconID = 125
	ICON_WAVE_SQUARE             IconID = 126
	ICON_WAVE_TRIANGULAR         IconID = 127
	ICON_CROSS_SMALL             IconID = 128
	ICON_PLAYER_PREVIOUS         IconID = 129
	ICON_PLAYER_PLAY_BACK        IconID = 130
	ICON_PLAYER_PLAY             IconID = 131
	ICON_PLAYER_PAUSE            IconID = 132
	ICON_PLAYER_STOP             IconID = 133
	ICON_PLAYER_NEXT             IconID = 134
	ICON_PLAYER_RECORD           IconID = 135
	ICON_MAGNET                  IconID = 136
	ICON_LOCK_CLOSE              IconID = 137
	ICON_LOCK_OPEN               IconID = 138
	ICON_CLOCK                   IconID = 139
	ICON_TOOLS                   IconID = 140
	ICON_GEAR                    IconID = 141
	ICON_GEAR_BIG                IconID = 142
	ICON_BIN                     IconID = 143
	ICON_HAND_POINTER            IconID = 144
	ICON_LASER                   IconID = 145
	ICON_COIN                    IconID = 146
	ICON_EXPLOSION               IconID = 147
	ICON_1UP                     IconID = 148
	ICON_PLAYER                  IconID = 149
	ICON_PLAYER_JUMP             IconID = 150
	ICON_KEY                     IconID = 151
	ICON_DEMON                   IconID = 152
	ICON_TEXT_POPUP              IconID = 153
	ICON_GEAR_EX                 IconID = 154
	ICON_CRACK                   IconID = 155
	ICON_CRACK_POINTS            IconID = 156
	ICON_STAR                    IconID = 157
	ICON_DOOR                    IconID = 158
	ICON_EXIT                    IconID = 159
	ICON_MODE_2D                 IconID = 160
	ICON_MODE_3D                 IconID = 161
	ICON_CUBE                    IconID = 162
	ICON_CUBE_FACE_TOP           IconID = 163
	ICON_CUBE_FACE_LEFT          IconID = 164
	ICON_CUBE_FACE_FRONT         IconID = 165
	ICON_CUBE_FACE_BOTTOM        IconID = 166
	ICON_CUBE_FACE_RIGHT         IconID = 167
	ICON_CUBE_FACE_BACK          IconID = 168
	ICON_CAMERA                  IconID = 169
	ICON_SPECIAL                 IconID = 170
	ICON_LINK_NET                IconID = 171
	ICON_LINK_BOXES              IconID = 172
	ICON_LINK_MULTI              IconID = 173
	ICON_LINK                    IconID = 174
	ICON_LINK_BROKE              IconID = 175
	ICON_TEXT_NOTES              IconID = 176
	ICON_NOTEBOOK                IconID = 177
	ICON_SUITCASE                IconID = 178
	ICON_SUITCASE_ZIP            IconID = 179
	ICON_MAILBOX                 IconID = 180
	ICON_MONITOR                 IconID = 181
	ICON_PRINTER                 IconID = 182
	ICON_PHOTO_CAMERA            IconID = 183
	ICON_PHOTO_CAMERA_FLASH      IconID = 184
	ICON_HOUSE                   IconID = 185
	ICON_HEART                   IconID = 186
	ICON_CORNER                  IconID = 187
	ICON_VERTICAL_BARS           IconID = 188
	ICON_VERTICAL_BARS_FILL      IconID = 189
	ICON_LIFE_BARS               IconID = 190
	ICON_INFO                    IconID = 191
	ICON_CROSSLINE               IconID = 192
	ICON_HELP                    IconID = 193
	ICON_FILETYPE_ALPHA          IconID = 194
	ICON_FILETYPE_HOME           IconID = 195
	ICON_LAYERS_VISIBLE          IconID = 196
	ICON_LAYERS                  IconID = 197
	ICON_WINDOW                  IconID = 198
	ICON_HIDPI                   IconID = 199
	ICON_FILETYPE_BINARY         IconID = 200
	ICON_HEX                     IconID = 201
	ICON_SHIELD                  IconID = 202
	ICON_FILE_NEW                IconID = 203
	ICON_FOLDER_ADD              IconID = 204
	ICON_ALARM                   IconID = 205
	ICON_CPU                     IconID = 206
	ICON_ROM                     IconID = 207
	ICON_STEP_OVER               IconID = 208
	ICON_STEP_INTO               IconID = 209
	ICON_STEP_OUT                IconID = 210
	ICON_RESTART                 IconID = 211
	ICON_BREAKPOINT_ON           IconID = 212
	ICON_BREAKPOINT_OFF          IconID = 213
	ICON_BURGER_MENU             IconID = 214
	ICON_CASE_SENSITIVE          IconID = 215
	ICON_REG_EXP                 IconID = 216
	ICON_FOLDER                  IconID = 217
	ICON_FILE                    IconID = 218
	ICON_SAND_TIMER              IconID = 219
	ICON_WARNING                 IconID = 220
	ICON_HELP_BOX                IconID = 221
	ICON_INFO_BOX                IconID = 222
	ICON_PRIORITY                IconID = 223
	ICON_LAYERS_ISO              IconID = 224
	ICON_LAYERS2                 IconID = 225
	ICON_MLAYERS                 IconID = 226
	ICON_MAPS                    IconID = 227
	ICON_HOT                     IconID = 228
	ICON_LABEL                   IconID = 229
	ICON_NAME_ID                 IconID = 230
	ICON_SLICING                 IconID = 231
	ICON_MANUAL_CONTROL          IconID = 232
	ICON_COLLISION               IconID = 233
	ICON_CIRCLE_ADD              IconID = 234
	ICON_CIRCLE_ADD_FILL         IconID = 235
	ICON_CIRCLE_WARNING          IconID = 236
	ICON_CIRCLE_WARNING_FILL     IconID = 237
	ICON_BOX_MORE                IconID = 238
	ICON_BOX_MORE_FILL           IconID = 239
	ICON_BOX_MINUS               IconID = 240
	ICON_BOX_MINUS_FILL          IconID = 241
	ICON_UNION                   IconID = 242
	ICON_INTERSECTION            IconID = 243
	ICON_DIFFERENCE              IconID = 244
	ICON_SPHERE                  IconID = 245
	ICON_CYLINDER                IconID = 246
	ICON_CONE                    IconID = 247
	ICON_ELLIPSOID               IconID = 248
	ICON_CAPSULE                 IconID = 249
	ICON_250                     IconID = 250
	ICON_251                     IconID = 251
	ICON_252                     IconID = 252
	ICON_253                     IconID = 253
	ICON_254                     IconID = 254
	ICON_255                     IconID = 255
)

//----------------------------------------------------------------------------------
// Gui Setup Functions Definition
//----------------------------------------------------------------------------------

// Enable gui global state
func Enable() {
	C.GuiEnable()
}

// Disable gui global state
func Disable() {
	C.GuiDisable()
}

// Lock gui global state
func Lock() {
	C.GuiLock()
}

// Unlock gui global state
func Unlock() {
	C.GuiUnlock()
}

// Check if gui is locked (global state)
func IsLocked() bool {
	return bool(C.GuiIsLocked())
}

// Set gui controls alpha global state
func SetAlpha(alpha float32) {
	calpha := C.float(alpha)
	C.GuiSetAlpha(calpha)
}

// Set gui state (global state)
func SetState(state PropertyValue) {
	cstate := C.int(state)
	C.GuiSetState(cstate)
}

// Get gui state (global state)
func GetState() PropertyValue {
	return PropertyValue(C.GuiGetState())
}

// Set custom gui font
func SetFont(font rl.Font) {
	cfont := (*C.Font)(unsafe.Pointer(&font))
	C.GuiSetFont(*cfont)
}

// Get custom gui font
func GetFont() rl.Font {
	ret := C.GuiGetFont()
	ptr := unsafe.Pointer(&ret)
	return *(*rl.Font)(ptr)
}

// Set control style property value
func SetStyle(control ControlID, property PropertyID, value PropertyValue) {
	ccontrol := C.int(control)
	cproperty := C.int(property)
	cvalue := C.int(value)
	C.GuiSetStyle(ccontrol, cproperty, cvalue)
}

// Get control style property value
func GetStyle(control ControlID, property PropertyID) PropertyValue {
	ccontrol := C.int(control)
	cproperty := C.int(property)
	return PropertyValue(C.GuiGetStyle(ccontrol, cproperty))
}

func GetColor(control ControlID, property PropertyID) rl.Color {
	color := C.GuiGetStyle(C.int(control), C.int(property))
	return rl.Color{R: uint8(color >> 24), G: uint8(color >> 16), B: uint8(color >> 8), A: uint8(color)}
}

//----------------------------------------------------------------------------------
// Gui Controls Functions Definition
//----------------------------------------------------------------------------------

// Window Box control
func WindowBox(bounds rl.Rectangle, title string) bool {
	var cbounds C.struct_Rectangle
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	var ctitle *C.char
	if len(title) > 0 {
		ctitle = C.CString(title)
		defer C.free(unsafe.Pointer(ctitle))
	}
	return C.GuiWindowBox(cbounds, ctitle) != 0
}

// Group Box control with text name
func GroupBox(bounds rl.Rectangle, text string) {
	var cbounds C.struct_Rectangle
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	C.GuiGroupBox(cbounds, ctext)
}

// Line control
func Line(bounds rl.Rectangle, text string) {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	C.GuiLine(cbounds, ctext)
}

// Panel control
func Panel(bounds rl.Rectangle, text string) {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	C.GuiPanel(cbounds, ctext)
}

// Tab Bar control, returns the current TAB closing requested, -1 otherwise
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

// Scroll Panel control
func ScrollPanel(bounds rl.Rectangle, text string, content rl.Rectangle, scroll *rl.Vector2, view *rl.Rectangle) {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
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
	defer func() {
		view.X = float32(cview.x)
		view.Y = float32(cview.y)
		view.Width = float32(cview.width)
		view.Height = float32(cview.height)
	}()

	C.GuiScrollPanel(cbounds, ctext, ccontent, &cscroll, &cview)
}

// Label control
func Label(bounds rl.Rectangle, text string) {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	C.GuiLabel(cbounds, ctext)
}

// Button control, returns true when clicked
func Button(bounds rl.Rectangle, text string) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	return C.GuiButton(cbounds, ctext) != 0
}

// LabelButton control, returns true when clicked
func LabelButton(bounds rl.Rectangle, text string) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	return C.GuiLabelButton(cbounds, ctext) != 0
}

// Toggle control, returns true when active
func Toggle(bounds rl.Rectangle, text string, active bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
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
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
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
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
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
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
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
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	cactive := C.int(active)
	C.GuiComboBox(cbounds, ctext, &cactive)
	return int32(cactive)
}

// DropdownBox control, returns true when clicked
func DropdownBox(bounds rl.Rectangle, text string, active *int32, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}

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

// TextBox control, updates input text, returns true on ENTER pressed or defocused
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
		*text = strings.Trim(string(bs), "\x00")
		// no need : C.free(unsafe.Pointer(ctext))
	}()

	ctextSize := C.int(textSize)
	ceditMode := C.bool(editMode)

	return C.GuiTextBox(cbounds, ctext, ctextSize, ceditMode) != 0
}

// Spinner control, sets value to the selected number and returns true when clicked.
func Spinner(bounds rl.Rectangle, text string, value *int32, minValue, maxValue int, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}

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

	return C.GuiSpinner(cbounds, ctext, &cvalue, cminValue, cmaxValue, ceditMode) != 0
}

// ValueBox control, updates input text with numbers
func ValueBox(bounds rl.Rectangle, text string, value *int32, minValue, maxValue int, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}

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

// Floating point Value Box control, updates input val_str with numbers
func ValueBoxFloat(bounds rl.Rectangle, text string, textValue *string, value *float32, editMode bool) bool {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}

	bs := []byte(*textValue)
	if len(bs) == 0 {
		bs = []byte{byte(0)}
	}
	if 0 < len(bs) && bs[len(bs)-1] != byte(0) { // minimalize allocation
		bs = append(bs, byte(0)) // for next input symbols
	}
	ctextValue := (*C.char)(unsafe.Pointer(&bs[0]))
	defer func() {
		*textValue = strings.Trim(string(bs), "\x00")
		// no need : C.free(unsafe.Pointer(ctext))
	}()

	if value == nil {
		value = new(float32)
	}
	cvalue := C.float(*value)
	defer func() {
		*value = float32(cvalue)
	}()

	return C.GuiValueBoxFloat(cbounds, ctext, ctextValue, &cvalue, C.bool(editMode)) != 0
}

// Slider control
func Slider(bounds rl.Rectangle, textLeft, textRight string, value, minValue, maxValue float32) float32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)

	var ctextLeft *C.char
	if len(textLeft) > 0 {
		ctextLeft = C.CString(textLeft)
		defer C.free(unsafe.Pointer(ctextLeft))
	}

	var ctextRight *C.char
	if len(textRight) > 0 {
		ctextRight = C.CString(textRight)
		defer C.free(unsafe.Pointer(ctextRight))
	}

	cvalue := C.float(value)
	cminValue := C.float(minValue)
	cmaxValue := C.float(maxValue)
	C.GuiSlider(cbounds, ctextLeft, ctextRight, &cvalue, cminValue, cmaxValue)
	return float32(cvalue)
}

// SliderBar control, returns selected value
func SliderBar(bounds rl.Rectangle, textLeft, textRight string, value, minValue, maxValue float32) float32 {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)

	var ctextLeft *C.char
	if len(textLeft) > 0 {
		ctextLeft = C.CString(textLeft)
		defer C.free(unsafe.Pointer(ctextLeft))
	}

	var ctextRight *C.char
	if len(textRight) > 0 {
		ctextRight = C.CString(textRight)
		defer C.free(unsafe.Pointer(ctextRight))
	}

	cvalue := C.float(value)
	cminValue := C.float(minValue)
	cmaxValue := C.float(maxValue)
	C.GuiSliderBar(cbounds, ctextLeft, ctextRight, &cvalue, cminValue, cmaxValue)
	return float32(cvalue)
}

// ProgressBar control, shows current progress value
func ProgressBar(bounds rl.Rectangle, textLeft, textRight string, value, minValue, maxValue float32) float32 {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)

	var ctextLeft *C.char
	if len(textLeft) > 0 {
		ctextLeft = C.CString(textLeft)
		defer C.free(unsafe.Pointer(ctextLeft))
	}

	var ctextRight *C.char
	if len(textRight) > 0 {
		ctextRight = C.CString(textRight)
		defer C.free(unsafe.Pointer(ctextRight))
	}

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
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	C.GuiStatusBar(cbounds, ctext)
}

// DummyRectangle control, intended for placeholding
func DummyRec(bounds rl.Rectangle, text string) {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	C.GuiDummyRec(cbounds, ctext)
}

// ListView control, returns selected list item index
func ListView(bounds rl.Rectangle, text string, scrollIndex *int32, active int32) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}

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

// ListView control with extended parameters
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

// ColorPanel control, Color (RGBA) variant
func ColorPanel(bounds rl.Rectangle, text string, color rl.Color) rl.Color {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
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

// ColorBarAlpha control, returns alpha value normalized [0..1]
func ColorBarAlpha(bounds rl.Rectangle, text string, alpha float32) float32 {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	calpha := C.float(alpha)
	C.GuiColorBarAlpha(cbounds, ctext, &calpha)
	return float32(calpha)
}

// ColorBarHue control, returns alpha value normalized [0..1]
func ColorBarHue(bounds rl.Rectangle, text string, value float32) float32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	cvalue := C.float(value)
	C.GuiColorBarHue(cbounds, ctext, &cvalue)
	return float32(cvalue)
}

// ColorPicker control (multiple color controls)
// NOTE: this picker converts RGB to HSV, which can cause the Hue control to jump. If you have this problem, consider using the HSV variant instead
func ColorPicker(bounds rl.Rectangle, text string, color rl.Color) rl.Color {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
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

// ColorPicker control that avoids conversion to RGB on each call (multiple color controls)
func ColorPickerHSV(bounds rl.Rectangle, text string, colorHSV *rl.Vector3) int32 {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)

	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}

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

// ColorPanel control that returns HSV color value
func ColorPanelHSV(bounds rl.Rectangle, text string, colorHSV *rl.Vector3) int32 {
	var cbounds C.struct_Rectangle
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)

	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}

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

// MessageBox control
func MessageBox(bounds rl.Rectangle, title, message, buttons string) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	var ctitle *C.char
	if len(title) > 0 {
		ctitle = C.CString(title)
		defer C.free(unsafe.Pointer(ctitle))
	}
	var cmessage *C.char
	if len(message) > 0 {
		cmessage = C.CString(message)
		defer C.free(unsafe.Pointer(cmessage))
	}
	cbuttons := C.CString(buttons)
	defer C.free(unsafe.Pointer(cbuttons))
	return int32(C.GuiMessageBox(cbounds, ctitle, cmessage, cbuttons))
}

// TextInputBox control, ask for text
func TextInputBox(bounds rl.Rectangle, title, message, buttons string, text *string, textMaxSize int32, secretViewActive *bool) int32 {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)

	var ctitle *C.char
	if len(title) > 0 {
		ctitle = C.CString(title)
		defer C.free(unsafe.Pointer(ctitle))
	}

	var cmessage *C.char
	if len(message) > 0 {
		cmessage = C.CString(message)
		defer C.free(unsafe.Pointer(cmessage))
	}

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

// Grid control, returns mouse cell position
func Grid(bounds rl.Rectangle, text string, spacing float32, subdivs int32, mouseCell *rl.Vector2) int32 {
	var cbounds C.struct_Rectangle
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)
	cbounds.x = C.float(bounds.X)
	var ctext *C.char
	if len(text) > 0 {
		ctext = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
	}
	cspacing := C.float(spacing)
	csubdivs := C.int(subdivs)
	var cmouseCell C.struct_Vector2
	cmouseCell.x = C.float(mouseCell.X)
	cmouseCell.y = C.float(mouseCell.Y)
	res := C.GuiGrid(cbounds, ctext, cspacing, csubdivs, &cmouseCell)
	mouseCell.X = float32(cmouseCell.x)
	mouseCell.Y = float32(cmouseCell.y)
	return int32(res)
}

//----------------------------------------------------------------------------------
// Tooltip management functions
// NOTE: Tooltips requires some global variables: tooltipPtr
//----------------------------------------------------------------------------------

// Enable gui tooltips (global state)
func EnableTooltip() {
	C.GuiEnableTooltip()
}

// Disable gui tooltips (global state)
func DisableTooltip() {
	C.GuiDisableTooltip()
}

// Set tooltip string
func SetTooltip(tooltip string) {
	ctooltip := C.CString(tooltip)
	defer C.free(unsafe.Pointer(ctooltip))
	C.GuiSetTooltip(ctooltip)
}

//----------------------------------------------------------------------------------
// Styles loading functions
//----------------------------------------------------------------------------------

// Load raygui style file (.rgs)
func LoadStyle(fileName string) {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	C.GuiLoadStyle(cfileName)
}

// Load style default over global style
func LoadStyleDefault() {
	C.GuiLoadStyleDefault()
}

// IconText gets text with icon id prepended (if supported)
func IconText(iconId IconID, text string) string {
	ciconId := C.int(iconId)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	return C.GoString(C.GuiIconText(ciconId, ctext))
}

// Load raygui icons file (.rgi)
func LoadIcons(fileName string, loadIconsName bool) {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	C.GuiLoadIcons(cfileName, C.bool(loadIconsName))
}

// Load icons from memory (Binary files only)
func LoadIconsFromMemory(data []byte, loadIconsName bool) {
	C.GuiLoadIconsFromMemory((*C.uchar)(unsafe.Pointer(&data[0])), C.int(len(data)), C.bool(loadIconsName))
}

// Draw icon using pixel size at specified position
func DrawIcon(iconId IconID, posX, posY, pixelSize int32, col color.RGBA) {
	C.GuiDrawIcon(C.int(iconId), C.int(posX), C.int(posY), C.int(pixelSize), *(*C.Color)(unsafe.Pointer(&col)))
}

// Set icon drawing size
func SetIconScale(scale int32) {
	C.GuiSetIconScale(C.int(scale))
}

// Get text width considering gui style and icon size (if required)
func GetTextWidth(text string) int32 {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	return int32(C.GuiGetTextWidth(ctext))
}

//----------------------------------------------------------------------------------
// Module Internal Functions Definition
//----------------------------------------------------------------------------------

// Load style from memory (Binary files only)
func LoadStyleFromMemory(data []byte) {
	C.GuiLoadStyleFromMemory((*C.uchar)(unsafe.Pointer(&data[0])), C.int(len(data)))
}

// ScrollBar control
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

// Color fade-in or fade-out, alpha value normalized [0..1]
// WARNING: It multiplies current alpha by alpha scale factor
func Fade(color rl.Color, alpha float32) rl.Color {
	ccolor := C.struct_Color{C.uchar(color.R), C.uchar(color.G), C.uchar(color.B), C.uchar(color.A)}
	calpha := C.float(alpha)
	cresult := C.GuiFade(ccolor, calpha)
	return rl.Color{R: uint8(cresult.r), G: uint8(cresult.g), B: uint8(cresult.b), A: uint8(cresult.a)}
}

//----------------------------------------------------------------------------------
// Additional Draw functions
//----------------------------------------------------------------------------------

func DrawRectangle(bounds rl.Rectangle, borderWidth int32, borderColor, fillColor rl.Color) {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)

	var cborderColor C.struct_Color
	cborderColor.r = C.uchar(borderColor.R)
	cborderColor.g = C.uchar(borderColor.G)
	cborderColor.b = C.uchar(borderColor.B)
	cborderColor.a = C.uchar(borderColor.A)

	var cfillColor C.struct_Color
	cfillColor.r = C.uchar(fillColor.R)
	cfillColor.g = C.uchar(fillColor.G)
	cfillColor.b = C.uchar(fillColor.B)
	cfillColor.a = C.uchar(fillColor.A)

	bw := C.int(borderWidth)

	C.GuiDrawRectangle(cbounds, bw, cborderColor, cfillColor)
}

// DrawText - static void GuiDrawText(const char *text, Rectangle textBounds, int alignment, Color tint);
func DrawText(text string, position rl.Rectangle, alignment int32, color rl.Color) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	var cposition C.struct_Rectangle
	cposition.x = C.float(position.X)
	cposition.y = C.float(position.Y)
	cposition.width = C.float(position.Width)
	cposition.height = C.float(position.Height)

	calignment := C.int(alignment)
	var ccolor C.struct_Color
	ccolor.r = C.uchar(color.R)
	ccolor.g = C.uchar(color.G)
	ccolor.b = C.uchar(color.B)
	ccolor.a = C.uchar(color.A)

	C.GuiDrawText(ctext, cposition, calignment, ccolor)
}

// GetTextBounds - static Rectangle GetTextBounds(int control, Rectangle bounds)
func GetTextBounds(control ControlID, bounds rl.Rectangle) rl.Rectangle {
	var cbounds C.struct_Rectangle
	cbounds.x = C.float(bounds.X)
	cbounds.y = C.float(bounds.Y)
	cbounds.width = C.float(bounds.Width)
	cbounds.height = C.float(bounds.Height)

	ccontrol := C.int(control)
	cretBounds := C.GetTextBounds(ccontrol, cbounds)
	return rl.Rectangle{
		X:      float32(cretBounds.x),
		Y:      float32(cretBounds.y),
		Width:  float32(cretBounds.width),
		Height: float32(cretBounds.height),
	}
}
