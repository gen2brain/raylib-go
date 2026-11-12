package raygui

import (
	"strings"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	iconTextPadding            = 4
	maxTextLines               = 128
	textSplitMaxItems          = 128
	textSplitMaxTextSize       = 1024
	iconMaxFontBacked          = 257
	iconFontAtlasPadding       = 1
	windowBoxStatusBarHeight   = 24
	windowBoxCloseButtonHeight = 18
	maxLineBufferSize          = 256
)

var blank = rl.Color{}

func at(s string, i int) byte {
	if i < 0 || i >= len(s) {
		return 0
	}
	return s[i]
}

func tail(s string, i int) string {
	if i < 0 {
		i = 0
	}
	if i >= len(s) {
		return ""
	}
	return s[i:]
}

func cstr(s string) string {
	if before, _, ok := strings.Cut(s, "\x00"); ok {
		return before
	}
	return s
}

func fontGlyphs(font rl.Font) []rl.GlyphInfo {
	if font.Chars == nil || font.CharsCount <= 0 {
		return nil
	}
	return unsafe.Slice(font.Chars, font.CharsCount)
}

func fontRecs(font rl.Font) []rl.Rectangle {
	if font.Recs == nil || font.CharsCount <= 0 {
		return nil
	}
	return unsafe.Slice(font.Recs, font.CharsCount)
}

func glyphWidthAt(index int32, scaleFactor float32) float32 {
	glyphs, recs := fontGlyphs(guiFont), fontRecs(guiFont)
	if glyphs[index].AdvanceX == 0 {
		return recs[index].Width * scaleFactor
	}
	return float32(glyphs[index].AdvanceX) * scaleFactor
}

func getColor(hexValue int32) rl.Color {
	v := uint32(hexValue)
	return rl.Color{R: uint8(v >> 24), G: uint8(v >> 16), B: uint8(v >> 8), A: uint8(v)}
}

func fade(color rl.Color, alpha float32) rl.Color {
	if alpha < 0.0 {
		alpha = 0.0
	} else if alpha > 1.0 {
		alpha = 1.0
	}
	return rl.Color{R: color.R, G: color.G, B: color.B, A: uint8(255.0 * alpha)}
}

func getCodepoint(text string, codepointSize *int32) int32 {
	codepoint := int32(0x3f)
	*codepointSize = 1
	if len(text) == 0 {
		return 0
	}

	octet := at(text, 0)

	if octet <= 0x7f {
		codepoint = int32(octet)
	} else if octet&0xe0 == 0xc0 {
		octet1 := at(text, 1)

		if octet1 == 0 || octet1>>6 != 2 {
			*codepointSize = 2
			return codepoint
		}

		if octet >= 0xc2 && octet <= 0xdf {
			codepoint = int32(octet&0x1f)<<6 | int32(octet1&0x3f)
			*codepointSize = 2
		}
	} else if octet&0xf0 == 0xe0 {
		octet1 := at(text, 1)

		if octet1 == 0 || octet1>>6 != 2 {
			*codepointSize = 2
			return codepoint
		}

		octet2 := at(text, 2)

		if octet2 == 0 || octet2>>6 != 2 {
			*codepointSize = 3
			return codepoint
		}

		if (octet == 0xe0 && !(octet1 >= 0xa0 && octet1 <= 0xbf)) || (octet == 0xed && !(octet1 >= 0x80 && octet1 <= 0x9f)) {
			*codepointSize = 2
			return codepoint
		}

		if octet >= 0xe0 && octet <= 0xef {
			codepoint = int32(octet&0xf)<<12 | int32(octet1&0x3f)<<6 | int32(octet2&0x3f)
			*codepointSize = 3
		}
	} else if octet&0xf8 == 0xf0 {
		if octet > 0xf4 {
			return codepoint
		}

		octet1 := at(text, 1)

		if octet1 == 0 || octet1>>6 != 2 {
			*codepointSize = 2
			return codepoint
		}

		octet2 := at(text, 2)

		if octet2 == 0 || octet2>>6 != 2 {
			*codepointSize = 3
			return codepoint
		}

		octet3 := at(text, 3)

		if octet3 == 0 || octet3>>6 != 2 {
			*codepointSize = 4
			return codepoint
		}

		if (octet == 0xf0 && !(octet1 >= 0x90 && octet1 <= 0xbf)) || (octet == 0xf4 && !(octet1 >= 0x80 && octet1 <= 0x8f)) {
			*codepointSize = 2
			return codepoint
		}

		if octet >= 0xf0 {
			codepoint = int32(octet&0x7)<<18 | int32(octet1&0x3f)<<12 | int32(octet2&0x3f)<<6 | int32(octet3&0x3f)
			*codepointSize = 4
		}
	}

	if codepoint > 0x10ffff {
		codepoint = 0x3f
	}

	return codepoint
}

