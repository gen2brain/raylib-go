/*******************************************************************************************
*
*   raylib [audio] example - Playing sound multiple times
*
*   Example originally created with raylib 4.6
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2023 Jeffery Myers (@JeffM2501)
*
********************************************************************************************/
package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	screenWidth  = 800
	screenHeight = 450
	maxSounds    = 10
)

var soundArray [maxSounds]rl.Sound
var currentSound int32

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [audio] example - playing sound multiple times")

	// Initialize audio device
	rl.InitAudioDevice()

	// Load WAV audio file into the first slot as the 'source' sound
	// this sound owns the sample data
	soundArray[0] = rl.LoadSound("sound.wav")

	for i := 1; i < maxSounds; i++ {
		// Load an alias of the sound into slots 1-9.
		// These do not own the sound data, but can be played
		soundArray[i] = rl.LoadSoundAlias(soundArray[0])
	}

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		if rl.IsKeyPressed(rl.KeySpace) {
			rl.PlaySound(soundArray[currentSound]) // play the next open sound slot
			currentSound++                         // increment the sound slot
			if currentSound >= maxSounds {         // if the sound slot is out of bounds, go back to 0
				currentSound = 0
			}

			// Note: a better way would be to look at the list for the first sound that is not playing and use that slot
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Press SPACE to PLAY a WAV sound!", 200, 180, 20, rl.LightGray)
		rl.EndDrawing()
	}

	rl.UnloadSound(soundArray[0]) // Unload source sound data
	rl.CloseAudioDevice()         // Close audio device

	rl.CloseWindow() // Close window and OpenGL context
}
