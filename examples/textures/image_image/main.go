package main

import (
	"image/png"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - texture from image.Image")

	r, err := os.Open("raylib_logo.png")
	if err != nil {
		rl.TraceLog(rl.LogError, err.Error())
	}
	defer r.Close()

	img, err := png.Decode(r)
	if err != nil {
		rl.TraceLog(rl.LogError, err.Error())
	}

	// Create rl.Image from Go image.Image and create texture
	im := rl.NewImageFromImage(img)
	texture := rl.LoadTextureFromImage(im)

	// Unload CPU (RAM) image data
	rl.UnloadImage(im)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyS) {
			rimg := rl.LoadImageFromTexture(texture)

			f, err := os.Create("image_saved.png")
			if err != nil {
				rl.TraceLog(rl.LogError, err.Error())
			}

			err = png.Encode(f, rimg.ToImage())
			if err != nil {
				rl.TraceLog(rl.LogError, err.Error())
			}

			f.Close()
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("PRESS S TO SAVE IMAGE FROM TEXTURE", 20, 20, 12, rl.LightGray)
		rl.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2, rl.White)
		rl.DrawText("this IS a texture loaded from an image.Image!", 285, 370, 10, rl.Gray)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)

	rl.CloseWindow()
}
