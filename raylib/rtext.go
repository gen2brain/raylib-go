package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import (
	"image/color"
	"unsafe"
)

// newGlyphInfoFromPointer - Returns new GlyphInfo from pointer
func newGlyphInfoFromPointer(ptr unsafe.Pointer) GlyphInfo {
	return *(*GlyphInfo)(ptr)
}

// cptr returns C pointer
func (c *GlyphInfo) cptr() *C.GlyphInfo {
	return (*C.GlyphInfo)(unsafe.Pointer(c))
}

// newFontFromPointer - Returns new Font from pointer
func newFontFromPointer(ptr unsafe.Pointer) Font {
	return *(*Font)(ptr)
}

// cptr returns C pointer
func (s *Font) cptr() *C.Font {
	return (*C.Font)(unsafe.Pointer(s))
}

// GetFontDefault - Get the default Font
func GetFontDefault() Font {
	ret := C.GetFontDefault()
	v := newFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadFont - Load a Font image into GPU memory (VRAM)
func LoadFont(fileName string) Font {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadFont(cfileName)
	v := newFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadFontEx - Load Font from file with extended parameters
func LoadFontEx(fileName string, fontSize int32, fontChars []rune, runesNumber ...int32) Font {
	var cfontChars *C.int
	var ccharsCount C.int

	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cfontSize := (C.int)(fontSize)
	if fontChars != nil {
		cfontChars = (*C.int)(unsafe.Pointer(&fontChars[0]))
		ccharsCount = (C.int)(len(fontChars))
	}
	if fontChars == nil {
		if len(runesNumber) > 0 {
			ccharsCount = (C.int)(runesNumber[0])
		}
	}
	ret := C.LoadFontEx(cfileName, cfontSize, cfontChars, ccharsCount)
	v := newFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadFontFromImage - Loads an Image font file (XNA style)
func LoadFontFromImage(image Image, key color.RGBA, firstChar int32) Font {
	cimage := image.cptr()
	ckey := colorCptr(key)
	cfirstChar := (C.int)(firstChar)
	ret := C.LoadFontFromImage(*cimage, *ckey, cfirstChar)
	v := newFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadFontFromMemory - Load font from memory buffer, fileType refers to extension: i.e. ".ttf"
func LoadFontFromMemory(fileType string, fileData []byte, fontSize int32, codepoints []rune) Font {
	cfileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cfileType))
	cfileData := (*C.uchar)(unsafe.Pointer(&fileData[0]))
	cdataSize := (C.int)(len(fileData))
	cfontSize := (C.int)(fontSize)
	cfontChars := (*C.int)(unsafe.SliceData(codepoints))
	ccharsCount := (C.int)(len(codepoints))
	ret := C.LoadFontFromMemory(cfileType, cfileData, cdataSize, cfontSize, cfontChars, ccharsCount)
	v := newFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// IsFontValid - Check if a font is valid (font data loaded, WARNING: GPU texture not checked)
func IsFontValid(font Font) bool {
	cfont := font.cptr()
	ret := C.IsFontValid(*cfont)
	v := bool(ret)
	return v
}

// LoadFontData - Load font data for further use
func LoadFontData(fileData []byte, fontSize int32, codePoints []rune, codepointCount, typ int32) []GlyphInfo {
	cfileData := (*C.uchar)(unsafe.Pointer(&fileData[0]))
	cdataSize := (C.int)(len(fileData))
	cfontSize := (C.int)(fontSize)
	ccodePoints := (*C.int)(unsafe.SliceData(codePoints))
	// In case no chars count provided, default to 95
	if codepointCount <= 0 {
		codepointCount = 95
	}
	ccodePointCount := (C.int)(codepointCount)
	ctype := (C.int)(typ)
	ret := C.LoadFontData(cfileData, cdataSize, cfontSize, ccodePoints, ccodePointCount, ctype)
	v := unsafe.Slice((*GlyphInfo)(unsafe.Pointer(ret)), ccodePointCount)
	return v
}

// UnloadFontData - Unload font chars info data (RAM)
func UnloadFontData(glyphs []GlyphInfo) {
	cglyphs := (*C.GlyphInfo)(unsafe.Pointer(&glyphs[0]))
	C.UnloadFontData(cglyphs, C.int(len(glyphs)))
}

// UnloadFont - Unload Font from GPU memory (VRAM)
func UnloadFont(font Font) {
	cfont := font.cptr()
	C.UnloadFont(*cfont)
}

// DrawFPS - Shows current FPS
func DrawFPS(posX int32, posY int32) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	C.DrawFPS(cposX, cposY)
}

// DrawText - Draw text (using default font)
func DrawText(text string, posX int32, posY int32, fontSize int32, col color.RGBA) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cfontSize := (C.int)(fontSize)
	ccolor := colorCptr(col)
	C.DrawText(ctext, cposX, cposY, cfontSize, *ccolor)
}

// DrawTextEx - Draw text using Font and additional parameters
func DrawTextEx(font Font, text string, position Vector2, fontSize float32, spacing float32, tint color.RGBA) {
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cposition := position.cptr()
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	ctint := colorCptr(tint)
	C.DrawTextEx(*cfont, ctext, *cposition, cfontSize, cspacing, *ctint)
}

// DrawTextPro - Draw text using Font and pro parameters (rotation)
func DrawTextPro(font Font, text string, position Vector2, origin Vector2, rotation, fontSize float32, spacing float32, tint color.RGBA) {
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cposition := position.cptr()
	crotation := (C.float)(rotation)
	corigin := origin.cptr()
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	ctint := colorCptr(tint)
	C.DrawTextPro(*cfont, ctext, *cposition, *corigin, crotation, cfontSize, cspacing, *ctint)
}

// SetTextLineSpacing - Set vertical line spacing when drawing with line-breaks
func SetTextLineSpacing(spacing int) {
	cspacing := (C.int)(spacing)
	C.SetTextLineSpacing(cspacing)
}

// MeasureText - Measure string width for default font
func MeasureText(text string, fontSize int32) int32 {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.int)(fontSize)
	ret := C.MeasureText(ctext, cfontSize)
	v := (int32)(ret)
	return v
}

