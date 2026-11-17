package icon

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadIconSquares(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 40, 20))
	for x := range 40 {
		for y := range 20 {
			src.Set(x, y, color.NRGBA{255, 0, 0, 255})
		}
	}
	path := filepath.Join(t.TempDir(), "icon.png")
	var buf bytes.Buffer
	if err := png.Encode(&buf, src); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	img, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if b := img.Bounds(); b.Dx() != 40 || b.Dy() != 40 {
		t.Fatalf("bounds = %v, want 40x40", b)
	}
	if _, _, _, a := img.At(20, 5).RGBA(); a != 0 {
		t.Errorf("padding alpha = %d, want 0", a)
	}
	if r, _, _, _ := img.At(20, 20).RGBA(); r != 0xffff {
		t.Errorf("content red = %d, want 0xffff", r)
	}
}

func TestEncodeICNS(t *testing.T) {
	data, err := ICNS(image.NewNRGBA(image.Rect(0, 0, 64, 64)))
	if err != nil {
		t.Fatal(err)
	}
	if string(data[:4]) != "icns" || int(binary.BigEndian.Uint32(data[4:])) != len(data) {
		t.Fatalf("bad header %q %d", data[:4], binary.BigEndian.Uint32(data[4:]))
	}

	sizes := map[string]int{}
	for off := 8; off < len(data); {
		typ := string(data[off : off+4])
		n := int(binary.BigEndian.Uint32(data[off+4:]))
		cfg, err := png.DecodeConfig(bytes.NewReader(data[off+8 : off+n]))
		if err != nil {
			t.Fatalf("%s: %v", typ, err)
		}
		sizes[typ] = cfg.Width
		off += n
	}
	want := map[string]int{"icp4": 16, "icp5": 32, "icp6": 64, "ic07": 128, "ic08": 256, "ic09": 512, "ic10": 1024, "ic11": 32, "ic12": 64, "ic13": 256, "ic14": 512}
	for typ, w := range want {
		if sizes[typ] != w {
			t.Errorf("%s width = %d, want %d", typ, sizes[typ], w)
		}
	}
}
