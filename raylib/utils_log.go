package rl

/*
#include "utils_log.h"
*/
import "C"

import "unsafe"

var internalTraceLogCallbackFun TraceLogCallbackFun = func(int, string) {}

// SetTraceLogCallback - set a call-back function for trace log
func SetTraceLogCallback(fn TraceLogCallbackFun) {
	internalTraceLogCallbackFun = fn
	C.setLogCallbackWrapper()
}

//export internalTraceLogCallbackGo
func internalTraceLogCallbackGo(logType C.int, cstr unsafe.Pointer, length C.int) {
	str := string(C.GoBytes(cstr, length))
	lt := int(logType)
	internalTraceLogCallbackFun(lt, str)
}
