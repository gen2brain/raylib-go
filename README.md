![logo](https://goo.gl/XlIcXz)
## raylib-go
[![TravisCI Build Status](https://travis-ci.org/gen2brain/raylib-go.svg?branch=master)](https://travis-ci.org/gen2brain/raylib-go)
[![AppVeyor Build status](https://ci.appveyor.com/api/projects/status/qv2iggrqtgl7xhr0?svg=true)](https://ci.appveyor.com/project/gen2brain/raylib-go)
[![GoDoc](https://godoc.org/github.com/gen2brain/raylib-go/raylib?status.svg)](https://godoc.org/github.com/gen2brain/raylib-go/raylib)
[![Go Report Card](https://goreportcard.com/badge/github.com/gen2brain/raylib-go)](https://goreportcard.com/report/github.com/gen2brain/raylib-go)
[![Examples](https://img.shields.io/badge/learn%20by-examples-0077b3.svg?style=flat-square)](https://github.com/gen2brain/raylib-go/tree/master/examples)

Golang bindings for [raylib](http://www.raylib.com/), a simple and easy-to-use library to learn videogames programming.

![screenshot](https://goo.gl/q6DAoy)

### Requirements

* [OpenAL Soft](http://kcat.strangesoft.net/openal.html)
NOTE: if you don't need audio you can use `-tags noaudio` during build, OpenAL will not be linked to binary, though none of the audio functions will be available.

* [GLFW](http://www.glfw.org/) is included as part of the Go package, but you need to make sure you have dependencies installed, see below.

##### Ubuntu

    apt-get install libopenal-dev libgl1-mesa-dev libxi-dev libxinerama-dev libxcursor-dev libxxf86vm-dev libxrandr-dev

##### Fedora

    dnf install openal-soft-devel mesa-libGL-devel libXi-devel libXcursor-devel libXrandr-devel libXinerama-devel

##### OS X

On OS X system OpenAL framework is used, you need Xcode or Command Line Tools for Xcode.

##### Windows ([MSYS2](https://msys2.github.io/))

    pacman -S mingw-w64-x86_64-openal mingw-w64-x86_64-gcc mingw-w64-x86_64-go git

On Windows, build binary in MSYS2 shell.

##### Android

[Android example](https://github.com/gen2brain/raylib-go/tree/master/examples/android/example).

##### Raspberry Pi

[RPi example](https://github.com/gen2brain/raylib-go/tree/master/examples/rpi/basic_window).

### Installation

    go get -v -u github.com/gen2brain/raylib-go/raylib

### Build tags

* `noaudio` - disables audio functions and doesn't link against OpenAL libraries
* `opengl21` - use OpenGL 2.1 backend (default is 3.3 on desktop)
* `opengl11` - use OpenGL 1.1 backend (pseudo OpenGL 1.1 style)
* `wayland` - builds against Wayland libraries
* `static` - link against OpenAL static libraries

### Documentation

Documentation on [GoDoc](https://godoc.org/github.com/gen2brain/raylib-go/raylib). Also check raylib [cheatsheet](http://www.raylib.com/cheatsheet/cheatsheet.html).

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
