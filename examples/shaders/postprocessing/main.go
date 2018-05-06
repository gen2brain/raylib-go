package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const MaxPostproShaders = 12

const (
	FxGrayscale = iota
	FxPosterization
	FxDreamVision
	FxPixelizer
	FxCrossHatching
	FxCrossStitching
	FxPredatorView
	FxScanlines
	FxFisheye
	FxSobel
	FxBloom
	FxBlur
)

var postproShaderText = []string{
	"GRAYSCALE",
	"POSTERIZATION",
	"DREAM_VISION",
	"PIXELIZER",
	"CROSS_HATCHING",
	"CROSS_STITCHING",
	"PREDATOR_VIEW",
	"SCANLINES",
	"FISHEYE",
	"SOBEL",
	"BLOOM",
	"BLUR",
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.SetConfigFlags(raylib.FlagMsaa4xHint | raylib.FlagVsyncHint) // Enable Multi Sampling Anti Aliasing 4x (if available)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - postprocessing shader")

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(3.0, 3.0, 3.0)
	camera.Target = raylib.NewVector3(0.0, 1.5, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0

	dwarf := raylib.LoadModel("dwarf.obj")                   // Load OBJ model
	texture := raylib.LoadTexture("dwarf_diffuse.png")       // Load model texture
	dwarf.Material.Maps[raylib.MapDiffuse].Texture = texture // Set dwarf model diffuse texture

	position := raylib.NewVector3(0.0, 0.0, 0.0) // Set model position

	// Load all postpro shaders
	// NOTE 1: All postpro shader use the base vertex shader (DEFAULT_VERTEX_SHADER)
	shaders := make([]raylib.Shader, MaxPostproShaders)
	shaders[FxGrayscale] = raylib.LoadShader("", "glsl330/grayscale.fs")
	shaders[FxPosterization] = raylib.LoadShader("", "glsl330/posterization.fs")
	shaders[FxDreamVision] = raylib.LoadShader("", "glsl330/dream_vision.fs")
	shaders[FxPixelizer] = raylib.LoadShader("", "glsl330/pixelizer.fs")
	shaders[FxCrossHatching] = raylib.LoadShader("", "glsl330/cross_hatching.fs")
	shaders[FxCrossStitching] = raylib.LoadShader("", "glsl330/cross_stitching.fs")
	shaders[FxPredatorView] = raylib.LoadShader("", "glsl330/predator.fs")
	shaders[FxScanlines] = raylib.LoadShader("", "glsl330/scanlines.fs")
	shaders[FxFisheye] = raylib.LoadShader("", "glsl330/fisheye.fs")
	shaders[FxSobel] = raylib.LoadShader("", "glsl330/sobel.fs")
	shaders[FxBlur] = raylib.LoadShader("", "glsl330/blur.fs")
	shaders[FxBloom] = raylib.LoadShader("", "glsl330/bloom.fs")

	currentShader := FxGrayscale

	// Create a RenderTexture2D to be used for render to texture
	target := raylib.LoadRenderTexture(screenWidth, screenHeight)

	raylib.SetCameraMode(camera, raylib.CameraOrbital) // Set free camera mode

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera) // Update camera

		if raylib.IsKeyPressed(raylib.KeyRight) {
			currentShader++
		} else if raylib.IsKeyPressed(raylib.KeyLeft) {
			currentShader--
		}

		if currentShader >= MaxPostproShaders {
			currentShader = 0
		} else if currentShader < 0 {
			currentShader = MaxPostproShaders - 1
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginTextureMode(target) // Enable drawing to texture

		raylib.BeginMode3D(camera)

		raylib.DrawModel(dwarf, position, 2.0, raylib.White) // Draw 3d model with texture

		raylib.DrawGrid(10, 1.0) // Draw a grid

		raylib.EndMode3D()

		raylib.EndTextureMode() // End drawing to texture (now we have a texture available for next passes)

		// Render previously generated texture using selected postpro shader
		raylib.BeginShaderMode(shaders[currentShader])

		// NOTE: Render texture must be y-flipped due to default OpenGL coordinates (left-bottom)
		raylib.DrawTextureRec(target.Texture, raylib.NewRectangle(0, 0, float32(target.Texture.Width), float32(-target.Texture.Height)), raylib.NewVector2(0, 0), raylib.White)

		raylib.EndShaderMode()

		raylib.DrawRectangle(0, 9, 580, 30, raylib.Fade(raylib.LightGray, 0.7))

		raylib.DrawText("(c) Dwarf 3D model by David Moreno", screenWidth-200, screenHeight-20, 10, raylib.DarkGray)

		raylib.DrawText("CURRENT POSTPRO SHADER:", 10, 15, 20, raylib.Black)
		raylib.DrawText(postproShaderText[currentShader], 330, 15, 20, raylib.Red)
		raylib.DrawText("< >", 540, 10, 30, raylib.DarkBlue)

		raylib.DrawFPS(700, 15)

		raylib.EndDrawing()
	}

	// Unload all postpro shaders
	for i := 0; i < MaxPostproShaders; i++ {
		raylib.UnloadShader(shaders[i])
	}

	raylib.UnloadTexture(texture)      // Unload texture
	raylib.UnloadModel(dwarf)          // Unload model
	raylib.UnloadRenderTexture(target) // Unload render texture

	raylib.CloseWindow()
}
