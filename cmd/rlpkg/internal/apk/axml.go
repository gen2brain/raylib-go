package apk

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"unicode/utf16"
)

const (
	axmlStringPool   = 0x0001
	axmlResourceMap  = 0x0180
	axmlStartElement = 0x0102
	axmlEndElement   = 0x0103

	axmlTypeString = 0x03
	axmlTypeInt    = 0x10

	attrName              = 0x01010003
	attrLabel             = 0x01010001
	attrValue             = 0x01010024
	attrVersionCode       = 0x0101021b
	attrVersionName       = 0x0101021c
	attrScreenOrientation = 0x0101001e
	attrGLESVersion       = 0x01010281
)

// Orientations maps the android:screenOrientation names to their values.
var Orientations = map[string]int32{
	"sensorLandscape": 6,
	"sensorPortrait":  7,
	"fullSensor":      10,
}

type axml struct {
	strings []string
	utf8    bool
	resMap  []uint32
	chunks  [][]byte
}

func parseAXML(data []byte) (*axml, error) {
	if len(data) < 8 || binary.LittleEndian.Uint16(data) != 0x0003 {
		return nil, errors.New("not a binary XML file")
	}
	a := &axml{}
	off := 8
	for off < len(data) {
		if off+8 > len(data) {
			return nil, errors.New("truncated binary XML")
		}
		typ := binary.LittleEndian.Uint16(data[off:])
		size := int(binary.LittleEndian.Uint32(data[off+4:]))
		if size < 8 || off+size > len(data) {
			return nil, errors.New("bad chunk size in binary XML")
		}
		chunk := data[off : off+size]
		switch typ {
		case axmlStringPool:
			if err := a.parsePool(chunk); err != nil {
				return nil, err
			}
		case axmlResourceMap:
			hs := int(binary.LittleEndian.Uint16(chunk[2:]))
			for i := hs; i+4 <= len(chunk); i += 4 {
				a.resMap = append(a.resMap, binary.LittleEndian.Uint32(chunk[i:]))
			}
			a.chunks = append(a.chunks, chunk)
		default:
			a.chunks = append(a.chunks, chunk)
		}
		off += size
	}
	if a.strings == nil {
		return nil, errors.New("binary XML without a string pool")
	}
	return a, nil
}

func (a *axml) parsePool(chunk []byte) error {
	hs := int(binary.LittleEndian.Uint16(chunk[2:]))
	count := int(binary.LittleEndian.Uint32(chunk[8:]))
	flags := binary.LittleEndian.Uint32(chunk[16:])
	start := int(binary.LittleEndian.Uint32(chunk[20:]))
	a.utf8 = flags&0x100 != 0
	for i := range count {
		off := start + int(binary.LittleEndian.Uint32(chunk[hs+4*i:]))
		if off >= len(chunk) {
			return errors.New("bad string offset in binary XML")
		}
		if a.utf8 {
			_, n := axmlLen8(chunk[off:])
			off += n
			blen, n := axmlLen8(chunk[off:])
			off += n
			a.strings = append(a.strings, string(chunk[off:off+blen]))
		} else {
			clen := int(binary.LittleEndian.Uint16(chunk[off:]))
			off += 2
			if clen&0x8000 != 0 {
				clen = (clen&0x7fff)<<16 | int(binary.LittleEndian.Uint16(chunk[off:]))
				off += 2
			}
			u := make([]uint16, clen)
			for j := range u {
				u[j] = binary.LittleEndian.Uint16(chunk[off+2*j:])
			}
			a.strings = append(a.strings, string(utf16.Decode(u)))
		}
	}
	return nil
}

func axmlLen8(b []byte) (int, int) {
	if b[0]&0x80 != 0 {
		return int(b[0]&0x7f)<<8 | int(b[1]), 2
	}
	return int(b[0]), 1
}

