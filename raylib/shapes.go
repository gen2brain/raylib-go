package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// DrawPixel - Draw a pixel
func DrawPixel(posX, posY int32, color Color) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	ccolor := color.cptr()
	C.DrawPixel(cposX, cposY, *ccolor)
}

// DrawPixelV - Draw a pixel (Vector version)
func DrawPixelV(position Vector2, color Color) {
	cposition := position.cptr()
	ccolor := color.cptr()
	C.DrawPixelV(*cposition, *ccolor)
}

// DrawLine - Draw a line
func DrawLine(startPosX, startPosY, endPosX, endPosY int32, color Color) {
	cstartPosX := (C.int)(startPosX)
	cstartPosY := (C.int)(startPosY)
	cendPosX := (C.int)(endPosX)
	cendPosY := (C.int)(endPosY)
	ccolor := color.cptr()
	C.DrawLine(cstartPosX, cstartPosY, cendPosX, cendPosY, *ccolor)
}

// DrawLineV - Draw a line (Vector version)
func DrawLineV(startPos, endPos Vector2, color Color) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	ccolor := color.cptr()
	C.DrawLineV(*cstartPos, *cendPos, *ccolor)
}

// DrawLineEx - Draw a line defining thickness
func DrawLineEx(startPos, endPos Vector2, thick float32, color Color) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	cthick := (C.float)(thick)
	ccolor := color.cptr()
	C.DrawLineEx(*cstartPos, *cendPos, cthick, *ccolor)
}

// DrawLineBezier - Draw a line using cubic-bezier curves in-out
func DrawLineBezier(startPos, endPos Vector2, thick float32, color Color) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	cthick := (C.float)(thick)
	ccolor := color.cptr()
	C.DrawLineBezier(*cstartPos, *cendPos, cthick, *ccolor)
}

// DrawCircle - Draw a color-filled circle
func DrawCircle(centerX, centerY int32, radius float32, color Color) {
	ccenterX := (C.int)(centerX)
	ccenterY := (C.int)(centerY)
	cradius := (C.float)(radius)
	ccolor := color.cptr()
	C.DrawCircle(ccenterX, ccenterY, cradius, *ccolor)
}

// DrawCircleGradient - Draw a gradient-filled circle
func DrawCircleGradient(centerX, centerY int32, radius float32, color1, color2 Color) {
	ccenterX := (C.int)(centerX)
	ccenterY := (C.int)(centerY)
	cradius := (C.float)(radius)
	ccolor1 := color1.cptr()
	ccolor2 := color2.cptr()
	C.DrawCircleGradient(ccenterX, ccenterY, cradius, *ccolor1, *ccolor2)
}

// DrawCircleV - Draw a color-filled circle (Vector version)
func DrawCircleV(center Vector2, radius float32, color Color) {
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	ccolor := color.cptr()
	C.DrawCircleV(*ccenter, cradius, *ccolor)
}

// DrawCircleLines - Draw circle outline
func DrawCircleLines(centerX, centerY int32, radius float32, color Color) {
	ccenterX := (C.int)(centerX)
	ccenterY := (C.int)(centerY)
	cradius := (C.float)(radius)
	ccolor := color.cptr()
	C.DrawCircleLines(ccenterX, ccenterY, cradius, *ccolor)
}

// DrawRectangle - Draw a color-filled rectangle
func DrawRectangle(posX, posY, width, height int32, color Color) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor := color.cptr()
	C.DrawRectangle(cposX, cposY, cwidth, cheight, *ccolor)
}

// DrawRectangleRec - Draw a color-filled rectangle
func DrawRectangleRec(rec Rectangle, color Color) {
	crec := rec.cptr()
	ccolor := color.cptr()
	C.DrawRectangleRec(*crec, *ccolor)
}

// DrawRectanglePro - Draw a color-filled rectangle with pro parameters
func DrawRectanglePro(rec Rectangle, origin Vector2, rotation float32, color Color) {
	crec := rec.cptr()
	corigin := origin.cptr()
	crotation := (C.float)(rotation)
	ccolor := color.cptr()
	C.DrawRectanglePro(*crec, *corigin, crotation, *ccolor)
}

// DrawRectangleGradientV - Draw a vertical-gradient-filled rectangle
func DrawRectangleGradientV(posX, posY, width, height int32, color1, color2 Color) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor1 := color1.cptr()
	ccolor2 := color2.cptr()
	C.DrawRectangleGradientV(cposX, cposY, cwidth, cheight, *ccolor1, *ccolor2)
}

// DrawRectangleGradientH - Draw a horizontal-gradient-filled rectangle
func DrawRectangleGradientH(posX, posY, width, height int32, color1, color2 Color) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor1 := color1.cptr()
	ccolor2 := color2.cptr()
	C.DrawRectangleGradientH(cposX, cposY, cwidth, cheight, *ccolor1, *ccolor2)
}

// DrawRectangleGradientEx - Draw a gradient-filled rectangle with custom vertex colors
func DrawRectangleGradientEx(rec Rectangle, color1, color2, color3, color4 Color) {
	crec := rec.cptr()
	ccolor1 := color1.cptr()
	ccolor2 := color2.cptr()
	ccolor3 := color3.cptr()
	ccolor4 := color4.cptr()
	C.DrawRectangleGradientEx(*crec, *ccolor1, *ccolor2, *ccolor3, *ccolor4)
}

