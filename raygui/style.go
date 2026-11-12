package raygui

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// LoadStyle - Load style file over global style variable (.rgs)
func LoadStyle(fileName string) {
	tryBinary := false
	if !guiStyleLoaded {
		LoadStyleDefault()
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	lines := rgsTextLines(data)
	if len(lines) > 0 && at(lines[0], 0) == '#' {
		var controlId, propertyId int64
		var propertyValue, version uint64

		for _, buffer := range lines {
			switch at(buffer, 0) {
			case 'v':
				scanStyleVersion(buffer, &version)
			case 'p':
				scanStyleProperty(buffer, &controlId, &propertyId, &propertyValue)
				if controlId >= 0 && controlId < RAYGUI_MAX_CONTROLS && propertyId >= 0 && propertyId < RAYGUI_MAX_PROPS_BASE+RAYGUI_MAX_PROPS_EXTENDED {
					SetStyle(ControlID(controlId), PropertyID(propertyId), PropertyValue(uint32(propertyValue)))
				}
			case 'f':
				fontSize, first, second := scanStyleFont(buffer)

				charmapFileName, fontFileName := first, second
				if version >= 600 {
					fontFileName, charmapFileName = first, second
				}

				var font rl.Font
				var codepoints []rune

				if at(charmapFileName, 0) != '0' {
					textData, _ := os.ReadFile(filepath.Dir(fileName) + "/" + charmapFileName)
					codepoints = loadCodepoints(cstr(string(textData)))
				}

				if fontFileName != "" {
					path := filepath.Dir(fileName) + "/" + fontFileName
					if len(codepoints) > 0 {
						font = rl.LoadFontEx(path, int32(fontSize), codepoints, int32(len(codepoints)))
					} else {
						font = rl.LoadFontEx(path, int32(fontSize), nil, 0)
					}
				}

				if font.Texture.ID == 0 {
					font = rl.GetFontDefault()
					SetStyle(DEFAULT, TEXT_SIZE, 10)
					SetStyle(DEFAULT, TEXT_SPACING, 1)
				}

				if font.Texture.ID > 0 && font.CharsCount > 0 {
					SetFont(font)
				}
			}
		}
	} else {
		tryBinary = true
	}

	if tryBinary && len(data) > 0 {
		LoadStyleFromMemory(data)
	}
}

func rgsTextLines(data []byte) []string {
	var lines []string
	for len(data) > 0 {
		n := bytes.IndexByte(data, '\n')
		if n < 0 {
			break
		}
		line := data[:n+1]
		for len(line) > maxLineBufferSize-1 {
			lines = append(lines, string(line[:maxLineBufferSize-1]))
			line = line[maxLineBufferSize-1:]
		}
		lines = append(lines, string(line))
		data = data[n+1:]
	}
	return lines
}

type scanner struct {
	s string
	i int
}

func (c *scanner) space() {
	for isSpace(int32(at(c.s, c.i))) && at(c.s, c.i) != 0 {
		c.i++
	}
}

func (c *scanner) lit(b byte) bool {
	if at(c.s, c.i) == b {
		c.i++
		return true
	}
	return false
}

func (c *scanner) digits(base uint64) (uint64, bool) {
	var v uint64
	n := 0
	for {
		ch := at(c.s, c.i)
		var d uint64
		switch {
		case ch >= '0' && ch <= '9':
			d = uint64(ch - '0')
		case base == 16 && ch >= 'a' && ch <= 'f':
			d = uint64(ch-'a') + 10
		case base == 16 && ch >= 'A' && ch <= 'F':
			d = uint64(ch-'A') + 10
		default:
			return v, n > 0
		}
		v = v*base + d
		c.i++
		n++
	}
}

func (c *scanner) int() (int64, bool) {
	c.space()
	neg := false
	if at(c.s, c.i) == '+' || at(c.s, c.i) == '-' {
		neg = at(c.s, c.i) == '-'
		c.i++
	}
	v, ok := c.digits(10)
	if neg {
		return -int64(v), ok
	}
	return int64(v), ok
}

func (c *scanner) uint() (uint64, bool) {
	v, ok := c.int()
	return uint64(uint32(v)), ok
}

func (c *scanner) hex() (uint64, bool) {
	c.space()
	neg := false
	if at(c.s, c.i) == '+' || at(c.s, c.i) == '-' {
		neg = at(c.s, c.i) == '-'
		c.i++
	}
	if at(c.s, c.i) == '0' && (at(c.s, c.i+1) == 'x' || at(c.s, c.i+1) == 'X') && isHexDigit(at(c.s, c.i+2)) {
		c.i += 2
	}
	v, ok := c.digits(16)
	if neg {
		v = -v
	}
	return v, ok
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func (c *scanner) word(max int) (string, bool) {
	c.space()
	start := c.i
	for c.i-start < max && at(c.s, c.i) != 0 && !isSpace(int32(at(c.s, c.i))) {
		c.i++
	}
	return c.s[start:c.i], c.i > start
}

func (c *scanner) line(max int) (string, bool) {
	start := c.i
	for c.i-start < max && at(c.s, c.i) != 0 && at(c.s, c.i) != '\r' && at(c.s, c.i) != '\n' {
		c.i++
	}
	return c.s[start:c.i], c.i > start
}

func scanStyleVersion(buffer string, version *uint64) {
	c := &scanner{s: buffer, i: 1}
	if v, ok := c.uint(); ok {
		*version = v
	}
}

func scanStyleProperty(buffer string, controlId, propertyId *int64, propertyValue *uint64) {
	c := &scanner{s: buffer, i: 1}
	if v, ok := c.int(); ok {
		*controlId = v
		c.space()
		if v, ok := c.int(); ok {
			*propertyId = v
			c.space()
			if c.lit('0') && c.lit('x') {
				if v, ok := c.hex(); ok {
					*propertyValue = v
				}
			}
		}
	}
}

func scanStyleFont(buffer string) (int64, string, string) {
	var fontSize int64
	var first, second string

	c := &scanner{s: buffer, i: 1}
	if v, ok := c.int(); ok {
		fontSize = v
		c.space()
		if v, ok := c.word(31); ok {
			first = v
			c.space()
			if v, ok := c.line(31); ok {
				second = v
			}
		}
	}

	return fontSize, first, second
}

func loadCodepoints(text string) []rune {
	var codepoints []rune
	var codepointSize int32
	for i := 0; i < len(text); i += int(codepointSize) {
		codepoints = append(codepoints, getCodepointNext(tail(text, i), &codepointSize))
	}
	return codepoints
}

func decompressData(compData []byte) []byte {
	data, _ := io.ReadAll(flate.NewReader(bytes.NewReader(compData)))
	return data
}

type rgsReader struct {
	data []byte
	pos  int
}

func (r *rgsReader) bytes(n int) []byte {
	out := make([]byte, n)
	if r.pos < len(r.data) {
		copy(out, r.data[r.pos:])
	}
	r.pos += n
	return out
}

func (r *rgsReader) int16() int16 {
	return int16(binary.LittleEndian.Uint16(r.bytes(2)))
}

func (r *rgsReader) int32() int32 {
	return int32(binary.LittleEndian.Uint32(r.bytes(4)))
}

func (r *rgsReader) uint32() uint32 {
	return binary.LittleEndian.Uint32(r.bytes(4))
}

func (r *rgsReader) float32() float32 {
	return math.Float32frombits(r.uint32())
}

// LoadStyleFromMemory - Load style from memory (binary only)
func LoadStyleFromMemory(fileData []byte) {
	r := &rgsReader{data: fileData}

	signature := r.bytes(4)
	version := r.int16()
	r.int16()
	propertyCount := r.int32()

	if string(signature) != "rGS " {
		return
	}

	for range propertyCount {
		controlId := r.int16()
		propertyId := r.int16()
		propertyValue := r.uint32()

		if controlId == 0 {
			SetStyle(0, PropertyID(propertyId), PropertyValue(propertyValue))

			if propertyId < RAYGUI_MAX_PROPS_BASE {
				for j := 1; j < RAYGUI_MAX_CONTROLS; j++ {
					SetStyle(ControlID(j), PropertyID(propertyId), PropertyValue(propertyValue))
				}
			}
		} else {
			SetStyle(ControlID(controlId), PropertyID(propertyId), PropertyValue(propertyValue))
		}
	}

	fontDataSize := r.int32()

	if fontDataSize > 0 {
		var font rl.Font

		if version >= 600 {
			r.bytes(32)
		}

		font.BaseSize = r.int32()
		font.CharsCount = r.int32()
		r.int32()

		fontWhiteRec := rl.Rectangle{X: r.float32(), Y: r.float32(), Width: r.float32(), Height: r.float32()}

		fontImageUncompSize := r.int32()
		fontImageCompSize := r.int32()

		var imFont rl.Image
		imFont.Mipmaps = 1
		imFont.Width = r.int32()
		imFont.Height = r.int32()
		imFont.Format = rl.PixelFormat(r.int32())

		var imData []byte
		if fontImageCompSize > 0 && fontImageCompSize != fontImageUncompSize {
			imData = decompressData(r.bytes(int(fontImageCompSize)))
			if int32(len(imData)) != fontImageUncompSize {
				rl.TraceLog(rl.LogWarning, "Uncompressed font atlas image data could be corrupted")
			}
		} else {
			imData = r.bytes(int(fontImageUncompSize))
		}

		glyphCount := int(font.CharsCount)

		recsDataSize := int32(glyphCount * 16)
		var recsDataCompressedSize int32

		if version >= 400 {
			recsDataCompressedSize = r.int32()
		}

		recs := make([]rl.Rectangle, glyphCount)
		if recsDataCompressedSize > 0 && recsDataCompressedSize != recsDataSize {
			recsData := decompressData(r.bytes(int(recsDataCompressedSize)))
			if int32(len(recsData)) != recsDataSize {
				rl.TraceLog(rl.LogWarning, "Uncompressed font recs data could be corrupted")
			}
			rr := &rgsReader{data: recsData}
			for i := range recs {
				recs[i] = rl.Rectangle{X: rr.float32(), Y: rr.float32(), Width: rr.float32(), Height: rr.float32()}
			}
		} else {
			for i := range recs {
				recs[i] = rl.Rectangle{X: r.float32(), Y: r.float32(), Width: r.float32(), Height: r.float32()}
			}
		}

		glyphsDataSize := int32(glyphCount * 16)
		var glyphsDataCompressedSize int32

		if version >= 400 {
			glyphsDataCompressedSize = r.int32()
		}

		glyphs := make([]rl.GlyphInfo, glyphCount)
		gr := r
		if glyphsDataCompressedSize > 0 && glyphsDataCompressedSize != glyphsDataSize {
			glyphsData := decompressData(r.bytes(int(glyphsDataCompressedSize)))
			if int32(len(glyphsData)) != glyphsDataSize {
				rl.TraceLog(rl.LogWarning, "Uncompressed font glyphs data could be corrupted")
			}
			gr = &rgsReader{data: glyphsData}
		}
		for i := range glyphs {
			glyphs[i].Value = gr.int32()
			glyphs[i].OffsetX = gr.int32()
			glyphs[i].OffsetY = gr.int32()
			glyphs[i].AdvanceX = gr.int32()
		}

		if glyphCount > 0 {
			font.Recs = &recs[0]
			font.Chars = &glyphs[0]
		}

		if len(imData) > 0 {
			imFont.Data = unsafe.Pointer(&imData[0])
		}

		font.Texture = rl.LoadTextureFromImage(&imFont)

		if font.Texture.ID != 0 {
			if fontWhiteRec.X > 0 && fontWhiteRec.Y > 0 && fontWhiteRec.Width > 0 && fontWhiteRec.Height > 0 {
				rl.SetShapesTexture(font.Texture, fontWhiteRec)
			}
		} else {
			font = rl.GetFontDefault()
		}

		SetFont(font)
	}
}

// LoadStyleDefault - Load style default over global style
func LoadStyleDefault() {
	guiStyleLoaded = true

	SetStyle(DEFAULT, BORDER_COLOR_NORMAL, 0x838383ff)
	SetStyle(DEFAULT, BASE_COLOR_NORMAL, 0xc9c9c9ff)
	SetStyle(DEFAULT, TEXT_COLOR_NORMAL, 0x686868ff)
	SetStyle(DEFAULT, BORDER_COLOR_FOCUSED, 0x5bb2d9ff)
	SetStyle(DEFAULT, BASE_COLOR_FOCUSED, 0xc9effeff)
	SetStyle(DEFAULT, TEXT_COLOR_FOCUSED, 0x6c9bbcff)
	SetStyle(DEFAULT, BORDER_COLOR_PRESSED, 0x0492c7ff)
	SetStyle(DEFAULT, BASE_COLOR_PRESSED, 0x97e8ffff)
	SetStyle(DEFAULT, TEXT_COLOR_PRESSED, 0x368bafff)
	SetStyle(DEFAULT, BORDER_COLOR_DISABLED, 0xb5c1c2ff)
	SetStyle(DEFAULT, BASE_COLOR_DISABLED, 0xe6e9e9ff)
	SetStyle(DEFAULT, TEXT_COLOR_DISABLED, 0xaeb7b8ff)
	SetStyle(DEFAULT, BORDER_WIDTH, 1)
	SetStyle(DEFAULT, TEXT_PADDING, 0)
	SetStyle(DEFAULT, TEXT_ALIGNMENT, TEXT_ALIGN_CENTER)

	SetStyle(DEFAULT, TEXT_SIZE, 10)
	SetStyle(DEFAULT, TEXT_SPACING, 1)
	SetStyle(DEFAULT, LINE_COLOR, 0x90abb5ff)
	SetStyle(DEFAULT, BACKGROUND_COLOR, 0xf5f5f5ff)
	SetStyle(DEFAULT, TEXT_LINE_SPACING, 12)
	SetStyle(DEFAULT, TEXT_ALIGNMENT_VERTICAL, TEXT_ALIGN_MIDDLE)

	SetStyle(LABEL, TEXT_ALIGNMENT, TEXT_ALIGN_LEFT)
	SetStyle(BUTTON, BORDER_WIDTH, 2)
	SetStyle(SLIDER, TEXT_PADDING, 4)
	SetStyle(PROGRESSBAR, TEXT_PADDING, 4)
	SetStyle(CHECKBOX, TEXT_PADDING, 4)
	SetStyle(CHECKBOX, TEXT_ALIGNMENT, TEXT_ALIGN_RIGHT)
	SetStyle(DROPDOWNBOX, TEXT_PADDING, 0)
	SetStyle(DROPDOWNBOX, TEXT_ALIGNMENT, TEXT_ALIGN_CENTER)
	SetStyle(TEXTBOX, TEXT_PADDING, 4)
	SetStyle(TEXTBOX, TEXT_ALIGNMENT, TEXT_ALIGN_LEFT)
	SetStyle(VALUEBOX, TEXT_PADDING, 0)
	SetStyle(VALUEBOX, TEXT_ALIGNMENT, TEXT_ALIGN_LEFT)
	SetStyle(STATUSBAR, TEXT_PADDING, 8)
	SetStyle(STATUSBAR, TEXT_ALIGNMENT, TEXT_ALIGN_LEFT)
	SetStyle(TABBAR, TAB_ITEMS_WIDTH, 160)

	SetStyle(TOGGLE, GROUP_PADDING, 2)
	SetStyle(SLIDER, SLIDER_WIDTH, 16)
	SetStyle(SLIDER, SLIDER_PADDING, 1)
	SetStyle(PROGRESSBAR, PROGRESS_PADDING, 1)
	SetStyle(CHECKBOX, CHECK_PADDING, 1)
	SetStyle(COMBOBOX, COMBO_BUTTON_WIDTH, 32)
	SetStyle(COMBOBOX, COMBO_BUTTON_SPACING, 2)
	SetStyle(DROPDOWNBOX, ARROW_PADDING, 16)
	SetStyle(DROPDOWNBOX, DROPDOWN_ITEMS_SPACING, 2)
	SetStyle(VALUEBOX, SPINNER_BUTTON_WIDTH, 24)
	SetStyle(VALUEBOX, SPINNER_BUTTON_SPACING, 2)
	SetStyle(SCROLLBAR, BORDER_WIDTH, 0)
	SetStyle(SCROLLBAR, ARROWS_VISIBLE, 0)
	SetStyle(SCROLLBAR, ARROWS_SIZE, 6)
	SetStyle(SCROLLBAR, SCROLL_SLIDER_PADDING, 0)
	SetStyle(SCROLLBAR, SCROLL_SLIDER_SIZE, 16)
	SetStyle(SCROLLBAR, SCROLL_PADDING, 0)
	SetStyle(SCROLLBAR, SCROLL_SPEED, 12)
	SetStyle(LISTVIEW, LIST_ITEMS_HEIGHT, 28)
	SetStyle(LISTVIEW, LIST_ITEMS_SPACING, 2)
	SetStyle(LISTVIEW, LIST_ITEMS_BORDER_WIDTH, 1)
	SetStyle(LISTVIEW, SCROLLBAR_WIDTH, 12)
	SetStyle(LISTVIEW, SCROLLBAR_SIDE, SCROLLBAR_RIGHT_SIDE)
	SetStyle(COLORPICKER, COLOR_SELECTOR_SIZE, 8)
	SetStyle(COLORPICKER, HUEBAR_WIDTH, 16)
	SetStyle(COLORPICKER, HUEBAR_PADDING, 8)
	SetStyle(COLORPICKER, HUEBAR_SELECTOR_HEIGHT, 8)
	SetStyle(COLORPICKER, HUEBAR_SELECTOR_OVERFLOW, 2)

	if guiFont.Texture.ID != rl.GetFontDefault().Texture.ID {
		if guiFont.Texture.ID > 0 {
			rl.UnloadTexture(guiFont.Texture)
		}

		guiFont = rl.GetFontDefault()

		if recs := fontRecs(guiFont); len(recs) > 95 {
			whiteChar := recs[95]
			rl.SetShapesTexture(guiFont.Texture, rl.Rectangle{X: whiteChar.X + 1, Y: whiteChar.Y + 1, Width: whiteChar.Width - 2, Height: whiteChar.Height - 2})
		}
	}
}

// IconText - Get text with icon id prepended (if supported)
func IconText(iconId IconID, text string) string {
	buffer := fmt.Sprintf("#%03d#", iconId)
	if text == "" {
		return buffer
	}

	if len(buffer)+len(text) > 1023 {
		text = text[:1023-len(buffer)]
	}

	return buffer + cstr(text)
}
