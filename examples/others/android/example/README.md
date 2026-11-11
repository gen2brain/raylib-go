### Android example

To compile example to shared library you will need [Android NDK](https://developer.android.com/ndk/downloads/index.html).
To build Android apk you will need [Android SDK](http://developer.android.com/sdk/index.html#Other) and JDK 17 or newer.

Export path to Android NDK, point to location where you have unpacked archive:

    export ANDROID_NDK_HOME=/opt/android-ndk

Add toolchain bin directory to PATH:

    export PATH=${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/linux-x86_64/bin:${PATH}

Export API version (the example targets minimum API 24):

    export ANDROID_API=24

And compile shared library:

    CC="aarch64-linux-android${ANDROID_API}-clang" \
    CGO_ENABLED=1 GOOS=android GOARCH=arm64 \
    go build -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-soname,libexample.so" \
    -o=android/libs/arm64-v8a/libexample.so

To build apk export path to Android SDK, point to location where you unpacked archive:

    export ANDROID_HOME=/opt/android-sdk

And build apk:

    ./gradlew assembleDebug

If everything is successfully built apk can be found in the android/build/outputs.


For 32-bit ARM use `CC="armv7a-linux-androideabi${ANDROID_API}-clang"`, set GOARCH to arm and put the library in `android/libs/armeabi-v7a`.

### Windows 

To build shared libraries on Windows you can inspect the file `androidcompile.bat` by opening it on a text editor, it can then be configured to generate the libraries when running the .bat file

See also [instructions](https://developer.android.com/build/building-cmdline#build_bundle) to build a android app bundle for appstore distribution
