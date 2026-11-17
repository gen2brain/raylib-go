package main

import (
	_ "embed"
	"fmt"
	"html"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

//go:embed wasm/index.html
var indexHTML string

func packageJS(c *config) error {
	dir := filepath.Join(c.out, c.exe)
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	env := []string{"GOOS=js", "GOARCH=wasm", "CGO_ENABLED=0"}
	if err := goBuildWith(c, env, filepath.Join(dir, "app.wasm"), nil, nil, nil); err != nil {
		return err
	}

	module := "raylib"
	if c.hasTag("es3") {
		module = "raylib-es3"
	}
	for _, ext := range []string{"js", "wasm"} {
		src := filepath.Join(c.rlDir, "libs", "wasm", module+"."+ext+".gz")
		data, err := readGzip(src)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, "raylib."+ext), data, 0o644); err != nil {
			return err
		}
	}
	loader, err := os.ReadFile(filepath.Join(c.rlDir, "libs", "wasm", "raylib-go.js"))
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "raylib-go.js"), loader, 0o644); err != nil {
		return err
	}

	page := strings.Replace(indexHTML, "<title>raylib-go</title>", "<title>"+html.EscapeString(c.name)+"</title>", 1)
	if len(c.files) > 0 {
		var files []tarFile
		for _, f := range c.files {
			files = append(files, tarFile{path.Clean(f.name), int64(f.mode), f.data})
		}
		if err := writeTarGz(filepath.Join(dir, "resources.tar.gz"), files); err != nil {
			return err
		}
		page = strings.Replace(page, "assets: null", "assets: 'resources.tar.gz'", 1)
	}
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(page), 0o644); err != nil {
		return err
	}

	goroot, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return err
	}
	root := strings.TrimSpace(string(goroot))
	execJS := filepath.Join(root, "lib", "wasm", "wasm_exec.js")
	if _, err := os.Stat(execJS); err != nil {
		execJS = filepath.Join(root, "misc", "wasm", "wasm_exec.js")
	}
	data, err := os.ReadFile(execJS)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "wasm_exec.js"), data, 0o644); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, dir)
	return nil
}
