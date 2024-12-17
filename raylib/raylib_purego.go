//go:build !cgo && windows
// +build !cgo,windows

package rl

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"reflect"
	"unsafe"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/windows"
)

var (
	// raylibDll is the pointer to the shared library
	raylibDll uintptr

	// audioCallbacks is needed to have a reference between Go functions (keys) created by the user
	// and C function pointers (values) created by purego.NewCallback
	audioCallbacks map[uintptr]uintptr
)

var initWindow func(width int32, height int32, title string)
var closeWindow func()
var windowShouldClose func() bool
var isWindowReady func() bool
var isWindowFullscreen func() bool
var isWindowHidden func() bool
var isWindowMinimized func() bool
var isWindowMaximized func() bool
var isWindowFocused func() bool
var isWindowResized func() bool
var isWindowState func(flag uint32) bool
var setWindowState func(flags uint32)
var clearWindowState func(flags uint32)
var toggleFullscreen func()
var toggleBorderlessWindowed func()
var maximizeWindow func()
var minimizeWindow func()
var restoreWindow func()
var setWindowIcon func(image uintptr)
var setWindowIcons func(images uintptr, count int32)
var setWindowTitle func(title string)
var setWindowPosition func(x int32, y int32)
var setWindowMonitor func(monitor int32)
var setWindowMinSize func(width int32, height int32)
var setWindowMaxSize func(width int32, height int32)
var setWindowSize func(width int32, height int32)
var setWindowOpacity func(opacity float32)
var setWindowFocused func()
var getWindowHandle func() unsafe.Pointer
var getScreenWidth func() int32
var getScreenHeight func() int32
var getRenderWidth func() int32
var getRenderHeight func() int32
var getMonitorCount func() int32
var getCurrentMonitor func() int32
var getMonitorPosition func(monitor int32) uintptr
var getMonitorWidth func(monitor int32) int32
var getMonitorHeight func(monitor int32) int32
var getMonitorPhysicalWidth func(monitor int32) int32
var getMonitorPhysicalHeight func(monitor int32) int32
var getMonitorRefreshRate func(monitor int32) int32
var getWindowPosition func() uintptr
var getWindowScaleDPI func() uintptr
var getMonitorName func(monitor int32) string
var setClipboardText func(text string)
var getClipboardText func() string
var getClipboardImage func(img uintptr)
var enableEventWaiting func()
var disableEventWaiting func()
var showCursor func()
var hideCursor func()
var isCursorHidden func() bool
var enableCursor func()
var disableCursor func()
var isCursorOnScreen func() bool
var clearBackground func(col uintptr)
var beginDrawing func()
var endDrawing func()
var beginMode2D func(camera uintptr)
var endMode2D func()
var beginMode3D func(camera uintptr)
var endMode3D func()
var beginTextureMode func(target uintptr)
var endTextureMode func()
var beginShaderMode func(shader uintptr)
var endShaderMode func()
var beginBlendMode func(mode int32)
var endBlendMode func()
var beginScissorMode func(x int32, y int32, width int32, height int32)
var endScissorMode func()
var beginVrStereoMode func(config uintptr)
var endVrStereoMode func()
var loadVrStereoConfig func(config uintptr, device uintptr)
var unloadVrStereoConfig func(config uintptr)
var loadShader func(shader uintptr, vsFileName uintptr, fsFileName uintptr)
var loadShaderFromMemory func(shader uintptr, vsCode uintptr, fsCode uintptr)
var isShaderValid func(shader uintptr) bool
var getShaderLocation func(shader uintptr, uniformName string) int32
var getShaderLocationAttrib func(shader uintptr, attribName string) int32
var setShaderValue func(shader uintptr, locIndex int32, value []float32, uniformType int32)
var setShaderValueV func(shader uintptr, locIndex int32, value []float32, uniformType int32, count int32)
var setShaderValueMatrix func(shader uintptr, locIndex int32, mat uintptr)
var setShaderValueTexture func(shader uintptr, locIndex int32, texture uintptr)
var unloadShader func(shader uintptr)
var getScreenToWorldRay func(ray uintptr, position uintptr, camera uintptr)
var getScreenToWorldRayEx func(ray uintptr, position uintptr, camera uintptr, width, height int32)
var getCameraMatrix func(mat uintptr, camera uintptr)
var getCameraMatrix2D func(mat uintptr, camera uintptr)
var getWorldToScreen func(position uintptr, camera uintptr) uintptr
var getScreenToWorld2D func(position uintptr, camera uintptr) uintptr
var getWorldToScreenEx func(position uintptr, camera uintptr, width int32, height int32) uintptr
var getWorldToScreen2D func(position uintptr, camera uintptr) uintptr
var setTargetFPS func(fps int32)
var getFrameTime func() float32
var getTime func() float64
var getFPS func() int32
var swapScreenBuffer func()
var pollInputEvents func()
var waitTime func(seconds float64)
var setRandomSeed func(seed uint32)
var getRandomValue func(minimum int32, maximum int32) int32
var loadRandomSequence func(count uint32, minimum int32, maximum int32) *int32
var unloadRandomSequence func(sequence *int32)
var takeScreenshot func(fileName string)
var setConfigFlags func(flags uint32)
var openURL func(url string)
var traceLog func(logLevel int32, text string)
var setTraceLogLevel func(logLevel int32)
var memAlloc func(size uint32) unsafe.Pointer
var memRealloc func(ptr unsafe.Pointer, size uint32) unsafe.Pointer
var memFree func(ptr unsafe.Pointer)
var setTraceLogCallback func(callback uintptr)
var isFileDropped func() bool
var loadDroppedFiles func(files uintptr)
var unloadDroppedFiles func(files uintptr)
var loadAutomationEventList func(automationEventList uintptr, fileName string) uintptr
var unloadAutomationEventList func(list uintptr)
var exportAutomationEventList func(list uintptr, fileName string) bool
var setAutomationEventList func(list uintptr)
var setAutomationEventBaseFrame func(frame int32)
var startAutomationEventRecording func()
var stopAutomationEventRecording func()
var playAutomationEvent func(event uintptr)
var isKeyPressed func(key int32) bool
var isKeyPressedRepeat func(key int32) bool
var isKeyDown func(key int32) bool
var isKeyReleased func(key int32) bool
var isKeyUp func(key int32) bool
var getKeyPressed func() int32
var getCharPressed func() int32
var setExitKey func(key int32)
var isGamepadAvailable func(gamepad int32) bool
var getGamepadName func(gamepad int32) string
var isGamepadButtonPressed func(gamepad int32, button int32) bool
var isGamepadButtonDown func(gamepad int32, button int32) bool
var isGamepadButtonReleased func(gamepad int32, button int32) bool
var isGamepadButtonUp func(gamepad int32, button int32) bool
var getGamepadButtonPressed func() int32
var getGamepadAxisCount func(gamepad int32) int32
var getGamepadAxisMovement func(gamepad int32, axis int32) float32
var setGamepadMappings func(mappings string) int32
var setGamepadVibration func(gamepad int32, leftMotor, rightMotor, duration float32)
var isMouseButtonPressed func(button int32) bool
var isMouseButtonDown func(button int32) bool
var isMouseButtonReleased func(button int32) bool
var isMouseButtonUp func(button int32) bool
var getMouseX func() int32
var getMouseY func() int32
var getMousePosition func() uintptr
var getMouseDelta func() uintptr
var setMousePosition func(x int32, y int32)
var setMouseOffset func(offsetX int32, offsetY int32)
var setMouseScale func(scaleX float32, scaleY float32)
var getMouseWheelMove func() float32
var getMouseWheelMoveV func() uintptr
var setMouseCursor func(cursor int32)
var getTouchX func() int32
var getTouchY func() int32
var getTouchPosition func(index int32) uintptr
var getTouchPointId func(index int32) int32
var getTouchPointCount func() int32
var setGesturesEnabled func(flags uint32)
var isGestureDetected func(gesture uint32) bool
var getGestureDetected func() int32
var getGestureHoldDuration func() float32
var getGestureDragVector func() uintptr
var getGestureDragAngle func() float32
var getGesturePinchVector func() uintptr
var getGesturePinchAngle func() float32
var setShapesTexture func(texture uintptr, source uintptr)
var getShapesTexture func(texture uintptr)
var getShapesTextureRectangle func(rec uintptr)
var drawPixel func(posX int32, posY int32, col uintptr)
var drawPixelV func(position uintptr, col uintptr)
var drawLine func(startPosX int32, startPosY int32, endPosX int32, endPosY int32, col uintptr)
var drawLineV func(startPos uintptr, endPos uintptr, col uintptr)
var drawLineEx func(startPos uintptr, endPos uintptr, thick float32, col uintptr)
var drawLineStrip func(points *Vector2, pointCount int32, col uintptr)
var drawLineBezier func(startPos uintptr, endPos uintptr, thick float32, col uintptr)
var drawCircle func(centerX int32, centerY int32, radius float32, col uintptr)
var drawCircleSector func(center uintptr, radius float32, startAngle float32, endAngle float32, segments int32, col uintptr)
var drawCircleSectorLines func(center uintptr, radius float32, startAngle float32, endAngle float32, segments int32, col uintptr)
var drawCircleGradient func(centerX int32, centerY int32, radius float32, inner uintptr, outer uintptr)
var drawCircleV func(center uintptr, radius float32, col uintptr)
var drawCircleLines func(centerX int32, centerY int32, radius float32, col uintptr)
var drawCircleLinesV func(center uintptr, radius float32, col uintptr)
var drawEllipse func(centerX int32, centerY int32, radiusH float32, radiusV float32, col uintptr)
var drawEllipseLines func(centerX int32, centerY int32, radiusH float32, radiusV float32, col uintptr)
var drawRing func(center uintptr, innerRadius float32, outerRadius float32, startAngle float32, endAngle float32, segments int32, col uintptr)
var drawRingLines func(center uintptr, innerRadius float32, outerRadius float32, startAngle float32, endAngle float32, segments int32, col uintptr)
var drawRectangle func(posX int32, posY int32, width int32, height int32, col uintptr)
var drawRectangleV func(position uintptr, size uintptr, col uintptr)
var drawRectangleRec func(rec uintptr, col uintptr)
var drawRectanglePro func(rec uintptr, origin uintptr, rotation float32, col uintptr)
var drawRectangleGradientV func(posX int32, posY int32, width int32, height int32, top uintptr, bottom uintptr)
var drawRectangleGradientH func(posX int32, posY int32, width int32, height int32, left uintptr, right uintptr)
var drawRectangleGradientEx func(rec uintptr, topLeft uintptr, bottomLeft uintptr, topRight uintptr, bottomRight uintptr)
var drawRectangleLines func(posX int32, posY int32, width int32, height int32, col uintptr)
var drawRectangleLinesEx func(rec uintptr, lineThick float32, col uintptr)
var drawRectangleRounded func(rec uintptr, roundness float32, segments int32, col uintptr)
var drawRectangleRoundedLines func(rec uintptr, roundness float32, segments int32, col uintptr)
var drawRectangleRoundedLinesEx func(rec uintptr, roundness float32, segments int32, lineThick float32, col uintptr)
var drawTriangle func(v1 uintptr, v2 uintptr, v3 uintptr, col uintptr)
var drawTriangleLines func(v1 uintptr, v2 uintptr, v3 uintptr, col uintptr)
var drawTriangleFan func(points *Vector2, pointCount int32, col uintptr)
var drawTriangleStrip func(points *Vector2, pointCount int32, col uintptr)
var drawPoly func(center uintptr, sides int32, radius float32, rotation float32, col uintptr)
var drawPolyLines func(center uintptr, sides int32, radius float32, rotation float32, col uintptr)
var drawPolyLinesEx func(center uintptr, sides int32, radius float32, rotation float32, lineThick float32, col uintptr)
var drawSplineLinear func(points *Vector2, pointCount int32, thick float32, col uintptr)
var drawSplineBasis func(points *Vector2, pointCount int32, thick float32, col uintptr)
var drawSplineCatmullRom func(points *Vector2, pointCount int32, thick float32, col uintptr)
var drawSplineBezierQuadratic func(points *Vector2, pointCount int32, thick float32, col uintptr)
var drawSplineBezierCubic func(points *Vector2, pointCount int32, thick float32, col uintptr)
var drawSplineSegmentLinear func(p1 uintptr, p2 uintptr, thick float32, col uintptr)
var drawSplineSegmentBasis func(p1 uintptr, p2 uintptr, p3 uintptr, p4 uintptr, thick float32, col uintptr)
var drawSplineSegmentCatmullRom func(p1 uintptr, p2 uintptr, p3 uintptr, p4 uintptr, thick float32, col uintptr)
var drawSplineSegmentBezierQuadratic func(p1 uintptr, c2 uintptr, p3 uintptr, thick float32, col uintptr)
var drawSplineSegmentBezierCubic func(p1 uintptr, c2 uintptr, c3 uintptr, p4 uintptr, thick float32, col uintptr)
var getSplinePointLinear func(startPos uintptr, endPos uintptr, t float32) uintptr
var getSplinePointBasis func(p1 uintptr, p2 uintptr, p3 uintptr, p4 uintptr, t float32) uintptr
var getSplinePointCatmullRom func(p1 uintptr, p2 uintptr, p3 uintptr, p4 uintptr, t float32) uintptr
var getSplinePointBezierQuad func(p1 uintptr, c2 uintptr, p3 uintptr, t float32) uintptr
var getSplinePointBezierCubic func(p1 uintptr, c2 uintptr, c3 uintptr, p4 uintptr, t float32) uintptr
var checkCollisionRecs func(rec1 uintptr, rec2 uintptr) bool
var checkCollisionCircles func(center1 uintptr, radius1 float32, center2 uintptr, radius2 float32) bool
var checkCollisionCircleRec func(center uintptr, radius float32, rec uintptr) bool
var checkCollisionCircleLine func(center uintptr, radius float32, p1, p2 uintptr) bool
var checkCollisionPointRec func(point uintptr, rec uintptr) bool
var checkCollisionPointCircle func(point uintptr, center uintptr, radius float32) bool
var checkCollisionPointTriangle func(point uintptr, p1 uintptr, p2 uintptr, p3 uintptr) bool
var checkCollisionPointPoly func(point uintptr, points *Vector2, pointCount int32) bool
var checkCollisionLines func(startPos1 uintptr, endPos1 uintptr, startPos2 uintptr, endPos2 uintptr, collisionPoint *Vector2) bool
var checkCollisionPointLine func(point uintptr, p1 uintptr, p2 uintptr, threshold int32) bool
var getCollisionRec func(rec uintptr, rec1 uintptr, rec2 uintptr)
var loadImage func(img uintptr, fileName string)
var loadImageRaw func(img uintptr, fileName string, width int32, height int32, format int32, headerSize int32)
var loadImageAnim func(img uintptr, fileName string, frames *int32)
var loadImageAnimFromMemory func(img uintptr, fileType string, fileData []byte, dataSize int32, frames *int32)
var loadImageFromMemory func(img uintptr, fileType string, fileData []byte, dataSize int32)
var loadImageFromTexture func(img uintptr, texture uintptr)
var loadImageFromScreen func(img uintptr)
var isImageValid func(image uintptr) bool
var unloadImage func(image uintptr)
var exportImage func(image uintptr, fileName string) bool
var exportImageToMemory func(image uintptr, fileType string, fileSize *int32) *byte
var genImageColor func(image uintptr, width int32, height int32, col uintptr)
var genImageGradientLinear func(image uintptr, width int32, height int32, direction int32, start uintptr, end uintptr)
var genImageGradientRadial func(image uintptr, width int32, height int32, density float32, inner uintptr, outer uintptr)
var genImageGradientSquare func(image uintptr, width int32, height int32, density float32, inner uintptr, outer uintptr)
var genImageChecked func(image uintptr, width int32, height int32, checksX int32, checksY int32, col1 uintptr, col2 uintptr)
var genImageWhiteNoise func(image uintptr, width int32, height int32, factor float32)
var genImagePerlinNoise func(image uintptr, width int32, height int32, offsetX int32, offsetY int32, scale float32)
var genImageCellular func(image uintptr, width int32, height int32, tileSize int32)
var genImageText func(image uintptr, width int32, height int32, text string)
var imageCopy func(retImage uintptr, image uintptr)
var imageFromImage func(retImage uintptr, image uintptr, rec uintptr)
var imageFromChannel func(retImage uintptr, image uintptr, selectedChannel int32)
var imageText func(retImage uintptr, text string, fontSize int32, col uintptr)
var imageTextEx func(retImage uintptr, font uintptr, text string, fontSize float32, spacing float32, tint uintptr)
var imageFormat func(image *Image, newFormat int32)
var imageToPOT func(image *Image, fill uintptr)
var imageCrop func(image *Image, crop uintptr)
var imageAlphaCrop func(image *Image, threshold float32)
var imageAlphaClear func(image *Image, col uintptr, threshold float32)
var imageAlphaMask func(image *Image, alphaMask uintptr)
var imageAlphaPremultiply func(image *Image)
var imageBlurGaussian func(image *Image, blurSize int32)
var imageKernelConvolution func(image *Image, kernel []float32, kernelSize int32)
var imageResize func(image *Image, newWidth int32, newHeight int32)
var imageResizeNN func(image *Image, newWidth int32, newHeight int32)
var imageResizeCanvas func(image *Image, newWidth int32, newHeight int32, offsetX int32, offsetY int32, fill uintptr)
var imageMipmaps func(image *Image)
var imageDither func(image *Image, rBpp int32, gBpp int32, bBpp int32, aBpp int32)
var imageFlipVertical func(image *Image)
var imageFlipHorizontal func(image *Image)
var imageRotate func(image *Image, degrees int32)
var imageRotateCW func(image *Image)
var imageRotateCCW func(image *Image)
var imageColorTint func(image *Image, col uintptr)
var imageColorInvert func(image *Image)
var imageColorGrayscale func(image *Image)
var imageColorContrast func(image *Image, contrast float32)
var imageColorBrightness func(image *Image, brightness int32)
var imageColorReplace func(image *Image, col uintptr, replace uintptr)
var loadImageColors func(image uintptr) *color.RGBA
var loadImagePalette func(image uintptr, maxPaletteSize int32, colorCount *int32) *color.RGBA
var unloadImageColors func(colors *color.RGBA)
var unloadImagePalette func(colors *color.RGBA)
var getImageAlphaBorder func(rec uintptr, image uintptr, threshold float32)
var getImageColor func(image uintptr, x int32, y int32) uintptr
var imageClearBackground func(dst *Image, col uintptr)
var imageDrawPixel func(dst *Image, posX int32, posY int32, col uintptr)
var imageDrawPixelV func(dst *Image, position uintptr, col uintptr)
var imageDrawLine func(dst *Image, startPosX int32, startPosY int32, endPosX int32, endPosY int32, col uintptr)
var imageDrawLineV func(dst *Image, start uintptr, end uintptr, col uintptr)
var imageDrawLineEx func(dst *Image, start uintptr, end uintptr, thick int32, col uintptr)
var imageDrawCircle func(dst *Image, centerX int32, centerY int32, radius int32, col uintptr)
var imageDrawCircleV func(dst *Image, center uintptr, radius int32, col uintptr)
var imageDrawCircleLines func(dst *Image, centerX int32, centerY int32, radius int32, col uintptr)
var imageDrawCircleLinesV func(dst *Image, center uintptr, radius int32, col uintptr)
var imageDrawRectangle func(dst *Image, posX int32, posY int32, width int32, height int32, col uintptr)
var imageDrawRectangleV func(dst *Image, position uintptr, size uintptr, col uintptr)
var imageDrawRectangleRec func(dst *Image, rec uintptr, col uintptr)
var imageDrawRectangleLines func(dst *Image, rec uintptr, thick int32, col uintptr)
var imageDrawTriangle func(dst *Image, v1, v2, v3 uintptr, col uintptr)
var imageDrawTriangleEx func(dst *Image, v1, v2, v3 uintptr, c1, c2, c3 uintptr)
var imageDrawTriangleLines func(dst *Image, v1, v2, v3 uintptr, col uintptr)
var imageDrawTriangleFan func(dst *Image, points *Vector2, pointCount int32, col uintptr)
var imageDrawTriangleStrip func(dst *Image, points *Vector2, pointCount int32, col uintptr)
var imageDraw func(dst *Image, src uintptr, srcRec uintptr, dstRec uintptr, tint uintptr)
var imageDrawText func(dst *Image, text string, posX int32, posY int32, fontSize int32, col uintptr)
var imageDrawTextEx func(dst *Image, font uintptr, text string, position uintptr, fontSize float32, spacing float32, tint uintptr)
var loadTexture func(texture uintptr, fileName string)
var loadTextureFromImage func(texture uintptr, image uintptr)
var loadTextureCubemap func(texture uintptr, image uintptr, layout int32)
var loadRenderTexture func(texture uintptr, width int32, height int32)
var isTextureValid func(texture uintptr) bool
var unloadTexture func(texture uintptr)
var isRenderTextureValid func(target uintptr) bool
var unloadRenderTexture func(target uintptr)
var updateTexture func(texture uintptr, pixels *color.RGBA)
var updateTextureRec func(texture uintptr, rec uintptr, pixels *color.RGBA)
var genTextureMipmaps func(texture *Texture2D)
var setTextureFilter func(texture uintptr, filter int32)
var setTextureWrap func(texture uintptr, wrap int32)
var drawTexture func(texture uintptr, posX int32, posY int32, tint uintptr)
var drawTextureV func(texture uintptr, position uintptr, tint uintptr)
var drawTextureEx func(texture uintptr, position uintptr, rotation float32, scale float32, tint uintptr)
var drawTextureRec func(texture uintptr, source uintptr, position uintptr, tint uintptr)
var drawTexturePro func(texture uintptr, source uintptr, dest uintptr, origin uintptr, rotation float32, tint uintptr)
var drawTextureNPatch func(texture uintptr, nPatchInfo uintptr, dest uintptr, origin uintptr, rotation float32, tint uintptr)
var fade func(col uintptr, alpha float32) uintptr
var colorToInt func(col uintptr) int32
var colorNormalize func(vector4 uintptr, col uintptr)
var colorFromNormalized func(normalized uintptr) uintptr
var colorToHSV func(vector3 uintptr, col uintptr)
var colorFromHSV func(hue float32, saturation float32, value float32) uintptr
var colorTint func(col uintptr, tint uintptr) uintptr
var colorBrightness func(col uintptr, factor float32) uintptr
var colorContrast func(col uintptr, contrast float32) uintptr
var colorAlpha func(col uintptr, alpha float32) uintptr
var colorAlphaBlend func(dst uintptr, src uintptr, tint uintptr) uintptr
var colorLerp func(col1, col2 uintptr, factor float32) uintptr
var getColor func(hexValue uint32) uintptr
var getPixelColor func(srcPtr unsafe.Pointer, format int32) uintptr
var setPixelColor func(dstPtr unsafe.Pointer, col uintptr, format int32)
var getPixelDataSize func(width int32, height int32, format int32) int32
var getFontDefault func(font uintptr)
var loadFont func(font uintptr, fileName string)
var loadFontEx func(font uintptr, fileName string, fontSize int32, codepoints []int32, codepointCount int32)
var loadFontFromImage func(font uintptr, image uintptr, key uintptr, firstChar int32)
var loadFontFromMemory func(font uintptr, fileType string, fileData []byte, dataSize int32, fontSize int32, codepoints []int32, codepointCount int32)
var isFontValid func(font uintptr) bool
var loadFontData func(fileData []byte, dataSize int32, fontSize int32, codepoints []int32, codepointCount int32, _type int32) *GlyphInfo
var genImageFontAtlas func(image uintptr, glyphs *GlyphInfo, glyphRecs []*Rectangle, glyphCount int32, fontSize int32, padding int32, packMethod int32)
var unloadFontData func(glyphs *GlyphInfo, glyphCount int32)
var unloadFont func(font uintptr)
var drawFPS func(posX int32, posY int32)
var drawText func(text string, posX int32, posY int32, fontSize int32, col uintptr)
var drawTextEx func(font uintptr, text string, position uintptr, fontSize float32, spacing float32, tint uintptr)
var drawTextPro func(font uintptr, text string, position uintptr, origin uintptr, rotation float32, fontSize float32, spacing float32, tint uintptr)
var drawTextCodepoint func(font uintptr, codepoint int32, position uintptr, fontSize float32, tint uintptr)
var drawTextCodepoints func(font uintptr, codepoints []int32, codepointCount int32, position uintptr, fontSize float32, spacing float32, tint uintptr)
var setTextLineSpacing func(spacing int32)
var measureText func(text string, fontSize int32) int32
var measureTextEx func(font uintptr, text string, fontSize float32, spacing float32) uintptr
var getGlyphIndex func(font uintptr, codepoint int32) int32
var getGlyphInfo func(glyphInfo uintptr, font uintptr, codepoint int32)
var getGlyphAtlasRec func(rec uintptr, font uintptr, codepoint int32)
var drawLine3D func(startPos uintptr, endPos uintptr, col uintptr)
var drawPoint3D func(position uintptr, col uintptr)
var drawCircle3D func(center uintptr, radius float32, rotationAxis uintptr, rotationAngle float32, col uintptr)
var drawTriangle3D func(v1 uintptr, v2 uintptr, v3 uintptr, col uintptr)
var drawTriangleStrip3D func(points *Vector3, pointCount int32, col uintptr)
var drawCube func(position uintptr, width float32, height float32, length float32, col uintptr)
var drawCubeV func(position uintptr, size uintptr, col uintptr)
var drawCubeWires func(position uintptr, width float32, height float32, length float32, col uintptr)
var drawCubeWiresV func(position uintptr, size uintptr, col uintptr)
var drawSphere func(centerPos uintptr, radius float32, col uintptr)
var drawSphereEx func(centerPos uintptr, radius float32, rings int32, slices int32, col uintptr)
var drawSphereWires func(centerPos uintptr, radius float32, rings int32, slices int32, col uintptr)
var drawCylinder func(position uintptr, radiusTop float32, radiusBottom float32, height float32, slices int32, col uintptr)
var drawCylinderEx func(startPos uintptr, endPos uintptr, startRadius float32, endRadius float32, sides int32, col uintptr)
var drawCylinderWires func(position uintptr, radiusTop float32, radiusBottom float32, height float32, slices int32, col uintptr)
var drawCylinderWiresEx func(startPos uintptr, endPos uintptr, startRadius float32, endRadius float32, sides int32, col uintptr)
var drawCapsule func(startPos uintptr, endPos uintptr, radius float32, slices int32, rings int32, col uintptr)
var drawCapsuleWires func(startPos uintptr, endPos uintptr, radius float32, slices int32, rings int32, col uintptr)
var drawPlane func(centerPos uintptr, size uintptr, col uintptr)
var drawRay func(ray uintptr, col uintptr)
var drawGrid func(slices int32, spacing float32)
var loadModel func(model uintptr, fileName string)
var loadModelFromMesh func(model uintptr, mesh uintptr)
var isModelValid func(model uintptr) bool
var unloadModel func(model uintptr)
var getModelBoundingBox func(boundingBox uintptr, model uintptr)
var drawModel func(model uintptr, position uintptr, scale float32, tint uintptr)
var drawModelEx func(model uintptr, position uintptr, rotationAxis uintptr, rotationAngle float32, scale uintptr, tint uintptr)
var drawModelWires func(model uintptr, position uintptr, scale float32, tint uintptr)
var drawModelWiresEx func(model uintptr, position uintptr, rotationAxis uintptr, rotationAngle float32, scale uintptr, tint uintptr)
var drawModelPoints func(model uintptr, position uintptr, scale float32, tint uintptr)
var drawModelPointsEx func(model uintptr, position uintptr, rotationAxis uintptr, rotationAngle float32, scale uintptr, tint uintptr)
var drawBoundingBox func(box uintptr, col uintptr)
var drawBillboard func(camera uintptr, texture uintptr, position uintptr, scale float32, tint uintptr)
var drawBillboardRec func(camera uintptr, texture uintptr, source uintptr, position uintptr, size uintptr, tint uintptr)
var drawBillboardPro func(camera uintptr, texture uintptr, source uintptr, position uintptr, up uintptr, size uintptr, origin uintptr, rotation float32, tint uintptr)
var uploadMesh func(mesh *Mesh, dynamic bool)
var updateMeshBuffer func(mesh uintptr, index int32, data []byte, dataSize int32, offset int32)
var unloadMesh func(mesh uintptr)
var drawMesh func(mesh uintptr, material uintptr, transform uintptr)
var drawMeshInstanced func(mesh uintptr, material uintptr, transforms []Matrix, instances int32)
var exportMesh func(mesh uintptr, fileName string) bool
var getMeshBoundingBox func(boundingBox uintptr, mesh uintptr)
var genMeshTangents func(mesh *Mesh)
var genMeshPoly func(mesh uintptr, sides int32, radius float32)
var genMeshPlane func(mesh uintptr, width float32, length float32, resX int32, resZ int32)
var genMeshCube func(mesh uintptr, width float32, height float32, length float32)
var genMeshSphere func(mesh uintptr, radius float32, rings int32, slices int32)
var genMeshHemiSphere func(mesh uintptr, radius float32, rings int32, slices int32)
var genMeshCylinder func(mesh uintptr, radius float32, height float32, slices int32)
var genMeshCone func(mesh uintptr, radius float32, height float32, slices int32)
var genMeshTorus func(mesh uintptr, radius float32, size float32, radSeg int32, sides int32)
var genMeshKnot func(mesh uintptr, radius float32, size float32, radSeg int32, sides int32)
var genMeshHeightmap func(mesh uintptr, heightmap uintptr, size uintptr)
var genMeshCubicmap func(mesh uintptr, cubicmap uintptr, cubeSize uintptr)
var loadMaterials func(fileName string, materialCount *int32) *Material
var loadMaterialDefault func(material uintptr)
var isMaterialValid func(material uintptr) bool
var unloadMaterial func(material uintptr)
var setMaterialTexture func(material *Material, mapType int32, texture uintptr)
var setModelMeshMaterial func(model *Model, meshId int32, materialId int32)
var loadModelAnimations func(fileName string, animCount *int32) *ModelAnimation
var updateModelAnimation func(model uintptr, anim uintptr, frame int32)
var updateModelAnimationBones func(model uintptr, anim uintptr, frame int32)
var unloadModelAnimation func(anim uintptr)
var unloadModelAnimations func(animations *ModelAnimation, animCount int32)
var isModelAnimationValid func(model uintptr, anim uintptr) bool
var checkCollisionSpheres func(center1 uintptr, radius1 float32, center2 uintptr, radius2 float32) bool
var checkCollisionBoxes func(box1 uintptr, box2 uintptr) bool
var checkCollisionBoxSphere func(box uintptr, center uintptr, radius float32) bool
var getRayCollisionSphere func(rayCollision uintptr, ray uintptr, center uintptr, radius float32)
var getRayCollisionBox func(rayCollision uintptr, ray uintptr, box uintptr)
var getRayCollisionMesh func(rayCollision uintptr, ray uintptr, mesh uintptr, transform uintptr)
var getRayCollisionTriangle func(rayCollision uintptr, ray uintptr, p1 uintptr, p2 uintptr, p3 uintptr)
var getRayCollisionQuad func(rayCollision uintptr, ray uintptr, p1 uintptr, p2 uintptr, p3 uintptr, p4 uintptr)
var initAudioDevice func()
var closeAudioDevice func()
var isAudioDeviceReady func() bool
var setMasterVolume func(volume float32)
var getMasterVolume func() float32
var loadWave func(wave uintptr, fileName string)
var loadWaveFromMemory func(wave uintptr, fileType string, fileData []byte, dataSize int32)
var isWaveValid func(wave uintptr) bool
var loadSound func(sound uintptr, fileName string)
var loadSoundFromWave func(sound uintptr, wave uintptr)
var loadSoundAlias func(sound uintptr, source uintptr)
var isSoundValid func(sound uintptr) bool
var updateSound func(sound uintptr, data []byte, sampleCount int32)
var unloadWave func(wave uintptr)
var unloadSound func(sound uintptr)
var unloadSoundAlias func(alias uintptr)
var exportWave func(wave uintptr, fileName string) bool
var playSound func(sound uintptr)
var stopSound func(sound uintptr)
var pauseSound func(sound uintptr)
var resumeSound func(sound uintptr)
var isSoundPlaying func(sound uintptr) bool
var setSoundVolume func(sound uintptr, volume float32)
var setSoundPitch func(sound uintptr, pitch float32)
var setSoundPan func(sound uintptr, pan float32)
var waveCopy func(copy uintptr, wave uintptr)
var waveCrop func(wave *Wave, initFrame int32, finalFrame int32)
var waveFormat func(wave *Wave, sampleRate int32, sampleSize int32, channels int32)
var loadWaveSamples func(wave uintptr) *float32
var unloadWaveSamples func(samples []float32)
var loadMusicStream func(music uintptr, fileName string)
var loadMusicStreamFromMemory func(sound uintptr, fileType string, data []byte, dataSize int32)
var isMusicValid func(music uintptr) bool
var unloadMusicStream func(music uintptr)
var playMusicStream func(music uintptr)
var isMusicStreamPlaying func(music uintptr) bool
var updateMusicStream func(music uintptr)
var stopMusicStream func(music uintptr)
var pauseMusicStream func(music uintptr)
var resumeMusicStream func(music uintptr)
var seekMusicStream func(music uintptr, position float32)
var setMusicVolume func(music uintptr, volume float32)
var setMusicPitch func(music uintptr, pitch float32)
var setMusicPan func(music uintptr, pan float32)
var getMusicTimeLength func(music uintptr) float32
var getMusicTimePlayed func(music uintptr) float32
var loadAudioStream func(audioStream uintptr, sampleRate uint32, sampleSize uint32, channels uint32)
var isAudioStreamValid func(stream uintptr) bool
var unloadAudioStream func(stream uintptr)
var updateAudioStream func(stream uintptr, data []float32, frameCount int32)
var isAudioStreamProcessed func(stream uintptr) bool
var playAudioStream func(stream uintptr)
var pauseAudioStream func(stream uintptr)
var resumeAudioStream func(stream uintptr)
var isAudioStreamPlaying func(stream uintptr) bool
var stopAudioStream func(stream uintptr)
var setAudioStreamVolume func(stream uintptr, volume float32)
var setAudioStreamPitch func(stream uintptr, pitch float32)
var setAudioStreamPan func(stream uintptr, pan float32)
var setAudioStreamBufferSizeDefault func(size int32)
var setAudioStreamCallback func(stream uintptr, callback uintptr)
var attachAudioStreamProcessor func(stream uintptr, processor uintptr)
var detachAudioStreamProcessor func(stream uintptr, processor uintptr)
var attachAudioMixedProcessor func(processor uintptr)
var detachAudioMixedProcessor func(processor uintptr)

