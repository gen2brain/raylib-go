//go:build !cgo && windows
// +build !cgo,windows

package rl

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

var rlMatrixMode func(mode int32)
var rlPushMatrix func()
var rlPopMatrix func()
var rlLoadIdentity func()
var rlTranslatef func(x float32, y float32, z float32)
var rlRotatef func(angle float32, x float32, y float32, z float32)
var rlScalef func(x float32, y float32, z float32)
var rlMultMatrixf func(matf *float32)
var rlFrustum func(left float64, right float64, bottom float64, top float64, znear float64, zfar float64)
var rlOrtho func(left float64, right float64, bottom float64, top float64, znear float64, zfar float64)
var rlViewport func(x int32, y int32, width int32, height int32)
var rlSetClipPlanes func(float64, float64)
var rlGetCullDistanceNear func() float64
var rlGetCullDistanceFar func() float64
var rlBegin func(mode int32)
var rlEnd func()
var rlVertex2i func(x int32, y int32)
var rlVertex2f func(x float32, y float32)
var rlVertex3f func(x float32, y float32, z float32)
var rlTexCoord2f func(x float32, y float32)
var rlNormal3f func(x float32, y float32, z float32)
var rlColor4ub func(r byte, g byte, b byte, a byte)
var rlColor3f func(x float32, y float32, z float32)
var rlColor4f func(x float32, y float32, z float32, w float32)
var rlEnableVertexArray func(vaoId uint32) bool
var rlDisableVertexArray func()
var rlEnableVertexBuffer func(id uint32)
var rlDisableVertexBuffer func()
var rlEnableVertexBufferElement func(id uint32)
var rlDisableVertexBufferElement func()
var rlEnableVertexAttribute func(index uint32)
var rlDisableVertexAttribute func(index uint32)
var rlActiveTextureSlot func(slot int32)
var rlEnableTexture func(id uint32)
var rlDisableTexture func()
var rlEnableTextureCubemap func(id uint32)
var rlDisableTextureCubemap func()
var rlTextureParameters func(id uint32, param int32, value int32)
var rlCubemapParameters func(id uint32, param int32, value int32)
var rlEnableShader func(id uint32)
var rlDisableShader func()
var rlEnableFramebuffer func(id uint32)
var rlDisableFramebuffer func()
var rlGetActiveFramebuffer func() uint32
var rlActiveDrawBuffers func(count int32)
var rlBlitFramebuffer func(srcX, srcY, srcWidth, srcHeight, dstX, dstY, dstWidth, dstHeight, bufferMask int32)
var rlBindFramebuffer func(target, framebuffer uint32)
var rlEnableColorBlend func()
var rlDisableColorBlend func()
var rlEnableDepthTest func()
var rlDisableDepthTest func()
var rlEnableDepthMask func()
var rlDisableDepthMask func()
var rlEnableBackfaceCulling func()
var rlDisableBackfaceCulling func()
var rlColorMask func(r, g, b, a bool)
var rlSetCullFace func(mode int32)
var rlEnableScissorTest func()
var rlDisableScissorTest func()
var rlScissor func(x int32, y int32, width int32, height int32)
var rlEnableWireMode func()
var rlEnablePointMode func()
var rlDisableWireMode func()
var rlSetLineWidth func(width float32)
var rlGetLineWidth func() float32
var rlEnableSmoothLines func()
var rlDisableSmoothLines func()
var rlEnableStereoRender func()
var rlDisableStereoRender func()
var rlIsStereoRenderEnabled func() bool
var rlClearColor func(r byte, g byte, b byte, a byte)
var rlClearScreenBuffers func()
var rlCheckErrors func()
var rlSetBlendMode func(mode int32)
var rlSetBlendFactors func(glSrcFactor int32, glDstFactor int32, glEquation int32)
var rlSetBlendFactorsSeparate func(glSrcRGB int32, glDstRGB int32, glSrcAlpha int32, glDstAlpha int32, glEqRGB int32, glEqAlpha int32)
var rlglInit func(width int32, height int32)
var rlglClose func()
var rlGetVersion func() int32
var rlSetFramebufferWidth func(width int32)
var rlGetFramebufferWidth func() int32
var rlSetFramebufferHeight func(height int32)
var rlGetFramebufferHeight func() int32
var rlGetTextureIdDefault func() uint32
var rlGetShaderIdDefault func() uint32
var rlLoadRenderBatch func(batch uintptr, numBuffers int32, bufferElements int32)
var rlUnloadRenderBatch func(batch uintptr)
var rlDrawRenderBatch func(batch *RenderBatch)
var rlSetRenderBatchActive func(batch *RenderBatch)
var rlDrawRenderBatchActive func()
var rlCheckRenderBatchLimit func(vCount int32) bool
var rlSetTexture func(id uint32)
var rlLoadVertexArray func() uint32
var rlUnloadVertexBuffer func(vboId uint32)
var rlSetVertexAttributeDivisor func(index uint32, divisor int32)
var rlLoadTextureDepth func(width int32, height int32, useRenderBuffer bool) uint32
var rlLoadFramebuffer func() uint32
var rlFramebufferAttach func(fboId uint32, texId uint32, attachType int32, texType int32, mipLevel int32)
var rlFramebufferComplete func(id uint32) bool
var rlUnloadFramebuffer func(id uint32)
var rlLoadShaderCode func(vsCode string, fsCode string) uint32
var rlCompileShader func(shaderCode string, _type int32) uint32
var rlLoadShaderProgram func(vShaderId uint32, fShaderId uint32) uint32
var rlUnloadShaderProgram func(id uint32)
var rlGetLocationUniform func(shaderId uint32, uniformName string) int32
var rlGetLocationAttrib func(shaderId uint32, attribName string) int32
var rlSetUniform func(locIndex int32, value unsafe.Pointer, uniformType, count int32)
var rlSetUniformMatrix func(locIndex int32, mat uintptr)
var rlSetUniformMatrices func(locIndex int32, mat *Matrix, count int32)
var rlSetUniformSampler func(locIndex int32, textureId uint32)
var rlLoadComputeShaderProgram func(shaderID uint32) uint32
var rlComputeShaderDispatch func(groupX uint32, groupY uint32, groupZ uint32)
var rlLoadShaderBuffer func(size uint32, data unsafe.Pointer, usageHint int32) uint32
var rlUnloadShaderBuffer func(id uint32)
var rlUpdateShaderBuffer func(id uint32, data unsafe.Pointer, dataSize uint32, offset uint32)
var rlBindShaderBuffer func(id uint32, index uint32)
var rlReadShaderBuffer func(id uint32, dest unsafe.Pointer, count uint32, offset uint32)
var rlCopyShaderBuffer func(destId uint32, srcId uint32, destOffset uint32, srcOffset uint32, count uint32)
var rlGetShaderBufferSize func(id uint32) uint32
var rlBindImageTexture func(id uint32, index uint32, format int32, readonly bool)
var rlGetMatrixModelview func(matrix uintptr)
var rlGetMatrixProjection func(matrix uintptr)
var rlGetMatrixTransform func(matrix uintptr)
var rlGetMatrixProjectionStereo func(matrix uintptr, eye int32)
var rlGetMatrixViewOffsetStereo func(matrix uintptr, eye int32)
var rlSetMatrixProjection func(proj uintptr)
var rlSetMatrixModelview func(view uintptr)
var rlSetMatrixProjectionStereo func(right uintptr, left uintptr)
var rlSetMatrixViewOffsetStereo func(right uintptr, left uintptr)
var rlLoadDrawCube func()
var rlLoadDrawQuad func()

