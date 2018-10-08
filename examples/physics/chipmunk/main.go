package main

import (
	"math"
	"math/rand"

	"github.com/gen2brain/raylib-go/raylib"

	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

const (
	ballRadius = 25
	ballMass   = 1
)

// Game type
type Game struct {
	Space       *chipmunk.Space
	Balls       []*chipmunk.Shape
	StaticLines []*chipmunk.Shape

	ticksToNextBall int
}

// NewGame - Start new game
func NewGame() (g Game) {
	g.Init()
	return
}

// Init - Initialize game
func (g *Game) Init() {
	g.createBodies()

	g.ticksToNextBall = 10
}

// Update - Update game
func (g *Game) Update() {
	g.ticksToNextBall--
	if g.ticksToNextBall == 0 {
		g.ticksToNextBall = rand.Intn(100) + 1
		g.addBall()
	}

	// Physics steps calculations
	g.step(rl.GetFrameTime())
}

// Draw - Draw game
func (g *Game) Draw() {
	for i := range g.StaticLines {
		x := g.StaticLines[i].GetAsSegment().A.X
		y := g.StaticLines[i].GetAsSegment().A.Y

		x2 := g.StaticLines[i].GetAsSegment().B.X
		y2 := g.StaticLines[i].GetAsSegment().B.Y

		rl.DrawLine(int32(x), int32(y), int32(x2), int32(y2), rl.DarkBlue)
	}

	for _, b := range g.Balls {
		pos := b.Body.Position()
		rl.DrawCircleLines(int32(pos.X), int32(pos.Y), float32(ballRadius), rl.DarkBlue)
	}
}

// createBodies sets up the chipmunk space and static bodies
func (g *Game) createBodies() {
	g.Space = chipmunk.NewSpace()
	g.Space.Gravity = vect.Vect{0, 900}

	staticBody := chipmunk.NewBodyStatic()
	g.StaticLines = []*chipmunk.Shape{
		chipmunk.NewSegment(vect.Vect{250.0, 240.0}, vect.Vect{550.0, 280.0}, 0),
		chipmunk.NewSegment(vect.Vect{550.0, 280.0}, vect.Vect{550.0, 180.0}, 0),
	}

	for _, segment := range g.StaticLines {
		segment.SetElasticity(0.6)
		staticBody.AddShape(segment)
	}

	g.Space.AddBody(staticBody)
}

// addBall adds ball to chipmunk space and body
func (g *Game) addBall() {
	x := rand.Intn(600-200) + 200
	ball := chipmunk.NewCircle(vect.Vector_Zero, float32(ballRadius))
	ball.SetElasticity(0.95)

	body := chipmunk.NewBody(vect.Float(ballMass), ball.Moment(float32(ballMass)))
	body.SetPosition(vect.Vect{vect.Float(x), 0.0})
	body.SetAngle(vect.Float(rand.Float32() * 2 * math.Pi))
	body.AddShape(ball)

	g.Space.AddBody(body)
	g.Balls = append(g.Balls, ball)
}

// step advances the physics engine and cleans up any balls that are off-screen
func (g *Game) step(dt float32) {
	g.Space.Step(vect.Float(dt))

	for i := 0; i < len(g.Balls); i++ {
		p := g.Balls[i].Body.Position()
		if p.Y < -100 {
			g.Space.RemoveBody(g.Balls[i].Body)
			g.Balls[i] = nil
			g.Balls = append(g.Balls[:i], g.Balls[i+1:]...)
			i-- // consider same index again
		}
	}
}

func main() {
	rl.InitWindow(800, 450, "raylib [physics] example - chipmunk")

	rl.SetTargetFPS(60)

	game := NewGame()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		game.Update()

		game.Draw()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
