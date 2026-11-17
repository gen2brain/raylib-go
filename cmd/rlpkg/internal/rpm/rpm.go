package rpm

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"path"
	"slices"
	"time"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/pgp"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/pkgtree"
)

type Package struct {
	Name        string
	Version     string
	Release     string
	Arch        string
	Summary     string
	Description string
	License     string
	Vendor      string
	URL         string
	Files       []pkgtree.File
	Sign        pgp.Signer
}

const (
	typeChar        = 1
	typeInt8        = 2
	typeInt16       = 3
	typeInt32       = 4
	typeInt64       = 5
	typeString      = 6
	typeBin         = 7
	typeStringArray = 8
	typeI18NString  = 9

	tagHeaderSignatures = 62
	tagHeaderImmutable  = 63
	tagHeaderI18NTable  = 100
	tagSigSize          = 1000
	tagSigPGP           = 1002
	tagSigMD5           = 1004
	tagSigGPG           = 1005
	tagSigReserved      = 1008
	tagSigDSA           = 267
	tagSigRSA           = 268
	tagSigSHA1          = 269
	tagSigSHA256        = 273

	tagName             = 1000
	tagVersion          = 1001
	tagRelease          = 1002
	tagSummary          = 1004
	tagDescription      = 1005
	tagBuildTime        = 1006
	tagBuildHost        = 1007
	tagSize             = 1009
	tagVendor           = 1011
	tagLicense          = 1014
	tagGroup            = 1016
	tagURL              = 1020
	tagOS               = 1021
	tagArch             = 1022
	tagFileSizes        = 1028
	tagFileModes        = 1030
	tagFileRdevs        = 1033
	tagFileMtimes       = 1034
	tagFileDigests      = 1035
	tagFileLinktos      = 1036
	tagFileFlags        = 1037
	tagFileUsername     = 1039
	tagFileGroupname    = 1040
	tagSourceRPM        = 1044
	tagFileVerifyFlags  = 1045
	tagProvideName      = 1047
	tagRequireFlags     = 1048
	tagRequireName      = 1049
	tagRequireVersion   = 1050
	tagFileDevices      = 1095
	tagFileInodes       = 1096
	tagFileLangs        = 1097
	tagProvideFlags     = 1112
	tagProvideVersion   = 1113
	tagDirIndexes       = 1116
	tagBasenames        = 1117
	tagDirnames         = 1118
	tagPayloadFormat    = 1124
	tagPayloadCompress  = 1125
	tagPayloadFlags     = 1126
	tagFileDigestAlgo   = 5011
	tagRPMFormat        = 5114
	tagPayloadDigest    = 5092
	tagPayloadDigestAlg = 5093
	tagPayloadDigestAlt = 5097

	senseLess   = 2
	senseEqual  = 8
	senseRPMLib = 1 << 24
)

type entry struct {
	tag   int32
	typ   int32
	count int32
	data  []byte
}

func str(tag int32, s string) entry { return entry{tag, typeString, 1, append([]byte(s), 0)} }

func i18n(tag int32, s string) entry { return entry{tag, typeI18NString, 1, append([]byte(s), 0)} }

func strs(tag int32, list []string) entry {
	var b []byte
	for _, s := range list {
		b = append(append(b, s...), 0)
	}
	return entry{tag, typeStringArray, int32(len(list)), b}
}

func int32s(tag int32, list []int32) entry {
	var b []byte
	for _, v := range list {
		b = binary.BigEndian.AppendUint32(b, uint32(v))
	}
	return entry{tag, typeInt32, int32(len(list)), b}
}

func int16s(tag int32, list []int16) entry {
	var b []byte
	for _, v := range list {
		b = binary.BigEndian.AppendUint16(b, uint16(v))
	}
	return entry{tag, typeInt16, int32(len(list)), b}
}

func bin(tag int32, b []byte) entry { return entry{tag, typeBin, int32(len(b)), b} }

func align(typ int32) int {
	switch typ {
	case typeInt16:
		return 2
	case typeInt32:
		return 4
	case typeInt64:
		return 8
	}
	return 1
}

