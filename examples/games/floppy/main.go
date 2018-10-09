package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	// Screen width
	screenWidth = 504
	// Screen height
	screenHeight = 896

	// Maximum number of pipes
	maxPipes = 100
	// Maximum number of particles
	maxParticles = 50

	// Pipes width
	pipesWidth = 60
	// Sprite size
	spriteSize = 48

	// Pipes speed
	pipesSpeedX = 2.5
	// Clouds speed
	cloudsSpeedX = 1

	// Gravity
	gravity = 1.2
)

// Floppy type
type Floppy struct {
	Position rl.Vector2
}

// Pipe type
type Pipe struct {
	Rec    rl.Rectangle
	Color  rl.Color
	Active bool
}

// Particle type
type Particle struct {
	Position rl.Vector2
	Color    rl.Color
	Alpha    float32
	Size     float32
	Rotation float32
	Active   bool
}

// Game type
type Game struct {
	FxFlap  rl.Sound
	FxSlap  rl.Sound
	FxPoint rl.Sound
	FxClick rl.Sound

	TxSprites rl.Texture2D
	TxSmoke   rl.Texture2D
	TxClouds  rl.Texture2D

	CloudRec rl.Rectangle
	FrameRec rl.Rectangle

	GameOver bool
	Dead     bool
	Pause    bool
	SuperFX  bool

	Score         int
	HiScore       int
	FramesCounter int32

	WindowShouldClose bool

	Floppy    Floppy
	Particles []Particle

	Pipes    []Pipe
	PipesPos []rl.Vector2
}

// NewGame - Start new game
func NewGame() (g Game) {
	g.Init()
	return
}

// On Android this sets callback function to be used for android_main
func init() {
	rl.SetCallbackFunc(main)
}

func main() {
	// Initialize game
	game := NewGame()
	game.GameOver = true

	// Initialize window
	rl.InitWindow(screenWidth, screenHeight, "Floppy Gopher")

	// Initialize audio
	rl.InitAudioDevice()

	// NOTE: Textures and Sounds MUST be loaded after Window/Audio initialization
	game.Load()

	// Limit FPS
	rl.SetTargetFPS(60)

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
	rl.CloseAudioDevice()

	// Close window
	rl.CloseWindow()

	// Exit
	os.Exit(0)
}

