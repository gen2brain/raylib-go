package main

import (
	"math"
	"math/rand"

	"github.com/gen2brain/raylib-go/raylib"

	box2d "github.com/neguse/go-box2d-lite/box2dlite"
)

// Game type
type Game struct {
	World    *box2d.World
	TimeStep float64
}

// NewGame - Start new game
func NewGame() (g Game) {
	g.Init()
	return
}

// Init - Initialize game
func (g *Game) Init() {
	gravity := box2d.Vec2{0.0, -10.0}
	g.World = box2d.NewWorld(gravity, 10) // 10 iterations
}

// Update - Update game
func (g *Game) Update() {
	// Keys 1-9 switch demos
	switch rl.GetKeyPressed() {
	case rl.KeyOne:
		g.Demo1()
	case rl.KeyTwo:
		g.Demo2()
	case rl.KeyThree:
		g.Demo3()
	case rl.KeyFour:
		g.Demo4()
	case rl.KeyFive:
		g.Demo5()
	case rl.KeySix:
		g.Demo6()
	case rl.KeySeven:
		g.Demo7()
	case rl.KeyEight:
		g.Demo8()
	case rl.KeyNine:
		g.Demo9()
	}

	g.TimeStep = float64(rl.GetFrameTime())

	// Physics steps calculations
	g.World.Step(g.TimeStep)
}

// Draw - Draw game
func (g *Game) Draw() {
	for _, b := range g.World.Bodies {
		g.DrawBody(b)
	}
	for _, j := range g.World.Joints {
		g.DrawJoint(j)
	}

	rl.DrawText("Use keys 1-9 to switch current demo", 20, 20, 10, rl.RayWhite)
}

// DrawBody - Draw body
func (g *Game) DrawBody(b *box2d.Body) {
	R := box2d.Mat22ByAngle(b.Rotation)
	x := b.Position
	h := box2d.MulSV(0.5, b.Width)

	o := box2d.Vec2{400, 400}
	S := box2d.Mat22{box2d.Vec2{20.0, 0.0}, box2d.Vec2{0.0, -20.0}}

	v1 := o.Add(S.MulV(x.Add(R.MulV(box2d.Vec2{-h.X, -h.Y}))))
	v2 := o.Add(S.MulV(x.Add(R.MulV(box2d.Vec2{h.X, -h.Y}))))
	v3 := o.Add(S.MulV(x.Add(R.MulV(box2d.Vec2{h.X, h.Y}))))
	v4 := o.Add(S.MulV(x.Add(R.MulV(box2d.Vec2{-h.X, h.Y}))))

	rl.DrawLine(int32(v1.X), int32(v1.Y), int32(v2.X), int32(v2.Y), rl.RayWhite)
	rl.DrawLine(int32(v2.X), int32(v2.Y), int32(v3.X), int32(v3.Y), rl.RayWhite)
	rl.DrawLine(int32(v3.X), int32(v3.Y), int32(v4.X), int32(v4.Y), rl.RayWhite)
	rl.DrawLine(int32(v4.X), int32(v4.Y), int32(v1.X), int32(v1.Y), rl.RayWhite)
}

// DrawJoint - Draw joint
func (g *Game) DrawJoint(j *box2d.Joint) {
	b1 := j.Body1
	b2 := j.Body2

	R1 := box2d.Mat22ByAngle(b1.Rotation)
	R2 := box2d.Mat22ByAngle(b2.Rotation)

	x1 := b1.Position
	p1 := x1.Add(R1.MulV(j.LocalAnchor1))

	x2 := b2.Position
	p2 := x2.Add(R2.MulV(j.LocalAnchor2))

	o := box2d.Vec2{400, 400}
	S := box2d.Mat22{box2d.Vec2{20.0, 0.0}, box2d.Vec2{0.0, -20.0}}

	x1 = o.Add(S.MulV(x1))
	p1 = o.Add(S.MulV(p1))
	x2 = o.Add(S.MulV(x2))
	p2 = o.Add(S.MulV(p2))

	rl.DrawLine(int32(x1.X), int32(x1.Y), int32(p1.X), int32(p1.Y), rl.RayWhite)
	rl.DrawLine(int32(x2.X), int32(x2.Y), int32(p2.X), int32(p2.Y), rl.RayWhite)
}

