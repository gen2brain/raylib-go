module examples

go 1.21

replace github.com/gen2brain/raylib-go/raylib => ../raylib

replace github.com/gen2brain/raylib-go/raygui => ../raygui

require (
	github.com/gen2brain/raylib-go/easings v0.0.0-20231021203613-2d673bb5f4b3
	github.com/gen2brain/raylib-go/physics v0.0.0-20231021203613-2d673bb5f4b3
	github.com/gen2brain/raylib-go/raygui v0.0.0-20231110085703-5830da3d8795
	github.com/gen2brain/raylib-go/raylib v0.0.0-20231110085703-5830da3d8795
	github.com/jakecoffman/cp v1.2.1
	github.com/neguse/go-box2d-lite v0.0.0-20170921151050-5d8ed9b7272b
)
