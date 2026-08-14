package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Custom Blend Mode Constants
const (
	rlglSrcAlpha int32 = 0x0302
	rlglMin      int32 = 0x8007
	rlglMax      int32 = 0x8008

	maxBoxes   = 20
	maxShadows = maxBoxes * 3 // Each box can cast up to 2 shadow volumes for edges + 1 for box itself
	maxLights  = 16
)

// Shadow geometry type
type ShadowGeometry struct {
	Vertices [4]rl.Vector2
}

// Light info type
type LightInfo struct {
	Active      bool
	Dirty       bool
	Valid       bool
	Position    rl.Vector2
	Mask        rl.RenderTexture2D
	OuterRadius float32
	Bounds      rl.Rectangle
	Shadows     [maxShadows]ShadowGeometry
	ShadowCount int
}

// Global lights state
var lights [maxLights]LightInfo

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - top down lights")
	defer rl.CloseWindow()

	// Initialize our 'world' of boxes
	boxes := setupBoxes()

	// Create a checkerboard ground texture
	img := rl.GenImageChecked(64, 64, 32, 32, rl.DarkBrown, rl.DarkGray)
	backgroundTexture := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	defer rl.UnloadTexture(backgroundTexture)

	// Create a global light mask to hold all the blended lights
	lightMask := rl.LoadRenderTexture(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()))
	defer rl.UnloadRenderTexture(lightMask)

	// Cleanup per-light render textures at exit
	defer func() {
		for i := 0; i < maxLights; i++ {
			if lights[i].Active {
				rl.UnloadRenderTexture(lights[i].Mask)
			}
		}
	}()

	// Setup initial light
	setupLight(0, 600, 400, 300)
	nextLight := 1

	showLines := false

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		// Drag light 0
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			mousePos := rl.GetMousePosition()
			moveLight(0, mousePos.X, mousePos.Y)
		}

		// Make a new light
		if rl.IsMouseButtonPressed(rl.MouseButtonRight) && nextLight < maxLights {
			mousePos := rl.GetMousePosition()
			setupLight(nextLight, mousePos.X, mousePos.Y, 200)
			nextLight++
		}

		// Toggle debug info
		if rl.IsKeyPressed(rl.KeyF1) {
			showLines = !showLines
		}

		// Update the lights and keep track if any were dirty
		dirtyLights := false
		for i := 0; i < maxLights; i++ {
			if updateLight(i, boxes) {
				dirtyLights = true
			}
		}

		// Update the light mask
		if dirtyLights {
			rl.BeginTextureMode(lightMask)

			rl.ClearBackground(rl.Black)

			// Force blend mode to only set the alpha of the destination
			rl.SetBlendFactors(rlglSrcAlpha, rlglSrcAlpha, rlglMin)
			rl.SetBlendMode(rl.BlendCustom)

			// Merge in all the light masks
			srcRec := rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), -float32(rl.GetScreenHeight()))
			for i := 0; i < maxLights; i++ {
				if lights[i].Active {
					rl.DrawTextureRec(lights[i].Mask.Texture, srcRec, rl.Vector2Zero(), rl.White)
				}
			}

			rl.DrawRenderBatchActive()

			// Go back to normal blend
			rl.SetBlendMode(rl.BlendAlpha)

			rl.EndTextureMode()
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		// Draw tile background
		bgRec := rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()))
		rl.DrawTextureRec(backgroundTexture, bgRec, rl.Vector2Zero(), rl.White)

		// Overlay shadows from all lights
		maskSrcRec := rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), -float32(rl.GetScreenHeight()))
		alpha := float32(1.0)
		if showLines {
			alpha = 0.75
		}
		rl.DrawTextureRec(lightMask.Texture, maskSrcRec, rl.Vector2Zero(), rl.ColorAlpha(rl.White, alpha))

		// Draw light positions
		for i := 0; i < maxLights; i++ {
			if lights[i].Active {
				lightColor := rl.White
				if i == 0 {
					lightColor = rl.Yellow
				}
				rl.DrawCircle(int32(lights[i].Position.X), int32(lights[i].Position.Y), 10, lightColor)
			}
		}

		if showLines {
			for s := 0; s < lights[0].ShadowCount; s++ {
				rl.DrawTriangleFan(lights[0].Shadows[s].Vertices[:], rl.DarkPurple)
			}

			for b := 0; b < len(boxes); b++ {
				if rl.CheckCollisionRecs(boxes[b], lights[0].Bounds) {
					rl.DrawRectangleRec(boxes[b], rl.Purple)
				}

				rl.DrawRectangleLines(
					int32(boxes[b].X),
					int32(boxes[b].Y),
					int32(boxes[b].Width),
					int32(boxes[b].Height),
					rl.DarkBlue,
				)
			}

			rl.DrawText("(F1) Hide Shadow Volumes", 10, 50, 10, rl.Green)
		} else {
			rl.DrawText("(F1) Show Shadow Volumes", 10, 50, 10, rl.Green)
		}

		rl.DrawFPS(screenWidth-80, 10)
		rl.DrawText("Drag to move light #1", 10, 10, 10, rl.DarkGreen)
		rl.DrawText("Right click to add new light", 10, 30, 10, rl.DarkGreen)

		rl.EndDrawing()
	}
}

// Move a light and mark it as dirty so that we update its mask next frame
func moveLight(slot int, x, y float32) {
	lights[slot].Dirty = true
	lights[slot].Position.X = x
	lights[slot].Position.Y = y

	lights[slot].Bounds.X = x - lights[slot].OuterRadius
	lights[slot].Bounds.Y = y - lights[slot].OuterRadius
}

