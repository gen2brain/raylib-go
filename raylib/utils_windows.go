//go:build windows
// +build windows

package rl

/*
#include "raylib.h"
*/
import "C"

import (
	"fmt"
	"os"
)

// SetTraceLog - Enable trace log message types (bit flags based)
func SetTraceLog(typeFlags int) {
	logTypeFlags = typeFlags

	ctypeFlags := (C.int)(typeFlags)
	C.SetTraceLogLevel(ctypeFlags)
}

// TraceLog - Show trace log messages (INFO, WARNING, ERROR, DEBUG)
func TraceLog(msgType int, text string, v ...interface{}) {
	switch msgType {
	case LogInfo:
		if logTypeFlags&LogInfo == 0 {
			fmt.Printf("INFO: "+text+"\n", v...)
		}
	case LogWarning:
		if logTypeFlags&LogWarning == 0 {
			fmt.Printf("WARNING: "+text+"\n", v...)
		}
	case LogError:
		if logTypeFlags&LogError == 0 {
			fmt.Printf("ERROR: "+text+"\n", v...)
		}
	case LogDebug:
		if logTypeFlags&LogDebug == 0 {
			fmt.Printf("DEBUG: "+text+"\n", v...)
		}
	}
}

// HomeDir - Returns user home directory
func HomeDir() string {
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}
