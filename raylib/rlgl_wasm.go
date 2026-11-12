//go:build js && wasm

package rl

import (
	"unsafe"
)

// MatrixMode - Choose the current matrix to be transformed
func MatrixMode(mode int32) {
	crlMatrixMode(mode)
}

// PushMatrix - Push the current matrix to stack
func PushMatrix() {
	crlPushMatrix()
}

// PopMatrix - Pop lattest inserted matrix from stack
func PopMatrix() {
	crlPopMatrix()
}

// LoadIdentity - Reset current matrix to identity matrix
func LoadIdentity() {
	crlLoadIdentity()
}

// Translatef - Multiply the current matrix by a translation matrix
func Translatef(x float32, y float32, z float32) {
	crlTranslatef(x, y, z)
}

// Rotatef - Multiply the current matrix by a rotation matrix
func Rotatef(angle float32, x float32, y float32, z float32) {
	crlRotatef(angle, x, y, z)
}

// Scalef - Multiply the current matrix by a scaling matrix
func Scalef(x float32, y float32, z float32) {
	crlScalef(x, y, z)
}

// MultMatrix - Multiply the current matrix by another matrix
func MultMatrix(m Matrix) {
	defer scratch.release(scratch.mark())
	crlMultMatrixf(scratchSlice(MatrixToFloat(m)))
}

// Frustum .
func Frustum(left float64, right float64, bottom float64, top float64, znear float64, zfar float64) {
	crlFrustum(left, right, bottom, top, znear, zfar)
}

// Ortho .
func Ortho(left float64, right float64, bottom float64, top float64, znear float64, zfar float64) {
	crlOrtho(left, right, bottom, top, znear, zfar)
}

// Viewport - Set the viewport area
func Viewport(x int32, y int32, width int32, height int32) {
	crlViewport(x, y, width, height)
}

// SetClipPlanes - Set clip planes distances
func SetClipPlanes(nearPlane, farPlane float64) {
	crlSetClipPlanes(nearPlane, farPlane)
}

// GetCullDistanceNear - Get cull plane distance near
func GetCullDistanceNear() float64 {
	return crlGetCullDistanceNear()
}

// GetCullDistanceFar - Get cull plane distance far
func GetCullDistanceFar() float64 {
	return crlGetCullDistanceFar()
}

// Begin - Initialize drawing mode (how to organize vertex)
func Begin(mode int32) {
	crlBegin(mode)
}

// End - Finish vertex providing
func End() {
	crlEnd()
}

// Vertex2i - Define one vertex (position) - 2 int
func Vertex2i(x int32, y int32) {
	crlVertex2i(x, y)
}

// Vertex2f - Define one vertex (position) - 2 float
func Vertex2f(x float32, y float32) {
	crlVertex2f(x, y)
}

// Vertex3f - Define one vertex (position) - 3 float
func Vertex3f(x float32, y float32, z float32) {
	crlVertex3f(x, y, z)
}

// TexCoord2f - Define one vertex (texture coordinate) - 2 float
func TexCoord2f(x float32, y float32) {
	crlTexCoord2f(x, y)
}

// Normal3f - Define one vertex (normal) - 3 float
func Normal3f(x float32, y float32, z float32) {
	crlNormal3f(x, y, z)
}

// Color4ub - Define one vertex (color) - 4 byte
func Color4ub(r uint8, g uint8, b uint8, a uint8) {
	crlColor4ub(uint32(r), uint32(g), uint32(b), uint32(a))
}

// Color3f - Define one vertex (color) - 3 float
func Color3f(x float32, y float32, z float32) {
	crlColor3f(x, y, z)
}

// Color4f - Define one vertex (color) - 4 float
func Color4f(x float32, y float32, z float32, w float32) {
	crlColor4f(x, y, z, w)
}

// EnableVertexArray - Enable vertex array (VAO, if supported)
func EnableVertexArray(vaoId uint32) bool {
	return crlEnableVertexArray(vaoId) != 0
}

// DisableVertexArray - Disable vertex array (VAO, if supported)
func DisableVertexArray() {
	crlDisableVertexArray()
}

// EnableVertexBuffer - Enable vertex buffer (VBO)
func EnableVertexBuffer(id uint32) {
	crlEnableVertexBuffer(id)
}

