package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

const numTextures = 24

// Formats
const (
	PngR8g8b8a8 = iota
	PvrGrayscale
	PvrGrayAlpha
	PvrR5g6b5
	PvrR5g5b5a1
	PvrR4g4b4a4
	DdsR5g6b5
	DdsR5g5b5a1
	DdsR4g4b4a4
	DdsR8g8b8a8
	DdsDxt1Rgb
	DdsDxt1Rgba
	DdsDxt3Rgba
	DdsDxt5Rgba
	PkmEtc1Rgb
	PkmEtc2Rgb
	PkmEtc2EacRgba
	KtxEtc1Rgb
	KtxEtc2Rgb
	KtxEtc2EacRgba
	Astc4x4Ldr
	Astc8x8Ldr
	PvrPvrtRgb
	PvrPvrtRgba
)

var formatText = []string{
	"PNG_R8G8B8A8",
	"PVR_GRAYSCALE",
	"PVR_GRAY_ALPHA",
	"PVR_R5G6B5",
	"PVR_R5G5B5A1",
	"PVR_R4G4B4A4",
	"DDS_R5G6B5",
	"DDS_R5G5B5A1",
	"DDS_R4G4B4A4",
	"DDS_R8G8B8A8",
	"DDS_DXT1_RGB",
	"DDS_DXT1_RGBA",
	"DDS_DXT3_RGBA",
	"DDS_DXT5_RGBA",
	"PKM_ETC1_RGB",
	"PKM_ETC2_RGB",
	"PKM_ETC2_EAC_RGBA",
	"KTX_ETC1_RGB",
	"KTX_ETC2_RGB",
	"KTX_ETC2_EAC_RGBA",
	"ASTC_4x4_LDR",
	"ASTC_8x8_LDR",
	"PVR_PVRT_RGB",
	"PVR_PVRT_RGBA",
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [textures] example - texture formats loading")

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)

	sonic := make([]raylib.Texture2D, numTextures)

	sonic[PngR8g8b8a8] = raylib.LoadTexture("texture_formats/sonic.png")

	// Load UNCOMPRESSED PVR texture data
	sonic[PvrGrayscale] = raylib.LoadTexture("texture_formats/sonic_GRAYSCALE.pvr")
	sonic[PvrGrayAlpha] = raylib.LoadTexture("texture_formats/sonic_L8A8.pvr")
	sonic[PvrR5g6b5] = raylib.LoadTexture("texture_formats/sonic_R5G6B5.pvr")
	sonic[PvrR5g5b5a1] = raylib.LoadTexture("texture_formats/sonic_R5G5B5A1.pvr")
	sonic[PvrR4g4b4a4] = raylib.LoadTexture("texture_formats/sonic_R4G4B4A4.pvr")

	// Load UNCOMPRESSED DDS texture data
	sonic[DdsR5g6b5] = raylib.LoadTexture("texture_formats/sonic_R5G6B5.dds")
	sonic[DdsR5g5b5a1] = raylib.LoadTexture("texture_formats/sonic_A1R5G5B5.dds")
	sonic[DdsR4g4b4a4] = raylib.LoadTexture("texture_formats/sonic_A4R4G4B4.dds")
	sonic[DdsR8g8b8a8] = raylib.LoadTexture("texture_formats/sonic_A8R8G8B8.dds")

	// Load COMPRESSED DXT DDS texture data (if supported)
	sonic[DdsDxt1Rgb] = raylib.LoadTexture("texture_formats/sonic_DXT1_RGB.dds")
	sonic[DdsDxt1Rgba] = raylib.LoadTexture("texture_formats/sonic_DXT1_RGBA.dds")
	sonic[DdsDxt3Rgba] = raylib.LoadTexture("texture_formats/sonic_DXT3_RGBA.dds")
	sonic[DdsDxt5Rgba] = raylib.LoadTexture("texture_formats/sonic_DXT5_RGBA.dds")

	// Load COMPRESSED ETC texture data (if supported)
	sonic[PkmEtc1Rgb] = raylib.LoadTexture("texture_formats/sonic_ETC1_RGB.pkm")
	sonic[PkmEtc2Rgb] = raylib.LoadTexture("texture_formats/sonic_ETC2_RGB.pkm")
	sonic[PkmEtc2EacRgba] = raylib.LoadTexture("texture_formats/sonic_ETC2_EAC_RGBA.pkm")

	sonic[KtxEtc1Rgb] = raylib.LoadTexture("texture_formats/sonic_ETC1_RGB.ktx")
	sonic[KtxEtc2Rgb] = raylib.LoadTexture("texture_formats/sonic_ETC2_RGB.ktx")
	sonic[KtxEtc2EacRgba] = raylib.LoadTexture("texture_formats/sonic_ETC2_EAC_RGBA.ktx")

	// Load COMPRESSED ASTC texture data (if supported)
	sonic[Astc4x4Ldr] = raylib.LoadTexture("texture_formats/sonic_ASTC_4x4_ldr.astc")
	sonic[Astc8x8Ldr] = raylib.LoadTexture("texture_formats/sonic_ASTC_8x8_ldr.astc")

	// Load COMPRESSED PVR texture data (if supported)
	sonic[PvrPvrtRgb] = raylib.LoadTexture("texture_formats/sonic_PVRT_RGB.pvr")
	sonic[PvrPvrtRgba] = raylib.LoadTexture("texture_formats/sonic_PVRT_RGBA.pvr")

	selectedFormat := PngR8g8b8a8

	selectRecs := make([]raylib.Rectangle, numTextures)

	for i := 0; i < numTextures; i++ {
		if i < numTextures/2 {
			selectRecs[i] = raylib.NewRectangle(40, int32(30+32*i), 150, 30)
		} else {
			selectRecs[i] = raylib.NewRectangle(40+152, int32(30+32*(i-numTextures/2)), 150, 30)
		}
	}

	// Texture sizes in KB
	var textureSizes = [numTextures]int{
		512 * 512 * 32 / 8 / 1024, //PNG_R8G8B8A8 (32 bpp)
		512 * 512 * 8 / 8 / 1024,  //PVR_GRAYSCALE (8 bpp)
		512 * 512 * 16 / 8 / 1024, //PVR_GRAY_ALPHA (16 bpp)
		512 * 512 * 16 / 8 / 1024, //PVR_R5G6B5 (16 bpp)
		512 * 512 * 16 / 8 / 1024, //PVR_R5G5B5A1 (16 bpp)
		512 * 512 * 16 / 8 / 1024, //PVR_R4G4B4A4 (16 bpp)
		512 * 512 * 16 / 8 / 1024, //DDS_R5G6B5 (16 bpp)
		512 * 512 * 16 / 8 / 1024, //DDS_R5G5B5A1 (16 bpp)
		512 * 512 * 16 / 8 / 1024, //DDS_R4G4B4A4 (16 bpp)
		512 * 512 * 32 / 8 / 1024, //DDS_R8G8B8A8 (32 bpp)
		512 * 512 * 4 / 8 / 1024,  //DDS_DXT1_RGB (4 bpp) -Compressed-
		512 * 512 * 4 / 8 / 1024,  //DDS_DXT1_RGBA (4 bpp) -Compressed-
		512 * 512 * 8 / 8 / 1024,  //DDS_DXT3_RGBA (8 bpp) -Compressed-
		512 * 512 * 8 / 8 / 1024,  //DDS_DXT5_RGBA (8 bpp) -Compressed-
		512 * 512 * 4 / 8 / 1024,  //PKM_ETC1_RGB (4 bpp) -Compressed-
		512 * 512 * 4 / 8 / 1024,  //PKM_ETC2_RGB (4 bpp) -Compressed-
		512 * 512 * 8 / 8 / 1024,  //PKM_ETC2_EAC_RGBA (8 bpp) -Compressed-
		512 * 512 * 4 / 8 / 1024,  //KTX_ETC1_RGB (4 bpp) -Compressed-
		512 * 512 * 4 / 8 / 1024,  //KTX_ETC2_RGB (4 bpp) -Compressed-
		512 * 512 * 8 / 8 / 1024,  //KTX_ETC2_EAC_RGBA (8 bpp) -Compressed-
		512 * 512 * 8 / 8 / 1024,  //ASTC_4x4_LDR (8 bpp) -Compressed-
		512 * 512 * 2 / 8 / 1024,  //ASTC_8x8_LDR (2 bpp) -Compressed-
		512 * 512 * 4 / 8 / 1024,  //PVR_PVRT_RGB (4 bpp) -Compressed-
		512 * 512 * 4 / 8 / 1024,  //PVR_PVRT_RGBA (4 bpp) -Compressed-
	}

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		// Update

		if raylib.IsKeyPressed(raylib.KeyDown) {
			selectedFormat++
			if selectedFormat >= numTextures {
				selectedFormat = 0
			}
		} else if raylib.IsKeyPressed(raylib.KeyUp) {
			selectedFormat--
			if selectedFormat < 0 {
				selectedFormat = numTextures - 1
			}
		} else if raylib.IsKeyPressed(raylib.KeyRight) {
			if selectedFormat < numTextures/2 {
				selectedFormat += numTextures / 2
			}
		} else if raylib.IsKeyPressed(raylib.KeyLeft) {
			if selectedFormat >= numTextures/2 {
				selectedFormat -= numTextures / 2
			}
		}

		// Draw

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		// Draw rectangles
		for i := 0; i < numTextures; i++ {
			if i == selectedFormat {
				raylib.DrawRectangleRec(selectRecs[i], raylib.SkyBlue)
				raylib.DrawRectangleLines(selectRecs[i].X, selectRecs[i].Y, selectRecs[i].Width, selectRecs[i].Height, raylib.Blue)
				raylib.DrawText(formatText[i], selectRecs[i].X+selectRecs[i].Width/2-raylib.MeasureText(formatText[i], 10)/2, selectRecs[i].Y+11, 10, raylib.DarkBlue)
			} else {
				raylib.DrawRectangleRec(selectRecs[i], raylib.LightGray)
				raylib.DrawRectangleLines(selectRecs[i].X, selectRecs[i].Y, selectRecs[i].Width, selectRecs[i].Height, raylib.Gray)
				raylib.DrawText(formatText[i], selectRecs[i].X+selectRecs[i].Width/2-raylib.MeasureText(formatText[i], 10)/2, selectRecs[i].Y+11, 10, raylib.DarkGray)
			}
		}

		// Draw selected texture
		if sonic[selectedFormat].ID != 0 {
			raylib.DrawTexture(sonic[selectedFormat], 350, -10, raylib.White)
		} else {
			raylib.DrawRectangleLines(488, 165, 200, 110, raylib.DarkGray)
			raylib.DrawText("FORMAT", 550, 180, 20, raylib.Maroon)
			raylib.DrawText("NOT SUPPORTED", 500, 210, 20, raylib.Maroon)
			raylib.DrawText("ON YOUR GPU", 520, 240, 20, raylib.Maroon)
		}

		raylib.DrawText("Select texture format (use cursor keys):", 40, 10, 10, raylib.DarkGray)
		raylib.DrawText("Required GPU memory size (VRAM):", 40, 427, 10, raylib.DarkGray)
		raylib.DrawText(fmt.Sprintf("%4.0d KB", textureSizes[selectedFormat]), 240, 420, 20, raylib.DarkBlue)

		raylib.EndDrawing()
	}

	for i := 0; i < numTextures; i++ {
		raylib.UnloadTexture(sonic[i])
	}

	raylib.CloseWindow()
}
