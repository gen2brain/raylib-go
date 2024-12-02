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

// SetShapesTexture - Define default texture used to draw shapes
func SetShapesTexture(texture Texture2D, source Rectangle) {
	ctexture := texture.cptr()
	csource := source.cptr()
	C.SetShapesTexture(*ctexture, *csource)
}

// GetShapesTexture - Get texture that is used for shapes drawing
func GetShapesTexture() Texture2D {
	ret := C.GetShapesTexture()
	v := newTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetShapesTextureRectangle - Get texture source rectangle that is used for shapes drawing
func GetShapesTextureRectangle() Rectangle {
	ret := C.GetShapesTextureRectangle()
	v := newRectangleFromPointer(unsafe.Pointer(&ret))
	return v
}

// DrawPixel - Draw a pixel
func DrawPixel(posX, posY int32, col color.RGBA) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	ccolor := colorCptr(col)
	C.DrawPixel(cposX, cposY, *ccolor)
}

// DrawPixelV - Draw a pixel (Vector version)
func DrawPixelV(position Vector2, col color.RGBA) {
	cposition := position.cptr()
	ccolor := colorCptr(col)
	C.DrawPixelV(*cposition, *ccolor)
}

// DrawLine - Draw a line
func DrawLine(startPosX, startPosY, endPosX, endPosY int32, col color.RGBA) {
	cstartPosX := (C.int)(startPosX)
	cstartPosY := (C.int)(startPosY)
	cendPosX := (C.int)(endPosX)
	cendPosY := (C.int)(endPosY)
	ccolor := colorCptr(col)
	C.DrawLine(cstartPosX, cstartPosY, cendPosX, cendPosY, *ccolor)
}

// DrawLineV - Draw a line (Vector version)
func DrawLineV(startPos, endPos Vector2, col color.RGBA) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	ccolor := colorCptr(col)
	C.DrawLineV(*cstartPos, *cendPos, *ccolor)
}

// DrawLineEx - Draw a line defining thickness
func DrawLineEx(startPos, endPos Vector2, thick float32, col color.RGBA) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawLineEx(*cstartPos, *cendPos, cthick, *ccolor)
}

// DrawLineStrip - Draw lines sequence
func DrawLineStrip(points []Vector2, col color.RGBA) {
	cpoints := (*C.Vector2)(unsafe.Pointer(&points[0]))
	cpointCount := (C.int)(len(points))
	ccolor := colorCptr(col)
	C.DrawLineStrip(cpoints, cpointCount, *ccolor)
}

// DrawLineBezier - Draw a line using cubic-bezier curves in-out
func DrawLineBezier(startPos, endPos Vector2, thick float32, col color.RGBA) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawLineBezier(*cstartPos, *cendPos, cthick, *ccolor)
}

// DrawCircle - Draw a color-filled circle
func DrawCircle(centerX, centerY int32, radius float32, col color.RGBA) {
	ccenterX := (C.int)(centerX)
	ccenterY := (C.int)(centerY)
	cradius := (C.float)(radius)
	ccolor := colorCptr(col)
	C.DrawCircle(ccenterX, ccenterY, cradius, *ccolor)
}

// DrawCircleSector - Draw a piece of a circle
func DrawCircleSector(center Vector2, radius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	cstartAngle := (C.float)(startAngle)
	cendAngle := (C.float)(endAngle)
	csegments := (C.int)(segments)
	ccolor := colorCptr(col)
	C.DrawCircleSector(*ccenter, cradius, cstartAngle, cendAngle, csegments, *ccolor)
}

// DrawCircleSectorLines -
func DrawCircleSectorLines(center Vector2, radius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	cstartAngle := (C.float)(startAngle)
	cendAngle := (C.float)(endAngle)
	csegments := (C.int)(segments)
	ccolor := colorCptr(col)
	C.DrawCircleSectorLines(*ccenter, cradius, cstartAngle, cendAngle, csegments, *ccolor)
}

