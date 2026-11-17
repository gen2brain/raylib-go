package apple

import (
	"crypto/sha1"
	"crypto/sha256"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type rule struct {
	pattern  string
	omit     bool
	optional bool
	nested   bool
	weight   float64
	re       *regexp.Regexp
}

func compile(rules []rule) []rule {
	for i := range rules {
		rules[i].re = regexp.MustCompile(rules[i].pattern)
	}
	return rules
}

// Apple's tables from Security/libsecurity_codesigning/lib/bundlediskrep.cpp.
var macOSRules = compile([]rule{
	{pattern: "^Resources/"},
	{pattern: "^Resources/.*\\.lproj/", optional: true, weight: 1000},
	{pattern: "^Resources/Base\\.lproj/", weight: 1010},
	{pattern: "^Resources/.*\\.lproj/locversion.plist$", omit: true, weight: 1100},
})

var macOSRules2 = compile([]rule{
	{pattern: ".*\\.dSYM($|/)", weight: 11},
	{pattern: "^(.*/)?\\.DS_Store$", omit: true, weight: 2000},
	{pattern: "^(Frameworks|SharedFrameworks|PlugIns|Plug-ins|XPCServices|Helpers|MacOS|Library/(Automator|Spotlight|LoginItems))/", nested: true, weight: 10},
	{pattern: "^.*"},
	{pattern: "^Info\\.plist$", omit: true, weight: 20},
	{pattern: "^PkgInfo$", omit: true, weight: 20},
	{pattern: "^Resources/", weight: 20},
	{pattern: "^Resources/.*\\.lproj/", optional: true, weight: 1000},
	{pattern: "^Resources/Base\\.lproj/", weight: 1010},
	{pattern: "^Resources/.*\\.lproj/locversion.plist$", omit: true, weight: 1100},
	{pattern: "^[^/]+$", nested: true, weight: 10},
	{pattern: "^embedded\\.provisionprofile$", weight: 20},
	{pattern: "^version\\.plist$", weight: 20},
})

var iosRules = compile([]rule{
	{pattern: "^.*"},
	{pattern: "^.*\\.lproj/", optional: true, weight: 1000},
	{pattern: "^Base\\.lproj/", weight: 1010},
	{pattern: "^.*\\.lproj/locversion.plist$", omit: true, weight: 1100},
})

var iosRules2 = compile([]rule{
	{pattern: ".*\\.dSYM($|/)", weight: 11},
	{pattern: "^(.*/)?\\.DS_Store$", omit: true, weight: 2000},
	{pattern: "^.*"},
	{pattern: "^.*\\.lproj/", optional: true, weight: 1000},
	{pattern: "^Base\\.lproj/", weight: 1010},
	{pattern: "^.*\\.lproj/locversion.plist$", omit: true, weight: 1100},
	{pattern: "^Info\\.plist$", omit: true, weight: 20},
	{pattern: "^PkgInfo$", omit: true, weight: 20},
	{pattern: "^embedded\\.provisionprofile$", weight: 20},
	{pattern: "^version\\.plist$", weight: 20},
})

type Nested struct {
	CDHash      []byte
	Requirement string
}

func rulesPlist(rules []rule) map[string]any {
	m := map[string]any{}
	for _, r := range rules {
		if !r.omit && !r.optional && !r.nested && r.weight == 0 {
			m[r.pattern] = true
			continue
		}
		d := map[string]any{}
		if r.omit {
			d["omit"] = true
		}
		if r.optional {
			d["optional"] = true
		}
		if r.nested {
			d["nested"] = true
		}
		if r.weight != 0 {
			d["weight"] = r.weight
		}
		m[r.pattern] = d
	}
	return m
}

func matchRule(rules []rule, rel string) *rule {
	var best *rule
	bestWeight := -1.0
	for i := range rules {
		r := &rules[i]
		if !r.re.MatchString(rel) {
			continue
		}
		w := r.weight
		if w == 0 {
			w = 1
		}
		if w > bestWeight {
			best, bestWeight = r, w
		}
	}
	return best
}

func CodeResources(root string, macOS bool, exclude string, nested map[string]Nested) ([]byte, error) {
	rules1, rules2 := iosRules, iosRules2
	if macOS {
		rules1, rules2 = macOSRules, macOSRules2
	}

	files := map[string]any{}
	files2 := map[string]any{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if rel == "." || rel == exclude || rel == "_CodeSignature" || strings.HasPrefix(rel, "_CodeSignature/") || rel == "CodeResources" {
			if d.IsDir() && rel != "." {
				return fs.SkipDir
			}
			return nil
		}
		if n, ok := nested[rel]; ok {
			files2[rel] = map[string]any{"cdhash": n.CDHash, "requirement": n.Requirement}
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}

		r2 := matchRule(rules2, rel)
		if r2 != nil && r2.omit {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		entry := map[string]any{}
		if info.Mode()&fs.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			entry["symlink"] = target
		} else {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			h1 := sha1.Sum(data)
			h2 := sha256.Sum256(data)
			entry["hash"] = h1[:]
			entry["hash2"] = h2[:]
			if r1 := matchRule(rules1, rel); r1 != nil && !r1.omit {
				if r1.optional {
					files[rel] = map[string]any{"hash": h1[:], "optional": true}
				} else {
					files[rel] = h1[:]
				}
			}
		}
		if r2 != nil && r2.optional {
			entry["optional"] = true
		}
		files2[rel] = entry
		return nil
	})
	if err != nil {
		return nil, err
	}

	return WritePlist(map[string]any{
		"files":  files,
		"files2": files2,
		"rules":  rulesPlist(rules1),
		"rules2": rulesPlist(rules2),
	}), nil
}
