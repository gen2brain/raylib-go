package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const screenW = int32(1280)
const screenH = int32(720)

type gameScreen int

const (
	LOGO = iota
	TITLE
	GAMEPLAY
	ENDING
)

func main() {

	rl.InitWindow(screenW, screenH, "raylib [core] example - basic screen manager")

	var currentScreen gameScreen
	currentScreen = LOGO
	frames := 0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		switch currentScreen {
		case LOGO:
			frames++
			if frames > 240 {
				currentScreen = TITLE
			}
		case TITLE:
			if rl.IsKeyPressed(rl.KeyEnter) {
				currentScreen = GAMEPLAY
			}
		case GAMEPLAY:
			if rl.IsKeyPressed(rl.KeyEnter) {
				currentScreen = ENDING
			}
		case ENDING:
			if rl.IsKeyPressed(rl.KeyEnter) {
				currentScreen = LOGO
				frames = 0
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		rec := rl.NewRectangle(0, 0, float32(screenW), float32(screenH))
		switch currentScreen {
		case LOGO:
			txt := "YOUR LOGO GOES HERE"
			txtlen := rl.MeasureText(txt, 50)
			rl.DrawText(txt, screenW/2-txtlen/2-3, screenH/2-50+3, 50, rl.Magenta)
			rl.DrawText(txt, screenW/2-txtlen/2-1, screenH/2-50+1, 50, rl.Black)
			rl.DrawText(txt, screenW/2-txtlen/2, screenH/2-50, 50, rl.White)
			txt = "this message disappears in " + fmt.Sprint(240-frames) + " frames"
			txtlen = rl.MeasureText(txt, 30)
			rl.DrawText(txt, screenW/2-txtlen/2-3, screenH/2+3, 30, rl.Magenta)
			rl.DrawText(txt, screenW/2-txtlen/2-1, screenH/2+1, 30, rl.Black)
			rl.DrawText(txt, screenW/2-txtlen/2, screenH/2, 30, rl.White)
		case TITLE:
			rl.DrawRectangleRec(rec, rl.DarkGreen)
			txt := "AN AMAZING TITLE GOES HERE"
			txtlen := rl.MeasureText(txt, 50)
			rl.DrawText(txt, screenW/2-txtlen/2-2, screenH/2-50+2, 50, rl.Black)
			rl.DrawText(txt, screenW/2-txtlen/2, screenH/2-50, 50, rl.White)
			txt = "press enter to move to next screen"
			txtlen = rl.MeasureText(txt, 30)
			rl.DrawText(txt, screenW/2-txtlen/2-2, screenH/2+2, 30, rl.Black)
			rl.DrawText(txt, screenW/2-txtlen/2, screenH/2, 30, rl.White)
		case GAMEPLAY:
			rl.DrawRectangleRec(rec, rl.DarkPurple)
			txt := "FUN GAMEPLAY GOES HERE"
			txtlen := rl.MeasureText(txt, 50)
			rl.DrawText(txt, screenW/2-txtlen/2-2, screenH/2-50+2, 50, rl.Black)
			rl.DrawText(txt, screenW/2-txtlen/2, screenH/2-50, 50, rl.White)
			txt = "press enter to move to next screen"
			txtlen = rl.MeasureText(txt, 30)
			rl.DrawText(txt, screenW/2-txtlen/2-2, screenH/2+2, 30, rl.Black)
			rl.DrawText(txt, screenW/2-txtlen/2, screenH/2, 30, rl.White)
		case ENDING:
			rl.DrawRectangleRec(rec, rl.DarkBlue)
			txt := "A DRAMATIC ENDING GOES HERE"
			txtlen := rl.MeasureText(txt, 50)
			rl.DrawText(txt, screenW/2-txtlen/2-2, screenH/2-50+2, 50, rl.Black)
			rl.DrawText(txt, screenW/2-txtlen/2, screenH/2-50, 50, rl.White)
			txt = "press enter to move to next screen"
			txtlen = rl.MeasureText(txt, 30)
			rl.DrawText(txt, screenW/2-txtlen/2-2, screenH/2+2, 30, rl.Black)
			rl.DrawText(txt, screenW/2-txtlen/2, screenH/2, 30, rl.White)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
