package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const (
	maxParticles = 200
)

type particle struct {
	Position raylib.Vector2
	Color    raylib.Color
	Alpha    float32
	Size     float32
	Rotation float32
	Active   bool
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	//raylib.SetConfigFlags(raylib.FlagVsyncHint)
	raylib.InitWindow(screenWidth, screenHeight, "raylib [textures] example - particles blending")

	// Particles pool, reuse them!
	mouseTail := make([]particle, maxParticles)

	// Initialize particles
	for i := 0; i < maxParticles; i++ {
		mouseTail[i].Position = raylib.NewVector2(0, 0)
		mouseTail[i].Color = raylib.NewColor(byte(raylib.GetRandomValue(0, 255)), byte(raylib.GetRandomValue(0, 255)), byte(raylib.GetRandomValue(0, 255)), 255)
		mouseTail[i].Alpha = 1.0
		mouseTail[i].Size = float32(raylib.GetRandomValue(1, 30)) / 20.0
		mouseTail[i].Rotation = float32(raylib.GetRandomValue(0, 360))
		mouseTail[i].Active = false
	}

	gravity := float32(3.0)

	smoke := raylib.LoadTexture("smoke.png")

	blending := raylib.BlendAlpha

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update

		// Activate one particle every frame and Update active particles
		// NOTE: Particles initial position should be mouse position when activated
		// NOTE: Particles fall down with gravity and rotation... and disappear after 2 seconds (alpha = 0)
		// NOTE: When a particle disappears, active = false and it can be reused.
		for i := 0; i < maxParticles; i++ {
			if !mouseTail[i].Active {
				mouseTail[i].Active = true
				mouseTail[i].Alpha = 1.0
				mouseTail[i].Position = raylib.GetMousePosition()
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

		if raylib.IsKeyPressed(raylib.KeySpace) {
			if blending == raylib.BlendAlpha {
				blending = raylib.BlendAdditive
			} else {
				blending = raylib.BlendAlpha
			}
		}

		// Draw

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.DarkGray)

		raylib.BeginBlendMode(blending)

		// Draw active particles
		for i := 0; i < maxParticles; i++ {
			if mouseTail[i].Active {
				raylib.DrawTexturePro(
					smoke,
					raylib.NewRectangle(0, 0, float32(smoke.Width), float32(smoke.Height)),
					raylib.NewRectangle(mouseTail[i].Position.X, mouseTail[i].Position.Y, float32(smoke.Width)*mouseTail[i].Size, float32(smoke.Height)*mouseTail[i].Size),
					raylib.NewVector2(float32(smoke.Width)*mouseTail[i].Size/2, float32(smoke.Height)*mouseTail[i].Size/2),
					mouseTail[i].Rotation,
					raylib.Fade(mouseTail[i].Color, mouseTail[i].Alpha),
				)
			}
		}

		raylib.EndBlendMode()

		raylib.DrawText("PRESS SPACE to CHANGE BLENDING MODE", 180, 20, 20, raylib.Black)

		if blending == raylib.BlendAlpha {
			raylib.DrawText("ALPHA BLENDING", 290, screenHeight-40, 20, raylib.Black)
		} else {
			raylib.DrawText("ADDITIVE BLENDING", 280, screenHeight-40, 20, raylib.RayWhite)
		}

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(smoke)

	raylib.CloseWindow()
}
