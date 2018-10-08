package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib [audio] example - music playing (streaming)")
	rl.InitAudioDevice()

	music := rl.LoadMusicStream("guitar_noodling.ogg")
	pause := false

	rl.PlayMusicStream(music)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
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

		// Get timePlayed scaled to bar dimensions (400 pixels)
		timePlayed := rl.GetMusicTimePlayed(music) / rl.GetMusicTimeLength(music) * 100 * 4

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("MUSIC SHOULD BE PLAYING!", 255, 150, 20, rl.LightGray)

		rl.DrawRectangle(200, 200, 400, 12, rl.LightGray)
		rl.DrawRectangle(200, 200, int32(timePlayed), 12, rl.Maroon)
		rl.DrawRectangleLines(200, 200, 400, 12, rl.Gray)

		rl.DrawText("PRESS SPACE TO RESTART MUSIC", 215, 250, 20, rl.LightGray)
		rl.DrawText("PRESS P TO PAUSE/RESUME MUSIC", 208, 280, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.UnloadMusicStream(music) // Unload music stream buffers from RAM
	rl.CloseAudioDevice()       // Close audio device (music streaming is automatically stopped)

	rl.CloseWindow()
}
