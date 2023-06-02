package rres

// #include <raylib.h>
// #define RRES_RAYLIB_IMPLEMENTATION
// #define RRES_SUPPORT_COMPRESSION_LZ4
// #define RRES_SUPPORT_ENCRYPTION_AES
// #define RRES_SUPPORT_ENCRYPTION_XCHACHA20
// #include <rres-raylib.h>
// #include <rres.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadImageFromResource(chunk ResourceChunk) rl.Image {
	cchunk := *(*C.rresResourceChunk)(unsafe.Pointer(&chunk))
	ret := C.LoadImageFromResource(cchunk)
	v := *(*rl.Image)(unsafe.Pointer(&ret))
	return v
}