// DrawCircleGradient - Draw a gradient-filled circle
func DrawCircleGradient(centerX, centerY int32, radius float32, inner, outer color.RGBA) {
	ccenterX := (C.int)(centerX)
	ccenterY := (C.int)(centerY)
	cradius := (C.float)(radius)
	cinner := colorCptr(inner)
	couter := colorCptr(outer)
	C.DrawCircleGradient(ccenterX, ccenterY, cradius, *cinner, *couter)
}

// DrawCircleV - Draw a color-filled circle (Vector version)
func DrawCircleV(center Vector2, radius float32, col color.RGBA) {
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	ccolor := colorCptr(col)
	C.DrawCircleV(*ccenter, cradius, *ccolor)
}

// DrawCircleLines - Draw circle outline
func DrawCircleLines(centerX, centerY int32, radius float32, col color.RGBA) {
	ccenterX := (C.int)(centerX)
	ccenterY := (C.int)(centerY)
	cradius := (C.float)(radius)
	ccolor := colorCptr(col)
	C.DrawCircleLines(ccenterX, ccenterY, cradius, *ccolor)
}

// DrawCircleLinesV - Draw circle outline (Vector version)
func DrawCircleLinesV(center Vector2, radius float32, col color.RGBA) {
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	ccolor := colorCptr(col)
	C.DrawCircleLinesV(*ccenter, cradius, *ccolor)
}

// DrawEllipse - Draw ellipse
func DrawEllipse(centerX, centerY int32, radiusH, radiusV float32, col color.RGBA) {
	ccenterX := (C.int)(centerX)
	ccenterY := (C.int)(centerY)
	cradiusH := (C.float)(radiusH)
	cradiusV := (C.float)(radiusV)
	ccolor := colorCptr(col)
	C.DrawEllipse(ccenterX, ccenterY, cradiusH, cradiusV, *ccolor)
}

// DrawEllipseLines - Draw ellipse outline
func DrawEllipseLines(centerX, centerY int32, radiusH, radiusV float32, col color.RGBA) {
	ccenterX := (C.int)(centerX)
	ccenterY := (C.int)(centerY)
	cradiusH := (C.float)(radiusH)
	cradiusV := (C.float)(radiusV)
	ccolor := colorCptr(col)
	C.DrawEllipseLines(ccenterX, ccenterY, cradiusH, cradiusV, *ccolor)
}

// DrawRing - Draw ring
func DrawRing(center Vector2, innerRadius, outerRadius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	ccenter := center.cptr()
	cinnerRadius := (C.float)(innerRadius)
	couterRadius := (C.float)(outerRadius)
	cstartAngle := (C.float)(startAngle)
	cendAngle := (C.float)(endAngle)
	csegments := (C.int)(segments)
	ccolor := colorCptr(col)
	C.DrawRing(*ccenter, cinnerRadius, couterRadius, cstartAngle, cendAngle, csegments, *ccolor)
}

// DrawRingLines - Draw ring outline
func DrawRingLines(center Vector2, innerRadius, outerRadius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	ccenter := center.cptr()
	cinnerRadius := (C.float)(innerRadius)
	couterRadius := (C.float)(outerRadius)
	cstartAngle := (C.float)(startAngle)
	cendAngle := (C.float)(endAngle)
	csegments := (C.int)(segments)
	ccolor := colorCptr(col)
	C.DrawRingLines(*ccenter, cinnerRadius, couterRadius, cstartAngle, cendAngle, csegments, *ccolor)
}

// DrawRectangle - Draw a color-filled rectangle
func DrawRectangle(posX, posY, width, height int32, col color.RGBA) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor := colorCptr(col)
	C.DrawRectangle(cposX, cposY, cwidth, cheight, *ccolor)
}

// DrawRectangleV - Draw a color-filled rectangle (Vector version)
func DrawRectangleV(position Vector2, size Vector2, col color.RGBA) {
	cposition := position.cptr()
	csize := size.cptr()
	ccolor := colorCptr(col)
	C.DrawRectangleV(*cposition, *csize, *ccolor)
}

// DrawRectangleRec - Draw a color-filled rectangle
func DrawRectangleRec(rec Rectangle, col color.RGBA) {
	crec := rec.cptr()
	ccolor := colorCptr(col)
	C.DrawRectangleRec(*crec, *ccolor)
}

