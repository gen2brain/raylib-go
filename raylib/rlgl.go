package rl

/*
#include "raylib.h"
#include "rlgl.h"
#include <stdlib.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

// cptr returns C pointer
func (s *Shader) cptr() *C.Shader {
	return (*C.Shader)(unsafe.Pointer(s))
}

// LoadShader - Load a custom shader and bind default locations
func LoadShader(vsFileName string, fsFileName string) Shader {
	cvsFileName := C.CString(vsFileName)
	defer C.free(unsafe.Pointer(cvsFileName))

	cfsFileName := C.CString(fsFileName)
	defer C.free(unsafe.Pointer(cfsFileName))

	if vsFileName == "" {
		cvsFileName = nil
	}

	if fsFileName == "" {
		cfsFileName = nil
	}

	ret := C.LoadShader(cvsFileName, cfsFileName)
	v := newShaderFromPointer(unsafe.Pointer(&ret))

	return v
}

// LoadShaderFromMemory - Load shader from code strings and bind default locations
func LoadShaderFromMemory(vsCode string, fsCode string) Shader {
	cvsCode := C.CString(vsCode)
	defer C.free(unsafe.Pointer(cvsCode))

	cfsCode := C.CString(fsCode)
	defer C.free(unsafe.Pointer(cfsCode))

	if vsCode == "" {
		cvsCode = nil
	}

	if fsCode == "" {
		cfsCode = nil
	}

	ret := C.LoadShaderFromMemory(cvsCode, cfsCode)
	v := newShaderFromPointer(unsafe.Pointer(&ret))

	return v
}

// UnloadShader - Unload a custom shader from memory
func UnloadShader(shader Shader) {
	cshader := shader.cptr()
	C.UnloadShader(*cshader)
}

// GetShaderLocation - Get shader uniform location
func GetShaderLocation(shader Shader, uniformName string) int32 {
	cshader := shader.cptr()
	cuniformName := C.CString(uniformName)
	defer C.free(unsafe.Pointer(cuniformName))

	ret := C.GetShaderLocation(*cshader, cuniformName)
	v := (int32)(ret)
	return v
}

// GetShaderLocationAttrib - Get shader attribute location
func GetShaderLocationAttrib(shader Shader, attribName string) int32 {
	cshader := shader.cptr()
	cuniformName := C.CString(attribName)
	defer C.free(unsafe.Pointer(cuniformName))

	ret := C.GetShaderLocationAttrib(*cshader, cuniformName)
	v := (int32)(ret)
	return v
}

// SetShaderValue - Set shader uniform value (float)
func SetShaderValue(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	cvalue := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&value)).Data)
	cuniformType := (C.int)(uniformType)
	C.SetShaderValue(*cshader, clocIndex, cvalue, cuniformType)
}

// SetShaderValueV - Set shader uniform value (float)
func SetShaderValueV(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType, count int32) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	cvalue := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&value)).Data)
	cuniformType := (C.int)(uniformType)
	ccount := (C.int)(count)
	C.SetShaderValueV(*cshader, clocIndex, cvalue, cuniformType, ccount)
}

// SetShaderValueMatrix - Set shader uniform value (matrix 4x4)
func SetShaderValueMatrix(shader Shader, locIndex int32, mat Matrix) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	cmat := mat.cptr()
	C.SetShaderValueMatrix(*cshader, clocIndex, *cmat)
}

// SetShaderValueTexture - Set shader uniform value for texture (sampler2d)
func SetShaderValueTexture(shader Shader, locIndex int32, texture Texture2D) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	ctexture := texture.cptr()
	C.SetShaderValueTexture(*cshader, clocIndex, *ctexture)
}

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

const (
	// Texture parameters (equivalent to OpenGL defines)
	RL_TEXTURE_WRAP_S     = 0x2802 // GL_TEXTURE_WRAP_S
	RL_TEXTURE_WRAP_T     = 0x2803 // GL_TEXTURE_WRAP_T
	RL_TEXTURE_MAG_FILTER = 0x2800 // GL_TEXTURE_MAG_FILTER
	RL_TEXTURE_MIN_FILTER = 0x2801 // GL_TEXTURE_MIN_FILTER

	RL_TEXTURE_FILTER_NEAREST            = 0x2600 // GL_NEAREST
	RL_TEXTURE_FILTER_LINEAR             = 0x2601 // GL_LINEAR
	RL_TEXTURE_FILTER_MIP_NEAREST        = 0x2700 // GL_NEAREST_MIPMAP_NEAREST
	RL_TEXTURE_FILTER_NEAREST_MIP_LINEAR = 0x2702 // GL_NEAREST_MIPMAP_LINEAR
	RL_TEXTURE_FILTER_LINEAR_MIP_NEAREST = 0x2701 // GL_LINEAR_MIPMAP_NEAREST
	RL_TEXTURE_FILTER_MIP_LINEAR         = 0x2703 // GL_LINEAR_MIPMAP_LINEAR
	RL_TEXTURE_FILTER_ANISOTROPIC        = 0x3000 // Anisotropic filter (custom identifier)
	RL_TEXTURE_MIPMAP_BIAS_RATIO         = 0x4000 // Texture mipmap bias, percentage ratio (custom identifier)

	RL_TEXTURE_WRAP_REPEAT        = 0x2901 // GL_REPEAT
	RL_TEXTURE_WRAP_CLAMP         = 0x812F // GL_CLAMP_TO_EDGE
	RL_TEXTURE_WRAP_MIRROR_REPEAT = 0x8370 // GL_MIRRORED_REPEAT
	RL_TEXTURE_WRAP_MIRROR_CLAMP  = 0x8742 // GL_MIRROR_CLAMP_EXT

	// Matrix modes (equivalent to OpenGL)
	RL_MODELVIEW  = 0x1700 // GL_MODELVIEW
	RL_PROJECTION = 0x1701 // GL_PROJECTION
	RL_TEXTURE    = 0x1702 // GL_TEXTURE

	// Primitive assembly draw modes
	RL_LINES     = 0x0001 // GL_LINES
	RL_TRIANGLES = 0x0004 // GL_TRIANGLES
	RL_QUADS     = 0x0007 // GL_QUADS

	// GL equivalent data types
	RL_UNSIGNED_BYTE = 0x1401 // GL_UNSIGNED_BYTE
	RL_FLOAT         = 0x1406 // GL_FLOAT

	// Buffer usage hint
	RL_STREAM_DRAW  = 0x88E0 // GL_STREAM_DRAW
	RL_STREAM_READ  = 0x88E1 // GL_STREAM_READ
	RL_STREAM_COPY  = 0x88E2 // GL_STREAM_COPY
	RL_STATIC_DRAW  = 0x88E4 // GL_STATIC_DRAW
	RL_STATIC_READ  = 0x88E5 // GL_STATIC_READ
	RL_STATIC_COPY  = 0x88E6 // GL_STATIC_COPY
	RL_DYNAMIC_DRAW = 0x88E8 // GL_DYNAMIC_DRAW
	RL_DYNAMIC_READ = 0x88E9 // GL_DYNAMIC_READ
	RL_DYNAMIC_COPY = 0x88EA // GL_DYNAMIC_COPY

	// GL Shader type
	RL_FRAGMENT_SHADER = 0x8B30 // GL_FRAGMENT_SHADER
	RL_VERTEX_SHADER   = 0x8B31 // GL_VERTEX_SHADER
	RL_COMPUTE_SHADER  = 0x91B9 // GL_COMPUTE_SHADER
)

// rlVertexBuffer - Dynamic vertex buffers (position + texcoords + colors + indices arrays)
type rlVertexBuffer struct {
	elementCount int32
	vertices     []float32
	texcoords    []float32
	colors       []uint8
	indices      []uint32
	vaoId        uint32
	vboId        [4]uint32
}

// rlDrawCall - Draw call type
// NOTE: Only texture changes register a new draw, other state-change-related elements are not
// used at this moment (vaoId, shaderId, matrices), raylib just forces a batch draw call if any
// of those state-change happens (this is done in core module)
type rlDrawCall struct {
	mode            int32
	vertexCount     int32
	vertexAlignment int32
	textureId       uint32
}

// rlRenderBatch type
type rlRenderBatch struct {
	bufferCount   int32
	currentBuffer int32
	vertexBuffer  []rlVertexBuffer
	draws         []rlDrawCall
	drawCounter   int32
	currentDepth  float32
}

// OpenGL version
const (
	RL_OPENGL_11    int32 = 1
	RL_OPENGL_21    int32 = 2
	RL_OPENGL_33    int32 = 3
	RL_OPENGL_43    int32 = 4
	RL_OPENGL_ES_20 int32 = 5
)

type rlGlVersion = int32

// RL_LOG_ALL - Trace log level
// NOTE: Organized by priority level
const (
	RL_LOG_ALL     int32 = 0
	RL_LOG_TRACE   int32 = 1
	RL_LOG_DEBUG   int32 = 2
	RL_LOG_INFO    int32 = 3
	RL_LOG_WARNING int32 = 4
	RL_LOG_ERROR   int32 = 5
	RL_LOG_FATAL   int32 = 6
	RL_LOG_NONE    int32 = 7
)

type rlTraceLogLevel = int32

// Texture pixel formats
// NOTE: Support depends on OpenGL version
const (
	RL_PIXELFORMAT_UNCOMPRESSED_GRAYSCALE    int32 = 1
	RL_PIXELFORMAT_UNCOMPRESSED_GRAY_ALPHA   int32 = 2
	RL_PIXELFORMAT_UNCOMPRESSED_R5G6B5       int32 = 3
	RL_PIXELFORMAT_UNCOMPRESSED_R8G8B8       int32 = 4
	RL_PIXELFORMAT_UNCOMPRESSED_R5G5B5A1     int32 = 5
	RL_PIXELFORMAT_UNCOMPRESSED_R4G4B4A4     int32 = 6
	RL_PIXELFORMAT_UNCOMPRESSED_R8G8B8A8     int32 = 7
	RL_PIXELFORMAT_UNCOMPRESSED_R32          int32 = 8
	RL_PIXELFORMAT_UNCOMPRESSED_R32G32B32    int32 = 9
	RL_PIXELFORMAT_UNCOMPRESSED_R32G32B32A32 int32 = 10
	RL_PIXELFORMAT_COMPRESSED_DXT1_RGB       int32 = 11
	RL_PIXELFORMAT_COMPRESSED_DXT1_RGBA      int32 = 12
	RL_PIXELFORMAT_COMPRESSED_DXT3_RGBA      int32 = 13
	RL_PIXELFORMAT_COMPRESSED_DXT5_RGBA      int32 = 14
	RL_PIXELFORMAT_COMPRESSED_ETC1_RGB       int32 = 15
	RL_PIXELFORMAT_COMPRESSED_ETC2_RGB       int32 = 16
	RL_PIXELFORMAT_COMPRESSED_ETC2_EAC_RGBA  int32 = 17
	RL_PIXELFORMAT_COMPRESSED_PVRT_RGB       int32 = 18
	RL_PIXELFORMAT_COMPRESSED_PVRT_RGBA      int32 = 19
	RL_PIXELFORMAT_COMPRESSED_ASTC_4x4_RGBA  int32 = 20
	RL_PIXELFORMAT_COMPRESSED_ASTC_8x8_RGBA  int32 = 21
)

type rlPixelFormat = int32

// Texture parameters: filter mode
// NOTE 1: Filtering considers mipmaps if available in the texture
// NOTE 2: Filter is accordingly set for minification and magnification
const (
	RL_TEXTURE_FILTER_POINT           int32 = 0
	RL_TEXTURE_FILTER_BILINEAR        int32 = 1
	RL_TEXTURE_FILTER_TRILINEAR       int32 = 2
	RL_TEXTURE_FILTER_ANISOTROPIC_4X  int32 = 3
	RL_TEXTURE_FILTER_ANISOTROPIC_8X  int32 = 4
	RL_TEXTURE_FILTER_ANISOTROPIC_16X int32 = 5
)

type rlTextureFilter = int32

// Color blending modes (pre-defined)
const (
	RL_BLEND_ALPHA             int32 = 0
	RL_BLEND_ADDITIVE          int32 = 1
	RL_BLEND_MULTIPLIED        int32 = 2
	RL_BLEND_ADD_COLORS        int32 = 3
	RL_BLEND_SUBTRACT_COLORS   int32 = 4
	RL_BLEND_ALPHA_PREMULTIPLY int32 = 5
	RL_BLEND_CUSTOM            int32 = 6
	RL_BLEND_CUSTOM_SEPARATE   int32 = 7
)

type rlBlendMode = int32

// Shader location point type
const (
	RL_SHADER_LOC_VERTEX_POSITION   int32 = 0
	RL_SHADER_LOC_VERTEX_TEXCOORD01 int32 = 1
	RL_SHADER_LOC_VERTEX_TEXCOORD02 int32 = 2
	RL_SHADER_LOC_VERTEX_NORMAL     int32 = 3
	RL_SHADER_LOC_VERTEX_TANGENT    int32 = 4
	RL_SHADER_LOC_VERTEX_COLOR      int32 = 5
	RL_SHADER_LOC_MATRIX_MVP        int32 = 6
	RL_SHADER_LOC_MATRIX_VIEW       int32 = 7
	RL_SHADER_LOC_MATRIX_PROJECTION int32 = 8
	RL_SHADER_LOC_MATRIX_MODEL      int32 = 9
	RL_SHADER_LOC_MATRIX_NORMAL     int32 = 10
	RL_SHADER_LOC_VECTOR_VIEW       int32 = 11
	RL_SHADER_LOC_COLOR_DIFFUSE     int32 = 12
	RL_SHADER_LOC_COLOR_SPECULAR    int32 = 13
	RL_SHADER_LOC_COLOR_AMBIENT     int32 = 14
	RL_SHADER_LOC_MAP_ALBEDO        int32 = 15
	RL_SHADER_LOC_MAP_METALNESS     int32 = 16
	RL_SHADER_LOC_MAP_NORMAL        int32 = 17
	RL_SHADER_LOC_MAP_ROUGHNESS     int32 = 18
	RL_SHADER_LOC_MAP_OCCLUSION     int32 = 19
	RL_SHADER_LOC_MAP_EMISSION      int32 = 20
	RL_SHADER_LOC_MAP_HEIGHT        int32 = 21
	RL_SHADER_LOC_MAP_CUBEMAP       int32 = 22
	RL_SHADER_LOC_MAP_IRRADIANCE    int32 = 23
	RL_SHADER_LOC_MAP_PREFILTER     int32 = 24
	RL_SHADER_LOC_MAP_BRDF          int32 = 25
)

type rlShaderLocationIndex = int32

// Shader uniform data type
const (
	RL_SHADER_UNIFORM_FLOAT     int32 = 0
	RL_SHADER_UNIFORM_VEC2      int32 = 1
	RL_SHADER_UNIFORM_VEC3      int32 = 2
	RL_SHADER_UNIFORM_VEC4      int32 = 3
	RL_SHADER_UNIFORM_INT       int32 = 4
	RL_SHADER_UNIFORM_IVEC2     int32 = 5
	RL_SHADER_UNIFORM_IVEC3     int32 = 6
	RL_SHADER_UNIFORM_IVEC4     int32 = 7
	RL_SHADER_UNIFORM_SAMPLER2D int32 = 8
)

type rlShaderUniformDataType = int32

// Shader attribute data types
const (
	RL_SHADER_ATTRIB_FLOAT int32 = 0
	RL_SHADER_ATTRIB_VEC2  int32 = 1
	RL_SHADER_ATTRIB_VEC3  int32 = 2
	RL_SHADER_ATTRIB_VEC4  int32 = 3
)

type rlShaderAttributeDataType = int32

// Framebuffer attachment type
// NOTE: By default up to 8 color channels defined but it can be more
const (
	RL_ATTACHMENT_COLOR_CHANNEL0 int32 = 0
	RL_ATTACHMENT_COLOR_CHANNEL1 int32 = 1
	RL_ATTACHMENT_COLOR_CHANNEL2 int32 = 2
	RL_ATTACHMENT_COLOR_CHANNEL3 int32 = 3
	RL_ATTACHMENT_COLOR_CHANNEL4 int32 = 4
	RL_ATTACHMENT_COLOR_CHANNEL5 int32 = 5
	RL_ATTACHMENT_COLOR_CHANNEL6 int32 = 6
	RL_ATTACHMENT_COLOR_CHANNEL7 int32 = 7
	RL_ATTACHMENT_DEPTH          int32 = 100
	RL_ATTACHMENT_STENCIL        int32 = 200
)

type rlFramebufferAttachType = int32

// Framebuffer texture attachment type
const (
	RL_ATTACHMENT_CUBEMAP_POSITIVE_X int32 = 0
	RL_ATTACHMENT_CUBEMAP_NEGATIVE_X int32 = 1
	RL_ATTACHMENT_CUBEMAP_POSITIVE_Y int32 = 2
	RL_ATTACHMENT_CUBEMAP_NEGATIVE_Y int32 = 3
	RL_ATTACHMENT_CUBEMAP_POSITIVE_Z int32 = 4
	RL_ATTACHMENT_CUBEMAP_NEGATIVE_Z int32 = 5
	RL_ATTACHMENT_TEXTURE2D          int32 = 100
	RL_ATTACHMENT_RENDERBUFFER       int32 = 200
)

type rlFramebufferAttachTextureType = int32

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

// Ortho - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:518
func Ortho(left float64, right float64, bottom float64, top float64, znear float64, zfar float64) {
	cleft := C.double(left)
	cright := C.double(right)
	cbottom := C.double(bottom)
	ctop := C.double(top)
	cznear := C.double(znear)
	czfar := C.double(zfar)
	C.rlOrtho(cleft, cright, cbottom, ctop, cznear, czfar)
}

// Viewport - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:519
// Set the viewport area
func Viewport(x int32, y int32, width int32, height int32) {
	cx := C.int(x)
	cy := C.int(y)
	cwidth := C.int(width)
	cheight := C.int(height)
	C.rlViewport(cx, cy, cwidth, cheight)
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

// ActiveDrawBuffers - Activate multiple draw color buffers
func ActiveDrawBuffers(count int32) {
	ccount := C.int(count)
	C.rlActiveDrawBuffers(ccount)
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
func SetBlendMode(mode int32) {
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

// glInit - Initialize rlgl (buffers, shaders, textures, states)
func glInit(width int32, height int32) {
	cwidth := C.int(width)
	cheight := C.int(height)
	C.rlglInit(cwidth, cheight)
}

// glClose - De-inititialize rlgl (buffers, shaders, textures)
func glClose() {
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

// LoadFramebuffer - Load an empty framebuffer
func LoadFramebuffer(width int32, height int32) uint32 {
	cwidth := C.int(width)
	cheight := C.int(height)
	return uint32(C.rlLoadFramebuffer(cwidth, cheight))
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

// CompileShader - Compile custom shader and return shader id (type: RL_VERTEX_SHADER, RL_FRAGMENT_SHADER, RL_COMPUTE_SHADER)
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

// SetUniformSampler - Set shader value sampler
func SetUniformSampler(locIndex int32, textureId uint32) {
	clocIndex := C.int(locIndex)
	ctextureId := C.uint(textureId)
	C.rlSetUniformSampler(clocIndex, ctextureId)
}

// ComputeShaderDispatch - Dispatch compute shader (equivalent to *draw* for graphics pilepine)
func ComputeShaderDispatch(groupX uint32, groupY uint32, groupZ uint32) {
	cgroupX := C.uint(groupX)
	cgroupY := C.uint(groupY)
	cgroupZ := C.uint(groupZ)
	C.rlComputeShaderDispatch(cgroupX, cgroupY, cgroupZ)
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
