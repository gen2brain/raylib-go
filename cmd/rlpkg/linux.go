package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/deb"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/icon"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/pgp"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/pkgtree"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/rpm"
)

var linuxIconSizes = []int{16, 22, 24, 32, 48, 64, 128, 256, 512}

func packageLinux(c *config) error {
	tmp, err := os.MkdirTemp("", "rlpkg-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	bin := filepath.Join(tmp, c.exe)
	if err := goBuild(c, c.goarch, bin, nil, nil); err != nil {
		return err
	}
	exe, err := os.ReadFile(bin)
	if err != nil {
		return err
	}

	img, err := c.loadIcon()
	if err != nil {
		return err
	}

	id := c.exe
	if c.idSet {
		id = c.id
	}
	top := fmt.Sprintf("%s-%s", c.exe, c.version)
	files := []tarFile{
		{path.Join(top, "bin", c.exe), 0o755, exe},
		{path.Join(top, "share", "applications", id+".desktop"), 0o644, desktopEntry(c, id)},
		{path.Join(top, "Makefile"), 0o644, makefile(c, id)},
	}
	for _, size := range linuxIconSizes {
		data, err := icon.PNG(icon.Resize(img, size))
		if err != nil {
			return err
		}
		dim := strconv.Itoa(size)
		files = append(files, tarFile{path.Join(top, "share", "icons", "hicolor", dim+"x"+dim, "apps", id+".png"), 0o644, data})
	}
	for _, f := range c.files {
		files = append(files, tarFile{path.Join(top, "share", c.exe, f.name), int64(f.mode), f.data})
	}

	var sign pgp.Signer
	if c.sign != "" {
		if sign, err = c.pgpSigner(); err != nil {
			return err
		}
	}

	for _, format := range c.formats {
		var archive string
		var err error
		switch format {
		case "targz":
			archive = filepath.Join(c.out, fmt.Sprintf("%s-%s-linux-%s.tar.gz", c.exe, c.version, c.goarch))
			if err = writeTarGz(archive, files); err == nil && sign != nil {
				err = signDetached(archive, sign)
			}
		case "deb":
			archive, err = writeDeb(c, id, exe, img, sign)
		case "rpm":
			archive, err = writeRPM(c, id, exe, img, sign)
		default:
			err = fmt.Errorf("unknown format %q (targz, deb, rpm)", format)
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, archive)
	}
	return nil
}

var debArch = map[string]string{"amd64": "amd64", "arm64": "arm64", "386": "i386", "arm": "armhf", "riscv64": "riscv64", "ppc64le": "ppc64el"}
var rpmArch = map[string]string{"amd64": "x86_64", "arm64": "aarch64", "386": "i686", "arm": "armv7hl", "riscv64": "riscv64", "ppc64le": "ppc64le"}

func installedFiles(c *config, id string, exe []byte, img image.Image) ([]pkgtree.File, error) {
	files := []pkgtree.File{
		{Path: "/usr/bin/" + c.exe, Mode: 0o755, Data: exe},
		{Path: "/usr/share/applications/" + id + ".desktop", Mode: 0o644, Data: installedDesktopEntry(c, id)},
	}
	for _, size := range linuxIconSizes {
		data, err := icon.PNG(icon.Resize(img, size))
		if err != nil {
			return nil, err
		}
		dim := strconv.Itoa(size)
		files = append(files, pkgtree.File{Path: "/usr/share/icons/hicolor/" + dim + "x" + dim + "/apps/" + id + ".png", Mode: 0o644, Data: data})
	}
	if len(c.files) > 0 {
		root := "/usr/share/" + c.exe
		dirs := []string{root}
		for _, f := range c.files {
			for d := path.Dir(f.name); d != "."; d = path.Dir(d) {
				if !slices.Contains(dirs, root+"/"+d) {
					dirs = append(dirs, root+"/"+d)
				}
			}
			files = append(files, pkgtree.File{Path: root + "/" + f.name, Mode: uint32(f.mode), Data: f.data})
		}
		for _, d := range dirs {
			files = append(files, pkgtree.File{Path: d, Mode: 0o755, Dir: true})
		}
	}
	return files, nil
}

func packageName(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '.' {
			b.WriteRune(r)
		} else {
			b.WriteByte('-')
		}
	}
	return b.String()
}

func writeDeb(c *config, id string, exe []byte, img image.Image, sign pgp.Signer) (string, error) {
	arch, ok := debArch[c.goarch]
	if !ok {
		return "", fmt.Errorf("no Debian architecture for %s", c.goarch)
	}
	files, err := installedFiles(c, id, exe, img)
	if err != nil {
		return "", err
	}
	name := packageName(c.exe)
	version := fmt.Sprintf("%s-%d", c.version, c.build)
	data, err := deb.Write(deb.Package{
		Name: name, Version: version, Arch: arch, Maintainer: c.vendor, Description: c.name, Section: "misc", Files: files, Sign: sign,
	})
	if err != nil {
		return "", err
	}
	out := filepath.Join(c.out, fmt.Sprintf("%s_%s_%s.deb", name, version, arch))
	return out, os.WriteFile(out, data, 0o644)
}