// DrawRectanglePro - Draw a color-filled rectangle with pro parameters
func DrawRectanglePro(rec Rectangle, origin Vector2, rotation float32, col color.RGBA) {
	crec := rec.cptr()
	corigin := origin.cptr()
	crotation := (C.float)(rotation)
	ccolor := colorCptr(col)
	C.DrawRectanglePro(*crec, *corigin, crotation, *ccolor)
}

// DrawRectangleGradientV - Draw a vertical-gradient-filled rectangle
func DrawRectangleGradientV(posX, posY, width, height int32, top, bottom color.RGBA) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ctop := colorCptr(top)
	cbottom := colorCptr(bottom)
	C.DrawRectangleGradientV(cposX, cposY, cwidth, cheight, *ctop, *cbottom)
}

// DrawRectangleGradientH - Draw a horizontal-gradient-filled rectangle
func DrawRectangleGradientH(posX, posY, width, height int32, left, right color.RGBA) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cleft := colorCptr(left)
	cright := colorCptr(right)
	C.DrawRectangleGradientH(cposX, cposY, cwidth, cheight, *cleft, *cright)
}

// DrawRectangleGradientEx - Draw a gradient-filled rectangle with custom vertex colors
func DrawRectangleGradientEx(rec Rectangle, topLeft, bottomLeft, topRight, bottomRight color.RGBA) {
	crec := rec.cptr()
	ctopLeft := colorCptr(topLeft)
	cbottomLeft := colorCptr(bottomLeft)
	ctopRight := colorCptr(topRight)
	cbottomRight := colorCptr(bottomRight)
	C.DrawRectangleGradientEx(*crec, *ctopLeft, *cbottomLeft, *ctopRight, *cbottomRight)
}

// DrawRectangleLines - Draw rectangle outline
func DrawRectangleLines(posX, posY, width, height int32, col color.RGBA) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor := colorCptr(col)
	C.DrawRectangleLines(cposX, cposY, cwidth, cheight, *ccolor)
}

// DrawRectangleLinesEx - Draw rectangle outline with extended parameters
func DrawRectangleLinesEx(rec Rectangle, lineThick float32, col color.RGBA) {
	crec := rec.cptr()
	clineThick := (C.float)(lineThick)
	ccolor := colorCptr(col)
	C.DrawRectangleLinesEx(*crec, clineThick, *ccolor)
}

// DrawRectangleRounded - Draw rectangle with rounded edges
func DrawRectangleRounded(rec Rectangle, roundness float32, segments int32, col color.RGBA) {
	crec := rec.cptr()
	croundness := (C.float)(roundness)
	csegments := (C.int)(segments)
	ccolor := colorCptr(col)
	C.DrawRectangleRounded(*crec, croundness, csegments, *ccolor)
}

// DrawRectangleRoundedLines - Draw rectangle lines with rounded edges
func DrawRectangleRoundedLines(rec Rectangle, roundness float32, segments int32, col color.RGBA) {
	crec := rec.cptr()
	croundness := (C.float)(roundness)
	csegments := (C.int)(segments)
	ccolor := colorCptr(col)
	C.DrawRectangleRoundedLines(*crec, croundness, csegments, *ccolor)
}

// DrawRectangleRoundedLinesEx - Draw rectangle with rounded edges outline
func DrawRectangleRoundedLinesEx(rec Rectangle, roundness float32, segments int32, lineThick float32, col color.RGBA) {
	crec := rec.cptr()
	croundness := (C.float)(roundness)
	csegments := (C.int)(segments)
	clineThick := (C.float)(lineThick)
	ccolor := colorCptr(col)
	C.DrawRectangleRoundedLinesEx(*crec, croundness, csegments, clineThick, *ccolor)
}

// DrawTriangle - Draw a color-filled triangle
func DrawTriangle(v1, v2, v3 Vector2, col color.RGBA) {
	cv1 := v1.cptr()
	cv2 := v2.cptr()
	cv3 := v3.cptr()
	ccolor := colorCptr(col)
	C.DrawTriangle(*cv1, *cv2, *cv3, *ccolor)
}

