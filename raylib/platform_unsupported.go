//go:build (software && (!cgo || nocgo || !(sdl || sdl3 || rgfw || memory || (drm && linux) || (win32 && windows)))) || (memory && (!cgo || nocgo))

package rl

var _ = softwareAndMemoryNeedCgoAndSoftwareNeedsTagSdlSdl3RgfwDrmWin32OrMemory
