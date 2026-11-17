package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/apple"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/icon"
)

var iosIcons = []struct {
	name string
	size int
}{
	{"Icon-60@2x.png", 120}, {"Icon-60@3x.png", 180},
	{"Icon-76@2x.png", 152}, {"Icon-83.5@2x.png", 167},
	{"Icon-1024.png", 1024},
}

type iosTarget struct {
	env     []string
	cc      string
	cflags  []string
	ldflags []string
}

func packageIOS(c *config) error {
	arch := c.goarch
	if arch != "arm64" && !(c.simulator && arch == "amd64") {
		return fmt.Errorf("unsupported ios arch %q", arch)
	}

	target, err := iosToolchain(c, arch)
	if err != nil {
		return err
	}

	tmp, err := os.MkdirTemp("", "rlpkg-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	app := filepath.Join(tmp, "Payload", c.exe+".app")
	if err := os.MkdirAll(app, 0o755); err != nil {
		return err
	}
	archive := filepath.Join(tmp, "lib"+c.exe+".a")
	if err := goBuildWith(c, target.env, archive, nil, target.ldflags, []string{"-buildmode=c-archive"}); err != nil {
		return err
	}
	ldflags, err := cgoLDFlags(c, target.env)
	if err != nil {
		return err
	}
	link := append(slices.Clone(target.cflags), "-o", filepath.Join(app, c.exe), archive)
	if err := run(target.cc, append(link, ldflags...)...); err != nil {
		return err
	}

	img, err := c.loadIcon()
	if err != nil {
		return err
	}
	for _, ic := range iosIcons {
		data, err := icon.PNG(icon.Resize(img, ic.size))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(app, ic.name), data, 0o644); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(app, "Info.plist"), iosInfoPlist(c), 0o644); err != nil {
		return err
	}
	if err := c.writeData(app, "_CodeSignature", "embedded.mobileprovision", "Frameworks", "PlugIns"); err != nil {
		return err
	}

	if !c.simulator {
		if err := signIOS(c, app, tmp); err != nil {
			return err
		}
	}

	suffix := ".ipa"
	if c.simulator {
		suffix = "-simulator.zip"
	}
	out := filepath.Join(c.out, c.exe+suffix)
	if err := zipDir(out, filepath.Join(tmp, "Payload")); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, out)

	if c.install && !c.simulator {
		if runtime.GOOS == "darwin" {
			return run("xcrun", "devicectl", "device", "install", "app", out)
		}
		return run("go-ios", "install", "--path="+out)
	}
	return nil
}

func iosToolchain(c *config, arch string) (*iosTarget, error) {
	env := []string{"GOOS=ios", "GOARCH=" + arch, "CGO_ENABLED=1", "IOS_DEPLOYMENT_TARGET=" + c.minOS}
	if runtime.GOOS == "darwin" {
		sdk := "iphoneos"
		triple := "arm64-apple-ios" + c.minOS
		if c.simulator {
			sdk = "iphonesimulator"
			triple = map[string]string{"arm64": "arm64", "amd64": "x86_64"}[arch] + "-apple-ios" + c.minOS + "-simulator"
		}
		sdkPath, err := exec.Command("xcrun", "-sdk", sdk, "--show-sdk-path").Output()
		if err != nil {
			return nil, fmt.Errorf("xcrun: %w", err)
		}
		cc, err := exec.Command("xcrun", "-sdk", sdk, "-find", "clang").Output()
		if err != nil {
			return nil, fmt.Errorf("xcrun: %w", err)
		}
		clang := strings.TrimSpace(string(cc))
		cflags := []string{"-isysroot", strings.TrimSpace(string(sdkPath)), "-target", triple}
		env = append(append(env, "CC="+clang, "CXX="+clang+"++"), cgoFlags(cflags...)...)
		return &iosTarget{env, clang, cflags, nil}, nil
	}

	cc := "arm64-apple-ios-clang"
	if c.simulator {
		cc = map[string]string{"arm64": "arm64", "amd64": "x86_64"}[arch] + "-apple-ios-simulator-clang"
	}
	if v := os.Getenv("CC"); v != "" {
		cc = v
	}
	ccPath, err := exec.LookPath(cc)
	if err != nil {
		return nil, fmt.Errorf("%s not found; install Xcode on macOS or set CC", cc)
	}
	ar := os.Getenv("AR")
	if ar == "" {
		ars, _ := filepath.Glob(filepath.Join(filepath.Dir(ccPath), "arm64-apple-darwin*-ar"))
		if len(ars) == 0 {
			return nil, fmt.Errorf("no Apple ar next to %s; set AR", ccPath)
		}
		ar = ars[0]
	}
	return &iosTarget{append(env, "CC="+cc, "CXX="+cc+"++"), cc, nil, []string{"-extar=" + ar}}, nil
}

