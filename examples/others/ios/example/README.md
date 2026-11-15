### iOS example

Requirements:

* macOS with [Xcode](https://developer.apple.com/xcode/) and Go
* Apple development certificate (`.p12`) and provisioning profile (`.mobileprovision`)
* [zsign](https://github.com/zhlynn/zsign) and [go-ios](https://github.com/danielpaulus/go-ios)

The program registers its `main` function with `rl.SetMain(main)`, see `main.go`.

Build the application, the output is `build/example.app` with the files from `assets` copied into it:

    ./build.sh

Sign it, the bundle identifier (`com.example.raylib` in `Info.plist`) must match your provisioning profile:

    mkdir Payload && cp -R build/example.app Payload/
    zsign -k certificate.p12 -p password -m profile.mobileprovision -b com.example.raylib Payload/example.app
    zip -qry example.ipa Payload

Install it on the device:

    go-ios install --path=example.ipa

OpenGL ES 2.0 is used by default, build with `./build.sh -tags es3` for OpenGL ES 3.0 or `./build.sh -tags software`
for the software renderer.
