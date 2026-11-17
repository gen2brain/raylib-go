package pkgtree

import "testing"

func TestWithDirs(t *testing.T) {
	out := WithDirs([]File{{Path: "/usr/share/a/b.txt"}, {Path: "/usr/bin/x"}})
	var paths []string
	for _, f := range out {
		paths = append(paths, f.Path)
	}
	want := []string{"/usr", "/usr/bin", "/usr/bin/x", "/usr/share", "/usr/share/a", "/usr/share/a/b.txt"}
	if len(paths) != len(want) {
		t.Fatalf("paths = %v", paths)
	}
	for i := range want {
		if paths[i] != want[i] {
			t.Fatalf("paths = %v", paths)
		}
	}
	if !out[0].Dir || out[2].Dir {
		t.Error("dir flags wrong")
	}
}
