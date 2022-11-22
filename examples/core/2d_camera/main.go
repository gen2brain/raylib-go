package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

var (
	maxBuildings int = 100
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - 2d camera")

	player := rl.NewRectangle(400, 280, 40, 40)

	buildings := make([]rl.Rectangle, maxBuildings)
	buildColors := make([]rl.Color, maxBuildings)

	spacing := float32(0)

	for i := 0; i < maxBuildings; i++ {
		r := rl.Rectangle{}
		r.Width = float32(rl.GetRandomValue(50, 200))
		r.Height = float32(rl.GetRandomValue(100, 800))
		r.Y = float32(screenHeight) - 130 - r.Height
		r.X = -6000 + spacing

		spacing += r.Width

		c := rl.NewColor(byte(rl.GetRandomValue(200, 240)), byte(rl.GetRandomValue(200, 240)), byte(rl.GetRandomValue(200, 250)), byte(255))

		buildings[i] = r
		buildColors[i] = c
	}

	camera := rl.Camera2D{}
	camera.Target = rl.NewVector2(float32(player.X+20), float32(player.Y+20))
	camera.Offset = rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2))
	camera.Rotation = 0.0
	camera.Zoom = 1.0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if rl.IsKeyDown(rl.KeyRight) {
			player.X += 2 // Player movement
		} else if rl.IsKeyDown(rl.KeyLeft) {
			player.X -= 2 // Player movement
		}

		// Camera target follows player
		camera.Target = rl.NewVector2(float32(player.X+20), float32(player.Y+20))

		// Camera rotation controls
		if rl.IsKeyDown(rl.KeyA) {
			camera.Rotation--
		} else if rl.IsKeyDown(rl.KeyS) {
			camera.Rotation++
		}

		// Limit camera rotation to 80 degrees (-40 to 40)
		if camera.Rotation > 40 {
			camera.Rotation = 40
		} else if camera.Rotation < -40 {
			camera.Rotation = -40
		}

		// Camera zoom controls
		camera.Zoom += float32(rl.GetMouseWheelMove()) * 0.05

		if camera.Zoom > 3.0 {
			camera.Zoom = 3.0
		} else if camera.Zoom < 0.1 {
			camera.Zoom = 0.1
		}

		// Camera reset (zoom and rotation)
		if rl.IsKeyPressed(rl.KeyR) {
			camera.Zoom = 1.0
			camera.Rotation = 0.0
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode2D(camera)

		rl.DrawRectangle(-6000, 320, 13000, 8000, rl.DarkGray)

		for i := 0; i < maxBuildings; i++ {
			rl.DrawRectangleRec(buildings[i], buildColors[i])
		}

		rl.DrawRectangleRec(player, rl.Red)

		rl.DrawRectangle(int32(camera.Target.X), -500, 1, screenHeight*4, rl.Green)
		rl.DrawRectangle(-500, int32(camera.Target.Y), screenWidth*4, 1, rl.Green)

		rl.EndMode2D()

		rl.DrawText("SCREEN AREA", 640, 10, 20, rl.Red)

		rl.DrawRectangle(0, 0, screenWidth, 5, rl.Red)
		rl.DrawRectangle(0, 5, 5, screenHeight-10, rl.Red)
		rl.DrawRectangle(screenWidth-5, 5, 5, screenHeight-10, rl.Red)
		rl.DrawRectangle(0, screenHeight-5, screenWidth, 5, rl.Red)

		rl.DrawRectangle(10, 10, 250, 113, rl.Fade(rl.SkyBlue, 0.5))
		rl.DrawRectangleLines(10, 10, 250, 113, rl.Blue)

		rl.DrawText("Free 2d camera controls:", 20, 20, 10, rl.Black)
		rl.DrawText("- Right/Left to move Offset", 40, 40, 10, rl.DarkGray)
		rl.DrawText("- Mouse Wheel to Zoom in-out", 40, 60, 10, rl.DarkGray)
		rl.DrawText("- A / S to Rotate", 40, 80, 10, rl.DarkGray)
		rl.DrawText("- R to reset Zoom and Rotation", 40, 100, 10, rl.DarkGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
