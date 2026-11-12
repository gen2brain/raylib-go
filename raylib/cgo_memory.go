//go:build memory && !android && !nocgo
// +build memory,!android,!nocgo

package rl

/*
#cgo CFLAGS: -DPLATFORM_MEMORY -DGRAPHICS_API_OPENGL_SOFTWARE
#cgo windows CFLAGS: -Iexternal
#cgo linux LDFLAGS: -lm -lpthread -ldl -lrt
#cgo freebsd LDFLAGS: -lm -lpthread
#cgo openbsd LDFLAGS: -lm -lpthread
#cgo windows LDFLAGS: -lwinmm -lole32
*/
import "C"
