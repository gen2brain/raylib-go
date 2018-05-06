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

	raylib.InitWindow(screenWidth, screenHeight, "raylib [core] example - 2d camera")

	player := raylib.NewRectangle(400, 280, 40, 40)

	buildings := make([]raylib.Rectangle, maxBuildings)
	buildColors := make([]raylib.Color, maxBuildings)

	spacing := float32(0)

	for i := 0; i < maxBuildings; i++ {
		r := raylib.Rectangle{}
		r.Width = float32(raylib.GetRandomValue(50, 200))
		r.Height = float32(raylib.GetRandomValue(100, 800))
		r.Y = float32(screenHeight) - 130 - r.Height
		r.X = -6000 + spacing

		spacing += r.Width

		c := raylib.NewColor(byte(raylib.GetRandomValue(200, 240)), byte(raylib.GetRandomValue(200, 240)), byte(raylib.GetRandomValue(200, 250)), byte(255))

		buildings[i] = r
		buildColors[i] = c
	}

	camera := raylib.Camera2D{}
	camera.Target = raylib.NewVector2(float32(player.X+20), float32(player.Y+20))
	camera.Offset = raylib.NewVector2(0, 0)
	camera.Rotation = 0.0
	camera.Zoom = 1.0

	raylib.SetTargetFPS(30)

	for !raylib.WindowShouldClose() {
		if raylib.IsKeyDown(raylib.KeyRight) {
			player.X += 2        // Player movement
			camera.Offset.X -= 2 // Camera displacement with player movement
		} else if raylib.IsKeyDown(raylib.KeyLeft) {
			player.X -= 2        // Player movement
			camera.Offset.X += 2 // Camera displacement with player movement
		}

		// Camera target follows player
		camera.Target = raylib.NewVector2(float32(player.X+20), float32(player.Y+20))

		// Camera rotation controls
		if raylib.IsKeyDown(raylib.KeyA) {
			camera.Rotation--
		} else if raylib.IsKeyDown(raylib.KeyS) {
			camera.Rotation++
		}

		// Limit camera rotation to 80 degrees (-40 to 40)
		if camera.Rotation > 40 {
			camera.Rotation = 40
		} else if camera.Rotation < -40 {
			camera.Rotation = -40
		}

		// Camera zoom controls
		camera.Zoom += float32(raylib.GetMouseWheelMove()) * 0.05

		if camera.Zoom > 3.0 {
			camera.Zoom = 3.0
		} else if camera.Zoom < 0.1 {
			camera.Zoom = 0.1
		}

		// Camera reset (zoom and rotation)
		if raylib.IsKeyPressed(raylib.KeyR) {
			camera.Zoom = 1.0
			camera.Rotation = 0.0
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.BeginMode2D(camera)

		raylib.DrawRectangle(-6000, 320, 13000, 8000, raylib.DarkGray)

		for i := 0; i < maxBuildings; i++ {
			raylib.DrawRectangleRec(buildings[i], buildColors[i])
		}

		raylib.DrawRectangleRec(player, raylib.Red)

		raylib.DrawRectangle(int32(camera.Target.X), -500, 1, screenHeight*4, raylib.Green)
		raylib.DrawRectangle(-500, int32(camera.Target.Y), screenWidth*4, 1, raylib.Green)

		raylib.EndMode2D()

		raylib.DrawText("SCREEN AREA", 640, 10, 20, raylib.Red)

		raylib.DrawRectangle(0, 0, screenWidth, 5, raylib.Red)
		raylib.DrawRectangle(0, 5, 5, screenHeight-10, raylib.Red)
		raylib.DrawRectangle(screenWidth-5, 5, 5, screenHeight-10, raylib.Red)
		raylib.DrawRectangle(0, screenHeight-5, screenWidth, 5, raylib.Red)

		raylib.DrawRectangle(10, 10, 250, 113, raylib.Fade(raylib.SkyBlue, 0.5))
		raylib.DrawRectangleLines(10, 10, 250, 113, raylib.Blue)

		raylib.DrawText("Free 2d camera controls:", 20, 20, 10, raylib.Black)
		raylib.DrawText("- Right/Left to move Offset", 40, 40, 10, raylib.DarkGray)
		raylib.DrawText("- Mouse Wheel to Zoom in-out", 40, 60, 10, raylib.DarkGray)
		raylib.DrawText("- A / S to Rotate", 40, 80, 10, raylib.DarkGray)
		raylib.DrawText("- R to reset Zoom and Rotation", 40, 100, 10, raylib.DarkGray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
