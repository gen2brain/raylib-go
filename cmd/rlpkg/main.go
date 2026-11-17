package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/apple"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/authenticode"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/icon"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/keys"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/pgp"
)

//go:embed icon.png
var defaultIcon []byte

const usageText = `rlpkg builds and packages raylib-go applications.

Usage:

	rlpkg package [flags] [package]
	rlpkg sign [flags] <file>
	rlpkg serve [--addr localhost:8000] <directory>
	rlpkg staple <app>

Targets (--os):

	windows  .exe with icon, version info and manifest, --format msix; Authenticode with --sign
	darwin   .app bundle and a zip of it; signed, notarized and stapled with --sign and --notary-*
	linux    .tar.gz with .desktop file, hicolor icons and a Makefile; --format deb,rpm; OpenPGP signatures with --sign
	android  .apk, needs only the NDK, --format aab for an App Bundle; signed
	ios      .ipa signed with --sign and --profile
	js       directory with the wasm module and loader, served by rlpkg serve

Signing is done by rlpkg itself on every platform; on macOS --signer codesign
uses a keychain identity instead, for linux --signer gpg a key of the gpg program.
Password of a .p12 file: RLPKG_P12_PASSWORD, of an OpenPGP key: RLPKG_GPG_PASSPHRASE.
sign takes an existing .exe, .msix, .app, .ipa, .apk, .deb, .rpm or Mach-O file, and
with an OpenPGP key any other file; run 'rlpkg sign -h' for its flags.

Flags:
`

type config struct {
	goos        string
	goarch      string
	out         string
	name        string
	id          string
	version     string
	build       int
	icon        string
	tags        []string
	ldflags     string
	cgo         bool
	debug       bool
	console     bool
	install     bool
	simulator   bool
	minOS       string
	profile     string
	orientation string
	category    string
	formats     []string
	vendor      string
	publisher   string
	copyright   string
	license     string
	permissions []string
	data        []string
	sign        string
	signer      string
	timestamp   bool
	tsaURL      string
	notaryKey   string
	notaryID    string
	notaryIssue string

	idSet   bool
	archSet bool
	pkg     string
	pkgDir  string
	exe     string
	rlDir   string
	files   []dataFile
}

type dataFile struct {
	name string
	mode os.FileMode
	data []byte
}

func main() {
	if len(os.Args) < 2 {
		usage(nil)
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "package":
		err = runPackage(os.Args[2:])
	case "sign":
		err = runSign(os.Args[2:])
	case "serve":
		err = runServe(os.Args[2:])
	case "staple":
		if len(os.Args) != 3 {
			err = errors.New("usage: rlpkg staple <app>")
		} else {
			err = apple.Staple(os.Args[2], nil)
		}
	case "help", "-h", "-help", "--help":
		usage(newFlagSet(&config{}))
	default:
		fmt.Fprintf(os.Stderr, "rlpkg: unknown command %q\n", os.Args[1])
		usage(nil)
		os.Exit(2)
	}

	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintln(os.Stderr, "rlpkg:", err)
		os.Exit(1)
	}
}

func usage(fs *flag.FlagSet) {
	fmt.Fprint(os.Stderr, usageText)
	if fs != nil {
		printDefaults(fs)
	} else {
		fmt.Fprintln(os.Stderr, "\trun 'rlpkg help' for the flag list")
	}
}

func printDefaults(fs *flag.FlagSet) {
	fs.VisitAll(func(f *flag.Flag) {
		name, usage := flag.UnquoteUsage(f)
		line := "  --" + f.Name
		if name != "" {
			line += " " + name
		}
		line += "\n    \t" + usage
		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" {
			if fmt.Sprintf("%T", f.Value) == "*flag.stringValue" {
				line += fmt.Sprintf(" (default %q)", f.DefValue)
			} else {
				line += " (default " + f.DefValue + ")"
			}
		}
		fmt.Fprintln(os.Stderr, line)
	})
}

