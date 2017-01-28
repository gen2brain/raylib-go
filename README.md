## raylib-go [![GoDoc](https://godoc.org/github.com/gen2brain/raylib-go/raylib?status.svg)](https://godoc.org/github.com/gen2brain/raylib-go/raylib)

Golang bindings for [raylib](http://www.raylib.com/), a simple and easy-to-use library to learn videogames programming.

![screenshot](https://goo.gl/q6DAoy)

### Requirements

* [GLFW3](http://www.glfw.org/) (desktop platform only, not needed on Android and RPi)
* [OpenAL Soft](http://kcat.strangesoft.net/openal.html)

##### Ubuntu

    apt-get install libglfw3-dev
    apt-get install openal-dev

##### Fedora
    
    dnf install glfw-devel
    dnf install openal-soft-devel

##### OS X

    brew install glfw3
    brew install openal-soft

### Installation

    go get -v github.com/gen2brain/raylib-go

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
