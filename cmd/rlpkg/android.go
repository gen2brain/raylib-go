package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/apk"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/icon"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/keys"
)

//go:embed android/template.apk
var templateAPK []byte

//go:embed android/template.aab
var templateAAB []byte

type androidABI struct {
	abi    string
	triple string
}

var androidABIs = map[string]androidABI{
	"arm64": {"arm64-v8a", "aarch64-linux-android"},
	"arm":   {"armeabi-v7a", "armv7a-linux-androideabi"},
	"amd64": {"x86_64", "x86_64-linux-android"},
	"386":   {"x86", "i686-linux-android"},
}

var androidPermissions = map[string][]string{
	"camera":        {"android.permission.CAMERA"},
	"microphone":    {"android.permission.RECORD_AUDIO"},
	"location":      {"android.permission.ACCESS_COARSE_LOCATION", "android.permission.ACCESS_FINE_LOCATION"},
	"notifications": {"android.permission.POST_NOTIFICATIONS"},
	"internet":      {"android.permission.INTERNET"},
	"network":       {"android.permission.ACCESS_NETWORK_STATE", "android.permission.CHANGE_NETWORK_STATE"},
	"vibrate":       {"android.permission.VIBRATE"},
}

var androidOrientations = map[string]string{
	"":          "sensorLandscape",
	"landscape": "sensorLandscape",
	"portrait":  "sensorPortrait",
	"any":       "fullSensor",
}

const androidAPI = "24"

type androidLib struct {
	abi  string
	data []byte
}

func packageAndroid(c *config) error {
	var archs []string
	for _, a := range splitList(c.goarch) {
		if _, ok := androidABIs[a]; !ok {
			return fmt.Errorf("unknown android arch %q (arm64, arm, amd64, 386)", a)
		}
		archs = append(archs, a)
	}

	toolchain, err := ndkToolchain()
	if err != nil {
		return err
	}

	tmp, err := os.MkdirTemp("", "rlpkg-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	lib := "lib" + c.exe + ".so"
	var libs []androidLib
	for _, arch := range archs {
		abi := androidABIs[arch]
		out := filepath.Join(tmp, abi.abi, lib)
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		clang := filepath.Join(toolchain, abi.triple+androidAPI+"-clang")
		env := []string{
			"GOOS=android", "GOARCH=" + arch, "CGO_ENABLED=1",
			"CC=" + clang, "CXX=" + clang + "++", "PKG_CONFIG=true",
		}
		ldflags := []string{"-extldflags=-Wl,-soname," + lib}
		if err := goBuildWith(c, env, out, nil, ldflags, []string{"-buildmode=c-shared"}); err != nil {
			return err
		}
		so, err := os.ReadFile(out)
		if err != nil {
			return err
		}
		libs = append(libs, androidLib{abi.abi, so})
	}

	var perms []string
	for _, p := range c.permissions {
		perms = append(perms, androidPermissions[p]...)
	}
	gles := uint32(0x20000)
	if c.hasTag("es3") {
		gles = 0x30000
	}
	values := apk.ManifestValues{
		Package: c.id, Label: c.name, Version: c.version, Build: c.build, Library: c.exe,
		GLESVersion: gles, Orientation: androidOrientations[c.orientation], Permissions: perms,
	}

	var id *keys.Identity
	if c.sign == "" {
		id, err = apk.DebugIdentity()
	} else {
		id, err = c.identity()
	}
	if err != nil {
		return err
	}

	for _, format := range c.formats {
		var out string
		if format == "aab" {
			out, err = packageAndroidBundle(c, values, lib, libs, id)
		} else {
			out, err = packageAndroidAPK(c, values, lib, libs, id)
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, out)
		if c.install && format == "apk" {
			if err := run("adb", "install", "-r", out); err != nil {
				return err
			}
		}
	}
	return nil
}

func readTemplate(name string, data []byte) ([]*apk.Entry, error) {
	entries, err := apk.ReadArchive(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	return entries, nil
}

func packageAndroidAPK(c *config, values apk.ManifestValues, lib string, libs []androidLib, id *keys.Identity) (string, error) {
	entries, err := readTemplate("template.apk", templateAPK)
	if err != nil {
		return "", err
	}

	manifest := apk.Find(entries, "AndroidManifest.xml")
	arsc := apk.Find(entries, "resources.arsc")
	if manifest == nil || arsc == nil {
		return "", errors.New("template.apk: not an APK template")
	}

	data, err := manifest.Content()
	if err != nil {
		return "", err
	}
	if data, err = apk.PatchManifest(data, values); err != nil {
		return "", fmt.Errorf("template.apk: %w", err)
	}
	if err := manifest.SetContent(data); err != nil {
		return "", err
	}

	table, err := arsc.Content()
	if err != nil {
		return "", err
	}
	if table, err = apk.RenamePackage(table, c.id); err != nil {
		return "", fmt.Errorf("template.apk: %w", err)
	}
	if err := arsc.SetContent(table); err != nil {
		return "", err
	}

	icons, err := apk.ResourcePaths(table, "mipmap", "ic_launcher")
	if err != nil {
		return "", err
	}
	if err := replaceLauncherIcons(c, entries, icons); err != nil {
		return "", err
	}

	for _, l := range libs {
		e, err := apk.NewEntry("lib/"+l.abi+"/"+lib, l.data, zip.Store, apk.AlignLib)
		if err != nil {
			return "", err
		}
		entries = append(entries, e)
	}
	if entries, err = addAndroidData(c, entries, "assets/"); err != nil {
		return "", err
	}

	if entries, err = apk.SignV1(entries, id); err != nil {
		return "", err
	}
	if data, err = apk.WriteZip(entries); err != nil {
		return "", err
	}
	if data, err = apk.SignV2(data, id); err != nil {
		return "", err
	}

	out := filepath.Join(c.out, c.exe+".apk")
	return out, os.WriteFile(out, data, 0o644)
}

func packageAndroidBundle(c *config, values apk.ManifestValues, lib string, libs []androidLib, id *keys.Identity) (string, error) {
	entries, err := readTemplate("template.aab", templateAAB)
	if err != nil {
		return "", err
	}

	manifest := apk.Find(entries, "base/manifest/AndroidManifest.xml")
	table := apk.Find(entries, "base/resources.pb")
	if manifest == nil || table == nil {
		return "", errors.New("template.aab: not an app bundle template")
	}

	data, err := manifest.Content()
	if err != nil {
		return "", err
	}
	if data, err = apk.PatchBundleManifest(data, values); err != nil {
		return "", fmt.Errorf("template.aab: %w", err)
	}
	if err := manifest.SetContent(data); err != nil {
		return "", err
	}

	if data, err = table.Content(); err != nil {
		return "", err
	}
	if data, err = apk.RenameBundlePackage(data, c.id); err != nil {
		return "", fmt.Errorf("template.aab: %w", err)
	}
	if err := table.SetContent(data); err != nil {
		return "", err
	}

	var icons []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name, "base/res/mipmap-") && strings.HasSuffix(e.Name, "/ic_launcher.png") {
			icons = append(icons, e.Name)
		}
	}
	if len(icons) == 0 {
		return "", errors.New("template.aab: no launcher icons")
	}
	if err := replaceLauncherIcons(c, entries, icons); err != nil {
		return "", err
	}

	var abis []string
	for _, l := range libs {
		e, err := apk.NewEntry("base/lib/"+l.abi+"/"+lib, l.data, zip.Deflate, 0)
		if err != nil {
			return "", err
		}
		entries = append(entries, e)
		abis = append(abis, l.abi)
	}
	if data, err = apk.BundleNativeTargeting(abis); err != nil {
		return "", err
	}
	e, err := apk.NewEntry("base/native.pb", data, zip.Deflate, 0)
	if err != nil {
		return "", err
	}
	entries = append(entries, e)
	if entries, err = addAndroidData(c, entries, "base/assets/"); err != nil {
		return "", err
	}

	if entries, err = apk.SignBundle(entries, id); err != nil {
		return "", err
	}
	if data, err = apk.WriteZip(entries); err != nil {
		return "", err
	}

	out := filepath.Join(c.out, c.exe+".aab")
	return out, os.WriteFile(out, data, 0o644)
}

