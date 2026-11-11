package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"image"
	"image/color"
	"unsafe"
)

// newImageFromPointer - Returns new Image from pointer
func newImageFromPointer(ptr unsafe.Pointer) *Image {
	return (*Image)(ptr)
}

// cptr returns C pointer
func (i *Image) cptr() *C.Image {
	return (*C.Image)(unsafe.Pointer(i))
}

// ToImage converts a Image to Go image.Image
func (i *Image) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(i.Width), int(i.Height)))

	// Get pixel data from image (RGBA 32bit)
	cimg := i.cptr()
	ret := C.LoadImageColors(*cimg)
	pixels := (*[1 << 24]uint8)(unsafe.Pointer(ret))[0 : i.Width*i.Height*4]

	img.Pix = pixels

	return img
}

// newTexture2DFromPointer - Returns new Texture2D from pointer
func newTexture2DFromPointer(ptr unsafe.Pointer) Texture2D {
	return *(*Texture2D)(ptr)
}

// cptr returns C pointer
func (t *Texture2D) cptr() *C.Texture2D {
	return (*C.Texture2D)(unsafe.Pointer(t))
}

// newRenderTexture2DFromPointer - Returns new RenderTexture2D from pointer
func newRenderTexture2DFromPointer(ptr unsafe.Pointer) RenderTexture2D {
	return *(*RenderTexture2D)(ptr)
}

// cptr returns C pointer
func (r *RenderTexture2D) cptr() *C.RenderTexture2D {
	return (*C.RenderTexture2D)(unsafe.Pointer(r))
}

