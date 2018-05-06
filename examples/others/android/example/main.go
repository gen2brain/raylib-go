package main

import (
	"os"
	"runtime"

	"github.com/gen2brain/raylib-go/raylib"
)

// Game states
const (
	Logo = iota
	Title
	GamePlay
	Ending
)

func init() {
	raylib.SetCallbackFunc(main)
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.SetConfigFlags(raylib.FlagVsyncHint)

	raylib.InitWindow(screenWidth, screenHeight, "Android example")

	raylib.InitAudioDevice()

	currentScreen := Logo
	windowShouldClose := false

	texture := raylib.LoadTexture("raylib_logo.png") // Load texture (placed on assets folder)
	fx := raylib.LoadSound("coin.wav")               // Load WAV audio file (placed on assets folder)
	ambient := raylib.LoadMusicStream("ambient.ogg") // Load music

	raylib.PlayMusicStream(ambient)

	framesCounter := 0 // Used to count frames

	//raylib.SetTargetFPS(60)

	for !windowShouldClose {
		raylib.UpdateMusicStream(ambient)

		if runtime.GOOS == "android" && raylib.IsKeyDown(raylib.KeyBack) || raylib.WindowShouldClose() {
			windowShouldClose = true
		}

		switch currentScreen {
		case Logo:
			framesCounter++ // Count frames

			// Wait for 4 seconds (240 frames) before jumping to Title screen
			if framesCounter > 240 {
				currentScreen = Title
			}
			break
		case Title:
			// Press enter to change to GamePlay screen
			if raylib.IsGestureDetected(raylib.GestureTap) {
				raylib.PlaySound(fx)
				currentScreen = GamePlay
			}
			break
		case GamePlay:
			// Press enter to change to Ending screen
			if raylib.IsGestureDetected(raylib.GestureTap) {
				raylib.PlaySound(fx)
				currentScreen = Ending
			}
			break
		case Ending:
			// Press enter to return to Title screen
			if raylib.IsGestureDetected(raylib.GestureTap) {
				raylib.PlaySound(fx)
				currentScreen = Title
			}
			break
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		switch currentScreen {
		case Logo:
			raylib.DrawText("LOGO SCREEN", 20, 20, 40, raylib.LightGray)
			raylib.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2, raylib.White)
			raylib.DrawText("WAIT for 4 SECONDS...", 290, 400, 20, raylib.Gray)
			break
		case Title:
			raylib.DrawRectangle(0, 0, screenWidth, screenHeight, raylib.Green)
			raylib.DrawText("TITLE SCREEN", 20, 20, 40, raylib.DarkGreen)
			raylib.DrawText("TAP SCREEN to JUMP to GAMEPLAY SCREEN", 160, 220, 20, raylib.DarkGreen)
			break
		case GamePlay:
			raylib.DrawRectangle(0, 0, screenWidth, screenHeight, raylib.Purple)
			raylib.DrawText("GAMEPLAY SCREEN", 20, 20, 40, raylib.Maroon)
			raylib.DrawText("TAP SCREEN to JUMP to ENDING SCREEN", 170, 220, 20, raylib.Maroon)
			break
		case Ending:
			raylib.DrawRectangle(0, 0, screenWidth, screenHeight, raylib.Blue)
			raylib.DrawText("ENDING SCREEN", 20, 20, 40, raylib.DarkBlue)
			raylib.DrawText("TAP SCREEN to RETURN to TITLE SCREEN", 160, 220, 20, raylib.DarkBlue)
			break
		default:
			break
		}

		raylib.EndDrawing()
	}

	raylib.UnloadSound(fx)            // Unload sound data
	raylib.UnloadMusicStream(ambient) // Unload music stream data
	raylib.CloseAudioDevice()         // Close audio device (music streaming is automatically stopped)
	raylib.UnloadTexture(texture)     // Unload texture data
	raylib.CloseWindow()              // Close window

	os.Exit(0)
}
