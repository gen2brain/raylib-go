#!/bin/bash

##### Bootstrap Go android/arm and android/arm64 with Android standalone toolchains.

if [[ -z "$1" ]]; then
    echo "Usage: bootstrap.sh <install prefix>"
    exit 1
fi

if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
    OS="windows"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="darwin"
elif [[ "$OSTYPE" == "linux-gnu" ]]; then
    OS="linux"
fi

MYCC=clang
MYCXX=clang++

###################################################

echo "##### Check requirements"
for x in "python" "tar" "unzip" "curl" "cmake" "grep" "sed" "awk"; do
    which ${x} || MISSING=true
done
if [[ "$MISSING" = true ]]; then
    exit 1
fi

API_VERSION_ARM=15
API_VERSION_ARM64=21

BUILD_TOOLS_VERSION=26.0.2

if [[ -z "$GO_VERSION" ]]; then
    # go1.9.2
    GO_VERSION=`curl -s https://golang.org/dl/ | grep 'id="go' | head -n1 | awk -F'"' '{print $4}'`
fi
if [[ -z "$NDK_VERSION" ]]; then
    # r15c
    NDK_VERSION=`curl -s https://developer.android.com/ndk/downloads/index.html | grep 'id="stable-downloads"' | awk -F'(' '{print $2}' | awk -F')' '{print $1}'`
fi
if [[ -z "$TOOLS_VERSION" ]]; then
    # 3859397
    TOOLS_VERSION=`curl -s https://developer.android.com/studio/index.html | grep sdk-tools-linux | grep -o '[0-9]*' | head -n1`
fi

INSTALL_PREFIX="$1"
export PATH=${INSTALL_PREFIX}/android-arm/bin:${INSTALL_PREFIX}/android-arm64/bin:${PATH}

BUILD_DIR=`mktemp -d`
mkdir -p ${BUILD_DIR}/bootstrap || exit 1

function finish {
    echo; echo "##### Remove build directory"
    rm -rf ${BUILD_DIR}
}

trap finish EXIT

###################################################

if [[ -z "$ANDROID_NDK_HOME" ]]; then
    echo; echo "##### Download NDK ${NDK_VERSION}"
    cd ${BUILD_DIR} && curl -L -O --progress-bar https://dl.google.com/android/repository/android-ndk-${NDK_VERSION}-${OS}-x86_64.zip; tar -xzf android-ndk-${NDK_VERSION}-${OS}-x86_64.zip >/dev/null 2>&1 || unzip -q android-ndk-${NDK_VERSION}-${OS}-x86_64.zip || exit -1 
    ANDROID_NDK_HOME=${BUILD_DIR}/android-ndk-${NDK_VERSION}
fi

if [[ "$OS" == "windows" && "$OSTYPE" == "msys" ]]; then
    # Fix for python in MSYS2 shell
    grep -q 'MSYS' ${ANDROID_NDK_HOME}/build/tools/make_standalone_toolchain.py || sed -i "s/platform\.system() == 'Windows'/platform.system() == 'Windows' or platform.system().startswith('MSYS')/" ${ANDROID_NDK_HOME}/build/tools/make_standalone_toolchain.py
    # Export path to mingw64
    export PATH=/mingw64/bin:${PATH}
fi

echo; echo "##### Make standalone Android toolchains"
${ANDROID_NDK_HOME}/build/tools/make_standalone_toolchain.py --api=${API_VERSION_ARM} --arch=arm --install-dir=${INSTALL_PREFIX}/android-arm -v  || exit 1
${ANDROID_NDK_HOME}/build/tools/make_standalone_toolchain.py --api=${API_VERSION_ARM64} --arch=arm64 --install-dir=${INSTALL_PREFIX}/android-arm64 -v || exit 1

###################################################

echo; echo "##### Download Go ${GO_VERSION} binaries"
if [[ "$OS" == "windows" ]]; then
    cd ${BUILD_DIR}/bootstrap && curl -L -O --progress-bar http://storage.googleapis.com/golang/${GO_VERSION}.${OS}-amd64.zip; tar -xzf ${GO_VERSION}.${OS}-amd64.zip >/dev/null 2>&1 || unzip -q ${GO_VERSION}.${OS}-amd64.zip || exit 1