func init() {
	raylibDll = loadLibrary()
	audioCallbacks = make(map[uintptr]uintptr)

	initRlglPurego()

	purego.RegisterLibFunc(&initWindow, raylibDll, "InitWindow")
	purego.RegisterLibFunc(&closeWindow, raylibDll, "CloseWindow")
	purego.RegisterLibFunc(&windowShouldClose, raylibDll, "WindowShouldClose")
	purego.RegisterLibFunc(&isWindowReady, raylibDll, "IsWindowReady")
	purego.RegisterLibFunc(&isWindowFullscreen, raylibDll, "IsWindowFullscreen")
	purego.RegisterLibFunc(&isWindowHidden, raylibDll, "IsWindowHidden")
	purego.RegisterLibFunc(&isWindowMinimized, raylibDll, "IsWindowMinimized")
	purego.RegisterLibFunc(&isWindowMaximized, raylibDll, "IsWindowMaximized")
	purego.RegisterLibFunc(&isWindowFocused, raylibDll, "IsWindowFocused")
	purego.RegisterLibFunc(&isWindowResized, raylibDll, "IsWindowResized")
	purego.RegisterLibFunc(&isWindowState, raylibDll, "IsWindowState")
	purego.RegisterLibFunc(&setWindowState, raylibDll, "SetWindowState")
	purego.RegisterLibFunc(&clearWindowState, raylibDll, "ClearWindowState")
	purego.RegisterLibFunc(&toggleFullscreen, raylibDll, "ToggleFullscreen")
	purego.RegisterLibFunc(&toggleBorderlessWindowed, raylibDll, "ToggleBorderlessWindowed")
	purego.RegisterLibFunc(&maximizeWindow, raylibDll, "MaximizeWindow")
	purego.RegisterLibFunc(&minimizeWindow, raylibDll, "MinimizeWindow")
	purego.RegisterLibFunc(&restoreWindow, raylibDll, "RestoreWindow")
	purego.RegisterLibFunc(&setWindowIcon, raylibDll, "SetWindowIcon")
	purego.RegisterLibFunc(&setWindowIcons, raylibDll, "SetWindowIcons")
	purego.RegisterLibFunc(&setWindowTitle, raylibDll, "SetWindowTitle")
	purego.RegisterLibFunc(&setWindowPosition, raylibDll, "SetWindowPosition")
	purego.RegisterLibFunc(&setWindowMonitor, raylibDll, "SetWindowMonitor")
	purego.RegisterLibFunc(&setWindowMinSize, raylibDll, "SetWindowMinSize")
	purego.RegisterLibFunc(&setWindowMaxSize, raylibDll, "SetWindowMaxSize")
	purego.RegisterLibFunc(&setWindowSize, raylibDll, "SetWindowSize")
	purego.RegisterLibFunc(&setWindowOpacity, raylibDll, "SetWindowOpacity")
	purego.RegisterLibFunc(&setWindowFocused, raylibDll, "SetWindowFocused")
	purego.RegisterLibFunc(&getWindowHandle, raylibDll, "GetWindowHandle")
	purego.RegisterLibFunc(&getScreenWidth, raylibDll, "GetScreenWidth")
	purego.RegisterLibFunc(&getScreenHeight, raylibDll, "GetScreenHeight")
	purego.RegisterLibFunc(&getRenderWidth, raylibDll, "GetRenderWidth")
	purego.RegisterLibFunc(&getRenderHeight, raylibDll, "GetRenderHeight")
	purego.RegisterLibFunc(&getMonitorCount, raylibDll, "GetMonitorCount")
	purego.RegisterLibFunc(&getCurrentMonitor, raylibDll, "GetCurrentMonitor")
	purego.RegisterLibFunc(&getMonitorPosition, raylibDll, "GetMonitorPosition")
	purego.RegisterLibFunc(&getMonitorWidth, raylibDll, "GetMonitorWidth")
	purego.RegisterLibFunc(&getMonitorHeight, raylibDll, "GetMonitorHeight")
	purego.RegisterLibFunc(&getMonitorPhysicalWidth, raylibDll, "GetMonitorPhysicalWidth")
	purego.RegisterLibFunc(&getMonitorPhysicalHeight, raylibDll, "GetMonitorPhysicalHeight")
	purego.RegisterLibFunc(&getMonitorRefreshRate, raylibDll, "GetMonitorRefreshRate")
	purego.RegisterLibFunc(&getWindowPosition, raylibDll, "GetWindowPosition")
	purego.RegisterLibFunc(&getWindowScaleDPI, raylibDll, "GetWindowScaleDPI")
	purego.RegisterLibFunc(&getMonitorName, raylibDll, "GetMonitorName")
	purego.RegisterLibFunc(&setClipboardText, raylibDll, "SetClipboardText")
	purego.RegisterLibFunc(&getClipboardText, raylibDll, "GetClipboardText")
	purego.RegisterLibFunc(&getClipboardImage, raylibDll, "GetClipboardImage")
	purego.RegisterLibFunc(&enableEventWaiting, raylibDll, "EnableEventWaiting")
	purego.RegisterLibFunc(&disableEventWaiting, raylibDll, "DisableEventWaiting")
	purego.RegisterLibFunc(&showCursor, raylibDll, "ShowCursor")
	purego.RegisterLibFunc(&hideCursor, raylibDll, "HideCursor")
	purego.RegisterLibFunc(&isCursorHidden, raylibDll, "IsCursorHidden")
	purego.RegisterLibFunc(&enableCursor, raylibDll, "EnableCursor")
	purego.RegisterLibFunc(&disableCursor, raylibDll, "DisableCursor")
	purego.RegisterLibFunc(&isCursorOnScreen, raylibDll, "IsCursorOnScreen")
	purego.RegisterLibFunc(&clearBackground, raylibDll, "ClearBackground")
	purego.RegisterLibFunc(&beginDrawing, raylibDll, "BeginDrawing")
	purego.RegisterLibFunc(&endDrawing, raylibDll, "EndDrawing")
	purego.RegisterLibFunc(&beginMode2D, raylibDll, "BeginMode2D")
	purego.RegisterLibFunc(&endMode2D, raylibDll, "EndMode2D")
	purego.RegisterLibFunc(&beginMode3D, raylibDll, "BeginMode3D")
	purego.RegisterLibFunc(&endMode3D, raylibDll, "EndMode3D")
	purego.RegisterLibFunc(&beginTextureMode, raylibDll, "BeginTextureMode")
	purego.RegisterLibFunc(&endTextureMode, raylibDll, "EndTextureMode")
	purego.RegisterLibFunc(&beginShaderMode, raylibDll, "BeginShaderMode")
	purego.RegisterLibFunc(&endShaderMode, raylibDll, "EndShaderMode")
	purego.RegisterLibFunc(&beginBlendMode, raylibDll, "BeginBlendMode")
	purego.RegisterLibFunc(&endBlendMode, raylibDll, "EndBlendMode")
	purego.RegisterLibFunc(&beginScissorMode, raylibDll, "BeginScissorMode")
	purego.RegisterLibFunc(&endScissorMode, raylibDll, "EndScissorMode")
	purego.RegisterLibFunc(&beginVrStereoMode, raylibDll, "BeginVrStereoMode")
	purego.RegisterLibFunc(&endVrStereoMode, raylibDll, "EndVrStereoMode")
	purego.RegisterLibFunc(&loadVrStereoConfig, raylibDll, "LoadVrStereoConfig")
	purego.RegisterLibFunc(&unloadVrStereoConfig, raylibDll, "UnloadVrStereoConfig")
	purego.RegisterLibFunc(&loadShader, raylibDll, "LoadShader")
	purego.RegisterLibFunc(&loadShaderFromMemory, raylibDll, "LoadShaderFromMemory")
	purego.RegisterLibFunc(&isShaderValid, raylibDll, "IsShaderValid")
	purego.RegisterLibFunc(&getShaderLocation, raylibDll, "GetShaderLocation")
	purego.RegisterLibFunc(&getShaderLocationAttrib, raylibDll, "GetShaderLocationAttrib")
	purego.RegisterLibFunc(&setShaderValue, raylibDll, "SetShaderValue")
	purego.RegisterLibFunc(&setShaderValueV, raylibDll, "SetShaderValueV")
	purego.RegisterLibFunc(&setShaderValueMatrix, raylibDll, "SetShaderValueMatrix")
	purego.RegisterLibFunc(&setShaderValueTexture, raylibDll, "SetShaderValueTexture")
	purego.RegisterLibFunc(&unloadShader, raylibDll, "UnloadShader")
	purego.RegisterLibFunc(&getScreenToWorldRay, raylibDll, "GetScreenToWorldRay")
	purego.RegisterLibFunc(&getScreenToWorldRayEx, raylibDll, "GetScreenToWorldRayEx")
	purego.RegisterLibFunc(&getCameraMatrix, raylibDll, "GetCameraMatrix")
	purego.RegisterLibFunc(&getCameraMatrix2D, raylibDll, "GetCameraMatrix2D")
	purego.RegisterLibFunc(&getWorldToScreen, raylibDll, "GetWorldToScreen")
	purego.RegisterLibFunc(&getScreenToWorld2D, raylibDll, "GetScreenToWorld2D")
	purego.RegisterLibFunc(&getWorldToScreenEx, raylibDll, "GetWorldToScreenEx")
	purego.RegisterLibFunc(&getWorldToScreen2D, raylibDll, "GetWorldToScreen2D")
	purego.RegisterLibFunc(&setTargetFPS, raylibDll, "SetTargetFPS")
	purego.RegisterLibFunc(&getFrameTime, raylibDll, "GetFrameTime")
	purego.RegisterLibFunc(&getTime, raylibDll, "GetTime")
	purego.RegisterLibFunc(&getFPS, raylibDll, "GetFPS")
	purego.RegisterLibFunc(&swapScreenBuffer, raylibDll, "SwapScreenBuffer")
	purego.RegisterLibFunc(&pollInputEvents, raylibDll, "PollInputEvents")
	purego.RegisterLibFunc(&waitTime, raylibDll, "WaitTime")
	purego.RegisterLibFunc(&setRandomSeed, raylibDll, "SetRandomSeed")
	purego.RegisterLibFunc(&getRandomValue, raylibDll, "GetRandomValue")
	purego.RegisterLibFunc(&loadRandomSequence, raylibDll, "LoadRandomSequence")
	purego.RegisterLibFunc(&unloadRandomSequence, raylibDll, "UnloadRandomSequence")
	purego.RegisterLibFunc(&takeScreenshot, raylibDll, "TakeScreenshot")
	purego.RegisterLibFunc(&setConfigFlags, raylibDll, "SetConfigFlags")
	purego.RegisterLibFunc(&openURL, raylibDll, "OpenURL")
	purego.RegisterLibFunc(&traceLog, raylibDll, "TraceLog")
	purego.RegisterLibFunc(&setTraceLogLevel, raylibDll, "SetTraceLogLevel")
	purego.RegisterLibFunc(&memAlloc, raylibDll, "MemAlloc")
	purego.RegisterLibFunc(&memRealloc, raylibDll, "MemRealloc")
	purego.RegisterLibFunc(&memFree, raylibDll, "MemFree")
	purego.RegisterLibFunc(&setTraceLogCallback, raylibDll, "SetTraceLogCallback")
	purego.RegisterLibFunc(&isFileDropped, raylibDll, "IsFileDropped")
	purego.RegisterLibFunc(&loadDroppedFiles, raylibDll, "LoadDroppedFiles")
	purego.RegisterLibFunc(&unloadDroppedFiles, raylibDll, "UnloadDroppedFiles")
	purego.RegisterLibFunc(&loadAutomationEventList, raylibDll, "LoadAutomationEventList")
	purego.RegisterLibFunc(&unloadAutomationEventList, raylibDll, "UnloadAutomationEventList")
	purego.RegisterLibFunc(&exportAutomationEventList, raylibDll, "ExportAutomationEventList")
	purego.RegisterLibFunc(&setAutomationEventList, raylibDll, "SetAutomationEventList")
	purego.RegisterLibFunc(&setAutomationEventBaseFrame, raylibDll, "SetAutomationEventBaseFrame")
	purego.RegisterLibFunc(&startAutomationEventRecording, raylibDll, "StartAutomationEventRecording")
	purego.RegisterLibFunc(&stopAutomationEventRecording, raylibDll, "StopAutomationEventRecording")
	purego.RegisterLibFunc(&playAutomationEvent, raylibDll, "PlayAutomationEvent")
	purego.RegisterLibFunc(&isKeyPressed, raylibDll, "IsKeyPressed")
	purego.RegisterLibFunc(&isKeyPressedRepeat, raylibDll, "IsKeyPressedRepeat")
	purego.RegisterLibFunc(&isKeyDown, raylibDll, "IsKeyDown")
	purego.RegisterLibFunc(&isKeyReleased, raylibDll, "IsKeyReleased")
	purego.RegisterLibFunc(&isKeyUp, raylibDll, "IsKeyUp")
	purego.RegisterLibFunc(&getKeyPressed, raylibDll, "GetKeyPressed")
	purego.RegisterLibFunc(&getCharPressed, raylibDll, "GetCharPressed")
	purego.RegisterLibFunc(&setExitKey, raylibDll, "SetExitKey")
	purego.RegisterLibFunc(&isGamepadAvailable, raylibDll, "IsGamepadAvailable")
	purego.RegisterLibFunc(&getGamepadName, raylibDll, "GetGamepadName")
	purego.RegisterLibFunc(&isGamepadButtonPressed, raylibDll, "IsGamepadButtonPressed")
	purego.RegisterLibFunc(&isGamepadButtonDown, raylibDll, "IsGamepadButtonDown")
	purego.RegisterLibFunc(&isGamepadButtonReleased, raylibDll, "IsGamepadButtonReleased")
	purego.RegisterLibFunc(&isGamepadButtonUp, raylibDll, "IsGamepadButtonUp")
	purego.RegisterLibFunc(&getGamepadButtonPressed, raylibDll, "GetGamepadButtonPressed")
	purego.RegisterLibFunc(&getGamepadAxisCount, raylibDll, "GetGamepadAxisCount")
	purego.RegisterLibFunc(&getGamepadAxisMovement, raylibDll, "GetGamepadAxisMovement")
	purego.RegisterLibFunc(&setGamepadMappings, raylibDll, "SetGamepadMappings")
	purego.RegisterLibFunc(&setGamepadVibration, raylibDll, "SetGamepadVibration")
	purego.RegisterLibFunc(&isMouseButtonPressed, raylibDll, "IsMouseButtonPressed")
	purego.RegisterLibFunc(&isMouseButtonDown, raylibDll, "IsMouseButtonDown")
	purego.RegisterLibFunc(&isMouseButtonReleased, raylibDll, "IsMouseButtonReleased")
	purego.RegisterLibFunc(&isMouseButtonUp, raylibDll, "IsMouseButtonUp")
	purego.RegisterLibFunc(&getMouseX, raylibDll, "GetMouseX")
	purego.RegisterLibFunc(&getMouseY, raylibDll, "GetMouseY")
	purego.RegisterLibFunc(&getMousePosition, raylibDll, "GetMousePosition")
	purego.RegisterLibFunc(&getMouseDelta, raylibDll, "GetMouseDelta")
	purego.RegisterLibFunc(&setMousePosition, raylibDll, "SetMousePosition")
	purego.RegisterLibFunc(&setMouseOffset, raylibDll, "SetMouseOffset")
	purego.RegisterLibFunc(&setMouseScale, raylibDll, "SetMouseScale")
	purego.RegisterLibFunc(&getMouseWheelMove, raylibDll, "GetMouseWheelMove")
	purego.RegisterLibFunc(&getMouseWheelMoveV, raylibDll, "GetMouseWheelMoveV")
	purego.RegisterLibFunc(&setMouseCursor, raylibDll, "SetMouseCursor")
	purego.RegisterLibFunc(&getTouchX, raylibDll, "GetTouchX")
	purego.RegisterLibFunc(&getTouchY, raylibDll, "GetTouchY")
	purego.RegisterLibFunc(&getTouchPosition, raylibDll, "GetTouchPosition")
	purego.RegisterLibFunc(&getTouchPointId, raylibDll, "GetTouchPointId")
	purego.RegisterLibFunc(&getTouchPointCount, raylibDll, "GetTouchPointCount")
	purego.RegisterLibFunc(&setGesturesEnabled, raylibDll, "SetGesturesEnabled")
	purego.RegisterLibFunc(&isGestureDetected, raylibDll, "IsGestureDetected")
	purego.RegisterLibFunc(&getGestureDetected, raylibDll, "GetGestureDetected")
	purego.RegisterLibFunc(&getGestureHoldDuration, raylibDll, "GetGestureHoldDuration")
	purego.RegisterLibFunc(&getGestureDragVector, raylibDll, "GetGestureDragVector")
	purego.RegisterLibFunc(&getGestureDragAngle, raylibDll, "GetGestureDragAngle")
	purego.RegisterLibFunc(&getGesturePinchVector, raylibDll, "GetGesturePinchVector")
	purego.RegisterLibFunc(&getGesturePinchAngle, raylibDll, "GetGesturePinchAngle")
	purego.RegisterLibFunc(&setShapesTexture, raylibDll, "SetShapesTexture")
	purego.RegisterLibFunc(&getShapesTexture, raylibDll, "GetShapesTexture")
	purego.RegisterLibFunc(&getShapesTextureRectangle, raylibDll, "GetShapesTextureRectangle")
	purego.RegisterLibFunc(&drawPixel, raylibDll, "DrawPixel")
	purego.RegisterLibFunc(&drawPixelV, raylibDll, "DrawPixelV")
	purego.RegisterLibFunc(&drawLine, raylibDll, "DrawLine")
	purego.RegisterLibFunc(&drawLineV, raylibDll, "DrawLineV")
	purego.RegisterLibFunc(&drawLineEx, raylibDll, "DrawLineEx")
	purego.RegisterLibFunc(&drawLineStrip, raylibDll, "DrawLineStrip")
	purego.RegisterLibFunc(&drawLineBezier, raylibDll, "DrawLineBezier")
	purego.RegisterLibFunc(&drawCircle, raylibDll, "DrawCircle")
	purego.RegisterLibFunc(&drawCircleSector, raylibDll, "DrawCircleSector")
	purego.RegisterLibFunc(&drawCircleSectorLines, raylibDll, "DrawCircleSectorLines")
	purego.RegisterLibFunc(&drawCircleGradient, raylibDll, "DrawCircleGradient")
	purego.RegisterLibFunc(&drawCircleV, raylibDll, "DrawCircleV")
	purego.RegisterLibFunc(&drawCircleLines, raylibDll, "DrawCircleLines")
	purego.RegisterLibFunc(&drawCircleLinesV, raylibDll, "DrawCircleLinesV")
	purego.RegisterLibFunc(&drawEllipse, raylibDll, "DrawEllipse")
	purego.RegisterLibFunc(&drawEllipseLines, raylibDll, "DrawEllipseLines")
	purego.RegisterLibFunc(&drawRing, raylibDll, "DrawRing")
	purego.RegisterLibFunc(&drawRingLines, raylibDll, "DrawRingLines")
	purego.RegisterLibFunc(&drawRectangle, raylibDll, "DrawRectangle")
	purego.RegisterLibFunc(&drawRectangleV, raylibDll, "DrawRectangleV")
	purego.RegisterLibFunc(&drawRectangleRec, raylibDll, "DrawRectangleRec")
	purego.RegisterLibFunc(&drawRectanglePro, raylibDll, "DrawRectanglePro")
	purego.RegisterLibFunc(&drawRectangleGradientV, raylibDll, "DrawRectangleGradientV")
	purego.RegisterLibFunc(&drawRectangleGradientH, raylibDll, "DrawRectangleGradientH")
	purego.RegisterLibFunc(&drawRectangleGradientEx, raylibDll, "DrawRectangleGradientEx")
	purego.RegisterLibFunc(&drawRectangleLines, raylibDll, "DrawRectangleLines")
	purego.RegisterLibFunc(&drawRectangleLinesEx, raylibDll, "DrawRectangleLinesEx")
	purego.RegisterLibFunc(&drawRectangleRounded, raylibDll, "DrawRectangleRounded")
	purego.RegisterLibFunc(&drawRectangleRoundedLines, raylibDll, "DrawRectangleRoundedLines")
	purego.RegisterLibFunc(&drawRectangleRoundedLinesEx, raylibDll, "DrawRectangleRoundedLinesEx")
	purego.RegisterLibFunc(&drawTriangle, raylibDll, "DrawTriangle")
	purego.RegisterLibFunc(&drawTriangleLines, raylibDll, "DrawTriangleLines")
	purego.RegisterLibFunc(&drawTriangleFan, raylibDll, "DrawTriangleFan")
	purego.RegisterLibFunc(&drawTriangleStrip, raylibDll, "DrawTriangleStrip")
	purego.RegisterLibFunc(&drawPoly, raylibDll, "DrawPoly")
	purego.RegisterLibFunc(&drawPolyLines, raylibDll, "DrawPolyLines")
	purego.RegisterLibFunc(&drawPolyLinesEx, raylibDll, "DrawPolyLinesEx")
	purego.RegisterLibFunc(&drawSplineLinear, raylibDll, "DrawSplineLinear")
	purego.RegisterLibFunc(&drawSplineBasis, raylibDll, "DrawSplineBasis")
	purego.RegisterLibFunc(&drawSplineCatmullRom, raylibDll, "DrawSplineCatmullRom")
	purego.RegisterLibFunc(&drawSplineBezierQuadratic, raylibDll, "DrawSplineBezierQuadratic")
	purego.RegisterLibFunc(&drawSplineBezierCubic, raylibDll, "DrawSplineBezierCubic")
	purego.RegisterLibFunc(&drawSplineSegmentLinear, raylibDll, "DrawSplineSegmentLinear")
	purego.RegisterLibFunc(&drawSplineSegmentBasis, raylibDll, "DrawSplineSegmentBasis")
	purego.RegisterLibFunc(&drawSplineSegmentCatmullRom, raylibDll, "DrawSplineSegmentCatmullRom")
	purego.RegisterLibFunc(&drawSplineSegmentBezierQuadratic, raylibDll, "DrawSplineSegmentBezierQuadratic")
	purego.RegisterLibFunc(&drawSplineSegmentBezierCubic, raylibDll, "DrawSplineSegmentBezierCubic")
	purego.RegisterLibFunc(&getSplinePointLinear, raylibDll, "GetSplinePointLinear")
	purego.RegisterLibFunc(&getSplinePointBasis, raylibDll, "GetSplinePointBasis")
	purego.RegisterLibFunc(&getSplinePointCatmullRom, raylibDll, "GetSplinePointCatmullRom")
	purego.RegisterLibFunc(&getSplinePointBezierQuad, raylibDll, "GetSplinePointBezierQuad")
	purego.RegisterLibFunc(&getSplinePointBezierCubic, raylibDll, "GetSplinePointBezierCubic")
	purego.RegisterLibFunc(&checkCollisionRecs, raylibDll, "CheckCollisionRecs")
	purego.RegisterLibFunc(&checkCollisionCircles, raylibDll, "CheckCollisionCircles")
	purego.RegisterLibFunc(&checkCollisionCircleRec, raylibDll, "CheckCollisionCircleRec")
	purego.RegisterLibFunc(&checkCollisionCircleLine, raylibDll, "CheckCollisionCircleLine")
	purego.RegisterLibFunc(&checkCollisionPointRec, raylibDll, "CheckCollisionPointRec")
	purego.RegisterLibFunc(&checkCollisionPointCircle, raylibDll, "CheckCollisionPointCircle")
	purego.RegisterLibFunc(&checkCollisionPointTriangle, raylibDll, "CheckCollisionPointTriangle")
	purego.RegisterLibFunc(&checkCollisionPointPoly, raylibDll, "CheckCollisionPointPoly")
	purego.RegisterLibFunc(&checkCollisionLines, raylibDll, "CheckCollisionLines")
	purego.RegisterLibFunc(&checkCollisionPointLine, raylibDll, "CheckCollisionPointLine")
	purego.RegisterLibFunc(&getCollisionRec, raylibDll, "GetCollisionRec")
	purego.RegisterLibFunc(&loadImage, raylibDll, "LoadImage")
	purego.RegisterLibFunc(&loadImageRaw, raylibDll, "LoadImageRaw")
	purego.RegisterLibFunc(&loadImageAnim, raylibDll, "LoadImageAnim")
	purego.RegisterLibFunc(&loadImageAnimFromMemory, raylibDll, "LoadImageAnimFromMemory")
	purego.RegisterLibFunc(&loadImageFromMemory, raylibDll, "LoadImageFromMemory")
	purego.RegisterLibFunc(&loadImageFromTexture, raylibDll, "LoadImageFromTexture")
	purego.RegisterLibFunc(&loadImageFromScreen, raylibDll, "LoadImageFromScreen")
	purego.RegisterLibFunc(&isImageValid, raylibDll, "IsImageValid")
	purego.RegisterLibFunc(&unloadImage, raylibDll, "UnloadImage")
	purego.RegisterLibFunc(&exportImage, raylibDll, "ExportImage")
	purego.RegisterLibFunc(&exportImageToMemory, raylibDll, "ExportImageToMemory")
	purego.RegisterLibFunc(&genImageColor, raylibDll, "GenImageColor")
	purego.RegisterLibFunc(&genImageGradientLinear, raylibDll, "GenImageGradientLinear")
	purego.RegisterLibFunc(&genImageGradientRadial, raylibDll, "GenImageGradientRadial")
	purego.RegisterLibFunc(&genImageGradientSquare, raylibDll, "GenImageGradientSquare")
	purego.RegisterLibFunc(&genImageChecked, raylibDll, "GenImageChecked")
	purego.RegisterLibFunc(&genImageWhiteNoise, raylibDll, "GenImageWhiteNoise")
	purego.RegisterLibFunc(&genImagePerlinNoise, raylibDll, "GenImagePerlinNoise")
	purego.RegisterLibFunc(&genImageCellular, raylibDll, "GenImageCellular")
	purego.RegisterLibFunc(&genImageText, raylibDll, "GenImageText")
	purego.RegisterLibFunc(&imageCopy, raylibDll, "ImageCopy")
	purego.RegisterLibFunc(&imageFromImage, raylibDll, "ImageFromImage")
	purego.RegisterLibFunc(&imageFromChannel, raylibDll, "ImageFromChannel")
	purego.RegisterLibFunc(&imageText, raylibDll, "ImageText")
	purego.RegisterLibFunc(&imageTextEx, raylibDll, "ImageTextEx")
	purego.RegisterLibFunc(&imageFormat, raylibDll, "ImageFormat")
	purego.RegisterLibFunc(&imageToPOT, raylibDll, "ImageToPOT")
	purego.RegisterLibFunc(&imageCrop, raylibDll, "ImageCrop")
	purego.RegisterLibFunc(&imageAlphaCrop, raylibDll, "ImageAlphaCrop")
	purego.RegisterLibFunc(&imageAlphaClear, raylibDll, "ImageAlphaClear")
	purego.RegisterLibFunc(&imageAlphaMask, raylibDll, "ImageAlphaMask")
	purego.RegisterLibFunc(&imageAlphaPremultiply, raylibDll, "ImageAlphaPremultiply")
	purego.RegisterLibFunc(&imageBlurGaussian, raylibDll, "ImageBlurGaussian")
	purego.RegisterLibFunc(&imageKernelConvolution, raylibDll, "ImageKernelConvolution")
	purego.RegisterLibFunc(&imageResize, raylibDll, "ImageResize")
	purego.RegisterLibFunc(&imageResizeNN, raylibDll, "ImageResizeNN")
	purego.RegisterLibFunc(&imageResizeCanvas, raylibDll, "ImageResizeCanvas")
	purego.RegisterLibFunc(&imageMipmaps, raylibDll, "ImageMipmaps")
	purego.RegisterLibFunc(&imageDither, raylibDll, "ImageDither")
	purego.RegisterLibFunc(&imageFlipVertical, raylibDll, "ImageFlipVertical")
	purego.RegisterLibFunc(&imageFlipHorizontal, raylibDll, "ImageFlipHorizontal")
	purego.RegisterLibFunc(&imageRotate, raylibDll, "ImageRotate")
	purego.RegisterLibFunc(&imageRotateCW, raylibDll, "ImageRotateCW")
	purego.RegisterLibFunc(&imageRotateCCW, raylibDll, "ImageRotateCCW")
	purego.RegisterLibFunc(&imageColorTint, raylibDll, "ImageColorTint")
	purego.RegisterLibFunc(&imageColorInvert, raylibDll, "ImageColorInvert")
	purego.RegisterLibFunc(&imageColorGrayscale, raylibDll, "ImageColorGrayscale")
	purego.RegisterLibFunc(&imageColorContrast, raylibDll, "ImageColorContrast")
	purego.RegisterLibFunc(&imageColorBrightness, raylibDll, "ImageColorBrightness")
	purego.RegisterLibFunc(&imageColorReplace, raylibDll, "ImageColorReplace")
	purego.RegisterLibFunc(&loadImageColors, raylibDll, "LoadImageColors")
	purego.RegisterLibFunc(&loadImagePalette, raylibDll, "LoadImagePalette")
	purego.RegisterLibFunc(&unloadImageColors, raylibDll, "UnloadImageColors")
	purego.RegisterLibFunc(&unloadImagePalette, raylibDll, "UnloadImagePalette")
	purego.RegisterLibFunc(&getImageAlphaBorder, raylibDll, "GetImageAlphaBorder")
	purego.RegisterLibFunc(&getImageColor, raylibDll, "GetImageColor")
	purego.RegisterLibFunc(&imageClearBackground, raylibDll, "ImageClearBackground")
	purego.RegisterLibFunc(&imageDrawPixel, raylibDll, "ImageDrawPixel")
	purego.RegisterLibFunc(&imageDrawPixelV, raylibDll, "ImageDrawPixelV")
	purego.RegisterLibFunc(&imageDrawLine, raylibDll, "ImageDrawLine")
	purego.RegisterLibFunc(&imageDrawLineV, raylibDll, "ImageDrawLineV")
	purego.RegisterLibFunc(&imageDrawLineEx, raylibDll, "ImageDrawLineEx")
	purego.RegisterLibFunc(&imageDrawCircle, raylibDll, "ImageDrawCircle")
	purego.RegisterLibFunc(&imageDrawCircleV, raylibDll, "ImageDrawCircleV")
	purego.RegisterLibFunc(&imageDrawCircleLines, raylibDll, "ImageDrawCircleLines")
	purego.RegisterLibFunc(&imageDrawCircleLinesV, raylibDll, "ImageDrawCircleLinesV")
	purego.RegisterLibFunc(&imageDrawRectangle, raylibDll, "ImageDrawRectangle")
	purego.RegisterLibFunc(&imageDrawRectangleV, raylibDll, "ImageDrawRectangleV")
	purego.RegisterLibFunc(&imageDrawRectangleRec, raylibDll, "ImageDrawRectangleRec")
	purego.RegisterLibFunc(&imageDrawRectangleLines, raylibDll, "ImageDrawRectangleLines")
	purego.RegisterLibFunc(&imageDrawTriangle, raylibDll, "ImageDrawTriangle")
	purego.RegisterLibFunc(&imageDrawTriangleEx, raylibDll, "ImageDrawTriangleEx")
	purego.RegisterLibFunc(&imageDrawTriangleLines, raylibDll, "ImageDrawTriangleLines")
	purego.RegisterLibFunc(&imageDrawTriangleFan, raylibDll, "ImageDrawTriangleFan")
	purego.RegisterLibFunc(&imageDrawTriangleStrip, raylibDll, "ImageDrawTriangleStrip")
	purego.RegisterLibFunc(&imageDraw, raylibDll, "ImageDraw")
	purego.RegisterLibFunc(&imageDrawText, raylibDll, "ImageDrawText")
	purego.RegisterLibFunc(&imageDrawTextEx, raylibDll, "ImageDrawTextEx")
	purego.RegisterLibFunc(&loadTexture, raylibDll, "LoadTexture")
	purego.RegisterLibFunc(&loadTextureFromImage, raylibDll, "LoadTextureFromImage")
	purego.RegisterLibFunc(&loadTextureCubemap, raylibDll, "LoadTextureCubemap")
	purego.RegisterLibFunc(&loadRenderTexture, raylibDll, "LoadRenderTexture")
	purego.RegisterLibFunc(&isTextureValid, raylibDll, "IsTextureValid")
	purego.RegisterLibFunc(&unloadTexture, raylibDll, "UnloadTexture")
	purego.RegisterLibFunc(&isRenderTextureValid, raylibDll, "IsRenderTextureValid")
	purego.RegisterLibFunc(&unloadRenderTexture, raylibDll, "UnloadRenderTexture")
	purego.RegisterLibFunc(&updateTexture, raylibDll, "UpdateTexture")
	purego.RegisterLibFunc(&updateTextureRec, raylibDll, "UpdateTextureRec")
	purego.RegisterLibFunc(&genTextureMipmaps, raylibDll, "GenTextureMipmaps")
	purego.RegisterLibFunc(&setTextureFilter, raylibDll, "SetTextureFilter")
	purego.RegisterLibFunc(&setTextureWrap, raylibDll, "SetTextureWrap")
	purego.RegisterLibFunc(&drawTexture, raylibDll, "DrawTexture")
	purego.RegisterLibFunc(&drawTextureV, raylibDll, "DrawTextureV")
	purego.RegisterLibFunc(&drawTextureEx, raylibDll, "DrawTextureEx")
	purego.RegisterLibFunc(&drawTextureRec, raylibDll, "DrawTextureRec")
	purego.RegisterLibFunc(&drawTexturePro, raylibDll, "DrawTexturePro")
	purego.RegisterLibFunc(&drawTextureNPatch, raylibDll, "DrawTextureNPatch")
	purego.RegisterLibFunc(&fade, raylibDll, "Fade")
	purego.RegisterLibFunc(&colorToInt, raylibDll, "ColorToInt")
	purego.RegisterLibFunc(&colorNormalize, raylibDll, "ColorNormalize")
	purego.RegisterLibFunc(&colorFromNormalized, raylibDll, "ColorFromNormalized")
	purego.RegisterLibFunc(&colorToHSV, raylibDll, "ColorToHSV")
	purego.RegisterLibFunc(&colorFromHSV, raylibDll, "ColorFromHSV")
	purego.RegisterLibFunc(&colorTint, raylibDll, "ColorTint")
	purego.RegisterLibFunc(&colorBrightness, raylibDll, "ColorBrightness")
	purego.RegisterLibFunc(&colorContrast, raylibDll, "ColorContrast")
	purego.RegisterLibFunc(&colorAlpha, raylibDll, "ColorAlpha")
	purego.RegisterLibFunc(&colorAlphaBlend, raylibDll, "ColorAlphaBlend")
	purego.RegisterLibFunc(&colorLerp, raylibDll, "ColorLerp")
	purego.RegisterLibFunc(&getColor, raylibDll, "GetColor")
	purego.RegisterLibFunc(&getPixelColor, raylibDll, "GetPixelColor")
	purego.RegisterLibFunc(&setPixelColor, raylibDll, "SetPixelColor")
	purego.RegisterLibFunc(&getPixelDataSize, raylibDll, "GetPixelDataSize")
	purego.RegisterLibFunc(&getFontDefault, raylibDll, "GetFontDefault")
	purego.RegisterLibFunc(&loadFont, raylibDll, "LoadFont")
	purego.RegisterLibFunc(&loadFontEx, raylibDll, "LoadFontEx")
	purego.RegisterLibFunc(&loadFontFromImage, raylibDll, "LoadFontFromImage")
	purego.RegisterLibFunc(&loadFontFromMemory, raylibDll, "LoadFontFromMemory")
	purego.RegisterLibFunc(&isFontValid, raylibDll, "IsFontValid")
	purego.RegisterLibFunc(&loadFontData, raylibDll, "LoadFontData")
	purego.RegisterLibFunc(&genImageFontAtlas, raylibDll, "GenImageFontAtlas")
	purego.RegisterLibFunc(&unloadFontData, raylibDll, "UnloadFontData")
	purego.RegisterLibFunc(&unloadFont, raylibDll, "UnloadFont")
	purego.RegisterLibFunc(&drawFPS, raylibDll, "DrawFPS")
	purego.RegisterLibFunc(&drawText, raylibDll, "DrawText")
	purego.RegisterLibFunc(&drawTextEx, raylibDll, "DrawTextEx")
	purego.RegisterLibFunc(&drawTextPro, raylibDll, "DrawTextPro")
	purego.RegisterLibFunc(&drawTextCodepoint, raylibDll, "DrawTextCodepoint")
	purego.RegisterLibFunc(&drawTextCodepoints, raylibDll, "DrawTextCodepoints")
	purego.RegisterLibFunc(&setTextLineSpacing, raylibDll, "SetTextLineSpacing")
	purego.RegisterLibFunc(&measureText, raylibDll, "MeasureText")
	purego.RegisterLibFunc(&measureTextEx, raylibDll, "MeasureTextEx")
	purego.RegisterLibFunc(&getGlyphIndex, raylibDll, "GetGlyphIndex")
	purego.RegisterLibFunc(&getGlyphInfo, raylibDll, "GetGlyphInfo")
	purego.RegisterLibFunc(&getGlyphAtlasRec, raylibDll, "GetGlyphAtlasRec")
	purego.RegisterLibFunc(&drawLine3D, raylibDll, "DrawLine3D")
	purego.RegisterLibFunc(&drawPoint3D, raylibDll, "DrawPoint3D")
	purego.RegisterLibFunc(&drawCircle3D, raylibDll, "DrawCircle3D")
	purego.RegisterLibFunc(&drawTriangle3D, raylibDll, "DrawTriangle3D")
	purego.RegisterLibFunc(&drawTriangleStrip3D, raylibDll, "DrawTriangleStrip3D")
	purego.RegisterLibFunc(&drawCube, raylibDll, "DrawCube")
	purego.RegisterLibFunc(&drawCubeV, raylibDll, "DrawCubeV")
	purego.RegisterLibFunc(&drawCubeWires, raylibDll, "DrawCubeWires")
	purego.RegisterLibFunc(&drawCubeWiresV, raylibDll, "DrawCubeWiresV")
	purego.RegisterLibFunc(&drawSphere, raylibDll, "DrawSphere")
	purego.RegisterLibFunc(&drawSphereEx, raylibDll, "DrawSphereEx")
	purego.RegisterLibFunc(&drawSphereWires, raylibDll, "DrawSphereWires")
	purego.RegisterLibFunc(&drawCylinder, raylibDll, "DrawCylinder")
	purego.RegisterLibFunc(&drawCylinderEx, raylibDll, "DrawCylinderEx")
	purego.RegisterLibFunc(&drawCylinderWires, raylibDll, "DrawCylinderWires")
	purego.RegisterLibFunc(&drawCylinderWiresEx, raylibDll, "DrawCylinderWiresEx")
	purego.RegisterLibFunc(&drawCapsule, raylibDll, "DrawCapsule")
	purego.RegisterLibFunc(&drawCapsuleWires, raylibDll, "DrawCapsuleWires")
	purego.RegisterLibFunc(&drawPlane, raylibDll, "DrawPlane")
	purego.RegisterLibFunc(&drawRay, raylibDll, "DrawRay")
	purego.RegisterLibFunc(&drawGrid, raylibDll, "DrawGrid")
	purego.RegisterLibFunc(&loadModel, raylibDll, "LoadModel")
	purego.RegisterLibFunc(&loadModelFromMesh, raylibDll, "LoadModelFromMesh")
	purego.RegisterLibFunc(&isModelValid, raylibDll, "IsModelValid")
	purego.RegisterLibFunc(&unloadModel, raylibDll, "UnloadModel")
	purego.RegisterLibFunc(&getModelBoundingBox, raylibDll, "GetModelBoundingBox")
	purego.RegisterLibFunc(&drawModel, raylibDll, "DrawModel")
	purego.RegisterLibFunc(&drawModelEx, raylibDll, "DrawModelEx")
	purego.RegisterLibFunc(&drawModelWires, raylibDll, "DrawModelWires")
	purego.RegisterLibFunc(&drawModelWiresEx, raylibDll, "DrawModelWiresEx")
	purego.RegisterLibFunc(&drawModelPoints, raylibDll, "DrawModelPoints")
	purego.RegisterLibFunc(&drawModelPointsEx, raylibDll, "DrawModelPointsEx")
	purego.RegisterLibFunc(&drawBoundingBox, raylibDll, "DrawBoundingBox")
	purego.RegisterLibFunc(&drawBillboard, raylibDll, "DrawBillboard")
	purego.RegisterLibFunc(&drawBillboardRec, raylibDll, "DrawBillboardRec")
	purego.RegisterLibFunc(&drawBillboardPro, raylibDll, "DrawBillboardPro")
	purego.RegisterLibFunc(&uploadMesh, raylibDll, "UploadMesh")
	purego.RegisterLibFunc(&updateMeshBuffer, raylibDll, "UpdateMeshBuffer")
	purego.RegisterLibFunc(&unloadMesh, raylibDll, "UnloadMesh")
	purego.RegisterLibFunc(&drawMesh, raylibDll, "DrawMesh")
	purego.RegisterLibFunc(&drawMeshInstanced, raylibDll, "DrawMeshInstanced")
	purego.RegisterLibFunc(&exportMesh, raylibDll, "ExportMesh")
	purego.RegisterLibFunc(&getMeshBoundingBox, raylibDll, "GetMeshBoundingBox")
	purego.RegisterLibFunc(&genMeshTangents, raylibDll, "GenMeshTangents")
	purego.RegisterLibFunc(&genMeshPoly, raylibDll, "GenMeshPoly")
	purego.RegisterLibFunc(&genMeshPlane, raylibDll, "GenMeshPlane")
	purego.RegisterLibFunc(&genMeshCube, raylibDll, "GenMeshCube")
	purego.RegisterLibFunc(&genMeshSphere, raylibDll, "GenMeshSphere")
	purego.RegisterLibFunc(&genMeshHemiSphere, raylibDll, "GenMeshHemiSphere")
	purego.RegisterLibFunc(&genMeshCylinder, raylibDll, "GenMeshCylinder")
	purego.RegisterLibFunc(&genMeshCone, raylibDll, "GenMeshCone")
	purego.RegisterLibFunc(&genMeshTorus, raylibDll, "GenMeshTorus")
	purego.RegisterLibFunc(&genMeshKnot, raylibDll, "GenMeshKnot")
	purego.RegisterLibFunc(&genMeshHeightmap, raylibDll, "GenMeshHeightmap")
	purego.RegisterLibFunc(&genMeshCubicmap, raylibDll, "GenMeshCubicmap")
	purego.RegisterLibFunc(&loadMaterials, raylibDll, "LoadMaterials")
	purego.RegisterLibFunc(&loadMaterialDefault, raylibDll, "LoadMaterialDefault")
	purego.RegisterLibFunc(&isMaterialValid, raylibDll, "IsMaterialValid")
	purego.RegisterLibFunc(&unloadMaterial, raylibDll, "UnloadMaterial")
	purego.RegisterLibFunc(&setMaterialTexture, raylibDll, "SetMaterialTexture")
	purego.RegisterLibFunc(&setModelMeshMaterial, raylibDll, "SetModelMeshMaterial")
	purego.RegisterLibFunc(&loadModelAnimations, raylibDll, "LoadModelAnimations")
	purego.RegisterLibFunc(&updateModelAnimation, raylibDll, "UpdateModelAnimation")
	purego.RegisterLibFunc(&updateModelAnimationBones, raylibDll, "UpdateModelAnimationBones")
	purego.RegisterLibFunc(&unloadModelAnimation, raylibDll, "UnloadModelAnimation")
	purego.RegisterLibFunc(&unloadModelAnimations, raylibDll, "UnloadModelAnimations")
	purego.RegisterLibFunc(&isModelAnimationValid, raylibDll, "IsModelAnimationValid")
	purego.RegisterLibFunc(&checkCollisionSpheres, raylibDll, "CheckCollisionSpheres")
	purego.RegisterLibFunc(&checkCollisionBoxes, raylibDll, "CheckCollisionBoxes")
	purego.RegisterLibFunc(&checkCollisionBoxSphere, raylibDll, "CheckCollisionBoxSphere")
	purego.RegisterLibFunc(&getRayCollisionSphere, raylibDll, "GetRayCollisionSphere")
	purego.RegisterLibFunc(&getRayCollisionBox, raylibDll, "GetRayCollisionBox")
	purego.RegisterLibFunc(&getRayCollisionMesh, raylibDll, "GetRayCollisionMesh")
	purego.RegisterLibFunc(&getRayCollisionTriangle, raylibDll, "GetRayCollisionTriangle")
	purego.RegisterLibFunc(&getRayCollisionQuad, raylibDll, "GetRayCollisionQuad")
	purego.RegisterLibFunc(&initAudioDevice, raylibDll, "InitAudioDevice")
	purego.RegisterLibFunc(&closeAudioDevice, raylibDll, "CloseAudioDevice")
	purego.RegisterLibFunc(&isAudioDeviceReady, raylibDll, "IsAudioDeviceReady")
	purego.RegisterLibFunc(&setMasterVolume, raylibDll, "SetMasterVolume")
	purego.RegisterLibFunc(&getMasterVolume, raylibDll, "GetMasterVolume")
	purego.RegisterLibFunc(&loadWave, raylibDll, "LoadWave")
	purego.RegisterLibFunc(&loadWaveFromMemory, raylibDll, "LoadWaveFromMemory")
	purego.RegisterLibFunc(&isWaveValid, raylibDll, "IsWaveValid")
	purego.RegisterLibFunc(&loadSound, raylibDll, "LoadSound")
	purego.RegisterLibFunc(&loadSoundFromWave, raylibDll, "LoadSoundFromWave")
	purego.RegisterLibFunc(&loadSoundAlias, raylibDll, "LoadSoundAlias")
	purego.RegisterLibFunc(&isSoundValid, raylibDll, "IsSoundValid")
	purego.RegisterLibFunc(&updateSound, raylibDll, "UpdateSound")
	purego.RegisterLibFunc(&unloadWave, raylibDll, "UnloadWave")
	purego.RegisterLibFunc(&unloadSound, raylibDll, "UnloadSound")
	purego.RegisterLibFunc(&unloadSoundAlias, raylibDll, "UnloadSoundAlias")
	purego.RegisterLibFunc(&exportWave, raylibDll, "ExportWave")
	purego.RegisterLibFunc(&playSound, raylibDll, "PlaySound")
	purego.RegisterLibFunc(&stopSound, raylibDll, "StopSound")
	purego.RegisterLibFunc(&pauseSound, raylibDll, "PauseSound")
	purego.RegisterLibFunc(&resumeSound, raylibDll, "ResumeSound")
	purego.RegisterLibFunc(&isSoundPlaying, raylibDll, "IsSoundPlaying")
	purego.RegisterLibFunc(&setSoundVolume, raylibDll, "SetSoundVolume")
	purego.RegisterLibFunc(&setSoundPitch, raylibDll, "SetSoundPitch")
	purego.RegisterLibFunc(&setSoundPan, raylibDll, "SetSoundPan")
	purego.RegisterLibFunc(&waveCopy, raylibDll, "WaveCopy")
	purego.RegisterLibFunc(&waveCrop, raylibDll, "WaveCrop")
	purego.RegisterLibFunc(&waveFormat, raylibDll, "WaveFormat")
	purego.RegisterLibFunc(&loadWaveSamples, raylibDll, "LoadWaveSamples")
	purego.RegisterLibFunc(&unloadWaveSamples, raylibDll, "UnloadWaveSamples")
	purego.RegisterLibFunc(&loadMusicStream, raylibDll, "LoadMusicStream")
	purego.RegisterLibFunc(&loadMusicStreamFromMemory, raylibDll, "LoadMusicStreamFromMemory")
	purego.RegisterLibFunc(&isMusicValid, raylibDll, "IsMusicValid")
	purego.RegisterLibFunc(&unloadMusicStream, raylibDll, "UnloadMusicStream")
	purego.RegisterLibFunc(&playMusicStream, raylibDll, "PlayMusicStream")
	purego.RegisterLibFunc(&isMusicStreamPlaying, raylibDll, "IsMusicStreamPlaying")
	purego.RegisterLibFunc(&updateMusicStream, raylibDll, "UpdateMusicStream")
	purego.RegisterLibFunc(&stopMusicStream, raylibDll, "StopMusicStream")
	purego.RegisterLibFunc(&pauseMusicStream, raylibDll, "PauseMusicStream")
	purego.RegisterLibFunc(&resumeMusicStream, raylibDll, "ResumeMusicStream")
	purego.RegisterLibFunc(&seekMusicStream, raylibDll, "SeekMusicStream")
	purego.RegisterLibFunc(&setMusicVolume, raylibDll, "SetMusicVolume")
	purego.RegisterLibFunc(&setMusicPitch, raylibDll, "SetMusicPitch")
	purego.RegisterLibFunc(&setMusicPan, raylibDll, "SetMusicPan")
	purego.RegisterLibFunc(&getMusicTimeLength, raylibDll, "GetMusicTimeLength")
	purego.RegisterLibFunc(&getMusicTimePlayed, raylibDll, "GetMusicTimePlayed")
	purego.RegisterLibFunc(&loadAudioStream, raylibDll, "LoadAudioStream")
	purego.RegisterLibFunc(&isAudioStreamValid, raylibDll, "IsAudioStreamValid")
	purego.RegisterLibFunc(&unloadAudioStream, raylibDll, "UnloadAudioStream")
	purego.RegisterLibFunc(&updateAudioStream, raylibDll, "UpdateAudioStream")
	purego.RegisterLibFunc(&isAudioStreamProcessed, raylibDll, "IsAudioStreamProcessed")
	purego.RegisterLibFunc(&playAudioStream, raylibDll, "PlayAudioStream")
	purego.RegisterLibFunc(&pauseAudioStream, raylibDll, "PauseAudioStream")
	purego.RegisterLibFunc(&resumeAudioStream, raylibDll, "ResumeAudioStream")
	purego.RegisterLibFunc(&isAudioStreamPlaying, raylibDll, "IsAudioStreamPlaying")
	purego.RegisterLibFunc(&stopAudioStream, raylibDll, "StopAudioStream")
	purego.RegisterLibFunc(&setAudioStreamVolume, raylibDll, "SetAudioStreamVolume")
	purego.RegisterLibFunc(&setAudioStreamPitch, raylibDll, "SetAudioStreamPitch")
	purego.RegisterLibFunc(&setAudioStreamPan, raylibDll, "SetAudioStreamPan")
	purego.RegisterLibFunc(&setAudioStreamBufferSizeDefault, raylibDll, "SetAudioStreamBufferSizeDefault")
	purego.RegisterLibFunc(&setAudioStreamCallback, raylibDll, "SetAudioStreamCallback")
	purego.RegisterLibFunc(&attachAudioStreamProcessor, raylibDll, "AttachAudioStreamProcessor")
	purego.RegisterLibFunc(&detachAudioStreamProcessor, raylibDll, "DetachAudioStreamProcessor")
	purego.RegisterLibFunc(&attachAudioMixedProcessor, raylibDll, "AttachAudioMixedProcessor")
	purego.RegisterLibFunc(&detachAudioMixedProcessor, raylibDll, "DetachAudioMixedProcessor")
}