func initRlglPurego() {
	purego.RegisterLibFunc(&rlMatrixMode, raylibDll, "rlMatrixMode")
	purego.RegisterLibFunc(&rlPushMatrix, raylibDll, "rlPushMatrix")
	purego.RegisterLibFunc(&rlPopMatrix, raylibDll, "rlPopMatrix")
	purego.RegisterLibFunc(&rlLoadIdentity, raylibDll, "rlLoadIdentity")
	purego.RegisterLibFunc(&rlTranslatef, raylibDll, "rlTranslatef")
	purego.RegisterLibFunc(&rlRotatef, raylibDll, "rlRotatef")
	purego.RegisterLibFunc(&rlScalef, raylibDll, "rlScalef")
	purego.RegisterLibFunc(&rlMultMatrixf, raylibDll, "rlMultMatrixf")
	purego.RegisterLibFunc(&rlFrustum, raylibDll, "rlFrustum")
	purego.RegisterLibFunc(&rlOrtho, raylibDll, "rlOrtho")
	purego.RegisterLibFunc(&rlViewport, raylibDll, "rlViewport")
	purego.RegisterLibFunc(&rlSetClipPlanes, raylibDll, "rlSetClipPlanes")
	purego.RegisterLibFunc(&rlGetCullDistanceNear, raylibDll, "rlGetCullDistanceNear")
	purego.RegisterLibFunc(&rlGetCullDistanceFar, raylibDll, "rlGetCullDistanceFar")
	purego.RegisterLibFunc(&rlBegin, raylibDll, "rlBegin")
	purego.RegisterLibFunc(&rlEnd, raylibDll, "rlEnd")
	purego.RegisterLibFunc(&rlVertex2i, raylibDll, "rlVertex2i")
	purego.RegisterLibFunc(&rlVertex2f, raylibDll, "rlVertex2f")
	purego.RegisterLibFunc(&rlVertex3f, raylibDll, "rlVertex3f")
	purego.RegisterLibFunc(&rlTexCoord2f, raylibDll, "rlTexCoord2f")
	purego.RegisterLibFunc(&rlNormal3f, raylibDll, "rlNormal3f")
	purego.RegisterLibFunc(&rlColor4ub, raylibDll, "rlColor4ub")
	purego.RegisterLibFunc(&rlColor3f, raylibDll, "rlColor3f")
	purego.RegisterLibFunc(&rlColor4f, raylibDll, "rlColor4f")
	purego.RegisterLibFunc(&rlEnableVertexArray, raylibDll, "rlEnableVertexArray")
	purego.RegisterLibFunc(&rlDisableVertexArray, raylibDll, "rlDisableVertexArray")
	purego.RegisterLibFunc(&rlEnableVertexBuffer, raylibDll, "rlEnableVertexBuffer")
	purego.RegisterLibFunc(&rlDisableVertexBuffer, raylibDll, "rlDisableVertexBuffer")
	purego.RegisterLibFunc(&rlEnableVertexBufferElement, raylibDll, "rlEnableVertexBufferElement")
	purego.RegisterLibFunc(&rlDisableVertexBufferElement, raylibDll, "rlDisableVertexBufferElement")
	purego.RegisterLibFunc(&rlEnableVertexAttribute, raylibDll, "rlEnableVertexAttribute")
	purego.RegisterLibFunc(&rlDisableVertexAttribute, raylibDll, "rlDisableVertexAttribute")
	purego.RegisterLibFunc(&rlActiveTextureSlot, raylibDll, "rlActiveTextureSlot")
	purego.RegisterLibFunc(&rlEnableTexture, raylibDll, "rlEnableTexture")
	purego.RegisterLibFunc(&rlDisableTexture, raylibDll, "rlDisableTexture")
	purego.RegisterLibFunc(&rlEnableTextureCubemap, raylibDll, "rlEnableTextureCubemap")
	purego.RegisterLibFunc(&rlDisableTextureCubemap, raylibDll, "rlDisableTextureCubemap")
	purego.RegisterLibFunc(&rlTextureParameters, raylibDll, "rlTextureParameters")
	purego.RegisterLibFunc(&rlCubemapParameters, raylibDll, "rlCubemapParameters")
	purego.RegisterLibFunc(&rlEnableShader, raylibDll, "rlEnableShader")
	purego.RegisterLibFunc(&rlDisableShader, raylibDll, "rlDisableShader")
	purego.RegisterLibFunc(&rlEnableFramebuffer, raylibDll, "rlEnableFramebuffer")
	purego.RegisterLibFunc(&rlDisableFramebuffer, raylibDll, "rlDisableFramebuffer")
	purego.RegisterLibFunc(&rlGetActiveFramebuffer, raylibDll, "rlGetActiveFramebuffer")
	purego.RegisterLibFunc(&rlActiveDrawBuffers, raylibDll, "rlActiveDrawBuffers")
	purego.RegisterLibFunc(&rlBlitFramebuffer, raylibDll, "rlBlitFramebuffer")
	purego.RegisterLibFunc(&rlBindFramebuffer, raylibDll, "rlBindFramebuffer")
	purego.RegisterLibFunc(&rlEnableColorBlend, raylibDll, "rlEnableColorBlend")
	purego.RegisterLibFunc(&rlDisableColorBlend, raylibDll, "rlDisableColorBlend")
	purego.RegisterLibFunc(&rlEnableDepthTest, raylibDll, "rlEnableDepthTest")
	purego.RegisterLibFunc(&rlDisableDepthTest, raylibDll, "rlDisableDepthTest")
	purego.RegisterLibFunc(&rlEnableDepthMask, raylibDll, "rlEnableDepthMask")
	purego.RegisterLibFunc(&rlDisableDepthMask, raylibDll, "rlDisableDepthMask")
	purego.RegisterLibFunc(&rlEnableBackfaceCulling, raylibDll, "rlEnableBackfaceCulling")
	purego.RegisterLibFunc(&rlDisableBackfaceCulling, raylibDll, "rlDisableBackfaceCulling")
	purego.RegisterLibFunc(&rlColorMask, raylibDll, "rlColorMask")
	purego.RegisterLibFunc(&rlSetCullFace, raylibDll, "rlSetCullFace")
	purego.RegisterLibFunc(&rlEnableScissorTest, raylibDll, "rlEnableScissorTest")
	purego.RegisterLibFunc(&rlDisableScissorTest, raylibDll, "rlDisableScissorTest")
	purego.RegisterLibFunc(&rlScissor, raylibDll, "rlScissor")
	purego.RegisterLibFunc(&rlEnableWireMode, raylibDll, "rlEnableWireMode")
	purego.RegisterLibFunc(&rlEnablePointMode, raylibDll, "rlEnablePointMode")
	purego.RegisterLibFunc(&rlDisableWireMode, raylibDll, "rlDisableWireMode")
	purego.RegisterLibFunc(&rlSetLineWidth, raylibDll, "rlSetLineWidth")
	purego.RegisterLibFunc(&rlGetLineWidth, raylibDll, "rlGetLineWidth")
	purego.RegisterLibFunc(&rlEnableSmoothLines, raylibDll, "rlEnableSmoothLines")
	purego.RegisterLibFunc(&rlDisableSmoothLines, raylibDll, "rlDisableSmoothLines")
	purego.RegisterLibFunc(&rlEnableStereoRender, raylibDll, "rlEnableStereoRender")
	purego.RegisterLibFunc(&rlDisableStereoRender, raylibDll, "rlDisableStereoRender")
	purego.RegisterLibFunc(&rlIsStereoRenderEnabled, raylibDll, "rlIsStereoRenderEnabled")
	purego.RegisterLibFunc(&rlClearColor, raylibDll, "rlClearColor")
	purego.RegisterLibFunc(&rlClearScreenBuffers, raylibDll, "rlClearScreenBuffers")
	purego.RegisterLibFunc(&rlCheckErrors, raylibDll, "rlCheckErrors")
	purego.RegisterLibFunc(&rlSetBlendMode, raylibDll, "rlSetBlendMode")
	purego.RegisterLibFunc(&rlSetBlendFactors, raylibDll, "rlSetBlendFactors")
	purego.RegisterLibFunc(&rlSetBlendFactorsSeparate, raylibDll, "rlSetBlendFactorsSeparate")
	purego.RegisterLibFunc(&rlglInit, raylibDll, "rlglInit")
	purego.RegisterLibFunc(&rlglClose, raylibDll, "rlglClose")
	purego.RegisterLibFunc(&rlGetVersion, raylibDll, "rlGetVersion")
	purego.RegisterLibFunc(&rlSetFramebufferWidth, raylibDll, "rlSetFramebufferWidth")
	purego.RegisterLibFunc(&rlGetFramebufferWidth, raylibDll, "rlGetFramebufferWidth")
	purego.RegisterLibFunc(&rlSetFramebufferHeight, raylibDll, "rlSetFramebufferHeight")
	purego.RegisterLibFunc(&rlGetFramebufferHeight, raylibDll, "rlGetFramebufferHeight")
	purego.RegisterLibFunc(&rlGetTextureIdDefault, raylibDll, "rlGetTextureIdDefault")
	purego.RegisterLibFunc(&rlGetShaderIdDefault, raylibDll, "rlGetShaderIdDefault")
	purego.RegisterLibFunc(&rlLoadRenderBatch, raylibDll, "rlLoadRenderBatch")
	purego.RegisterLibFunc(&rlUnloadRenderBatch, raylibDll, "rlUnloadRenderBatch")
	purego.RegisterLibFunc(&rlDrawRenderBatch, raylibDll, "rlDrawRenderBatch")
	purego.RegisterLibFunc(&rlSetRenderBatchActive, raylibDll, "rlSetRenderBatchActive")
	purego.RegisterLibFunc(&rlDrawRenderBatchActive, raylibDll, "rlDrawRenderBatchActive")
	purego.RegisterLibFunc(&rlCheckRenderBatchLimit, raylibDll, "rlCheckRenderBatchLimit")
	purego.RegisterLibFunc(&rlSetTexture, raylibDll, "rlSetTexture")
	purego.RegisterLibFunc(&rlLoadVertexArray, raylibDll, "rlLoadVertexArray")
	purego.RegisterLibFunc(&rlUnloadVertexBuffer, raylibDll, "rlUnloadVertexBuffer")
	purego.RegisterLibFunc(&rlSetVertexAttributeDivisor, raylibDll, "rlSetVertexAttributeDivisor")
	purego.RegisterLibFunc(&rlLoadTextureDepth, raylibDll, "rlLoadTextureDepth")
	purego.RegisterLibFunc(&rlLoadFramebuffer, raylibDll, "rlLoadFramebuffer")
	purego.RegisterLibFunc(&rlFramebufferAttach, raylibDll, "rlFramebufferAttach")
	purego.RegisterLibFunc(&rlFramebufferComplete, raylibDll, "rlFramebufferComplete")
	purego.RegisterLibFunc(&rlUnloadFramebuffer, raylibDll, "rlUnloadFramebuffer")
	purego.RegisterLibFunc(&rlLoadShaderCode, raylibDll, "rlLoadShaderCode")
	purego.RegisterLibFunc(&rlCompileShader, raylibDll, "rlCompileShader")
	purego.RegisterLibFunc(&rlLoadShaderProgram, raylibDll, "rlLoadShaderProgram")
	purego.RegisterLibFunc(&rlUnloadShaderProgram, raylibDll, "rlUnloadShaderProgram")
	purego.RegisterLibFunc(&rlGetLocationUniform, raylibDll, "rlGetLocationUniform")
	purego.RegisterLibFunc(&rlGetLocationAttrib, raylibDll, "rlGetLocationAttrib")
	purego.RegisterLibFunc(&rlSetUniform, raylibDll, "rlSetUniform")
	purego.RegisterLibFunc(&rlSetUniformMatrix, raylibDll, "rlSetUniformMatrix")
	purego.RegisterLibFunc(&rlSetUniformMatrices, raylibDll, "rlSetUniformMatrices")
	purego.RegisterLibFunc(&rlSetUniformSampler, raylibDll, "rlSetUniformSampler")
	purego.RegisterLibFunc(&rlLoadComputeShaderProgram, raylibDll, "rlLoadComputeShaderProgram")
	purego.RegisterLibFunc(&rlComputeShaderDispatch, raylibDll, "rlComputeShaderDispatch")
	purego.RegisterLibFunc(&rlLoadShaderBuffer, raylibDll, "rlLoadShaderBuffer")
	purego.RegisterLibFunc(&rlUnloadShaderBuffer, raylibDll, "rlUnloadShaderBuffer")
	purego.RegisterLibFunc(&rlUpdateShaderBuffer, raylibDll, "rlUpdateShaderBuffer")
	purego.RegisterLibFunc(&rlBindShaderBuffer, raylibDll, "rlBindShaderBuffer")
	purego.RegisterLibFunc(&rlReadShaderBuffer, raylibDll, "rlReadShaderBuffer")
	purego.RegisterLibFunc(&rlCopyShaderBuffer, raylibDll, "rlCopyShaderBuffer")
	purego.RegisterLibFunc(&rlGetShaderBufferSize, raylibDll, "rlGetShaderBufferSize")
	purego.RegisterLibFunc(&rlBindImageTexture, raylibDll, "rlBindImageTexture")
	purego.RegisterLibFunc(&rlGetMatrixModelview, raylibDll, "rlGetMatrixModelview")
	purego.RegisterLibFunc(&rlGetMatrixProjection, raylibDll, "rlGetMatrixProjection")
	purego.RegisterLibFunc(&rlGetMatrixTransform, raylibDll, "rlGetMatrixTransform")
	purego.RegisterLibFunc(&rlGetMatrixProjectionStereo, raylibDll, "rlGetMatrixProjectionStereo")
	purego.RegisterLibFunc(&rlGetMatrixViewOffsetStereo, raylibDll, "rlGetMatrixViewOffsetStereo")
	purego.RegisterLibFunc(&rlSetMatrixProjection, raylibDll, "rlSetMatrixProjection")
	purego.RegisterLibFunc(&rlSetMatrixModelview, raylibDll, "rlSetMatrixModelview")
	purego.RegisterLibFunc(&rlSetMatrixProjectionStereo, raylibDll, "rlSetMatrixProjectionStereo")
	purego.RegisterLibFunc(&rlSetMatrixViewOffsetStereo, raylibDll, "rlSetMatrixViewOffsetStereo")
	purego.RegisterLibFunc(&rlLoadDrawCube, raylibDll, "rlLoadDrawCube")
	purego.RegisterLibFunc(&rlLoadDrawQuad, raylibDll, "rlLoadDrawQuad")
}

