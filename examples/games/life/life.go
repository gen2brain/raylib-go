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
	Cols          int32
	Rows          int32
	FramesCounter int32
	Playing       bool
	Cells         [][]*Cell
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := Game{}
	game.Init(false)

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
func (g *Game) Init(clear bool) {
	g.ScreenWidth = 1024
	g.ScreenHeight = 768
	g.FramesCounter = 0

	g.Cols = g.ScreenWidth / squareSize
	g.Rows = g.ScreenHeight / squareSize

	g.Cells = make([][]*Cell, g.Cols+1)
	for i := int32(0); i <= g.Cols; i++ {
		g.Cells[i] = make([]*Cell, g.Rows+1)
	}

	for x := int32(0); x <= g.Cols; x++ {
		for y := int32(0); y <= g.Rows; y++ {
			g.Cells[x][y] = &Cell{}
			g.Cells[x][y].Position = raylib.NewVector2((float32(x) * squareSize), (float32(y)*squareSize)+1)
			g.Cells[x][y].Size = raylib.NewVector2(squareSize-1, squareSize-1)
			if rand.Float64() < 0.1 && clear == false {
				g.Cells[x][y].Alive = true
			}
		}
	}
}

// Input - Game input
func (g *Game) Input() {
	// control
	if raylib.IsKeyPressed(raylib.KeyR) {
		g.Init(false)
	}
	if raylib.IsKeyPressed(raylib.KeyC) {
		g.Init(true)
	}
	if raylib.IsKeyDown(raylib.KeyRight) && !g.Playing {
		g.Update()
	}
	if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
		g.Click(raylib.GetMouseX(), raylib.GetMouseY())
	}
	if raylib.IsKeyPressed(raylib.KeySpace) {
		g.Playing = !g.Playing
	}

	g.FramesCounter++
}

// Click - Toggle if a cell is alive or dead on click
func (g *Game) Click(x, y int32) {
	for i := int32(0); i <= g.Cols; i++ {
		for j := int32(0); j <= g.Rows; j++ {
			cell := g.Cells[i][j].Position
			if int32(cell.X) < x && int32(cell.X)+squareSize > x && int32(cell.Y) < y && int32(cell.Y)+squareSize > y {
				g.Cells[i][j].Alive = !g.Cells[i][j].Alive
				g.Cells[i][j].Next = g.Cells[i][j].Alive
			}
		}
	}
}

// Update - Update game
func (g *Game) Update() {
	for i := int32(0); i <= g.Cols; i++ {
		for j := int32(0); j <= g.Rows; j++ {
			NeighborCount := g.CountNeighbors(i, j)
			if g.Cells[i][j].Alive {
				if NeighborCount < 2 {
					g.Cells[i][j].Next = false
				} else if NeighborCount > 3 {
					g.Cells[i][j].Next = false
				} else {
					g.Cells[i][j].Next = true
				}
			} else {
				if NeighborCount == 3 {
					g.Cells[i][j].Next = true
					g.Cells[i][j].Visited = true
				}
			}
		}
	}
	for i := int32(0); i <= g.Cols; i++ {
		for j := int32(0); j < g.Rows; j++ {
			g.Cells[i][j].Alive = g.Cells[i][j].Next
		}
	}
}

// CountNeighbors - Counts how many neighbous a cell has
func (g *Game) CountNeighbors(x, y int32) int {
	count := 0

	for i := int32(-1); i < 2; i++ {
		for j := int32(-1); j < 2; j++ {
			col := (x + i + (g.Cols)) % (g.Cols)
			row := (y + j + (g.Rows)) % (g.Rows)
			if g.Cells[col][row].Alive {
				count++
			}
		}
	}

	if g.Cells[x][y].Alive {
		count--
	}

	return count
}

// Draw - Draw game
func (g *Game) Draw() {
	raylib.BeginDrawing()
	raylib.ClearBackground(raylib.RayWhite)

	// Draw cells
	for x := int32(0); x <= g.Cols; x++ {
		for y := int32(0); y <= g.Rows; y++ {
			if g.Cells[x][y].Alive {
				raylib.DrawRectangleV(g.Cells[x][y].Position, g.Cells[x][y].Size, raylib.Blue)
			} else if g.Cells[x][y].Visited {
				raylib.DrawRectangleV(g.Cells[x][y].Position, g.Cells[x][y].Size, raylib.Color{R: 128, G: 177, B: 136, A: 255})
			}
		}
	}

	// Draw grid lines
	for i := int32(0); i < g.Cols+1; i++ {
		raylib.DrawLineV(
			raylib.NewVector2(float32(squareSize*i), 0),
			raylib.NewVector2(float32(squareSize*i), float32(g.ScreenHeight)),
			raylib.LightGray,
		)
	}

	for i := int32(0); i < g.Rows+1; i++ {
		raylib.DrawLineV(
			raylib.NewVector2(0, float32(squareSize*i)),
			raylib.NewVector2(float32(g.ScreenWidth), float32(squareSize*i)),
			raylib.LightGray,
		)
	}

	raylib.EndDrawing()
}
