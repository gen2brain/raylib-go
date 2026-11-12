package raygui

import (
	"encoding/binary"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// GetIcons - Get raygui icons data pointer
func GetIcons() []uint32 {
	return guiIconsPtr
}

// LoadIcons - Load raygui icons file (.rgi) into internal icons data
func LoadIcons(fileName string, loadIconsName bool) []string {
	fileData, err := os.ReadFile(fileName)
	if err != nil || len(fileData) == 0 {
		return nil
	}

	return LoadIconsFromMemory(fileData, loadIconsName)
}

// LoadIconsFromMemory - Load icons from memory (Binary files only)
func LoadIconsFromMemory(fileData []byte, loadIconsName bool) []string {
	r := &rgsReader{data: fileData}

	var guiIconsName []string

	signature := r.bytes(4)
	r.int16()
	r.int16()
	iconCount := int(r.int16())
	iconSize := int(r.int16())

	if string(signature) == "rGI " {
		if loadIconsName {
			guiIconsName = make([]string, RAYGUI_ICON_MAX_ICONS)
			for i := range iconCount {
				name := cstr(string(r.bytes(RAYGUI_ICON_MAX_NAME_LENGTH)))
				if i < RAYGUI_ICON_MAX_ICONS {
					guiIconsName[i] = name
				}
			}
		} else {
			r.pos += iconCount * RAYGUI_ICON_MAX_NAME_LENGTH
		}

		iconDataSize := iconCount * (iconSize * iconSize / 32)
		data := r.bytes(iconDataSize * 4)

		guiIconsPtr = make([]uint32, iconDataSize)
		for i := range guiIconsPtr {
			guiIconsPtr[i] = binary.LittleEndian.Uint32(data[i*4:])
		}
	}

	return guiIconsName
}

// DrawIcon - Draw icon using pixel size at specified position
func DrawIcon(iconId IconID, posX, posY, pixelSize int32, color rl.Color) {
	for i, y := 0, int32(0); i < RAYGUI_ICON_SIZE*RAYGUI_ICON_SIZE/32; i++ {
		idx := int(iconId)*RAYGUI_ICON_DATA_ELEMENTS + i

		var data uint32
		if idx < len(guiIconsPtr) {
			data = guiIconsPtr[idx]
		}

		for k := range int32(32) {
			if data&(1<<uint(k)) != 0 {
				guiDrawRectangle(rl.Rectangle{X: float32(posX + (k%RAYGUI_ICON_SIZE)*pixelSize), Y: float32(posY + y*pixelSize), Width: float32(pixelSize), Height: float32(pixelSize)}, 0, blank, color)
			}

			if k == 15 || k == 31 {
				y++
			}
		}
	}
}

// SetIconScale - Set default icon drawing size
func SetIconScale(scale int32) {
	if scale >= 1 {
		guiIconScale = scale
	}
}

// GetTextWidth - Get text width considering gui style and icon size (if required)
func GetTextWidth(text string) int32 {
	var maxWidth int32

	p := 0
	for at(text, p) != 0 && p < maxLineBufferSize {
		lineWidth := getLineWidth(tail(text, p))
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}

		for at(text, p) != 0 && at(text, p) != '\n' && p < maxLineBufferSize {
			p++
		}

		if at(text, p) == '\n' {
			p++
		}
	}

	return maxWidth
}
