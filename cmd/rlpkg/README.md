## rlpkg

Builds and packages a raylib-go program for distribution: a Windows executable with icon, version info and manifest, or an `.msix`,
a macOS `.app`, a Linux tarball, `.deb` or `.rpm`, an Android `.apk` or `.aab`, an iOS `.ipa` or a WebAssembly site.

Every format is written by rlpkg itself, including the signatures, so a Go toolchain is all a host needs.
Cross-compiling C code and installing on devices need the platform tools named below.

```sh
go install github.com/gen2brain/raylib-go/cmd/rlpkg@latest
```

### Quick start

Run it from the module of the program you package:

```sh
rlpkg package --os linux ./game                       # game-1.0.0-linux-amd64.tar.gz
rlpkg package --os windows --arch amd64 ./game        # game.exe
rlpkg package --os darwin --arch universal ./game     # game.app + zip
rlpkg package --os android --install ./game           # game.apk, installed with adb
rlpkg package --os js ./game && rlpkg serve game      # http://localhost:8000/
```

### Usage

```
rlpkg package [flags] [package]
rlpkg sign [flags] <file>
rlpkg serve [--addr localhost:8000] <directory>
rlpkg staple <app>
```

`package` is the main package to build, `.` by default. Output goes to `--out` (default: the current directory). `rlpkg help` prints the flag list.

 | Flag                                                 | Default             | Meaning                                                                                                        |
 |------------------------------------------------------|---------------------|----------------------------------------------------------------------------------------------------------------|
 | `--os`                                               | host                | `windows`, `darwin`, `linux`, `android`, `ios`, `js` (or `wasm`)                                               |
 | `--arch`                                             | host                | target architecture; `universal` on darwin, a comma list on android, `arm64` on android and ios by default     |
 | `--out`                                              | `.`                 | output directory                                                                                               |
 | `--name`                                             | executable name     | application name shown to the user                                                                             |
 | `--id`                                               | `com.example.<exe>` | application identifier; linux: the desktop file name                                                           |
 | `--version`                                          | `1.0.0`             | version string                                                                                                 |
 | `--build`                                            | `1`                 | build number: `versionCode`, `CFBundleVersion`, last Windows version field, deb/rpm revision                   |
 | `--icon`                                             | the raylib logo     | square PNG, 1024 px recommended                                                                                |
 | `--data`                                             |                     | file or directory shipped with the application, repeatable; a directory keeps its name                         |
 | `--tags`                                             |                     | Go build tags (`es3,opengl43,...`)                                                                             |
 | `--ldflags`                                          |                     | extra linker flags                                                                                             |
 | `--cgo`                                              |                     | build with cgo (default: `go env CGO_ENABLED` for the host platform, off when cross-compiling)                 |
 | `--debug`                                            |                     | keep symbols, debug info and file paths; otherwise `-trimpath -ldflags "-s -w"`                                |
 | `--console`                                          |                     | windows: console application                                                                                   |
 | `--orientation`                                      | see Android, iOS    | android, ios: `landscape`, `portrait`, `any`                                                                   |
 | `--format`                                           | linux: `targz`      | linux: comma list of `targz`, `deb`, `rpm`; windows: `msix`, written next to the `.exe`; android: `apk`, `aab` |
 | `--category`                                         | `Game`              | linux: desktop entry categories                                                                                |
 | `--vendor`                                           | application name    | Windows company name, Debian maintainer (`Name <email>`), RPM vendor                                           |
 | `--publisher`                                        | certificate subject | windows msix: publisher of the package identity                                                                |
 | `--copyright`                                        |                     | Windows copyright line                                                                                         |
 | `--license`                                          | `Unknown`           | rpm license name                                                                                               |
 | `--permissions`                                      |                     | comma list of `camera`, `microphone`, `location`, `notifications`, `internet`, `network`, `vibrate`            |
 | `--sign`                                             |                     | signing identity, see Signing                                                                                  |
 | `--signer`                                           | `rlpkg`             | `rlpkg`; darwin, ios: `codesign` (the default on a Mac); linux: `gpg`                                          |
 | `--timestamp`                                        | `true`              | add a trusted timestamp when signing with a certificate                                                        |
 | `--timestamp-url`                                    | DigiCert            | windows: RFC 3161 server                                                                                       |
 | `--install`                                          |                     | android, ios: install on the connected device                                                                  |
 | `--simulator`                                        |                     | ios: build for the simulator                                                                                   |
 | `--minos`                                            | `15.0`              | ios: minimum iOS version                                                                                       |
 | `--profile`                                          |                     | ios: provisioning profile                                                                                      |
 | `--notary-key`, `--notary-key-id`, `--notary-issuer` |                     | darwin: App Store Connect API key, enables notarization                                                        |

