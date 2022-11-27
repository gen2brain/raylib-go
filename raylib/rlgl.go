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

//
//	Package - transpiled by c4go
//
//	If you have found any issues, please raise an issue at:
//	https://github.com/Konstantin8105/c4go/
//

// Warning (*ast.FunctionDecl): {prefix: n:rlMultMatrixf,t1:void (float *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:516 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlMultMatrixf`. cannot parse C type: `float *`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadExtensions,t1:void (void *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:607 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadExtensions`. cannot parse C type: `void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadRenderBatch,t1:rlRenderBatch (int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:621 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadRenderBatch`. field type is pointer: `rlVertexBuffer *`
// Warning (*ast.FunctionDecl): {prefix: n:rlUnloadRenderBatch,t1:void (rlRenderBatch),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:622 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlUnloadRenderBatch`. field type is pointer: `rlVertexBuffer *`
// Warning (*ast.FunctionDecl): {prefix: n:rlDrawRenderBatch,t1:void (rlRenderBatch *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:623 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlDrawRenderBatch`. cannot parse C type: `rlRenderBatch *`
// Warning (*ast.FunctionDecl): {prefix: n:rlSetRenderBatchActive,t1:void (rlRenderBatch *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:624 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlSetRenderBatchActive`. cannot parse C type: `rlRenderBatch *`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadVertexBuffer,t1:unsigned int (const void *, int, _Bool),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:634 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadVertexBuffer`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadVertexBufferElement,t1:unsigned int (const void *, int, _Bool),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:635 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadVertexBufferElement`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlUpdateVertexBuffer,t1:void (unsigned int, const void *, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:636 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlUpdateVertexBuffer`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlUpdateVertexBufferElements,t1:void (unsigned int, const void *, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:637 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlUpdateVertexBufferElements`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlSetVertexAttribute,t1:void (unsigned int, int, int, _Bool, int, const void *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:640 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlSetVertexAttribute`. cannot parse C type: `_Bool`
// Warning (*ast.FunctionDecl): {prefix: n:rlSetVertexAttributeDefault,t1:void (int, const void *, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:642 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlSetVertexAttributeDefault`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlDrawVertexArrayElements,t1:void (int, int, const void *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:644 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlDrawVertexArrayElements`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlDrawVertexArrayElementsInstanced,t1:void (int, int, const void *, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:646 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlDrawVertexArrayElementsInstanced`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadTexture,t1:unsigned int (const void *, int, int, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:649 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadTexture`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadTextureDepth,t1:unsigned int (int, int, _Bool),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:650 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadTextureDepth`. cannot parse C type: `_Bool`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadTextureCubemap,t1:unsigned int (const void *, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:651 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadTextureCubemap`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlUpdateTexture,t1:void (unsigned int, int, int, int, int, int, const void *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:652 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlUpdateTexture`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlGetGlTextureFormats,t1:void (int, unsigned int *, unsigned int *, unsigned int *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:653 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlGetGlTextureFormats`. cannot parse C type: `unsigned int *`
// Warning (*ast.FunctionDecl): {prefix: n:rlGenTextureMipmaps,t1:void (unsigned int, int, int, int, int *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:656 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlGenTextureMipmaps`. cannot parse C type: `int *`
// Warning (*ast.FunctionDecl): {prefix: n:rlSetUniform,t1:void (int, const void *, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:673 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlSetUniform`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlSetShader,t1:void (unsigned int, int *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:676 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlSetShader`. cannot parse C type: `int *`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadShaderBuffer,t1:unsigned int (unsigned int, const void *, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:683 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadShaderBuffer`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlUpdateShaderBuffer,t1:void (unsigned int, const void *, unsigned int, unsigned int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:685 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlUpdateShaderBuffer`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlReadShaderBuffer,t1:void (unsigned int, void *, unsigned int, unsigned int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:687 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlReadShaderBuffer`. cannot parse C type: `void *`

