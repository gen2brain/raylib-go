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
	Position raylib.Vector2
}

// Pipe type
type Pipe struct {
	Rec    raylib.Rectangle
	Color  raylib.Color
	Active bool
}

// Particle type
type Particle struct {
	Position raylib.Vector2
	Color    raylib.Color
	Alpha    float32
	Size     float32
	Rotation float32
	Active   bool
}

// Game type
type Game struct {
	FxFlap  raylib.Sound
	FxSlap  raylib.Sound
	FxPoint raylib.Sound
	FxClick raylib.Sound

	TxSprites raylib.Texture2D
	TxSmoke   raylib.Texture2D
	TxClouds  raylib.Texture2D

	CloudRec raylib.Rectangle
	FrameRec raylib.Rectangle

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
	PipesPos []raylib.Vector2
}

// NewGame - Start new game
func NewGame() (g Game) {
	g.Init()
	return
}

// On Android this sets callback function to be used for android_main
func init() {
	raylib.SetCallbackFunc(main)
}

func main() {
	// Initialize game
	game := NewGame()
	game.GameOver = true

	// Initialize window
	raylib.InitWindow(screenWidth, screenHeight, "Floppy Gopher")

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

// Init - Initialize game
func (g *Game) Init() {

	// Gopher
	g.Floppy = Floppy{raylib.NewVector2(80, float32(screenHeight)/2-spriteSize/2)}

	// Sprite rectangle
	g.FrameRec = raylib.NewRectangle(0, 0, spriteSize, spriteSize)

	// Cloud rectangle
	g.CloudRec = raylib.NewRectangle(0, 0, float32(screenWidth), float32(g.TxClouds.Height))

	// Initialize particles
	g.Particles = make([]Particle, maxParticles)
	for i := 0; i < maxParticles; i++ {
		g.Particles[i].Position = raylib.NewVector2(0, 0)
		g.Particles[i].Color = raylib.RayWhite
		g.Particles[i].Alpha = 1.0
		g.Particles[i].Size = float32(raylib.GetRandomValue(1, 30)) / 20.0
		g.Particles[i].Rotation = float32(raylib.GetRandomValue(0, 360))
		g.Particles[i].Active = false
	}

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
		g.Pipes[i].Rec.X = g.PipesPos[i/2].X
		g.Pipes[i].Rec.Y = g.PipesPos[i/2].Y
		g.Pipes[i].Rec.Width = pipesWidth
		g.Pipes[i].Rec.Height = 550
		g.Pipes[i].Color = colors[raylib.GetRandomValue(0, int32(len(colors)-1))]

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
	g.FxFlap = raylib.LoadSound("sounds/flap.wav")
	g.FxSlap = raylib.LoadSound("sounds/slap.wav")
	g.FxPoint = raylib.LoadSound("sounds/point.wav")
	g.FxClick = raylib.LoadSound("sounds/click.wav")
	g.TxSprites = raylib.LoadTexture("images/sprite.png")
	g.TxSmoke = raylib.LoadTexture("images/smoke.png")
	g.TxClouds = raylib.LoadTexture("images/clouds.png")
}

// Unload - Unload resources
func (g *Game) Unload() {
	raylib.UnloadSound(g.FxFlap)
	raylib.UnloadSound(g.FxSlap)
	raylib.UnloadSound(g.FxPoint)
	raylib.UnloadSound(g.FxClick)
	raylib.UnloadTexture(g.TxSprites)
	raylib.UnloadTexture(g.TxSmoke)
	raylib.UnloadTexture(g.TxClouds)
}

// Update - Update game
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
				if raylib.IsKeyDown(raylib.KeySpace) || raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
					raylib.PlaySound(g.FxFlap)

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
					if raylib.CheckCollisionRecs(raylib.NewRectangle(g.Floppy.Position.X, g.Floppy.Position.Y, spriteSize, spriteSize), g.Pipes[i].Rec) {
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

				// Switch dead sprite
				if g.FramesCounter >= 8 {
					g.FrameRec.X = spriteSize * 5
				} else {
					g.FrameRec.X = spriteSize * 4
				}
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
	raylib.BeginDrawing()

	raylib.ClearBackground(raylib.SkyBlue)

	if !g.GameOver {
		// Draw clouds
		raylib.DrawTextureRec(g.TxClouds, g.CloudRec, raylib.NewVector2(0, float32(screenHeight-g.TxClouds.Height)), raylib.RayWhite)

		// Draw rotated clouds
		raylib.DrawTexturePro(g.TxClouds, raylib.NewRectangle(-g.CloudRec.X, 0, float32(g.TxClouds.Width), float32(g.TxClouds.Height)),
			raylib.NewRectangle(0, 0, float32(g.TxClouds.Width), float32(g.TxClouds.Height)), raylib.NewVector2(float32(g.TxClouds.Width), float32(g.TxClouds.Height)), 180, raylib.White)

		// Draw Gopher
		raylib.DrawTextureRec(g.TxSprites, g.FrameRec, g.Floppy.Position, raylib.RayWhite)

		// Draw active particles
		if !g.Dead {
			for i := 0; i < maxParticles; i++ {
				if g.Particles[i].Active {
					raylib.DrawTexturePro(
						g.TxSmoke,
						raylib.NewRectangle(0, 0, float32(g.TxSmoke.Width), float32(g.TxSmoke.Height)),
						raylib.NewRectangle(g.Particles[i].Position.X, g.Particles[i].Position.Y, float32(g.TxSmoke.Width)*g.Particles[i].Size, float32(g.TxSmoke.Height)*g.Particles[i].Size),
						raylib.NewVector2(float32(g.TxSmoke.Width)*g.Particles[i].Size/2, float32(g.TxSmoke.Height)*g.Particles[i].Size/2),
						g.Particles[i].Rotation,
						raylib.Fade(g.Particles[i].Color, g.Particles[i].Alpha),
					)
				}
			}
		}

		// Draw pipes
		for i := 0; i < maxPipes; i++ {
			raylib.DrawRectangle(int32(g.Pipes[i*2].Rec.X), int32(g.Pipes[i*2].Rec.Y), int32(g.Pipes[i*2].Rec.Width), int32(g.Pipes[i*2].Rec.Height), g.Pipes[i*2].Color)
			raylib.DrawRectangle(int32(g.Pipes[i*2+1].Rec.X), int32(g.Pipes[i*2+1].Rec.Y), int32(g.Pipes[i*2+1].Rec.Width), int32(g.Pipes[i*2+1].Rec.Height), g.Pipes[i*2].Color)

			// Draw borders
			raylib.DrawRectangleLines(int32(g.Pipes[i*2].Rec.X), int32(g.Pipes[i*2].Rec.Y), int32(g.Pipes[i*2].Rec.Width), int32(g.Pipes[i*2].Rec.Height), raylib.Black)
			raylib.DrawRectangleLines(int32(g.Pipes[i*2+1].Rec.X), int32(g.Pipes[i*2+1].Rec.Y), int32(g.Pipes[i*2+1].Rec.Width), int32(g.Pipes[i*2+1].Rec.Height), raylib.Black)
		}

		// Draw Super Flashing FX (one frame only)
		if g.SuperFX {
			raylib.DrawRectangle(0, 0, screenWidth, screenHeight, raylib.White)
			g.SuperFX = false
		}

		// Draw HI-SCORE
		raylib.DrawText(fmt.Sprintf("%02d", g.Score), 20, 20, 32, raylib.Black)
		raylib.DrawText(fmt.Sprintf("HI-SCORE: %02d", g.HiScore), 20, 64, 20, raylib.Black)

		if g.Pause {
			// Draw PAUSED text
			raylib.DrawText("PAUSED", screenWidth/2-raylib.MeasureText("PAUSED", 24)/2, screenHeight/2-50, 20, raylib.Black)
		}
	} else {
		// Draw text
		raylib.DrawText("Floppy Gopher", raylib.GetScreenWidth()/2-raylib.MeasureText("Floppy Gopher", 40)/2, raylib.GetScreenHeight()/2-150, 40, raylib.RayWhite)

		if runtime.GOOS == "android" {
			raylib.DrawText("[TAP] TO PLAY", raylib.GetScreenWidth()/2-raylib.MeasureText("[TAP] TO PLAY", 20)/2, raylib.GetScreenHeight()/2-50, 20, raylib.Black)
		} else {
			raylib.DrawText("[ENTER] TO PLAY", raylib.GetScreenWidth()/2-raylib.MeasureText("[ENTER] TO PLAY", 20)/2, raylib.GetScreenHeight()/2-50, 20, raylib.Black)
		}

		// Draw Gopher
		raylib.DrawTextureRec(g.TxSprites, g.FrameRec, raylib.NewVector2(float32(raylib.GetScreenWidth()/2-spriteSize/2), float32(raylib.GetScreenHeight()/2)), raylib.RayWhite)
	}

	raylib.EndDrawing()
}
