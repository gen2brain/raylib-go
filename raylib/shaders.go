package rl

/*
#include "raylib.h"
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

// LoadShaderCode - Load shader from code strings and bind default locations
func LoadShaderCode(vsCode string, fsCode string) Shader {
	cvsCode := C.CString(vsCode)
	defer C.free(unsafe.Pointer(cvsCode))

	cfsCode := C.CString(fsCode)
	defer C.free(unsafe.Pointer(cfsCode))

	ret := C.LoadShaderCode(cvsCode, cfsCode)
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
	ret := C.GetShaderDefault()
	v := newShaderFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetTextureDefault - Get default texture
func GetTextureDefault() *Texture2D {
	ret := C.GetTextureDefault()
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
	cvalue := (*C.float)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&value)).Data))
	csize := (C.int)(size)
	C.SetShaderValue(*cshader, cuniformLoc, cvalue, csize)
}

// SetShaderValuei - Set shader uniform value (int)
func SetShaderValuei(shader Shader, uniformLoc int32, value []int32, size int32) {
	cshader := shader.cptr()
	cuniformLoc := (C.int)(uniformLoc)
	cvalue := (*C.int)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&value)).Data))
	csize := (C.int)(size)
	C.SetShaderValuei(*cshader, cuniformLoc, cvalue, csize)
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
	C.SetMatrixProjection(*cproj)
}

// SetMatrixModelview - Set a custom modelview matrix (replaces internal modelview matrix)
func SetMatrixModelview(view Matrix) {
	cview := view.cptr()
	C.SetMatrixModelview(*cview)
}

// GenTextureCubemap - Generate cubemap texture from HDR texture
func GenTextureCubemap(shader Shader, skyHDR Texture2D, size int) Texture2D {
	cshader := shader.cptr()
	cskyHDR := skyHDR.cptr()
	csize := (C.int)(size)

	ret := C.GenTextureCubemap(*cshader, *cskyHDR, csize)
	v := newTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenTextureIrradiance - Generate irradiance texture using cubemap data
func GenTextureIrradiance(shader Shader, cubemap Texture2D, size int) Texture2D {
	cshader := shader.cptr()
	ccubemap := cubemap.cptr()
	csize := (C.int)(size)

	ret := C.GenTextureIrradiance(*cshader, *ccubemap, csize)
	v := newTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenTexturePrefilter - Generate prefilter texture using cubemap data
func GenTexturePrefilter(shader Shader, cubemap Texture2D, size int) Texture2D {
	cshader := shader.cptr()
	ccubemap := cubemap.cptr()
	csize := (C.int)(size)

	ret := C.GenTexturePrefilter(*cshader, *ccubemap, csize)
	v := newTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenTextureBRDF - Generate BRDF texture using cubemap data
func GenTextureBRDF(shader Shader, cubemap Texture2D, size int) Texture2D {
	cshader := shader.cptr()
	ccubemap := cubemap.cptr()
	csize := (C.int)(size)

	ret := C.GenTextureBRDF(*cshader, *ccubemap, csize)
	v := newTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
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

// GetVrDeviceInfo - Get VR device information for some standard devices
func GetVrDeviceInfo(vrDevice VrDevice) VrDeviceInfo {
	cvrDevice := (C.int)(vrDevice)
	ret := C.GetVrDeviceInfo(cvrDevice)
	v := newVrDeviceInfoFromPointer(unsafe.Pointer(&ret))
	return v
}

// InitVrSimulator - Init VR simulator for selected device
func InitVrSimulator(vrDeviceInfo VrDeviceInfo) {
	cvrDeviceInfo := vrDeviceInfo.cptr()
	C.InitVrSimulator(*cvrDeviceInfo)
}

// CloseVrSimulator - Close VR simulator for current device
func CloseVrSimulator() {
	C.CloseVrSimulator()
}

// IsVrSimulatorReady - Detect if VR simulator is ready
func IsVrSimulatorReady() bool {
	ret := C.IsVrSimulatorReady()
	v := bool(int(ret) == 1)
	return v
}

// UpdateVrTracking - Update VR tracking (position and orientation) and camera
func UpdateVrTracking(camera *Camera) {
	ccamera := camera.cptr()
	C.UpdateVrTracking(ccamera)
}

// ToggleVrMode - Enable/Disable VR experience (device or simulator)
func ToggleVrMode() {
	C.ToggleVrMode()
}

// BeginVrDrawing - Begin VR simulator stereo rendering
func BeginVrDrawing() {
	C.BeginVrDrawing()
}

// EndVrDrawing - End VR simulator stereo rendering
func EndVrDrawing() {
	C.EndVrDrawing()
}
