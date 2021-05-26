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

	rl.SetConfigFlags(rl.FlagMsaa4xHint) // Enable Multi Sampling Anti Aliasing 4x (if available)

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - postprocessing shader")

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(2.0, 3.0, 2.0)
	camera.Target = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0

	obj := rl.LoadModel("church.obj")                            // Load OBJ model
	texture := rl.LoadTexture("church_diffuse.png")              // Load model texture
	rl.SetMaterialTexture(obj.Materials, rl.MapDiffuse, texture) // Set obj model diffuse texture

	position := rl.NewVector3(0.0, 0.0, 0.0) // Set model position

	// Load all postpro shaders
	// NOTE 1: All postpro shader use the base vertex shader (DEFAULT_VERTEX_SHADER)
	shaders := make([]rl.Shader, MaxPostproShaders)
	shaders[FxGrayscale] = rl.LoadShader("", "glsl330/grayscale.fs")
	shaders[FxPosterization] = rl.LoadShader("", "glsl330/posterization.fs")
	shaders[FxDreamVision] = rl.LoadShader("", "glsl330/dream_vision.fs")
	shaders[FxPixelizer] = rl.LoadShader("", "glsl330/pixelizer.fs")
	shaders[FxCrossHatching] = rl.LoadShader("", "glsl330/cross_hatching.fs")
	shaders[FxCrossStitching] = rl.LoadShader("", "glsl330/cross_stitching.fs")
	shaders[FxPredatorView] = rl.LoadShader("", "glsl330/predator.fs")
	shaders[FxScanlines] = rl.LoadShader("", "glsl330/scanlines.fs")
	shaders[FxFisheye] = rl.LoadShader("", "glsl330/fisheye.fs")
	shaders[FxSobel] = rl.LoadShader("", "glsl330/sobel.fs")
	shaders[FxBlur] = rl.LoadShader("", "glsl330/blur.fs")
	shaders[FxBloom] = rl.LoadShader("", "glsl330/bloom.fs")

	currentShader := FxGrayscale

	// Create a RenderTexture2D to be used for render to texture
	target := rl.LoadRenderTexture(screenWidth, screenHeight)

	rl.SetCameraMode(camera, rl.CameraOrbital) // Set free camera mode

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera) // Update camera

		if rl.IsKeyPressed(rl.KeyRight) {
			currentShader++
		} else if rl.IsKeyPressed(rl.KeyLeft) {
			currentShader--
		}

		if currentShader >= MaxPostproShaders {
			currentShader = 0
		} else if currentShader < 0 {
			currentShader = MaxPostproShaders - 1
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginTextureMode(target) // Enable drawing to texture

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(obj, position, 0.1, rl.White) // Draw 3d model with texture

		rl.DrawGrid(10, 1.0) // Draw a grid

		rl.EndMode3D()

		rl.EndTextureMode() // End drawing to texture (now we have a texture available for next passes)

		// Render previously generated texture using selected postpro shader
		rl.BeginShaderMode(shaders[currentShader])

		// NOTE: Render texture must be y-flipped due to default OpenGL coordinates (left-bottom)
		rl.DrawTextureRec(target.Texture, rl.NewRectangle(0, 0, float32(target.Texture.Width), float32(-target.Texture.Height)), rl.NewVector2(0, 0), rl.White)

		rl.EndShaderMode()

		rl.DrawRectangle(0, 9, 580, 30, rl.Fade(rl.LightGray, 0.7))

		rl.DrawText("(c) Church 3D model by Alberto Cano", screenWidth-200, screenHeight-20, 10, rl.DarkGray)

		rl.DrawText("CURRENT POSTPRO SHADER:", 10, 15, 20, rl.Black)
		rl.DrawText(postproShaderText[currentShader], 330, 15, 20, rl.Red)
		rl.DrawText("< >", 540, 10, 30, rl.DarkBlue)

		rl.DrawFPS(700, 15)

		rl.EndDrawing()
	}

	// Unload all postpro shaders
	for i := 0; i < MaxPostproShaders; i++ {
		rl.UnloadShader(shaders[i])
	}

	rl.UnloadTexture(texture)      // Unload texture
	rl.UnloadModel(obj)            // Unload model
	rl.UnloadRenderTexture(target) // Unload render texture

	rl.CloseWindow()
}
