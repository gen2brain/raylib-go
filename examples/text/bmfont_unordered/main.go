package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [text] example - bmfont unordered loading and drawing")

	// NOTE: Using chars outside the [32..127] limits!
	// NOTE: If a character is not found in the font, it just renders a space
	msg := "ASCII extended characters:\n¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆ\nÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæ\nçèéêëìíîïðñòóôõö÷øùúûüýþÿ"

	// NOTE: Loaded font has an unordered list of characters (chars in the range 32..255)
	font := raylib.LoadFont("fonts/pixantiqua.fnt") // BMFont (AngelCode)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("Font name:       PixAntiqua", 40, 50, 20, raylib.Gray)
		raylib.DrawText(fmt.Sprintf("Font base size:           %d", font.BaseSize), 40, 80, 20, raylib.Gray)
		raylib.DrawText(fmt.Sprintf("Font chars number:     %d", font.CharsCount), 40, 110, 20, raylib.Gray)

		raylib.DrawTextEx(font, msg, raylib.NewVector2(40, 180), float32(font.BaseSize), 0, raylib.Maroon)

		raylib.EndDrawing()
	}

	raylib.UnloadFont(font) // AngelCode Font unloading

	raylib.CloseWindow()
}