// NewImageFromImage - Returns new Image from Go image.Image
func NewImageFromImage(img image.Image) *Image {
	size := img.Bounds().Size()

	cx := (C.int)(size.X)
	cy := (C.int)(size.Y)
	ccolor := colorCptr(White)
	ret := C.GenImageColor(cx, cy, *ccolor)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			rcolor := NewColor(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
			ccolor = colorCptr(rcolor)

			cx = (C.int)(x)
			cy = (C.int)(y)
			C.ImageDrawPixel(&ret, cx, cy, *ccolor)
		}
	}
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImage - Load image from file into CPU memory (RAM)
func LoadImage(fileName string) *Image {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadImage(cfileName)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImageRaw - Load image from RAW file data
func LoadImageRaw(fileName string, width, height int32, format PixelFormat, headerSize int32) *Image {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cformat := (C.int)(format)
	cheaderSize := (C.int)(headerSize)
	ret := C.LoadImageRaw(cfileName, cwidth, cheight, cformat, cheaderSize)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImageAnim - Load image sequence from file (frames appended to image.data)
func LoadImageAnim(fileName string, frames *int32) *Image {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cframes := (*C.int)(frames)
	ret := C.LoadImageAnim(cfileName, cframes)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImageAnimFromMemory - Load image sequence from memory buffer
func LoadImageAnimFromMemory(fileType string, fileData []byte, dataSize int32, frames *int32) *Image {
	cfileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cfileType))
	cfileData := (*C.uchar)(unsafe.Pointer(&fileData[0]))
	cdataSize := (C.int)(dataSize)
	cframes := (*C.int)(frames)
	ret := C.LoadImageAnimFromMemory(cfileType, cfileData, cdataSize, cframes)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImageFromMemory - Load image from memory buffer, fileType refers to extension: i.e. '.png'
func LoadImageFromMemory(fileType string, fileData []byte, dataSize int32) *Image {
	cfileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cfileType))
	cfileData := (*C.uchar)(unsafe.Pointer(&fileData[0]))
	cdataSize := (C.int)(dataSize)
	ret := C.LoadImageFromMemory(cfileType, cfileData, cdataSize)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImageFromTexture - Load image from GPU texture data
func LoadImageFromTexture(texture Texture2D) *Image {
	ctexture := texture.cptr()
	ret := C.LoadImageFromTexture(*ctexture)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImageFromScreen - Load image from screen buffer and (screenshot)
func LoadImageFromScreen() *Image {
	ret := C.LoadImageFromScreen()
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// IsImageValid - Check if an image is valid (data and parameters)
func IsImageValid(image *Image) bool {
	cimage := image.cptr()
	ret := C.IsImageValid(*cimage)
	v := bool(ret)
	return v
}

// LoadTexture - Load texture from file into GPU memory (VRAM)
func LoadTexture(fileName string) Texture2D {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadTexture(cfileName)
	v := newTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadTextureFromImage - Load texture from image data
func LoadTextureFromImage(image *Image) Texture2D {
	cimage := image.cptr()
	ret := C.LoadTextureFromImage(*cimage)
	v := newTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadRenderTexture - Load texture for rendering (framebuffer)
func LoadRenderTexture(width, height int32) RenderTexture2D {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ret := C.LoadRenderTexture(cwidth, cheight)
	v := newRenderTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadRenderTextureEx - Load texture for rendering (framebuffer), with specific format
func LoadRenderTextureEx(width, height int32, format PixelFormat) RenderTexture2D {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cformat := (C.int)(format)
	ret := C.LoadRenderTextureEx(cwidth, cheight, cformat)
	v := newRenderTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadTextureCubemap - Load cubemap from image, multiple image cubemap layouts supported
func LoadTextureCubemap(image *Image, layout int32) Texture2D {
	cimage := image.cptr()
	clayout := (C.int)(layout)
	ret := C.LoadTextureCubemap(*cimage, clayout)
	v := newTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// UnloadImage - Unload image from CPU memory (RAM)
func UnloadImage(image *Image) {
	cimage := image.cptr()
	C.UnloadImage(*cimage)
}

// IsTextureValid - Check if a texture is valid (loaded in GPU)
func IsTextureValid(texture Texture2D) bool {
	ctexture := texture.cptr()
	ret := C.IsTextureValid(*ctexture)
	v := bool(ret)
	return v
}

// UnloadTexture - Unload texture from GPU memory (VRAM)
func UnloadTexture(texture Texture2D) {
	ctexture := texture.cptr()
	C.UnloadTexture(*ctexture)
}

// IsRenderTextureValid - Check if a render texture is valid (loaded in GPU)
func IsRenderTextureValid(target RenderTexture2D) bool {
	ctarget := target.cptr()
	ret := C.IsRenderTextureValid(*ctarget)
	v := bool(ret)
	return v
}

// UnloadRenderTexture - Unload render texture from GPU memory (VRAM)
func UnloadRenderTexture(target RenderTexture2D) {
	ctarget := target.cptr()
	C.UnloadRenderTexture(*ctarget)
}

// LoadImageColors - Load color data from image as a Color array (RGBA - 32bit)
//
// NOTE: Memory allocated should be freed using UnloadImageColors()
func LoadImageColors(img *Image) []color.RGBA {
	cimg := img.cptr()
	ret := C.LoadImageColors(*cimg)
	return (*[1 << 24]color.RGBA)(unsafe.Pointer(ret))[0 : img.Width*img.Height]
}

// LoadImagePalette - Load colors palette from image as a Color array (RGBA - 32bit)
//
// NOTE: Memory allocated should be freed using UnloadImagePalette()
func LoadImagePalette(img Image, maxPaletteSize int32) []color.RGBA {
	var colorCount C.int
	cimg := img.cptr()
	ret := C.LoadImagePalette(*cimg, C.int(maxPaletteSize), &colorCount)
	return unsafe.Slice((*color.RGBA)(unsafe.Pointer(ret)), colorCount)
}

// UnloadImageColors - Unload color data loaded with LoadImageColors()
func UnloadImageColors(colors []color.RGBA) {
	C.UnloadImageColors((*C.Color)(unsafe.Pointer(&colors[0])))
}

// UnloadImagePalette - Unload colors palette loaded with LoadImagePalette()
func UnloadImagePalette(colors []color.RGBA) {
	C.UnloadImagePalette((*C.Color)(unsafe.Pointer(&colors[0])))
}

// UpdateTexture - Update GPU texture with new data (pixels should be able to fill texture)
func UpdateTexture[T []color.RGBA | *image.RGBA | []byte](texture Texture2D, pixels T) {
	var cpixels unsafe.Pointer
	switch p := any(pixels).(type) {
	case []color.RGBA:
		cpixels = unsafe.Pointer(&p[0])
	case *image.RGBA:
		cpixels = unsafe.Pointer(&p.Pix[0])
	case []byte:
		cpixels = unsafe.Pointer(&p[0])
	}
	ctexture := texture.cptr()
	C.UpdateTexture(*ctexture, cpixels)
}

// UpdateTextureRec - Update GPU texture rectangle with new data (pixels and rec should fit in texture)
func UpdateTextureRec[T []color.RGBA | *image.RGBA | []byte](texture Texture2D, rec Rectangle, pixels T) {
	var cpixels unsafe.Pointer
	switch p := any(pixels).(type) {
	case []color.RGBA:
		cpixels = unsafe.Pointer(&p[0])
	case *image.RGBA:
		cpixels = unsafe.Pointer(&p.Pix[0])
	case []byte:
		cpixels = unsafe.Pointer(&p[0])
	}
	ctexture := texture.cptr()
	crec := rec.cptr()
	C.UpdateTextureRec(*ctexture, *crec, cpixels)
}

// ExportImage - Export image data to file, returns true on success
func ExportImage(image Image, fileName string) bool {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cimage := image.cptr()
	return bool(C.ExportImage(*cimage, cfileName))
}

// ExportImageToMemory - Export image to memory buffer
//
// The returned memory is a Go-managed slice. It doesn't need to be freed.
func ExportImageToMemory(image Image, fileType string) []byte {
	cfileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cfileType))
	cimage := image.cptr()

	var size C.int
	ret := C.ExportImageToMemory(*cimage, cfileType, &size)
	defer C.free(unsafe.Pointer(ret)) // free the memory

	if ret == nil {
		return nil
	}
	result := make([]byte, size)
	copy(result, unsafe.Slice((*byte)(unsafe.Pointer(ret)), size))
	return result
}

// ImageCopy - Create an image duplicate (useful for transformations)
func ImageCopy(image *Image) *Image {
	cimage := image.cptr()
	ret := C.ImageCopy(*cimage)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// ImageFromImage - Create an image from another image piece
func ImageFromImage(image Image, rec Rectangle) Image {
	cimage := image.cptr()
	crec := rec.cptr()
	ret := C.ImageFromImage(*cimage, *crec)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return *v
}

// ImageFromChannel - Create an image from a selected channel of another image (GRAYSCALE)
func ImageFromChannel(image Image, selectedChannel int32) Image {
	cimage := image.cptr()
	cselectedChannel := C.int(selectedChannel)
	ret := C.ImageFromChannel(*cimage, cselectedChannel)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return *v
}

// ImageText - Create an image from text (default font)
func ImageText(text string, fontSize int32, col color.RGBA) *Image {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.int)(fontSize)
	ccolor := colorCptr(col)
	ret := C.ImageText(ctext, cfontSize, *ccolor)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// ImageTextEx - Create an image from text (custom sprite font)
func ImageTextEx(font Font, text string, fontSize, spacing float32, tint color.RGBA) *Image {
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	ctint := colorCptr(tint)
	ret := C.ImageTextEx(*cfont, ctext, cfontSize, cspacing, *ctint)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// ImageFormat - Convert image data to desired format
func ImageFormat(image *Image, newFormat PixelFormat) {
	cimage := image.cptr()
	cnewFormat := (C.int)(newFormat)
	C.ImageFormat(cimage, cnewFormat)
}

// ImageToPOT - Convert image to POT (power-of-two)
func ImageToPOT(image *Image, fill color.RGBA) {
	cimage := image.cptr()
	cfill := colorCptr(fill)
	C.ImageToPOT(cimage, *cfill)
}

// ImageCrop - Crop an image to a defined rectangle
func ImageCrop(image *Image, crop Rectangle) {
	cimage := image.cptr()
	ccrop := crop.cptr()
	C.ImageCrop(cimage, *ccrop)
}

// ImageAlphaCrop - Crop image depending on alpha value
func ImageAlphaCrop(image *Image, threshold float32) {
	cimage := image.cptr()
	cthreshold := (C.float)(threshold)
	C.ImageAlphaCrop(cimage, cthreshold)
}

// ImageAlphaClear - Clear alpha channel to desired color
func ImageAlphaClear(image *Image, col color.RGBA, threshold float32) {
	cimage := image.cptr()
	ccolor := colorCptr(col)
	cthreshold := (C.float)(threshold)
	C.ImageAlphaClear(cimage, *ccolor, cthreshold)
}

// ImageAlphaMask - Apply alpha mask to image
func ImageAlphaMask(image, alphaMask *Image) {
	cimage := image.cptr()
	calphaMask := alphaMask.cptr()
	C.ImageAlphaMask(cimage, *calphaMask)
}

// ImageAlphaPremultiply - Premultiply alpha channel
func ImageAlphaPremultiply(image *Image) {
	cimage := image.cptr()
	C.ImageAlphaPremultiply(cimage)
}

// ImageBlurGaussian - Apply Gaussian blur using a box blur approximation
func ImageBlurGaussian(image *Image, blurSize int32) {
	cimage := image.cptr()
	cblurSize := C.int(blurSize)
	C.ImageBlurGaussian(cimage, cblurSize)
}

// ImageKernelConvolution - Apply custom square convolution kernel to image
func ImageKernelConvolution(image *Image, kernel []float32) {
	cimage := image.cptr()
	ckernel := (*C.float)(unsafe.Pointer(&kernel[0]))
	C.ImageKernelConvolution(cimage, ckernel, C.int(len(kernel)))
}

// ImageResize - Resize image (Bicubic scaling algorithm)
func ImageResize(image *Image, newWidth, newHeight int32) {
	cimage := image.cptr()
	cnewWidth := (C.int)(newWidth)
	cnewHeight := (C.int)(newHeight)
	C.ImageResize(cimage, cnewWidth, cnewHeight)
}

// ImageResizeNN - Resize image (Nearest-Neighbor scaling algorithm)
func ImageResizeNN(image *Image, newWidth, newHeight int32) {
	cimage := image.cptr()
	cnewWidth := (C.int)(newWidth)
	cnewHeight := (C.int)(newHeight)
	C.ImageResizeNN(cimage, cnewWidth, cnewHeight)
}

// ImageResizeCanvas - Resize canvas and fill with color
func ImageResizeCanvas(image *Image, newWidth, newHeight, offsetX, offsetY int32, fill color.RGBA) {
	cimage := image.cptr()
	cnewWidth := (C.int)(newWidth)
	cnewHeight := (C.int)(newHeight)
	coffsetX := (C.int)(offsetX)
	coffsetY := (C.int)(offsetY)
	ccolor := colorCptr(fill)
	C.ImageResizeCanvas(cimage, cnewWidth, cnewHeight, coffsetX, coffsetY, *ccolor)
}

// ImageMipmaps - Compute all mipmap levels for a provided image
func ImageMipmaps(image *Image) {
	cimage := image.cptr()
	C.ImageMipmaps(cimage)
}

// ImageDither - Dither image data to 16bpp or lower (Floyd-Steinberg dithering)
func ImageDither(image *Image, rBpp, gBpp, bBpp, aBpp int32) {
	cimage := image.cptr()
	crBpp := (C.int)(rBpp)
	cgBpp := (C.int)(gBpp)
	cbBpp := (C.int)(bBpp)
	caBpp := (C.int)(aBpp)
	C.ImageDither(cimage, crBpp, cgBpp, cbBpp, caBpp)
}

// ImageFlipVertical - Flip image vertically
func ImageFlipVertical(image *Image) {
	cimage := image.cptr()
	C.ImageFlipVertical(cimage)
}

// ImageFlipHorizontal - Flip image horizontally
func ImageFlipHorizontal(image *Image) {
	cimage := image.cptr()
	C.ImageFlipHorizontal(cimage)
}

// ImageRotate - Rotate image by input angle in degrees (-359 to 359)
func ImageRotate(image *Image, degrees int32) {
	cimage := image.cptr()
	cdegrees := (C.int)(degrees)
	C.ImageRotate(cimage, cdegrees)
}

// ImageRotateCW - Rotate image clockwise 90deg
func ImageRotateCW(image *Image) {
	cimage := image.cptr()
	C.ImageRotateCW(cimage)
}

// ImageRotateCCW - Rotate image counter-clockwise 90deg
func ImageRotateCCW(image *Image) {
	cimage := image.cptr()
	C.ImageRotateCCW(cimage)
}

// ImageColorTint - Modify image color: tint
func ImageColorTint(image *Image, col color.RGBA) {
	cimage := image.cptr()
	ccolor := colorCptr(col)
	C.ImageColorTint(cimage, *ccolor)
}

// ImageColorInvert - Modify image color: invert
func ImageColorInvert(image *Image) {
	cimage := image.cptr()
	C.ImageColorInvert(cimage)
}

// ImageColorGrayscale - Modify image color: grayscale
func ImageColorGrayscale(image *Image) {
	cimage := image.cptr()
	C.ImageColorGrayscale(cimage)
}

// ImageColorContrast - Modify image color: contrast (-100 to 100)
func ImageColorContrast(image *Image, contrast int32) {
	cimage := image.cptr()
	ccontrast := (C.int)(contrast)
	C.ImageColorContrast(cimage, ccontrast)
}

// ImageColorBrightness - Modify image color: brightness (-255 to 255)
func ImageColorBrightness(image *Image, brightness int32) {
	cimage := image.cptr()
	cbrightness := (C.int)(brightness)
	C.ImageColorBrightness(cimage, cbrightness)
}

// ImageColorReplace - Modify image color: replace color
func ImageColorReplace(image *Image, col, replace color.RGBA) {
	cimage := image.cptr()
	ccolor := colorCptr(col)
	creplace := colorCptr(replace)
	C.ImageColorReplace(cimage, *ccolor, *creplace)
}

// GetImageAlphaBorder - Get image alpha border rectangle
func GetImageAlphaBorder(image Image, threshold float32) Rectangle {
	cimage := image.cptr()
	cthreshold := C.float(threshold)
	ret := C.GetImageAlphaBorder(*cimage, cthreshold)
	v := newRectangleFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetImageColor - Get image pixel color at (x, y) position
func GetImageColor(image Image, x, y int32) color.RGBA {
	cimage := image.cptr()
	cx := (C.int)(x)
	cy := (C.int)(y)

	ret := C.GetImageColor(*cimage, cx, cy)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// ImageClearBackground - Clear image background with given color
func ImageClearBackground(dst *Image, col color.RGBA) {
	cdst := dst.cptr()
	ccolor := colorCptr(col)
	C.ImageClearBackground(cdst, *ccolor)
}

// ImageDrawImage - Draw an image within an image
func ImageDrawImage(dst *Image, src *Image, posX, posY int32, tint color.RGBA) {
	cdst := dst.cptr()
	csrc := src.cptr()
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	ctint := colorCptr(tint)
	C.ImageDrawImage(cdst, *csrc, cposX, cposY, *ctint)
}

// ImageDrawImageEx - Draw an image with scaling and rotation within an image
func ImageDrawImageEx(dst *Image, src *Image, position Vector2, rotation, scale float32, tint color.RGBA) {
	cdst := dst.cptr()
	csrc := src.cptr()
	cposition := position.cptr()
	crotation := (C.float)(rotation)
	cscale := (C.float)(scale)
	ctint := colorCptr(tint)
	C.ImageDrawImageEx(cdst, *csrc, *cposition, crotation, cscale, *ctint)
}

// ImageDrawImageRec - Draw a part of an image defined by a rectangle within an image
func ImageDrawImageRec(dst *Image, src *Image, srcRec Rectangle, position Vector2, tint color.RGBA) {
	cdst := dst.cptr()
	csrc := src.cptr()
	csrcRec := srcRec.cptr()
	cposition := position.cptr()
	ctint := colorCptr(tint)
	C.ImageDrawImageRec(cdst, *csrc, *csrcRec, *cposition, *ctint)
}

// ImageDrawImagePro - Draw a part of an image defined by a rectangle into destination rectangle, with scaling and rotation, within an image
func ImageDrawImagePro(dst *Image, src *Image, srcRec, dstRec Rectangle, origin Vector2, rotation float32, tint color.RGBA) {
	cdst := dst.cptr()
	csrc := src.cptr()
	csrcRec := srcRec.cptr()
	cdstRec := dstRec.cptr()
	corigin := origin.cptr()
	crotation := (C.float)(rotation)
	ctint := colorCptr(tint)
	C.ImageDrawImagePro(cdst, *csrc, *csrcRec, *cdstRec, *corigin, crotation, *ctint)
}

// ImageDrawLine - Draw line within an image
func ImageDrawLine(dst *Image, startPosX, startPosY, endPosX, endPosY int32, col color.RGBA) {
	cdst := dst.cptr()
	cstartPosX := (C.int)(startPosX)
	cstartPosY := (C.int)(startPosY)
	cendPosX := (C.int)(endPosX)
	cendPosY := (C.int)(endPosY)
	ccolor := colorCptr(col)
	C.ImageDrawLine(cdst, cstartPosX, cstartPosY, cendPosX, cendPosY, *ccolor)
}

// ImageDrawLineV - Draw line within an image (Vector version)
func ImageDrawLineV(dst *Image, start, end Vector2, col color.RGBA) {
	cdst := dst.cptr()
	cstart := start.cptr()
	cend := end.cptr()
	ccolor := colorCptr(col)
	C.ImageDrawLineV(cdst, *cstart, *cend, *ccolor)
}

// ImageDrawLineEx - Draw a line defining thickness within an image
func ImageDrawLineEx(dst *Image, start, end Vector2, thick int32, col color.RGBA) {
	cdst := dst.cptr()
	cstart := start.cptr()
	cend := end.cptr()
	cthick := C.int(thick)
	ccolor := colorCptr(col)
	C.ImageDrawLineEx(cdst, *cstart, *cend, cthick, *ccolor)
}

// ImageDrawLineStrip - Draw a lines sequence within an image
func ImageDrawLineStrip(dst *Image, points []Vector2, col color.RGBA) {
	cdst := dst.cptr()
	cpoints := (*C.Vector2)(unsafe.Pointer(unsafe.SliceData(points)))
	cpointCount := (C.int)(len(points))
	ccolor := colorCptr(col)
	C.ImageDrawLineStrip(cdst, cpoints, cpointCount, *ccolor)
}

// ImageDrawCircle - Draw a filled circle within an image
func ImageDrawCircle(dst *Image, centerX, centerY, radius int32, col color.RGBA) {
	cdst := dst.cptr()
	ccenterX := (C.int)(centerX)
	ccenterY := (C.int)(centerY)
	cradius := (C.int)(radius)
	ccolor := colorCptr(col)
	C.ImageDrawCircle(cdst, ccenterX, ccenterY, cradius, *ccolor)
}

// ImageDrawCircleV - Draw a filled circle within an image (Vector version)
func ImageDrawCircleV(dst *Image, center Vector2, radius int32, col color.RGBA) {
	cdst := dst.cptr()
	ccenter := center.cptr()
	cradius := (C.int)(radius)
	ccolor := colorCptr(col)
	C.ImageDrawCircleV(cdst, *ccenter, cradius, *ccolor)
}

// ImageDrawCircleLines - Draw circle outline within an image
func ImageDrawCircleLines(dst *Image, centerX, centerY, radius int32, col color.RGBA) {
	cdst := dst.cptr()
	ccenterX := (C.int)(centerX)
	ccenterY := (C.int)(centerY)
	cradius := (C.int)(radius)
	ccolor := colorCptr(col)
	C.ImageDrawCircleLines(cdst, ccenterX, ccenterY, cradius, *ccolor)
}

// ImageDrawCircleLinesV - Draw circle outline within an image (Vector version)
func ImageDrawCircleLinesV(dst *Image, center Vector2, radius int32, col color.RGBA) {
	cdst := dst.cptr()
	ccenter := center.cptr()
	cradius := (C.int)(radius)
	ccolor := colorCptr(col)
	C.ImageDrawCircleLinesV(cdst, *ccenter, cradius, *ccolor)
}

// ImageDrawCircleGradient - Draw a gradient-filled circle within an image
func ImageDrawCircleGradient(dst *Image, center Vector2, radius float32, inner, outer color.RGBA) {
	cdst := dst.cptr()
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	cinner := colorCptr(inner)
	couter := colorCptr(outer)
	C.ImageDrawCircleGradient(cdst, *ccenter, cradius, *cinner, *couter)
}

// ImageDrawPixel - Draw pixel within an image
func ImageDrawPixel(dst *Image, posX, posY int32, col color.RGBA) {
	cdst := dst.cptr()
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	ccolor := colorCptr(col)
	C.ImageDrawPixel(cdst, cposX, cposY, *ccolor)
}

// ImageDrawPixelV - Draw pixel within an image (Vector version)
func ImageDrawPixelV(dst *Image, position Vector2, col color.RGBA) {
	cdst := dst.cptr()
	cposition := position.cptr()
	ccolor := colorCptr(col)
	C.ImageDrawPixelV(cdst, *cposition, *ccolor)
}

// ImageDrawRectangle - Draw rectangle within an image
func ImageDrawRectangle(dst *Image, posX, posY, width, height int32, col color.RGBA) {
	cdst := dst.cptr()
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor := colorCptr(col)
	C.ImageDrawRectangle(cdst, cposX, cposY, cwidth, cheight, *ccolor)
}

// ImageDrawRectangleV - Draw rectangle within an image (Vector version)
func ImageDrawRectangleV(dst *Image, position, size Vector2, col color.RGBA) {
	cdst := dst.cptr()
	cposition := position.cptr()
	csize := size.cptr()
	ccolor := colorCptr(col)
	C.ImageDrawRectangleV(cdst, *cposition, *csize, *ccolor)
}

// ImageDrawRectangleLines - Draw rectangle lines within an image
func ImageDrawRectangleLines(dst *Image, posX, posY, width, height int32, col color.RGBA) {
	cdst := dst.cptr()
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor := colorCptr(col)
	C.ImageDrawRectangleLines(cdst, cposX, cposY, cwidth, cheight, *ccolor)
}

// ImageDrawRectangleLinesEx - Draw rectangle lines within an image with line thickness
func ImageDrawRectangleLinesEx(dst *Image, rec Rectangle, thick int, col color.RGBA) {
	cdst := dst.cptr()
	crec := rec.cptr()
	cthick := (C.int)(thick)
	ccolor := colorCptr(col)
	C.ImageDrawRectangleLinesEx(cdst, *crec, cthick, *ccolor)
}

// ImageDrawRectanglePro - Draw a color-filled rectangle with pro parameters within an image
func ImageDrawRectanglePro(dst *Image, rec Rectangle, origin Vector2, rotation float32, col color.RGBA) {
	cdst := dst.cptr()
	crec := rec.cptr()
	corigin := origin.cptr()
	crotation := (C.float)(rotation)
	ccolor := colorCptr(col)
	C.ImageDrawRectanglePro(cdst, *crec, *corigin, crotation, *ccolor)
}

// ImageDrawRectangleGradientEx - Draw rectangle with gradient colors within an image, counter-clockwise color order
func ImageDrawRectangleGradientEx(dst *Image, rec Rectangle, col1, col2, col3, col4 color.RGBA) {
	cdst := dst.cptr()
	crec := rec.cptr()
	ccol1 := colorCptr(col1)
	ccol2 := colorCptr(col2)
	ccol3 := colorCptr(col3)
	ccol4 := colorCptr(col4)
	C.ImageDrawRectangleGradientEx(cdst, *crec, *ccol1, *ccol2, *ccol3, *ccol4)
}

// ImageDrawTriangle - Draw triangle within an image
func ImageDrawTriangle(dst *Image, v1, v2, v3 Vector2, col color.RGBA) {
	cdst := dst.cptr()
	cv1 := v1.cptr()
	cv2 := v2.cptr()
	cv3 := v3.cptr()
	ccol := colorCptr(col)
	C.ImageDrawTriangle(cdst, *cv1, *cv2, *cv3, *ccol)
}

// ImageDrawTriangleGradient - Draw triangle with interpolated colors within an image
func ImageDrawTriangleGradient(dst *Image, v1, v2, v3 Vector2, c1, c2, c3 color.RGBA) {
	cdst := dst.cptr()
	cv1 := v1.cptr()
	cv2 := v2.cptr()
	cv3 := v3.cptr()
	cc1 := colorCptr(c1)
	cc2 := colorCptr(c2)
	cc3 := colorCptr(c3)
	C.ImageDrawTriangleGradient(cdst, *cv1, *cv2, *cv3, *cc1, *cc2, *cc3)
}

// ImageDrawTriangleLines - Draw triangle outline within an image
func ImageDrawTriangleLines(dst *Image, v1, v2, v3 Vector2, col color.RGBA) {
	cdst := dst.cptr()
	cv1 := v1.cptr()
	cv2 := v2.cptr()
	cv3 := v3.cptr()
	ccol := colorCptr(col)
	C.ImageDrawTriangleLines(cdst, *cv1, *cv2, *cv3, *ccol)
}

// ImageDrawTriangleFan - Draw a triangle fan defined by points within an image (first vertex is the center)
func ImageDrawTriangleFan(dst *Image, points []Vector2, col color.RGBA) {
	cdst := dst.cptr()
	cpoints := (*C.Vector2)(unsafe.Pointer(&points[0]))
	pointCount := C.int(len(points))
	ccol := colorCptr(col)
	C.ImageDrawTriangleFan(cdst, cpoints, pointCount, *ccol)
}

// ImageDrawTriangleStrip - Draw a triangle strip defined by points within an image
func ImageDrawTriangleStrip(dst *Image, points []Vector2, col color.RGBA) {
	cdst := dst.cptr()
	cpoints := (*C.Vector2)(unsafe.Pointer(&points[0]))
	pointCount := C.int(len(points))
	ccol := colorCptr(col)
	C.ImageDrawTriangleStrip(cdst, cpoints, pointCount, *ccol)
}

// ImageDrawRectangleRec - Draw rectangle within an image
func ImageDrawRectangleRec(dst *Image, rec Rectangle, col color.RGBA) {
	cdst := dst.cptr()
	crec := rec.cptr()
	ccolor := colorCptr(col)
	C.ImageDrawRectangleRec(cdst, *crec, *ccolor)
}

// ImageDrawText - Draw text (using default font) within an image (destination)
func ImageDrawText(dst *Image, posX, posY int32, text string, fontSize int32, col color.RGBA) {
	cdst := dst.cptr()
	posx := (C.int)(posX)
	posy := (C.int)(posY)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.int)(fontSize)
	ccolor := colorCptr(col)
	C.ImageDrawText(cdst, ctext, posx, posy, cfontSize, *ccolor)
}

// ImageDrawTextEx - Draw text (custom sprite font) within an image (destination)
func ImageDrawTextEx(dst *Image, position Vector2, font Font, text string, fontSize, spacing float32, tint color.RGBA) {
	cdst := dst.cptr()
	cposition := position.cptr()
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	ctint := colorCptr(tint)
	C.ImageDrawTextEx(cdst, *cfont, ctext, *cposition, cfontSize, cspacing, *ctint)
}

// ImageDrawTextPro - Draw text using Font and pro parameters (rotation)
func ImageDrawTextPro(dst *Image, font Font, text string, position, origin Vector2, rotation, fontSize, spacing float32, tint color.RGBA) {
	cdst := dst.cptr()
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cposition := position.cptr()
	corigin := origin.cptr()
	crotation := (C.float)(rotation)
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	ctint := colorCptr(tint)
	C.ImageDrawTextPro(cdst, *cfont, ctext, *cposition, *corigin, crotation, cfontSize, cspacing, *ctint)
}

// GenImageColor - Generate image: plain color
func GenImageColor(width, height int, col color.RGBA) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor := colorCptr(col)

	ret := C.GenImageColor(cwidth, cheight, *ccolor)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageGradientLinear - Generate image: linear gradient, direction in degrees [0..360], 0=Vertical gradient
func GenImageGradientLinear(width, height, direction int, start, end color.RGBA) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cdensity := (C.int)(direction)
	cstart := colorCptr(start)
	cend := colorCptr(end)

	ret := C.GenImageGradientLinear(cwidth, cheight, cdensity, *cstart, *cend)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageGradientRadial - Generate image: radial gradient
func GenImageGradientRadial(width, height int, density float32, inner, outer color.RGBA) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cdensity := (C.float)(density)
	cinner := colorCptr(inner)
	couter := colorCptr(outer)

	ret := C.GenImageGradientRadial(cwidth, cheight, cdensity, *cinner, *couter)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageGradientSquare - Generate image: square gradient
func GenImageGradientSquare(width, height int, density float32, inner, outer color.RGBA) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cdensity := (C.float)(density)
	cinner := colorCptr(inner)
	couter := colorCptr(outer)

	ret := C.GenImageGradientSquare(cwidth, cheight, cdensity, *cinner, *couter)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageChecked - Generate image: checked
func GenImageChecked(width, height, checksX, checksY int, col1, col2 color.RGBA) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cchecksX := (C.int)(checksX)
	cchecksY := (C.int)(checksY)
	ccol1 := colorCptr(col1)
	ccol2 := colorCptr(col2)

	ret := C.GenImageChecked(cwidth, cheight, cchecksX, cchecksY, *ccol1, *ccol2)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageWhiteNoise - Generate image: white noise
func GenImageWhiteNoise(width, height int, factor float32) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cfactor := (C.float)(factor)

	ret := C.GenImageWhiteNoise(cwidth, cheight, cfactor)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImagePerlinNoise - Generate image: perlin noise
func GenImagePerlinNoise(width, height, offsetX, offsetY int, scale float32) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	coffsetX := (C.int)(offsetX)
	coffsetY := (C.int)(offsetY)
	cscale := (C.float)(scale)

	ret := C.GenImagePerlinNoise(cwidth, cheight, coffsetX, coffsetY, cscale)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageCellular - Generate image: cellular algorithm, bigger tileSize means bigger cells
func GenImageCellular(width, height, tileSize int) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ctileSize := (C.int)(tileSize)

	ret := C.GenImageCellular(cwidth, cheight, ctileSize)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageText - Generate image: grayscale image from text data
func GenImageText(width, height int, text string) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	ret := C.GenImageText(cwidth, cheight, ctext)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenTextureMipmaps - Generate GPU mipmaps for a texture
func GenTextureMipmaps(texture *Texture2D) {
	ctexture := texture.cptr()
	C.GenTextureMipmaps(ctexture)
}

// SetTextureFilter - Set texture scaling filter mode
func SetTextureFilter(texture Texture2D, filter TextureFilterMode) {
	ctexture := texture.cptr()
	cfilter := (C.int)(filter)
	C.SetTextureFilter(*ctexture, cfilter)
}

// SetTextureWrap - Set texture wrapping mode
func SetTextureWrap(texture Texture2D, wrap TextureWrapMode) {
	ctexture := texture.cptr()
	cwrap := (C.int)(wrap)
	C.SetTextureWrap(*ctexture, cwrap)
}

// DrawTexture - Draw a Texture2D
func DrawTexture(texture Texture2D, posX int32, posY int32, tint color.RGBA) {
	ctexture := texture.cptr()
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	ctint := colorCptr(tint)
	C.DrawTexture(*ctexture, cposX, cposY, *ctint)
}

// DrawTextureV - Draw a Texture2D with position defined as Vector2
func DrawTextureV(texture Texture2D, position Vector2, tint color.RGBA) {
	ctexture := texture.cptr()
	cposition := position.cptr()
	ctint := colorCptr(tint)
	C.DrawTextureV(*ctexture, *cposition, *ctint)
}

// DrawTextureEx - Draw a Texture2D with extended parameters
func DrawTextureEx(texture Texture2D, position Vector2, rotation, scale float32, tint color.RGBA) {
	ctexture := texture.cptr()
	cposition := position.cptr()
	crotation := (C.float)(rotation)
	cscale := (C.float)(scale)
	ctint := colorCptr(tint)
	C.DrawTextureEx(*ctexture, *cposition, crotation, cscale, *ctint)
}

// DrawTextureRec - Draw a part of a texture defined by a rectangle
func DrawTextureRec(texture Texture2D, source Rectangle, position Vector2, tint color.RGBA) {
	ctexture := texture.cptr()
	csource := source.cptr()
	cposition := position.cptr()
	ctint := colorCptr(tint)
	C.DrawTextureRec(*ctexture, *csource, *cposition, *ctint)
}

// DrawTexturePro - Draw a part of a texture defined by a rectangle with 'pro' parameters
func DrawTexturePro(texture Texture2D, source, dest Rectangle, origin Vector2, rotation float32, tint color.RGBA) {
	ctexture := texture.cptr()
	csource := source.cptr()
	cdest := dest.cptr()
	corigin := origin.cptr()
	crotation := (C.float)(rotation)
	ctint := colorCptr(tint)
	C.DrawTexturePro(*ctexture, *csource, *cdest, *corigin, crotation, *ctint)
}

// cptr returns C pointer
func (n *NPatchInfo) cptr() *C.NPatchInfo {
	return (*C.NPatchInfo)(unsafe.Pointer(n))
}

// DrawTextureNPatch - Draws a texture (or part of it) that stretches or shrinks nicely
func DrawTextureNPatch(texture Texture2D, nPatchInfo NPatchInfo, dest Rectangle, origin Vector2, rotation float32, tint color.RGBA) {
	ctexture := texture.cptr()
	cnPatchInfo := nPatchInfo.cptr()
	cdest := dest.cptr()
	corigin := origin.cptr()
	crotation := (C.float)(rotation)
	ctint := colorCptr(tint)
	C.DrawTextureNPatch(*ctexture, *cnPatchInfo, *cdest, *corigin, crotation, *ctint)
}
