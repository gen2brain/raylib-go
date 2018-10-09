package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const (
	snakeLength = 256
	squareSize  = 31
)

// Snake type
type Snake struct {
	Position rl.Vector2
	Size     rl.Vector2
	Speed    rl.Vector2
	Color    rl.Color
}

// Food type
type Food struct {
	Position rl.Vector2
	Size     rl.Vector2
	Active   bool
	Color    rl.Color
}

// Game type
type Game struct {
	ScreenWidth  int32
	ScreenHeight int32

	FramesCounter int32
	GameOver      bool
	Pause         bool

	Fruit         Food
	Snake         []Snake
	SnakePosition []rl.Vector2
	AllowMove     bool
	Offset        rl.Vector2
	CounterTail   int
}

func main() {
	game := Game{}
	game.Init()

	rl.InitWindow(game.ScreenWidth, game.ScreenHeight, "sample game: snake")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		game.Update()

		game.Draw()
	}

	rl.CloseWindow()
}

// Init - Initialize game
func (g *Game) Init() {
	g.ScreenWidth = 800
	g.ScreenHeight = 450

	g.FramesCounter = 0
	g.GameOver = false
	g.Pause = false

	g.CounterTail = 1
	g.AllowMove = false

	g.Offset = rl.Vector2{}
	g.Offset.X = float32(g.ScreenWidth % squareSize)
	g.Offset.Y = float32(g.ScreenHeight % squareSize)

	g.Snake = make([]Snake, snakeLength)

	for i := 0; i < snakeLength; i++ {
		g.Snake[i].Position = rl.NewVector2(g.Offset.X/2, g.Offset.Y/2)
		g.Snake[i].Size = rl.NewVector2(squareSize, squareSize)
		g.Snake[i].Speed = rl.NewVector2(squareSize, 0)

		if i == 0 {
			g.Snake[i].Color = rl.DarkBlue
		} else {
			g.Snake[i].Color = rl.Blue
		}
	}

	g.SnakePosition = make([]rl.Vector2, snakeLength)

	for i := 0; i < snakeLength; i++ {
		g.SnakePosition[i] = rl.NewVector2(0.0, 0.0)
	}

	g.Fruit.Size = rl.NewVector2(squareSize, squareSize)
	g.Fruit.Color = rl.SkyBlue
	g.Fruit.Active = false
}

