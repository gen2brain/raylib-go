### Raspberry Pi example

To compile example for Raspberry Pi you will need RPi toolchain and userspace libraries, and Go must be cross compiled for linux/arm with that toolchain.

There is a bootstrap.sh script that you can use to cross compile Go and OpenAL.

Compile Go and OpenAL, /usr/local is prefix where Go and RPi toolchain and libraries will be installed:

    ./bootstrap.sh /usr/local

After build is complete point GOROOT to new Go installation in /usr/local, and add toolchain bin directory to PATH:

    export GOROOT=/usr/local/go
    export PATH=/usr/local/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/bin:${PATH}

And compile example:

    CGO_CFLAGS="-I/usr/local/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/include -I/usr/local/vc/include -I/usr/local/vc/include/interface/vcos -I/usr/local/vc/include/interface/vmcs_host/linux -I/usr/local/vc/include/interface/vcos/pthreads" \
    CGO_LDFLAGS="-L/usr/local/vc/lib -L/usr/local/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/lib -ldl" \
    CC=arm-linux-gnueabihf-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 ${GOROOT}/bin/go build -v -x