// DisableVertexBuffer - Disable vertex buffer (VBO)
func DisableVertexBuffer() {
	crlDisableVertexBuffer()
}

// EnableVertexBufferElement - Enable vertex buffer element (VBO element)
func EnableVertexBufferElement(id uint32) {
	crlEnableVertexBufferElement(id)
}

// DisableVertexBufferElement - Disable vertex buffer element (VBO element)
func DisableVertexBufferElement() {
	crlDisableVertexBufferElement()
}

// EnableVertexAttribute - Enable vertex attribute index
func EnableVertexAttribute(index uint32) {
	crlEnableVertexAttribute(index)
}

// DisableVertexAttribute - Disable vertex attribute index
func DisableVertexAttribute(index uint32) {
	crlDisableVertexAttribute(index)
}

// EnableStatePointer - Enable attribute state pointer
func EnableStatePointer(vertexAttribType int32, buffer unsafe.Pointer) {
	crlEnableStatePointer(vertexAttribType, uint32(uintptr(buffer)))
}

// DisableStatePointer - Disable attribute state pointer
func DisableStatePointer(vertexAttribType int32) {
	crlDisableStatePointer(vertexAttribType)
}

// ActiveTextureSlot - Select and active a texture slot
func ActiveTextureSlot(slot int32) {
	crlActiveTextureSlot(slot)
}

// EnableTexture - Enable texture
func EnableTexture(id uint32) {
	crlEnableTexture(id)
}

// DisableTexture - Disable texture
func DisableTexture() {
	crlDisableTexture()
}

// EnableTextureCubemap - Enable texture cubemap
func EnableTextureCubemap(id uint32) {
	crlEnableTextureCubemap(id)
}

// DisableTextureCubemap - Disable texture cubemap
func DisableTextureCubemap() {
	crlDisableTextureCubemap()
}

// TextureParameters - Set texture parameters (filter, wrap)
func TextureParameters(id uint32, param int32, value int32) {
	crlTextureParameters(id, param, value)
}

// CubemapParameters - Set cubemap parameters (filter, wrap)
func CubemapParameters(id uint32, param int32, value int32) {
	crlCubemapParameters(id, param, value)
}

// EnableShader - Enable shader program
func EnableShader(id uint32) {
	crlEnableShader(id)
}

// DisableShader - Disable shader program
func DisableShader() {
	crlDisableShader()
}

// EnableFramebuffer - Enable render texture (fbo)
func EnableFramebuffer(id uint32) {
	crlEnableFramebuffer(id)
}

// DisableFramebuffer - Disable render texture (fbo), return to default framebuffer
func DisableFramebuffer() {
	crlDisableFramebuffer()
}

// GetActiveFramebuffer - Get the currently active render texture (fbo), 0 for default framebuffer
func GetActiveFramebuffer() uint32 {
	return crlGetActiveFramebuffer()
}

// ActiveDrawBuffers - Activate multiple draw color buffers
func ActiveDrawBuffers(count int32) {
	crlActiveDrawBuffers(count)
}

// BlitFramebuffer - Blit active framebuffer to main framebuffer
func BlitFramebuffer(srcX, srcY, srcWidth, srcHeight, dstX, dstY, dstWidth, dstHeight, bufferMask int32) {
	crlBlitFramebuffer(srcX, srcY, srcWidth, srcHeight, dstX, dstY, dstWidth, dstHeight, bufferMask)
}

// BindFramebuffer - Bind framebuffer (FBO)
func BindFramebuffer(target, framebuffer uint32) {
	crlBindFramebuffer(target, framebuffer)
}

// EnableColorBlend - Enable color blending
func EnableColorBlend() {
	crlEnableColorBlend()
}

// DisableColorBlend - Disable color blending
func DisableColorBlend() {
	crlDisableColorBlend()
}

// EnableDepthTest - Enable depth test
func EnableDepthTest() {
	crlEnableDepthTest()
}

// DisableDepthTest - Disable depth test
func DisableDepthTest() {
	crlDisableDepthTest()
}

// EnableDepthMask - Enable depth write
func EnableDepthMask() {
	crlEnableDepthMask()
}

// DisableDepthMask - Disable depth write
func DisableDepthMask() {
	crlDisableDepthMask()
}

