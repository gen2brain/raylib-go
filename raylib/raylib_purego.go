//go:build !cgo
// +build !cgo

package rl

import (
	"image/color"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	// raylibDll is the pointer to the shared library
	raylibDll uintptr
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
var loadShader func(shader uintptr, vsFileName string, fsFileName string)
var loadShaderFromMemory func(shader uintptr, vsCode string, fsCode string)
var isShaderReady func(shader uintptr) bool
var getShaderLocation func(shader uintptr, uniformName string) int32
var getShaderLocationAttrib func(shader uintptr, attribName string) int32
var setShaderValue func(shader uintptr, locIndex int32, value []float32, uniformType int32)
var setShaderValueV func(shader uintptr, locIndex int32, value []float32, uniformType int32, count int32)
var setShaderValueMatrix func(shader uintptr, locIndex int32, mat uintptr)
var setShaderValueTexture func(shader uintptr, locIndex int32, texture uintptr)
var unloadShader func(shader uintptr)
var getMouseRay func(ray uintptr, mousePosition uintptr, camera uintptr)
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
var traceLog func(logLevel int32, text string, args ...any)
var setTraceLogLevel func(logLevel int32)
var memAlloc func(size uint32) unsafe.Pointer
var memRealloc func(ptr unsafe.Pointer, size uint32) unsafe.Pointer
var memFree func(ptr unsafe.Pointer)
var setTraceLogCallback func(callback uintptr)

// var setLoadFileDataCallback func(callback uintptr)
// var setSaveFileDataCallback func(callback uintptr)
// var setLoadFileTextCallback func(callback uintptr)
// var setSaveFileTextCallback func(callback uintptr)
var loadFileData func(fileName string, dataSize *int32) *byte
var unloadFileData func(data *byte)
var saveFileData func(fileName string, data unsafe.Pointer, dataSize int32) bool
var exportDataAsCode func(data *byte, dataSize int32, fileName string) bool

// var loadFileText func(fileName string) string
// var unloadFileText func(text string)
// var saveFileText func(fileName string, text string) bool
// var fileExists func(fileName string) bool
// var directoryExists func(dirPath string) bool
// var isFileExtension func(fileName string, ext string) bool
// var getFileLength func(fileName string) int32
// var getFileExtension func(fileName string) string
// var getFileName func(filePath string) string
// var getFileNameWithoutExt func(filePath string) string
// var getDirectoryPath func(filePath string) string
// var getPrevDirectoryPath func(dirPath string) string
var getWorkingDirectory func() string
var getApplicationDirectory func() string
var changeDirectory func(dir string) bool

// var isPathFile func(path string) bool
// var loadDirectoryFiles func(dirPath string) FilePathList
// var loadDirectoryFilesEx func(basePath string, filter string, scanSubdirs bool) FilePathList
// var unloadDirectoryFiles func(files uintptr)
var isFileDropped func() bool
var loadDroppedFiles func(files uintptr)
var unloadDroppedFiles func(files uintptr)

// var getFileModTime func(fileName string) int64
// var compressData func(data []byte, dataSize int32, compDataSize []int32) []byte
// var decompressData func(compData []byte, compDataSize int32, dataSize []int32) []byte
// var encodeDataBase64 func(data []byte, dataSize int32, outputSize []int32) string
// var decodeDataBase64 func(data []byte, outputSize []int32) []byte
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

// var updateCamera func(camera uintptr, mode int32)
// var updateCameraPro func(camera uintptr, movement uintptr, rotation uintptr, zoom float32)
var setShapesTexture func(texture uintptr, source uintptr)
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
var drawCircleGradient func(centerX int32, centerY int32, radius float32, color1 uintptr, color2 uintptr)
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
var drawRectangleGradientV func(posX int32, posY int32, width int32, height int32, color1 uintptr, color2 uintptr)
var drawRectangleGradientH func(posX int32, posY int32, width int32, height int32, color1 uintptr, color2 uintptr)
var drawRectangleGradientEx func(rec uintptr, col1 uintptr, col2 uintptr, col3 uintptr, col4 uintptr)
var drawRectangleLines func(posX int32, posY int32, width int32, height int32, col uintptr)
var drawRectangleLinesEx func(rec uintptr, lineThick float32, col uintptr)
var drawRectangleRounded func(rec uintptr, roundness float32, segments int32, col uintptr)
var drawRectangleRoundedLines func(rec uintptr, roundness float32, segments int32, lineThick float32, col uintptr)
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
var getSplinePointLinear func(startPos uintptr, endPos uintptr, t float32) Vector2
var getSplinePointBasis func(p1 uintptr, p2 uintptr, p3 uintptr, p4 uintptr, t float32) Vector2
var getSplinePointCatmullRom func(p1 uintptr, p2 uintptr, p3 uintptr, p4 uintptr, t float32) Vector2
var getSplinePointBezierQuad func(p1 uintptr, c2 uintptr, p3 uintptr, t float32) Vector2
var getSplinePointBezierCubic func(p1 uintptr, c2 uintptr, c3 uintptr, p4 uintptr, t float32) Vector2
var checkCollisionRecs func(rec1 uintptr, rec2 uintptr) bool
var checkCollisionCircles func(center1 uintptr, radius1 float32, center2 uintptr, radius2 float32) bool
var checkCollisionCircleRec func(center uintptr, radius float32, rec uintptr) bool
var checkCollisionPointRec func(point uintptr, rec uintptr) bool
var checkCollisionPointCircle func(point uintptr, center uintptr, radius float32) bool
var checkCollisionPointTriangle func(point uintptr, p1 uintptr, p2 uintptr, p3 uintptr) bool
var checkCollisionPointPoly func(point uintptr, points *Vector2, pointCount int32) bool
var checkCollisionLines func(startPos1 uintptr, endPos1 uintptr, startPos2 uintptr, endPos2 uintptr, collisionPoint *Vector2) bool
var checkCollisionPointLine func(point uintptr, p1 uintptr, p2 uintptr, threshold int32) bool
var getCollisionRec func(rec uintptr, rec1 uintptr, rec2 uintptr)
var loadImage func(img uintptr, fileName string)
var loadImageRaw func(img uintptr, fileName string, width int32, height int32, format int32, headerSize int32)
var loadImageSvg func(img uintptr, fileNameOrString string, width int32, height int32)
var loadImageAnim func(img uintptr, fileName string, frames []int32)
var loadImageFromMemory func(img uintptr, fileType string, fileData []byte, dataSize int32)
var loadImageFromTexture func(img uintptr, texture uintptr)
var loadImageFromScreen func(img uintptr)
var isImageReady func(image uintptr) bool
var unloadImage func(image uintptr)
var exportImage func(image uintptr, fileName string) bool
var exportImageToMemory func(image uintptr, fileType string, fileSize *int32) *byte
var exportImageAsCode func(image uintptr, fileName string) bool
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
var imageDrawCircle func(dst *Image, centerX int32, centerY int32, radius int32, col uintptr)
var imageDrawCircleV func(dst *Image, center uintptr, radius int32, col uintptr)
var imageDrawCircleLines func(dst *Image, centerX int32, centerY int32, radius int32, col uintptr)
var imageDrawCircleLinesV func(dst *Image, center uintptr, radius int32, col uintptr)
var imageDrawRectangle func(dst *Image, posX int32, posY int32, width int32, height int32, col uintptr)
var imageDrawRectangleV func(dst *Image, position uintptr, size uintptr, col uintptr)
var imageDrawRectangleRec func(dst *Image, rec uintptr, col uintptr)
var imageDrawRectangleLines func(dst *Image, rec uintptr, thick int32, col uintptr)
var imageDraw func(dst *Image, src uintptr, srcRec uintptr, dstRec uintptr, tint uintptr)
var imageDrawText func(dst *Image, text string, posX int32, posY int32, fontSize int32, col uintptr)
var imageDrawTextEx func(dst *Image, font uintptr, text string, position uintptr, fontSize float32, spacing float32, tint uintptr)
var loadTexture func(texture uintptr, fileName string)
var loadTextureFromImage func(texture uintptr, image uintptr)
var loadTextureCubemap func(texture uintptr, image uintptr, layout int32)
var loadRenderTexture func(texture uintptr, width int32, height int32)
var isTextureReady func(texture uintptr) bool
var unloadTexture func(texture uintptr)
var isRenderTextureReady func(target uintptr) bool
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
var getColor func(hexValue uint32) uintptr
var getPixelColor func(srcPtr unsafe.Pointer, format int32) uintptr
var setPixelColor func(dstPtr unsafe.Pointer, col uintptr, format int32)
var getPixelDataSize func(width int32, height int32, format int32) int32
var getFontDefault func(font uintptr)
var loadFont func(font uintptr, fileName string)
var loadFontEx func(font uintptr, fileName string, fontSize int32, codepoints []int32, codepointCount int32)
var loadFontFromImage func(font uintptr, image uintptr, key uintptr, firstChar int32)
var loadFontFromMemory func(font uintptr, fileType string, fileData []byte, dataSize int32, fontSize int32, codepoints []int32, codepointCount int32)
var isFontReady func(font uintptr) bool
var loadFontData func(fileData []byte, dataSize int32, fontSize int32, codepoints []int32, codepointCount int32, _type int32) *GlyphInfo
var genImageFontAtlas func(image uintptr, glyphs *GlyphInfo, glyphRecs []*Rectangle, glyphCount int32, fontSize int32, padding int32, packMethod int32)
var unloadFontData func(glyphs *GlyphInfo, glyphCount int32)
var unloadFont func(font uintptr)
var exportFontAsCode func(font uintptr, fileName string) bool
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

// var loadUTF8 func(codepoints []int32, length int32) string
// var unloadUTF8 func(text string)
// var loadCodepoints func(text string, count []int32) []int32
// var unloadCodepoints func(codepoints []int32)
// var getCodepointCount func(text string) int32
// var getCodepoint func(text string, codepointSize []int32) int32
// var getCodepointNext func(text string, codepointSize []int32) int32
// var getCodepointPrevious func(text string, codepointSize []int32) int32
// var codepointToUTF8 func(codepoint int32, utf8Size []int32) string
// var textCopy func(dst string, src string) int32
// var textIsEqual func(text1 string, text2 string) bool
// var textLength func(text string) uint32
// var textFormat func(text string, args uintptr) string
// var textSubtext func(text string, position int32, length int32) string
// var textReplace func(text string, replace string, by string) string
// var textInsert func(text string, insert string, position int32) string
// var textJoin func(textList []string, count int32, delimiter string) string
// var textSplit func(text string, delimiter int8, count []int32) []string
// var textAppend func(text string, _append string, position []int32)
// var textFindIndex func(text string, find string) int32
// var textToUpper func(text string) string
// var textToLower func(text string) string
// var textToPascal func(text string) string
// var textToInteger func(text string) int32
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
var loadModel func(fileName string) Model
var loadModelFromMesh func(mesh uintptr) Model
var isModelReady func(model uintptr) bool
var unloadModel func(model uintptr)
var getModelBoundingBox func(model uintptr) BoundingBox
var drawModel func(model uintptr, position uintptr, scale float32, tint uintptr)
var drawModelEx func(model uintptr, position uintptr, rotationAxis uintptr, rotationAngle float32, scale uintptr, tint uintptr)
var drawModelWires func(model uintptr, position uintptr, scale float32, tint uintptr)
var drawModelWiresEx func(model uintptr, position uintptr, rotationAxis uintptr, rotationAngle float32, scale uintptr, tint uintptr)
var drawBoundingBox func(box uintptr, col uintptr)
var drawBillboard func(camera uintptr, texture uintptr, position uintptr, size float32, tint uintptr)
var drawBillboardRec func(camera uintptr, texture uintptr, source uintptr, position uintptr, size uintptr, tint uintptr)
var drawBillboardPro func(camera uintptr, texture uintptr, source uintptr, position uintptr, up uintptr, size uintptr, origin uintptr, rotation float32, tint uintptr)
var uploadMesh func(mesh uintptr, dynamic bool)
var updateMeshBuffer func(mesh uintptr, index int32, data uintptr, dataSize int32, offset int32)
var unloadMesh func(mesh uintptr)
var drawMesh func(mesh uintptr, material uintptr, transform uintptr)
var drawMeshInstanced func(mesh uintptr, material uintptr, transforms uintptr, instances int32)
var exportMesh func(mesh uintptr, fileName string) bool
var getMeshBoundingBox func(mesh uintptr) BoundingBox
var genMeshTangents func(mesh uintptr)
var genMeshPoly func(sides int32, radius float32) Mesh
var genMeshPlane func(width float32, length float32, resX int32, resZ int32) Mesh
var genMeshCube func(width float32, height float32, length float32) Mesh
var genMeshSphere func(radius float32, rings int32, slices int32) Mesh
var genMeshHemiSphere func(radius float32, rings int32, slices int32) Mesh
var genMeshCylinder func(radius float32, height float32, slices int32) Mesh
var genMeshCone func(radius float32, height float32, slices int32) Mesh
var genMeshTorus func(radius float32, size float32, radSeg int32, sides int32) Mesh
var genMeshKnot func(radius float32, size float32, radSeg int32, sides int32) Mesh
var genMeshHeightmap func(heightmap uintptr, size uintptr) Mesh
var genMeshCubicmap func(cubicmap uintptr, cubeSize uintptr) Mesh
var loadMaterials func(fileName string, materialCount []int32) *Material
var loadMaterialDefault func() Material
var isMaterialReady func(material uintptr) bool
var unloadMaterial func(material uintptr)
var setMaterialTexture func(material uintptr, mapType int32, texture uintptr)
var setModelMeshMaterial func(model uintptr, meshId int32, materialId int32)
var loadModelAnimations func(fileName string, animCount []int32) *ModelAnimation
var updateModelAnimation func(model uintptr, anim uintptr, frame int32)
var unloadModelAnimation func(anim uintptr)
var unloadModelAnimations func(animations uintptr, animCount int32)
var isModelAnimationValid func(model uintptr, anim uintptr) bool
var checkCollisionSpheres func(center1 uintptr, radius1 float32, center2 uintptr, radius2 float32) bool
var checkCollisionBoxes func(box1 uintptr, box2 uintptr) bool
var checkCollisionBoxSphere func(box uintptr, center uintptr, radius float32) bool
var getRayCollisionSphere func(ray uintptr, center uintptr, radius float32) RayCollision
var getRayCollisionBox func(ray uintptr, box uintptr) RayCollision
var getRayCollisionMesh func(ray uintptr, mesh uintptr, transform uintptr) RayCollision
var getRayCollisionTriangle func(ray uintptr, p1 uintptr, p2 uintptr, p3 uintptr) RayCollision
var getRayCollisionQuad func(ray uintptr, p1 uintptr, p2 uintptr, p3 uintptr, p4 uintptr) RayCollision
var initAudioDevice func()
var closeAudioDevice func()
var isAudioDeviceReady func() bool
var setMasterVolume func(volume float32)
var getMasterVolume func() float32
var loadWave func(fileName string) Wave
var loadWaveFromMemory func(fileType string, fileData []byte, dataSize int32) Wave
var isWaveReady func(wave uintptr) bool
var loadSound func(sound uintptr, fileName string)
var loadSoundFromWave func(sound uintptr, wave uintptr)
var loadSoundAlias func(sound uintptr, source uintptr)
var isSoundReady func(sound uintptr) bool
var updateSound func(sound uintptr, data uintptr, sampleCount int32)
var unloadWave func(wave uintptr)
var unloadSound func(sound uintptr)
var unloadSoundAlias func(alias uintptr)
var exportWave func(wave uintptr, fileName string) bool
var exportWaveAsCode func(wave uintptr, fileName string) bool
var playSound func(sound uintptr)
var stopSound func(sound uintptr)
var pauseSound func(sound uintptr)
var resumeSound func(sound uintptr)
var isSoundPlaying func(sound uintptr) bool
var setSoundVolume func(sound uintptr, volume float32)
var setSoundPitch func(sound uintptr, pitch float32)
var setSoundPan func(sound uintptr, pan float32)
var waveCopy func(wave uintptr) Wave
var waveCrop func(wave uintptr, initSample int32, finalSample int32)
var waveFormat func(wave uintptr, sampleRate int32, sampleSize int32, channels int32)
var loadWaveSamples func(wave uintptr) []float32
var unloadWaveSamples func(samples []float32)
var loadMusicStream func(music uintptr, fileName string)
var loadMusicStreamFromMemory func(sound uintptr, fileType string, data []byte, dataSize int32)
var isMusicReady func(music uintptr) bool
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
var isAudioStreamReady func(stream uintptr) bool
var unloadAudioStream func(stream uintptr)
var updateAudioStream func(stream uintptr, data uintptr, frameCount int32)
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
	purego.RegisterLibFunc(&isShaderReady, raylibDll, "IsShaderReady")
	purego.RegisterLibFunc(&getShaderLocation, raylibDll, "GetShaderLocation")
	purego.RegisterLibFunc(&getShaderLocationAttrib, raylibDll, "GetShaderLocationAttrib")
	purego.RegisterLibFunc(&setShaderValue, raylibDll, "SetShaderValue")
	purego.RegisterLibFunc(&setShaderValueV, raylibDll, "SetShaderValueV")
	purego.RegisterLibFunc(&setShaderValueMatrix, raylibDll, "SetShaderValueMatrix")
	purego.RegisterLibFunc(&setShaderValueTexture, raylibDll, "SetShaderValueTexture")
	purego.RegisterLibFunc(&unloadShader, raylibDll, "UnloadShader")
	purego.RegisterLibFunc(&getMouseRay, raylibDll, "GetMouseRay")
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
	// purego.RegisterLibFunc(&setLoadFileDataCallback, raylibDll, "SetLoadFileDataCallback")
	// purego.RegisterLibFunc(&setSaveFileDataCallback, raylibDll, "SetSaveFileDataCallback")
	// purego.RegisterLibFunc(&setLoadFileTextCallback, raylibDll, "SetLoadFileTextCallback")
	// purego.RegisterLibFunc(&setSaveFileTextCallback, raylibDll, "SetSaveFileTextCallback")
	purego.RegisterLibFunc(&loadFileData, raylibDll, "LoadFileData")
	purego.RegisterLibFunc(&unloadFileData, raylibDll, "UnloadFileData")
	purego.RegisterLibFunc(&saveFileData, raylibDll, "SaveFileData")
	purego.RegisterLibFunc(&exportDataAsCode, raylibDll, "ExportDataAsCode")
	// purego.RegisterLibFunc(&loadFileText, raylibDll, "LoadFileText")
	// purego.RegisterLibFunc(&unloadFileText, raylibDll, "UnloadFileText")
	// purego.RegisterLibFunc(&saveFileText, raylibDll, "SaveFileText")
	// purego.RegisterLibFunc(&fileExists, raylibDll, "FileExists")
	// purego.RegisterLibFunc(&directoryExists, raylibDll, "DirectoryExists")
	// purego.RegisterLibFunc(&isFileExtension, raylibDll, "IsFileExtension")
	// purego.RegisterLibFunc(&getFileLength, raylibDll, "GetFileLength")
	// purego.RegisterLibFunc(&getFileExtension, raylibDll, "GetFileExtension")
	// purego.RegisterLibFunc(&getFileName, raylibDll, "GetFileName")
	// purego.RegisterLibFunc(&getFileNameWithoutExt, raylibDll, "GetFileNameWithoutExt")
	// purego.RegisterLibFunc(&getDirectoryPath, raylibDll, "GetDirectoryPath")
	// purego.RegisterLibFunc(&getPrevDirectoryPath, raylibDll, "GetPrevDirectoryPath")
	purego.RegisterLibFunc(&getWorkingDirectory, raylibDll, "GetWorkingDirectory")
	purego.RegisterLibFunc(&getApplicationDirectory, raylibDll, "GetApplicationDirectory")
	purego.RegisterLibFunc(&changeDirectory, raylibDll, "ChangeDirectory")
	// purego.RegisterLibFunc(&isPathFile, raylibDll, "IsPathFile")
	// purego.RegisterLibFunc(&loadDirectoryFiles, raylibDll, "LoadDirectoryFiles")
	// purego.RegisterLibFunc(&loadDirectoryFilesEx, raylibDll, "LoadDirectoryFilesEx")
	// purego.RegisterLibFunc(&unloadDirectoryFiles, raylibDll, "UnloadDirectoryFiles")
	purego.RegisterLibFunc(&isFileDropped, raylibDll, "IsFileDropped")
	purego.RegisterLibFunc(&loadDroppedFiles, raylibDll, "LoadDroppedFiles")
	purego.RegisterLibFunc(&unloadDroppedFiles, raylibDll, "UnloadDroppedFiles")
	// purego.RegisterLibFunc(&getFileModTime, raylibDll, "GetFileModTime")
	// purego.RegisterLibFunc(&compressData, raylibDll, "CompressData")
	// purego.RegisterLibFunc(&decompressData, raylibDll, "DecompressData")
	// purego.RegisterLibFunc(&encodeDataBase64, raylibDll, "EncodeDataBase64")
	// purego.RegisterLibFunc(&decodeDataBase64, raylibDll, "DecodeDataBase64")
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
	// purego.RegisterLibFunc(&updateCamera, raylibDll, "UpdateCamera")
	// purego.RegisterLibFunc(&updateCameraPro, raylibDll, "UpdateCameraPro")
	purego.RegisterLibFunc(&setShapesTexture, raylibDll, "SetShapesTexture")
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
	purego.RegisterLibFunc(&checkCollisionPointRec, raylibDll, "CheckCollisionPointRec")
	purego.RegisterLibFunc(&checkCollisionPointCircle, raylibDll, "CheckCollisionPointCircle")
	purego.RegisterLibFunc(&checkCollisionPointTriangle, raylibDll, "CheckCollisionPointTriangle")
	purego.RegisterLibFunc(&checkCollisionPointPoly, raylibDll, "CheckCollisionPointPoly")
	purego.RegisterLibFunc(&checkCollisionLines, raylibDll, "CheckCollisionLines")
	purego.RegisterLibFunc(&checkCollisionPointLine, raylibDll, "CheckCollisionPointLine")
	purego.RegisterLibFunc(&getCollisionRec, raylibDll, "GetCollisionRec")
	purego.RegisterLibFunc(&loadImage, raylibDll, "LoadImage")
	purego.RegisterLibFunc(&loadImageRaw, raylibDll, "LoadImageRaw")
	purego.RegisterLibFunc(&loadImageSvg, raylibDll, "LoadImageSvg")
	purego.RegisterLibFunc(&loadImageAnim, raylibDll, "LoadImageAnim")
	purego.RegisterLibFunc(&loadImageFromMemory, raylibDll, "LoadImageFromMemory")
	purego.RegisterLibFunc(&loadImageFromTexture, raylibDll, "LoadImageFromTexture")
	purego.RegisterLibFunc(&loadImageFromScreen, raylibDll, "LoadImageFromScreen")
	purego.RegisterLibFunc(&isImageReady, raylibDll, "IsImageReady")
	purego.RegisterLibFunc(&unloadImage, raylibDll, "UnloadImage")
	purego.RegisterLibFunc(&exportImage, raylibDll, "ExportImage")
	purego.RegisterLibFunc(&exportImageToMemory, raylibDll, "ExportImageToMemory")
	purego.RegisterLibFunc(&exportImageAsCode, raylibDll, "ExportImageAsCode")
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
	purego.RegisterLibFunc(&imageDrawCircle, raylibDll, "ImageDrawCircle")
	purego.RegisterLibFunc(&imageDrawCircleV, raylibDll, "ImageDrawCircleV")
	purego.RegisterLibFunc(&imageDrawCircleLines, raylibDll, "ImageDrawCircleLines")
	purego.RegisterLibFunc(&imageDrawCircleLinesV, raylibDll, "ImageDrawCircleLinesV")
	purego.RegisterLibFunc(&imageDrawRectangle, raylibDll, "ImageDrawRectangle")
	purego.RegisterLibFunc(&imageDrawRectangleV, raylibDll, "ImageDrawRectangleV")
	purego.RegisterLibFunc(&imageDrawRectangleRec, raylibDll, "ImageDrawRectangleRec")
	purego.RegisterLibFunc(&imageDrawRectangleLines, raylibDll, "ImageDrawRectangleLines")
	purego.RegisterLibFunc(&imageDraw, raylibDll, "ImageDraw")
	purego.RegisterLibFunc(&imageDrawText, raylibDll, "ImageDrawText")
	purego.RegisterLibFunc(&imageDrawTextEx, raylibDll, "ImageDrawTextEx")
	purego.RegisterLibFunc(&loadTexture, raylibDll, "LoadTexture")
	purego.RegisterLibFunc(&loadTextureFromImage, raylibDll, "LoadTextureFromImage")
	purego.RegisterLibFunc(&loadTextureCubemap, raylibDll, "LoadTextureCubemap")
	purego.RegisterLibFunc(&loadRenderTexture, raylibDll, "LoadRenderTexture")
	purego.RegisterLibFunc(&isTextureReady, raylibDll, "IsTextureReady")
	purego.RegisterLibFunc(&unloadTexture, raylibDll, "UnloadTexture")
	purego.RegisterLibFunc(&isRenderTextureReady, raylibDll, "IsRenderTextureReady")
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
	purego.RegisterLibFunc(&getColor, raylibDll, "GetColor")
	purego.RegisterLibFunc(&getPixelColor, raylibDll, "GetPixelColor")
	purego.RegisterLibFunc(&setPixelColor, raylibDll, "SetPixelColor")
	purego.RegisterLibFunc(&getPixelDataSize, raylibDll, "GetPixelDataSize")
	purego.RegisterLibFunc(&getFontDefault, raylibDll, "GetFontDefault")
	purego.RegisterLibFunc(&loadFont, raylibDll, "LoadFont")
	purego.RegisterLibFunc(&loadFontEx, raylibDll, "LoadFontEx")
	purego.RegisterLibFunc(&loadFontFromImage, raylibDll, "LoadFontFromImage")
	purego.RegisterLibFunc(&loadFontFromMemory, raylibDll, "LoadFontFromMemory")
	purego.RegisterLibFunc(&isFontReady, raylibDll, "IsFontReady")
	purego.RegisterLibFunc(&loadFontData, raylibDll, "LoadFontData")
	purego.RegisterLibFunc(&genImageFontAtlas, raylibDll, "GenImageFontAtlas")
	purego.RegisterLibFunc(&unloadFontData, raylibDll, "UnloadFontData")
	purego.RegisterLibFunc(&unloadFont, raylibDll, "UnloadFont")
	purego.RegisterLibFunc(&exportFontAsCode, raylibDll, "ExportFontAsCode")
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
	// purego.RegisterLibFunc(&loadUTF8, raylibDll, "LoadUTF8")
	// purego.RegisterLibFunc(&unloadUTF8, raylibDll, "UnloadUTF8")
	// purego.RegisterLibFunc(&loadCodepoints, raylibDll, "LoadCodepoints")
	// purego.RegisterLibFunc(&unloadCodepoints, raylibDll, "UnloadCodepoints")
	// purego.RegisterLibFunc(&getCodepointCount, raylibDll, "GetCodepointCount")
	// purego.RegisterLibFunc(&getCodepoint, raylibDll, "GetCodepoint")
	// purego.RegisterLibFunc(&getCodepointNext, raylibDll, "GetCodepointNext")
	// purego.RegisterLibFunc(&getCodepointPrevious, raylibDll, "GetCodepointPrevious")
	// purego.RegisterLibFunc(&codepointToUTF8, raylibDll, "CodepointToUTF8")
	// purego.RegisterLibFunc(&textCopy, raylibDll, "TextCopy")
	// purego.RegisterLibFunc(&textIsEqual, raylibDll, "TextIsEqual")
	// purego.RegisterLibFunc(&textLength, raylibDll, "TextLength")
	// purego.RegisterLibFunc(&textFormat, raylibDll, "TextFormat")
	// purego.RegisterLibFunc(&textSubtext, raylibDll, "TextSubtext")
	// purego.RegisterLibFunc(&textReplace, raylibDll, "TextReplace")
	// purego.RegisterLibFunc(&textInsert, raylibDll, "TextInsert")
	// purego.RegisterLibFunc(&textJoin, raylibDll, "TextJoin")
	// purego.RegisterLibFunc(&textSplit, raylibDll, "TextSplit")
	// purego.RegisterLibFunc(&textAppend, raylibDll, "TextAppend")
	// purego.RegisterLibFunc(&textFindIndex, raylibDll, "TextFindIndex")
	// purego.RegisterLibFunc(&textToUpper, raylibDll, "TextToUpper")
	// purego.RegisterLibFunc(&textToLower, raylibDll, "TextToLower")
	// purego.RegisterLibFunc(&textToPascal, raylibDll, "TextToPascal")
	// purego.RegisterLibFunc(&textToInteger, raylibDll, "TextToInteger")
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
	purego.RegisterLibFunc(&isModelReady, raylibDll, "IsModelReady")
	purego.RegisterLibFunc(&unloadModel, raylibDll, "UnloadModel")
	purego.RegisterLibFunc(&getModelBoundingBox, raylibDll, "GetModelBoundingBox")
	purego.RegisterLibFunc(&drawModel, raylibDll, "DrawModel")
	purego.RegisterLibFunc(&drawModelEx, raylibDll, "DrawModelEx")
	purego.RegisterLibFunc(&drawModelWires, raylibDll, "DrawModelWires")
	purego.RegisterLibFunc(&drawModelWiresEx, raylibDll, "DrawModelWiresEx")
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
	purego.RegisterLibFunc(&isMaterialReady, raylibDll, "IsMaterialReady")
	purego.RegisterLibFunc(&unloadMaterial, raylibDll, "UnloadMaterial")
	purego.RegisterLibFunc(&setMaterialTexture, raylibDll, "SetMaterialTexture")
	purego.RegisterLibFunc(&setModelMeshMaterial, raylibDll, "SetModelMeshMaterial")
	purego.RegisterLibFunc(&loadModelAnimations, raylibDll, "LoadModelAnimations")
	purego.RegisterLibFunc(&updateModelAnimation, raylibDll, "UpdateModelAnimation")
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
	purego.RegisterLibFunc(&isWaveReady, raylibDll, "IsWaveReady")
	purego.RegisterLibFunc(&loadSound, raylibDll, "LoadSound")
	purego.RegisterLibFunc(&loadSoundFromWave, raylibDll, "LoadSoundFromWave")
	purego.RegisterLibFunc(&loadSoundAlias, raylibDll, "LoadSoundAlias")
	purego.RegisterLibFunc(&isSoundReady, raylibDll, "IsSoundReady")
	purego.RegisterLibFunc(&updateSound, raylibDll, "UpdateSound")
	purego.RegisterLibFunc(&unloadWave, raylibDll, "UnloadWave")
	purego.RegisterLibFunc(&unloadSound, raylibDll, "UnloadSound")
	purego.RegisterLibFunc(&unloadSoundAlias, raylibDll, "UnloadSoundAlias")
	purego.RegisterLibFunc(&exportWave, raylibDll, "ExportWave")
	purego.RegisterLibFunc(&exportWaveAsCode, raylibDll, "ExportWaveAsCode")
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
	purego.RegisterLibFunc(&isMusicReady, raylibDll, "IsMusicReady")
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
	purego.RegisterLibFunc(&isAudioStreamReady, raylibDll, "IsAudioStreamReady")
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
func SetWindowPosition(x int32, y int32) {
	setWindowPosition(x, y)
}

// SetWindowMonitor - Set monitor for the current window
func SetWindowMonitor(monitor int32) {
	setWindowMonitor(monitor)
}

// SetWindowMinSize - Set window minimum dimensions (for FLAG_WINDOW_RESIZABLE)
func SetWindowMinSize(width int32, height int32) {
	setWindowMinSize(width, height)
}

// SetWindowMaxSize - Set window maximum dimensions (for FLAG_WINDOW_RESIZABLE)
func SetWindowMaxSize(width int32, height int32) {
	setWindowMaxSize(width, height)
}

// SetWindowSize - Set window dimensions
func SetWindowSize(width int32, height int32) {
	setWindowSize(width, height)
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
func GetScreenWidth() int32 {
	return getScreenWidth()
}

// GetScreenHeight - Get current screen height
func GetScreenHeight() int32 {
	return getScreenHeight()
}

// GetRenderWidth - Get current render width (it considers HiDPI)
func GetRenderWidth() int32 {
	return getRenderWidth()
}

// GetRenderHeight - Get current render height (it considers HiDPI)
func GetRenderHeight() int32 {
	return getRenderHeight()
}

// GetMonitorCount - Get number of connected monitors
func GetMonitorCount() int32 {
	return getMonitorCount()
}

// GetCurrentMonitor - Get current connected monitor
func GetCurrentMonitor() int32 {
	return getCurrentMonitor()
}

// GetMonitorPosition - Get specified monitor position
func GetMonitorPosition(monitor int32) Vector2 {
	ret := getMonitorPosition(monitor)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetMonitorWidth - Get specified monitor width (current video mode used by monitor)
func GetMonitorWidth(monitor int32) int32 {
	return getMonitorWidth(monitor)
}

// GetMonitorHeight - Get specified monitor height (current video mode used by monitor)
func GetMonitorHeight(monitor int32) int32 {
	return getMonitorHeight(monitor)
}

// GetMonitorPhysicalWidth - Get specified monitor physical width in millimetres
func GetMonitorPhysicalWidth(monitor int32) int32 {
	return getMonitorPhysicalWidth(monitor)
}

// GetMonitorPhysicalHeight - Get specified monitor physical height in millimetres
func GetMonitorPhysicalHeight(monitor int32) int32 {
	return getMonitorPhysicalHeight(monitor)
}

// GetMonitorRefreshRate - Get specified monitor refresh rate
func GetMonitorRefreshRate(monitor int32) int32 {
	return getMonitorRefreshRate(monitor)
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
func GetMonitorName(monitor int32) string {
	return getMonitorName(monitor)
}

// SetClipboardText - Set clipboard text content
func SetClipboardText(text string) {
	setClipboardText(text)
}

// GetClipboardText - Get clipboard text content
func GetClipboardText() string {
	return getClipboardText()
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
	beginShaderMode(*(*uintptr)(unsafe.Pointer(&shader)))
}

// EndShaderMode - End custom shader drawing (use default shader)
func EndShaderMode() {
	endShaderMode()
}

// BeginBlendMode - Begin blending mode (alpha, additive, multiplied, subtract, custom)
func BeginBlendMode(mode int32) {
	beginBlendMode(mode)
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
func BeginVrStereoMode(config VrStereoConfig) {}

// EndVrStereoMode - End stereo rendering (requires VR simulator)
func EndVrStereoMode() {}

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
	loadShader(uintptr(unsafe.Pointer(&shader)), vsFileName, fsFileName)
	return shader
}

// LoadShaderFromMemory - Load shader from code strings and bind default locations
func LoadShaderFromMemory(vsCode string, fsCode string) Shader {
	var shader Shader
	loadShaderFromMemory(uintptr(unsafe.Pointer(&shader)), vsCode, fsCode)
	return shader
}

// IsShaderReady - Check if a shader is ready
func IsShaderReady(shader Shader) bool {
	return isShaderReady(uintptr(unsafe.Pointer(&shader)))
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
func SetShaderValue(shader Shader, locIndex int32, value []float32, uniformType int32) {
	setShaderValue(uintptr(unsafe.Pointer(&shader)), locIndex, value, uniformType)
}

// SetShaderValueV - Set shader uniform value vector
func SetShaderValueV(shader Shader, locIndex int32, value []float32, uniformType int32, count int32) {
	setShaderValueV(uintptr(unsafe.Pointer(&shader)), locIndex, value, uniformType, count)
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
func GetMouseRay(mousePosition Vector2, camera Camera) Ray {
	var ray Ray
	getMouseRay(uintptr(unsafe.Pointer(&ray)), *(*uintptr)(unsafe.Pointer(&mousePosition)), uintptr(unsafe.Pointer(&camera)))
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
func TraceLog(logLevel int32, text string, args ...any) {
	traceLog(logLevel, text, args...)
}

// SetTraceLogLevel - Set the current threshold (minimum) log level
func SetTraceLogLevel(logLevel int32) {
	setTraceLogLevel(logLevel)
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
//
// REVIEW NEEDED! 2023-11-15 JupiterRider: The argument list paramter isn't impelmented yet.
func SetTraceLogCallback(fn TraceLogCallbackFun) {
	setTraceLogCallback(traceLogCallbackWrapper(fn))
}

// // SetLoadFileDataCallback - Set custom file binary data loader
// func SetLoadFileDataCallback(callback LoadFileDataCallback) {}

// // SetSaveFileDataCallback - Set custom file binary data saver
// func SetSaveFileDataCallback(callback SaveFileDataCallback) {}

// // SetLoadFileTextCallback - Set custom file text data loader
// func SetLoadFileTextCallback(callback LoadFileTextCallback) {}

// // SetSaveFileTextCallback - Set custom file text data saver
// func SetSaveFileTextCallback(callback SaveFileTextCallback) {}

// LoadFileData - Load file data as byte array (read)
//
// Note: Because the slice is allocated in C, use UnloadFileData when it isn't needed anymore.
func LoadFileData(fileName string, dataSize *int32) []byte {
	ret := loadFileData(fileName, dataSize)
	return unsafe.Slice(ret, int(*dataSize))
}

// UnloadFileData - Unload file data allocated by LoadFileData()
func UnloadFileData(data []byte) {
	unloadFileData(unsafe.SliceData(data))
}

// SaveFileData - Save data to file from byte array (write), returns true on success
//
// Note: As an alternative, you could use go's os.WriteFile function.
func SaveFileData(fileName string, data []byte, dataSize int32) bool {
	return saveFileData(fileName, unsafe.Pointer(unsafe.SliceData(data)), dataSize)
}

// ExportDataAsCode - Export data to code (.h), returns true on success
func ExportDataAsCode(data []byte, dataSize int32, fileName string) bool {
	return exportDataAsCode(unsafe.SliceData(data), dataSize, fileName)
}

// // LoadFileText - Load text data from file (read), returns a '\0' terminated string
// func LoadFileText(fileName string) string {}

// // UnloadFileText - Unload file text data allocated by LoadFileText()
// func UnloadFileText(text string) {}

// // SaveFileText - Save text data to file (write), string must be '\0' terminated, returns true on success
// func SaveFileText(fileName string, text string) bool {}

// // FileExists - Check if file exists
// func FileExists(fileName string) bool {}

// // DirectoryExists - Check if a directory path exists
// func DirectoryExists(dirPath string) bool {}

// // IsFileExtension - Check file extension (including point: .png, .wav)
// func IsFileExtension(fileName string, ext string) bool {}

// // GetFileLength - Get file length in bytes (NOTE: GetFileSize() conflicts with windows.h)
// func GetFileLength(fileName string) int32 {}

// // GetFileExtension - Get pointer to extension for a filename string (includes dot: '.png')
// func GetFileExtension(fileName string) string {}

// // GetFileName - Get pointer to filename for a path string
// func GetFileName(filePath string) string {}

// // GetFileNameWithoutExt - Get filename string without extension (uses static string)
// func GetFileNameWithoutExt(filePath string) string {}

// // GetDirectoryPath - Get full path for a given fileName with path (uses static string)
// func GetDirectoryPath(filePath string) string {}

// // GetPrevDirectoryPath - Get previous directory path for a given path (uses static string)
// func GetPrevDirectoryPath(dirPath string) string {}

// GetWorkingDirectory - Get current working directory (uses static string)
func GetWorkingDirectory() string {
	return getWorkingDirectory()
}

// GetApplicationDirectory - Get the directory of the running application (uses static string)
func GetApplicationDirectory() string {
	return getApplicationDirectory()
}

// ChangeDirectory - Change working directory, return true on success
func ChangeDirectory(dir string) bool {
	return changeDirectory(dir)
}

// // IsPathFile - Check if a given path is a file or a directory
// func IsPathFile(path string) bool {}

// // LoadDirectoryFiles - Load directory filepaths
// func LoadDirectoryFiles(dirPath string) FilePathList {}

// // LoadDirectoryFilesEx - Load directory filepaths with extension filtering and recursive directory scan
// func LoadDirectoryFilesEx(basePath string, filter string, scanSubdirs bool) FilePathList {}

// // UnloadDirectoryFiles - Unload filepaths
// func UnloadDirectoryFiles(files FilePathList) {}

// IsFileDropped - Check if a file has been dropped into window
//
// REVIEW NEEDED! 2023-11-12 JupiterRider: This funtions always returns true.
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

// // UnloadDroppedFiles - Unload dropped filepaths
// func UnloadDroppedFiles(files FilePathList) {}

// // GetFileModTime - Get file modification time (last write time)
// func GetFileModTime(fileName string) int64 {}

// // CompressData - Compress data (DEFLATE algorithm), memory must be MemFree()
// func CompressData(data []byte, dataSize int32, compDataSize []int32) []byte {}

// // DecompressData - Decompress data (DEFLATE algorithm), memory must be MemFree()
// func DecompressData(compData []byte, compDataSize int32, dataSize []int32) []byte {}

// // EncodeDataBase64 - Encode data to Base64 string, memory must be MemFree()
// func EncodeDataBase64(data []byte, dataSize int32, outputSize []int32) string {}

// // DecodeDataBase64 - Decode Base64 string data, memory must be MemFree()
// func DecodeDataBase64(data []byte, outputSize []int32) []byte {}

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
func SetAutomationEventBaseFrame(frame int32) {
	setAutomationEventBaseFrame(frame)
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

// IsMouseButtonPressed - Check if a mouse button has been pressed once
func IsMouseButtonPressed(button int32) bool {
	return isMouseButtonPressed(button)
}

// IsMouseButtonDown - Check if a mouse button is being pressed
func IsMouseButtonDown(button int32) bool {
	return isMouseButtonDown(button)
}

// IsMouseButtonReleased - Check if a mouse button has been released once
func IsMouseButtonReleased(button int32) bool {
	return isMouseButtonReleased(button)
}

// IsMouseButtonUp - Check if a mouse button is NOT being pressed
func IsMouseButtonUp(button int32) bool {
	return isMouseButtonUp(button)
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
func IsGestureDetected(gesture uint32) bool {
	return isGestureDetected(gesture)
}

// GetGestureDetected - Get latest detected gesture
func GetGestureDetected() int32 {
	return getGestureDetected()
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

// // UpdateCamera - Update camera position for selected mode
// func UpdateCamera(camera *Camera, mode int32) {}

// // UpdateCameraPro - Update camera movement/rotation
// func UpdateCameraPro(camera *Camera, movement Vector3, rotation Vector3, zoom float32) {}

// SetShapesTexture - Set texture and rectangle to be used on shapes drawing
func SetShapesTexture(texture Texture2D, source Rectangle) {
	setShapesTexture(uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(&source)))
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
func DrawCircleGradient(centerX int32, centerY int32, radius float32, color1 color.RGBA, color2 color.RGBA) {
	drawCircleGradient(centerX, centerY, radius, *(*uintptr)(unsafe.Pointer(&color1)), *(*uintptr)(unsafe.Pointer(&color2)))
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
func DrawRectangleGradientV(posX int32, posY int32, width int32, height int32, color1 color.RGBA, color2 color.RGBA) {
	drawRectangleGradientV(posX, posY, width, height, *(*uintptr)(unsafe.Pointer(&color1)), *(*uintptr)(unsafe.Pointer(&color2)))
}

// DrawRectangleGradientH - Draw a horizontal-gradient-filled rectangle
func DrawRectangleGradientH(posX int32, posY int32, width int32, height int32, color1 color.RGBA, color2 color.RGBA) {
	drawRectangleGradientH(posX, posY, width, height, *(*uintptr)(unsafe.Pointer(&color1)), *(*uintptr)(unsafe.Pointer(&color2)))
}

// DrawRectangleGradientEx - Draw a gradient-filled rectangle with custom vertex colors
func DrawRectangleGradientEx(rec Rectangle, col1 color.RGBA, col2 color.RGBA, col3 color.RGBA, col4 color.RGBA) {
	drawRectangleGradientEx(uintptr(unsafe.Pointer(&rec)), *(*uintptr)(unsafe.Pointer(&col1)), *(*uintptr)(unsafe.Pointer(&col2)), *(*uintptr)(unsafe.Pointer(&col3)), *(*uintptr)(unsafe.Pointer(&col4)))
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

// DrawRectangleRoundedLines - Draw rectangle with rounded edges outline
func DrawRectangleRoundedLines(rec Rectangle, roundness float32, segments int32, lineThick float32, col color.RGBA) {
	drawRectangleRoundedLines(uintptr(unsafe.Pointer(&rec)), roundness, segments, lineThick, *(*uintptr)(unsafe.Pointer(&col)))
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
func LoadImage(fileName string) Image {
	var img Image
	loadImage(uintptr(unsafe.Pointer(&img)), fileName)
	return img
}

// LoadImageRaw - Load image from RAW file data
func LoadImageRaw(fileName string, width int32, height int32, format int32, headerSize int32) Image {
	var img Image
	loadImageRaw(uintptr(unsafe.Pointer(&img)), fileName, width, height, format, headerSize)
	return img
}

// LoadImageSvg - Load image from SVG file data or string with specified size
func LoadImageSvg(fileNameOrString string, width int32, height int32) Image {
	var img Image
	loadImageSvg(uintptr(unsafe.Pointer(&img)), fileNameOrString, width, height)
	return img
}

// LoadImageAnim - Load image sequence from file (frames appended to image.data)
func LoadImageAnim(fileName string, frames []int32) Image {
	var img Image
	loadImageAnim(uintptr(unsafe.Pointer(&img)), fileName, frames)
	return img
}

// LoadImageFromMemory - Load image from memory buffer, fileType refers to extension: i.e. '.png'
func LoadImageFromMemory(fileType string, fileData []byte, dataSize int32) Image {
	var img Image
	loadImageFromMemory(uintptr(unsafe.Pointer(&img)), fileType, fileData, dataSize)
	return img
}

// LoadImageFromTexture - Load image from GPU texture data
func LoadImageFromTexture(texture Texture2D) Image {
	var img Image
	loadImageFromTexture(uintptr(unsafe.Pointer(&img)), uintptr(unsafe.Pointer(&texture)))
	return img
}

// LoadImageFromScreen - Load image from screen buffer and (screenshot)
func LoadImageFromScreen() Image {
	var img Image
	loadImageFromScreen(uintptr(unsafe.Pointer(&img)))
	return img
}

// IsImageReady - Check if an image is ready
func IsImageReady(image Image) bool {
	return isImageReady(uintptr(unsafe.Pointer(&image)))
}

// UnloadImage - Unload image from CPU memory (RAM)
func UnloadImage(image Image) {
	unloadImage(uintptr(unsafe.Pointer(&image)))
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

// ExportImageAsCode - Export image as code file defining an array of bytes, returns true on success
func ExportImageAsCode(image Image, fileName string) bool {
	return exportImageAsCode(uintptr(unsafe.Pointer(&image)), fileName)
}

// GenImageColor - Generate image: plain color
func GenImageColor(width int32, height int32, col color.RGBA) Image {
	var image Image
	genImageColor(uintptr(unsafe.Pointer(&image)), width, height, *(*uintptr)(unsafe.Pointer(&col)))
	return image
}

// GenImageGradientLinear - Generate image: linear gradient, direction in degrees [0..360], 0=Vertical gradient
func GenImageGradientLinear(width int32, height int32, direction int32, start color.RGBA, end color.RGBA) Image {
	var image Image
	genImageGradientLinear(uintptr(unsafe.Pointer(&image)), width, height, direction, *(*uintptr)(unsafe.Pointer(&start)), *(*uintptr)(unsafe.Pointer(&end)))
	return image
}

// GenImageGradientRadial - Generate image: radial gradient
func GenImageGradientRadial(width int32, height int32, density float32, inner color.RGBA, outer color.RGBA) Image {
	var image Image
	genImageGradientRadial(uintptr(unsafe.Pointer(&image)), width, height, density, *(*uintptr)(unsafe.Pointer(&inner)), *(*uintptr)(unsafe.Pointer(&outer)))
	return image
}

// GenImageGradientSquare - Generate image: square gradient
func GenImageGradientSquare(width int32, height int32, density float32, inner color.RGBA, outer color.RGBA) Image {
	var image Image
	genImageGradientSquare(uintptr(unsafe.Pointer(&image)), width, height, density, *(*uintptr)(unsafe.Pointer(&inner)), *(*uintptr)(unsafe.Pointer(&outer)))
	return image
}

// GenImageChecked - Generate image: checked
func GenImageChecked(width int32, height int32, checksX int32, checksY int32, col1 color.RGBA, col2 color.RGBA) Image {
	var image Image
	genImageChecked(uintptr(unsafe.Pointer(&image)), width, height, checksX, checksY, *(*uintptr)(unsafe.Pointer(&col1)), *(*uintptr)(unsafe.Pointer(&col2)))
	return image
}

// GenImageWhiteNoise - Generate image: white noise
func GenImageWhiteNoise(width int32, height int32, factor float32) Image {
	var image Image
	genImageWhiteNoise(uintptr(unsafe.Pointer(&image)), width, height, factor)
	return image
}

// GenImagePerlinNoise - Generate image: perlin noise
func GenImagePerlinNoise(width int32, height int32, offsetX int32, offsetY int32, scale float32) Image {
	var image Image
	genImagePerlinNoise(uintptr(unsafe.Pointer(&image)), width, height, offsetX, offsetY, scale)
	return image
}

// GenImageCellular - Generate image: cellular algorithm, bigger tileSize means bigger cells
func GenImageCellular(width int32, height int32, tileSize int32) Image {
	var image Image
	genImageCellular(uintptr(unsafe.Pointer(&image)), width, height, tileSize)
	return image
}

// GenImageText - Generate image: grayscale image from text data
func GenImageText(width int32, height int32, text string) Image {
	var image Image
	genImageText(uintptr(unsafe.Pointer(&image)), width, height, text)
	return image
}

// ImageCopy - Create an image duplicate (useful for transformations)
func ImageCopy(image Image) Image {
	var retImage Image
	imageCopy(uintptr(unsafe.Pointer(&retImage)), uintptr(unsafe.Pointer(&image)))
	return retImage
}

// ImageFromImage - Create an image from another image piece
func ImageFromImage(image Image, rec Rectangle) Image {
	var retImage Image
	imageFromImage(uintptr(unsafe.Pointer(&retImage)), uintptr(unsafe.Pointer(&image)), uintptr(unsafe.Pointer(&rec)))
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
func ImageFormat(image *Image, newFormat int32) {
	imageFormat(image, newFormat)
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
func ImageAlphaMask(image *Image, alphaMask Image) {
	imageAlphaMask(image, uintptr(unsafe.Pointer(&alphaMask)))
}

// ImageAlphaPremultiply - Premultiply alpha channel
func ImageAlphaPremultiply(image *Image) {
	imageAlphaPremultiply(image)
}

// ImageBlurGaussian - Apply Gaussian blur using a box blur approximation
func ImageBlurGaussian(image *Image, blurSize int32) {
	imageBlurGaussian(image, blurSize)
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
func LoadImageColors(image Image) []color.RGBA {
	ret := loadImageColors(uintptr(unsafe.Pointer(&image)))
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
func ImageDrawLineV(dst *Image, start Vector2, end Vector2, col color.RGBA) {
	imageDrawLineV(dst, *(*uintptr)(unsafe.Pointer(&start)), *(*uintptr)(unsafe.Pointer(&end)), *(*uintptr)(unsafe.Pointer(&col)))
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
func ImageDrawRectangleLines(dst *Image, rec Rectangle, thick int32, col color.RGBA) {
	imageDrawRectangleLines(dst, uintptr(unsafe.Pointer(&rec)), thick, *(*uintptr)(unsafe.Pointer(&col)))
}

// ImageDraw - Draw a source image within a destination image (tint applied to source)
func ImageDraw(dst *Image, src Image, srcRec Rectangle, dstRec Rectangle, tint color.RGBA) {
	imageDraw(dst, uintptr(unsafe.Pointer(&src)), uintptr(unsafe.Pointer(&srcRec)), uintptr(unsafe.Pointer(&dstRec)), *(*uintptr)(unsafe.Pointer(&tint)))
}

// ImageDrawText - Draw text (using default font) within an image (destination)
func ImageDrawText(dst *Image, text string, posX int32, posY int32, fontSize int32, col color.RGBA) {
	imageDrawText(dst, text, posX, posY, fontSize, *(*uintptr)(unsafe.Pointer(&col)))

}

// ImageDrawTextEx - Draw text (custom sprite font) within an image (destination)
func ImageDrawTextEx(dst *Image, font Font, text string, position Vector2, fontSize float32, spacing float32, tint color.RGBA) {
	imageDrawTextEx(dst, uintptr(unsafe.Pointer(&font)), text, *(*uintptr)(unsafe.Pointer(&position)), fontSize, spacing, *(*uintptr)(unsafe.Pointer(&tint)))
}

// LoadTexture - Load texture from file into GPU memory (VRAM)
func LoadTexture(fileName string) Texture2D {
	var texture Texture2D
	loadTexture(uintptr(unsafe.Pointer(&texture)), fileName)
	return texture
}

// LoadTextureFromImage - Load texture from image data
func LoadTextureFromImage(image Image) Texture2D {
	var texture Texture2D
	loadTextureFromImage(uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(&image)))
	return texture
}

// LoadTextureCubemap - Load cubemap from image, multiple image cubemap layouts supported
func LoadTextureCubemap(image Image, layout int32) Texture2D {
	var texture Texture2D
	loadTextureCubemap(uintptr(unsafe.Pointer(&texture)), uintptr(unsafe.Pointer(&image)), layout)
	return texture
}

// LoadRenderTexture - Load texture for rendering (framebuffer)
func LoadRenderTexture(width int32, height int32) RenderTexture2D {
	var texture RenderTexture2D
	loadRenderTexture(uintptr(unsafe.Pointer(&texture)), width, height)
	return texture
}

// IsTextureReady - Check if a texture is ready
func IsTextureReady(texture Texture2D) bool {
	return isTextureReady(uintptr(unsafe.Pointer(&texture)))
}

// UnloadTexture - Unload texture from GPU memory (VRAM)
func UnloadTexture(texture Texture2D) {
	unloadTexture(uintptr(unsafe.Pointer(&texture)))
}

// IsRenderTextureReady - Check if a render texture is ready
func IsRenderTextureReady(target RenderTexture2D) bool {
	return isRenderTextureReady(uintptr(unsafe.Pointer(&target)))
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
func SetTextureFilter(texture Texture2D, filter int32) {
	setTextureFilter(uintptr(unsafe.Pointer(&texture)), filter)
}

// SetTextureWrap - Set texture wrapping mode
func SetTextureWrap(texture Texture2D, wrap int32) {
	setTextureWrap(uintptr(unsafe.Pointer(&texture)), wrap)
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

// ColorToInt - Get hexadecimal value for a Color
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
	ret := colorFromNormalized(*(*uintptr)(unsafe.Pointer(&normalized)))
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

// GetColor - Get Color structure from hexadecimal value
func GetColor(hexValue uint32) color.RGBA {
	ret := getColor(hexValue)
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
func LoadFontEx(fileName string, fontSize int32, codepoints []rune) Font {
	var font Font
	codepointCount := int32(len(codepoints))
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

// IsFontReady - Check if a font is ready
func IsFontReady(font Font) bool {
	return isFontReady(uintptr(unsafe.Pointer(&font)))
}

// LoadFontData - Load font data for further use
func LoadFontData(fileData []byte, fontSize int32, codepoints []rune, typ int32) []GlyphInfo {
	dataSize := int32(len(fileData))
	codepointCount := int32(len(codepoints))
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

// ExportFontAsCode - Export font as code file, returns true on success
func ExportFontAsCode(font Font, fileName string) bool {
	return exportFontAsCode(uintptr(unsafe.Pointer(&font)), fileName)
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
func SetTextLineSpacing(spacing int32) {
	setTextLineSpacing(spacing)
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

// // LoadUTF8 - Load UTF-8 text encoded from codepoints array
// func LoadUTF8(codepoints []int32, length int32) string {
// 	return ""
// }

// // UnloadUTF8 - Unload UTF-8 text encoded from codepoints array
// func UnloadUTF8(text string) {}

// // LoadCodepoints - Load all codepoints from a UTF-8 text string, codepoints count returned by parameter
// func LoadCodepoints(text string, count []int32) []int32 {
// 	return nil
// }

// // UnloadCodepoints - Unload codepoints data from memory
// func UnloadCodepoints(codepoints []int32) {}

// // GetCodepointCount - Get total number of codepoints in a UTF-8 encoded string
// func GetCodepointCount(text string) int32 {
// 	return 0
// }

// // GetCodepoint - Get next codepoint in a UTF-8 encoded string, 0x3f('?') is returned on failure
// func GetCodepoint(text string, codepointSize []int32) int32 {
// 	return 0
// }

// // GetCodepointNext - Get next codepoint in a UTF-8 encoded string, 0x3f('?') is returned on failure
// func GetCodepointNext(text string, codepointSize []int32) int32 {
// 	return 0
// }

// // GetCodepointPrevious - Get previous codepoint in a UTF-8 encoded string, 0x3f('?') is returned on failure
// func GetCodepointPrevious(text string, codepointSize []int32) int32 {
// 	return 0
// }

// // CodepointToUTF8 - Encode one codepoint into UTF-8 byte array (array length returned as parameter)
// func CodepointToUTF8(codepoint int32, utf8Size []int32) string {
// 	return ""
// }

// // TextCopy - Copy one string to another, returns bytes copied
// func TextCopy(dst string, src string) int32 {
// 	return 0
// }

// // TextIsEqual - Check if two text string are equal
// func TextIsEqual(text1 string, text2 string) bool {
// 	return false
// }

// // TextLength - Get text length, checks for '\0' ending
// func TextLength(text string) uint32 {
// 	return 0
// }

// // TextFormat - Text formatting with variables (sprintf() style)
// func TextFormat(text string, args ...any) string {
// 	return ""
// }

// // TextSubtext - Get a piece of a text string
// func TextSubtext(text string, position int32, length int32) string {
// 	return ""
// }

// // TextReplace - Replace text string (WARNING: memory must be freed!)
// func TextReplace(text string, replace string, by string) string {
// 	return ""
// }

// // TextInsert - Insert text in a position (WARNING: memory must be freed!)
// func TextInsert(text string, insert string, position int32) string {
// 	return ""
// }

// // TextJoin - Join text strings with delimiter
// func TextJoin(textList []string, count int32, delimiter string) string {
// 	return ""
// }

// // TextSplit - Split text into multiple strings
// func TextSplit(text string, delimiter int8, count []int32) []string {
// 	return nil
// }

// // TextAppend - Append text at specific position and move cursor!
// func TextAppend(text string, _append string, position []int32) {}

// // TextFindIndex - Find first text occurrence within a string
// func TextFindIndex(text string, find string) int32 {
// 	return 0
// }

// // TextToUpper - Get upper case version of provided string
// func TextToUpper(text string) string {
// 	return ""
// }

// // TextToLower - Get lower case version of provided string
// func TextToLower(text string) string {
// 	return ""
// }

// // TextToPascal - Get Pascal case notation version of provided string
// func TextToPascal(text string) string {
// 	return ""
// }

// // TextToInteger - Get integer value from text (negative values not supported)
// func TextToInteger(text string) int32 {
// 	return 0
// }

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
func DrawSphereEx(centerPos Vector3, radius float32, rings int32, slices int32, col color.RGBA) {}

// DrawSphereWires - Draw sphere wires
func DrawSphereWires(centerPos Vector3, radius float32, rings int32, slices int32, col color.RGBA) {
}

// DrawCylinder - Draw a cylinder/cone
func DrawCylinder(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, col color.RGBA) {
}

// DrawCylinderEx - Draw a cylinder with base at startPos and top at endPos
func DrawCylinderEx(startPos Vector3, endPos Vector3, startRadius float32, endRadius float32, sides int32, col color.RGBA) {
}

// DrawCylinderWires - Draw a cylinder/cone wires
func DrawCylinderWires(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, col color.RGBA) {
}

// DrawCylinderWiresEx - Draw a cylinder wires with base at startPos and top at endPos
func DrawCylinderWiresEx(startPos Vector3, endPos Vector3, startRadius float32, endRadius float32, sides int32, col color.RGBA) {
}

// DrawCapsule - Draw a capsule with the center of its sphere caps at startPos and endPos
func DrawCapsule(startPos Vector3, endPos Vector3, radius float32, slices int32, rings int32, col color.RGBA) {
}

// DrawCapsuleWires - Draw capsule wireframe with the center of its sphere caps at startPos and endPos
func DrawCapsuleWires(startPos Vector3, endPos Vector3, radius float32, slices int32, rings int32, col color.RGBA) {
}

// DrawPlane - Draw a plane XZ
func DrawPlane(centerPos Vector3, size Vector2, col color.RGBA) {}

// DrawRay - Draw a ray line
func DrawRay(ray Ray, col color.RGBA) {}

// DrawGrid - Draw a grid (centered at (0, 0, 0))
func DrawGrid(slices int32, spacing float32) {}

// LoadModel - Load model from files (meshes and materials)
func LoadModel(fileName string) Model {
	return Model{}
}

// LoadModelFromMesh - Load model from generated mesh (default material)
func LoadModelFromMesh(mesh Mesh) Model {
	return Model{}
}

// IsModelReady - Check if a model is ready
func IsModelReady(model Model) bool {
	return false
}

// UnloadModel - Unload model (including meshes) from memory (RAM and/or VRAM)
func UnloadModel(model Model) {}

// GetModelBoundingBox - Compute model bounding box limits (considers all meshes)
func GetModelBoundingBox(model Model) BoundingBox {
	return BoundingBox{}
}

// DrawModel - Draw a model (with texture if set)
func DrawModel(model Model, position Vector3, scale float32, tint color.RGBA) {}

// DrawModelEx - Draw a model with extended parameters
func DrawModelEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
}

// DrawModelWires - Draw a model wires (with texture if set)
func DrawModelWires(model Model, position Vector3, scale float32, tint color.RGBA) {}

// DrawModelWiresEx - Draw a model wires (with texture if set) with extended parameters
func DrawModelWiresEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
}

// DrawBoundingBox - Draw bounding box (wires)
func DrawBoundingBox(box BoundingBox, col color.RGBA) {}

// DrawBillboard - Draw a billboard texture
func DrawBillboard(camera Camera, texture Texture2D, position Vector3, size float32, tint color.RGBA) {
}

// DrawBillboardRec - Draw a billboard texture defined by source
func DrawBillboardRec(camera Camera, texture Texture2D, source Rectangle, position Vector3, size Vector2, tint color.RGBA) {
}

// DrawBillboardPro - Draw a billboard texture defined by source and rotation
func DrawBillboardPro(camera Camera, texture Texture2D, source Rectangle, position Vector3, up Vector3, size Vector2, origin Vector2, rotation float32, tint color.RGBA) {
}

// UploadMesh - Upload mesh vertex data in GPU and provide VAO/VBO ids
func UploadMesh(mesh *Mesh, dynamic bool) {}

// UpdateMeshBuffer - Update mesh vertex data in GPU for a specific buffer index
func UpdateMeshBuffer(mesh Mesh, index int32, data unsafe.Pointer, dataSize int32, offset int32) {}

// UnloadMesh - Unload mesh data from CPU and GPU
func UnloadMesh(mesh Mesh) {}

// DrawMesh - Draw a 3d mesh with material and transform
func DrawMesh(mesh Mesh, material Material, transform Matrix) {}

// DrawMeshInstanced - Draw multiple mesh instances with material and different transforms
func DrawMeshInstanced(mesh Mesh, material Material, transforms *Matrix, instances int32) {}

// ExportMesh - Export mesh data to file, returns true on success
func ExportMesh(mesh Mesh, fileName string) bool {
	return false
}

// GetMeshBoundingBox - Compute mesh bounding box limits
func GetMeshBoundingBox(mesh Mesh) BoundingBox {
	return BoundingBox{}
}

// GenMeshTangents - Compute mesh tangents
func GenMeshTangents(mesh *Mesh) {}

// GenMeshPoly - Generate polygonal mesh
func GenMeshPoly(sides int32, radius float32) Mesh {
	return Mesh{}
}

// GenMeshPlane - Generate plane mesh (with subdivisions)
func GenMeshPlane(width float32, length float32, resX int32, resZ int32) Mesh {
	return Mesh{}
}

// GenMeshCube - Generate cuboid mesh
func GenMeshCube(width float32, height float32, length float32) Mesh {
	return Mesh{}
}

// GenMeshSphere - Generate sphere mesh (standard sphere)
func GenMeshSphere(radius float32, rings int32, slices int32) Mesh {
	return Mesh{}
}

// GenMeshHemiSphere - Generate half-sphere mesh (no bottom cap)
func GenMeshHemiSphere(radius float32, rings int32, slices int32) Mesh {
	return Mesh{}
}

// GenMeshCylinder - Generate cylinder mesh
func GenMeshCylinder(radius float32, height float32, slices int32) Mesh {
	return Mesh{}
}

// GenMeshCone - Generate cone/pyramid mesh
func GenMeshCone(radius float32, height float32, slices int32) Mesh {
	return Mesh{}
}

// GenMeshTorus - Generate torus mesh
func GenMeshTorus(radius float32, size float32, radSeg int32, sides int32) Mesh {
	return Mesh{}
}

// GenMeshKnot - Generate trefoil knot mesh
func GenMeshKnot(radius float32, size float32, radSeg int32, sides int32) Mesh {
	return Mesh{}
}

// GenMeshHeightmap - Generate heightmap mesh from image data
func GenMeshHeightmap(heightmap Image, size Vector3) Mesh {
	return Mesh{}
}

// GenMeshCubicmap - Generate cubes-based map mesh from image data
func GenMeshCubicmap(cubicmap Image, cubeSize Vector3) Mesh {
	return Mesh{}
}

// LoadMaterials - Load materials from model file
func LoadMaterials(fileName string, materialCount []int32) *Material {
	return &Material{}
}

// LoadMaterialDefault - Load default material (Supports: DIFFUSE, SPECULAR, NORMAL maps)
func LoadMaterialDefault() Material {
	return Material{}
}

// IsMaterialReady - Check if a material is ready
func IsMaterialReady(material Material) bool {
	return false
}

// UnloadMaterial - Unload material from GPU memory (VRAM)
func UnloadMaterial(material Material) {}

// SetMaterialTexture - Set texture for a material map type (MATERIAL_MAP_DIFFUSE, MATERIAL_MAP_SPECULAR...)
func SetMaterialTexture(material *Material, mapType int32, texture Texture2D) {}

// SetModelMeshMaterial - Set material for a mesh
func SetModelMeshMaterial(model *Model, meshId int32, materialId int32) {}

// LoadModelAnimations - Load model animations from file
func LoadModelAnimations(fileName string, animCount []int32) *ModelAnimation {
	return &ModelAnimation{}
}

// UpdateModelAnimation - Update model animation pose
func UpdateModelAnimation(model Model, anim ModelAnimation, frame int32) {}

// UnloadModelAnimation - Unload animation data
func UnloadModelAnimation(anim ModelAnimation) {}

// UnloadModelAnimations - Unload animation array data
func UnloadModelAnimations(animations *ModelAnimation, animCount int32) {}

// IsModelAnimationValid - Check model animation skeleton match
func IsModelAnimationValid(model Model, anim ModelAnimation) bool {
	return false
}

// CheckCollisionSpheres - Check collision between two spheres
func CheckCollisionSpheres(center1 Vector3, radius1 float32, center2 Vector3, radius2 float32) bool {
	return false
}

// CheckCollisionBoxes - Check collision between two bounding boxes
func CheckCollisionBoxes(box1 BoundingBox, box2 BoundingBox) bool {
	return false
}

// CheckCollisionBoxSphere - Check collision between box and sphere
func CheckCollisionBoxSphere(box BoundingBox, center Vector3, radius float32) bool {
	return false
}

// GetRayCollisionSphere - Get collision info between ray and sphere
func GetRayCollisionSphere(ray Ray, center Vector3, radius float32) RayCollision {
	return RayCollision{}
}

// GetRayCollisionBox - Get collision info between ray and box
func GetRayCollisionBox(ray Ray, box BoundingBox) RayCollision {
	return RayCollision{}
}

// GetRayCollisionMesh - Get collision info between ray and mesh
func GetRayCollisionMesh(ray Ray, mesh Mesh, transform Matrix) RayCollision {
	return RayCollision{}
}

// GetRayCollisionTriangle - Get collision info between ray and triangle
func GetRayCollisionTriangle(ray Ray, p1 Vector3, p2 Vector3, p3 Vector3) RayCollision {
	return RayCollision{}
}

// GetRayCollisionQuad - Get collision info between ray and quad
func GetRayCollisionQuad(ray Ray, p1 Vector3, p2 Vector3, p3 Vector3, p4 Vector3) RayCollision {
	return RayCollision{}
}

// InitAudioDevice - Initialize audio device and context
func InitAudioDevice() {}

// CloseAudioDevice - Close the audio device and context
func CloseAudioDevice() {}

// IsAudioDeviceReady - Check if audio device has been initialized successfully
func IsAudioDeviceReady() bool {
	return true
}

// SetMasterVolume - Set master volume (listener)
func SetMasterVolume(volume float32) {}

// GetMasterVolume - Get master volume (listener)
func GetMasterVolume() float32 {
	return 0
}

// LoadWave - Load wave data from file
func LoadWave(fileName string) Wave {
	return Wave{}
}

// LoadWaveFromMemory - Load wave from memory buffer, fileType refers to extension: i.e. '.wav'
func LoadWaveFromMemory(fileType string, fileData []byte, dataSize int32) Wave {
	return Wave{}
}

// IsWaveReady - Checks if wave data is ready
func IsWaveReady(wave Wave) bool {
	return true
}

// // LoadSound - Load sound from file
// func LoadSound(fileName string) Sound {
// 	return Sound{}
// }

// // LoadSoundFromWave - Load sound from wave data
// func LoadSoundFromWave(wave Wave) Sound {}

// // LoadSoundAlias - Create a new sound that shares the same sample data as the source sound, does not own the sound data
// func LoadSoundAlias(source Sound) Sound {}

// // IsSoundReady - Checks if a sound is ready
// func IsSoundReady(sound Sound) bool {}

// // UpdateSound - Update sound buffer with new data
// func UpdateSound(sound Sound, data unsafe.Pointer, sampleCount int32) {}

// // UnloadWave - Unload wave data
// func UnloadWave(wave Wave) {}

// // UnloadSound - Unload sound
// func UnloadSound(sound Sound) {}

// // UnloadSoundAlias - Unload a sound alias (does not deallocate sample data)
// func UnloadSoundAlias(alias Sound) {}

// // ExportWave - Export wave data to file, returns true on success
// func ExportWave(wave Wave, fileName string) bool {}

// // ExportWaveAsCode - Export wave sample data to code (.h), returns true on success
// func ExportWaveAsCode(wave Wave, fileName string) bool {}

// // PlaySound - Play a sound
// func PlaySound(sound Sound) {}

// // StopSound - Stop playing a sound
// func StopSound(sound Sound) {}

// // PauseSound - Pause a sound
// func PauseSound(sound Sound) {}

// // ResumeSound - Resume a paused sound
// func ResumeSound(sound Sound) {}

// // IsSoundPlaying - Check if a sound is currently playing
// func IsSoundPlaying(sound Sound) bool {}

// // SetSoundVolume - Set volume for a sound (1.0 is max level)
// func SetSoundVolume(sound Sound, volume float32) {}

// // SetSoundPitch - Set pitch for a sound (1.0 is base level)
// func SetSoundPitch(sound Sound, pitch float32) {}

// // SetSoundPan - Set pan for a sound (0.5 is center)
// func SetSoundPan(sound Sound, pan float32) {}

// // WaveCopy - Copy a wave to a new wave
// func WaveCopy(wave Wave) Wave {}

// // WaveCrop - Crop a wave to defined samples range
// func WaveCrop(wave *Wave, initSample int32, finalSample int32) {}

// // WaveFormat - Convert wave data to desired format
// func WaveFormat(wave *Wave, sampleRate int32, sampleSize int32, channels int32) {}

// // LoadWaveSamples - Load samples data from wave as a 32bit float data array
// func LoadWaveSamples(wave Wave) []float32 {}

// // UnloadWaveSamples - Unload samples data loaded with LoadWaveSamples()
// func UnloadWaveSamples(samples []float32) {}

// // LoadMusicStream - Load music stream from file
// func LoadMusicStream(fileName string) Music {}

// // LoadMusicStreamFromMemory - Load music stream from data
// func LoadMusicStreamFromMemory(fileType string, data []byte, dataSize int32) Music {}

// // IsMusicReady - Checks if a music stream is ready
// func IsMusicReady(music Music) bool {}

// // UnloadMusicStream - Unload music stream
// func UnloadMusicStream(music Music) {}

// // PlayMusicStream - Start music playing
// func PlayMusicStream(music Music) {}

// // IsMusicStreamPlaying - Check if music is playing
// func IsMusicStreamPlaying(music Music) bool {}

// // UpdateMusicStream - Updates buffers for music streaming
// func UpdateMusicStream(music Music) {}

// // StopMusicStream - Stop music playing
// func StopMusicStream(music Music) {}

// // PauseMusicStream - Pause music playing
// func PauseMusicStream(music Music) {}

// // ResumeMusicStream - Resume playing paused music
// func ResumeMusicStream(music Music) {}

// // SeekMusicStream - Seek music to a position (in seconds)
// func SeekMusicStream(music Music, position float32) {}

// // SetMusicVolume - Set volume for music (1.0 is max level)
// func SetMusicVolume(music Music, volume float32) {}

// // SetMusicPitch - Set pitch for a music (1.0 is base level)
// func SetMusicPitch(music Music, pitch float32) {}

// // SetMusicPan - Set pan for a music (0.5 is center)
// func SetMusicPan(music Music, pan float32) {}

// // GetMusicTimeLength - Get music time length (in seconds)
// func GetMusicTimeLength(music Music) float32 {}

// // GetMusicTimePlayed - Get current music time played (in seconds)
// func GetMusicTimePlayed(music Music) float32 {}

// // LoadAudioStream - Load audio stream (to stream raw audio pcm data)
// func LoadAudioStream(sampleRate uint32, sampleSize uint32, channels uint32) AudioStream {}

// // IsAudioStreamReady - Checks if an audio stream is ready
// func IsAudioStreamReady(stream AudioStream) bool {}

// // UnloadAudioStream - Unload audio stream and free memory
// func UnloadAudioStream(stream AudioStream) {}

// // UpdateAudioStream - Update audio stream buffers with data
// func UpdateAudioStream(stream AudioStream, data unsafe.Pointer, frameCount int32) {}

// // IsAudioStreamProcessed - Check if any audio stream buffers requires refill
// func IsAudioStreamProcessed(stream AudioStream) bool {}

// // PlayAudioStream - Play audio stream
// func PlayAudioStream(stream AudioStream) {}

// // PauseAudioStream - Pause audio stream
// func PauseAudioStream(stream AudioStream) {}

// // ResumeAudioStream - Resume audio stream
// func ResumeAudioStream(stream AudioStream) {}

// // IsAudioStreamPlaying - Check if audio stream is playing
// func IsAudioStreamPlaying(stream AudioStream) bool {}

// // StopAudioStream - Stop audio stream
// func StopAudioStream(stream AudioStream) {}

// // SetAudioStreamVolume - Set volume for audio stream (1.0 is max level)
// func SetAudioStreamVolume(stream AudioStream, volume float32) {}

// // SetAudioStreamPitch - Set pitch for audio stream (1.0 is base level)
// func SetAudioStreamPitch(stream AudioStream, pitch float32) {}

// // SetAudioStreamPan - Set pan for audio stream (0.5 is centered)
// func SetAudioStreamPan(stream AudioStream, pan float32) {}

// // SetAudioStreamBufferSizeDefault - Default size for new audio streams
// func SetAudioStreamBufferSizeDefault(size int32) {}

// // SetAudioStreamCallback - Audio thread callback to request new data
// func SetAudioStreamCallback(stream AudioStream, callback AudioCallback) {}

// // AttachAudioStreamProcessor - Attach audio stream processor to stream, receives the samples as <float>s
// func AttachAudioStreamProcessor(stream AudioStream, processor AudioCallback) {}

// // DetachAudioStreamProcessor - Detach audio stream processor from stream
// func DetachAudioStreamProcessor(stream AudioStream, processor AudioCallback) {}

// // AttachAudioMixedProcessor - Attach audio stream processor to the entire audio pipeline, receives the samples as <float>s
// func AttachAudioMixedProcessor(processor AudioCallback) {}

// // DetachAudioMixedProcessor - Detach audio stream processor from the entire audio pipeline
// func DetachAudioMixedProcessor(processor AudioCallback) {}
