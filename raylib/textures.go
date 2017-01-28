package raylib

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type TextureFormat int32

// Texture formats
// NOTE: Support depends on OpenGL version and platform
const (
	// 8 bit per pixel (no alpha)
	UncompressedGrayscale TextureFormat = C.UNCOMPRESSED_GRAYSCALE
	// 16 bpp (2 channels)
	UncompressedGrayAlpha TextureFormat = C.UNCOMPRESSED_GRAY_ALPHA
	// 16 bpp
	UncompressedR5g6b5 TextureFormat = C.UNCOMPRESSED_R5G6B5
	// 24 bpp
	UncompressedR8g8b8 TextureFormat = C.UNCOMPRESSED_R8G8B8
	// 16 bpp (1 bit alpha)
	UncompressedR5g5b5a1 TextureFormat = C.UNCOMPRESSED_R5G5B5A1
	// 16 bpp (4 bit alpha)
	UncompressedR4g4b4a4 TextureFormat = C.UNCOMPRESSED_R4G4B4A4
	// 32 bpp
	UncompressedR8g8b8a8 TextureFormat = C.UNCOMPRESSED_R8G8B8A8
	// 4 bpp (no alpha)
	CompressedDxt1Rgb TextureFormat = C.COMPRESSED_DXT1_RGB
	// 4 bpp (1 bit alpha)
	CompressedDxt1Rgba TextureFormat = C.COMPRESSED_DXT1_RGBA
	// 8 bpp
	CompressedDxt3Rgba TextureFormat = C.COMPRESSED_DXT3_RGBA
	// 8 bpp
	CompressedDxt5Rgba TextureFormat = C.COMPRESSED_DXT5_RGBA
	// 4 bpp
	CompressedEtc1Rgb TextureFormat = C.COMPRESSED_ETC1_RGB
	// 4 bpp
	CompressedEtc2Rgb TextureFormat = C.COMPRESSED_ETC2_RGB
	// 8 bpp
	CompressedEtc2EacRgba TextureFormat = C.COMPRESSED_ETC2_EAC_RGBA
	// 4 bpp
	CompressedPvrtRgb TextureFormat = C.COMPRESSED_PVRT_RGB
	// 4 bpp
	CompressedPvrtRgba TextureFormat = C.COMPRESSED_PVRT_RGBA
	// 8 bpp
	CompressedAstc4x4Rgba TextureFormat = C.COMPRESSED_ASTC_4x4_RGBA
	// 2 bpp
	CompressedAstc8x8Rgba TextureFormat = C.COMPRESSED_ASTC_8x8_RGBA
)

type TextureFilterMode int32

// Texture parameters: filter mode
// NOTE 1: Filtering considers mipmaps if available in the texture
// NOTE 2: Filter is accordingly set for minification and magnification
const (
	// No filter, just pixel aproximation
	FilterPoint TextureFilterMode = C.FILTER_POINT
	// Linear filtering
	FilterBilinear TextureFilterMode = C.FILTER_BILINEAR
	// Trilinear filtering (linear with mipmaps)
	FilterTrilinear TextureFilterMode = C.FILTER_TRILINEAR
	// Anisotropic filtering 4x
	FilterAnisotropic4x TextureFilterMode = C.FILTER_ANISOTROPIC_4X
	// Anisotropic filtering 8x
	FilterAnisotropic8x TextureFilterMode = C.FILTER_ANISOTROPIC_8X
	// Anisotropic filtering 16x
	FilterAnisotropic16x TextureFilterMode = C.FILTER_ANISOTROPIC_16X
)

type TextureWrapMode int32

// Texture parameters: wrap mode
const (
	WrapRepeat TextureWrapMode = C.WRAP_REPEAT
	WrapClamp  TextureWrapMode = C.WRAP_CLAMP
	WrapMirror TextureWrapMode = C.WRAP_MIRROR
)

// Image type, bpp always RGBA (32bit)
// NOTE: Data stored in CPU memory (RAM)
type Image struct {
	// Image raw data
	Data unsafe.Pointer
	// Image base width
	Width int32
	// Image base height
	Height int32
	// Mipmap levels, 1 by default
	Mipmaps int32
	// Data format (TextureFormat)
	Format TextureFormat
}

