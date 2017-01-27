package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [textures] example - texture loading and drawing")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	guybrush := raylib.LoadTexture("guybrush.png") // Texture loading

	position := raylib.NewVector2(350.0, 240.0)
	frameRec := raylib.NewRectangle(0, 0, guybrush.Width/7, guybrush.Height)
	currentFrame := int32(0)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if raylib.IsKeyPressed(raylib.KeyRight) {
			currentFrame++

			if currentFrame > 6 {
				currentFrame = 0
			}

			frameRec.X = currentFrame * guybrush.Width / 7
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawTexture(guybrush, 35, 40, raylib.White)
		raylib.DrawRectangleLines(35, 40, guybrush.Width, guybrush.Height, raylib.Lime)

		raylib.DrawTextureRec(guybrush, frameRec, position, raylib.White) // Draw part of the texture

		raylib.DrawRectangleLines(35+frameRec.X, 40+frameRec.Y, frameRec.Width, frameRec.Height, raylib.Red)

		raylib.DrawText("PRESS RIGHT KEY TO", 540, 310, 10, raylib.Gray)
		raylib.DrawText("CHANGE DRAWING RECTANGLE", 520, 330, 10, raylib.Gray)

		raylib.DrawText("Guybrush Ulysses Threepwood,", 100, 300, 10, raylib.Gray)
		raylib.DrawText("main character of the Monkey Island series", 80, 320, 10, raylib.Gray)
		raylib.DrawText("of computer adventure games by LucasArts.", 80, 340, 10, raylib.Gray)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(guybrush)

	raylib.CloseWindow()
}
