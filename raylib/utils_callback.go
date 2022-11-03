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
	internalLoadFileTextCallback func(name string) (string, error)
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

func SetLoadFileTextCallback(fn func(name string) (string, error)) {
	internalLoadFileTextCallback = fn
	C.setLoadFileTextCallbackWrapper()
}

//export loadFileTextCallbackGo
func loadFileTextCallbackGo(cstr unsafe.Pointer, slen C.int, outstrlen *C.int, ref **C.char) {
	str := string(C.GoBytes(cstr, slen))
	data, err := internalLoadFileTextCallback(str)
	if err != nil {
		//TODO: handle error
		return
	}
	*outstrlen = C.int(len(data))
	*ref = (*C.char)(C.CString(data))
}
