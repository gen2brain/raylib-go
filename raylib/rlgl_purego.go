//go:build !cgo
// +build !cgo

package rl

import (
	"unsafe"

	"github.com/gen2brain/raylib-go/raylib/internal/convert"
	"github.com/jupiterrider/ffi"
)

var (
	// Matrix operations

	rlMatrixMode          = dll.MustPrep("rlMatrixMode", &ffi.TypeVoid, &ffi.TypeSint32)
	rlPushMatrix          = dll.MustPrep("rlPushMatrix", &ffi.TypeVoid)
	rlPopMatrix           = dll.MustPrep("rlPopMatrix", &ffi.TypeVoid)
	rlLoadIdentity        = dll.MustPrep("rlLoadIdentity", &ffi.TypeVoid)
	rlTranslatef          = dll.MustPrep("rlTranslatef", &ffi.TypeVoid, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
	rlRotatef             = dll.MustPrep("rlRotatef", &ffi.TypeVoid, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
	rlScalef              = dll.MustPrep("rlScalef", &ffi.TypeVoid, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
	rlMultMatrixf         = dll.MustPrep("rlMultMatrixf", &ffi.TypeVoid, &ffi.TypePointer)
	rlFrustum             = dll.MustPrep("rlFrustum", &ffi.TypeVoid, &ffi.TypeDouble, &ffi.TypeDouble, &ffi.TypeDouble, &ffi.TypeDouble, &ffi.TypeDouble, &ffi.TypeDouble)
	rlOrtho               = dll.MustPrep("rlOrtho", &ffi.TypeVoid, &ffi.TypeDouble, &ffi.TypeDouble, &ffi.TypeDouble, &ffi.TypeDouble, &ffi.TypeDouble, &ffi.TypeDouble)
	rlViewport            = dll.MustPrep("rlViewport", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	rlSetClipPlanes       = dll.MustPrep("rlSetClipPlanes", &ffi.TypeVoid, &ffi.TypeDouble, &ffi.TypeDouble)
	rlGetCullDistanceNear = dll.MustPrep("rlGetCullDistanceNear", &ffi.TypeDouble)
	rlGetCullDistanceFar  = dll.MustPrep("rlGetCullDistanceFar", &ffi.TypeDouble)

	// Vertex level operations

	rlBegin      = dll.MustPrep("rlBegin", &ffi.TypeVoid, &ffi.TypeSint32)
	rlEnd        = dll.MustPrep("rlEnd", &ffi.TypeVoid)
	rlVertex2i   = dll.MustPrep("rlVertex2i", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32)
	rlVertex2f   = dll.MustPrep("rlVertex2f", &ffi.TypeVoid, &ffi.TypeFloat, &ffi.TypeFloat)
	rlVertex3f   = dll.MustPrep("rlVertex3f", &ffi.TypeVoid, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
	rlTexCoord2f = dll.MustPrep("rlTexCoord2f", &ffi.TypeVoid, &ffi.TypeFloat, &ffi.TypeFloat)
	rlNormal3f   = dll.MustPrep("rlNormal3f", &ffi.TypeVoid, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
	rlColor4ub   = dll.MustPrep("rlColor4ub", &ffi.TypeVoid, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8)
	rlColor3f    = dll.MustPrep("rlColor3f", &ffi.TypeVoid, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
	rlColor4f    = dll.MustPrep("rlColor4f", &ffi.TypeVoid, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)

	// Vertex buffers state

	rlEnableVertexArray          = dll.MustPrep("rlEnableVertexArray", &ffi.TypeUint8, &ffi.TypeUint32)
	rlDisableVertexArray         = dll.MustPrep("rlDisableVertexArray", &ffi.TypeVoid)
	rlEnableVertexBuffer         = dll.MustPrep("rlEnableVertexBuffer", &ffi.TypeVoid, &ffi.TypeUint32)
	rlDisableVertexBuffer        = dll.MustPrep("rlDisableVertexBuffer", &ffi.TypeVoid)
	rlEnableVertexBufferElement  = dll.MustPrep("rlEnableVertexBufferElement", &ffi.TypeVoid, &ffi.TypeUint32)
	rlDisableVertexBufferElement = dll.MustPrep("rlDisableVertexBufferElement", &ffi.TypeVoid)
	rlEnableVertexAttribute      = dll.MustPrep("rlEnableVertexAttribute", &ffi.TypeVoid, &ffi.TypeUint32)
	rlDisableVertexAttribute     = dll.MustPrep("rlDisableVertexAttribute", &ffi.TypeVoid, &ffi.TypeUint32)
	rlEnableStatePointer         = dll.MustPrep("rlEnableStatePointer", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer)
	rlDisableStatePointer        = dll.MustPrep("rlDisableStatePointer", &ffi.TypeVoid, &ffi.TypeSint32)

	// Textures state

	rlActiveTextureSlot     = dll.MustPrep("rlActiveTextureSlot", &ffi.TypeVoid, &ffi.TypeSint32)
	rlEnableTexture         = dll.MustPrep("rlEnableTexture", &ffi.TypeVoid, &ffi.TypeUint32)
	rlDisableTexture        = dll.MustPrep("rlDisableTexture", &ffi.TypeVoid)
	rlEnableTextureCubemap  = dll.MustPrep("rlEnableTextureCubemap", &ffi.TypeVoid, &ffi.TypeUint32)
	rlDisableTextureCubemap = dll.MustPrep("rlDisableTextureCubemap", &ffi.TypeVoid)
	rlTextureParameters     = dll.MustPrep("rlTextureParameters", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypeSint32, &ffi.TypeSint32)
	rlCubemapParameters     = dll.MustPrep("rlCubemapParameters", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypeSint32, &ffi.TypeSint32)

	// Shader state

	rlEnableShader  = dll.MustPrep("rlEnableShader", &ffi.TypeVoid, &ffi.TypeUint32)
	rlDisableShader = dll.MustPrep("rlDisableShader", &ffi.TypeVoid)

	// Framebuffer state

	rlEnableFramebuffer    = dll.MustPrep("rlEnableFramebuffer", &ffi.TypeVoid, &ffi.TypeUint32)
	rlDisableFramebuffer   = dll.MustPrep("rlDisableFramebuffer", &ffi.TypeVoid)
	rlGetActiveFramebuffer = dll.MustPrep("rlGetActiveFramebuffer", &ffi.TypeUint32)
	rlActiveDrawBuffers    = dll.MustPrep("rlActiveDrawBuffers", &ffi.TypeVoid, &ffi.TypeSint32)
	rlBlitFramebuffer      = dll.MustPrep("rlBlitFramebuffer", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	rlBindFramebuffer      = dll.MustPrep("rlBindFramebuffer", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypeUint32)

	// General render state

	rlEnableColorBlend        = dll.MustPrep("rlEnableColorBlend", &ffi.TypeVoid)
	rlDisableColorBlend       = dll.MustPrep("rlDisableColorBlend", &ffi.TypeVoid)
	rlEnableDepthTest         = dll.MustPrep("rlEnableDepthTest", &ffi.TypeVoid)
	rlDisableDepthTest        = dll.MustPrep("rlDisableDepthTest", &ffi.TypeVoid)
	rlEnableDepthMask         = dll.MustPrep("rlEnableDepthMask", &ffi.TypeVoid)
	rlDisableDepthMask        = dll.MustPrep("rlDisableDepthMask", &ffi.TypeVoid)
	rlEnableBackfaceCulling   = dll.MustPrep("rlEnableBackfaceCulling", &ffi.TypeVoid)
	rlDisableBackfaceCulling  = dll.MustPrep("rlDisableBackfaceCulling", &ffi.TypeVoid)
	rlColorMask               = dll.MustPrep("rlColorMask", &ffi.TypeVoid, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8)
	rlSetCullFace             = dll.MustPrep("rlSetCullFace", &ffi.TypeVoid, &ffi.TypeSint32)
	rlEnableScissorTest       = dll.MustPrep("rlEnableScissorTest", &ffi.TypeVoid)
	rlDisableScissorTest      = dll.MustPrep("rlDisableScissorTest", &ffi.TypeVoid)
	rlScissor                 = dll.MustPrep("rlScissor", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	rlEnablePointMode         = dll.MustPrep("rlEnablePointMode", &ffi.TypeVoid)
	rlDisablePointMode        = dll.MustPrep("rlDisablePointMode", &ffi.TypeVoid)
	rlSetPointSize            = dll.MustPrep("rlSetPointSize", &ffi.TypeVoid, &ffi.TypeFloat)
	rlGetPointSize            = dll.MustPrep("rlGetPointSize", &ffi.TypeFloat)
	rlEnableWireMode          = dll.MustPrep("rlEnableWireMode", &ffi.TypeVoid)
	rlDisableWireMode         = dll.MustPrep("rlDisableWireMode", &ffi.TypeVoid)
	rlSetLineWidth            = dll.MustPrep("rlSetLineWidth", &ffi.TypeVoid, &ffi.TypeFloat)
	rlGetLineWidth            = dll.MustPrep("rlGetLineWidth", &ffi.TypeFloat)
	rlEnableSmoothLines       = dll.MustPrep("rlEnableSmoothLines", &ffi.TypeVoid)
	rlDisableSmoothLines      = dll.MustPrep("rlDisableSmoothLines", &ffi.TypeVoid)
	rlEnableStereoRender      = dll.MustPrep("rlEnableStereoRender", &ffi.TypeVoid)
	rlDisableStereoRender     = dll.MustPrep("rlDisableStereoRender", &ffi.TypeVoid)
	rlIsStereoRenderEnabled   = dll.MustPrep("rlIsStereoRenderEnabled", &ffi.TypeUint8)
	rlClearColor              = dll.MustPrep("rlClearColor", &ffi.TypeVoid, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8)
	rlClearScreenBuffers      = dll.MustPrep("rlClearScreenBuffers", &ffi.TypeVoid)
	rlCheckErrors             = dll.MustPrep("rlCheckErrors", &ffi.TypeVoid)
	rlSetBlendMode            = dll.MustPrep("rlSetBlendMode", &ffi.TypeVoid, &ffi.TypeSint32)
	rlSetBlendFactors         = dll.MustPrep("rlSetBlendFactors", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	rlSetBlendFactorsSeparate = dll.MustPrep("rlSetBlendFactorsSeparate", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)

	// rlgl initialization functions

	rlglInit               = dll.MustPrep("rlglInit", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32)
	rlglClose              = dll.MustPrep("rlglClose", &ffi.TypeVoid)
	rlGetProcAddress       = dll.MustPrep("rlGetProcAddress", &ffi.TypePointer, &ffi.TypePointer)
	rlGetVersion           = dll.MustPrep("rlGetVersion", &ffi.TypeSint32)
	rlSetFramebufferWidth  = dll.MustPrep("rlSetFramebufferWidth", &ffi.TypeVoid, &ffi.TypeSint32)
	rlGetFramebufferWidth  = dll.MustPrep("rlGetFramebufferWidth", &ffi.TypeSint32)
	rlSetFramebufferHeight = dll.MustPrep("rlSetFramebufferHeight", &ffi.TypeVoid, &ffi.TypeSint32)
	rlGetFramebufferHeight = dll.MustPrep("rlGetFramebufferHeight", &ffi.TypeSint32)
	rlGetTextureIdDefault  = dll.MustPrep("rlGetTextureIdDefault", &ffi.TypeUint32)
	rlGetShaderIdDefault   = dll.MustPrep("rlGetShaderIdDefault", &ffi.TypeUint32)
	rlGetShaderLocsDefault = dll.MustPrep("rlGetShaderLocsDefault", &ffi.TypePointer)

	// Render batch management

	rlLoadRenderBatch       = dll.MustPrep("rlLoadRenderBatch", &typeRenderBatch, &ffi.TypeSint32, &ffi.TypeSint32)
	rlUnloadRenderBatch     = dll.MustPrep("rlUnloadRenderBatch", &ffi.TypeVoid, &typeRenderBatch)
	rlDrawRenderBatch       = dll.MustPrep("rlDrawRenderBatch", &ffi.TypeVoid, &ffi.TypePointer)
	rlSetRenderBatchActive  = dll.MustPrep("rlSetRenderBatchActive", &ffi.TypeVoid, &ffi.TypePointer)
	rlDrawRenderBatchActive = dll.MustPrep("rlDrawRenderBatchActive", &ffi.TypeVoid)
	rlCheckRenderBatchLimit = dll.MustPrep("rlCheckRenderBatchLimit", &ffi.TypeUint8, &ffi.TypeSint32)
	rlSetTexture            = dll.MustPrep("rlSetTexture", &ffi.TypeVoid, &ffi.TypeUint32)

	// Vertex buffers management

	rlLoadVertexArray                  = dll.MustPrep("rlLoadVertexArray", &ffi.TypeUint32)
	rlLoadVertexBuffer                 = dll.MustPrep("rlLoadVertexBuffer", &ffi.TypeUint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeUint8)
	rlLoadVertexBufferElement          = dll.MustPrep("rlLoadVertexBufferElement", &ffi.TypeUint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeUint8)
	rlUpdateVertexBuffer               = dll.MustPrep("rlUpdateVertexBuffer", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32)
	rlUpdateVertexBufferElements       = dll.MustPrep("rlUpdateVertexBufferElements", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32)
	rlUnloadVertexArray                = dll.MustPrep("rlUnloadVertexArray", &ffi.TypeVoid, &ffi.TypeUint32)
	rlUnloadVertexBuffer               = dll.MustPrep("rlUnloadVertexBuffer", &ffi.TypeVoid, &ffi.TypeUint32)
	rlSetVertexAttribute               = dll.MustPrep("rlSetVertexAttribute", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeUint8, &ffi.TypeSint32, &ffi.TypeSint32)
	rlSetVertexAttributeDivisor        = dll.MustPrep("rlSetVertexAttributeDivisor", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypeSint32)
	rlSetVertexAttributeDefault        = dll.MustPrep("rlSetVertexAttributeDefault", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32)
	rlDrawVertexArray                  = dll.MustPrep("rlDrawVertexArray", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32)
	rlDrawVertexArrayElements          = dll.MustPrep("rlDrawVertexArrayElements", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer)
	rlDrawVertexArrayInstanced         = dll.MustPrep("rlDrawVertexArrayInstanced", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	rlDrawVertexArrayElementsInstanced = dll.MustPrep("rlDrawVertexArrayElementsInstanced", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32)

	// Textures management

	rlLoadTextureDepth = dll.MustPrep("rlLoadTextureDepth", &ffi.TypeUint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeUint8)

	// Framebuffer management (fbo)

	rlLoadFramebuffer     = dll.MustPrep("rlLoadFramebuffer", &ffi.TypeUint32)
	rlFramebufferAttach   = dll.MustPrep("rlFramebufferAttach", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	rlFramebufferComplete = dll.MustPrep("rlFramebufferComplete", &ffi.TypeUint8, &ffi.TypeUint32)
	rlUnloadFramebuffer   = dll.MustPrep("rlUnloadFramebuffer", &ffi.TypeVoid, &ffi.TypeUint32)
	rlCopyFramebuffer     = dll.MustPrep("rlCopyFramebuffer", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer)
	rlResizeFramebuffer   = dll.MustPrep("rlResizeFramebuffer", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32)

	// Shaders management

	rlLoadShader               = dll.MustPrep("rlLoadShader", &ffi.TypeUint32, &ffi.TypePointer, &ffi.TypeSint32)
	rlLoadShaderProgram        = dll.MustPrep("rlLoadShaderProgram", &ffi.TypeUint32, &ffi.TypePointer, &ffi.TypePointer)
	rlLoadShaderProgramEx      = dll.MustPrep("rlLoadShaderProgramEx", &ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeUint32)
	rlLoadShaderProgramCompute = dll.MustPrep("rlLoadShaderProgramCompute", &ffi.TypeUint32, &ffi.TypeUint32)
	rlUnloadShader             = dll.MustPrep("rlUnloadShader", &ffi.TypeVoid, &ffi.TypeUint32)
	rlUnloadShaderProgram      = dll.MustPrep("rlUnloadShaderProgram", &ffi.TypeVoid, &ffi.TypeUint32)
	rlGetLocationUniform       = dll.MustPrep("rlGetLocationUniform", &ffi.TypeSint32, &ffi.TypeUint32, &ffi.TypePointer)
	rlGetLocationAttrib        = dll.MustPrep("rlGetLocationAttrib", &ffi.TypeSint32, &ffi.TypeUint32, &ffi.TypePointer)
	rlSetUniform               = dll.MustPrep("rlSetUniform", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32)
	rlSetUniformMatrix         = dll.MustPrep("rlSetUniformMatrix", &ffi.TypeVoid, &ffi.TypeSint32, &typeMatrix)
	rlSetUniformMatrices       = dll.MustPrep("rlSetUniformMatrices", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32)
	rlSetUniformSampler        = dll.MustPrep("rlSetUniformSampler", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeUint32)
	rlSetShader                = dll.MustPrep("rlSetShader", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypePointer)

	// Compute shader management

	rlComputeShaderDispatch = dll.MustPrep("rlComputeShaderDispatch", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeUint32)

	// Shader buffer storage object management (ssbo)

	rlLoadShaderBuffer    = dll.MustPrep("rlLoadShaderBuffer", &ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypePointer, &ffi.TypeSint32)
	rlUnloadShaderBuffer  = dll.MustPrep("rlUnloadShaderBuffer", &ffi.TypeVoid, &ffi.TypeUint32)
	rlUpdateShaderBuffer  = dll.MustPrep("rlUpdateShaderBuffer", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypePointer, &ffi.TypeUint32, &ffi.TypeUint32)
	rlBindShaderBuffer    = dll.MustPrep("rlBindShaderBuffer", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypeUint32)
	rlReadShaderBuffer    = dll.MustPrep("rlReadShaderBuffer", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypePointer, &ffi.TypeUint32, &ffi.TypeUint32)
	rlCopyShaderBuffer    = dll.MustPrep("rlCopyShaderBuffer", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeUint32)
	rlGetShaderBufferSize = dll.MustPrep("rlGetShaderBufferSize", &ffi.TypeUint32, &ffi.TypeUint32)

	// Buffer management

	rlBindImageTexture = dll.MustPrep("rlBindImageTexture", &ffi.TypeVoid, &ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeSint32, &ffi.TypeUint8)

	// Matrix state management

	rlGetMatrixModelview        = dll.MustPrep("rlGetMatrixModelview", &typeMatrix)
	rlGetMatrixProjection       = dll.MustPrep("rlGetMatrixProjection", &typeMatrix)
	rlGetMatrixTransform        = dll.MustPrep("rlGetMatrixTransform", &typeMatrix)
	rlGetMatrixProjectionStereo = dll.MustPrep("rlGetMatrixProjectionStereo", &typeMatrix, &ffi.TypeSint32)
	rlGetMatrixViewOffsetStereo = dll.MustPrep("rlGetMatrixViewOffsetStereo", &typeMatrix, &ffi.TypeSint32)
	rlSetMatrixProjection       = dll.MustPrep("rlSetMatrixProjection", &ffi.TypeVoid, &typeMatrix)
	rlSetMatrixModelview        = dll.MustPrep("rlSetMatrixModelview", &ffi.TypeVoid, &typeMatrix)
	rlSetMatrixProjectionStereo = dll.MustPrep("rlSetMatrixProjectionStereo", &ffi.TypeVoid, &typeMatrix, &typeMatrix)
	rlSetMatrixViewOffsetStereo = dll.MustPrep("rlSetMatrixViewOffsetStereo", &ffi.TypeVoid, &typeMatrix, &typeMatrix)

	// Quick and dirty cube/quad buffers load->draw->unload

	rlLoadDrawCube = dll.MustPrep("rlLoadDrawCube", &ffi.TypeVoid)
	rlLoadDrawQuad = dll.MustPrep("rlLoadDrawQuad", &ffi.TypeVoid)
)

// MatrixMode - Choose the current matrix to be transformed
func MatrixMode(mode int32) {
	rlMatrixMode.Call(nil, &mode)
}

// PushMatrix - Push the current matrix to stack
func PushMatrix() {
	rlPushMatrix.Call(nil)
}

// PopMatrix - Pop lattest inserted matrix from stack
func PopMatrix() {
	rlPopMatrix.Call(nil)
}

// LoadIdentity - Reset current matrix to identity matrix
func LoadIdentity() {
	rlLoadIdentity.Call(nil)
}

// Translatef - Multiply the current matrix by a translation matrix
func Translatef(x float32, y float32, z float32) {
	rlTranslatef.Call(nil, &x, &y, &z)
}

// Rotatef - Multiply the current matrix by a rotation matrix
func Rotatef(angle float32, x float32, y float32, z float32) {
	rlRotatef.Call(nil, &angle, &x, &y, &z)
}

// Scalef - Multiply the current matrix by a scaling matrix
func Scalef(x float32, y float32, z float32) {
	rlScalef.Call(nil, &x, &y, &z)
}

// MultMatrix - Multiply the current matrix by another matrix
func MultMatrix(m Matrix) {
	f := unsafe.SliceData(MatrixToFloat(m))
	rlMultMatrixf.Call(nil, &f)
}

// Frustum .
func Frustum(left float64, right float64, bottom float64, top float64, znear float64, zfar float64) {
	rlFrustum.Call(nil, &left, &right, &bottom, &top, &znear, &zfar)
}

// Ortho .
func Ortho(left float64, right float64, bottom float64, top float64, znear float64, zfar float64) {
	rlOrtho.Call(nil, &left, &right, &bottom, &top, &znear, &zfar)
}

// Viewport - Set the viewport area
func Viewport(x int32, y int32, width int32, height int32) {
	rlViewport.Call(nil, &x, &y, &width, &height)
}

// SetClipPlanes - Set clip planes distances
func SetClipPlanes(nearPlane, farPlane float64) {
	rlSetClipPlanes.Call(nil, &nearPlane, &farPlane)
}

// GetCullDistanceNear - Get cull plane distance near
func GetCullDistanceNear() float64 {
	var ret float64
	rlGetCullDistanceNear.Call(&ret)
	return ret
}

// GetCullDistanceFar - Get cull plane distance far
func GetCullDistanceFar() float64 {
	var ret float64
	rlGetCullDistanceFar.Call(&ret)
	return ret
}

// Begin - Initialize drawing mode (how to organize vertex)
func Begin(mode int32) {
	rlBegin.Call(nil, &mode)
}

// End - Finish vertex providing
func End() {
	rlEnd.Call(nil)
}

// Vertex2i - Define one vertex (position) - 2 int
func Vertex2i(x int32, y int32) {
	rlVertex2i.Call(nil, &x, &y)
}

// Vertex2f - Define one vertex (position) - 2 float
func Vertex2f(x float32, y float32) {
	rlVertex2f.Call(nil, &x, &y)
}

// Vertex3f - Define one vertex (position) - 3 float
func Vertex3f(x float32, y float32, z float32) {
	rlVertex3f.Call(nil, &x, &y, &z)
}

// TexCoord2f - Define one vertex (texture coordinate) - 2 float
func TexCoord2f(x float32, y float32) {
	rlTexCoord2f.Call(nil, &x, &y)
}

// Normal3f - Define one vertex (normal) - 3 float
func Normal3f(x float32, y float32, z float32) {
	rlNormal3f.Call(nil, &x, &y, &z)
}

// Color4ub - Define one vertex (color) - 4 byte
func Color4ub(r uint8, g uint8, b uint8, a uint8) {
	rlColor4ub.Call(nil, &r, &g, &b, &a)
}

// Color3f - Define one vertex (color) - 3 float
func Color3f(x float32, y float32, z float32) {
	rlColor3f.Call(nil, &x, &y, &z)
}

// Color4f - Define one vertex (color) - 4 float
func Color4f(x float32, y float32, z float32, w float32) {
	rlColor4f.Call(nil, &x, &y, &z, &w)
}

// EnableVertexArray - Enable vertex array (VAO, if supported)
func EnableVertexArray(vaoId uint32) bool {
	var ret ffi.Arg
	rlEnableVertexArray.Call(&ret, &vaoId)
	return ret.Bool()
}

// DisableVertexArray - Disable vertex array (VAO, if supported)
func DisableVertexArray() {
	rlDisableVertexArray.Call(nil)
}

// EnableVertexBuffer - Enable vertex buffer (VBO)
func EnableVertexBuffer(id uint32) {
	rlEnableVertexBuffer.Call(nil, &id)
}

// DisableVertexBuffer - Disable vertex buffer (VBO)
func DisableVertexBuffer() {
	rlDisableVertexBuffer.Call(nil)
}

// EnableVertexBufferElement - Enable vertex buffer element (VBO element)
func EnableVertexBufferElement(id uint32) {
	rlEnableVertexBufferElement.Call(nil, &id)
}

// DisableVertexBufferElement - Disable vertex buffer element (VBO element)
func DisableVertexBufferElement() {
	rlDisableVertexBufferElement.Call(nil)
}

// EnableVertexAttribute - Enable vertex attribute index
func EnableVertexAttribute(index uint32) {
	rlEnableVertexAttribute.Call(nil, &index)
}

// DisableVertexAttribute - Disable vertex attribute index
func DisableVertexAttribute(index uint32) {
	rlDisableVertexAttribute.Call(nil, &index)
}

// EnableStatePointer - Enable attribute state pointer
func EnableStatePointer(vertexAttribType int32, buffer unsafe.Pointer) {
	rlEnableStatePointer.Call(nil, &vertexAttribType, &buffer)
}

// DisableStatePointer - Disable attribute state pointer
func DisableStatePointer(vertexAttribType int32) {
	rlDisableStatePointer.Call(nil, &vertexAttribType)
}

// ActiveTextureSlot - Select and active a texture slot
func ActiveTextureSlot(slot int32) {
	rlActiveTextureSlot.Call(nil, &slot)
}

// EnableTexture - Enable texture
func EnableTexture(id uint32) {
	rlEnableTexture.Call(nil, &id)
}

// DisableTexture - Disable texture
func DisableTexture() {
	rlDisableTexture.Call(nil)
}

// EnableTextureCubemap - Enable texture cubemap
func EnableTextureCubemap(id uint32) {
	rlEnableTextureCubemap.Call(nil, &id)
}

// DisableTextureCubemap - Disable texture cubemap
func DisableTextureCubemap() {
	rlDisableTextureCubemap.Call(nil)
}

// TextureParameters - Set texture parameters (filter, wrap)
func TextureParameters(id uint32, param int32, value int32) {
	rlTextureParameters.Call(nil, &id, &param, &value)
}

// CubemapParameters - Set cubemap parameters (filter, wrap)
func CubemapParameters(id uint32, param int32, value int32) {
	rlCubemapParameters.Call(nil, &id, &param, &value)
}

// EnableShader - Enable shader program
func EnableShader(id uint32) {
	rlEnableShader.Call(nil, &id)
}

// DisableShader - Disable shader program
func DisableShader() {
	rlDisableShader.Call(nil)
}

// EnableFramebuffer - Enable render texture (fbo)
func EnableFramebuffer(id uint32) {
	rlEnableFramebuffer.Call(nil, &id)
}

// DisableFramebuffer - Disable render texture (fbo), return to default framebuffer
func DisableFramebuffer() {
	rlDisableFramebuffer.Call(nil)
}

// GetActiveFramebuffer - Get the currently active render texture (fbo), 0 for default framebuffer
func GetActiveFramebuffer() uint32 {
	var ret ffi.Arg
	rlGetActiveFramebuffer.Call(&ret)
	return uint32(ret)
}

// ActiveDrawBuffers - Activate multiple draw color buffers
func ActiveDrawBuffers(count int32) {
	rlActiveDrawBuffers.Call(nil, &count)
}

// BlitFramebuffer - Blit active framebuffer to main framebuffer
func BlitFramebuffer(srcX, srcY, srcWidth, srcHeight, dstX, dstY, dstWidth, dstHeight, bufferMask int32) {
	rlBlitFramebuffer.Call(nil, &srcX, &srcY, &srcWidth, &srcHeight, &dstX, &dstY, &dstWidth, &dstHeight, &bufferMask)
}

// BindFramebuffer - Bind framebuffer (FBO)
func BindFramebuffer(target, framebuffer uint32) {
	rlBindFramebuffer.Call(nil, &target, &framebuffer)
}

// EnableColorBlend - Enable color blending
func EnableColorBlend() {
	rlEnableColorBlend.Call(nil)
}

// DisableColorBlend - Disable color blending
func DisableColorBlend() {
	rlDisableColorBlend.Call(nil)
}

// EnableDepthTest - Enable depth test
func EnableDepthTest() {
	rlEnableDepthTest.Call(nil)
}

// DisableDepthTest - Disable depth test
func DisableDepthTest() {
	rlDisableDepthTest.Call(nil)
}

// EnableDepthMask - Enable depth write
func EnableDepthMask() {
	rlEnableDepthMask.Call(nil)
}

// DisableDepthMask - Disable depth write
func DisableDepthMask() {
	rlDisableDepthMask.Call(nil)
}

// EnableBackfaceCulling - Enable backface culling
func EnableBackfaceCulling() {
	rlEnableBackfaceCulling.Call(nil)
}

// DisableBackfaceCulling - Disable backface culling
func DisableBackfaceCulling() {
	rlDisableBackfaceCulling.Call(nil)
}

// ColorMask - Color mask control
func ColorMask(r, g, b, a bool) {
	rlColorMask.Call(nil, &r, &g, &b, &a)
}

// SetCullFace - Set face culling mode
func SetCullFace(mode int32) {
	rlSetCullFace.Call(nil, &mode)
}

// EnableScissorTest - Enable scissor test
func EnableScissorTest() {
	rlEnableScissorTest.Call(nil)
}

// DisableScissorTest - Disable scissor test
func DisableScissorTest() {
	rlDisableScissorTest.Call(nil)
}

// Scissor - Scissor test
func Scissor(x int32, y int32, width int32, height int32) {
	rlScissor.Call(nil, &x, &y, &width, &height)
}

// EnablePointMode - Enable point mode
func EnablePointMode() {
	rlEnablePointMode.Call(nil)
}

// DisablePointMode - Disable point mode
func DisablePointMode() {
	rlDisablePointMode.Call(nil)
}

// SetPointSize - Set the point drawing size
func SetPointSize(size float32) {
	rlSetPointSize.Call(nil, &size)
}

// GetPointSize - Get the point drawing size
func GetPointSize() float32 {
	var ret float32
	rlGetPointSize.Call(&ret)
	return ret
}

// EnableWireMode - Enable wire mode
func EnableWireMode() {
	rlEnableWireMode.Call(nil)
}

// DisableWireMode - Disable wire mode
func DisableWireMode() {
	rlDisableWireMode.Call(nil)
}

// SetLineWidth - Set the line drawing width
func SetLineWidth(width float32) {
	rlSetLineWidth.Call(nil, &width)
}

// GetLineWidth - Get the line drawing width
func GetLineWidth() float32 {
	var ret float32
	rlGetLineWidth.Call(&ret)
	return ret
}

// EnableSmoothLines - Enable line aliasing
func EnableSmoothLines() {
	rlEnableSmoothLines.Call(nil)
}

// DisableSmoothLines - Disable line aliasing
func DisableSmoothLines() {
	rlDisableSmoothLines.Call(nil)
}

// EnableStereoRender - Enable stereo rendering
func EnableStereoRender() {
	rlEnableStereoRender.Call(nil)
}

// DisableStereoRender - Disable stereo rendering
func DisableStereoRender() {
	rlDisableStereoRender.Call(nil)
}

// IsStereoRenderEnabled - Check if stereo render is enabled
func IsStereoRenderEnabled() bool {
	var ret ffi.Arg
	rlIsStereoRenderEnabled.Call(&ret)
	return ret.Bool()
}

// ClearColor - Clear color buffer with color
func ClearColor(r, g, b, a uint8) {
	rlClearColor.Call(nil, &r, &g, &b, &a)
}

// ClearScreenBuffers - Clear used screen buffers (color and depth)
func ClearScreenBuffers() {
	rlClearScreenBuffers.Call(nil)
}

// CheckErrors - Check and log OpenGL error codes
func CheckErrors() {
	rlCheckErrors.Call(nil)
}

// SetBlendMode - Set blending mode
func SetBlendMode(mode BlendMode) {
	rlSetBlendMode.Call(nil, &mode)
}

// SetBlendFactors - Set blending mode factor and equation (using OpenGL factors)
func SetBlendFactors(glSrcFactor, glDstFactor, glEquation int32) {
	rlSetBlendFactors.Call(nil, &glSrcFactor, &glDstFactor, &glEquation)
}

// SetBlendFactorsSeparate - Set blending mode factors and equations separately (using OpenGL factors)
func SetBlendFactorsSeparate(glSrcRGB, glDstRGB, glSrcAlpha, glDstAlpha, glEqRGB, glEqAlpha int32) {
	rlSetBlendFactorsSeparate.Call(nil, &glSrcRGB, &glDstRGB, &glSrcAlpha, &glDstAlpha, &glEqRGB, &glEqAlpha)
}

// GlInit - Initialize rlgl (buffers, shaders, textures, states)
func GlInit(width int32, height int32) {
	rlglInit.Call(nil, &width, &height)
}

// GlClose - De-inititialize rlgl (buffers, shaders, textures)
func GlClose() {
	rlglClose.Call(nil)
}

// GetProcAddress - Get OpenGL procedure address
func GetProcAddress(procName string) unsafe.Pointer {
	procNamePtr := convert.ToBytePtr(procName)
	var ret unsafe.Pointer
	rlGetProcAddress.Call(&ret, &procNamePtr)
	return ret
}

// GetVersion - Get current OpenGL version
func GetVersion() int32 {
	var ret ffi.Arg
	rlGetVersion.Call(&ret)
	return int32(ret)
}

// SetFramebufferWidth - Set current framebuffer width
func SetFramebufferWidth(width int32) {
	rlSetFramebufferWidth.Call(nil, &width)
}

// GetFramebufferWidth - Get default framebuffer width
func GetFramebufferWidth() int32 {
	var ret ffi.Arg
	rlGetFramebufferWidth.Call(&ret)
	return int32(ret)
}

// SetFramebufferHeight - Set current framebuffer height
func SetFramebufferHeight(height int32) {
	rlSetFramebufferHeight.Call(nil, &height)
}

// GetFramebufferHeight - Get default framebuffer height
func GetFramebufferHeight() int32 {
	var ret ffi.Arg
	rlGetFramebufferHeight.Call(&ret)
	return int32(ret)
}

// GetTextureIdDefault - Get default texture id
func GetTextureIdDefault() uint32 {
	var ret ffi.Arg
	rlGetTextureIdDefault.Call(&ret)
	return uint32(ret)
}

// GetShaderIdDefault - Get default shader id
func GetShaderIdDefault() uint32 {
	var ret ffi.Arg
	rlGetShaderIdDefault.Call(&ret)
	return uint32(ret)
}

// GetShaderLocsDefault - Get default shader locations
func GetShaderLocsDefault() []int32 {
	var ret *int32
	rlGetShaderLocsDefault.Call(&ret)
	// the default value of RL_MAX_SHADER_LOCATIONS is 32
	return unsafe.Slice(ret, 32)
}

// LoadRenderBatch - Load a render batch system
func LoadRenderBatch(numBuffers int32, bufferElements int32) RenderBatch {
	var ret RenderBatch
	rlLoadRenderBatch.Call(&ret, &numBuffers, &bufferElements)
	return ret
}

// UnloadRenderBatch - Unload render batch system
func UnloadRenderBatch(batch RenderBatch) {
	rlUnloadRenderBatch.Call(nil, &batch)
}

// DrawRenderBatch - Draw render batch data (Update->Draw->Reset)
func DrawRenderBatch(batch *RenderBatch) {
	rlDrawRenderBatch.Call(nil, &batch)
}

// SetRenderBatchActive - Set the active render batch for rlgl (NULL for default internal)
func SetRenderBatchActive(batch *RenderBatch) {
	rlSetRenderBatchActive.Call(nil, &batch)
}

// DrawRenderBatchActive - Update and draw internal render batch
func DrawRenderBatchActive() {
	rlDrawRenderBatchActive.Call(nil)
}

// CheckRenderBatchLimit - Check internal buffer overflow for a given number of vertex
func CheckRenderBatchLimit(vCount int32) bool {
	var ret ffi.Arg
	rlCheckRenderBatchLimit.Call(&ret, &vCount)
	return ret.Bool()
}

// SetTexture - Set current texture for render batch and check buffers limits
func SetTexture(id uint32) {
	rlSetTexture.Call(nil, &id)
}

// LoadVertexArray - Load vertex array (vao) if supported
func LoadVertexArray() uint32 {
	var ret ffi.Arg
	rlLoadVertexArray.Call(&ret)
	return uint32(ret)
}

// LoadVertexBuffer - Load a vertex buffer object
func LoadVertexBuffer[T any](buffer []T, dynamic bool) uint32 {
	if len(buffer) == 0 {
		return 0
	}
	size := int32(int(unsafe.Sizeof(buffer[0])) * len(buffer))
	bufferPtr := unsafe.SliceData(buffer)
	var ret ffi.Arg
	rlLoadVertexBuffer.Call(&ret, &bufferPtr, &size, &dynamic)
	return uint32(ret)
}

// LoadVertexBufferElement - Load vertex buffer elements object
func LoadVertexBufferElement[T any](buffer []T, dynamic bool) uint32 {
	if len(buffer) == 0 {
		return 0
	}
	size := int32(int(unsafe.Sizeof(buffer[0])) * len(buffer))
	bufferPtr := unsafe.SliceData(buffer)
	var ret ffi.Arg
	rlLoadVertexBufferElement.Call(&ret, &bufferPtr, &size, &dynamic)
	return uint32(ret)
}

// UpdateVertexBuffer - Update vertex buffer object data on GPU buffer
func UpdateVertexBuffer[T any](bufferId uint32, data []T, offset int32) {
	if len(data) == 0 {
		return
	}
	dataSize := int32(int(unsafe.Sizeof(data[0])) * len(data))
	dataPtr := unsafe.SliceData(data)
	rlUpdateVertexBuffer.Call(nil, &bufferId, &dataPtr, &dataSize, &offset)
}

// UpdateVertexBufferElements - Update vertex buffer elements data on GPU buffer
func UpdateVertexBufferElements[T any](id uint32, data []T, offset int32) {
	if len(data) == 0 {
		return
	}
	dataSize := int32(int(unsafe.Sizeof(data[0])) * len(data))
	dataPtr := unsafe.SliceData(data)
	rlUpdateVertexBufferElements.Call(nil, &id, &dataPtr, &dataSize, &offset)
}

// UnloadVertexArray - Unload vertex array (vao)
func UnloadVertexArray(vaoId uint32) {
	rlUnloadVertexArray.Call(nil, &vaoId)
}

// UnloadVertexBuffer - Unload vertex buffer object
func UnloadVertexBuffer(vboId uint32) {
	rlUnloadVertexBuffer.Call(nil, &vboId)
}

// SetVertexAttribute - Set vertex attribute data configuration
func SetVertexAttribute(index uint32, compSize int32, attrType int32, normalized bool, stride int32, offset int32) {
	rlSetVertexAttribute.Call(nil, &index, &compSize, &attrType, &normalized, &stride, &offset)
}

// SetVertexAttributeDivisor - Set vertex attribute data divisor
func SetVertexAttributeDivisor(index uint32, divisor int32) {
	rlSetVertexAttributeDivisor.Call(nil, &index, &divisor)
}

// SetVertexAttributeDefault - Set vertex attribute default value, when attribute to provided
func SetVertexAttributeDefault(locIndex int32, value unsafe.Pointer, attribType int32, count int32) {
	rlSetVertexAttributeDefault.Call(nil, &locIndex, &value, &attribType, &count)
}

// DrawVertexArray - Draw vertex array (currently active vao)
func DrawVertexArray(offset int32, count int32) {
	rlDrawVertexArray.Call(nil, &offset, &count)
}

// DrawVertexArrayElements - Draw vertex array elements
func DrawVertexArrayElements(offset int32, count int32, buffer unsafe.Pointer) {
	rlDrawVertexArrayElements.Call(nil, &offset, &count, &buffer)
}

// DrawVertexArrayInstanced - Draw vertex array (currently active vao) with instancing
func DrawVertexArrayInstanced(offset, count, instances int32) {
	rlDrawVertexArrayInstanced.Call(nil, &offset, &count, &instances)
}

// DrawVertexArrayElementsInstanced - Draw vertex array elements with instancing
func DrawVertexArrayElementsInstanced(offset int32, count int32, buffer unsafe.Pointer, instances int32) {
	rlDrawVertexArrayElementsInstanced.Call(nil, &offset, &count, &buffer, &instances)
}

// LoadTextureDepth - Load depth texture/renderbuffer (to be attached to fbo)
func LoadTextureDepth(width, height int32, useRenderBuffer bool) uint32 {
	var ret ffi.Arg
	rlLoadTextureDepth.Call(&ret, &width, &height, &useRenderBuffer)
	return uint32(ret)
}

// LoadFramebuffer - Load an empty framebuffer
func LoadFramebuffer() uint32 {
	var ret ffi.Arg
	rlLoadFramebuffer.Call(&ret)
	return uint32(ret)
}

// FramebufferAttach - Attach texture/renderbuffer to a framebuffer
func FramebufferAttach(id, texId uint32, attachType, texType, mipLevel int32) {
	rlFramebufferAttach.Call(nil, &id, &texId, &attachType, &texType, &mipLevel)
}

// FramebufferComplete - Verify framebuffer is complete
func FramebufferComplete(id uint32) bool {
	var ret ffi.Arg
	rlFramebufferComplete.Call(&ret, &id)
	return ret.Bool()
}

// UnloadFramebuffer - Delete framebuffer from GPU
func UnloadFramebuffer(id uint32) {
	rlUnloadFramebuffer.Call(nil, &id)
}

// CopyFramebuffer - Copy framebuffer pixel data to internal buffer
//
// WARNING: Copy and resize framebuffer functionality only defined for software backend
func CopyFramebuffer(x, y, width, height, format int32, pixels unsafe.Pointer) {
	rlCopyFramebuffer.Call(nil, &x, &y, &width, &height, &format, &pixels)
}

// ResizeFramebuffer - Resize internal framebuffer
//
// WARNING: Copy and resize framebuffer functionality only defined for software backend
func ResizeFramebuffer(width, height int32) {
	rlResizeFramebuffer.Call(nil, &width, &height)
}

// LoadShaderId - Load (compile) shader and return shader id (type: [VertexShader], [FragmentShader], [ComputeShader])
func LoadShaderId(code string, shaderType int32) uint32 {
	codePtr := convert.ToBytePtr(code)
	var ret ffi.Arg
	rlLoadShader.Call(&ret, &codePtr, &shaderType)
	return uint32(ret)
}

// LoadShaderProgram - Load shader from code strings
func LoadShaderProgram(vsCode, fsCode string) uint32 {
	vsCodePtr := convert.ToBytePtrNullable(vsCode)
	fsCodePtr := convert.ToBytePtrNullable(fsCode)
	var ret ffi.Arg
	rlLoadShaderProgram.Call(&ret, &vsCodePtr, &fsCodePtr)
	return uint32(ret)
}

// LoadShaderProgramEx - Load shader program, using already loaded shader ids
func LoadShaderProgramEx(vsId, fsId uint32) uint32 {
	var ret ffi.Arg
	rlLoadShaderProgramEx.Call(&ret, &vsId, &fsId)
	return uint32(ret)
}

// LoadShaderProgramCompute - Load compute shader program
func LoadShaderProgramCompute(csId uint32) uint32 {
	var ret ffi.Arg
	rlLoadShaderProgramCompute.Call(&ret, &csId)
	return uint32(ret)
}

// UnloadShaderId - Unload shader, loaded with [LoadShaderId]
func UnloadShaderId(id uint32) {
	rlUnloadShader.Call(nil, &id)
}

// UnloadShaderProgram - Unload shader program
func UnloadShaderProgram(id uint32) {
	rlUnloadShaderProgram.Call(nil, &id)
}

// GetLocationUniform - Get shader location uniform, requires shader program id
func GetLocationUniform(id uint32, uniformName string) int32 {
	uniformNamePtr := convert.ToBytePtr(uniformName)
	var ret ffi.Arg
	rlGetLocationUniform.Call(&ret, &id, &uniformNamePtr)
	return int32(ret)
}

// GetLocationAttrib - Get shader location attribute, requires shader program id
func GetLocationAttrib(id uint32, attribName string) int32 {
	attribNamePtr := convert.ToBytePtr(attribName)
	var ret ffi.Arg
	rlGetLocationAttrib.Call(&ret, &id, &attribNamePtr)
	return int32(ret)
}

// SetUniform - Set shader value uniform ([]float32, []int32, []uint32)
func SetUniform[T ~float32 | ~int32 | ~uint32](locIndex int32, value []T, uniformType, count int32) {
	valuePtr := unsafe.SliceData(value)
	rlSetUniform.Call(nil, &locIndex, &valuePtr, &uniformType, &count)
}

// SetUniformMatrix - Set shader value matrix
func SetUniformMatrix(locIndex int32, mat Matrix) {
	rlSetUniformMatrix.Call(nil, &locIndex, &mat)
}

// SetUniformMatrices - Set shader value matrices
func SetUniformMatrices(locIndex int32, mat []Matrix) {
	matPtr := unsafe.SliceData(mat)
	count := int32(len(mat))
	rlSetUniformMatrices.Call(nil, &locIndex, &matPtr, &count)
}

// SetUniformSampler - Set shader value sampler
func SetUniformSampler(locIndex int32, textureId uint32) {
	rlSetUniformSampler.Call(nil, &locIndex, &textureId)
}

// SetShader - Set shader currently active (id and locations)
func SetShader(id uint32, locs *int32) {
	rlSetShader.Call(nil, &id, &locs)
}

// ComputeShaderDispatch - Dispatch compute shader (equivalent to *draw* for graphics pilepine)
func ComputeShaderDispatch(groupX uint32, groupY uint32, groupZ uint32) {
	rlComputeShaderDispatch.Call(nil, &groupX, &groupY, &groupZ)
}

// LoadShaderBuffer - Load shader storage buffer object (SSBO)
func LoadShaderBuffer(size uint32, data unsafe.Pointer, usageHint int32) uint32 {
	var ret ffi.Arg
	rlLoadShaderBuffer.Call(&ret, &size, &data, &usageHint)
	return uint32(ret)
}

// UnloadShaderBuffer - Unload shader storage buffer object (SSBO)
func UnloadShaderBuffer(id uint32) {
	rlUnloadShaderBuffer.Call(nil, &id)
}

// UpdateShaderBuffer - Update SSBO buffer data
func UpdateShaderBuffer(id uint32, data unsafe.Pointer, dataSize uint32, offset uint32) {
	rlUpdateShaderBuffer.Call(nil, &id, &data, &dataSize, &offset)
}

// BindShaderBuffer - Bind SSBO buffer
func BindShaderBuffer(id uint32, index uint32) {
	rlBindShaderBuffer.Call(nil, &id, &index)
}

// ReadShaderBuffer - Read SSBO buffer data (GPU->CPU)
func ReadShaderBuffer(id uint32, dest unsafe.Pointer, count uint32, offset uint32) {
	rlReadShaderBuffer.Call(nil, &id, &dest, &count, &offset)
}

// CopyShaderBuffer - Copy SSBO data between buffers
func CopyShaderBuffer(destId uint32, srcId uint32, destOffset uint32, srcOffset uint32, count uint32) {
	rlCopyShaderBuffer.Call(nil, &destId, &srcId, &destOffset, &srcOffset, &count)
}

// GetShaderBufferSize - Get SSBO buffer size
func GetShaderBufferSize(id uint32) uint32 {
	var ret ffi.Arg
	rlGetShaderBufferSize.Call(&ret, &id)
	return uint32(ret)
}

// BindImageTexture - Bind image texture
func BindImageTexture(id uint32, index uint32, format int32, readonly bool) {
	rlBindImageTexture.Call(nil, &id, &index, &format, &readonly)
}

// GetMatrixModelview - Get internal modelview matrix
func GetMatrixModelview() Matrix {
	var ret Matrix
	rlGetMatrixModelview.Call(&ret)
	return ret
}

// GetMatrixProjection - Get internal projection matrix
func GetMatrixProjection() Matrix {
	var ret Matrix
	rlGetMatrixProjection.Call(&ret)
	return ret
}

// GetMatrixTransform - Get internal accumulated transform matrix
func GetMatrixTransform() Matrix {
	var ret Matrix
	rlGetMatrixTransform.Call(&ret)
	return ret
}

// GetMatrixProjectionStereo - Get internal projection matrix for stereo render (selected eye)
func GetMatrixProjectionStereo(eye int32) Matrix {
	var ret Matrix
	rlGetMatrixProjectionStereo.Call(&ret, &eye)
	return ret
}

// GetMatrixViewOffsetStereo - Get internal view offset matrix for stereo render (selected eye)
func GetMatrixViewOffsetStereo(eye int32) Matrix {
	var ret Matrix
	rlGetMatrixViewOffsetStereo.Call(&ret, &eye)
	return ret
}

// SetMatrixProjection - Set a custom projection matrix (replaces internal projection matrix)
func SetMatrixProjection(proj Matrix) {
	rlSetMatrixProjection.Call(nil, &proj)
}

// SetMatrixModelview - Set a custom modelview matrix (replaces internal modelview matrix)
func SetMatrixModelview(view Matrix) {
	rlSetMatrixModelview.Call(nil, &view)
}

// SetMatrixProjectionStereo - Set eyes projection matrices for stereo rendering
func SetMatrixProjectionStereo(right Matrix, left Matrix) {
	rlSetMatrixProjectionStereo.Call(nil, &right, &left)
}

// SetMatrixViewOffsetStereo - Set eyes view offsets matrices for stereo rendering
func SetMatrixViewOffsetStereo(right Matrix, left Matrix) {
	rlSetMatrixViewOffsetStereo.Call(nil, &right, &left)
}

// LoadDrawCube - Load and draw a cube
func LoadDrawCube() {
	rlLoadDrawCube.Call(nil)
}

// LoadDrawQuad - Load and draw a quad
func LoadDrawQuad() {
	rlLoadDrawQuad.Call(nil)
}
