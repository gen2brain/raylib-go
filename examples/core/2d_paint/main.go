package main

import (
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const colorCount = 23

func main() {
	const screenWidth, screenHeight int32 = 800, 450

	rl.InitWindow(screenWidth, screenHeight, "raypaint")

	// Colors to choose from
	var colors = [colorCount]rl.Color{
		rl.RayWhite, rl.Yellow, rl.Gold, rl.Orange, rl.Pink, rl.Red, rl.Maroon, rl.Green, rl.Lime, rl.DarkGreen,
		rl.SkyBlue, rl.Blue, rl.DarkBlue, rl.Purple, rl.Violet, rl.DarkPurple, rl.Beige, rl.Brown, rl.DarkBrown,
		rl.LightGray, rl.Gray, rl.DarkGray, rl.Black,
	}

	// Define colorsRecs data (for every rectangle)
	var colorRecs = [colorCount]rl.Rectangle{}

	for i := 0; i < colorCount; i++ {
		colorRecs[i].X = float32(10 + 30*i + 2*i)
		colorRecs[i].Y = 10
		colorRecs[i].Width = 30
		colorRecs[i].Height = 30
	}

	colorSelected := 0
	colorSelectedPrev := colorSelected
	colorMouseHover := 0
	brushSize := 20

	var btnSaveRec = rl.Rectangle{750, 10, 40, 30}
	btnSaveMouseHover := false
	showSaveMessage := false
	saveMessageCounter := 0

	checkSaveHover := func() rl.Color {
		if btnSaveMouseHover {
			return rl.Red
		}
		return rl.Black
	}

	// Create a RenderTexture2D to use as a canvas
	var target rl.RenderTexture2D = rl.LoadRenderTexture(screenWidth, screenHeight)

	// Clear render texture before entering the game loop
	rl.BeginTextureMode(target)
	rl.ClearBackground(colors[0])
	rl.EndTextureMode()

	rl.SetTargetFPS(120)

	// Main game loop
	for !rl.WindowShouldClose() {
		mousePos := rl.GetMousePosition()

		// Move between colors with keys
		if rl.IsKeyPressed(rl.KeyRight) {
			colorSelected++
		} else if rl.IsKeyPressed(rl.KeyLeft) {
			colorSelected--
		}

		if colorSelected >= colorCount {
			colorSelected = colorCount - 1
		} else if colorSelected < 0 {
			colorSelected = 0
		}

		// Choose color with mouse
		for i := 0; i < colorCount; i++ {
			if rl.CheckCollisionPointRec(mousePos, colorRecs[i]) {
				colorMouseHover = i
				break
			} else {
				colorMouseHover = -1
			}
		}

		if colorMouseHover >= 0 && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			colorSelected = colorMouseHover
			colorSelectedPrev = colorSelected
		}

		// Change brush size
		brushSize += int(rl.GetMouseWheelMove() * 5)
		if brushSize < 2 {
			brushSize = 2
		}

		if brushSize > 50 {
			brushSize = 50
		}

		if rl.IsKeyPressed(rl.KeyC) {
			// Clear render texture to clear color
			rl.BeginTextureMode(target)
			rl.ClearBackground(colors[0])
			rl.EndTextureMode()
		}

		if rl.IsMouseButtonDown(rl.MouseLeftButton) || rl.GetGestureDetected() == rl.GestureDrag {
			// Clear render texture to clear color
			rl.BeginTextureMode(target)

			if mousePos.Y > 50 {
				rl.DrawCircle(int32(mousePos.X), int32(mousePos.Y), float32(brushSize), colors[colorSelected])
			}

			rl.EndTextureMode()
		}

		if rl.IsMouseButtonDown(rl.MouseRightButton) {
			colorSelected = 0

			// Erase circle from render texture
			rl.BeginTextureMode(target)

			if mousePos.Y > 50 {
				rl.DrawCircle(int32(mousePos.X), int32(mousePos.Y), float32(brushSize), colors[0])
			}

			rl.EndTextureMode()
		} else {
			colorSelected = colorSelectedPrev
		}

		if rl.CheckCollisionPointRec(mousePos, btnSaveRec) {
			btnSaveMouseHover = true
		} else {
			btnSaveMouseHover = false
		}

		if btnSaveMouseHover && rl.IsMouseButtonReleased(rl.MouseLeftButton) || rl.IsKeyPressed(rl.KeyS) {
			image := rl.LoadImageFromTexture(target.Texture)
			rl.ImageFlipVertical(*&image)
			rl.ExportImage(*image, "export.png")
			rl.UnloadImage(image)
			showSaveMessage = true
		}

		if showSaveMessage {
			// On saving, show a full screen message for 2 seconds
			saveMessageCounter++
			if saveMessageCounter > 240 {
				showSaveMessage = false
				saveMessageCounter = 0
			}
		}

		// Draw
		//----------------------------------------------------------------------------------
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		// NOTE: Render texture must be y-flipped due to default OpenGL coordinates (left-bottom)
		rl.DrawTextureRec(target.Texture, rl.Rectangle{0, 0, float32(target.Texture.Width), float32(-target.Texture.Height)}, rl.Vector2{0, 0}, rl.White)

		if mousePos.Y > 50 {
			if rl.IsMouseButtonDown(rl.MouseRightButton) {
				rl.DrawCircleLines(int32(mousePos.X), int32(mousePos.Y), float32(brushSize), rl.Gray)
			} else {
				rl.DrawCircle(rl.GetMouseX(), rl.GetMouseY(), float32(brushSize), colors[colorSelected])
			}
		}

		// Draw top panel
		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), 50, rl.RayWhite)
		rl.DrawLine(0, 50, int32(rl.GetScreenWidth()), 50, rl.LightGray)

		// Draw color selection rectangles
		for i := 0; i < colorCount; i++ {
			rl.DrawRectangleRec(colorRecs[i], colors[i])
		}

		rl.DrawRectangleLines(10, 10, 30, 30, rl.LightGray)

		if colorMouseHover >= 0 {
			rl.DrawRectangleRec(colorRecs[colorMouseHover], rl.Fade(rl.White, 0.0))
		}

		rl.DrawRectangleLinesEx(rl.Rectangle{
			colorRecs[colorSelected].X - 2, colorRecs[colorSelected].Y - 2, colorRecs[colorSelected].Width + 4, colorRecs[colorSelected].Height + 4,
		}, 2, rl.Black)

		// Draw save image button
		rl.DrawRectangleLinesEx(btnSaveRec, 2, checkSaveHover())
		rl.DrawText("SAVE!", 755, 20, 10, checkSaveHover())

		if showSaveMessage {
			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.RayWhite, 0.8))
			rl.DrawRectangle(0, 150, int32(rl.GetScreenWidth()), 80, rl.Black)
			rl.DrawText("IMAGE SAVED:  export.png", 150, 180, 20, rl.RayWhite)
		}

		rl.EndDrawing()

	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	rl.UnloadRenderTexture(target)

	rl.CloseWindow()

	os.Exit(0)

}
