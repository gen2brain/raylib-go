#!/bin/sh
# Builds build/example.app for iOS devices, requires macOS with Xcode.
# Arguments are passed to go build, for example: ./build.sh -tags es3
set -e

export SDKROOT="$(xcrun --sdk iphoneos --show-sdk-path)"
export CGO_CFLAGS="-miphoneos-version-min=15.0"
export CGO_LDFLAGS="-miphoneos-version-min=15.0"

mkdir -p build/example.app

CGO_ENABLED=1 GOOS=ios GOARCH=arm64 CC=clang \
	go build -buildmode=c-archive -ldflags="-s -w" -o build/libexample.a "$@"

clang -arch arm64 -miphoneos-version-min=15.0 -o build/example.app/example build/libexample.a \
	-framework UIKit -framework Foundation -framework CoreFoundation -framework CoreGraphics -framework QuartzCore \
	-framework OpenGLES -framework AVFoundation -framework AudioToolbox -framework CoreAudio \
	-framework GameController -framework CoreHaptics -framework Security -lresolv

cp Info.plist build/example.app/
cp -R assets/. build/example.app/
