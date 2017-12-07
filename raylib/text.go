package raylib

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// CharInfo - SpriteFont character info
type CharInfo struct {
	// Character value (Unicode)
	Value int32
	// Character rectangle in sprite font
	Rec Rectangle
	// Character offset X when drawing
	OffsetX int32
	// Character offset Y when drawing
	OffsetY int32
	// Character advance position X
	AdvanceX int32
}

func (c *CharInfo) cptr() *C.CharInfo {
	return (*C.CharInfo)(unsafe.Pointer(c))
}

// NewCharInfo - Returns new SpriteFont
func NewCharInfo(value int32, rec Rectangle, offsetX, offsetY, advanceX int32) CharInfo {
	return CharInfo{value, rec, offsetX, offsetY, advanceX}
}

// newCharInfoFromPointer - Returns new SpriteFont from pointer
func newCharInfoFromPointer(ptr unsafe.Pointer) CharInfo {
	return *(*CharInfo)(ptr)
}

// SpriteFont type, includes texture and charSet array data
type SpriteFont struct {
	// Font texture
	Texture Texture2D
	// Base size (default chars height)
	BaseSize int32
	// Number of characters
	CharsCount int32
	// Characters info data
	Chars *CharInfo
}

func (s *SpriteFont) cptr() *C.SpriteFont {
	return (*C.SpriteFont)(unsafe.Pointer(s))
}

// NewSpriteFont - Returns new SpriteFont
func NewSpriteFont(texture Texture2D, baseSize, charsCount int32, chars *CharInfo) SpriteFont {
	return SpriteFont{texture, baseSize, charsCount, chars}
}

// newSpriteFontFromPointer - Returns new SpriteFont from pointer
func newSpriteFontFromPointer(ptr unsafe.Pointer) SpriteFont {
	return *(*SpriteFont)(ptr)
}

// GetDefaultFont - Get the default SpriteFont
func GetDefaultFont() SpriteFont {
	ret := C.GetDefaultFont()
	v := newSpriteFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadSpriteFont - Load a SpriteFont image into GPU memory (VRAM)
func LoadSpriteFont(fileName string) SpriteFont {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadSpriteFont(cfileName)
	v := newSpriteFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadSpriteFontEx - Load SpriteFont from file with extended parameters
func LoadSpriteFontEx(fileName string, fontSize int32, charsCount int32, fontChars *int32) SpriteFont {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cfontSize := (C.int)(fontSize)
	ccharsCount := (C.int)(charsCount)
	cfontChars := (*C.int)(unsafe.Pointer(fontChars))
	ret := C.LoadSpriteFontEx(cfileName, cfontSize, ccharsCount, cfontChars)
	v := newSpriteFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// UnloadSpriteFont - Unload SpriteFont from GPU memory (VRAM)
func UnloadSpriteFont(spriteFont SpriteFont) {
	cspriteFont := spriteFont.cptr()
	C.UnloadSpriteFont(*cspriteFont)
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

// DrawTextEx - Draw text using SpriteFont and additional parameters
func DrawTextEx(spriteFont SpriteFont, text string, position Vector2, fontSize float32, spacing int32, tint Color) {
	cspriteFont := spriteFont.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cposition := position.cptr()
	cfontSize := (C.float)(fontSize)
	cspacing := (C.int)(spacing)
	ctint := tint.cptr()
	C.DrawTextEx(*cspriteFont, ctext, *cposition, cfontSize, cspacing, *ctint)
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

// MeasureTextEx - Measure string size for SpriteFont
func MeasureTextEx(spriteFont SpriteFont, text string, fontSize float32, spacing int32) Vector2 {
	cspriteFont := spriteFont.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.float)(fontSize)
	cspacing := (C.int)(spacing)
	ret := C.MeasureTextEx(*cspriteFont, ctext, cfontSize, cspacing)
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// DrawFPS - Shows current FPS
func DrawFPS(posX int32, posY int32) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	C.DrawFPS(cposX, cposY)
}
