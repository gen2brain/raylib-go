package main

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"fmt"
	"image"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/authenticode"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/icon"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/keys"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/msix"
	"github.com/tc-hib/winres"
	"github.com/tc-hib/winres/version"
)

//go:embed windows/app.manifest
var windowsManifest []byte

func packageWindows(c *config) error {
	existing, err := filepath.Glob(filepath.Join(c.pkgDir, "*.syso"))
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return fmt.Errorf("%s already has resources (%s); remove them, rlpkg writes its own", c.pkg, filepath.Base(existing[0]))
	}

	img, err := c.loadIcon()
	if err != nil {
		return err
	}
	ico, err := winres.NewIconFromResizedImage(img, []int{256, 64, 48, 40, 32, 24, 20, 16})
	if err != nil {
		return err
	}
	rs := &winres.ResourceSet{}
	if err := rs.SetIcon(winres.Name("GLFW_ICON"), ico); err != nil {
		return err
	}
	rs.SetVersionInfo(versionInfo(c))
	if err := rs.Set(winres.RT_MANIFEST, winres.ID(1), winres.LCIDDefault, windowsManifest); err != nil {
		return err
	}

	var obj bytes.Buffer
	if err := rs.WriteObject(&obj, winres.Arch(c.goarch)); err != nil {
		return err
	}
	syso := filepath.Join(c.pkgDir, "rlpkg_windows_"+c.goarch+".syso")
	if err := os.WriteFile(syso, obj.Bytes(), 0o644); err != nil {
		return err
	}
	defer os.Remove(syso)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	defer signal.Stop(stop)
	go func() {
		if _, ok := <-stop; ok {
			os.Remove(syso)
			os.Exit(130)
		}
	}()

	var ldflags []string
	if !c.console {
		ldflags = append(ldflags, "-H=windowsgui")
	}
	out := filepath.Join(c.out, c.exe+".exe")
	if err := goBuild(c, c.goarch, out, nil, ldflags); err != nil {
		return err
	}

	if c.sign != "" {
		if err := signWindows(c, out); err != nil {
			return err
		}
	}
	fmt.Fprintln(os.Stderr, out)
	if len(c.files) > 0 {
		if err := zipWindows(c, out); err != nil {
			return err
		}
	}
	if slices.Contains(c.formats, "msix") {
		return packageMSIX(c, out, img, ldflags)
	}
	return nil
}

func zipWindows(c *config, exe string) error {
	tmp, err := os.MkdirTemp("", "rlpkg-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	stage := filepath.Join(tmp, c.exe)
	if err := os.MkdirAll(stage, 0o755); err != nil {
		return err
	}
	data, err := os.ReadFile(exe)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(stage, c.exe+".exe"), data, 0o755); err != nil {
		return err
	}
	if err := c.writeData(stage); err != nil {
		return err
	}
	archive := filepath.Join(c.out, fmt.Sprintf("%s-%s-windows-%s.zip", c.exe, c.version, c.goarch))
	if err := zipDir(archive, stage); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, archive)
	return nil
}

var msixArch = map[string]string{"amd64": "x64", "386": "x86", "arm64": "arm64", "arm": "arm"}

var msixLogos = []struct {
	name string
	size int
}{{"StoreLogo", 50}, {"Square44x44Logo", 44}, {"Square150x150Logo", 150}}

var msixCapabilities = map[string]string{
	"internet":   `<Capability Name="internetClient"/>`,
	"camera":     `<DeviceCapability Name="webcam"/>`,
	"microphone": `<DeviceCapability Name="microphone"/>`,
	"location":   `<DeviceCapability Name="location"/>`,
}

func packageMSIX(c *config, exe string, img image.Image, ldflags []string) error {
	arch, ok := msixArch[c.goarch]
	if !ok {
		return fmt.Errorf("no MSIX architecture for %s", c.goarch)
	}
	if !msix.ValidName(c.id) {
		return fmt.Errorf("--id %q is not an MSIX identity name: 3 to 50 letters, digits, periods and dashes", c.id)
	}
	var signer *authenticode.Signer
	publisher := c.publisher
	if c.sign != "" {
		var id *keys.Identity
		var err error
		if signer, id, err = windowsSigner(c); err != nil {
			return err
		}
		subject, err := msix.Publisher(id.Cert)
		if err != nil {
			return err
		}
		if publisher != "" && publisher != subject {
			return fmt.Errorf("--publisher %q is not the certificate subject %q", publisher, subject)
		}
		publisher = subject
	}
	if publisher == "" {
		publisher = msix.CommonName(c.vendor)
	}
	if signer == nil {
		fmt.Fprintln(os.Stderr, "rlpkg: the .msix is unsigned, Windows installs it only after it is signed (--sign, rlpkg sign, or the Store)")
	}

	var files []msix.File
	if !c.cgo {
		tmp, err := os.MkdirTemp("", "rlpkg-")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmp)
		exe = filepath.Join(tmp, c.exe+".exe")
		if err := goBuild(c, c.goarch, exe, []string{"extlib"}, ldflags); err != nil {
			return err
		}
		if c.sign != "" {
			if err := signWindows(c, exe); err != nil {
				return err
			}
		}
		lib, err := readGzip(filepath.Join(c.rlDir, "libs", "windows_"+c.goarch, "raylib.dll.gz"))
		if err != nil {
			return err
		}
		if c.sign != "" {
			if lib, err = signWindowsData(c, lib); err != nil {
				return err
			}
		}
		files = append(files, msix.File{Name: "raylib.dll", Data: lib})
	}
	data, err := os.ReadFile(exe)
	if err != nil {
		return err
	}
	files = append([]msix.File{{Name: c.exe + ".exe", Data: data}}, files...)
	for _, logo := range msixLogos {
		png, err := icon.PNG(icon.Resize(img, logo.size))
		if err != nil {
			return err
		}
		files = append(files, msix.File{Name: "Assets/" + logo.name + ".png", Data: png})
	}
	for _, f := range c.files {
		for _, existing := range files {
			if strings.EqualFold(existing.Name, f.name) {
				return fmt.Errorf("--data: %s is already in the package", f.name)
			}
		}
		files = append(files, msix.File{Name: f.name, Data: f.data})
	}
	pkg, err := msix.Write(msixManifest(c, arch, publisher), files, signer)
	if err != nil {
		return err
	}
	out := filepath.Join(c.out, fmt.Sprintf("%s-%s-%s.msix", c.exe, c.version, arch))
	if err := os.WriteFile(out, pkg, 0o644); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, out)
	return nil
}

