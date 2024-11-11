/*******************************************************************************************
*
*   raylib [text] example - Draw 3d
*
*   NOTE: Draw a 2D text in 3D space, each letter is drawn in a quad (or 2 quads if back face is set)
*   where the texture coordinates of each quad map to the texture coordinates of the glyphs
*   inside the font texture.
*
*   A more efficient approach, i believe, would be to render the text in a render texture and
*   map that texture to a plane and render that, or maybe a shader but my method allows more
*   flexibility...for example to change position of each letter individually to make some think
*   like a wavy text effect.
*
*   Special thanks to:
*        @Nighten for the DrawTextStyle() code https://github.com/NightenDushi/Raylib_DrawTextStyle
*        Chris Camacho (codifies - http://bedroomcoders.co.uk/) for the alpha discard shader
*
*   Example originally created with raylib 3.5, last time updated with raylib 4.0
*
*   Example contributed by Vlad Adrian (@demizdor) and reviewed by Ramon Santamaria (@raysan5)
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2021-2024 Vlad Adrian (@demizdor)
*
********************************************************************************************/
package main

import (
	"fmt"
	"math"
	"path/filepath"
	"unicode/utf8"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Globals
const (
	letterBoundarySize = 0.25
	textMaxLayers      = 32
	maxTextLength      = 64
	screenWidth        = 800
	screenHeight       = 450
)

var letterBoundaryColor = rl.Violet
var showLetterBoundary, showTextBoundary = false, false

// Data Types definition

// WaveTextConfig is a configuration structure for waving the text
type WaveTextConfig struct {
	waveRange, waveSpeed, waveOffset rl.Vector3
}

// Program main entry point
func main() {
	// Initialization
	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - draw 2D text in 3D")

	spin := true        // Spin the camera?
	multicolor := false // Multicolor mode

	// Define the camera to look into our 3d world
	camera := rl.Camera3D{
		Position: rl.Vector3{
			X: -10.0,
			Y: 15.0,
			Z: -10.0,
		}, // Camera position
		Target:     rl.Vector3{},         // Camera looking at point
		Up:         rl.Vector3{Y: 1.0},   // Camera up vector (rotation towards target)
		Fovy:       45.0,                 // Camera field-of-view Y
		Projection: rl.CameraPerspective, // Camera projection type
	}

	cameraMode := rl.CameraOrbital

	cubePosition := rl.Vector3{Y: 1.0}
	cubeSize := rl.Vector3{
		X: 2.0,
		Y: 2.0,
		Z: 2.0,
	}

	// Use the default font
	font := rl.GetFontDefault()
	var fontSize, fontSpacing, lineSpacing float32 = 8.0, 0.5, -1

	// Set the text (using markdown!)
	text := "Hello ~~World~~ in 3D!"
	tbox := rl.Vector3{}
	var layers, quads int32 = 1, 0
	var layerDistance float32 = 0.01

	wcfg := WaveTextConfig{
		waveSpeed: rl.Vector3{
			X: 3,
			Y: 3,
			Z: 0.5,
		},
		waveOffset: rl.Vector3{
			X: 0.35,
			Y: 0.35,
			Z: 0.35,
		},
		waveRange: rl.Vector3{
			X: 0.45,
			Y: 0.45,
			Z: 0.45,
		},
	}

	var time float32

	// Set up a light and dark color
	light, dark := rl.Maroon, rl.Red

	// Load the alpha discard shader
	alphaDiscard := rl.LoadShader("", "alpha_discard.fs")

	// Array filled with multiple random colors (when multicolor mode is set)
	var multi [textMaxLayers]rl.Color

	rl.DisableCursor() // Limit cursor to relative movement inside the window

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		rl.UpdateCamera(&camera, cameraMode)

		// Handle font files dropped
		if rl.IsFileDropped() {
			droppedFiles := rl.LoadDroppedFiles()

			// NOTE: We only support first ttf file dropped
			rl.UnloadFont(font)
			if filepath.Ext(droppedFiles[0]) == ".ttf" {
				font = rl.LoadFontEx(droppedFiles[0], int32(fontSize), nil, 0)
			} else if filepath.Ext(droppedFiles[0]) == ".fnt" {
				font = rl.LoadFont(droppedFiles[0])
				fontSize = float32(font.BaseSize)
			}

			rl.UnloadDroppedFiles() // Unload file paths from memory
		}

		// Handle Events
		if rl.IsKeyPressed(rl.KeyF1) {
			showLetterBoundary = !showLetterBoundary
		}
		if rl.IsKeyPressed(rl.KeyF2) {
			showTextBoundary = !showTextBoundary
		}
		if rl.IsKeyPressed(rl.KeyF3) {
			// Handle camera change
			spin = !spin
			// we need to reset the camera when changing modes
			camera = rl.Camera3D{
				Target:     rl.Vector3{},         // Camera looking at point
				Up:         rl.Vector3{Y: 1.0},   // Camera up vector (rotation towards target)
				Fovy:       45.0,                 // Camera field-of-view Y
				Projection: rl.CameraPerspective, // Camera projection type
			}

			if spin {
				camera.Position = rl.Vector3{
					X: -10.0,
					Y: 15.0,
					Z: -10.0,
				} // Camera position
				cameraMode = rl.CameraOrbital
			} else {
				camera.Position = rl.Vector3{
					X: 10.0,
					Y: 10.0,
					Z: -10.0,
				} // Camera position
				cameraMode = rl.CameraFree
			}
		}

		// Handle clicking the cube
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			// TODO : Missing function, see issue https://github.com/gen2brain/raylib-go/issues/457
			//ray := rl.GetScreenToWorldRay(rl.GetMousePosition(), camera)
			ray := rl.GetMouseRay(rl.GetMousePosition(), camera)

			// Check collision between ray and box
			v1 := rl.Vector3{
				X: cubePosition.X - cubeSize.X/2,
				Y: cubePosition.Y - cubeSize.Y/2,
				Z: cubePosition.Z - cubeSize.Z/2,
			}
			v2 := rl.Vector3{
				X: cubePosition.X + cubeSize.X/2,
				Y: cubePosition.Y + cubeSize.Y/2,
				Z: cubePosition.Z + cubeSize.Z/2,
			}
			collision := rl.GetRayCollisionBox(ray, rl.BoundingBox{
				Min: v1,
				Max: v2,
			})

			if collision.Hit {
				// Generate new random colors
				light = generateRandomColor(0.5, 0.78)
				dark = generateRandomColor(0.4, 0.58)
			}
		}

		// Handle text layers changes
		if rl.IsKeyPressed(rl.KeyHome) && layers > 1 {
			layers--
		} else if rl.IsKeyPressed(rl.KeyEnd) && layers < textMaxLayers {
			layers++
		}

		// Handle text changes
		if rl.IsKeyPressed(rl.KeyLeft) {
			fontSize -= 0.5
		} else if rl.IsKeyPressed(rl.KeyRight) {
			fontSize += 0.5
		} else if rl.IsKeyPressed(rl.KeyUp) {
			fontSpacing -= 0.1
		} else if rl.IsKeyPressed(rl.KeyDown) {
			fontSpacing += 0.1
		} else if rl.IsKeyPressed(rl.KeyPageUp) {
			lineSpacing -= 0.5
		} else if rl.IsKeyPressed(rl.KeyPageDown) {
			lineSpacing += 0.5
		} else if rl.IsKeyDown(rl.KeyInsert) {
			layerDistance -= 0.001
		} else if rl.IsKeyDown(rl.KeyDelete) {
			layerDistance += 0.001
		} else if rl.IsKeyPressed(rl.KeyTab) {
			multicolor = !multicolor // Enable /disable multicolor mode

			if multicolor {
				// Fill color array with random colors
				for i := 0; i < textMaxLayers; i++ {
					multi[i] = generateRandomColor(0.5, 0.8)
					multi[i].A = uint8(rl.GetRandomValue(0, 255))
				}
			}
		}

		// Handle text input
		ch := rl.GetCharPressed()
		if rl.IsKeyPressed(rl.KeyBackspace) {
			// Remove last char
			text = text[:len(text)-1]
		} else if rl.IsKeyPressed(rl.KeyEnter) {
			// handle newline
			text += "\n"
		} else if ch != 0 && len(text) < maxTextLength {
			// append only printable chars
			text += string(ch)
		}

		// Measure 3D text so we can center it
		tbox = measureTextWave3D(font, text, fontSize, fontSpacing, lineSpacing)

		quads = 0                 // Reset quad counter
		time += rl.GetFrameTime() // Update timer needed by `drawTextWave3D()`

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)
		rl.DrawCubeV(cubePosition, cubeSize, dark)
		rl.DrawCubeWires(cubePosition, 2.1, 2.1, 2.1, light)

		rl.DrawGrid(10, 2.0)

		// Use a shader to handle the depth buffer issue with transparent textures
		// NOTE: more info at https://bedroomcoders.co.uk/raylib-billboards-advanced-use/
		rl.BeginShaderMode(alphaDiscard)

		// Draw the 3D text above the red cube
		rl.PushMatrix()
		rl.Rotatef(90.0, 1.0, 0.0, 0.0)
		rl.Rotatef(90.0, 0.0, 0.0, -1.0)

		for i := int32(0); i < layers; i++ {
			clr := light
			if multicolor {
				clr = multi[i]
			}

			vec := rl.Vector3{
				X: -tbox.X / 2.0,
				Y: layerDistance * float32(i),
				Z: -4.5,
			}
			drawTextWave3D(font, text, vec, fontSize, fontSpacing, lineSpacing, true, &wcfg, time, clr)
		}

		// Draw the text boundary if set
		if showTextBoundary {
			rl.DrawCubeWiresV(rl.Vector3{Z: -4.5 + tbox.Z/2}, tbox, dark)
		}
		rl.PopMatrix()

		// Don't draw the letter boundaries for the 3D text below
		slb := showLetterBoundary
		showLetterBoundary = false

		// Draw 3D options (use default font)
		rl.PushMatrix()
		rl.Rotatef(180.0, 0.0, 1.0, 0.0)
		opt := fmt.Sprintf("< SIZE: %2.1f >", fontSize)
		quads += int32(len(opt))
		m := measureText3D(rl.GetFontDefault(), opt, 8.0, 1.0, 0.0)
		pos := rl.Vector3{
			X: -m.X / 2.0,
			Y: 0.01,
			Z: 2.0,
		}
		drawText3D(rl.GetFontDefault(), opt, pos, 8.0, 1.0, 0.0, false, rl.Blue)
		pos.Z += 0.5 + m.Z

		opt = fmt.Sprintf("< SPACING: %2.1f >", fontSpacing)
		quads += int32(len(opt))
		m = measureText3D(rl.GetFontDefault(), opt, 8.0, 1.0, 0.0)
		pos.X = -m.X / 2.0
		drawText3D(rl.GetFontDefault(), opt, pos, 8.0, 1.0, 0.0, false, rl.Blue)
		pos.Z += 0.5 + m.Z

		opt = fmt.Sprintf("< LINE: %2.1f >", lineSpacing)
		quads += int32(len(opt))
		m = measureText3D(rl.GetFontDefault(), opt, 8.0, 1.0, 0.0)
		pos.X = -m.X / 2.0
		drawText3D(rl.GetFontDefault(), opt, pos, 8.0, 1.0, 0.0, false, rl.Blue)
		pos.Z += 1.0 + m.Z

		lbox := "OFF"
		if slb {
			lbox = "ON"
		}
		opt = fmt.Sprintf("< LBOX: %3s >", lbox)
		quads += int32(len(opt))
		m = measureText3D(rl.GetFontDefault(), opt, 8.0, 1.0, 0.0)
		pos.X = -m.X / 2.0
		drawText3D(rl.GetFontDefault(), opt, pos, 8.0, 1.0, 0.0, false, rl.Red)
		pos.Z += 0.5 + m.Z

		tb := "OFF"
		if showTextBoundary {
			tb = "ON"
		}
		opt = fmt.Sprintf("< TBOX: %3s >", tb)
		quads += int32(len(opt))
		m = measureText3D(rl.GetFontDefault(), opt, 8.0, 1.0, 0.0)
		pos.X = -m.X / 2.0
		drawText3D(rl.GetFontDefault(), opt, pos, 8.0, 1.0, 0.0, false, rl.Red)
		pos.Z += 0.5 + m.Z

		opt = fmt.Sprintf("< LAYER DISTANCE: %.3f >", layerDistance)
		quads += int32(len(opt))
		m = measureText3D(rl.GetFontDefault(), opt, 8.0, 1.0, 0.0)
		pos.X = -m.X / 2.0
		drawText3D(rl.GetFontDefault(), opt, pos, 8.0, 1.0, 0.0, false, rl.DarkPurple)
		rl.PopMatrix()

		// Draw 3D info text (use default font)
		opt = "All the text displayed here is in 3D"
		quads += 36
		m = measureText3D(rl.GetFontDefault(), opt, 10.0, 0.5, 0.0)
		pos = rl.Vector3{
			X: -m.X / 2.0,
			Y: 0.01,
			Z: 2.0,
		}
		drawText3D(rl.GetFontDefault(), opt, pos, 10.0, 0.5, 0.0, false, rl.DarkBlue)
		pos.Z += 1.5 + m.Z

		opt = "press [Left]/[Right] to change the font size"
		quads += 44
		m = measureText3D(rl.GetFontDefault(), opt, 6.0, 0.5, 0.0)
		pos.X = -m.X / 2.0
		drawText3D(rl.GetFontDefault(), opt, pos, 6.0, 0.5, 0.0, false, rl.DarkBlue)
		pos.Z += 0.5 + m.Z

		opt = "press [Up]/[Down] to change the font spacing"
		quads += 44
		m = measureText3D(rl.GetFontDefault(), opt, 6.0, 0.5, 0.0)
		pos.X = -m.X / 2.0
		drawText3D(rl.GetFontDefault(), opt, pos, 6.0, 0.5, 0.0, false, rl.DarkBlue)
		pos.Z += 0.5 + m.Z

		opt = "press [PgUp]/[PgDown] to change the line spacing"
		quads += 48
		m = measureText3D(rl.GetFontDefault(), opt, 6.0, 0.5, 0.0)
		pos.X = -m.X / 2.0
		drawText3D(rl.GetFontDefault(), opt, pos, 6.0, 0.5, 0.0, false, rl.DarkBlue)
		pos.Z += 0.5 + m.Z

		opt = "press [F1] to toggle the letter boundary"
		quads += 39
		m = measureText3D(rl.GetFontDefault(), opt, 6.0, 0.5, 0.0)
		pos.X = -m.X / 2.0
		drawText3D(rl.GetFontDefault(), opt, pos, 6.0, 0.5, 0.0, false, rl.DarkBlue)
		pos.Z += 0.5 + m.Z

		opt = "press [F2] to toggle the text boundary"
		quads += 37
		m = measureText3D(rl.GetFontDefault(), opt, 6.0, 0.5, 0.0)
		pos.X = -m.X / 2.0
		drawText3D(rl.GetFontDefault(), opt, pos, 6.0, 0.5, 0.0, false, rl.DarkBlue)

		showLetterBoundary = slb
		rl.EndShaderMode()

		rl.EndMode3D()

		// Draw 2D info text & stats
		msg := `Drag & drop a font file to change the font!
Type something, see what happens!

Press [F3] to toggle the camera`
		rl.DrawText(msg, 10, 35, 10, rl.Black)

		cam := "FREE"
		if spin {
			cam = "ORBITAL"
		}
		quads += int32(len(text)) * 2 * layers
		tmp := fmt.Sprintf("%2d layer(s) | %s camera | %4d quads (%4d verts)", layers, cam, quads, quads*4)
		width := rl.MeasureText(tmp, 10)
		rl.DrawText(tmp, screenWidth-20-width, 10, 10, rl.DarkGreen)

		tmp = "[Home]/[End] to add/remove 3D text layers"
		width = rl.MeasureText(tmp, 10)
		rl.DrawText(tmp, screenWidth-20-width, 25, 10, rl.DarkGray)

		tmp = "[Insert]/[Delete] to increase/decrease distance between layers"
		width = rl.MeasureText(tmp, 10)
		rl.DrawText(tmp, screenWidth-20-width, 40, 10, rl.DarkGray)

		tmp = "click the [CUBE] for a random color"
		width = rl.MeasureText(tmp, 10)
		rl.DrawText(tmp, screenWidth-20-width, 55, 10, rl.DarkGray)

		tmp = "[Tab] to toggle multicolor mode"
		width = rl.MeasureText(tmp, 10)
		rl.DrawText(tmp, screenWidth-20-width, 70, 10, rl.DarkGray)

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadFont(font)
	rl.CloseWindow() // Close window and OpenGL context

}