func (i *Image) cptr() *C.Image {
	return (*C.Image)(unsafe.Pointer(i))
}

// Returns new Image
func NewImage(data unsafe.Pointer, width, height, mipmaps int32, format TextureFormat) *Image {
	return &Image{data, width, height, mipmaps, format}
}

// Returns new Image from pointer
func NewImageFromPointer(ptr unsafe.Pointer) *Image {
	return (*Image)(ptr)
}

// Texture2D type, bpp always RGBA (32bit)
// NOTE: Data stored in GPU memory
type Texture2D struct {
	// OpenGL texture id
	Id uint32
	// Texture base width
	Width int32
	// Texture base height
	Height int32
	// Mipmap levels, 1 by default
	Mipmaps int32
	// Data format (TextureFormat)
	Format TextureFormat
}

func (t *Texture2D) cptr() *C.Texture2D {
	return (*C.Texture2D)(unsafe.Pointer(t))
}

// Returns new Texture2D
func NewTexture2D(id uint32, width, height, mipmaps int32, format TextureFormat) Texture2D {
	return Texture2D{id, width, height, mipmaps, format}
}

// Returns new Texture2D from pointer
func NewTexture2DFromPointer(ptr unsafe.Pointer) Texture2D {
	return *(*Texture2D)(ptr)
}

// RenderTexture2D type, for texture rendering
type RenderTexture2D struct {
	// Render texture (fbo) id
	Id uint32
	// Color buffer attachment texture
	Texture Texture2D
	// Depth buffer attachment texture
	Depth Texture2D
}

func (r *RenderTexture2D) cptr() *C.RenderTexture2D {
	return (*C.RenderTexture2D)(unsafe.Pointer(r))
}

// Returns new RenderTexture2D
func NewRenderTexture2D(id uint32, texture, depth Texture2D) RenderTexture2D {
	return RenderTexture2D{id, texture, depth}
}

// Returns new RenderTexture2D from pointer
func NewRenderTexture2DFromPointer(ptr unsafe.Pointer) RenderTexture2D {
	return *(*RenderTexture2D)(ptr)
}