func msixManifest(c *config, arch, publisher string) []byte {
	esc := func(s string) string {
		var b bytes.Buffer
		xml.EscapeText(&b, []byte(s))
		return b.String()
	}
	var capabilities, devices string
	for _, p := range c.permissions {
		if strings.HasPrefix(msixCapabilities[p], "<Capability") {
			capabilities += "\n    " + msixCapabilities[p]
		} else if msixCapabilities[p] != "" {
			devices += "\n    " + msixCapabilities[p]
		}
	}
	return fmt.Appendf(nil, `<?xml version="1.0" encoding="utf-8"?>
<Package xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10" xmlns:uap="http://schemas.microsoft.com/appx/manifest/uap/windows10" xmlns:rescap="http://schemas.microsoft.com/appx/manifest/foundation/windows10/restrictedcapabilities" IgnorableNamespaces="uap rescap">
  <Identity Name="%s" Publisher="%s" Version="%s" ProcessorArchitecture="%s"/>
  <Properties>
    <DisplayName>%s</DisplayName>
    <PublisherDisplayName>%s</PublisherDisplayName>
    <Logo>Assets\StoreLogo.png</Logo>
  </Properties>
  <Resources>
    <Resource Language="en-us"/>
  </Resources>
  <Dependencies>
    <TargetDeviceFamily Name="Windows.Desktop" MinVersion="10.0.17763.0" MaxVersionTested="10.0.26100.0"/>
  </Dependencies>
  <Capabilities>%s
    <rescap:Capability Name="runFullTrust"/>%s
  </Capabilities>
  <Applications>
    <Application Id="App" Executable="%s.exe" EntryPoint="Windows.FullTrustApplication">
      <uap:VisualElements DisplayName="%s" Description="%s" BackgroundColor="transparent" Square150x150Logo="Assets\Square150x150Logo.png" Square44x44Logo="Assets\Square44x44Logo.png"/>
    </Application>
  </Applications>
</Package>
`, esc(c.id), esc(publisher), windowsVersion(c.version, c.build), arch, esc(c.name), esc(c.vendor),
		capabilities, devices, esc(c.exe), esc(c.name), esc(c.name))
}

func windowsSigner(c *config) (*authenticode.Signer, *keys.Identity, error) {
	id, err := c.identity()
	if err != nil {
		return nil, nil, err
	}
	signer := &authenticode.Signer{Key: id.Key, Chain: id.Certificates()}
	if c.timestamp {
		signer.TimestampURL = c.tsaURL
	}
	return signer, id, nil
}

func signMSIX(c *config, path string) error {
	signer, _, err := windowsSigner(c)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	signed, err := msix.Sign(data, signer)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return os.WriteFile(path, signed, 0o644)
}

func signWindows(c *config, exe string) error {
	data, err := os.ReadFile(exe)
	if err != nil {
		return err
	}
	signed, err := signWindowsData(c, data)
	if err != nil {
		return err
	}
	return os.WriteFile(exe, signed, 0o755)
}

func signWindowsData(c *config, data []byte) ([]byte, error) {
	signer, _, err := windowsSigner(c)
	if err != nil {
		return nil, err
	}
	return authenticode.Sign(data, signer)
}

func versionInfo(c *config) version.Info {
	fileVersion := windowsVersion(c.version, c.build)

	vi := version.Info{}
	vi.Set(version.LangDefault, version.ProductName, c.name)
	vi.Set(version.LangDefault, version.FileDescription, c.name)
	vi.Set(version.LangDefault, version.InternalName, c.exe)
	vi.Set(version.LangDefault, version.OriginalFilename, c.exe+".exe")
	vi.Set(version.LangDefault, version.CompanyName, c.vendor)
	if c.copyright != "" {
		vi.Set(version.LangDefault, version.LegalCopyright, c.copyright)
	}
	vi.SetFileVersion(fileVersion)
	vi.SetProductVersion(fileVersion)
	vi.Set(version.LangDefault, version.ProductVersion, c.version)
	return vi
}

func windowsVersion(v string, build int) string {
	parts := make([]string, 0, 4)
	for f := range strings.SplitSeq(v, ".") {
		if len(parts) == 3 {
			break
		}
		n, err := strconv.Atoi(f)
		if err != nil {
			break
		}
		parts = append(parts, strconv.Itoa(n))
	}
	for len(parts) < 3 {
		parts = append(parts, "0")
	}
	return strings.Join(append(parts, strconv.Itoa(build)), ".")
}