// Module Functions Definitions

// drawTextCodepoint3D draws codepoint at specified position in 3D space
func drawTextCodepoint3D(font rl.Font, codepoint rune, position rl.Vector3, fontSize float32, backface bool,
	tint rl.Color) {
	// Character index position in sprite font
	// NOTE: In case a codepoint is not available in the font, index returned points to '?'
	index := rl.GetGlyphIndex(font, codepoint)
	scale := fontSize / float32(font.BaseSize)

	// Character destination rectangle on screen
	// NOTE: We consider charsPadding on drawing
	glyphs := unsafe.Slice(font.Chars, font.CharsCount)
	position.X += float32(glyphs[index].OffsetX-font.CharsPadding) / float32(font.BaseSize) * scale
	position.Z += float32(glyphs[index].OffsetY-font.CharsPadding) / float32(font.BaseSize) * scale

	// Character source rectangle from font texture atlas
	// NOTE: We consider chars padding when drawing, it could be required for outline/glow shader effects
	recs := unsafe.Slice(font.Recs, font.CharsCount)
	srcRec := rl.Rectangle{
		X:      recs[index].X - float32(font.CharsPadding),
		Y:      recs[index].Y - float32(font.CharsPadding),
		Width:  recs[index].Width + 2.0*float32(font.CharsPadding),
		Height: recs[index].Height + 2.0*float32(font.CharsPadding),
	}

	width := (recs[index].Width + 2.0*float32(font.CharsPadding)) / float32(font.BaseSize) * scale
	height := (recs[index].Height + 2.0*float32(font.CharsPadding)) / float32(font.BaseSize) * scale

	if font.Texture.ID > 0 {
		var x, y, z float32

		// normalized texture coordinates of the glyph inside the font texture (0.0f -> 1.0f)
		tx := srcRec.X / float32(font.Texture.Width)
		ty := srcRec.Y / float32(font.Texture.Height)
		tw := (srcRec.X + srcRec.Width) / float32(font.Texture.Width)
		th := (srcRec.Y + srcRec.Height) / float32(font.Texture.Height)

		if showLetterBoundary {
			pos := rl.Vector3{
				X: position.X + width/2,
				Y: position.Y,
				Z: position.Z + height/2,
			}
			size := rl.Vector3{
				X: width,
				Y: letterBoundarySize,
				Z: height,
			}
			rl.DrawCubeWiresV(pos, size, letterBoundaryColor)
		}

		var limit int32 = 4
		if backface {
			limit = 8
		}
		rl.CheckRenderBatchLimit(limit)
		rl.SetTexture(font.Texture.ID)

		rl.PushMatrix()
		rl.Translatef(position.X, position.Y, position.Z)

		rl.Begin(rl.Quads)
		rl.Color4ub(tint.R, tint.G, tint.B, tint.A)

		// Front Face
		rl.Normal3f(0.0, 1.0, 0.0) // Normal Pointing Up
		rl.TexCoord2f(tx, ty)
		rl.Vertex3f(x, y, z) // Top Left Of The Texture and Quad
		rl.TexCoord2f(tx, th)
		rl.Vertex3f(x, y, z+height) // Bottom Left Of The Texture and Quad
		rl.TexCoord2f(tw, th)
		rl.Vertex3f(x+width, y, z+height) // Bottom Right Of The Texture and Quad
		rl.TexCoord2f(tw, ty)
		rl.Vertex3f(x+width, y, z) // Top Right Of The Texture and Quad

		if backface {
			// Back Face
			rl.Normal3f(0.0, -1.0, 0.0) // Normal Pointing Down
			rl.TexCoord2f(tx, ty)
			rl.Vertex3f(x, y, z) // Top Right Of The Texture and Quad
			rl.TexCoord2f(tw, ty)
			rl.Vertex3f(x+width, y, z) // Top Left Of The Texture and Quad
			rl.TexCoord2f(tw, th)
			rl.Vertex3f(x+width, y, z+height) // Bottom Left Of The Texture and Quad
			rl.TexCoord2f(tx, th)
			rl.Vertex3f(x, y, z+height) // Bottom Right Of The Texture and Quad
		}
		rl.End()
		rl.PopMatrix()
		rl.SetTexture(0)
	}
}