func (a *axml) encodePool() []byte {
	var body []byte
	offsets := make([]byte, 4*len(a.strings))
	for i, s := range a.strings {
		binary.LittleEndian.PutUint32(offsets[4*i:], uint32(len(body)))
		if a.utf8 {
			body = append(body, axmlPut8(len(utf16.Encode([]rune(s))))...)
			body = append(body, axmlPut8(len(s))...)
			body = append(body, s...)
			body = append(body, 0)
		} else {
			u := utf16.Encode([]rune(s))
			if len(u) >= 0x8000 {
				body = binary.LittleEndian.AppendUint16(body, uint16(len(u)>>16|0x8000))
			}
			body = binary.LittleEndian.AppendUint16(body, uint16(len(u)))
			for _, c := range u {
				body = binary.LittleEndian.AppendUint16(body, c)
			}
			body = append(body, 0, 0)
		}
	}
	for len(body)%4 != 0 {
		body = append(body, 0)
	}

	hdr := make([]byte, 28)
	binary.LittleEndian.PutUint16(hdr, axmlStringPool)
	binary.LittleEndian.PutUint16(hdr[2:], 28)
	binary.LittleEndian.PutUint32(hdr[4:], uint32(28+len(offsets)+len(body)))
	binary.LittleEndian.PutUint32(hdr[8:], uint32(len(a.strings)))
	if a.utf8 {
		binary.LittleEndian.PutUint32(hdr[16:], 0x100)
	}
	binary.LittleEndian.PutUint32(hdr[20:], uint32(28+len(offsets)))
	out := append(hdr, offsets...)
	return append(out, body...)
}

func axmlPut8(n int) []byte {
	if n >= 0x80 {
		return []byte{byte(n>>8) | 0x80, byte(n)}
	}
	return []byte{byte(n)}
}

func (a *axml) bytes() []byte {
	pool := a.encodePool()
	size := 8 + len(pool)
	for _, c := range a.chunks {
		size += len(c)
	}
	out := make([]byte, 8, size)
	binary.LittleEndian.PutUint16(out, 0x0003)
	binary.LittleEndian.PutUint16(out[2:], 8)
	binary.LittleEndian.PutUint32(out[4:], uint32(size))
	out = append(out, pool...)
	for _, c := range a.chunks {
		out = append(out, c...)
	}
	return out
}

func (a *axml) addString(s string) uint32 {
	a.strings = append(a.strings, s)
	return uint32(len(a.strings) - 1)
}

func (a *axml) attrRes(nameIndex uint32) uint32 {
	if int(nameIndex) < len(a.resMap) {
		return a.resMap[nameIndex]
	}
	return 0
}

func (a *axml) elementName(chunk []byte) string {
	if binary.LittleEndian.Uint16(chunk) != axmlStartElement && binary.LittleEndian.Uint16(chunk) != axmlEndElement {
		return ""
	}
	i := binary.LittleEndian.Uint32(chunk[20:])
	if int(i) < len(a.strings) {
		return a.strings[i]
	}
	return ""
}

func (a *axml) attributes(chunk []byte) [][]byte {
	if binary.LittleEndian.Uint16(chunk) != axmlStartElement {
		return nil
	}
	start := 16 + int(binary.LittleEndian.Uint16(chunk[24:]))
	size := int(binary.LittleEndian.Uint16(chunk[26:]))
	count := int(binary.LittleEndian.Uint16(chunk[28:]))
	var attrs [][]byte
	for i := range count {
		off := start + i*size
		if off+size > len(chunk) {
			break
		}
		attrs = append(attrs, chunk[off:off+size])
	}
	return attrs
}

func (a *axml) findAttr(chunk []byte, res uint32, plain string) []byte {
	for _, at := range a.attributes(chunk) {
		name := binary.LittleEndian.Uint32(at[4:])
		if res != 0 && a.attrRes(name) == res {
			return at
		}
		if plain != "" && int(name) < len(a.strings) && a.strings[name] == plain && a.attrRes(name) == 0 {
			return at
		}
	}
	return nil
}

func (a *axml) attrString(at []byte) (string, bool) {
	raw := binary.LittleEndian.Uint32(at[8:])
	if raw == 0xffffffff || int(raw) >= len(a.strings) {
		return "", false
	}
	return a.strings[raw], true
}

func (a *axml) setAttrString(at []byte, s string) {
	i := a.addString(s)
	binary.LittleEndian.PutUint32(at[8:], i)
	at[15] = axmlTypeString
	binary.LittleEndian.PutUint32(at[16:], i)
}

func (a *axml) setAttrInt(at []byte, v uint32) {
	binary.LittleEndian.PutUint32(at[8:], 0xffffffff)
	at[15] = axmlTypeInt
	binary.LittleEndian.PutUint32(at[16:], v)
}

