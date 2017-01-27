package raylib

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// SpriteFont type, includes texture and charSet array data
type SpriteFont struct {
	// Font texture
	Texture Texture2D
	// Base size (default chars height)
	Size int32
	// Number of characters
	NumChars int32
	// Characters values array
	CharValues []int32
	// Characters rectangles within the texture
	CharRecs []Rectangle
	// Characters offsets (on drawing)
	CharOffsets []Vector2
	// Characters x advance (on drawing)
	CharAdvanceX []int32
}

func (s *SpriteFont) cptr() *C.SpriteFont {
	return (*C.SpriteFont)(unsafe.Pointer(s))
}

// Returns new SpriteFont
func NewSpriteFont(texture Texture2D, size, numChars int32, charValues []int32, charRecs []Rectangle, charOffsets []Vector2, charAdvanceX []int32) SpriteFont {
	return SpriteFont{texture, size, numChars, charValues, charRecs, charOffsets, charAdvanceX}
}

// Returns new SpriteFont from pointer
func NewSpriteFontFromPointer(ptr unsafe.Pointer) SpriteFont {
	return *(*SpriteFont)(ptr)
}

// Get the default SpriteFont
func GetDefaultFont() SpriteFont {
	ret := C.GetDefaultFont()
	v := NewSpriteFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// Load a SpriteFont image into GPU memory
func LoadSpriteFont(fileName string) SpriteFont {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadSpriteFont(cfileName)
	v := NewSpriteFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// Load a SpriteFont from TTF font with parameters
func LoadSpriteFontTTF(fileName string, fontSize int32, numChars int32, fontChars *int32) SpriteFont {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cfontSize := (C.int)(fontSize)
	cnumChars := (C.int)(numChars)
	cfontChars := (*C.int)(unsafe.Pointer(fontChars))
	ret := C.LoadSpriteFontTTF(cfileName, cfontSize, cnumChars, cfontChars)
	v := NewSpriteFontFromPointer(unsafe.Pointer(&ret))
	return v
}

// Unload SpriteFont from GPU memory
func UnloadSpriteFont(spriteFont SpriteFont) {
	cspriteFont := spriteFont.cptr()
	C.UnloadSpriteFont(*cspriteFont)
}

// Draw text (using default font)
func DrawText(text string, posX int32, posY int32, fontSize int32, color Color) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cfontSize := (C.int)(fontSize)
	ccolor := color.cptr()
	C.DrawText(ctext, cposX, cposY, cfontSize, *ccolor)
}

// Draw text using SpriteFont and additional parameters
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

// Measure string width for default font
func MeasureText(text string, fontSize int32) int32 {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.int)(fontSize)
	ret := C.MeasureText(ctext, cfontSize)
	v := (int32)(ret)
	return v
}

// Measure string size for SpriteFont
func MeasureTextEx(spriteFont SpriteFont, text string, fontSize float32, spacing int32) Vector2 {
	cspriteFont := spriteFont.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.float)(fontSize)
	cspacing := (C.int)(spacing)
	ret := C.MeasureTextEx(*cspriteFont, ctext, cfontSize, cspacing)
	v := NewVector2FromPointer(unsafe.Pointer(&ret))
	return v
}
