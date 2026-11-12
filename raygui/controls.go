package raygui

import (
	"fmt"
	"math"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	groupBoxLineThick      = 1
	lineMarginText         = 12
	lineTextPadding        = 4
	panelBorderWidth       = 1
	minScrollBarWidth      = 40
	minScrollBarHeight     = 40
	minMouseWheelSpeed     = 20
	toggleGroupMaxItemText = 256
)

func lineColor(state PropertyValue) rl.Color {
	if state == STATE_DISABLED {
		return getColor(getStyle(DEFAULT, BORDER_COLOR_DISABLED))
	}
	return getColor(getStyle(DEFAULT, LINE_COLOR))
}

// WindowBox - Window Box control, shows a window that can be closed
func WindowBox(bounds rl.Rectangle, title string) int32 {
	result := RESULT_NONE

	statusBarHeight := int32(windowBoxStatusBarHeight)
	statusBorderWidth := getStyle(STATUSBAR, BORDER_WIDTH)

	statusBar := rl.Rectangle{X: bounds.X, Y: bounds.Y, Width: bounds.Width, Height: float32(statusBarHeight)}
	if bounds.Height < float32(statusBarHeight)*2.0 {
		bounds.Height = float32(statusBarHeight) * 2.0
	}

	vPadding := float32(statusBarHeight)/2.0 - windowBoxCloseButtonHeight/2.0
	windowPanel := rl.Rectangle{X: bounds.X, Y: bounds.Y + float32(statusBarHeight) - float32(statusBorderWidth), Width: bounds.Width, Height: bounds.Height - float32(statusBarHeight) + float32(statusBorderWidth)}
	closeButtonRec := rl.Rectangle{X: statusBar.X + statusBar.Width - float32(statusBorderWidth) - windowBoxCloseButtonHeight - vPadding,
		Y: statusBar.Y + vPadding, Width: windowBoxCloseButtonHeight, Height: windowBoxCloseButtonHeight}

	Panel(windowPanel, "")

	tempTextAlignment := getStyle(STATUSBAR, TEXT_ALIGNMENT)
	setStyle(STATUSBAR, TEXT_ALIGNMENT, int32(TEXT_ALIGN_LEFT))
	StatusBar(statusBar, title)
	setStyle(STATUSBAR, TEXT_ALIGNMENT, tempTextAlignment)

	tempBorderWidth := getStyle(BUTTON, BORDER_WIDTH)
	tempTextAlignment = getStyle(BUTTON, TEXT_ALIGNMENT)
	setStyle(BUTTON, BORDER_WIDTH, 1)
	setStyle(BUTTON, TEXT_ALIGNMENT, int32(TEXT_ALIGN_CENTER))
	result = Button(closeButtonRec, IconText(ICON_CROSS_SMALL, ""))
	setStyle(BUTTON, BORDER_WIDTH, tempBorderWidth)
	setStyle(BUTTON, TEXT_ALIGNMENT, tempTextAlignment)

	return result
}

// GroupBox - Group Box control with text name
func GroupBox(bounds rl.Rectangle, text string) int32 {
	result := RESULT_NONE
	state := guiState

	guiDrawRectangle(rl.Rectangle{X: bounds.X, Y: bounds.Y, Width: groupBoxLineThick, Height: bounds.Height}, 0, blank, lineColor(state))
	guiDrawRectangle(rl.Rectangle{X: bounds.X, Y: bounds.Y + bounds.Height - 1, Width: bounds.Width, Height: groupBoxLineThick}, 0, blank, lineColor(state))
	guiDrawRectangle(rl.Rectangle{X: bounds.X + bounds.Width - 1, Y: bounds.Y, Width: groupBoxLineThick, Height: bounds.Height}, 0, blank, lineColor(state))

	Line(rl.Rectangle{X: bounds.X, Y: bounds.Y - float32(getStyle(DEFAULT, TEXT_SIZE)/2), Width: bounds.Width, Height: float32(getStyle(DEFAULT, TEXT_SIZE))}, text)

	return result
}

// Line - Line separator control, could contain text
func Line(bounds rl.Rectangle, text string) int32 {
	result := RESULT_NONE
	state := guiState

	color := lineColor(state)

	if text == "" {
		guiDrawRectangle(rl.Rectangle{X: bounds.X, Y: bounds.Y + bounds.Height/2, Width: bounds.Width, Height: 1}, 0, blank, color)
	} else {
		var textBounds rl.Rectangle
		textBounds.Width = float32(GetTextWidth(text)) + 2
		textBounds.Height = bounds.Height
		textBounds.X = bounds.X + lineMarginText
		textBounds.Y = bounds.Y

		guiDrawRectangle(rl.Rectangle{X: bounds.X, Y: bounds.Y + bounds.Height/2, Width: lineMarginText - lineTextPadding, Height: 1}, 0, blank, color)
		guiDrawText(text, textBounds, int32(TEXT_ALIGN_LEFT), color)
		guiDrawRectangle(rl.Rectangle{X: bounds.X + 12 + textBounds.Width + 4, Y: bounds.Y + bounds.Height/2, Width: bounds.Width - textBounds.Width - lineMarginText - lineTextPadding, Height: 1}, 0, blank, color)
	}

	return result
}

// Panel - Panel control, useful to group controls
func Panel(bounds rl.Rectangle, text string) int32 {
	result := RESULT_NONE
	state := guiState

	statusBar := rl.Rectangle{X: bounds.X, Y: bounds.Y, Width: bounds.Width, Height: windowBoxStatusBarHeight}
	if text != "" && bounds.Height < windowBoxStatusBarHeight*2.0 {
		bounds.Height = windowBoxStatusBarHeight * 2.0
	}

	if text != "" {
		bounds.Y += windowBoxStatusBarHeight - 1
		bounds.Height -= windowBoxStatusBarHeight - 1
	}

	if text != "" {
		result = StatusBar(statusBar, text)
	}

	baseColor := getColor(getStyle(DEFAULT, BACKGROUND_COLOR))
	if state == STATE_DISABLED {
		baseColor = getColor(getStyle(DEFAULT, BASE_COLOR_DISABLED))
	}
	guiDrawRectangle(bounds, panelBorderWidth, lineColor(state), baseColor)

	return result
}

// ScrollPanel - Scroll Panel control
func ScrollPanel(bounds rl.Rectangle, text string, content rl.Rectangle, scroll *rl.Vector2, view *rl.Rectangle) int32 {
	result := RESULT_NONE
	state := guiState

	var temp rl.Rectangle
	if view == nil {
		view = &temp
	}

	var scrollPos rl.Vector2
	if scroll != nil {
		scrollPos = *scroll
	}

	statusBar := rl.Rectangle{X: bounds.X, Y: bounds.Y, Width: bounds.Width, Height: windowBoxStatusBarHeight}
	if bounds.Height < windowBoxStatusBarHeight*2.0 {
		bounds.Height = windowBoxStatusBarHeight * 2.0
	}

	if text != "" {
		bounds.Y += windowBoxStatusBarHeight - 1
		bounds.Height -= windowBoxStatusBarHeight - 1
	}

	bw := getStyle(DEFAULT, BORDER_WIDTH)
	sbw := getStyle(LISTVIEW, SCROLLBAR_WIDTH)
	leftSide := GetStyle(LISTVIEW, SCROLLBAR_SIDE) == SCROLLBAR_LEFT_SIDE

	hasHorizontalScrollBar := content.Width > bounds.Width-float32(2*bw)
	hasVerticalScrollBar := content.Height > bounds.Height-float32(2*bw)

	if !hasHorizontalScrollBar {
		hasHorizontalScrollBar = hasVerticalScrollBar && content.Width > bounds.Width-float32(2*bw)-float32(sbw)
	}
	if !hasVerticalScrollBar {
		hasVerticalScrollBar = hasHorizontalScrollBar && content.Height > bounds.Height-float32(2*bw)-float32(sbw)
	}

	var horizontalScrollBarWidth, verticalScrollBarWidth int32
	if hasHorizontalScrollBar {
		horizontalScrollBarWidth = sbw
	}
	if hasVerticalScrollBar {
		verticalScrollBarWidth = sbw
	}

	hx := bounds.X
	if leftSide {
		hx = bounds.X + float32(verticalScrollBarWidth)
	}
	horizontalScrollBar := rl.Rectangle{
		X:      hx + float32(bw),
		Y:      bounds.Y + bounds.Height - float32(horizontalScrollBarWidth) - float32(bw),
		Width:  bounds.Width - float32(verticalScrollBarWidth) - float32(2*bw),
		Height: float32(horizontalScrollBarWidth),
	}

	vx := bounds.X + bounds.Width - float32(verticalScrollBarWidth) - float32(bw)
	if leftSide {
		vx = bounds.X + float32(bw)
	}
	verticalScrollBar := rl.Rectangle{
		X:      vx,
		Y:      bounds.Y + float32(bw),
		Width:  float32(verticalScrollBarWidth),
		Height: bounds.Height - float32(horizontalScrollBarWidth) - float32(2*bw),
	}

	if horizontalScrollBar.Width < minScrollBarWidth {
		horizontalScrollBar.Width = minScrollBarWidth
	}
	if verticalScrollBar.Height < minScrollBarHeight {
		verticalScrollBar.Height = minScrollBarHeight
	}

	if leftSide {
		*view = rl.Rectangle{X: bounds.X + float32(verticalScrollBarWidth) + float32(bw), Y: bounds.Y + float32(bw), Width: bounds.Width - float32(2*bw) - float32(verticalScrollBarWidth), Height: bounds.Height - float32(2*bw) - float32(horizontalScrollBarWidth)}
	} else {
		*view = rl.Rectangle{X: bounds.X + float32(bw), Y: bounds.Y + float32(bw), Width: bounds.Width - float32(2*bw) - float32(verticalScrollBarWidth), Height: bounds.Height - float32(2*bw) - float32(horizontalScrollBarWidth)}
	}

	if view.Width > content.Width {
		view.Width = content.Width
	}
	if view.Height > content.Height {
		view.Height = content.Height
	}

	var leftOffset float32
	if leftSide {
		leftOffset = float32(-verticalScrollBarWidth)
	}
	horizontalMin := leftOffset - float32(bw)
	horizontalMax := float32(-bw)
	if hasHorizontalScrollBar {
		var leftWidth float32
		if leftSide {
			leftWidth = float32(verticalScrollBarWidth)
		}
		horizontalMax = content.Width - bounds.Width + float32(verticalScrollBarWidth) + float32(bw) - leftWidth
	}
	verticalMin := -float32(bw)
	verticalMax := float32(-bw)
	if hasVerticalScrollBar {
		verticalMax = content.Height - bounds.Height + float32(horizontalScrollBarWidth) + float32(bw)
	}

	if state != STATE_DISABLED && !guiLocked {
		mousePoint := guiPointerPosition()

		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if guiButtonDown() {
				state = STATE_PRESSED
			} else {
				state = STATE_FOCUSED
			}

			scrollDelta := guiScrollDelta()

			scrollSpeed := rl.Vector2{X: content.Width / bounds.Width, Y: content.Height / bounds.Height}
			if scrollSpeed.X < minMouseWheelSpeed {
				scrollSpeed.X = minMouseWheelSpeed
			}
			if scrollSpeed.Y < minMouseWheelSpeed {
				scrollSpeed.Y = minMouseWheelSpeed
			}

			if hasHorizontalScrollBar && (guiKeyDown(rl.KeyLeftControl) || guiKeyDown(rl.KeyLeftShift)) {
				scrollPos.X += scrollDelta * scrollSpeed.X
			} else {
				scrollPos.Y += scrollDelta * scrollSpeed.Y
			}
		}
	}

	if scrollPos.X > -horizontalMin {
		scrollPos.X = -horizontalMin
	}
	if scrollPos.X < -horizontalMax {
		scrollPos.X = -horizontalMax
	}
	if scrollPos.Y > -verticalMin {
		scrollPos.Y = -verticalMin
	}
	if scrollPos.Y < -verticalMax {
		scrollPos.Y = -verticalMax
	}

	if text != "" {
		result = StatusBar(statusBar, text)
	}

	guiDrawRectangle(bounds, 0, blank, getColor(getStyle(DEFAULT, BACKGROUND_COLOR)))

	slider := getStyle(SCROLLBAR, SCROLL_SLIDER_SIZE)

	if hasHorizontalScrollBar {
		setStyle(SCROLLBAR, SCROLL_SLIDER_SIZE, int32(((bounds.Width-float32(2*bw)-float32(verticalScrollBarWidth))/float32(int32(content.Width)))*float32(int32(bounds.Width)-2*bw-verticalScrollBarWidth)))
		scrollPos.X = float32(-guiScrollBar(horizontalScrollBar, int32(-scrollPos.X), int32(horizontalMin), int32(horizontalMax)))
	} else {
		scrollPos.X = 0.0
	}

	if hasVerticalScrollBar {
		setStyle(SCROLLBAR, SCROLL_SLIDER_SIZE, int32(((bounds.Height-float32(2*bw)-float32(horizontalScrollBarWidth))/float32(int32(content.Height)))*float32(int32(bounds.Height)-2*bw-horizontalScrollBarWidth)))
		scrollPos.Y = float32(-guiScrollBar(verticalScrollBar, int32(-scrollPos.Y), int32(verticalMin), int32(verticalMax)))
	} else {
		scrollPos.Y = 0.0
	}

	if hasHorizontalScrollBar && hasVerticalScrollBar {
		cx := horizontalScrollBar.X + horizontalScrollBar.Width + 2
		if leftSide {
			cx = bounds.X + float32(bw) + 2
		}
		corner := rl.Rectangle{X: cx, Y: verticalScrollBar.Y + verticalScrollBar.Height + 2, Width: float32(horizontalScrollBarWidth) - 4, Height: float32(verticalScrollBarWidth) - 4}
		guiDrawRectangle(corner, 0, blank, getColor(getStyle(LISTVIEW, textElement+PropertyID(state*3))))
	}

	guiDrawRectangle(bounds, getStyle(LISTVIEW, BORDER_WIDTH), getColor(getStyle(LISTVIEW, borderElement+PropertyID(state*3))), blank)

	setStyle(SCROLLBAR, SCROLL_SLIDER_SIZE, slider)

	if scroll != nil {
		*scroll = scrollPos
	}

	return result
}