// Update - Update game
func (g *Game) Update() {
	if !g.GameOver {
		if rl.IsKeyPressed(rl.KeyP) {
			g.Pause = !g.Pause
		}

		if !g.Pause {
			// control
			if rl.IsKeyPressed(rl.KeyRight) && g.Snake[0].Speed.X == 0 && g.AllowMove {
				g.Snake[0].Speed = rl.NewVector2(squareSize, 0)
				g.AllowMove = false
			}
			if rl.IsKeyPressed(rl.KeyLeft) && g.Snake[0].Speed.X == 0 && g.AllowMove {
				g.Snake[0].Speed = rl.NewVector2(-squareSize, 0)
				g.AllowMove = false
			}
			if rl.IsKeyPressed(rl.KeyUp) && g.Snake[0].Speed.Y == 0 && g.AllowMove {
				g.Snake[0].Speed = rl.NewVector2(0, -squareSize)
				g.AllowMove = false
			}
			if rl.IsKeyPressed(rl.KeyDown) && g.Snake[0].Speed.Y == 0 && g.AllowMove {
				g.Snake[0].Speed = rl.NewVector2(0, squareSize)
				g.AllowMove = false
			}

			// movement
			for i := 0; i < g.CounterTail; i++ {
				g.SnakePosition[i] = g.Snake[i].Position
			}

			if g.FramesCounter%5 == 0 {
				for i := 0; i < g.CounterTail; i++ {
					if i == 0 {
						g.Snake[0].Position.X += g.Snake[0].Speed.X
						g.Snake[0].Position.Y += g.Snake[0].Speed.Y
						g.AllowMove = true
					} else {
						g.Snake[i].Position = g.SnakePosition[i-1]
					}
				}
			}

			// wall behaviour
			if ((g.Snake[0].Position.X) > (float32(g.ScreenWidth) - g.Offset.X)) ||
				((g.Snake[0].Position.Y) > (float32(g.ScreenHeight) - g.Offset.Y)) ||
				(g.Snake[0].Position.X < 0) || (g.Snake[0].Position.Y < 0) {
				g.GameOver = true
			}

			// collision with yourself
			for i := 1; i < g.CounterTail; i++ {
				if (g.Snake[0].Position.X == g.Snake[i].Position.X) && (g.Snake[0].Position.Y == g.Snake[i].Position.Y) {
					g.GameOver = true
				}
			}

			if !g.Fruit.Active {
				g.Fruit.Active = true
				g.Fruit.Position = rl.NewVector2(
					float32(rl.GetRandomValue(0, (g.ScreenWidth/squareSize)-1)*squareSize+int32(g.Offset.X)/2),
					float32(rl.GetRandomValue(0, (g.ScreenHeight/squareSize)-1)*squareSize+int32(g.Offset.Y)/2),
				)

				for i := 0; i < g.CounterTail; i++ {
					for (g.Fruit.Position.X == g.Snake[i].Position.X) && (g.Fruit.Position.Y == g.Snake[i].Position.Y) {
						g.Fruit.Position = rl.NewVector2(
							float32(rl.GetRandomValue(0, (g.ScreenWidth/squareSize)-1)*squareSize),
							float32(rl.GetRandomValue(0, (g.ScreenHeight/squareSize)-1)*squareSize),
						)
						i = 0
					}
				}
			}

			// collision
			if rl.CheckCollisionRecs(
				rl.NewRectangle(g.Snake[0].Position.X, g.Snake[0].Position.Y, g.Snake[0].Size.X, g.Snake[0].Size.Y),
				rl.NewRectangle(g.Fruit.Position.X, g.Fruit.Position.Y, g.Fruit.Size.X, g.Fruit.Size.Y),
			) {
				g.Snake[g.CounterTail].Position = g.SnakePosition[g.CounterTail-1]
				g.CounterTail += 1
				g.Fruit.Active = false
			}

			g.FramesCounter++
		}
	} else {
		if rl.IsKeyPressed(rl.KeyEnter) {
			g.Init()
			g.GameOver = false
		}
	}
}

// Draw - Draw game
func (g *Game) Draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	if !g.GameOver {
		// Draw grid lines
		for i := int32(0); i < g.ScreenWidth/squareSize+1; i++ {
			rl.DrawLineV(
				rl.NewVector2(float32(squareSize*i)+g.Offset.X/2, g.Offset.Y/2),
				rl.NewVector2(float32(squareSize*i)+g.Offset.X/2, float32(g.ScreenHeight)-g.Offset.Y/2),
				rl.LightGray,
			)
		}

		for i := int32(0); i < g.ScreenHeight/squareSize+1; i++ {
			rl.DrawLineV(
				rl.NewVector2(g.Offset.X/2, float32(squareSize*i)+g.Offset.Y/2),
				rl.NewVector2(float32(g.ScreenWidth)-g.Offset.X/2, float32(squareSize*i)+g.Offset.Y/2),
				rl.LightGray,
			)
		}

		// Draw snake
		for i := 0; i < g.CounterTail; i++ {
			rl.DrawRectangleV(g.Snake[i].Position, g.Snake[i].Size, g.Snake[i].Color)
		}

		// Draw fruit to pick
		rl.DrawRectangleV(g.Fruit.Position, g.Fruit.Size, g.Fruit.Color)

		if g.Pause {
			rl.DrawText("GAME PAUSED", g.ScreenWidth/2-rl.MeasureText("GAME PAUSED", 40)/2, g.ScreenHeight/2-40, 40, rl.Gray)
		}
	} else {
		rl.DrawText("PRESS [ENTER] TO PLAY AGAIN", int32(rl.GetScreenWidth())/2-rl.MeasureText("PRESS [ENTER] TO PLAY AGAIN", 20)/2, int32(rl.GetScreenHeight())/2-50, 20, rl.Gray)
	}

	rl.EndDrawing()
}
