package msix

import (
	"bytes"
	"compress/flate"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/authenticode"
)

const (
	manifestName     = "AppxManifest.xml"
	blockMapName     = "AppxBlockMap.xml"
	contentTypesName = "[Content_Types].xml"
	signatureName    = "AppxSignature.p7x"

	blockSize = 64 * 1024

	signatureOverride = `<Override ContentType="application/vnd.ms-appx.signature" PartName="/AppxSignature.p7x"/>`
)

type File struct {
	Name string
	Data []byte
}

var contentTypes = map[string]string{
	"exe": "application/x-msdownload",
	"dll": "application/x-msdownload",
	"png": "image/png",
	"xml": "application/xml",
	"txt": "text/plain",
}

var stored = map[string]bool{"png": true}

type entry struct {
	name       string
	method     uint16
	crc        uint32
	compressed uint64
	size       uint64
	offset     uint64
	modTime    uint16
	modDate    uint16
	flags      uint16
	needed     uint16
	extra      []byte
	comment    []byte
	madeBy     uint16
	internal   uint16
	external   uint32
	raw        []byte
}

func dosTime(t time.Time) (uint16, uint16) {
	return uint16(t.Hour()<<11 | t.Minute()<<5 | t.Second()/2), uint16((t.Year()-1980)<<9 | int(t.Month())<<5 | t.Day())
}

