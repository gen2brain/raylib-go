#!/bin/sh

if [ -z "$1" ]; then
    echo "Usage: bootstrap.sh <install prefix>"
    exit 1
fi

GO_OS="linux"
GO_ARCH="amd64"
GO_VERSION=`curl -s https://golang.org/dl/ | grep 'id="go' | head -n1 | awk -F'"' '{print $4}'`

OPENAL_VERSION="1.17.2"

INSTALL_PREFIX="$1"
export PATH=${INSTALL_PREFIX}/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/bin:${PATH}

BUILD_DIR=`mktemp -d`
mkdir -p ${BUILD_DIR}/bootstrap

echo "##### Download toolchain"
cd ${BUILD_DIR} && git clone --depth=1 --branch=master https://github.com/raspberrypi/tools.git
cp -r -f ${BUILD_DIR}/tools/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64 ${INSTALL_PREFIX}

echo "##### Download userspace libraries"
cd ${BUILD_DIR} && git clone --depth=1 --branch=master https://github.com/raspberrypi/firmware.git
cp -r -f ${BUILD_DIR}/firmware/hardfp/opt/vc ${INSTALL_PREFIX}

echo "##### Download Go binaries"
cd ${BUILD_DIR}/bootstrap && curl -s -L http://storage.googleapis.com/golang/${GO_VERSION}.${GO_OS}-${GO_ARCH}.tar.gz | tar xz

echo "##### Download Go source"
cd ${BUILD_DIR} && curl -s -L http://storage.googleapis.com/golang/${GO_VERSION}.src.tar.gz | tar xz && cd ${BUILD_DIR}/go/src

echo "##### Compile Go for host"
GOROOT_BOOTSTRAP=${BUILD_DIR}/bootstrap/go ./make.bash || exit 1

echo "##### Compile Go for arm-linux-gnueabihf"
GOROOT_BOOTSTRAP=${BUILD_DIR}/bootstrap/go CC_FOR_TARGET=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 ./make.bash --no-clean || exit 1

cp -r -f ${BUILD_DIR}/go ${INSTALL_PREFIX}

echo "##### Compile OpenAL"

cd ${BUILD_DIR} && curl -s -L  http://kcat.strangesoft.net/openal-releases/openal-soft-${OPENAL_VERSION}.tar.bz2 | tar -xj

cat << EOF > ${BUILD_DIR}/openal-soft-${OPENAL_VERSION}/linux-rpi.cmake
set(TOOLCHAIN_PREFIX arm-linux-gnueabihf)
set(CMAKE_C_COMPILER \${TOOLCHAIN_PREFIX}-gcc)
set(CMAKE_FIND_ROOT_PATH \${INSTALL_PREFIX}/gcc-linaro-arm-linux-gnueabihf-raspbian-x64)
set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
EOF

mkdir -p ${BUILD_DIR}/openal-soft-${OPENAL_VERSION}/build-rpi
cd ${BUILD_DIR}/openal-soft-${OPENAL_VERSION}/build-rpi
cmake -DLIBTYPE=STATIC -DCMAKE_TOOLCHAIN_FILE=../linux-rpi.cmake -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX}/gcc-linaro-arm-linux-gnueabihf-raspbian-x64 ..
make -j $(nproc) && make install

echo "##### Remove build directory"
rm -rf ${BUILD_DIR}
