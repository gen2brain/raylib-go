### Raspberry Pi example

To cross compile example for Raspberry Pi you will need [RPi toolchain](https://github.com/raspberrypi/tools/tree/master/arm-bcm2708) and [userspace libraries](https://github.com/raspberrypi/firmware) (opt/vc).

Export path to RPi toolchain:

    export RPI_HOME=/opt/tools/arm-bcm2708/arm-linux-gnueabihf

Add toolchain bin directory to PATH:

    export PATH=${RPI_HOME}/bin:${PATH}

And compile example:

    CC=arm-linux-gnueabihf-gcc \
    CGO_CFLAGS="-I/opt/vc/include -I/opt/vc/include/interface/vcos -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads --sysroot=${RPI_HOME}/arm-linux-gnueabihf/sysroot" \
    CGO_LDFLAGS="-L/opt/vc/lib -L/opt/vc/lib64 --sysroot=${RPI_HOME}/arm-linux-gnueabihf/sysroot" \
    CGO_ENABLED=1 GOOS=linux GOARCH=arm go build
