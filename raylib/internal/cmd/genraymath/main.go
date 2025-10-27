package main

import (
	"bytes"
	"cmp"
	"embed"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var (
	//go:embed *.tmpl
	templatesFS embed.FS
)

var (
	// don't generate tests or bindings for these functions
	skipTestAndBinding = []string{
		// not in raymath.h
		"Mat2MultiplyVector2",
		"Mat2Radians",
		"Mat2Transpose",
		"MatrixNormalize",
		"Mat2Set",
		"Vector2Cross",

		"MatrixToFloat",  // MatrixToFloatV tested
		"Vector3ToFloat", // Vector3ToFloatV tested
	}

	// these functions are tested manually
	skipTest = []string{
		// pointer params makes code gen more complicated - simpler just to write manually for now
		"Vector3OrthoNormalize",
		"QuaternionToAxisAngle",
		"MatrixDecompose",
	}

	// some outputs can differ a lot when the input is big (especially for rotations) so we need to
	// normalize test inputs
	normalize = []string{
		"Vector3RotateByAxisAngle",
		"Vector2Rotate",
		"QuaternionSlerp",
	}

	// the C versions of these functions use double params instead of floats
	useDouble = []string{
		"MatrixFrustum",
		"MatrixPerspective",
		"MatrixOrtho",
	}
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		goSrc         = flag.String("src", "raymath.go", "path to raylib-go's raymath.go")
		shouldFormat  = flag.Bool("format", true, "run Go formatter on output")
		inlineMethods = flag.Bool("inline-methods", false, "inline function bodies when generating methods")

		skipGenTests   = flag.Bool("skip-tests", false, "don't generate tests or bindings")
		skipGenMethods = flag.Bool("skip-methods", false, "don't generate methods")

		fuzzTime     = flag.Duration("fuzztime", 0, "run fuzz tests for given time after generating")
		fuzzFrom     = flag.String("fuzzfrom", "", "resume fuzzing from given function name")
		fuzzContinue = flag.Bool("fuzzcontinue", false, "don't exit on fuzz test error")
	)
	flag.Parse()

	if *fuzzTime > 0 {
		*skipGenTests = true
		*skipGenMethods = true
	}

	t := template.New("")
	t.Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"testvar": func(name string) string {
			switch name {
			// make sure these variable names aren't used for tests/fuzzing/benchmarks
			case "b", "t", "f":
				return name + "1"
			default:
				return name
			}
		},
	})
	templates, err := t.ParseFS(templatesFS, "*.tmpl")
	if err != nil {
		return err
	}

	data, err := os.ReadFile(*goSrc)
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, *goSrc, data, parser.SkipObjectResolution|parser.ParseComments)
	if err != nil {
		return err
	}
	var funcs []funcInfo
	ast.Inspect(file, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		t, err := parseFunc(fn, fset, *inlineMethods)
		if err != nil {
			if !errors.Is(err, errSkip) {
				log.Printf("%s: %v", fn.Name, err)
			}
			return false
		}
		funcs = append(funcs, t)
		return false
	})
	slices.SortFunc(funcs, func(a, b funcInfo) int {
		return cmp.Compare(a.Name, b.Name)
	})

	var outputs []*output
	if !*skipGenTests {
		outputs = append(outputs,
			&output{name: "binding.go"}, // C bindings can't live in the test unfortunately
			&output{name: "test.go"},
		)
	}
	if !*skipGenMethods {
		outputs = append(outputs,
			&output{name: "methods.go"},
		)
	}
	for _, v := range outputs {
		var buf bytes.Buffer
		err := templates.ExecuteTemplate(&buf, v.name+".tmpl", funcs)
		if err != nil {
			return err
		}
		v.data = buf.Bytes()
		if *shouldFormat {
			v.data, err = format.Source(v.data)
			if err != nil {
				return err
			}
		}
	}

	// only write once we're sure the output data is good
	for _, v := range outputs {
		name := "raymath_generated_" + v.name
		path := filepath.Join(filepath.Dir(*goSrc), name)
		err = os.WriteFile(path, v.data, 0600)
		if err != nil {
			return err
		}
	}

	if *fuzzTime <= 0 {
		return nil
	}

	// helper for running all the fuzz tests
	var total int
	for _, fn := range funcs {
		if fn.SkipTest || len(fn.Params) == 0 {
			continue
		}
		total++
	}
	var failed []string
	ok := *fuzzFrom == ""
	var n int
	for _, fn := range funcs {
		if fn.SkipTest || len(fn.Params) == 0 {
			continue
		}
		n++
		if ok || fn.Name == *fuzzFrom {
			ok = true
		} else {
			continue
		}
		log.Printf("running fuzz test %d/%d for %q", n, total, fn.Name)
		testName := "^Fuzz" + fn.Name + "$"
		cmd := exec.Command("go", "test",
			"-run", "^$",
			"-fuzz", testName,
			"-fuzztime", (*fuzzTime).String(),
		)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Dir = filepath.Dir(*goSrc)
		err := cmd.Run()
		if err != nil {
			if !*fuzzContinue {
				return err
			}
			failed = append(failed, fn.Name)
		}
	}
	if len(failed) != 0 {
		log.Println(failed)
		return errors.New("failed")
	}
	return nil
}

type output struct {
	name string
	data []byte
}

var errSkip = errors.New("skip")

