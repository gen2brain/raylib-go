/*******************************************************************************************
*
*   raylib [audio] example - Mixed audio processing
*
*   Example originally created with raylib 4.2, last time updated with raylib 4.2
*
*   Example contributed by hkc (@hatkidchan) and reviewed by Ramon Santamaria (@raysan5)
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2023 hkc (@hatkidchan)
*
********************************************************************************************/
package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

var exponent float32 = 1.0     // Audio exponentiation value
var averageVolume [400]float32 // Average volume history

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [audio] example - processing mixed output")

	rl.InitAudioDevice() // Initialize audio device
	rl.AttachAudioMixedProcessor(ProcessAudio)

	music := rl.LoadMusicStream("country.mp3")
	sound := rl.LoadSound("coin.wav")

	rl.PlayMusicStream(music)

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		rl.UpdateMusicStream(music) // Update music buffer with new stream data

		// Modify processing variables
		if rl.IsKeyPressed(rl.KeyLeft) {
			exponent -= 0.05
		}
		if rl.IsKeyPressed(rl.KeyRight) {
			exponent += 0.05
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			exponent -= 0.25
		}
		if rl.IsKeyPressed(rl.KeyUp) {
			exponent += 0.25
		}

		// Make sure that exponent stays between 0.5 and 3
		exponent = clamp(exponent, 0.5, 3)

		if rl.IsKeyPressed(rl.KeySpace) {
			rl.PlaySound(sound)
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("MUSIC SHOULD BE PLAYING!", 255, 150, 20, rl.LightGray)

		rl.DrawText(fmt.Sprintf("EXPONENT = %.2f", exponent), 215, 180, 20, rl.LightGray)

		rl.DrawRectangle(199, 199, 402, 34, rl.LightGray)
		for i := int32(0); i < 400; i++ {
			rl.DrawLine(201+i, 232-int32(averageVolume[i]*32), 201+i, 232, rl.Maroon)
		}
		rl.DrawRectangleLines(199, 199, 402, 34, rl.Gray)

		rl.DrawText("PRESS SPACE TO PLAY OTHER SOUND", 200, 250, 20, rl.LightGray)
		rl.DrawText("USE LEFT AND RIGHT ARROWS TO ALTER DISTORTION", 140, 280, 20, rl.LightGray)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadMusicStream(music)                // Unload music stream buffers from RAM
	rl.DetachAudioMixedProcessor(ProcessAudio) // Disconnect audio processor
	rl.CloseAudioDevice()                      // Close audio device (music streaming is automatically stopped)

	rl.CloseWindow() // Close window and OpenGL context
}

// ProcessAudio is the audio processing function
func ProcessAudio(buffer []float32, frames int) {
	var average float32 // Temporary average volume
	maxFrame := frames / 2

	// Each frame has 2 samples (left and right),
	// so we should loop `frames / 2` times
	for frame := 0; frame < maxFrame; frame++ {
		left := &buffer[frame*2+0]  // Left channel
		right := &buffer[frame*2+1] // Right channel

		// Modify left and right channel samples with exponent
		*left = pow(abs(*left), exponent) * sgn(*left)
		*right = pow(abs(*right), exponent) * sgn(*right)

		// Accumulating average volume
		average += abs(*left) / float32(maxFrame)
		average += abs(*right) / float32(maxFrame)
	}

	// Shift average volume history buffer to the left
	for i := 0; i < 399; i++ {
		averageVolume[i] = averageVolume[i+1]
	}

	// Add the new average value
	averageVolume[399] = average
}

// Helper functions to make the code shorter
// (using less type conversions)
// (Golang: Please make the `math` package generic! This is ridiculous :-)
func abs(value float32) float32 {
	return float32(math.Abs(float64(value)))
}

func pow(value, exponent float32) float32 {
	return float32(math.Pow(float64(value), float64(exponent)))
}

func sgn(value float32) float32 {
	if value < 0 {
		return -1
	} else if value > 0 {
		return 1
	}
	return 0
}
func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
