package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	screenW       = int32(1280)
	screenH       = int32(720)
	blendCountMax = 4
	blendMode     = 0
)

func main() {

	rl.InitWindow(screenW, screenH, "raylib [textures] example - blend modes")

	bgImg := rl.LoadImage("cyberpunk_street_background.png")
	bgTex := rl.LoadTextureFromImage(bgImg)

	fgImg := rl.LoadImage("cyberpunk_street_foreground.png")
	fgTex := rl.LoadTextureFromImage(fgImg)

	rl.UnloadImage(bgImg)
	rl.UnloadImage(fgImg)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		if rl.IsKeyPressed(rl.KeySpace) {
			if blendMode >= blendCountMax-1 {
				blendMode = 0
			} else {
				blendMode++
			}

		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(bgTex, screenW/2-bgTex.Width/2, screenH/2-bgTex.Height/2, rl.White)
		rl.BeginBlendMode(rl.BlendMode(blendMode))
		rl.DrawTexture(fgTex, screenW/2-fgTex.Width/2, screenH/2-fgTex.Height/2, rl.White)
		rl.EndBlendMode()

		txt := "Press SPACE to change blend modes"
		txtlen := rl.MeasureText(txt, 20)
		rl.DrawText(txt, (screenW/2)-txtlen/2, screenH-70, 20, rl.Black)
		switch rl.BlendMode(blendMode) {
		case rl.BlendAlpha:
			txt = "Current rl.BlendAlpha"
			txtlen = rl.MeasureText(txt, 20)
			rl.DrawText(txt, (screenW/2)-txtlen/2, screenH-50, 20, rl.Black)
		case rl.BlendAdditive:
			txt = "Current rl.BlendAdditive"
			txtlen = rl.MeasureText(txt, 20)
			rl.DrawText(txt, (screenW/2)-txtlen/2, screenH-50, 20, rl.Black)
		case rl.BlendMultiplied:
			txt = "Current rl.BlendMultiplied"
			txtlen = rl.MeasureText(txt, 20)
			rl.DrawText(txt, (screenW/2)-txtlen/2, screenH-50, 20, rl.Black)
		case rl.BlendAddColors:
			txt = "Current rl.BlendAddColors"
			txtlen = rl.MeasureText(txt, 20)
			rl.DrawText(txt, (screenW/2)-txtlen/2, screenH-50, 20, rl.Black)
		}

		txt = "(c) Cyberpunk Street Environment by Luis Zuno (@ansimuz)"
		txtlen = rl.MeasureText(txt, 10)
		rl.DrawText(txt, (screenW/2)-txtlen/2, screenH-25, 10, rl.Black)

		rl.EndDrawing()
	}

	rl.UnloadTexture(fgTex)
	rl.UnloadTexture(bgTex)

	rl.CloseWindow()
}