func getCodepointNext(text string, codepointSize *int32) int32 {
	codepoint := int32(0x3f)
	*codepointSize = 1

	b0, b1, b2, b3 := at(text, 0), at(text, 1), at(text, 2), at(text, 3)

	if 0xf8&b0 == 0xf0 {
		if (b1&0xc0)^0x80 != 0 || (b2&0xc0)^0x80 != 0 || (b3&0xc0)^0x80 != 0 {
			return codepoint
		}
		codepoint = int32(0x07&b0)<<18 | int32(0x3f&b1)<<12 | int32(0x3f&b2)<<6 | int32(0x3f&b3)
		*codepointSize = 4
	} else if 0xf0&b0 == 0xe0 {
		if (b1&0xc0)^0x80 != 0 || (b2&0xc0)^0x80 != 0 {
			return codepoint
		}
		codepoint = int32(0x0f&b0)<<12 | int32(0x3f&b1)<<6 | int32(0x3f&b2)
		*codepointSize = 3
	} else if 0xe0&b0 == 0xc0 {
		if (b1&0xc0)^0x80 != 0 {
			return codepoint
		}
		codepoint = int32(0x1f&b0)<<6 | int32(0x3f&b1)
		*codepointSize = 2
	} else if 0x80&b0 == 0 {
		codepoint = int32(b0)
		*codepointSize = 1
	}

	return codepoint
}

func getCodepointPrevious(text string, pos int, codepointSize *int32) int32 {
	codepoint := int32(0x3f)
	*codepointSize = 1

	p := pos
	for {
		p--
		if !(0x80&at(text, p) != 0 && 0xc0&at(text, p) == 0x80) {
			break
		}
	}

	var cpSize int32
	codepoint = getCodepointNext(tail(text, p), &cpSize)
	if codepoint != 0 {
		*codepointSize = cpSize
	}

	return codepoint
}

func codepointToUTF8(codepoint int32) string {
	var utf8 [4]byte
	size := 0

	if codepoint <= 0x7f {
		utf8[0] = byte(codepoint)
		size = 1
	} else if codepoint <= 0x7ff {
		utf8[0] = byte(((codepoint >> 6) & 0x1f) | 0xc0)
		utf8[1] = byte((codepoint & 0x3f) | 0x80)
		size = 2
	} else if codepoint <= 0xffff {
		utf8[0] = byte(((codepoint >> 12) & 0x0f) | 0xe0)
		utf8[1] = byte(((codepoint >> 6) & 0x3f) | 0x80)
		utf8[2] = byte((codepoint & 0x3f) | 0x80)
		size = 3
	} else if codepoint <= 0x10ffff {
		utf8[0] = byte(((codepoint >> 18) & 0x07) | 0xf0)
		utf8[1] = byte(((codepoint >> 12) & 0x3f) | 0x80)
		utf8[2] = byte(((codepoint >> 6) & 0x3f) | 0x80)
		utf8[3] = byte((codepoint & 0x3f) | 0x80)
		size = 4
	}

	return string(utf8[:size])
}

func textToInteger(text string) int32 {
	var value int32
	var sign int32 = 1

	i := 0
	if at(text, 0) == '+' || at(text, 0) == '-' {
		if at(text, 0) == '-' {
			sign = -1
		}
		i++
	}

	for ; at(text, i) >= '0' && at(text, i) <= '9'; i++ {
		value = value*10 + int32(at(text, i)-'0')
	}

	return value * sign
}

func textToFloat(text string) float32 {
	var value float32
	var sign float32 = 1.0

	i := 0
	if at(text, 0) == '+' || at(text, 0) == '-' {
		if at(text, 0) == '-' {
			sign = -1.0
		}
		i++
	}

	for ; at(text, i) >= '0' && at(text, i) <= '9'; i++ {
		value = value*10.0 + float32(at(text, i)-'0')
	}

	c := at(text, i)
	i++
	if c == '.' {
		var divisor float32 = 10.0
		for ; at(text, i) >= '0' && at(text, i) <= '9'; i++ {
			value += float32(at(text, i)-'0') / divisor
			divisor = divisor * 10.0
		}
	}

	return value * sign
}