// EnableBackfaceCulling - Enable backface culling
func EnableBackfaceCulling() {
	crlEnableBackfaceCulling()
}

// DisableBackfaceCulling - Disable backface culling
func DisableBackfaceCulling() {
	crlDisableBackfaceCulling()
}

// ColorMask - Color mask control
func ColorMask(r, g, b, a bool) {
	crlColorMask(boolToInt32(r), boolToInt32(g), boolToInt32(b), boolToInt32(a))
}

// SetCullFace - Set face culling mode
func SetCullFace(mode int32) {
	crlSetCullFace(mode)
}

// EnableScissorTest - Enable scissor test
func EnableScissorTest() {
	crlEnableScissorTest()
}

// DisableScissorTest - Disable scissor test
func DisableScissorTest() {
	crlDisableScissorTest()
}

// Scissor - Scissor test
func Scissor(x int32, y int32, width int32, height int32) {
	crlScissor(x, y, width, height)
}

// EnablePointMode - Enable point mode
func EnablePointMode() {
	crlEnablePointMode()
}

// DisablePointMode - Disable point mode
func DisablePointMode() {
	crlDisablePointMode()
}

// SetPointSize - Set the point drawing size
func SetPointSize(size float32) {
	crlSetPointSize(size)
}

// GetPointSize - Get the point drawing size
func GetPointSize() float32 {
	return crlGetPointSize()
}

// EnableWireMode - Enable wire mode
func EnableWireMode() {
	crlEnableWireMode()
}

// DisableWireMode - Disable wire mode
func DisableWireMode() {
	crlDisableWireMode()
}

// SetLineWidth - Set the line drawing width
func SetLineWidth(width float32) {
	crlSetLineWidth(width)
}

// GetLineWidth - Get the line drawing width
func GetLineWidth() float32 {
	return crlGetLineWidth()
}

// EnableSmoothLines - Enable line aliasing
func EnableSmoothLines() {
	crlEnableSmoothLines()
}

// DisableSmoothLines - Disable line aliasing
func DisableSmoothLines() {
	crlDisableSmoothLines()
}

// EnableStereoRender - Enable stereo rendering
func EnableStereoRender() {
	crlEnableStereoRender()
}

// DisableStereoRender - Disable stereo rendering
func DisableStereoRender() {
	crlDisableStereoRender()
}

// IsStereoRenderEnabled - Check if stereo render is enabled
func IsStereoRenderEnabled() bool {
	return crlIsStereoRenderEnabled() != 0
}

// ClearColor - Clear color buffer with color
func ClearColor(r, g, b, a uint8) {
	crlClearColor(uint32(r), uint32(g), uint32(b), uint32(a))
}

// ClearScreenBuffers - Clear used screen buffers (color and depth)
func ClearScreenBuffers() {
	crlClearScreenBuffers()
}

// CheckErrors - Check and log OpenGL error codes
func CheckErrors() {
	crlCheckErrors()
}

// SetBlendMode - Set blending mode
func SetBlendMode(mode BlendMode) {
	crlSetBlendMode(int32(mode))
}

// SetBlendFactors - Set blending mode factor and equation (using OpenGL factors)
func SetBlendFactors(glSrcFactor, glDstFactor, glEquation int32) {
	crlSetBlendFactors(glSrcFactor, glDstFactor, glEquation)
}

// SetBlendFactorsSeparate - Set blending mode factors and equations separately (using OpenGL factors)
func SetBlendFactorsSeparate(glSrcRGB, glDstRGB, glSrcAlpha, glDstAlpha, glEqRGB, glEqAlpha int32) {
	crlSetBlendFactorsSeparate(glSrcRGB, glDstRGB, glSrcAlpha, glDstAlpha, glEqRGB, glEqAlpha)
}

// GlInit - Initialize rlgl (buffers, shaders, textures, states)
func GlInit(width int32, height int32) {
	crlglInit(width, height)
}

// GlClose - De-inititialize rlgl (buffers, shaders, textures)
func GlClose() {
	crlglClose()
}

// GetProcAddress - Get OpenGL procedure address
func GetProcAddress(procName string) unsafe.Pointer {
	return nil
}

// GetVersion - Get current OpenGL version
func GetVersion() int32 {
	return crlGetVersion()
}

