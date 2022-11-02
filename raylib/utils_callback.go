package rl

/*
#include "utils_callback.h"
*/
import "C"

import (
	"unsafe"
)

var (
	internalLoadFileDataCallback func(name string) ([]byte, error)
)

func SetLoadFileDataCallback(fn func(name string) ([]byte, error)) {
	internalLoadFileDataCallback = fn
	C.setLoadFileDataCallbackWrapper()
}

//export loadFileDataCallbackGo
func loadFileDataCallbackGo(cstr unsafe.Pointer, slen C.int, bytesRead *C.int, ref **C.uchar) {
	str := string(C.GoBytes(cstr, slen))
	data, err := internalLoadFileDataCallback(str)
	if err != nil {
		//TODO: handle error
		return
	}
	*bytesRead = C.int(len(data))
	*ref = (*C.uchar)(unsafe.Pointer(&data[0]))
}
