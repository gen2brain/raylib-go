#!/usr/bin/env bash

CHROOT="/usr/x86_64-pc-linux-gnu-static"
MINGW="/usr/i686-w64-mingw32"
RPI="/usr/armv6j-hardfloat-linux-gnueabi"
ANDROID="/opt/android-toolchain-arm7"

CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=c-archive -o release/linux/rreslib.a -ldflags "-s -w"

CC="i686-w64-mingw32-gcc" CXX="i686-w64-mingw32-g++" \
CGO_ENABLED=1 GOOS=windows GOARCH=386 go build -buildmode=c-archive -o release/win32/rreslib.a -ldflags "-s -w"

CC="armv6j-hardfloat-linux-gnueabi-gcc" CXX="armv6j-hardfloat-linux-gnueabi-g++" \
CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -buildmode=c-archive -o release/rpi/rreslib.a -ldflags "-s -w"

PATH="$PATH:$ANDROID/bin" \
CC="arm-linux-androideabi-gcc" CXX="arm-linux-androideabi-g++" \
CGO_ENABLED=1 GOOS=android GOARCH=arm go build -buildmode=c-shared -o release/android/armeabi-v7a/rreslib.so -ldflags "-s -w"