// Load an image into CPU memory (RAM)
func LoadImage(fileName string) *Image {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadImage(cfileName)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// Load image data from Color array data (RGBA - 32bit)
func LoadImageEx(pixels []Color, width int32, height int32) *Image {
	cpixels := pixels[0].cptr()
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ret := C.LoadImageEx(cpixels, cwidth, cheight)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// Load image data from RAW file
func LoadImageRaw(fileName string, width int32, height int32, format TextureFormat, headerSize int32) *Image {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cformat := (C.int)(format)
	cheaderSize := (C.int)(headerSize)
	ret := C.LoadImageRaw(cfileName, cwidth, cheight, cformat, cheaderSize)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// Load an image as texture into GPU memory
func LoadTexture(fileName string) Texture2D {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadTexture(cfileName)
	v := NewTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// Load a texture from image data
func LoadTextureFromImage(image *Image) Texture2D {
	cimage := image.cptr()
	ret := C.LoadTextureFromImage(*cimage)
	v := NewTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// Load a texture to be used for rendering
func LoadRenderTexture(width int32, height int32) RenderTexture2D {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ret := C.LoadRenderTexture(cwidth, cheight)
	v := NewRenderTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// Unload image from CPU memory (RAM)
func UnloadImage(image *Image) {
	cimage := image.cptr()
	C.UnloadImage(*cimage)
}

// Unload texture from GPU memory
func UnloadTexture(texture Texture2D) {
	ctexture := texture.cptr()
	C.UnloadTexture(*ctexture)
}

// Unload render texture from GPU memory
func UnloadRenderTexture(target RenderTexture2D) {
	ctarget := target.cptr()
	C.UnloadRenderTexture(*ctarget)
}

// Get pixel data from image
func GetImageData(image *Image) unsafe.Pointer {
	cimage := image.cptr()
	ret := C.GetImageData(*cimage)
	return unsafe.Pointer(ret)
}

// Get pixel data from GPU texture and return an Image
func GetTextureData(texture Texture2D) *Image {
	ctexture := texture.cptr()
	ret := C.GetTextureData(*ctexture)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// Update GPU texture with new data
func UpdateTexture(texture Texture2D, pixels unsafe.Pointer) {
	ctexture := texture.cptr()
	cpixels := (unsafe.Pointer)(unsafe.Pointer(pixels))
	C.UpdateTexture(*ctexture, cpixels)
}

// Convert image to POT (power-of-two)
func ImageToPOT(image *Image, fillColor Color) {
	cimage := image.cptr()
	cfillColor := fillColor.cptr()
	C.ImageToPOT(cimage, *cfillColor)
}

// Convert image data to desired format
func ImageFormat(image *Image, newFormat int32) {
	cimage := image.cptr()
	cnewFormat := (C.int)(newFormat)
	C.ImageFormat(cimage, cnewFormat)
}

// Apply alpha mask to image
func ImageAlphaMask(image *Image, alphaMask *Image) {
	cimage := image.cptr()
	calphaMask := alphaMask.cptr()
	C.ImageAlphaMask(cimage, *calphaMask)
}

// Dither image data to 16bpp or lower (Floyd-Steinberg dithering)
func ImageDither(image *Image, rBpp int32, gBpp int32, bBpp int32, aBpp int32) {
	cimage := image.cptr()
	crBpp := (C.int)(rBpp)
	cgBpp := (C.int)(gBpp)
	cbBpp := (C.int)(bBpp)
	caBpp := (C.int)(aBpp)
	C.ImageDither(cimage, crBpp, cgBpp, cbBpp, caBpp)
}

// Create an image duplicate (useful for transformations)
func ImageCopy(image *Image) *Image {
	cimage := image.cptr()
	ret := C.ImageCopy(*cimage)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// Crop an image to a defined rectangle
func ImageCrop(image *Image, crop Rectangle) {
	cimage := image.cptr()
	ccrop := crop.cptr()
	C.ImageCrop(cimage, *ccrop)
}

// Resize an image (bilinear filtering)
func ImageResize(image *Image, newWidth int32, newHeight int32) {
	cimage := image.cptr()
	cnewWidth := (C.int)(newWidth)
	cnewHeight := (C.int)(newHeight)
	C.ImageResize(cimage, cnewWidth, cnewHeight)
}

// Resize an image (Nearest-Neighbor scaling algorithm)
func ImageResizeNN(image *Image, newWidth int32, newHeight int32) {
	cimage := image.cptr()
	cnewWidth := (C.int)(newWidth)
	cnewHeight := (C.int)(newHeight)
	C.ImageResizeNN(cimage, cnewWidth, cnewHeight)
}

// Create an image from text (default font)
func ImageText(text string, fontSize int32, color Color) *Image {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.int)(fontSize)
	ccolor := color.cptr()
	ret := C.ImageText(ctext, cfontSize, *ccolor)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// Create an image from text (custom sprite font)
func ImageTextEx(font SpriteFont, text string, fontSize float32, spacing int32, tint Color) *Image {
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.float)(fontSize)
	cspacing := (C.int)(spacing)
	ctint := tint.cptr()
	ret := C.ImageTextEx(*cfont, ctext, cfontSize, cspacing, *ctint)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// Draw a source image within a destination image
func ImageDraw(dst *Image, src *Image, srcRec Rectangle, dstRec Rectangle) {
	cdst := dst.cptr()
	csrc := src.cptr()
	csrcRec := srcRec.cptr()
	cdstRec := dstRec.cptr()
	C.ImageDraw(cdst, *csrc, *csrcRec, *cdstRec)
}

// Draw text (default font) within an image (destination)
func ImageDrawText(dst *Image, position Vector2, text string, fontSize int32, color Color) {
	cdst := dst.cptr()
	cposition := position.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.int)(fontSize)
	ccolor := color.cptr()
	C.ImageDrawText(cdst, *cposition, ctext, cfontSize, *ccolor)
}

// Draw text (custom sprite font) within an image (destination)
func ImageDrawTextEx(dst *Image, position Vector2, font SpriteFont, text string, fontSize float32, spacing int32, color Color) {
	cdst := dst.cptr()
	cposition := position.cptr()
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.float)(fontSize)
	cspacing := (C.int)(spacing)
	ccolor := color.cptr()
	C.ImageDrawTextEx(cdst, *cposition, *cfont, ctext, cfontSize, cspacing, *ccolor)
}

// Flip image vertically
func ImageFlipVertical(image *Image) {
	cimage := image.cptr()
	C.ImageFlipVertical(cimage)
}

// Flip image horizontally
func ImageFlipHorizontal(image *Image) {
	cimage := image.cptr()
	C.ImageFlipHorizontal(cimage)
}

// Modify image color: tint
func ImageColorTint(image *Image, color Color) {
	cimage := image.cptr()
	ccolor := color.cptr()
	C.ImageColorTint(cimage, *ccolor)
}

// Modify image color: invert
func ImageColorInvert(image *Image) {
	cimage := image.cptr()
	C.ImageColorInvert(cimage)
}

// Modify image color: grayscale
func ImageColorGrayscale(image *Image) {
	cimage := image.cptr()
	C.ImageColorGrayscale(cimage)
}

// Modify image color: contrast (-100 to 100)
func ImageColorContrast(image *Image, contrast float32) {
	cimage := image.cptr()
	ccontrast := (C.float)(contrast)
	C.ImageColorContrast(cimage, ccontrast)
}

// Modify image color: brightness (-255 to 255)
func ImageColorBrightness(image *Image, brightness int32) {
	cimage := image.cptr()
	cbrightness := (C.int)(brightness)
	C.ImageColorBrightness(cimage, cbrightness)
}

// Generate GPU mipmaps for a texture
func GenTextureMipmaps(texture *Texture2D) {
	ctexture := texture.cptr()
	C.GenTextureMipmaps(ctexture)
}

// Set texture scaling filter mode
func SetTextureFilter(texture Texture2D, filterMode TextureFilterMode) {
	ctexture := texture.cptr()
	cfilterMode := (C.int)(filterMode)
	C.SetTextureFilter(*ctexture, cfilterMode)
}

// Set texture wrapping mode
func SetTextureWrap(texture Texture2D, wrapMode TextureWrapMode) {
	ctexture := texture.cptr()
	cwrapMode := (C.int)(wrapMode)
	C.SetTextureWrap(*ctexture, cwrapMode)
}

// Draw a Texture2D
func DrawTexture(texture Texture2D, posX int32, posY int32, tint Color) {
	ctexture := texture.cptr()
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	ctint := tint.cptr()
	C.DrawTexture(*ctexture, cposX, cposY, *ctint)
}

// Draw a Texture2D with position defined as Vector2
func DrawTextureV(texture Texture2D, position Vector2, tint Color) {
	ctexture := texture.cptr()
	cposition := position.cptr()
	ctint := tint.cptr()
	C.DrawTextureV(*ctexture, *cposition, *ctint)
}

// Draw a Texture2D with extended parameters
func DrawTextureEx(texture Texture2D, position Vector2, rotation float32, scale float32, tint Color) {
	ctexture := texture.cptr()
	cposition := position.cptr()
	crotation := (C.float)(rotation)
	cscale := (C.float)(scale)
	ctint := tint.cptr()
	C.DrawTextureEx(*ctexture, *cposition, crotation, cscale, *ctint)
}

// Draw a part of a texture defined by a rectangle
func DrawTextureRec(texture Texture2D, sourceRec Rectangle, position Vector2, tint Color) {
	ctexture := texture.cptr()
	csourceRec := sourceRec.cptr()
	cposition := position.cptr()
	ctint := tint.cptr()
	C.DrawTextureRec(*ctexture, *csourceRec, *cposition, *ctint)
}

// Draw a part of a texture defined by a rectangle with 'pro' parameters
func DrawTexturePro(texture Texture2D, sourceRec Rectangle, destRec Rectangle, origin Vector2, rotation float32, tint Color) {
	ctexture := texture.cptr()
	csourceRec := sourceRec.cptr()
	cdestRec := destRec.cptr()
	corigin := origin.cptr()
	crotation := (C.float)(rotation)
	ctint := tint.cptr()
	C.DrawTexturePro(*ctexture, *csourceRec, *cdestRec, *corigin, crotation, *ctint)
}