// SetFramebufferWidth - Set current framebuffer width
func SetFramebufferWidth(width int32) {
	crlSetFramebufferWidth(width)
}

// GetFramebufferWidth - Get default framebuffer width
func GetFramebufferWidth() int32 {
	return crlGetFramebufferWidth()
}

// SetFramebufferHeight - Set current framebuffer height
func SetFramebufferHeight(height int32) {
	crlSetFramebufferHeight(height)
}

// GetFramebufferHeight - Get default framebuffer height
func GetFramebufferHeight() int32 {
	return crlGetFramebufferHeight()
}

// GetTextureIdDefault - Get default texture id
func GetTextureIdDefault() uint32 {
	return crlGetTextureIdDefault()
}

// GetShaderIdDefault - Get default shader id
func GetShaderIdDefault() uint32 {
	return crlGetShaderIdDefault()
}

// GetShaderLocsDefault - Get default shader locations
func GetShaderLocsDefault() []int32 {
	return goSlice[int32](crlGetShaderLocsDefault(), MaxShaderLocations)
}

// LoadRenderBatch - Load a render batch system
func LoadRenderBatch(numBuffers int32, bufferElements int32) RenderBatch {
	crlLoadRenderBatch(numBuffers, bufferElements)
	return renderBatchFromC(readRet[cRlRenderBatch]())
}

// UnloadRenderBatch - Unload render batch system
func UnloadRenderBatch(batch RenderBatch) {
	if p := renderBatchC(&batch); p != 0 {
		crlUnloadRenderBatch(p)
		cFree(p)
		delete(renderBatches, batch.VertexBuffer)
	}
}

// DrawRenderBatch - Draw render batch data (Update->Draw->Reset)
func DrawRenderBatch(batch *RenderBatch) {
	crlDrawRenderBatch(renderBatchC(batch))
}

// SetRenderBatchActive - Set the active render batch for rlgl (NULL for default internal)
func SetRenderBatchActive(batch *RenderBatch) {
	crlSetRenderBatchActive(renderBatchC(batch))
}

// DrawRenderBatchActive - Update and draw internal render batch
func DrawRenderBatchActive() {
	crlDrawRenderBatchActive()
}

// CheckRenderBatchLimit - Check internal buffer overflow for a given number of vertex
func CheckRenderBatchLimit(vCount int32) bool {
	return crlCheckRenderBatchLimit(vCount) != 0
}

// SetTexture - Set current texture for render batch and check buffers limits
func SetTexture(id uint32) {
	crlSetTexture(id)
}

// LoadVertexArray - Load vertex array (vao) if supported
func LoadVertexArray() uint32 {
	return crlLoadVertexArray()
}

// LoadVertexBuffer - Load a vertex buffer object
func LoadVertexBuffer[T any](buffer []T, dynamic bool) uint32 {
	if len(buffer) == 0 {
		return 0
	}
	defer scratch.release(scratch.mark())
	size := int32(int(unsafe.Sizeof(buffer[0])) * len(buffer))
	return crlLoadVertexBuffer(scratchSlice(buffer), size, boolToInt32(dynamic))
}

// LoadVertexBufferElement - Load vertex buffer elements object
func LoadVertexBufferElement[T any](buffer []T, dynamic bool) uint32 {
	if len(buffer) == 0 {
		return 0
	}
	defer scratch.release(scratch.mark())
	size := int32(int(unsafe.Sizeof(buffer[0])) * len(buffer))
	return crlLoadVertexBufferElement(scratchSlice(buffer), size, boolToInt32(dynamic))
}

// UpdateVertexBuffer - Update vertex buffer object data on GPU buffer
func UpdateVertexBuffer[T any](bufferId uint32, data []T, offset int32) {
	if len(data) == 0 {
		return
	}
	defer scratch.release(scratch.mark())
	dataSize := int32(int(unsafe.Sizeof(data[0])) * len(data))
	crlUpdateVertexBuffer(bufferId, scratchSlice(data), dataSize, offset)
}

// UpdateVertexBufferElements - Update vertex buffer elements data on GPU buffer
func UpdateVertexBufferElements[T any](id uint32, data []T, offset int32) {
	if len(data) == 0 {
		return
	}
	defer scratch.release(scratch.mark())
	dataSize := int32(int(unsafe.Sizeof(data[0])) * len(data))
	crlUpdateVertexBufferElements(id, scratchSlice(data), dataSize, offset)
}

