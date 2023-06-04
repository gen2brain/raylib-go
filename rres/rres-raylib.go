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

// LoadTextFromResource - Load text data from rres resource chunk
func LoadTextFromResource(chunk ResourceChunk) string {
	cchunk := *(*C.rresResourceChunk)(unsafe.Pointer(&chunk))
	ret := C.LoadTextFromResource(cchunk)
	defer C.free(unsafe.Pointer(ret))
	v := C.GoString(ret)
	return v
}

// LoadImageFromResource - Load Image data from rres resource chunk
func LoadImageFromResource(chunk ResourceChunk) rl.Image {
	cchunk := *(*C.rresResourceChunk)(unsafe.Pointer(&chunk))
	ret := C.LoadImageFromResource(cchunk)
	v := *(*rl.Image)(unsafe.Pointer(&ret))
	return v
}

// LoadWaveFromResource - Load Wave data from rres resource chunk
func LoadWaveFromResource(chunk ResourceChunk) rl.Wave {
	cchunk := *(*C.rresResourceChunk)(unsafe.Pointer(&chunk))
	ret := C.LoadWaveFromResource(cchunk)
	v := *(*rl.Wave)(unsafe.Pointer(&ret))
	return v
}

// LoadFontFromResource - Load Font data from rres resource multiple chunks
func LoadFontFromResource(multi ResourceMulti) rl.Font {
	cmulti := *(*C.rresResourceMulti)(unsafe.Pointer(&multi))
	ret := C.LoadFontFromResource(cmulti)
	v := *(*rl.Font)(unsafe.Pointer(&ret))
	return v
}

// LoadMeshFromResource - Load Mesh data from rres resource multiple chunks
func LoadMeshFromResource(multi ResourceMulti) rl.Mesh {
	cmulti := *(*C.rresResourceMulti)(unsafe.Pointer(&multi))
	ret := C.LoadMeshFromResource(cmulti)
	v := *(*rl.Mesh)(unsafe.Pointer(&ret))
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

// SetBaseDirectory - Set base directory for externally linked data
//
// NOTE: When resource chunk contains an external link (FourCC: LINK, Type: RRES_DATA_LINK),
// a base directory is required to be prepended to link path
//
// If not provided, the application path is prepended to link by default
func SetBaseDirectory(baseDir string) {
	cbaseDir := C.CString(baseDir)
	defer C.free(unsafe.Pointer(cbaseDir))
	C.SetBaseDirectory(cbaseDir)
}