// Matrix - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:295
//
// *
// *   rlgl v4.2 - A multi-OpenGL abstraction layer with an immediate-mode style API
// *
// *   An abstraction layer for multiple OpenGL versions (1.1, 2.1, 3.3 Core, 4.3 Core, ES 2.0)
// *   that provides a pseudo-OpenGL 1.1 immediate-mode style API (rlVertex, rlTranslate, rlRotate...)
// *
// *   When chosing an OpenGL backend different than OpenGL 1.1, some internal buffer are
// *   initialized on rlglInit() to accumulate vertex data.
// *
// *   When an internal state change is required all the stored vertex data is renderer in batch,
// *   additioanlly, rlDrawRenderBatchActive() could be called to force flushing of the batch.
// *
// *   Some additional resources are also loaded for convenience, here the complete list:
// *      - Default batch (RLGL.defaultBatch): RenderBatch system to accumulate vertex data
// *      - Default texture (RLGL.defaultTextureId): 1x1 white pixel R8G8B8A8
// *      - Default shader (RLGL.State.defaultShaderId, RLGL.State.defaultShaderLocs)
// *
// *   Internal buffer (and additional resources) must be manually unloaded calling rlglClose().
// *
// *
// *   CONFIGURATION:
// *
// *   #define GRAPHICS_API_OPENGL_11
// *   #define GRAPHICS_API_OPENGL_21
// *   #define GRAPHICS_API_OPENGL_33
// *   #define GRAPHICS_API_OPENGL_43
// *   #define GRAPHICS_API_OPENGL_ES2
// *       Use selected OpenGL graphics backend, should be supported by platform
// *       Those preprocessor defines are only used on rlgl module, if OpenGL version is
// *       required by any other module, use rlGetVersion() to check it
// *
// *   #define RLGL_IMPLEMENTATION
// *       Generates the implementation of the library into the included file.
// *       If not defined, the library is in header only mode and can be included in other headers
// *       or source files without problems. But only ONE file should hold the implementation.
// *
// *   #define RLGL_RENDER_TEXTURES_HINT
// *       Enable framebuffer objects (fbo) support (enabled by default)
// *       Some GPUs could not support them despite the OpenGL version
// *
// *   #define RLGL_SHOW_GL_DETAILS_INFO
// *       Show OpenGL extensions and capabilities detailed logs on init
// *
// *   #define RLGL_ENABLE_OPENGL_DEBUG_CONTEXT
// *       Enable debug context (only available on OpenGL 4.3)
// *
// *   rlgl capabilities could be customized just defining some internal
// *   values before library inclusion (default values listed):
// *
// *   #define RL_DEFAULT_BATCH_BUFFER_ELEMENTS   8192    // Default internal render batch elements limits
// *   #define RL_DEFAULT_BATCH_BUFFERS              1    // Default number of batch buffers (multi-buffering)
// *   #define RL_DEFAULT_BATCH_DRAWCALLS          256    // Default number of batch draw calls (by state changes: mode, texture)
// *   #define RL_DEFAULT_BATCH_MAX_TEXTURE_UNITS    4    // Maximum number of textures units that can be activated on batch drawing (SetShaderValueTexture())
// *
// *   #define RL_MAX_MATRIX_STACK_SIZE             32    // Maximum size of internal Matrix stack
// *   #define RL_MAX_SHADER_LOCATIONS              32    // Maximum number of shader locations supported
// *   #define RL_CULL_DISTANCE_NEAR              0.01    // Default projection matrix near cull distance
// *   #define RL_CULL_DISTANCE_FAR             1000.0    // Default projection matrix far cull distance
// *
// *   When loading a shader, the following vertex attribute and uniform
// *   location names are tried to be set automatically:
// *
// *   #define RL_DEFAULT_SHADER_ATTRIB_NAME_POSITION     "vertexPosition"    // Binded by default to shader location: 0
// *   #define RL_DEFAULT_SHADER_ATTRIB_NAME_TEXCOORD     "vertexTexCoord"    // Binded by default to shader location: 1
// *   #define RL_DEFAULT_SHADER_ATTRIB_NAME_NORMAL       "vertexNormal"      // Binded by default to shader location: 2
// *   #define RL_DEFAULT_SHADER_ATTRIB_NAME_COLOR        "vertexColor"       // Binded by default to shader location: 3
// *   #define RL_DEFAULT_SHADER_ATTRIB_NAME_TANGENT      "vertexTangent"     // Binded by default to shader location: 4
// *   #define RL_DEFAULT_SHADER_ATTRIB_NAME_TEXCOORD2    "vertexTexCoord2"   // Binded by default to shader location: 5
// *   #define RL_DEFAULT_SHADER_UNIFORM_NAME_MVP         "mvp"               // model-view-projection matrix
// *   #define RL_DEFAULT_SHADER_UNIFORM_NAME_VIEW        "matView"           // view matrix
// *   #define RL_DEFAULT_SHADER_UNIFORM_NAME_PROJECTION  "matProjection"     // projection matrix
// *   #define RL_DEFAULT_SHADER_UNIFORM_NAME_MODEL       "matModel"          // model matrix
// *   #define RL_DEFAULT_SHADER_UNIFORM_NAME_NORMAL      "matNormal"         // normal matrix (transpose(inverse(matModelView))
// *   #define RL_DEFAULT_SHADER_UNIFORM_NAME_COLOR       "colDiffuse"        // color diffuse (base tint color, multiplied by texture color)
// *   #define RL_DEFAULT_SHADER_SAMPLER2D_NAME_TEXTURE0  "texture0"          // texture0 (texture slot active 0)
// *   #define RL_DEFAULT_SHADER_SAMPLER2D_NAME_TEXTURE1  "texture1"          // texture1 (texture slot active 1)
// *   #define RL_DEFAULT_SHADER_SAMPLER2D_NAME_TEXTURE2  "texture2"          // texture2 (texture slot active 2)
// *
// *   DEPENDENCIES:
// *
// *      - OpenGL libraries (depending on platform and OpenGL version selected)
// *      - GLAD OpenGL extensions loading library (only for OpenGL 3.3 Core, 4.3 Core)
// *
// *
// *   LICENSE: zlib/libpng
// *
// *   Copyright (c) 2014-2022 Ramon Santamaria (@raysan5)
// *
// *   This software is provided "as-is", without any express or implied warranty. In no event
// *   will the authors be held liable for any damages arising from the use of this software.
// *
// *   Permission is granted to anyone to use this software for any purpose, including commercial
// *   applications, and to alter it and redistribute it freely, subject to the following restrictions:
// *
// *     1. The origin of this software must not be misrepresented; you must not claim that you
// *     wrote the original software. If you use this software in a product, an acknowledgment
// *     in the product documentation would be appreciated but is not required.
// *
// *     2. Altered source versions must be plainly marked as such, and must not be misrepresented
// *     as being the original software.
// *
// *     3. This notice may not be removed or altered from any source distribution.
// *
//
// Function specifiers in case library is build/used as a shared library (Windows)
// NOTE: Microsoft specifiers to tell compiler that symbols are imported/exported from a .dll
// Function specifiers definition
// Support TRACELOG macros
// Allow custom memory allocators
// Security check in case no GRAPHICS_API_OPENGL_* defined
// Security check in case multiple GRAPHICS_API_OPENGL_* defined
// OpenGL 2.1 uses most of OpenGL 3.3 Core functionality
// WARNING: Specific parts are checked with #if defines
// OpenGL 4.3 uses OpenGL 3.3 Core functionality
// Support framebuffer objects by default
// NOTE: Some driver implementation do not support it, despite they should
// ----------------------------------------------------------------------------------
// Defines and Macros
// ----------------------------------------------------------------------------------
// Default internal render batch elements limits
// This is the maximum amount of elements (quads) per batch
// NOTE: Be careful with text, every letter maps to a quad
// Internal Matrix stack
// Shader limits
// Projection matrix culling
// Texture parameters (equivalent to OpenGL defines)
// Matrix modes (equivalent to OpenGL)
// Primitive assembly draw modes
// GL equivalent data types
// Buffer usage hint
// GL Shader Types

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

// ----------------------------------------------------------------------------------
// Types and Structures Definition
// ----------------------------------------------------------------------------------
// Matrix, 4x4 components, column major, OpenGL style, right handed
// type Matrix struct {
// 	m0  float32
// 	m4  float32
// 	m8  float32
// 	m12 float32
// 	m1  float32
// 	m5  float32
// 	m9  float32
// 	m13 float32
// 	m2  float32
// 	m6  float32
// 	m10 float32
// 	m14 float32
// 	m3  float32
// 	m7  float32
// 	m11 float32
// 	m15 float32
// }

// rlVertexBuffer - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:305
// Matrix first row (4 components)
// Matrix second row (4 components)
// Matrix third row (4 components)
// Matrix fourth row (4 components)
// Dynamic vertex buffers (position + texcoords + colors + indices arrays)
type rlVertexBuffer struct {
	elementCount int32
	vertices     []float32
	texcoords    []float32
	colors       []uint8
	indices      []uint32
	vaoId        uint32
	vboId        [4]uint32
}

// rlDrawCall - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:325
// Number of elements in the buffer (QUADS)
// Vertex position (XYZ - 3 components per vertex) (shader-location = 0)
// Vertex texture coordinates (UV - 2 components per vertex) (shader-location = 1)
// Vertex colors (RGBA - 4 components per vertex) (shader-location = 3)
// Vertex indices (in case vertex data comes indexed) (6 indices per quad)
// OpenGL Vertex Array Object id
// OpenGL Vertex Buffer Objects id (4 types of vertex data)
// Draw call type
// NOTE: Only texture changes register a new draw, other state-change-related elements are not
// used at this moment (vaoId, shaderId, matrices), raylib just forces a batch draw call if any
// of those state-change happens (this is done in core module)
type rlDrawCall struct {
	mode            int32
	vertexCount     int32
	vertexAlignment int32
	textureId       uint32
}