// SetMatrixProjection - Set a custom projection matrix (replaces internal projection matrix)
func SetMatrixProjection(proj Matrix) {
	rlSetMatrixProjection(uintptr(unsafe.Pointer(&proj)))
}

// SetMatrixModelview - Set a custom modelview matrix (replaces internal modelview matrix)
func SetMatrixModelview(view Matrix) {
	rlSetMatrixModelview(uintptr(unsafe.Pointer(&view)))
}

// MatrixMode - Choose the current matrix to be transformed
func MatrixMode(mode int32) {
	rlMatrixMode(mode)
}

// PushMatrix - Push the current matrix to stack
func PushMatrix() {
	rlPushMatrix()
}

// PopMatrix - Pop lattest inserted matrix from stack
func PopMatrix() {
	rlPopMatrix()
}

// LoadIdentity - Reset current matrix to identity matrix
func LoadIdentity() {
	rlLoadIdentity()
}

// Translatef - Multiply the current matrix by a translation matrix
func Translatef(x float32, y float32, z float32) {
	rlTranslatef(x, y, z)
}

// Rotatef - Multiply the current matrix by a rotation matrix
func Rotatef(angle float32, x float32, y float32, z float32) {
	rlRotatef(angle, x, y, z)
}

// Scalef - Multiply the current matrix by a scaling matrix
func Scalef(x float32, y float32, z float32) {
	rlScalef(x, y, z)
}

