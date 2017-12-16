// +build android,!js

package raylib

/*
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
	"os"
	"unsafe"
)

// TraceLog - Trace log messages showing (INFO, WARNING, ERROR, DEBUG)
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

// HomeDir - Returns user home directory
// NOTE: On Android this returns internal data path and must be called after InitWindow
func HomeDir() string {
	return C.GoString(C.get_internal_storage_path())
}