// rlRenderBatch - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:338
// Drawing mode: LINES, TRIANGLES, QUADS
// Number of vertex of the draw
// Number of vertex required for index alignment (LINES, TRIANGLES)
// unsigned int vaoId;       // Vertex array id to be used on the draw -> Using RLGL.currentBatch->vertexBuffer.vaoId
// unsigned int shaderId;    // Shader id to be used on the draw -> Using RLGL.currentShaderId
// Texture id to be used on the draw -> Use to create new draw call if changes
// Matrix projection;        // Projection matrix for this draw -> Using RLGL.projection by default
// Matrix modelview;         // Modelview matrix for this draw -> Using RLGL.modelview by default
// rlRenderBatch type
type rlRenderBatch struct {
	bufferCount   int32
	currentBuffer int32
	vertexBuffer  []rlVertexBuffer
	draws         []rlDrawCall
	drawCounter   int32
	currentDepth  float32
}

// RL_OPENGL_11 - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:349
// Number of vertex buffers (multi-buffering support)
// Current buffer tracking in case of multi-buffering
// Dynamic buffer(s) for vertex data
// Draw calls array, depends on textureId
// Draw calls counter
// Current depth value for next draw
// OpenGL version
const (
	RL_OPENGL_11    int32 = 1
	RL_OPENGL_21    int32 = 2
	RL_OPENGL_33    int32 = 3
	RL_OPENGL_43    int32 = 4
	RL_OPENGL_ES_20 int32 = 5
)

// rlGlVersion - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:349
type rlGlVersion = int32

// RL_LOG_ALL - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:359
// OpenGL 1.1
// OpenGL 2.1 (GLSL 120)
// OpenGL 3.3 (GLSL 330)
// OpenGL 4.3 (using GLSL 330)
// OpenGL ES 2.0 (GLSL 100)
// Trace log level
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

// rlTraceLogLevel - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:359
type rlTraceLogLevel = int32

// RL_PIXELFORMAT_UNCOMPRESSED_GRAYSCALE - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:372
// Display all logs
// Trace logging, intended for internal use only
// Debug logging, used for internal debugging, it should be disabled on release builds
// Info logging, used for program execution info
// Warning logging, used on recoverable failures
// Error logging, used on unrecoverable failures
// Fatal logging, used to abort program: exit(EXIT_FAILURE)
// Disable logging
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

// rlPixelFormat - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:372
type rlPixelFormat = int32

// RL_TEXTURE_FILTER_POINT - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:399
// 8 bit per pixel (no alpha)
// 8*2 bpp (2 channels)
// 16 bpp
// 24 bpp
// 16 bpp (1 bit alpha)
// 16 bpp (4 bit alpha)
// 32 bpp
// 32 bpp (1 channel - float)
// 32*3 bpp (3 channels - float)
// 32*4 bpp (4 channels - float)
// 4 bpp (no alpha)
// 4 bpp (1 bit alpha)
// 8 bpp
// 8 bpp
// 4 bpp
// 4 bpp
// 8 bpp
// 4 bpp
// 4 bpp
// 8 bpp
// 2 bpp
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

// rlTextureFilter - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:399
type rlTextureFilter = int32

// RL_BLEND_ALPHA - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:409
// No filter, just pixel approximation
// Linear filtering
// Trilinear filtering (linear with mipmaps)
// Anisotropic filtering 4x
// Anisotropic filtering 8x
// Anisotropic filtering 16x
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

// rlBlendMode - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:409
type rlBlendMode = int32

// RL_SHADER_LOC_VERTEX_POSITION - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:421
// Blend textures considering alpha (default)
// Blend textures adding colors
// Blend textures multiplying colors
// Blend textures adding colors (alternative)
// Blend textures subtracting colors (alternative)
// Blend premultiplied textures considering alpha
// Blend textures using custom src/dst factors (use rlSetBlendFactors())
// Blend textures using custom src/dst factors (use rlSetBlendFactorsSeparate())
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

// rlShaderLocationIndex - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:421
type rlShaderLocationIndex = int32

// RL_SHADER_UNIFORM_FLOAT - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:454
// Shader location: vertex attribute: position
// Shader location: vertex attribute: texcoord01
// Shader location: vertex attribute: texcoord02
// Shader location: vertex attribute: normal
// Shader location: vertex attribute: tangent
// Shader location: vertex attribute: color
// Shader location: matrix uniform: model-view-projection
// Shader location: matrix uniform: view (camera transform)
// Shader location: matrix uniform: projection
// Shader location: matrix uniform: model (transform)
// Shader location: matrix uniform: normal
// Shader location: vector uniform: view
// Shader location: vector uniform: diffuse color
// Shader location: vector uniform: specular color
// Shader location: vector uniform: ambient color
// Shader location: sampler2d texture: albedo (same as: RL_SHADER_LOC_MAP_DIFFUSE)
// Shader location: sampler2d texture: metalness (same as: RL_SHADER_LOC_MAP_SPECULAR)
// Shader location: sampler2d texture: normal
// Shader location: sampler2d texture: roughness
// Shader location: sampler2d texture: occlusion
// Shader location: sampler2d texture: emission
// Shader location: sampler2d texture: height
// Shader location: samplerCube texture: cubemap
// Shader location: samplerCube texture: irradiance
// Shader location: samplerCube texture: prefilter
// Shader location: sampler2d texture: brdf
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

// rlShaderUniformDataType - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:454
type rlShaderUniformDataType = int32

// RL_SHADER_ATTRIB_FLOAT - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:467
// Shader uniform type: float
// Shader uniform type: vec2 (2 float)
// Shader uniform type: vec3 (3 float)
// Shader uniform type: vec4 (4 float)
// Shader uniform type: int
// Shader uniform type: ivec2 (2 int)
// Shader uniform type: ivec3 (3 int)
// Shader uniform type: ivec4 (4 int)
// Shader uniform type: sampler2d
// Shader attribute data types
const (
	RL_SHADER_ATTRIB_FLOAT int32 = 0
	RL_SHADER_ATTRIB_VEC2  int32 = 1
	RL_SHADER_ATTRIB_VEC3  int32 = 2
	RL_SHADER_ATTRIB_VEC4  int32 = 3
)

// rlShaderAttributeDataType - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:467
type rlShaderAttributeDataType = int32

// RL_ATTACHMENT_COLOR_CHANNEL0 - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:476
// Shader attribute type: float
// Shader attribute type: vec2 (2 float)
// Shader attribute type: vec3 (3 float)
// Shader attribute type: vec4 (4 float)
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

// rlFramebufferAttachType - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:476
type rlFramebufferAttachType = int32

// RL_ATTACHMENT_CUBEMAP_POSITIVE_X - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:490
// Framebuffer attachmment type: color 0
// Framebuffer attachmment type: color 1
// Framebuffer attachmment type: color 2
// Framebuffer attachmment type: color 3
// Framebuffer attachmment type: color 4
// Framebuffer attachmment type: color 5
// Framebuffer attachmment type: color 6
// Framebuffer attachmment type: color 7
// Framebuffer attachmment type: depth
// Framebuffer attachmment type: stencil
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

