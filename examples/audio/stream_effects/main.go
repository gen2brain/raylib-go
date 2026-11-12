/*******************************************************************************************
*
*   raylib [audio] example - stream effects
*
*   Example originally created with raylib 4.2, last time updated with raylib 5.0
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2022-2025 Ramon Santamaria (@raysan5)
*
********************************************************************************************/
package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

var (
	delayBuffer     []float32
	delayReadIndex  = 2
	delayWriteIndex = 0

	low [2]float32
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [audio] example - stream effects")

	rl.InitAudioDevice() // Initialize audio device

	music := rl.LoadMusicStream("country.mp3")

	// Allocate buffer for the delay effect
	delayBuffer = make([]float32, 48000*2) // 1 second delay (device sampleRate*channels)

	rl.PlayMusicStream(music)

	var timePlayed float32     // Time played normalized [0.0f..1.0f]
	pause := false             // Music playing paused
	enableEffectLPF := false   // Enable effect low-pass-filter
	enableEffectDelay := false // Enable effect delay (1 second)

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		rl.UpdateMusicStream(music) // Update music buffer with new stream data

		// Restart music playing (stop and play)
		if rl.IsKeyPressed(rl.KeySpace) {
			rl.StopMusicStream(music)
			rl.PlayMusicStream(music)
		}

		// Pause/Resume music playing
		if rl.IsKeyPressed(rl.KeyP) {
			pause = !pause

			if pause {
				rl.PauseMusicStream(music)
			} else {
				rl.ResumeMusicStream(music)
			}
		}

		// Add/Remove effect: lowpass filter
		if rl.IsKeyPressed(rl.KeyF) {
			enableEffectLPF = !enableEffectLPF
			if enableEffectLPF {
				rl.AttachAudioStreamProcessor(music.Stream, AudioProcessEffectLPF)
			} else {
				rl.DetachAudioStreamProcessor(music.Stream, AudioProcessEffectLPF)
			}
		}

		// Add/Remove effect: delay
		if rl.IsKeyPressed(rl.KeyD) {
			enableEffectDelay = !enableEffectDelay
			if enableEffectDelay {
				rl.AttachAudioStreamProcessor(music.Stream, AudioProcessEffectDelay)
			} else {
				rl.DetachAudioStreamProcessor(music.Stream, AudioProcessEffectDelay)
			}
		}

		// Get normalized time played for current music stream
		timePlayed = rl.GetMusicTimePlayed(music) / rl.GetMusicTimeLength(music)

		if timePlayed > 1 {
			timePlayed = 1 // Make sure time played is no longer than music
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("MUSIC SHOULD BE PLAYING!", 245, 150, 20, rl.LightGray)

		rl.DrawRectangle(200, 180, 400, 12, rl.LightGray)
		rl.DrawRectangle(200, 180, int32(timePlayed*400), 12, rl.Maroon)
		rl.DrawRectangleLines(200, 180, 400, 12, rl.Gray)

		rl.DrawText("PRESS SPACE TO RESTART MUSIC", 215, 230, 20, rl.LightGray)
		rl.DrawText("PRESS P TO PAUSE/RESUME MUSIC", 208, 260, 20, rl.LightGray)

		rl.DrawText(fmt.Sprintf("PRESS F TO TOGGLE LPF EFFECT: %s", onOff(enableEffectLPF)), 200, 320, 20, rl.Gray)
		rl.DrawText(fmt.Sprintf("PRESS D TO TOGGLE DELAY EFFECT: %s", onOff(enableEffectDelay)), 180, 350, 20, rl.Gray)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadMusicStream(music) // Unload music stream buffers from RAM

	rl.CloseAudioDevice() // Close audio device (music streaming is automatically stopped)

	rl.CloseWindow() // Close window and OpenGL context
}

// AudioProcessEffectLPF - Audio effect: lowpass filter
func AudioProcessEffectLPF(buffer []float32, frames int) {
	const cutoff = 70.0 / 44100.0              // 70 Hz lowpass filter
	const k = cutoff / (cutoff + 0.1591549431) // RC filter formula

	for i := 0; i < frames*2; i += 2 {
		l := buffer[i]
		r := buffer[i+1]

		low[0] += k * (l - low[0])
		low[1] += k * (r - low[1])
		buffer[i] = low[0]
		buffer[i+1] = low[1]
	}
}

// AudioProcessEffectDelay - Audio effect: delay
func AudioProcessEffectDelay(buffer []float32, frames int) {
	for i := 0; i < frames*2; i += 2 {
		leftDelay := delayBuffer[delayReadIndex]
		rightDelay := delayBuffer[delayReadIndex+1]
		delayReadIndex += 2

		if delayReadIndex == len(delayBuffer) {
			delayReadIndex = 0
		}

		buffer[i] = 0.5*buffer[i] + 0.5*leftDelay
		buffer[i+1] = 0.5*buffer[i+1] + 0.5*rightDelay

		delayBuffer[delayWriteIndex] = buffer[i]
		delayBuffer[delayWriteIndex+1] = buffer[i+1]
		delayWriteIndex += 2
		if delayWriteIndex == len(delayBuffer) {
			delayWriteIndex = 0
		}
	}
}

func onOff(b bool) string {
	if b {
		return "ON"
	}
	return "OFF"
}
