module examples

go 1.22

replace github.com/gen2brain/raylib-go/raylib => ../raylib

replace github.com/gen2brain/raylib-go/raygui => ../raygui

replace github.com/gen2brain/raylib-go/easings => ../easings

replace github.com/gen2brain/raylib-go/physics => ../physics

require (
	github.com/gen2brain/raylib-go/easings v0.0.0-00010101000000-000000000000
	github.com/gen2brain/raylib-go/physics v0.0.0-00010101000000-000000000000
	github.com/gen2brain/raylib-go/raygui v0.0.0-00010101000000-000000000000
	github.com/gen2brain/raylib-go/raylib v0.62.0
	github.com/jakecoffman/cp v1.2.1
	github.com/neguse/go-box2d-lite v0.0.0-20170921151050-5d8ed9b7272b
)

require (
	github.com/ebitengine/purego v0.10.2 // indirect
	github.com/jupiterrider/ffi v0.8.0 // indirect
)
