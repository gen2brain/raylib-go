### iOS example

Requirements:

* macOS with [Xcode](https://developer.apple.com/xcode/) and Go
* Apple development certificate and provisioning profile

The program registers its `main` function with `rl.SetMain(main)`, see `main.go`.

Build, sign and install it on the connected device with [rlpkg](https://github.com/gen2brain/raylib-go/tree/master/cmd/rlpkg),
the bundle identifier must match your provisioning profile:

    rlpkg package --os ios --id com.example.raylib --sign certificate.p12 --profile profile.mobileprovision --install \
        --data assets/raylib_logo.png --data assets/coin.wav --data assets/ambient.ogg .

Or build the application with `./build.sh`, the output is `build/example.app` with the files from `assets` copied into it,
then sign it and install it:

    rlpkg sign --sign certificate.p12 --profile profile.mobileprovision build/example.app
    xcrun devicectl device install app --device <device> build/example.app

OpenGL ES 2.0 is used by default, build with `./build.sh -tags es3` (`--tags es3`) for OpenGL ES 3.0 or `-tags software`
for the software renderer.
