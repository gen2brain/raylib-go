/*******************************************************************************************
*
*   raylib [core] example - Storage save/load values
*
*   Example originally created with raylib 1.4, last time updated with raylib 4.2
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2015-2024 Ramon Santamaria (@raysan5)
*
********************************************************************************************/
package main

import (
	"encoding/binary"
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth     = 800
	screenHeight    = 450
	storageDataFile = "storage.data"
	// NOTE: Storage positions must start with 0, directly related to file memory layout
	storagePositionScore   = 0
	storagePositionHiscore = 1
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - storage save/load values")

	var score int32
	var hiscore int32
	var framesCounter int32

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		if rl.IsKeyPressed(rl.KeyR) {
			score = rl.GetRandomValue(1000, 2000)
			hiscore = rl.GetRandomValue(2000, 4000)
		}

		if rl.IsKeyPressed(rl.KeyEnter) {
			_ = SaveStorageValue(storagePositionScore, score)
			_ = SaveStorageValue(storagePositionHiscore, hiscore)
		} else if rl.IsKeyPressed(rl.KeySpace) {
			// NOTE: If requested position could not be found, value 0 is returned
			score, _ = LoadStorageValue(storagePositionScore)
			hiscore, _ = LoadStorageValue(storagePositionHiscore)
		}

		framesCounter++

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawText(fmt.Sprintf("SCORE: %d", score), 280, 130, 40, rl.Maroon)
		rl.DrawText(fmt.Sprintf("HI-SCORE: %d", hiscore), 210, 200, 50, rl.Black)

		rl.DrawText(fmt.Sprintf("frames: %d", framesCounter), 10, 10, 20, rl.Lime)

		rl.DrawText("Press R to generate random numbers", 220, 40, 20, rl.LightGray)
		rl.DrawText("Press ENTER to SAVE values", 250, 310, 20, rl.LightGray)
		rl.DrawText("Press SPACE to LOAD values", 252, 350, 20, rl.LightGray)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.CloseWindow() // Close window and OpenGL context
}

// SaveStorageValue saves an integer value to storage file (to defined position)
// NOTE: Storage positions is directly related to file memory layout (4 bytes each integer)
func SaveStorageValue(position uint32, value int32) error {
	// Load the file data
	fileData, err := os.ReadFile(storageDataFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load file: %v", err)
	}

	dataSize := len(fileData)
	requiredSize := int((position + 1) * 4) // Each int32 is 4 bytes
	newFileData := make([]byte, requiredSize)

	// Copy existing file data to newFileData
	if dataSize > 0 {
		copy(newFileData, fileData)
	}

	// Update the value at the specified position
	binary.LittleEndian.PutUint32(newFileData[position*4:], uint32(value))

	// Save the updated data back to the file
	err = os.WriteFile(storageDataFile, newFileData, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

// LoadStorageValue loads an integer value from storage file (from defined position)
// NOTE: If requested position could not be found, value 0 is returned
func LoadStorageValue(position uint32) (int32, error) {
	// Load the file data
	fileData, err := os.ReadFile(storageDataFile)
	if err != nil {
		return 0, fmt.Errorf("failed to load file: %v", err)
	}

	dataSize := len(fileData)
	if dataSize < int((position+1)*4) {
		return 0, fmt.Errorf("position out of bounds")
	}

	// Read the value at the specified position
	value := int32(binary.LittleEndian.Uint32(fileData[position*4:]))
	return value, nil
}