// DrawTriangleLines - Draw triangle outline
func DrawTriangleLines(v1, v2, v3 Vector2, col color.RGBA) {
	cv1 := v1.cptr()
	cv2 := v2.cptr()
	cv3 := v3.cptr()
	ccolor := colorCptr(col)
	C.DrawTriangleLines(*cv1, *cv2, *cv3, *ccolor)
}

// DrawTriangleFan - Draw a triangle fan defined by points
func DrawTriangleFan(points []Vector2, col color.RGBA) {
	cpoints := (*C.Vector2)(unsafe.Pointer(&points[0]))
	cpointsCount := (C.int)(len(points))
	ccolor := colorCptr(col)
	C.DrawTriangleFan(cpoints, cpointsCount, *ccolor)
}

// DrawTriangleStrip - Draw a triangle strip defined by points
func DrawTriangleStrip(points []Vector2, col color.RGBA) {
	cpoints := (*C.Vector2)(unsafe.Pointer(&points[0]))
	cpointsCount := (C.int)(int32(len(points)))
	ccolor := colorCptr(col)
	C.DrawTriangleStrip(cpoints, cpointsCount, *ccolor)
}

// DrawPoly - Draw a regular polygon (Vector version)
func DrawPoly(center Vector2, sides int32, radius, rotation float32, col color.RGBA) {
	ccenter := center.cptr()
	csides := (C.int)(sides)
	cradius := (C.float)(radius)
	crotation := (C.float)(rotation)
	ccolor := colorCptr(col)
	C.DrawPoly(*ccenter, csides, cradius, crotation, *ccolor)
}

// DrawPolyLines - Draw a polygon outline of n sides
func DrawPolyLines(center Vector2, sides int32, radius, rotation float32, col color.RGBA) {
	ccenter := center.cptr()
	csides := (C.int)(sides)
	cradius := (C.float)(radius)
	crotation := (C.float)(rotation)
	ccolor := colorCptr(col)
	C.DrawPolyLines(*ccenter, csides, cradius, crotation, *ccolor)
}

// DrawPolyLinesEx - Draw a polygon outline of n sides with extended parameters
func DrawPolyLinesEx(center Vector2, sides int32, radius float32, rotation float32, lineThick float32, col color.RGBA) {
	ccenter := center.cptr()
	csides := (C.int)(sides)
	cradius := (C.float)(radius)
	crotation := (C.float)(rotation)
	clineThick := (C.float)(lineThick)
	ccolor := colorCptr(col)
	C.DrawPolyLinesEx(*ccenter, csides, cradius, crotation, clineThick, *ccolor)
}

// DrawSplineLinear - Draw spline: Linear, minimum 2 points
func DrawSplineLinear(points []Vector2, thick float32, col color.RGBA) {
	cpoints := (*C.Vector2)(unsafe.Pointer(&points[0]))
	cpointCount := (C.int)(len(points))
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawSplineLinear(cpoints, cpointCount, cthick, *ccolor)
}

// DrawSplineBasis - Draw spline: B-Spline, minimum 4 points
func DrawSplineBasis(points []Vector2, thick float32, col color.RGBA) {
	cpoints := (*C.Vector2)(unsafe.Pointer(&points[0]))
	cpointCount := (C.int)(len(points))
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawSplineBasis(cpoints, cpointCount, cthick, *ccolor)
}

// DrawSplineCatmullRom - Draw spline: Catmull-Rom, minimum 4 points
func DrawSplineCatmullRom(points []Vector2, thick float32, col color.RGBA) {
	cpoints := (*C.Vector2)(unsafe.Pointer(&points[0]))
	cpointCount := (C.int)(len(points))
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawSplineCatmullRom(cpoints, cpointCount, cthick, *ccolor)
}

// DrawSplineBezierQuadratic - Draw spline: Quadratic Bezier, minimum 3 points (1 control point): [p1, c2, p3, c4...]
func DrawSplineBezierQuadratic(points []Vector2, thick float32, col color.RGBA) {
	cpoints := (*C.Vector2)(unsafe.Pointer(&points[0]))
	cpointCount := (C.int)(len(points))
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawSplineBezierQuadratic(cpoints, cpointCount, cthick, *ccolor)
}