// UnloadVertexArray - Unload vertex array (vao)
func UnloadVertexArray(vaoId uint32) {
	crlUnloadVertexArray(vaoId)
}

// UnloadVertexBuffer - Unload vertex buffer object
func UnloadVertexBuffer(vboId uint32) {
	crlUnloadVertexBuffer(vboId)
}

// SetVertexAttribute - Set vertex attribute data configuration
func SetVertexAttribute(index uint32, compSize int32, attrType int32, normalized bool, stride int32, offset int32) {
	crlSetVertexAttribute(index, compSize, attrType, boolToInt32(normalized), stride, offset)
}

// SetVertexAttributeDivisor - Set vertex attribute data divisor
func SetVertexAttributeDivisor(index uint32, divisor int32) {
	crlSetVertexAttributeDivisor(index, divisor)
}

// SetVertexAttributeDefault - Set vertex attribute default value, when attribute to provided
func SetVertexAttributeDefault(locIndex int32, value unsafe.Pointer, attribType int32, count int32) {
	defer scratch.release(scratch.mark())
	crlSetVertexAttributeDefault(locIndex, scratch.bytes(value, uint32(max(count, 0))*4), attribType, count)
}

// DrawVertexArray - Draw vertex array (currently active vao)
func DrawVertexArray(offset int32, count int32) {
	crlDrawVertexArray(offset, count)
}

// DrawVertexArrayElements - Draw vertex array elements
func DrawVertexArrayElements(offset int32, count int32, buffer unsafe.Pointer) {
	crlDrawVertexArrayElements(offset, count, uint32(uintptr(buffer)))
}

// DrawVertexArrayInstanced - Draw vertex array (currently active vao) with instancing
func DrawVertexArrayInstanced(offset, count, instances int32) {
	crlDrawVertexArrayInstanced(offset, count, instances)
}

// DrawVertexArrayElementsInstanced - Draw vertex array elements with instancing
func DrawVertexArrayElementsInstanced(offset int32, count int32, buffer unsafe.Pointer, instances int32) {
	crlDrawVertexArrayElementsInstanced(offset, count, uint32(uintptr(buffer)), instances)
}

// LoadTextureDepth - Load depth texture/renderbuffer (to be attached to fbo)
func LoadTextureDepth(width, height int32, useRenderBuffer bool) uint32 {
	return crlLoadTextureDepth(width, height, boolToInt32(useRenderBuffer))
}

// LoadFramebuffer - Load an empty framebuffer
func LoadFramebuffer() uint32 {
	return crlLoadFramebuffer()
}

// FramebufferAttach - Attach texture/renderbuffer to a framebuffer
func FramebufferAttach(id, texId uint32, attachType, texType, mipLevel int32) {
	crlFramebufferAttach(id, texId, attachType, texType, mipLevel)
}

// FramebufferComplete - Verify framebuffer is complete
func FramebufferComplete(id uint32) bool {
	return crlFramebufferComplete(id) != 0
}

// UnloadFramebuffer - Delete framebuffer from GPU
func UnloadFramebuffer(id uint32) {
	crlUnloadFramebuffer(id)
}

// CopyFramebuffer - Copy framebuffer pixel data to internal buffer
//
// WARNING: Copy and resize framebuffer functionality only defined for software backend
func CopyFramebuffer(x, y, width, height, format int32, pixels unsafe.Pointer) {
	defer scratch.release(scratch.mark())
	n := uint32(cGetPixelDataSize(width, height, format))
	p := scratch.alloc(n)
	crlCopyFramebuffer(x, y, width, height, format, p)
	if pixels != nil {
		copyOut(pixels, p, n)
	}
}

// ResizeFramebuffer - Resize internal framebuffer
//
// WARNING: Copy and resize framebuffer functionality only defined for software backend
func ResizeFramebuffer(width, height int32) {
	crlResizeFramebuffer(width, height)
}

// LoadShaderId - Load (compile) shader and return shader id (type: [VertexShader], [FragmentShader], [ComputeShader])
func LoadShaderId(code string, shaderType int32) uint32 {
	defer scratch.release(scratch.mark())
	return crlLoadShader(scratch.str(code), shaderType)
}

