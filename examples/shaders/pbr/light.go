package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"unsafe"
)

type LightType int32

const (
	LightTypeDirectional LightType = iota
	LightTypePoint
	LightTypeSpot
)

type Light struct {
	Shader        rl.Shader
	combineStatus bool
	lightType     LightType
	position      rl.Vector3
	direction     rl.Vector3
	lightColor    rl.Color
	enabled       int32
	enargy        float32
	cutOff        float32
	outerCutOff   float32
	constant      float32
	linear        float32
	quadratic     float32
	shiny         float32
	specularStr   float32
	// shader locations
	enabledLoc     int32
	typeLoc        int32
	posLoc         int32
	dirLoc         int32
	colorLoc       int32
	viewPosLoc     int32
	enargyLoc      int32
	cutOffLoc      int32
	outerCutOffLoc int32
	constantLoc    int32
	linearLoc      int32
	quadraticLoc   int32
	shinyLoc       int32
	specularStrLoc int32
}

var MaxLights = 5
var lightCount = 0

func (lt *Light) NewLight(
	lightType LightType,
	position, direction rl.Vector3,
	color rl.Color,
	enargy float32, shader *rl.Shader) Light {

	light := Light{
		Shader: *shader,
	}

	if lightCount < MaxLights {
		light.enabled = 1
		light.lightType = lightType
		light.position = position
		light.direction = direction
		light.lightColor = color
		light.enargy = enargy
		light.cutOff = 12.5
		light.outerCutOff = 17.5
		light.constant = 1.0
		light.linear = 0.09
		light.quadratic = 0.032
		light.shiny = 32.0
		light.specularStr = 0.9
		light.enabledLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].enabled", lightCount))
		light.typeLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].type", lightCount))
		light.posLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].position", lightCount))
		light.dirLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].direction", lightCount))
		light.colorLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].lightColor", lightCount))
		light.enargyLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].enargy", lightCount))
		light.cutOffLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].cutOff", lightCount))
		light.outerCutOffLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].outerCutOff", lightCount))
		light.constantLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].constant", lightCount))
		light.linearLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].linear", lightCount))
		light.quadraticLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].quadratic", lightCount))
		light.shinyLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].shiny", lightCount))
		light.specularStrLoc = rl.GetShaderLocation(*shader, fmt.Sprintf("lights[%d].specularStr", lightCount))

		light.UpdateValues()
		lightCount++
	}
	return light
}

func (lt *Light) UpdateValues() {
	// Send to shader light enabled state and type
	rl.SetShaderValue(lt.Shader, lt.enabledLoc, unsafe.Slice((*float32)(unsafe.Pointer(&lt.enabled)), 4), rl.ShaderUniformInt)
	rl.SetShaderValue(lt.Shader, lt.typeLoc, unsafe.Slice((*float32)(unsafe.Pointer(&lt.lightType)), 4), rl.ShaderUniformInt)

	// Send to shader light position values
	rl.SetShaderValue(lt.Shader, lt.posLoc, []float32{lt.position.X, lt.position.Y, lt.position.Z}, rl.ShaderUniformVec3)

	// Send to shader light direction values
	rl.SetShaderValue(lt.Shader, lt.dirLoc, []float32{lt.direction.X, lt.direction.Y, lt.direction.Z}, rl.ShaderUniformVec3)

	// Send to shader light color values
	rl.SetShaderValue(lt.Shader, lt.colorLoc,
		[]float32{float32(lt.lightColor.R) / 255, float32(lt.lightColor.G) / 255, float32(lt.lightColor.B) / 255},
		rl.ShaderUniformVec3)

	// Send to shader light enargy values
	rl.SetShaderValue(lt.Shader, lt.enargyLoc, []float32{lt.enargy}, rl.ShaderUniformFloat)

	// Send to shader light spot values
	rl.SetShaderValue(lt.Shader, lt.cutOffLoc, []float32{lt.cutOff}, rl.ShaderUniformFloat)
	rl.SetShaderValue(lt.Shader, lt.outerCutOffLoc, []float32{lt.outerCutOff}, rl.ShaderUniformFloat)

	// Send to shader light pointLight values
	rl.SetShaderValue(lt.Shader, lt.constantLoc, []float32{lt.constant}, rl.ShaderUniformFloat)
	rl.SetShaderValue(lt.Shader, lt.linearLoc, []float32{lt.linear}, rl.ShaderUniformFloat)
	rl.SetShaderValue(lt.Shader, lt.quadraticLoc, []float32{lt.quadratic}, rl.ShaderUniformFloat)

	// Send to shader light shiness values
	rl.SetShaderValue(lt.Shader, lt.shinyLoc, []float32{lt.shiny}, rl.ShaderUniformFloat)
	rl.SetShaderValue(lt.Shader, lt.specularStrLoc, []float32{lt.specularStr}, rl.ShaderUniformFloat)

}