Without `--cgo` the program is built with `CGO_ENABLED=0` and the purego backend, so the desktop targets build from any host with only Go installed.
With `--cgo` raylib is compiled into the executable; cross builds then need the matching C toolchain in `CC`.
Android and iOS are always cgo builds.

`--data` files are loaded with relative paths, as on the desktop: they go to `Contents/Resources` (darwin, the working directory of a bundle), the `.app` root (ios),
`assets` (android), `resources.tar.gz` (js, unpacked before the program starts), `share/<exe>` (targz, installed by its Makefile), `/usr/share/<exe>` (deb, rpm) and the package root (msix).
The Linux desktop entry starts the program in that directory. On Windows the executable and the files are also written to `<exe>-<version>-windows-<arch>.zip`.

### Targets

#### Windows

`<exe>.exe` with the icon, version resource (product name, version, company from `--vendor`, copyright) and an application manifest embedded.
Built as a GUI application unless `--console`. A cross build with `--cgo` needs llvm-mingw or MinGW-w64 in `CC`.

```sh
rlpkg package --os windows --arch amd64 --name "My Game" --version 1.2.0 --icon icon.png ./cmd/game
```

`--format msix` also writes `<exe>-<version>-<arch>.msix`, a full trust package with the identity `--id` and the version `<version>.<build>`.
Its publisher is the subject of the `--sign` certificate, or `--publisher` for an unsigned package; Windows installs only a signed one.
A purego build carries `raylib.dll` next to the executable.

```sh
rlpkg package --os windows --arch amd64 --format msix --id com.example.game --sign codesign.p12 ./cmd/game
```

#### macOS

`<Name>.app` and a zip of it. `--arch universal` builds amd64 and arm64 and merges them. A purego build carries `libraylib.dylib` in `Contents/Frameworks`; a cgo build (the default on a Mac) has raylib linked in.
`LSMinimumSystemVersion` is the minimum macOS of the Go toolchain. Camera, microphone and location usage descriptions and entitlements are added for the listed `--permissions`. Without `--sign` the bundle is signed ad-hoc.

```sh
rlpkg package --os darwin --arch universal --name "My Game" --id com.example.game ./cmd/game
```

#### Linux

`<exe>-<version>-linux-<arch>.tar.gz` with `bin/<exe>`, a `.desktop` file, hicolor icons from 16 to 512 px and a `Makefile` (`install`, `uninstall`, `user-install`, `user-uninstall`, `PREFIX`, `DESTDIR`).

```sh
rlpkg package --os linux --arch arm64 ./cmd/game
tar xzf game-1.0.0-linux-arm64.tar.gz && make -C game-1.0.0 user-install
```

`--format deb,rpm` writes `<name>_<version>-<build>_<arch>.deb` and
`<name>-<version>-<build>.<arch>.rpm` instead, installing under `/usr`. With `--sign` all three are signed with an OpenPGP key, see Signing.

```sh
rlpkg package --os linux --format deb,rpm --vendor "Jane Doe <jane@example.com>" --license MIT ./cmd/game
```

#### Android

`<exe>.apk`, a `NativeActivity` application without Java code. Only the Android NDK is needed: set `ANDROID_NDK_HOME`, or install it under `$ANDROID_SDK_ROOT/ndk`.
The program registers its `main` function with `rl.SetMain(main)` in `init`. It is compiled as a shared library per ABI; the manifest, launcher icons and native libraries are written into the template, then the APK is aligned and signed (v1 and v2).
OpenGL ES 3.0 is required with the `es3` tag, else 2.0. The orientation is `landscape` unless set (either landscape side, `any` follows the sensor).
No permission is requested unless listed. `--install` runs `adb install`.
`--format aab` writes `<exe>.aab` instead, an Android App Bundle for Google Play, signed with the upload key (v1); `--format apk,aab` writes both.

```sh
rlpkg package --os android --arch arm64,amd64 --name "My Game" --id com.example.game --install ./cmd/game
rlpkg package --os android --arch arm64,arm --format aab --sign upload.p12 ./cmd/game
```

#### iOS

`<exe>.ipa` (or `<exe>-simulator.zip` with `--simulator`) with the icons, `Info.plist` and the provisioning profile embedded.
Needs macOS with Xcode. The program registers its `main` function with `rl.SetMain(main)` in `init`; it is built as a C archive and linked with the frameworks its packages use.
The orientation is any unless set. `--sign` and `--profile` are both required for a device build; the entitlements come from the profile.
`--install` installs it on the connected device with `xcrun devicectl`.

```sh
rlpkg package --os ios --name "My Game" --id com.example.game --sign dev.p12 --profile dev.mobileprovision --install ./cmd/game
```

#### WebAssembly

A directory `<exe>/` ready to serve: `index.html`, the raylib module (`raylib.js`, `raylib.wasm`, WebGL 2 with the `es3` tag), the loader `raylib-go.js`, the program as `app.wasm` and Go's `wasm_exec.js`.
`rlpkg serve` is a static server for it; any web server works.

