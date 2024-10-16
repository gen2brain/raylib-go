package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

func main() {
	// Set custom logger
	rl.SetTraceLogCallback(customLog)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - Custom Logging")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Check the console output to see the custom logger in action!", 60, 200, 20, rl.LightGray)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

// Custom logging function
func customLog(msgType int, text string) {
	now := time.Now()

	switch rl.TraceLogLevel(msgType) {
	case rl.LogInfo:
		fmt.Println("[INFO] : ", now, text)
	case rl.LogError:
		fmt.Println("[ERROR] : ", now, text)
	case rl.LogWarning:
		fmt.Println("[WARNING] : ", now, text)
	case rl.LogDebug:
		fmt.Println("[DEBUG] : ", now, text)
	default:
		fmt.Println("[UNKNOWN] : ", now, text)
	}
}
