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

	// Get image size
	size := img.Bounds().Size()

	// Dynamic memory allocation to store pixels data (Color type)
	pixels := make([]raylib.Color, size.X*size.Y)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			pixels[x+y*size.Y] = raylib.NewColor(uint8(r), uint8(g), uint8(b), uint8(a))
		}
	}

	// Load pixels data into an image structure and create texture
	imEx := raylib.LoadImageEx(pixels, int32(size.X), int32(size.Y))
	texture := raylib.LoadTextureFromImage(imEx)

	// Unload CPU (RAM) image data
	raylib.UnloadImage(imEx)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)
		raylib.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2, raylib.White)
		raylib.DrawText("this IS a texture loaded from an image.Image!", 285, 370, 10, raylib.Gray)

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(texture)

	raylib.CloseWindow()
}
