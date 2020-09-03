### Android example

To compile example to shared library you will need [Android NDK](https://developer.android.com/ndk/downloads/index.html).
To build Android apk you will need [Android SDK](http://developer.android.com/sdk/index.html#Other).

Export path to Android NDK, point to location where you have unpacked archive:

    export ANDROID_NDK_HOME=/opt/android-ndk

Add toolchain bin directory to PATH:

    export PATH=${ANDROID_NDK_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/bin:${PATH}

Export sysroot path:
    
    export ANDROID_SYSROOT=${ANDROID_NDK_HOME}/platforms/android-16/arch-arm

And compile shared library:

    CC=arm-linux-androideabi-gcc \
    CGO_CFLAGS="-I${ANDROID_SYSROOT}/usr/include --sysroot=${ANDROID_SYSROOT} -D__ANDROID_API__=16" \
    CGO_LDFLAGS="-L${ANDROID_SYSROOT}/usr/lib --sysroot=${ANDROID_SYSROOT}" \
    CGO_ENABLED=1 GOOS=android GOARCH=arm \
    go build -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-soname,libexample.so" \
    -o=android/libs/armeabi-v7a/libexample.so

To build apk export path to Android SDK, point to location where you unpacked archive:

    export ANDROID_HOME=/opt/android-sdk

And build apk with ant:

    cd android
    ant clean debug

Or with gradle:

    ./gradlew assembleDebug

If everything is successfully built apk can be found in bin/ directory or in the android/build/outputs in case `gradle` is used.
