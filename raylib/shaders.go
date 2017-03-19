package raylib

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "reflect"

// VrDevice type
type VrDevice int32

// Head Mounted Display devices
const (
	HmdDefaultDevice     VrDevice = C.HMD_DEFAULT_DEVICE
	HmdOculusRiftDk2     VrDevice = C.HMD_OCULUS_RIFT_DK2
	HmdOculusRiftCv1     VrDevice = C.HMD_OCULUS_RIFT_CV1
	HmdValveHtcVive      VrDevice = C.HMD_VALVE_HTC_VIVE
	HmdSamsungGearVr     VrDevice = C.HMD_SAMSUNG_GEAR_VR
	HmdGoogleCardboard   VrDevice = C.HMD_GOOGLE_CARDBOARD
	HmdSonyPlaystationVr VrDevice = C.HMD_SONY_PLAYSTATION_VR
	HmdRazerOsvr         VrDevice = C.HMD_RAZER_OSVR
	HmdFoveVr            VrDevice = C.HMD_FOVE_VR
)

// BlendMode type
type BlendMode int32

// Color blending modes (pre-defined)
const (
	BlendAlpha      BlendMode = C.BLEND_ALPHA
	BlendAdditive   BlendMode = C.BLEND_ADDITIVE
	BlendMultiplied BlendMode = C.BLEND_MULTIPLIED
)

// Shader type (generic shader)
type Shader struct {
	// Shader program id
	ID uint32
	// Vertex attribute location point (default-location = 0)
	VertexLoc int32
	// Texcoord attribute location point (default-location = 1)
	TexcoordLoc int32
	// Texcoord2 attribute location point (default-location = 5)
	Texcoord2Loc int32
	// Normal attribute location point (default-location = 2)
	NormalLoc int32
	// Tangent attribute location point (default-location = 4)
	TangentLoc int32
	// Color attibute location point (default-location = 3)
	ColorLoc int32
	// ModelView-Projection matrix uniform location point (vertex shader)
	MvpLoc int32
	// Diffuse color uniform location point (fragment shader)
	ColDiffuseLoc int32
	// Ambient color uniform location point (fragment shader)
	ColAmbientLoc int32
	// Specular color uniform location point (fragment shader)
	ColSpecularLoc int32
	// Map texture uniform location point (default-texture-unit = 0)
	MapTexture0Loc int32
	// Map texture uniform location point (default-texture-unit = 1)
	MapTexture1Loc int32
	// Map texture uniform location point (default-texture-unit = 2)
	MapTexture2Loc int32
}

func (s *Shader) cptr() *C.Shader {
	return (*C.Shader)(unsafe.Pointer(s))
}

// NewShader - Returns new Shader
func NewShader(id uint32, vertexLoc, texcoordLoc, texcoord2Loc, normalLoc, tangentLoc, colorLoc, mvpLoc, colDiffuseLoc, colAmbientLoc, colSpecularLoc, mapTexture0Loc, mapTexture1Loc, mapTexture2Loc int32) Shader {
	return Shader{id, vertexLoc, texcoordLoc, texcoord2Loc, normalLoc, tangentLoc, colorLoc, mvpLoc, colDiffuseLoc, colAmbientLoc, colSpecularLoc, mapTexture0Loc, mapTexture1Loc, mapTexture2Loc}
}

// NewShaderFromPointer - Returns new Shader from pointer
func NewShaderFromPointer(ptr unsafe.Pointer) Shader {
	return *(*Shader)(ptr)
}

// LoadShader - Load a custom shader and bind default locations
func LoadShader(vsFileName string, fsFileName string) Shader {
	cvsFileName := C.CString(vsFileName)
	cfsFileName := C.CString(fsFileName)
	defer C.free(unsafe.Pointer(cvsFileName))
	defer C.free(unsafe.Pointer(cfsFileName))

	ret := C.LoadShader(cvsFileName, cfsFileName)
	v := NewShaderFromPointer(unsafe.Pointer(&ret))
	return v
}

// UnloadShader - Unload a custom shader from memory
func UnloadShader(shader Shader) {
	cshader := shader.cptr()
	C.UnloadShader(*cshader)
}

// GetDefaultShader - Get default shader
func GetDefaultShader() Shader {
	ret := C.GetDefaultShader()
	v := NewShaderFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetDefaultTexture - Get default texture
func GetDefaultTexture() *Texture2D {
	ret := C.GetDefaultTexture()
	v := NewTexture2DFromPointer(unsafe.Pointer(&ret))
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

// InitVrSimulator - Init VR simulator for selected device
func InitVrSimulator(vdDevice VrDevice) {
	cvdDevice := (C.int)(vdDevice)
	C.InitVrSimulator(cvdDevice)
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
