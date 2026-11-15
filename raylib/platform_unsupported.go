//go:build (software && (!cgo || nocgo || !(sdl || sdl3 || rgfw || memory || (drm && linux) || (wayland && linux) || (cocoa && darwin) || (win32 && windows) || ios))) || (memory && (!cgo || nocgo))

package rl

var _ = softwareAndMemoryNeedCgoAndSoftwareNeedsTagSdlSdl3RgfwDrmWaylandCocoaWin32OrMemoryOrIos