// InitWindow - Initialize window and OpenGL context
func InitWindow(width int32, height int32, title string) {
	initWindow(width, height, title)
}

// CloseWindow - Close window and unload OpenGL context
func CloseWindow() {
	closeWindow()
}

// WindowShouldClose - Check if application should close (KEY_ESCAPE pressed or windows close icon clicked)
func WindowShouldClose() bool {
	return windowShouldClose()
}

// IsWindowReady - Check if window has been initialized successfully
func IsWindowReady() bool {
	return isWindowReady()
}

// IsWindowFullscreen - Check if window is currently fullscreen
func IsWindowFullscreen() bool {
	return isWindowFullscreen()
}

// IsWindowHidden - Check if window is currently hidden (only PLATFORM_DESKTOP)
func IsWindowHidden() bool {
	return isWindowHidden()
}

// IsWindowMinimized - Check if window is currently minimized (only PLATFORM_DESKTOP)
func IsWindowMinimized() bool {
	return isWindowMinimized()
}

// IsWindowMaximized - Check if window is currently maximized (only PLATFORM_DESKTOP)
func IsWindowMaximized() bool {
	return isWindowMaximized()
}

// IsWindowFocused - Check if window is currently focused (only PLATFORM_DESKTOP)
func IsWindowFocused() bool {
	return isWindowFocused()
}