func getLineWidth(text string) int32 {
	var textSize rl.Vector2
	textIconOffset := 0

	if at(text, 0) != 0 {
		if at(text, 0) == '#' {
			pos := 1
			for pos < 4 && at(text, pos) >= '0' && at(text, pos) <= '9' {
				pos++
			}
			if at(text, pos) == '#' {
				textIconOffset = pos + 1
			}
		}

		text = tail(text, textIconOffset)

		GetStyle(DEFAULT, TEXT_SIZE)

		textSize = guiMeasureText(text)

		if textIconOffset > 0 {
			textSize.X += RAYGUI_ICON_SIZE + iconTextPadding
		}
	}

	return int32(textSize.X)
}

func getTextBounds(control ControlID, bounds rl.Rectangle) rl.Rectangle {
	textBounds := bounds

	textBounds.X = bounds.X + float32(getStyle(control, BORDER_WIDTH))
	textBounds.Y = bounds.Y + float32(getStyle(control, BORDER_WIDTH)+getStyle(control, TEXT_PADDING))
	textBounds.Width = bounds.Width - float32(2*getStyle(control, BORDER_WIDTH)) - float32(2*getStyle(control, TEXT_PADDING))
	textBounds.Height = bounds.Height - float32(2*getStyle(control, BORDER_WIDTH)) - float32(2*getStyle(control, TEXT_PADDING))

	if GetStyle(control, TEXT_ALIGNMENT) == TEXT_ALIGN_RIGHT {
		textBounds.X -= float32(getStyle(control, TEXT_PADDING))
	} else {
		textBounds.X += float32(getStyle(control, TEXT_PADDING))
	}

	return textBounds
}

func getTextIcon(text string, iconId *int32) string {
	*iconId = -1
	if at(text, 0) == '#' {
		var iconValue []byte

		pos := 1
		for pos < 4 && at(text, pos) >= '0' && at(text, pos) <= '9' {
			iconValue = append(iconValue, at(text, pos))
			pos++
		}

		if at(text, pos) == '#' {
			rawIconId := textToInteger(string(iconValue))
			if rawIconId < RAYGUI_ICON_MAX_ICONS {
				*iconId = rawIconId

				if *iconId >= 0 {
					text = tail(text, pos+1)
				}
			}
		}
	}

	return text
}

func getTextLines(text string) []string {
	textLength := len(text)

	lines := []string{text}

	for i := 0; i < textLength && len(lines) < maxTextLines; i++ {
		if text[i] == '\n' && i+1 < textLength {
			lines = append(lines, text[i+1:])
		}
	}

	return lines
}

func getNextSpaceWidth(text string, nextSpaceIndex *int32) float32 {
	var width float32
	var codepointByteCount int32
	scaleFactor := float32(getStyle(DEFAULT, TEXT_SIZE)) / float32(guiFont.BaseSize)

	for i := 0; at(text, i) != 0; i++ {
		if at(text, i) != ' ' {
			codepoint := getCodepoint(tail(text, i), &codepointByteCount)
			index := rl.GetGlyphIndex(guiFont, codepoint)
			glyphWidth := glyphWidthAt(index, scaleFactor)
			width += glyphWidth + float32(getStyle(DEFAULT, TEXT_SPACING))
		} else {
			*nextSpaceIndex = int32(i)
			break
		}
	}

	return width
}

func textValignPixelOffset(h float32) int32 {
	return int32(h) % 2
}

