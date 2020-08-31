package rl

/*
#include "utils_log.h"
*/
import "C"

import "unsafe"

// TraceLogCallbackFun - function that will recive the trace log messages
type TraceLogCallbackFun func(int, string)

var internalTraceLogCallbackFun TraceLogCallbackFun = func(int, string) {}

// SetTraceLogCallback - set a call-back function for trace log
func SetTraceLogCallback(fn TraceLogCallbackFun) {
	internalTraceLogCallbackFun = fn
	C.setLogCallbackWrapper()
}

//export internalTraceLogCallbackGo
func internalTraceLogCallbackGo(logType C.int, cstr unsafe.Pointer, len C.int) {
	str := string(C.GoBytes(cstr, len))
	lt := int(logType)
	internalTraceLogCallbackFun(lt, str)
}