// rlFramebufferAttachTextureType - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:490
type rlFramebufferAttachTextureType = int32

// MatrixMode - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:509
// Framebuffer texture attachment type: cubemap, +X side
// Framebuffer texture attachment type: cubemap, -X side
// Framebuffer texture attachment type: cubemap, +Y side
// Framebuffer texture attachment type: cubemap, -Y side
// Framebuffer texture attachment type: cubemap, +Z side
// Framebuffer texture attachment type: cubemap, -Z side
// Framebuffer texture attachment type: texture2d
// Framebuffer texture attachment type: renderbuffer
// ------------------------------------------------------------------------------------
// Functions Declaration - Matrix operations
// ------------------------------------------------------------------------------------
// Choose the current matrix to be transformed
func MatrixMode(mode int32) {
	cmode := C.int(mode)
	C.rlMatrixMode(cmode)
}

// PushMatrix - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:510
// Push the current matrix to stack
func PushMatrix() {
	C.rlPushMatrix()
}

// PopMatrix - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:511
// Pop lattest inserted matrix from stack
func PopMatrix() {
	C.rlPopMatrix()
}

// LoadIdentity - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:512
// Reset current matrix to identity matrix
func LoadIdentity() {
	C.rlLoadIdentity()
}

// Translatef - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:513
// Multiply the current matrix by a translation matrix
func Translatef(x float32, y float32, z float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlTranslatef(cx, cy, cz)
}

// Rotatef - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:514
// Multiply the current matrix by a rotation matrix
func Rotatef(angle float32, x float32, y float32, z float32) {
	cangle := C.float(angle)
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlRotatef(cangle, cx, cy, cz)
}

// Scalef - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:515
// Multiply the current matrix by a scaling matrix
func Scalef(x float32, y float32, z float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlScalef(cx, cy, cz)
}

// Frustum - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:517
// Multiply the current matrix by another matrix
// Warning (*ast.FunctionDecl): {prefix: n:rlMultMatrixf,t1:void (float *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:516 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlMultMatrixf`. cannot parse C type: `float *`
// func Frustum(left float64, right float64, bottom float64, top float64, znear float64, zfar float64) {
// 	cleft := C.double(left)
// 	cright := C.double(right)
// 	cbottom := C.double(bottom)
// 	ctop := C.double(top)
// 	cznear := C.double(znear)
// 	czfar := C.double(zfar)
// 	C.rlFrustum(cleft, cright, cbottom, ctop, cznear, czfar)
// }

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

// Begin - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:524
// ------------------------------------------------------------------------------------
// Functions Declaration - Vertex level operations
// ------------------------------------------------------------------------------------
// Initialize drawing mode (how to organize vertex)
func Begin(mode int32) {
	cmode := C.int(mode)
	C.rlBegin(cmode)
}

// End - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:525
// Finish vertex providing
func End() {
	C.rlEnd()
}

// Vertex2i - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:526
// Define one vertex (position) - 2 int
func Vertex2i(x int32, y int32) {
	cx := C.int(x)
	cy := C.int(y)
	C.rlVertex2i(cx, cy)
}

// Vertex2f - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:527
// Define one vertex (position) - 2 float
func Vertex2f(x float32, y float32) {
	cx := C.float(x)
	cy := C.float(y)
	C.rlVertex2f(cx, cy)
}

// Vertex3f - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:528
// Define one vertex (position) - 3 float
func Vertex3f(x float32, y float32, z float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlVertex3f(cx, cy, cz)
}

// TexCoord2f - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:529
// Define one vertex (texture coordinate) - 2 float
func TexCoord2f(x float32, y float32) {
	cx := C.float(x)
	cy := C.float(y)
	C.rlTexCoord2f(cx, cy)
}

// Normal3f - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:530
// Define one vertex (normal) - 3 float
func Normal3f(x float32, y float32, z float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlNormal3f(cx, cy, cz)
}

// Color4ub - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:531
// Define one vertex (color) - 4 byte
func Color4ub(r uint8, g uint8, b uint8, a uint8) {
	cr := C.uchar(r)
	cg := C.uchar(g)
	cb := C.uchar(b)
	ca := C.uchar(a)
	C.rlColor4ub(cr, cg, cb, ca)
}

// Color3f - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:532
// Define one vertex (color) - 3 float
func Color3f(x float32, y float32, z float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	C.rlColor3f(cx, cy, cz)
}

// Color4f - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:533
// Define one vertex (color) - 4 float
func Color4f(x float32, y float32, z float32, w float32) {
	cx := C.float(x)
	cy := C.float(y)
	cz := C.float(z)
	cw := C.float(w)
	C.rlColor4f(cx, cy, cz, cw)
}

// EnableVertexArray - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:542
// ------------------------------------------------------------------------------------
// Functions Declaration - OpenGL style functions (common to 1.1, 3.3+, ES2)
// NOTE: This functions are used to completely abstract raylib code from OpenGL layer,
// some of them are direct wrappers over OpenGL calls, some others are custom
// ------------------------------------------------------------------------------------
// Vertex buffers state
// Enable vertex array (VAO, if supported)
func EnableVertexArray(vaoId uint32) bool {
	cvaoId := C.uint(vaoId)
	return bool(C.rlEnableVertexArray(cvaoId))
}

// DisableVertexArray - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:543
// Disable vertex array (VAO, if supported)
func DisableVertexArray() {
	C.rlDisableVertexArray()
}

// EnableVertexBuffer - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:544
// Enable vertex buffer (VBO)
func EnableVertexBuffer(id uint32) {
	cid := C.uint(id)
	C.rlEnableVertexBuffer(cid)
}

// DisableVertexBuffer - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:545
// Disable vertex buffer (VBO)
func DisableVertexBuffer() {
	C.rlDisableVertexBuffer()
}

// EnableVertexBufferElement - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:546
// Enable vertex buffer element (VBO element)
func EnableVertexBufferElement(id uint32) {
	cid := C.uint(id)
	C.rlEnableVertexBufferElement(cid)
}

// DisableVertexBufferElement - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:547
// Disable vertex buffer element (VBO element)
func DisableVertexBufferElement() {
	C.rlDisableVertexBufferElement()
}

// EnableVertexAttribute - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:548
// Enable vertex attribute index
func EnableVertexAttribute(index uint32) {
	cindex := C.uint(index)
	C.rlEnableVertexAttribute(cindex)
}

// DisableVertexAttribute - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:549
// Disable vertex attribute index
func DisableVertexAttribute(index uint32) {
	cindex := C.uint(index)
	C.rlDisableVertexAttribute(cindex)
}

// ActiveTextureSlot - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:556
// Textures state
// Select and active a texture slot
func ActiveTextureSlot(slot int32) {
	cslot := C.int(slot)
	C.rlActiveTextureSlot(cslot)
}