// Demo1 - Single box
func (g *Game) Demo1() {
	g.World.Clear()

	var b1, b2 box2d.Body

	b1.Set(&box2d.Vec2{100.0, 20.0}, math.MaxFloat64)
	b1.Position = box2d.Vec2{0.0, -0.5 * b1.Width.Y}
	g.World.AddBody(&b1)

	b2.Set(&box2d.Vec2{1.0, 1.0}, 200.0)
	b2.Position = box2d.Vec2{0.0, 4.0}
	g.World.AddBody(&b2)
}

// Demo2 - A simple pendulum
func (g *Game) Demo2() {
	g.World.Clear()

	var b2, b1 box2d.Body
	var j box2d.Joint

	b1.Set(&box2d.Vec2{100.0, 20.0}, math.MaxFloat64)
	b1.Friction = 0.2
	b1.Position = box2d.Vec2{0.0, -0.5 * b1.Width.Y}
	b1.Rotation = 0.0
	g.World.AddBody(&b1)

	b2.Set(&box2d.Vec2{1.0, 1.0}, 100.0)
	b2.Friction = 0.2
	b2.Position = box2d.Vec2{9.0, 11.0}
	b2.Rotation = 0.0
	g.World.AddBody(&b2)

	j.Set(&b1, &b2, &box2d.Vec2{0.0, 11.0})
	g.World.AddJoint(&j)
}

// Demo3 - Varying friction coefficients
func (g *Game) Demo3() {
	g.World.Clear()

	{
		var b box2d.Body
		b.Set(&box2d.Vec2{100.0, 20.0}, math.MaxFloat64)
		b.Position = box2d.Vec2{0.0, -0.5 * b.Width.Y}
		g.World.AddBody(&b)
	}

	{
		var b box2d.Body
		b.Set(&box2d.Vec2{13.0, 0.25}, math.MaxFloat64)
		b.Position = box2d.Vec2{-2.0, 11.0}
		b.Rotation = -0.25
		g.World.AddBody(&b)
	}

	{
		var b box2d.Body
		b.Set(&box2d.Vec2{0.25, 1.0}, math.MaxFloat64)
		b.Position = box2d.Vec2{5.25, 9.5}
		g.World.AddBody(&b)
	}

	{
		var b box2d.Body
		b.Set(&box2d.Vec2{13.0, 0.25}, math.MaxFloat64)
		b.Position = box2d.Vec2{2.0, 7.0}
		b.Rotation = 0.25
		g.World.AddBody(&b)
	}

	{
		var b box2d.Body
		b.Set(&box2d.Vec2{0.25, 1.0}, math.MaxFloat64)
		b.Position = box2d.Vec2{-5.25, 5.5}
		g.World.AddBody(&b)
	}

	frictions := []float64{0.75, 0.5, 0.35, 0.1, 0.0}
	for i := 0; i < 5; i++ {
		var b box2d.Body
		b.Set(&box2d.Vec2{0.5, 0.5}, 25.0)
		b.Friction = frictions[i]
		b.Position = box2d.Vec2{-7.5 + 2.0*float64(i), 14.0}
		g.World.AddBody(&b)
	}

}

// Demo4 - A vertical stack
func (g *Game) Demo4() {
	g.World.Clear()

	{
		var b box2d.Body
		b.Set(&box2d.Vec2{100.0, 20.0}, math.MaxFloat64)
		b.Friction = 0.2
		b.Position = box2d.Vec2{0.0, -0.5 * b.Width.Y}
		b.Rotation = 0.0
		g.World.AddBody(&b)
	}

	for i := 0; i < 10; i++ {
		var b box2d.Body
		b.Set(&box2d.Vec2{1.0, 1.0}, 1.0)
		b.Friction = 0.2
		x := rand.Float64()*0.2 - 0.1
		b.Position = box2d.Vec2{x, 0.51 + 1.05*float64(i)}
		g.World.AddBody(&b)
	}

}