func encodeName(name string) string {
	var b strings.Builder
	for i := range len(name) {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', strings.IndexByte("-._~/()!$&'+,;=@", c) >= 0:
			b.WriteByte(c)
		default:
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

func extension(name string) string {
	return strings.ToLower(strings.TrimPrefix(path.Ext(name), "."))
}

func deflateBlocks(data []byte) ([]byte, []int, error) {
	var out bytes.Buffer
	var sizes []int
	for len(data) > 0 {
		n := min(len(data), blockSize)
		start := out.Len()
		w, err := flate.NewWriter(&out, flate.BestCompression)
		if err != nil {
			return nil, nil, err
		}
		if _, err := w.Write(data[:n]); err != nil {
			return nil, nil, err
		}
		if err := w.Flush(); err != nil {
			return nil, nil, err
		}
		sizes = append(sizes, out.Len()-start)
		data = data[n:]
	}
	out.Write([]byte{0x03, 0x00})
	return out.Bytes(), sizes, nil
}

type writer struct {
	out     bytes.Buffer
	entries []entry
	now     time.Time
}

func (w *writer) localHeader(e entry) []byte {
	b := []byte{'P', 'K', 3, 4}
	b = binary.LittleEndian.AppendUint16(b, e.needed)
	b = binary.LittleEndian.AppendUint16(b, e.flags)
	b = binary.LittleEndian.AppendUint16(b, e.method)
	b = binary.LittleEndian.AppendUint16(b, e.modTime)
	b = binary.LittleEndian.AppendUint16(b, e.modDate)
	b = binary.LittleEndian.AppendUint32(b, e.crc)
	b = binary.LittleEndian.AppendUint32(b, uint32(e.compressed))
	b = binary.LittleEndian.AppendUint32(b, uint32(e.size))
	b = binary.LittleEndian.AppendUint16(b, uint16(len(e.name)))
	b = binary.LittleEndian.AppendUint16(b, 0)
	return append(b, e.name...)
}

func (w *writer) add(name string, data []byte, compress bool) ([]int, error) {
	if uint64(len(data)) >= 0xffffffff {
		return nil, fmt.Errorf("%s: larger than 4 GB", name)
	}
	if uint64(w.out.Len()) >= 0xffffffff {
		return nil, errors.New("package larger than 4 GB")
	}
	e := entry{name: name, crc: crc32.ChecksumIEEE(data), size: uint64(len(data)), offset: uint64(w.out.Len()), needed: 20, madeBy: 45}
	e.modTime, e.modDate = dosTime(w.now)
	body := data
	var sizes []int
	if compress {
		var err error
		if body, sizes, err = deflateBlocks(data); err != nil {
			return nil, err
		}
		e.method = 8
	}
	e.compressed = uint64(len(body))
	w.out.Write(w.localHeader(e))
	w.out.Write(body)
	w.entries = append(w.entries, e)
	return sizes, nil
}

func centralHeader(e entry, offset uint64) []byte {
	if e.raw != nil {
		return e.raw
	}
	b := []byte{'P', 'K', 1, 2}
	b = binary.LittleEndian.AppendUint16(b, e.madeBy)
	b = binary.LittleEndian.AppendUint16(b, e.needed)
	b = binary.LittleEndian.AppendUint16(b, e.flags)
	b = binary.LittleEndian.AppendUint16(b, e.method)
	b = binary.LittleEndian.AppendUint16(b, e.modTime)
	b = binary.LittleEndian.AppendUint16(b, e.modDate)
	b = binary.LittleEndian.AppendUint32(b, e.crc)
	b = binary.LittleEndian.AppendUint32(b, uint32(min(e.compressed, 0xffffffff)))
	b = binary.LittleEndian.AppendUint32(b, uint32(min(e.size, 0xffffffff)))
	b = binary.LittleEndian.AppendUint16(b, uint16(len(e.name)))
	b = binary.LittleEndian.AppendUint16(b, uint16(len(e.extra)))
	b = binary.LittleEndian.AppendUint16(b, uint16(len(e.comment)))
	b = binary.LittleEndian.AppendUint16(b, 0)
	b = binary.LittleEndian.AppendUint16(b, e.internal)
	b = binary.LittleEndian.AppendUint32(b, e.external)
	b = binary.LittleEndian.AppendUint32(b, uint32(min(offset, 0xffffffff)))
	b = append(b, e.name...)
	b = append(b, e.extra...)
	return append(b, e.comment...)
}

func directory(entries []entry, offset uint64) []byte {
	var b []byte
	for _, e := range entries {
		b = append(b, centralHeader(e, e.offset)...)
	}
	size := uint64(len(b))
	b = append(b, 'P', 'K', 6, 6)
	b = binary.LittleEndian.AppendUint64(b, 44)
	b = binary.LittleEndian.AppendUint16(b, 45)
	b = binary.LittleEndian.AppendUint16(b, 45)
	b = binary.LittleEndian.AppendUint32(b, 0)
	b = binary.LittleEndian.AppendUint32(b, 0)
	b = binary.LittleEndian.AppendUint64(b, uint64(len(entries)))
	b = binary.LittleEndian.AppendUint64(b, uint64(len(entries)))
	b = binary.LittleEndian.AppendUint64(b, size)
	b = binary.LittleEndian.AppendUint64(b, offset)
	b = append(b, 'P', 'K', 6, 7)
	b = binary.LittleEndian.AppendUint32(b, 0)
	b = binary.LittleEndian.AppendUint64(b, offset+size)
	b = binary.LittleEndian.AppendUint32(b, 1)
	b = append(b, 'P', 'K', 5, 6, 0, 0, 0, 0, 0xff, 0xff, 0xff, 0xff)
	b = binary.LittleEndian.AppendUint32(b, 0xffffffff)
	b = binary.LittleEndian.AppendUint32(b, 0xffffffff)
	return binary.LittleEndian.AppendUint16(b, 0)
}

func escape(s string) string {
	var b bytes.Buffer
	xml.EscapeText(&b, []byte(s))
	return b.String()
}

func Write(manifest []byte, files []File, signer *authenticode.Signer) ([]byte, error) {
	w := &writer{now: time.Now()}
	var blockMap, types bytes.Buffer
	blockMap.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>` + "\r\n")
	blockMap.WriteString(`<BlockMap xmlns="http://schemas.microsoft.com/appx/2010/blockmap" HashMethod="http://www.w3.org/2001/04/xmlenc#sha256">`)
	types.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\r\n")
	types.WriteString(`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`)

	seen := map[string]bool{}
	for _, f := range append(files, File{manifestName, manifest}) {
		if f.Name != manifestName && strings.EqualFold(f.Name, manifestName) || strings.EqualFold(f.Name, blockMapName) ||
			strings.EqualFold(f.Name, contentTypesName) || strings.EqualFold(f.Name, signatureName) {
			return nil, fmt.Errorf("%s: reserved name", f.Name)
		}
		ext := extension(f.Name)
		name := encodeName(f.Name)
		sizes, err := w.add(name, f.Data, !stored[ext])
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(&blockMap, `<File Name="%s" Size="%d" LfhSize="%d">`, escape(strings.ReplaceAll(f.Name, "/", `\`)), len(f.Data), 30+len(name))
		for i := 0; i*blockSize < len(f.Data); i++ {
			sum := sha256.Sum256(f.Data[i*blockSize : min((i+1)*blockSize, len(f.Data))])
			fmt.Fprintf(&blockMap, `<Block Hash="%s"`, base64.StdEncoding.EncodeToString(sum[:]))
			if sizes != nil {
				fmt.Fprintf(&blockMap, ` Size="%d"`, sizes[i])
			}
			blockMap.WriteString("/>")
		}
		blockMap.WriteString("</File>")

		switch {
		case f.Name == manifestName:
			fmt.Fprintf(&types, `<Override ContentType="application/vnd.ms-appx.manifest+xml" PartName="/%s"/>`, manifestName)
		case ext == "" || contentTypes[ext] == "":
			fmt.Fprintf(&types, `<Override ContentType="application/octet-stream" PartName="/%s"/>`, escape(name))
		case !seen[ext]:
			seen[ext] = true
			fmt.Fprintf(&types, `<Default ContentType="%s" Extension="%s"/>`, contentTypes[ext], ext)
		}
	}
	blockMap.WriteString("</BlockMap>")
	fmt.Fprintf(&types, `<Override ContentType="application/vnd.ms-appx.blockmap+xml" PartName="/%s"/>`, blockMapName)
	types.WriteString("</Types>")

	if _, err := w.add(blockMapName, blockMap.Bytes(), true); err != nil {
		return nil, err
	}
	if _, err := w.add(contentTypesName, types.Bytes(), true); err != nil {
		return nil, err
	}
	data := append(w.out.Bytes(), directory(w.entries, uint64(w.out.Len()))...)
	if signer == nil {
		return data, nil
	}
	return Sign(data, signer)
}

func checkPublisher(manifest []byte, signer *authenticode.Signer) error {
	if signer == nil || len(signer.Chain) == 0 {
		return errors.New("a certificate is required")
	}
	var pkg struct {
		Identity struct {
			Publisher string `xml:"Publisher,attr"`
		}
	}
	if err := xml.Unmarshal(manifest, &pkg); err != nil {
		return fmt.Errorf("%s: %w", manifestName, err)
	}
	subject, err := Publisher(signer.Chain[0])
	if err != nil {
		return err
	}
	if canonical(pkg.Identity.Publisher) != canonical(subject) {
		return fmt.Errorf("the package publisher %q is not the certificate subject %q", pkg.Identity.Publisher, subject)
	}
	return nil
}

type archive struct {
	entries []entry
	end     uint64
}

func parse(data []byte) (*archive, error) {
	bad := errors.New("not an MSIX package")
	le := binary.LittleEndian
	eocd := bytes.LastIndex(data, []byte{'P', 'K', 5, 6})
	if eocd < 20 || eocd+22 > len(data) || eocd+22+int(le.Uint16(data[eocd+20:])) != len(data) {
		return nil, bad
	}
	count, size, offset := uint64(le.Uint16(data[eocd+10:])), uint64(le.Uint32(data[eocd+12:])), uint64(le.Uint32(data[eocd+16:]))
	if loc := eocd - 20; bytes.Equal(data[loc:loc+4], []byte{'P', 'K', 6, 7}) {
		rec := le.Uint64(data[loc+8:])
		if rec+56 > uint64(len(data)) || !bytes.Equal(data[rec:rec+4], []byte{'P', 'K', 6, 6}) {
			return nil, bad
		}
		count, size, offset = le.Uint64(data[rec+32:]), le.Uint64(data[rec+40:]), le.Uint64(data[rec+48:])
	}
	if offset+size > uint64(len(data)) || offset+size < offset {
		return nil, bad
	}

	a := &archive{end: offset}
	dir := data[offset : offset+size]
	for range count {
		if len(dir) < 46 || !bytes.Equal(dir[:4], []byte{'P', 'K', 1, 2}) {
			return nil, bad
		}
		nameLen, extraLen, commentLen := int(le.Uint16(dir[28:])), int(le.Uint16(dir[30:])), int(le.Uint16(dir[32:]))
		if 46+nameLen+extraLen+commentLen > len(dir) {
			return nil, bad
		}
		e := entry{
			madeBy: le.Uint16(dir[4:]), needed: le.Uint16(dir[6:]), flags: le.Uint16(dir[8:]), method: le.Uint16(dir[10:]),
			modTime: le.Uint16(dir[12:]), modDate: le.Uint16(dir[14:]), crc: le.Uint32(dir[16:]),
			compressed: uint64(le.Uint32(dir[20:])), size: uint64(le.Uint32(dir[24:])),
			internal: le.Uint16(dir[36:]), external: le.Uint32(dir[38:]), offset: uint64(le.Uint32(dir[42:])),
			name:    string(dir[46 : 46+nameLen]),
			extra:   bytes.Clone(dir[46+nameLen : 46+nameLen+extraLen]),
			comment: bytes.Clone(dir[46+nameLen+extraLen : 46+nameLen+extraLen+commentLen]),
			raw:     bytes.Clone(dir[:46+nameLen+extraLen+commentLen]),
		}
		for x := e.extra; len(x) >= 4; {
			id, n := le.Uint16(x), int(le.Uint16(x[2:]))
			if 4+n > len(x) {
				return nil, bad
			}
			if id == 1 {
				v := x[4 : 4+n]
				for _, field := range []*uint64{&e.size, &e.compressed, &e.offset} {
					if *field == 0xffffffff {
						if len(v) < 8 {
							return nil, bad
						}
						*field, v = le.Uint64(v), v[8:]
					}
				}
			}
			x = x[4+n:]
		}
		if e.offset >= offset {
			return nil, bad
		}
		a.entries = append(a.entries, e)
		dir = dir[46+nameLen+extraLen+commentLen:]
	}
	return a, nil
}

func (a *archive) read(data []byte, e entry) ([]byte, error) {
	le := binary.LittleEndian
	if e.offset+30 > a.end || !bytes.Equal(data[e.offset:e.offset+4], []byte{'P', 'K', 3, 4}) {
		return nil, fmt.Errorf("%s: bad local header", e.name)
	}
	start := e.offset + 30 + uint64(le.Uint16(data[e.offset+26:])) + uint64(le.Uint16(data[e.offset+28:]))
	if start+e.compressed > a.end || start+e.compressed < start {
		return nil, fmt.Errorf("%s: truncated", e.name)
	}
	body := data[start : start+e.compressed]
	if e.method == 0 {
		return body, nil
	}
	if e.method != 8 {
		return nil, fmt.Errorf("%s: unsupported compression", e.name)
	}
	out, err := io.ReadAll(io.LimitReader(flate.NewReader(bytes.NewReader(body)), int64(e.size)+1))
	if err != nil || uint64(len(out)) != e.size {
		return nil, fmt.Errorf("%s: bad data", e.name)
	}
	return out, nil
}

func Sign(data []byte, signer *authenticode.Signer) ([]byte, error) {
	a, err := parse(data)
	if err != nil {
		return nil, err
	}
	if n := len(a.entries); n > 0 && a.entries[n-1].name == signatureName {
		a.end = a.entries[n-1].offset
		a.entries = a.entries[:n-1]
	}
	n := len(a.entries)
	if n < 3 || a.entries[n-1].name != contentTypesName {
		return nil, errors.New("not an MSIX package: the content types are not the last file")
	}
	var blockMap, manifest []byte
	for _, e := range a.entries {
		switch e.name {
		case signatureName:
			return nil, errors.New("the signature is not the last file")
		case blockMapName:
			if blockMap, err = a.read(data, e); err != nil {
				return nil, err
			}
		case manifestName:
			if manifest, err = a.read(data, e); err != nil {
				return nil, err
			}
		}
	}
	if blockMap == nil || manifest == nil {
		return nil, errors.New("not an MSIX package: no block map or manifest")
	}
	if err := checkPublisher(manifest, signer); err != nil {
		return nil, err
	}
	last := a.entries[n-1]
	types, err := a.read(data, last)
	if err != nil {
		return nil, err
	}
	if !bytes.Contains(types, []byte("/"+signatureName)) {
		end := bytes.LastIndex(types, []byte("</Types>"))
		if end < 0 {
			return nil, errors.New("bad content types")
		}
		types = append(append(bytes.Clone(types[:end]), signatureOverride...), "</Types>"...)
	}

	w := &writer{now: time.Now(), entries: slices.Clip(a.entries[:n-1])}
	w.out.Write(data[:last.offset])
	if _, err := w.add(last.name, types, true); err != nil {
		return nil, err
	}

	digests := []byte("APPX")
	for _, part := range []struct {
		tag  string
		data []byte
	}{
		{"AXPC", w.out.Bytes()},
		{"AXCD", directory(w.entries, uint64(w.out.Len()))},
		{"AXCT", types},
		{"AXBM", blockMap},
	} {
		sum := sha256.Sum256(part.data)
		digests = append(append(digests, part.tag...), sum[:]...)
	}
	p7, err := authenticode.SignAppx(digests, signer)
	if err != nil {
		return nil, err
	}
	if _, err := w.add(signatureName, append([]byte("PKCX"), p7...), true); err != nil {
		return nil, err
	}
	return append(w.out.Bytes(), directory(w.entries, uint64(w.out.Len()))...), nil
}