// EnableTexture - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:557
// Enable texture
func EnableTexture(id uint32) {
	cid := C.uint(id)
	C.rlEnableTexture(cid)
}

// DisableTexture - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:558
// Disable texture
func DisableTexture() {
	C.rlDisableTexture()
}

// EnableTextureCubemap - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:559
// Enable texture cubemap
func EnableTextureCubemap(id uint32) {
	cid := C.uint(id)
	C.rlEnableTextureCubemap(cid)
}

// DisableTextureCubemap - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:560
// Disable texture cubemap
func DisableTextureCubemap() {
	C.rlDisableTextureCubemap()
}

// TextureParameters - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:561
// Set texture parameters (filter, wrap)
func TextureParameters(id uint32, param int32, value int32) {
	cid := C.uint(id)
	cparam := C.int(param)
	cvalue := C.int(value)
	C.rlTextureParameters(cid, cparam, cvalue)
}

// EnableShader - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:564
// Shader state
// Enable shader program
func EnableShader(id uint32) {
	cid := C.uint(id)
	C.rlEnableShader(cid)
}

// DisableShader - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:565
// Disable shader program
func DisableShader() {
	C.rlDisableShader()
}

// EnableFramebuffer - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:568
// Framebuffer state
// Enable render texture (fbo)
func EnableFramebuffer(id uint32) {
	cid := C.uint(id)
	C.rlEnableFramebuffer(cid)
}

// DisableFramebuffer - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:569
// Disable render texture (fbo), return to default framebuffer
func DisableFramebuffer() {
	C.rlDisableFramebuffer()
}

// ActiveDrawBuffers - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:570
// Activate multiple draw color buffers
func ActiveDrawBuffers(count int32) {
	ccount := C.int(count)
	C.rlActiveDrawBuffers(ccount)
}

// EnableColorBlend - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:573
// General render state
// Enable color blending
func EnableColorBlend() {
	C.rlEnableColorBlend()
}

// DisableColorBlend - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:574
// Disable color blending
func DisableColorBlend() {
	C.rlDisableColorBlend()
}

// EnableDepthTest - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:575
// Enable depth test
func EnableDepthTest() {
	C.rlEnableDepthTest()
}

// DisableDepthTest - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:576
// Disable depth test
func DisableDepthTest() {
	C.rlDisableDepthTest()
}

// EnableDepthMask - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:577
// Enable depth write
func EnableDepthMask() {
	C.rlEnableDepthMask()
}

// DisableDepthMask - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:578
// Disable depth write
func DisableDepthMask() {
	C.rlDisableDepthMask()
}

// EnableBackfaceCulling - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:579
// Enable backface culling
func EnableBackfaceCulling() {
	C.rlEnableBackfaceCulling()
}

// DisableBackfaceCulling - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:580
// Disable backface culling
func DisableBackfaceCulling() {
	C.rlDisableBackfaceCulling()
}

// EnableScissorTest - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:581
// Enable scissor test
func EnableScissorTest() {
	C.rlEnableScissorTest()
}

// DisableScissorTest - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:582
// Disable scissor test
func DisableScissorTest() {
	C.rlDisableScissorTest()
}

// Scissor - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:583
// Scissor test
func Scissor(x int32, y int32, width int32, height int32) {
	cx := C.int(x)
	cy := C.int(y)
	cwidth := C.int(width)
	cheight := C.int(height)
	C.rlScissor(cx, cy, cwidth, cheight)
}

// EnableWireMode - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:584
// Enable wire mode
func EnableWireMode() {
	C.rlEnableWireMode()
}

// DisableWireMode - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:585
// Disable wire mode
func DisableWireMode() {
	C.rlDisableWireMode()
}

// SetLineWidth - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:586
// Set the line drawing width
func SetLineWidth(width float32) {
	cwidth := C.float(width)
	C.rlSetLineWidth(cwidth)
}

// GetLineWidth - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:587
// Get the line drawing width
func GetLineWidth() float32 {
	return float32(C.rlGetLineWidth())
}

// EnableSmoothLines - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:588
// Enable line aliasing
func EnableSmoothLines() {
	C.rlEnableSmoothLines()
}

// DisableSmoothLines - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:589
// Disable line aliasing
func DisableSmoothLines() {
	C.rlDisableSmoothLines()
}

// EnableStereoRender - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:590
// Enable stereo rendering
func EnableStereoRender() {
	C.rlEnableStereoRender()
}

// DisableStereoRender - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:591
// Disable stereo rendering
func DisableStereoRender() {
	C.rlDisableStereoRender()
}

// IsStereoRenderEnabled - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:592
// Check if stereo render is enabled
func IsStereoRenderEnabled() bool {
	return bool(C.rlIsStereoRenderEnabled())
}

// ClearColor - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:594
// Clear color buffer with color
func ClearColor(r uint8, g uint8, b uint8, a uint8) {
	cr := C.uchar(r)
	cg := C.uchar(g)
	cb := C.uchar(b)
	ca := C.uchar(a)
	C.rlClearColor(cr, cg, cb, ca)
}

// ClearScreenBuffers - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:595
// Clear used screen buffers (color and depth)
func ClearScreenBuffers() {
	C.rlClearScreenBuffers()
}

// CheckErrors - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:596
// Check and log OpenGL error codes
func CheckErrors() {
	C.rlCheckErrors()
}

// SetBlendMode - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:597
// Set blending mode
func SetBlendMode(mode int32) {
	cmode := C.int(mode)
	C.rlSetBlendMode(cmode)
}

// SetBlendFactors - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:598
// Set blending mode factor and equation (using OpenGL factors)
func SetBlendFactors(glSrcFactor int32, glDstFactor int32, glEquation int32) {
	cglSrcFactor := C.int(glSrcFactor)
	cglDstFactor := C.int(glDstFactor)
	cglEquation := C.int(glEquation)
	C.rlSetBlendFactors(cglSrcFactor, cglDstFactor, cglEquation)
}

// SetBlendFactorsSeparate - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:599
// Set blending mode factors and equations separately (using OpenGL factors)
func SetBlendFactorsSeparate(glSrcRGB int32, glDstRGB int32, glSrcAlpha int32, glDstAlpha int32, glEqRGB int32, glEqAlpha int32) {
	cglSrcRGB := C.int(glSrcRGB)
	cglDstRGB := C.int(glDstRGB)
	cglSrcAlpha := C.int(glSrcAlpha)
	cglDstAlpha := C.int(glDstAlpha)
	cglEqRGB := C.int(glEqRGB)
	cglEqAlpha := C.int(glEqAlpha)
	C.rlSetBlendFactorsSeparate(cglSrcRGB, cglDstRGB, cglSrcAlpha, cglDstAlpha, cglEqRGB, cglEqAlpha)
}

