package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const maxTextures = 4

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - color correction")

	textures := [maxTextures]rl.Texture2D{
		rl.LoadTexture("parrots.png"),
		rl.LoadTexture("cat.png"),
		rl.LoadTexture("mandrill.png"),
		rl.LoadTexture("fudesumi.png"),
	}

	shdrColorCorrection := rl.LoadShader("", "color_correction.fs")

	imageIndex := int32(0)
	resetButtonClicked := false

	contrast := float32(0.0)
	saturation := float32(0.0)
	brightness := float32(0.0)

	// Get shader locations
	contrastLoc := rl.GetShaderLocation(shdrColorCorrection, "contrast")
	saturationLoc := rl.GetShaderLocation(shdrColorCorrection, "saturation")
	brightnessLoc := rl.GetShaderLocation(shdrColorCorrection, "brightness")

	// Set initial shader values
	rl.SetShaderValue(shdrColorCorrection, contrastLoc, []float32{contrast}, rl.ShaderUniformFloat)
	rl.SetShaderValue(shdrColorCorrection, saturationLoc, []float32{saturation}, rl.ShaderUniformFloat)
	rl.SetShaderValue(shdrColorCorrection, brightnessLoc, []float32{brightness}, rl.ShaderUniformFloat)

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		// Select texture to draw via keys
		if rl.IsKeyPressed(rl.KeyOne) {
			imageIndex = 0
		} else if rl.IsKeyPressed(rl.KeyTwo) {
			imageIndex = 1
		} else if rl.IsKeyPressed(rl.KeyThree) {
			imageIndex = 2
		} else if rl.IsKeyPressed(rl.KeyFour) {
			imageIndex = 3
		}

		// Reset values to 0
		if rl.IsKeyPressed(rl.KeyR) || resetButtonClicked {
			contrast = 0.0
			saturation = 0.0
			brightness = 0.0
		}

		// Send values to shader
		rl.SetShaderValue(shdrColorCorrection, contrastLoc, []float32{contrast}, rl.ShaderUniformFloat)
		rl.SetShaderValue(shdrColorCorrection, saturationLoc, []float32{saturation}, rl.ShaderUniformFloat)
		rl.SetShaderValue(shdrColorCorrection, brightnessLoc, []float32{brightness}, rl.ShaderUniformFloat)

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginShaderMode(shdrColorCorrection)

		currTex := textures[imageIndex]
		posX := int32(580/2) - currTex.Width/2
		posY := int32(rl.GetScreenHeight()/2) - currTex.Height/2
		rl.DrawTexture(currTex, posX, posY, rl.White)

		rl.EndShaderMode()

		rl.DrawLine(580, 0, 580, int32(rl.GetScreenHeight()), rl.NewColor(218, 218, 218, 255))
		rl.DrawRectangle(580, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.NewColor(232, 232, 232, 255))

		// Draw UI info text
		rl.DrawText("Color Correction", 585, 40, 20, rl.Gray)

		rl.DrawText("Picture", 602, 75, 10, rl.Gray)
		rl.DrawText("Press [1] - [4] to Change Picture", 600, 230, 8, rl.Gray)
		rl.DrawText("Press [R] to Reset Values", 600, 250, 8, rl.Gray)

		// Draw GUI controls using pure Raylib helpers
		drawGuiToggleGroup(rl.NewRectangle(645, 70, 20, 20), []string{"1", "2", "3", "4"}, &imageIndex)

		drawGuiSliderBar(rl.NewRectangle(645, 100, 120, 20), "Contrast", fmt.Sprintf("%.0f", contrast), &contrast, -100.0, 100.0)
		drawGuiSliderBar(rl.NewRectangle(645, 130, 120, 20), "Saturation", fmt.Sprintf("%.0f", saturation), &saturation, -100.0, 100.0)
		drawGuiSliderBar(rl.NewRectangle(645, 160, 120, 20), "Brightness", fmt.Sprintf("%.0f", brightness), &brightness, -100.0, 100.0)

		resetButtonClicked = drawGuiButton(rl.NewRectangle(645, 190, 40, 20), "Reset")

		rl.DrawFPS(710, 10)

		rl.EndDrawing()
	}

	for i := 0; i < maxTextures; i++ {
		rl.UnloadTexture(textures[i])
	}
	rl.UnloadShader(shdrColorCorrection)

	rl.CloseWindow()
}

func drawGuiToggleGroup(bounds rl.Rectangle, options []string, active *int32) {
	mousePos := rl.GetMousePosition()
	for i, option := range options {
		rec := rl.NewRectangle(bounds.X+float32(i)*(bounds.Width+5), bounds.Y, bounds.Width, bounds.Height)
		isSelected := *active == int32(i)
		isHovered := rl.CheckCollisionPointRec(mousePos, rec)

		bgColor := rl.LightGray
		textColor := rl.DarkGray
		if isSelected {
			bgColor = rl.Blue
			textColor = rl.White
		} else if isHovered {
			bgColor = rl.Gray
			textColor = rl.White
		}

		rl.DrawRectangleRec(rec, bgColor)
		rl.DrawRectangleLinesEx(rec, 1, rl.DarkGray)
		rl.DrawText(option, int32(rec.X+rec.Width/2-3), int32(rec.Y+rec.Height/2-5), 10, textColor)

		if isHovered && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			*active = int32(i)
		}
	}
}

func drawGuiSliderBar(bounds rl.Rectangle, label, textRight string, value *float32, minValue, maxValue float32) {
	rl.DrawText(label, int32(bounds.X-60), int32(bounds.Y+5), 10, rl.DarkGray)

	rl.DrawRectangleRec(bounds, rl.LightGray)
	rl.DrawRectangleLinesEx(bounds, 1, rl.DarkGray)

	// Handle mouse dragging
	mousePos := rl.GetMousePosition()
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(mousePos, bounds) {
		pct := (mousePos.X - bounds.X) / bounds.Width
		if pct < 0 {
			pct = 0
		}
		if pct > 1 {
			pct = 1
		}
		*value = minValue + pct*(maxValue-minValue)
	}

	// Draw filled bar
	pct := (*value - minValue) / (maxValue - minValue)
	fillWidth := bounds.Width * pct
	if fillWidth > 0 {
		rl.DrawRectangleRec(rl.NewRectangle(bounds.X, bounds.Y, fillWidth, bounds.Height), rl.SkyBlue)
	}

	rl.DrawText(textRight, int32(bounds.X+bounds.Width+10), int32(bounds.Y+5), 10, rl.DarkGray)
}

func drawGuiButton(bounds rl.Rectangle, text string) bool {
	mousePos := rl.GetMousePosition()
	isHovered := rl.CheckCollisionPointRec(mousePos, bounds)

	bgColor := rl.LightGray
	textColor := rl.DarkGray
	if isHovered {
		bgColor = rl.Gray
		textColor = rl.White
	}

	rl.DrawRectangleRec(bounds, bgColor)
	rl.DrawRectangleLinesEx(bounds, 1, rl.DarkGray)
	rl.DrawText(text, int32(bounds.X+bounds.Width/2-13), int32(bounds.Y+bounds.Height/2-5), 10, textColor)

	return isHovered && rl.IsMouseButtonPressed(rl.MouseButtonLeft)
}