// MultMatrix - Multiply the current matrix by another matrix
func MultMatrix(m Matrix) {
	f := MatrixToFloat(m)
	rlMultMatrixf((*float32)(&f[0]))
}

// Frustum .
func Frustum(left float64, right float64, bottom float64, top float64, znear float64, zfar float64) {
	rlFrustum(left, right, bottom, top, znear, zfar)
}

// Ortho .
func Ortho(left float64, right float64, bottom float64, top float64, znear float64, zfar float64) {
	rlOrtho(left, right, bottom, top, znear, zfar)
}

// Viewport - Set the viewport area
func Viewport(x int32, y int32, width int32, height int32) {
	rlViewport(x, y, width, height)
}

// SetClipPlanes - Set clip planes distances
func SetClipPlanes(nearPlane, farPlane float64) {
	rlSetClipPlanes(nearPlane, farPlane)
}

// GetCullDistanceNear - Get cull plane distance near
func GetCullDistanceNear() float64 {
	return rlGetCullDistanceNear()
}

// GetCullDistanceFar - Get cull plane distance far
func GetCullDistanceFar() float64 {
	return rlGetCullDistanceFar()
}

// Begin - Initialize drawing mode (how to organize vertex)
func Begin(mode int32) {
	rlBegin(mode)
}

