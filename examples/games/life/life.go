package main

import (
	"math/rand"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	squareSize = 8
)

// Cell type
type Cell struct {
	Position raylib.Vector2
	Size     raylib.Vector2
	Alive    bool
	Next     bool
	Visited  bool
}

// Game type
type Game struct {
	ScreenWidth   int32
	ScreenHeight  int32
	FramesCounter int32
	Playing       bool
	Cells         [][]*Cell
}

func main() {
	game := Game{}
	game.Init()
	raylib.InitWindow(game.ScreenWidth, game.ScreenHeight, "Conway's Game of Life")
	raylib.SetTargetFPS(20)
	for !raylib.WindowShouldClose() {
		if game.Playing {
			game.Update()
		}
		game.Input()

		game.Draw()
	}
	raylib.CloseWindow()
}

// Init - Initialize game
func (g *Game) Init() {
	g.ScreenWidth = 1024
	g.ScreenHeight = 768
	g.FramesCounter = 0

	g.Cells = make([][]*Cell, g.ScreenWidth/squareSize+1)
	for i := int32(0); i <= g.ScreenWidth/squareSize; i++ {
		g.Cells[i] = make([]*Cell, g.ScreenHeight/squareSize+1)
	}

	for x := int32(0); x <= g.ScreenWidth/squareSize; x++ {
		for y := int32(0); y <= g.ScreenHeight/squareSize; y++ {
			g.Cells[x][y] = &Cell{}
			g.Cells[x][y].Position = raylib.NewVector2((float32(x) * squareSize), (float32(y)*squareSize)+1)
			g.Cells[x][y].Size = raylib.NewVector2(squareSize-1, squareSize-1)
			rand.Seed(time.Now().UnixNano())
			if rand.Float64() < 0.1 {
				g.Cells[x][y].Alive = true
			}
		}
	}
}

// Input - Game input
func (g *Game) Input() {
	// control
	if raylib.IsKeyPressed(raylib.KeyR) {
		g.Init()
	}
	if raylib.IsKeyDown(raylib.KeyRight) && !g.Playing {
		g.Update()
	}
	if raylib.IsKeyPressed(raylib.KeySpace) {
		g.Playing = !g.Playing
	}

	g.FramesCounter++
}

// Update - Update game
func (g *Game) Update() {
	for i := int32(0); i <= g.ScreenWidth/squareSize; i++ {
		for j := int32(0); j <= g.ScreenHeight/squareSize; j++ {
			NeighbourCount := 0
			if j-1 >= 0 {
				if g.Cells[i][j-1].Alive {
					NeighbourCount++
				}
			}
			if j+1 <= g.ScreenHeight/squareSize {
				if g.Cells[i][j+1].Alive {
					NeighbourCount++
				}
			}
			if i-1 >= 0 {
				if g.Cells[i-1][j].Alive {
					NeighbourCount++
				}
			}
			if i+1 <= g.ScreenWidth/squareSize {
				if g.Cells[i+1][j].Alive {
					NeighbourCount++
				}
			}
			if i-1 >= 0 && j-1 >= 0 {
				if g.Cells[i-1][j-1].Alive {
					NeighbourCount++
				}
			}
			if i-1 >= 0 && j+1 <= g.ScreenHeight/squareSize {
				if g.Cells[i-1][j+1].Alive {
					NeighbourCount++
				}
			}
			if i+1 <= g.ScreenWidth/squareSize && j-1 >= 0 {
				if g.Cells[i+1][j-1].Alive {
					NeighbourCount++
				}
			}
			if i+1 <= g.ScreenWidth/squareSize && j+1 <= g.ScreenHeight/squareSize {
				if g.Cells[i+1][j+1].Alive {
					NeighbourCount++
				}
			}
			if g.Cells[i][j].Alive {
				if NeighbourCount < 2 {
					g.Cells[i][j].Next = false
				} else if NeighbourCount > 3 {
					g.Cells[i][j].Next = false
				} else {
					g.Cells[i][j].Next = true
				}
			} else {
				if NeighbourCount == 3 {
					g.Cells[i][j].Next = true
					g.Cells[i][j].Visited = true
				}
			}
		}
	}
	for i := int32(0); i <= g.ScreenWidth/squareSize; i++ {
		for j := int32(0); j < g.ScreenHeight/squareSize; j++ {
			g.Cells[i][j].Alive = g.Cells[i][j].Next
		}
	}
}

// Draw - Draw game
func (g *Game) Draw() {
	raylib.BeginDrawing()
	raylib.ClearBackground(raylib.RayWhite)

	// Draw cells
	for x := int32(0); x <= g.ScreenWidth/squareSize; x++ {
		for y := int32(0); y <= g.ScreenHeight/squareSize; y++ {
			if g.Cells[x][y].Alive {
				raylib.DrawRectangleV(g.Cells[x][y].Position, g.Cells[x][y].Size, raylib.Blue)
			} else if g.Cells[x][y].Visited {
				raylib.DrawRectangleV(g.Cells[x][y].Position, g.Cells[x][y].Size, raylib.Color{R: 128, G: 177, B: 136, A: 255})
			}
		}
	}

	// Draw grid lines
	for i := int32(0); i < g.ScreenWidth/squareSize+1; i++ {
		raylib.DrawLineV(
			raylib.NewVector2(float32(squareSize*i), 0),
			raylib.NewVector2(float32(squareSize*i), float32(g.ScreenHeight)),
			raylib.LightGray,
		)
	}

	for i := int32(0); i < g.ScreenHeight/squareSize+1; i++ {
		raylib.DrawLineV(
			raylib.NewVector2(0, float32(squareSize*i)),
			raylib.NewVector2(float32(g.ScreenWidth), float32(squareSize*i)),
			raylib.LightGray,
		)
	}

	raylib.EndDrawing()
}
