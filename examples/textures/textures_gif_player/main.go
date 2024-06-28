package main

import (
	"fmt"
	"image/color"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	const MAX_FRAME_DELAY int32 = 20
	const MIN_FRAME_DELAY int32 = 1

	// Initialization
	const screenWidth int32 = 800
	const screenHeight int32 = 450

	rl.InitWindow(screenWidth, screenHeight, "raylib [textures] example - gif playing")

	var animFrames int32 = 0

	// Load all GIF animation frames into a single Image
	// NOTE: GIF data is always loaded as RGBA (32bit) by default
	// NOTE: Frames are just appended one after another in image.data memory
	var imScarfyAnim *rl.Image = rl.LoadImageAnim("scarfy_run.gif", &animFrames)

	// Load texture from image
	// NOTE: We will update this texture when required with next frame data
	// WARNING: It's not recommended to use this technique for sprites animation,
	// use spritesheets instead, like illustrated in textures_sprite_anim example
	var texScarfyAnim rl.Texture2D = rl.LoadTextureFromImage(imScarfyAnim)
	var texScarfyAnimSize int32 = imScarfyAnim.Width * imScarfyAnim.Height

	var nextFrameDataOffset uint32 = 0 // Current byte offset to next frame in image.data
	var currentAnimFrame int32 = 0     // Current animation frame to load and draw
	var frameDelay int32 = 8           // Frame delay to switch between animation frames
	var frameCounter int32 = 0         // General frames counter

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		frameCounter++

		if frameCounter >= frameDelay {
			// Move to next frame
			// NOTE: If final frame is reached we return to first frame
			currentAnimFrame++
			if currentAnimFrame >= animFrames {
				currentAnimFrame = 0
			}

			// Get memory offset position for next frame data in image.data
			nextFrameDataOffset = uint32(imScarfyAnim.Width * imScarfyAnim.Height * int32(4) * currentAnimFrame)
			// Update GPU texture data with next frame image data
			// WARNING: Data size (frame size) and pixel format must match already created texture
			// here we needed to make the Data as public
			rl.UpdateTexture(texScarfyAnim,
				unsafe.Slice((*color.RGBA)(unsafe.Pointer(uintptr(imScarfyAnim.Data)+uintptr(nextFrameDataOffset))), texScarfyAnimSize))

			frameCounter = 0
		}

		// Control frames delay
		if rl.IsKeyPressed(rl.KeyRight) {
			frameDelay++
		} else if rl.IsKeyPressed(rl.KeyLeft) {
			frameDelay--
		}

		if frameDelay > MAX_FRAME_DELAY {
			frameDelay = MAX_FRAME_DELAY
		} else if frameDelay < MIN_FRAME_DELAY {
			frameDelay = MIN_FRAME_DELAY
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText(fmt.Sprintf("TOTAL GIF FRAMES:  %02d", animFrames), 50, 30, 20, rl.LightGray)
		rl.DrawText(fmt.Sprintf("CURRENT FRAME: %02d", currentAnimFrame), 50, 60, 20, rl.Gray)
		rl.DrawText(fmt.Sprintf("CURRENT FRAME IMAGE.DATA OFFSET: %02d", nextFrameDataOffset), 50, 90, 20, rl.Gray)

		rl.DrawText("FRAMES DELAY: ", 100, 305, 10, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("%02d frames", frameDelay), 620, 305, 10, rl.DarkGray)
		rl.DrawText("PRESS RIGHT/LEFT KEYS to CHANGE SPEED!", 290, 350, 10, rl.DarkGray)

		for i := int32(0); i < MAX_FRAME_DELAY; i++ {
			if i < frameDelay {
				rl.DrawRectangle(190+21*i, 300, 20, 20, rl.Red)
			}
			rl.DrawRectangleLines(190+21*i, 300, 20, 20, rl.Maroon)
		}

		rl.DrawTexture(texScarfyAnim, int32(rl.GetScreenWidth()/2)-texScarfyAnim.Width/2, 140, rl.White)

		rl.DrawText("(c) Scarfy sprite by Eiden Marsal", screenWidth-200, screenHeight-20, 10, rl.Gray)

		rl.EndDrawing()
	}

	// De-Initialization
	defer rl.UnloadTexture(texScarfyAnim) // Unload texture
	defer rl.UnloadImage(imScarfyAnim)    // Unload image (contains all frames)

	defer rl.CloseWindow() // Close window and OpenGL context

}