// End - Finish vertex providing
func End() {
	rlEnd()
}

// Vertex2i - Define one vertex (position) - 2 int
func Vertex2i(x int32, y int32) {
	rlVertex2i(x, y)
}

// Vertex2f - Define one vertex (position) - 2 float
func Vertex2f(x float32, y float32) {
	rlVertex2f(x, y)
}

// Vertex3f - Define one vertex (position) - 3 float
func Vertex3f(x float32, y float32, z float32) {
	rlVertex3f(x, y, z)
}

// TexCoord2f - Define one vertex (texture coordinate) - 2 float
func TexCoord2f(x float32, y float32) {
	rlTexCoord2f(x, y)
}

// Normal3f - Define one vertex (normal) - 3 float
func Normal3f(x float32, y float32, z float32) {
	rlNormal3f(x, y, z)
}

// Color4ub - Define one vertex (color) - 4 byte
func Color4ub(r uint8, g uint8, b uint8, a uint8) {
	rlColor4ub(r, g, b, a)
}

// Color3f - Define one vertex (color) - 3 float
func Color3f(x float32, y float32, z float32) {
	rlColor3f(x, y, z)
}

// Color4f - Define one vertex (color) - 4 float
func Color4f(x float32, y float32, z float32, w float32) {
	rlColor4f(x, y, z, w)
}

// EnableVertexArray - Enable vertex array (VAO, if supported)
func EnableVertexArray(vaoId uint32) bool {
	return rlEnableVertexArray(vaoId)
}

// DisableVertexArray - Disable vertex array (VAO, if supported)
func DisableVertexArray() {
	rlDisableVertexArray()
}

// EnableVertexBuffer - Enable vertex buffer (VBO)
func EnableVertexBuffer(id uint32) {
	rlEnableVertexBuffer(id)
}

// DisableVertexBuffer - Disable vertex buffer (VBO)
func DisableVertexBuffer() {
	rlDisableVertexBuffer()
}

// EnableVertexBufferElement - Enable vertex buffer element (VBO element)
func EnableVertexBufferElement(id uint32) {
	rlEnableVertexBufferElement(id)
}

// DisableVertexBufferElement - Disable vertex buffer element (VBO element)
func DisableVertexBufferElement() {
	rlDisableVertexBufferElement()
}

// EnableVertexAttribute - Enable vertex attribute index
func EnableVertexAttribute(index uint32) {
	rlEnableVertexAttribute(index)
}

// DisableVertexAttribute - Disable vertex attribute index
func DisableVertexAttribute(index uint32) {
	rlDisableVertexAttribute(index)
}