// IsWindowResized - Check if window has been resized last frame
func IsWindowResized() bool {
	return isWindowResized()
}

// IsWindowState - Check if one specific window flag is enabled
func IsWindowState(flag uint32) bool {
	return isWindowState(flag)
}

// SetWindowState - Set window configuration state using flags (only PLATFORM_DESKTOP)
func SetWindowState(flags uint32) {
	setWindowState(flags)
}

// ClearWindowState - Clear window configuration state flags
func ClearWindowState(flags uint32) {
	clearWindowState(flags)
}

// ToggleFullscreen - Toggle window state: fullscreen/windowed (only PLATFORM_DESKTOP)
func ToggleFullscreen() {
	toggleFullscreen()
}

// ToggleBorderlessWindowed - Toggle window state: borderless windowed (only PLATFORM_DESKTOP)
func ToggleBorderlessWindowed() {
	toggleBorderlessWindowed()
}

// MaximizeWindow - Set window state: maximized, if resizable (only PLATFORM_DESKTOP)
func MaximizeWindow() {
	maximizeWindow()
}

// MinimizeWindow - Set window state: minimized, if resizable (only PLATFORM_DESKTOP)
func MinimizeWindow() {
	minimizeWindow()
}

// RestoreWindow - Set window state: not minimized/maximized (only PLATFORM_DESKTOP)
func RestoreWindow() {
	restoreWindow()
}

// SetWindowIcon - Set icon for window (single image, RGBA 32bit, only PLATFORM_DESKTOP)
func SetWindowIcon(image Image) {
	setWindowIcon(uintptr(unsafe.Pointer(&image)))
}

// SetWindowIcons - Set icon for window (multiple images, RGBA 32bit, only PLATFORM_DESKTOP)
func SetWindowIcons(images []Image, count int32) {
	setWindowIcons(uintptr(unsafe.Pointer(&images[0])), int32(len(images)))
}

// SetWindowTitle - Set title for window (only PLATFORM_DESKTOP and PLATFORM_WEB)
func SetWindowTitle(title string) {
	setWindowTitle(title)
}

// SetWindowPosition - Set window position on screen (only PLATFORM_DESKTOP)
func SetWindowPosition(x int, y int) {
	setWindowPosition(int32(x), int32(y))
}

// SetWindowMonitor - Set monitor for the current window
func SetWindowMonitor(monitor int) {
	setWindowMonitor(int32(monitor))
}

// SetWindowMinSize - Set window minimum dimensions (for FLAG_WINDOW_RESIZABLE)
func SetWindowMinSize(width int, height int) {
	setWindowMinSize(int32(width), int32(height))
}

// SetWindowMaxSize - Set window maximum dimensions (for FLAG_WINDOW_RESIZABLE)
func SetWindowMaxSize(width int, height int) {
	setWindowMaxSize(int32(width), int32(height))
}

// SetWindowSize - Set window dimensions
func SetWindowSize(width int, height int) {
	setWindowSize(int32(width), int32(height))
}

// SetWindowOpacity - Set window opacity [0.0f..1.0f] (only PLATFORM_DESKTOP)
func SetWindowOpacity(opacity float32) {
	setWindowOpacity(opacity)
}

// SetWindowFocused - Set window focused (only PLATFORM_DESKTOP)
func SetWindowFocused() {
	setWindowFocused()
}

// GetWindowHandle - Get native window handle
func GetWindowHandle() unsafe.Pointer {
	return getWindowHandle()
}

// GetScreenWidth - Get current screen width
func GetScreenWidth() int {
	return int(getScreenWidth())
}

// GetScreenHeight - Get current screen height
func GetScreenHeight() int {
	return int(getScreenHeight())
}

// GetRenderWidth - Get current render width (it considers HiDPI)
func GetRenderWidth() int {
	return int(getRenderWidth())
}

// GetRenderHeight - Get current render height (it considers HiDPI)
func GetRenderHeight() int {
	return int(getRenderHeight())
}

// GetMonitorCount - Get number of connected monitors
func GetMonitorCount() int {
	return int(getMonitorCount())
}

// GetCurrentMonitor - Get current monitor where window is placed
func GetCurrentMonitor() int {
	return int(getCurrentMonitor())
}

