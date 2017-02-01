## raylib-go [![Build Status](https://travis-ci.org/gen2brain/raylib-go.svg?branch=master)](https://travis-ci.org/gen2brain/raylib-go) [![GoDoc](https://godoc.org/github.com/gen2brain/raylib-go/raylib?status.svg)](https://godoc.org/github.com/gen2brain/raylib-go/raylib) [![Go Report Card](https://goreportcard.com/badge/github.com/gen2brain/raylib-go)](https://goreportcard.com/report/github.com/gen2brain/raylib-go)

Golang bindings for [raylib](http://www.raylib.com/), a simple and easy-to-use library to learn videogames programming.

![screenshot](https://goo.gl/q6DAoy)

### Requirements

* [GLFW3](http://www.glfw.org/) (desktop platform only, not needed on Android and RPi)
* [OpenAL Soft](http://kcat.strangesoft.net/openal.html) (on OS X system framework is used)

GLFW version 3.2 is required.

##### Ubuntu

    apt-get install libglfw3-dev libopenal-dev libxi-dev libxinerama-dev libxcursor-dev libxxf86vm-dev
    
On older Ubuntu releases you will need to compile GLFW, instructions are in [travis file](https://github.com/gen2brain/raylib-go/blob/master/.travis.yml).

##### Fedora
    
    dnf install glfw-devel openal-soft-devel mesa-libGL-devel libXi-devel

##### OS X

    brew install glfw3

##### Windows ([MSYS2](https://msys2.github.io/))

    pacman -S mingw-w64-x86_64-openal mingw-w64-x86_64-glfw mingw-w64-x86_64-gcc mingw-w64-x86_64-go git 

##### Android

[Android example](https://github.com/gen2brain/raylib-go/tree/master/examples/android/example).

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
