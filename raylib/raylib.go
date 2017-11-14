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

// CallQueueCap is the capacity of the call queue.
//
// The default value is 16.
var CallQueueCap = 16

// callInMain calls a function in the main thread. It is only properly initialized inside raylib.Main(..).
// As a default, it panics.
var callInMain = func(f func()) {
	panic("raylib.Main(main func()) must be called before raylib.Do(f func())")
}

func init() {
	// Make sure the main goroutine is bound to the main thread.
	runtime.LockOSThread()
}

// Main entry point. Run this function at the beginning of main(), and pass your own main body to it as a function. E.g.:
//
// 	func main() {
// 		raylib.Main(func() {
// 			// Your code here....
// 			// [....]
//
// 			// Calls to raylib can be made by any goroutine, but always guarded by raylib.Do()
// 			raylib.Do(func() {
// 				raylib.DrawTexture(..)
// 			})
// 		})
// 	}
func Main(main func()) {
	// Queue of functions that are thread-sensitive
	callQueue := make(chan func(), CallQueueCap)

	done := make(chan bool, 1)

	// Properly initialize callInMain for use by raylib.Do(..)
	callInMain = func(f func()) {
		callQueue <- func() {
			f()
			done <- true
		}
		<-done
	}

	go func() {
		main()

		close(callQueue)
	}()

	for f := range callQueue {
		f()
	}
}

// Do queues function f on the main thread and blocks until the function f finishes
func Do(f func()) {
	callInMain(f)
}
