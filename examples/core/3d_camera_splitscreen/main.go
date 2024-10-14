package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	screenW    = int32(800)
	screenH    = int32(450)
	cam1, cam2 rl.Camera3D
)

func main() {

	rl.InitWindow(screenW, screenH, "raylib [core] example - 3d camera split screen")

	// Setup player 1 camera and screen
	cam1.Fovy = 45
	cam1.Up.Y = 1
	cam1.Target.Y = 1
	cam1.Position.Z = -3
	cam1.Position.Y = 1

	// Setup player two camera and screen
	cam2.Fovy = 45
	cam2.Up.Y = 1
	cam2.Target.Y = 3
	cam2.Position.X = -3
	cam2.Position.Y = 3

	screenCam1 := rl.LoadRenderTexture(screenW/2, screenH)
	screenCam2 := rl.LoadRenderTexture(screenW/2, screenH)

	splitScreenRec := rl.NewRectangle(0, 0, float32(screenCam1.Texture.Width), -float32(screenCam1.Texture.Height))

	count := float32(5)
	spacing := float32(4)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// If anyone moves this frame, how far will they move based on the time
		// since the last frame this moves things at 10 world units per second,
		// regardless of the actual FPS
		frameOffset := 10 * rl.GetFrameTime()

		// Move Player1 forward and backwards (no turning)
		if rl.IsKeyDown(rl.KeyW) {
			cam1.Position.Z += frameOffset
			cam1.Target.Z += frameOffset
		} else if rl.IsKeyDown(rl.KeyS) {
			cam1.Position.Z -= frameOffset
			cam1.Target.Z -= frameOffset
		}

		// Move Player2 forward and backwards (no turning)
		if rl.IsKeyDown(rl.KeyUp) {
			cam2.Position.X += frameOffset
			cam2.Target.X += frameOffset
		} else if rl.IsKeyDown(rl.KeyDown) {
			cam2.Position.X -= frameOffset
			cam2.Target.X -= frameOffset
		}

		// Draw Player1 view to the render texture
		rl.BeginTextureMode(screenCam1)
		rl.ClearBackground(rl.SkyBlue)
		rl.BeginMode3D(cam1)

		rl.DrawPlane(rl.Vector3Zero(), rl.NewVector2(50, 50), rl.Beige)

		for x := -count * spacing; x <= count*spacing; x += spacing {
			for z := -count * spacing; z <= count*spacing; z += spacing {
				rl.DrawCube(rl.NewVector3(x-0.5, 1.5, z), 1, 1, 1, rl.Lime)
				rl.DrawCube(rl.NewVector3(x-0.5, 0.5, z), 0.25, 1, 0.25, rl.Brown)
			}
		}

		rl.DrawCube(cam1.Position, 1, 1, 1, rl.Red)
		rl.DrawCube(cam2.Position, 1, 1, 1, rl.Blue)

		rl.EndMode3D()

		rl.DrawRectangle(0, 0, screenW/2, 40, rl.Fade(rl.RayWhite, 0.8))
		rl.DrawText("PLAYER1: W/S to move", 10, 10, 20, rl.Maroon)
		rl.EndTextureMode()

		// Draw Player2 view to the render texture
		rl.BeginTextureMode(screenCam2)
		rl.ClearBackground(rl.SkyBlue)
		rl.BeginMode3D(cam2)

		rl.DrawPlane(rl.Vector3Zero(), rl.NewVector2(50, 50), rl.Beige)

		for x := -count * spacing; x <= count*spacing; x += spacing {
			for z := -count * spacing; z <= count*spacing; z += spacing {
				rl.DrawCube(rl.NewVector3(x, 1.5, z), 1, 1, 1, rl.Lime)
				rl.DrawCube(rl.NewVector3(x, 0.5, z), 0.25, 1, 0.25, rl.Brown)
			}
		}

		rl.DrawCube(cam1.Position, 1, 1, 1, rl.Red)
		rl.DrawCube(cam2.Position, 1, 1, 1, rl.Blue)

		rl.EndMode3D()

		rl.DrawRectangle(0, 0, screenW/2, 40, rl.Fade(rl.RayWhite, 0.8))
		rl.DrawText("PLAYER2: UP/DOWN to move", 10, 10, 20, rl.Maroon)
		rl.EndTextureMode()

		// Draw both views render textures to the screen side by side
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.DrawTextureRec(screenCam1.Texture, splitScreenRec, rl.NewVector2(0, 0), rl.White)
		rl.DrawTextureRec(screenCam2.Texture, splitScreenRec, rl.NewVector2(float32(screenW/2), 0), rl.White)
		rl.DrawRectangle((screenW/2)-2, 0, 4, screenH, rl.LightGray)

		rl.EndDrawing()

	}

	rl.UnloadRenderTexture(screenCam1)
	rl.UnloadRenderTexture(screenCam2)

	rl.CloseWindow()
}
