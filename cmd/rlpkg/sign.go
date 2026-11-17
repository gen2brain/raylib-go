package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/apk"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/apple"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/authenticode"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/deb"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/pgp"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/rpm"
)

func runSign(args []string) error {
	c := &config{signer: "rlpkg"}
	var entitlements string
	fs := flag.NewFlagSet("sign", flag.ContinueOnError)
	fs.StringVar(&c.sign, "sign", "", "signing `identity`: a .p12 or PEM file with the key and certificate, or an exported OpenPGP secret key (default: android a debug key, darwin ad-hoc)")
	fs.StringVar(&c.signer, "signer", "rlpkg", "signing `tool`: rlpkg, or gpg with a gpg key in --sign")
	fs.StringVar(&c.profile, "profile", "", "ios: provisioning profile `file` (.mobileprovision), embedded and the source of the entitlements")
	fs.StringVar(&entitlements, "entitlements", "", "darwin: entitlements plist `file`")
	fs.BoolVar(&c.timestamp, "timestamp", true, "add a trusted timestamp when signing with a certificate")
	fs.StringVar(&c.tsaURL, "timestamp-url", authenticode.DefaultTimestampURL, "windows: RFC 3161 timestamp server `url`")
	fs.StringVar(&c.notaryKey, "notary-key", "", "darwin: App Store Connect API private key `file` (AuthKey_<id>.p8), notarizes and staples")
	fs.StringVar(&c.notaryID, "notary-key-id", "", "darwin: App Store Connect API key `id` (default: from the key file name)")
	fs.StringVar(&c.notaryIssue, "notary-issuer", "", "darwin: App Store Connect issuer `id`")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage:\n\n\trlpkg sign [flags] <file>\n\nSigns an existing .exe, .msix, .app, .ipa, .apk, .aab, .deb, .rpm, or Mach-O executable or dylib in place.\nWith an OpenPGP key any other file gets a detached <file>.asc signature.\n\nFlags:\n")
		printDefaults(fs)
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return errors.New("sign needs one file")
	}
	path := fs.Arg(0)
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if c.signer != "rlpkg" && c.signer != "gpg" {
		return fmt.Errorf("unknown signer %q", c.signer)
	}
	ext := strings.ToLower(filepath.Ext(path))
	if !info.IsDir() && (ext == ".deb" || ext == ".rpm" || c.pgpKey()) {
		return signPGP(c, path, ext)
	}

	switch {
	case info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".app"):
		if _, err := os.Stat(filepath.Join(path, "Contents", "Info.plist")); err == nil {
			return signMacApp(c, path, entitlements)
		}
		return signIOSApp(c, path)
	case info.IsDir():
		return fmt.Errorf("%s: not an .app bundle", path)
	}
	switch ext {
	case ".exe", ".dll":
		if c.sign == "" {
			return errors.New("--sign is required for Windows executables")
		}
		return signWindows(c, path)
	case ".msix", ".appx":
		if c.sign == "" {
			return errors.New("--sign is required for MSIX packages")
		}
		return signMSIX(c, path)
	case ".ipa":
		return signIPA(c, path)
	case ".apk", ".aab":
		return signAPK(c, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !apple.IsMachO(data) {
		return fmt.Errorf("%s: not a signable file", path)
	}
	signer, err := appleSigner(c)
	if err != nil {
		return err
	}
	id := strings.TrimSuffix(filepath.Base(path), ".dylib")
	signed, _, err := apple.SignMachO(data, signer, apple.MachOOptions{Identifier: id, HardenedRuntime: signer.Key != nil})
	if err != nil {
		return err
	}
	return os.WriteFile(path, signed, info.Mode())
}

func (c *config) pgpKey() bool {
	return c.signer == "gpg" || slices.Contains([]string{".asc", ".gpg", ".pgp"}, strings.ToLower(filepath.Ext(c.sign)))
}

func signPGP(c *config, path, ext string) error {
	sign, err := c.pgpSigner()
	if err != nil {
		return err
	}
	if ext != ".deb" && ext != ".rpm" {
		return signDetached(path, sign)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if ext == ".deb" {
		data, err = deb.Sign(data, sign)
	} else {
		data, err = rpm.Sign(data, sign)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return os.WriteFile(path, data, 0o644)
}

func signDetached(path string, sign pgp.Signer) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sig, err := sign(data)
	if err != nil {
		return err
	}
	if sig, err = pgp.Armor(sig); err != nil {
		return err
	}
	return os.WriteFile(path+".asc", sig, 0o644)
}

func signMacApp(c *config, app, entitlements string) error {
	notarize := c.notaryKey != ""
	if notarize && (c.sign == "" || c.notaryIssue == "") {
		return errors.New("notarization needs --sign, --notary-key and --notary-issuer")
	}
	var entData []byte
	if entitlements != "" {
		var err error
		if entData, err = os.ReadFile(entitlements); err != nil {
			return err
		}
	}
	signer, err := appleSigner(c)
	if err != nil {
		return err
	}
	os.Remove(filepath.Join(app, "Contents", "CodeResources"))
	cdhash, err := apple.SignBundle(app, apple.BundleOptions{
		MacOS: true, Signer: signer, Entitlements: entData, HardenedRuntime: signer.Key != nil,
	})
	if err != nil {
		return err
	}
	if !notarize {
		return nil
	}
	tmp, err := os.MkdirTemp("", "rlpkg-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	return notarizeDarwin(c, app, tmp, cdhash)
}

func signIOSApp(c *config, app string) error {
	if c.sign == "" || c.profile == "" {
		return errors.New("iOS apps need --sign and --profile")
	}
	infoData, err := os.ReadFile(filepath.Join(app, "Info.plist"))
	if err != nil {
		return err
	}
	info, err := apple.ParsePlist(infoData)
	if err != nil {
		return fmt.Errorf("Info.plist: %w", err)
	}
	bundleID, _ := info.(map[string]any)["CFBundleIdentifier"].(string)
	if bundleID == "" {
		return errors.New("Info.plist has no CFBundleIdentifier")
	}
	entitlements, err := embedProfile(c.profile, app, bundleID)
	if err != nil {
		return err
	}
	signer, err := appleSigner(c)
	if err != nil {
		return err
	}
	_, err = apple.SignBundle(app, apple.BundleOptions{Signer: signer, Entitlements: entitlements})
	return err
}

func signIPA(c *config, ipa string) error {
	tmp, err := os.MkdirTemp("", "rlpkg-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	if err := unzipDir(ipa, tmp); err != nil {
		return err
	}
	apps, _ := filepath.Glob(filepath.Join(tmp, "Payload", "*.app"))
	if len(apps) != 1 {
		return fmt.Errorf("%s: no Payload/*.app", ipa)
	}
	if err := signIOSApp(c, apps[0]); err != nil {
		return err
	}
	return zipDir(ipa, filepath.Join(tmp, "Payload"))
}

func signAPK(c *config, path string) error {
	entries, err := apk.ReadTemplate(path)
	if err != nil {
		return err
	}
	id, err := apk.DebugIdentity()
	if c.sign != "" {
		id, err = c.identity()
	}
	if err != nil {
		return err
	}
	bundle := strings.EqualFold(filepath.Ext(path), ".aab")
	if bundle {
		entries, err = apk.SignBundle(entries, id)
	} else {
		entries, err = apk.SignV1(entries, id)
	}
	if err != nil {
		return err
	}
	data, err := apk.WriteZip(entries)
	if err != nil {
		return err
	}
	if !bundle {
		if data, err = apk.SignV2(data, id); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0o644)
}

func unzipDir(archive, dir string) error {
	zr, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, f := range zr.File {
		rel := filepath.FromSlash(f.Name)
		if rel == "" || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
			return fmt.Errorf("%s: bad entry %q", archive, f.Name)
		}
		dst := filepath.Join(dir, rel)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(dst, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		r, err := f.Open()
		if err != nil {
			return err
		}
		w, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, f.Mode().Perm()|0o600)
		if err != nil {
			r.Close()
			return err
		}
		_, err = io.Copy(w, r)
		r.Close()
		if cerr := w.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			return err
		}
	}
	return nil
}