func newFlagSet(c *config) *flag.FlagSet {
	fs := flag.NewFlagSet("package", flag.ContinueOnError)
	fs.StringVar(&c.goos, "os", runtime.GOOS, "target `os`: windows, darwin, linux, android, ios, js")
	fs.StringVar(&c.goarch, "arch", runtime.GOARCH, "target `arch`; darwin also accepts universal, android a comma-separated list; android and ios use arm64 unless set")
	fs.StringVar(&c.out, "out", ".", "output `directory`")
	fs.StringVar(&c.name, "name", "", "application `name` (default: the executable name)")
	fs.StringVar(&c.id, "id", "", "application `identifier` in reverse DNS form (default: com.example.<executable>; linux: the executable name)")
	fs.StringVar(&c.version, "version", "1.0.0", "application `version`")
	fs.IntVar(&c.build, "build", 1, "build `number`")
	fs.StringVar(&c.icon, "icon", "", "square PNG `file`, 1024 px recommended (default: the raylib logo)")
	fs.Func("tags", "comma-separated list of build `tags`", func(s string) error {
		c.tags = splitList(s)
		return nil
	})
	fs.StringVar(&c.ldflags, "ldflags", "", "extra linker `flags`")
	fs.BoolVar(&c.cgo, "cgo", false, "build with cgo (default: go env CGO_ENABLED for the host platform, off when cross-compiling)")
	fs.BoolVar(&c.debug, "debug", false, "keep symbols, debug info and file paths")
	fs.BoolVar(&c.console, "console", false, "windows: build a console application")
	fs.BoolVar(&c.install, "install", false, "android, ios: install on the connected device")
	fs.BoolVar(&c.simulator, "simulator", false, "ios: build for the simulator")
	fs.StringVar(&c.minOS, "minos", "15.0", "ios: minimum iOS `version`")
	fs.StringVar(&c.profile, "profile", "", "ios: provisioning profile `file` (.mobileprovision)")
	fs.StringVar(&c.orientation, "orientation", "", "android, ios: screen `orientation`: landscape, portrait, any (default: android landscape, ios any)")
	fs.StringVar(&c.category, "category", "Game", "linux: desktop entry `categories`, separated by ;")
	fs.Func("format", "comma-separated `formats`, linux: targz, deb, rpm (default targz); windows: msix, written next to the .exe; android: apk, aab (default apk)", func(s string) error {
		c.formats = splitList(s)
		return nil
	})
	fs.StringVar(&c.vendor, "vendor", "", "`vendor`: Windows company name, Debian maintainer, RPM vendor (default: the application name)")
	fs.StringVar(&c.publisher, "publisher", "", "windows msix: publisher `name` of the package identity (default: the subject of the --sign certificate; the Store assigns one)")
	fs.StringVar(&c.copyright, "copyright", "", "copyright `line` for Windows")
	fs.StringVar(&c.license, "license", "Unknown", "license `name` for rpm")
	fs.Func("permissions", "comma-separated `list`: camera, microphone, location, notifications, internet, network, vibrate", func(s string) error {
		c.permissions = splitList(s)
		return nil
	})
	fs.Func("data", "`file` or directory shipped with the application, loaded with relative paths; a directory keeps its name; repeatable", func(s string) error {
		c.data = append(c.data, s)
		return nil
	})
	fs.StringVar(&c.sign, "sign", "", "signing `identity`: a .p12 or PEM file with the key and certificate, linux: an exported OpenPGP secret key; a keychain identity with --signer codesign, a gpg key with --signer gpg (default: android a debug key, darwin ad-hoc, others unsigned)")
	fs.StringVar(&c.signer, "signer", "", "signing `tool`: rlpkg, darwin, ios: codesign, linux: gpg (default: codesign on macOS, else rlpkg)")
	fs.BoolVar(&c.timestamp, "timestamp", true, "add a trusted timestamp when signing with a certificate")
	fs.StringVar(&c.tsaURL, "timestamp-url", authenticode.DefaultTimestampURL, "windows: RFC 3161 timestamp server `url`")
	fs.StringVar(&c.notaryKey, "notary-key", "", "darwin: App Store Connect API private key `file` (AuthKey_<id>.p8), enables notarization")
	fs.StringVar(&c.notaryID, "notary-key-id", "", "darwin: App Store Connect API key `id` (default: from the key file name)")
	fs.StringVar(&c.notaryIssue, "notary-issuer", "", "darwin: App Store Connect issuer `id`")
	fs.Usage = func() { usage(fs) }
	return fs
}

