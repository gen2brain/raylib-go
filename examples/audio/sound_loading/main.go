package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib [audio] example - sound loading and playing")

	rl.InitAudioDevice()

	fxWav := rl.LoadSound("weird.wav")
	fxOgg := rl.LoadSound("tanatana.ogg")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeySpace) {
			rl.PlaySound(fxWav)
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			rl.PlaySound(fxOgg)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("Press SPACE to PLAY the WAV sound!", 200, 180, 20, rl.LightGray)
		rl.DrawText("Press ENTER to PLAY the OGG sound!", 200, 220, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.UnloadSound(fxWav)
	rl.UnloadSound(fxOgg)

	rl.CloseAudioDevice()

	rl.CloseWindow()
}