// Label - Label control
func Label(bounds rl.Rectangle, text string) int32 {
	result := RESULT_NONE
	state := guiState

	guiDrawText(text, getTextBounds(LABEL, bounds), getStyle(LABEL, TEXT_ALIGNMENT), getColor(getStyle(LABEL, textElement+PropertyID(state*3))))

	return result
}

// Button - Button control, returns true when clicked
func Button(bounds rl.Rectangle, text string) int32 {
	result := RESULT_NONE
	state := guiState

	if state != STATE_DISABLED && !guiLocked && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()

		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if guiButtonDown() {
				state = STATE_PRESSED
			} else {
				state = STATE_FOCUSED
			}

			if guiButtonReleased() {
				result = RESULT_PRESSED
			}
		}
	}

	guiDrawRectangle(bounds, getStyle(BUTTON, BORDER_WIDTH), getColor(getStyle(BUTTON, borderElement+PropertyID(state*3))), getColor(getStyle(BUTTON, baseElement+PropertyID(state*3))))
	guiDrawText(text, getTextBounds(BUTTON, bounds), getStyle(BUTTON, TEXT_ALIGNMENT), getColor(getStyle(BUTTON, textElement+PropertyID(state*3))))

	if state == STATE_FOCUSED {
		drawTooltip(bounds)
	}

	return result
}

// LabelButton - Label button control, returns true when clicked
func LabelButton(bounds rl.Rectangle, text string) int32 {
	result := RESULT_NONE
	state := guiState

	textWidth := float32(GetTextWidth(text))
	if bounds.Width-float32(2*getStyle(LABEL, BORDER_WIDTH))-float32(2*getStyle(LABEL, TEXT_PADDING)) < textWidth {
		bounds.Width = textWidth + float32(2*getStyle(LABEL, BORDER_WIDTH)) + float32(2*getStyle(LABEL, TEXT_PADDING)) + 2
	}

	if state != STATE_DISABLED && !guiLocked && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()

		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if guiButtonDown() {
				state = STATE_PRESSED
			} else {
				state = STATE_FOCUSED
			}

			if guiButtonReleased() {
				result = RESULT_PRESSED
			}
		}
	}

	guiDrawText(text, getTextBounds(LABEL, bounds), getStyle(LABEL, TEXT_ALIGNMENT), getColor(getStyle(LABEL, textElement+PropertyID(state*3))))

	return result
}

// Toggle - Toggle Button control
func Toggle(bounds rl.Rectangle, text string, active *bool) int32 {
	result := RESULT_NONE
	state := guiState

	var temp bool
	if active == nil {
		active = &temp
	}

	if state != STATE_DISABLED && !guiLocked && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()

		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if guiButtonDown() {
				state = STATE_PRESSED
			} else if guiButtonReleased() {
				state = STATE_NORMAL
				*active = !*active
				result = RESULT_CHANGED
			} else {
				state = STATE_FOCUSED
			}
		}
	}

	if state == STATE_NORMAL {
		borderProp, baseProp, textProp := borderElement+PropertyID(state*3), baseElement+PropertyID(state*3), textElement+PropertyID(state*3)
		if *active {
			borderProp, baseProp, textProp = BORDER_COLOR_PRESSED, BASE_COLOR_PRESSED, TEXT_COLOR_PRESSED
		}
		guiDrawRectangle(bounds, getStyle(TOGGLE, BORDER_WIDTH), getColor(getStyle(TOGGLE, borderProp)), getColor(getStyle(TOGGLE, baseProp)))
		guiDrawText(text, getTextBounds(TOGGLE, bounds), getStyle(TOGGLE, TEXT_ALIGNMENT), getColor(getStyle(TOGGLE, textProp)))
	} else {
		guiDrawRectangle(bounds, getStyle(TOGGLE, BORDER_WIDTH), getColor(getStyle(TOGGLE, borderElement+PropertyID(state*3))), getColor(getStyle(TOGGLE, baseElement+PropertyID(state*3))))
		guiDrawText(text, getTextBounds(TOGGLE, bounds), getStyle(TOGGLE, TEXT_ALIGNMENT), getColor(getStyle(TOGGLE, textElement+PropertyID(state*3))))
	}

	if state == STATE_FOCUSED {
		drawTooltip(bounds)
	}

	return result
}

// ToggleGroup - Toggle Group control
func ToggleGroup(bounds rl.Rectangle, text string, active *int32) int32 {
	result := RESULT_NONE

	var temp int32
	var prevActive int32
	if active != nil {
		prevActive = *active
	}
	if active == nil {
		active = &temp
	}

	toggle := false
	itemReady := false

	initBoundsX := bounds.X
	initBoundsY := bounds.Y

	if getStyle(TOGGLE, GROUP_WIDTH_FULL) != 0 {
		itemCount := 1
		for c := 0; at(text, c) != 0; c++ {
			if at(text, c) == ';' {
				itemCount++
			}
		}
		bounds.Width /= float32(itemCount)
	}

	var itemText strings.Builder
	var itemIndex, row, col int32
	for c, exit := 0, false; !exit; c++ {
		ch := at(text, c)
		if ch == '\n' {
			row++
			col = 0
			itemReady = true
		} else if ch == ';' {
			col++
			itemReady = true
		} else if ch == 0 {
			itemReady = true
			exit = true
		} else if itemText.Len() < toggleGroupMaxItemText-1 {
			itemText.WriteByte(ch)
		}

		if itemReady {
			if itemIndex == *active {
				toggle = true
				Toggle(bounds, itemText.String(), &toggle)
			} else {
				toggle = false
				Toggle(bounds, itemText.String(), &toggle)
				if toggle {
					*active = itemIndex
				}
			}

			bounds.X = initBoundsX + float32(col)*(bounds.Width+float32(getStyle(TOGGLE, GROUP_PADDING)))
			bounds.Y = initBoundsY + float32(row)*(bounds.Height+float32(getStyle(TOGGLE, GROUP_PADDING)))
			itemIndex++
			itemReady = false
			itemText.Reset()
		}
	}

	if prevActive != *active {
		result = RESULT_CHANGED
	}

	return result
}

// ToggleSlider - Toggle Slider control
func ToggleSlider(bounds rl.Rectangle, text string, active *int32) int32 {
	result := RESULT_NONE
	state := guiState

	var temp int32
	var prevActive int32
	if active != nil {
		prevActive = *active
	}
	if active == nil {
		active = &temp
	}

	var items []string
	if text != "" {
		items = guiTextSplit(text, ';')
	}
	itemCount := int32(len(items))

	slider := rl.Rectangle{
		Y:      bounds.Y + float32(getStyle(SLIDER, BORDER_WIDTH)+getStyle(SLIDER, SLIDER_PADDING)),
		Width:  (bounds.Width - float32(2*getStyle(SLIDER, BORDER_WIDTH)) - float32((itemCount+1)*getStyle(SLIDER, SLIDER_PADDING))) / float32(itemCount),
		Height: bounds.Height - float32(2*getStyle(SLIDER, BORDER_WIDTH)) - float32(2*getStyle(SLIDER, SLIDER_PADDING)),
	}

	if state != STATE_DISABLED && !guiLocked {
		mousePoint := guiPointerPosition()

		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if guiButtonDown() {
				state = STATE_PRESSED
			} else if guiButtonReleased() {
				state = STATE_PRESSED
				*active++
				result = 1
			} else {
				state = STATE_FOCUSED
			}
		}

		if *active != 0 && state != STATE_FOCUSED {
			state = STATE_PRESSED
		}
	}

	if *active >= itemCount {
		*active = 0
	}
	slider.X = bounds.X + float32(getStyle(SLIDER, BORDER_WIDTH)) + float32((*active+1)*getStyle(SLIDER, SLIDER_PADDING)) + float32(*active)*slider.Width

	if prevActive != *active {
		result = RESULT_CHANGED
	}

	guiDrawRectangle(bounds, getStyle(SLIDER, BORDER_WIDTH), getColor(getStyle(TOGGLE, borderElement+PropertyID(state*3))), getColor(getStyle(TOGGLE, BASE_COLOR_NORMAL)))

	if state == STATE_NORMAL {
		guiDrawRectangle(slider, 0, blank, getColor(getStyle(SLIDER, BASE_COLOR_PRESSED)))
	} else if state == STATE_FOCUSED {
		guiDrawRectangle(slider, 0, blank, getColor(getStyle(SLIDER, BASE_COLOR_FOCUSED)))
	} else if state == STATE_PRESSED {
		guiDrawRectangle(slider, 0, blank, getColor(getStyle(SLIDER, BASE_COLOR_PRESSED)))
	}

	if text != "" {
		var textBounds rl.Rectangle
		textBounds.Width = float32(GetTextWidth(text))
		textBounds.Height = float32(getStyle(DEFAULT, TEXT_SIZE))
		textBounds.X = slider.X + slider.Width/2 - textBounds.Width/2
		textBounds.Y = bounds.Y + bounds.Height/2 - float32(getStyle(DEFAULT, TEXT_SIZE)/2)

		guiDrawText(item(items, *active), textBounds, getStyle(TOGGLE, TEXT_ALIGNMENT), fade(getColor(getStyle(TOGGLE, textElement+PropertyID(state*3))), guiAlpha))
	}

	return result
}

// CheckBox - Check Box control, returns true when active
func CheckBox(bounds rl.Rectangle, text string, checked *bool) int32 {
	result := RESULT_NONE
	state := guiState

	var temp bool
	if checked == nil {
		checked = &temp
	}

	var textBounds rl.Rectangle

	if text != "" {
		textBounds.Width = float32(GetTextWidth(text)) + 2
		textBounds.Height = float32(getStyle(DEFAULT, TEXT_SIZE))
		textBounds.X = bounds.X + bounds.Width + float32(getStyle(CHECKBOX, TEXT_PADDING))
		textBounds.Y = bounds.Y + bounds.Height/2 - float32(getStyle(DEFAULT, TEXT_SIZE)/2)
		if GetStyle(CHECKBOX, TEXT_ALIGNMENT) == TEXT_ALIGN_LEFT {
			textBounds.X = bounds.X - textBounds.Width - float32(getStyle(CHECKBOX, TEXT_PADDING))
		}
	}

	if state != STATE_DISABLED && !guiLocked && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()

		tx := bounds.X
		if GetStyle(CHECKBOX, TEXT_ALIGNMENT) == TEXT_ALIGN_LEFT {
			tx = textBounds.X
		}
		totalBounds := rl.Rectangle{
			X:      tx,
			Y:      bounds.Y,
			Width:  bounds.Width + textBounds.Width + float32(getStyle(CHECKBOX, TEXT_PADDING)),
			Height: bounds.Height,
		}

		if rl.CheckCollisionPointRec(mousePoint, totalBounds) {
			if guiButtonDown() {
				state = STATE_PRESSED
			} else {
				state = STATE_FOCUSED
			}

			if guiButtonReleased() {
				*checked = !*checked
				result = RESULT_CHANGED
			}
		}
	}

	guiDrawRectangle(bounds, getStyle(CHECKBOX, BORDER_WIDTH), getColor(getStyle(CHECKBOX, borderElement+PropertyID(state*3))), blank)

	if *checked {
		bw, cp := getStyle(CHECKBOX, BORDER_WIDTH), getStyle(CHECKBOX, CHECK_PADDING)
		check := rl.Rectangle{X: bounds.X + float32(bw+cp),
			Y:      bounds.Y + float32(bw+cp),
			Width:  bounds.Width - float32(2*(bw+cp)),
			Height: bounds.Height - float32(2*(bw+cp))}
		guiDrawRectangle(check, 0, blank, getColor(getStyle(CHECKBOX, textElement+PropertyID(state*3))))
	}

	align := int32(TEXT_ALIGN_RIGHT)
	if GetStyle(CHECKBOX, TEXT_ALIGNMENT) == TEXT_ALIGN_RIGHT {
		align = int32(TEXT_ALIGN_LEFT)
	}
	guiDrawText(text, textBounds, align, getColor(getStyle(LABEL, textElement+PropertyID(state*3))))

	return result
}

func item(items []string, i int32) string {
	if i < 0 || int(i) >= len(items) {
		return ""
	}
	return items[i]
}

