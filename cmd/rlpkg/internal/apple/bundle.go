package apple

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type BundleOptions struct {
	MacOS           bool
	Signer          *Signer
	Entitlements    []byte
	HardenedRuntime bool
}

func SignBundle(app string, o BundleOptions) ([]byte, error) {
	contents := app
	if o.MacOS {
		contents = filepath.Join(app, "Contents")
	}
	infoData, err := os.ReadFile(filepath.Join(contents, "Info.plist"))
	if err != nil {
		return nil, err
	}
	info, err := ParsePlist(infoData)
	if err != nil {
		return nil, fmt.Errorf("Info.plist: %w", err)
	}
	dict, _ := info.(map[string]any)
	exe, _ := dict["CFBundleExecutable"].(string)
	bundleID, _ := dict["CFBundleIdentifier"].(string)
	if exe == "" || bundleID == "" {
		return nil, errors.New("Info.plist needs CFBundleExecutable and CFBundleIdentifier")
	}
	mainExe := filepath.Join(contents, exe)
	if o.MacOS {
		mainExe = filepath.Join(contents, "MacOS", exe)
	}

	os.RemoveAll(filepath.Join(contents, "_CodeSignature"))

	nested := map[string]Nested{}
	var nestedPaths []string
	err = filepath.WalkDir(contents, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || path == mainExe {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil || !IsMachO(data) {
			return nil
		}
		nestedPaths = append(nestedPaths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(nestedPaths, func(i, j int) bool { return len(nestedPaths[i]) > len(nestedPaths[j]) })
	for _, path := range nestedPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		id := strings.TrimSuffix(filepath.Base(path), ".dylib")
		signed, cdhashes, err := SignMachO(data, o.Signer, MachOOptions{Identifier: id, HardenedRuntime: o.HardenedRuntime})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if err := os.WriteFile(path, signed, info.Mode()); err != nil {
			return nil, err
		}
		rel, _ := filepath.Rel(contents, path)
		nested[filepath.ToSlash(rel)] = Nested{CDHash: cdhashes[0][:20], Requirement: Requirement(id, o.Signer, cdhashes)}
	}

	mainRel, _ := filepath.Rel(contents, mainExe)
	resources, err := CodeResources(contents, o.MacOS, filepath.ToSlash(mainRel), nested)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(contents, "_CodeSignature"), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(contents, "_CodeSignature", "CodeResources"), resources, 0o644); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(mainExe)
	if err != nil {
		return nil, err
	}
	resHash := sha256.Sum256(resources)
	signed, cdhashes, err := SignMachO(data, o.Signer, MachOOptions{
		Identifier: bundleID, Entitlements: o.Entitlements, HardenedRuntime: o.HardenedRuntime,
		InfoPlist: infoData, ResourcesHash: resHash[:],
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", mainExe, err)
	}
	return cdhashes[0], os.WriteFile(mainExe, signed, 0o755)
}
