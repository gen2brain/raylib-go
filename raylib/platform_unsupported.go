//go:build (software && (!cgo || nocgo || !(sdl || sdl3 || rgfw || memory || (drm && linux) || (wayland && linux) || (cocoa && darwin) || (win32 && windows)))) || (memory && (!cgo || nocgo))

package rl

var _ = softwareAndMemoryNeedCgoAndSoftwareNeedsTagSdlSdl3RgfwDrmWaylandCocoaWin32OrMemory
