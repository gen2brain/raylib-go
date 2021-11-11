//go:build android
// +build android

package rl

/*
#include "raylib.h"
#include <stdlib.h>

void log_info(const char *msg);
void log_warn(const char *msg);
void log_error(const char *msg);
void log_debug(const char *msg);

extern char* get_internal_storage_path();
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// SetTraceLog - Enable trace log message types (bit flags based)
func SetTraceLog(typeFlags int) {
	logTypeFlags = typeFlags

	ctypeFlags := (C.int)(typeFlags)
	C.SetTraceLogLevel(ctypeFlags)
}

// TraceLog - Trace log messages showing (INFO, WARNING, ERROR, DEBUG)
func TraceLog(msgType int, text string, v ...interface{}) {
	switch msgType {
	case LogInfo:
		if logTypeFlags&LogInfo == 0 {
			msg := C.CString(fmt.Sprintf("INFO: "+text, v...))
			defer C.free(unsafe.Pointer(msg))
			C.log_info(msg)
		}
	case LogWarning:
		if logTypeFlags&LogWarning == 0 {
			msg := C.CString(fmt.Sprintf("WARNING: "+text, v...))
			defer C.free(unsafe.Pointer(msg))
			C.log_warn(msg)
		}
	case LogError:
		if logTypeFlags&LogError == 0 {
			msg := C.CString(fmt.Sprintf("ERROR: "+text, v...))
			defer C.free(unsafe.Pointer(msg))
			C.log_error(msg)
		}
	case LogDebug:
		if logTypeFlags&LogDebug == 0 {
			msg := C.CString(fmt.Sprintf("DEBUG: "+text, v...))
			defer C.free(unsafe.Pointer(msg))
			C.log_debug(msg)
		}
	}
}

// HomeDir - Returns user home directory
// NOTE: On Android this returns internal data path and must be called after InitWindow
func HomeDir() string {
	return C.GoString(C.get_internal_storage_path())
}