// if you want more 5 light in the light.fs change #define MAX_LIGHTS "your number"
func (lt *Light) SetMaxLight(max int) {
	MaxLights = max
}

func (lt *Light) SetConfigSpotLight(light *Light, cutOff, outerCutOff float32) {
	light.cutOff = cutOff
	light.outerCutOff = outerCutOff
	light.UpdateValues()
}

func (lt *Light) SetConfigPointLight(light *Light, constant, linear, quadratic float32) {
	light.constant = constant
	light.linear = linear
	light.quadratic = quadratic
	light.UpdateValues()
}

func (lt *Light) SetConfigShiness(light *Light, shiny, specularStr float32) {
	light.shiny = shiny
	light.specularStr = specularStr
	light.UpdateValues()
}

func (lt *Light) SetMaterialTexture(materials []*rl.Material, texture []*rl.Texture2D) {
	for index, material := range materials {
		rl.SetMaterialTexture(material, rl.MapDiffuse, *texture[index])
	}
}

func (lt *Light) SetFlashlightTexture(materials []*rl.Material, texure *rl.Texture2D) {

	lt.Shader.UpdateLocation(rl.ShaderLocMapOcclusion, rl.GetShaderLocation(lt.Shader, "flashlight"))
	for _, material := range materials {
		rl.SetMaterialTexture(material, rl.MapOcclusion, *texure)
	}
}

func (lt *Light) Init(ambientStrength float32, ambientColor rl.Vector3) {
	if !lt.combineStatus {
		lt.configShader()
	}
	lt.viewPosLoc = rl.GetShaderLocation(lt.Shader, "viewPos")
	rl.SetShaderValue(lt.Shader, rl.GetShaderLocation(lt.Shader, "ambientColor"), []float32{ambientColor.X, ambientColor.Y, ambientColor.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(lt.Shader, rl.GetShaderLocation(lt.Shader, "ambientStrength"), []float32{ambientStrength}, rl.ShaderUniformFloat)
}

func (lt *Light) DisableLight(light *Light) {
	light.enabled *= -1
	light.UpdateValues()
}

func (lt *Light) EnableLight(light *Light) {
	light.enabled = 1
	light.UpdateValues()
}

func (lt *Light) DrawSpherelight(light *Light) {
	if light.enabled == 1 {
		rl.DrawSphereEx(light.position, 0.2, 8, 8, light.lightColor)
	} else {
		rl.DrawSphereWires(light.position, 0.2, 8, 8, rl.Fade(light.lightColor, 0.3))
	}
}

func (lt *Light) UpdateReflect(cameraPos rl.Vector3) {
	rl.SetShaderValue(lt.Shader, lt.viewPosLoc, []float32{cameraPos.X, cameraPos.Y, cameraPos.Z}, rl.ShaderUniformVec3)
}

func (lt *Light) configShader() {
		lt.Shader = rl.LoadShader("pbr.vs","./pbr.fs")
}

// exce before init or set manually
func (lt *Light) SetCombineShader(CombineShader *rl.Shader) {
	lt.combineStatus = true
	lt.Shader = *CombineShader
}

func (lt *Light) Unload() {
	rl.UnloadShader(lt.Shader)
}