func writeRPM(c *config, id string, exe []byte, img image.Image, sign pgp.Signer) (string, error) {
	arch, ok := rpmArch[c.goarch]
	if !ok {
		return "", fmt.Errorf("no RPM architecture for %s", c.goarch)
	}
	files, err := installedFiles(c, id, exe, img)
	if err != nil {
		return "", err
	}
	name := packageName(c.exe)
	release := strconv.Itoa(c.build)
	data, err := rpm.Write(rpm.Package{
		Name: name, Version: c.version, Release: release, Arch: arch, Summary: c.name, Description: c.name,
		License: c.license, Vendor: c.vendor, Files: files, Sign: sign,
	})
	if err != nil {
		return "", err
	}
	out := filepath.Join(c.out, fmt.Sprintf("%s-%s-%s.%s.rpm", name, c.version, release, arch))
	return out, os.WriteFile(out, data, 0o644)
}

type tarFile struct {
	name string
	mode int64
	data []byte
}

func writeTarGz(archive string, files []tarFile) error {
	f, err := os.Create(archive)
	if err != nil {
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	now := time.Now()

	var dirs []string
	for _, file := range files {
		for dir := path.Dir(file.name); dir != "."; dir = path.Dir(dir) {
			if !slices.Contains(dirs, dir) {
				dirs = append(dirs, dir)
			}
		}
	}
	slices.Sort(dirs)
	for _, dir := range dirs {
		if err := tw.WriteHeader(&tar.Header{Typeflag: tar.TypeDir, Name: dir + "/", Mode: 0o755, ModTime: now}); err != nil {
			return err
		}
	}
	for _, file := range files {
		hdr := &tar.Header{Typeflag: tar.TypeReg, Name: file.name, Mode: file.mode, Size: int64(len(file.data)), ModTime: now}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write(file.data); err != nil {
			return err
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}
	if err := gw.Close(); err != nil {
		return err
	}
	return f.Close()
}

func desktopEntry(c *config, id string) []byte {
	category := strings.TrimSuffix(c.category, ";") + ";"
	name := strings.NewReplacer("\n", " ", "\r", " ").Replace(c.name)
	var dir string
	if len(c.files) > 0 {
		dir = "Path=@DATADIR@\n"
	}
	return fmt.Appendf(nil, `[Desktop Entry]
Type=Application
Name=%s
Exec=@BINDIR@/%s
%sIcon=%s
Categories=%s
Terminal=false
`, name, c.exe, dir, id, category)
}

func installedDesktopEntry(c *config, id string) []byte {
	return []byte(strings.NewReplacer("@BINDIR@", "/usr/bin", "@DATADIR@", "/usr/share/"+c.exe).Replace(string(desktopEntry(c, id))))
}

func makefile(c *config, id string) []byte {
	var sizes []string
	for _, size := range linuxIconSizes {
		sizes = append(sizes, strconv.Itoa(size))
	}
	var datadir, install, uninstall string
	if len(c.files) > 0 {
		datadir = " -e 's|@DATADIR@|$(PREFIX)/share/$(EXE)|'"
		install = "\tmkdir -p $(DESTDIR)$(PREFIX)/share/$(EXE)\n\tcp -R share/$(EXE)/. $(DESTDIR)$(PREFIX)/share/$(EXE)/\n"
		uninstall = "\tfind share/$(EXE) -type f | while IFS= read -r f; do rm -f \"$(DESTDIR)$(PREFIX)/$$f\"; done\n" +
			"\tfind share/$(EXE) -depth -type d | while IFS= read -r d; do rmdir \"$(DESTDIR)$(PREFIX)/$$d\" 2>/dev/null || true; done\n"
	}
	return fmt.Appendf(nil, `PREFIX ?= /usr/local

EXE = %s
ID = %s
SIZES = %s

install:
	install -Dm755 bin/$(EXE) $(DESTDIR)$(PREFIX)/bin/$(EXE)
	mkdir -p $(DESTDIR)$(PREFIX)/share/applications
	sed -e 's|@BINDIR@|$(PREFIX)/bin|'%s share/applications/$(ID).desktop > $(DESTDIR)$(PREFIX)/share/applications/$(ID).desktop
	for s in $(SIZES); do install -Dm644 share/icons/hicolor/$${s}x$${s}/apps/$(ID).png $(DESTDIR)$(PREFIX)/share/icons/hicolor/$${s}x$${s}/apps/$(ID).png; done
%s
uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/$(EXE)
	rm -f $(DESTDIR)$(PREFIX)/share/applications/$(ID).desktop
	for s in $(SIZES); do rm -f $(DESTDIR)$(PREFIX)/share/icons/hicolor/$${s}x$${s}/apps/$(ID).png; done
%s
user-install:
	$(MAKE) install PREFIX=$(HOME)/.local

user-uninstall:
	$(MAKE) uninstall PREFIX=$(HOME)/.local

.PHONY: install uninstall user-install user-uninstall
`, c.exe, id, strings.Join(sizes, " "), datadir, install, uninstall)
}
