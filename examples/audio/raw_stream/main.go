package main

import (
	"fmt"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	nSamples   = 44000 * 10
	sampleRate = 6000
	bufferSize = nSamples
	frequency  = 440
	targetFPS  = 240
)

func main() {
	rl.InitWindow(800, 450, "raylib [audio] example - raw audio streaming")
	position := rl.NewVector2(0, 0)

	rl.SetAudioStreamBufferSizeDefault(bufferSize)

	rl.InitAudioDevice()

	// Init raw audio stream (sample rate: <nSamples>, sample size: 32bit-float, channels: 1-mono)
	stream := rl.LoadAudioStream(nSamples, 32, 1)

	//// Create sine wave to play
	data := make([]float32, nSamples)

	for i := 0; i < nSamples; i++ {
		t := float32(i) / float32(nSamples)
		data[i] = float32(math.Sin(float64((2 * rl.Pi * frequency * t))))
	}

	// NOTE: The buffer can only be updated when it has been processed.  Time between buffer processing and next load and causes clipping
	rl.PlayAudioStream(stream)

	startTime := time.Now()
	rl.SetTargetFPS(targetFPS)

	for !rl.WindowShouldClose() {
		// Refill audio stream if buffer is processed
		if rl.IsAudioStreamProcessed(stream) {
			elapsedTime := time.Since(startTime).Seconds()
			currentSampleIndex := int(math.Mod(elapsedTime*float64(sampleRate), float64(nSamples)))
			nextSampleIndex := currentSampleIndex + bufferSize

			if nextSampleIndex > nSamples {
				nextSampleIndex = bufferSize - (nSamples - currentSampleIndex)
				rl.UpdateAudioStream(stream, append(data[currentSampleIndex:], data[:nextSampleIndex]...))
			} else {
				rl.UpdateAudioStream(stream, data[currentSampleIndex:nextSampleIndex])
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawText(fmt.Sprintf("%d Hz SINE WAVE SHOULD BE PLAYING!", frequency), 200, 140, 20, rl.LightGray)

		// NOTE: Draw a part of the sine wave (only screen width)
		for i := 0; i < int(rl.GetScreenWidth()); i++ {
			position.X = float32(i)
			position.Y = 250 + 10*data[i*nSamples/int(20*rl.GetScreenWidth())]

			rl.DrawPixelV(position, rl.Red)
		}

		rl.EndDrawing()
	}

	rl.UnloadAudioStream(stream) // Close raw audio stream and delete buffers from RAM

	rl.CloseAudioDevice() // Close audio device (music streaming is automatically stopped)

	rl.CloseWindow()
}
