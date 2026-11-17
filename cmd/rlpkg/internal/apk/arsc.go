package apk

import (
	"encoding/binary"
	"errors"
	"fmt"
	"unicode/utf16"
)

const (
	arscTable        = 0x0002
	arscPackage      = 0x0200
	arscType         = 0x0201
	arscFlagSparse   = 0x01
	arscFlagOffset16 = 0x02
	arscAppPackageID = 0x7f
	arscNameOffset   = 12
	arscNameLen      = 128
)

func arscPool(data []byte) ([]string, int, error) {
	if len(data) < 28 || binary.LittleEndian.Uint16(data) != axmlStringPool {
		return nil, 0, errors.New("string pool expected in resources.arsc")
	}
	size := int(binary.LittleEndian.Uint32(data[4:]))
	if size > len(data) {
		return nil, 0, errors.New("truncated string pool in resources.arsc")
	}
	a := &axml{}
	if err := a.parsePool(data[:size]); err != nil {
		return nil, 0, err
	}
	return a.strings, size, nil
}

func RenamePackage(arsc []byte, name string) ([]byte, error) {
	if len(arsc) < 12 || binary.LittleEndian.Uint16(arsc) != arscTable {
		return nil, errors.New("not a resource table")
	}
	u := utf16.Encode([]rune(name))
	if len(u) >= arscNameLen {
		return nil, errors.New("package name too long for resources.arsc")
	}
	out := append([]byte(nil), arsc...)
	hs := int(binary.LittleEndian.Uint16(out[2:]))
	_, size, err := arscPool(out[hs:])
	if err != nil {
		return nil, err
	}

	renamed := false
	off := hs + size
	for off+8 <= len(out) {
		typ := binary.LittleEndian.Uint16(out[off:])
		csize := int(binary.LittleEndian.Uint32(out[off+4:]))
		if csize < 8 || off+csize > len(out) {
			return nil, errors.New("bad chunk in resources.arsc")
		}
		if typ == arscPackage && csize >= arscNameOffset+2*arscNameLen && binary.LittleEndian.Uint32(out[off+8:]) == arscAppPackageID {
			field := out[off+arscNameOffset : off+arscNameOffset+2*arscNameLen]
			clear(field)
			for i, c := range u {
				binary.LittleEndian.PutUint16(field[2*i:], c)
			}
			renamed = true
		}
		off += csize
	}
	if !renamed {
		return nil, errors.New("application package not found in resources.arsc")
	}
	return out, nil
}

func ResourcePaths(arsc []byte, resType, name string) ([]string, error) {
	if len(arsc) < 12 || binary.LittleEndian.Uint16(arsc) != arscTable {
		return nil, errors.New("not a resource table")
	}
	hs := int(binary.LittleEndian.Uint16(arsc[2:]))
	global, size, err := arscPool(arsc[hs:])
	if err != nil {
		return nil, err
	}

	var paths []string
	off := hs + size
	for off+8 <= len(arsc) {
		typ := binary.LittleEndian.Uint16(arsc[off:])
		csize := int(binary.LittleEndian.Uint32(arsc[off+4:]))
		if csize < 8 || off+csize > len(arsc) {
			return nil, errors.New("bad chunk in resources.arsc")
		}
		if typ == arscPackage {
			p, err := packageResourcePaths(arsc[off:off+csize], global, resType, name)
			if err != nil {
				return nil, err
			}
			paths = append(paths, p...)
		}
		off += csize
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("resource %s/%s not found in resources.arsc", resType, name)
	}
	return paths, nil
}

func packageResourcePaths(pkg []byte, global []string, resType, name string) ([]string, error) {
	hs := int(binary.LittleEndian.Uint16(pkg[2:]))
	typeOff := int(binary.LittleEndian.Uint32(pkg[0x10c:]))
	keyOff := int(binary.LittleEndian.Uint32(pkg[0x114:]))
	types, _, err := arscPool(pkg[typeOff:])
	if err != nil {
		return nil, err
	}
	keys, _, err := arscPool(pkg[keyOff:])
	if err != nil {
		return nil, err
	}

	var paths []string
	off := hs
	for off+8 <= len(pkg) {
		typ := binary.LittleEndian.Uint16(pkg[off:])
		chs := int(binary.LittleEndian.Uint16(pkg[off+2:]))
		csize := int(binary.LittleEndian.Uint32(pkg[off+4:]))
		if csize < 8 || off+csize > len(pkg) {
			return nil, errors.New("bad chunk in resources.arsc package")
		}
		if typ == arscType {
			chunk := pkg[off : off+csize]
			id := int(chunk[8])
			flags := chunk[9]
			count := int(binary.LittleEndian.Uint32(chunk[12:]))
			start := int(binary.LittleEndian.Uint32(chunk[16:]))
			if id >= 1 && id <= len(types) && types[id-1] == resType && flags&arscFlagSparse == 0 {
				for i := range count {
					var eoff int
					if flags&arscFlagOffset16 != 0 {
						v := binary.LittleEndian.Uint16(chunk[chs+2*i:])
						if v == 0xffff {
							continue
						}
						eoff = int(v) * 4
					} else {
						v := binary.LittleEndian.Uint32(chunk[chs+4*i:])
						if v == 0xffffffff {
							continue
						}
						eoff = int(v)
					}
					e := start + eoff
					if e+16 > len(chunk) {
						continue
					}
					esize := int(binary.LittleEndian.Uint16(chunk[e:]))
					key := int(binary.LittleEndian.Uint32(chunk[e+4:]))
					if binary.LittleEndian.Uint16(chunk[e+2:])&0x0001 != 0 || key >= len(keys) || keys[key] != name {
						continue
					}
					val := chunk[e+esize:]
					if len(val) < 8 || val[3] != axmlTypeString {
						continue
					}
					s := int(binary.LittleEndian.Uint32(val[4:]))
					if s < len(global) {
						paths = append(paths, global[s])
					}
				}
			}
		}
		off += csize
	}
	return paths, nil
}
