package main

import (
	//"bytes"

	"github.com/gen2brain/raylib-go/raylib"
)

const numTextures = 4

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	raylib.InitWindow(screenWidth, screenHeight, "raylib [core] example - resources loading")

	raylib.InitAudioDevice()

	// OpenAsset() will also work on Android (reads files from assets/)
	reader, err := raylib.OpenAsset("data.rres")
	if err != nil {
		raylib.TraceLog(raylib.LogWarning, "[%s] rRES raylib resource file could not be opened: %v", "data.rres", err)
	}

	defer reader.Close()

	// bindata
	//b := MustAsset("data.rres")
	//reader := bytes.NewReader(b)

	res := raylib.LoadResource(reader, 0, nil)
	wav := raylib.LoadWaveEx(res.Data, int32(res.Param1), int32(res.Param2), int32(res.Param3), int32(res.Param4))
	snd := raylib.LoadSoundFromWave(wav)
	raylib.UnloadWave(wav)

	textures := make([]raylib.Texture2D, numTextures)
	for i := 0; i < numTextures; i++ {
		r := raylib.LoadResource(reader, i+1, nil)
		image := raylib.LoadImagePro(r.Data, int32(r.Param1), int32(r.Param2), raylib.TextureFormat(r.Param3))
		textures[i] = raylib.LoadTextureFromImage(image)
		raylib.UnloadImage(image)
	}

	currentTexture := 0

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if raylib.IsKeyPressed(raylib.KeySpace) {
			raylib.PlaySound(snd)
		}

		if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
			currentTexture = (currentTexture + 1) % numTextures // Cycle between the textures
		}

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawTexture(textures[currentTexture], screenWidth/2-textures[currentTexture].Width/2, screenHeight/2-textures[currentTexture].Height/2, raylib.RayWhite)

		raylib.DrawText("MOUSE LEFT BUTTON to CYCLE TEXTURES", 40, 410, 10, raylib.Gray)
		raylib.DrawText("SPACE to PLAY SOUND", 40, 430, 10, raylib.Gray)

		switch currentTexture {
		case 0:
			raylib.DrawText("GIF", 272, 70, 20, raylib.Gray)
			break
		case 1:
			raylib.DrawText("JPEG", 272, 70, 20, raylib.Gray)
			break
		case 2:
			raylib.DrawText("PNG", 272, 70, 20, raylib.Gray)
			break
		case 3:
			raylib.DrawText("TGA", 272, 70, 20, raylib.Gray)
			break
		default:
			break
		}

		raylib.EndDrawing()
	}

	raylib.UnloadSound(snd)

	for _, t := range textures {
		raylib.UnloadTexture(t)
	}

	raylib.CloseAudioDevice()

	raylib.CloseWindow()
}
