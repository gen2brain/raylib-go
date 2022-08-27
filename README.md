![logo](https://goo.gl/XlIcXz)
## raylib-go
[![Build Status](https://github.com/gen2brain/raylib-go/actions/workflows/build.yml/badge.svg)](https://github.com/gen2brain/raylib-go/actions)
[![GoDoc](https://godoc.org/github.com/gen2brain/raylib-go/raylib?status.svg)](https://godoc.org/github.com/gen2brain/raylib-go/raylib)
[![Go Report Card](https://goreportcard.com/badge/github.com/gen2brain/raylib-go)](https://goreportcard.com/report/github.com/gen2brain/raylib-go)
[![Examples](https://img.shields.io/badge/learn%20by-examples-0077b3.svg?style=flat-square)](https://github.com/gen2brain/raylib-go/tree/master/examples)

Golang bindings for [raylib](http://www.raylib.com/), a simple and easy-to-use library to enjoy videogames programming.

### Requirements

##### Ubuntu

###### X11

    apt-get install libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev 

###### Wayland 

    apt-get install libgl1-mesa-dev libwayland-dev libxkbcommon-dev 

##### Fedora

###### X11

    dnf install mesa-libGL-devel libXi-devel libXcursor-devel libXrandr-devel libXinerama-devel

###### Wayland 

    dnf install mesa-libGL-devel wayland-devel libxkbcommon-devel

##### macOS

On macOS you need Xcode or Command Line Tools for Xcode.

##### Windows

On Windows you need C compiler, like [Mingw-w64](https://mingw-w64.org) or [TDM-GCC](http://tdm-gcc.tdragon.net/).
You can also build binary in [MSYS2](https://msys2.github.io/) shell.

##### Android

[Android example](https://github.com/gen2brain/raylib-go/tree/master/examples/others/android/example).

##### Raspberry Pi

[RPi example](https://github.com/gen2brain/raylib-go/tree/master/examples/others/rpi/basic_window).

### Installation

    go get -v -u github.com/gen2brain/raylib-go/raylib

### Build tags

* `drm` - build for Linux native mode, including Raspberry Pi 4 and other devices (PLATFORM_DRM)
* `rpi` - build for Raspberry Pi platform (PLATFORM_RPI)
* `wayland` - build against Wayland libraries
* `noaudio` - disables audio functions
* `opengl43` - uses OpenGL 4.3 backend
* `opengl21` - uses OpenGL 2.1 backend (default is 3.3 on desktop)
* `opengl11` - uses OpenGL 1.1 backend (pseudo OpenGL 1.1 style)

### Documentation

Documentation on [GoDoc](https://godoc.org/github.com/gen2brain/raylib-go/raylib). Also check raylib [cheatsheet](http://www.raylib.com/cheatsheet/cheatsheet.html).

### Example

```go
package main

import "github.com/gen2brain/raylib-go/raylib"

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
```

Check more [examples](https://github.com/gen2brain/raylib-go/tree/master/examples) organized by raylib modules.


### License

raylib-go is licensed under an unmodified zlib/libpng license. View [LICENSE](https://github.com/gen2brain/raylib-go/blob/master/LICENSE).
