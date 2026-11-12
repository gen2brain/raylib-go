@echo off
setlocal

cd /d "%~dp0"
if "%OUT%"=="" set OUT=web
if not exist "%OUT%" mkdir "%OUT%"

for /f "delims=" %%i in ('go list -m -f "{{.Dir}}" github.com/gen2brain/raylib-go/raylib') do set RL=%%i
for /f "delims=" %%i in ('go env GOROOT') do set GR=%%i

set GOOS=js
set GOARCH=wasm
set CGO_ENABLED=0
go build -tags "%TAGS%" -ldflags="-s -w" -o "%OUT%\app.wasm" . || exit /b 1

if exist "%GR%\lib\wasm\wasm_exec.js" (
  copy /y "%GR%\lib\wasm\wasm_exec.js" "%OUT%" >nul
) else (
  copy /y "%GR%\misc\wasm\wasm_exec.js" "%OUT%" >nul
)

set MODULE=raylib
set T= %TAGS% 
set T=%T:,= %
if not "%T: es3 =%"=="%T%" set MODULE=raylib-es3

call :gunzip "%RL%\libs\wasm\%MODULE%.js.gz" "%OUT%\raylib.js" || exit /b 1
call :gunzip "%RL%\libs\wasm\%MODULE%.wasm.gz" "%OUT%\raylib.wasm" || exit /b 1
copy /y "%RL%\libs\wasm\raylib-go.js" "%OUT%" >nul
copy /y index.html "%OUT%" >nul

if exist resources tar -czf "%OUT%\resources.tar.gz" resources
exit /b 0

:gunzip
powershell -NoProfile -Command "$i=[IO.File]::OpenRead('%~1'); $o=[IO.File]::Create('%~2'); $g=New-Object IO.Compression.GZipStream($i,[IO.Compression.CompressionMode]::Decompress); $g.CopyTo($o); $g.Close(); $o.Close()"
exit /b %errorlevel%
