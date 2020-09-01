package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialization
	//--------------------------------------------------------------------------------------
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [audio] example - Multichannel sound playing")

	rl.InitAudioDevice() // Initialize audio device

	fxWav := rl.LoadSound("sound.wav")  // Load WAV audio file
	fxOgg := rl.LoadSound("target.ogg") // Load OGG audio file

	rl.SetSoundVolume(fxWav, 0.2)

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second
	//--------------------------------------------------------------------------------------

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		//----------------------------------------------------------------------------------
		if rl.IsKeyPressed(rl.KeyEnter) {
			rl.PlaySoundMulti(fxWav)
		} // Play a new wav sound instance
		if rl.IsKeyPressed(rl.KeySpace) {
			rl.PlaySoundMulti(fxOgg)
		} // Play a new ogg sound instance
		//----------------------------------------------------------------------------------

		// Draw
		//----------------------------------------------------------------------------------
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("MULTICHANNEL SOUND PLAYING", 20, 20, 20, rl.Gray)
		rl.DrawText("Press SPACE to play new ogg instance!", 200, 120, 20, rl.LightGray)
		rl.DrawText("Press ENTER to play new wav instance!", 200, 180, 20, rl.LightGray)

		rl.DrawText(fmt.Sprintf("CONCURRENT SOUNDS PLAYING: %02d", rl.GetSoundsPlaying()), 220, 280, 20, rl.Red)

		rl.EndDrawing()
		//----------------------------------------------------------------------------------
	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	rl.StopSoundMulti() // We must stop the buffer pool before unloading

	rl.UnloadSound(fxWav) // Unload sound data
	rl.UnloadSound(fxOgg) // Unload sound data

	rl.CloseAudioDevice() // Close audio device

	rl.CloseWindow() // Close window and OpenGL context
	//--------------------------------------------------------------------------------------
}
