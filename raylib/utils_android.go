// +build android

package raylib

/*
#include <android/log.h>
#include <stdlib.h>

void log_info(const char *msg) {
    __android_log_print(ANDROID_LOG_INFO, "raylib", msg);
}

void log_warn(const char *msg) {
    __android_log_print(ANDROID_LOG_WARN, "raylib", msg);
}

void log_error(const char *msg) {
    __android_log_print(ANDROID_LOG_ERROR, "raylib", msg);
}

void log_debug(const char *msg) {
    __android_log_print(ANDROID_LOG_DEBUG, "raylib", msg);
}

void log_info(const char *msg);
void log_warn(const char *msg);
void log_error(const char *msg);
void log_debug(const char *msg);
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

// Log message types
const (
	LogInfo = iota
	LogError
	LogWarning
	LogDebug
)

var traceDebugMsgs = false

// Set debug messages
func SetDebug(enabled bool) {
	traceDebugMsgs = enabled
}

// Trace log
func TraceLog(msgType int, text string, v ...interface{}) {
	switch msgType {
	case LogInfo:
		msg := C.CString(fmt.Sprintf("INFO: "+text, v...))
		defer C.free(unsafe.Pointer(msg))
		C.log_info(msg)
	case LogError:
		msg := C.CString(fmt.Sprintf("ERROR: "+text, v...))
		defer C.free(unsafe.Pointer(msg))
		C.log_error(msg)
		os.Exit(1)
	case LogWarning:
		msg := C.CString(fmt.Sprintf("WARNING: "+text, v...))
		defer C.free(unsafe.Pointer(msg))
		C.log_warn(msg)
	case LogDebug:
		if traceDebugMsgs {
			msg := C.CString(fmt.Sprintf("DEBUG: "+text, v...))
			defer C.free(unsafe.Pointer(msg))
			C.log_debug(msg)
		}
	}
}
