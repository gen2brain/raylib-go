#!/bin/sh

if [ -z "$ANDROID_NDK_HOME" ]; then
    echo "You must define ANDROID_NDK_HOME before starting. It should point to your NDK directories."
    exit 1
fi

if [ -z "$1" ]; then
    echo "Usage: bootstrap.sh <install prefix>"
    exit 1
fi

GO_OS="linux"
GO_ARCH="amd64"
GO_VERSION=`curl -s https://golang.org/dl/ | grep 'id="go' | head -n1 | awk -F'"' '{print $4}'`

OPENAL_VERSION="1.17.2"

INSTALL_PREFIX="$1"
export PATH=${INSTALL_PREFIX}/android-arm7/bin:${INSTALL_PREFIX}/android-arm64/bin:${PATH}

BUILD_DIR=`mktemp -d`
mkdir -p ${BUILD_DIR}/bootstrap

echo "##### Make standalone android toolchains"
${ANDROID_NDK_HOME}/build/tools/make-standalone-toolchain.sh --platform=android-9 --install-dir=${INSTALL_PREFIX}/android-arm7 --toolchain=arm-linux-androideabi-4.9
${ANDROID_NDK_HOME}/build/tools/make-standalone-toolchain.sh --platform=android-21 --install-dir=${INSTALL_PREFIX}/android-arm64 --toolchain=aarch64-linux-android-4.9

echo "##### Download Go binaries"
cd ${BUILD_DIR}/bootstrap && curl -s -L http://storage.googleapis.com/golang/${GO_VERSION}.${GO_OS}-${GO_ARCH}.tar.gz | tar xz

echo "##### Download Go source"
cd ${BUILD_DIR} && curl -s -L http://storage.googleapis.com/golang/${GO_VERSION}.src.tar.gz | tar xz && cd ${BUILD_DIR}/go/src

echo "##### Compile Go for host"
GOROOT_BOOTSTRAP=${BUILD_DIR}/bootstrap/go ./make.bash || exit 1

echo "##### Compile Go for arm-linux-androideabi"
GOROOT_BOOTSTRAP=${BUILD_DIR}/bootstrap/go CC_FOR_TARGET=arm-linux-androideabi-gcc GOOS=android GOARCH=arm CGO_ENABLED=1 ./make.bash --no-clean || exit 1

echo "##### Compile Go for aarch64-linux-android"
GOROOT_BOOTSTRAP=${BUILD_DIR}/bootstrap/go CC_FOR_TARGET=aarch64-linux-android-gcc GOOS=android GOARCH=arm64 CGO_ENABLED=1 ./make.bash --no-clean || exit 1

cp -r -f ${BUILD_DIR}/go ${INSTALL_PREFIX}

echo "##### Compile OpenAL"

cd ${BUILD_DIR} && curl -s -L  http://kcat.strangesoft.net/openal-releases/openal-soft-${OPENAL_VERSION}.tar.bz2 | tar -xj

cat << EOF > ${BUILD_DIR}/openal-soft-${OPENAL_VERSION}/android-arm7.cmake
set(TOOLCHAIN_PREFIX arm-linux-androideabi)
set(CMAKE_C_COMPILER \${TOOLCHAIN_PREFIX}-gcc)
set(CMAKE_FIND_ROOT_PATH \${INSTALL_PREFIX}/android-arm7)
set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
EOF

cat << EOF > ${BUILD_DIR}/openal-soft-${OPENAL_VERSION}/android-arm64.cmake
set(TOOLCHAIN_PREFIX aarch64-linux-android)
set(CMAKE_C_COMPILER \${TOOLCHAIN_PREFIX}-gcc)
set(CMAKE_FIND_ROOT_PATH \${INSTALL_PREFIX}/android-arm64)
set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
EOF

mkdir -p ${BUILD_DIR}/openal-soft-${OPENAL_VERSION}/build-arm7
cd ${BUILD_DIR}/openal-soft-${OPENAL_VERSION}/build-arm7
cmake -DLIBTYPE=STATIC -DCMAKE_TOOLCHAIN_FILE=../android-arm7.cmake -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX}/android-arm7 ..
make -j $(nproc) && make install

mkdir -p ${BUILD_DIR}/openal-soft-${OPENAL_VERSION}/build-arm64
cd ${BUILD_DIR}/openal-soft-${OPENAL_VERSION}/build-arm64
cmake -DLIBTYPE=STATIC -DCMAKE_TOOLCHAIN_FILE=../android-arm64.cmake -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX}/android-arm64 ..
make -j $(nproc) && make install

echo "##### Compile android_native_app_glue"

mkdir -p ${BUILD_DIR}/native_app_glue/jni
cp -r ${ANDROID_NDK_HOME}/sources/android/native_app_glue/* ${BUILD_DIR}/native_app_glue/jni/
echo "APP_ABI := armeabi-v7a arm64-v8a" > ${BUILD_DIR}/native_app_glue/jni/Application.mk

cd ${BUILD_DIR}/native_app_glue && ${ANDROID_NDK_HOME}/ndk-build

cp ${BUILD_DIR}/native_app_glue/obj/local/armeabi-v7a/libandroid_native_app_glue.a ${INSTALL_PREFIX}/android-arm7/lib/
cp ${BUILD_DIR}/native_app_glue/obj/local/arm64-v8a/libandroid_native_app_glue.a ${INSTALL_PREFIX}/android-arm64/lib/

cp ${ANDROID_NDK_HOME}/sources/android/native_app_glue/android_native_app_glue.h ${INSTALL_PREFIX}/android-arm7/include/
cp ${ANDROID_NDK_HOME}/sources/android/native_app_glue/android_native_app_glue.h ${INSTALL_PREFIX}/android-arm64/include/

echo "##### Remove build directory"
rm -rf ${BUILD_DIR}