func guiDrawText(text string, textBounds rl.Rectangle, alignment int32, tint rl.Color) {
	if at(text, 0) == 0 {
		return
	}

	lines := getTextLines(text)
	lineCount := int32(len(lines))

	alignmentVertical := getStyle(DEFAULT, TEXT_ALIGNMENT_VERTICAL)
	wrapMode := getStyle(DEFAULT, TEXT_WRAP_MODE)

	totalHeight := float32(lineCount*getStyle(DEFAULT, TEXT_SIZE) + (lineCount-1)*getStyle(DEFAULT, TEXT_LINE_SPACING))
	var posOffsetY float32

	for i := range lines {
		var iconId int32
		lines[i] = getTextIcon(lines[i], &iconId)

		textBoundsPosition := rl.Vector2{X: textBounds.X, Y: textBounds.Y}
		var textBoundsWidthOffset float32

		textSizeX := getLineWidth(lines[i])

		if iconId >= 0 {
			textSizeX += RAYGUI_ICON_SIZE * guiIconScale

			if at(lines[i], 0) != 0 {
				textSizeX += iconTextPadding
			}
		}

		switch PropertyValue(alignment) {
		case TEXT_ALIGN_LEFT:
			textBoundsPosition.X = textBounds.X
		case TEXT_ALIGN_CENTER:
			textBoundsPosition.X = textBounds.X + textBounds.Width/2 - float32(textSizeX/2)
		case TEXT_ALIGN_RIGHT:
			textBoundsPosition.X = textBounds.X + textBounds.Width - float32(textSizeX)
		}

		if float32(textSizeX) > textBounds.Width && at(lines[i], 0) != 0 {
			textBoundsPosition.X = textBounds.X
		}

		switch PropertyValue(alignmentVertical) {
		case TEXT_ALIGN_TOP:
			textBoundsPosition.Y = textBounds.Y + posOffsetY
		case TEXT_ALIGN_MIDDLE:
			textBoundsPosition.Y = textBounds.Y + posOffsetY + textBounds.Height/2 - totalHeight/2 + float32(textValignPixelOffset(textBounds.Height))
		case TEXT_ALIGN_BOTTOM:
			textBoundsPosition.Y = textBounds.Y + posOffsetY + textBounds.Height - totalHeight + float32(textValignPixelOffset(textBounds.Height))
		}

		textBoundsPosition.X = float32(int32(textBoundsPosition.X))
		textBoundsPosition.Y = float32(int32(textBoundsPosition.Y))

		if iconId >= 0 {
			DrawIcon(IconID(iconId), int32(textBoundsPosition.X), int32(textBounds.Y+textBounds.Height/2-float32(uint32(RAYGUI_ICON_SIZE*guiIconScale)/2)+float32(textValignPixelOffset(textBounds.Height))), guiIconScale, tint)
			textBoundsPosition.X += float32(RAYGUI_ICON_SIZE*guiIconScale + iconTextPadding)
			textBoundsWidthOffset = float32(RAYGUI_ICON_SIZE*guiIconScale + iconTextPadding)
		}

		line := lines[i]
		lineSize := 0
		for c := 0; at(line, c) != 0 && at(line, c) != '\n' && at(line, c) != '\r'; c++ {
			lineSize++
		}
		scaleFactor := float32(getStyle(DEFAULT, TEXT_SIZE)) / float32(guiFont.BaseSize)

		lastSpaceIndex := 0
		tempWrapCharMode := false

		var textOffsetY int32
		var textOffsetX float32
		var glyphWidth float32

		ellipsisWidth := GetTextWidth("...")
		textOverflow := false
		var codepointSize int32
		for c := 0; c < lineSize; c += int(codepointSize) {
			codepoint := getCodepointNext(tail(line, c), &codepointSize)
			index := rl.GetGlyphIndex(guiFont, codepoint)

			if codepoint == 0x3f {
				codepointSize = 1
			}

			glyphWidth = glyphWidthAt(index, scaleFactor)

			if PropertyValue(wrapMode) == TEXT_WRAP_CHAR {
				if textOffsetX+glyphWidth > textBounds.Width-textBoundsWidthOffset {
					textOffsetX = 0.0
					textOffsetY += getStyle(DEFAULT, TEXT_SIZE) + getStyle(DEFAULT, TEXT_LINE_SPACING)

					if tempWrapCharMode {
						wrapMode = int32(TEXT_WRAP_WORD)
						tempWrapCharMode = false
					}
				}
			} else if PropertyValue(wrapMode) == TEXT_WRAP_WORD {
				if codepoint == 32 {
					lastSpaceIndex = c
				}

				var nextSpaceIndex int32
				nextSpaceWidth := getNextSpaceWidth(tail(line, c), &nextSpaceIndex)

				var nextSpaceIndex2 int32
				nextWordSize := getNextSpaceWidth(tail(line, lastSpaceIndex+1), &nextSpaceIndex2)

				if nextWordSize > textBounds.Width-textBoundsWidthOffset {
					tempWrapCharMode = true
					wrapMode = int32(TEXT_WRAP_CHAR)
				} else if textOffsetX+nextSpaceWidth > textBounds.Width-textBoundsWidthOffset {
					textOffsetX = 0.0
					textOffsetY += getStyle(DEFAULT, TEXT_SIZE) + getStyle(DEFAULT, TEXT_LINE_SPACING)
				}
			}

			if codepoint == '\n' {
				break
			}

			if codepoint != ' ' && codepoint != '\t' {
				pos := rl.Vector2{X: textBoundsPosition.X + textOffsetX, Y: textBoundsPosition.Y + float32(textOffsetY)}
				textSize := float32(getStyle(DEFAULT, TEXT_SIZE))
				if PropertyValue(wrapMode) == TEXT_WRAP_NONE {
					if float32(textSizeX) > textBounds.Width {
						if textOffsetX <= textBounds.Width-glyphWidth-textBoundsWidthOffset-float32(ellipsisWidth) {
							rl.DrawTextCodepoint(guiFont, codepoint, pos, textSize, guiFade(tint, guiAlpha))
						} else if !textOverflow {
							textOverflow = true

							for j := int32(0); j < ellipsisWidth; j += ellipsisWidth / 3 {
								rl.DrawTextCodepoint(guiFont, '.', rl.Vector2{X: textBoundsPosition.X + textOffsetX + float32(j), Y: textBoundsPosition.Y + float32(textOffsetY)}, textSize, guiFade(tint, guiAlpha))
							}
						}
					} else {
						rl.DrawTextCodepoint(guiFont, codepoint, pos, textSize, guiFade(tint, guiAlpha))
					}
				} else if PropertyValue(wrapMode) == TEXT_WRAP_CHAR || PropertyValue(wrapMode) == TEXT_WRAP_WORD {
					if textBoundsPosition.Y+float32(textOffsetY) <= textBounds.Y+textBounds.Height-textSize {
						rl.DrawTextCodepoint(guiFont, codepoint, pos, textSize, guiFade(tint, guiAlpha))
					}
				}
			}

			textOffsetX += glyphWidthAt(index, scaleFactor) + float32(getStyle(DEFAULT, TEXT_SPACING))
		}

		if PropertyValue(wrapMode) == TEXT_WRAP_NONE {
			posOffsetY += float32(getStyle(DEFAULT, TEXT_SIZE) + getStyle(DEFAULT, TEXT_LINE_SPACING))
		} else if PropertyValue(wrapMode) == TEXT_WRAP_CHAR || PropertyValue(wrapMode) == TEXT_WRAP_WORD {
			posOffsetY += float32(textOffsetY + getStyle(DEFAULT, TEXT_SIZE))
		}
	}
}

