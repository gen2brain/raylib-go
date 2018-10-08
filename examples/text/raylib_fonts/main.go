package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const maxFonts = 8

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - raylib fonts")

	fonts := make([]rl.Font, maxFonts)
	fonts[0] = rl.LoadFont("fonts/alagard.png")
	fonts[1] = rl.LoadFont("fonts/pixelplay.png")
	fonts[2] = rl.LoadFont("fonts/mecha.png")
	fonts[3] = rl.LoadFont("fonts/setback.png")
	fonts[4] = rl.LoadFont("fonts/romulus.png")
	fonts[5] = rl.LoadFont("fonts/pixantiqua.png")
	fonts[6] = rl.LoadFont("fonts/alpha_beta.png")
	fonts[7] = rl.LoadFont("fonts/jupiter_crash.png")

	messages := []string{
		"ALAGARD FONT designed by Hewett Tsoi",
		"PIXELPLAY FONT designed by Aleksander Shevchuk",
		"MECHA FONT designed by Captain Falcon",
		"SETBACK FONT designed by Brian Kent (AEnigma)",
		"ROMULUS FONT designed by Hewett Tsoi",
		"PIXANTIQUA FONT designed by Gerhard Grossmann",
		"ALPHA_BETA FONT designed by Brian Kent (AEnigma)",
		"JUPITER_CRASH FONT designed by Brian Kent (AEnigma)",
	}

	spacings := []float32{2, 4, 8, 4, 3, 4, 4, 1}
	positions := make([]rl.Vector2, maxFonts)

	var i int32
	for i = 0; i < maxFonts; i++ {
		x := screenWidth/2 - int32(rl.MeasureTextEx(fonts[i], messages[i], float32(fonts[i].BaseSize*2), spacings[i]).X/2)
		y := 60 + fonts[i].BaseSize + 45*i
		positions[i] = rl.NewVector2(float32(x), float32(y))
	}

	// Small Y position corrections
	positions[3].Y += 8
	positions[4].Y += 2
	positions[7].Y -= 8

	colors := []rl.Color{rl.Maroon, rl.Orange, rl.DarkGreen, rl.DarkBlue, rl.DarkPurple, rl.Lime, rl.Gold, rl.DarkBrown}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("free fonts included with raylib", 250, 20, 20, rl.DarkGray)
		rl.DrawLine(220, 50, 590, 50, rl.DarkGray)

		for i = 0; i < maxFonts; i++ {
			rl.DrawTextEx(fonts[i], messages[i], positions[i], float32(fonts[i].BaseSize*2), spacings[i], colors[i])
		}

		rl.EndDrawing()
	}

	for i = 0; i < maxFonts; i++ {
		rl.UnloadFont(fonts[i])
	}

	rl.CloseWindow()
}
