package main

import (
	//"bytes"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/rres"
)

const numTextures = 4

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - resources loading")

	rl.InitAudioDevice()

	// OpenAsset() will also work on Android (reads files from assets/)
	reader, err := rl.OpenAsset("data.rres")
	if err != nil {
		rl.TraceLog(rl.LogWarning, "[%s] rRES raylib resource file could not be opened: %v", "data.rres", err)
	}

	defer reader.Close()

	// bindata
	//b := MustAsset("data.rres")
	//reader := bytes.NewReader(b)

	//res := rres.LoadResource(reader, 0, []byte("passwordpassword"))
	//wav := rl.LoadWaveEx(res.Data, int32(res.Param1), int32(res.Param2), int32(res.Param3), int32(res.Param4))
	//snd := rl.LoadSoundFromWave(wav)
	//rl.UnloadWave(wav)

	textures := make([]rl.Texture2D, numTextures)
	for i := 0; i < numTextures; i++ {
		r := rres.LoadResource(reader, i+1, []byte("passwordpassword"))
		image := rl.LoadImagePro(r.Data, int32(r.Param1), int32(r.Param2), rl.PixelFormat(r.Param3))
		textures[i] = rl.LoadTextureFromImage(image)
		rl.UnloadImage(image)
	}

	currentTexture := 0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeySpace) {
			//rl.PlaySound(snd)
		}

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			currentTexture = (currentTexture + 1) % numTextures // Cycle between the textures
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(textures[currentTexture], screenWidth/2-textures[currentTexture].Width/2, screenHeight/2-textures[currentTexture].Height/2, rl.RayWhite)

		rl.DrawText("MOUSE LEFT BUTTON to CYCLE TEXTURES", 40, 410, 10, rl.Gray)
		rl.DrawText("SPACE to PLAY SOUND", 40, 430, 10, rl.Gray)

		switch currentTexture {
		case 0:
			rl.DrawText("GIF", 272, 70, 20, rl.Gray)
			break
		case 1:
			rl.DrawText("JPEG", 272, 70, 20, rl.Gray)
			break
		case 2:
			rl.DrawText("PNG", 272, 70, 20, rl.Gray)
			break
		case 3:
			rl.DrawText("TGA", 272, 70, 20, rl.Gray)
			break
		default:
			break
		}

		rl.EndDrawing()
	}

	//rl.UnloadSound(snd)

	for _, t := range textures {
		rl.UnloadTexture(t)
	}

	rl.CloseAudioDevice()

	rl.CloseWindow()
}
