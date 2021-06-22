package main

import (
	"fmt"
	"reflect"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LightType int32

const (
	LightTypeDirectional LightType = iota
	LightTypePoint
)

type Light struct {
	shader    rl.Shader
	lightType LightType
	position  rl.Vector3
	target    rl.Vector3
	color     rl.Color
	enabled   int32

	// shader locations
	enabledLoc int32
	typeLoc    int32
	posLoc     int32
	targetLoc  int32
	colorLoc   int32
}

const maxLightsCount = 4

var lightCount = 0

func NewLight(
	lightType LightType,
	position, target rl.Vector3,
	color rl.Color,
	shader rl.Shader) Light {

	light := Light{
		shader: shader,
	}

	if lightCount < maxLightsCount {
		light.enabled = 1
		light.lightType = lightType
		light.position = position
		light.target = target
		light.color = color

		light.enabledLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].enabled", lightCount))
		light.typeLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].type", lightCount))
		light.posLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].position", lightCount))
		light.targetLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].target", lightCount))
		light.colorLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].color", lightCount))
		light.UpdateValues()

		lightCount++
	}

	return light
}

// Send light properties to shader
func (lt *Light) UpdateValues() {
	// not pretty -_-
	// need nicer api
	sh := &reflect.SliceHeader{
		Len: 4,
		Cap: 4,
	}

	// Send to shader light enabled state and type
	sh.Data = uintptr(unsafe.Pointer(&lt.enabled))
	rl.SetShaderValue(lt.shader, lt.enabledLoc, *(*[]float32)(unsafe.Pointer(sh)), rl.ShaderUniformInt)

	// Send to shader light position values
	sh.Data = uintptr(unsafe.Pointer(&lt.lightType))
	rl.SetShaderValue(lt.shader, lt.posLoc, []float32{lt.position.X, lt.position.Y, lt.position.Z}, rl.ShaderUniformVec3)

	// Send to shader light target target values
	rl.SetShaderValue(lt.shader, lt.targetLoc, []float32{lt.target.X, lt.target.Y, lt.target.Z}, rl.ShaderUniformVec3)

	// Send to shader light color values
	rl.SetShaderValue(lt.shader, lt.colorLoc,
		[]float32{float32(lt.color.R) / 255, float32(lt.color.G) / 255, float32(lt.color.B) / 255, float32(lt.color.A) / 255},
		rl.ShaderUniformVec4)
}
