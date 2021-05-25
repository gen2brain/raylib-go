package rl

/*
#include "raylib.h"
#include "rlgl.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "reflect"

// cptr returns C pointer
func (v *VrDeviceInfo) cptr() *C.VrDeviceInfo {
	return (*C.VrDeviceInfo)(unsafe.Pointer(v))
}

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

	var v Shader
	if vsFileName == "" {
		ret := C.LoadShader(nil, cfsFileName)
		v = newShaderFromPointer(unsafe.Pointer(&ret))
	} else {
		ret := C.LoadShader(cvsFileName, cfsFileName)
		v = newShaderFromPointer(unsafe.Pointer(&ret))
	}

	return v
}

// LoadShaderFromMemory - Load shader from code strings and bind default locations
func LoadShaderFromMemory(vsCode string, fsCode string) Shader {
	cvsCode := C.CString(vsCode)
	defer C.free(unsafe.Pointer(cvsCode))

	cfsCode := C.CString(fsCode)
	defer C.free(unsafe.Pointer(cfsCode))

	ret := C.LoadShaderFromMemory(cvsCode, cfsCode)
	v := newShaderFromPointer(unsafe.Pointer(&ret))

	return v
}

// UnloadShader - Unload a custom shader from memory
func UnloadShader(shader Shader) {
	cshader := shader.cptr()
	C.UnloadShader(*cshader)
}

// GetShaderDefault - Get default shader
func GetShaderDefault() Shader {
	ret := C.rlGetShaderDefault()
	v := newShaderFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetTextureDefault - Get default texture
func GetTextureDefault() *Texture2D {
	ret := C.rlGetTextureDefault()
	v := newTexture2DFromPointer(unsafe.Pointer(&ret))
	return &v
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

// SetShaderValue - Set shader uniform value (float)
func SetShaderValue(shader Shader, uniformLoc int32, value []float32, size int32) {
	cshader := shader.cptr()
	cuniformLoc := (C.int)(uniformLoc)
	cvalue := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&value)).Data)
	csize := (C.int)(size)
	C.SetShaderValue(*cshader, cuniformLoc, cvalue, csize)
}

// SetShaderValueMatrix - Set shader uniform value (matrix 4x4)
func SetShaderValueMatrix(shader Shader, uniformLoc int32, mat Matrix) {
	cshader := shader.cptr()
	cuniformLoc := (C.int)(uniformLoc)
	cmat := mat.cptr()
	C.SetShaderValueMatrix(*cshader, cuniformLoc, *cmat)
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