// ActiveTextureSlot - Select and active a texture slot
func ActiveTextureSlot(slot int32) {
	rlActiveTextureSlot(slot)
}

// EnableTexture - Enable texture
func EnableTexture(id uint32) {
	rlEnableTexture(id)
}

// DisableTexture - Disable texture
func DisableTexture() {
	rlDisableTexture()
}

// EnableTextureCubemap - Enable texture cubemap
func EnableTextureCubemap(id uint32) {
	rlEnableTextureCubemap(id)
}

// DisableTextureCubemap - Disable texture cubemap
func DisableTextureCubemap() {
	rlDisableTextureCubemap()
}

// TextureParameters - Set texture parameters (filter, wrap)
func TextureParameters(id uint32, param int32, value int32) {
	rlTextureParameters(id, param, value)
}

// CubemapParameters - Set cubemap parameters (filter, wrap)
func CubemapParameters(id uint32, param int32, value int32) {
	rlCubemapParameters(id, param, value)
}

// EnableShader - Enable shader program
func EnableShader(id uint32) {
	rlEnableShader(id)
}

// DisableShader - Disable shader program
func DisableShader() {
	rlDisableShader()
}

// EnableFramebuffer - Enable render texture (fbo)
func EnableFramebuffer(id uint32) {
	rlEnableFramebuffer(id)
}

// DisableFramebuffer - Disable render texture (fbo), return to default framebuffer
func DisableFramebuffer() {
	rlDisableFramebuffer()
}

// GetActiveFramebuffer - Get the currently active render texture (fbo), 0 for default framebuffer
func GetActiveFramebuffer() uint32 {
	return rlGetActiveFramebuffer()
}

// ActiveDrawBuffers - Activate multiple draw color buffers
func ActiveDrawBuffers(count int32) {
	rlActiveDrawBuffers(count)
}

// BlitFramebuffer - Blit active framebuffer to main framebuffer
func BlitFramebuffer(srcX, srcY, srcWidth, srcHeight, dstX, dstY, dstWidth, dstHeight, bufferMask int32) {
	rlBlitFramebuffer(srcX, srcY, srcWidth, srcHeight, dstX, dstY, dstWidth, dstHeight, bufferMask)
}

// BindFramebuffer - Bind framebuffer (FBO)
func BindFramebuffer(target, framebuffer uint32) {
	rlBindFramebuffer(target, framebuffer)
}

// EnableColorBlend - Enable color blending
func EnableColorBlend() {
	rlEnableColorBlend()
}

// DisableColorBlend - Disable color blending
func DisableColorBlend() {
	rlDisableColorBlend()
}

// EnableDepthTest - Enable depth test
func EnableDepthTest() {
	rlEnableDepthTest()
}

// DisableDepthTest - Disable depth test
func DisableDepthTest() {
	rlDisableDepthTest()
}

// EnableDepthMask - Enable depth write
func EnableDepthMask() {
	rlEnableDepthMask()
}

// DisableDepthMask - Disable depth write
func DisableDepthMask() {
	rlDisableDepthMask()
}

// EnableBackfaceCulling - Enable backface culling
func EnableBackfaceCulling() {
	rlEnableBackfaceCulling()
}

// DisableBackfaceCulling - Disable backface culling
func DisableBackfaceCulling() {
	rlDisableBackfaceCulling()
}

// ColorMask - Color mask control
func ColorMask(r, g, b, a bool) {
	rlColorMask(r, g, b, a)
}

// SetCullFace - Set face culling mode
func SetCullFace(mode int32) {
	rlSetCullFace(mode)
}

// EnableScissorTest - Enable scissor test
func EnableScissorTest() {
	rlEnableScissorTest()
}

// DisableScissorTest - Disable scissor test
func DisableScissorTest() {
	rlDisableScissorTest()
}

// Scissor - Scissor test
func Scissor(x int32, y int32, width int32, height int32) {
	rlScissor(x, y, width, height)
}

// EnableWireMode - Enable wire mode
func EnableWireMode() {
	rlEnableWireMode()
}

// EnablePointMode - Enable point mode
func EnablePointMode() {
	rlEnablePointMode()
}

// DisableWireMode - Disable wire mode
func DisableWireMode() {
	rlDisableWireMode()
}

// SetLineWidth - Set the line drawing width
func SetLineWidth(width float32) {
	rlSetLineWidth(width)
}

// GetLineWidth - Get the line drawing width
func GetLineWidth() float32 {
	return rlGetLineWidth()
}

// EnableSmoothLines - Enable line aliasing
func EnableSmoothLines() {
	rlEnableSmoothLines()
}

// DisableSmoothLines - Disable line aliasing
func DisableSmoothLines() {
	rlDisableSmoothLines()
}

// EnableStereoRender - Enable stereo rendering
func EnableStereoRender() {
	rlEnableStereoRender()
}

// DisableStereoRender - Disable stereo rendering
func DisableStereoRender() {
	rlDisableStereoRender()
}

// IsStereoRenderEnabled - Check if stereo render is enabled
func IsStereoRenderEnabled() bool {
	return rlIsStereoRenderEnabled()
}

// ClearColor - Clear color buffer with color
func ClearColor(r uint8, g uint8, b uint8, a uint8) {
	rlClearColor(r, g, b, a)
}

// ClearScreenBuffers - Clear used screen buffers (color and depth)
func ClearScreenBuffers() {
	rlClearScreenBuffers()
}

// CheckErrors - Check and log OpenGL error codes
func CheckErrors() {
	rlCheckErrors()
}

// SetBlendMode - Set blending mode
func SetBlendMode(mode BlendMode) {
	rlSetBlendMode(int32(mode))
}

