package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [textures] examples - texture source and destination rectangles")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	scarfy := raylib.LoadTexture("scarfy.png") // Texture loading

	frameWidth := float32(scarfy.Width) / 7
	frameHeight := float32(scarfy.Height)

	// NOTE: Source rectangle (part of the texture to use for drawing)
	sourceRec := raylib.NewRectangle(0, 0, frameWidth, frameHeight)

	// NOTE: Destination rectangle (screen rectangle where drawing part of texture)
	destRec := raylib.NewRectangle(float32(screenWidth)/2, float32(screenHeight)/2, frameWidth*2, frameHeight*2)

	// NOTE: Origin of the texture (rotation/scale point), it's relative to destination rectangle size
	origin := raylib.NewVector2(float32(frameWidth), float32(frameHeight))

	rotation := float32(0)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update
		rotation++

		// Draw
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		// NOTE: Using DrawTexturePro() we can easily rotate and scale the part of the texture we draw
		// sourceRec defines the part of the texture we use for drawing
		// destRec defines the rectangle where our texture part will fit (scaling it to fit)
		// origin defines the point of the texture used as reference for rotation and scaling
		// rotation defines the texture rotation (using origin as rotation point)
		raylib.DrawTexturePro(scarfy, sourceRec, destRec, origin, rotation, raylib.White)

		raylib.DrawLine(int32(destRec.X), 0, int32(destRec.X), screenHeight, raylib.Gray)
		raylib.DrawLine(0, int32(destRec.Y), screenWidth, int32(destRec.Y), raylib.Gray)

		raylib.DrawText("(c) Scarfy sprite by Eiden Marsal", screenWidth-200, screenHeight-20, 10, raylib.Gray)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(scarfy)

	raylib.CloseWindow()
}