// ComboBox - Combo Box control
func ComboBox(bounds rl.Rectangle, text string, active *int32) int32 {
	result := RESULT_NONE
	state := guiState

	var temp int32
	var prevActive int32
	if active != nil {
		prevActive = *active
	}
	if active == nil {
		active = &temp
	}

	bounds.Width -= float32(getStyle(COMBOBOX, COMBO_BUTTON_WIDTH) + getStyle(COMBOBOX, COMBO_BUTTON_SPACING))

	selector := rl.Rectangle{X: bounds.X + bounds.Width + float32(getStyle(COMBOBOX, COMBO_BUTTON_SPACING)),
		Y: bounds.Y, Width: float32(getStyle(COMBOBOX, COMBO_BUTTON_WIDTH)), Height: bounds.Height}

	items := guiTextSplit(text, ';')
	itemCount := int32(len(items))

	if *active < 0 {
		*active = 0
	} else if *active > itemCount-1 {
		*active = itemCount - 1
	}

	if state != STATE_DISABLED && !guiLocked && itemCount > 1 && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()

		if rl.CheckCollisionPointRec(mousePoint, bounds) || rl.CheckCollisionPointRec(mousePoint, selector) {
			if guiButtonDown() {
				state = STATE_PRESSED
			} else {
				state = STATE_FOCUSED
			}

			if guiButtonPressed() {
				*active += 1
				if *active >= itemCount {
					*active = 0
				}
			}
		}
	}

	if prevActive != *active {
		result = RESULT_CHANGED
	}

	guiDrawRectangle(bounds, getStyle(COMBOBOX, BORDER_WIDTH), getColor(getStyle(COMBOBOX, borderElement+PropertyID(state*3))), getColor(getStyle(COMBOBOX, baseElement+PropertyID(state*3))))
	guiDrawText(item(items, *active), getTextBounds(COMBOBOX, bounds), getStyle(COMBOBOX, TEXT_ALIGNMENT), getColor(getStyle(COMBOBOX, textElement+PropertyID(state*3))))

	tempBorderWidth := getStyle(BUTTON, BORDER_WIDTH)
	tempTextAlign := getStyle(BUTTON, TEXT_ALIGNMENT)
	setStyle(BUTTON, BORDER_WIDTH, 1)
	setStyle(BUTTON, TEXT_ALIGNMENT, int32(TEXT_ALIGN_CENTER))

	Button(selector, fmt.Sprintf("%d/%d", *active+1, itemCount))

	setStyle(BUTTON, TEXT_ALIGNMENT, tempTextAlign)
	setStyle(BUTTON, BORDER_WIDTH, tempBorderWidth)

	return result
}

// DropdownBox - Dropdown Box control
func DropdownBox(bounds rl.Rectangle, text string, active *int32, editMode bool) int32 {
	result := RESULT_NONE
	state := guiState

	var temp int32
	var prevActive int32
	if active != nil {
		prevActive = *active
	}
	if active == nil {
		active = &temp
	}

	itemSelected := *active
	itemFocused := int32(-1)

	direction := 0
	if getStyle(DROPDOWNBOX, DROPDOWN_ROLL_UP) == 1 {
		direction = 1
	}

	items := guiTextSplit(text, ';')
	itemCount := int32(len(items))

	spacing := float32(getStyle(DROPDOWNBOX, DROPDOWN_ITEMS_SPACING))
	boundsOpen := bounds
	boundsOpen.Height = float32(itemCount+1) * (bounds.Height + spacing)
	if direction == 1 {
		boundsOpen.Y -= float32(itemCount)*(bounds.Height+spacing) + spacing
	}

	itemBounds := bounds

	if state != STATE_DISABLED && (editMode || !guiLocked) && itemCount > 1 && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()

		if editMode {
			state = STATE_PRESSED

			if !rl.CheckCollisionPointRec(mousePoint, boundsOpen) {
				if guiButtonPressed() || guiButtonReleased() {
					result = 1
				}
			}

			if rl.CheckCollisionPointRec(mousePoint, bounds) && guiButtonPressed() {
				result = 1
			}

			for i := range itemCount {
				if direction == 0 {
					itemBounds.Y += bounds.Height + spacing
				} else {
					itemBounds.Y -= bounds.Height + spacing
				}

				if rl.CheckCollisionPointRec(mousePoint, itemBounds) {
					itemFocused = i
					if guiButtonReleased() {
						itemSelected = i
						result = 1
					}
					break
				}
			}

			itemBounds = bounds
		} else {
			if rl.CheckCollisionPointRec(mousePoint, bounds) {
				if guiButtonPressed() {
					result = 1
					state = STATE_PRESSED
				} else {
					state = STATE_FOCUSED
				}
			}
		}
	}

	if editMode {
		Panel(boundsOpen, "")
	}

	guiDrawRectangle(bounds, getStyle(DROPDOWNBOX, BORDER_WIDTH), getColor(getStyle(DROPDOWNBOX, borderElement+PropertyID(state*3))), getColor(getStyle(DROPDOWNBOX, baseElement+PropertyID(state*3))))
	guiDrawText(item(items, itemSelected), getTextBounds(DROPDOWNBOX, bounds), getStyle(DROPDOWNBOX, TEXT_ALIGNMENT), getColor(getStyle(DROPDOWNBOX, textElement+PropertyID(state*3))))

	if editMode {
		for i := range itemCount {
			if direction == 0 {
				itemBounds.Y += bounds.Height + spacing
			} else {
				itemBounds.Y -= bounds.Height + spacing
			}

			if i == itemSelected {
				guiDrawRectangle(itemBounds, getStyle(DROPDOWNBOX, BORDER_WIDTH), getColor(getStyle(DROPDOWNBOX, BORDER_COLOR_PRESSED)), getColor(getStyle(DROPDOWNBOX, BASE_COLOR_PRESSED)))
				guiDrawText(items[i], getTextBounds(DROPDOWNBOX, itemBounds), getStyle(DROPDOWNBOX, TEXT_ALIGNMENT), getColor(getStyle(DROPDOWNBOX, TEXT_COLOR_PRESSED)))
			} else if i == itemFocused {
				guiDrawRectangle(itemBounds, getStyle(DROPDOWNBOX, BORDER_WIDTH), getColor(getStyle(DROPDOWNBOX, BORDER_COLOR_FOCUSED)), getColor(getStyle(DROPDOWNBOX, BASE_COLOR_FOCUSED)))
				guiDrawText(items[i], getTextBounds(DROPDOWNBOX, itemBounds), getStyle(DROPDOWNBOX, TEXT_ALIGNMENT), getColor(getStyle(DROPDOWNBOX, TEXT_COLOR_FOCUSED)))
			} else {
				guiDrawText(items[i], getTextBounds(DROPDOWNBOX, itemBounds), getStyle(DROPDOWNBOX, TEXT_ALIGNMENT), getColor(getStyle(DROPDOWNBOX, TEXT_COLOR_NORMAL)))
			}
		}
	}

	if getStyle(DROPDOWNBOX, DROPDOWN_ARROW_HIDDEN) == 0 {
		arrow := IconText(ICON_ARROW_DOWN_FILL, "")
		if direction != 0 {
			arrow = IconText(ICON_ARROW_UP_FILL, "")
		}
		guiDrawText(arrow, rl.Rectangle{X: bounds.X + bounds.Width - float32(getStyle(DROPDOWNBOX, ARROW_PADDING)), Y: bounds.Y + bounds.Height/2 - 6, Width: 10, Height: 10},
			int32(TEXT_ALIGN_CENTER), getColor(getStyle(DROPDOWNBOX, textElement+PropertyID(state*3))))
	}

	*active = itemSelected

	if prevActive != *active {
		result = RESULT_CHANGED
	}

	return result
}

const (
	textBoxAutoCursorCooldown = 20
	textBoxAutoCursorDelay    = 1
)

func isPunct(c int32) bool {
	return (c >= 33 && c <= 47) || (c >= 58 && c <= 64) || (c >= 91 && c <= 96) || (c >= 123 && c <= 126)
}

func isSpace(c int32) bool {
	return c == ' ' || (c >= '\t' && c <= '\r')
}

type textBuffer struct {
	b []byte
}

func newTextBuffer(text *string, textSize int32) *textBuffer {
	n := int(textSize)
	if text != nil && len(*text)+1 > n {
		n = len(*text) + 1
	}
	if n < 1 {
		n = 1
	}
	t := &textBuffer{b: make([]byte, n)}
	if text != nil {
		copy(t.b, *text)
	}
	return t
}

func (t *textBuffer) str() string {
	return cstr(string(t.b))
}

func (t *textBuffer) from(i int32) string {
	if i < 0 {
		i = 0
	}
	if int(i) >= len(t.b) {
		return ""
	}
	return cstr(string(t.b[i:]))
}

func (t *textBuffer) set(i int32, c byte) {
	if i >= 0 && int(i) < len(t.b) {
		t.b[i] = c
	}
}

func (t *textBuffer) get(i int32) byte {
	if i >= 0 && int(i) < len(t.b) {
		return t.b[i]
	}
	return 0
}