// DrawSplineBezierCubic - Draw spline: Cubic Bezier, minimum 4 points (2 control points): [p1, c2, c3, p4, c5, c6...]
func DrawSplineBezierCubic(points []Vector2, thick float32, col color.RGBA) {
	cpoints := (*C.Vector2)(unsafe.Pointer(&points[0]))
	cpointCount := (C.int)(len(points))
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawSplineBezierCubic(cpoints, cpointCount, cthick, *ccolor)
}

// DrawSplineSegmentLinear - Draw spline segment: Linear, 2 points
func DrawSplineSegmentLinear(p1, p2 Vector2, thick float32, col color.RGBA) {
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawSplineSegmentLinear(*cp1, *cp2, cthick, *ccolor)
}

// DrawSplineSegmentBasis - Draw spline segment: B-Spline, 4 points
func DrawSplineSegmentBasis(p1, p2, p3, p4 Vector2, thick float32, col color.RGBA) {
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	cp4 := p4.cptr()
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawSplineSegmentBasis(*cp1, *cp2, *cp3, *cp4, cthick, *ccolor)
}

// DrawSplineSegmentCatmullRom - Draw spline segment: Catmull-Rom, 4 points
func DrawSplineSegmentCatmullRom(p1, p2, p3, p4 Vector2, thick float32, col color.RGBA) {
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	cp4 := p4.cptr()
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawSplineSegmentCatmullRom(*cp1, *cp2, *cp3, *cp4, cthick, *ccolor)
}

// DrawSplineSegmentBezierQuadratic - Draw spline segment: Quadratic Bezier, 2 points, 1 control point
func DrawSplineSegmentBezierQuadratic(p1, p2, p3 Vector2, thick float32, col color.RGBA) {
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawSplineSegmentBezierQuadratic(*cp1, *cp2, *cp3, cthick, *ccolor)
}

// DrawSplineSegmentBezierCubic - Draw spline segment: Cubic Bezier, 2 points, 2 control points
func DrawSplineSegmentBezierCubic(p1, p2, p3, p4 Vector2, thick float32, col color.RGBA) {
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	cp4 := p4.cptr()
	cthick := (C.float)(thick)
	ccolor := colorCptr(col)
	C.DrawSplineSegmentBezierCubic(*cp1, *cp2, *cp3, *cp4, cthick, *ccolor)
}

