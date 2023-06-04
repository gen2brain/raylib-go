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

// LoadDataFromResource - Load raw data from rres resource chunk
//
// NOTE: Chunk data must be provided uncompressed/unencrypted
func LoadDataFromResource(chunk ResourceChunk) []byte {
	cchunk := *(*C.rresResourceChunk)(unsafe.Pointer(&chunk))
	var csize C.uint
	ret := C.LoadDataFromResource(cchunk, &csize)
	defer C.free(ret)
	v := C.GoBytes(ret, C.int(csize))
	return v
}

// UnpackResourceChunk - Unpack resource chunk data (decompres/decrypt data)
//
// NOTE: Function return 0 on success or other value on failure
func UnpackResourceChunk(chunk *ResourceChunk) ErrorType {
	cchunk := (*C.rresResourceChunk)(unsafe.Pointer(chunk))
	ret := C.UnpackResourceChunk(cchunk)
	v := ErrorType(ret)
	return v
}

func LoadImageFromResource(chunk ResourceChunk) rl.Image {
	cchunk := *(*C.rresResourceChunk)(unsafe.Pointer(&chunk))
	ret := C.LoadImageFromResource(cchunk)
	v := *(*rl.Image)(unsafe.Pointer(&ret))
	return v
}