// Init - Initialize game
func (g *Game) Init() {

	// Gopher
	g.Floppy = Floppy{rl.NewVector2(80, float32(screenHeight)/2-spriteSize/2)}

	// Sprite rectangle
	g.FrameRec = rl.NewRectangle(0, 0, spriteSize, spriteSize)

	// Cloud rectangle
	g.CloudRec = rl.NewRectangle(0, 0, float32(screenWidth), float32(g.TxClouds.Height))

	// Initialize particles
	g.Particles = make([]Particle, maxParticles)
	for i := 0; i < maxParticles; i++ {
		g.Particles[i].Position = rl.NewVector2(0, 0)
		g.Particles[i].Color = rl.RayWhite
		g.Particles[i].Alpha = 1.0
		g.Particles[i].Size = float32(rl.GetRandomValue(1, 30)) / 20.0
		g.Particles[i].Rotation = float32(rl.GetRandomValue(0, 360))
		g.Particles[i].Active = false
	}

	// Pipes positions
	g.PipesPos = make([]rl.Vector2, maxPipes)
	for i := 0; i < maxPipes; i++ {
		g.PipesPos[i].X = float32(480 + 360*i)
		g.PipesPos[i].Y = -float32(rl.GetRandomValue(0, 240))
	}

	// Pipes colors
	colors := []rl.Color{
		rl.Orange, rl.Red, rl.Gold, rl.Lime,
		rl.Violet, rl.Brown, rl.LightGray, rl.Blue,
		rl.Yellow, rl.Green, rl.Purple, rl.Beige,
	}

	// Pipes
	g.Pipes = make([]Pipe, maxPipes*2)
	for i := 0; i < maxPipes*2; i += 2 {
		g.Pipes[i].Rec.X = g.PipesPos[i/2].X
		g.Pipes[i].Rec.Y = g.PipesPos[i/2].Y
		g.Pipes[i].Rec.Width = pipesWidth
		g.Pipes[i].Rec.Height = 550
		g.Pipes[i].Color = colors[rl.GetRandomValue(0, int32(len(colors)-1))]

		g.Pipes[i+1].Rec.X = g.PipesPos[i/2].X
		g.Pipes[i+1].Rec.Y = 1200 + g.PipesPos[i/2].Y - 550
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

// Load - Load resources
func (g *Game) Load() {
	g.FxFlap = rl.LoadSound("sounds/flap.wav")
	g.FxSlap = rl.LoadSound("sounds/slap.wav")
	g.FxPoint = rl.LoadSound("sounds/point.wav")
	g.FxClick = rl.LoadSound("sounds/click.wav")
	g.TxSprites = rl.LoadTexture("images/sprite.png")
	g.TxSmoke = rl.LoadTexture("images/smoke.png")
	g.TxClouds = rl.LoadTexture("images/clouds.png")
}

// Unload - Unload resources
func (g *Game) Unload() {
	rl.UnloadSound(g.FxFlap)
	rl.UnloadSound(g.FxSlap)
	rl.UnloadSound(g.FxPoint)
	rl.UnloadSound(g.FxClick)
	rl.UnloadTexture(g.TxSprites)
	rl.UnloadTexture(g.TxSmoke)
	rl.UnloadTexture(g.TxClouds)
}

// Update - Update game
func (g *Game) Update() {
	if rl.WindowShouldClose() {
		g.WindowShouldClose = true
	}

	if !g.GameOver {
		if rl.IsKeyPressed(rl.KeyP) || rl.IsKeyPressed(rl.KeyBack) {
			rl.PlaySound(g.FxClick)

			if runtime.GOOS == "android" && g.Pause {
				g.WindowShouldClose = true
			}

			g.Pause = !g.Pause
		}

		if !g.Pause {
			if !g.Dead {
				// Scroll pipes
				for i := 0; i < maxPipes; i++ {
					g.PipesPos[i].X -= float32(pipesSpeedX)
				}

				for i := 0; i < maxPipes*2; i += 2 {
					g.Pipes[i].Rec.X = g.PipesPos[i/2].X
					g.Pipes[i+1].Rec.X = g.PipesPos[i/2].X
				}

				// Scroll clouds
				g.CloudRec.X += cloudsSpeedX
				if g.CloudRec.X > float32(g.TxClouds.Width) {
					g.CloudRec.X = 0
				}

				// Movement/Controls
				if rl.IsKeyDown(rl.KeySpace) || rl.IsMouseButtonDown(rl.MouseLeftButton) {
					rl.PlaySound(g.FxFlap)

					// Activate one particle every frame
					for i := 0; i < maxParticles; i++ {
						if !g.Particles[i].Active {
							g.Particles[i].Active = true
							g.Particles[i].Alpha = 1.0
							g.Particles[i].Position = g.Floppy.Position
							g.Particles[i].Position.X += spriteSize / 2
							g.Particles[i].Position.Y += spriteSize / 2
							i = maxParticles
						}
					}

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
					// Switch run sprites every 8 frames
					g.FramesCounter++
					if g.FramesCounter >= 8 {
						g.FramesCounter = 0
						g.FrameRec.X = spriteSize
					} else {
						g.FrameRec.X = 0
					}

					// Floppy fall down
					g.Floppy.Position.Y += gravity
				}

				// Update active particles
				for i := 0; i < maxParticles; i++ {
					if g.Particles[i].Active {
						g.Particles[i].Position.X -= 1.0
						g.Particles[i].Alpha -= 0.05

						if g.Particles[i].Alpha <= 0.0 {
							g.Particles[i].Active = false
						}

						g.Particles[i].Rotation += 3.0
					}
				}

				// Check Collisions
				for i := 0; i < maxPipes*2; i++ {
					if rl.CheckCollisionRecs(rl.NewRectangle(g.Floppy.Position.X, g.Floppy.Position.Y, spriteSize, spriteSize), g.Pipes[i].Rec) {
						// OMG You killed Gopher you bastard!
						g.Dead = true

						rl.PlaySound(g.FxSlap)
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

						rl.PlaySound(g.FxPoint)
					}
				}
			} else {
				// Wait 60 frames before GameOver
				g.FramesCounter++
				if g.FramesCounter >= 60 {
					g.GameOver = true
				}

				// Switch dead sprite
				if g.FramesCounter >= 8 {
					g.FrameRec.X = spriteSize * 5
				} else {
					g.FrameRec.X = spriteSize * 4
				}
			}
		} else {
			if rl.IsMouseButtonDown(rl.MouseLeftButton) {
				g.Pause = !g.Pause
			}
		}
	} else {
		if rl.IsKeyPressed(rl.KeyEnter) || rl.IsMouseButtonDown(rl.MouseLeftButton) {
			rl.PlaySound(g.FxClick)

			// Return of the Gopher!
			g.Init()
		} else if runtime.GOOS == "android" && rl.IsKeyDown(rl.KeyBack) {
			g.WindowShouldClose = true
		}

		// Switch flap sprites
		g.FramesCounter++
		if g.FramesCounter >= 8 {
			g.FramesCounter = 0
			g.FrameRec.X = spriteSize
		} else {
			g.FrameRec.X = 0
		}

	}
}

// Draw - Draw game
func (g *Game) Draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.SkyBlue)

	if !g.GameOver {
		// Draw clouds
		rl.DrawTextureRec(g.TxClouds, g.CloudRec, rl.NewVector2(0, float32(screenHeight-g.TxClouds.Height)), rl.RayWhite)

		// Draw rotated clouds
		rl.DrawTexturePro(g.TxClouds, rl.NewRectangle(-g.CloudRec.X, 0, float32(g.TxClouds.Width), float32(g.TxClouds.Height)),
			rl.NewRectangle(0, 0, float32(g.TxClouds.Width), float32(g.TxClouds.Height)), rl.NewVector2(float32(g.TxClouds.Width), float32(g.TxClouds.Height)), 180, rl.White)

		// Draw Gopher
		rl.DrawTextureRec(g.TxSprites, g.FrameRec, g.Floppy.Position, rl.RayWhite)

		// Draw active particles
		if !g.Dead {
			for i := 0; i < maxParticles; i++ {
				if g.Particles[i].Active {
					rl.DrawTexturePro(
						g.TxSmoke,
						rl.NewRectangle(0, 0, float32(g.TxSmoke.Width), float32(g.TxSmoke.Height)),
						rl.NewRectangle(g.Particles[i].Position.X, g.Particles[i].Position.Y, float32(g.TxSmoke.Width)*g.Particles[i].Size, float32(g.TxSmoke.Height)*g.Particles[i].Size),
						rl.NewVector2(float32(g.TxSmoke.Width)*g.Particles[i].Size/2, float32(g.TxSmoke.Height)*g.Particles[i].Size/2),
						g.Particles[i].Rotation,
						rl.Fade(g.Particles[i].Color, g.Particles[i].Alpha),
					)
				}
			}
		}

		// Draw pipes
		for i := 0; i < maxPipes; i++ {
			rl.DrawRectangle(int32(g.Pipes[i*2].Rec.X), int32(g.Pipes[i*2].Rec.Y), int32(g.Pipes[i*2].Rec.Width), int32(g.Pipes[i*2].Rec.Height), g.Pipes[i*2].Color)
			rl.DrawRectangle(int32(g.Pipes[i*2+1].Rec.X), int32(g.Pipes[i*2+1].Rec.Y), int32(g.Pipes[i*2+1].Rec.Width), int32(g.Pipes[i*2+1].Rec.Height), g.Pipes[i*2].Color)

			// Draw borders
			rl.DrawRectangleLines(int32(g.Pipes[i*2].Rec.X), int32(g.Pipes[i*2].Rec.Y), int32(g.Pipes[i*2].Rec.Width), int32(g.Pipes[i*2].Rec.Height), rl.Black)
			rl.DrawRectangleLines(int32(g.Pipes[i*2+1].Rec.X), int32(g.Pipes[i*2+1].Rec.Y), int32(g.Pipes[i*2+1].Rec.Width), int32(g.Pipes[i*2+1].Rec.Height), rl.Black)
		}

		// Draw Super Flashing FX (one frame only)
		if g.SuperFX {
			rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.White)
			g.SuperFX = false
		}

		// Draw HI-SCORE
		rl.DrawText(fmt.Sprintf("%02d", g.Score), 20, 20, 32, rl.Black)
		rl.DrawText(fmt.Sprintf("HI-SCORE: %02d", g.HiScore), 20, 64, 20, rl.Black)

		if g.Pause {
			// Draw PAUSED text
			rl.DrawText("PAUSED", screenWidth/2-rl.MeasureText("PAUSED", 24)/2, screenHeight/2-50, 20, rl.Black)
		}
	} else {
		// Draw text
		rl.DrawText("Floppy Gopher", int32(rl.GetScreenWidth())/2-rl.MeasureText("Floppy Gopher", 40)/2, int32(rl.GetScreenHeight())/2-150, 40, rl.RayWhite)

		if runtime.GOOS == "android" {
			rl.DrawText("[TAP] TO PLAY", int32(rl.GetScreenWidth())/2-rl.MeasureText("[TAP] TO PLAY", 20)/2, int32(rl.GetScreenHeight())/2-50, 20, rl.Black)
		} else {
			rl.DrawText("[ENTER] TO PLAY", int32(rl.GetScreenWidth())/2-rl.MeasureText("[ENTER] TO PLAY", 20)/2, int32(rl.GetScreenHeight())/2-50, 20, rl.Black)
		}

		// Draw Gopher
		rl.DrawTextureRec(g.TxSprites, g.FrameRec, rl.NewVector2(float32(rl.GetScreenWidth()/2-spriteSize/2), float32(rl.GetScreenHeight()/2)), rl.RayWhite)
	}

	rl.EndDrawing()
}
