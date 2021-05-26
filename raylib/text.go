package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// cptr returns C pointer
func (c *CharInfo) cptr() *C.CharInfo {
	return (*C.CharInfo)(unsafe.Pointer(c))
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
func LoadFontEx(fileName string, fontSize int32, fontChars *int32, charsCount int32) Font {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cfontSize := (C.int)(fontSize)
	cfontChars := (*C.int)(unsafe.Pointer(fontChars))
	ccharsCount := (C.int)(charsCount)
	ret := C.LoadFontEx(cfileName, cfontSize, cfontChars, ccharsCount)
	v := newFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadFontFromImage - Loads an Image font file (XNA style)
func LoadFontFromImage(image Image, key Color, firstChar int32) Font {
	cimage := image.cptr()
	ckey := key.cptr()
	cfirstChar := (C.int)(firstChar)
	ret := C.LoadFontFromImage(*cimage, *ckey, cfirstChar)
	v := newFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadFontFromMemory - Load font from memory buffer, fileType refers to extension: i.e. ".ttf"
func LoadFontFromMemory(fileType string, fileData []byte, dataSize int32, fontSize int32, fontChars *int32, charsCount int32) Font {
	cfileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cfileType))
	cfileData := (*C.uchar)(unsafe.Pointer(&fileData[0]))
	cdataSize := (C.int)(dataSize)
	cfontSize := (C.int)(fontSize)
	cfontChars := (*C.int)(unsafe.Pointer(fontChars))
	ccharsCount := (C.int)(charsCount)
	ret := C.LoadFontFromMemory(cfileType, cfileData, cdataSize, cfontSize, cfontChars, ccharsCount)
	v := newFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadFontData - Load font data for further use
func LoadFontData(fileData []byte, dataSize int32, fontSize int32, fontChars *int32, charsCount, typ int32) *CharInfo {
	cfileData := (*C.uchar)(unsafe.Pointer(&fileData[0]))
	cdataSize := (C.int)(dataSize)
	cfontSize := (C.int)(fontSize)
	cfontChars := (*C.int)(unsafe.Pointer(fontChars))
	ccharsCount := (C.int)(charsCount)
	ctype := (C.int)(typ)
	ret := C.LoadFontData(cfileData, cdataSize, cfontSize, cfontChars, ccharsCount, ctype)
	v := newCharInfoFromPointer(unsafe.Pointer(&ret))
	return &v
}

// UnloadFont - Unload Font from GPU memory (VRAM)
func UnloadFont(font Font) {
	cfont := font.cptr()
	C.UnloadFont(*cfont)
}

// DrawText - Draw text (using default font)
func DrawText(text string, posX int32, posY int32, fontSize int32, color Color) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cfontSize := (C.int)(fontSize)
	ccolor := color.cptr()
	C.DrawText(ctext, cposX, cposY, cfontSize, *ccolor)
}

// DrawTextEx - Draw text using Font and additional parameters
func DrawTextEx(font Font, text string, position Vector2, fontSize float32, spacing float32, tint Color) {
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cposition := position.cptr()
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	ctint := tint.cptr()
	C.DrawTextEx(*cfont, ctext, *cposition, cfontSize, cspacing, *ctint)
}

// DrawTextRec - Draw text using font inside rectangle limits
func DrawTextRec(font Font, text string, rec Rectangle, fontSize, spacing float32, wordWrap bool, tint Color) {
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	crec := rec.cptr()
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	cwordWrap := (C.bool)(wordWrap)
	ctint := tint.cptr()
	C.DrawTextRec(*cfont, ctext, *crec, cfontSize, cspacing, cwordWrap, *ctint)
}

// DrawTextRecEx - Draw text using font inside rectangle limits with support for text selection
func DrawTextRecEx(font Font, text string, rec Rectangle, fontSize, spacing float32, wordWrap bool, tint Color, selectStart, selectLength int32, selectText, selectBack Color) {
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	crec := rec.cptr()
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	cwordWrap := (C.bool)(wordWrap)
	ctint := tint.cptr()
	cselectStart := (C.int)(selectStart)
	cselectLength := (C.int)(selectLength)
	cselectText := selectText.cptr()
	cselectBack := selectBack.cptr()
	C.DrawTextRecEx(*cfont, ctext, *crec, cfontSize, cspacing, cwordWrap, *ctint, cselectStart, cselectLength, *cselectText, *cselectBack)
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

// GetGlyphIndex - Returns index position for a unicode character on spritefont
func GetGlyphIndex(font Font, character int32) int32 {
	cfont := font.cptr()
	ccharacter := (C.int)(character)
	ret := C.GetGlyphIndex(*cfont, ccharacter)
	v := (int32)(ret)
	return v
}

// DrawFPS - Shows current FPS
func DrawFPS(posX int32, posY int32) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	C.DrawFPS(cposX, cposY)
}
