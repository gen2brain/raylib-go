//go:build (software && (!cgo || nocgo || ((glfw || glfw_x11 || glfw_wayland) && !(sdl || sdl3 || rgfw || memory || (drm && linux) || (android && !termux) || ios)))) || (memory && (!cgo || nocgo))

package rl

var _ = softwareAndMemoryNeedCgoAndSoftwareIsNotSupportedByGlfw