// MeasureTextEx - Measure string size for Font
func MeasureTextEx(font Font, text string, fontSize float32, spacing float32) Vector2 {
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	ret := C.MeasureTextEx(*cfont, ctext, cfontSize, cspacing)
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// GetGlyphIndex - Get glyph index position in font for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphIndex(font Font, codepoint int32) int32 {
	cfont := font.cptr()
	ccodepoint := (C.int)(codepoint)
	ret := C.GetGlyphIndex(*cfont, ccodepoint)
	v := (int32)(ret)
	return v
}

// GetGlyphInfo - Get glyph font info data for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphInfo(font Font, codepoint int32) GlyphInfo {
	cfont := font.cptr()
	ccodepoint := (C.int)(codepoint)
	ret := C.GetGlyphInfo(*cfont, ccodepoint)
	v := newGlyphInfoFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetGlyphAtlasRec - Get glyph rectangle in font atlas for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphAtlasRec(font Font, codepoint int32) Rectangle {
	cfont := font.cptr()
	ccodepoint := (C.int)(codepoint)
	ret := C.GetGlyphAtlasRec(*cfont, ccodepoint)
	v := newRectangleFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageFontAtlas - Generate image font atlas using chars info
func GenImageFontAtlas(glyphs []GlyphInfo, glyphRecs []*Rectangle, fontSize int32, padding int32, packMethod int32) Image {
	cglyphs := (*C.GlyphInfo)(unsafe.Pointer(&glyphs[0]))
	cglyphRecs := (**C.Rectangle)(unsafe.Pointer(&glyphRecs[0]))
	cglyphCount := C.int(len(glyphs))
	ret := C.GenImageFontAtlas(cglyphs, cglyphRecs, cglyphCount, C.int(fontSize), C.int(padding), C.int(packMethod))
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return *v
}

// DrawTextCodepoint - Draw one character (codepoint)
func DrawTextCodepoint(font Font, codepoint rune, position Vector2, fontSize float32, tint color.RGBA) {
	cfont := font.cptr()
	ccodepoint := (C.int)(codepoint)
	cposition := position.cptr()
	cfontSize := (C.float)(fontSize)
	ctint := colorCptr(tint)
	C.DrawTextCodepoint(*cfont, ccodepoint, *cposition, cfontSize, *ctint)
}

// DrawTextCodepoints - Draw multiple character (codepoint)
func DrawTextCodepoints(font Font, codepoints []rune, position Vector2, fontSize float32, spacing float32, tint color.RGBA) {
	cfont := font.cptr()
	ccodepoints := (*C.int)(unsafe.SliceData(codepoints))
	ccodepointCount := C.int(len(codepoints))
	cposition := position.cptr()
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	ctint := colorCptr(tint)
	C.DrawTextCodepoints(*cfont, ccodepoints, ccodepointCount, *cposition, cfontSize, cspacing, *ctint)
}