else
    cd ${BUILD_DIR}/bootstrap && curl -L --progress-bar http://storage.googleapis.com/golang/${GO_VERSION}.${OS}-amd64.tar.gz | tar -xz || exit 1
fi

echo; echo "##### Download Go ${GO_VERSION} source"
cd ${BUILD_DIR} && curl -L --progress-bar http://storage.googleapis.com/golang/${GO_VERSION}.src.tar.gz | tar -xz && cd ${BUILD_DIR}/go/src || exit 1

if [[ "$OS" == "windows" && "$OSTYPE" == "msys" ]]; then
    # Fix/Hack for MSYS2 shell
    sed -i '/^eval/c\export $(./cmd/dist/dist env -p | grep -v PATH | xargs)' ${BUILD_DIR}/go/src/make.bash
    # Fix/Hack for MSYS2 shell
    sed -i 's/old_bin_files go gofmt/old_bin_files go.exe gofmt.exe/' ${BUILD_DIR}/go/src/make.bash
fi

echo; echo "##### Compile Go ${GO_VERSION} for host"
GOROOT_BOOTSTRAP=${BUILD_DIR}/bootstrap/go ./make.bash || exit 1

echo; echo "##### Compile Go ${GO_VERSION} for arm-linux-androideabi"
GOROOT_BOOTSTRAP=${BUILD_DIR}/bootstrap/go CC_FOR_TARGET=arm-linux-androideabi-${MYCC} GOOS=android GOARCH=arm CGO_ENABLED=1 ./make.bash --no-clean || exit 1

echo; echo "##### Compile Go ${GO_VERSION} for aarch64-linux-android"
GOROOT_BOOTSTRAP=${BUILD_DIR}/bootstrap/go CC_FOR_TARGET=aarch64-linux-android-${MYCC} GOOS=android GOARCH=arm64 CGO_ENABLED=1 ./make.bash --no-clean || exit 1

cp -r -f ${BUILD_DIR}/go ${INSTALL_PREFIX}

###################################################

echo; echo "##### Compile android_native_app_glue"
mkdir -p ${BUILD_DIR}/native_app_glue/jni
cp -r ${ANDROID_NDK_HOME}/sources/android/native_app_glue/* ${BUILD_DIR}/native_app_glue/jni/
echo "APP_ABI := armeabi-v7a arm64-v8a" > ${BUILD_DIR}/native_app_glue/jni/Application.mk

cd ${BUILD_DIR}/native_app_glue && ${ANDROID_NDK_HOME}/ndk-build V=1 || exit 1

cp ${BUILD_DIR}/native_app_glue/obj/local/armeabi-v7a/libandroid_native_app_glue.a ${INSTALL_PREFIX}/android-arm/lib/
cp ${BUILD_DIR}/native_app_glue/obj/local/arm64-v8a/libandroid_native_app_glue.a ${INSTALL_PREFIX}/android-arm64/lib/

cp ${ANDROID_NDK_HOME}/sources/android/native_app_glue/android_native_app_glue.h ${INSTALL_PREFIX}/android-arm/include/
cp ${ANDROID_NDK_HOME}/sources/android/native_app_glue/android_native_app_glue.h ${INSTALL_PREFIX}/android-arm64/include/

###################################################

echo; echo "##### Download sdk tools"
cd ${BUILD_DIR} && curl -L -O --progress-bar https://dl.google.com/android/repository/sdk-tools-${OS}-${TOOLS_VERSION}.zip; tar -xzf sdk-tools-${OS}-${TOOLS_VERSION}.zip >/dev/null 2>&1 || unzip -q sdk-tools-${OS}-${TOOLS_VERSION}.zip || exit -1 
mkdir -p ${INSTALL_PREFIX}/android-tools
cp -r -f ${BUILD_DIR}/tools/* ${INSTALL_PREFIX}/android-tools/

echo; echo "##### Install build tools and platforms"
yes | ${INSTALL_PREFIX}/android-tools/bin/sdkmanager --licenses
${INSTALL_PREFIX}/android-tools/bin/sdkmanager "build-tools;${BUILD_TOOLS_VERSION}"
${INSTALL_PREFIX}/android-tools/bin/sdkmanager "platforms;android-${API_VERSION_ARM}"
${INSTALL_PREFIX}/android-tools/bin/sdkmanager "platforms;android-${API_VERSION_ARM64}"