func parseFunc(fn *ast.FuncDecl, fset *token.FileSet, inlineMethods bool) (funcInfo, error) {
	var ret funcInfo

	if !fn.Name.IsExported() {
		return ret, errSkip
	}

	ret.Name = fn.Name.Name
	if slices.Contains(skipTestAndBinding, ret.Name) {
		ret.SkipBinding = true
		ret.SkipTest = true
	}
	if slices.Contains(skipTest, ret.Name) {
		ret.SkipTest = true
	}

	if fn.Type.Params != nil {
		for _, v := range fn.Type.Params.List {
			typeName, ok := exprString(v.Type)
			if !ok {
				return ret, fmt.Errorf("invalid param type: %#v", v.Type)
			}
			if !ret.SkipTest && strings.HasPrefix(typeName, "*") {
				log.Printf("skipping %q with pointer param", ret.Name)
				ret.SkipTest = true
			}
			names := make([]string, 0, len(v.Names))
			for _, name := range v.Names {
				names = append(names, name.Name)
			}
			ret.Params = append(ret.Params, param{
				Names:     names,
				TypeName:  typeName,
				Normalize: strings.HasPrefix(fn.Name.Name, typeName) && slices.Contains(normalize, fn.Name.Name),
				UseDouble: typeName == "float32" && slices.Contains(useDouble, fn.Name.Name),
			})
		}
	}
	if fn.Type.Results != nil {
		if len(fn.Type.Results.List) > 1 {
			return ret, fmt.Errorf("invalid number of results: %v", len(fn.Type.Results.List))
		}
		v := fn.Type.Results.List[0]
		typeName, ok := exprString(v.Type)
		if !ok {
			return ret, fmt.Errorf("invalid result type: %#v", v.Type)
		}
		if !ret.SkipTest && strings.HasPrefix(typeName, "*") {
			log.Printf("skipping %q with pointer return type", ret.Name)
			ret.SkipTest = true
		}
		ret.ReturnType = typeName
	}
	if ret.SkipBinding {
		ret.SkipTest = true
	}

	if fn.Doc != nil {
		ret.Comments = make([]string, 0, len(fn.Doc.List))
		for _, v := range fn.Doc.List {
			text := v.Text
			if ret.CanBeMethod() {
				text = strings.ReplaceAll(text, ret.Name, ret.MethodName())
			}
			ret.Comments = append(ret.Comments, text)
		}
	}
	if inlineMethods && ret.CanBeMethod() {
		// rename first parameter
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if ident.Name == ret.Params[0].Names[0] {
				ident.Name = ret.Self()
			}
			return true
		})

		var buf bytes.Buffer
		// some comments are missing but oh well
		err := format.Node(&buf, fset, fn.Body)
		if err != nil {
			return ret, err
		}
		ret.Body = buf.String()
	}
	return ret, nil
}

type funcInfo struct {
	Name        string
	Params      []param
	ReturnType  string
	Comments    []string
	SkipTest    bool
	SkipBinding bool
	Body        string
}

func (f funcInfo) CanBeMethod() bool {
	if len(f.Params) == 0 {
		return false
	}
	return strings.HasPrefix(f.Name, f.Params[0].TypeName)
}

func (f funcInfo) Self() string {
	return strings.ToLower(string(f.Name[0]))
}

func (f funcInfo) Struct() string {
	return f.Params[0].TypeName
}

func (f funcInfo) MethodName() string {
	return strings.TrimPrefix(f.Name, f.Params[0].TypeName)
}

func (f funcInfo) ReturnValue() string {
	t := f.ReturnType
	switch t {
	case "float32":
		return "float32(ret)"
	case "bool":
		return "ret != 0"
	default:
		return "*(*" + t + ")(unsafe.Pointer(&ret))"
	}
}

func (f funcInfo) TestNotEqual(a, b string) string {
	t := f.ReturnType
	if t == "bool" {
		return a + " != " + b
	}

	rightBracketIndex := strings.IndexRune(t, ']')
	if rightBracketIndex != -1 { // slice or array
		t = t[rightBracketIndex+1:] + "Slice"
		if rightBracketIndex != 1 { // array
			a += "[:]"
			b += "[:]"
		}
	}

	first, size := utf8.DecodeRuneInString(t)
	t = string(unicode.ToUpper(first)) + t[size:]
	return fmt.Sprintf("!test%sEquals(%s, %s)", t, a, b)
}

type param struct {
	Names     []string
	TypeName  string
	Normalize bool
	UseDouble bool
}

func (p param) CType(arg string) string {
	name := p.TypeName
	if p.UseDouble {
		name = "float64"
	}
	name, isPointer := strings.CutPrefix(name, "*")

	var cname string
	switch name {
	case "float32":
		cname = "float"
	case "float64":
		cname = "double"
	default:
		if isPointer {
			return "(*C." + name + ")(unsafe.Pointer(" + arg + "))"
		}
		return "*(*C." + name + ")(unsafe.Pointer(&" + arg + "))"
	}

	if isPointer {
		return "(*C." + cname + ")(" + arg + ")"
	}
	return "C." + cname + "(" + arg + ")"
}

func exprString(expr ast.Expr) (string, bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, true
	case *ast.ArrayType:
		s, ok := exprString(t.Elt)
		if !ok {
			return "", false
		}
		var length string
		if lit, ok := t.Len.(*ast.BasicLit); ok {
			length = lit.Value
		}
		return "[" + length + "]" + s, true
	case *ast.StarExpr:
		s, ok := exprString(t.X)
		if !ok {
			return "", false
		}
		return "*" + s, true
	default:
		return "", false
	}
}