// TextBox - Text Box control, updates input text
func TextBox(bounds rl.Rectangle, text *string, textSize int32, editMode bool) int32 {
	result := RESULT_NONE
	state := guiState

	multiline := false
	wrapMode := GetStyle(DEFAULT, TEXT_WRAP_MODE)

	textBounds := getTextBounds(TEXTBOX, bounds)
	buf := newTextBuffer(text, textSize)
	textLength := int32(len(buf.str()))
	thisCursorIndex := min(textBoxCursorIndex, textLength)
	textWidth := GetTextWidth(buf.str()) - GetTextWidth(buf.from(thisCursorIndex))
	var textIndexOffset int32

	cursor := rl.Rectangle{
		X:      textBounds.X + float32(textWidth) + float32(getStyle(DEFAULT, TEXT_SPACING)),
		Y:      textBounds.Y + textBounds.Height/2 - float32(getStyle(DEFAULT, TEXT_SIZE)),
		Width:  2,
		Height: float32(getStyle(DEFAULT, TEXT_SIZE)) * 2,
	}

	if cursor.Height >= bounds.Height {
		cursor.Height = bounds.Height - float32(getStyle(TEXTBOX, BORDER_WIDTH)*2)
	}
	if cursor.Y < bounds.Y+float32(getStyle(TEXTBOX, BORDER_WIDTH)) {
		cursor.Y = bounds.Y + float32(getStyle(TEXTBOX, BORDER_WIDTH))
	}

	mouseCursor := cursor
	mouseCursor.X = -1
	mouseCursor.Width = 1

	ctrlDown := func() bool { return guiKeyDown(rl.KeyLeftControl) || guiKeyDown(rl.KeyRightControl) }

	if state != STATE_DISABLED && getStyle(TEXTBOX, TEXT_READONLY) == 0 && !guiLocked && !guiControlExclusiveMode && wrapMode == TEXT_WRAP_NONE {
		mousePosition := guiPointerPosition()

		if editMode {
			if guiKeyDown(rl.KeyLeft) || guiKeyDown(rl.KeyRight) || guiKeyDown(rl.KeyUp) || guiKeyDown(rl.KeyDown) || guiKeyDown(rl.KeyBackspace) || guiKeyDown(rl.KeyDelete) {
				autoCursorCounter++
			} else {
				autoCursorCounter = 0
			}

			autoCursorShouldTrigger := autoCursorCounter > textBoxAutoCursorCooldown && autoCursorCounter%textBoxAutoCursorDelay == 0

			state = STATE_PRESSED

			if textBoxCursorIndex > textLength {
				textBoxCursorIndex = textLength
			}

			for float32(textWidth) >= textBounds.Width {
				var nextCodepointSize int32
				getCodepointNext(buf.from(textIndexOffset), &nextCodepointSize)
				textIndexOffset += nextCodepointSize
				textWidth = GetTextWidth(buf.from(textIndexOffset)) - GetTextWidth(buf.from(textBoxCursorIndex))
			}

			codepoint := guiInputKey()
			if multiline && guiKeyPressed(rl.KeyEnter) {
				codepoint = '\n'
			}

			charEncoded := codepointToUTF8(codepoint)
			codepointSize := int32(len(charEncoded))

			if guiKeyPressed(rl.KeyV) && ctrlDown() {
				pasteText := rl.GetClipboardText()

				var pasteLength int32
				var pasteCodepointSize int32

				for {
					pasteCodepoint := getCodepointNext(tail(pasteText, int(pasteLength)), &pasteCodepointSize)
					if textLength+pasteLength+pasteCodepointSize >= textSize {
						break
					}
					if !(multiline && pasteCodepoint == '\n') && !(pasteCodepoint >= 32) {
						break
					}
					pasteLength += pasteCodepointSize
				}

				if pasteLength > 0 {
					for i := textLength + pasteLength; i >= textBoxCursorIndex+pasteLength; i-- {
						buf.set(i, buf.get(i-pasteLength))
					}

					for i := int32(0); i < pasteLength; i++ {
						buf.set(textBoxCursorIndex+i, at(pasteText, int(i)))
					}

					textBoxCursorIndex += pasteLength
					textLength += pasteLength
					buf.set(textLength, 0)
				}
			} else if ((multiline && codepoint == '\n') || codepoint >= 32) && textLength+codepointSize < textSize {
				for i := textLength + codepointSize; i >= textBoxCursorIndex+codepointSize; i-- {
					buf.set(i, buf.get(i-codepointSize))
				}

				for i := int32(0); i < codepointSize; i++ {
					buf.set(textBoxCursorIndex+i, charEncoded[i])
				}

				textBoxCursorIndex += codepointSize
				textLength += codepointSize

				buf.set(textLength, 0)
			}

			if textLength > 0 && guiKeyPressed(rl.KeyHome) {
				textBoxCursorIndex = 0
			}

			if textLength > textBoxCursorIndex && guiKeyPressed(rl.KeyEnd) {
				textBoxCursorIndex = textLength
			}

			if textLength > textBoxCursorIndex && guiKeyPressed(rl.KeyDelete) && ctrlDown() {
				offset := textBoxCursorIndex
				var accCodepointSize int32
				var nextCodepointSize int32

				nextCodepoint := getCodepointNext(buf.from(offset), &nextCodepointSize)
				puctuation := isPunct(nextCodepoint & 0xff)
				for offset < textLength {
					if (puctuation && !isPunct(nextCodepoint&0xff)) || (!puctuation && (isSpace(nextCodepoint&0xff) || isPunct(nextCodepoint&0xff))) {
						break
					}
					offset += nextCodepointSize
					accCodepointSize += nextCodepointSize
					nextCodepoint = getCodepointNext(buf.from(offset), &nextCodepointSize)
				}

				for offset < textLength {
					if !isSpace(nextCodepoint & 0xff) {
						break
					}
					offset += nextCodepointSize
					accCodepointSize += nextCodepointSize
					nextCodepoint = getCodepointNext(buf.from(offset), &nextCodepointSize)
				}

				for i := offset; i <= textLength; i++ {
					buf.set(i-accCodepointSize, buf.get(i))
				}

				textLength -= accCodepointSize
			} else if textLength > textBoxCursorIndex && (guiKeyPressed(rl.KeyDelete) || (guiKeyDown(rl.KeyDelete) && autoCursorShouldTrigger)) {
				var nextCodepointSize int32
				getCodepointNext(buf.from(textBoxCursorIndex), &nextCodepointSize)

				for i := textBoxCursorIndex + nextCodepointSize; i <= textLength; i++ {
					buf.set(i-nextCodepointSize, buf.get(i))
				}

				textLength -= nextCodepointSize
			}

			if textBoxCursorIndex > 0 && guiKeyPressed(rl.KeyBackspace) && ctrlDown() {
				offset := textBoxCursorIndex
				var accCodepointSize int32
				var prevCodepointSize int32
				var prevCodepoint int32

				for offset > 0 {
					prevCodepoint = getCodepointPrevious(buf.str(), int(offset), &prevCodepointSize)
					if !isSpace(prevCodepoint & 0xff) {
						break
					}

					offset -= prevCodepointSize
					accCodepointSize += prevCodepointSize
				}

				puctuation := isPunct(prevCodepoint & 0xff)
				for offset > 0 {
					prevCodepoint = getCodepointPrevious(buf.str(), int(offset), &prevCodepointSize)
					if (puctuation && !isPunct(prevCodepoint&0xff)) || (!puctuation && (isSpace(prevCodepoint&0xff) || isPunct(prevCodepoint&0xff))) {
						break
					}

					offset -= prevCodepointSize
					accCodepointSize += prevCodepointSize
				}

				for i := textBoxCursorIndex; i <= textLength; i++ {
					buf.set(i-accCodepointSize, buf.get(i))
				}

				textLength -= accCodepointSize
				textBoxCursorIndex -= accCodepointSize
			} else if textBoxCursorIndex > 0 && (guiKeyPressed(rl.KeyBackspace) || (guiKeyDown(rl.KeyBackspace) && autoCursorShouldTrigger)) {
				var prevCodepointSize int32

				getCodepointPrevious(buf.str(), int(textBoxCursorIndex), &prevCodepointSize)

				for i := textBoxCursorIndex; i <= textLength; i++ {
					buf.set(i-prevCodepointSize, buf.get(i))
				}

				textLength -= prevCodepointSize
				textBoxCursorIndex -= prevCodepointSize
			}

			if textBoxCursorIndex > 0 && guiKeyPressed(rl.KeyLeft) && ctrlDown() {
				offset := textBoxCursorIndex
				var prevCodepointSize int32
				var prevCodepoint int32

				for offset > 0 {
					prevCodepoint = getCodepointPrevious(buf.str(), int(offset), &prevCodepointSize)
					if !isSpace(prevCodepoint & 0xff) {
						break
					}

					offset -= prevCodepointSize
				}

				puctuation := isPunct(prevCodepoint & 0xff)
				for offset > 0 {
					prevCodepoint = getCodepointPrevious(buf.str(), int(offset), &prevCodepointSize)
					if (puctuation && !isPunct(prevCodepoint&0xff)) || (!puctuation && (isSpace(prevCodepoint&0xff) || isPunct(prevCodepoint&0xff))) {
						break
					}

					offset -= prevCodepointSize
				}

				textBoxCursorIndex = offset
			} else if textBoxCursorIndex > 0 && (guiKeyPressed(rl.KeyLeft) || (guiKeyDown(rl.KeyLeft) && autoCursorShouldTrigger)) {
				var prevCodepointSize int32
				getCodepointPrevious(buf.str(), int(textBoxCursorIndex), &prevCodepointSize)

				textBoxCursorIndex -= prevCodepointSize
			} else if textLength > textBoxCursorIndex && guiKeyPressed(rl.KeyRight) && ctrlDown() {
				offset := textBoxCursorIndex
				var nextCodepointSize int32

				nextCodepoint := getCodepointNext(buf.from(offset), &nextCodepointSize)
				puctuation := isPunct(nextCodepoint & 0xff)
				for offset < textLength {
					if (puctuation && !isPunct(nextCodepoint&0xff)) || (!puctuation && (isSpace(nextCodepoint&0xff) || isPunct(nextCodepoint&0xff))) {
						break
					}

					offset += nextCodepointSize
					nextCodepoint = getCodepointNext(buf.from(offset), &nextCodepointSize)
				}

				for offset < textLength {
					if !isSpace(nextCodepoint & 0xff) {
						break
					}

					offset += nextCodepointSize
					nextCodepoint = getCodepointNext(buf.from(offset), &nextCodepointSize)
				}

				textBoxCursorIndex = offset
			} else if textLength > textBoxCursorIndex && (guiKeyPressed(rl.KeyRight) || (guiKeyDown(rl.KeyRight) && autoCursorShouldTrigger)) {
				var nextCodepointSize int32
				getCodepointNext(buf.from(textBoxCursorIndex), &nextCodepointSize)

				textBoxCursorIndex += nextCodepointSize
			}

			if rl.CheckCollisionPointRec(mousePosition, textBounds) {
				scaleFactor := float32(getStyle(DEFAULT, TEXT_SIZE)) / float32(guiFont.BaseSize)
				var glyphWidth float32
				var widthToMouseX float32
				var mouseCursorIndex int32

				for i := textIndexOffset; i < textLength; i += codepointSize {
					codepoint = getCodepointNext(buf.from(i), &codepointSize)
					codepointIndex := rl.GetGlyphIndex(guiFont, codepoint)

					glyphWidth = glyphWidthAt(codepointIndex, scaleFactor)

					if mousePosition.X <= textBounds.X+(widthToMouseX+glyphWidth/2) {
						mouseCursor.X = textBounds.X + widthToMouseX
						mouseCursorIndex = i
						break
					}

					widthToMouseX += glyphWidth + float32(getStyle(DEFAULT, TEXT_SPACING))
				}

				textEndWidth := GetTextWidth(buf.from(textIndexOffset))
				if guiPointerPosition().X >= textBounds.X+float32(textEndWidth)-glyphWidth/2 {
					mouseCursor.X = textBounds.X + float32(textEndWidth)
					mouseCursorIndex = textLength
				}

				if mouseCursor.X >= 0 && guiButtonPressed() {
					cursor.X = mouseCursor.X
					textBoxCursorIndex = mouseCursorIndex
				}
			} else {
				mouseCursor.X = -1
			}

			cursor.X = bounds.X + float32(getStyle(TEXTBOX, TEXT_PADDING)) + float32(GetTextWidth(buf.from(textIndexOffset))) - float32(GetTextWidth(buf.from(textBoxCursorIndex))) + float32(getStyle(DEFAULT, TEXT_SPACING))

			if (!multiline && guiKeyPressed(rl.KeyEnter)) || (!rl.CheckCollisionPointRec(mousePosition, bounds) && guiButtonPressed()) {
				textBoxCursorIndex = 0
				autoCursorCounter = 0
				result = RESULT_PRESSED
			}
		} else {
			if rl.CheckCollisionPointRec(mousePosition, bounds) {
				state = STATE_FOCUSED

				if guiButtonPressed() {
					textBoxCursorIndex = textLength
					autoCursorCounter = 0
					result = RESULT_PRESSED
				}
			}
		}
	}

	if state == STATE_PRESSED {
		guiDrawRectangle(bounds, getStyle(TEXTBOX, BORDER_WIDTH), getColor(getStyle(TEXTBOX, borderElement+PropertyID(state*3))), getColor(getStyle(TEXTBOX, BASE_COLOR_PRESSED)))
	} else if state == STATE_DISABLED {
		guiDrawRectangle(bounds, getStyle(TEXTBOX, BORDER_WIDTH), getColor(getStyle(TEXTBOX, borderElement+PropertyID(state*3))), getColor(getStyle(TEXTBOX, BASE_COLOR_DISABLED)))
	} else {
		guiDrawRectangle(bounds, getStyle(TEXTBOX, BORDER_WIDTH), getColor(getStyle(TEXTBOX, borderElement+PropertyID(state*3))), blank)
	}

	guiDrawText(buf.from(textIndexOffset), textBounds, getStyle(TEXTBOX, TEXT_ALIGNMENT), getColor(getStyle(TEXTBOX, textElement+PropertyID(state*3))))

	if editMode && getStyle(TEXTBOX, TEXT_READONLY) == 0 {
		guiDrawRectangle(cursor, 0, blank, getColor(getStyle(TEXTBOX, BORDER_COLOR_PRESSED)))

		if mouseCursor.X >= 0 {
			guiDrawRectangle(mouseCursor, 0, blank, getColor(getStyle(TEXTBOX, BORDER_COLOR_PRESSED)))
		}
	} else if state == STATE_FOCUSED {
		drawTooltip(bounds)
	}

	if text != nil {
		*text = buf.str()
	}

	return result
}

const valueBoxMaxChars = 32

func sideTextBounds(control ControlID, bounds rl.Rectangle, text string) rl.Rectangle {
	var textBounds rl.Rectangle
	if text != "" {
		textBounds.Width = float32(GetTextWidth(text)) + 2
		textBounds.Height = float32(getStyle(DEFAULT, TEXT_SIZE))
		textBounds.X = bounds.X + bounds.Width + float32(getStyle(control, TEXT_PADDING))
		textBounds.Y = bounds.Y + bounds.Height/2 - float32(getStyle(DEFAULT, TEXT_SIZE)/2)
		if GetStyle(control, TEXT_ALIGNMENT) == TEXT_ALIGN_LEFT {
			textBounds.X = bounds.X - textBounds.Width - float32(getStyle(control, TEXT_PADDING))
		}
	}
	return textBounds
}

func sideTextAlign(control ControlID) int32 {
	if GetStyle(control, TEXT_ALIGNMENT) == TEXT_ALIGN_RIGHT {
		return int32(TEXT_ALIGN_LEFT)
	}
	return int32(TEXT_ALIGN_RIGHT)
}

// Spinner - Spinner control
func Spinner(bounds rl.Rectangle, text string, value *int32, minValue, maxValue int32, editMode bool) int32 {
	state := guiState

	var temp int32
	if value == nil {
		value = &temp
	}
	tempValue := *value

	bw, bs := getStyle(VALUEBOX, SPINNER_BUTTON_WIDTH), getStyle(VALUEBOX, SPINNER_BUTTON_SPACING)
	valueBoxBounds := rl.Rectangle{
		X:      bounds.X + float32(bw+bs),
		Y:      bounds.Y,
		Width:  bounds.Width - float32(2*(bw+bs)),
		Height: bounds.Height,
	}

	leftButtonBound := rl.Rectangle{X: bounds.X, Y: bounds.Y, Width: float32(bw), Height: bounds.Height}
	rightButtonBound := rl.Rectangle{X: bounds.X + bounds.Width - float32(bw), Y: bounds.Y, Width: float32(bw), Height: bounds.Height}

	textBounds := sideTextBounds(VALUEBOX, bounds, text)

	if state != STATE_DISABLED && !guiLocked && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()

		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if guiButtonDown() {
				state = STATE_PRESSED
			} else {
				state = STATE_FOCUSED
			}
		}
	}

	if Button(leftButtonBound, IconText(ICON_ARROW_LEFT_FILL, "")) != 0 {
		tempValue--
	}
	if Button(rightButtonBound, IconText(ICON_ARROW_RIGHT_FILL, "")) != 0 {
		tempValue++
	}

	if !editMode {
		if tempValue < minValue {
			tempValue = minValue
		}
		if tempValue > maxValue {
			tempValue = maxValue
		}
	}

	result := ValueBox(valueBoxBounds, "", &tempValue, minValue, maxValue, editMode)

	tempBorderWidth := getStyle(BUTTON, BORDER_WIDTH)
	tempTextAlign := getStyle(BUTTON, TEXT_ALIGNMENT)
	setStyle(BUTTON, BORDER_WIDTH, getStyle(VALUEBOX, BORDER_WIDTH))
	setStyle(BUTTON, TEXT_ALIGNMENT, int32(TEXT_ALIGN_CENTER))
	setStyle(BUTTON, TEXT_ALIGNMENT, tempTextAlign)
	setStyle(BUTTON, BORDER_WIDTH, tempBorderWidth)

	guiDrawText(text, textBounds, sideTextAlign(VALUEBOX), getColor(getStyle(LABEL, textElement+PropertyID(state*3))))

	*value = tempValue

	return result
}

func valueBoxBase(state PropertyValue) rl.Color {
	baseColor := blank
	if state == STATE_PRESSED {
		baseColor = getColor(getStyle(VALUEBOX, BASE_COLOR_PRESSED))
	} else if state == STATE_DISABLED {
		baseColor = getColor(getStyle(VALUEBOX, BASE_COLOR_DISABLED))
	}
	return baseColor
}

