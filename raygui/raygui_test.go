package raygui

import "testing"

func TestScanStyleProperty(t *testing.T) {
	lines := []string{"p 00 00 0x878787ff\n", "p 01 02 878787ff\n", "p 03 04 0x1234abcz\n", "p 5 6\n", "p x\n", "p 07 18 0X00FF00ff\n", "p -1 2 0xff\n", "p 2 3 0x\n"}
	want := [][3]int64{{0, 0, 2273806335}, {1, 2, 2273806335}, {3, 4, 19090108}, {5, 6, 19090108}, {5, 6, 19090108}, {7, 18, 19090108}, {-1, 2, 255}, {2, 3, 255}}

	var cid, pid int64
	var val uint64
	for i, line := range lines {
		scanStyleProperty(line, &cid, &pid, &val)
		if got := [3]int64{cid, pid, int64(uint32(val))}; got != want[i] {
			t.Errorf("%q: got %v, want %v", line, got, want[i])
		}
	}
}

func TestScanStyleVersion(t *testing.T) {
	lines := []string{"v 600\n", "v600\n", "v 400x\n", "v abc\n", "v -1\n"}
	want := []uint64{600, 600, 400, 400, 4294967295}

	var ver uint64 = 7
	for i, line := range lines {
		scanStyleVersion(line, &ver)
		if ver != want[i] {
			t.Errorf("%q: got %d, want %d", line, ver, want[i])
		}
	}
}

func TestScanStyleFont(t *testing.T) {
	tests := []struct {
		line          string
		size          int64
		first, second string
	}{
		{"f 14 Kyrou7Wide.ttf charset.txt\n", 14, "Kyrou7Wide.ttf", "charset.txt"},
		{"f abc x.ttf\n", 0, "", ""},
		{"f 10-5 a b\n", 10, "-5", "a b"},
		{"f 16 font.ttf\r\n", 16, "font.ttf", ""},
		{"f 12 averyveryveryveryveryverylongfontname.ttf charmap with spaces.txt\n", 12, "averyveryveryveryveryverylongfo", "ntname.ttf charmap with spaces."},
	}

	for _, tt := range tests {
		size, first, second := scanStyleFont(tt.line)
		if size != tt.size || first != tt.first || second != tt.second {
			t.Errorf("%q: got %d %q %q, want %d %q %q", tt.line, size, first, second, tt.size, tt.first, tt.second)
		}
	}
}

func TestGetCodepointPreviousAtStart(t *testing.T) {
	var size int32
	if cp := getCodepointPrevious("\xa3x", 1, &size); cp != 0x3f || size != 1 {
		t.Errorf("got %#x size %d, want 0x3f size 1", cp, size)
	}
}