func runPackage(args []string) error {
	c := &config{}
	fs := newFlagSet(c)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 1 {
		return errors.New("only one package can be packaged at a time")
	}

	c.pkg = "."
	if fs.NArg() == 1 {
		c.pkg = fs.Arg(0)
	}

	cgoSet := false
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "cgo":
			cgoSet = true
		case "id":
			c.idSet = true
		case "arch":
			c.archSet = true
		}
	})
	if !cgoSet {
		c.cgo = defaultCgo(c.goos, c.goarch)
	}
	if c.goos == "wasm" {
		c.goos = "js"
	}
	if c.goos == "js" {
		c.cgo = false
		c.goarch = "wasm"
	}
	if c.goos == "android" || c.goos == "ios" {
		c.cgo = true
		if !c.archSet {
			c.goarch = "arm64"
		}
	}

	if c.signer == "" {
		c.signer = "rlpkg"
		if runtime.GOOS == "darwin" {
			c.signer = "codesign"
		}
	}
	if !slices.Contains([]string{"rlpkg", "codesign", "gpg"}, c.signer) {
		return fmt.Errorf("unknown signer %q", c.signer)
	}
	if c.signer == "gpg" && c.goos != "linux" {
		return errors.New("--signer gpg is for linux packages")
	}
	if c.signer == "codesign" && runtime.GOOS != "darwin" {
		return errors.New("--signer codesign needs a macOS host")
	}

	for _, p := range c.permissions {
		if !slices.Contains([]string{"camera", "microphone", "location", "notifications", "internet", "network", "vibrate"}, p) {
			return fmt.Errorf("unknown permission %q", p)
		}
	}
	if c.orientation != "" && c.goos != "android" && c.goos != "ios" {
		return errors.New("--orientation is for android and ios packages")
	}
	if c.orientation != "" && !slices.Contains([]string{"landscape", "portrait", "any"}, c.orientation) {
		return fmt.Errorf("unknown orientation %q (landscape, portrait, any)", c.orientation)
	}

	if err := c.loadData(); err != nil {
		return err
	}

	if err := c.resolve(); err != nil {
		return err
	}

	if err := os.MkdirAll(c.out, 0o755); err != nil {
		return err
	}

	switch c.goos {
	case "windows":
		return packageWindows(c)
	case "darwin":
		return packageDarwin(c)
	case "linux":
		return packageLinux(c)
	case "android":
		return packageAndroid(c)
	case "ios":
		return packageIOS(c)
	case "js":
		return packageJS(c)
	}
	return fmt.Errorf("unsupported target os %q", c.goos)
}

func (c *config) loadData() error {
	seen := map[string]string{}
	add := func(name, src string) error {
		if prev, ok := seen[strings.ToLower(name)]; ok {
			return fmt.Errorf("--data: %s and %s both become %s", prev, src, name)
		}
		seen[strings.ToLower(name)] = src
		info, err := os.Stat(src)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		mode := os.FileMode(0o644)
		if info.Mode()&0o111 != 0 {
			mode = 0o755
		}
		c.files = append(c.files, dataFile{name, mode, data})
		return nil
	}

	for _, root := range c.data {
		info, err := os.Stat(root)
		if err != nil {
			return fmt.Errorf("--data: %w", err)
		}
		base := filepath.Base(filepath.Clean(root))
		if !info.IsDir() {
			if err := add(base, root); err != nil {
				return err
			}
			continue
		}
		err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(root, p)
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() {
				if info, err := os.Stat(p); err != nil || !info.Mode().IsRegular() {
					return nil
				}
			}
			return add(path.Join(base, filepath.ToSlash(rel)), p)
		})
		if err != nil {
			return fmt.Errorf("--data: %w", err)
		}
	}
	return nil
}

func (c *config) writeData(root string, reserved ...string) error {
	for _, f := range c.files {
		top, _, _ := strings.Cut(f.name, "/")
		for _, r := range reserved {
			if strings.EqualFold(top, r) {
				return fmt.Errorf("--data: %s is reserved in the package", f.name)
			}
		}
		dst := filepath.Join(root, filepath.FromSlash(f.name))
		if _, err := os.Lstat(dst); err == nil {
			return fmt.Errorf("--data: %s is already in the package", f.name)
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, f.data, f.mode); err != nil {
			return err
		}
	}
	return nil
}

func (c *config) resolve() error {
	out, err := goList(c, "-f", "{{.Name}}\t{{.Dir}}", c.pkg)
	if err != nil {
		return err
	}
	name, dir, _ := strings.Cut(strings.TrimSpace(out), "\t")
	if name != "main" {
		return fmt.Errorf("%s is not a main package", c.pkg)
	}
	c.pkgDir = dir
	c.exe = filepath.Base(dir)

	c.rlDir, err = raylibModuleDir(c)
	if err != nil {
		return err
	}

	if c.name == "" {
		c.name = c.exe
	}
	if c.id == "" {
		c.id = "com.example." + identifierPart(c.exe)
	}
	if c.vendor == "" {
		c.vendor = c.name
	}
	if len(c.formats) == 0 && c.goos == "linux" {
		c.formats = []string{"targz"}
	}
	if len(c.formats) == 0 && c.goos == "android" {
		c.formats = []string{"apk"}
	}
	if c.goos == "android" {
		for _, format := range c.formats {
			if format != "apk" && format != "aab" {
				return fmt.Errorf("unknown format %q (apk, aab)", format)
			}
		}
		if c.install && !slices.Contains(c.formats, "apk") {
			return errors.New("--install needs the apk format")
		}
	}
	if c.goos == "windows" {
		for _, format := range c.formats {
			if format != "msix" {
				return fmt.Errorf("unknown format %q (msix)", format)
			}
		}
	}
	return nil
}

