package main

import (
	"embed"
	"fmt"
	"os"
	"runtime"

	"github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/*
var content embed.FS

func init() {
	rl.SetMain(main)
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [others] example - asset loading")

	var assets *rl.Asset
	switch runtime.GOOS {
	case "android":
		assets = rl.NewAsset("")
	case "js":
		assets = rl.NewAssetFromFS(content, "assets")
	default:
		assets = rl.NewAsset("assets")
	}

	data, err := assets.ReadFile("raylib_logo.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	image := rl.LoadImageFromMemory(".png", data, int32(len(data)))
	texture := rl.LoadTextureFromImage(image)

	rl.UnloadImage(image)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2, rl.White)

		rl.DrawText("this IS a texture loaded from an image!", 300, 370, 10, rl.Gray)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)

	rl.CloseWindow()
}