// SetBlendFactors - Set blending mode factor and equation (using OpenGL factors)
func SetBlendFactors(glSrcFactor int32, glDstFactor int32, glEquation int32) {
	rlSetBlendFactors(glSrcFactor, glDstFactor, glEquation)
}

// SetBlendFactorsSeparate - Set blending mode factors and equations separately (using OpenGL factors)
func SetBlendFactorsSeparate(glSrcRGB int32, glDstRGB int32, glSrcAlpha int32, glDstAlpha int32, glEqRGB int32, glEqAlpha int32) {
	rlSetBlendFactorsSeparate(glSrcRGB, glDstRGB, glSrcAlpha, glDstAlpha, glEqRGB, glEqAlpha)
}

// GlInit - Initialize rlgl (buffers, shaders, textures, states)
func GlInit(width int32, height int32) {
	rlglInit(width, height)
}

// GlClose - De-inititialize rlgl (buffers, shaders, textures)
func GlClose() {
	rlglClose()
}

// GetVersion - Get current OpenGL version
func GetVersion() int32 {
	return rlGetVersion()
}

// SetFramebufferWidth - Set current framebuffer width
func SetFramebufferWidth(width int32) {
	rlSetFramebufferWidth(width)
}

// GetFramebufferWidth - Get default framebuffer width
func GetFramebufferWidth() int32 {
	return rlGetFramebufferWidth()
}

// SetFramebufferHeight - Set current framebuffer height
func SetFramebufferHeight(height int32) {
	rlSetFramebufferHeight(height)
}

// GetFramebufferHeight - Get default framebuffer height
func GetFramebufferHeight() int32 {
	return rlGetFramebufferHeight()
}

// GetTextureIdDefault - Get default texture id
func GetTextureIdDefault() uint32 {
	return rlGetTextureIdDefault()
}

// GetShaderIdDefault - Get default shader id
func GetShaderIdDefault() uint32 {
	return rlGetShaderIdDefault()
}

// LoadRenderBatch - Load a render batch system
func LoadRenderBatch(numBuffers int32, bufferElements int32) RenderBatch {
	var batch RenderBatch
	rlLoadRenderBatch(uintptr(unsafe.Pointer(&batch)), numBuffers, bufferElements)
	return batch
}

// UnloadRenderBatch - Unload render batch system
func UnloadRenderBatch(batch RenderBatch) {
	rlUnloadRenderBatch(uintptr(unsafe.Pointer(&batch)))
}

// DrawRenderBatch - Draw render batch data (Update->Draw->Reset)
func DrawRenderBatch(batch *RenderBatch) {
	rlDrawRenderBatch(batch)
}

// rlSetRenderBatchActive - Set the active render batch for rlgl (NULL for default internal)
func SetRenderBatchActive(batch *RenderBatch) {
	rlSetRenderBatchActive(batch)
}

// DrawRenderBatchActive - Update and draw internal render batch
func DrawRenderBatchActive() {
	rlDrawRenderBatchActive()
}

// CheckRenderBatchLimit - Check internal buffer overflow for a given number of vertex
func CheckRenderBatchLimit(vCount int32) bool {
	return rlCheckRenderBatchLimit(vCount)
}

// SetTexture - Set current texture for render batch and check buffers limits
func SetTexture(id uint32) {
	rlSetTexture(id)
}

// LoadVertexArray - Load vertex array (vao) if supported
func LoadVertexArray() uint32 {
	return rlLoadVertexArray()
}

// UnloadVertexBuffer .
func UnloadVertexBuffer(vboId uint32) {
	rlUnloadVertexBuffer(vboId)
}

// SetVertexAttributeDivisor .
func SetVertexAttributeDivisor(index uint32, divisor int32) {
	rlSetVertexAttributeDivisor(index, divisor)
}

// LoadTextureDepth - Load depth texture/renderbuffer (to be attached to fbo)
func LoadTextureDepth(width, height int32, useRenderBuffer bool) uint32{
	return rlLoadTextureDepth(width, height, useRenderBuffer)
}

// LoadFramebuffer - Load an empty framebuffer
func LoadFramebuffer() uint32 {
	return rlLoadFramebuffer()
}

// FramebufferAttach - Attach texture/renderbuffer to a framebuffer
func FramebufferAttach(fboId uint32, texId uint32, attachType int32, texType int32, mipLevel int32) {
	rlFramebufferAttach(fboId, texId, attachType, texType, mipLevel)
}

// FramebufferComplete - Verify framebuffer is complete
func FramebufferComplete(id uint32) bool {
	return rlFramebufferComplete(id)
}

// UnloadFramebuffer - Delete framebuffer from GPU
func UnloadFramebuffer(id uint32) {
	rlUnloadFramebuffer(id)
}

// LoadShaderCode - Load shader from code strings
func LoadShaderCode(vsCode string, fsCode string) uint32 {
	return rlLoadShaderCode(vsCode, fsCode)
}

// CompileShader - Compile custom shader and return shader id (type: VERTEX_SHADER, FRAGMENT_SHADER, COMPUTE_SHADER)
func CompileShader(shaderCode string, type_ int32) uint32 {
	return rlCompileShader(shaderCode, type_)
}

// LoadShaderProgram - Load custom shader program
func LoadShaderProgram(vShaderId uint32, fShaderId uint32) uint32 {
	return rlLoadShaderProgram(vShaderId, fShaderId)
}

// UnloadShaderProgram - Unload shader program
func UnloadShaderProgram(id uint32) {
	rlUnloadShaderProgram(id)
}

