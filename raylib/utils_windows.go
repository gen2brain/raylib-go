// +build windows,!js

package raylib

import (
	"fmt"
	"os"
)

// TraceLog - Trace log messages showing (INFO, WARNING, ERROR, DEBUG)
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

// HomeDir - Returns user home directory
func HomeDir() string {
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}