func addAndroidData(c *config, entries []*apk.Entry, dir string) ([]*apk.Entry, error) {
	for _, f := range c.files {
		name := dir + f.name
		if apk.Find(entries, name) != nil {
			return nil, fmt.Errorf("--data: %s is already in the package", f.name)
		}
		e, err := apk.NewEntry(name, f.data, zip.Deflate, 0)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func replaceLauncherIcons(c *config, entries []*apk.Entry, paths []string) error {
	img, err := c.loadIcon()
	if err != nil {
		return err
	}
	for _, p := range paths {
		e := apk.Find(entries, p)
		if e == nil {
			return fmt.Errorf("template icon %s missing", p)
		}
		old, err := e.Content()
		if err != nil {
			return err
		}
		cfg, err := png.DecodeConfig(bytes.NewReader(old))
		if err != nil {
			return fmt.Errorf("template icon %s: %w", p, err)
		}
		data, err := icon.PNG(icon.Resize(img, cfg.Width))
		if err != nil {
			return err
		}
		if err := e.SetContent(data); err != nil {
			return err
		}
	}
	return nil
}

func ndkToolchain() (string, error) {
	var roots []string
	for _, v := range []string{"ANDROID_NDK_HOME", "ANDROID_NDK_ROOT"} {
		if d := os.Getenv(v); d != "" {
			roots = append(roots, d)
		}
	}
	for _, v := range []string{"ANDROID_SDK_ROOT", "ANDROID_HOME"} {
		if d := os.Getenv(v); d != "" {
			versions, _ := filepath.Glob(filepath.Join(d, "ndk", "*"))
			slices.Sort(versions)
			slices.Reverse(versions)
			roots = append(roots, versions...)
		}
	}

	host := runtime.GOOS + "-x86_64"
	for _, root := range roots {
		dir := filepath.Join(root, "toolchains", "llvm", "prebuilt", host, "bin")
		if _, err := os.Stat(filepath.Join(dir, "clang"+exeSuffix())); err == nil {
			return dir, nil
		}
	}
	return "", errors.New("no Android NDK found; set ANDROID_NDK_HOME")
}

func exeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}