func (c *config) loadIcon() (image.Image, error) {
	if c.icon != "" {
		return icon.Load(c.icon)
	}
	return png.Decode(bytes.NewReader(defaultIcon))
}

func (c *config) identity() (*keys.Identity, error) {
	if c.sign == "" {
		return nil, errors.New("--sign is required")
	}
	return keys.Load(c.sign)
}

func (c *config) pgpSigner() (pgp.Signer, error) {
	if c.sign == "" {
		return nil, errors.New("--sign is required")
	}
	if c.signer == "gpg" {
		return pgp.GPG(c.sign), nil
	}
	return pgp.Load(c.sign)
}

func appleSigner(c *config) (*apple.Signer, error) {
	if c.sign == "" {
		return &apple.Signer{}, nil
	}
	id, err := c.identity()
	if err != nil {
		return nil, err
	}
	chain, err := apple.Chain(id.Cert, id.Extra)
	if err != nil {
		return nil, err
	}
	return &apple.Signer{Key: id.Key, Chain: chain, Timestamp: c.timestamp}, nil
}

func defaultCgo(goos, goarch string) bool {
	if goos == "darwin" && goarch == "universal" {
		goarch = runtime.GOARCH
	}
	if goos != runtime.GOOS || goarch != runtime.GOARCH {
		return os.Getenv("CGO_ENABLED") == "1"
	}
	out, err := exec.Command("go", "env", "CGO_ENABLED").Output()
	return err == nil && strings.TrimSpace(string(out)) == "1"
}

func goList(c *config, args ...string) (string, error) {
	goarch, _, _ := strings.Cut(c.goarch, ",")
	if goarch == "universal" {
		goarch = "arm64"
	}
	list := []string{"list"}
	if len(c.tags) > 0 {
		list = append(list, "-tags", strings.Join(c.tags, ","))
	}
	args = append(list, args...)
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(), c.goEnv(goarch)...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	return string(out), err
}

func raylibModuleDir(c *config) (string, error) {
	const module = "github.com/gen2brain/raylib-go/raylib"
	out, err := goList(c, "-m", "-f", "{{.Dir}}", module)
	if err != nil {
		return "", fmt.Errorf("%s is not a dependency of %s", module, c.pkg)
	}
	if dir := strings.TrimSpace(out); dir != "" {
		return dir, nil
	}
	out, err = goList(c, "-mod=mod", "-m", "-f", "{{.Dir}}", module)
	if dir := strings.TrimSpace(out); err == nil && dir != "" {
		return dir, nil
	}
	cmd := exec.Command("go", "mod", "download", "-json", module)
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cannot download %s: %w", module, err)
	}
	var info struct{ Dir, Error string }
	if err := json.Unmarshal(data, &info); err != nil || info.Dir == "" {
		return "", fmt.Errorf("cannot download %s: %s", module, info.Error)
	}
	return info.Dir, nil
}

func (c *config) goEnv(goarch string) []string {
	return []string{"GOOS=" + c.goos, "GOARCH=" + goarch, "CGO_ENABLED=" + strconv.Itoa(boolInt(c.cgo))}
}

func cgoFlags(flags ...string) []string {
	var env []string
	for _, name := range []string{"CGO_CFLAGS", "CGO_LDFLAGS"} {
		value := os.Getenv(name)
		if value == "" {
			value = "-O2 -g"
		}
		env = append(env, name+"="+value+" "+strings.Join(flags, " "))
	}
	return env
}

func goBuild(c *config, goarch, out string, tags, ldflags []string) error {
	return goBuildWith(c, c.goEnv(goarch), out, tags, ldflags, nil)
}

func goBuildWith(c *config, env []string, out string, tags, ldflags, extra []string) error {
	args := append([]string{"build", "-o", out}, extra...)
	tags = append(slices.Clone(c.tags), tags...)
	if len(tags) > 0 {
		args = append(args, "-tags", strings.Join(tags, ","))
	}
	if !c.debug {
		args = append(args, "-trimpath")
		ldflags = append(ldflags, "-s", "-w")
	}
	if c.ldflags != "" {
		ldflags = append(ldflags, c.ldflags)
	}
	if len(ldflags) > 0 {
		args = append(args, "-ldflags", strings.Join(ldflags, " "))
	}
	args = append(args, c.pkg)

	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *config) hasPermission(p string) bool {
	return slices.Contains(c.permissions, p)
}

func (c *config) hasTag(tag string) bool {
	return slices.Contains(c.tags, tag)
}

func splitList(s string) []string {
	var list []string
	for f := range strings.SplitSeq(s, ",") {
		if f = strings.TrimSpace(f); f != "" {
			list = append(list, f)
		}
	}
	return list
}

func identifierPart(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	id := b.String()
	if id == "" || id[0] >= '0' && id[0] <= '9' {
		id = "app" + id
	}
	return id
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}