func header(regionTag int32, entries []entry) []byte {
	slices.SortFunc(entries, func(a, b entry) int { return int(a.tag - b.tag) })

	var data bytes.Buffer
	var index bytes.Buffer
	for _, e := range entries {
		for data.Len()%align(e.typ) != 0 {
			data.WriteByte(0)
		}
		index.Write(binary.BigEndian.AppendUint32(nil, uint32(e.tag)))
		index.Write(binary.BigEndian.AppendUint32(nil, uint32(e.typ)))
		index.Write(binary.BigEndian.AppendUint32(nil, uint32(data.Len())))
		index.Write(binary.BigEndian.AppendUint32(nil, uint32(e.count)))
		data.Write(e.data)
	}

	nindex := int32(len(entries) + 1)
	trailer := binary.BigEndian.AppendUint32(nil, uint32(regionTag))
	trailer = binary.BigEndian.AppendUint32(trailer, typeBin)
	trailer = binary.BigEndian.AppendUint32(trailer, uint32(-nindex*16))
	trailer = binary.BigEndian.AppendUint32(trailer, 16)
	region := binary.BigEndian.AppendUint32(nil, uint32(regionTag))
	region = binary.BigEndian.AppendUint32(region, typeBin)
	region = binary.BigEndian.AppendUint32(region, uint32(data.Len()))
	region = binary.BigEndian.AppendUint32(region, 16)
	data.Write(trailer)

	var out bytes.Buffer
	out.Write([]byte{0x8e, 0xad, 0xe8, 0x01, 0, 0, 0, 0})
	out.Write(binary.BigEndian.AppendUint32(nil, uint32(nindex)))
	out.Write(binary.BigEndian.AppendUint32(nil, uint32(data.Len())))
	out.Write(region)
	out.Write(index.Bytes())
	out.Write(data.Bytes())
	return out.Bytes()
}

