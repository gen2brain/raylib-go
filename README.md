![logo](https://goo.gl/jhpm6K)
## raylib-go [![Build Status](https://travis-ci.org/gen2brain/raylib-go.svg?branch=master)](https://travis-ci.org/gen2brain/raylib-go) [![GoDoc](https://godoc.org/github.com/gen2brain/raylib-go/raylib?status.svg)](https://godoc.org/github.com/gen2brain/raylib-go/raylib) [![Go Report Card](https://goreportcard.com/badge/github.com/gen2brain/raylib-go)](https://goreportcard.com/report/github.com/gen2brain/raylib-go)

Golang bindings for [raylib](http://www.raylib.com/), a simple and easy-to-use library to learn videogames programming.

![screenshot](https://goo.gl/q6DAoy)

### Requirements

* [OpenAL Soft](http://kcat.strangesoft.net/openal.html)

##### Ubuntu

    apt-get install libopenal-dev libxi-dev libxinerama-dev libxcursor-dev libxxf86vm-dev

##### Fedora

    dnf install openal-soft-devel mesa-libGL-devel libXi-devel libXcursor-devel libXrandr-devel libXinerama-devel

##### OS X

On OS X system OpenAL framework is used.

##### Windows ([MSYS2](https://msys2.github.io/))

    pacman -S mingw-w64-x86_64-openal mingw-w64-x86_64-gcc mingw-w64-x86_64-go git

##### Android

[Android example](https://github.com/gen2brain/raylib-go/tree/master/examples/android/example).

##### Raspberry Pi

[RPi example](https://github.com/gen2brain/raylib-go/tree/master/examples/rpi/basic_window).

### Installation

    go get -v github.com/gen2brain/raylib-go/raylib

### Example

```go
package main

import "github.com/gen2brain/raylib-go/raylib"

func main() {
	raylib.InitWindow(800, 450, "raylib [core] example - basic window")

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("Congrats! You created your first window!", 190, 200, 20, raylib.LightGray)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
```

Check more [examples](https://github.com/gen2brain/raylib-go/tree/master/examples) organized by raylib modules.


### License

raylib-go is licensed under an unmodified zlib/libpng license. View [LICENSE](https://github.com/gen2brain/raylib-go/blob/master/LICENSE).