// ValueBox - Value Box control, updates input text with numbers
func ValueBox(bounds rl.Rectangle, text string, value *int32, minValue, maxValue int32, editMode bool) int32 {
	result := RESULT_NONE
	state := guiState

	var temp int32
	if value == nil {
		value = &temp
	}

	var textValue [valueBoxMaxChars + 1]byte
	s := fmt.Sprintf("%d", *value)
	if len(s) > valueBoxMaxChars {
		s = s[:valueBoxMaxChars]
	}
	copy(textValue[:], s)
	str := func() string { return cstr(string(textValue[:])) }

	textBounds := sideTextBounds(VALUEBOX, bounds, text)

	if state != STATE_DISABLED && !guiLocked && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()
		valueHasChanged := false

		if editMode {
			state = STATE_PRESSED

			keyCount := len(str())

			if guiKeyPressed(rl.KeyMinus) {
				if textValue[0] == '-' {
					for i := 0; i < keyCount; i++ {
						textValue[i] = textValue[i+1]
					}

					keyCount--
					valueHasChanged = true
				} else if keyCount < valueBoxMaxChars {
					if keyCount == 0 {
						textValue[0] = '0'
						textValue[1] = 0
						keyCount++
					}

					for i := keyCount; i > -1; i-- {
						textValue[i+1] = textValue[i]
					}

					textValue[0] = '-'
					keyCount++
					valueHasChanged = true
				}
			}

			if keyCount >= 0 && keyCount < valueBoxMaxChars && float32(GetTextWidth(str())) < bounds.Width {
				key := guiInputKey()

				if key >= 48 && key <= 57 {
					textValue[keyCount] = byte(key)
					keyCount++
					valueHasChanged = true
				}
			}

			if keyCount > 0 && guiKeyPressed(rl.KeyBackspace) {
				keyCount--
				textValue[keyCount] = 0
				valueHasChanged = true
			}

			if valueHasChanged {
				*value = textToInteger(str())
			}

			if guiKeyPressed(rl.KeyEnter) || guiKeyPressed(rl.KeyKpEnter) || (!rl.CheckCollisionPointRec(mousePoint, bounds) && guiButtonPressed()) {
				if *value > maxValue {
					*value = maxValue
				} else if *value < minValue {
					*value = minValue
				}

				result = RESULT_PRESSED
			}
		} else {
			if *value > maxValue {
				*value = maxValue
			} else if *value < minValue {
				*value = minValue
			}

			if rl.CheckCollisionPointRec(mousePoint, bounds) {
				state = STATE_FOCUSED
				if guiButtonPressed() {
					result = RESULT_PRESSED
				}
			}
		}
	}

	guiDrawRectangle(bounds, getStyle(VALUEBOX, BORDER_WIDTH), getColor(getStyle(VALUEBOX, borderElement+PropertyID(state*3))), valueBoxBase(state))
	guiDrawText(str(), getTextBounds(VALUEBOX, bounds), int32(TEXT_ALIGN_CENTER), getColor(getStyle(VALUEBOX, textElement+PropertyID(state*3))))

	if editMode {
		cursor := rl.Rectangle{X: bounds.X + float32(GetTextWidth(str())/2) + bounds.Width/2 + 1,
			Y:      bounds.Y + float32(getStyle(TEXTBOX, BORDER_WIDTH)) + 2,
			Width:  2,
			Height: bounds.Height - float32(getStyle(TEXTBOX, BORDER_WIDTH)*2) - 4}
		if cursor.Height > bounds.Height {
			cursor.Height = bounds.Height - float32(getStyle(TEXTBOX, BORDER_WIDTH)*2)
		}
		guiDrawRectangle(cursor, 0, blank, getColor(getStyle(VALUEBOX, BORDER_COLOR_PRESSED)))
	}

	guiDrawText(text, textBounds, sideTextAlign(VALUEBOX), getColor(getStyle(LABEL, textElement+PropertyID(state*3))))

	return result
}

// ValueBoxFloat - Value box control for float values
func ValueBoxFloat(bounds rl.Rectangle, text string, textValue *string, value *float32, editMode bool) int32 {
	result := RESULT_NONE
	state := guiState

	var tempValue float32
	if value == nil {
		value = &tempValue
	}
	var tempText string
	if textValue == nil {
		textValue = &tempText
	}

	buf := newTextBuffer(textValue, valueBoxMaxChars+1)
	tv := buf.b

	textBounds := sideTextBounds(VALUEBOX, bounds, text)

	if state != STATE_DISABLED && !guiLocked && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()
		valueHasChanged := false

		if editMode {
			state = STATE_PRESSED

			keyCount := len(buf.str())

			if guiKeyPressed(rl.KeyMinus) {
				if tv[0] == '-' {
					for i := 0; i < keyCount; i++ {
						tv[i] = tv[i+1]
					}

					keyCount--
					valueHasChanged = true
				} else if keyCount < valueBoxMaxChars-1 {
					if keyCount == 0 {
						tv[0] = '0'
						tv[1] = 0
						keyCount++
					}

					for i := keyCount; i > -1; i-- {
						tv[i+1] = tv[i]
					}

					tv[0] = '-'
					keyCount++
					valueHasChanged = true
				}
			}

			if keyCount < valueBoxMaxChars {
				if float32(GetTextWidth(buf.str())) < bounds.Width {
					key := guiInputKey()
					if (key >= 48 && key <= 57) || key == '.' || (keyCount == 0 && key == '+') || (keyCount == 0 && key == '-') {
						tv[keyCount] = byte(key)
						keyCount++
						valueHasChanged = true
					}
				}
			}

			if guiKeyPressed(rl.KeyBackspace) {
				if keyCount > 0 {
					keyCount--
					tv[keyCount] = 0
					valueHasChanged = true
				}
			}

			if valueHasChanged {
				*value = textToFloat(buf.str())
			}

			if guiKeyPressed(rl.KeyEnter) || guiKeyPressed(rl.KeyKpEnter) || (!rl.CheckCollisionPointRec(mousePoint, bounds) && guiButtonPressed()) {
				result = RESULT_PRESSED
			}
		} else {
			if rl.CheckCollisionPointRec(mousePoint, bounds) {
				state = STATE_FOCUSED
				if guiButtonPressed() {
					result = RESULT_PRESSED
				}
			}
		}
	}

	guiDrawRectangle(bounds, getStyle(VALUEBOX, BORDER_WIDTH), getColor(getStyle(VALUEBOX, borderElement+PropertyID(state*3))), valueBoxBase(state))
	guiDrawText(buf.str(), getTextBounds(VALUEBOX, bounds), int32(TEXT_ALIGN_CENTER), getColor(getStyle(VALUEBOX, textElement+PropertyID(state*3))))

	if editMode {
		cursor := rl.Rectangle{X: bounds.X + float32(GetTextWidth(buf.str())/2) + bounds.Width/2 + 1,
			Y: bounds.Y + float32(2*getStyle(VALUEBOX, BORDER_WIDTH)), Width: 4,
			Height: bounds.Height - float32(4*getStyle(VALUEBOX, BORDER_WIDTH))}
		guiDrawRectangle(cursor, 0, blank, getColor(getStyle(VALUEBOX, BORDER_COLOR_PRESSED)))
	}

	guiDrawText(text, textBounds, sideTextAlign(VALUEBOX), getColor(getStyle(LABEL, textElement+PropertyID(state*3))))

	*textValue = buf.str()

	return result
}

func leftRightText(control ControlID, bounds rl.Rectangle, textLeft, textRight string, state PropertyValue) {
	if textLeft != "" {
		var textBounds rl.Rectangle
		textBounds.Width = float32(GetTextWidth(textLeft))
		textBounds.Height = float32(getStyle(DEFAULT, TEXT_SIZE))
		textBounds.X = bounds.X - textBounds.Width - float32(getStyle(control, TEXT_PADDING))
		textBounds.Y = bounds.Y + bounds.Height/2 - float32(getStyle(DEFAULT, TEXT_SIZE)/2)

		guiDrawText(textLeft, textBounds, int32(TEXT_ALIGN_RIGHT), getColor(getStyle(LABEL, textElement+PropertyID(state*3))))
	}

	if textRight != "" {
		var textBounds rl.Rectangle
		textBounds.Width = float32(GetTextWidth(textRight))
		textBounds.Height = float32(getStyle(DEFAULT, TEXT_SIZE))
		textBounds.X = bounds.X + bounds.Width + float32(getStyle(control, TEXT_PADDING))
		textBounds.Y = bounds.Y + bounds.Height/2 - float32(getStyle(DEFAULT, TEXT_SIZE)/2)

		guiDrawText(textRight, textBounds, int32(TEXT_ALIGN_LEFT), getColor(getStyle(LABEL, textElement+PropertyID(state*3))))
	}
}

// Slider - Slider control
func Slider(bounds rl.Rectangle, textLeft, textRight string, value *float32, minValue, maxValue float32) int32 {
	result := RESULT_NONE
	state := guiState

	temp := (maxValue - minValue) / 2.0
	if value == nil {
		value = &temp
	}
	prevValue := *value

	sliderWidth := getStyle(SLIDER, SLIDER_WIDTH)
	bw := getStyle(SLIDER, BORDER_WIDTH)

	slider := rl.Rectangle{X: bounds.X, Y: bounds.Y + float32(bw+getStyle(SLIDER, SLIDER_PADDING)),
		Height: bounds.Height - float32(2*bw) - float32(2*getStyle(SLIDER, SLIDER_PADDING))}

	if state != STATE_DISABLED && !guiLocked {
		mousePoint := guiPointerPosition()

		if guiControlExclusiveMode {
			if guiButtonDown() {
				if checkBoundsID(bounds, guiControlExclusiveRec) {
					state = STATE_PRESSED
					*value = (maxValue-minValue)*((mousePoint.X-bounds.X-float32(sliderWidth/2))/(bounds.Width-float32(sliderWidth))) + minValue
				}
			} else {
				guiControlExclusiveMode = false
				guiControlExclusiveRec = rl.Rectangle{}
			}
		} else if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if guiButtonDown() {
				state = STATE_PRESSED
				guiControlExclusiveMode = true
				guiControlExclusiveRec = bounds

				if !rl.CheckCollisionPointRec(mousePoint, slider) {
					*value = (maxValue-minValue)*((mousePoint.X-bounds.X-float32(sliderWidth/2))/(bounds.Width-float32(sliderWidth))) + minValue
				}
			} else {
				state = STATE_FOCUSED
			}
		}

		if *value > maxValue {
			*value = maxValue
		} else if *value < minValue {
			*value = minValue
		}
	}

	sliderValue := ((*value - minValue) / (maxValue - minValue)) * (bounds.Width - float32(sliderWidth) - float32(2*bw))
	if sliderWidth > 0 {
		slider.X += sliderValue
		slider.Width = float32(sliderWidth)
		if slider.X <= bounds.X+float32(bw) {
			slider.X = bounds.X + float32(bw)
		} else if slider.X+slider.Width >= bounds.X+bounds.Width {
			slider.X = bounds.X + bounds.Width - slider.Width - float32(bw)
		}
	} else if sliderWidth == 0 {
		slider.X += float32(bw)
		slider.Width = sliderValue
		if slider.Width > bounds.Width {
			slider.Width = bounds.Width - float32(2*bw)
		}
	}

	if prevValue != *value {
		result = RESULT_CHANGED
	}

	baseProp := BASE_COLOR_NORMAL
	if state == STATE_DISABLED {
		baseProp = BASE_COLOR_DISABLED
	}
	guiDrawRectangle(bounds, bw, getColor(getStyle(SLIDER, borderElement+PropertyID(state*3))), getColor(getStyle(SLIDER, baseProp)))

	switch state {
	case STATE_NORMAL:
		guiDrawRectangle(slider, 0, blank, getColor(getStyle(SLIDER, BASE_COLOR_PRESSED)))
	case STATE_FOCUSED:
		guiDrawRectangle(slider, 0, blank, getColor(getStyle(SLIDER, TEXT_COLOR_FOCUSED)))
	case STATE_PRESSED:
		guiDrawRectangle(slider, 0, blank, getColor(getStyle(SLIDER, TEXT_COLOR_PRESSED)))
	case STATE_DISABLED:
		guiDrawRectangle(slider, 0, blank, getColor(getStyle(SLIDER, TEXT_COLOR_DISABLED)))
	}

	leftRightText(SLIDER, bounds, textLeft, textRight, state)

	return result
}

// SliderBar - Slider Bar control
func SliderBar(bounds rl.Rectangle, textLeft, textRight string, value *float32, minValue, maxValue float32) int32 {
	preSliderWidth := getStyle(SLIDER, SLIDER_WIDTH)
	setStyle(SLIDER, SLIDER_WIDTH, 0)
	result := Slider(bounds, textLeft, textRight, value, minValue, maxValue)
	setStyle(SLIDER, SLIDER_WIDTH, preSliderWidth)

	return result
}

