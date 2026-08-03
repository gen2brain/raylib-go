package main

import (
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LightType int32

const (
	LightTypeDirectional LightType = iota
	LightTypePoint
)

type Light struct {
	Shader    rl.Shader
	LightType LightType
	Position  rl.Vector3
	Target    rl.Vector3
	Color     rl.Color
	Enabled   int32

	// Shader location indexes
	enabledLoc int32
	typeLoc    int32
	posLoc     int32
	targetLoc  int32
	colorLoc   int32
}

const maxLightsCount = 4

var lightCount = 0

// CreateLight initializes a light and sends initial values to the shader
func CreateLight(
	lightType LightType,
	position, target rl.Vector3,
	color rl.Color,
	shader rl.Shader,
) Light {
	light := Light{
		Shader: shader,
	}

	if lightCount < maxLightsCount {
		light.Enabled = 1
		light.LightType = lightType
		light.Position = position
		light.Target = target
		light.Color = color

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

// NewLight alias for CreateLight
func NewLight(
	lightType LightType,
	position, target rl.Vector3,
	color rl.Color,
	shader rl.Shader,
) Light {
	return CreateLight(lightType, position, target, color, shader)
}

// UpdateValues sends updated light attributes to the GLSL shader
func (lt *Light) UpdateValues() {
	rl.SetShaderValue(lt.Shader, lt.enabledLoc, unsafe.Slice((*float32)(unsafe.Pointer(&lt.Enabled)), 1), rl.ShaderUniformInt)
	rl.SetShaderValue(lt.Shader, lt.typeLoc, unsafe.Slice((*float32)(unsafe.Pointer(&lt.LightType)), 1), rl.ShaderUniformInt)

	rl.SetShaderValue(lt.Shader, lt.posLoc, []float32{lt.Position.X, lt.Position.Y, lt.Position.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(lt.Shader, lt.targetLoc, []float32{lt.Target.X, lt.Target.Y, lt.Target.Z}, rl.ShaderUniformVec3)

	rl.SetShaderValue(lt.Shader, lt.colorLoc,
		[]float32{float32(lt.Color.R) / 255.0, float32(lt.Color.G) / 255.0, float32(lt.Color.B) / 255.0, float32(lt.Color.A) / 255.0},
		rl.ShaderUniformVec4)
}