func cgoLDFlags(c *config, env []string) ([]string, error) {
	args := []string{"list", "-deps", "-f", `{{join .CgoLDFLAGS "\n"}}`}
	if len(c.tags) > 0 {
		args = append(args, "-tags", strings.Join(c.tags, ","))
	}
	cmd := exec.Command("go", append(args, c.pkg)...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var flags, seen []string
	fields := strings.Fields(string(out))
	for i := 0; i < len(fields); i++ {
		flag := fields[i]
		if flag == "-framework" && i+1 < len(fields) {
			i++
			flag += " " + fields[i]
		}
		if slices.Contains(seen, flag) {
			continue
		}
		seen = append(seen, flag)
		flags = append(flags, strings.Fields(flag)...)
	}
	return flags, nil
}

func iosInfoPlist(c *config) []byte {
	var b bytes.Buffer
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
`)
	str := func(key, value string) {
		b.WriteString("\t<key>" + key + "</key>\n\t<string>")
		xml.EscapeText(&b, []byte(value))
		b.WriteString("</string>\n")
	}
	array := func(key string, values ...string) {
		b.WriteString("\t<key>" + key + "</key>\n\t<array>\n")
		for _, v := range values {
			b.WriteString("\t\t<string>" + v + "</string>\n")
		}
		b.WriteString("\t</array>\n")
	}
	platform := "iPhoneOS"
	if c.simulator {
		platform = "iPhoneSimulator"
	}
	str("CFBundleDevelopmentRegion", "en")
	str("CFBundleDisplayName", c.name)
	str("CFBundleExecutable", c.exe)
	str("CFBundleIdentifier", c.id)
	str("CFBundleInfoDictionaryVersion", "6.0")
	str("CFBundleName", c.name)
	str("CFBundlePackageType", "APPL")
	str("CFBundleShortVersionString", c.version)
	str("CFBundleVersion", strconv.Itoa(c.build))
	array("CFBundleSupportedPlatforms", platform)
	b.WriteString("\t<key>LSRequiresIPhoneOS</key>\n\t<true/>\n")
	str("MinimumOSVersion", c.minOS)
	b.WriteString(`	<key>UIDeviceFamily</key>
	<array>
		<integer>1</integer>
		<integer>2</integer>
	</array>
	<key>CFBundleIcons</key>
	<dict>
		<key>CFBundlePrimaryIcon</key>
		<dict>
			<key>CFBundleIconFiles</key>
			<array>
				<string>Icon-60</string>
			</array>
		</dict>
	</dict>
	<key>CFBundleIcons~ipad</key>
	<dict>
		<key>CFBundlePrimaryIcon</key>
		<dict>
			<key>CFBundleIconFiles</key>
			<array>
				<string>Icon-60</string>
				<string>Icon-76</string>
				<string>Icon-83.5</string>
			</array>
		</dict>
	</dict>
	<key>UIApplicationSceneManifest</key>
	<dict>
		<key>UIApplicationSupportsMultipleScenes</key>
		<false/>
	</dict>
	<key>UILaunchScreen</key>
	<dict/>
	<key>UIRequiredDeviceCapabilities</key>
	<array>
		<string>arm64</string>
	</array>
`)
	landscape := []string{"UIInterfaceOrientationLandscapeLeft", "UIInterfaceOrientationLandscapeRight"}
	switch c.orientation {
	case "landscape":
		array("UISupportedInterfaceOrientations", landscape...)
		array("UISupportedInterfaceOrientations~ipad", landscape...)
	case "portrait":
		array("UISupportedInterfaceOrientations", "UIInterfaceOrientationPortrait")
		array("UISupportedInterfaceOrientations~ipad", "UIInterfaceOrientationPortrait", "UIInterfaceOrientationPortraitUpsideDown")
	default:
		array("UISupportedInterfaceOrientations", append([]string{"UIInterfaceOrientationPortrait"}, landscape...)...)
		array("UISupportedInterfaceOrientations~ipad", append([]string{"UIInterfaceOrientationPortrait", "UIInterfaceOrientationPortraitUpsideDown"}, landscape...)...)
	}
	if c.hasPermission("camera") {
		str("NSCameraUsageDescription", c.name+" uses the camera.")
	}
	if c.hasPermission("microphone") {
		str("NSMicrophoneUsageDescription", c.name+" uses the microphone.")
	}
	if c.hasPermission("location") {
		str("NSLocationWhenInUseUsageDescription", c.name+" uses your location.")
	}
	b.WriteString("</dict>\n</plist>\n")
	return b.Bytes()
}

var (
	profileEntitlements = regexp.MustCompile(`(?s)<key>Entitlements</key>\s*(<dict>.*?</dict>)`)
	profileAppID        = regexp.MustCompile(`(<key>application-identifier</key>\s*<string>[A-Z0-9]+\.)\*(</string>)`)
)

func profileEntitlementsPlist(profile []byte, bundleID string) ([]byte, error) {
	m := profileEntitlements.FindSubmatch(profile)
	if m == nil {
		return nil, errors.New("provisioning profile has no Entitlements")
	}
	dict := profileAppID.ReplaceAll(m[1], []byte("${1}"+bundleID+"${2}"))
	return fmt.Appendf(nil, `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
%s
</plist>
`, dict), nil
}

func embedProfile(path, app, bundleID string) ([]byte, error) {
	profile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(app, "embedded.mobileprovision"), profile, 0o644); err != nil {
		return nil, err
	}
	entitlements, err := profileEntitlementsPlist(profile, bundleID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return entitlements, nil
}

func signIOS(c *config, app, tmp string) error {
	if c.sign == "" || c.profile == "" {
		fmt.Fprintln(os.Stderr, "rlpkg: --sign and --profile not given, the app is not signed")
		return nil
	}
	entitlements, err := embedProfile(c.profile, app, c.id)
	if err != nil {
		return err
	}
	entFile := filepath.Join(tmp, "entitlements.plist")
	if err := os.WriteFile(entFile, entitlements, 0o644); err != nil {
		return err
	}

	if c.signer == "codesign" {
		return run("codesign", "--force", "--sign", c.sign, "--timestamp=none", "--entitlements", entFile, app)
	}
	signer, err := appleSigner(c)
	if err != nil {
		return err
	}
	_, err = apple.SignBundle(app, apple.BundleOptions{Signer: signer, Entitlements: entitlements})
	return err
}