// ProgressBar - Progress Bar control
func ProgressBar(bounds rl.Rectangle, textLeft, textRight string, value *float32, minValue, maxValue float32) int32 {
	result := RESULT_NONE
	state := guiState

	temp := (maxValue - minValue) / 2.0
	if value == nil {
		value = &temp
	}
	prevValue := *value

	bw := getStyle(PROGRESSBAR, BORDER_WIDTH)
	fbw := float32(bw)
	progress := rl.Rectangle{X: bounds.X + fbw,
		Y:      bounds.Y + float32(bw+getStyle(PROGRESSBAR, PROGRESS_PADDING)),
		Height: bounds.Height - fbw - float32(2*getStyle(PROGRESSBAR, PROGRESS_PADDING)) - 1}

	if *value > maxValue {
		*value = maxValue
	}

	if state != STATE_DISABLED {
		progress.Width = (*value / (maxValue - minValue)) * (bounds.Width - float32(2*bw))
	}

	if prevValue != *value {
		result = RESULT_CHANGED
	}

	pw := float32(int32(progress.Width))

	if state == STATE_DISABLED {
		guiDrawRectangle(bounds, bw, getColor(getStyle(PROGRESSBAR, borderElement+PropertyID(state*3))), blank)
	} else {
		if *value > minValue {
			guiDrawRectangle(rl.Rectangle{X: bounds.X, Y: bounds.Y, Width: pw + fbw, Height: fbw}, 0, blank, getColor(getStyle(PROGRESSBAR, BORDER_COLOR_FOCUSED)))
			guiDrawRectangle(rl.Rectangle{X: bounds.X, Y: bounds.Y + 1, Width: fbw, Height: bounds.Height - 2}, 0, blank, getColor(getStyle(PROGRESSBAR, BORDER_COLOR_FOCUSED)))
			guiDrawRectangle(rl.Rectangle{X: bounds.X, Y: bounds.Y + bounds.Height - 1, Width: pw + fbw, Height: fbw}, 0, blank, getColor(getStyle(PROGRESSBAR, BORDER_COLOR_FOCUSED)))
		} else {
			guiDrawRectangle(rl.Rectangle{X: bounds.X, Y: bounds.Y, Width: fbw, Height: bounds.Height + fbw - 1}, 0, blank, getColor(getStyle(PROGRESSBAR, BORDER_COLOR_NORMAL)))
		}

		if *value >= maxValue {
			guiDrawRectangle(rl.Rectangle{X: bounds.X + progress.Width + fbw, Y: bounds.Y, Width: fbw, Height: bounds.Height + fbw - 1}, 0, blank, getColor(getStyle(PROGRESSBAR, BORDER_COLOR_FOCUSED)))
		} else {
			guiDrawRectangle(rl.Rectangle{X: bounds.X + pw + fbw, Y: bounds.Y, Width: bounds.Width - fbw - pw - 1, Height: fbw}, 0, blank, getColor(getStyle(PROGRESSBAR, BORDER_COLOR_NORMAL)))
			guiDrawRectangle(rl.Rectangle{X: bounds.X + pw + fbw, Y: bounds.Y + bounds.Height - 1, Width: bounds.Width - fbw - pw - 1, Height: fbw}, 0, blank, getColor(getStyle(PROGRESSBAR, BORDER_COLOR_NORMAL)))
			guiDrawRectangle(rl.Rectangle{X: bounds.X + bounds.Width - fbw, Y: bounds.Y, Width: fbw, Height: bounds.Height + fbw - 1}, 0, blank, getColor(getStyle(PROGRESSBAR, BORDER_COLOR_NORMAL)))
		}

		if getStyle(PROGRESSBAR, PROGRESS_SIDE) == 0 {
			guiDrawRectangle(progress, 0, blank, getColor(getStyle(PROGRESSBAR, BASE_COLOR_PRESSED)))
		} else {
			progress.X = bounds.X + bounds.Width - progress.Width - fbw
			guiDrawRectangle(progress, 0, blank, getColor(getStyle(PROGRESSBAR, BASE_COLOR_PRESSED)))
		}
	}

	leftRightText(PROGRESSBAR, bounds, textLeft, textRight, state)

	return result
}

// StatusBar - Status Bar control, shows info text
func StatusBar(bounds rl.Rectangle, text string) int32 {
	result := RESULT_NONE
	state := guiState

	if state != STATE_DISABLED && !guiLocked && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()

		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if guiButtonReleased() {
				result = RESULT_PRESSED
			}
		}
	}

	guiDrawRectangle(bounds, getStyle(STATUSBAR, BORDER_WIDTH), getColor(getStyle(STATUSBAR, borderElement+PropertyID(state*3))), getColor(getStyle(STATUSBAR, baseElement+PropertyID(state*3))))
	guiDrawText(text, getTextBounds(STATUSBAR, bounds), getStyle(STATUSBAR, TEXT_ALIGNMENT), getColor(getStyle(STATUSBAR, textElement+PropertyID(state*3))))

	return result
}

// DummyRec - Dummy control for placeholders
func DummyRec(bounds rl.Rectangle, text string) int32 {
	result := RESULT_NONE
	state := guiState

	if state != STATE_DISABLED && !guiLocked && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()

		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if guiButtonDown() {
				state = STATE_PRESSED
			} else {
				state = STATE_FOCUSED
			}

			if guiButtonReleased() {
				result = RESULT_PRESSED
			}
		}
	}

	baseProp, textProp := BASE_COLOR_NORMAL, TEXT_COLOR_NORMAL
	if state == STATE_DISABLED {
		baseProp, textProp = BASE_COLOR_DISABLED, TEXT_COLOR_DISABLED
	}
	guiDrawRectangle(bounds, 0, blank, getColor(getStyle(DEFAULT, baseProp)))
	guiDrawText(text, getTextBounds(DEFAULT, bounds), int32(TEXT_ALIGN_CENTER), getColor(getStyle(BUTTON, textProp)))

	return result
}

// ListView - List View control
func ListView(bounds rl.Rectangle, text string, scrollIndex *int32, active *int32) int32 {
	var items []string
	if text != "" {
		items = guiTextSplit(text, ';')
	}

	return ListViewEx(bounds, items, scrollIndex, active, nil)
}

// ListViewEx - List View with extended parameters
func ListViewEx(bounds rl.Rectangle, text []string, scrollIndex *int32, active *int32, focus *int32) int32 {
	result := RESULT_NONE
	state := guiState

	count := int32(len(text))

	itemFocused := int32(-1)
	if focus != nil {
		itemFocused = *focus
	}
	itemSelected := int32(-1)
	if active != nil {
		itemSelected = *active
	}
	var prevActive int32
	if active != nil {
		prevActive = *active
	}

	itemsHeight, itemsSpacing := getStyle(LISTVIEW, LIST_ITEMS_HEIGHT), getStyle(LISTVIEW, LIST_ITEMS_SPACING)

	useScrollBar := float32((itemsHeight+itemsSpacing)*count) > bounds.Height

	var itemBounds rl.Rectangle
	itemBounds.X = bounds.X + float32(itemsSpacing)
	itemBounds.Y = bounds.Y + float32(itemsSpacing+getStyle(DEFAULT, BORDER_WIDTH))
	itemBounds.Width = bounds.Width - float32(2*itemsSpacing) - float32(getStyle(DEFAULT, BORDER_WIDTH))
	itemBounds.Height = float32(itemsHeight)
	if useScrollBar {
		itemBounds.Width -= float32(getStyle(LISTVIEW, SCROLLBAR_WIDTH))
	}

	visibleItems := min(int32(bounds.Height)/(itemsHeight+itemsSpacing), count)

	var startIndex int32
	if scrollIndex != nil {
		startIndex = *scrollIndex
	}
	if startIndex < 0 || startIndex > count-visibleItems {
		startIndex = 0
	}
	endIndex := startIndex + visibleItems

	if state != STATE_DISABLED && !guiLocked && !guiControlExclusiveMode {
		mousePoint := guiPointerPosition()

		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			state = STATE_FOCUSED

			for i := range visibleItems {
				if rl.CheckCollisionPointRec(mousePoint, itemBounds) {
					itemFocused = startIndex + i
					if guiButtonPressed() {
						if itemSelected == startIndex+i {
							itemSelected = -1
						} else {
							itemSelected = startIndex + i
						}
					}
					break
				}

				itemBounds.Y += float32(itemsHeight + itemsSpacing)
			}

			if useScrollBar {
				scrollDelta := guiScrollDelta()
				startIndex -= int32(scrollDelta)

				if startIndex < 0 {
					startIndex = 0
				} else if startIndex > count-visibleItems {
					startIndex = count - visibleItems
				}

				endIndex = min(startIndex+visibleItems, count)
			}
		} else {
			itemFocused = -1
		}

		itemBounds.Y = bounds.Y + float32(itemsSpacing+getStyle(DEFAULT, BORDER_WIDTH))
	}

	guiDrawRectangle(bounds, getStyle(LISTVIEW, BORDER_WIDTH), getColor(getStyle(LISTVIEW, borderElement+PropertyID(state*3))), getColor(getStyle(DEFAULT, BACKGROUND_COLOR)))

	for i := int32(0); i < visibleItems && text != nil; i++ {
		itemText := item(text, startIndex+i)

		if getStyle(LISTVIEW, LIST_ITEMS_BORDER_NORMAL) != 0 {
			guiDrawRectangle(itemBounds, getStyle(LISTVIEW, LIST_ITEMS_BORDER_WIDTH), getColor(getStyle(LISTVIEW, BORDER_COLOR_NORMAL)), blank)
		}

		if state == STATE_DISABLED {
			if startIndex+i == itemSelected {
				guiDrawRectangle(itemBounds, getStyle(LISTVIEW, LIST_ITEMS_BORDER_WIDTH), getColor(getStyle(LISTVIEW, BORDER_COLOR_DISABLED)), getColor(getStyle(LISTVIEW, BASE_COLOR_DISABLED)))
			}

			guiDrawText(itemText, getTextBounds(LISTVIEW, itemBounds), getStyle(LISTVIEW, TEXT_ALIGNMENT), getColor(getStyle(LISTVIEW, TEXT_COLOR_DISABLED)))
		} else {
			if startIndex+i == itemSelected && active != nil {
				guiDrawRectangle(itemBounds, getStyle(LISTVIEW, LIST_ITEMS_BORDER_WIDTH), getColor(getStyle(LISTVIEW, BORDER_COLOR_PRESSED)), getColor(getStyle(LISTVIEW, BASE_COLOR_PRESSED)))
				guiDrawText(itemText, getTextBounds(LISTVIEW, itemBounds), getStyle(LISTVIEW, TEXT_ALIGNMENT), getColor(getStyle(LISTVIEW, TEXT_COLOR_PRESSED)))
			} else if startIndex+i == itemFocused {
				guiDrawRectangle(itemBounds, getStyle(LISTVIEW, LIST_ITEMS_BORDER_WIDTH), getColor(getStyle(LISTVIEW, BORDER_COLOR_FOCUSED)), getColor(getStyle(LISTVIEW, BASE_COLOR_FOCUSED)))
				guiDrawText(itemText, getTextBounds(LISTVIEW, itemBounds), getStyle(LISTVIEW, TEXT_ALIGNMENT), getColor(getStyle(LISTVIEW, TEXT_COLOR_FOCUSED)))
			} else {
				guiDrawText(itemText, getTextBounds(LISTVIEW, itemBounds), getStyle(LISTVIEW, TEXT_ALIGNMENT), getColor(getStyle(LISTVIEW, TEXT_COLOR_NORMAL)))
			}
		}

		itemBounds.Y += float32(itemsHeight + itemsSpacing)
	}

	if useScrollBar {
		scrollBarBounds := rl.Rectangle{
			X:      bounds.X + bounds.Width - float32(getStyle(LISTVIEW, BORDER_WIDTH)) - float32(getStyle(LISTVIEW, SCROLLBAR_WIDTH)),
			Y:      bounds.Y + float32(getStyle(LISTVIEW, BORDER_WIDTH)),
			Width:  float32(getStyle(LISTVIEW, SCROLLBAR_WIDTH)),
			Height: bounds.Height - float32(2*getStyle(DEFAULT, BORDER_WIDTH)),
		}

		percentVisible := float32(endIndex-startIndex) / float32(count)
		sliderSize := bounds.Height * percentVisible

		prevSliderSize := getStyle(SCROLLBAR, SCROLL_SLIDER_SIZE)
		prevScrollSpeed := getStyle(SCROLLBAR, SCROLL_SPEED)
		setStyle(SCROLLBAR, SCROLL_SLIDER_SIZE, int32(sliderSize))
		setStyle(SCROLLBAR, SCROLL_SPEED, count-visibleItems)

		startIndex = guiScrollBar(scrollBarBounds, startIndex, 0, count-visibleItems)

		setStyle(SCROLLBAR, SCROLL_SPEED, prevScrollSpeed)
		setStyle(SCROLLBAR, SCROLL_SLIDER_SIZE, prevSliderSize)
	}

	if active != nil {
		*active = itemSelected
	}
	if focus != nil {
		*focus = itemFocused
	}
	if scrollIndex != nil {
		*scrollIndex = startIndex
	}

	if active != nil && prevActive != *active {
		result = RESULT_CHANGED
	}

	return result
}

// TabBar - Tab Bar control
func TabBar(bounds rl.Rectangle, text string, hscroll *int32, active *int32) int32 {
	var items []string
	if text != "" {
		items = guiTextSplit(text, ';')
	}

	return TabBarEx(bounds, items, hscroll, active, nil)
}

