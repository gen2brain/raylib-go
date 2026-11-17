package pkgtree

import (
	"path"
	"slices"
	"strings"
)

type File struct {
	Path string
	Mode uint32
	Data []byte
	Link string
	Dir  bool
}

func WithDirs(files []File) []File {
	seen := map[string]bool{}
	var out []File
	for _, f := range files {
		seen[f.Path] = true
	}
	for _, f := range files {
		for d := path.Dir(f.Path); d != "/" && d != "."; d = path.Dir(d) {
			if seen[d] {
				continue
			}
			seen[d] = true
			out = append(out, File{Path: d, Mode: 0o755, Dir: true})
		}
	}
	out = append(out, files...)
	slices.SortStableFunc(out, func(a, b File) int { return strings.Compare(a.Path, b.Path) })
	return out
}
