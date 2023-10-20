package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	screenW = int32(1280)
	screenH = int32(500)
)

func main() {

	rl.InitWindow(screenW, screenH, "raylib [textures] example - background scrolling")

	background := rl.LoadTexture("cyberpunk_street_background.png")
	midground := rl.LoadTexture("cyberpunk_street_midground.png")
	foreground := rl.LoadTexture("cyberpunk_street_foreground.png")

	scrollBack := float32(0)
	scrollMid := float32(0)
	scrollFore := float32(0)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		scrollBack -= 0.1
		scrollMid -= 0.5
		scrollFore -= 1

		if scrollBack <= -float32(background.Width)*2 {
			scrollBack = 0
		}

		if scrollMid <= -float32(midground.Width)*2 {
			scrollMid = 0
		}

		if scrollFore <= -float32(foreground.Width)*2 {
			scrollFore = 0
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTextureEx(background, rl.NewVector2(scrollBack, 20), 0, 2, rl.White)
		rl.DrawTextureEx(background, rl.NewVector2(float32(background.Width*2)+scrollBack, 20), 0, 2, rl.White)

		rl.DrawTextureEx(midground, rl.NewVector2(scrollMid, 20), 0, 2, rl.White)
		rl.DrawTextureEx(midground, rl.NewVector2(float32(midground.Width*2)+scrollMid, 20), 0, 2, rl.White)

		rl.DrawTextureEx(foreground, rl.NewVector2(scrollFore, 20), 0, 2, rl.White)
		rl.DrawTextureEx(foreground, rl.NewVector2(float32(foreground.Width*2)+scrollFore, 20), 0, 2, rl.White)

		txt := "BACKGROUND SCROLLING & PARALLAX"
		txtlen := rl.MeasureText(txt, 20)
		rl.DrawText(txt, (screenW/2)-txtlen/2, screenH-50, 20, rl.Black)
		txt = "(c) Cyberpunk Street Environment by Luis Zuno (@ansimuz)"
		txtlen = rl.MeasureText(txt, 10)
		rl.DrawText(txt, (screenW/2)-txtlen/2, screenH-25, 10, rl.Black)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
