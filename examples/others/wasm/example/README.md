### Wasm example

Desktop code runs in the browser unchanged, built with `GOOS=js GOARCH=wasm`. raylib is a prebuilt module from
`raylib/libs/wasm`, no C toolchain is needed.

Build with `./build.sh` (`build.bat` on Windows), output goes to `web/` (set `OUT` to change it):

* `app.wasm` - the Go program
* `raylib.js`, `raylib.wasm` - raylib
* `raylib-go.js`, `wasm_exec.js`, `index.html` - the loader and the page
* `resources.tar.gz` - the `resources` directory

The default is WebGL 1, set `TAGS=es3` for WebGL 2.

Serve the directory with any web server, for example:

    python3 -m http.server -d web

The resources are unpacked into an in-memory filesystem before the program starts, so `rl.LoadTexture("resources/raylib_logo.png")`
and Go's `os` functions work as on the desktop.

Alternatively, embed the files with `go:embed` and read them with `rl.NewAssetFromFS`, see the
[asset example](https://github.com/gen2brain/raylib-go/tree/master/examples/others/asset); no `resources.tar.gz` is needed then.
