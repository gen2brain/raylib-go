package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	exitName      = "exit_key"
	coverage_test = `
package main

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	go func() {
		time.Sleep(time.Second)
		exit_key = true
	}()
	main()
}`
)

// coverage packages
var coverpkgs = [...]string{
	"github.com/gen2brain/raylib-go/raylib",
	"github.com/gen2brain/raylib-go/raygui",
}

func main() {
	// only for linux
	switch runtime.GOOS {
	case "linux", "darwin":
		// do nothing
	default:
		log.Println("not supported OS")
		return
	}

	view := flag.Bool("v", false, "view html result of coverage test")
	flag.Parse()

	// test of test
	if !strings.Contains(coverage_test, exitName) {
		log.Fatalln("source of main_test have not name: ", exitName)
		return
	}

	// create templorary folder
	tempDir, err := os.MkdirTemp("", "coverage")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Prepare test folder: ", tempDir)

	// do not delete templorary folder
	// defer os.Remove(tempDir.Name())

	// 	// copy example folder
	// 	log.Println("copy `examples` folder")
	// 	if _, err := exec.Command("cp", "-r", "./examples/", tempDir).Output(); err != nil {
	// 		log.Fatalf("%v", err)
	// 		return
	// 	}

	// walking
	log.Println("walking by temp folder with `main.go` files")
	var testPaths []string
	if err := filepath.Walk("./examples/",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !strings.HasSuffix(path, "main.go") {
				return nil
			}
			index := strings.LastIndex(path, "/")
			testPaths = append(testPaths, path[:index])
			return nil
		}); err != nil {
		log.Fatalf("%v", err)
		return
	}

	// TODO: create test in parallel
	for index, testPath := range testPaths {
		if 0 < flag.NArg() && !strings.Contains(testPath, os.Args[len(os.Args)-1]) {
			continue
		}
		if err := coverage(index, testPath, tempDir); err != nil {
			log.Fatalf("%v", err)
			return
		}
	}

	// combine coverage files
	log.Println("combine coverage files")
	var combineFiles []string
	combineName := "all.coverage"
	if err := filepath.Walk(tempDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			combineFiles = append(combineFiles, path)
			fmt.Println(path)
			return nil
		}); err != nil {
		log.Fatalf("%v", err)
		return
	}
	cmd := exec.Command("gocovmerge", combineFiles...)
	cmd.Dir = "./raygui/"
	var out []byte
	out, err = cmd.Output()
	if err != nil {
		log.Fatalf("gocovmerge: %v\nRun:\ngo install github.com/wadey/gocovmerge@latest", err)
		return
	}
	if err = os.WriteFile("./raygui/"+combineName, []byte(out), 0644); err != nil {
		log.Fatalf("writefile: %v", err)
		return
	}

	// run coverage
	if *view {
		log.Println("run coverage")
		cmd = exec.Command(
			"go", "tool", "cover",
			"-html="+combineName,
		)
		cmd.Dir = "./raygui/"
		if _, err = cmd.Output(); err != nil {
			log.Fatalf("%v", err)
			return
		}
	} else {
		fmt.Fprintf(os.Stdout, "if you see the result, then run:\ngo run coverage.go -v\n")
	}
}

func coverage(index int, testPath, coverPath string) (err error) {
	name := strings.ReplaceAll(testPath, coverPath, "")
	log.Printf("run: %s", name)

	// has file main.go
	if _, err = os.Stat(testPath + "/main.go"); errors.Is(err, os.ErrNotExist) {
		log.Printf("file `main.go` not exist")
		return
	}

	// has exitName
	{
		var dat []byte
		dat, err = os.ReadFile(testPath + "/main.go")
		if err != nil {
			log.Fatalf("readfile: %v", err)
			return
		}
		if !strings.Contains(string(dat), exitName) {
			log.Printf("file `main.go` have not `%s`", exitName)
			return
		}
	}

	// create main_test file
	log.Println("create `main_test.go` file")
	if err = os.WriteFile(testPath+"/coverage_test.go", []byte(coverage_test), 0644); err != nil {
		log.Fatalf("writefile: %v", err)
		return
	}

	// testing
	log.Println("testing")
	for ic, cpkg := range coverpkgs {
		cmd := exec.Command(
			"go", "test",
			fmt.Sprintf("-coverprofile=%s/%d.%d.coverage", coverPath, index, ic),
			fmt.Sprintf("-coverpkg=%s", cpkg),
		)
		cmd.Dir = testPath
		if _, err = cmd.Output(); err != nil {
			log.Fatalf("%v", err)
			return
		}
	}
	return
}
