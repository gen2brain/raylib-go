package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	maxSamples          = 22050
	maxSamplesPerUpdate = 4096
)

func main() {
	rl.InitWindow(800, 450, "raylib [audio] example - raw audio streaming")

	rl.InitAudioDevice()

	// Init raw audio stream (sample rate: 22050, sample size: 32bit-float, channels: 1-mono)
	stream := rl.LoadAudioStream(22050, 32, 1)

	//// Fill audio stream with some samples (sine wave)
	data := make([]float32, maxSamples)

	for i := 0; i < maxSamples; i++ {
		data[i] = float32(math.Sin(float64((2*rl.Pi*float32(i))/2) * rl.Deg2rad))
	}

	// NOTE: The generated MAX_SAMPLES do not fit to close a perfect loop
	// for that reason, there is a clip everytime audio stream is looped
	rl.PlayAudioStream(stream)

	totalSamples := int32(maxSamples)
	samplesLeft := int32(totalSamples)

	position := rl.NewVector2(0, 0)

	rl.SetTargetFPS(30)

	for !rl.WindowShouldClose() {
		// Refill audio stream if required
		if rl.IsAudioStreamProcessed(stream) {
			numSamples := int32(0)
			if samplesLeft >= maxSamplesPerUpdate {
				numSamples = maxSamplesPerUpdate
			} else {
				numSamples = samplesLeft
			}

			rl.UpdateAudioStream(stream, data[totalSamples-samplesLeft:], numSamples)

			samplesLeft -= numSamples

			// Reset samples feeding (loop audio)
			if samplesLeft <= 0 {
				samplesLeft = totalSamples
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("SINE WAVE SHOULD BE PLAYING!", 240, 140, 20, rl.LightGray)

		// NOTE: Draw a part of the sine wave (only screen width)
		for i := 0; i < int(rl.GetScreenWidth()); i++ {
			position.X = float32(i)
			position.Y = 250 + 50*data[i]

			rl.DrawPixelV(position, rl.Red)
		}

		rl.EndDrawing()
	}

	rl.UnloadAudioStream(stream) // Close raw audio stream and delete buffers from RAM

	rl.CloseAudioDevice() // Close audio device (music streaming is automatically stopped)

	rl.CloseWindow()
}
