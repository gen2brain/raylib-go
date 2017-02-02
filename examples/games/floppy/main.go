package main

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	// Maximum number of pipes
	maxPipes = 100
	// Pipes width
	pipesWidth = 60
	// Sprite size
	spriteSize = 48
)

// Floppy
type Floppy struct {
	Position raylib.Vector2
}

// Pipe
type Pipe struct {
	Rec    raylib.Rectangle
	Color  raylib.Color
	Active bool
}

// Game
type Game struct {
	ScreenWidth  int32
	ScreenHeight int32

	FxFlap  raylib.Sound
	FxSlap  raylib.Sound
	FxPoint raylib.Sound
	FxClick raylib.Sound

	Texture  raylib.Texture2D
	FrameRec raylib.Rectangle

	GameOver bool
	Dead     bool
	Pause    bool
	SuperFX  bool

	Score             int
	HiScore           int
	FramesCounter     int32
	WindowShouldClose bool

	Floppy Floppy

	Pipes       []Pipe
	PipesPos    []raylib.Vector2
	PipesSpeedX int32
}

// New Game
func NewGame() (g Game) {
	g.Init()
	return
}

// On Android this sets callback function to be used for android_main
func init() {
	raylib.SetCallbackFunc(run)
}

// Main function, not used on Android, used on desktop platform
func main() {
	run(nil)
}

// Callback function on Android, not needed on desktop platform
func run(app unsafe.Pointer) {
	// Initialize game
	game := NewGame()

	// Initialize window
	if runtime.GOOS != "android" {
		raylib.InitWindow(game.ScreenWidth, game.ScreenHeight, "Floppy Gopher")
	} else {
		raylib.InitWindow(game.ScreenWidth, game.ScreenHeight, app)
	}

	// Initialize audio
	raylib.InitAudioDevice()

	// NOTE: Textures and Sounds MUST be loaded after Window/Audio initialization
	game.Load()

	// Limit FPS
	raylib.SetTargetFPS(60)

	// Main loop
	for !game.WindowShouldClose {
		// Update game
		game.Update()

		// Draw game
		game.Draw()
	}

	// Free resources
	game.Unload()

	// Close audio
	raylib.CloseAudioDevice()

	// Close window
	raylib.CloseWindow()

	// Exit
	os.Exit(0)
}

// Initialize game
func (g *Game) Init() {
	// Window resolution
	g.ScreenWidth = 504
	g.ScreenHeight = 896

	g.Floppy = Floppy{}
	g.Floppy.Position = raylib.NewVector2(80, float32(g.ScreenHeight)/2-spriteSize/2)
	g.PipesSpeedX = 2

	// Sprite rectangle
	g.FrameRec = raylib.NewRectangle(0, 0, spriteSize, spriteSize)

	// Pipes positions
	g.PipesPos = make([]raylib.Vector2, maxPipes)
	for i := 0; i < maxPipes; i++ {
		g.PipesPos[i].X = float32(480 + 360*i)
		g.PipesPos[i].Y = -float32(raylib.GetRandomValue(0, 240))
	}

	// Pipes colors
	colors := []raylib.Color{
		raylib.Orange, raylib.Red, raylib.Gold, raylib.Lime,
		raylib.Violet, raylib.Brown, raylib.LightGray, raylib.Blue,
		raylib.Yellow, raylib.Green, raylib.Purple, raylib.Beige,
	}

	// Pipes
	g.Pipes = make([]Pipe, maxPipes*2)
	for i := 0; i < maxPipes*2; i += 2 {
		g.Pipes[i].Rec.X = int32(g.PipesPos[i/2].X)
		g.Pipes[i].Rec.Y = int32(g.PipesPos[i/2].Y)
		g.Pipes[i].Rec.Width = pipesWidth
		g.Pipes[i].Rec.Height = 550
		g.Pipes[i].Color = colors[raylib.GetRandomValue(0, int32(len(colors)-1))]

		g.Pipes[i+1].Rec.X = int32(g.PipesPos[i/2].X)
		g.Pipes[i+1].Rec.Y = int32(1200 + g.PipesPos[i/2].Y - 550)
		g.Pipes[i+1].Rec.Width = pipesWidth
		g.Pipes[i+1].Rec.Height = 550

		g.Pipes[i/2].Active = true
	}

	g.Score = 0
	g.FramesCounter = 0
	g.WindowShouldClose = false

	g.GameOver = false
	g.Dead = false
	g.SuperFX = false
	g.Pause = false
}

// Load resources
func (g *Game) Load() {
	g.FxFlap = raylib.LoadSound("sounds/flap.wav")
	g.FxSlap = raylib.LoadSound("sounds/slap.wav")
	g.FxPoint = raylib.LoadSound("sounds/point.wav")
	g.FxClick = raylib.LoadSound("sounds/click.wav")
	g.Texture = raylib.LoadTexture("images/sprite.png")
}

// Unload resources
func (g *Game) Unload() {
	raylib.UnloadSound(g.FxFlap)
	raylib.UnloadSound(g.FxSlap)
	raylib.UnloadSound(g.FxPoint)
	raylib.UnloadSound(g.FxClick)
	raylib.UnloadTexture(g.Texture)
}