// Compute a shadow volume for the edge
func computeShadowVolumeForEdge(slot int, sp, ep rl.Vector2) {
	if lights[slot].ShadowCount >= maxShadows {
		return
	}

	extension := lights[slot].OuterRadius * 2.0

	spVector := rl.Vector2Normalize(rl.Vector2Subtract(sp, lights[slot].Position))
	spProjection := rl.Vector2Add(sp, rl.Vector2Scale(spVector, extension))

	epVector := rl.Vector2Normalize(rl.Vector2Subtract(ep, lights[slot].Position))
	epProjection := rl.Vector2Add(ep, rl.Vector2Scale(epVector, extension))

	idx := lights[slot].ShadowCount
	lights[slot].Shadows[idx].Vertices[0] = sp
	lights[slot].Shadows[idx].Vertices[1] = ep
	lights[slot].Shadows[idx].Vertices[2] = epProjection
	lights[slot].Shadows[idx].Vertices[3] = spProjection

	lights[slot].ShadowCount++
}

// Setup a light
func setupLight(slot int, x, y, radius float32) {
	lights[slot].Active = true
	lights[slot].Valid = false
	lights[slot].Mask = rl.LoadRenderTexture(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()))
	lights[slot].OuterRadius = radius

	lights[slot].Bounds.Width = radius * 2.0
	lights[slot].Bounds.Height = radius * 2.0

	moveLight(slot, x, y)

	drawLightMask(slot)
}

// See if a light needs to update its mask
func updateLight(slot int, boxes []rl.Rectangle) bool {
	if !lights[slot].Active || !lights[slot].Dirty {
		return false
	}

	lights[slot].Dirty = false
	lights[slot].ShadowCount = 0
	lights[slot].Valid = false

	for i := 0; i < len(boxes); i++ {
		if rl.CheckCollisionPointRec(lights[slot].Position, boxes[i]) {
			return false
		}

		if !rl.CheckCollisionRecs(lights[slot].Bounds, boxes[i]) {
			continue
		}

		// Top
		sp := rl.NewVector2(boxes[i].X, boxes[i].Y)
		ep := rl.NewVector2(boxes[i].X+boxes[i].Width, boxes[i].Y)

		if lights[slot].Position.Y > ep.Y {
			computeShadowVolumeForEdge(slot, sp, ep)
		}

		// Right
		sp = ep
		ep.Y += boxes[i].Height
		if lights[slot].Position.X < ep.X {
			computeShadowVolumeForEdge(slot, sp, ep)
		}

		// Bottom
		sp = ep
		ep.X -= boxes[i].Width
		if lights[slot].Position.Y < ep.Y {
			computeShadowVolumeForEdge(slot, sp, ep)
		}

		// Left
		sp = ep
		ep.Y -= boxes[i].Height
		if lights[slot].Position.X > ep.X {
			computeShadowVolumeForEdge(slot, sp, ep)
		}

		// The box itself
		idx := lights[slot].ShadowCount
		lights[slot].Shadows[idx].Vertices[0] = rl.NewVector2(boxes[i].X, boxes[i].Y)
		lights[slot].Shadows[idx].Vertices[1] = rl.NewVector2(boxes[i].X, boxes[i].Y+boxes[i].Height)
		lights[slot].Shadows[idx].Vertices[2] = rl.NewVector2(boxes[i].X+boxes[i].Width, boxes[i].Y+boxes[i].Height)
		lights[slot].Shadows[idx].Vertices[3] = rl.NewVector2(boxes[i].X+boxes[i].Width, boxes[i].Y)
		lights[slot].ShadowCount++
	}

	lights[slot].Valid = true

	drawLightMask(slot)

	return true
}

// Draw the light and shadows to the mask for a light
func drawLightMask(slot int) {
	rl.BeginTextureMode(lights[slot].Mask)

	rl.ClearBackground(rl.White)

	// Force blend mode to only set alpha of destination
	rl.SetBlendFactors(rlglSrcAlpha, rlglSrcAlpha, rlglMin)
	rl.SetBlendMode(rl.BlendCustom)

	if lights[slot].Valid {
		rl.DrawCircleGradient(
			lights[slot].Position,
			lights[slot].OuterRadius,
			rl.ColorAlpha(rl.White, 0),
			rl.White,
		)
	}

	rl.DrawRenderBatchActive()

	// Cut out shadows from light radius by forcing alpha to max
	rl.SetBlendMode(rl.BlendAlpha)
	rl.SetBlendFactors(rlglSrcAlpha, rlglSrcAlpha, rlglMax)
	rl.SetBlendMode(rl.BlendCustom)

	for i := 0; i < lights[slot].ShadowCount; i++ {
		rl.DrawTriangleFan(lights[slot].Shadows[i].Vertices[:], rl.White)
	}

	rl.DrawRenderBatchActive()

	rl.SetBlendMode(rl.BlendAlpha)

	rl.EndTextureMode()
}

// Set up some boxes
func setupBoxes() []rl.Rectangle {
	boxes := make([]rl.Rectangle, maxBoxes)
	boxes[0] = rl.NewRectangle(150, 80, 40, 40)
	boxes[1] = rl.NewRectangle(1200, 700, 40, 40)
	boxes[2] = rl.NewRectangle(200, 600, 40, 40)
	boxes[3] = rl.NewRectangle(1000, 50, 40, 40)
	boxes[4] = rl.NewRectangle(500, 350, 40, 40)

	for i := 5; i < maxBoxes; i++ {
		boxes[i] = rl.NewRectangle(
			float32(rl.GetRandomValue(0, int32(rl.GetScreenWidth()))),
			float32(rl.GetRandomValue(0, int32(rl.GetScreenHeight()))),
			float32(rl.GetRandomValue(10, 100)),
			float32(rl.GetRandomValue(10, 100)),
		)
	}

	return boxes
}