func guiDrawRectangle(rec rl.Rectangle, borderWidth int32, borderColor, color rl.Color) {
	if color.A > 0 {
		rl.DrawRectangle(int32(rec.X), int32(rec.Y), int32(rec.Width), int32(rec.Height), guiFade(color, guiAlpha))
	}

	if borderWidth > 0 {
		rl.DrawRectangle(int32(rec.X), int32(rec.Y), int32(rec.Width), borderWidth, guiFade(borderColor, guiAlpha))
		rl.DrawRectangle(int32(rec.X), int32(rec.Y)+borderWidth, borderWidth, int32(rec.Height)-2*borderWidth, guiFade(borderColor, guiAlpha))
		rl.DrawRectangle(int32(rec.X)+int32(rec.Width)-borderWidth, int32(rec.Y)+borderWidth, borderWidth, int32(rec.Height)-2*borderWidth, guiFade(borderColor, guiAlpha))
		rl.DrawRectangle(int32(rec.X), int32(rec.Y)+int32(rec.Height)-borderWidth, int32(rec.Width), borderWidth, guiFade(borderColor, guiAlpha))
	}
}

func guiMeasureText(text string) rl.Vector2 {
	var textSize rl.Vector2

	if guiFont.Texture.ID > 0 {
		size := 0
		for i := range maxLineBufferSize {
			if at(text, i) != 0 && at(text, i) != '\n' {
				size++
			} else {
				break
			}
		}

		scaleFactor := float32(getStyle(DEFAULT, TEXT_SIZE)) / float32(guiFont.BaseSize)
		textSize.Y = float32(guiFont.BaseSize) * scaleFactor

		var codepointSize int32
		for i := 0; i < size; i += int(codepointSize) {
			codepoint := getCodepointNext(tail(text, i), &codepointSize)
			codepointIndex := rl.GetGlyphIndex(guiFont, codepoint)

			textSize.X += glyphWidthAt(codepointIndex, scaleFactor) + float32(getStyle(DEFAULT, TEXT_SPACING))
		}
	}

	return textSize
}