// GetLocationUniform - Get shader location uniform
func GetLocationUniform(shaderId uint32, uniformName string) int32 {
	return rlGetLocationUniform(shaderId, uniformName)
}

// GetLocationAttrib - Get shader location attribute
func GetLocationAttrib(shaderId uint32, attribName string) int32 {
	return rlGetLocationAttrib(shaderId, attribName)
}

// SetUniform - Set shader value uniform
func SetUniform(locIndex int32, value []float32, uniformType int32) {
	rlSetUniform(locIndex, unsafe.Pointer(unsafe.SliceData(value)), uniformType, int32(len(value)))
}

// SetUniformMatrix - Set shader value matrix
func SetUniformMatrix(locIndex int32, mat Matrix) {
	rlSetUniformMatrix(locIndex, uintptr(unsafe.Pointer(&mat)))
}

// SetUniformMatrices - Set shader value matrices
func SetUniformMatrices(locIndex int32, mat []Matrix) {
	rlSetUniformMatrices(locIndex, unsafe.SliceData(mat), int32(len(mat)))
}

// SetUniformSampler - Set shader value sampler
func SetUniformSampler(locIndex int32, textureId uint32) {
	rlSetUniformSampler(locIndex, textureId)
}

// LoadComputeShaderProgram - Load compute shader program
func LoadComputeShaderProgram(shaderID uint32) uint32 {
	return rlLoadComputeShaderProgram(shaderID)
}

// ComputeShaderDispatch - Dispatch compute shader (equivalent to *draw* for graphics pilepine)
func ComputeShaderDispatch(groupX uint32, groupY uint32, groupZ uint32) {
	rlComputeShaderDispatch(groupX, groupY, groupZ)
}

// LoadShaderBuffer loads a shader storage buffer object (SSBO)
func LoadShaderBuffer(size uint32, data unsafe.Pointer, usageHint int32) uint32 {
	return rlLoadShaderBuffer(size, data, usageHint)
}

// UnloadShaderBuffer - Unload shader storage buffer object (SSBO)
func UnloadShaderBuffer(id uint32) {
	rlUnloadShaderBuffer(id)
}

// UpdateShaderBuffer - Update SSBO buffer data
func UpdateShaderBuffer(id uint32, data unsafe.Pointer, dataSize uint32, offset uint32) {
	rlUpdateShaderBuffer(id, data, dataSize, offset)
}

// BindShaderBuffer - Bind SSBO buffer
func BindShaderBuffer(id uint32, index uint32) {
	rlBindShaderBuffer(id, index)
}

// ReadShaderBuffer - Read SSBO buffer data (GPU->CPU)
func ReadShaderBuffer(id uint32, dest unsafe.Pointer, count uint32, offset uint32) {
	rlReadShaderBuffer(id, dest, count, offset)
}

// CopyShaderBuffer - Copy SSBO data between buffers
func CopyShaderBuffer(destId uint32, srcId uint32, destOffset uint32, srcOffset uint32, count uint32) {
	rlCopyShaderBuffer(destId, srcId, destOffset, srcOffset, count)
}

// GetShaderBufferSize - Get SSBO buffer size
func GetShaderBufferSize(id uint32) uint32 {
	return rlGetShaderBufferSize(id)
}

// BindImageTexture - Bind image texture
func BindImageTexture(id uint32, index uint32, format int32, readonly bool) {
	rlBindImageTexture(id, index, format, readonly)
}

// GetMatrixModelview - Get internal modelview matrix
func GetMatrixModelview() Matrix {
	var matrix Matrix
	rlGetMatrixModelview(uintptr(unsafe.Pointer(&matrix)))
	return matrix
}

// GetMatrixProjection - Get internal projection matrix
func GetMatrixProjection() Matrix {
	var matrix Matrix
	rlGetMatrixProjection(uintptr(unsafe.Pointer(&matrix)))
	return matrix
}

// GetMatrixTransform - Get internal accumulated transform matrix
func GetMatrixTransform() Matrix {
	var matrix Matrix
	rlGetMatrixTransform(uintptr(unsafe.Pointer(&matrix)))
	return matrix
}

// GetMatrixProjectionStereo - Get internal projection matrix for stereo render (selected eye)
func GetMatrixProjectionStereo(eye int32) Matrix {
	var matrix Matrix
	rlGetMatrixProjectionStereo(uintptr(unsafe.Pointer(&matrix)), eye)
	return matrix
}

// GetMatrixViewOffsetStereo - Get internal view offset matrix for stereo render (selected eye)
func GetMatrixViewOffsetStereo(eye int32) Matrix {
	var matrix Matrix
	rlGetMatrixViewOffsetStereo(uintptr(unsafe.Pointer(&matrix)), eye)
	return matrix
}

// SetMatrixProjectionStereo - Set eyes projection matrices for stereo rendering
func SetMatrixProjectionStereo(right Matrix, left Matrix) {
	rlSetMatrixProjectionStereo(uintptr(unsafe.Pointer(&right)), uintptr(unsafe.Pointer(&left)))
}

// SetMatrixViewOffsetStereo - Set eyes view offsets matrices for stereo rendering
func SetMatrixViewOffsetStereo(right Matrix, left Matrix) {
	rlSetMatrixViewOffsetStereo(uintptr(unsafe.Pointer(&right)), uintptr(unsafe.Pointer(&left)))
}

// LoadDrawCube - Load and draw a cube
func LoadDrawCube() {
	rlLoadDrawCube()
}

// LoadDrawQuad - Load and draw a quad
func LoadDrawQuad() {
	rlLoadDrawQuad()
}
