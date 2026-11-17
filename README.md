![logo](logo.png)
## raylib-go
[![Build Status](https://github.com/gen2brain/raylib-go/actions/workflows/build.yml/badge.svg)](https://github.com/gen2brain/raylib-go/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/gen2brain/raylib-go/raylib.svg)](https://pkg.go.dev/github.com/gen2brain/raylib-go/raylib)
[![Examples](https://img.shields.io/badge/learn%20by-examples-0077b3.svg?style=flat-square)](https://github.com/gen2brain/raylib-go/tree/master/examples)

Golang bindings for [raylib](https://www.raylib.com/), a simple and easy-to-use library to enjoy videogames programming.

raylib C source code is included and compiled together with bindings. Note that the first build can take a few minutes.

It is also possible to use raylib-go without cgo (see below).

### Requirements

#### cgo

##### Ubuntu

    apt-get install libgl1-mesa-dev libegl1-mesa-dev libgbm-dev libdrm-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev libxext-dev libwayland-dev libxkbcommon-dev

##### Fedora

    dnf install mesa-libGL-devel mesa-libEGL-devel mesa-libgbm-devel libdrm-devel libXi-devel libXcursor-devel libXrandr-devel libXinerama-devel libXext-devel wayland-devel libxkbcommon-devel

##### macOS

On macOS, you need Xcode or Command Line Tools for Xcode (if you have `brew` installed, you already have this).

##### Windows

On Windows you need C compiler, like [Mingw-w64](https://mingw-w64.org) or [TDM-GCC](http://tdm-gcc.tdragon.net/).
You can also build binary in [MSYS2](https://msys2.github.io/) shell.

To remove console window, build with `-ldflags "-H=windowsgui"`.

##### Android

[Android example](https://github.com/gen2brain/raylib-go/tree/master/examples/others/android/example).

##### iOS

[iOS example](https://github.com/gen2brain/raylib-go/tree/master/examples/others/ios/example).

##### Wasm

[Wasm example](https://github.com/gen2brain/raylib-go/tree/master/examples/others/wasm/example).

#### purego (without cgo)

You can use raylib-go on `Linux`, `macOS`, `Windows` and `FreeBSD` without the need of [cgo](https://go.dev/blog/cgo), i.e. with `CGO_ENABLED=0`,
or with the `nocgo` build tag when cgo is enabled. 32-bit is not supported.

The shared libraries of raylib are embedded for Linux, macOS and Windows on amd64 and arm64.
With the build tag `extlib` nothing is embedded and the system raylib library is loaded instead.

On FreeBSD nothing is embedded, the system raylib library is loaded; build with `-gcflags="github.com/ebitengine/purego/internal/fakecgo=-std"`.


### Installation

    go get github.com/gen2brain/raylib-go/raylib

### Build tags

Platforms (cgo, default is the native platform: Win32 on Windows, Cocoa on macOS, X11 on BSD, and on Linux Wayland, X11 and DRM selected at runtime):

| Tag       | Description                                                                                                                                   |
|-----------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| `glfw`    | [GLFW](https://www.glfw.org/) backend (PLATFORM_DESKTOP_GLFW)                                                                                 |
| `drm`     | Linux native [DRM](https://en.wikipedia.org/wiki/Direct_Rendering_Manager) mode only, including Raspberry Pi and other devices (PLATFORM_DRM) |
| `sdl`     | [SDL](https://github.com/libsdl-org/SDL) backend (PLATFORM_DESKTOP_SDL)                                                                       |
| `sdl3`    | [SDL3](https://github.com/libsdl-org/SDL) backend (PLATFORM_DESKTOP_SDL3)                                                                     |
| `rgfw`    | [RGFW](https://github.com/ColleagueRiley/RGFW) backend (PLATFORM_DESKTOP_RGFW)                                                                |
| `wayland` | Native Wayland backend only, Linux only (PLATFORM_DESKTOP_WAYLAND)                                                                            |
| `x11`     | Native X11 backend only, Linux and BSD (PLATFORM_DESKTOP_X11)                                                                                 |
| `memory`  | Headless, renders to memory (PLATFORM_MEMORY)                                                                                                 |
| `termux`  | [Termux](https://termux.dev/) on Android (X11)                                                                                                |

On Linux the default platform (PLATFORM_DESKTOP_LINUX) uses Wayland if `WAYLAND_DISPLAY` is set, X11 if `DISPLAY` is set, and DRM otherwise;
the environment variable `RAYLIB_PLATFORM` (`wayland`, `x11` or `drm`) selects one.

Graphics API (cgo, default is OpenGL 3.3 on desktop, OpenGL ES 2.0 on Linux arm64):

| Tag        | Description                                                                                   |
|------------|-----------------------------------------------------------------------------------------------|
| `opengl43` | OpenGL 4.3                                                                                    |
| `opengl33` | OpenGL 3.3                                                                                    |
| `opengl21` | OpenGL 2.1                                                                                    |
| `opengl11` | OpenGL 1.1 (pseudo OpenGL 1.1 style)                                                          |
| `es2`      | OpenGL ES 2.0 (can be used to link against [Google's ANGLE](https://github.com/google/angle)) |
| `es3`      | OpenGL ES 3.0 (experimental)                                                                  |
| `software` | Software renderer (rlsw)                                                                      |

Options:

| Tag            | Description                                                                   |
|----------------|-------------------------------------------------------------------------------|
| `glfw_x11`     | GLFW with X11 only (implies `glfw`)                                           |
| `glfw_wayland` | GLFW with Wayland only (implies `glfw`)                                       |
| `noaudio`      | Disable audio functions                                                       |
| `nocgo`        | Use the purego backend even when cgo is enabled                               |
| `extlib`       | Load the system raylib library instead of the embedded one (used with purego) |

### Documentation

Documentation on [pkg.go.dev](https://pkg.go.dev/github.com/gen2brain/raylib-go/raylib). Also check raylib [cheatsheet](https://www.raylib.com/cheatsheet/cheatsheet.html). If you have problems or need assistance there is an active community on the [Raylib Discord Server](https://discord.gg/raylib) that can help.

### Example

```go
package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
  rl.InitWindow(800, 450, "raylib [core] example - basic window")
  defer rl.CloseWindow()

  rl.SetTargetFPS(60)

  for !rl.WindowShouldClose() {
    rl.BeginDrawing()

    rl.ClearBackground(rl.RayWhite)
    rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

    rl.EndDrawing()
  }
}
```

Check more [examples](https://github.com/gen2brain/raylib-go/tree/master/examples) organized by raylib modules.

### Cross-compile (on Linux with cgo)

To cross-compile for Windows, install [MinGW](https://www.mingw-w64.org/) toolchain.

```
$ CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags "-s -w"
$ file basic_window.exe
basic_window.exe: PE32+ executable for MS Windows 10.00 (console), x86-64, 5 sections

$ CGO_ENABLED=1 CC=i686-w64-mingw32-gcc GOOS=windows GOARCH=386 go build -ldflags "-s -w"
$ file basic_window.exe
basic_window.exe: PE32 executable for MS Windows 10.00 (console), Intel i386, 4 sections
```

To cross-compile for Windows on ARM, install [LLVM MinGW](https://github.com/mstorsjo/llvm-mingw) toolchain.

```
$ CGO_ENABLED=1 CC=aarch64-w64-mingw32-clang GOOS=windows GOARCH=arm64 go build -ldflags "-s -w"
$ file basic_window.exe
basic_window.exe: PE32+ executable for MS Windows 10.00 (console), ARM64, 6 sections
```

To cross-compile for macOS, install [OSXCross](https://github.com/tpoechtrager/osxcross) toolchain.
Ready-made SDK tarballs are available [here](https://github.com/joseluisq/macosx-sdks).

```
$ CGO_ENABLED=1 CC=x86_64-apple-darwin25.1-clang GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w"
$ file basic_window
basic_window: Mach-O 64-bit x86_64 executable, flags:<NOUNDEFS|DYLDLINK|TWOLEVEL|WEAK_DEFINES|BINDS_TO_WEAK|PIE>

$ CGO_ENABLED=1 CC=aarch64-apple-darwin25.1-clang GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w"
$ file basic_window
basic_window: Mach-O 64-bit arm64 executable, flags:<NOUNDEFS|DYLDLINK|TWOLEVEL|WEAK_DEFINES|BINDS_TO_WEAK|PIE>
```

### License

raylib-go is licensed under an unmodified zlib/libpng license. View [LICENSE](https://github.com/gen2brain/raylib-go/blob/master/LICENSE).
