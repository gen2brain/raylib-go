package main

import (
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestWindowsVersion(t *testing.T) {
	tests := []struct {
		version string
		build   int
		want    string
	}{
		{"1.0.0", 1, "1.0.0.1"},
		{"2.5", 7, "2.5.0.7"},
		{"3", 0, "3.0.0.0"},
		{"1.2.3.4", 9, "1.2.3.9"},
		{"1.2-beta", 3, "1.0.0.3"},
	}
	for _, tt := range tests {
		if got := windowsVersion(tt.version, tt.build); got != tt.want {
			t.Errorf("windowsVersion(%q, %d) = %q, want %q", tt.version, tt.build, got, tt.want)
		}
	}
}

func TestIdentifierPart(t *testing.T) {
	tests := map[string]string{
		"button":   "button",
		"My-App_2": "myapp2",
		"2048":     "app2048",
		"---":      "app",
	}
	for in, want := range tests {
		if got := identifierPart(in); got != want {
			t.Errorf("identifierPart(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestInfoPlistEscapesAndPermissions(t *testing.T) {
	c := &config{name: "A & B", exe: "ab", id: "com.example.ab", version: "1.0", build: 2, permissions: []string{"camera"}}
	plist := string(infoPlist(c, "11.0"))
	for _, want := range []string{"<string>A &amp; B</string>", "<key>NSCameraUsageDescription</key>", "<string>11.0</string>"} {
		if !bytes.Contains([]byte(plist), []byte(want)) {
			t.Errorf("plist missing %q", want)
		}
	}
	if bytes.Contains([]byte(plist), []byte("NSMicrophoneUsageDescription")) {
		t.Error("plist has microphone usage without the permission")
	}
	if ent := entitlementsPlist(&config{}); ent != nil {
		t.Error("entitlements without permissions")
	}
}

func TestMaxVersion(t *testing.T) {
	tests := [][3]string{
		{"11.0", "13.0", "13.0"},
		{"13.0", "11.0", "13.0"},
		{"10.15", "10.9", "10.15"},
		{"12", "12.1", "12.1"},
	}
	for _, tt := range tests {
		if got := maxVersion(tt[0], tt[1]); got != tt[2] {
			t.Errorf("maxVersion(%q, %q) = %q, want %q", tt[0], tt[1], got, tt[2])
		}
	}
}

func TestProfileEntitlementsPlist(t *testing.T) {
	profile := []byte("garbage<plist><dict><key>Entitlements</key>\n<dict>\n\t<key>application-identifier</key>\n\t<string>ABC123.*</string>\n\t<key>get-task-allow</key>\n\t<true/>\n</dict><key>Name</key><string>x</string></dict></plist>trailer")
	out, err := profileEntitlementsPlist(profile, "com.example.app")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"<string>ABC123.com.example.app</string>", "<key>get-task-allow</key>", "<plist version=\"1.0\">"} {
		if !bytes.Contains(out, []byte(want)) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if bytes.Contains(out, []byte("<key>Name</key>")) {
		t.Error("copied more than the Entitlements dict")
	}
	if _, err := profileEntitlementsPlist([]byte("nothing"), "a"); err == nil {
		t.Error("expected an error without Entitlements")
	}
}

func TestMSIXManifest(t *testing.T) {
	c := &config{name: "A & B", exe: "ab", id: "com.example.ab", vendor: "V", version: "1.2", build: 3, permissions: []string{"camera", "internet"}}
	manifest := msixManifest(c, "x64", `CN=V, O="V, Inc."`)
	var pkg struct {
		Identity struct {
			Publisher string `xml:"Publisher,attr"`
			Version   string `xml:"Version,attr"`
		}
	}
	if err := xml.Unmarshal(manifest, &pkg); err != nil {
		t.Fatal(err)
	}
	if pkg.Identity.Publisher != `CN=V, O="V, Inc."` || pkg.Identity.Version != "1.2.0.3" {
		t.Errorf("identity: %+v", pkg.Identity)
	}
	for _, want := range []string{"<DisplayName>A &amp; B</DisplayName>", `<Capability Name="internetClient"/>`, `<DeviceCapability Name="webcam"/>`, `Executable="ab.exe"`} {
		if !bytes.Contains(manifest, []byte(want)) {
			t.Errorf("manifest missing %q", want)
		}
	}
}

func TestIOSInfoPlistOrientation(t *testing.T) {
	tests := map[string][]string{
		"landscape": {"UIInterfaceOrientationLandscapeLeft", "UIInterfaceOrientationLandscapeRight"},
		"portrait":  {"UIInterfaceOrientationPortrait"},
		"":          {"UIInterfaceOrientationPortrait", "UIInterfaceOrientationLandscapeLeft", "UIInterfaceOrientationLandscapeRight"},
	}
	for orientation, want := range tests {
		c := &config{name: "A", exe: "a", id: "com.example.a", version: "1.0", build: 1, minOS: "15.0", orientation: orientation}
		var plist struct{}
		data := iosInfoPlist(c)
		if err := xml.Unmarshal(data, &plist); err != nil {
			t.Fatal(err)
		}
		_, after, _ := strings.Cut(string(data), "<key>UISupportedInterfaceOrientations</key>")
		list, _, _ := strings.Cut(after, "</array>")
		var got []string
		for _, line := range strings.Split(list, "\n") {
			if v, ok := strings.CutPrefix(strings.TrimSpace(line), "<string>"); ok {
				got = append(got, strings.TrimSuffix(v, "</string>"))
			}
		}
		if !slices.Equal(got, want) {
			t.Errorf("orientation %q: %v, want %v", orientation, got, want)
		}
	}
}

func TestWindowsManifest(t *testing.T) {
	var assembly struct{}
	if err := xml.Unmarshal(windowsManifest, &assembly); err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(windowsManifest, []byte("dpiAware")) {
		t.Error("the manifest sets DPI awareness, raylib sets it at run time")
	}
}

func TestLoadData(t *testing.T) {
	dir := t.TempDir()
	write := func(name, data string) {
		p := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("gopher.png", "png")
	write("images/a.png", "a")
	write("images/sub/b.png", "b")
	write("other/gopher.png", "dup")

	c := &config{data: []string{filepath.Join(dir, "gopher.png"), filepath.Join(dir, "images")}}
	if err := c.loadData(); err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, f := range c.files {
		names = append(names, f.name)
	}
	slices.Sort(names)
	want := []string{"gopher.png", "images/a.png", "images/sub/b.png"}
	if !slices.Equal(names, want) {
		t.Errorf("names = %v, want %v", names, want)
	}

	c = &config{data: []string{filepath.Join(dir, "gopher.png"), filepath.Join(dir, "other", "gopher.png")}}
	if err := c.loadData(); err == nil {
		t.Error("two entries with the same name: no error")
	}
}

func TestMakefileData(t *testing.T) {
	c := &config{exe: "demo"}
	if mk := string(makefile(c, "demo")); strings.Contains(mk, "share/$(EXE)") {
		t.Errorf("makefile without data installs share/$(EXE):\n%s", mk)
	}
	c.files = []dataFile{{"gopher.png", 0o644, []byte("png")}}
	mk := string(makefile(c, "demo"))
	for _, want := range []string{"cp -R share/$(EXE)/. $(DESTDIR)$(PREFIX)/share/$(EXE)/", "find share/$(EXE) -type f", "find share/$(EXE) -depth -type d"} {
		if !strings.Contains(mk, want) {
			t.Errorf("makefile lacks %q:\n%s", want, mk)
		}
	}
}

func TestDesktopEntryPath(t *testing.T) {
	c := &config{name: "Demo", exe: "demo", category: "Game"}
	if entry := string(installedDesktopEntry(c, "demo")); strings.Contains(entry, "Path=") || !strings.Contains(entry, "Exec=/usr/bin/demo\n") {
		t.Errorf("entry without data:\n%s", entry)
	}
	c.files = []dataFile{{"gopher.png", 0o644, []byte("png")}}
	if entry := string(installedDesktopEntry(c, "demo")); !strings.Contains(entry, "Path=/usr/share/demo\n") {
		t.Errorf("entry with data lacks the working directory:\n%s", entry)
	}
	if mk := string(makefile(c, "demo")); !strings.Contains(mk, "s|@DATADIR@|$(PREFIX)/share/$(EXE)|") {
		t.Errorf("makefile does not substitute @DATADIR@:\n%s", mk)
	}
}

func TestEntitlementsLibraryValidation(t *testing.T) {
	const key = "com.apple.security.cs.disable-library-validation"
	if ent := entitlementsPlist(&config{cgo: true, sign: "dev.p12"}); bytes.Contains(ent, []byte(key)) {
		t.Error("cgo build disables library validation")
	}
	if ent := entitlementsPlist(&config{}); ent != nil {
		t.Error("ad-hoc purego build has entitlements")
	}
	if ent := entitlementsPlist(&config{sign: "dev.p12"}); !bytes.Contains(ent, []byte(key)) {
		t.Error("signed purego build keeps library validation, libffi from the user cache cannot load")
	}
}