func drawTooltip(controlRec rl.Rectangle) {
	if !guiLocked && guiTooltip && guiTooltipPtr != "" && !guiControlExclusiveMode {
		textSize := guiMeasureText(guiTooltipPtr)

		if controlRec.X+textSize.X+16 > float32(rl.GetScreenWidth()) {
			controlRec.X -= textSize.X + 16 - controlRec.Width
		}

		lineCount := int32(len(getTextLines(guiTooltipPtr)))
		if controlRec.Y+controlRec.Height+textSize.Y+4+float32(getStyle(DEFAULT, TEXT_LINE_SPACING)*lineCount) > float32(rl.GetScreenHeight()) {
			controlRec.Y -= controlRec.Height + textSize.Y + 4 + float32(getStyle(DEFAULT, TEXT_LINE_SPACING)*lineCount)
		}

		Panel(rl.Rectangle{X: controlRec.X, Y: controlRec.Y + controlRec.Height + 4, Width: textSize.X + 16, Height: textSize.Y + float32(getStyle(DEFAULT, TEXT_LINE_SPACING)*lineCount)}, "")

		textPadding := getStyle(LABEL, TEXT_PADDING)
		textAlignment := getStyle(LABEL, TEXT_ALIGNMENT)
		setStyle(LABEL, TEXT_PADDING, 0)
		setStyle(LABEL, TEXT_ALIGNMENT, int32(TEXT_ALIGN_CENTER))
		Label(rl.Rectangle{X: controlRec.X, Y: controlRec.Y + controlRec.Height + 4, Width: textSize.X + 16, Height: textSize.Y + 8.0*float32(lineCount)}, guiTooltipPtr)
		setStyle(LABEL, TEXT_ALIGNMENT, textAlignment)
		setStyle(LABEL, TEXT_PADDING, textPadding)
	}
}