// TabBarEx - Tab Bar control with extended parameters
func TabBarEx(bounds rl.Rectangle, text []string, hscroll *int32, active *int32, focus *int32) int32 {
	result := RESULT_NONE

	count := int32(len(text))

	var temp int32
	if active == nil {
		active = &temp
	}

	tabItemsWidth := getStyle(TABBAR, TAB_ITEMS_WIDTH)
	tabBounds := rl.Rectangle{X: bounds.X, Y: bounds.Y, Width: float32(tabItemsWidth), Height: bounds.Height}

	if *active < 0 {
		*active = 0
	} else if *active > count-1 {
		*active = count - 1
	}
	prevActive := *active

	offsetX := max(*active*tabItemsWidth-int32(rl.GetScreenWidth()), 0)

	toggle := false

	for i := range count {
		tabBounds.X = bounds.X + float32((tabItemsWidth+4)*i+offsetX)

		if tabBounds.X < float32(rl.GetScreenWidth()) {
			textAlignment := getStyle(TOGGLE, TEXT_ALIGNMENT)
			textPadding := getStyle(TOGGLE, TEXT_PADDING)
			setStyle(TOGGLE, TEXT_ALIGNMENT, int32(TEXT_ALIGN_LEFT))
			setStyle(TOGGLE, TEXT_PADDING, 8)

			if i == *active {
				toggle = true
				Toggle(tabBounds, text[i], &toggle)
			} else {
				toggle = false
				Toggle(tabBounds, text[i], &toggle)
				if toggle {
					*active = i
				}
			}

			if rl.CheckCollisionPointRec(guiPointerPosition(), tabBounds) && guiButtonPressedMid() {
				result = RESULT_TAB_CLOSE
			}

			setStyle(TOGGLE, TEXT_PADDING, textPadding)
			setStyle(TOGGLE, TEXT_ALIGNMENT, textAlignment)

			if getStyle(TABBAR, TAB_CLOSE_BUTTON) != 0 {
				tempBorderWidth := getStyle(BUTTON, BORDER_WIDTH)
				tempTextAlignment := getStyle(BUTTON, TEXT_ALIGNMENT)
				setStyle(BUTTON, BORDER_WIDTH, 1)
				setStyle(BUTTON, TEXT_ALIGNMENT, int32(TEXT_ALIGN_CENTER))

				if Button(rl.Rectangle{X: tabBounds.X + tabBounds.Width - 14 - 5, Y: tabBounds.Y + 5, Width: 14, Height: 14}, IconText(ICON_CROSS_SMALL, "")) != 0 {
					result = RESULT_TAB_CLOSE
				}

				setStyle(BUTTON, BORDER_WIDTH, tempBorderWidth)
				setStyle(BUTTON, TEXT_ALIGNMENT, tempTextAlignment)
			}
		}
	}

	var tabBorderOffset float32
	if getStyle(TABBAR, TAB_LINE_SIDE) != 1 {
		tabBorderOffset = bounds.Height - 1
	}
	guiDrawRectangle(rl.Rectangle{X: bounds.X, Y: bounds.Y + tabBorderOffset, Width: bounds.Width, Height: 1}, 0, blank, getColor(getStyle(TABBAR, BORDER_COLOR_NORMAL)))

	if result != RESULT_TAB_CLOSE && prevActive != *active {
		result = RESULT_CHANGED
	}

	return result
}

const colorBarAlphaCheckedSize = 10

// ColorPanel - Color Panel control
func ColorPanel(bounds rl.Rectangle, text string, color *rl.Color) int32 {
	var temp rl.Color
	if color == nil {
		color = &temp
	}

	vcolor := rl.Vector3{X: float32(color.R) / 255.0, Y: float32(color.G) / 255.0, Z: float32(color.B) / 255.0}
	hsv := convertRGBtoHSV(vcolor)
	prevHsv := hsv

	result := ColorPanelHSV(bounds, text, &hsv)

	if hsv.X != prevHsv.X || hsv.Y != prevHsv.Y || hsv.Z != prevHsv.Z {
		rgb := convertHSVtoRGB(hsv)

		*color = rl.Color{R: uint8(255.0 * rgb.X), G: uint8(255.0 * rgb.Y), B: uint8(255.0 * rgb.Z), A: color.A}
	}

	return result
}

// ColorBarAlpha - Color Bar Alpha control
func ColorBarAlpha(bounds rl.Rectangle, text string, alpha *float32) int32 {
	result := RESULT_NONE
	state := guiState

	var temp float32
	if alpha == nil {
		alpha = &temp
	}
	prevAlpha := *alpha

	selectorHeight, selectorOverflow := getStyle(COLORPICKER, HUEBAR_SELECTOR_HEIGHT), getStyle(COLORPICKER, HUEBAR_SELECTOR_OVERFLOW)
	selector := rl.Rectangle{X: bounds.X + (*alpha)*bounds.Width - float32(selectorHeight/2),
		Y:      bounds.Y - float32(selectorOverflow),
		Width:  float32(selectorHeight),
		Height: bounds.Height + float32(selectorOverflow*2)}

	if state != STATE_DISABLED && !guiLocked {
		mousePoint := guiPointerPosition()

		if guiControlExclusiveMode {
			if guiButtonDown() {
				if checkBoundsID(bounds, guiControlExclusiveRec) {
					state = STATE_PRESSED

					*alpha = (mousePoint.X - bounds.X) / bounds.Width
					if *alpha <= 0.0 {
						*alpha = 0.0
					}
					if *alpha >= 1.0 {
						*alpha = 1.0
					}
				}
			} else {
				guiControlExclusiveMode = false
				guiControlExclusiveRec = rl.Rectangle{}
			}
		} else if rl.CheckCollisionPointRec(mousePoint, bounds) || rl.CheckCollisionPointRec(mousePoint, selector) {
			if guiButtonDown() {
				state = STATE_PRESSED
				guiControlExclusiveMode = true
				guiControlExclusiveRec = bounds

				*alpha = (mousePoint.X - bounds.X) / bounds.Width
				if *alpha <= 0.0 {
					*alpha = 0.0
				}
				if *alpha >= 1.0 {
					*alpha = 1.0
				}
			} else {
				state = STATE_FOCUSED
			}
		}
	}

	if prevAlpha != *alpha {
		result = RESULT_CHANGED
	}

	if state != STATE_DISABLED {
		checksX := int32(bounds.Width) / colorBarAlphaCheckedSize
		checksY := int32(bounds.Height) / colorBarAlphaCheckedSize

		for x := range checksX {
			for y := range checksY {
				check := rl.Rectangle{X: bounds.X + float32(x*colorBarAlphaCheckedSize), Y: bounds.Y + float32(y*colorBarAlphaCheckedSize), Width: colorBarAlphaCheckedSize, Height: colorBarAlphaCheckedSize}
				c := fade(getColor(getStyle(COLORPICKER, BASE_COLOR_DISABLED)), 0.4)
				if (x+y)%2 != 0 {
					c = fade(getColor(getStyle(COLORPICKER, BORDER_COLOR_DISABLED)), 0.4)
				}
				guiDrawRectangle(check, 0, blank, c)
			}
		}

		rl.DrawRectangleGradientEx(bounds, rl.Color{R: 255, G: 255, B: 255, A: 0}, rl.Color{R: 255, G: 255, B: 255, A: 0}, fade(rl.Color{A: 255}, guiAlpha), fade(rl.Color{A: 255}, guiAlpha))
	} else {
		rl.DrawRectangleGradientEx(bounds, fade(getColor(getStyle(COLORPICKER, BASE_COLOR_DISABLED)), 0.1), fade(getColor(getStyle(COLORPICKER, BASE_COLOR_DISABLED)), 0.1), fade(getColor(getStyle(COLORPICKER, BORDER_COLOR_DISABLED)), guiAlpha), fade(getColor(getStyle(COLORPICKER, BORDER_COLOR_DISABLED)), guiAlpha))
	}

	guiDrawRectangle(bounds, getStyle(COLORPICKER, BORDER_WIDTH), getColor(getStyle(COLORPICKER, borderElement+PropertyID(state*3))), blank)

	guiDrawRectangle(selector, 0, blank, getColor(getStyle(COLORPICKER, borderElement+PropertyID(state*3))))

	return result
}

// ColorBarHue - Color Bar Hue control
func ColorBarHue(bounds rl.Rectangle, text string, value *float32) int32 {
	result := RESULT_NONE
	state := guiState

	var temp float32
	if value == nil {
		value = &temp
	}
	prevHue := *value

	selectorHeight, selectorOverflow := getStyle(COLORPICKER, HUEBAR_SELECTOR_HEIGHT), getStyle(COLORPICKER, HUEBAR_SELECTOR_OVERFLOW)
	selector := rl.Rectangle{X: bounds.X - float32(selectorOverflow), Y: bounds.Y + (*value)/360.0*bounds.Height - float32(selectorHeight/2), Width: bounds.Width + float32(selectorOverflow*2), Height: float32(selectorHeight)}

	if state != STATE_DISABLED && !guiLocked {
		mousePoint := guiPointerPosition()

		if guiControlExclusiveMode {
			if guiButtonDown() {
				if checkBoundsID(bounds, guiControlExclusiveRec) {
					state = STATE_PRESSED

					*value = (mousePoint.Y - bounds.Y) * 360 / bounds.Height
					if *value <= 0.0 {
						*value = 0.0
					}
					if *value >= 359.0 {
						*value = 359.0
					}
				}
			} else {
				guiControlExclusiveMode = false
				guiControlExclusiveRec = rl.Rectangle{}
			}
		} else if rl.CheckCollisionPointRec(mousePoint, bounds) || rl.CheckCollisionPointRec(mousePoint, selector) {
			if guiButtonDown() {
				state = STATE_PRESSED
				guiControlExclusiveMode = true
				guiControlExclusiveRec = bounds

				*value = (mousePoint.Y - bounds.Y) * 360 / bounds.Height
				if *value <= 0.0 {
					*value = 0.0
				}
				if *value >= 359.0 {
					*value = 359.0
				}
			} else {
				state = STATE_FOCUSED
			}
		}
	}

	if prevHue != *value {
		result = RESULT_CHANGED
	}

	if state != STATE_DISABLED {
		stops := []rl.Color{
			{R: 255, G: 0, B: 0, A: 255},
			{R: 255, G: 255, B: 0, A: 255},
			{R: 0, G: 255, B: 0, A: 255},
			{R: 0, G: 255, B: 255, A: 255},
			{R: 0, G: 0, B: 255, A: 255},
			{R: 255, G: 0, B: 255, A: 255},
			{R: 255, G: 0, B: 0, A: 255},
		}
		for i := range 6 {
			c1, c2 := fade(stops[i], guiAlpha), fade(stops[i+1], guiAlpha)
			rl.DrawRectangleGradientEx(rl.Rectangle{X: bounds.X, Y: bounds.Y + float32(i)*(bounds.Height/6.0), Width: bounds.Width, Height: bounds.Height / 6.0}, c1, c2, c2, c1)
		}
	} else {
		rl.DrawRectangleGradientEx(bounds,
			fade(fade(getColor(getStyle(COLORPICKER, BASE_COLOR_DISABLED)), 0.1), guiAlpha), fade(getColor(getStyle(COLORPICKER, BORDER_COLOR_DISABLED)), guiAlpha),
			fade(getColor(getStyle(COLORPICKER, BORDER_COLOR_DISABLED)), guiAlpha), fade(fade(getColor(getStyle(COLORPICKER, BASE_COLOR_DISABLED)), 0.1), guiAlpha))
	}

	guiDrawRectangle(bounds, getStyle(COLORPICKER, BORDER_WIDTH), getColor(getStyle(COLORPICKER, borderElement+PropertyID(state*3))), blank)

	guiDrawRectangle(selector, 0, blank, getColor(getStyle(COLORPICKER, borderElement+PropertyID(state*3))))

	return result
}

// ColorPicker - Color Picker control (multiple color controls)
func ColorPicker(bounds rl.Rectangle, text string, color *rl.Color) int32 {
	temp := rl.Color{R: 200, G: 0, B: 0, A: 255}
	if color == nil {
		color = &temp
	}

	ColorPanel(bounds, "", color)

	boundsHue := rl.Rectangle{X: bounds.X + bounds.Width + float32(getStyle(COLORPICKER, HUEBAR_PADDING)), Y: bounds.Y, Width: float32(getStyle(COLORPICKER, HUEBAR_WIDTH)), Height: bounds.Height}

	hsv := convertRGBtoHSV(rl.Vector3{X: float32(color.R) / 255.0, Y: float32(color.G) / 255.0, Z: float32(color.B) / 255.0})

	result := ColorBarHue(boundsHue, "", &hsv.X)

	rgb := convertHSVtoRGB(hsv)

	*color = rl.Color{R: uint8(math.Round(float64(rgb.X * 255.0))), G: uint8(math.Round(float64(rgb.Y * 255.0))), B: uint8(math.Round(float64(rgb.Z * 255.0))), A: color.A}

	return result
}

// ColorPickerHSV - Color Picker control that avoids conversion to RGB on each call (multiple color controls)
func ColorPickerHSV(bounds rl.Rectangle, text string, colorHsv *rl.Vector3) int32 {
	var tempHsv rl.Vector3

	if colorHsv == nil {
		tempColor := rl.Vector3{X: 200.0 / 255.0, Y: 0.0, Z: 0.0}
		tempHsv = convertRGBtoHSV(tempColor)
		colorHsv = &tempHsv
	}

	result := ColorPanelHSV(bounds, "", colorHsv)

	boundsHue := rl.Rectangle{X: bounds.X + bounds.Width + float32(getStyle(COLORPICKER, HUEBAR_PADDING)), Y: bounds.Y, Width: float32(getStyle(COLORPICKER, HUEBAR_WIDTH)), Height: bounds.Height}

	if result == RESULT_NONE {
		result = ColorBarHue(boundsHue, "", &colorHsv.X)
	}

	return result
}

