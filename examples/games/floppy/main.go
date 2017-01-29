package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	maxTubes     = 100
	floppyRadius = 24.0
	tubesWidth   = 80
)

type Floppy struct {
	Position raylib.Vector2
	Radius   float32
	Color    raylib.Color
}

type Tubes struct {
	Rec    raylib.Rectangle
	Color  raylib.Color
	Active bool
}

type Game struct {
	ScreenWidth  int32
	ScreenHeight int32

	FramesCounter int32
	GameOver      bool
	Pause         bool
	Score         int
	HiScore       int

	Floppy      Floppy
	Tubes       []Tubes
	TubesPos    []raylib.Vector2
	TubesSpeedX int32
	SuperFX     bool
}

func main() {
	game := Game{}
	game.Init()

	raylib.InitWindow(game.ScreenWidth, game.ScreenHeight, "sample game: floppy")

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		game.Update()

		game.Draw()
	}

	raylib.CloseWindow()
}

// Initialize game
func (g *Game) Init() {
	g.ScreenWidth = 800
	g.ScreenHeight = 450

	g.Floppy = Floppy{}
	g.Floppy.Radius = floppyRadius
	g.Floppy.Position = raylib.NewVector2(80, float32(g.ScreenHeight)/2-g.Floppy.Radius)
	g.TubesSpeedX = 2

	g.TubesPos = make([]raylib.Vector2, maxTubes)

	for i := 0; i < maxTubes; i++ {
		g.TubesPos[i].X = float32(400 + 280*i)
		g.TubesPos[i].Y = -float32(raylib.GetRandomValue(0, 120))
	}

	g.Tubes = make([]Tubes, maxTubes*2)

	for i := 0; i < maxTubes*2; i += 2 {
		g.Tubes[i].Rec.X = int32(g.TubesPos[i/2].X)
		g.Tubes[i].Rec.Y = int32(g.TubesPos[i/2].Y)
		g.Tubes[i].Rec.Width = tubesWidth
		g.Tubes[i].Rec.Height = 255

		g.Tubes[i+1].Rec.X = int32(g.TubesPos[i/2].X)
		g.Tubes[i+1].Rec.Y = int32(600 + g.TubesPos[i/2].Y - 255)
		g.Tubes[i+1].Rec.Width = tubesWidth
		g.Tubes[i+1].Rec.Height = 255

		g.Tubes[i/2].Active = true
	}

	g.Score = 0
	g.FramesCounter = 0

	g.GameOver = false
	g.SuperFX = false
	g.Pause = false
}

// Update game
func (g *Game) Update() {
	if !g.GameOver {
		if raylib.IsKeyPressed(raylib.KeyP) {
			g.Pause = !g.Pause
		}

		if !g.Pause {
			for i := 0; i < maxTubes; i++ {
				g.TubesPos[i].X -= float32(g.TubesSpeedX)
			}

			for i := 0; i < maxTubes*2; i += 2 {
				g.Tubes[i].Rec.X = int32(g.TubesPos[i/2].X)
				g.Tubes[i+1].Rec.X = int32(g.TubesPos[i/2].X)
			}

			if raylib.IsKeyDown(raylib.KeySpace) && !g.GameOver {
				g.Floppy.Position.Y -= 3
			} else {
				g.Floppy.Position.Y += 1
			}

			// Check Collisions
			for i := 0; i < maxTubes*2; i++ {
				if raylib.CheckCollisionCircleRec(g.Floppy.Position, g.Floppy.Radius, g.Tubes[i].Rec) {
					g.GameOver = true
					g.Pause = false
				} else if (g.TubesPos[i/2].X < g.Floppy.Position.X) && g.Tubes[i/2].Active && !g.GameOver {
					g.Score += 100
					g.Tubes[i/2].Active = false

					g.SuperFX = true

					if g.Score > g.HiScore {
						g.HiScore = g.Score
					}
				}
			}
		}
	} else {
		if raylib.IsKeyPressed(raylib.KeyEnter) {
			g.Init()
			g.GameOver = false
		}
	}
}

// Draw game
func (g *Game) Draw() {
	raylib.BeginDrawing()

	raylib.ClearBackground(raylib.RayWhite)

	if !g.GameOver {
		raylib.DrawCircle(int32(g.Floppy.Position.X), int32(g.Floppy.Position.Y), g.Floppy.Radius, raylib.DarkGray)

		// Draw tubes
		for i := 0; i < maxTubes; i++ {
			raylib.DrawRectangle(g.Tubes[i*2].Rec.X, g.Tubes[i*2].Rec.Y, g.Tubes[i*2].Rec.Width, g.Tubes[i*2].Rec.Height, raylib.Gray)
			raylib.DrawRectangle(g.Tubes[i*2+1].Rec.X, g.Tubes[i*2+1].Rec.Y, g.Tubes[i*2+1].Rec.Width, g.Tubes[i*2+1].Rec.Height, raylib.Gray)
		}

		// Draw flashing fx (one frame only)
		if g.SuperFX {
			raylib.DrawRectangle(0, 0, g.ScreenWidth, g.ScreenHeight, raylib.White)
			g.SuperFX = false
		}

		raylib.DrawText(fmt.Sprintf("%04d", g.Score), 20, 20, 40, raylib.Gray)
		raylib.DrawText(fmt.Sprintf("HI-SCORE: %04d", g.HiScore), 20, 70, 20, raylib.LightGray)

		if g.Pause {
			raylib.DrawText("GAME PAUSED", g.ScreenWidth/2-raylib.MeasureText("GAME PAUSED", 40)/2, g.ScreenHeight/2-40, 40, raylib.Gray)
		}
	} else {
		raylib.DrawText("PRESS [ENTER] TO PLAY AGAIN", raylib.GetScreenWidth()/2-raylib.MeasureText("PRESS [ENTER] TO PLAY AGAIN", 20)/2, raylib.GetScreenHeight()/2-50, 20, raylib.Gray)
	}

	raylib.EndDrawing()
}