// drawText3D draws a 2D text in 3D space
func drawText3D(font rl.Font, text string, position rl.Vector3, fontSize, fontSpacing,
	lineSpacing float32, backface bool, tint rl.Color) {
	length := int32(len(text)) // Total length in bytes of the text,
	// scanned by codepoints in loop

	// TextOffsetY : Offset between lines (on line break '\n')
	// TextOffsetX : Offset X to next character to draw
	var textOffsetY, textOffsetX float32

	scale := fontSize / float32(font.BaseSize)
	for i := int32(0); i < length; {

		// Get next codepoint from byte string and glyph index in font
		codepoint, codepointByteCount := getCodepoint(text, i)
		index := rl.GetGlyphIndex(font, codepoint)

		// NOTE: Normally we exit the decoding sequence as soon as a bad byte is found (and return 0x3f)
		// but we need to draw all the bad bytes using the '?' symbol moving one byte
		if codepoint == 0x3f {
			codepointByteCount = 1
		}

		if codepoint == 0 || codepoint == '\n' {
			// NOTE: Fixed line spacing of 1.5 line-height
			// TODO: Support custom line spacing defined by user
			textOffsetY += scale + lineSpacing/float32(font.BaseSize)*scale
			textOffsetX = 0.0
		} else {
			if (codepoint != ' ') && (codepoint != '\t') {
				vec := rl.Vector3{
					X: position.X + textOffsetX,
					Y: position.Y,
					Z: position.Z + textOffsetY,
				}
				drawTextCodepoint3D(font, codepoint, vec, fontSize, backface, tint)
			}

			textOffsetX += getTextWidth(font, index, fontSpacing, scale)
		}

		i += codepointByteCount // Move text bytes counter to next codepoint
	}
}

