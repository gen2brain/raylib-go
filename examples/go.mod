module examples

go 1.22

replace github.com/gen2brain/raylib-go/raylib => ../raylib

replace github.com/gen2brain/raylib-go/raygui => ../raygui

replace github.com/gen2brain/raylib-go/easings => ../easings

replace github.com/gen2brain/raylib-go/physics => ../physics

require (
	github.com/gen2brain/raylib-go/easings v0.62.0
	github.com/gen2brain/raylib-go/physics v0.62.0
	github.com/gen2brain/raylib-go/raygui v0.62.0
	github.com/gen2brain/raylib-go/raylib v0.62.0
	github.com/jakecoffman/cp v1.2.1
	github.com/neguse/go-box2d-lite v0.0.0-20260315142500-8c06e40d0e7b
)

require (
	github.com/ebitengine/purego v0.10.2 // indirect
	github.com/jupiterrider/ffi v0.8.0 // indirect
)