func guiScrollBar(bounds rl.Rectangle, value, minValue, maxValue int32) int32 {
	state := guiState

	isVertical := !(bounds.Width > bounds.Height)

	var spinnerSize int32
	if getStyle(SCROLLBAR, ARROWS_VISIBLE) != 0 {
		if isVertical {
			spinnerSize = int32(bounds.Width) - 2*getStyle(SCROLLBAR, BORDER_WIDTH)
		} else {
			spinnerSize = int32(bounds.Height) - 2*getStyle(SCROLLBAR, BORDER_WIDTH)
		}
	}

	var arrowUpLeft, arrowDownRight, scrollbar, slider rl.Rectangle

	if value > maxValue {
		value = maxValue
	}
	if value < minValue {
		value = minValue
	}

	valueRange := maxValue - minValue
	if valueRange <= 0 {
		valueRange = 1
	}

	sliderSize := max(getStyle(SCROLLBAR, SCROLL_SLIDER_SIZE), 1)

	bw := getStyle(SCROLLBAR, BORDER_WIDTH)
	arrowUpLeft = rl.Rectangle{X: bounds.X + float32(bw), Y: bounds.Y + float32(bw), Width: float32(spinnerSize), Height: float32(spinnerSize)}

	if isVertical {
		arrowDownRight = rl.Rectangle{X: bounds.X + float32(bw), Y: bounds.Y + bounds.Height - float32(spinnerSize) - float32(bw), Width: float32(spinnerSize), Height: float32(spinnerSize)}
		scrollbar = rl.Rectangle{X: bounds.X + float32(bw+getStyle(SCROLLBAR, SCROLL_PADDING)), Y: arrowUpLeft.Y + arrowUpLeft.Height, Width: bounds.Width - float32(2*(bw+getStyle(SCROLLBAR, SCROLL_PADDING))), Height: bounds.Height - arrowUpLeft.Height - arrowDownRight.Height - float32(2*bw)}

		if float32(sliderSize) >= scrollbar.Height {
			sliderSize = int32(scrollbar.Height) - 2
		}
		slider = rl.Rectangle{
			X:      bounds.X + float32(bw+getStyle(SCROLLBAR, SCROLL_SLIDER_PADDING)),
			Y:      scrollbar.Y + float32(int32((float32(value-minValue)/float32(valueRange))*(scrollbar.Height-float32(sliderSize)))),
			Width:  bounds.Width - float32(2*(bw+getStyle(SCROLLBAR, SCROLL_SLIDER_PADDING))),
			Height: float32(sliderSize),
		}
	} else {
		arrowDownRight = rl.Rectangle{X: bounds.X + bounds.Width - float32(spinnerSize) - float32(bw), Y: bounds.Y + float32(bw), Width: float32(spinnerSize), Height: float32(spinnerSize)}
		scrollbar = rl.Rectangle{X: arrowUpLeft.X + arrowUpLeft.Width, Y: bounds.Y + float32(bw+getStyle(SCROLLBAR, SCROLL_PADDING)), Width: bounds.Width - arrowUpLeft.Width - arrowDownRight.Width - float32(2*bw), Height: bounds.Height - float32(2*(bw+getStyle(SCROLLBAR, SCROLL_PADDING)))}

		if float32(sliderSize) >= scrollbar.Width {
			sliderSize = int32(scrollbar.Width) - 2
		}
		slider = rl.Rectangle{
			X:      scrollbar.X + float32(int32((float32(value-minValue)/float32(valueRange))*(scrollbar.Width-float32(sliderSize)))),
			Y:      bounds.Y + float32(bw+getStyle(SCROLLBAR, SCROLL_SLIDER_PADDING)),
			Width:  float32(sliderSize),
			Height: bounds.Height - float32(2*(bw+getStyle(SCROLLBAR, SCROLL_SLIDER_PADDING))),
		}
	}

	if state != STATE_DISABLED && !guiLocked {
		mousePoint := guiPointerPosition()

		if guiControlExclusiveMode {
			if guiButtonDown() && !rl.CheckCollisionPointRec(mousePoint, arrowUpLeft) && !rl.CheckCollisionPointRec(mousePoint, arrowDownRight) {
				if checkBoundsID(bounds, guiControlExclusiveRec) {
					state = STATE_PRESSED
					if isVertical {
						value = int32((mousePoint.Y-scrollbar.Y-slider.Height/2)*float32(valueRange)/(scrollbar.Height-slider.Height) + float32(minValue))
					} else {
						value = int32((mousePoint.X-scrollbar.X-slider.Width/2)*float32(valueRange)/(scrollbar.Width-slider.Width) + float32(minValue))
					}
				}
			} else {
				guiControlExclusiveMode = false
				guiControlExclusiveRec = rl.Rectangle{}
			}
		} else if rl.CheckCollisionPointRec(mousePoint, bounds) {
			state = STATE_FOCUSED

			scrollDelta := guiScrollDelta()
			if scrollDelta != 0 {
				value += int32(scrollDelta)
			}

			if guiButtonPressed() {
				guiControlExclusiveMode = true
				guiControlExclusiveRec = bounds

				if rl.CheckCollisionPointRec(mousePoint, arrowUpLeft) {
					value -= valueRange / getStyle(SCROLLBAR, SCROLL_SPEED)
				} else if rl.CheckCollisionPointRec(mousePoint, arrowDownRight) {
					value += valueRange / getStyle(SCROLLBAR, SCROLL_SPEED)
				} else if !rl.CheckCollisionPointRec(mousePoint, slider) {
					if isVertical {
						value = int32((mousePoint.Y-scrollbar.Y-slider.Height/2)*float32(valueRange)/(scrollbar.Height-slider.Height) + float32(minValue))
					} else {
						value = int32((mousePoint.X-scrollbar.X-slider.Width/2)*float32(valueRange)/(scrollbar.Width-slider.Width) + float32(minValue))
					}
				}

				state = STATE_PRESSED
			}
		}

		if value > maxValue {
			value = maxValue
		}
		if value < minValue {
			value = minValue
		}
	}

	guiDrawRectangle(bounds, bw, getColor(getStyle(LISTVIEW, borderElement+PropertyID(state*3))), getColor(getStyle(DEFAULT, BORDER_COLOR_DISABLED)))
	guiDrawRectangle(scrollbar, 0, blank, getColor(getStyle(BUTTON, BASE_COLOR_NORMAL)))
	guiDrawRectangle(slider, 0, blank, getColor(getStyle(SLIDER, borderElement+PropertyID(state*3))))

	if getStyle(SCROLLBAR, ARROWS_VISIBLE) != 0 {
		side := bounds.Height
		if isVertical {
			side = bounds.Width
		}
		upLeft, downRight := IconText(ICON_ARROW_LEFT_FILL, ""), IconText(ICON_ARROW_RIGHT_FILL, "")
		if isVertical {
			upLeft, downRight = IconText(ICON_ARROW_UP_FILL, ""), IconText(ICON_ARROW_DOWN_FILL, "")
		}
		guiDrawText(upLeft, rl.Rectangle{X: arrowUpLeft.X, Y: arrowUpLeft.Y, Width: side, Height: side}, int32(TEXT_ALIGN_CENTER), getColor(getStyle(SCROLLBAR, textElement+PropertyID(state*3))))
		guiDrawText(downRight, rl.Rectangle{X: arrowDownRight.X, Y: arrowDownRight.Y, Width: side, Height: side}, int32(TEXT_ALIGN_CENTER), getColor(getStyle(SCROLLBAR, textElement+PropertyID(state*3))))
	}

	return value
}