// measureText3D measures a text in 3D. For some reason `MeasureTextEx()`
// just doesn't seem to work, so I had to use this instead.
func measureText3D(font rl.Font, text string, fontSize, fontSpacing, lineSpacing float32) rl.Vector3 {
	length := int32(len(text))
	var tempLen, lenCounter int32 // Used to count longer text line num chars

	var (
		tempTextWidth float32 // Used to count longer text line width
		scale         = fontSize / float32(font.BaseSize)
		textHeight    = scale
		textWidth     float32
	)

	var (
		letter rune  // Current character
		index  int32 // Index position in sprite font
	)

	for i := int32(0); i < length; i++ {
		lenCounter++

		var next int32
		letter, next = getCodepoint(text, i)
		index = rl.GetGlyphIndex(font, letter)

		// NOTE: normally we exit the decoding sequence as soon as a bad byte is found (and return 0x3f)
		// but we need to draw all the bad bytes using the '?' symbol so to not skip any we set next = 1
		if letter == 0x3f {
			next = 1
		}
		i += next - 1

		if letter != 0 && letter != '\n' {
			textWidth += getTextWidth(font, index, fontSpacing, scale)
		} else {
			if tempTextWidth < textWidth {
				tempTextWidth = textWidth
			}
			lenCounter = 0
			textWidth = 0.0
			textHeight += scale + lineSpacing/float32(font.BaseSize)*scale
		}

		if tempLen < lenCounter {
			tempLen = lenCounter
		}
	}

	if tempTextWidth < textWidth {
		tempTextWidth = textWidth
	}

	// Adds chars spacing to measure
	vec := rl.Vector3{
		X: tempTextWidth + float32(tempLen-1)*fontSpacing/float32(font.BaseSize)*scale,
		Y: 0.25,
		Z: textHeight,
	}
	return vec
}

