package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	screenW = int32(800)
	screenH = int32(450)

	colors = []rl.Color{rl.RayWhite, rl.Yellow, rl.Gold, rl.Pink, rl.Red, rl.Maroon, rl.Green, rl.Lime, rl.DarkGreen, rl.SkyBlue, rl.Blue, rl.DarkBlue, rl.Purple, rl.Violet, rl.DarkPurple, rl.Beige, rl.Brown, rl.DarkBrown, rl.LightGray, rl.Gray, rl.DarkGray, rl.Black}

	colorRecs []rl.Rectangle
)

func main() {

	rl.InitWindow(screenW, screenH, "raylib [textures] example - mouse painting")

	for i := 0; i < len(colors); i++ {
		colorRecs = append(colorRecs, rl.NewRectangle(10+(30*float32(i))+(2*float32(i)), 10, 30, 30))
	}
	colorSelected, colorMouseHover := 1, 0
	colorPrev := colorSelected
	brushSize := float32(20)
	mousePressed := false

	btnSaveRec := rl.NewRectangle(750, 10, 40, 30)
	btnSaveMouseHover, showSaveMsg := false, false
	saveMsgCount := 0

	target := rl.LoadRenderTexture(screenW, screenH)

	rl.BeginTextureMode(target)
	rl.ClearBackground(colors[0])
	rl.EndTextureMode()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		mousePos := rl.GetMousePosition()

		if rl.IsKeyPressed(rl.KeyRight) {
			colorSelected++
		} else if rl.IsKeyPressed(rl.KeyLeft) {
			colorSelected--
		}

		if colorSelected >= len(colors) {
			colorSelected = len(colors) - 1
		} else if colorSelected < 0 {
			colorSelected = 0
		}

		for i := 0; i < len(colorRecs); i++ {
			if rl.CheckCollisionPointRec(mousePos, colorRecs[i]) {
				colorMouseHover = i
			} else {
				colorMouseHover = -1
			}

			if colorMouseHover >= 0 && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				colorSelected = colorMouseHover
				colorPrev = colorSelected
			}
		}

		if rl.IsKeyPressed(rl.KeyUp) {
			brushSize += 5
		} else if rl.IsKeyPressed(rl.KeyDown) {
			brushSize -= 5
		}
		brushSize += rl.GetMouseWheelMove() * 5
		if brushSize < 2 {
			brushSize = 2
		}
		if brushSize > 50 {
			brushSize = 50
		}

		if rl.IsKeyPressed(rl.KeyC) {
			rl.BeginTextureMode(target)
			rl.ClearBackground(colors[0])
			rl.EndTextureMode()
		}

		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			rl.BeginTextureMode(target)
			if mousePos.Y > 50 {
				rl.DrawCircle(int32(mousePos.X), int32(mousePos.Y), brushSize, colors[colorSelected])
			}
			rl.EndTextureMode()
		}

		if rl.IsMouseButtonDown(rl.MouseButtonRight) {
			if !mousePressed {
				colorPrev = colorSelected
				colorSelected = 0
			}

			mousePressed = true

			rl.BeginTextureMode(target)
			if mousePos.Y > 50 {
				rl.DrawCircle(int32(mousePos.X), int32(mousePos.Y), brushSize, colors[0])
			}
			rl.EndTextureMode()
		} else if rl.IsMouseButtonReleased(rl.MouseButtonRight) && mousePressed {
			colorSelected = colorPrev
			mousePressed = false
		}

		if rl.CheckCollisionPointRec(mousePos, btnSaveRec) {
			btnSaveMouseHover = true
		} else {
			btnSaveMouseHover = false
		}

		if btnSaveMouseHover && rl.IsMouseButtonReleased(rl.MouseButtonLeft) || rl.IsKeyPressed(rl.KeyS) {
			image := rl.LoadImageFromTexture(target.Texture)
			rl.ImageFlipVertical(image)
			rl.ExportImage(*image, "raylib_mouse_painting.png")
			rl.UnloadImage(image)
			showSaveMsg = true
		}

		if showSaveMsg {
			saveMsgCount++
			if saveMsgCount > 240 {
				showSaveMsg = false
				saveMsgCount = 0
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTextureRec(target.Texture, rl.NewRectangle(0, 0, float32(target.Texture.Width), -float32(target.Texture.Height)), rl.Vector2Zero(), rl.White)

		if mousePos.Y > 50 {
			if rl.IsMouseButtonDown(rl.MouseButtonRight) {
				rl.DrawCircleLines(int32(mousePos.X), int32(mousePos.Y), brushSize, rl.Gray)
			} else {
				rl.DrawCircle(rl.GetMouseX(), rl.GetMouseY(), brushSize, colors[colorSelected])
			}
		}

		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), 50, rl.RayWhite)
		rl.DrawLine(0, 50, int32(rl.GetScreenWidth()), 50, rl.LightGray)

		for i := 0; i < len(colors); i++ {
			rl.DrawRectangleRec(colorRecs[i], colors[i])
		}

		rl.DrawRectangleLines(10, 10, 30, 30, rl.LightGray)

		if colorMouseHover >= 0 {
			rl.DrawRectangleRec(colorRecs[colorMouseHover], rl.Fade(rl.White, 0.6))
		}

		rl.DrawRectangleLinesEx(rl.NewRectangle(colorRecs[colorSelected].X-2, colorRecs[colorSelected].Y-2, colorRecs[colorSelected].Width+4, colorRecs[colorSelected].Height+4), 2, rl.Black)

		if btnSaveMouseHover {
			rl.DrawRectangleLinesEx(btnSaveRec, 2, rl.Red)
			rl.DrawText("SAVE!", 755, 20, 10, rl.Red)
		} else {
			rl.DrawRectangleLinesEx(btnSaveRec, 2, rl.Black)
			rl.DrawText("SAVE!", 755, 20, 10, rl.Black)
		}

		if showSaveMsg {
			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.RayWhite, 0.8))
			rl.DrawRectangle(0, screenH/2-40, int32(rl.GetScreenWidth()), 80, rl.Black)
			txt := "IMG SAVED"
			txtlen := rl.MeasureText(txt, 20)
			rl.DrawText(txt, screenW/2-txtlen/2, screenH/2-10, 20, rl.RayWhite)
		}

		rl.DrawText("hold left mouse to draw right mouse to erase | mouse wheel up/down arrows change brush size  |  right/left arrows change color | c key to clear", 10, screenH-15, 10, rl.Black)

		rl.EndDrawing()
	}

	rl.UnloadRenderTexture(target)

	rl.CloseWindow()
}
