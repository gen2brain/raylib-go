package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"log"
)

func main() {
	rl.SetTraceLogCallback(func(logType int, str string) {
		level := ""
		switch logType {
		case rl.LogDebug:
			level = "Debug"
		case rl.LogError:
			level = "Error"
		case rl.LogInfo:
			level = "Info"
		case rl.LogTrace:
			level = "Trace"
		case rl.LogWarning:
			level = "Warning"
		case rl.LogFatal:
			level = "Fatal"
		}
		if logType != rl.LogFatal {
			log.Printf("%s - %s", level, str)
		} else {
			log.Fatalf("%s - %s", level, str)
		}
	})

	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(800, 450, "raylib [utils] example - SetTraceLogCallback")

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("The raylib trace log is controlled in GO!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
