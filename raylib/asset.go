//go:build !android
// +build !android

package rl

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

// Open implements fs.FS interface - opens the named file for reading
func (a *Asset) Open(name string) (fs.File, error) {
	if a.fsys != nil {
		return a.fsys.Open(a.fsPath(name))
	}

	f, err := os.Open(a.osPath(name))
	if err != nil {
		return nil, err
	}

	return &desktopAsset{f}, nil
}

// ReadFile implements fs.ReadFileFS interface - reads the entire file
func (a *Asset) ReadFile(name string) ([]byte, error) {
	if a.fsys != nil {
		return fs.ReadFile(a.fsys, a.fsPath(name))
	}

	return os.ReadFile(a.osPath(name))
}

// ReadDir implements fs.ReadDirFS interface - reads the directory
func (a *Asset) ReadDir(name string) ([]fs.DirEntry, error) {
	if a.fsys != nil {
		return fs.ReadDir(a.fsys, a.fsPath(name))
	}

	return os.ReadDir(a.osPath(name))
}

func (a *Asset) fsPath(name string) string {
	if a.root == "" {
		return name
	}

	return path.Join(a.root, name)
}

func (a *Asset) osPath(name string) string {
	if a.root == "" {
		return name
	}

	return filepath.Join(a.root, name)
}

// desktopAsset wraps os.File to implement AssetFile interface
type desktopAsset struct {
	*os.File
}
