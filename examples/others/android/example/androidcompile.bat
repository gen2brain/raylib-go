@REM Set your desired api. Max is 31
@set ANDROID_API=31
@REM Your library name. If you change it here you should also change it in your android manifest...
@set LIBRARY_NAME=example
@REM Set your android sdk folder here
@set ANDROID_HOME=F:/AndroidSDK 
@REM Set your android NDK folder here. WARNING: NDK 27 is not supported yet. See https://github.com/raysan5/raylib/issues/4166
@set ANDROID_NDK_HOME=F:/AndroidSDK/ndk/23.2.8568313
@REM The target architecture for android. See https://developer.android.com/ndk/guides/abis. 
@REM Valid options are: armeabi-v7a,arm64-v8a,x86,x86_64 or "all" if you want to build for all architectures.
@set TARGET_ARCH="all"
@REM Automatic setup. Should work by default. Do not change anything below here
@set PATH=%ANDROID_NDK_HOME%/toolchains/llvm/prebuilt/windows-x86_64/bin;%PATH%
@set ANDROID_SYSROOT=%ANDROID_NDK_HOME%/toolchains/llvm/prebuilt/windows-x86_64/sysroot
@set ANDROID_TOOLCHAIN=%ANDROID_NDK_HOME%/toolchains/llvm/prebuilt/windows-x86_64
@IF %TARGET_ARCH% == "all" (
    @GOTO:BUILDALL
) else (
    @GOTO:MAIN
)

:COMPILE
        @echo compiling for platform %FL%
        @set CGO_CFLAGS="-I%ANDROID_SYSROOT%/usr/include -I%ANDROID_SYSROOT%/usr/include/%TRIPLE% --sysroot=%ANDROID_SYSROOT% -D__ANDROID_API__=%ANDROID_API%"
        @set CGO_LDFLAGS="-L%ANDROID_SYSROOT%/usr/lib/%TRIPLE%/%ANDROID_API% -L%ANDROID_TOOLCHAIN%/%TRIPLE%/lib --sysroot=%ANDROID_SYSROOT%"
        @set CGO_ENABLED=1
        @set GOOS=android
        @set GOARCH=arm
        @go build -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-soname,lib%LIBRARY_NAME%.so" -o=android/libs/%FL%/lib%LIBRARY_NAME%.so
@EXIT /B

:BUILDALL
    @set TARGET_ARCH="armeabi-v7a"
    @CALL:MAIN
    @set TARGET_ARCH="arm64-v8a"
    @CALL:MAIN
    @set TARGET_ARCH="x86"
    @CALL:MAIN
    @set TARGET_ARCH="x86_64"
    @CALL:MAIN
@EXIT /B

:MAIN
    @IF %TARGET_ARCH% == "armeabi-v7a" (
        @set CC="armv7a-linux-androideabi%ANDROID_API%-clang"
        @set TRIPLE=arm-linux-androideabi
        @set FL=armeabi-v7a
        @set GOARCH=arm
        @CALL:COMPILE )
    @IF %TARGET_ARCH% == "arm64-v8a" (
        @set CC="armv7a-linux-androideabi%ANDROID_API%-clang"
        @set TRIPLE=aarch64-linux-android
        @set FL=arm64-v8a
        @set GOARCH=arm64
        @CALL:COMPILE )
    @IF %TARGET_ARCH% == "x86" (
        @set CC="armv7a-linux-androideabi%ANDROID_API%-clang"
        @set TRIPLE=i686-linux-android
        @set FL=x86
        @set GOARCH=arm
        @CALL:COMPILE )
    @IF %TARGET_ARCH% == "x86_64" (
        @set CC="armv7a-linux-androideabi%ANDROID_API%-clang"
        @set TRIPLE=x86_64-linux-android
        @set FL=x86_64
        @set GOARCH=arm64
        @CALL:COMPILE )
@EXIT /B
