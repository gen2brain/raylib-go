module example

go 1.19

require (
	github.com/gen2brain/raylib-go/raygui v0.0.0-20221122151443-e8a384ed1346
	github.com/gen2brain/raylib-go/raylib v0.0.0-20221122155035-fe6d2c0ed32a
)

replace github.com/gen2brain/raylib-go/raylib => ../../../raylib

replace github.com/gen2brain/raylib-go/raygui => ../../../raygui3_5