// drawTextWave3D draws a 2D text in 3D space and wave the parts that start with `~~` and end with `~~`.
// This is a modified version of the original code by @Nighten found here https://github.com/NightenDushi/Raylib_DrawTextStyle
func drawTextWave3D(font rl.Font, text string, position rl.Vector3, fontSize, fontSpacing,
	lineSpacing float32, backface bool, config *WaveTextConfig, time float32, tint rl.Color) {
	length := int32(len(text)) // Total length in bytes of the text, scanned by codepoints in loop

	// TextOffsetY : Offset between lines (on line break '\n')
	// TextOffsetX : Offset X to next character to draw
	var textOffsetY, textOffsetX float32
	var wave bool

	scale := fontSize / float32(font.BaseSize)

	for i, k := int32(0), int32(0); i < length; k++ {
		// Get next codepoint from byte string and glyph index in font
		codepoint, codepointByteCount := getCodepoint(text, i)
		index := rl.GetGlyphIndex(font, codepoint)

		// NOTE: Normally we exit the decoding sequence as soon as a bad byte is found (and return 0x3f)
		// but we need to draw all the bad bytes using the '?' symbol moving one byte
		if codepoint == 0x3f {
			codepointByteCount = 1
		}

		if codepoint == 0 || codepoint == '\n' {
			// NOTE: Fixed line spacing of 1.5 line-height
			// TODO: Support custom line spacing defined by user
			textOffsetY += scale + lineSpacing/float32(font.BaseSize)*scale
			textOffsetX = 0.0
			k = 0
		} else if codepoint == '~' {
			var r rune
			r, codepointByteCount = getCodepoint(text, i+1)
			if r == '~' {
				codepointByteCount += 1
				wave = !wave
			}
		} else {
			if (codepoint != ' ') && (codepoint != '\t') {
				pos := position
				if wave { // Apply the wave effect
					kk := float32(k)
					pos.X += sin(time*config.waveSpeed.X-kk*config.waveOffset.X) * config.waveRange.X
					pos.Y += sin(time*config.waveSpeed.Y-kk*config.waveOffset.Y) * config.waveRange.Y
					pos.Z += sin(time*config.waveSpeed.Z-kk*config.waveOffset.Z) * config.waveRange.Z
				}

				vec := rl.Vector3{
					X: pos.X + textOffsetX,
					Y: pos.Y,
					Z: pos.Z + textOffsetY,
				}
				drawTextCodepoint3D(font, codepoint, vec, fontSize, backface, tint)
			}

			textOffsetX += getTextWidth(font, index, fontSpacing, scale)
		}

		i += codepointByteCount // Move text bytes counter to next codepoint
	}
}