// GetSplinePointLinear - Get (evaluate) spline point: Linear
func GetSplinePointLinear(p1, p2 Vector2, t float32) Vector2 {
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	ct := (C.float)(t)
	ret := C.GetSplinePointLinear(*cp1, *cp2, ct)
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// GetSplinePointBasis - Get (evaluate) spline point: B-Spline
func GetSplinePointBasis(p1, p2, p3, p4 Vector2, t float32) Vector2 {
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	cp4 := p4.cptr()
	ct := (C.float)(t)
	ret := C.GetSplinePointBasis(*cp1, *cp2, *cp3, *cp4, ct)
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// GetSplinePointCatmullRom - Get (evaluate) spline point: Catmull-Rom
func GetSplinePointCatmullRom(p1, p2, p3, p4 Vector2, t float32) Vector2 {
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	cp4 := p4.cptr()
	ct := (C.float)(t)
	ret := C.GetSplinePointCatmullRom(*cp1, *cp2, *cp3, *cp4, ct)
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// GetSplinePointBezierQuad - Get (evaluate) spline point: Quadratic Bezier
func GetSplinePointBezierQuad(p1, p2, p3 Vector2, t float32) Vector2 {
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	ct := (C.float)(t)
	ret := C.GetSplinePointBezierQuad(*cp1, *cp2, *cp3, ct)
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// GetSplinePointBezierCubic - Get (evaluate) spline point: Cubic Bezier
func GetSplinePointBezierCubic(p1, p2, p3, p4 Vector2, t float32) Vector2 {
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	cp4 := p4.cptr()
	ct := (C.float)(t)
	ret := C.GetSplinePointBezierCubic(*cp1, *cp2, *cp3, *cp4, ct)
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// CheckCollisionRecs - Check collision between two rectangles
func CheckCollisionRecs(rec1, rec2 Rectangle) bool {
	crec1 := rec1.cptr()
	crec2 := rec2.cptr()
	ret := C.CheckCollisionRecs(*crec1, *crec2)
	v := bool(ret)
	return v
}

// CheckCollisionCircles - Check collision between two circles
func CheckCollisionCircles(center1 Vector2, radius1 float32, center2 Vector2, radius2 float32) bool {
	ccenter1 := center1.cptr()
	cradius1 := (C.float)(radius1)
	ccenter2 := center2.cptr()
	cradius2 := (C.float)(radius2)
	ret := C.CheckCollisionCircles(*ccenter1, cradius1, *ccenter2, cradius2)
	v := bool(ret)
	return v
}

// CheckCollisionCircleRec - Check collision between circle and rectangle
func CheckCollisionCircleRec(center Vector2, radius float32, rec Rectangle) bool {
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	crec := rec.cptr()
	ret := C.CheckCollisionCircleRec(*ccenter, cradius, *crec)
	v := bool(ret)
	return v
}

// CheckCollisionCircleLine - Check if circle collides with a line created betweeen two points [p1] and [p2]
func CheckCollisionCircleLine(center Vector2, radius float32, p1, p2 Vector2) bool {
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	ret := C.CheckCollisionCircleLine(*ccenter, cradius, *cp1, *cp2)
	v := bool(ret)
	return v
}

// CheckCollisionPointRec - Check if point is inside rectangle
func CheckCollisionPointRec(point Vector2, rec Rectangle) bool {
	cpoint := point.cptr()
	crec := rec.cptr()
	ret := C.CheckCollisionPointRec(*cpoint, *crec)
	v := bool(ret)
	return v
}

// CheckCollisionPointCircle - Check if point is inside circle
func CheckCollisionPointCircle(point Vector2, center Vector2, radius float32) bool {
	cpoint := point.cptr()
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	ret := C.CheckCollisionPointCircle(*cpoint, *ccenter, cradius)
	v := bool(ret)
	return v
}

// CheckCollisionPointTriangle - Check if point is inside a triangle
func CheckCollisionPointTriangle(point, p1, p2, p3 Vector2) bool {
	cpoint := point.cptr()
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	ret := C.CheckCollisionPointTriangle(*cpoint, *cp1, *cp2, *cp3)
	v := bool(ret)
	return v
}

// CheckCollisionPointPoly - Check if point is within a polygon described by array of vertices
//
// NOTE: Based on http://jeffreythompson.org/collision-detection/poly-point.php
func CheckCollisionPointPoly(point Vector2, points []Vector2) bool {
	cpoint := point.cptr()
	cpoints := (&points[0]).cptr()
	cpointCount := C.int(len(points))
	ret := C.CheckCollisionPointPoly(*cpoint, cpoints, cpointCount)
	v := bool(ret)
	return v
}

// CheckCollisionLines - Check the collision between two lines defined by two points each, returns collision point by reference
func CheckCollisionLines(startPos1, endPos1, startPos2, endPos2 Vector2, point *Vector2) bool {
	cstartPos1 := startPos1.cptr()
	cendPos1 := endPos1.cptr()
	cstartPos2 := startPos2.cptr()
	cendPos2 := endPos2.cptr()
	cpoint := point.cptr()
	ret := C.CheckCollisionLines(*cstartPos1, *cendPos1, *cstartPos2, *cendPos2, cpoint)
	v := bool(ret)
	return v
}

// CheckCollisionPointLine - Check if point belongs to line created between two points [p1] and [p2] with defined margin in pixels [threshold]
func CheckCollisionPointLine(point, p1, p2 Vector2, threshold int32) bool {
	cpoint := point.cptr()
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cthreshold := (C.int)(threshold)
	ret := C.CheckCollisionPointLine(*cpoint, *cp1, *cp2, cthreshold)
	v := bool(ret)
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
