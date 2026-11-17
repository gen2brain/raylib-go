package apk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Field numbers of aapt2's Resources.proto.
const (
	xmlNodeElement = 1

	xmlElementName      = 3
	xmlElementAttribute = 4
	xmlElementChild     = 5

	xmlAttributeName  = 2
	xmlAttributeValue = 3
	xmlAttributeItem  = 6

	itemString    = 2
	itemPrimitive = 7
	stringValue   = 1
	primitiveInt  = 6
	primitiveHex  = 7

	tablePackage = 2
	packageName  = 2
)

func visitBundleXML(node []byte, fn func(element []protoField) ([]protoField, error)) ([]byte, error) {
	fields, err := protoParse(node)
	if err != nil {
		return nil, err
	}
	i := protoFind(fields, xmlNodeElement)
	if i < 0 {
		return node, nil
	}
	element, err := protoParse(fields[i].b)
	if err != nil {
		return nil, err
	}
	if element, err = fn(element); err != nil {
		return nil, err
	}
	for j := range element {
		if element[j].num == xmlElementChild {
			if element[j].b, err = visitBundleXML(element[j].b, fn); err != nil {
				return nil, err
			}
		}
	}
	fields[i].b = protoEncode(element)
	return protoEncode(fields), nil
}

func bundleAttr(element []protoField, name string) (int, []protoField) {
	for i, f := range element {
		if f.num != xmlElementAttribute {
			continue
		}
		attr, err := protoParse(f.b)
		if err == nil && protoGetString(attr, xmlAttributeName) == name {
			return i, attr
		}
	}
	return -1, nil
}

func bundleSetAttr(element []protoField, name, value string, item []protoField) bool {
	i, attr := bundleAttr(element, name)
	if i < 0 {
		return false
	}
	attr = protoSetString(attr, xmlAttributeValue, value)
	if item != nil {
		if j := protoFind(attr, xmlAttributeItem); j >= 0 {
			attr[j].b = protoEncode(item)
		} else {
			attr = append(attr, protoMessage(xmlAttributeItem, item))
		}
	}
	element[i].b = protoEncode(attr)
	return true
}

func bundleSetInt(element []protoField, name, value string, v int64) bool {
	i, attr := bundleAttr(element, name)
	if i < 0 {
		return false
	}
	j := protoFind(attr, xmlAttributeItem)
	if j < 0 {
		return false
	}
	item, err := protoParse(attr[j].b)
	if err != nil {
		return false
	}
	k := protoFind(item, itemPrimitive)
	if k < 0 {
		return false
	}
	prim, err := protoParse(item[k].b)
	if err != nil {
		return false
	}
	n := protoFind(prim, primitiveInt)
	if n < 0 {
		n = protoFind(prim, primitiveHex)
	}
	if n < 0 {
		return false
	}
	prim[n].v = uint64(v)
	if prim[n].num == primitiveHex {
		prim[n].v = uint64(uint32(v))
	}
	item[k].b = protoEncode(prim)
	attr[j].b = protoEncode(item)
	element[i].b = protoEncode(protoSetString(attr, xmlAttributeValue, value))
	return true
}

func bundleRenameValues(element []protoField, oldPkg, newPkg string) error {
	rename := func(s string) (string, bool) {
		if s == oldPkg {
			return newPkg, true
		}
		if strings.HasPrefix(s, oldPkg+".") {
			return newPkg + strings.TrimPrefix(s, oldPkg), true
		}
		return s, false
	}
	for i, f := range element {
		if f.num != xmlElementAttribute {
			continue
		}
		attr, err := protoParse(f.b)
		if err != nil {
			return err
		}
		s, ok := rename(protoGetString(attr, xmlAttributeValue))
		if !ok {
			continue
		}
		attr = protoSetString(attr, xmlAttributeValue, s)
		if j := protoFind(attr, xmlAttributeItem); j >= 0 {
			item, err := protoParse(attr[j].b)
			if err != nil {
				return err
			}
			if k := protoFind(item, itemString); k >= 0 {
				item[k] = protoMessage(itemString, []protoField{protoString(stringValue, s)})
				attr[j].b = protoEncode(item)
			}
		}
		element[i].b = protoEncode(attr)
	}
	return nil
}

func bundleSetPermissions(manifest []protoField, perms []string) ([]protoField, error) {
	model := -1
	for i, f := range manifest {
		if f.num != xmlElementChild {
			continue
		}
		node, err := protoParse(f.b)
		if err != nil {
			return nil, err
		}
		j := protoFind(node, xmlNodeElement)
		if j < 0 {
			continue
		}
		element, err := protoParse(node[j].b)
		if err != nil {
			return nil, err
		}
		if protoGetString(element, xmlElementName) == "uses-permission" {
			model = i
			break
		}
	}
	if model < 0 {
		if len(perms) == 0 {
			return manifest, nil
		}
		return nil, errors.New("no uses-permission element to clone")
	}

	var extra []protoField
	for _, p := range perms {
		node, err := visitBundleXML(manifest[model].b, func(element []protoField) ([]protoField, error) {
			bundleSetAttr(element, "name", p, nil)
			return element, nil
		})
		if err != nil {
			return nil, err
		}
		extra = append(extra, protoField{num: xmlElementChild, wire: wireBytes, b: node})
	}
	return append(manifest[:model:model], append(extra, manifest[model+1:]...)...), nil
}

