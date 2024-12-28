package rl

/*
#include "raylib.h"
#include "rlgl.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// SetMatrixProjection - Set a custom projection matrix (replaces internal projection matrix)
func SetMatrixProjection(proj Matrix) {
	cproj := proj.cptr()
	C.rlSetMatrixProjection(*cproj)
}

// SetMatrixModelview - Set a custom modelview matrix (replaces internal modelview matrix)
func SetMatrixModelview(view Matrix) {
	cview := view.cptr()
	C.rlSetMatrixModelview(*cview)
}

// BeginShaderMode - Begin custom shader drawing
func BeginShaderMode(shader Shader) {
	cshader := shader.cptr()
	C.BeginShaderMode(*cshader)
}

// EndShaderMode - End custom shader drawing (use default shader)
func EndShaderMode() {
	C.EndShaderMode()
}

// BeginBlendMode - Begin blending mode (alpha, additive, multiplied)
func BeginBlendMode(mode BlendMode) {
	cmode := (C.int)(mode)
	C.BeginBlendMode(cmode)
}

// EndBlendMode - End blending mode (reset to default: alpha blending)
func EndBlendMode() {
	C.EndBlendMode()
}

// MatrixMode - Choose the current matrix to be transformed
func MatrixMode(mode int32) {
	cmode := C.int(mode)
	C.rlMatrixMode(cmode)
}

// PushMatrix - Push the current matrix to stack
func PushMatrix() {
	C.rlPushMatrix()
}

// PopMatrix - Pop lattest inserted matrix from stack
func PopMatrix() {
	C.rlPopMatrix()
}

// LoadIdentity - Reset current matrix to identity matrix
func LoadIdentity() {
	C.rlLoadIdentity()
}

// Translatef - Multiply the current matrix by a translation matrix
func Translatef(x float32, y float32, z float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlTranslatef(cx, cy, cz)
}

// Rotatef - Multiply the current matrix by a rotation matrix
func Rotatef(angle float32, x float32, y float32, z float32) {
	cangle := C.float(angle)
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlRotatef(cangle, cx, cy, cz)
}

// Scalef - Multiply the current matrix by a scaling matrix
func Scalef(x float32, y float32, z float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlScalef(cx, cy, cz)
}

// MultMatrix - Multiply the current matrix by another matrix
func MultMatrix(m Matrix) {
	f := MatrixToFloat(m)
	C.rlMultMatrixf((*C.float)(&f[0]))
}

// Frustum .
func Frustum(left float64, right float64, bottom float64, top float64, znear float64, zfar float64) {
	cleft := C.double(left)
	cright := C.double(right)
	cbottom := C.double(bottom)
	ctop := C.double(top)
	cznear := C.double(znear)
	czfar := C.double(zfar)
	C.rlFrustum(cleft, cright, cbottom, ctop, cznear, czfar)
}

// Ortho .
func Ortho(left float64, right float64, bottom float64, top float64, znear float64, zfar float64) {
	cleft := C.double(left)
	cright := C.double(right)
	cbottom := C.double(bottom)
	ctop := C.double(top)
	cznear := C.double(znear)
	czfar := C.double(zfar)
	C.rlOrtho(cleft, cright, cbottom, ctop, cznear, czfar)
}

// Viewport - Set the viewport area
func Viewport(x int32, y int32, width int32, height int32) {
	cx := C.int(x)
	cy := C.int(y)
	cwidth := C.int(width)
	cheight := C.int(height)
	C.rlViewport(cx, cy, cwidth, cheight)
}

// SetClipPlanes - Set clip planes distances
func SetClipPlanes(nearPlane, farPlane float64) {
	C.rlSetClipPlanes(C.double(nearPlane), C.double(farPlane))
}

// GetCullDistanceNear - Get cull plane distance near
func GetCullDistanceNear() float64 {
	ret := C.rlGetCullDistanceNear()
	return float64(ret)
}

// GetCullDistanceFar - Get cull plane distance far
func GetCullDistanceFar() float64 {
	ret := C.rlGetCullDistanceFar()
	return float64(ret)
}

// Begin - Initialize drawing mode (how to organize vertex)
func Begin(mode int32) {
	cmode := C.int(mode)
	C.rlBegin(cmode)
}

// End - Finish vertex providing
func End() {
	C.rlEnd()
}

// Vertex2i - Define one vertex (position) - 2 int
func Vertex2i(x int32, y int32) {
	cx := C.int(x)
	cy := C.int(y)
	C.rlVertex2i(cx, cy)
}

// Vertex2f - Define one vertex (position) - 2 float
func Vertex2f(x float32, y float32) {
	cx := C.float(x)
	cy := C.float(y)
	C.rlVertex2f(cx, cy)
}

// Vertex3f - Define one vertex (position) - 3 float
func Vertex3f(x float32, y float32, z float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlVertex3f(cx, cy, cz)
}

// TexCoord2f - Define one vertex (texture coordinate) - 2 float
func TexCoord2f(x float32, y float32) {
	cx := C.float(x)
	cy := C.float(y)
	C.rlTexCoord2f(cx, cy)
}

// Normal3f - Define one vertex (normal) - 3 float
func Normal3f(x float32, y float32, z float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlNormal3f(cx, cy, cz)
}

// Color4ub - Define one vertex (color) - 4 byte
func Color4ub(r uint8, g uint8, b uint8, a uint8) {
	cr := C.uchar(r)
	cg := C.uchar(g)
	cb := C.uchar(b)
	ca := C.uchar(a)
	C.rlColor4ub(cr, cg, cb, ca)
}

// Color3f - Define one vertex (color) - 3 float
func Color3f(x float32, y float32, z float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlColor3f(cx, cy, cz)
}

// Color4f - Define one vertex (color) - 4 float
func Color4f(x float32, y float32, z float32, w float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	cw := C.float(w)
	C.rlColor4f(cx, cy, cz, cw)
}

// EnableVertexArray - Enable vertex array (VAO, if supported)
func EnableVertexArray(vaoId uint32) bool {
	cvaoId := C.uint(vaoId)
	return bool(C.rlEnableVertexArray(cvaoId))
}

// DisableVertexArray - Disable vertex array (VAO, if supported)
func DisableVertexArray() {
	C.rlDisableVertexArray()
}

// EnableVertexBuffer - Enable vertex buffer (VBO)
func EnableVertexBuffer(id uint32) {
	cid := C.uint(id)
	C.rlEnableVertexBuffer(cid)
}

// DisableVertexBuffer - Disable vertex buffer (VBO)
func DisableVertexBuffer() {
	C.rlDisableVertexBuffer()
}

// EnableVertexBufferElement - Enable vertex buffer element (VBO element)
func EnableVertexBufferElement(id uint32) {
	cid := C.uint(id)
	C.rlEnableVertexBufferElement(cid)
}

// DisableVertexBufferElement - Disable vertex buffer element (VBO element)
func DisableVertexBufferElement() {
	C.rlDisableVertexBufferElement()
}

// EnableVertexAttribute - Enable vertex attribute index
func EnableVertexAttribute(index uint32) {
	cindex := C.uint(index)
	C.rlEnableVertexAttribute(cindex)
}

// DisableVertexAttribute - Disable vertex attribute index
func DisableVertexAttribute(index uint32) {
	cindex := C.uint(index)
	C.rlDisableVertexAttribute(cindex)
}

// ActiveTextureSlot - Select and active a texture slot
func ActiveTextureSlot(slot int32) {
	cslot := C.int(slot)
	C.rlActiveTextureSlot(cslot)
}

// EnableTexture - Enable texture
func EnableTexture(id uint32) {
	cid := C.uint(id)
	C.rlEnableTexture(cid)
}

// DisableTexture - Disable texture
func DisableTexture() {
	C.rlDisableTexture()
}

// EnableTextureCubemap - Enable texture cubemap
func EnableTextureCubemap(id uint32) {
	cid := C.uint(id)
	C.rlEnableTextureCubemap(cid)
}

// DisableTextureCubemap - Disable texture cubemap
func DisableTextureCubemap() {
	C.rlDisableTextureCubemap()
}

// TextureParameters - Set texture parameters (filter, wrap)
func TextureParameters(id uint32, param int32, value int32) {
	cid := C.uint(id)
	cparam := C.int(param)
	cvalue := C.int(value)
	C.rlTextureParameters(cid, cparam, cvalue)
}

// CubemapParameters - Set cubemap parameters (filter, wrap)
func CubemapParameters(id uint32, param int32, value int32) {
	cid := C.uint(id)
	cparam := C.int(param)
	cvalue := C.int(value)
	C.rlCubemapParameters(cid, cparam, cvalue)
}

// EnableShader - Enable shader program
func EnableShader(id uint32) {
	cid := C.uint(id)
	C.rlEnableShader(cid)
}

// DisableShader - Disable shader program
func DisableShader() {
	C.rlDisableShader()
}

// EnableFramebuffer - Enable render texture (fbo)
func EnableFramebuffer(id uint32) {
	cid := C.uint(id)
	C.rlEnableFramebuffer(cid)
}

// DisableFramebuffer - Disable render texture (fbo), return to default framebuffer
func DisableFramebuffer() {
	C.rlDisableFramebuffer()
}

// GetActiveFramebuffer - Get the currently active render texture (fbo), 0 for default framebuffer
func GetActiveFramebuffer() uint32 {
	return uint32(C.rlGetActiveFramebuffer())
}

// ActiveDrawBuffers - Activate multiple draw color buffers
func ActiveDrawBuffers(count int32) {
	ccount := C.int(count)
	C.rlActiveDrawBuffers(ccount)
}

// BlitFramebuffer - Blit active framebuffer to main framebuffer
func BlitFramebuffer(srcX, srcY, srcWidth, srcHeight, dstX, dstY, dstWidth, dstHeight, bufferMask int32) {
	C.rlBlitFramebuffer(C.int(srcX), C.int(srcY), C.int(srcWidth), C.int(srcHeight), C.int(dstX), C.int(dstY), C.int(dstWidth), C.int(dstHeight), C.int(bufferMask))
}

// BindFramebuffer - Bind framebuffer (FBO)
func BindFramebuffer(target, framebuffer uint32) {
	C.rlBindFramebuffer(C.uint(target), C.uint(framebuffer))
}

// EnableColorBlend - Enable color blending
func EnableColorBlend() {
	C.rlEnableColorBlend()
}

// DisableColorBlend - Disable color blending
func DisableColorBlend() {
	C.rlDisableColorBlend()
}

// EnableDepthTest - Enable depth test
func EnableDepthTest() {
	C.rlEnableDepthTest()
}

// DisableDepthTest - Disable depth test
func DisableDepthTest() {
	C.rlDisableDepthTest()
}

// EnableDepthMask - Enable depth write
func EnableDepthMask() {
	C.rlEnableDepthMask()
}

// DisableDepthMask - Disable depth write
func DisableDepthMask() {
	C.rlDisableDepthMask()
}

// EnableBackfaceCulling - Enable backface culling
func EnableBackfaceCulling() {
	C.rlEnableBackfaceCulling()
}

// DisableBackfaceCulling - Disable backface culling
func DisableBackfaceCulling() {
	C.rlDisableBackfaceCulling()
}

// ColorMask - Color mask control
func ColorMask(r, g, b, a bool) {
	C.rlColorMask(C.bool(r), C.bool(g), C.bool(b), C.bool(a))
}

// SetCullFace - Set face culling mode
func SetCullFace(mode int32) {
	cmode := C.int(mode)
	C.rlSetCullFace(cmode)
}

// EnableScissorTest - Enable scissor test
func EnableScissorTest() {
	C.rlEnableScissorTest()
}

// DisableScissorTest - Disable scissor test
func DisableScissorTest() {
	C.rlDisableScissorTest()
}

// Scissor - Scissor test
func Scissor(x int32, y int32, width int32, height int32) {
	cx := C.int(x)
	cy := C.int(y)
	cwidth := C.int(width)
	cheight := C.int(height)
	C.rlScissor(cx, cy, cwidth, cheight)
}

// EnableWireMode - Enable wire mode
func EnableWireMode() {
	C.rlEnableWireMode()
}

// EnablePointMode - Enable point mode
func EnablePointMode() {
	C.rlEnablePointMode()
}

// DisableWireMode - Disable wire mode
func DisableWireMode() {
	C.rlDisableWireMode()
}

// SetLineWidth - Set the line drawing width
func SetLineWidth(width float32) {
	cwidth := C.float(width)
	C.rlSetLineWidth(cwidth)
}

// GetLineWidth - Get the line drawing width
func GetLineWidth() float32 {
	return float32(C.rlGetLineWidth())
}

// EnableSmoothLines - Enable line aliasing
func EnableSmoothLines() {
	C.rlEnableSmoothLines()
}

// DisableSmoothLines - Disable line aliasing
func DisableSmoothLines() {
	C.rlDisableSmoothLines()
}

// EnableStereoRender - Enable stereo rendering
func EnableStereoRender() {
	C.rlEnableStereoRender()
}

// DisableStereoRender - Disable stereo rendering
func DisableStereoRender() {
	C.rlDisableStereoRender()
}

// IsStereoRenderEnabled - Check if stereo render is enabled
func IsStereoRenderEnabled() bool {
	return bool(C.rlIsStereoRenderEnabled())
}

// ClearColor - Clear color buffer with color
func ClearColor(r uint8, g uint8, b uint8, a uint8) {
	cr := C.uchar(r)
	cg := C.uchar(g)
	cb := C.uchar(b)
	ca := C.uchar(a)
	C.rlClearColor(cr, cg, cb, ca)
}

// ClearScreenBuffers - Clear used screen buffers (color and depth)
func ClearScreenBuffers() {
	C.rlClearScreenBuffers()
}

// CheckErrors - Check and log OpenGL error codes
func CheckErrors() {
	C.rlCheckErrors()
}

// SetBlendMode - Set blending mode
func SetBlendMode(mode BlendMode) {
	cmode := C.int(mode)
	C.rlSetBlendMode(cmode)
}

// SetBlendFactors - Set blending mode factor and equation (using OpenGL factors)
func SetBlendFactors(glSrcFactor int32, glDstFactor int32, glEquation int32) {
	cglSrcFactor := C.int(glSrcFactor)
	cglDstFactor := C.int(glDstFactor)
	cglEquation := C.int(glEquation)
	C.rlSetBlendFactors(cglSrcFactor, cglDstFactor, cglEquation)
}

// SetBlendFactorsSeparate - Set blending mode factors and equations separately (using OpenGL factors)
func SetBlendFactorsSeparate(glSrcRGB int32, glDstRGB int32, glSrcAlpha int32, glDstAlpha int32, glEqRGB int32, glEqAlpha int32) {
	cglSrcRGB := C.int(glSrcRGB)
	cglDstRGB := C.int(glDstRGB)
	cglSrcAlpha := C.int(glSrcAlpha)
	cglDstAlpha := C.int(glDstAlpha)
	cglEqRGB := C.int(glEqRGB)
	cglEqAlpha := C.int(glEqAlpha)
	C.rlSetBlendFactorsSeparate(cglSrcRGB, cglDstRGB, cglSrcAlpha, cglDstAlpha, cglEqRGB, cglEqAlpha)
}

// GlInit - Initialize rlgl (buffers, shaders, textures, states)
func GlInit(width int32, height int32) {
	cwidth := C.int(width)
	cheight := C.int(height)
	C.rlglInit(cwidth, cheight)
}

// GlClose - De-inititialize rlgl (buffers, shaders, textures)
func GlClose() {
	C.rlglClose()
}

// GetVersion - Get current OpenGL version
func GetVersion() int32 {
	return int32(C.rlGetVersion())
}

// SetFramebufferWidth - Set current framebuffer width
func SetFramebufferWidth(width int32) {
	cwidth := C.int(width)
	C.rlSetFramebufferWidth(cwidth)
}

// GetFramebufferWidth - Get default framebuffer width
func GetFramebufferWidth() int32 {
	return int32(C.rlGetFramebufferWidth())
}

// SetFramebufferHeight - Set current framebuffer height
func SetFramebufferHeight(height int32) {
	cheight := C.int(height)
	C.rlSetFramebufferHeight(cheight)
}

// GetFramebufferHeight - Get default framebuffer height
func GetFramebufferHeight() int32 {
	return int32(C.rlGetFramebufferHeight())
}

// GetTextureIdDefault - Get default texture id
func GetTextureIdDefault() uint32 {
	return uint32(C.rlGetTextureIdDefault())
}

// GetShaderIdDefault - Get default shader id
func GetShaderIdDefault() uint32 {
	return uint32(C.rlGetShaderIdDefault())
}

// LoadRenderBatch - Load a render batch system
func LoadRenderBatch(numBuffers int32, bufferElements int32) RenderBatch {
	ret := C.rlLoadRenderBatch(C.int(numBuffers), C.int(bufferElements))
	return *(*RenderBatch)(unsafe.Pointer(&ret))
}

// UnloadRenderBatch - Unload render batch system
func UnloadRenderBatch(batch RenderBatch) {
	C.rlUnloadRenderBatch(*(*C.rlRenderBatch)(unsafe.Pointer(&batch)))
}

// DrawRenderBatch - Draw render batch data (Update->Draw->Reset)
func DrawRenderBatch(batch *RenderBatch) {
	C.rlDrawRenderBatch((*C.rlRenderBatch)(unsafe.Pointer(batch)))
}

// SetRenderBatchActive - Set the active render batch for rlgl (NULL for default internal)
func SetRenderBatchActive(batch *RenderBatch) {
	C.rlSetRenderBatchActive((*C.rlRenderBatch)(unsafe.Pointer(batch)))
}

// DrawRenderBatchActive - Update and draw internal render batch
func DrawRenderBatchActive() {
	C.rlDrawRenderBatchActive()
}

// CheckRenderBatchLimit - Check internal buffer overflow for a given number of vertex
func CheckRenderBatchLimit(vCount int32) bool {
	cvCount := C.int(vCount)
	return bool(C.rlCheckRenderBatchLimit(cvCount))
}

// SetTexture - Set current texture for render batch and check buffers limits
func SetTexture(id uint32) {
	cid := C.uint(id)
	C.rlSetTexture(cid)
}

// LoadVertexArray - Load vertex array (vao) if supported
func LoadVertexArray() uint32 {
	return uint32(C.rlLoadVertexArray())
}

// UnloadVertexBuffer .
func UnloadVertexBuffer(vboId uint32) {
	cvboId := C.uint(vboId)
	C.rlUnloadVertexBuffer(cvboId)
}

// SetVertexAttributeDivisor .
func SetVertexAttributeDivisor(index uint32, divisor int32) {
	cindex := C.uint(index)
	cdivisor := C.int(divisor)
	C.rlSetVertexAttributeDivisor(cindex, cdivisor)
}

// LoadTextureDepth - Load depth texture/renderbuffer (to be attached to fbo)
func LoadTextureDepth(width, height int32, useRenderBuffer bool) uint32 {
	cwidth := C.int(width)
	cheight := C.int(height)
	cuseRenderBuffer := C.bool(useRenderBuffer)
	cid := C.rlLoadTextureDepth(cwidth, cheight, cuseRenderBuffer)
	return uint32(cid)
}

// LoadFramebuffer - Load an empty framebuffer
func LoadFramebuffer() uint32 {
	return uint32(C.rlLoadFramebuffer())
}

// FramebufferAttach - Attach texture/renderbuffer to a framebuffer
func FramebufferAttach(fboId uint32, texId uint32, attachType int32, texType int32, mipLevel int32) {
	cfboId := C.uint(fboId)
	ctexId := C.uint(texId)
	cattachType := C.int(attachType)
	ctexType := C.int(texType)
	cmipLevel := C.int(mipLevel)
	C.rlFramebufferAttach(cfboId, ctexId, cattachType, ctexType, cmipLevel)
}

// FramebufferComplete - Verify framebuffer is complete
func FramebufferComplete(id uint32) bool {
	cid := C.uint(id)
	return bool(C.rlFramebufferComplete(cid))
}

// UnloadFramebuffer - Delete framebuffer from GPU
func UnloadFramebuffer(id uint32) {
	cid := C.uint(id)
	C.rlUnloadFramebuffer(cid)
}

// LoadShaderCode - Load shader from code strings
func LoadShaderCode(vsCode string, fsCode string) uint32 {
	cvsCode := C.CString(vsCode)
	defer C.free(unsafe.Pointer(cvsCode))
	cfsCode := C.CString(fsCode)
	defer C.free(unsafe.Pointer(cfsCode))
	return uint32(C.rlLoadShaderCode(cvsCode, cfsCode))
}

// CompileShader - Compile custom shader and return shader id (type: VERTEX_SHADER, FRAGMENT_SHADER, COMPUTE_SHADER)
func CompileShader(shaderCode string, type_ int32) uint32 {
	cshaderCode := C.CString(shaderCode)
	defer C.free(unsafe.Pointer(cshaderCode))
	ctype_ := C.int(type_)
	return uint32(C.rlCompileShader(cshaderCode, ctype_))
}

// LoadShaderProgram - Load custom shader program
func LoadShaderProgram(vShaderId uint32, fShaderId uint32) uint32 {
	cvShaderId := C.uint(vShaderId)
	cfShaderId := C.uint(fShaderId)
	return uint32(C.rlLoadShaderProgram(cvShaderId, cfShaderId))
}

// UnloadShaderProgram - Unload shader program
func UnloadShaderProgram(id uint32) {
	cid := C.uint(id)
	C.rlUnloadShaderProgram(cid)
}

// GetLocationUniform - Get shader location uniform
func GetLocationUniform(shaderId uint32, uniformName string) int32 {
	cshaderId := C.uint(shaderId)
	cuniformName := C.CString(uniformName)
	defer C.free(unsafe.Pointer(cuniformName))
	return int32(C.rlGetLocationUniform(cshaderId, cuniformName))
}

// GetLocationAttrib - Get shader location attribute
func GetLocationAttrib(shaderId uint32, attribName string) int32 {
	cshaderId := C.uint(shaderId)
	cattribName := C.CString(attribName)
	defer C.free(unsafe.Pointer(cattribName))
	return int32(C.rlGetLocationAttrib(cshaderId, cattribName))
}

// SetUniform - Set shader value uniform
func SetUniform(locIndex int32, value []float32, uniformType int32) {
	C.rlSetUniform(C.int(locIndex), unsafe.Pointer(unsafe.SliceData(value)), C.int(uniformType), C.int((len(value))))
}

// SetUniformMatrix - Set shader value matrix
func SetUniformMatrix(locIndex int32, mat Matrix) {
	cmat := (*C.Matrix)(unsafe.Pointer(&mat))
	C.rlSetUniformMatrix(C.int(locIndex), *cmat)
}

// SetUniformMatrices - Set shader value matrices
func SetUniformMatrices(locIndex int32, mat []Matrix) {
	cmat := (*C.Matrix)(unsafe.Pointer(unsafe.SliceData(mat)))
	C.rlSetUniformMatrices(C.int(locIndex), cmat, C.int(len(mat)))
}

// SetUniformSampler - Set shader value sampler
func SetUniformSampler(locIndex int32, textureId uint32) {
	clocIndex := C.int(locIndex)
	ctextureId := C.uint(textureId)
	C.rlSetUniformSampler(clocIndex, ctextureId)
}

// LoadComputeShaderProgram -
func LoadComputeShaderProgram(shaderID uint32) uint32 {
	cshaderID := C.uint(shaderID)
	return uint32(C.rlLoadComputeShaderProgram(cshaderID))
}

// ComputeShaderDispatch - Dispatch compute shader (equivalent to *draw* for graphics pilepine)
func ComputeShaderDispatch(groupX uint32, groupY uint32, groupZ uint32) {
	cgroupX := C.uint(groupX)
	cgroupY := C.uint(groupY)
	cgroupZ := C.uint(groupZ)
	C.rlComputeShaderDispatch(cgroupX, cgroupY, cgroupZ)
}

// LoadShaderBuffer - Load shader storage buffer object (SSBO)
func LoadShaderBuffer(size uint32, data unsafe.Pointer, usageHint int32) uint32 {
	csize := C.uint(size)
	cdata := data
	cusageHint := C.int(usageHint)
	return uint32(C.rlLoadShaderBuffer(csize, cdata, cusageHint))
}

// UnloadShaderBuffer - Unload shader storage buffer object (SSBO)
func UnloadShaderBuffer(id uint32) {
	cid := C.uint(id)
	C.rlUnloadShaderBuffer(cid)
}

// UpdateShaderBuffer - Update SSBO buffer data
func UpdateShaderBuffer(id uint32, data unsafe.Pointer, dataSize uint32, offset uint32) {
	cid := C.uint(id)
	cdata := data
	cdataSize := C.uint(dataSize)
	coffset := C.uint(offset)
	C.rlUpdateShaderBuffer(cid, cdata, cdataSize, coffset)
}

// BindShaderBuffer - Bind SSBO buffer
func BindShaderBuffer(id uint32, index uint32) {
	cid := C.uint(id)
	cindex := C.uint(index)
	C.rlBindShaderBuffer(cid, cindex)
}

// ReadShaderBuffer - Read SSBO buffer data (GPU->CPU)
func ReadShaderBuffer(id uint32, dest unsafe.Pointer, count uint32, offset uint32) {
	cid := C.uint(id)
	cdest := dest
	ccount := C.uint(count)
	coffset := C.uint(offset)
	C.rlReadShaderBuffer(cid, cdest, ccount, coffset)
}

// CopyShaderBuffer - Copy SSBO data between buffers
func CopyShaderBuffer(destId uint32, srcId uint32, destOffset uint32, srcOffset uint32, count uint32) {
	cdestId := C.uint(destId)
	csrcId := C.uint(srcId)
	cdestOffset := C.uint(destOffset)
	csrcOffset := C.uint(srcOffset)
	ccount := C.uint(count)
	C.rlCopyShaderBuffer(cdestId, csrcId, cdestOffset, csrcOffset, ccount)
}

// GetShaderBufferSize - Get SSBO buffer size
func GetShaderBufferSize(id uint32) uint32 {
	cid := C.uint(id)
	return uint32(C.rlGetShaderBufferSize(cid))
}

// BindImageTexture - Bind image texture
func BindImageTexture(id uint32, index uint32, format int32, readonly bool) {
	cid := C.uint(id)
	cindex := C.uint(index)
	cformat := C.int(format)
	creadonly := C.bool(readonly)
	C.rlBindImageTexture(cid, cindex, cformat, creadonly)
}

// GetMatrixModelview - Get internal modelview matrix
func GetMatrixModelview() Matrix {
	cResult := C.rlGetMatrixModelview()
	var goRes Matrix
	goRes.M4 = float32(cResult.m4)
	goRes.M2 = float32(cResult.m2)
	goRes.M14 = float32(cResult.m14)
	goRes.M3 = float32(cResult.m3)
	goRes.M7 = float32(cResult.m7)
	goRes.M12 = float32(cResult.m12)
	goRes.M6 = float32(cResult.m6)
	goRes.M15 = float32(cResult.m15)
	goRes.M5 = float32(cResult.m5)
	goRes.M9 = float32(cResult.m9)
	goRes.M13 = float32(cResult.m13)
	goRes.M10 = float32(cResult.m10)
	goRes.M11 = float32(cResult.m11)
	goRes.M0 = float32(cResult.m0)
	goRes.M8 = float32(cResult.m8)
	goRes.M1 = float32(cResult.m1)
	return goRes
}

// GetMatrixProjection - Get internal projection matrix
func GetMatrixProjection() Matrix {
	cResult := C.rlGetMatrixProjection()
	var goRes Matrix
	goRes.M13 = float32(cResult.m13)
	goRes.M10 = float32(cResult.m10)
	goRes.M11 = float32(cResult.m11)
	goRes.M15 = float32(cResult.m15)
	goRes.M5 = float32(cResult.m5)
	goRes.M9 = float32(cResult.m9)
	goRes.M1 = float32(cResult.m1)
	goRes.M0 = float32(cResult.m0)
	goRes.M8 = float32(cResult.m8)
	goRes.M14 = float32(cResult.m14)
	goRes.M3 = float32(cResult.m3)
	goRes.M7 = float32(cResult.m7)
	goRes.M4 = float32(cResult.m4)
	goRes.M2 = float32(cResult.m2)
	goRes.M12 = float32(cResult.m12)
	goRes.M6 = float32(cResult.m6)
	return goRes
}

// GetMatrixTransform - Get internal accumulated transform matrix
func GetMatrixTransform() Matrix {
	cResult := C.rlGetMatrixTransform()
	var goRes Matrix
	goRes.M0 = float32(cResult.m0)
	goRes.M8 = float32(cResult.m8)
	goRes.M1 = float32(cResult.m1)
	goRes.M7 = float32(cResult.m7)
	goRes.M4 = float32(cResult.m4)
	goRes.M2 = float32(cResult.m2)
	goRes.M14 = float32(cResult.m14)
	goRes.M3 = float32(cResult.m3)
	goRes.M12 = float32(cResult.m12)
	goRes.M6 = float32(cResult.m6)
	goRes.M11 = float32(cResult.m11)
	goRes.M15 = float32(cResult.m15)
	goRes.M5 = float32(cResult.m5)
	goRes.M9 = float32(cResult.m9)
	goRes.M13 = float32(cResult.m13)
	goRes.M10 = float32(cResult.m10)
	return goRes
}

// GetMatrixProjectionStereo - Get internal projection matrix for stereo render (selected eye)
func GetMatrixProjectionStereo(eye int32) Matrix {
	ceye := C.int(eye)
	cResult := C.rlGetMatrixProjectionStereo(ceye)
	var goRes Matrix
	goRes.M12 = float32(cResult.m12)
	goRes.M6 = float32(cResult.m6)
	goRes.M15 = float32(cResult.m15)
	goRes.M5 = float32(cResult.m5)
	goRes.M9 = float32(cResult.m9)
	goRes.M13 = float32(cResult.m13)
	goRes.M10 = float32(cResult.m10)
	goRes.M11 = float32(cResult.m11)
	goRes.M0 = float32(cResult.m0)
	goRes.M8 = float32(cResult.m8)
	goRes.M1 = float32(cResult.m1)
	goRes.M4 = float32(cResult.m4)
	goRes.M2 = float32(cResult.m2)
	goRes.M14 = float32(cResult.m14)
	goRes.M3 = float32(cResult.m3)
	goRes.M7 = float32(cResult.m7)
	return goRes
}

// GetMatrixViewOffsetStereo - Get internal view offset matrix for stereo render (selected eye)
func GetMatrixViewOffsetStereo(eye int32) Matrix {
	ceye := C.int(eye)
	cResult := C.rlGetMatrixViewOffsetStereo(ceye)
	var goRes Matrix
	goRes.M0 = float32(cResult.m0)
	goRes.M8 = float32(cResult.m8)
	goRes.M1 = float32(cResult.m1)
	goRes.M4 = float32(cResult.m4)
	goRes.M2 = float32(cResult.m2)
	goRes.M14 = float32(cResult.m14)
	goRes.M3 = float32(cResult.m3)
	goRes.M7 = float32(cResult.m7)
	goRes.M12 = float32(cResult.m12)
	goRes.M6 = float32(cResult.m6)
	goRes.M5 = float32(cResult.m5)
	goRes.M9 = float32(cResult.m9)
	goRes.M13 = float32(cResult.m13)
	goRes.M10 = float32(cResult.m10)
	goRes.M11 = float32(cResult.m11)
	goRes.M15 = float32(cResult.m15)
	return goRes
}

// SetMatrixProjectionStereo - Set eyes projection matrices for stereo rendering
func SetMatrixProjectionStereo(right Matrix, left Matrix) {
	var cright C.struct_Matrix
	cright.m12 = C.float(right.M12)
	cright.m6 = C.float(right.M6)
	cright.m5 = C.float(right.M5)
	cright.m9 = C.float(right.M9)
	cright.m13 = C.float(right.M13)
	cright.m10 = C.float(right.M10)
	cright.m11 = C.float(right.M11)
	cright.m15 = C.float(right.M15)
	cright.m0 = C.float(right.M0)
	cright.m8 = C.float(right.M8)
	cright.m1 = C.float(right.M1)
	cright.m4 = C.float(right.M4)
	cright.m2 = C.float(right.M2)
	cright.m14 = C.float(right.M14)
	cright.m3 = C.float(right.M3)
	cright.m7 = C.float(right.M7)
	var cleft C.struct_Matrix
	cleft.m10 = C.float(left.M10)
	cleft.m11 = C.float(left.M11)
	cleft.m15 = C.float(left.M15)
	cleft.m5 = C.float(left.M5)
	cleft.m9 = C.float(left.M9)
	cleft.m13 = C.float(left.M13)
	cleft.m0 = C.float(left.M0)
	cleft.m8 = C.float(left.M8)
	cleft.m1 = C.float(left.M1)
	cleft.m3 = C.float(left.M3)
	cleft.m7 = C.float(left.M7)
	cleft.m4 = C.float(left.M4)
	cleft.m2 = C.float(left.M2)
	cleft.m14 = C.float(left.M14)
	cleft.m12 = C.float(left.M12)
	cleft.m6 = C.float(left.M6)
	C.rlSetMatrixProjectionStereo(cright, cleft)
}

// SetMatrixViewOffsetStereo - Set eyes view offsets matrices for stereo rendering
func SetMatrixViewOffsetStereo(right Matrix, left Matrix) {
	var cright C.struct_Matrix
	cright.m12 = C.float(right.M12)
	cright.m6 = C.float(right.M6)
	cright.m5 = C.float(right.M5)
	cright.m9 = C.float(right.M9)
	cright.m13 = C.float(right.M13)
	cright.m10 = C.float(right.M10)
	cright.m11 = C.float(right.M11)
	cright.m15 = C.float(right.M15)
	cright.m0 = C.float(right.M0)
	cright.m8 = C.float(right.M8)
	cright.m1 = C.float(right.M1)
	cright.m4 = C.float(right.M4)
	cright.m2 = C.float(right.M2)
	cright.m14 = C.float(right.M14)
	cright.m3 = C.float(right.M3)
	cright.m7 = C.float(right.M7)
	var cleft C.struct_Matrix
	cleft.m12 = C.float(left.M12)
	cleft.m6 = C.float(left.M6)
	cleft.m5 = C.float(left.M5)
	cleft.m9 = C.float(left.M9)
	cleft.m13 = C.float(left.M13)
	cleft.m10 = C.float(left.M10)
	cleft.m11 = C.float(left.M11)
	cleft.m15 = C.float(left.M15)
	cleft.m0 = C.float(left.M0)
	cleft.m8 = C.float(left.M8)
	cleft.m1 = C.float(left.M1)
	cleft.m4 = C.float(left.M4)
	cleft.m2 = C.float(left.M2)
	cleft.m14 = C.float(left.M14)
	cleft.m3 = C.float(left.M3)
	cleft.m7 = C.float(left.M7)
	C.rlSetMatrixViewOffsetStereo(cright, cleft)
}

// LoadDrawCube - Load and draw a cube
func LoadDrawCube() {
	C.rlLoadDrawCube()
}

// LoadDrawQuad - Load and draw a quad
func LoadDrawQuad() {
	C.rlLoadDrawQuad()
}
