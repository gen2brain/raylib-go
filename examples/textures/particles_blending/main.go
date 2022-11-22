package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const (
	maxParticles = 200
)

type particle struct {
	Position rl.Vector2
	Color    rl.Color
	Alpha    float32
	Size     float32
	Rotation float32
	Active   bool
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	//rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - particles blending")

	// Particles pool, reuse them!
	mouseTail := make([]particle, maxParticles)

	// Initialize particles
	for i := 0; i < maxParticles; i++ {
		mouseTail[i].Position = rl.NewVector2(0, 0)
		mouseTail[i].Color = rl.NewColor(byte(rl.GetRandomValue(0, 255)), byte(rl.GetRandomValue(0, 255)), byte(rl.GetRandomValue(0, 255)), 255)
		mouseTail[i].Alpha = 1.0
		mouseTail[i].Size = float32(rl.GetRandomValue(1, 30)) / 20.0
		mouseTail[i].Rotation = float32(rl.GetRandomValue(0, 360))
		mouseTail[i].Active = false
	}

	gravity := float32(3.0)

	smoke := rl.LoadTexture("smoke.png")

	blending := rl.BlendAlpha

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update

		// Activate one particle every frame and Update active particles
		// NOTE: Particles initial position should be mouse position when activated
		// NOTE: Particles fall down with gravity and rotation... and disappear after 2 seconds (alpha = 0)
		// NOTE: When a particle disappears, active = false and it can be reused.
		for i := 0; i < maxParticles; i++ {
			if !mouseTail[i].Active {
				mouseTail[i].Active = true
				mouseTail[i].Alpha = 1.0
				mouseTail[i].Position = rl.GetMousePosition()
				i = maxParticles
			}
		}

		for i := 0; i < maxParticles; i++ {
			if mouseTail[i].Active {
				mouseTail[i].Position.Y += gravity
				mouseTail[i].Alpha -= 0.01

				if mouseTail[i].Alpha <= 0.0 {
					mouseTail[i].Active = false
				}

				mouseTail[i].Rotation += 5.0
			}
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			if blending == rl.BlendAlpha {
				blending = rl.BlendAdditive
			} else {
				blending = rl.BlendAlpha
			}
		}

		// Draw

		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkGray)

		rl.BeginBlendMode(blending)

		// Draw active particles
		for i := 0; i < maxParticles; i++ {
			if mouseTail[i].Active {
				rl.DrawTexturePro(
					smoke,
					rl.NewRectangle(0, 0, float32(smoke.Width), float32(smoke.Height)),
					rl.NewRectangle(mouseTail[i].Position.X, mouseTail[i].Position.Y, float32(smoke.Width)*mouseTail[i].Size, float32(smoke.Height)*mouseTail[i].Size),
					rl.NewVector2(float32(smoke.Width)*mouseTail[i].Size/2, float32(smoke.Height)*mouseTail[i].Size/2),
					mouseTail[i].Rotation,
					rl.Fade(mouseTail[i].Color, mouseTail[i].Alpha),
				)
			}
		}

		rl.EndBlendMode()

		rl.DrawText("PRESS SPACE to CHANGE BLENDING MODE", 180, 20, 20, rl.Black)

		if blending == rl.BlendAlpha {
			rl.DrawText("ALPHA BLENDING", 290, screenHeight-40, 20, rl.Black)
		} else {
			rl.DrawText("ADDITIVE BLENDING", 280, screenHeight-40, 20, rl.RayWhite)
		}

		rl.EndDrawing()
	}

	rl.UnloadTexture(smoke)

	rl.CloseWindow()
}
