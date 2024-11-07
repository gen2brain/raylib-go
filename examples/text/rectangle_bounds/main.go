/*******************************************************************************************
*
*   raylib [text] example - Rectangle bounds
*
*   Example originally created with raylib 2.5, last time updated with raylib 4.0
*
*   Example contributed by Vlad Adrian (@demizdor) and reviewed by Ramon Santamaria (@raysan5)
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2018-2024 Vlad Adrian (@demizdor) and Ramon Santamaria (@raysan5)
*
********************************************************************************************/
package main

import (
	"unicode/utf8"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	MeasureState = 0
	DrawState    = 1
)

func main() {

	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - draw text inside a rectangle")

	text := `Text cannot escape	this container	...word wrap also works when active so here's a long text for testing.
		
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Nec ullamcorper sit amet risus nullam eget felis eget.`

	resizing, wordWrap := false, true

	container := rl.Rectangle{
		X:      25.0,
		Y:      25.0,
		Width:  screenWidth - 50.0,
		Height: screenHeight - 250.0,
	}
	resizer := rl.Rectangle{
		X:      container.X + container.Width - 17,
		Y:      container.Y + container.Height - 17,
		Width:  14,
		Height: 14,
	}

	// Minimum width and height for the container rectangle
	minWidth := float32(60.0)
	minHeight := float32(60.0)
	maxWidth := screenWidth - float32(50.0)
	maxHeight := screenHeight - float32(160.0)

	lastMouse := rl.Vector2{}   // Stores last mouse coordinates
	borderColor := rl.Maroon    // Container border color
	font := rl.GetFontDefault() // Get default system font

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update
		if rl.IsKeyPressed(rl.KeySpace) {
			wordWrap = !wordWrap
		}

		mouse := rl.GetMousePosition()

		// Check if the mouse is inside the container and toggle border color
		if rl.CheckCollisionPointRec(mouse, container) {
			borderColor = rl.Fade(rl.Maroon, 0.4)
		} else if !resizing {
			borderColor = rl.Maroon
		}

		// Container resizing logic
		if resizing {
			if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
				resizing = false
			}

			width := container.Width + (mouse.X - lastMouse.X)
			height := container.Height + (mouse.Y - lastMouse.Y)

			container.Width = rl.Clamp(width, minWidth, maxWidth)
			container.Height = rl.Clamp(height, minHeight, maxHeight)
		} else {
			// Check if we're resizing
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(mouse, resizer) {
				resizing = true
			}
		}

		// Move resizer rectangle properly
		resizer.X = container.X + container.Width - 17
		resizer.Y = container.Y + container.Height - 17

		lastMouse = mouse // Update mouse

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawRectangleLinesEx(container, 3, borderColor) // Draw container border

		// Draw text in container (add some padding)
		DrawTextBoxed(font, text, rl.Rectangle{
			X:      container.X + 4,
			Y:      container.Y + 4,
			Width:  container.Width - 4,
			Height: container.Height - 4,
		}, 20.0, 2.0, wordWrap, rl.Gray)

		rl.DrawRectangleRec(resizer, borderColor) // Draw the resize box

		// Draw bottom info
		rl.DrawRectangle(0, screenHeight-54, screenWidth, 54, rl.Gray)
		rl.DrawRectangleRec(rl.Rectangle{
			X:      382.0,
			Y:      screenHeight - 34.0,
			Width:  12.0,
			Height: 12.0,
		}, rl.Maroon)

		rl.DrawText("Word Wrap: ", 313, screenHeight-115, 20, rl.Black)
		if wordWrap {
			rl.DrawText("ON", 447, screenHeight-115, 20, rl.Red)
		} else {
			rl.DrawText("OFF", 447, screenHeight-115, 20, rl.Black)
		}

		rl.DrawText("Press [SPACE] to toggle word wrap", 218, screenHeight-86, 20, rl.Gray)
		rl.DrawText("Click hold & drag the    to resize the container", 155, screenHeight-38, 20, rl.RayWhite)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.CloseWindow() // Close window and OpenGL context
}

//--------------------------------------------------------------------------------------
// Module functions definition
//--------------------------------------------------------------------------------------

// DrawTextBoxed draws text using font inside rectangle limits
func DrawTextBoxed(font rl.Font, text string, rec rl.Rectangle, fontSize, spacing float32, wordWrap bool, tint rl.Color) {
	DrawTextBoxedSelectable(font, text, rec, fontSize, spacing, wordWrap, tint, 0, 0, rl.White, rl.White)
}