// GetMonitorPosition - Get specified monitor position
func GetMonitorPosition(monitor int) Vector2 {
	ret := getMonitorPosition(int32(monitor))
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetMonitorWidth - Get specified monitor width (current video mode used by monitor)
func GetMonitorWidth(monitor int) int {
	return int(getMonitorWidth(int32(monitor)))
}

// GetMonitorHeight - Get specified monitor height (current video mode used by monitor)
func GetMonitorHeight(monitor int) int {
	return int(getMonitorHeight(int32(monitor)))
}

// GetMonitorPhysicalWidth - Get specified monitor physical width in millimetres
func GetMonitorPhysicalWidth(monitor int) int {
	return int(getMonitorPhysicalWidth(int32(monitor)))
}

// GetMonitorPhysicalHeight - Get specified monitor physical height in millimetres
func GetMonitorPhysicalHeight(monitor int) int {
	return int(getMonitorPhysicalHeight(int32(monitor)))
}

// GetMonitorRefreshRate - Get specified monitor refresh rate
func GetMonitorRefreshRate(monitor int) int {
	return int(getMonitorRefreshRate(int32(monitor)))
}

// GetWindowPosition - Get window position XY on monitor
func GetWindowPosition() Vector2 {
	ret := getWindowPosition()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetWindowScaleDPI - Get window scale DPI factor
func GetWindowScaleDPI() Vector2 {
	ret := getWindowScaleDPI()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetMonitorName - Get the human-readable, UTF-8 encoded name of the specified monitor
func GetMonitorName(monitor int) string {
	return getMonitorName(int32(monitor))
}

// SetClipboardText - Set clipboard text content
func SetClipboardText(text string) {
	setClipboardText(text)
}

// GetClipboardText - Get clipboard text content
func GetClipboardText() string {
	return getClipboardText()
}

// GetClipboardImage - Get clipboard image content
//
// Only works with SDL3 backend or Windows with RGFW/GLFW
func GetClipboardImage() Image {
	var img Image
	getClipboardImage(uintptr(unsafe.Pointer(&img)))
	return img
}

// EnableEventWaiting - Enable waiting for events on EndDrawing(), no automatic event polling
func EnableEventWaiting() {
	enableEventWaiting()
}

// DisableEventWaiting - Disable waiting for events on EndDrawing(), automatic events polling
func DisableEventWaiting() {
	disableEventWaiting()
}

// ShowCursor - Shows cursor
func ShowCursor() {
	showCursor()
}

// HideCursor - Hides cursor
func HideCursor() {
	hideCursor()
}

// IsCursorHidden - Check if cursor is not visible
func IsCursorHidden() bool {
	return isCursorHidden()
}

// EnableCursor - Enables cursor (unlock cursor)
func EnableCursor() {
	enableCursor()
}

// DisableCursor - Disables cursor (lock cursor)
func DisableCursor() {
	disableCursor()
}

// IsCursorOnScreen - Check if cursor is on the screen
func IsCursorOnScreen() bool {
	return isCursorOnScreen()
}

// ClearBackground - Set background color (framebuffer clear color)
func ClearBackground(col color.RGBA) {
	clearBackground(*(*uintptr)(unsafe.Pointer(&col)))
}

// BeginDrawing - Setup canvas (framebuffer) to start drawing
func BeginDrawing() {
	beginDrawing()
}

// EndDrawing - End canvas drawing and swap buffers (double buffering)
func EndDrawing() {
	endDrawing()
}

// BeginMode2D - Begin 2D mode with custom camera (2D)
func BeginMode2D(camera Camera2D) {
	beginMode2D(uintptr(unsafe.Pointer(&camera)))
}

// EndMode2D - Ends 2D mode with custom camera
func EndMode2D() {
	endMode2D()
}

// BeginMode3D - Begin 3D mode with custom camera (3D)
func BeginMode3D(camera Camera3D) {
	beginMode3D(uintptr(unsafe.Pointer(&camera)))
}

// EndMode3D - Ends 3D mode and returns to default 2D orthographic mode
func EndMode3D() {
	endMode3D()
}

// BeginTextureMode - Begin drawing to render texture
func BeginTextureMode(target RenderTexture2D) {
	beginTextureMode(uintptr(unsafe.Pointer(&target)))
}

// EndTextureMode - Ends drawing to render texture
func EndTextureMode() {
	endTextureMode()
}

// BeginShaderMode - Begin custom shader drawing
func BeginShaderMode(shader Shader) {
	beginShaderMode(uintptr(unsafe.Pointer(&shader)))
}

// EndShaderMode - End custom shader drawing (use default shader)
func EndShaderMode() {
	endShaderMode()
}

// BeginBlendMode - Begin blending mode (alpha, additive, multiplied, subtract, custom)
func BeginBlendMode(mode BlendMode) {
	beginBlendMode(int32(mode))
}

// EndBlendMode - End blending mode (reset to default: alpha blending)
func EndBlendMode() {
	endBlendMode()
}

// BeginScissorMode - Begin scissor mode (define screen area for following drawing)
func BeginScissorMode(x int32, y int32, width int32, height int32) {
	beginScissorMode(x, y, width, height)
}

// EndScissorMode - End scissor mode
func EndScissorMode() {
	endScissorMode()
}

// BeginVrStereoMode - Begin stereo rendering (requires VR simulator)
func BeginVrStereoMode(config VrStereoConfig) {
	beginVrStereoMode(uintptr(unsafe.Pointer(&config)))
}

// EndVrStereoMode - End stereo rendering (requires VR simulator)
func EndVrStereoMode() {
	endVrStereoMode()
}

// LoadVrStereoConfig - Load VR stereo config for VR simulator device parameters
func LoadVrStereoConfig(device VrDeviceInfo) VrStereoConfig {
	var config VrStereoConfig
	loadVrStereoConfig(uintptr(unsafe.Pointer(&config)), uintptr(unsafe.Pointer(&device)))
	return config
}

// UnloadVrStereoConfig - Unload VR stereo config
func UnloadVrStereoConfig(config VrStereoConfig) {
	unloadVrStereoConfig(uintptr(unsafe.Pointer(&config)))
}

// LoadShader - Load shader from files and bind default locations
func LoadShader(vsFileName string, fsFileName string) Shader {
	var shader Shader
	var cvsFileName, cfsFileName *byte
	if vsFileName != "" {
		var err error
		cvsFileName, err = windows.BytePtrFromString(vsFileName)
		if err != nil {
			panic(err)
		}
	}
	if fsFileName != "" {
		var err error
		cfsFileName, err = windows.BytePtrFromString(fsFileName)
		if err != nil {
			panic(err)
		}
	}
	loadShader(uintptr(unsafe.Pointer(&shader)), uintptr(unsafe.Pointer(cvsFileName)), uintptr(unsafe.Pointer(cfsFileName)))
	return shader
}

// LoadShaderFromMemory - Load shader from code strings and bind default locations
func LoadShaderFromMemory(vsCode string, fsCode string) Shader {
	var shader Shader
	var cvsCode, cfsCode *byte
	if vsCode != "" {
		var err error
		cvsCode, err = windows.BytePtrFromString(vsCode)
		if err != nil {
			panic(err)
		}
	}
	if fsCode != "" {
		var err error
		cfsCode, err = windows.BytePtrFromString(fsCode)
		if err != nil {
			panic(err)
		}
	}
	loadShaderFromMemory(uintptr(unsafe.Pointer(&shader)), uintptr(unsafe.Pointer(cvsCode)), uintptr(unsafe.Pointer(cfsCode)))
	return shader
}

// IsShaderValid - Check if a shader is valid (loaded on GPU)
func IsShaderValid(shader Shader) bool {
	return isShaderValid(uintptr(unsafe.Pointer(&shader)))
}

// GetShaderLocation - Get shader uniform location
func GetShaderLocation(shader Shader, uniformName string) int32 {
	return getShaderLocation(uintptr(unsafe.Pointer(&shader)), uniformName)
}

// GetShaderLocationAttrib - Get shader attribute location
func GetShaderLocationAttrib(shader Shader, attribName string) int32 {
	return getShaderLocationAttrib(uintptr(unsafe.Pointer(&shader)), attribName)
}

// SetShaderValue - Set shader uniform value
func SetShaderValue(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType) {
	setShaderValue(uintptr(unsafe.Pointer(&shader)), locIndex, value, int32(uniformType))
}

// SetShaderValueV - Set shader uniform value vector
func SetShaderValueV(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType, count int32) {
	setShaderValueV(uintptr(unsafe.Pointer(&shader)), locIndex, value, int32(uniformType), count)
}

// SetShaderValueMatrix - Set shader uniform value (matrix 4x4)
func SetShaderValueMatrix(shader Shader, locIndex int32, mat Matrix) {
	setShaderValueMatrix(uintptr(unsafe.Pointer(&shader)), locIndex, uintptr(unsafe.Pointer(&mat)))
}

// SetShaderValueTexture - Set shader uniform value for texture (sampler2d)
func SetShaderValueTexture(shader Shader, locIndex int32, texture Texture2D) {
	setShaderValueTexture(uintptr(unsafe.Pointer(&shader)), locIndex, uintptr(unsafe.Pointer(&texture)))
}

// UnloadShader - Unload shader from GPU memory (VRAM)
func UnloadShader(shader Shader) {
	unloadShader(uintptr(unsafe.Pointer(&shader)))
}

// GetMouseRay - Get a ray trace from mouse position
//
// Deprecated: Use [GetScreenToWorldRay] instead.
func GetMouseRay(mousePosition Vector2, camera Camera) Ray {
	return GetScreenToWorldRay(mousePosition, camera)
}

// GetScreenToWorldRay - Get a ray trace from screen position (i.e mouse)
func GetScreenToWorldRay(position Vector2, camera Camera) Ray {
	var ray Ray
	getScreenToWorldRay(uintptr(unsafe.Pointer(&ray)), *(*uintptr)(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&camera)))
	return ray
}

// GetScreenToWorldRayEx - Get a ray trace from screen position (i.e mouse) in a viewport
func GetScreenToWorldRayEx(position Vector2, camera Camera, width, height int32) Ray {
	var ray Ray
	getScreenToWorldRayEx(uintptr(unsafe.Pointer(&ray)), *(*uintptr)(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&camera)), width, height)
	return ray
}

// GetCameraMatrix - Get camera transform matrix (view matrix)
func GetCameraMatrix(camera Camera) Matrix {
	var mat Matrix
	getCameraMatrix(uintptr(unsafe.Pointer(&mat)), uintptr(unsafe.Pointer(&camera)))
	return mat
}

// GetCameraMatrix2D - Get camera 2d transform matrix
func GetCameraMatrix2D(camera Camera2D) Matrix {
	var mat Matrix
	getCameraMatrix2D(uintptr(unsafe.Pointer(&mat)), uintptr(unsafe.Pointer(&camera)))
	return mat
}

// GetWorldToScreen - Get the screen space position for a 3d world space position
func GetWorldToScreen(position Vector3, camera Camera) Vector2 {
	ret := getWorldToScreen(uintptr(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&camera)))
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetScreenToWorld2D - Get the world space position for a 2d camera screen space position
func GetScreenToWorld2D(position Vector2, camera Camera2D) Vector2 {
	ret := getScreenToWorld2D(*(*uintptr)(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&camera)))
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetWorldToScreenEx - Get size position for a 3d world space position
func GetWorldToScreenEx(position Vector3, camera Camera, width int32, height int32) Vector2 {
	ret := getWorldToScreenEx(uintptr(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&camera)), width, height)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetWorldToScreen2D - Get the screen space position for a 2d camera world space position
func GetWorldToScreen2D(position Vector2, camera Camera2D) Vector2 {
	ret := getWorldToScreen2D(*(*uintptr)(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&camera)))
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// SetTargetFPS - Set target FPS (maximum)
func SetTargetFPS(fps int32) {
	setTargetFPS(fps)
}

// GetFrameTime - Get time in seconds for last frame drawn (delta time)
func GetFrameTime() float32 {
	return getFrameTime()
}

// GetTime - Get elapsed time in seconds since InitWindow()
func GetTime() float64 {
	return getTime()
}

// GetFPS - Get current FPS
func GetFPS() int32 {
	return getFPS()
}

// Custom frame control functions
// NOTE: SwapScreenBuffer and PollInputEvents are intended for advanced users that want full control over the frame processing
// By default EndDrawing() does this job: draws everything + SwapScreenBuffer() + manage frame timing + PollInputEvents()
// To avoid that behaviour and control frame processes manually you must recompile raylib with SUPPORT_CUSTOM_FRAME_CONTROL enabled in config.h

// SwapScreenBuffer - Swap back buffer with front buffer (screen drawing)
func SwapScreenBuffer() {
	swapScreenBuffer()
}

// PollInputEvents - Register all input events
func PollInputEvents() {
	pollInputEvents()
}

// WaitTime - Wait for some time (halt program execution)
func WaitTime(seconds float64) {
	waitTime(seconds)
}

// SetRandomSeed - Set the seed for the random number generator
//
// Note: You can use go's math/rand package instead
func SetRandomSeed(seed uint32) {
	setRandomSeed(seed)
}

// GetRandomValue - Get a random value between min and max (both included)
//
// Note: You can use go's math/rand package instead
func GetRandomValue(minimum int32, maximum int32) int32 {
	return getRandomValue(minimum, maximum)
}

// LoadRandomSequence - Load random values sequence, no values repeated
//
// Note: Use UnloadRandomSequence if you don't need the sequence any more. You can use go's math/rand.Perm function instead.
func LoadRandomSequence(count uint32, minimum int32, maximum int32) []int32 {
	ret := loadRandomSequence(count, minimum, maximum)
	return unsafe.Slice(ret, 10)
}

// UnloadRandomSequence - Unload random values sequence
func UnloadRandomSequence(sequence []int32) {
	unloadRandomSequence(unsafe.SliceData(sequence))
}

// TakeScreenshot - Takes a screenshot of current screen (filename extension defines format)
func TakeScreenshot(fileName string) {
	takeScreenshot(fileName)
}

// SetConfigFlags - Setup init configuration flags (view FLAGS)
func SetConfigFlags(flags uint32) {
	setConfigFlags(flags)
}

// OpenURL - Open URL with default system browser (if available)
func OpenURL(url string) {
	openURL(url)
}

// TraceLog - Show trace log messages (LOG_DEBUG, LOG_INFO, LOG_WARNING, LOG_ERROR...)
func TraceLog(logLevel TraceLogLevel, text string, args ...any) {
	traceLog(int32(logLevel), fmt.Sprintf(text, args...))
}

// SetTraceLogLevel - Set the current threshold (minimum) log level
func SetTraceLogLevel(logLevel TraceLogLevel) {
	setTraceLogLevel(int32(logLevel))
}

// MemAlloc - Internal memory allocator
func MemAlloc(size uint32) unsafe.Pointer {
	return memAlloc(size)
}

// MemRealloc - Internal memory reallocator
func MemRealloc(ptr unsafe.Pointer, size uint32) unsafe.Pointer {
	return memRealloc(ptr, size)
}

// MemFree - Internal memory free
func MemFree(ptr unsafe.Pointer) {
	memFree(ptr)
}

// SetTraceLogCallback - Set custom trace log
func SetTraceLogCallback(fn TraceLogCallbackFun) {
	setTraceLogCallback(traceLogCallbackWrapper(fn))
}

// IsFileDropped - Check if a file has been dropped into window
func IsFileDropped() bool {
	return isFileDropped()
}

// LoadDroppedFiles - Load dropped filepaths
func LoadDroppedFiles() []string {
	var filePathList = struct {
		capacity uint32
		count    uint32
		paths    **byte
	}{}
	loadDroppedFiles(uintptr(unsafe.Pointer(&filePathList)))
	defer unloadDroppedFiles(uintptr(unsafe.Pointer(&filePathList)))

	tmpslice := (*[1 << 24]*byte)(unsafe.Pointer(filePathList.paths))[:filePathList.count:filePathList.count]

	gostrings := make([]string, filePathList.count)
	for i, s := range tmpslice {
		gostrings[i] = func(p *byte) string {
			if p == nil || *p == 0 {
				return ""
			}

			n := 0
			for ptr := unsafe.Pointer(p); *(*byte)(ptr) != 0; n++ {
				ptr = unsafe.Pointer(uintptr(ptr) + 1)
			}

			return string(unsafe.Slice(p, n))
		}(s)
	}

	return gostrings
}

// UnloadDroppedFiles - Unload dropped filepaths
func UnloadDroppedFiles() {}

// LoadAutomationEventList - Load automation events list from file, NULL for empty list, capacity = MAX_AUTOMATION_EVENTS
func LoadAutomationEventList(fileName string) AutomationEventList {
	var automationEventList AutomationEventList
	loadAutomationEventList(uintptr(unsafe.Pointer(&automationEventList)), fileName)
	return automationEventList
}

// UnloadAutomationEventList - Unload automation events list from file
func UnloadAutomationEventList(list *AutomationEventList) {
	unloadAutomationEventList(uintptr(unsafe.Pointer(&list)))
}

// ExportAutomationEventList - Export automation events list as text file
func ExportAutomationEventList(list AutomationEventList, fileName string) bool {
	return exportAutomationEventList(uintptr(unsafe.Pointer(&list)), fileName)
}

// SetAutomationEventList - Set automation event list to record to
func SetAutomationEventList(list *AutomationEventList) {
	setAutomationEventList(uintptr(unsafe.Pointer(&list)))
}

// SetAutomationEventBaseFrame - Set automation event internal base frame to start recording
func SetAutomationEventBaseFrame(frame int) {
	setAutomationEventBaseFrame(int32(frame))
}

// StartAutomationEventRecording - Start recording automation events (AutomationEventList must be set)
func StartAutomationEventRecording() {
	startAutomationEventRecording()
}

// StopAutomationEventRecording - Stop recording automation events
func StopAutomationEventRecording() {
	stopAutomationEventRecording()
}

// PlayAutomationEvent - Play a recorded automation event
func PlayAutomationEvent(event AutomationEvent) {
	playAutomationEvent(uintptr(unsafe.Pointer(&event)))
}

// IsKeyPressed - Check if a key has been pressed once
func IsKeyPressed(key int32) bool {
	return isKeyPressed(key)
}

// IsKeyPressedRepeat - Check if a key has been pressed again (Only PLATFORM_DESKTOP)
func IsKeyPressedRepeat(key int32) bool {
	return isKeyPressedRepeat(key)
}

// IsKeyDown - Check if a key is being pressed
func IsKeyDown(key int32) bool {
	return isKeyDown(key)
}

// IsKeyReleased - Check if a key has been released once
func IsKeyReleased(key int32) bool {
	return isKeyReleased(key)
}

// IsKeyUp - Check if a key is NOT being pressed
func IsKeyUp(key int32) bool {
	return isKeyUp(key)
}

// GetKeyPressed - Get key pressed (keycode), call it multiple times for keys queued, returns 0 when the queue is empty
func GetKeyPressed() int32 {
	return getKeyPressed()
}

// GetCharPressed - Get char pressed (unicode), call it multiple times for chars queued, returns 0 when the queue is empty
func GetCharPressed() int32 {
	return getCharPressed()
}

// SetExitKey - Set a custom key to exit program (default is ESC)
func SetExitKey(key int32) {
	setExitKey(key)
}

// IsGamepadAvailable - Check if a gamepad is available
func IsGamepadAvailable(gamepad int32) bool {
	return isGamepadAvailable(gamepad)
}

// GetGamepadName - Get gamepad internal name id
func GetGamepadName(gamepad int32) string {
	return getGamepadName(gamepad)
}

// IsGamepadButtonPressed - Check if a gamepad button has been pressed once
func IsGamepadButtonPressed(gamepad int32, button int32) bool {
	return isGamepadButtonPressed(gamepad, button)
}

// IsGamepadButtonDown - Check if a gamepad button is being pressed
func IsGamepadButtonDown(gamepad int32, button int32) bool {
	return isGamepadButtonDown(gamepad, button)
}

// IsGamepadButtonReleased - Check if a gamepad button has been released once
func IsGamepadButtonReleased(gamepad int32, button int32) bool {
	return isGamepadButtonReleased(gamepad, button)
}

// IsGamepadButtonUp - Check if a gamepad button is NOT being pressed
func IsGamepadButtonUp(gamepad int32, button int32) bool {
	return isGamepadButtonUp(gamepad, button)
}

// GetGamepadButtonPressed - Get the last gamepad button pressed
func GetGamepadButtonPressed() int32 {
	return getGamepadButtonPressed()
}

// GetGamepadAxisCount - Get gamepad axis count for a gamepad
func GetGamepadAxisCount(gamepad int32) int32 {
	return getGamepadAxisCount(gamepad)
}

// GetGamepadAxisMovement - Get axis movement value for a gamepad axis
func GetGamepadAxisMovement(gamepad int32, axis int32) float32 {
	return getGamepadAxisMovement(gamepad, axis)
}

// SetGamepadMappings - Set internal gamepad mappings (SDL_GameControllerDB)
func SetGamepadMappings(mappings string) int32 {
	return setGamepadMappings(mappings)
}

// SetGamepadVibration - Set gamepad vibration for both motors (duration in seconds)
func SetGamepadVibration(gamepad int32, leftMotor, rightMotor, duration float32) {
	setGamepadVibration(gamepad, leftMotor, rightMotor, duration)
}

// IsMouseButtonPressed - Check if a mouse button has been pressed once
func IsMouseButtonPressed(button MouseButton) bool {
	return isMouseButtonPressed(int32(button))
}

// IsMouseButtonDown - Check if a mouse button is being pressed
func IsMouseButtonDown(button MouseButton) bool {
	return isMouseButtonDown(int32(button))
}

// IsMouseButtonReleased - Check if a mouse button has been released once
func IsMouseButtonReleased(button MouseButton) bool {
	return isMouseButtonReleased(int32(button))
}

// IsMouseButtonUp - Check if a mouse button is NOT being pressed
func IsMouseButtonUp(button MouseButton) bool {
	return isMouseButtonUp(int32(button))
}

// GetMouseX - Get mouse position X
func GetMouseX() int32 {
	return getMouseX()
}

// GetMouseY - Get mouse position Y
func GetMouseY() int32 {
	return getMouseY()
}

// GetMousePosition - Get mouse position XY
func GetMousePosition() Vector2 {
	ret := getMousePosition()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetMouseDelta - Get mouse delta between frames
func GetMouseDelta() Vector2 {
	ret := getMouseDelta()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// SetMousePosition - Set mouse position XY
func SetMousePosition(x int32, y int32) {
	setMousePosition(x, y)
}

// SetMouseOffset - Set mouse offset
func SetMouseOffset(offsetX int32, offsetY int32) {
	setMouseOffset(offsetX, offsetY)
}

// SetMouseScale - Set mouse scaling
func SetMouseScale(scaleX float32, scaleY float32) {
	setMouseScale(scaleX, scaleY)
}

// GetMouseWheelMove - Get mouse wheel movement for X or Y, whichever is larger
func GetMouseWheelMove() float32 {
	return getMouseWheelMove()
}

// GetMouseWheelMoveV - Get mouse wheel movement for both X and Y
func GetMouseWheelMoveV() Vector2 {
	ret := getMouseWheelMoveV()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// SetMouseCursor - Set mouse cursor
func SetMouseCursor(cursor int32) {
	setMouseCursor(cursor)
}

// GetTouchX - Get touch position X for touch point 0 (relative to screen size)
func GetTouchX() int32 {
	return getTouchX()
}

// GetTouchY - Get touch position Y for touch point 0 (relative to screen size)
func GetTouchY() int32 {
	return getTouchY()
}

// GetTouchPosition - Get touch position XY for a touch point index (relative to screen size)
func GetTouchPosition(index int32) Vector2 {
	ret := getTouchPosition(index)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetTouchPointId - Get touch point identifier for given index
func GetTouchPointId(index int32) int32 {
	return getTouchPointId(index)
}

// GetTouchPointCount - Get number of touch points
func GetTouchPointCount() int32 {
	return getTouchPointCount()
}

// SetGesturesEnabled - Enable a set of gestures using flags
func SetGesturesEnabled(flags uint32) {
	setGesturesEnabled(flags)
}

// IsGestureDetected - Check if a gesture have been detected
func IsGestureDetected(gesture Gestures) bool {
	return isGestureDetected(uint32(gesture))
}

// GetGestureDetected - Get latest detected gesture
func GetGestureDetected() Gestures {
	return Gestures(getGestureDetected())
}

// GetGestureHoldDuration - Get gesture hold time in milliseconds
func GetGestureHoldDuration() float32 {
	return getGestureHoldDuration()
}

// GetGestureDragVector - Get gesture drag vector
func GetGestureDragVector() Vector2 {
	ret := getGestureDragVector()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetGestureDragAngle - Get gesture drag angle
func GetGestureDragAngle() float32 {
	return getGestureDragAngle()
}

// GetGesturePinchVector - Get gesture pinch delta
func GetGesturePinchVector() Vector2 {
	ret := getGesturePinchVector()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetGesturePinchAngle - Get gesture pinch angle
func GetGesturePinchAngle() float32 {
	return getGesturePinchAngle()
}

// SetShapesTexture - Set texture and rectangle to be used on shapes drawing
func SetShapesTexture(texture Texture2D, source Rectangle) {
	setShapesTexture(uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(&source)))
}

// GetShapesTexture - Get texture that is used for shapes drawing
func GetShapesTexture() Texture2D {
	var texture Texture2D
	getShapesTexture(uintptr(unsafe.Pointer(&texture)))
	return texture
}

// GetShapesTextureRectangle - Get texture source rectangle that is used for shapes drawing
func GetShapesTextureRectangle() Rectangle {
	var rec Rectangle
	getShapesTextureRectangle(uintptr(unsafe.Pointer(&rec)))
	return rec
}

// DrawPixel - Draw a pixel
func DrawPixel(posX int32, posY int32, col color.RGBA) {
	drawPixel(posX, posY, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawPixelV - Draw a pixel (Vector version)
func DrawPixelV(position Vector2, col color.RGBA) {
	drawPixelV(*(*uintptr)(unsafe.Pointer(&position)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawLine - Draw a line
func DrawLine(startPosX int32, startPosY int32, endPosX int32, endPosY int32, col color.RGBA) {
	drawLine(startPosX, startPosY, endPosX, endPosY, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawLineV - Draw a line (using gl lines)
func DrawLineV(startPos Vector2, endPos Vector2, col color.RGBA) {
	drawLineV(*(*uintptr)(unsafe.Pointer(&startPos)), *(*uintptr)(unsafe.Pointer(&endPos)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawLineEx - Draw a line (using triangles/quads)
func DrawLineEx(startPos Vector2, endPos Vector2, thick float32, col color.RGBA) {
	drawLineEx(*(*uintptr)(unsafe.Pointer(&startPos)), *(*uintptr)(unsafe.Pointer(&endPos)), thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawLineStrip - Draw lines sequence (using gl lines)
func DrawLineStrip(points []Vector2, col color.RGBA) {
	pointCount := int32(len(points))
	drawLineStrip((unsafe.SliceData(points)), pointCount, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawLineBezier - Draw line segment cubic-bezier in-out interpolation
func DrawLineBezier(startPos Vector2, endPos Vector2, thick float32, col color.RGBA) {
	drawLineBezier(*(*uintptr)(unsafe.Pointer(&startPos)), *(*uintptr)(unsafe.Pointer(&endPos)), thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCircle - Draw a color-filled circle
func DrawCircle(centerX int32, centerY int32, radius float32, col color.RGBA) {
	drawCircle(centerX, centerY, radius, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCircleSector - Draw a piece of a circle
func DrawCircleSector(center Vector2, radius float32, startAngle float32, endAngle float32, segments int32, col color.RGBA) {
	drawCircleSector(*(*uintptr)(unsafe.Pointer(&center)), radius, startAngle, endAngle, segments, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCircleSectorLines - Draw circle sector outline
func DrawCircleSectorLines(center Vector2, radius float32, startAngle float32, endAngle float32, segments int32, col color.RGBA) {
	drawCircleSectorLines(*(*uintptr)(unsafe.Pointer(&center)), radius, startAngle, endAngle, segments, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCircleGradient - Draw a gradient-filled circle
func DrawCircleGradient(centerX int32, centerY int32, radius float32, inner color.RGBA, outer color.RGBA) {
	drawCircleGradient(centerX, centerY, radius, *(*uintptr)(unsafe.Pointer(&inner)), *(*uintptr)(unsafe.Pointer(&outer)))
}

// DrawCircleV - Draw a color-filled circle (Vector version)
func DrawCircleV(center Vector2, radius float32, col color.RGBA) {
	drawCircleV(*(*uintptr)(unsafe.Pointer(&center)), radius, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCircleLines - Draw circle outline
func DrawCircleLines(centerX int32, centerY int32, radius float32, col color.RGBA) {
	drawCircleLines(centerX, centerY, radius, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCircleLinesV - Draw circle outline (Vector version)
func DrawCircleLinesV(center Vector2, radius float32, col color.RGBA) {
	drawCircleLinesV(*(*uintptr)(unsafe.Pointer(&center)), radius, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawEllipse - Draw ellipse
func DrawEllipse(centerX int32, centerY int32, radiusH float32, radiusV float32, col color.RGBA) {
	drawEllipse(centerX, centerY, radiusH, radiusV, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawEllipseLines - Draw ellipse outline
func DrawEllipseLines(centerX int32, centerY int32, radiusH float32, radiusV float32, col color.RGBA) {
	drawEllipseLines(centerX, centerY, radiusH, radiusV, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRing - Draw ring
func DrawRing(center Vector2, innerRadius float32, outerRadius float32, startAngle float32, endAngle float32, segments int32, col color.RGBA) {
	drawRing(*(*uintptr)(unsafe.Pointer(&center)), innerRadius, outerRadius, startAngle, endAngle, segments, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRingLines - Draw ring outline
func DrawRingLines(center Vector2, innerRadius float32, outerRadius float32, startAngle float32, endAngle float32, segments int32, col color.RGBA) {
	drawRingLines(*(*uintptr)(unsafe.Pointer(&center)), innerRadius, outerRadius, startAngle, endAngle, segments, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRectangle - Draw a color-filled rectangle
func DrawRectangle(posX int32, posY int32, width int32, height int32, col color.RGBA) {
	drawRectangle(posX, posY, width, height, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRectangleV - Draw a color-filled rectangle (Vector version)
func DrawRectangleV(position Vector2, size Vector2, col color.RGBA) {
	drawRectangleV(*(*uintptr)(unsafe.Pointer(&position)), *(*uintptr)(unsafe.Pointer(&size)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRectangleRec - Draw a color-filled rectangle
func DrawRectangleRec(rec Rectangle, col color.RGBA) {
	drawRectangleRec(uintptr(unsafe.Pointer(&rec)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRectanglePro - Draw a color-filled rectangle with pro parameters
func DrawRectanglePro(rec Rectangle, origin Vector2, rotation float32, col color.RGBA) {
	drawRectanglePro(uintptr(unsafe.Pointer(&rec)), *(*uintptr)(unsafe.Pointer(&origin)), rotation, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRectangleGradientV - Draw a vertical-gradient-filled rectangle
func DrawRectangleGradientV(posX int32, posY int32, width int32, height int32, top color.RGBA, bottom color.RGBA) {
	drawRectangleGradientV(posX, posY, width, height, *(*uintptr)(unsafe.Pointer(&top)), *(*uintptr)(unsafe.Pointer(&bottom)))
}

// DrawRectangleGradientH - Draw a horizontal-gradient-filled rectangle
func DrawRectangleGradientH(posX int32, posY int32, width int32, height int32, left color.RGBA, right color.RGBA) {
	drawRectangleGradientH(posX, posY, width, height, *(*uintptr)(unsafe.Pointer(&left)), *(*uintptr)(unsafe.Pointer(&right)))
}

// DrawRectangleGradientEx - Draw a gradient-filled rectangle with custom vertex colors
func DrawRectangleGradientEx(rec Rectangle, topLeft color.RGBA, bottomLeft color.RGBA, topRight color.RGBA, bottomRight color.RGBA) {
	drawRectangleGradientEx(uintptr(unsafe.Pointer(&rec)), *(*uintptr)(unsafe.Pointer(&topLeft)), *(*uintptr)(unsafe.Pointer(&bottomLeft)), *(*uintptr)(unsafe.Pointer(&topRight)), *(*uintptr)(unsafe.Pointer(&bottomRight)))
}

// DrawRectangleLines - Draw rectangle outline
func DrawRectangleLines(posX int32, posY int32, width int32, height int32, col color.RGBA) {
	drawRectangleLines(posX, posY, width, height, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRectangleLinesEx - Draw rectangle outline with extended parameters
func DrawRectangleLinesEx(rec Rectangle, lineThick float32, col color.RGBA) {
	drawRectangleLinesEx(uintptr(unsafe.Pointer(&rec)), lineThick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRectangleRounded - Draw rectangle with rounded edges
func DrawRectangleRounded(rec Rectangle, roundness float32, segments int32, col color.RGBA) {
	drawRectangleRounded(uintptr(unsafe.Pointer(&rec)), roundness, segments, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRectangleRoundedLines - Draw rectangle lines with rounded edges
func DrawRectangleRoundedLines(rec Rectangle, roundness float32, segments int32, col color.RGBA) {
	drawRectangleRoundedLines(uintptr(unsafe.Pointer(&rec)), roundness, segments, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRectangleRoundedLinesEx - Draw rectangle with rounded edges outline
func DrawRectangleRoundedLinesEx(rec Rectangle, roundness float32, segments int32, lineThick float32, col color.RGBA) {
	drawRectangleRoundedLinesEx(uintptr(unsafe.Pointer(&rec)), roundness, segments, lineThick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawTriangle - Draw a color-filled triangle (vertex in counter-clockwise order!)
func DrawTriangle(v1 Vector2, v2 Vector2, v3 Vector2, col color.RGBA) {
	drawTriangle(*(*uintptr)(unsafe.Pointer(&v1)), *(*uintptr)(unsafe.Pointer(&v2)), *(*uintptr)(unsafe.Pointer(&v3)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawTriangleLines - Draw triangle outline (vertex in counter-clockwise order!)
func DrawTriangleLines(v1 Vector2, v2 Vector2, v3 Vector2, col color.RGBA) {
	drawTriangleLines(*(*uintptr)(unsafe.Pointer(&v1)), *(*uintptr)(unsafe.Pointer(&v2)), *(*uintptr)(unsafe.Pointer(&v3)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawTriangleFan - Draw a triangle fan defined by points (first vertex is the center)
func DrawTriangleFan(points []Vector2, col color.RGBA) {
	pointCount := int32(len(points))
	drawTriangleFan(unsafe.SliceData(points), pointCount, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawTriangleStrip - Draw a triangle strip defined by points
func DrawTriangleStrip(points []Vector2, col color.RGBA) {
	pointCount := int32(len(points))
	drawTriangleStrip(unsafe.SliceData(points), pointCount, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawPoly - Draw a regular polygon (Vector version)
func DrawPoly(center Vector2, sides int32, radius float32, rotation float32, col color.RGBA) {
	drawPoly(*(*uintptr)(unsafe.Pointer(&center)), sides, radius, rotation, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawPolyLines - Draw a polygon outline of n sides
func DrawPolyLines(center Vector2, sides int32, radius float32, rotation float32, col color.RGBA) {
	drawPolyLines(*(*uintptr)(unsafe.Pointer(&center)), sides, radius, rotation, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawPolyLinesEx - Draw a polygon outline of n sides with extended parameters
func DrawPolyLinesEx(center Vector2, sides int32, radius float32, rotation float32, lineThick float32, col color.RGBA) {
	drawPolyLinesEx(*(*uintptr)(unsafe.Pointer(&center)), sides, radius, rotation, lineThick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSplineLinear - Draw spline: Linear, minimum 2 points
func DrawSplineLinear(points []Vector2, thick float32, col color.RGBA) {
	pointCount := int32(len(points))
	drawSplineLinear(unsafe.SliceData(points), pointCount, thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSplineBasis - Draw spline: B-Spline, minimum 4 points
func DrawSplineBasis(points []Vector2, thick float32, col color.RGBA) {
	pointCount := int32(len(points))
	drawSplineBasis(unsafe.SliceData(points), pointCount, thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSplineCatmullRom - Draw spline: Catmull-Rom, minimum 4 points
func DrawSplineCatmullRom(points []Vector2, thick float32, col color.RGBA) {
	pointCount := int32(len(points))
	drawSplineCatmullRom(unsafe.SliceData(points), pointCount, thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSplineBezierQuadratic - Draw spline: Quadratic Bezier, minimum 3 points (1 control point): [p1, c2, p3, c4...]
func DrawSplineBezierQuadratic(points []Vector2, thick float32, col color.RGBA) {
	pointCount := int32(len(points))
	drawSplineBezierQuadratic(unsafe.SliceData(points), pointCount, thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSplineBezierCubic - Draw spline: Cubic Bezier, minimum 4 points (2 control points): [p1, c2, c3, p4, c5, c6...]
func DrawSplineBezierCubic(points []Vector2, thick float32, col color.RGBA) {
	pointCount := int32(len(points))
	drawSplineBezierCubic(unsafe.SliceData(points), pointCount, thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSplineSegmentLinear - Draw spline segment: Linear, 2 points
func DrawSplineSegmentLinear(p1 Vector2, p2 Vector2, thick float32, col color.RGBA) {
	drawSplineSegmentLinear(*(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&p2)), thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSplineSegmentBasis - Draw spline segment: B-Spline, 4 points
func DrawSplineSegmentBasis(p1 Vector2, p2 Vector2, p3 Vector2, p4 Vector2, thick float32, col color.RGBA) {
	drawSplineSegmentBasis(*(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&p2)), *(*uintptr)(unsafe.Pointer(&p3)), *(*uintptr)(unsafe.Pointer(&p4)), thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSplineSegmentCatmullRom - Draw spline segment: Catmull-Rom, 4 points
func DrawSplineSegmentCatmullRom(p1 Vector2, p2 Vector2, p3 Vector2, p4 Vector2, thick float32, col color.RGBA) {
	drawSplineSegmentCatmullRom(*(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&p2)), *(*uintptr)(unsafe.Pointer(&p3)), *(*uintptr)(unsafe.Pointer(&p4)), thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSplineSegmentBezierQuadratic - Draw spline segment: Quadratic Bezier, 2 points, 1 control point
func DrawSplineSegmentBezierQuadratic(p1 Vector2, c2 Vector2, p3 Vector2, thick float32, col color.RGBA) {
	drawSplineSegmentBezierQuadratic(*(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&c2)), *(*uintptr)(unsafe.Pointer(&p3)), thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSplineSegmentBezierCubic - Draw spline segment: Cubic Bezier, 2 points, 2 control points
func DrawSplineSegmentBezierCubic(p1 Vector2, c2 Vector2, c3 Vector2, p4 Vector2, thick float32, col color.RGBA) {
	drawSplineSegmentBezierCubic(*(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&c2)), *(*uintptr)(unsafe.Pointer(&c3)), *(*uintptr)(unsafe.Pointer(&p4)), thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// GetSplinePointLinear - Get (evaluate) spline point: Linear
func GetSplinePointLinear(startPos Vector2, endPos Vector2, t float32) Vector2 {
	ret := getSplinePointLinear(*(*uintptr)(unsafe.Pointer(&startPos)), *(*uintptr)(unsafe.Pointer(&endPos)), t)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetSplinePointBasis - Get (evaluate) spline point: B-Spline
func GetSplinePointBasis(p1 Vector2, p2 Vector2, p3 Vector2, p4 Vector2, t float32) Vector2 {
	ret := getSplinePointBasis(*(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&p2)), *(*uintptr)(unsafe.Pointer(&p3)), *(*uintptr)(unsafe.Pointer(&p4)), t)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetSplinePointCatmullRom - Get (evaluate) spline point: Catmull-Rom
func GetSplinePointCatmullRom(p1 Vector2, p2 Vector2, p3 Vector2, p4 Vector2, t float32) Vector2 {
	ret := getSplinePointCatmullRom(*(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&p2)), *(*uintptr)(unsafe.Pointer(&p3)), *(*uintptr)(unsafe.Pointer(&p4)), t)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetSplinePointBezierQuad - Get (evaluate) spline point: Quadratic Bezier
func GetSplinePointBezierQuad(p1 Vector2, c2 Vector2, p3 Vector2, t float32) Vector2 {
	ret := getSplinePointBezierQuad(*(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&c2)), *(*uintptr)(unsafe.Pointer(&p3)), t)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetSplinePointBezierCubic - Get (evaluate) spline point: Cubic Bezier
func GetSplinePointBezierCubic(p1 Vector2, c2 Vector2, c3 Vector2, p4 Vector2, t float32) Vector2 {
	ret := getSplinePointBezierCubic(*(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&c2)), *(*uintptr)(unsafe.Pointer(&c3)), *(*uintptr)(unsafe.Pointer(&p4)), t)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// CheckCollisionRecs - Check collision between two rectangles
func CheckCollisionRecs(rec1 Rectangle, rec2 Rectangle) bool {
	return checkCollisionRecs(uintptr(unsafe.Pointer(&rec1)), uintptr(unsafe.Pointer(&rec2)))
}

// CheckCollisionCircles - Check collision between two circles
func CheckCollisionCircles(center1 Vector2, radius1 float32, center2 Vector2, radius2 float32) bool {
	return checkCollisionCircles(*(*uintptr)(unsafe.Pointer(&center1)), radius1, *(*uintptr)(unsafe.Pointer(&center2)), radius2)
}

// CheckCollisionCircleRec - Check collision between circle and rectangle
func CheckCollisionCircleRec(center Vector2, radius float32, rec Rectangle) bool {
	return checkCollisionCircleRec(*(*uintptr)(unsafe.Pointer(&center)), radius, uintptr(unsafe.Pointer(&rec)))
}

// CheckCollisionCircleLine - Check if circle collides with a line created betweeen two points [p1] and [p2]
func CheckCollisionCircleLine(center Vector2, radius float32, p1, p2 Vector2) bool {
	return checkCollisionCircleLine(*(*uintptr)(unsafe.Pointer(&center)), radius, *(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&p2)))
}

// CheckCollisionPointRec - Check if point is inside rectangle
func CheckCollisionPointRec(point Vector2, rec Rectangle) bool {
	return checkCollisionPointRec(*(*uintptr)(unsafe.Pointer(&point)), uintptr(unsafe.Pointer(&rec)))
}

// CheckCollisionPointCircle - Check if point is inside circle
func CheckCollisionPointCircle(point Vector2, center Vector2, radius float32) bool {
	return checkCollisionPointCircle(*(*uintptr)(unsafe.Pointer(&point)), *(*uintptr)(unsafe.Pointer(&center)), radius)
}

// CheckCollisionPointTriangle - Check if point is inside a triangle
func CheckCollisionPointTriangle(point Vector2, p1 Vector2, p2 Vector2, p3 Vector2) bool {
	return checkCollisionPointTriangle(*(*uintptr)(unsafe.Pointer(&point)), *(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&p2)), *(*uintptr)(unsafe.Pointer(&p3)))
}

// CheckCollisionPointPoly - Check if point is within a polygon described by array of vertices
func CheckCollisionPointPoly(point Vector2, points []Vector2) bool {
	pointCount := int32(len(points))
	return checkCollisionPointPoly(*(*uintptr)(unsafe.Pointer(&point)), unsafe.SliceData(points), pointCount)
}

// CheckCollisionLines - Check the collision between two lines defined by two points each, returns collision point by reference
func CheckCollisionLines(startPos1 Vector2, endPos1 Vector2, startPos2 Vector2, endPos2 Vector2, collisionPoint *Vector2) bool {
	return checkCollisionLines(*(*uintptr)(unsafe.Pointer(&startPos1)), *(*uintptr)(unsafe.Pointer(&endPos1)), *(*uintptr)(unsafe.Pointer(&startPos2)), *(*uintptr)(unsafe.Pointer(&endPos2)), collisionPoint)
}

// CheckCollisionPointLine - Check if point belongs to line created between two points [p1] and [p2] with defined margin in pixels [threshold]
func CheckCollisionPointLine(point Vector2, p1 Vector2, p2 Vector2, threshold int32) bool {
	return checkCollisionPointLine(*(*uintptr)(unsafe.Pointer(&point)), *(*uintptr)(unsafe.Pointer(&p1)), *(*uintptr)(unsafe.Pointer(&p2)), threshold)
}

// GetCollisionRec - Get collision rectangle for two rectangles collision
func GetCollisionRec(rec1 Rectangle, rec2 Rectangle) Rectangle {
	var rec Rectangle
	getCollisionRec(uintptr(unsafe.Pointer(&rec)), uintptr(unsafe.Pointer(&rec1)), uintptr(unsafe.Pointer(&rec2)))
	return rec
}

// LoadImage - Load image from file into CPU memory (RAM)
func LoadImage(fileName string) *Image {
	var img Image
	loadImage(uintptr(unsafe.Pointer(&img)), fileName)
	return &img
}

// LoadImageRaw - Load image from RAW file data
func LoadImageRaw(fileName string, width int32, height int32, format PixelFormat, headerSize int32) *Image {
	var img Image
	loadImageRaw(uintptr(unsafe.Pointer(&img)), fileName, width, height, int32(format), headerSize)
	return &img
}

// LoadImageAnim - Load image sequence from file (frames appended to image.data)
func LoadImageAnim(fileName string, frames *int32) *Image {
	var img Image
	loadImageAnim(uintptr(unsafe.Pointer(&img)), fileName, frames)
	return &img
}

// LoadImageAnimFromMemory - Load image sequence from memory buffer
func LoadImageAnimFromMemory(fileType string, fileData []byte, dataSize int32, frames *int32) *Image {
	var img Image
	loadImageAnimFromMemory(uintptr(unsafe.Pointer(&img)), fileType, fileData, dataSize, frames)
	return &img
}

// LoadImageFromMemory - Load image from memory buffer, fileType refers to extension: i.e. '.png'
func LoadImageFromMemory(fileType string, fileData []byte, dataSize int32) *Image {
	var img Image
	loadImageFromMemory(uintptr(unsafe.Pointer(&img)), fileType, fileData, dataSize)
	return &img
}

// LoadImageFromTexture - Load image from GPU texture data
func LoadImageFromTexture(texture Texture2D) *Image {
	var img Image
	loadImageFromTexture(uintptr(unsafe.Pointer(&img)), uintptr(unsafe.Pointer(&texture)))
	return &img
}

// LoadImageFromScreen - Load image from screen buffer and (screenshot)
func LoadImageFromScreen() *Image {
	var img Image
	loadImageFromScreen(uintptr(unsafe.Pointer(&img)))
	return &img
}

// IsImageValid - Check if an image is valid (data and parameters)
func IsImageValid(image *Image) bool {
	return isImageValid(uintptr(unsafe.Pointer(image)))
}

// UnloadImage - Unload image from CPU memory (RAM)
func UnloadImage(image *Image) {
	unloadImage(uintptr(unsafe.Pointer(image)))
}

// ExportImage - Export image data to file, returns true on success
func ExportImage(image Image, fileName string) bool {
	return exportImage(uintptr(unsafe.Pointer(&image)), fileName)
}

// ExportImageToMemory - Export image to memory buffer
func ExportImageToMemory(image Image, fileType string) []byte {
	var fileSize int32
	ret := exportImageToMemory(uintptr(unsafe.Pointer(&image)), fileType, &fileSize)
	return unsafe.Slice(ret, fileSize)
}

// GenImageColor - Generate image: plain color
func GenImageColor(width int, height int, col color.RGBA) *Image {
	var image Image
	genImageColor(uintptr(unsafe.Pointer(&image)), int32(width), int32(height), *(*uintptr)(unsafe.Pointer(&col)))
	return &image
}

// GenImageGradientLinear - Generate image: linear gradient, direction in degrees [0..360], 0=Vertical gradient
func GenImageGradientLinear(width int, height int, direction int, start color.RGBA, end color.RGBA) *Image {
	var image Image
	genImageGradientLinear(uintptr(unsafe.Pointer(&image)), int32(width), int32(height), int32(direction), *(*uintptr)(unsafe.Pointer(&start)), *(*uintptr)(unsafe.Pointer(&end)))
	return &image
}

// GenImageGradientRadial - Generate image: radial gradient
func GenImageGradientRadial(width int, height int, density float32, inner color.RGBA, outer color.RGBA) *Image {
	var image Image
	genImageGradientRadial(uintptr(unsafe.Pointer(&image)), int32(width), int32(height), density, *(*uintptr)(unsafe.Pointer(&inner)), *(*uintptr)(unsafe.Pointer(&outer)))
	return &image
}

// GenImageGradientSquare - Generate image: square gradient
func GenImageGradientSquare(width int, height int, density float32, inner color.RGBA, outer color.RGBA) *Image {
	var image Image
	genImageGradientSquare(uintptr(unsafe.Pointer(&image)), int32(width), int32(height), density, *(*uintptr)(unsafe.Pointer(&inner)), *(*uintptr)(unsafe.Pointer(&outer)))
	return &image
}

// GenImageChecked - Generate image: checked
func GenImageChecked(width int, height int, checksX int, checksY int, col1 color.RGBA, col2 color.RGBA) *Image {
	var image Image
	genImageChecked(uintptr(unsafe.Pointer(&image)), int32(width), int32(height), int32(checksX), int32(checksY), *(*uintptr)(unsafe.Pointer(&col1)), *(*uintptr)(unsafe.Pointer(&col2)))
	return &image
}

// GenImageWhiteNoise - Generate image: white noise
func GenImageWhiteNoise(width int, height int, factor float32) *Image {
	var image Image
	genImageWhiteNoise(uintptr(unsafe.Pointer(&image)), int32(width), int32(height), factor)
	return &image
}

// GenImagePerlinNoise - Generate image: perlin noise
func GenImagePerlinNoise(width int, height int, offsetX int32, offsetY int32, scale float32) *Image {
	var image Image
	genImagePerlinNoise(uintptr(unsafe.Pointer(&image)), int32(width), int32(height), offsetX, offsetY, scale)
	return &image
}

// GenImageCellular - Generate image: cellular algorithm, bigger tileSize means bigger cells
func GenImageCellular(width int, height int, tileSize int) *Image {
	var image Image
	genImageCellular(uintptr(unsafe.Pointer(&image)), int32(width), int32(height), int32(tileSize))
	return &image
}

// GenImageText - Generate image: grayscale image from text data
func GenImageText(width int, height int, text string) Image {
	var image Image
	genImageText(uintptr(unsafe.Pointer(&image)), int32(width), int32(height), text)
	return image
}

// ImageCopy - Create an image duplicate (useful for transformations)
func ImageCopy(image *Image) *Image {
	var retImage Image
	imageCopy(uintptr(unsafe.Pointer(&retImage)), uintptr(unsafe.Pointer(image)))
	return &retImage
}

// ImageFromImage - Create an image from another image piece
func ImageFromImage(image Image, rec Rectangle) Image {
	var retImage Image
	imageFromImage(uintptr(unsafe.Pointer(&retImage)), uintptr(unsafe.Pointer(&image)), uintptr(unsafe.Pointer(&rec)))
	return retImage
}

// ImageFromChannel - Create an image from a selected channel of another image (GRAYSCALE)
func ImageFromChannel(image Image, selectedChannel int32) Image {
	var retImage Image
	imageFromChannel(uintptr(unsafe.Pointer(&retImage)), uintptr(unsafe.Pointer(&image)), selectedChannel)
	return retImage
}

// ImageText - Create an image from text (default font)
func ImageText(text string, fontSize int32, col color.RGBA) Image {
	var retImage Image
	imageText(uintptr(unsafe.Pointer(&retImage)), text, fontSize, *(*uintptr)(unsafe.Pointer(&col)))
	return retImage
}

// ImageTextEx - Create an image from text (custom sprite font)
func ImageTextEx(font Font, text string, fontSize float32, spacing float32, tint color.RGBA) Image {
	var retImage Image
	imageTextEx(uintptr(unsafe.Pointer(&retImage)), uintptr(unsafe.Pointer(&font)), text, fontSize, spacing, *(*uintptr)(unsafe.Pointer(&tint)))
	return retImage
}

// ImageFormat - Convert image data to desired format
func ImageFormat(image *Image, newFormat PixelFormat) {
	imageFormat(image, int32(newFormat))
}

// ImageToPOT - Convert image to POT (power-of-two)
func ImageToPOT(image *Image, fill color.RGBA) {
	imageToPOT(image, *(*uintptr)(unsafe.Pointer(&fill)))
}

// ImageCrop - Crop an image to a defined rectangle
func ImageCrop(image *Image, crop Rectangle) {
	imageCrop(image, uintptr(unsafe.Pointer(&crop)))
}

// ImageAlphaCrop - Crop image depending on alpha value
func ImageAlphaCrop(image *Image, threshold float32) {
	imageAlphaCrop(image, threshold)
}

// ImageAlphaClear - Clear alpha channel to desired color
func ImageAlphaClear(image *Image, col color.RGBA, threshold float32) {
	imageAlphaClear(image, *(*uintptr)(unsafe.Pointer(&col)), threshold)
}

// ImageAlphaMask - Apply alpha mask to image
func ImageAlphaMask(image *Image, alphaMask *Image) {
	imageAlphaMask(image, uintptr(unsafe.Pointer(alphaMask)))
}

// ImageAlphaPremultiply - Premultiply alpha channel
func ImageAlphaPremultiply(image *Image) {
	imageAlphaPremultiply(image)
}

// ImageBlurGaussian - Apply Gaussian blur using a box blur approximation
func ImageBlurGaussian(image *Image, blurSize int32) {
	imageBlurGaussian(image, blurSize)
}

// ImageKernelConvolution - Apply custom square convolution kernel to image
func ImageKernelConvolution(image *Image, kernel []float32) {
	imageKernelConvolution(image, kernel, int32(len(kernel)))
}

// ImageResize - Resize image (Bicubic scaling algorithm)
func ImageResize(image *Image, newWidth int32, newHeight int32) {
	imageResize(image, newWidth, newHeight)
}

// ImageResizeNN - Resize image (Nearest-Neighbor scaling algorithm)
func ImageResizeNN(image *Image, newWidth int32, newHeight int32) {
	imageResizeNN(image, newWidth, newHeight)
}

// ImageResizeCanvas - Resize canvas and fill with color
func ImageResizeCanvas(image *Image, newWidth int32, newHeight int32, offsetX int32, offsetY int32, fill color.RGBA) {
	imageResizeCanvas(image, newWidth, newHeight, offsetX, offsetY, *(*uintptr)(unsafe.Pointer(&fill)))
}

// ImageMipmaps - Compute all mipmap levels for a provided image
func ImageMipmaps(image *Image) {
	imageMipmaps(image)
}

// ImageDither - Dither image data to 16bpp or lower (Floyd-Steinberg dithering)
func ImageDither(image *Image, rBpp int32, gBpp int32, bBpp int32, aBpp int32) {
	imageDither(image, rBpp, gBpp, bBpp, aBpp)
}

// ImageFlipVertical - Flip image vertically
func ImageFlipVertical(image *Image) {
	imageFlipVertical(image)
}

// ImageFlipHorizontal - Flip image horizontally
func ImageFlipHorizontal(image *Image) {
	imageFlipHorizontal(image)
}

// ImageRotate - Rotate image by input angle in degrees (-359 to 359)
func ImageRotate(image *Image, degrees int32) {
	imageRotate(image, degrees)
}

// ImageRotateCW - Rotate image clockwise 90deg
func ImageRotateCW(image *Image) {
	imageRotateCW(image)
}

// ImageRotateCCW - Rotate image counter-clockwise 90deg
func ImageRotateCCW(image *Image) {
	imageRotateCCW(image)
}

// ImageColorTint - Modify image color: tint
func ImageColorTint(image *Image, col color.RGBA) {
	imageColorTint(image, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageColorInvert - Modify image color: invert
func ImageColorInvert(image *Image) {
	imageColorInvert(image)
}

// ImageColorGrayscale - Modify image color: grayscale
func ImageColorGrayscale(image *Image) {
	imageColorGrayscale(image)
}

// ImageColorContrast - Modify image color: contrast (-100 to 100)
func ImageColorContrast(image *Image, contrast float32) {
	imageColorContrast(image, contrast)
}

// ImageColorBrightness - Modify image color: brightness (-255 to 255)
func ImageColorBrightness(image *Image, brightness int32) {
	imageColorBrightness(image, brightness)
}

// ImageColorReplace - Modify image color: replace color
func ImageColorReplace(image *Image, col color.RGBA, replace color.RGBA) {
	imageColorReplace(image, *(*uintptr)(unsafe.Pointer(&col)), *(*uintptr)(unsafe.Pointer(&replace)))
}

// LoadImageColors - Load color data from image as a Color array (RGBA - 32bit)
//
// NOTE: Memory allocated should be freed using UnloadImageColors()
func LoadImageColors(image *Image) []color.RGBA {
	ret := loadImageColors(uintptr(unsafe.Pointer(image)))
	return unsafe.Slice(ret, image.Width*image.Height)
}

// LoadImagePalette - Load colors palette from image as a Color array (RGBA - 32bit)
//
// NOTE: Memory allocated should be freed using UnloadImagePalette()
func LoadImagePalette(image Image, maxPaletteSize int32) []color.RGBA {
	var colorCount int32
	ret := loadImagePalette(uintptr(unsafe.Pointer(&image)), maxPaletteSize, &colorCount)
	return unsafe.Slice(ret, colorCount)
}

// UnloadImageColors - Unload color data loaded with LoadImageColors()
func UnloadImageColors(colors []color.RGBA) {
	unloadImageColors(unsafe.SliceData(colors))
}

// UnloadImagePalette - Unload colors palette loaded with LoadImagePalette()
func UnloadImagePalette(colors []color.RGBA) {
	unloadImagePalette(unsafe.SliceData(colors))
}

// GetImageAlphaBorder - Get image alpha border rectangle
func GetImageAlphaBorder(image Image, threshold float32) Rectangle {
	var rec Rectangle
	getImageAlphaBorder(uintptr(unsafe.Pointer(&rec)), uintptr(unsafe.Pointer(&image)), threshold)
	return rec
}

// GetImageColor - Get image pixel color at (x, y) position
func GetImageColor(image Image, x int32, y int32) color.RGBA {
	ret := getImageColor(uintptr(unsafe.Pointer(&image)), x, y)
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// ImageClearBackground - Clear image background with given color
func ImageClearBackground(dst *Image, col color.RGBA) {
	imageClearBackground(dst, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawPixel - Draw pixel within an image
func ImageDrawPixel(dst *Image, posX int32, posY int32, col color.RGBA) {
	imageDrawPixel(dst, posX, posY, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawPixelV - Draw pixel within an image (Vector version)
func ImageDrawPixelV(dst *Image, position Vector2, col color.RGBA) {
	imageDrawPixelV(dst, *(*uintptr)(unsafe.Pointer(&position)), *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawLine - Draw line within an image
func ImageDrawLine(dst *Image, startPosX int32, startPosY int32, endPosX int32, endPosY int32, col color.RGBA) {
	imageDrawLine(dst, startPosX, startPosY, endPosX, endPosY, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawLineV - Draw line within an image (Vector version)
func ImageDrawLineV(dst *Image, start, end Vector2, col color.RGBA) {
	imageDrawLineV(dst, *(*uintptr)(unsafe.Pointer(&start)), *(*uintptr)(unsafe.Pointer(&end)), *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawLineEx - Draw a line defining thickness within an image
func ImageDrawLineEx(dst *Image, start, end Vector2, thick int32, col color.RGBA) {
	imageDrawLineEx(dst, *(*uintptr)(unsafe.Pointer(&start)), *(*uintptr)(unsafe.Pointer(&end)), thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawCircle - Draw a filled circle within an image
func ImageDrawCircle(dst *Image, centerX int32, centerY int32, radius int32, col color.RGBA) {
	imageDrawCircle(dst, centerX, centerY, radius, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawCircleV - Draw a filled circle within an image (Vector version)
func ImageDrawCircleV(dst *Image, center Vector2, radius int32, col color.RGBA) {
	imageDrawCircleV(dst, *(*uintptr)(unsafe.Pointer(&center)), radius, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawCircleLines - Draw circle outline within an image
func ImageDrawCircleLines(dst *Image, centerX int32, centerY int32, radius int32, col color.RGBA) {
	imageDrawCircleLines(dst, centerX, centerY, radius, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawCircleLinesV - Draw circle outline within an image (Vector version)
func ImageDrawCircleLinesV(dst *Image, center Vector2, radius int32, col color.RGBA) {
	imageDrawCircleLinesV(dst, *(*uintptr)(unsafe.Pointer(&center)), radius, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawRectangle - Draw rectangle within an image
func ImageDrawRectangle(dst *Image, posX int32, posY int32, width int32, height int32, col color.RGBA) {
	imageDrawRectangle(dst, posX, posY, width, height, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawRectangleV - Draw rectangle within an image (Vector version)
func ImageDrawRectangleV(dst *Image, position Vector2, size Vector2, col color.RGBA) {
	imageDrawRectangleV(dst, *(*uintptr)(unsafe.Pointer(&position)), *(*uintptr)(unsafe.Pointer(&size)), *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawRectangleRec - Draw rectangle within an image
func ImageDrawRectangleRec(dst *Image, rec Rectangle, col color.RGBA) {
	imageDrawRectangleRec(dst, uintptr(unsafe.Pointer(&rec)), *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawRectangleLines - Draw rectangle lines within an image
func ImageDrawRectangleLines(dst *Image, rec Rectangle, thick int, col color.RGBA) {
	imageDrawRectangleLines(dst, uintptr(unsafe.Pointer(&rec)), int32(thick), *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawTriangle - Draw triangle within an image
func ImageDrawTriangle(dst *Image, v1, v2, v3 Vector2, col color.RGBA) {
	imageDrawTriangle(dst, *(*uintptr)(unsafe.Pointer(&v1)), *(*uintptr)(unsafe.Pointer(&v2)), *(*uintptr)(unsafe.Pointer(&v3)), *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawTriangleEx - Draw triangle with interpolated colors within an image
func ImageDrawTriangleEx(dst *Image, v1, v2, v3 Vector2, c1, c2, c3 color.RGBA) {
	imageDrawTriangleEx(dst, *(*uintptr)(unsafe.Pointer(&v1)), *(*uintptr)(unsafe.Pointer(&v2)), *(*uintptr)(unsafe.Pointer(&v3)), *(*uintptr)(unsafe.Pointer(&c1)), *(*uintptr)(unsafe.Pointer(&c2)), *(*uintptr)(unsafe.Pointer(&c3)))
}

// ImageDrawTriangleLines - Draw triangle outline within an image
func ImageDrawTriangleLines(dst *Image, v1, v2, v3 Vector2, col color.RGBA) {
	imageDrawTriangleLines(dst, *(*uintptr)(unsafe.Pointer(&v1)), *(*uintptr)(unsafe.Pointer(&v2)), *(*uintptr)(unsafe.Pointer(&v3)), *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawTriangleFan - Draw a triangle fan defined by points within an image (first vertex is the center)
func ImageDrawTriangleFan(dst *Image, points []Vector2, col color.RGBA) {
	pointCount := int32(len(points))
	imageDrawTriangleFan(dst, (unsafe.SliceData(points)), pointCount, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDrawTriangleStrip - Draw a triangle strip defined by points within an image
func ImageDrawTriangleStrip(dst *Image, points []Vector2, col color.RGBA) {
	pointCount := int32(len(points))
	imageDrawTriangleStrip(dst, (unsafe.SliceData(points)), pointCount, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDraw - Draw a source image within a destination image (tint applied to source)
func ImageDraw(dst *Image, src *Image, srcRec Rectangle, dstRec Rectangle, tint color.RGBA) {
	imageDraw(dst, uintptr(unsafe.Pointer(src)), uintptr(unsafe.Pointer(&srcRec)), uintptr(unsafe.Pointer(&dstRec)), *(*uintptr)(unsafe.Pointer(&tint)))
}

// ImageDrawText - Draw text (using default font) within an image (destination)
func ImageDrawText(dst *Image, posX int32, posY int32, text string, fontSize int32, col color.RGBA) {
	imageDrawText(dst, text, posX, posY, fontSize, *(*uintptr)(unsafe.Pointer(&col)))

}

// ImageDrawTextEx - Draw text (custom sprite font) within an image (destination)
func ImageDrawTextEx(dst *Image, position Vector2, font Font, text string, fontSize float32, spacing float32, tint color.RGBA) {
	imageDrawTextEx(dst, uintptr(unsafe.Pointer(&font)), text, *(*uintptr)(unsafe.Pointer(&position)), fontSize, spacing, *(*uintptr)(unsafe.Pointer(&tint)))
}

// LoadTexture - Load texture from file into GPU memory (VRAM)
func LoadTexture(fileName string) Texture2D {
	var texture Texture2D
	loadTexture(uintptr(unsafe.Pointer(&texture)), fileName)
	return texture
}

// LoadTextureFromImage - Load texture from image data
func LoadTextureFromImage(image *Image) Texture2D {
	var texture Texture2D
	loadTextureFromImage(uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(image)))
	return texture
}

// LoadTextureCubemap - Load cubemap from image, multiple image cubemap layouts supported
func LoadTextureCubemap(image *Image, layout int32) Texture2D {
	var texture Texture2D
	loadTextureCubemap(uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(image)), layout)
	return texture
}

// LoadRenderTexture - Load texture for rendering (framebuffer)
func LoadRenderTexture(width int32, height int32) RenderTexture2D {
	var texture RenderTexture2D
	loadRenderTexture(uintptr(unsafe.Pointer(&texture)), width, height)
	return texture
}

// IsTextureValid - Check if a texture is valid (loaded in GPU)
func IsTextureValid(texture Texture2D) bool {
	return isTextureValid(uintptr(unsafe.Pointer(&texture)))
}

// UnloadTexture - Unload texture from GPU memory (VRAM)
func UnloadTexture(texture Texture2D) {
	unloadTexture(uintptr(unsafe.Pointer(&texture)))
}

// IsRenderTextureValid - Check if a render texture is valid (loaded in GPU)
func IsRenderTextureValid(target RenderTexture2D) bool {
	return isRenderTextureValid(uintptr(unsafe.Pointer(&target)))
}

// UnloadRenderTexture - Unload render texture from GPU memory (VRAM)
func UnloadRenderTexture(target RenderTexture2D) {
	unloadRenderTexture(uintptr(unsafe.Pointer(&target)))
}

// UpdateTexture - Update GPU texture with new data
func UpdateTexture(texture Texture2D, pixels []color.RGBA) {
	updateTexture(uintptr(unsafe.Pointer(&texture)), unsafe.SliceData(pixels))
}

// UpdateTextureRec - Update GPU texture rectangle with new data
func UpdateTextureRec(texture Texture2D, rec Rectangle, pixels []color.RGBA) {
	updateTextureRec(uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(&rec)), unsafe.SliceData(pixels))
}

// GenTextureMipmaps - Generate GPU mipmaps for a texture
func GenTextureMipmaps(texture *Texture2D) {
	genTextureMipmaps(texture)
}

// SetTextureFilter - Set texture scaling filter mode
func SetTextureFilter(texture Texture2D, filter TextureFilterMode) {
	setTextureFilter(uintptr(unsafe.Pointer(&texture)), int32(filter))
}

// SetTextureWrap - Set texture wrapping mode
func SetTextureWrap(texture Texture2D, wrap TextureWrapMode) {
	setTextureWrap(uintptr(unsafe.Pointer(&texture)), int32(wrap))
}

// DrawTexture - Draw a Texture2D
func DrawTexture(texture Texture2D, posX int32, posY int32, tint color.RGBA) {
	drawTexture(uintptr(unsafe.Pointer(&texture)), posX, posY, *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawTextureV - Draw a Texture2D with position defined as Vector2
func DrawTextureV(texture Texture2D, position Vector2, tint color.RGBA) {
	drawTextureV(uintptr(unsafe.Pointer(&texture)), *(*uintptr)(unsafe.Pointer(&position)), *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawTextureEx - Draw a Texture2D with extended parameters
func DrawTextureEx(texture Texture2D, position Vector2, rotation float32, scale float32, tint color.RGBA) {
	drawTextureEx(uintptr(unsafe.Pointer(&texture)), *(*uintptr)(unsafe.Pointer(&position)), rotation, scale, *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawTextureRec - Draw a part of a texture defined by a rectangle
func DrawTextureRec(texture Texture2D, source Rectangle, position Vector2, tint color.RGBA) {
	drawTextureRec(uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(&source)), *(*uintptr)(unsafe.Pointer(&position)), *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawTexturePro - Draw a part of a texture defined by a rectangle with 'pro' parameters
func DrawTexturePro(texture Texture2D, source Rectangle, dest Rectangle, origin Vector2, rotation float32, tint color.RGBA) {
	drawTexturePro(uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(&source)), uintptr(unsafe.Pointer(&dest)), *(*uintptr)(unsafe.Pointer(&origin)), rotation, *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawTextureNPatch - Draws a texture (or part of it) that stretches or shrinks nicely
func DrawTextureNPatch(texture Texture2D, nPatchInfo NPatchInfo, dest Rectangle, origin Vector2, rotation float32, tint color.RGBA) {
	drawTextureNPatch(uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(&nPatchInfo)), uintptr(unsafe.Pointer(&dest)), *(*uintptr)(unsafe.Pointer(&origin)), rotation, *(*uintptr)(unsafe.Pointer(&tint)))
}

// Fade - Get color with alpha applied, alpha goes from 0.0f to 1.0f
func Fade(col color.RGBA, alpha float32) color.RGBA {
	ret := fade(*(*uintptr)(unsafe.Pointer(&col)), alpha)
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// ColorToInt - Get hexadecimal value for a Color (0xRRGGBBAA)
func ColorToInt(col color.RGBA) int32 {
	return colorToInt(*(*uintptr)(unsafe.Pointer(&col)))
}

// ColorNormalize - Get Color normalized as float [0..1]
func ColorNormalize(col color.RGBA) Vector4 {
	var vector4 Vector4
	colorNormalize(uintptr(unsafe.Pointer(&vector4)), *(*uintptr)(unsafe.Pointer(&col)))
	return vector4
}

// ColorFromNormalized - Get Color from normalized values [0..1]
func ColorFromNormalized(normalized Vector4) color.RGBA {
	ret := colorFromNormalized(uintptr(unsafe.Pointer(&normalized)))
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// ColorToHSV - Get HSV values for a Color, hue [0..360], saturation/value [0..1]
func ColorToHSV(col color.RGBA) Vector3 {
	var vector3 Vector3
	colorToHSV(uintptr(unsafe.Pointer(&vector3)), *(*uintptr)(unsafe.Pointer(&col)))
	return vector3
}

// ColorFromHSV - Get a Color from HSV values, hue [0..360], saturation/value [0..1]
func ColorFromHSV(hue float32, saturation float32, value float32) color.RGBA {
	ret := colorFromHSV(hue, saturation, value)
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// ColorTint - Get color multiplied with another color
func ColorTint(col color.RGBA, tint color.RGBA) color.RGBA {
	ret := colorTint(*(*uintptr)(unsafe.Pointer(&col)), *(*uintptr)(unsafe.Pointer(&tint)))
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// ColorBrightness - Get color with brightness correction, brightness factor goes from -1.0f to 1.0f
func ColorBrightness(col color.RGBA, factor float32) color.RGBA {
	ret := colorBrightness(*(*uintptr)(unsafe.Pointer(&col)), factor)
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// ColorContrast - Get color with contrast correction, contrast values between -1.0f and 1.0f
func ColorContrast(col color.RGBA, contrast float32) color.RGBA {
	ret := colorContrast(*(*uintptr)(unsafe.Pointer(&col)), contrast)
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// ColorAlpha - Get color with alpha applied, alpha goes from 0.0f to 1.0f
func ColorAlpha(col color.RGBA, alpha float32) color.RGBA {
	ret := colorAlpha(*(*uintptr)(unsafe.Pointer(&col)), alpha)
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// ColorAlphaBlend - Get src alpha-blended into dst color with tint
func ColorAlphaBlend(dst color.RGBA, src color.RGBA, tint color.RGBA) color.RGBA {
	ret := colorAlphaBlend(*(*uintptr)(unsafe.Pointer(&dst)), *(*uintptr)(unsafe.Pointer(&src)), *(*uintptr)(unsafe.Pointer(&tint)))
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// ColorLerp - Get color lerp interpolation between two colors, factor [0.0f..1.0f]
func ColorLerp(col1, col2 color.RGBA, factor float32) color.RGBA {
	ret := colorLerp(*(*uintptr)(unsafe.Pointer(&col1)), *(*uintptr)(unsafe.Pointer(&col2)), factor)
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// GetColor - Get Color structure from hexadecimal value
func GetColor(hexValue uint) color.RGBA {
	ret := getColor(uint32(hexValue))
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// GetPixelColor - Get Color from a source pixel pointer of certain format
func GetPixelColor(srcPtr unsafe.Pointer, format int32) color.RGBA {
	ret := getPixelColor(srcPtr, format)
	return *(*color.RGBA)(unsafe.Pointer(&ret))
}

// SetPixelColor - Set color formatted into destination pixel pointer
func SetPixelColor(dstPtr unsafe.Pointer, col color.RGBA, format int32) {
	setPixelColor(dstPtr, *(*uintptr)(unsafe.Pointer(&col)), format)
}

// GetPixelDataSize - Get pixel data size in bytes for certain format
func GetPixelDataSize(width int32, height int32, format int32) int32 {
	return getPixelDataSize(width, height, format)
}

// GetFontDefault - Get the default Font
func GetFontDefault() Font {
	var font Font
	getFontDefault(uintptr(unsafe.Pointer(&font)))
	return font
}

// LoadFont - Load font from file into GPU memory (VRAM)
func LoadFont(fileName string) Font {
	var font Font
	loadFont(uintptr(unsafe.Pointer(&font)), fileName)
	return font
}

// LoadFontEx - Load font from file with extended parameters, use NULL for codepoints and 0 for codepointCount to load the default character setFont
func LoadFontEx(fileName string, fontSize int32, codepoints []rune, runesNumber ...int32) Font {
	var font Font
	codepointCount := int32(len(codepoints))
	if len(runesNumber) > 0 {
		codepointCount = int32(runesNumber[0])
	}
	loadFontEx(uintptr(unsafe.Pointer(&font)), fileName, fontSize, codepoints, codepointCount)
	return font
}

// LoadFontFromImage - Load font from Image (XNA style)
func LoadFontFromImage(image Image, key color.RGBA, firstChar rune) Font {
	var font Font
	loadFontFromImage(uintptr(unsafe.Pointer(&font)), uintptr(unsafe.Pointer(&image)), *(*uintptr)(unsafe.Pointer(&key)), firstChar)
	return font
}

// LoadFontFromMemory - Load font from memory buffer, fileType refers to extension: i.e. '.ttf'
func LoadFontFromMemory(fileType string, fileData []byte, fontSize int32, codepoints []rune) Font {
	var font Font
	dataSize := int32(len(fileData))
	codepointCount := int32(len(codepoints))
	loadFontFromMemory(uintptr(unsafe.Pointer(&font)), fileType, fileData, dataSize, fontSize, codepoints, codepointCount)
	return font
}

// IsFontValid - Check if a font is valid (font data loaded, WARNING: GPU texture not checked)
func IsFontValid(font Font) bool {
	return isFontValid(uintptr(unsafe.Pointer(&font)))
}

// LoadFontData - Load font data for further use
func LoadFontData(fileData []byte, fontSize int32, codepoints []rune, codepointCount, typ int32) []GlyphInfo {
	dataSize := int32(len(fileData))
	// In case no chars count provided, default to 95
	if codepointCount <= 0 {
		codepointCount = 95
	}
	ret := loadFontData(fileData, dataSize, fontSize, codepoints, codepointCount, typ)
	return unsafe.Slice(ret, codepointCount)
}

// GenImageFontAtlas - Generate image font atlas using chars info
func GenImageFontAtlas(glyphs []GlyphInfo, glyphRecs []*Rectangle, fontSize int32, padding int32, packMethod int32) Image {
	var image Image
	glyphCount := int32(len(glyphs))
	genImageFontAtlas(uintptr(unsafe.Pointer(&image)), unsafe.SliceData(glyphs), glyphRecs, glyphCount, fontSize, padding, packMethod)
	return image
}

// UnloadFontData - Unload font chars info data (RAM)
func UnloadFontData(glyphs []GlyphInfo) {
	glyphCount := int32(len(glyphs))
	unloadFontData(unsafe.SliceData(glyphs), glyphCount)
}

// UnloadFont - Unload font from GPU memory (VRAM)
func UnloadFont(font Font) {
	unloadFont(uintptr(unsafe.Pointer(&font)))
}

// DrawFPS - Draw current FPS
func DrawFPS(posX int32, posY int32) {
	drawFPS(posX, posY)
}

// DrawText - Draw text (using default font)
func DrawText(text string, posX int32, posY int32, fontSize int32, col color.RGBA) {
	drawText(text, posX, posY, fontSize, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawTextEx - Draw text using font and additional parameters
func DrawTextEx(font Font, text string, position Vector2, fontSize float32, spacing float32, tint color.RGBA) {
	drawTextEx(uintptr(unsafe.Pointer(&font)), text, *(*uintptr)(unsafe.Pointer(&position)), fontSize, spacing, *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawTextPro - Draw text using Font and pro parameters (rotation)
func DrawTextPro(font Font, text string, position Vector2, origin Vector2, rotation float32, fontSize float32, spacing float32, tint color.RGBA) {
	drawTextPro(uintptr(unsafe.Pointer(&font)), text, *(*uintptr)(unsafe.Pointer(&position)), *(*uintptr)(unsafe.Pointer(&origin)), rotation, fontSize, spacing, *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawTextCodepoint - Draw one character (codepoint)
func DrawTextCodepoint(font Font, codepoint rune, position Vector2, fontSize float32, tint color.RGBA) {
	drawTextCodepoint(uintptr(unsafe.Pointer(&font)), codepoint, *(*uintptr)(unsafe.Pointer(&position)), fontSize, *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawTextCodepoints - Draw multiple character (codepoint)
func DrawTextCodepoints(font Font, codepoints []rune, position Vector2, fontSize float32, spacing float32, tint color.RGBA) {
	codepointCount := int32(len(codepoints))
	drawTextCodepoints(uintptr(unsafe.Pointer(&font)), codepoints, codepointCount, *(*uintptr)(unsafe.Pointer(&position)), fontSize, spacing, *(*uintptr)(unsafe.Pointer(&tint)))
}

// SetTextLineSpacing - Set vertical line spacing when drawing with line-breaks
func SetTextLineSpacing(spacing int) {
	setTextLineSpacing(int32(spacing))
}

// MeasureText - Measure string width for default font
func MeasureText(text string, fontSize int32) int32 {
	return measureText(text, fontSize)
}

// MeasureTextEx - Measure string size for Font
func MeasureTextEx(font Font, text string, fontSize float32, spacing float32) Vector2 {
	ret := measureTextEx(uintptr(unsafe.Pointer(&font)), text, fontSize, spacing)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetGlyphIndex - Get glyph index position in font for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphIndex(font Font, codepoint rune) int32 {
	return getGlyphIndex(uintptr(unsafe.Pointer(&font)), codepoint)
}

// GetGlyphInfo - Get glyph font info data for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphInfo(font Font, codepoint rune) GlyphInfo {
	var glyphInfo GlyphInfo
	getGlyphInfo(uintptr(unsafe.Pointer(&glyphInfo)), uintptr(unsafe.Pointer(&font)), codepoint)
	return glyphInfo
}

// GetGlyphAtlasRec - Get glyph rectangle in font atlas for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphAtlasRec(font Font, codepoint rune) Rectangle {
	var rec Rectangle
	getGlyphAtlasRec(uintptr(unsafe.Pointer(&rec)), uintptr(unsafe.Pointer(&font)), codepoint)
	return rec
}

// DrawLine3D - Draw a line in 3D world space
func DrawLine3D(startPos Vector3, endPos Vector3, col color.RGBA) {
	drawLine3D(uintptr(unsafe.Pointer(&startPos)), uintptr(unsafe.Pointer(&endPos)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawPoint3D - Draw a point in 3D space, actually a small line
func DrawPoint3D(position Vector3, col color.RGBA) {
	drawPoint3D(uintptr(unsafe.Pointer(&position)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCircle3D - Draw a circle in 3D world space
func DrawCircle3D(center Vector3, radius float32, rotationAxis Vector3, rotationAngle float32, col color.RGBA) {
	drawCircle3D(uintptr(unsafe.Pointer(&center)), radius, uintptr(unsafe.Pointer(&rotationAxis)), rotationAngle, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawTriangle3D - Draw a color-filled triangle (vertex in counter-clockwise order!)
func DrawTriangle3D(v1 Vector3, v2 Vector3, v3 Vector3, col color.RGBA) {
	drawTriangle3D(uintptr(unsafe.Pointer(&v1)), uintptr(unsafe.Pointer(&v2)), uintptr(unsafe.Pointer(&v3)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawTriangleStrip3D - Draw a triangle strip defined by points
func DrawTriangleStrip3D(points []Vector3, col color.RGBA) {
	pointCount := int32(len(points))
	drawTriangleStrip3D(unsafe.SliceData(points), pointCount, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCube - Draw cube
func DrawCube(position Vector3, width float32, height float32, length float32, col color.RGBA) {
	drawCube(uintptr(unsafe.Pointer(&position)), width, height, length, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCubeV - Draw cube (Vector version)
func DrawCubeV(position Vector3, size Vector3, col color.RGBA) {
	drawCubeV(uintptr(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&size)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCubeWires - Draw cube wires
func DrawCubeWires(position Vector3, width float32, height float32, length float32, col color.RGBA) {
	drawCubeWires(uintptr(unsafe.Pointer(&position)), width, height, length, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCubeWiresV - Draw cube wires (Vector version)
func DrawCubeWiresV(position Vector3, size Vector3, col color.RGBA) {
	drawCubeWiresV(uintptr(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&size)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSphere - Draw sphere
func DrawSphere(centerPos Vector3, radius float32, col color.RGBA) {
	drawSphere(uintptr(unsafe.Pointer(&centerPos)), radius, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSphereEx - Draw sphere with extended parameters
func DrawSphereEx(centerPos Vector3, radius float32, rings int32, slices int32, col color.RGBA) {
	drawSphereEx(uintptr(unsafe.Pointer(&centerPos)), radius, rings, slices, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawSphereWires - Draw sphere wires
func DrawSphereWires(centerPos Vector3, radius float32, rings int32, slices int32, col color.RGBA) {
	drawSphereWires(uintptr(unsafe.Pointer(&centerPos)), radius, rings, slices, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCylinder - Draw a cylinder/cone
func DrawCylinder(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, col color.RGBA) {
	drawCylinder(uintptr(unsafe.Pointer(&position)), radiusTop, radiusBottom, height, slices, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCylinderEx - Draw a cylinder with base at startPos and top at endPos
func DrawCylinderEx(startPos Vector3, endPos Vector3, startRadius float32, endRadius float32, sides int32, col color.RGBA) {
	drawCylinderEx(uintptr(unsafe.Pointer(&startPos)), uintptr(unsafe.Pointer(&endPos)), startRadius, endRadius, sides, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCylinderWires - Draw a cylinder/cone wires
func DrawCylinderWires(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, col color.RGBA) {
	drawCylinderWires(uintptr(unsafe.Pointer(&position)), radiusTop, radiusBottom, height, slices, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCylinderWiresEx - Draw a cylinder wires with base at startPos and top at endPos
func DrawCylinderWiresEx(startPos Vector3, endPos Vector3, startRadius float32, endRadius float32, sides int32, col color.RGBA) {
	drawCylinderWiresEx(uintptr(unsafe.Pointer(&startPos)), uintptr(unsafe.Pointer(&endPos)), startRadius, endRadius, sides, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCapsule - Draw a capsule with the center of its sphere caps at startPos and endPos
func DrawCapsule(startPos Vector3, endPos Vector3, radius float32, slices int32, rings int32, col color.RGBA) {
	drawCapsule(uintptr(unsafe.Pointer(&startPos)), uintptr(unsafe.Pointer(&endPos)), radius, slices, rings, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawCapsuleWires - Draw capsule wireframe with the center of its sphere caps at startPos and endPos
func DrawCapsuleWires(startPos Vector3, endPos Vector3, radius float32, slices int32, rings int32, col color.RGBA) {
	drawCapsuleWires(uintptr(unsafe.Pointer(&startPos)), uintptr(unsafe.Pointer(&endPos)), radius, slices, rings, *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawPlane - Draw a plane XZ
func DrawPlane(centerPos Vector3, size Vector2, col color.RGBA) {
	drawPlane(uintptr(unsafe.Pointer(&centerPos)), *(*uintptr)(unsafe.Pointer(&size)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawRay - Draw a ray line
func DrawRay(ray Ray, col color.RGBA) {
	drawRay(uintptr(unsafe.Pointer(&ray)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawGrid - Draw a grid (centered at (0, 0, 0))
func DrawGrid(slices int32, spacing float32) {
	drawGrid(slices, spacing)
}

// LoadModel - Load model from files (meshes and materials)
func LoadModel(fileName string) Model {
	var model Model
	loadModel(uintptr(unsafe.Pointer(&model)), fileName)
	return model
}

// LoadModelFromMesh - Load model from generated mesh (default material)
func LoadModelFromMesh(mesh Mesh) Model {
	var model Model
	loadModelFromMesh(uintptr(unsafe.Pointer(&model)), uintptr(unsafe.Pointer(&mesh)))
	return model
}

// IsModelValid - Check if a model is valid (loaded in GPU, VAO/VBOs)
func IsModelValid(model Model) bool {
	return isModelValid(uintptr(unsafe.Pointer(&model)))
}

// UnloadModel - Unload model (including meshes) from memory (RAM and/or VRAM)
func UnloadModel(model Model) {
	unloadModel(uintptr(unsafe.Pointer(&model)))
}

// GetModelBoundingBox - Compute model bounding box limits (considers all meshes)
func GetModelBoundingBox(model Model) BoundingBox {
	var boundingBox BoundingBox
	getModelBoundingBox(uintptr(unsafe.Pointer(&boundingBox)), uintptr(unsafe.Pointer(&model)))
	return boundingBox
}

// DrawModel - Draw a model (with texture if set)
func DrawModel(model Model, position Vector3, scale float32, tint color.RGBA) {
	drawModel(uintptr(unsafe.Pointer(&model)), uintptr(unsafe.Pointer(&position)), scale, *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawModelEx - Draw a model with extended parameters
func DrawModelEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
	drawModelEx(uintptr(unsafe.Pointer(&model)), uintptr(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&rotationAxis)), rotationAngle, uintptr(unsafe.Pointer(&scale)), *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawModelWires - Draw a model wires (with texture if set)
func DrawModelWires(model Model, position Vector3, scale float32, tint color.RGBA) {
	drawModelWires(uintptr(unsafe.Pointer(&model)), uintptr(unsafe.Pointer(&position)), scale, *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawModelWiresEx - Draw a model wires (with texture if set) with extended parameters
func DrawModelWiresEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
	drawModelWiresEx(uintptr(unsafe.Pointer(&model)), uintptr(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&rotationAxis)), rotationAngle, uintptr(unsafe.Pointer(&scale)), *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawModelPoints - Draw a model as points
func DrawModelPoints(model Model, position Vector3, scale float32, tint color.RGBA) {
	drawModelPoints(uintptr(unsafe.Pointer(&model)), uintptr(unsafe.Pointer(&position)), scale, *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawModelPointsEx - Draw a model as points with extended parameters
func DrawModelPointsEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
	drawModelPointsEx(uintptr(unsafe.Pointer(&model)), uintptr(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&rotationAxis)), rotationAngle, uintptr(unsafe.Pointer(&scale)), *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawBoundingBox - Draw bounding box (wires)
func DrawBoundingBox(box BoundingBox, col color.RGBA) {
	drawBoundingBox(uintptr(unsafe.Pointer(&box)), *(*uintptr)(unsafe.Pointer(&col)))
}

// DrawBillboard - Draw a billboard texture
func DrawBillboard(camera Camera, texture Texture2D, position Vector3, scale float32, tint color.RGBA) {
	drawBillboard(uintptr(unsafe.Pointer(&camera)), uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(&position)), scale, *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawBillboardRec - Draw a billboard texture defined by source
func DrawBillboardRec(camera Camera, texture Texture2D, source Rectangle, position Vector3, size Vector2, tint color.RGBA) {
	drawBillboardRec(uintptr(unsafe.Pointer(&camera)), uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(&source)), uintptr(unsafe.Pointer(&position)), *(*uintptr)(unsafe.Pointer(&size)), *(*uintptr)(unsafe.Pointer(&tint)))
}

// DrawBillboardPro - Draw a billboard texture defined by source and rotation
func DrawBillboardPro(camera Camera, texture Texture2D, source Rectangle, position Vector3, up Vector3, size Vector2, origin Vector2, rotation float32, tint color.RGBA) {
	drawBillboardPro(uintptr(unsafe.Pointer(&camera)), uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(&source)), uintptr(unsafe.Pointer(&position)), uintptr(unsafe.Pointer(&up)), *(*uintptr)(unsafe.Pointer(&size)), *(*uintptr)(unsafe.Pointer(&origin)), rotation, *(*uintptr)(unsafe.Pointer(&tint)))
}

// UploadMesh - Upload mesh vertex data in GPU and provide VAO/VBO ids
func UploadMesh(mesh *Mesh, dynamic bool) {
	uploadMesh(mesh, dynamic)
}

// UpdateMeshBuffer - Update mesh vertex data in GPU for a specific buffer index
func UpdateMeshBuffer(mesh Mesh, index int32, data []byte, offset int) {
	dataSize := int32(len(data))
	updateMeshBuffer(uintptr(unsafe.Pointer(&mesh)), index, data, dataSize, int32(offset))
}

// UnloadMesh - Unload mesh data from CPU and GPU
func UnloadMesh(mesh *Mesh) {
	unloadMesh(uintptr(unsafe.Pointer(mesh)))
}

// DrawMesh - Draw a 3d mesh with material and transform
func DrawMesh(mesh Mesh, material Material, transform Matrix) {
	drawMesh(uintptr(unsafe.Pointer(&mesh)), uintptr(unsafe.Pointer(&material)), uintptr(unsafe.Pointer(&transform)))
}

// DrawMeshInstanced - Draw multiple mesh instances with material and different transforms
func DrawMeshInstanced(mesh Mesh, material Material, transforms []Matrix, instances int32) {
	drawMeshInstanced(uintptr(unsafe.Pointer(&mesh)), uintptr(unsafe.Pointer(&material)), transforms, instances)
}

// ExportMesh - Export mesh data to file, returns true on success
func ExportMesh(mesh Mesh, fileName string) bool {
	return exportMesh(uintptr(unsafe.Pointer(&mesh)), fileName)
}

// GetMeshBoundingBox - Compute mesh bounding box limits
func GetMeshBoundingBox(mesh Mesh) BoundingBox {
	var boundingBox BoundingBox
	getMeshBoundingBox(uintptr(unsafe.Pointer(&boundingBox)), uintptr(unsafe.Pointer(&mesh)))
	return boundingBox
}

// GenMeshTangents - Compute mesh tangents
func GenMeshTangents(mesh *Mesh) {
	genMeshTangents(mesh)
}

// GenMeshPoly - Generate polygonal mesh
func GenMeshPoly(sides int, radius float32) Mesh {
	var mesh Mesh
	genMeshPoly(uintptr(unsafe.Pointer(&mesh)), int32(sides), radius)
	return mesh
}

// GenMeshPlane - Generate plane mesh (with subdivisions)
func GenMeshPlane(width float32, length float32, resX int, resZ int) Mesh {
	var mesh Mesh
	genMeshPlane(uintptr(unsafe.Pointer(&mesh)), width, length, int32(resX), int32(resZ))
	return mesh
}

// GenMeshCube - Generate cuboid mesh
func GenMeshCube(width float32, height float32, length float32) Mesh {
	var mesh Mesh
	genMeshCube(uintptr(unsafe.Pointer(&mesh)), width, height, length)
	return mesh
}

// GenMeshSphere - Generate sphere mesh (standard sphere)
func GenMeshSphere(radius float32, rings int, slices int) Mesh {
	var mesh Mesh
	genMeshSphere(uintptr(unsafe.Pointer(&mesh)), radius, int32(rings), int32(slices))
	return mesh
}

// GenMeshHemiSphere - Generate half-sphere mesh (no bottom cap)
func GenMeshHemiSphere(radius float32, rings int, slices int) Mesh {
	var mesh Mesh
	genMeshHemiSphere(uintptr(unsafe.Pointer(&mesh)), radius, int32(rings), int32(slices))
	return mesh
}

// GenMeshCylinder - Generate cylinder mesh
func GenMeshCylinder(radius float32, height float32, slices int) Mesh {
	var mesh Mesh
	genMeshCylinder(uintptr(unsafe.Pointer(&mesh)), radius, height, int32(slices))
	return mesh
}

// GenMeshCone - Generate cone/pyramid mesh
func GenMeshCone(radius float32, height float32, slices int) Mesh {
	var mesh Mesh
	genMeshCone(uintptr(unsafe.Pointer(&mesh)), radius, height, int32(slices))
	return mesh
}

// GenMeshTorus - Generate torus mesh
func GenMeshTorus(radius float32, size float32, radSeg int, sides int) Mesh {
	var mesh Mesh
	genMeshTorus(uintptr(unsafe.Pointer(&mesh)), radius, size, int32(radSeg), int32(sides))
	return mesh
}

// GenMeshKnot - Generate trefoil knot mesh
func GenMeshKnot(radius float32, size float32, radSeg int, sides int) Mesh {
	var mesh Mesh
	genMeshKnot(uintptr(unsafe.Pointer(&mesh)), radius, size, int32(radSeg), int32(sides))
	return mesh
}

// GenMeshHeightmap - Generate heightmap mesh from image data
func GenMeshHeightmap(heightmap Image, size Vector3) Mesh {
	var mesh Mesh
	genMeshHeightmap(uintptr(unsafe.Pointer(&mesh)), uintptr(unsafe.Pointer(&heightmap)), uintptr(unsafe.Pointer(&size)))
	return mesh
}

// GenMeshCubicmap - Generate cubes-based map mesh from image data
func GenMeshCubicmap(cubicmap Image, cubeSize Vector3) Mesh {
	var mesh Mesh
	genMeshCubicmap(uintptr(unsafe.Pointer(&mesh)), uintptr(unsafe.Pointer(&cubicmap)), uintptr(unsafe.Pointer(&cubeSize)))
	return mesh
}

// LoadMaterials - Load materials from model file
func LoadMaterials(fileName string) []Material {
	var materialCount int32
	ret := loadMaterials(fileName, &materialCount)
	return unsafe.Slice(ret, materialCount)
}

// LoadMaterialDefault - Load default material (Supports: DIFFUSE, SPECULAR, NORMAL maps)
func LoadMaterialDefault() Material {
	var material Material
	loadMaterialDefault(uintptr(unsafe.Pointer(&material)))
	return material
}

// IsMaterialValid - Check if a material is valid (shader assigned, map textures loaded in GPU)
func IsMaterialValid(material Material) bool {
	return isMaterialValid(uintptr(unsafe.Pointer(&material)))
}

// UnloadMaterial - Unload material from GPU memory (VRAM)
func UnloadMaterial(material Material) {
	unloadMaterial(uintptr(unsafe.Pointer(&material)))
}

// SetMaterialTexture - Set texture for a material map type (MATERIAL_MAP_DIFFUSE, MATERIAL_MAP_SPECULAR...)
func SetMaterialTexture(material *Material, mapType int32, texture Texture2D) {
	setMaterialTexture(material, mapType, uintptr(unsafe.Pointer(&texture)))
}

// SetModelMeshMaterial - Set material for a mesh
func SetModelMeshMaterial(model *Model, meshId int32, materialId int32) {
	setModelMeshMaterial(model, meshId, materialId)
}

// LoadModelAnimations - Load model animations from file
func LoadModelAnimations(fileName string) []ModelAnimation {
	var animCount int32
	ret := loadModelAnimations(fileName, &animCount)
	return unsafe.Slice(ret, animCount)
}

// UpdateModelAnimation - Update model animation pose (CPU)
func UpdateModelAnimation(model Model, anim ModelAnimation, frame int32) {
	updateModelAnimation(uintptr(unsafe.Pointer(&model)), uintptr(unsafe.Pointer(&anim)), frame)
}

// UpdateModelAnimationBones - Update model animation mesh bone matrices (GPU skinning)
func UpdateModelAnimationBones(model Model, anim ModelAnimation, frame int32) {
	updateModelAnimationBones(uintptr(unsafe.Pointer(&model)), uintptr(unsafe.Pointer(&anim)), frame)
}

// UnloadModelAnimation - Unload animation data
func UnloadModelAnimation(anim ModelAnimation) {
	unloadModelAnimation(uintptr(unsafe.Pointer(&anim)))
}

// UnloadModelAnimations - Unload animation array data
func UnloadModelAnimations(animations []ModelAnimation) {
	animCount := int32(len(animations))
	unloadModelAnimations(unsafe.SliceData(animations), animCount)
}

// IsModelAnimationValid - Check model animation skeleton match
func IsModelAnimationValid(model Model, anim ModelAnimation) bool {
	return isModelAnimationValid(uintptr(unsafe.Pointer(&model)), uintptr(unsafe.Pointer(&anim)))
}

// CheckCollisionSpheres - Check collision between two spheres
func CheckCollisionSpheres(center1 Vector3, radius1 float32, center2 Vector3, radius2 float32) bool {
	return checkCollisionSpheres(uintptr(unsafe.Pointer(&center1)), radius1, uintptr(unsafe.Pointer(&center2)), radius2)
}

// CheckCollisionBoxes - Check collision between two bounding boxes
func CheckCollisionBoxes(box1 BoundingBox, box2 BoundingBox) bool {
	return checkCollisionBoxes(uintptr(unsafe.Pointer(&box1)), uintptr(unsafe.Pointer(&box2)))
}

// CheckCollisionBoxSphere - Check collision between box and sphere
func CheckCollisionBoxSphere(box BoundingBox, center Vector3, radius float32) bool {
	return checkCollisionBoxSphere(uintptr(unsafe.Pointer(&box)), uintptr(unsafe.Pointer(&center)), radius)
}

// GetRayCollisionSphere - Get collision info between ray and sphere
func GetRayCollisionSphere(ray Ray, center Vector3, radius float32) RayCollision {
	var rayCollision RayCollision
	getRayCollisionSphere(uintptr(unsafe.Pointer(&rayCollision)), uintptr(unsafe.Pointer(&ray)), uintptr(unsafe.Pointer(&center)), radius)
	return rayCollision
}

// GetRayCollisionBox - Get collision info between ray and box
func GetRayCollisionBox(ray Ray, box BoundingBox) RayCollision {
	var rayCollision RayCollision
	getRayCollisionBox(uintptr(unsafe.Pointer(&rayCollision)), uintptr(unsafe.Pointer(&ray)), uintptr(unsafe.Pointer(&box)))
	return rayCollision
}

// GetRayCollisionMesh - Get collision info between ray and mesh
func GetRayCollisionMesh(ray Ray, mesh Mesh, transform Matrix) RayCollision {
	var rayCollision RayCollision
	getRayCollisionMesh(uintptr(unsafe.Pointer(&rayCollision)), uintptr(unsafe.Pointer(&ray)), uintptr(unsafe.Pointer(&mesh)), uintptr(unsafe.Pointer(&transform)))
	return rayCollision
}

// GetRayCollisionTriangle - Get collision info between ray and triangle
func GetRayCollisionTriangle(ray Ray, p1 Vector3, p2 Vector3, p3 Vector3) RayCollision {
	var rayCollision RayCollision
	getRayCollisionTriangle(uintptr(unsafe.Pointer(&rayCollision)), uintptr(unsafe.Pointer(&ray)), uintptr(unsafe.Pointer(&p1)), uintptr(unsafe.Pointer(&p2)), uintptr(unsafe.Pointer(&p3)))
	return rayCollision
}

// GetRayCollisionQuad - Get collision info between ray and quad
func GetRayCollisionQuad(ray Ray, p1 Vector3, p2 Vector3, p3 Vector3, p4 Vector3) RayCollision {
	var rayCollision RayCollision
	getRayCollisionQuad(uintptr(unsafe.Pointer(&rayCollision)), uintptr(unsafe.Pointer(&ray)), uintptr(unsafe.Pointer(&p1)), uintptr(unsafe.Pointer(&p2)), uintptr(unsafe.Pointer(&p3)), uintptr(unsafe.Pointer(&p4)))
	return rayCollision
}

// InitAudioDevice - Initialize audio device and context
func InitAudioDevice() {
	initAudioDevice()
}

// CloseAudioDevice - Close the audio device and context
func CloseAudioDevice() {
	closeAudioDevice()
}

// IsAudioDeviceReady - Check if audio device has been initialized successfully
func IsAudioDeviceReady() bool {
	return isAudioDeviceReady()
}

// SetMasterVolume - Set master volume (listener)
func SetMasterVolume(volume float32) {
	setMasterVolume(volume)
}

// GetMasterVolume - Get master volume (listener)
func GetMasterVolume() float32 {
	return getMasterVolume()
}

// LoadWave - Load wave data from file
func LoadWave(fileName string) Wave {
	var wave Wave
	loadWave(uintptr(unsafe.Pointer(&wave)), fileName)
	return wave
}

// LoadWaveFromMemory - Load wave from memory buffer, fileType refers to extension: i.e. '.wav'
func LoadWaveFromMemory(fileType string, fileData []byte, dataSize int32) Wave {
	var wave Wave
	loadWaveFromMemory(uintptr(unsafe.Pointer(&wave)), fileType, fileData, dataSize)
	return wave
}

// IsWaveValid - Checks if wave data is valid (data loaded and parameters)
func IsWaveValid(wave Wave) bool {
	return isWaveValid(uintptr(unsafe.Pointer(&wave)))
}

// LoadSound - Load sound from file
func LoadSound(fileName string) Sound {
	var sound Sound
	loadSound(uintptr(unsafe.Pointer(&sound)), fileName)
	return sound
}

// LoadSoundFromWave - Load sound from wave data
func LoadSoundFromWave(wave Wave) Sound {
	var sound Sound
	loadSoundFromWave(uintptr(unsafe.Pointer(&sound)), uintptr(unsafe.Pointer(&wave)))
	return sound
}

// LoadSoundAlias - Create a new sound that shares the same sample data as the source sound, does not own the sound data
func LoadSoundAlias(source Sound) Sound {
	var sound Sound
	loadSoundAlias(uintptr(unsafe.Pointer(&sound)), uintptr(unsafe.Pointer(&source)))
	return sound
}

// IsSoundValid - Checks if a sound is valid (data loaded and buffers initialized)
func IsSoundValid(sound Sound) bool {
	return isSoundValid(uintptr(unsafe.Pointer(&sound)))
}

// UpdateSound - Update sound buffer with new data
func UpdateSound(sound Sound, data []byte, sampleCount int32) {
	updateSound(uintptr(unsafe.Pointer(&sound)), data, sampleCount)
}

// UnloadWave - Unload wave data
func UnloadWave(wave Wave) {
	unloadWave(uintptr(unsafe.Pointer(&wave)))
}

// UnloadSound - Unload sound
func UnloadSound(sound Sound) {
	unloadSound(uintptr(unsafe.Pointer(&sound)))
}

// UnloadSoundAlias - Unload a sound alias (does not deallocate sample data)
func UnloadSoundAlias(alias Sound) {
	unloadSoundAlias(uintptr(unsafe.Pointer(&alias)))
}

// ExportWave - Export wave data to file, returns true on success
func ExportWave(wave Wave, fileName string) bool {
	return exportWave(uintptr(unsafe.Pointer(&wave)), fileName)
}

// PlaySound - Play a sound
func PlaySound(sound Sound) {
	playSound(uintptr(unsafe.Pointer(&sound)))
}

// StopSound - Stop playing a sound
func StopSound(sound Sound) {
	stopSound(uintptr(unsafe.Pointer(&sound)))
}

// PauseSound - Pause a sound
func PauseSound(sound Sound) {
	pauseSound(uintptr(unsafe.Pointer(&sound)))
}

// ResumeSound - Resume a paused sound
func ResumeSound(sound Sound) {
	resumeSound(uintptr(unsafe.Pointer(&sound)))
}

// IsSoundPlaying - Check if a sound is currently playing
func IsSoundPlaying(sound Sound) bool {
	return isSoundPlaying(uintptr(unsafe.Pointer(&sound)))
}

// SetSoundVolume - Set volume for a sound (1.0 is max level)
func SetSoundVolume(sound Sound, volume float32) {
	setSoundVolume(uintptr(unsafe.Pointer(&sound)), volume)
}

// SetSoundPitch - Set pitch for a sound (1.0 is base level)
func SetSoundPitch(sound Sound, pitch float32) {
	setSoundPitch(uintptr(unsafe.Pointer(&sound)), pitch)
}

// SetSoundPan - Set pan for a sound (0.5 is center)
func SetSoundPan(sound Sound, pan float32) {
	setSoundPan(uintptr(unsafe.Pointer(&sound)), pan)
}

// WaveCopy - Copy a wave to a new wave
func WaveCopy(wave Wave) Wave {
	var copy Wave
	waveCopy(uintptr(unsafe.Pointer(&copy)), uintptr(unsafe.Pointer(&wave)))
	return copy
}

// WaveCrop - Crop a wave to defined frames range
func WaveCrop(wave *Wave, initFrame int32, finalFrame int32) {
	waveCrop(wave, initFrame, finalFrame)
}

// WaveFormat - Convert wave data to desired format
func WaveFormat(wave *Wave, sampleRate int32, sampleSize int32, channels int32) {
	waveFormat(wave, sampleRate, sampleRate, channels)
}

// LoadWaveSamples - Load samples data from wave as a 32bit float data array
func LoadWaveSamples(wave Wave) []float32 {
	ret := loadWaveSamples(uintptr(unsafe.Pointer(&wave)))
	return unsafe.Slice(ret, wave.FrameCount*wave.Channels)
}

// UnloadWaveSamples - Unload samples data loaded with LoadWaveSamples()
func UnloadWaveSamples(samples []float32) {
	unloadWaveSamples(samples)
}

// LoadMusicStream - Load music stream from file
func LoadMusicStream(fileName string) Music {
	var music Music
	loadMusicStream(uintptr(unsafe.Pointer(&music)), fileName)
	return music
}

// LoadMusicStreamFromMemory - Load music stream from data
func LoadMusicStreamFromMemory(fileType string, data []byte, dataSize int32) Music {
	var music Music
	loadMusicStreamFromMemory(uintptr(unsafe.Pointer(&music)), fileType, data, dataSize)
	return music
}

// IsMusicValid - Checks if a music stream is valid (context and buffers initialized)
func IsMusicValid(music Music) bool {
	return isMusicValid(uintptr(unsafe.Pointer(&music)))
}

// UnloadMusicStream - Unload music stream
func UnloadMusicStream(music Music) {
	unloadMusicStream(uintptr(unsafe.Pointer(&music)))
}

// PlayMusicStream - Start music playing
func PlayMusicStream(music Music) {
	playMusicStream(uintptr(unsafe.Pointer(&music)))
}

// IsMusicStreamPlaying - Check if music is playing
func IsMusicStreamPlaying(music Music) bool {
	return isMusicStreamPlaying(uintptr(unsafe.Pointer(&music)))
}

// UpdateMusicStream - Updates buffers for music streaming
func UpdateMusicStream(music Music) {
	updateMusicStream(uintptr(unsafe.Pointer(&music)))
}

// StopMusicStream - Stop music playing
func StopMusicStream(music Music) {
	stopMusicStream(uintptr(unsafe.Pointer(&music)))
}

// PauseMusicStream - Pause music playing
func PauseMusicStream(music Music) {
	pauseMusicStream(uintptr(unsafe.Pointer(&music)))
}

// ResumeMusicStream - Resume playing paused music
func ResumeMusicStream(music Music) {
	resumeMusicStream(uintptr(unsafe.Pointer(&music)))
}

// SeekMusicStream - Seek music to a position (in seconds)
func SeekMusicStream(music Music, position float32) {
	seekMusicStream(uintptr(unsafe.Pointer(&music)), position)
}

// SetMusicVolume - Set volume for music (1.0 is max level)
func SetMusicVolume(music Music, volume float32) {
	setMusicVolume(uintptr(unsafe.Pointer(&music)), volume)
}

// SetMusicPitch - Set pitch for a music (1.0 is base level)
func SetMusicPitch(music Music, pitch float32) {
	setMusicPitch(uintptr(unsafe.Pointer(&music)), pitch)
}

// SetMusicPan - Set pan for a music (0.5 is center)
func SetMusicPan(music Music, pan float32) {
	setMusicPan(uintptr(unsafe.Pointer(&music)), pan)
}

// GetMusicTimeLength - Get music time length (in seconds)
func GetMusicTimeLength(music Music) float32 {
	return getMusicTimeLength(uintptr(unsafe.Pointer(&music)))
}

// GetMusicTimePlayed - Get current music time played (in seconds)
func GetMusicTimePlayed(music Music) float32 {
	return getMusicTimePlayed(uintptr(unsafe.Pointer(&music)))
}

// LoadAudioStream - Load audio stream (to stream raw audio pcm data)
func LoadAudioStream(sampleRate uint32, sampleSize uint32, channels uint32) AudioStream {
	var audioStream AudioStream
	loadAudioStream(uintptr(unsafe.Pointer(&audioStream)), sampleRate, sampleSize, channels)
	return audioStream
}

// IsAudioStreamValid - Checks if an audio stream is valid (buffers initialized)
func IsAudioStreamValid(stream AudioStream) bool {
	return isAudioStreamValid(uintptr(unsafe.Pointer(&stream)))
}

// UnloadAudioStream - Unload audio stream and free memory
func UnloadAudioStream(stream AudioStream) {
	unloadAudioStream(uintptr(unsafe.Pointer(&stream)))
}

// UpdateAudioStream - Update audio stream buffers with data
func UpdateAudioStream(stream AudioStream, data []float32) {
	frameCount := int32(len(data))
	updateAudioStream(uintptr(unsafe.Pointer(&stream)), data, frameCount)
}

// IsAudioStreamProcessed - Check if any audio stream buffers requires refill
func IsAudioStreamProcessed(stream AudioStream) bool {
	return isAudioStreamProcessed(uintptr(unsafe.Pointer(&stream)))
}

// PlayAudioStream - Play audio stream
func PlayAudioStream(stream AudioStream) {
	playAudioStream(uintptr(unsafe.Pointer(&stream)))
}

// PauseAudioStream - Pause audio stream
func PauseAudioStream(stream AudioStream) {
	pauseAudioStream(uintptr(unsafe.Pointer(&stream)))
}

// ResumeAudioStream - Resume audio stream
func ResumeAudioStream(stream AudioStream) {
	resumeAudioStream(uintptr(unsafe.Pointer(&stream)))
}

// IsAudioStreamPlaying - Check if audio stream is playing
func IsAudioStreamPlaying(stream AudioStream) bool {
	return isAudioStreamPlaying(uintptr(unsafe.Pointer(&stream)))
}

// StopAudioStream - Stop audio stream
func StopAudioStream(stream AudioStream) {
	stopAudioStream(uintptr(unsafe.Pointer(&stream)))
}

// SetAudioStreamVolume - Set volume for audio stream (1.0 is max level)
func SetAudioStreamVolume(stream AudioStream, volume float32) {
	setAudioStreamVolume(uintptr(unsafe.Pointer(&stream)), volume)
}

// SetAudioStreamPitch - Set pitch for audio stream (1.0 is base level)
func SetAudioStreamPitch(stream AudioStream, pitch float32) {
	setAudioStreamPitch(uintptr(unsafe.Pointer(&stream)), pitch)
}

// SetAudioStreamPan - Set pan for audio stream (0.5 is centered)
func SetAudioStreamPan(stream AudioStream, pan float32) {
	setAudioStreamPan(uintptr(unsafe.Pointer(&stream)), pan)
}

// SetAudioStreamBufferSizeDefault - Default size for new audio streams
func SetAudioStreamBufferSizeDefault(size int32) {
	setAudioStreamBufferSizeDefault(size)
}

// SetAudioStreamCallback - Audio thread callback to request new data
func SetAudioStreamCallback(stream AudioStream, callback AudioCallback) {
	fn := purego.NewCallback(func(bufferData unsafe.Pointer, frames int32) uintptr {
		callback(unsafe.Slice((*float32)(bufferData), frames), int(frames))
		return 0
	})
	setAudioStreamCallback(uintptr(unsafe.Pointer(&stream)), fn)
}

// AttachAudioStreamProcessor - Attach audio stream processor to stream, receives the samples as <float>s
func AttachAudioStreamProcessor(stream AudioStream, processor AudioCallback) {
	fn := purego.NewCallback(func(bufferData unsafe.Pointer, frames int32) uintptr {
		processor(unsafe.Slice((*float32)(bufferData), frames), int(frames))
		return 0
	})
	ptr := uintptr(reflect.ValueOf(processor).UnsafePointer())
	audioCallbacks[ptr] = fn
	attachAudioStreamProcessor(uintptr(unsafe.Pointer(&stream)), fn)
}

// DetachAudioStreamProcessor - Detach audio stream processor from stream
func DetachAudioStreamProcessor(stream AudioStream, processor AudioCallback) {
	ptr := uintptr(reflect.ValueOf(processor).UnsafePointer())
	fn := audioCallbacks[ptr]
	detachAudioStreamProcessor(uintptr(unsafe.Pointer(&stream)), fn)
}

// AttachAudioMixedProcessor - Attach audio stream processor to the entire audio pipeline, receives the samples as <float>s
func AttachAudioMixedProcessor(processor AudioCallback) {
	fn := purego.NewCallback(func(bufferData unsafe.Pointer, frames int32) uintptr {
		processor(unsafe.Slice((*float32)(bufferData), frames), int(frames))
		return 0
	})
	ptr := uintptr(reflect.ValueOf(processor).UnsafePointer())
	audioCallbacks[ptr] = fn
	attachAudioMixedProcessor(fn)
}

// DetachAudioMixedProcessor - Detach audio stream processor from the entire audio pipeline
func DetachAudioMixedProcessor(processor AudioCallback) {
	ptr := uintptr(reflect.ValueOf(processor).UnsafePointer())
	fn := audioCallbacks[ptr]
	detachAudioMixedProcessor(fn)
}

// SetCallbackFunc - Sets callback function
func SetCallbackFunc(func()) {
}

// NewImageFromImage - Returns new Image from Go image.Image
func NewImageFromImage(img image.Image) *Image {
	size := img.Bounds().Size()

	ret := GenImageColor(size.X, size.Y, White)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			col := img.At(x, y)
			r, g, b, a := col.RGBA()
			rcolor := NewColor(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
			ImageDrawPixel(ret, int32(x), int32(y), rcolor)
		}
	}

	return ret
}

// ToImage converts a Image to Go image.Image
func (i *Image) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(i.Width), int(i.Height)))

	// Get pixel data from image (RGBA 32bit)
	ret := LoadImageColors(i)
	pixels := (*[1 << 24]uint8)(unsafe.Pointer(unsafe.SliceData(ret)))[0 : i.Width*i.Height*4]

	img.Pix = pixels

	return img
}

// OpenAsset - Open asset
func OpenAsset(name string) (Asset, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// HomeDir - Returns user home directory
// NOTE: On Android this returns internal data path and must be called after InitWindow
func HomeDir() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir
	}
	return ""
}