// measureTextWave3D measures a text in 3D ignoring the `~~` chars.
func measureTextWave3D(font rl.Font, text string, fontSize, fontSpacing, lineSpacing float32) rl.Vector3 {
	length := int32(len(text))
	var tempLen, lenCounter int32 // Used to count longer text line num chars

	var (
		tempTextWidth float32 = 0.0 // Used to count longer text line width
		scale                 = fontSize / float32(font.BaseSize)
		textHeight            = scale
		textWidth     float32 = 0.0
	)

	var letter, index int32 // Current character and Index position in sprite font

	for i := int32(0); i < length; i++ {
		lenCounter++
		var next int32

		letter, next = getCodepoint(text, i)
		index = rl.GetGlyphIndex(font, letter)

		// NOTE: normally we exit the decoding sequence as soon as a bad byte is found (and return 0x3f)
		// but we need to draw all the bad bytes using the '?' symbol so to not skip any we set next = 1
		if letter == 0x3f {
			next = 1
		}
		i += next - 1

		if letter != 0 && letter != '\n' {
			r, _ := getCodepoint(text, i+1)
			if letter == '~' && r == '~' {
				i++
			} else {
				textWidth += getTextWidth(font, index, fontSpacing, scale)
			}
		} else {
			if tempTextWidth < textWidth {
				tempTextWidth = textWidth
			}
			lenCounter = 0
			textWidth = 0.0
			textHeight += scale + lineSpacing/float32(font.BaseSize)*scale
		}

		if tempLen < lenCounter {
			tempLen = lenCounter
		}
	}

	if tempTextWidth < textWidth {
		tempTextWidth = textWidth
	}

	vec := rl.Vector3{
		X: tempTextWidth + float32(tempLen-1)*fontSpacing/float32(font.BaseSize)*scale, // Adds chars spacing to measure
		Y: 0.25,
		Z: textHeight,
	}

	return vec
}