// Update game
func (g *Game) Update() {
	if raylib.WindowShouldClose() {
		g.WindowShouldClose = true
	}

	if !g.GameOver {
		if raylib.IsKeyPressed(raylib.KeyP) || raylib.IsKeyPressed(raylib.KeyBack) {
			raylib.PlaySound(g.FxClick)

			if runtime.GOOS == "android" && g.Pause {
				g.WindowShouldClose = true
			}

			g.Pause = !g.Pause
		}

		if !g.Pause {
			if !g.Dead {
				// Scroll X
				for i := 0; i < maxPipes; i++ {
					g.PipesPos[i].X -= float32(g.PipesSpeedX)
				}

				for i := 0; i < maxPipes*2; i += 2 {
					g.Pipes[i].Rec.X = int32(g.PipesPos[i/2].X)
					g.Pipes[i+1].Rec.X = int32(g.PipesPos[i/2].X)
				}

				// Movement/Controls
				if raylib.IsKeyDown(raylib.KeySpace) || raylib.IsMouseButtonDown(raylib.MouseLeftButton) && !g.GameOver {
					raylib.PlaySound(g.FxFlap)

					// Switch flap sprites every 8 frames
					g.FramesCounter++
					if g.FramesCounter >= 8 {
						g.FramesCounter = 0
						g.FrameRec.X = spriteSize * 3
					} else {
						g.FrameRec.X = spriteSize * 2
					}

					// Floppy go up
					g.Floppy.Position.Y -= 3
				} else {
					// Default sprite
					//g.FrameRec.X = spriteSize

					// Switch flap sprites every 8 frames
					g.FramesCounter++
					if g.FramesCounter >= 8 {
						g.FramesCounter = 0
						g.FrameRec.X = spriteSize
					} else {
						g.FrameRec.X = 0
					}

					// Floppy fall down
					g.Floppy.Position.Y += 1
				}
			}

			if !g.Dead {
				// Check Collisions
				for i := 0; i < maxPipes*2; i++ {
					if raylib.CheckCollisionRecs(raylib.NewRectangle(int32(g.Floppy.Position.X), int32(g.Floppy.Position.Y), spriteSize, spriteSize), g.Pipes[i].Rec) {
						// OMG You killed Gopher you bastard!
						g.Dead = true

						raylib.PlaySound(g.FxSlap)
					} else if (g.PipesPos[i/2].X < g.Floppy.Position.X-spriteSize) && g.Pipes[i/2].Active && !g.GameOver {
						// Score point
						g.Score += 1
						g.Pipes[i/2].Active = false

						// Flash screen
						g.SuperFX = true

						// Update HiScore
						if g.Score > g.HiScore {
							g.HiScore = g.Score
						}

						raylib.PlaySound(g.FxPoint)
					}
				}
			} else {
				// Wait 60 frames before GameOver
				g.FramesCounter++
				if g.FramesCounter >= 60 {
					g.GameOver = true
				}

				g.FrameRec.X = spriteSize * 4
			}
		} else {
			if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
				g.Pause = !g.Pause
			}
		}
	} else {
		if raylib.IsKeyPressed(raylib.KeyEnter) || raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			raylib.PlaySound(g.FxClick)

			// Return of the Gopher!
			g.Init()
		} else if runtime.GOOS == "android" && raylib.IsKeyDown(raylib.KeyBack) {
			g.WindowShouldClose = true
		}

	}
}

// Draw game
func (g *Game) Draw() {
	raylib.BeginDrawing()

	raylib.ClearBackground(raylib.SkyBlue)

	if !g.GameOver {
		// Draw Gopher
		raylib.DrawTextureRec(g.Texture, g.FrameRec, g.Floppy.Position, raylib.RayWhite) // Draw part of the texture

		// Draw pipes
		for i := 0; i < maxPipes; i++ {
			raylib.DrawRectangle(g.Pipes[i*2].Rec.X, g.Pipes[i*2].Rec.Y, g.Pipes[i*2].Rec.Width, g.Pipes[i*2].Rec.Height, g.Pipes[i*2].Color)
			raylib.DrawRectangle(g.Pipes[i*2+1].Rec.X, g.Pipes[i*2+1].Rec.Y, g.Pipes[i*2+1].Rec.Width, g.Pipes[i*2+1].Rec.Height, g.Pipes[i*2].Color)

			// Draw borders
			raylib.DrawRectangleLines(g.Pipes[i*2].Rec.X, g.Pipes[i*2].Rec.Y, g.Pipes[i*2].Rec.Width, g.Pipes[i*2].Rec.Height, raylib.Black)
			raylib.DrawRectangleLines(g.Pipes[i*2+1].Rec.X, g.Pipes[i*2+1].Rec.Y, g.Pipes[i*2+1].Rec.Width, g.Pipes[i*2+1].Rec.Height, raylib.Black)
		}

		// Draw Super Flashing FX (one frame only)
		if g.SuperFX {
			raylib.DrawRectangle(0, 0, g.ScreenWidth, g.ScreenHeight, raylib.White)
			g.SuperFX = false
		}

		// Draw HI-SCORE
		raylib.DrawText(fmt.Sprintf("%02d", g.Score), 20, 20, 32, raylib.RayWhite)
		raylib.DrawText(fmt.Sprintf("HI-SCORE: %02d", g.HiScore), 20, 64, 20, raylib.RayWhite)

		if g.Pause {
			// Draw PAUSED text
			raylib.DrawText("PAUSED", g.ScreenWidth/2-raylib.MeasureText("PAUSED", 24)/2, g.ScreenHeight/2-50, 20, raylib.RayWhite)
		}
	} else {
		// Draw PLAY AGAIN text
		if runtime.GOOS == "android" {
			raylib.DrawText("[TAP] TO PLAY AGAIN", raylib.GetScreenWidth()/2-raylib.MeasureText("[TAP] TO PLAY AGAIN", 20)/2, raylib.GetScreenHeight()/2-50, 20, raylib.RayWhite)
		} else {
			raylib.DrawText("[ENTER] TO PLAY AGAIN", raylib.GetScreenWidth()/2-raylib.MeasureText("[ENTER] TO PLAY AGAIN", 20)/2, raylib.GetScreenHeight()/2-50, 20, raylib.RayWhite)
		}
	}

	raylib.EndDrawing()
}
