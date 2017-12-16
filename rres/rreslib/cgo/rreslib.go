package main

import "C"

import (
	"unsafe"

	"github.com/gen2brain/raylib-go/rres/rreslib"
)

//export Compress
func Compress(compType C.int, data *C.uchar, uncompSize C.ulong, outCompSize *C.ulong) *C.uchar {
	d := C.GoBytes(unsafe.Pointer(data), C.int(uncompSize))

	c, err := rreslib.Compress(d, int(compType))
	if err != nil {
		return nil
	}

	*outCompSize = C.ulong(len(c))
	return (*C.uchar)(unsafe.Pointer(&c[0]))
}

//export Decompress
func Decompress(compType C.int, data *C.uchar, uncompSize C.ulong, compSize C.ulong) *C.uchar {
	d := C.GoBytes(unsafe.Pointer(data), C.int(compSize))

	c, err := rreslib.Decompress(d, int(compType))
	if err != nil {
		return nil
	}

	if len(c) != int(uncompSize) {
		return nil
	}

	return (*C.uchar)(unsafe.Pointer(&c[0]))
}

//export Encrypt
func Encrypt(cryptoType C.int, key *C.char, data *C.uchar, dataSize C.ulong, outDataSize *C.ulong) *C.uchar {
	k := []byte(C.GoString(key))
	d := C.GoBytes(unsafe.Pointer(data), C.int(dataSize))

	c, err := rreslib.Encrypt(k, d, int(cryptoType))
	if err != nil {
		return nil
	}

	*outDataSize = C.ulong(len(c))
	return (*C.uchar)(unsafe.Pointer(&c[0]))
}

//export Decrypt
func Decrypt(cryptoType C.int, key *C.char, data *C.uchar, dataSize C.ulong) *C.uchar {
	k := []byte(C.GoString(key))
	d := C.GoBytes(unsafe.Pointer(data), C.int(dataSize))

	c, err := rreslib.Decrypt(k, d, int(cryptoType))
	if err != nil {
		return nil
	}

	return (*C.uchar)(unsafe.Pointer(&c[0]))
}

// We need the main() so we can compile as C library
func main() {
}
