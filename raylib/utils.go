// +build !android

package raylib

import (
	"fmt"
	"os"
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
		fmt.Printf("INFO: "+text+"\n", v...)
	case LogError:
		fmt.Printf("ERROR: "+text+"\n", v...)
		os.Exit(1)
	case LogWarning:
		fmt.Printf("WARNING: "+text+"\n", v...)
	case LogDebug:
		if traceDebugMsgs {
			fmt.Printf("DEBUG: "+text+"\n", v...)
		}
	}
}
