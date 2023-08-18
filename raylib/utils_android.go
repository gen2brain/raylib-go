//go:build android
// +build android

package rl

/*
#include "stdlib.h"
#include "raylib.h"
void TraceLogWrapper(int logLevel, const char *text)
{
	TraceLog(logLevel, text);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// SetTraceLogLevel - Set the current threshold (minimum) log level
func SetTraceLogLevel(logLevel TraceLogLevel) {
	clogLevel := (C.int)(logLevel)
	C.SetTraceLogLevel(clogLevel)
}

// TraceLog - Show trace log messages (LOG_DEBUG, LOG_INFO, LOG_WARNING, LOG_ERROR...)
func TraceLog(logLevel TraceLogLevel, text string, v ...interface{}) {
	ctext := C.CString(fmt.Sprintf(text, v...))
	defer C.free(unsafe.Pointer(ctext))
	clogLevel := (C.int)(logLevel)
	C.TraceLogWrapper(clogLevel, ctext)
}

// HomeDir - Returns user home directory
// NOTE: On Android this returns internal data path and must be called after InitWindow
func HomeDir() string {
	return getInternalStoragePath()
}
