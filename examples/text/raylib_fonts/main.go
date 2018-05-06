package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const maxFonts = 8

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [text] example - raylib fonts")

	fonts := make([]raylib.Font, maxFonts)
	fonts[0] = raylib.LoadFont("fonts/alagard.png")
	fonts[1] = raylib.LoadFont("fonts/pixelplay.png")
	fonts[2] = raylib.LoadFont("fonts/mecha.png")
	fonts[3] = raylib.LoadFont("fonts/setback.png")
	fonts[4] = raylib.LoadFont("fonts/romulus.png")
	fonts[5] = raylib.LoadFont("fonts/pixantiqua.png")
	fonts[6] = raylib.LoadFont("fonts/alpha_beta.png")
	fonts[7] = raylib.LoadFont("fonts/jupiter_crash.png")

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
	positions := make([]raylib.Vector2, maxFonts)

	var i int32
	for i = 0; i < maxFonts; i++ {
		x := screenWidth/2 - int32(raylib.MeasureTextEx(fonts[i], messages[i], float32(fonts[i].BaseSize*2), spacings[i]).X/2)
		y := 60 + fonts[i].BaseSize + 45*i
		positions[i] = raylib.NewVector2(float32(x), float32(y))
	}

	// Small Y position corrections
	positions[3].Y += 8
	positions[4].Y += 2
	positions[7].Y -= 8

	colors := []raylib.Color{raylib.Maroon, raylib.Orange, raylib.DarkGreen, raylib.DarkBlue, raylib.DarkPurple, raylib.Lime, raylib.Gold, raylib.DarkBrown}

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)
		raylib.DrawText("free fonts included with raylib", 250, 20, 20, raylib.DarkGray)
		raylib.DrawLine(220, 50, 590, 50, raylib.DarkGray)

		for i = 0; i < maxFonts; i++ {
			raylib.DrawTextEx(fonts[i], messages[i], positions[i], float32(fonts[i].BaseSize*2), spacings[i], colors[i])
		}

		raylib.EndDrawing()
	}

	for i = 0; i < maxFonts; i++ {
		raylib.UnloadFont(fonts[i])
	}

	raylib.CloseWindow()
}
