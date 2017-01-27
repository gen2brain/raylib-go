package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	raylib.InitWindow(800, 450, "raylib [audio] example - sound loading and playing")

	raylib.InitAudioDevice()

	fxWav := raylib.LoadSound("weird.wav")
	fxOgg := raylib.LoadSound("tanatana.ogg")

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if raylib.IsKeyPressed(raylib.KeySpace) {
			raylib.PlaySound(fxWav)
		}
		if raylib.IsKeyPressed(raylib.KeyEnter) {
			raylib.PlaySound(fxOgg)
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("Press SPACE to PLAY the WAV sound!", 200, 180, 20, raylib.LightGray)
		raylib.DrawText("Press ENTER to PLAY the OGG sound!", 200, 220, 20, raylib.LightGray)

		raylib.EndDrawing()
	}

	raylib.UnloadSound(fxWav)
	raylib.UnloadSound(fxOgg)

	raylib.CloseAudioDevice()

	raylib.CloseWindow()
}
