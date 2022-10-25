//go:build android
// +build android

package rl

// HomeDir - Returns user home directory
// NOTE: On Android this returns internal data path and must be called after InitWindow
func HomeDir() string {
	return getInternalStoragePath()
}