func guiTextSplit(text string, delimiter byte) []string {
	var buffer [textSplitMaxTextSize]byte
	starts := []int{0}

	n := 0
	for i := range textSplitMaxTextSize {
		buffer[i] = at(text, i)
		n = i + 1
		if buffer[i] == 0 {
			break
		} else if buffer[i] == delimiter || buffer[i] == '\n' {
			starts = append(starts, i+1)
			buffer[i] = 0

			if len(starts) >= textSplitMaxItems {
				break
			}
		}
	}

	items := make([]string, len(starts))
	for k, s := range starts {
		if s >= n {
			items[k] = ""
			continue
		}
		items[k] = cstr(string(buffer[s:n]))
	}

	return items
}

func convertRGBtoHSV(rgb rl.Vector3) rl.Vector3 {
	var hsv rl.Vector3

	minV := rgb.X
	if rgb.Y < minV {
		minV = rgb.Y
	}
	if rgb.Z < minV {
		minV = rgb.Z
	}

	maxV := rgb.X
	if rgb.Y > maxV {
		maxV = rgb.Y
	}
	if rgb.Z > maxV {
		maxV = rgb.Z
	}

	hsv.Z = maxV
	delta := maxV - minV

	if delta < 0.00001 {
		hsv.Y = 0.0
		hsv.X = 0.0
		return hsv
	}

	if maxV > 0.0 {
		hsv.Y = delta / maxV
	} else {
		hsv.Y = 0.0
		hsv.X = 0.0
		return hsv
	}

	if rgb.X >= maxV {
		hsv.X = (rgb.Y - rgb.Z) / delta
	} else {
		if rgb.Y >= maxV {
			hsv.X = 2.0 + (rgb.Z-rgb.X)/delta
		} else {
			hsv.X = 4.0 + (rgb.X-rgb.Y)/delta
		}
	}

	hsv.X *= 60.0

	if hsv.X < 0.0 {
		hsv.X += 360.0
	}

	return hsv
}

func convertHSVtoRGB(hsv rl.Vector3) rl.Vector3 {
	var rgb rl.Vector3

	if hsv.Y <= 0.0 {
		rgb.X = hsv.Z
		rgb.Y = hsv.Z
		rgb.Z = hsv.Z
		return rgb
	}

	hh := hsv.X
	if hh >= 360.0 {
		hh = 0.0
	}
	hh /= 60.0

	i := int64(hh)
	ff := hh - float32(i)
	p := hsv.Z * (1.0 - hsv.Y)
	q := hsv.Z * (1.0 - (hsv.Y * ff))
	t := hsv.Z * (1.0 - (hsv.Y * (1.0 - ff)))

	switch i {
	case 0:
		rgb.X, rgb.Y, rgb.Z = hsv.Z, t, p
	case 1:
		rgb.X, rgb.Y, rgb.Z = q, hsv.Z, p
	case 2:
		rgb.X, rgb.Y, rgb.Z = p, hsv.Z, t
	case 3:
		rgb.X, rgb.Y, rgb.Z = p, q, hsv.Z
	case 4:
		rgb.X, rgb.Y, rgb.Z = t, p, hsv.Z
	default:
		rgb.X, rgb.Y, rgb.Z = hsv.Z, p, q
	}

	return rgb
}

func guiFade(color rl.Color, alpha float32) rl.Color {
	if alpha < 0.0 {
		alpha = 0.0
	} else if alpha > 1.0 {
		alpha = 1.0
	}

	return rl.Color{R: color.R, G: color.G, B: color.B, A: uint8(float32(color.A) * alpha)}
}
