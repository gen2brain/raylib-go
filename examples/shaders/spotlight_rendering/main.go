package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	maxSpots = 3 // NOTE: It must be the same as defined in the shader
	maxStars = 400

	screenWidth  = 800
	screenHeight = 450
)

// Spot data
type Spot struct {
	Position    rl.Vector2
	Speed       rl.Vector2
	Inner       float32
	Radius      float32
	PositionLoc int32
	InnerLoc    int32
	RadiusLoc   int32
}

// Stars in the star field
type Star struct {
	Position rl.Vector2
	Speed    rl.Vector2
}

func main() {
	// Initialization
	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - spotlight rendering")

	rl.HideCursor()

	// Load texture directly from root folder
	texRay := rl.LoadTexture("raysan.png")

	var stars [maxStars]Star
	for n := 0; n < maxStars; n++ {
		resetStar(&stars[n])
	}

	// Progress all the stars on, so they don't all start in the centre
	for m := 0; m < screenWidth/2; m++ {
		for n := 0; n < maxStars; n++ {
			updateStar(&stars[n])
		}
	}

	frameCounter := 0

	// Load spotlight shader directly from root folder (use default vertex shader)
	shdrSpot := rl.LoadShader("", "spotlight.fs")

	// Get the locations of spots in the shader
	var spots [maxSpots]Spot

	for i := 0; i < maxSpots; i++ {
		posName := fmt.Sprintf("spots[%d].pos", i)
		innerName := fmt.Sprintf("spots[%d].inner", i)
		radiusName := fmt.Sprintf("spots[%d].radius", i)

		spots[i].PositionLoc = rl.GetShaderLocation(shdrSpot, posName)
		spots[i].InnerLoc = rl.GetShaderLocation(shdrSpot, innerName)
		spots[i].RadiusLoc = rl.GetShaderLocation(shdrSpot, radiusName)
	}

	// Tell the shader how wide the screen is
	wLoc := rl.GetShaderLocation(shdrSpot, "screenWidth")
	sw := float32(rl.GetScreenWidth())
	rl.SetShaderValue(shdrSpot, wLoc, []float32{sw}, rl.ShaderUniformFloat)

	// Randomize the locations and velocities of the spotlights
	for i := 0; i < maxSpots; i++ {
		spots[i].Position.X = float32(rl.GetRandomValue(64, screenWidth-64))
		spots[i].Position.Y = float32(rl.GetRandomValue(64, screenHeight-64))
		spots[i].Speed = rl.NewVector2(0, 0)

		for (math.Abs(float64(spots[i].Speed.X)) + math.Abs(float64(spots[i].Speed.Y))) < 2 {
			spots[i].Speed.X = float32(rl.GetRandomValue(-400, 40)) / 25.0
			spots[i].Speed.Y = float32(rl.GetRandomValue(-400, 40)) / 25.0
		}

		spots[i].Inner = 28.0 * float32(i+1)
		spots[i].Radius = 48.0 * float32(i+1)

		rl.SetShaderValue(shdrSpot, spots[i].PositionLoc, []float32{spots[i].Position.X, spots[i].Position.Y}, rl.ShaderUniformVec2)
		rl.SetShaderValue(shdrSpot, spots[i].InnerLoc, []float32{spots[i].Inner}, rl.ShaderUniformFloat)
		rl.SetShaderValue(shdrSpot, spots[i].RadiusLoc, []float32{spots[i].Radius}, rl.ShaderUniformFloat)
	}

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		frameCounter++

		// Move the stars, resetting them if they go offscreen
		for n := 0; n < maxStars; n++ {
			updateStar(&stars[n])
		}

		// Update the spots, send them to the shader
		for i := 0; i < maxSpots; i++ {
			if i == 0 {
				mp := rl.GetMousePosition()
				spots[i].Position.X = mp.X
				spots[i].Position.Y = float32(screenHeight) - mp.Y
			} else {
				spots[i].Position.X += spots[i].Speed.X
				spots[i].Position.Y += spots[i].Speed.Y

				if spots[i].Position.X < 64 {
					spots[i].Speed.X = -spots[i].Speed.X
				}
				if spots[i].Position.X > float32(screenWidth-64) {
					spots[i].Speed.X = -spots[i].Speed.X
				}
				if spots[i].Position.Y < 64 {
					spots[i].Speed.Y = -spots[i].Speed.Y
				}
				if spots[i].Position.Y > float32(screenHeight-64) {
					spots[i].Speed.Y = -spots[i].Speed.Y
				}
			}

			rl.SetShaderValue(shdrSpot, spots[i].PositionLoc, []float32{spots[i].Position.X, spots[i].Position.Y}, rl.ShaderUniformVec2)
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkBlue)

		// Draw stars
		for n := 0; n < maxStars; n++ {
			rl.DrawRectangle(int32(stars[n].Position.X), int32(stars[n].Position.Y), 2, 2, rl.White)
		}

		// Draw animated textures
		for i := 0; i < 16; i++ {
			posX := int32((float64(screenWidth) / 2.0) + math.Cos(float64(frameCounter+i*8)/51.45)*(float64(screenWidth)/2.2) - 32)
			posY := int32((float64(screenHeight) / 2.0) + math.Sin(float64(frameCounter+i*8)/17.87)*(float64(screenHeight)/4.2))
			rl.DrawTexture(texRay, posX, posY, rl.White)
		}

		// Draw spot lights overlay
		rl.BeginShaderMode(shdrSpot)
		rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.White)
		rl.EndShaderMode()

		rl.DrawFPS(10, 10)

		rl.DrawText("Move the mouse!", 10, 30, 20, rl.Green)
		rl.DrawText("Pitch Black", int32(screenWidth*0.2), screenHeight/2, 20, rl.Green)
		rl.DrawText("Dark", int32(screenWidth*0.66), screenHeight/2, 20, rl.Green)

		rl.EndDrawing()
	}
	rl.UnloadShader(shdrSpot)
	rl.UnloadTexture(texRay)
	rl.CloseWindow()
}

// Reset star to screen center with random velocity
func resetStar(star *Star) {
	star.Position = rl.NewVector2(float32(rl.GetScreenWidth())/2.0, float32(rl.GetScreenHeight())/2.0)

	star.Speed.X = float32(rl.GetRandomValue(-1000, 1000)) / 100.0
	star.Speed.Y = float32(rl.GetRandomValue(-1000, 1000)) / 100.0

	for (math.Abs(float64(star.Speed.X)) + math.Abs(float64(star.Speed.Y))) <= 1.0 {
		star.Speed.X = float32(rl.GetRandomValue(-1000, 1000)) / 100.0
		star.Speed.Y = float32(rl.GetRandomValue(-1000, 1000)) / 100.0
	}

	star.Position = rl.Vector2Add(star.Position, rl.Vector2Multiply(star.Speed, rl.NewVector2(8.0, 8.0)))
}

// Update star position based on speed
func updateStar(star *Star) {
	star.Position = rl.Vector2Add(star.Position, star.Speed)

	if star.Position.X < 0 || star.Position.X > float32(rl.GetScreenWidth()) ||
		star.Position.Y < 0 || star.Position.Y > float32(rl.GetScreenHeight()) {
		resetStar(star)
	}
}