// ColorPanelHSV - Color Panel control that updates Hue-Saturation-Value color value, used by ColorPickerHSV()
func ColorPanelHSV(bounds rl.Rectangle, text string, colorHsv *rl.Vector3) int32 {
	result := RESULT_NONE
	state := guiState

	var temp rl.Vector3
	if colorHsv == nil {
		colorHsv = &temp
	}
	prevColorHsv := *colorHsv

	var pickerSelector rl.Vector2

	colWhite := rl.Color{R: 255, G: 255, B: 255, A: 255}
	colBlack := rl.Color{R: 0, G: 0, B: 0, A: 255}

	pickerSelector.X = bounds.X + colorHsv.Y*bounds.Width
	pickerSelector.Y = bounds.Y + (1.0-colorHsv.Z)*bounds.Height

	maxHue := rl.Vector3{X: colorHsv.X, Y: 1.0, Z: 1.0}
	rgbHue := convertHSVtoRGB(maxHue)
	maxHueCol := rl.Color{R: uint8(255.0 * rgbHue.X), G: uint8(255.0 * rgbHue.Y), B: uint8(255.0 * rgbHue.Z), A: 255}

	if state != STATE_DISABLED && !guiLocked {
		mousePoint := guiPointerPosition()

		if guiControlExclusiveMode {
			if guiButtonDown() {
				if checkBoundsID(bounds, guiControlExclusiveRec) {
					pickerSelector = mousePoint

					if pickerSelector.X < bounds.X {
						pickerSelector.X = bounds.X
					}
					if pickerSelector.X > bounds.X+bounds.Width {
						pickerSelector.X = bounds.X + bounds.Width
					}
					if pickerSelector.Y < bounds.Y {
						pickerSelector.Y = bounds.Y
					}
					if pickerSelector.Y > bounds.Y+bounds.Height {
						pickerSelector.Y = bounds.Y + bounds.Height
					}

					colorPick := rl.Vector2{X: pickerSelector.X - bounds.X, Y: pickerSelector.Y - bounds.Y}

					colorPick.X /= bounds.Width
					colorPick.Y /= bounds.Height

					colorHsv.Y = colorPick.X
					colorHsv.Z = 1.0 - colorPick.Y
				}
			} else {
				guiControlExclusiveMode = false
				guiControlExclusiveRec = rl.Rectangle{}
			}
		} else if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if guiButtonDown() {
				state = STATE_PRESSED
				guiControlExclusiveMode = true
				guiControlExclusiveRec = bounds
				pickerSelector = mousePoint

				colorPick := rl.Vector2{X: pickerSelector.X - bounds.X, Y: pickerSelector.Y - bounds.Y}

				colorPick.X /= bounds.Width
				colorPick.Y /= bounds.Height

				colorHsv.Y = colorPick.X
				colorHsv.Z = 1.0 - colorPick.Y
			} else {
				state = STATE_FOCUSED
			}
		}
	}

	if prevColorHsv.X != colorHsv.X || prevColorHsv.Y != colorHsv.Y || prevColorHsv.Z != colorHsv.Z {
		result = RESULT_CHANGED
	}

	if state != STATE_DISABLED {
		rl.DrawRectangleGradientEx(bounds, fade(colWhite, guiAlpha), fade(colWhite, guiAlpha), fade(maxHueCol, guiAlpha), fade(maxHueCol, guiAlpha))
		rl.DrawRectangleGradientEx(bounds, fade(colBlack, 0), fade(colBlack, guiAlpha), fade(colBlack, guiAlpha), fade(colBlack, 0))

		selectorSize := getStyle(COLORPICKER, COLOR_SELECTOR_SIZE)
		selector := rl.Rectangle{X: pickerSelector.X - float32(selectorSize/2), Y: pickerSelector.Y - float32(selectorSize/2), Width: float32(selectorSize), Height: float32(selectorSize)}
		guiDrawRectangle(selector, 0, blank, colWhite)
	} else {
		rl.DrawRectangleGradientEx(bounds, fade(fade(getColor(getStyle(COLORPICKER, BASE_COLOR_DISABLED)), 0.1), guiAlpha), fade(fade(colBlack, 0.6), guiAlpha), fade(fade(colBlack, 0.6), guiAlpha), fade(fade(getColor(getStyle(COLORPICKER, BORDER_COLOR_DISABLED)), 0.6), guiAlpha))
	}

	guiDrawRectangle(bounds, getStyle(COLORPICKER, BORDER_WIDTH), getColor(getStyle(COLORPICKER, borderElement+PropertyID(state*3))), blank)

	return result
}

const (
	messageBoxButtonHeight    = 24
	messageBoxButtonPadding   = 12
	textInputBoxButtonHeight  = 24
	textInputBoxButtonPadding = 12
	textInputBoxHeight        = 26
	gridAlpha                 = 0.15
)

var textInputBoxEditMode bool

// MessageBox - Message Box control, displays a message
func MessageBox(bounds rl.Rectangle, title, message, btnText string, btnActive *int32) int32 {
	result := RESULT_NONE

	var temp int32
	if btnActive == nil {
		btnActive = &temp
	}

	btnTextList := guiTextSplit(btnText, ';')
	buttonCount := int32(len(btnTextList))

	var buttonBounds rl.Rectangle
	buttonBounds.X = bounds.X + messageBoxButtonPadding
	buttonBounds.Y = bounds.Y + bounds.Height - messageBoxButtonHeight - messageBoxButtonPadding
	buttonBounds.Width = (bounds.Width - float32(messageBoxButtonPadding*(buttonCount+1))) / float32(buttonCount)
	buttonBounds.Height = messageBoxButtonHeight

	var textBounds rl.Rectangle
	textBounds.X = bounds.X + messageBoxButtonPadding
	textBounds.Y = bounds.Y + windowBoxStatusBarHeight + messageBoxButtonPadding
	textBounds.Width = bounds.Width - messageBoxButtonPadding*2
	textBounds.Height = bounds.Height - windowBoxStatusBarHeight - 3*messageBoxButtonPadding - messageBoxButtonHeight

	if WindowBox(bounds, title) == RESULT_PRESSED {
		*btnActive = 0
		result = RESULT_PRESSED
	}

	prevTextAlignment := getStyle(LABEL, TEXT_ALIGNMENT)
	setStyle(LABEL, TEXT_ALIGNMENT, int32(TEXT_ALIGN_CENTER))
	Label(textBounds, message)
	setStyle(LABEL, TEXT_ALIGNMENT, prevTextAlignment)

	prevTextAlignment = getStyle(BUTTON, TEXT_ALIGNMENT)
	setStyle(BUTTON, TEXT_ALIGNMENT, int32(TEXT_ALIGN_CENTER))

	for i := range buttonCount {
		if Button(buttonBounds, btnTextList[i]) != 0 {
			*btnActive = i + 1
			result = RESULT_PRESSED
		}

		buttonBounds.X += buttonBounds.Width + messageBoxButtonPadding
	}

	setStyle(BUTTON, TEXT_ALIGNMENT, prevTextAlignment)

	return result
}

// TextInputBox - Text Input Box control, ask for text, supports secret
func TextInputBox(bounds rl.Rectangle, title, message string, text *string, textSize int32, btnText string, btnActive *int32, secretViewActive *bool) int32 {
	result := RESULT_NONE

	var temp int32
	if btnActive == nil {
		btnActive = &temp
	}

	btnTextList := guiTextSplit(btnText, ';')
	buttonCount := int32(len(btnTextList))

	var buttonBounds rl.Rectangle
	buttonBounds.X = bounds.X + textInputBoxButtonPadding
	buttonBounds.Y = bounds.Y + bounds.Height - textInputBoxButtonHeight - textInputBoxButtonPadding
	buttonBounds.Width = (bounds.Width - float32(textInputBoxButtonPadding*(buttonCount+1))) / float32(buttonCount)
	buttonBounds.Height = textInputBoxButtonHeight

	messageInputHeight := int32(bounds.Height) - windowBoxStatusBarHeight - getStyle(STATUSBAR, BORDER_WIDTH) - textInputBoxButtonHeight - 2*textInputBoxButtonPadding

	var textBounds rl.Rectangle
	if message != "" {
		messageTextSize := GetTextWidth(message) + 2

		textBounds.X = bounds.X + bounds.Width/2 - float32(messageTextSize/2)
		textBounds.Y = bounds.Y + windowBoxStatusBarHeight + float32(messageInputHeight/4) - float32(getStyle(DEFAULT, TEXT_SIZE))/2
		textBounds.Width = float32(messageTextSize)
		textBounds.Height = float32(getStyle(DEFAULT, TEXT_SIZE))
	}

	var textBoxBounds rl.Rectangle
	textBoxBounds.X = bounds.X + textInputBoxButtonPadding
	textBoxBounds.Y = bounds.Y + windowBoxStatusBarHeight - textInputBoxHeight/2
	if message == "" {
		textBoxBounds.Y = bounds.Y + 24 + textInputBoxButtonPadding
	} else {
		textBoxBounds.Y += float32(messageInputHeight/2 + messageInputHeight/4)
	}
	textBoxBounds.Width = bounds.Width - textInputBoxButtonPadding*2
	textBoxBounds.Height = textInputBoxHeight

	if WindowBox(bounds, title) == RESULT_PRESSED {
		*btnActive = 0
		result = RESULT_PRESSED
	}

	if message != "" {
		prevTextAlignment := getStyle(LABEL, TEXT_ALIGNMENT)
		setStyle(LABEL, TEXT_ALIGNMENT, int32(TEXT_ALIGN_CENTER))
		Label(textBounds, message)
		setStyle(LABEL, TEXT_ALIGNMENT, prevTextAlignment)
	}

	prevTextBoxAlignment := getStyle(TEXTBOX, TEXT_ALIGNMENT)
	setStyle(TEXTBOX, TEXT_ALIGNMENT, int32(TEXT_ALIGN_LEFT))

	if secretViewActive != nil {
		stars := "****************"
		boxText := &stars
		if *secretViewActive || textInputBoxEditMode {
			boxText = text
		}
		if TextBox(rl.Rectangle{X: textBoxBounds.X, Y: textBoxBounds.Y, Width: textBoxBounds.Width - 4 - textInputBoxHeight, Height: textBoxBounds.Height}, boxText, textSize, textInputBoxEditMode) == RESULT_PRESSED {
			textInputBoxEditMode = !textInputBoxEditMode
		}

		eye := IconText(ICON_EYE_OFF, "")
		if *secretViewActive {
			eye = IconText(ICON_EYE_ON, "")
		}
		Toggle(rl.Rectangle{X: textBoxBounds.X + textBoxBounds.Width - textInputBoxHeight, Y: textBoxBounds.Y, Width: textInputBoxHeight, Height: textInputBoxHeight}, eye, secretViewActive)
	} else {
		if TextBox(textBoxBounds, text, textSize, textInputBoxEditMode) == RESULT_PRESSED {
			textInputBoxEditMode = !textInputBoxEditMode
		}
	}

	setStyle(TEXTBOX, TEXT_ALIGNMENT, prevTextBoxAlignment)

	prevBtnTextAlignment := getStyle(BUTTON, TEXT_ALIGNMENT)
	setStyle(BUTTON, TEXT_ALIGNMENT, int32(TEXT_ALIGN_CENTER))

	for i := range buttonCount {
		if Button(buttonBounds, btnTextList[i]) != 0 {
			*btnActive = i + 1
			result = RESULT_PRESSED
		}

		buttonBounds.X += buttonBounds.Width + messageBoxButtonPadding
	}

	if result == RESULT_PRESSED {
		textInputBoxEditMode = false
	}

	setStyle(BUTTON, TEXT_ALIGNMENT, prevBtnTextAlignment)

	return result
}

// Grid - Grid control
func Grid(bounds rl.Rectangle, text string, spacing float32, subdivs int32, mouseCell *rl.Vector2) int32 {
	result := RESULT_NONE
	state := guiState

	mousePoint := guiPointerPosition()
	currentMouseCell := rl.Vector2{X: -1, Y: -1}

	spaceWidth := spacing / float32(subdivs)
	linesV := int32(bounds.Width/spaceWidth) + 1
	linesH := int32(bounds.Height/spaceWidth) + 1

	color := getStyle(DEFAULT, LINE_COLOR)

	if state != STATE_DISABLED && !guiLocked && !guiControlExclusiveMode {
		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			currentMouseCell.X = float32(math.Floor(float64((mousePoint.X - bounds.X) / spacing)))
			currentMouseCell.Y = float32(math.Floor(float64((mousePoint.Y - bounds.Y) / spacing)))
			result = RESULT_PRESSED
		}
	}

	if state == STATE_DISABLED {
		color = getStyle(DEFAULT, BORDER_COLOR_DISABLED)
	}

	if subdivs > 0 {
		gridLineColor := func(i int32) rl.Color {
			if i%subdivs == 0 {
				return guiFade(getColor(color), gridAlpha*4)
			}
			return guiFade(getColor(color), gridAlpha)
		}

		for i := range linesV {
			lineV := rl.Rectangle{X: bounds.X + spacing*float32(i)/float32(subdivs), Y: bounds.Y, Width: 1, Height: bounds.Height + 1}
			guiDrawRectangle(lineV, 0, blank, gridLineColor(i))
		}

		for i := range linesH {
			lineH := rl.Rectangle{X: bounds.X, Y: bounds.Y + spacing*float32(i)/float32(subdivs), Width: bounds.Width + 1, Height: 1}
			guiDrawRectangle(lineH, 0, blank, gridLineColor(i))
		}
	}

	if mouseCell != nil {
		*mouseCell = currentMouseCell
	}

	return result
}

// EnableTooltip - Enable gui tooltips (global state)
func EnableTooltip() {
	guiTooltip = true
}

// DisableTooltip - Disable gui tooltips (global state)
func DisableTooltip() {
	guiTooltip = false
}

// SetTooltip - Set tooltip string
func SetTooltip(tooltip string) {
	guiTooltipPtr = tooltip
}
