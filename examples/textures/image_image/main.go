package main

import (
	"image/png"
	"os"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [textures] example - texture from image.Image")

	r, err := os.Open("raylib_logo.png")
	if err != nil {
		raylib.TraceLog(raylib.LogError, err.Error())
	}
	defer r.Close()

	img, err := png.Decode(r)
	if err != nil {
		raylib.TraceLog(raylib.LogError, err.Error())
	}

	// Create raylib.Image from Go image.Image and create texture
	im := raylib.NewImageFromImage(img)
	texture := raylib.LoadTextureFromImage(im)

	// Unload CPU (RAM) image data
	raylib.UnloadImage(im)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if raylib.IsKeyPressed(raylib.KeyS) {
			rimg := raylib.GetTextureData(texture)

			f, err := os.Create("image_saved.png")
			if err != nil {
				raylib.TraceLog(raylib.LogError, err.Error())
			}

			err = png.Encode(f, rimg.ToImage())
			if err != nil {
				raylib.TraceLog(raylib.LogError, err.Error())
			}

			f.Close()
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("PRESS S TO SAVE IMAGE FROM TEXTURE", 20, 20, 12, raylib.LightGray)
		raylib.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2, raylib.White)
		raylib.DrawText("this IS a texture loaded from an image.Image!", 285, 370, 10, raylib.Gray)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(texture)

	raylib.CloseWindow()
}