func Write(p Package) ([]byte, error) {
	now := int32(time.Now().Unix())
	files := p.Files

	var dirnames []string
	dirIndex := map[string]int32{}
	var basenames, digests, linktos, users, groups, langs []string
	var sizes, mtimes, flags, verify, devices, inodes, dirIndexes []int32
	var modes, rdevs []int16
	var total int32

	var cpio bytes.Buffer
	for i, f := range files {
		dir := path.Dir(f.Path) + "/"
		if dir == "//" {
			dir = "/"
		}
		if _, ok := dirIndex[dir]; !ok {
			dirIndex[dir] = int32(len(dirnames))
			dirnames = append(dirnames, dir)
		}
		dirIndexes = append(dirIndexes, dirIndex[dir])
		basenames = append(basenames, path.Base(f.Path))
		users = append(users, "root")
		groups = append(groups, "root")
		langs = append(langs, "")
		mtimes = append(mtimes, now)
		flags = append(flags, 0)
		verify = append(verify, -1)
		devices = append(devices, 1)
		inodes = append(inodes, int32(i+1))
		rdevs = append(rdevs, 0)

		mode := f.Mode & 0o7777
		var body []byte
		switch {
		case f.Dir:
			mode |= 0o40000
			sizes = append(sizes, 4096)
			digests = append(digests, "")
			linktos = append(linktos, "")
		case f.Link != "":
			mode |= 0o120000
			body = []byte(f.Link)
			sizes = append(sizes, int32(len(body)))
			digests = append(digests, "")
			linktos = append(linktos, f.Link)
		default:
			mode |= 0o100000
			body = f.Data
			sizes = append(sizes, int32(len(body)))
			digests = append(digests, fmt.Sprintf("%x", sha256.Sum256(body)))
			linktos = append(linktos, "")
			total += int32(len(body))
		}
		modes = append(modes, int16(mode))
		cpioEntry(&cpio, "."+f.Path, mode, int32(i+1), now, body)
	}
	cpioEntry(&cpio, "TRAILER!!!", 0, 0, 0, nil)
	for cpio.Len()%512 != 0 {
		cpio.WriteByte(0)
	}

	var payload bytes.Buffer
	gw, err := gzip.NewWriterLevel(&payload, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	gw.Write(cpio.Bytes())
	if err := gw.Close(); err != nil {
		return nil, err
	}

	evr := p.Version + "-" + p.Release
	hdr := header(tagHeaderImmutable, []entry{
		strs(tagHeaderI18NTable, []string{"C"}),
		str(tagName, p.Name),
		str(tagVersion, p.Version),
		str(tagRelease, p.Release),
		i18n(tagSummary, p.Summary),
		i18n(tagDescription, p.Description),
		int32s(tagBuildTime, []int32{now}),
		str(tagBuildHost, "rlpkg"),
		int32s(tagSize, []int32{total}),
		str(tagVendor, p.Vendor),
		str(tagLicense, p.License),
		i18n(tagGroup, "Applications"),
		str(tagURL, p.URL),
		str(tagOS, "linux"),
		str(tagArch, p.Arch),
		int32s(tagFileSizes, sizes),
		int16s(tagFileModes, modes),
		int16s(tagFileRdevs, rdevs),
		int32s(tagFileMtimes, mtimes),
		strs(tagFileDigests, digests),
		strs(tagFileLinktos, linktos),
		int32s(tagFileFlags, flags),
		strs(tagFileUsername, users),
		strs(tagFileGroupname, groups),
		str(tagSourceRPM, p.Name+"-"+evr+".src.rpm"),
		int32s(tagFileVerifyFlags, verify),
		strs(tagProvideName, []string{p.Name}),
		int32s(tagProvideFlags, []int32{senseEqual}),
		strs(tagProvideVersion, []string{evr}),
		strs(tagRequireName, []string{"rpmlib(CompressedFileNames)", "rpmlib(FileDigests)", "rpmlib(PayloadFilesHavePrefix)"}),
		int32s(tagRequireFlags, []int32{senseRPMLib | senseLess | senseEqual, senseRPMLib | senseLess | senseEqual, senseRPMLib | senseLess | senseEqual}),
		strs(tagRequireVersion, []string{"3.0.4-1", "4.6.0-1", "4.0-1"}),
		int32s(tagFileDevices, devices),
		int32s(tagFileInodes, inodes),
		strs(tagFileLangs, langs),
		int32s(tagDirIndexes, dirIndexes),
		strs(tagBasenames, basenames),
		strs(tagDirnames, dirnames),
		str(tagPayloadFormat, "cpio"),
		str(tagPayloadCompress, "gzip"),
		str(tagPayloadFlags, "9"),
		int32s(tagFileDigestAlgo, []int32{8}),
		strs(tagPayloadDigest, []string{fmt.Sprintf("%x", sha256.Sum256(payload.Bytes()))}),
		int32s(tagPayloadDigestAlg, []int32{8}),
		strs(tagPayloadDigestAlt, []string{fmt.Sprintf("%x", sha256.Sum256(cpio.Bytes()))}),
	})

	sumMD5 := md5.New()
	sumMD5.Write(hdr)
	sumMD5.Write(payload.Bytes())
	sigs, err := signatures(hdr, payload.Bytes(), p.Sign)
	if err != nil {
		return nil, err
	}
	sig := header(tagHeaderSignatures, append(sigs,
		int32s(tagSigSize, []int32{int32(len(hdr) + payload.Len())}),
		bin(tagSigMD5, sumMD5.Sum(nil)),
		str(tagSigSHA1, fmt.Sprintf("%x", sha1.Sum(hdr))),
		str(tagSigSHA256, fmt.Sprintf("%x", sha256.Sum256(hdr))),
	))

	lead := []byte{0xed, 0xab, 0xee, 0xdb, 3, 0, 0, 0}
	lead = binary.BigEndian.AppendUint16(lead, 1)
	name := make([]byte, 66)
	copy(name, p.Name+"-"+evr)
	lead = append(lead, name...)
	lead = binary.BigEndian.AppendUint16(lead, 1)
	lead = binary.BigEndian.AppendUint16(lead, 5)
	lead = append(lead, make([]byte, 16)...)
	return assemble(lead, sig, hdr, payload.Bytes()), nil
}

func assemble(lead, sig, hdr, payload []byte) []byte {
	var out bytes.Buffer
	out.Write(lead)
	out.Write(sig)
	for out.Len()%8 != 0 {
		out.WriteByte(0)
	}
	out.Write(hdr)
	out.Write(payload)
	return out.Bytes()
}

func signatures(hdr, payload []byte, sign pgp.Signer) ([]entry, error) {
	if sign == nil {
		return nil, nil
	}
	hdrSig, err := sign(hdr)
	if err != nil {
		return nil, err
	}
	allSig, err := sign(slices.Concat(hdr, payload))
	if err != nil {
		return nil, err
	}
	rsa, err := pgp.IsRSA(hdrSig)
	if err != nil {
		return nil, err
	}
	if rsa {
		return []entry{bin(tagSigRSA, hdrSig), bin(tagSigPGP, allSig)}, nil
	}
	return []entry{bin(tagSigDSA, hdrSig), bin(tagSigGPG, allSig)}, nil
}

func Sign(data []byte, sign pgp.Signer) ([]byte, error) {
	const leadSize = 96
	if len(data) < leadSize || !bytes.Equal(data[:4], []byte{0xed, 0xab, 0xee, 0xdb}) {
		return nil, errors.New("not an RPM package")
	}
	sigEntries, sigLen, err := parseHeader(data[leadSize:])
	if err != nil {
		return nil, fmt.Errorf("signature header: %w", err)
	}
	rest := leadSize + (sigLen+7)/8*8
	if rest > len(data) {
		return nil, errors.New("truncated RPM package")
	}
	hdrEntries, hdrLen, err := parseHeader(data[rest:])
	if err != nil {
		return nil, fmt.Errorf("header: %w", err)
	}
	for _, e := range hdrEntries {
		if e.tag == tagRPMFormat && e.typ == typeInt32 && e.count == 1 && binary.BigEndian.Uint32(e.data) > 4 {
			return nil, errors.New("RPM v6 packages are not supported")
		}
	}
	hdr, payload := data[rest:rest+hdrLen], data[rest+hdrLen:]

	sigEntries = slices.DeleteFunc(sigEntries, func(e entry) bool {
		return slices.Contains([]int32{tagSigPGP, tagSigGPG, tagSigDSA, tagSigRSA, tagSigReserved}, e.tag)
	})
	sigs, err := signatures(hdr, payload, sign)
	if err != nil {
		return nil, err
	}
	return assemble(data[:leadSize], header(tagHeaderSignatures, append(sigs, sigEntries...)), hdr, payload), nil
}

func parseHeader(data []byte) ([]entry, int, error) {
	if len(data) < 16 || !bytes.Equal(data[:4], []byte{0x8e, 0xad, 0xe8, 0x01}) {
		return nil, 0, errors.New("bad magic")
	}
	nindex := int(binary.BigEndian.Uint32(data[8:]))
	hsize := int(binary.BigEndian.Uint32(data[12:]))
	if nindex < 0 || hsize < 0 || nindex > len(data)/16 || 16+nindex*16+hsize > len(data) {
		return nil, 0, errors.New("truncated")
	}
	store := data[16+nindex*16 : 16+nindex*16+hsize]

	var entries []entry
	for i := range nindex {
		index := data[16+i*16:]
		e := entry{
			tag:   int32(binary.BigEndian.Uint32(index)),
			typ:   int32(binary.BigEndian.Uint32(index[4:])),
			count: int32(binary.BigEndian.Uint32(index[12:])),
		}
		offset := int(binary.BigEndian.Uint32(index[8:]))
		if e.tag == tagHeaderSignatures || e.tag == tagHeaderImmutable {
			continue
		}
		if offset < 0 || offset > len(store) || e.count < 0 {
			return nil, 0, errors.New("bad entry")
		}
		var size int
		switch e.typ {
		case typeInt16:
			size = 2 * int(e.count)
		case typeInt32:
			size = 4 * int(e.count)
		case typeInt64:
			size = 8 * int(e.count)
		case typeBin, typeChar, typeInt8:
			size = int(e.count)
		case typeString, typeStringArray, typeI18NString:
			for range e.count {
				end := bytes.IndexByte(store[offset+size:], 0)
				if end < 0 {
					return nil, 0, errors.New("bad string")
				}
				size += end + 1
			}
		default:
			return nil, 0, fmt.Errorf("unsupported entry type %d", e.typ)
		}
		if size < 0 || offset+size > len(store) {
			return nil, 0, errors.New("bad entry")
		}
		e.data = bytes.Clone(store[offset : offset+size])
		entries = append(entries, e)
	}
	return entries, 16 + nindex*16 + hsize, nil
}

func cpioEntry(out *bytes.Buffer, name string, mode uint32, ino, mtime int32, body []byte) {
	nlink := 1
	if mode&0o40000 != 0 {
		nlink = 2
	}
	fmt.Fprintf(out, "070701%08x%08x%08x%08x%08x%08x%08x%08x%08x%08x%08x%08x%08x",
		ino, mode, 0, 0, nlink, mtime, len(body), 0, 0, 0, 0, len(name)+1, 0)
	out.WriteString(name)
	out.WriteByte(0)
	for out.Len()%4 != 0 {
		out.WriteByte(0)
	}
	out.Write(body)
	for out.Len()%4 != 0 {
		out.WriteByte(0)
	}
}
