// +build android

package raylib

/*
#include <android/log.h>
#include <stdlib.h>

void logInfo(const char *msg) {
    __android_log_print(ANDROID_LOG_INFO, "raylib", msg);
}

void logWarn(const char *msg) {
    __android_log_print(ANDROID_LOG_WARN, "raylib", msg);
}

void logError(const char *msg) {
    __android_log_print(ANDROID_LOG_ERROR, "raylib", msg);
}

void logDebug(const char *msg) {
    __android_log_print(ANDROID_LOG_DEBUG, "raylib", msg);
}

void logInfo(const char *msg);
void logWarn(const char *msg);
void logError(const char *msg);
void logDebug(const char *msg);
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

func TraceLog(msgType int, text string, v ...interface{}) {
	switch msgType {
	case LogInfo:
		msg := C.CString(fmt.Sprintf("INFO: "+text, v...))
		defer C.free(unsafe.Pointer(msg))
		C.logInfo(msg)
	case LogError:
		msg := C.CString(fmt.Sprintf("ERROR: "+text, v...))
		defer C.free(unsafe.Pointer(msg))
		C.logError(msg)
		os.Exit(1)
	case LogWarning:
		msg := C.CString(fmt.Sprintf("WARNING: "+text, v...))
		defer C.free(unsafe.Pointer(msg))
		C.logWarn(msg)
	case LogDebug:
		if traceDebugMsgs {
			msg := C.CString(fmt.Sprintf("DEBUG: "+text, v...))
			defer C.free(unsafe.Pointer(msg))
			C.logDebug(msg)
		}
	}
}

func SetDebug(enabled bool) {
	traceDebugMsgs = enabled
}