// Demo5 - A pyramid
func (g *Game) Demo5() {
	g.World.Clear()

	{
		var b box2d.Body
		b.Set(&box2d.Vec2{100.0, 20.0}, math.MaxFloat64)
		b.Friction = 0.2
		b.Position = box2d.Vec2{0.0, -0.5 * b.Width.Y}
		b.Rotation = 0.0
		g.World.AddBody(&b)
	}

	x := box2d.Vec2{-6.0, 0.75}

	for i := 0; i < 12; i++ {
		y := x
		for j := i; j < 12; j++ {
			var b box2d.Body
			b.Set(&box2d.Vec2{1.0, 1.0}, 10.0)
			b.Friction = 0.2
			b.Position = y
			g.World.AddBody(&b)

			y = y.Add(box2d.Vec2{1.125, 0.0})
		}

		x = x.Add(box2d.Vec2{0.5625, 2.0})
	}
}

// Demo6 - A teeter
func (g *Game) Demo6() {
	g.World.Clear()

	var b1, b2, b3, b4, b5 box2d.Body
	b1.Set(&box2d.Vec2{100.0, 20.0}, math.MaxFloat64)
	b1.Position = box2d.Vec2{0.0, -0.5 * b1.Width.Y}
	g.World.AddBody(&b1)

	b2.Set(&box2d.Vec2{12.0, 0.25}, 100)
	b2.Position = box2d.Vec2{0.0, 1.0}
	g.World.AddBody(&b2)

	b3.Set(&box2d.Vec2{0.5, 0.5}, 25.0)
	b3.Position = box2d.Vec2{-5.0, 2.0}
	g.World.AddBody(&b3)

	b4.Set(&box2d.Vec2{0.5, 0.5}, 25.0)
	b4.Position = box2d.Vec2{-5.5, 2.0}
	g.World.AddBody(&b4)

	b5.Set(&box2d.Vec2{1.0, 1.0}, 100)
	b5.Position = box2d.Vec2{5.5, 15.0}
	g.World.AddBody(&b5)

	{
		var j box2d.Joint
		j.Set(&b1, &b2, &box2d.Vec2{0.0, 1.0})
		g.World.AddJoint(&j)
	}

}

// Demo7 - A suspension bridge
func (g *Game) Demo7() {
	g.World.Clear()

	var ba []*box2d.Body

	{
		var b box2d.Body
		b.Set(&box2d.Vec2{100.0, 20.0}, math.MaxFloat64)
		b.Friction = 0.2
		b.Position = box2d.Vec2{0.0, -0.5 * b.Width.Y}
		b.Rotation = 0.0
		g.World.AddBody(&b)
		ba = append(ba, &b)
	}

	const numPlunks = 15
	const mass = 50.0

	for i := 0; i < numPlunks; i++ {
		var b box2d.Body
		b.Set(&box2d.Vec2{1.0, 0.25}, mass)
		b.Friction = 0.2
		b.Position = box2d.Vec2{-8.5 + 1.25*float64(i), 5.0}
		g.World.AddBody(&b)
		ba = append(ba, &b)
	}

	// Tuning
	const frequencyHz = 2.0
	const dampingRatio = 0.7

	// Frequency in radians
	const omega = 2.0 * math.Pi * frequencyHz

	// Damping coefficient
	const d = 2.0 * mass * dampingRatio * omega

	// Spring stifness
	const k = mass * omega * omega

	// Magic formulas
	softness := 1.0 / (d + g.TimeStep*k)
	biasFactor := g.TimeStep * k / (d + g.TimeStep*k)

	for i := 0; i <= numPlunks; i++ {
		var j box2d.Joint
		j.Set(ba[i], ba[(i+1)%(numPlunks+1)], &box2d.Vec2{-9.125 + 1.25*float64(i), 5.0})
		j.Softness = softness
		j.BiasFactor = biasFactor
		g.World.AddJoint(&j)
	}

}

