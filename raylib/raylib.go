/*
Package raylib - Go bindings for raylib, a simple and easy-to-use library to learn videogames programming.

raylib is highly inspired by Borland BGI graphics lib and by XNA framework.

raylib could be useful for prototyping, tools development, graphic applications, embedded systems and education.

NOTE for ADVENTURERS: raylib is a programming library to learn videogames programming; no fancy interface, no visual helpers, no auto-debugging... just coding in the most pure spartan-programmers way.

Example:

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


*/
package raylib

import (
	"runtime"
)

func init() {
	// Make sure the main goroutine is bound to the main thread.
	runtime.LockOSThread()
}
