package raygui

/*
#include <stdlib.h>
// Code from https://blog.pytool.com/language/golang/cgo/go-cgo-string/

static char** make_str_array(int size) {
	return calloc(sizeof(char*), size);
}

static int len_str_array(char **arr) {
	int i = 0;
	while (arr[i] != NULL) i++;
	return i+1; // NULL does count
}

static void set_str_array(char **arr, int idx, char *s) {
	arr[idx] = s;
}

static void free_str_array(char **arr, int size) {
	int i;
	for (i = 0; i < size; i++) {
		free(arr[i]);
	}
	free(arr);
}
*/
import "C"

import (
	"unsafe"
)

// CStringArray represents an array of pointers to NULL terminated C strings,
// the array itself is terminated with a NULL
type CStringArray struct {
	Pointer unsafe.Pointer
	Length  int
}

// NewCStringArray returns an instance of CStringArray
func NewCStringArray() *CStringArray {
	return &CStringArray{}
}

// NewCStringArrayFromSlice makes an instance of CStringArray then copy the
// input slice to it.
func NewCStringArrayFromSlice(ss []string) *CStringArray {
	var arr CStringArray
	arr.Copy(ss)
	return &arr
}

func NewCStringArrayFromPointer(p unsafe.Pointer) *CStringArray {
	return &CStringArray{
		Length:  int(C.len_str_array((**C.char)(p))),
		Pointer: p,
	}
}

// ToSlice converts CStringArray to Go slice of strings
func (arr *CStringArray) ToSlice() []string {
	if arr.Length == 0 || arr.Pointer == nil {
		return []string{}
	}

	var ss []string
	var cs **C.char
	defer C.free(unsafe.Pointer(cs))
	p := uintptr(arr.Pointer)
	for {
		cs = (**C.char)(unsafe.Pointer(p))
		if *cs == nil { // skip NULL - the last element
			break
		}
		ss = append(ss, C.GoString(*cs))
		p += unsafe.Sizeof(p)
	}

	return ss
}

// Copy converts Go slice of strings to C underlying struct of CStringArray
func (arr *CStringArray) Copy(ss []string) {
	arr.Length = len(ss) + 1 // one more element for NULL at the end
	arr.Pointer = unsafe.Pointer(C.make_str_array(C.int(arr.Length)))

	for i, s := range ss {
		cs := C.CString(s) // will be free by Free() method
		C.set_str_array((**C.char)(arr.Pointer), C.int(i), cs)
	}
}

// Free frees C underlying struct of CStringArray
// MUST call this method after using CStringArray
// Exception: If you use NewCStringArrayFromPointer() to create CStringArray object
// and you use other way to free C underlying structure pointed by the pointer,
// then don't need to call Free()
func (arr *CStringArray) Free() {
	C.free_str_array((**C.char)(arr.Pointer), C.int(arr.Length))
}
