package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	//go:embed assets/*
	assets embed.FS
)

// Gopher type
type Gopher struct {
	Position rl.Vector2
	Speed    rl.Vector2
	Color    rl.Color
}

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(960)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - Go Embed")

	texture := rl.LoadTexture("res://assets/gopher_pixel.png")

	gophers := make([]*Gopher, 0)
	gophersCount := 0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			// Create more gophers
			for i := 0; i < 100; i++ {
				b := &Gopher{}
				b.Position = rl.GetMousePosition()
				b.Speed.X = float32(rl.GetRandomValue(250, 500)) / 60.0
				b.Speed.Y = float32(rl.GetRandomValue(250, 500)-500) / 60.0

				gophers = append(gophers, b)
				gophersCount++
			}
		}

		// Update gophers
		for _, b := range gophers {
			b.Position.X += b.Speed.X
			b.Position.Y += b.Speed.Y

			if (b.Position.X > float32(screenWidth)) || (b.Position.X < 0) {
				b.Speed.X *= -1
			}

			if (b.Position.Y > float32(screenHeight)) || (b.Position.Y < 0) {
				b.Speed.Y *= -1
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		for _, b := range gophers {
			// NOTE: When internal QUADS batch limit is reached, a draw call is launched and
			// batching buffer starts being filled again; before launching the draw call,
			// updated vertex data from internal buffer is send to GPU... it seems it generates
			// a stall and consequently a frame drop, limiting number of gophers drawn at 60 fps
			rl.DrawTexture(texture, int32(b.Position.X), int32(b.Position.Y), rl.RayWhite)
		}

		rl.DrawRectangle(0, 0, screenWidth, 40, rl.LightGray)
		rl.DrawText("raylib gophermark", 10, 10, 20, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("gophers: %d", gophersCount), 400, 10, 20, rl.Red)

		rl.DrawFPS(260, 10)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)
	rl.CloseWindow()
}

func init() {
	rl.SetLoadFileDataCallback(func(name string) ([]byte, error) {
		if strings.HasPrefix(name, "res://") {
			return assets.ReadFile(name[6:])
		}
		return os.ReadFile(name)
	})
	rl.SetLoadFileTextCallback(func(name string) (string, error) {
		if strings.HasPrefix(name, "res://") {
			b, err := assets.ReadFile(name[6:])
			return string(b), err
		}
		fd, err := os.ReadFile(name)
		if err != nil {
			return "", err
		}
		return string(fd), nil
	})
}
