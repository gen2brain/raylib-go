package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	raylib.InitWindow(800, 450, "raylib [audio] example - music playing (streaming)")
	raylib.InitAudioDevice()

	music := raylib.LoadMusicStream("guitar_noodling.ogg")
	pause := false

	raylib.PlayMusicStream(music)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.UpdateMusicStream(music) // Update music buffer with new stream data

		// Restart music playing (stop and play)
		if raylib.IsKeyPressed(raylib.KeySpace) {
			raylib.StopMusicStream(music)
			raylib.PlayMusicStream(music)
		}

		// Pause/Resume music playing
		if raylib.IsKeyPressed(raylib.KeyP) {
			pause = !pause

			if pause {
				raylib.PauseMusicStream(music)
			} else {
				raylib.ResumeMusicStream(music)
			}
		}

		// Get timePlayed scaled to bar dimensions (400 pixels)
		timePlayed := raylib.GetMusicTimePlayed(music) / raylib.GetMusicTimeLength(music) * 100 * 4

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)
		raylib.DrawText("MUSIC SHOULD BE PLAYING!", 255, 150, 20, raylib.LightGray)

		raylib.DrawRectangle(200, 200, 400, 12, raylib.LightGray)
		raylib.DrawRectangle(200, 200, int32(timePlayed), 12, raylib.Maroon)
		raylib.DrawRectangleLines(200, 200, 400, 12, raylib.Gray)

		raylib.DrawText("PRESS SPACE TO RESTART MUSIC", 215, 250, 20, raylib.LightGray)
		raylib.DrawText("PRESS P TO PAUSE/RESUME MUSIC", 208, 280, 20, raylib.LightGray)

		raylib.EndDrawing()
	}

	raylib.UnloadMusicStream(music) // Unload music stream buffers from RAM
	raylib.CloseAudioDevice()       // Close audio device (music streaming is automatically stopped)

	raylib.CloseWindow()
}