// glInit - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:605
// ------------------------------------------------------------------------------------
// Functions Declaration - rlgl functionality
// ------------------------------------------------------------------------------------
// rlgl initialization functions
// Initialize rlgl (buffers, shaders, textures, states)
func glInit(width int32, height int32) {
	cwidth := C.int(width)
	cheight := C.int(height)
	C.rlglInit(cwidth, cheight)
}

// glClose - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:606
// De-inititialize rlgl (buffers, shaders, textures)
func glClose() {
	C.rlglClose()
}

// GetVersion - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:608
// Load OpenGL extensions (loader function required)
// Get current OpenGL version
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadExtensions,t1:void (void *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:607 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadExtensions`. cannot parse C type: `void *`
func GetVersion() int32 {
	return int32(C.rlGetVersion())
}

// SetFramebufferWidth - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:609
// Set current framebuffer width
func SetFramebufferWidth(width int32) {
	cwidth := C.int(width)
	C.rlSetFramebufferWidth(cwidth)
}

// GetFramebufferWidth - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:610
// Get default framebuffer width
func GetFramebufferWidth() int32 {
	return int32(C.rlGetFramebufferWidth())
}

// SetFramebufferHeight - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:611
// Set current framebuffer height
func SetFramebufferHeight(height int32) {
	cheight := C.int(height)
	C.rlSetFramebufferHeight(cheight)
}

// GetFramebufferHeight - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:612
// Get default framebuffer height
func GetFramebufferHeight() int32 {
	return int32(C.rlGetFramebufferHeight())
}

// GetTextureIdDefault - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:614
// Get default texture id
func GetTextureIdDefault() uint32 {
	return uint32(C.rlGetTextureIdDefault())
}

// GetShaderIdDefault - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:615
// Get default shader id
func GetShaderIdDefault() uint32 {
	return uint32(C.rlGetShaderIdDefault())
}

// GetShaderLocsDefault - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:616
// Get default shader locations
// func GetShaderLocsDefault() []int32 {
// 	return C.rlGetShaderLocsDefault()
// }

// DrawRenderBatchActive - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:625
// Render batch management
// NOTE: rlgl provides a default render batch to behave like OpenGL 1.1 immediate mode
// but this render batch API is exposed in case of custom batches are required
// Load a render batch system
// Unload render batch system
// Draw render batch data (Update->Draw->Reset)
// Set the active render batch for rlgl (NULL for default internal)
// Update and draw internal render batch
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadRenderBatch,t1:rlRenderBatch (int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:621 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadRenderBatch`. field type is pointer: `rlVertexBuffer *`
// Warning (*ast.FunctionDecl): {prefix: n:rlUnloadRenderBatch,t1:void (rlRenderBatch),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:622 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlUnloadRenderBatch`. field type is pointer: `rlVertexBuffer *`
// Warning (*ast.FunctionDecl): {prefix: n:rlDrawRenderBatch,t1:void (rlRenderBatch *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:623 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlDrawRenderBatch`. cannot parse C type: `rlRenderBatch *`
// Warning (*ast.FunctionDecl): {prefix: n:rlSetRenderBatchActive,t1:void (rlRenderBatch *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:624 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlSetRenderBatchActive`. cannot parse C type: `rlRenderBatch *`
// func DrawRenderBatchActive() {
// 	C.rlDrawRenderBatchActive()
// }

// CheckRenderBatchLimit - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:626
// Check internal buffer overflow for a given number of vertex
func CheckRenderBatchLimit(vCount int32) bool {
	cvCount := C.int(vCount)
	return bool(C.rlCheckRenderBatchLimit(cvCount))
}

// SetTexture - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:628
// Set current texture for render batch and check buffers limits
func SetTexture(id uint32) {
	cid := C.uint(id)
	C.rlSetTexture(cid)
}

// LoadVertexArray - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:633
// ------------------------------------------------------------------------------------------------------------------------
// Vertex buffers management
// Load vertex array (vao) if supported
func LoadVertexArray() uint32 {
	return uint32(C.rlLoadVertexArray())
}

// UnloadVertexArray - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:638
// Load a vertex buffer attribute
// Load a new attributes element buffer
// Update GPU buffer with new data
// Update vertex buffer elements with new data
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadVertexBuffer,t1:unsigned int (const void *, int, _Bool),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:634 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadVertexBuffer`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadVertexBufferElement,t1:unsigned int (const void *, int, _Bool),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:635 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadVertexBufferElement`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlUpdateVertexBuffer,t1:void (unsigned int, const void *, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:636 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlUpdateVertexBuffer`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlUpdateVertexBufferElements,t1:void (unsigned int, const void *, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:637 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlUpdateVertexBufferElements`. cannot parse C type: `const void *`
// func UnloadVertexArray(vaoId uint32) {
// 	cvaoId := C.uint(vaoId)
// 	C.rlUnloadVertexArray(cvaoId)
// }

// UnloadVertexBuffer - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:639
func UnloadVertexBuffer(vboId uint32) {
	cvboId := C.uint(vboId)
	C.rlUnloadVertexBuffer(cvboId)
}

// SetVertexAttributeDivisor - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:641
// Warning (*ast.FunctionDecl): {prefix: n:rlSetVertexAttribute,t1:void (unsigned int, int, int, _Bool, int, const void *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:640 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlSetVertexAttribute`. cannot parse C type: `_Bool`
func SetVertexAttributeDivisor(index uint32, divisor int32) {
	cindex := C.uint(index)
	cdivisor := C.int(divisor)
	C.rlSetVertexAttributeDivisor(cindex, cdivisor)
}

// DrawVertexArray - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:643
// Set vertex attribute default value
// Warning (*ast.FunctionDecl): {prefix: n:rlSetVertexAttributeDefault,t1:void (int, const void *, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:642 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlSetVertexAttributeDefault`. cannot parse C type: `const void *`
// func DrawVertexArray(offset int32, count int32) {
// 	coffset := C.int(offset)
// 	ccount := C.int(count)
// 	C.rlDrawVertexArray(coffset, ccount)
// }

// DrawVertexArrayInstanced - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:645
// Warning (*ast.FunctionDecl): {prefix: n:rlDrawVertexArrayElements,t1:void (int, int, const void *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:644 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlDrawVertexArrayElements`. cannot parse C type: `const void *`
// func DrawVertexArrayInstanced(offset int32, count int32, instances int32) {
// 	coffset := C.int(offset)
// 	ccount := C.int(count)
// 	cinstances := C.int(instances)
// 	C.rlDrawVertexArrayInstanced(coffset, ccount, cinstances)
// }

// GetPixelFormatName - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:654
// Textures management
// Load texture in GPU
// Load depth texture/renderbuffer (to be attached to fbo)
// Load texture cubemap
// Update GPU texture with new data
// Get OpenGL internal formats
// Get name string for pixel format
// Warning (*ast.FunctionDecl): {prefix: n:rlDrawVertexArrayElementsInstanced,t1:void (int, int, const void *, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:646 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlDrawVertexArrayElementsInstanced`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadTexture,t1:unsigned int (const void *, int, int, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:649 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadTexture`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadTextureDepth,t1:unsigned int (int, int, _Bool),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:650 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadTextureDepth`. cannot parse C type: `_Bool`
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadTextureCubemap,t1:unsigned int (const void *, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:651 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadTextureCubemap`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlUpdateTexture,t1:void (unsigned int, int, int, int, int, int, const void *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:652 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlUpdateTexture`. cannot parse C type: `const void *`
// Warning (*ast.FunctionDecl): {prefix: n:rlGetGlTextureFormats,t1:void (int, unsigned int *, unsigned int *, unsigned int *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:653 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlGetGlTextureFormats`. cannot parse C type: `unsigned int *`
// func GetPixelFormatName(format uint32) []byte {
// 	cformat := C.uint(format)
// 	return C.rlGetPixelFormatName(cformat)
// }

// UnloadTexture - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:655
// Unload texture from GPU memory
// func UnloadTexture(id uint32) {
// 	cid := C.uint(id)
// 	C.rlUnloadTexture(cid)
// }

// ReadTexturePixels - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:657
// Generate mipmap data for selected texture
// Read texture pixel data
// Warning (*ast.FunctionDecl): {prefix: n:rlGenTextureMipmaps,t1:void (unsigned int, int, int, int, int *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:656 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlGenTextureMipmaps`. cannot parse C type: `int *`
// func ReadTexturePixels(id uint32, width int32, height int32, format int32) interface{} {
// 	cid := C.uint(id)
// 	cwidth := C.int(width)
// 	cheight := C.int(height)
// 	cformat := C.int(format)
// 	return C.rlReadTexturePixels(cid, cwidth, cheight, cformat)
// }

// ReadScreenPixels - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:658
// Read screen pixel data (color buffer)
// func ReadScreenPixels(width int32, height int32) []uint8 {
// 	cwidth := C.int(width)
// 	cheight := C.int(height)
// 	return C.rlReadScreenPixels(cwidth, cheight)
// }

// LoadFramebuffer - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:661
// Framebuffer management (fbo)
// Load an empty framebuffer
func LoadFramebuffer(width int32, height int32) uint32 {
	cwidth := C.int(width)
	cheight := C.int(height)
	return uint32(C.rlLoadFramebuffer(cwidth, cheight))
}

// FramebufferAttach - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:662
// Attach texture/renderbuffer to a framebuffer
func FramebufferAttach(fboId uint32, texId uint32, attachType int32, texType int32, mipLevel int32) {
	cfboId := C.uint(fboId)
	ctexId := C.uint(texId)
	cattachType := C.int(attachType)
	ctexType := C.int(texType)
	cmipLevel := C.int(mipLevel)
	C.rlFramebufferAttach(cfboId, ctexId, cattachType, ctexType, cmipLevel)
}

// FramebufferComplete - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:663
// Verify framebuffer is complete
func FramebufferComplete(id uint32) bool {
	cid := C.uint(id)
	return bool(C.rlFramebufferComplete(cid))
}

// UnloadFramebuffer - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:664
// Delete framebuffer from GPU
func UnloadFramebuffer(id uint32) {
	cid := C.uint(id)
	C.rlUnloadFramebuffer(cid)
}

// LoadShaderCode - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:667
// Shaders management
// Load shader from code strings
func LoadShaderCode(vsCode string, fsCode string) uint32 {
	cvsCode := C.CString(vsCode)
	defer C.free(unsafe.Pointer(cvsCode))
	cfsCode := C.CString(fsCode)
	defer C.free(unsafe.Pointer(cfsCode))
	return uint32(C.rlLoadShaderCode(cvsCode, cfsCode))
}

// CompileShader - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:668
// Compile custom shader and return shader id (type: RL_VERTEX_SHADER, RL_FRAGMENT_SHADER, RL_COMPUTE_SHADER)
func CompileShader(shaderCode string, type_ int32) uint32 {
	cshaderCode := C.CString(shaderCode)
	defer C.free(unsafe.Pointer(cshaderCode))
	ctype_ := C.int(type_)
	return uint32(C.rlCompileShader(cshaderCode, ctype_))
}

// LoadShaderProgram - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:669
// Load custom shader program
func LoadShaderProgram(vShaderId uint32, fShaderId uint32) uint32 {
	cvShaderId := C.uint(vShaderId)
	cfShaderId := C.uint(fShaderId)
	return uint32(C.rlLoadShaderProgram(cvShaderId, cfShaderId))
}

// UnloadShaderProgram - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:670
// Unload shader program
func UnloadShaderProgram(id uint32) {
	cid := C.uint(id)
	C.rlUnloadShaderProgram(cid)
}

// GetLocationUniform - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:671
// Get shader location uniform
func GetLocationUniform(shaderId uint32, uniformName string) int32 {
	cshaderId := C.uint(shaderId)
	cuniformName := C.CString(uniformName)
	defer C.free(unsafe.Pointer(cuniformName))
	return int32(C.rlGetLocationUniform(cshaderId, cuniformName))
}

// GetLocationAttrib - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:672
// Get shader location attribute
func GetLocationAttrib(shaderId uint32, attribName string) int32 {
	cshaderId := C.uint(shaderId)
	cattribName := C.CString(attribName)
	defer C.free(unsafe.Pointer(cattribName))
	return int32(C.rlGetLocationAttrib(cshaderId, cattribName))
}

// SetUniformMatrix - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:674
// Set shader value uniform
// Set shader value matrix
// Warning (*ast.FunctionDecl): {prefix: n:rlSetUniform,t1:void (int, const void *, int, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:673 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlSetUniform`. cannot parse C type: `const void *`
// func SetUniformMatrix(locIndex int32, mat Matrix) {
// 	clocIndex := C.int(locIndex)
// 	var cmat C.struct_Matrix
// 	cmat.m13 = C.float(mat.m13)
// 	cmat.m10 = C.float(mat.m10)
// 	cmat.m11 = C.float(mat.m11)
// 	cmat.m15 = C.float(mat.m15)
// 	cmat.m5 = C.float(mat.m5)
// 	cmat.m9 = C.float(mat.m9)
// 	cmat.m1 = C.float(mat.m1)
// 	cmat.m0 = C.float(mat.m0)
// 	cmat.m8 = C.float(mat.m8)
// 	cmat.m14 = C.float(mat.m14)
// 	cmat.m3 = C.float(mat.m3)
// 	cmat.m7 = C.float(mat.m7)
// 	cmat.m4 = C.float(mat.m4)
// 	cmat.m2 = C.float(mat.m2)
// 	cmat.m12 = C.float(mat.m12)
// 	cmat.m6 = C.float(mat.m6)
// 	C.rlSetUniformMatrix(clocIndex, cmat)
// }

// SetUniformSampler - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:675
// Set shader value sampler
func SetUniformSampler(locIndex int32, textureId uint32) {
	clocIndex := C.int(locIndex)
	ctextureId := C.uint(textureId)
	C.rlSetUniformSampler(clocIndex, ctextureId)
}

// LoadComputeShaderProgram - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:679
// Set shader currently active (id and locations)
// Compute shader management
// Load compute shader program
// Warning (*ast.FunctionDecl): {prefix: n:rlSetShader,t1:void (unsigned int, int *),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:676 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlSetShader`. cannot parse C type: `int *`
// func LoadComputeShaderProgram(shaderId uint32) uint32 {
// 	cshaderId := C.uint(shaderId)
// 	return C.rlLoadComputeShaderProgram(cshaderId)
// }

// ComputeShaderDispatch - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:680
// Dispatch compute shader (equivalent to *draw* for graphics pilepine)
func ComputeShaderDispatch(groupX uint32, groupY uint32, groupZ uint32) {
	cgroupX := C.uint(groupX)
	cgroupY := C.uint(groupY)
	cgroupZ := C.uint(groupZ)
	C.rlComputeShaderDispatch(cgroupX, cgroupY, cgroupZ)
}

// UnloadShaderBuffer - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:684
// Shader buffer storage object management (ssbo)
// Load shader storage buffer object (SSBO)
// Unload shader storage buffer object (SSBO)
// Warning (*ast.FunctionDecl): {prefix: n:rlLoadShaderBuffer,t1:unsigned int (unsigned int, const void *, int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:683 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlLoadShaderBuffer`. cannot parse C type: `const void *`
// func UnloadShaderBuffer(ssboId uint32) {
// 	cssboId := C.uint(ssboId)
// 	C.rlUnloadShaderBuffer(cssboId)
// }

// BindShaderBuffer - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:686
// Update SSBO buffer data
// Bind SSBO buffer
// Warning (*ast.FunctionDecl): {prefix: n:rlUpdateShaderBuffer,t1:void (unsigned int, const void *, unsigned int, unsigned int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:685 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlUpdateShaderBuffer`. cannot parse C type: `const void *`
// func BindShaderBuffer(id uint32, index uint32) {
// 	cid := C.uint(id)
// 	cindex := C.uint(index)
// 	C.rlBindShaderBuffer(cid, cindex)
// }

// CopyShaderBuffer - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:688
// Read SSBO buffer data (GPU->CPU)
// Copy SSBO data between buffers
// Warning (*ast.FunctionDecl): {prefix: n:rlReadShaderBuffer,t1:void (unsigned int, void *, unsigned int, unsigned int),t2:}.  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:687 :cannot transpileFunctionDecl. cannot bindingFunctionDecl func `rlReadShaderBuffer`. cannot parse C type: `void *`
// func CopyShaderBuffer(destId uint32, srcId uint32, destOffset uint32, srcOffset uint32, count uint32) {
// 	cdestId := C.uint(destId)
// 	csrcId := C.uint(srcId)
// 	cdestOffset := C.uint(destOffset)
// 	csrcOffset := C.uint(srcOffset)
// 	ccount := C.uint(count)
// 	C.rlCopyShaderBuffer(cdestId, csrcId, cdestOffset, csrcOffset, ccount)
// }

// GetShaderBufferSize - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:689
// Get SSBO buffer size
func GetShaderBufferSize(id uint32) uint32 {
	cid := C.uint(id)
	return uint32(C.rlGetShaderBufferSize(cid))
}

// BindImageTexture - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:692
// Buffer management
// Bind image texture
func BindImageTexture(id uint32, index uint32, format uint32, readonly int32) {
	cid := C.uint(id)
	cindex := C.uint(index)
	cformat := C.uint(format)
	creadonly := C.int(readonly)
	C.rlBindImageTexture(cid, cindex, cformat, creadonly)
}

// GetMatrixModelview - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:695
// Matrix state management
// Get internal modelview matrix
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

// GetMatrixProjection - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:696
// Get internal projection matrix
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

// GetMatrixTransform - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:697
// Get internal accumulated transform matrix
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

// GetMatrixProjectionStereo - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:698
// Get internal projection matrix for stereo render (selected eye)
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

// GetMatrixViewOffsetStereo - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:699
// Get internal view offset matrix for stereo render (selected eye)
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

// SetMatrixProjection - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:700
// Set a custom projection matrix (replaces internal projection matrix)
// func SetMatrixProjection(proj Matrix) {
// 	var cproj C.struct_Matrix
// 	cproj.m0 = C.float(proj.m0)
// 	cproj.m8 = C.float(proj.m8)
// 	cproj.m1 = C.float(proj.m1)
// 	cproj.m4 = C.float(proj.m4)
// 	cproj.m2 = C.float(proj.m2)
// 	cproj.m14 = C.float(proj.m14)
// 	cproj.m3 = C.float(proj.m3)
// 	cproj.m7 = C.float(proj.m7)
// 	cproj.m12 = C.float(proj.m12)
// 	cproj.m6 = C.float(proj.m6)
// 	cproj.m15 = C.float(proj.m15)
// 	cproj.m5 = C.float(proj.m5)
// 	cproj.m9 = C.float(proj.m9)
// 	cproj.m13 = C.float(proj.m13)
// 	cproj.m10 = C.float(proj.m10)
// 	cproj.m11 = C.float(proj.m11)
// 	C.rlSetMatrixProjection(cproj)
// }

// SetMatrixModelview - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:701
// Set a custom modelview matrix (replaces internal modelview matrix)
// func SetMatrixModelview(view Matrix) {
// 	var cview C.struct_Matrix
// 	cview.m3 = C.float(view.m3)
// 	cview.m7 = C.float(view.m7)
// 	cview.m4 = C.float(view.m4)
// 	cview.m2 = C.float(view.m2)
// 	cview.m14 = C.float(view.m14)
// 	cview.m12 = C.float(view.m12)
// 	cview.m6 = C.float(view.m6)
// 	cview.m10 = C.float(view.m10)
// 	cview.m11 = C.float(view.m11)
// 	cview.m15 = C.float(view.m15)
// 	cview.m5 = C.float(view.m5)
// 	cview.m9 = C.float(view.m9)
// 	cview.m13 = C.float(view.m13)
// 	cview.m0 = C.float(view.m0)
// 	cview.m8 = C.float(view.m8)
// 	cview.m1 = C.float(view.m1)
// 	C.rlSetMatrixModelview(cview)
// }

// SetMatrixProjectionStereo - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:702
// Set eyes projection matrices for stereo rendering
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

// SetMatrixViewOffsetStereo - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:703
// Set eyes view offsets matrices for stereo rendering
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

// LoadDrawCube - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:706
// Quick and dirty cube/quad buffers load->draw->unload
// Load and draw a cube
func LoadDrawCube() {
	C.rlLoadDrawCube()
}

// LoadDrawQuad - transpiled function from  GOPATH/src/github.com/Konstantin8105/raylib-go/raylib/rlgl.h:707
// Load and draw a quad
func LoadDrawQuad() {
	C.rlLoadDrawQuad()
}

// type _Bool int32

//
//*
//*   RLGL IMPLEMENTATION
//*
//
