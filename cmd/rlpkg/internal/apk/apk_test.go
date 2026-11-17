package apk

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode/utf16"
)

func templateAPK(t *testing.T) string {
	_, file, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(file), "..", "..", "android", "template.apk")
	if _, err := os.Stat(path); err != nil {
		t.Skip(err)
	}
	return path
}

func templateBundle(t *testing.T) []*Entry {
	_, file, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(file), "..", "..", "android", "template.aab")
	if _, err := os.Stat(path); err != nil {
		t.Skip(err)
	}
	entries, err := ReadTemplate(path)
	if err != nil {
		t.Fatal(err)
	}
	return entries
}

func bundleContent(t *testing.T, entries []*Entry, name string) []byte {
	e := Find(entries, name)
	if e == nil {
		t.Fatalf("%s missing", name)
	}
	data, err := e.Content()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestProtoRoundTrip(t *testing.T) {
	entries := templateBundle(t)
	for _, name := range []string{"base/manifest/AndroidManifest.xml", "base/resources.pb"} {
		data := bundleContent(t, entries, name)
		fields, err := protoParse(data)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !bytes.Equal(protoEncode(fields), data) {
			t.Errorf("%s: re-encoded bytes differ", name)
		}
	}
}

func TestPatchBundleManifest(t *testing.T) {
	data := bundleContent(t, templateBundle(t), "base/manifest/AndroidManifest.xml")
	values := ManifestValues{
		Package: "org.example.demo", Label: "Demo App", Version: "2.0", Build: 9, Library: "demo",
		GLESVersion: 0x30000, Orientation: "sensorPortrait", Permissions: []string{"android.permission.CAMERA"},
	}
	out, err := PatchBundleManifest(data, values)
	if err != nil {
		t.Fatal(err)
	}

	attrs := map[string][]string{}
	var code, gles, orientation uint64
	_, err = visitBundleXML(out, func(element []protoField) ([]protoField, error) {
		name := protoGetString(element, xmlElementName)
		for _, f := range element {
			if f.num != xmlElementAttribute {
				continue
			}
			attr, err := protoParse(f.b)
			if err != nil {
				return nil, err
			}
			key := name + "@" + protoGetString(attr, xmlAttributeName)
			attrs[key] = append(attrs[key], protoGetString(attr, xmlAttributeValue))
			primitive := func(num int) uint64 {
				item, _ := protoParse(attr[protoFind(attr, xmlAttributeItem)].b)
				prim, _ := protoParse(item[protoFind(item, itemPrimitive)].b)
				return prim[protoFind(prim, num)].v
			}
			switch key {
			case "manifest@versionCode":
				code = primitive(primitiveInt)
			case "uses-feature@glEsVersion":
				gles = primitive(primitiveHex)
			case "activity@screenOrientation":
				orientation = primitive(primitiveInt)
			}
		}
		return element, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	for key, want := range map[string]string{
		"manifest@package": "org.example.demo", "manifest@versionName": "2.0", "manifest@versionCode": "9",
		"application@label": "Demo App",
	} {
		if got := attrs[key]; len(got) != 1 || got[0] != want {
			t.Errorf("%s = %q", key, got)
		}
	}
	if code != 9 {
		t.Errorf("versionCode item = %d", code)
	}
	if gles != 0x30000 {
		t.Errorf("glEsVersion item = %#x", gles)
	}
	if orientation != 7 {
		t.Errorf("screenOrientation item = %d", orientation)
	}
	if got := attrs["meta-data@value"]; !slices.Equal(got, []string{"demo"}) {
		t.Errorf("android.app.lib_name = %q", got)
	}
	if got := attrs["uses-permission@name"]; !slices.Equal(got, []string{"android.permission.CAMERA"}) {
		t.Errorf("permissions = %q", got)
	}
	if _, err := PatchBundleManifest(out, ManifestValues{Package: "a.b", Label: "x", Version: "1", Library: "x"}); err != nil {
		t.Errorf("repatch: %v", err)
	}
}

func TestRenameBundlePackage(t *testing.T) {
	out, err := RenameBundlePackage(bundleContent(t, templateBundle(t), "base/resources.pb"), "org.example.demo")
	if err != nil {
		t.Fatal(err)
	}
	fields, err := protoParse(out)
	if err != nil {
		t.Fatal(err)
	}
	pkg, err := protoParse(fields[protoFind(fields, tablePackage)].b)
	if err != nil {
		t.Fatal(err)
	}
	if got := protoGetString(pkg, packageName); got != "org.example.demo" {
		t.Errorf("package = %q", got)
	}
}

func TestBundleNativeTargeting(t *testing.T) {
	data, err := BundleNativeTargeting([]string{"arm64-v8a", "x86_64"})
	if err != nil {
		t.Fatal(err)
	}
	dirs, err := protoParse(data)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, d := range dirs {
		dir, _ := protoParse(d.b)
		targeting, _ := protoParse(dir[protoFind(dir, 2)].b)
		abi, _ := protoParse(targeting[protoFind(targeting, 1)].b)
		got = append(got, protoGetString(dir, 1)+"="+strconv.FormatUint(abi[protoFind(abi, 1)].v, 10))
	}
	if want := []string{"lib/arm64-v8a=3", "lib/x86_64=5"}; !slices.Equal(got, want) {
		t.Errorf("directories = %q", got)
	}
	if _, err := BundleNativeTargeting([]string{"mips"}); err == nil {
		t.Error("unknown ABI accepted")
	}
}

func TestPatchManifest(t *testing.T) {
	entries, err := ReadTemplate(templateAPK(t))
	if err != nil {
		t.Fatal(err)
	}
	m := Find(entries, "AndroidManifest.xml")
	data, err := m.Content()
	if err != nil {
		t.Fatal(err)
	}
	out, err := PatchManifest(data, ManifestValues{
		Package: "org.example.demo", Label: "Demo App", Version: "2.0", Build: 9, Library: "demo",
		GLESVersion: 0x30000, Orientation: "sensorPortrait", Permissions: []string{"android.permission.CAMERA"},
	})
	if err != nil {
		t.Fatal(err)
	}

	a, err := parseAXML(out)
	if err != nil {
		t.Fatal(err)
	}
	get := func(elem string, res uint32, plain string) string {
		for _, c := range a.startElements(elem) {
			if at := a.findAttr(c, res, plain); at != nil {
				s, _ := a.attrString(at)
				return s
			}
		}
		return ""
	}
	if got := get("manifest", 0, "package"); got != "org.example.demo" {
		t.Errorf("package = %q", got)
	}
	if got := get("manifest", attrVersionName, ""); got != "2.0" {
		t.Errorf("versionName = %q", got)
	}
	if got := get("application", attrLabel, ""); got != "Demo App" {
		t.Errorf("label = %q", got)
	}
	var perms []string
	for _, c := range a.startElements("uses-permission") {
		s, _ := a.attrString(a.findAttr(c, attrName, ""))
		perms = append(perms, s)
	}
	if !slices.Equal(perms, []string{"android.permission.CAMERA"}) {
		t.Errorf("permissions = %q", perms)
	}
	if got := get("meta-data", attrValue, ""); got != "demo" {
		t.Errorf("android.app.lib_name = %q", got)
	}
	data32 := func(elem string, res uint32) uint32 {
		at := a.findAttr(a.startElements(elem)[0], res, "")
		if at == nil {
			t.Fatalf("%s: attribute %#x missing", elem, res)
		}
		return binary.LittleEndian.Uint32(at[16:])
	}
	if got := data32("manifest", attrVersionCode); got != 9 {
		t.Errorf("versionCode = %d", got)
	}
	if got := data32("uses-feature", attrGLESVersion); got != 0x30000 {
		t.Errorf("glEsVersion = %#x", got)
	}
	if got := data32("activity", attrScreenOrientation); got != 7 {
		t.Errorf("screenOrientation = %d", got)
	}
	if _, err := PatchManifest(out, ManifestValues{Package: "a.b", Label: "x", Version: "1", Library: "x"}); err != nil {
		t.Errorf("repatch: %v", err)
	}

	none, err := PatchManifest(data, ManifestValues{Package: "a.b", Label: "x", Version: "1", Library: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := parseAXML(none); len(b.startElements("uses-permission")) != 0 {
		t.Error("template permission kept without --permissions")
	}
}

func TestResourcePaths(t *testing.T) {
	entries, err := ReadTemplate(templateAPK(t))
	if err != nil {
		t.Fatal(err)
	}
	table, err := Find(entries, "resources.arsc").Content()
	if err != nil {
		t.Fatal(err)
	}
	paths, err := ResourcePaths(table, "mipmap", "ic_launcher")
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 5 {
		t.Fatalf("paths = %v, want 5", paths)
	}
	for _, p := range paths {
		if Find(entries, p) == nil {
			t.Errorf("%s not in the template", p)
		}
	}
}

func TestRenamePackage(t *testing.T) {
	entries, err := ReadTemplate(templateAPK(t))
	if err != nil {
		t.Fatal(err)
	}
	table, err := Find(entries, "resources.arsc").Content()
	if err != nil {
		t.Fatal(err)
	}
	out, err := RenamePackage(table, "com.example.renamed")
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != len(table) {
		t.Fatalf("size %d, want %d", len(out), len(table))
	}
	hs := int(binary.LittleEndian.Uint16(out[2:]))
	_, size, err := arscPool(out[hs:])
	if err != nil {
		t.Fatal(err)
	}
	pkg := out[hs+size:]
	if binary.LittleEndian.Uint16(pkg) != arscPackage {
		t.Fatal("package chunk expected after the string pool")
	}
	var name []uint16
	for i := 0; i < arscNameLen; i++ {
		c := binary.LittleEndian.Uint16(pkg[arscNameOffset+2*i:])
		if c == 0 {
			break
		}
		name = append(name, c)
	}
	if got := string(utf16.Decode(name)); got != "com.example.renamed" {
		t.Errorf("package name = %q", got)
	}
	if _, err := ResourcePaths(out, "mipmap", "ic_launcher"); err != nil {
		t.Error(err)
	}
	if _, err := RenamePackage(table, strings.Repeat("a", arscNameLen)); err == nil {
		t.Error("over-long name accepted")
	}
}

func TestWriteZipAlignment(t *testing.T) {
	stored, _ := NewEntry("a/stored.bin", bytes.Repeat([]byte{1}, 1000), 0, Align)
	lib, _ := NewEntry("lib/arm64-v8a/libx.so", bytes.Repeat([]byte{2}, 5000), 0, AlignLib)
	deflated, _ := NewEntry("text.txt", bytes.Repeat([]byte("hello "), 200), 8, 0)
	data, err := WriteZip([]*Entry{deflated, stored, lib})
	if err != nil {
		t.Fatal(err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range zr.File {
		off, err := f.DataOffset()
		if err != nil {
			t.Fatal(err)
		}
		align := map[string]int64{"a/stored.bin": Align, "lib/arm64-v8a/libx.so": AlignLib, "text.txt": 1}[f.Name]
		if off%align != 0 {
			t.Errorf("%s data offset %d not aligned to %d", f.Name, off, align)
		}
		r, _ := f.Open()
		got, _ := io.ReadAll(r)
		if int64(len(got)) != int64(f.UncompressedSize64) {
			t.Errorf("%s size mismatch", f.Name)
		}
	}
}

func TestSignAPK(t *testing.T) {
	signer, err := DebugIdentity()
	if err != nil {
		t.Fatal(err)
	}
	entries, err := ReadTemplate(templateAPK(t))
	if err != nil {
		t.Fatal(err)
	}
	lib, _ := NewEntry("lib/arm64-v8a/libx.so", bytes.Repeat([]byte{7}, 70000), 0, AlignLib)
	entries, err = SignV1(append(entries, lib), signer)
	if err != nil {
		t.Fatal(err)
	}
	apk, err := WriteZip(entries)
	if err != nil {
		t.Fatal(err)
	}
	apk, err = SignV2(apk, signer)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(apk, []byte("APK Sig Block 42")) {
		t.Fatal("no signing block")
	}
	if _, err := zip.NewReader(bytes.NewReader(apk), int64(len(apk))); err != nil {
		t.Fatalf("signed apk unreadable: %v", err)
	}

	tools, _ := filepath.Glob(filepath.Join(os.Getenv("ANDROID_HOME"), "build-tools", "*", "apksigner"))
	if len(tools) == 0 {
		t.Skip("apksigner not found")
	}
	path := filepath.Join(t.TempDir(), "t.apk")
	os.WriteFile(path, apk, 0o644)
	out, err := exec.Command(tools[len(tools)-1], "verify", "-v", "--min-sdk-version", "23", path).CombinedOutput()
	if err != nil {
		t.Fatalf("apksigner: %v\n%s", err, out)
	}
	for _, want := range []string{"v1 scheme (JAR signing): true", "v2 scheme (APK Signature Scheme v2): true"} {
		if !bytes.Contains(out, []byte(want)) {
			t.Errorf("apksigner output lacks %q:\n%s", want, out)
		}
	}
}

func TestPatchManifestNoPermissionsTwice(t *testing.T) {
	entries, err := ReadTemplate(templateAPK(t))
	if err != nil {
		t.Fatal(err)
	}
	data, err := Find(entries, "AndroidManifest.xml").Content()
	if err != nil {
		t.Fatal(err)
	}
	values := ManifestValues{Package: "a.b", Label: "x", Version: "1", Library: "x"}
	out, err := PatchManifest(data, values)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := PatchManifest(out, values); err != nil {
		t.Errorf("repatch without permissions: %v", err)
	}
	bundle := bundleContent(t, templateBundle(t), "base/manifest/AndroidManifest.xml")
	if out, err = PatchBundleManifest(bundle, values); err != nil {
		t.Fatal(err)
	}
	if _, err := PatchBundleManifest(out, values); err != nil {
		t.Errorf("bundle repatch without permissions: %v", err)
	}
}