// LoadShaderProgram - Load shader from code strings
func LoadShaderProgram(vsCode, fsCode string) uint32 {
	defer scratch.release(scratch.mark())
	return crlLoadShaderProgram(scratch.strNullable(vsCode), scratch.strNullable(fsCode))
}

// LoadShaderProgramEx - Load shader program, using already loaded shader ids
func LoadShaderProgramEx(vsId, fsId uint32) uint32 {
	return crlLoadShaderProgramEx(vsId, fsId)
}

// LoadShaderProgramCompute - Load compute shader program
func LoadShaderProgramCompute(csId uint32) uint32 {
	return crlLoadShaderProgramCompute(csId)
}

// UnloadShaderId - Unload shader, loaded with [LoadShaderId]
func UnloadShaderId(id uint32) {
	crlUnloadShader(id)
}

// UnloadShaderProgram - Unload shader program
func UnloadShaderProgram(id uint32) {
	crlUnloadShaderProgram(id)
}

// GetLocationUniform - Get shader location uniform, requires shader program id
func GetLocationUniform(id uint32, uniformName string) int32 {
	defer scratch.release(scratch.mark())
	return crlGetLocationUniform(id, scratch.str(uniformName))
}

// GetLocationAttrib - Get shader location attribute, requires shader program id
func GetLocationAttrib(id uint32, attribName string) int32 {
	defer scratch.release(scratch.mark())
	return crlGetLocationAttrib(id, scratch.str(attribName))
}

// SetUniform - Set shader value uniform ([]float32, []int32, []uint32)
func SetUniform[T ~float32 | ~int32 | ~uint32](locIndex int32, value []T, uniformType, count int32) {
	defer scratch.release(scratch.mark())
	crlSetUniform(locIndex, scratchSlice(value), uniformType, count)
}

// SetUniformMatrix - Set shader value matrix
func SetUniformMatrix(locIndex int32, mat Matrix) {
	crlSetUniformMatrix(locIndex, mat.M0, mat.M4, mat.M8, mat.M12, mat.M1, mat.M5, mat.M9, mat.M13, mat.M2, mat.M6, mat.M10, mat.M14, mat.M3, mat.M7, mat.M11, mat.M15)
}

// SetUniformMatrices - Set shader value matrices
func SetUniformMatrices(locIndex int32, mat []Matrix) {
	defer scratch.release(scratch.mark())
	crlSetUniformMatrices(locIndex, scratchSlice(mat), int32(len(mat)))
}

// SetUniformSampler - Set shader value sampler
func SetUniformSampler(locIndex int32, textureId uint32) {
	crlSetUniformSampler(locIndex, textureId)
}

// SetShader - Set shader currently active (id and locations)
func SetShader(id uint32, locs *int32) {
	crlSetShader(id, shaderToC(Shader{ID: id, Locs: locs}).Locs)
}

// ComputeShaderDispatch - Dispatch compute shader (equivalent to *draw* for graphics pilepine)
func ComputeShaderDispatch(groupX uint32, groupY uint32, groupZ uint32) {
	crlComputeShaderDispatch(groupX, groupY, groupZ)
}

// LoadShaderBuffer - Load shader storage buffer object (SSBO)
func LoadShaderBuffer(size uint32, data unsafe.Pointer, usageHint int32) uint32 {
	defer scratch.release(scratch.mark())
	return crlLoadShaderBuffer(size, scratch.bytes(data, size), usageHint)
}

// UnloadShaderBuffer - Unload shader storage buffer object (SSBO)
func UnloadShaderBuffer(id uint32) {
	crlUnloadShaderBuffer(id)
}

// UpdateShaderBuffer - Update SSBO buffer data
func UpdateShaderBuffer(id uint32, data unsafe.Pointer, dataSize uint32, offset uint32) {
	defer scratch.release(scratch.mark())
	crlUpdateShaderBuffer(id, scratch.bytes(data, dataSize), dataSize, offset)
}

// BindShaderBuffer - Bind SSBO buffer
func BindShaderBuffer(id uint32, index uint32) {
	crlBindShaderBuffer(id, index)
}

