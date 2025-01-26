package main

import (
	"fmt"
	"reflect"
	"unsafe"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	LIGHT_DIRECTIONAL int32 = iota
	LIGHT_POINT
	LIGHT_SPOT
)

type light struct {
	Enabled   int32
	Type      int32
	Target    rl.Vector3
	Position  rl.Vector3
	Color     [4]float32
	Intensity float32
	// Shader light parameters locations
	enabledLoc   int32
	typeLoc      int32
	targetLoc    int32
	positionLoc  int32
	colorLoc     int32
	intensityLoc int32
}

type PhysicRender struct {
	Shader        *rl.Shader
	combineStatus bool
}

// make light and pbr Shader
func (ph *PhysicRender) Init() {

	if !ph.combineStatus {
		ph.configShader()
	}

	ph.Shader.UpdateLocation(rl.ShaderLocMapAlbedo, rl.GetShaderLocation(*ph.Shader, "albedoMap"))
	ph.Shader.UpdateLocation(rl.ShaderLocMapMetalness, rl.GetShaderLocation(*ph.Shader, "mraMap"))
	ph.Shader.UpdateLocation(rl.ShaderLocMapNormal, rl.GetShaderLocation(*ph.Shader, "normalMap"))
	ph.Shader.UpdateLocation(rl.ShaderLocMapEmission, rl.GetShaderLocation(*ph.Shader, "emissiveMap"))
	ph.Shader.UpdateLocation(12, rl.GetShaderLocation(*ph.Shader, "albedoColor"))
	ph.Shader.UpdateLocation(rl.ShaderLocVectorView, rl.GetShaderLocation(*ph.Shader, "viewPos"))

	ambientColor := rl.Color{R: 122, G: 36, B: 26, A: 100}
	ambientColorNormalized := rl.NewVector3(float32(ambientColor.R)/255.0, float32(ambientColor.G)/255.0, float32(ambientColor.B)/255.0)
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "ambientColor"), []float32{ambientColorNormalized.X, ambientColorNormalized.Y, ambientColorNormalized.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "ambient"), []float32{float32(0.03)}, rl.ShaderUniformFloat)

	lightCountLoc := rl.GetShaderLocation(*ph.Shader, "numOfLights")
	rl.SetShaderValue(*ph.Shader, lightCountLoc, generiteIntForGlsl(int32(MaxLights)), rl.ShaderUniformInt)

}

func (ph *PhysicRender) UpadteByCamera(pos rl.Vector3) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "viewPos"), []float32{pos.X, pos.Y, pos.Z}, rl.ShaderUniformVec3)
}

func (ph *PhysicRender) CreateLight(typeLight int32, position rl.Vector3, target rl.Vector3, color rl.Color, intensity float32) light {
	light := light{}
	if lightCount < MaxLights {
		light.Enabled = 1
		light.Position = position
		light.Type = typeLight
		light.Target = target
		light.Color[0] = float32(float32(color.R) / 255.0)
		light.Color[1] = float32(float32(color.G) / 255.0)
		light.Color[2] = float32(float32(color.B) / 255.0)
		light.Color[3] = float32(float32(color.A) / 255.0)
		light.Intensity = intensity
		// NOTE: Shader parameters names for lights must match the requested ones
		light.enabledLoc = rl.GetShaderLocation(*ph.Shader, fmt.Sprintf("lights[%d].enabled", lightCount))
		light.positionLoc = rl.GetShaderLocation(*ph.Shader, fmt.Sprintf("lights[%d].position", lightCount))
		light.colorLoc = rl.GetShaderLocation(*ph.Shader, fmt.Sprintf("lights[%d].color", lightCount))
		light.intensityLoc = rl.GetShaderLocation(*ph.Shader, fmt.Sprintf("lights[%d].intensity", lightCount))
		light.typeLoc = rl.GetShaderLocation(*ph.Shader, fmt.Sprintf("lights[%d].type", lightCount))
		//light.targetLoc = rl.GetShaderLocation(ph.Shader, fmt.Sprintf("lights[%d].target", lightCount))

		ph.UpdateLight(*ph.Shader, light)

		lightCount++

	}

	return light
}