// DrawTextBoxedSelectable draws text using font inside rectangle limits with support for text selection
func DrawTextBoxedSelectable(font rl.Font, text string, rec rl.Rectangle, fontSize, spacing float32,
	wordWrap bool, tint rl.Color, selectStart, selectLength int32, selectTint, selectBackTint rl.Color) {

	length := int32(len(text)) // Total length in bytes of the text, scanned by codepoints in loop

	// TextOffsetY : Offset between lines (on line break '\n')
	// TextOffsetX : Offset X to next character to draw
	var textOffsetY, textOffsetX float32

	scaleFactor := fontSize / float32(font.BaseSize) // Character rectangle scaling factor

	// Word/character wrapping mechanism variables
	state := DrawState
	if wordWrap {
		state = MeasureState
	}

	// StartLine : Index where to begin drawing (where a line begins)
	// EndLine : Index where to stop drawing (where a line ends)
	// LastK : Holds last value of the character position
	var startLine, endLine, lastk int32 = -1, -1, -1

	for i, k := int32(0), int32(0); i < length; i, k = i+1, k+1 {
		// Get next codepoint from byte string and glyph index in font
		codepoint, width := utf8.DecodeRuneInString(text[i:])
		codepointByteCount := int32(width)
		index := rl.GetGlyphIndex(font, codepoint)

		// NOTE: Normally we exit the decoding sequence as soon as a bad byte is found (and return 0x3f)
		// but we need to draw all the bad bytes using the '?' symbol moving one byte
		if codepoint == 0x3f {
			codepointByteCount = 1
		}
		i += codepointByteCount - 1

		var glyphWidth float32
		if codepoint != '\n' {
			chars := unsafe.Slice(font.Chars, font.CharsCount)
			if chars[index].AdvanceX == 0 {
				glyphWidth = unsafe.Slice(font.Recs, font.CharsCount)[index].Width * scaleFactor
			} else {
				glyphWidth = float32(chars[index].AdvanceX) * scaleFactor
			}

			if i+1 < length {
				glyphWidth = glyphWidth + spacing
			}
		}

		// NOTE: When wordWrap is ON we first measure how much of the text we can draw before going outside the rec container
		// We store this info in startLine and endLine, then we change states, draw the text between those two variables
		// and change states again and again recursively until the end of the text (or until we get outside the container).
		// When wordWrap is OFF we don't need the measure state so we go to the drawing state immediately
		// and begin drawing on the next line before we can get outside the container.
		if state == MeasureState {
			// TODO: There are multiple types of spaces in UNICODE, maybe it's a good idea to add support for more
			// Ref: http://jkorpela.fi/chars/spaces.html
			if (codepoint == ' ') || (codepoint == '\t') || (codepoint == '\n') {
				endLine = i
			}

			if (textOffsetX + glyphWidth) > rec.Width {
				if endLine < 1 {
					endLine = i
				}

				if i == endLine {
					endLine -= codepointByteCount
				}
				if (startLine + codepointByteCount) == endLine {
					endLine = i - codepointByteCount
				}

				state = 1 - state // Toggle state between MeasureState and DrawState
			} else if (i + 1) == length {
				endLine = i
				state = 1 - state // Toggle state between MeasureState and DrawState
			} else if codepoint == '\n' {
				state = 1 - state // Toggle state between MeasureState and DrawState
			}

			if state == DrawState {
				textOffsetX = 0
				i = startLine
				glyphWidth = 0

				// Save character position when we switch states
				tmp := lastk
				lastk = k - 1
				k = tmp
			}
		} else {
			if codepoint == '\n' {
				if !wordWrap {
					textOffsetY += float32(font.BaseSize+font.BaseSize/2) * scaleFactor
					textOffsetX = 0
				}
			} else {
				if !wordWrap && ((textOffsetX + glyphWidth) > rec.Width) {
					textOffsetY += float32(font.BaseSize+font.BaseSize/2) * scaleFactor
					textOffsetX = 0
				}

				// When text overflows rectangle height limit, just stop drawing
				if (textOffsetY + float32(font.BaseSize)*scaleFactor) > rec.Height {
					break
				}

				// Draw selection background
				isGlyphSelected := false
				if (selectStart >= 0) && (k >= selectStart) && (k < (selectStart + selectLength)) {
					rl.DrawRectangleRec(rl.Rectangle{
						X:      rec.X + textOffsetX - 1,
						Y:      rec.Y + textOffsetY,
						Width:  glyphWidth,
						Height: float32(font.BaseSize) * scaleFactor,
					}, selectBackTint)
					isGlyphSelected = true
				}

				// Draw current character glyph
				if (codepoint != ' ') && (codepoint != '\t') {
					col := tint
					if isGlyphSelected {
						col = selectTint
					}
					pos := rl.Vector2{
						X: rec.X + textOffsetX,
						Y: rec.Y + textOffsetY,
					}
					rl.DrawTextEx(font, string(codepoint), pos, fontSize, 0, col)
				}
			}

			if wordWrap && (i == endLine) {
				textOffsetY += float32(font.BaseSize+font.BaseSize/2) * scaleFactor
				textOffsetX = 0
				startLine = endLine
				endLine = -1
				glyphWidth = 0
				selectStart += lastk - k
				k = lastk

				state = 1 - state // Toggle state between MeasureState and DrawState
			}
		}

		if (textOffsetX != 0) || (codepoint != ' ') { // avoid leading spaces
			textOffsetX += glyphWidth
		}
	}
}
