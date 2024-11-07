/*******************************************************************************************
*
*   raylib [text] example - Unicode
*
*   Example originally created with raylib 2.5, last time updated with raylib 4.0
*
*   Example contributed by Vlad Adrian (@demizdor) and reviewed by Ramon Santamaria (@raysan5)
*
*   Example licensed under an unmodified zlib/libpng license, which is an OSI-certified,
*   BSD-like license that allows static linking with closed source software
*
*   Copyright (c) 2019-2024 Vlad Adrian (@demizdor) and Ramon Santamaria (@raysan5)
*
********************************************************************************************/
package main

import (
	"fmt"
	"unicode/utf8"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth    = 800
	screenHeight   = 450
	emojiPerWidth  = 8
	emojiPerHeight = 4
	MeasureState   = 0
	DrawState      = 1
)

type Messages struct {
	text     string
	language string
}

type Emoji struct {
	index   int32    // Index inside `emojiCodepoints`
	message int32    // Message index
	color   rl.Color // Emoji color
}

// Arrays that holds the random emojis
var emoji [emojiPerWidth * emojiPerHeight]Emoji
var hovered, selected int32 = -1, -1

func main() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "raylib [text] example - unicode")

	// Load the font resources
	// NOTE: fontAsian is for asian languages,
	// fontEmoji is the emojis and fontDefault is used for everything else
	fontDefault := rl.LoadFont("resources/dejavu.fnt")
	fontAsian := rl.LoadFont("resources/noto_cjk.fnt")
	fontEmoji := rl.LoadFont("resources/symbola.fnt")

	hoveredPos := rl.Vector2{}
	selectedPos := rl.Vector2{}

	RandomizeEmoji()

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update

		// Add a new set of emojis when SPACE is pressed
		if rl.IsKeyPressed(rl.KeySpace) {
			RandomizeEmoji()
		}

		// Set the selected emoji
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && (hovered != -1) && (hovered != selected) {
			selected = hovered
			selectedPos = hoveredPos
		}

		mouse := rl.GetMousePosition()
		position := rl.Vector2{
			X: 28.8,
			Y: 10.0,
		}
		hovered = -1

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw random emojis in the background
		for i := int32(0); i < emojiPerWidth*emojiPerHeight; i++ {
			txt := emojiCodepoints[emoji[i].index : emoji[i].index+4]
			emojiRect := rl.Rectangle{
				X:      position.X,
				Y:      position.Y,
				Width:  float32(fontEmoji.BaseSize),
				Height: float32(fontEmoji.BaseSize),
			}

			if !rl.CheckCollisionPointRec(mouse, emojiRect) {
				col := rl.Fade(rl.LightGray, 0.4)
				if selected == i {
					col = emoji[i].color
				}
				rl.DrawTextEx(fontEmoji, txt, position, float32(fontEmoji.BaseSize), 1.0, col)
			} else {
				rl.DrawTextEx(fontEmoji, txt, position, float32(fontEmoji.BaseSize), 1.0, emoji[i].color)
				hovered = i
				hoveredPos = position
			}

			if (i != 0) && (i%emojiPerWidth == 0) {
				position.Y += float32(fontEmoji.BaseSize) + 24.25
				position.X = 28.8
			} else {
				position.X += float32(fontEmoji.BaseSize) + 28.8
			}
		}

		// Draw the message when an emoji is selected
		if selected != -1 {
			message := emoji[selected].message
			horizontalPadding := 20
			verticalPadding := 30
			font := fontDefault

			// Set correct font for asian languages
			if messages[message].language == "Chinese" ||
				messages[message].language == "Korean" ||
				messages[message].language == "Japanese" {
				font = fontAsian
			}

			// Calculate size for the message box (approximate the height and width)
			sz := rl.MeasureTextEx(font, messages[message].text, float32(font.BaseSize), 1.0)
			if sz.X > 300 {
				sz.Y *= sz.X / 300
				sz.X = 300
			} else if sz.X < 160 {
				sz.X = 160
			}

			msgRect := rl.Rectangle{
				X:      selectedPos.X - 38.8,
				Y:      selectedPos.Y,
				Width:  float32(2*horizontalPadding) + sz.X,
				Height: float32(2*verticalPadding) + sz.Y,
			}
			msgRect.Y -= msgRect.Height

			// Coordinates for the chat bubble triangle
			a := rl.Vector2{
				X: selectedPos.X,
				Y: msgRect.Y + msgRect.Height,
			}
			b := rl.Vector2{
				X: a.X + 8,
				Y: a.Y + 10,
			}
			c := rl.Vector2{
				X: a.X + 10,
				Y: a.Y,
			}

			// Don't go outside the screen
			if msgRect.X < 10 {
				msgRect.X += 28
			}
			if msgRect.Y < 10 {
				msgRect.Y = selectedPos.Y + 84
				a.Y = msgRect.Y
				c.Y = a.Y
				b.Y = a.Y - 10

				// Swap values so we can actually render the triangle :(
				a, b = b, a
			}

			if msgRect.X+msgRect.Width > screenWidth {
				msgRect.X -= (msgRect.X + msgRect.Width) - screenWidth + 10
			}

			// Draw chat bubble
			rl.DrawRectangleRec(msgRect, emoji[selected].color)
			rl.DrawTriangle(a, b, c, emoji[selected].color)

			// Draw the main text message
			textRect := rl.Rectangle{
				X:      msgRect.X + float32(horizontalPadding)/2,
				Y:      msgRect.Y + float32(verticalPadding)/2,
				Width:  msgRect.Width - float32(horizontalPadding),
				Height: msgRect.Height,
			}
			DrawTextBoxed(font, messages[message].text, textRect, float32(font.BaseSize), 1.0, true, rl.White)

			// Draw the info text below the main message
			size := len(messages[message].text)
			length := utf8.RuneCountInString(messages[message].text)
			info := fmt.Sprintf("%s %d characters %d bytes", messages[message].language, length, size)
			sz = rl.MeasureTextEx(rl.GetFontDefault(), info, 10, 1.0)

			rl.DrawText(info, int32(textRect.X+textRect.Width-sz.X), int32(msgRect.Y+msgRect.Height-sz.Y-2), 10,
				rl.RayWhite)

		}

		// Draw the info text
		rl.DrawText("These emojis have something to tell you, click each to find out!", (screenWidth-650)/2,
			screenHeight-40, 20, rl.Gray)
		rl.DrawText("Each emoji is a unicode character from a font, not a texture... Press [SPACEBAR] to refresh",
			(screenWidth-484)/2, screenHeight-16, 10, rl.Gray)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadFont(fontDefault) // Unload font resource
	rl.UnloadFont(fontAsian)   // Unload font resource
	rl.UnloadFont(fontEmoji)   // Unload font resource

	rl.CloseWindow() // Close window and OpenGL context
}

// RandomizeEmoji fills the emoji array with random emoji (only those emojis present in fontEmoji)
func RandomizeEmoji() {
	hovered, selected = -1, -1
	start := rl.GetRandomValue(45, 360)

	for i := int32(0); i < emojiPerWidth*emojiPerHeight; i++ {
		// 0-179 emoji codepoints (from emoji char array) each 4bytes + null char
		emoji[i].index = rl.GetRandomValue(0, 179) * 5

		// Generate a random color for this emoji
		emoji[i].color = rl.Fade(rl.ColorFromHSV(float32((start*(i+1))%360), 0.6, 0.85), 0.8)

		// Set a random message for this emoji
		emoji[i].message = rl.GetRandomValue(0, int32(len(messages)-1))
	}
}

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
				lastk, k = k-1, lastk
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

		textOffsetX += glyphWidth
	}
}
