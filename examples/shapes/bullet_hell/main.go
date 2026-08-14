package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const maxBullets = 500000 // Max bullets to be processed

type Bullet struct {
	Position     rl.Vector2
	Acceleration rl.Vector2
	Disabled     bool
	Color        rl.Color
}

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - bullet hell")
	defer rl.CloseWindow()

	// Bullets definition
	bullets := make([]Bullet, maxBullets)
	bulletCount := 0
	bulletDisabledCount := 0
	bulletRadius := float32(10)
	bulletSpeed := float32(3.0)
	bulletRows := 6
	bulletColor := [2]rl.Color{rl.Red, rl.Blue}

	// Spawner variables
	baseDirection := float32(0)
	angleIncrement := 5
	spawnCooldown := float32(2)
	spawnCooldownTimer := spawnCooldown

	// Magic circle
	magicCircleRotation := float32(0)

	// Used on performance drawing
	bulletTexture := rl.LoadRenderTexture(24, 24)
	defer rl.UnloadRenderTexture(bulletTexture)

	// Pre-render circle to texture for performance mode
	rl.BeginTextureMode(bulletTexture)
	rl.DrawCircle(12, 12, bulletRadius, rl.White)
	rl.DrawCircleLines(12, 12, bulletRadius, rl.Black)
	rl.EndTextureMode()

	drawInPerformanceMode := true

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		// Reset the bullet index when capacity is reached
		if bulletCount >= maxBullets {
			bulletCount = 0
			bulletDisabledCount = 0
		}

		spawnCooldownTimer--
		if spawnCooldownTimer < 0 {
			spawnCooldownTimer = spawnCooldown

			// Spawn bullets
			degreesPerRow := float32(360.0) / float32(bulletRows)
			for row := 0; row < bulletRows; row++ {
				if bulletCount < maxBullets {
					bullets[bulletCount].Position = rl.NewVector2(float32(screenWidth)/2.0, float32(screenHeight)/2.0)
					bullets[bulletCount].Disabled = false
					bullets[bulletCount].Color = bulletColor[row%2]

					bulletDirection := baseDirection + (degreesPerRow * float32(row))

					rad := float64(bulletDirection * rl.Deg2rad)
					bullets[bulletCount].Acceleration = rl.NewVector2(
						bulletSpeed*float32(math.Cos(rad)),
						bulletSpeed*float32(math.Sin(rad)),
					)

					bulletCount++
				}
			}

			baseDirection += float32(angleIncrement)
		}

		// Update bullets position based on acceleration
		for i := 0; i < bulletCount; i++ {
			if !bullets[i].Disabled {
				bullets[i].Position.X += bullets[i].Acceleration.X
				bullets[i].Position.Y += bullets[i].Acceleration.Y

				// Disable bullet if out of screen
				if bullets[i].Position.X < -bulletRadius*2 ||
					bullets[i].Position.X > float32(screenWidth)+bulletRadius*2 ||
					bullets[i].Position.Y < -bulletRadius*2 ||
					bullets[i].Position.Y > float32(screenHeight)+bulletRadius*2 {
					bullets[i].Disabled = true
					bulletDisabledCount++
				}
			}
		}

		// Input logic
		if (rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD)) && bulletRows < 359 {
			bulletRows++
		}
		if (rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA)) && bulletRows > 1 {
			bulletRows--
		}
		if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
			bulletSpeed += 0.25
		}
		if (rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS)) && bulletSpeed > 0.50 {
			bulletSpeed -= 0.25
		}
		if rl.IsKeyPressed(rl.KeyZ) && spawnCooldown > 1 {
			spawnCooldown--
		}
		if rl.IsKeyPressed(rl.KeyX) {
			spawnCooldown++
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			drawInPerformanceMode = !drawInPerformanceMode
		}

		if rl.IsKeyDown(rl.KeySpace) {
			angleIncrement = (angleIncrement + 1) % 360
		}

		if rl.IsKeyPressed(rl.KeyC) {
			bulletCount = 0
			bulletDisabledCount = 0
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		// Draw magic circle
		magicCircleRotation++
		centerRec := rl.NewRectangle(float32(screenWidth)/2.0, float32(screenHeight)/2.0, 120, 120)
		centerOrigin := rl.NewVector2(60.0, 60.0)

		rl.DrawRectanglePro(centerRec, centerOrigin, magicCircleRotation, rl.Purple)
		rl.DrawRectanglePro(centerRec, centerOrigin, magicCircleRotation+45, rl.Purple)
		rl.DrawCircleLines(int32(screenWidth/2), int32(screenHeight/2), 70, rl.Black)
		rl.DrawCircleLines(int32(screenWidth/2), int32(screenHeight/2), 50, rl.Black)
		rl.DrawCircleLines(int32(screenWidth/2), int32(screenHeight/2), 30, rl.Black)

		// Draw bullets
		if drawInPerformanceMode {
			texWidth := float32(bulletTexture.Texture.Width)
			texHeight := float32(bulletTexture.Texture.Height)

			for i := 0; i < bulletCount; i++ {
				if !bullets[i].Disabled {
					rl.DrawTexture(
						bulletTexture.Texture,
						int32(bullets[i].Position.X-texWidth*0.5),
						int32(bullets[i].Position.Y-texHeight*0.5),
						bullets[i].Color,
					)
				}
			}
		} else {
			for i := 0; i < bulletCount; i++ {
				if !bullets[i].Disabled {
					rl.DrawCircleV(bullets[i].Position, bulletRadius, bullets[i].Color)
					// Use DrawCircleLines instead of DrawCircleLinesV to avoid purego FFI bug
					rl.DrawCircleLines(int32(bullets[i].Position.X), int32(bullets[i].Position.Y), bulletRadius, rl.Black)
				}
			}
		}

		// Draw UI
		bgOverlay := rl.NewColor(0, 0, 0, 200)

		rl.DrawRectangle(10, 10, 280, 150, bgOverlay)
		rl.DrawText("Controls:", 20, 20, 10, rl.LightGray)
		rl.DrawText("- Right/Left or A/D: Change rows number", 40, 40, 10, rl.LightGray)
		rl.DrawText("- Up/Down or W/S: Change bullet speed", 40, 60, 10, rl.LightGray)
		rl.DrawText("- Z or X: Change spawn cooldown", 40, 80, 10, rl.LightGray)
		rl.DrawText("- Space (Hold): Change the angle increment", 40, 100, 10, rl.LightGray)
		rl.DrawText("- Enter: Switch draw method (Performance)", 40, 120, 10, rl.LightGray)
		rl.DrawText("- C: Clear bullets", 40, 140, 10, rl.LightGray)

		rl.DrawRectangle(610, 10, 170, 30, bgOverlay)
		if drawInPerformanceMode {
			rl.DrawText("Draw method: DrawTexture(*)", 620, 20, 10, rl.Green)
		} else {
			rl.DrawText("Draw method: DrawCircle(*)", 620, 20, 10, rl.Red)
		}

		rl.DrawRectangle(135, 410, 530, 30, bgOverlay)
		rl.DrawText(
			fmt.Sprintf("[ FPS: %d, Bullets: %d, Rows: %d, Bullet speed: %.2f, Angle increment per frame: %d, Cooldown: %.0f ]",
				rl.GetFPS(), bulletCount-bulletDisabledCount, bulletRows, bulletSpeed, angleIncrement, spawnCooldown),
			155, 420, 10, rl.Green,
		)

		rl.EndDrawing()
	}
}