// ReadShaderBuffer - Read SSBO buffer data (GPU->CPU)
func ReadShaderBuffer(id uint32, dest unsafe.Pointer, count uint32, offset uint32) {
	defer scratch.release(scratch.mark())
	p := scratch.alloc(count)
	crlReadShaderBuffer(id, p, count, offset)
	if dest != nil {
		copyOut(dest, p, count)
	}
}

// CopyShaderBuffer - Copy SSBO data between buffers
func CopyShaderBuffer(destId uint32, srcId uint32, destOffset uint32, srcOffset uint32, count uint32) {
	crlCopyShaderBuffer(destId, srcId, destOffset, srcOffset, count)
}

// GetShaderBufferSize - Get SSBO buffer size
func GetShaderBufferSize(id uint32) uint32 {
	return crlGetShaderBufferSize(id)
}

// BindImageTexture - Bind image texture
func BindImageTexture(id uint32, index uint32, format int32, readonly bool) {
	crlBindImageTexture(id, index, format, boolToInt32(readonly))
}

// GetMatrixModelview - Get internal modelview matrix
func GetMatrixModelview() Matrix {
	crlGetMatrixModelview()
	return readRet[Matrix]()
}

// GetMatrixProjection - Get internal projection matrix
func GetMatrixProjection() Matrix {
	crlGetMatrixProjection()
	return readRet[Matrix]()
}

// GetMatrixTransform - Get internal accumulated transform matrix
func GetMatrixTransform() Matrix {
	crlGetMatrixTransform()
	return readRet[Matrix]()
}

// GetMatrixProjectionStereo - Get internal projection matrix for stereo render (selected eye)
func GetMatrixProjectionStereo(eye int32) Matrix {
	crlGetMatrixProjectionStereo(eye)
	return readRet[Matrix]()
}

// GetMatrixViewOffsetStereo - Get internal view offset matrix for stereo render (selected eye)
func GetMatrixViewOffsetStereo(eye int32) Matrix {
	crlGetMatrixViewOffsetStereo(eye)
	return readRet[Matrix]()
}

// SetMatrixProjection - Set a custom projection matrix (replaces internal projection matrix)
func SetMatrixProjection(proj Matrix) {
	crlSetMatrixProjection(proj.M0, proj.M4, proj.M8, proj.M12, proj.M1, proj.M5, proj.M9, proj.M13, proj.M2, proj.M6, proj.M10, proj.M14, proj.M3, proj.M7, proj.M11, proj.M15)
}

// SetMatrixModelview - Set a custom modelview matrix (replaces internal modelview matrix)
func SetMatrixModelview(view Matrix) {
	crlSetMatrixModelview(view.M0, view.M4, view.M8, view.M12, view.M1, view.M5, view.M9, view.M13, view.M2, view.M6, view.M10, view.M14, view.M3, view.M7, view.M11, view.M15)
}

// SetMatrixProjectionStereo - Set eyes projection matrices for stereo rendering
func SetMatrixProjectionStereo(right Matrix, left Matrix) {
	crlSetMatrixProjectionStereo(right.M0, right.M4, right.M8, right.M12, right.M1, right.M5, right.M9, right.M13, right.M2, right.M6, right.M10, right.M14, right.M3, right.M7, right.M11, right.M15, left.M0, left.M4, left.M8, left.M12, left.M1, left.M5, left.M9, left.M13, left.M2, left.M6, left.M10, left.M14, left.M3, left.M7, left.M11, left.M15)
}

// SetMatrixViewOffsetStereo - Set eyes view offsets matrices for stereo rendering
func SetMatrixViewOffsetStereo(right Matrix, left Matrix) {
	crlSetMatrixViewOffsetStereo(right.M0, right.M4, right.M8, right.M12, right.M1, right.M5, right.M9, right.M13, right.M2, right.M6, right.M10, right.M14, right.M3, right.M7, right.M11, right.M15, left.M0, left.M4, left.M8, left.M12, left.M1, left.M5, left.M9, left.M13, left.M2, left.M6, left.M10, left.M14, left.M3, left.M7, left.M11, left.M15)
}

// LoadDrawCube - Load and draw a cube
func LoadDrawCube() {
	crlLoadDrawCube()
}

// LoadDrawQuad - Load and draw a quad
func LoadDrawQuad() {
	crlLoadDrawQuad()
}