func (ph *PhysicRender) UpdateLight(Shader rl.Shader, light light) {
	rl.SetShaderValue(Shader, light.enabledLoc, generiteIntForGlsl(light.Enabled), rl.ShaderUniformInt)
	rl.SetShaderValue(Shader, light.positionLoc, []float32{light.Position.X, light.Position.Y, light.Position.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(Shader, light.colorLoc, []float32{light.Color[0], light.Color[1], light.Color[2], light.Color[3]}, rl.ShaderUniformVec4)
	rl.SetShaderValue(Shader, light.intensityLoc, []float32{light.Intensity}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Shader, light.typeLoc, generiteIntForGlsl(light.Type), rl.ShaderUniformInt)
	//rl.SetShaderValue(Shader, light.targetLoc,[]float32{light.Target.X, light.Target.Y, light.Target.Z}, rl.ShaderUniformVec3)
}

func (ph *PhysicRender) MaxLights(value int) {
	MaxLights = value
}

func (ph *PhysicRender) TextureMapAlbedo(modelMaterials *rl.Material, texture rl.Texture2D) {
	rl.SetMaterialTexture(modelMaterials, rl.MapAlbedo, texture)
}

func (ph *PhysicRender) TextureMapMetalness(modelMaterials *rl.Material, texture rl.Texture2D) {
	rl.SetMaterialTexture(modelMaterials, rl.MapMetalness, texture)
}

func (ph *PhysicRender) TextureMapRoughness(modelMaterials *rl.Material, texture rl.Texture2D) {
	rl.SetMaterialTexture(modelMaterials, rl.MapRoughness, texture)
}

func (ph *PhysicRender) TextureMapNormal(modelMaterials *rl.Material, texture rl.Texture2D) {
	rl.SetMaterialTexture(modelMaterials, rl.MapNormal, texture)
}

func (ph *PhysicRender) TextureMapOcclusion(modelMaterials *rl.Material, texture rl.Texture2D) {
	rl.SetMaterialTexture(modelMaterials, rl.MapOcclusion, texture)
}

func (ph *PhysicRender) DrawSphereLoctionLight(li light, color rl.Color) {
	rl.DrawSphereEx(li.Position, 0.2, 8, 8, color)
}

func (ph *PhysicRender) AlbedoColorModel(model *rl.Model, color rl.Color) {
	model.GetMaterials()[0].GetMap(rl.MapAlbedo).Color = color
}

func (ph *PhysicRender) EmissiveColor(color rl.Color) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "emissiveColor"), []float32{float32(color.R), float32(color.G), float32(color.B), float32(color.A)}, rl.ShaderUniformVec4)
}

func (ph *PhysicRender) NormalValue(value float32) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "emissiveColor"), []float32{value}, rl.ShaderUniformFloat)
}

func (ph *PhysicRender) MetallicValue(value float32) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "metallicValue"), []float32{value}, rl.ShaderUniformFloat)
}

func (ph *PhysicRender) RoughnessValue(value float32) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "roughnessValue"), []float32{value}, rl.ShaderUniformFloat)
}

func (ph *PhysicRender) AoValue(value float32) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "aoValue"), []float32{value}, rl.ShaderUniformFloat)
}

func (ph *PhysicRender) EmissivePower(value float32) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "emissivePower"), []float32{value}, rl.ShaderUniformFloat)
}

func (ph *PhysicRender) AmbientColor(colorAmbient rl.Vector3, ambientValue float32) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "ambientColor"), []float32{colorAmbient.X, colorAmbient.Y, colorAmbient.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "ambient"), []float32{ambientValue}, rl.ShaderUniformFloat)
}

func (ph *PhysicRender) SetTiling(value rl.Vector2) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "tiling"), []float32{value.X, value.Y}, rl.ShaderUniformVec2)
}

func (ph *PhysicRender) SetOffset(value rl.Vector2) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "offset"), []float32{value.X, value.Y}, rl.ShaderUniformVec2)
}

func (ph *PhysicRender) SetTilingFlashlight(value rl.Vector2) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "tilingFlashlight"), []float32{value.X, value.Y}, rl.ShaderUniformVec2)
}

func (ph *PhysicRender) SetOffsetFlashlight(value rl.Vector2) {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "offsetFlashlight"), []float32{value.X, value.Y}, rl.ShaderUniformVec2)
}

func (ph *PhysicRender) UseTexAlbedo() {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "useTexAlbedo"), generiteIntForGlsl(1), rl.ShaderUniformInt)
}

func (ph *PhysicRender) UseTexNormal() {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "useTexNormal"), generiteIntForGlsl(1), rl.ShaderUniformInt)
}

func (ph *PhysicRender) UseTexMRA() {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "useTexMRA"), generiteIntForGlsl(1), rl.ShaderUniformInt)
}

func (ph *PhysicRender) UseTexEmissive() {
	rl.SetShaderValue(*ph.Shader, rl.GetShaderLocation(*ph.Shader, "useTexEmissive"), generiteIntForGlsl(1), rl.ShaderUniformInt)
}

func (ph *PhysicRender) configShader() {
	sh := rl.LoadShader("./pbr.vs", "./pbr.fs")
	ph.Shader = &sh
}

// exce before init or set manually
func (ph *PhysicRender) SetCombineShader(CombineShader *rl.Shader) {
	ph.combineStatus = true
	ph.Shader = CombineShader
}

func (ph *PhysicRender) Unload() {
	rl.UnloadShader(*ph.Shader)
}

func generiteIntForGlsl(value int32) []float32 {
	data := &reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&value)),
		Len:  4,
		Cap:  4,
	}
	return *(*[]float32)(unsafe.Pointer(data))
}