func (a *axml) setAttrData(at []byte, v uint32) {
	binary.LittleEndian.PutUint32(at[8:], 0xffffffff)
	binary.LittleEndian.PutUint32(at[16:], v)
}

func (a *axml) startElements(name string) [][]byte {
	var list [][]byte
	for _, c := range a.chunks {
		if binary.LittleEndian.Uint16(c) == axmlStartElement && a.elementName(c) == name {
			list = append(list, c)
		}
	}
	return list
}

type ManifestValues struct {
	Package     string
	Label       string
	Version     string
	Build       int
	Library     string
	GLESVersion uint32
	Orientation string
	Permissions []string
}

func PatchManifest(data []byte, v ManifestValues) ([]byte, error) {
	a, err := parseAXML(data)
	if err != nil {
		return nil, err
	}

	manifest := a.startElements("manifest")
	if len(manifest) != 1 {
		return nil, errors.New("manifest element not found")
	}
	pkgAttr := a.findAttr(manifest[0], 0, "package")
	if pkgAttr == nil {
		return nil, errors.New("manifest has no package attribute")
	}
	oldPkg, _ := a.attrString(pkgAttr)

	for _, c := range a.chunks {
		for _, at := range a.attributes(c) {
			s, ok := a.attrString(at)
			if !ok {
				continue
			}
			if s == oldPkg {
				a.setAttrString(at, v.Package)
			} else if strings.HasPrefix(s, oldPkg+".") {
				a.setAttrString(at, v.Package+strings.TrimPrefix(s, oldPkg))
			}
		}
	}

	if at := a.findAttr(manifest[0], attrVersionCode, ""); at != nil {
		a.setAttrInt(at, uint32(v.Build))
	}
	if at := a.findAttr(manifest[0], attrVersionName, ""); at != nil {
		a.setAttrString(at, v.Version)
	}

	app := a.startElements("application")
	if len(app) != 1 {
		return nil, errors.New("application element not found")
	}
	if at := a.findAttr(app[0], attrLabel, ""); at != nil {
		a.setAttrString(at, v.Label)
	}

	found := false
	for _, md := range a.startElements("meta-data") {
		if name := a.findAttr(md, attrName, ""); name != nil {
			if s, _ := a.attrString(name); s == "android.app.lib_name" {
				if at := a.findAttr(md, attrValue, ""); at != nil {
					a.setAttrString(at, v.Library)
					found = true
				}
			}
		}
	}
	if !found {
		return nil, errors.New("android.app.lib_name meta-data not found")
	}

	if v.GLESVersion != 0 {
		features := a.startElements("uses-feature")
		at := []byte(nil)
		if len(features) == 1 {
			at = a.findAttr(features[0], attrGLESVersion, "")
		}
		if at == nil {
			return nil, errors.New("uses-feature glEsVersion not found")
		}
		a.setAttrData(at, v.GLESVersion)
	}

	if v.Orientation != "" {
		value, ok := Orientations[v.Orientation]
		if !ok {
			return nil, fmt.Errorf("unknown screen orientation %q", v.Orientation)
		}
		activities := a.startElements("activity")
		at := []byte(nil)
		if len(activities) == 1 {
			at = a.findAttr(activities[0], attrScreenOrientation, "")
		}
		if at == nil {
			return nil, errors.New("activity screenOrientation not found")
		}
		a.setAttrData(at, uint32(value))
	}

	if err := a.setPermissions(v.Permissions); err != nil {
		return nil, err
	}
	return a.bytes(), nil
}

func (a *axml) setPermissions(perms []string) error {
	model := -1
	for i, c := range a.chunks {
		if binary.LittleEndian.Uint16(c) == axmlStartElement && a.elementName(c) == "uses-permission" && i+1 < len(a.chunks) {
			if binary.LittleEndian.Uint16(a.chunks[i+1]) == axmlEndElement && len(a.attributes(c)) == 1 {
				model = i
				break
			}
		}
	}
	if model < 0 {
		if len(perms) == 0 {
			return nil
		}
		return errors.New("no uses-permission element to clone")
	}
	start, end := a.chunks[model], a.chunks[model+1]

	var extra [][]byte
	for _, p := range perms {
		s := append([]byte(nil), start...)
		a.setAttrString(a.attributes(s)[0], p)
		extra = append(extra, s, append([]byte(nil), end...))
	}
	a.chunks = append(a.chunks[:model], append(extra, a.chunks[model+2:]...)...)
	return nil
}
