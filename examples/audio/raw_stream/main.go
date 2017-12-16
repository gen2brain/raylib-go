package main

import (
	"math"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	maxSamples          = 22050
	maxSamplesPerUpdate = 4096
)

func main() {
	raylib.InitWindow(800, 450, "raylib [audio] example - raw audio streaming")

	raylib.InitAudioDevice()

	// Init raw audio stream (sample rate: 22050, sample size: 32bit-float, channels: 1-mono)
	stream := raylib.InitAudioStream(22050, 32, 1)

	//// Fill audio stream with some samples (sine wave)
	data := make([]float32, maxSamples)

	for i := 0; i < maxSamples; i++ {
		data[i] = float32(math.Sin(float64((2*raylib.Pi*float32(i))/2) * raylib.Deg2rad))
	}

	// NOTE: The generated MAX_SAMPLES do not fit to close a perfect loop
	// for that reason, there is a clip everytime audio stream is looped
	raylib.PlayAudioStream(stream)

	totalSamples := int32(maxSamples)
	samplesLeft := int32(totalSamples)

	position := raylib.NewVector2(0, 0)

	raylib.SetTargetFPS(30)

	for !raylib.WindowShouldClose() {
		// Refill audio stream if required
		if raylib.IsAudioBufferProcessed(stream) {
			numSamples := int32(0)
			if samplesLeft >= maxSamplesPerUpdate {
				numSamples = maxSamplesPerUpdate
			} else {
				numSamples = samplesLeft
			}

			raylib.UpdateAudioStream(stream, data[totalSamples-samplesLeft:], numSamples)

			samplesLeft -= numSamples

			// Reset samples feeding (loop audio)
			if samplesLeft <= 0 {
				samplesLeft = totalSamples
			}
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)
		raylib.DrawText("SINE WAVE SHOULD BE PLAYING!", 240, 140, 20, raylib.LightGray)

		// NOTE: Draw a part of the sine wave (only screen width)
		for i := 0; i < int(raylib.GetScreenWidth()); i++ {
			position.X = float32(i)
			position.Y = 250 + 50*data[i]

			raylib.DrawPixelV(position, raylib.Red)
		}

		raylib.EndDrawing()
	}

	raylib.CloseAudioStream(stream) // Close raw audio stream and delete buffers from RAM

	raylib.CloseAudioDevice() // Close audio device (music streaming is automatically stopped)

	raylib.CloseWindow()
}