```sh
rlpkg package --os js --name "My Game" ./cmd/game
rlpkg serve game
```

### Signing

`--sign` takes a `.p12`/`.pfx` file (password in `RLPKG_P12_PASSWORD`, empty by default) or a PEM file holding the private key, the certificate and any intermediate certificates.
With `--signer codesign` it is a keychain identity instead. Signatures carry a trusted timestamp unless `--timestamp=false`.

#### Windows

Authenticode: SHA-256, the full certificate chain from the file, an RFC 3161 countersignature from `--timestamp-url`.
Windows only trusts certificates from a CA in Microsoft's root program; a self-signed one produces a valid signature that SmartScreen still warns about.

An `.msix` and the executable in it get the same signature; a self-signed certificate has to be in the machine's Trusted People store.

```sh
RLPKG_P12_PASSWORD=secret rlpkg package --os windows --arch amd64 --sign codesign.p12 ./cmd/game
```

#### macOS

The certificate must be Apple-issued (Developer ID Application for distribution outside the App Store; Apple Development for testing).
The bundle gets the hardened runtime, entitlements from `--permissions`, sealed resources and nested dylibs, and a timestamp from Apple's service.
A purego bundle also gets `com.apple.security.cs.disable-library-validation`, the purego backend loads libffi from the user cache directory.

```sh
rlpkg package --os darwin --arch universal --sign devid.p12 ./cmd/game
```

Notarization needs a Developer ID certificate and an App Store Connect API key (`AuthKey_<id>.p8` with its issuer id).
rlpkg uploads the zipped bundle, waits for Apple's verdict (minutes to half an hour; a rejection prints Apple's log) and staples the ticket.
`rlpkg staple <app>` fetches and staples the ticket for a bundle notarized earlier.

```sh
rlpkg package --os darwin --arch universal --sign devid.p12 --notary-key AuthKey_ABC123DEFG.p8 --notary-issuer 69a6de7f-... ./cmd/game
```

On a Mac `--signer codesign` is the default and uses `codesign`, `notarytool` and `stapler` with a keychain identity in `--sign`;
`--signer rlpkg` selects the built-in signer there too. Apple's `security import` rejects `.p12` files written with current OpenSSL defaults; re-export them with `openssl pkcs12 -export -legacy` for the keychain.

#### iOS

The certificate must match the provisioning profile. On a Mac `codesign` is the default signer.

```sh
rlpkg package --os ios --sign dev.p12 --profile dev.mobileprovision --install ./cmd/game
```

#### Linux

OpenPGP, SHA-256: a detached `<archive>.asc` for the tarball, the header and legacy signatures in an `.rpm`, a debsigs origin signature in a `.deb` (dpkg and apt do not check it).
`--sign` takes a key exported with `gpg --export-secret-keys` (passphrase in `RLPKG_GPG_PASSPHRASE`, empty by default); with `--signer gpg` it names a key of the `gpg` program instead.

```sh
gpg --armor --export-secret-keys jane@example.com > release.asc
RLPKG_GPG_PASSPHRASE=secret rlpkg package --os linux --format targz,deb,rpm --sign release.asc ./cmd/game
rlpkg package --os linux --format rpm --signer gpg --sign jane@example.com ./cmd/game
```

#### Android

Without `--sign` a debug key is generated once and kept in the user cache directory. For release, convert a Java keystore once:

```sh
keytool -importkeystore -srckeystore release.jks -destkeystore release.p12 -deststoretype PKCS12
RLPKG_P12_PASSWORD=secret rlpkg package --os android --sign release.p12 ./cmd/game
```

### Signing existing files

`rlpkg sign` signs a file made elsewhere, in place, with the same signers: a Windows `.exe`, `.dll` or `.msix`, a macOS or iOS `.app` directory, an `.ipa`, an `.apk` or `.aab`, a `.deb`, an `.rpm`, or a bare Mach-O executable or dylib.
It takes `--sign`, `--signer` (`gpg`), `--profile` (iOS), `--entitlements` (macOS), `--timestamp`, `--timestamp-url` and the `--notary-*` flags; a macOS bundle with `--notary-key` is notarized and stapled after signing.

```sh
rlpkg sign --sign codesign.p12 game.exe
rlpkg sign --sign devid.p12 --notary-key AuthKey_ABC123DEFG.p8 --notary-issuer 69a6de7f-... "My Game.app"
rlpkg sign --sign dev.p12 --profile dev.mobileprovision game.ipa
rlpkg sign game.apk                                   # debug key
```

With an OpenPGP key (`--signer gpg`, or an `.asc`, `.gpg` or `.pgp` key file) a `.deb` or `.rpm` is re-signed in place, RPM v6 excepted, and any other file gets a detached `<file>.asc`.