// DrawRectangleV - Draw a color-filled rectangle (Vector version)
func DrawRectangleV(position Vector2, size Vector2, color Color) {
	cposition := position.cptr()
	csize := size.cptr()
	ccolor := color.cptr()
	C.DrawRectangleV(*cposition, *csize, *ccolor)
}

// DrawRectangleLines - Draw rectangle outline
func DrawRectangleLines(posX, posY, width, height int32, color Color) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor := color.cptr()
	C.DrawRectangleLines(cposX, cposY, cwidth, cheight, *ccolor)
}

// DrawRectangleLinesEx - Draw rectangle outline with extended parameters
func DrawRectangleLinesEx(rec Rectangle, lineThick int32, color Color) {
	crec := rec.cptr()
	clineThick := (C.int)(lineThick)
	ccolor := color.cptr()
	C.DrawRectangleLinesEx(*crec, clineThick, *ccolor)
}

// DrawTriangle - Draw a color-filled triangle
func DrawTriangle(v1, v2, v3 Vector2, color Color) {
	cv1 := v1.cptr()
	cv2 := v2.cptr()
	cv3 := v3.cptr()
	ccolor := color.cptr()
	C.DrawTriangle(*cv1, *cv2, *cv3, *ccolor)
}

// DrawTriangleLines - Draw triangle outline
func DrawTriangleLines(v1, v2, v3 Vector2, color Color) {
	cv1 := v1.cptr()
	cv2 := v2.cptr()
	cv3 := v3.cptr()
	ccolor := color.cptr()
	C.DrawTriangleLines(*cv1, *cv2, *cv3, *ccolor)
}

// DrawPoly - Draw a regular polygon (Vector version)
func DrawPoly(center Vector2, sides int32, radius, rotation float32, color Color) {
	ccenter := center.cptr()
	csides := (C.int)(sides)
	cradius := (C.float)(radius)
	crotation := (C.float)(rotation)
	ccolor := color.cptr()
	C.DrawPoly(*ccenter, csides, cradius, crotation, *ccolor)
}

// DrawPolyEx - Draw a closed polygon defined by points
func DrawPolyEx(points []Vector2, numPoints int32, color Color) {
	cpoints := points[0].cptr()
	cnumPoints := (C.int)(numPoints)
	ccolor := color.cptr()
	C.DrawPolyEx(cpoints, cnumPoints, *ccolor)
}

// DrawPolyExLines - Draw polygon lines
func DrawPolyExLines(points []Vector2, numPoints int32, color Color) {
	cpoints := points[0].cptr()
	cnumPoints := (C.int)(numPoints)
	ccolor := color.cptr()
	C.DrawPolyExLines(cpoints, cnumPoints, *ccolor)
}

// CheckCollisionRecs - Check collision between two rectangles
func CheckCollisionRecs(rec1, rec2 Rectangle) bool {
	crec1 := rec1.cptr()
	crec2 := rec2.cptr()
	ret := C.CheckCollisionRecs(*crec1, *crec2)
	v := bool(int(ret) == 1)
	return v
}

// CheckCollisionCircles - Check collision between two circles
func CheckCollisionCircles(center1 Vector2, radius1 float32, center2 Vector2, radius2 float32) bool {
	ccenter1 := center1.cptr()
	cradius1 := (C.float)(radius1)
	ccenter2 := center2.cptr()
	cradius2 := (C.float)(radius2)
	ret := C.CheckCollisionCircles(*ccenter1, cradius1, *ccenter2, cradius2)
	v := bool(int(ret) == 1)
	return v
}

// CheckCollisionCircleRec - Check collision between circle and rectangle
func CheckCollisionCircleRec(center Vector2, radius float32, rec Rectangle) bool {
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	crec := rec.cptr()
	ret := C.CheckCollisionCircleRec(*ccenter, cradius, *crec)
	v := bool(int(ret) == 1)
	return v
}

// GetCollisionRec - Get collision rectangle for two rectangles collision
func GetCollisionRec(rec1, rec2 Rectangle) Rectangle {
	crec1 := rec1.cptr()
	crec2 := rec2.cptr()
	ret := C.GetCollisionRec(*crec1, *crec2)
	v := newRectangleFromPointer(unsafe.Pointer(&ret))
	return v
}

// CheckCollisionPointRec - Check if point is inside rectangle
func CheckCollisionPointRec(point Vector2, rec Rectangle) bool {
	cpoint := point.cptr()
	crec := rec.cptr()
	ret := C.CheckCollisionPointRec(*cpoint, *crec)
	v := bool(int(ret) == 1)
	return v
}

// CheckCollisionPointCircle - Check if point is inside circle
func CheckCollisionPointCircle(point Vector2, center Vector2, radius float32) bool {
	cpoint := point.cptr()
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	ret := C.CheckCollisionPointCircle(*cpoint, *ccenter, cradius)
	v := bool(int(ret) == 1)
	return v
}

// CheckCollisionPointTriangle - Check if point is inside a triangle
func CheckCollisionPointTriangle(point, p1, p2, p3 Vector2) bool {
	cpoint := point.cptr()
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	ret := C.CheckCollisionPointTriangle(*cpoint, *cp1, *cp2, *cp3)
	v := bool(int(ret) == 1)
	return v
}
