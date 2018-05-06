package raylib

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

// GetDefaultFont - Get the default Font
func GetDefaultFont() Font {
	ret := C.GetDefaultFont()
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
func LoadFontEx(fileName string, fontSize int32, charsCount int32, fontChars *int32) Font {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cfontSize := (C.int)(fontSize)
	ccharsCount := (C.int)(charsCount)
	cfontChars := (*C.int)(unsafe.Pointer(fontChars))
	ret := C.LoadFontEx(cfileName, cfontSize, ccharsCount, cfontChars)
	v := newFontFromPointer(unsafe.Pointer(&ret))
	return v
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

// DrawFPS - Shows current FPS
func DrawFPS(posX int32, posY int32) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	C.DrawFPS(cposX, cposY)
}