// PatchBundleManifest is PatchManifest for the protobuf manifest of an Android App Bundle.
func PatchBundleManifest(data []byte, v ManifestValues) ([]byte, error) {
	root, err := protoParse(data)
	if err != nil {
		return nil, err
	}
	ri := protoFind(root, xmlNodeElement)
	if ri < 0 {
		return nil, errors.New("manifest element not found")
	}
	manifest, err := protoParse(root[ri].b)
	if err != nil {
		return nil, err
	}
	if protoGetString(manifest, xmlElementName) != "manifest" {
		return nil, errors.New("manifest element not found")
	}
	_, pkg := bundleAttr(manifest, "package")
	if pkg == nil {
		return nil, errors.New("manifest has no package attribute")
	}
	oldPkg := protoGetString(pkg, xmlAttributeValue)

	if manifest, err = bundleSetPermissions(manifest, v.Permissions); err != nil {
		return nil, err
	}
	root[ri].b = protoEncode(manifest)

	orientation, ok := Orientations[v.Orientation]
	if v.Orientation != "" && !ok {
		return nil, errors.New("unknown screen orientation " + v.Orientation)
	}

	found, gles, screen := false, v.GLESVersion == 0, v.Orientation == ""
	out, err := visitBundleXML(protoEncode(root), func(element []protoField) ([]protoField, error) {
		if err := bundleRenameValues(element, oldPkg, v.Package); err != nil {
			return nil, err
		}
		switch protoGetString(element, xmlElementName) {
		case "manifest":
			code := []protoField{protoMessage(itemPrimitive, []protoField{{num: primitiveInt, wire: wireVarint, v: uint64(uint32(v.Build))}})}
			bundleSetAttr(element, "versionCode", strconv.Itoa(v.Build), code)
			bundleSetAttr(element, "versionName", v.Version, nil)
		case "application":
			bundleSetAttr(element, "label", v.Label, nil)
		case "meta-data":
			if _, name := bundleAttr(element, "name"); name != nil && protoGetString(name, xmlAttributeValue) == "android.app.lib_name" {
				found = bundleSetAttr(element, "value", v.Library, nil)
			}
		case "uses-feature":
			if v.GLESVersion != 0 {
				gles = bundleSetInt(element, "glEsVersion", fmt.Sprintf("0x%08x", v.GLESVersion), int64(v.GLESVersion))
			}
		case "activity":
			if v.Orientation != "" {
				screen = bundleSetInt(element, "screenOrientation", v.Orientation, int64(orientation))
			}
		}
		return element, nil
	})
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("android.app.lib_name meta-data not found")
	}
	if !gles {
		return nil, errors.New("uses-feature glEsVersion not found")
	}
	if !screen {
		return nil, errors.New("activity screenOrientation not found")
	}
	return out, nil
}

// RenameBundlePackage is RenamePackage for the protobuf resource table of an Android App Bundle.
func RenameBundlePackage(table []byte, name string) ([]byte, error) {
	fields, err := protoParse(table)
	if err != nil {
		return nil, err
	}
	n := 0
	for i := range fields {
		if fields[i].num != tablePackage {
			continue
		}
		pkg, err := protoParse(fields[i].b)
		if err != nil {
			return nil, err
		}
		fields[i].b = protoEncode(protoSetString(pkg, packageName, name))
		n++
	}
	if n == 0 {
		return nil, errors.New("resource table has no package")
	}
	return protoEncode(fields), nil
}

var bundleAbiAlias = map[string]uint64{
	"armeabi-v7a": 2,
	"arm64-v8a":   3,
	"x86":         4,
	"x86_64":      5,
}

// BundleNativeTargeting returns base/native.pb, without which bundletool makes no ABI splits.
func BundleNativeTargeting(abis []string) ([]byte, error) {
	var dirs []protoField
	for _, abi := range abis {
		alias, ok := bundleAbiAlias[abi]
		if !ok {
			return nil, errors.New("unknown ABI " + abi)
		}
		targeting := protoMessage(2, []protoField{protoMessage(1, []protoField{{num: 1, wire: wireVarint, v: alias}})})
		dirs = append(dirs, protoMessage(1, []protoField{protoString(1, "lib/"+abi), targeting}))
	}
	return protoEncode(dirs), nil
}
