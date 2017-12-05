### Android example

To compile example to shared library you will need [Android NDK](https://developer.android.com/ndk/downloads/index.html).
To build Android apk you will need [Android SDK](http://developer.android.com/sdk/index.html#Other).
Download and unpack archives somewhere.

Go must be cross compiled for android. There is a bootstrap.sh script that you can use to compile Go for android/arm and android/arm64.

Export path to Android NDK, point to location where you have unpacked archive:

    export ANDROID_NDK_HOME=/opt/android-ndk

Compile Go and android_native_app_glue, /usr/local is prefix where Go and Android toolchains will be installed:

    ./bootstrap.sh /usr/local

After build is complete point GOROOT to new Go installation in /usr/local, and add toolchain bin directory to PATH:

    export GOROOT=/usr/local/go
    export PATH=/usr/local/android-arm7/bin:${PATH}

And compile shared library:

    CGO_CFLAGS="-I/usr/local/android-arm7/include" CGO_LDFLAGS="-L/usr/local/android-arm7/lib" \
    CC=arm-linux-androideabi-gcc CGO_ENABLED=1 GOOS=android GOARCH=arm \
    ${GOROOT}/bin/go build -v -x -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-soname,libexample.so" \
    -o=android/libs/armeabi-v7a/libexample.so

To build apk export path to Android SDK, point to location where you unpacked archive:

    export ANDROID_HOME=/opt/android-sdk

And build apk with ant:

    cd android
    ant clean debug

If everything is successfully built apk can be found in bin/ directory.
