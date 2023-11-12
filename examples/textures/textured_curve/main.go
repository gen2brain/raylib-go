package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	texRoad       rl.Texture2D
	showCurve     = true
	curveW        = float32(50)
	curveSegments = 24

	curveStartPos, curveStartPosTangent, curveEndPos, curveEndPosTangent rl.Vector2

	curveSelectedPoint *rl.Vector2

	screenW = int32(800)
	screenH = int32(450)
)

func main() {

	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagMsaa4xHint)

	rl.InitWindow(screenW, screenH, "raylib [textures] example - textured curve")

	texRoad = rl.LoadTexture("road.png")
	rl.SetTextureFilter(texRoad, rl.TextureFilterMode(rl.FilterBilinear))

	curveStartPos = rl.NewVector2(80, 100)
	curveStartPosTangent = rl.NewVector2(100, 300)

	curveEndPos = rl.NewVector2(700, 350)
	curveEndPosTangent = rl.NewVector2(600, 100)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		upCurve()
		upOptions()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		drawTexturedCurve()
		drawCurve()

		rl.DrawText("Drag points to move curve, press SPACE to show/hide base curve", 10, 10, 10, rl.Black)
		rl.DrawText("Curve width: "+fmt.Sprintf("%.0f", curveW)+" use UP/DOWN arrows to adjust", 10, 30, 10, rl.Black)
		rl.DrawText("Curve segments: "+fmt.Sprint(curveSegments)+" use RIGHT/LEFT arrows to adjust", 10, 50, 10, rl.Black)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texRoad)

	rl.CloseWindow()
}

func upCurve() {

	if !rl.IsMouseButtonDown(rl.MouseLeftButton) {
		curveSelectedPoint = &rl.Vector2{}
	}

	*curveSelectedPoint = rl.Vector2Add(*curveSelectedPoint, rl.GetMouseDelta())

	mouse := rl.GetMousePosition()

	if rl.CheckCollisionPointCircle(mouse, curveStartPos, 6) {
		curveSelectedPoint = &curveStartPos
	} else if rl.CheckCollisionPointCircle(mouse, curveStartPosTangent, 6) {
		curveSelectedPoint = &curveStartPosTangent
	} else if rl.CheckCollisionPointCircle(mouse, curveEndPos, 6) {
		curveSelectedPoint = &curveEndPos
	} else if rl.CheckCollisionPointCircle(mouse, curveEndPosTangent, 6) {
		curveSelectedPoint = &curveEndPosTangent
	}

}
func upOptions() {

	if rl.IsKeyPressed(rl.KeySpace) {
		showCurve = !showCurve
	}
	if rl.IsKeyPressed(rl.KeyUp) {
		curveW += 2
	}
	if rl.IsKeyPressed(rl.KeyDown) {
		curveW -= 2
	}
	if curveW < 2 {
		curveW = 2
	}
	if rl.IsKeyPressed(rl.KeyLeft) {
		curveSegments -= 2
	}
	if rl.IsKeyPressed(rl.KeyRight) {
		curveSegments += 2
	}
	if curveSegments < 2 {
		curveSegments = 2
	}

}
func drawTexturedCurve() {

	step := float32(1) / float32(curveSegments)
	previous := curveStartPos
	previousTangent := rl.Vector2Zero()
	previousV := float32(0)
	tangentSet := false
	current := rl.Vector2Zero()
	t := float32(0)

	for i := 0; i < curveSegments; i++ {
		t = step * float32(i)
		a := float32(math.Pow(1-float64(t), 3))
		b := 3 * float32(math.Pow(1-float64(t), 2)) * t
		c := 3 * (1 - t) * float32(math.Pow(float64(t), 2))
		d := float32(math.Pow(float64(t), 3))

		current.Y = a*curveStartPos.Y + b*curveStartPosTangent.Y + c*curveEndPosTangent.Y + d*curveEndPos.Y
		current.X = a*curveStartPos.X + b*curveStartPosTangent.X + c*curveEndPosTangent.X + d*curveEndPos.X

		delta := rl.NewVector2(current.X-previous.X, current.Y-previous.Y)
		normal := rl.Vector2Normalize(rl.NewVector2(-delta.Y, delta.X))
		v := previousV + rl.Vector2Length(delta)

		if !tangentSet {
			previousTangent = normal
			tangentSet = true
		}

		prevPosNormal := rl.Vector2Add(previous, rl.Vector2Scale(previousTangent, curveW))
		prevNegNormal := rl.Vector2Add(previous, rl.Vector2Scale(previousTangent, -curveW))

		currentPosNormal := rl.Vector2Add(current, rl.Vector2Scale(normal, curveW))
		currentNegNormal := rl.Vector2Add(current, rl.Vector2Scale(normal, -curveW))

		rl.SetTexture(texRoad.ID)
		rl.Begin(rl.Quads)

		rl.Color4ub(255, 255, 255, 255)
		rl.Normal3f(0, 0, 1)

		rl.TexCoord2f(0, previousV)
		rl.Vertex2f(prevNegNormal.X, prevNegNormal.Y)

		rl.TexCoord2f(1, previousV)
		rl.Vertex2f(prevPosNormal.X, prevPosNormal.Y)

		rl.TexCoord2f(1, v)
		rl.Vertex2f(currentPosNormal.X, currentPosNormal.Y)

		rl.TexCoord2f(0, v)
		rl.Vertex2f(currentNegNormal.X, currentNegNormal.Y)

		rl.End()

		previous = current
		previousTangent = normal
		previousV = v

	}

}
func drawCurve() {

	if showCurve {
		rl.DrawSplineSegmentBezierCubic(curveStartPos, curveEndPos, curveStartPosTangent, curveEndPosTangent, 2, rl.Blue)
	}
	rl.DrawLineV(curveStartPos, curveStartPosTangent, rl.SkyBlue)
	rl.DrawLineV(curveEndPos, curveEndPosTangent, rl.Purple)
	mouse := rl.GetMousePosition()

	if rl.CheckCollisionPointCircle(mouse, curveStartPos, 6) {
		rl.DrawCircleV(curveStartPos, 7, rl.Yellow)
	}
	rl.DrawCircleV(curveStartPos, 5, rl.Red)

	if rl.CheckCollisionPointCircle(mouse, curveStartPosTangent, 6) {
		rl.DrawCircleV(curveStartPosTangent, 7, rl.Yellow)
	}
	rl.DrawCircleV(curveStartPosTangent, 5, rl.Maroon)

	if rl.CheckCollisionPointCircle(mouse, curveEndPos, 6) {
		rl.DrawCircleV(curveEndPos, 7, rl.Yellow)
	}
	rl.DrawCircleV(curveEndPosTangent, 5, rl.Green)

	if rl.CheckCollisionPointCircle(mouse, curveEndPosTangent, 6) {
		rl.DrawCircleV(curveEndPosTangent, 7, rl.Yellow)
	}
	rl.DrawCircleV(curveEndPosTangent, 5, rl.DarkGreen)

}
