package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game settings
const (
	screenWidth  = 800
	screenHeight = 600
)

// Paddle struct
type Paddle struct {
	Pos    rl.Vector2
	Width  float32
	Height float32
	Speed  float32
}

// Ball struct
type Ball struct {
	Pos    rl.Vector2
	Speed  rl.Vector2
	Radius float32
}

// Game struct
type Game struct {
	paddle1      Paddle
	paddle2      Paddle
	ball         Ball
	player1Score int
	player2Score int
}

// Check for collision between ball and paddle, and determine side
func CheckCollisionSide(ball Ball, paddle Paddle) string {
	closestX := rl.Clamp(ball.Pos.X, paddle.Pos.X, paddle.Pos.X+paddle.Width)
	closestY := rl.Clamp(ball.Pos.Y, paddle.Pos.Y, paddle.Pos.Y+paddle.Height)

	dx := ball.Pos.X - closestX
	dy := ball.Pos.Y - closestY

	distanceSq := dx*dx + dy*dy
	radiusSq := ball.Radius * ball.Radius

	if distanceSq <= radiusSq {
		// Determine where it hit
		if closestX == paddle.Pos.X || closestX == paddle.Pos.X+paddle.Width {
			return "side"
		}
		if closestY == paddle.Pos.Y || closestY == paddle.Pos.Y+paddle.Height {
			return "topbottom"
		}
		return "side" // fallback
	}

	return "none"
}


// Initialize the game state
func (g *Game) Init() {
	g.paddle1 = Paddle{
		Pos:    rl.NewVector2(30, float32(screenHeight/2-50)),
		Width:  10,
		Height: 100,
		Speed:  5,
	}

	g.paddle2 = Paddle{
		Pos:    rl.NewVector2(float32(screenWidth-40), float32(screenHeight/2-50)),
		Width:  10,
		Height: 100,
		Speed:  5,
	}

	g.ball = Ball{
		Pos:    rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		Speed:  rl.NewVector2(5, 4),
		Radius: 10,
	}

	g.player1Score = 0
	g.player2Score = 0
}

// Game update logic
func (g *Game) Update() {
	// Player 1 - W/S
	if rl.IsKeyDown(rl.KeyW) && g.paddle1.Pos.Y > 0 {
		g.paddle1.Pos.Y -= g.paddle1.Speed
	}
	if rl.IsKeyDown(rl.KeyS) && g.paddle1.Pos.Y < screenHeight-g.paddle1.Height {
		g.paddle1.Pos.Y += g.paddle1.Speed
	}

	// Player 2 - Up/Down
	if rl.IsKeyDown(rl.KeyUp) && g.paddle2.Pos.Y > 0 {
		g.paddle2.Pos.Y -= g.paddle2.Speed
	}
	if rl.IsKeyDown(rl.KeyDown) && g.paddle2.Pos.Y < screenHeight-g.paddle2.Height {
		g.paddle2.Pos.Y += g.paddle2.Speed
	}

	// Ball movement
	g.ball.Pos.X += g.ball.Speed.X
	g.ball.Pos.Y += g.ball.Speed.Y

	// Ball bounce off top/bottom walls
	if g.ball.Pos.Y <= g.ball.Radius || g.ball.Pos.Y >= float32(screenHeight)-g.ball.Radius {
		g.ball.Speed.Y *= -1
	}

	// Ball collisions with paddle sides only
	if CheckCollisionSide(g.ball, g.paddle1) == "side" {
		g.ball.Speed.X *= -1
		// Reposition ball outside paddle1
		g.ball.Pos.X = g.paddle1.Pos.X + g.paddle1.Width + g.ball.Radius
	}

	if CheckCollisionSide(g.ball, g.paddle2) == "side" {
		g.ball.Speed.X *= -1
		// Reposition ball outside paddle2
		g.ball.Pos.X = g.paddle2.Pos.X - g.ball.Radius
	}


	// Scoring
	if g.ball.Pos.X < 0 {
		g.player2Score++
		g.ball.Pos = rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2))
		g.ball.Speed.X *= -1
	}
	if g.ball.Pos.X > float32(screenWidth) {
		g.player1Score++
		g.ball.Pos = rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2))
		g.ball.Speed.X *= -1
	}
}

// Draw game elements
func (g *Game) Draw() {
	// Draw paddles
	rl.DrawRectangleV(g.paddle1.Pos, rl.NewVector2(g.paddle1.Width, g.paddle1.Height), rl.Black)
	rl.DrawRectangleV(g.paddle2.Pos, rl.NewVector2(g.paddle2.Width, g.paddle2.Height), rl.Black)

	// Draw ball
	rl.DrawCircleV(g.ball.Pos, g.ball.Radius, rl.Black)

	// Draw scores
	rl.DrawText(fmt.Sprintf("Player 1 : %d", g.player1Score), 20, 20, 20, rl.DarkGray)
	rl.DrawText(fmt.Sprintf("Player 2 : %d", g.player2Score), screenWidth-140, 20, 20, rl.DarkGray)
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Pong in Go!")
	rl.SetTargetFPS(60)

	var game Game
	game.Init()

	for !rl.WindowShouldClose() {
		game.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		game.Draw()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
