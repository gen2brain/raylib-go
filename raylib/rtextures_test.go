package rl_test

import (
	"os"
	"testing"

	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestExportImage(t *testing.T) {
	var tests = []struct {
		name     string
		image    rl.Image
		fileName string
		want     bool
	}{
		{
			name:     "ValidImageExport",
			image:    *rl.GenImageColor(100, 100, color.RGBA{255, 0, 0, 255}),
			fileName: "test_image.png",
			want:     true,
		},
		{
			name:     "InvalidImageExport",
			image:    rl.Image{},
			fileName: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success := rl.ExportImage(tt.image, tt.fileName)
			defer os.Remove(tt.fileName)

			if success != tt.want {
				t.Errorf("ExportImage() result = %v, wantResult %v", success, tt.want)
				return
			}

			if _, err := os.Stat(tt.fileName); os.IsNotExist(err) && tt.want {
				t.Errorf("ExportImage() failed to create file %s", tt.fileName)
			}
		})
	}
}