// generateRandomColor generates a nice color with a random hue
func generateRandomColor(s, v float32) rl.Color {
	const Phi = float64(0.618033988749895) // Golden ratio conjugate
	h := float64(rl.GetRandomValue(0, 360))
	h = math.Mod(h+h*Phi, 360.0)
	return rl.ColorFromHSV(float32(h), s, v)
}

// getCodepoint returns the rune starting at index, plus the length of that rune in bytes
func getCodepoint(s string, index int32) (rune, int32) {
	if index == int32(len(s)) {
		return 0, 0
	}
	r := []rune(s[index:])
	return r[0], int32(utf8.RuneLen(r[0]))
}

// sin is just a convenience function to avoid a bunch of type conversions
func sin(value float32) float32 {
	return float32(math.Sin(float64(value)))
}

func getTextWidth(font rl.Font, index int32, spacing, scale float32) float32 {
	glyphs := unsafe.Slice(font.Chars, font.CharsCount)
	if glyphs[index].AdvanceX == 0 {
		recs := unsafe.Slice(font.Recs, font.CharsCount)
		return (recs[index].Width + spacing) / float32(font.BaseSize) * scale
	} else {
		return (float32(glyphs[index].AdvanceX) + spacing) / float32(font.BaseSize) * scale
	}
}
