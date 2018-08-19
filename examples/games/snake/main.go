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
	Position raylib.Vector2
	Size     raylib.Vector2
	Speed    raylib.Vector2
	Color    raylib.Color
}

// Food type
type Food struct {
	Position raylib.Vector2
	Size     raylib.Vector2
	Active   bool
	Color    raylib.Color
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
	SnakePosition []raylib.Vector2
	AllowMove     bool
	Offset        raylib.Vector2
	CounterTail   int
}

func main() {
	game := Game{}
	game.Init()

	raylib.InitWindow(game.ScreenWidth, game.ScreenHeight, "sample game: snake")

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		game.Update()

		game.Draw()
	}

	raylib.CloseWindow()
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

	g.Offset = raylib.Vector2{}
	g.Offset.X = float32(g.ScreenWidth % squareSize)
	g.Offset.Y = float32(g.ScreenHeight % squareSize)

	g.Snake = make([]Snake, snakeLength)

	for i := 0; i < snakeLength; i++ {
		g.Snake[i].Position = raylib.NewVector2(g.Offset.X/2, g.Offset.Y/2)
		g.Snake[i].Size = raylib.NewVector2(squareSize, squareSize)
		g.Snake[i].Speed = raylib.NewVector2(squareSize, 0)

		if i == 0 {
			g.Snake[i].Color = raylib.DarkBlue
		} else {
			g.Snake[i].Color = raylib.Blue
		}
	}

	g.SnakePosition = make([]raylib.Vector2, snakeLength)

	for i := 0; i < snakeLength; i++ {
		g.SnakePosition[i] = raylib.NewVector2(0.0, 0.0)
	}

	g.Fruit.Size = raylib.NewVector2(squareSize, squareSize)
	g.Fruit.Color = raylib.SkyBlue
	g.Fruit.Active = false
}

// Update - Update game
func (g *Game) Update() {
	if !g.GameOver {
		if raylib.IsKeyPressed(raylib.KeyP) {
			g.Pause = !g.Pause
		}

		if !g.Pause {
			// control
			if raylib.IsKeyPressed(raylib.KeyRight) && g.Snake[0].Speed.X == 0 && g.AllowMove {
				g.Snake[0].Speed = raylib.NewVector2(squareSize, 0)
				g.AllowMove = false
			}
			if raylib.IsKeyPressed(raylib.KeyLeft) && g.Snake[0].Speed.X == 0 && g.AllowMove {
				g.Snake[0].Speed = raylib.NewVector2(-squareSize, 0)
				g.AllowMove = false
			}
			if raylib.IsKeyPressed(raylib.KeyUp) && g.Snake[0].Speed.Y == 0 && g.AllowMove {
				g.Snake[0].Speed = raylib.NewVector2(0, -squareSize)
				g.AllowMove = false
			}
			if raylib.IsKeyPressed(raylib.KeyDown) && g.Snake[0].Speed.Y == 0 && g.AllowMove {
				g.Snake[0].Speed = raylib.NewVector2(0, squareSize)
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
				g.Fruit.Position = raylib.NewVector2(
					float32(raylib.GetRandomValue(0, (g.ScreenWidth/squareSize)-1)*squareSize+int32(g.Offset.X)/2),
					float32(raylib.GetRandomValue(0, (g.ScreenHeight/squareSize)-1)*squareSize+int32(g.Offset.Y)/2),
				)

				for i := 0; i < g.CounterTail; i++ {
					for (g.Fruit.Position.X == g.Snake[i].Position.X) && (g.Fruit.Position.Y == g.Snake[i].Position.Y) {
						g.Fruit.Position = raylib.NewVector2(
							float32(raylib.GetRandomValue(0, (g.ScreenWidth/squareSize)-1)*squareSize),
							float32(raylib.GetRandomValue(0, (g.ScreenHeight/squareSize)-1)*squareSize),
						)
						i = 0
					}
				}
			}

			// collision
			if raylib.CheckCollisionRecs(
				raylib.NewRectangle(g.Snake[0].Position.X, g.Snake[0].Position.Y, g.Snake[0].Size.X, g.Snake[0].Size.Y),
				raylib.NewRectangle(g.Fruit.Position.X, g.Fruit.Position.Y, g.Fruit.Size.X, g.Fruit.Size.Y),
			) {
				g.Snake[g.CounterTail].Position = g.SnakePosition[g.CounterTail-1]
				g.CounterTail += 1
				g.Fruit.Active = false
			}

			g.FramesCounter++
		}
	} else {
		if raylib.IsKeyPressed(raylib.KeyEnter) {
			g.Init()
			g.GameOver = false
		}
	}
}

// Draw - Draw game
func (g *Game) Draw() {
	raylib.BeginDrawing()

	raylib.ClearBackground(raylib.RayWhite)

	if !g.GameOver {
		// Draw grid lines
		for i := int32(0); i < g.ScreenWidth/squareSize+1; i++ {
			raylib.DrawLineV(
				raylib.NewVector2(float32(squareSize*i)+g.Offset.X/2, g.Offset.Y/2),
				raylib.NewVector2(float32(squareSize*i)+g.Offset.X/2, float32(g.ScreenHeight)-g.Offset.Y/2),
				raylib.LightGray,
			)
		}

		for i := int32(0); i < g.ScreenHeight/squareSize+1; i++ {
			raylib.DrawLineV(
				raylib.NewVector2(g.Offset.X/2, float32(squareSize*i)+g.Offset.Y/2),
				raylib.NewVector2(float32(g.ScreenWidth)-g.Offset.X/2, float32(squareSize*i)+g.Offset.Y/2),
				raylib.LightGray,
			)
		}

		// Draw snake
		for i := 0; i < g.CounterTail; i++ {
			raylib.DrawRectangleV(g.Snake[i].Position, g.Snake[i].Size, g.Snake[i].Color)
		}

		// Draw fruit to pick
		raylib.DrawRectangleV(g.Fruit.Position, g.Fruit.Size, g.Fruit.Color)

		if g.Pause {
			raylib.DrawText("GAME PAUSED", g.ScreenWidth/2-raylib.MeasureText("GAME PAUSED", 40)/2, g.ScreenHeight/2-40, 40, raylib.Gray)
		}
	} else {
		raylib.DrawText("PRESS [ENTER] TO PLAY AGAIN", raylib.GetScreenWidth()/2-raylib.MeasureText("PRESS [ENTER] TO PLAY AGAIN", 20)/2, raylib.GetScreenHeight()/2-50, 20, raylib.Gray)
	}

	raylib.EndDrawing()
}