// Demo8 - Dominos
func (g *Game) Demo8() {
	g.World.Clear()

	var b1 box2d.Body
	b1.Set(&box2d.Vec2{100.0, 20.0}, math.MaxFloat64)
	b1.Position = box2d.Vec2{0.0, -0.5 * b1.Width.Y}
	g.World.AddBody(&b1)

	{
		var b box2d.Body
		b.Set(&box2d.Vec2{12.0, 0.5}, math.MaxFloat64)
		b.Position = box2d.Vec2{-1.5, 10.0}
		g.World.AddBody(&b)
	}

	for i := 0; i < 10; i++ {
		var b box2d.Body
		b.Set(&box2d.Vec2{0.2, 2.0}, 10.0)
		b.Position = box2d.Vec2{-6.0 + 1.0*float64(i), 11.125}
		b.Friction = 0.1
		g.World.AddBody(&b)
	}

	{
		var b box2d.Body
		b.Set(&box2d.Vec2{14.0, 0.5}, math.MaxFloat64)
		b.Position = box2d.Vec2{1.0, 6.0}
		b.Rotation = 0.3
		g.World.AddBody(&b)
	}

	var b2 box2d.Body
	b2.Set(&box2d.Vec2{0.5, 3.0}, math.MaxFloat64)
	b2.Position = box2d.Vec2{-7.0, 4.0}
	g.World.AddBody(&b2)

	var b3 box2d.Body
	b3.Set(&box2d.Vec2{12.0, 0.25}, 20.0)
	b3.Position = box2d.Vec2{-0.9, 1.0}
	g.World.AddBody(&b3)

	{
		var j box2d.Joint
		j.Set(&b1, &b3, &box2d.Vec2{-2.0, 1.0})
		g.World.AddJoint(&j)
	}

	var b4 box2d.Body
	b4.Set(&box2d.Vec2{0.5, 0.5}, 10.0)
	b4.Position = box2d.Vec2{-10.0, 15.0}
	g.World.AddBody(&b4)

	{
		var j box2d.Joint
		j.Set(&b2, &b4, &box2d.Vec2{-7.0, 15.0})
		g.World.AddJoint(&j)
	}

	var b5 box2d.Body
	b5.Set(&box2d.Vec2{2.0, 2.0}, 20.0)
	b5.Position = box2d.Vec2{6.0, 2.5}
	b5.Friction = 0.1
	g.World.AddBody(&b5)

	{
		var j box2d.Joint
		j.Set(&b1, &b5, &box2d.Vec2{6.0, 2.6})
		g.World.AddJoint(&j)
	}

	var b6 box2d.Body
	b6.Set(&box2d.Vec2{2.0, 0.2}, 10.0)
	b6.Position = box2d.Vec2{6.0, 3.6}
	g.World.AddBody(&b6)

	{
		var j box2d.Joint
		j.Set(&b5, &b6, &box2d.Vec2{7.0, 3.5})
		g.World.AddJoint(&j)
	}

}

// Demo9 - A multi-pendulum
func (g *Game) Demo9() {
	g.World.Clear()

	var b1 *box2d.Body

	{
		var b box2d.Body
		b.Set(&box2d.Vec2{100.0, 20.0}, math.MaxFloat64)
		b.Position = box2d.Vec2{0.0, -0.5 * b.Width.Y}
		g.World.AddBody(&b)
		b1 = &b
	}

	const mass = 10.0

	// Tuning
	const frequencyHz = 4.0
	const dampingRatio = 0.7

	// Frequency in radians
	const omega = 2.0 * math.Pi * frequencyHz

	// Damping coefficient
	const d = 2.0 * mass * dampingRatio * omega

	// Spring stiffness
	const k = mass * omega * omega

	// Magic formulas
	softness := 1.0 / (d + g.TimeStep*k)
	biasFactor := g.TimeStep * k / (d + g.TimeStep*k)

	const y = 12.0

	for i := 0; i < 15; i++ {
		x := box2d.Vec2{0.5 + float64(i), y}

		var b box2d.Body
		b.Set(&box2d.Vec2{0.75, 0.25}, mass)
		b.Friction = 0.2
		b.Position = x
		b.Rotation = 0.0
		g.World.AddBody(&b)

		var j box2d.Joint
		j.Set(b1, &b, &box2d.Vec2{float64(i), y})
		j.Softness = softness
		j.BiasFactor = biasFactor
		g.World.AddJoint(&j)

		b1 = &b
	}

}

func main() {
	rl.InitWindow(800, 450, "raylib [physics] example - box2d")

	rl.SetTargetFPS(60)

	game := NewGame()

	game.Demo1()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		game.Update()

		game.Draw()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
