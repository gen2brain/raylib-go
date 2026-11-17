package main

import (
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"
)

type gzipResponse struct {
	http.ResponseWriter
	gz *gzip.Writer
}

func (w gzipResponse) Write(b []byte) (int, error) { return w.gz.Write(b) }

func (w gzipResponse) WriteHeader(code int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(code)
}

func runServe(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	addr := fs.String("addr", "localhost:8000", "listen `address`")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage:\n\n\trlpkg serve [--addr localhost:8000] <directory>\n\n")
		printDefaults(fs)
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return errors.New("serve needs the directory produced by package --os js")
	}
	dir := fs.Arg(0)
	if _, err := os.Stat(dir + "/index.html"); err != nil {
		return fmt.Errorf("%s: no index.html", dir)
	}

	mime.AddExtensionType(".wasm", "application/wasm")
	mime.AddExtensionType(".js", "text/javascript")

	files := http.FileServer(http.Dir(dir))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Add("Vary", "Accept-Encoding")
			r.Header.Del("Range")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			w = gzipResponse{ResponseWriter: w, gz: gz}
		}
		files.ServeHTTP(w, r)
	})

	fmt.Fprintf(os.Stderr, "serving %s at http://%s/\n", dir, *addr)
	return http.ListenAndServe(*addr, handler)
}
