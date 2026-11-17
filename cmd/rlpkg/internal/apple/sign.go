package apple

import (
	"bytes"
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"debug/macho"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	fat "github.com/gen2brain/raylib-go/cmd/rlpkg/internal/macho"
)

type Signer struct {
	Key       crypto.Signer
	Chain     []*x509.Certificate
	Timestamp bool
}

type MachOOptions struct {
	Identifier      string
	Entitlements    []byte
	HardenedRuntime bool
	InfoPlist       []byte
	ResourcesHash   []byte
}

func SignMachO(data []byte, s *Signer, o MachOOptions) ([]byte, [][]byte, error) {
	if len(data) >= 4 && data[0] == 0xca && data[1] == 0xfe && data[2] == 0xba && data[3] == 0xbe {
		ff, err := macho.NewFatFile(bytes.NewReader(data))
		if err != nil {
			return nil, nil, err
		}
		var slices, cdhashes [][]byte
		for _, a := range ff.Arches {
			signed, h, err := signSlice(data[a.Offset:a.Offset+a.Size], s, o)
			if err != nil {
				return nil, nil, err
			}
			slices = append(slices, signed)
			cdhashes = append(cdhashes, h)
		}
		out, err := fat.Universal(slices)
		return out, cdhashes, err
	}
	signed, h, err := signSlice(data, s, o)
	return signed, [][]byte{h}, err
}

const (
	lcSegment64      = 0x19
	lcCodeSignature  = 0x1d
	lcBuildVersion   = 0x32
	linkeditPageSize = 0x4000
)

func signSlice(data []byte, s *Signer, o MachOOptions) ([]byte, []byte, error) {
	f, err := macho.NewFile(bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}
	if f.Magic != macho.Magic64 {
		return nil, nil, errors.New("only 64-bit Mach-O files are supported")
	}
	bo := f.ByteOrder

	const hdrSize = 32
	sizeofcmds := int(bo.Uint32(data[20:]))
	csCmd, leCmd := -1, -1
	var csDataOff, csDataSize uint32
	var leFileOff, leFileSize uint64
	var textOff, textSize uint64
	firstSection := uint32(len(data))
	var sdk uint32
	off := hdrSize
	for i := 0; i < int(f.Ncmd); i++ {
		cmd := bo.Uint32(data[off:])
		size := int(bo.Uint32(data[off+4:]))
		switch cmd {
		case lcCodeSignature:
			csCmd = off
			csDataOff = bo.Uint32(data[off+8:])
			csDataSize = bo.Uint32(data[off+12:])
		case lcBuildVersion:
			sdk = bo.Uint32(data[off+16:])
		case lcSegment64:
			name := strings.TrimRight(string(data[off+8:off+24]), "\x00")
			switch name {
			case "__LINKEDIT":
				leCmd = off
				leFileOff = bo.Uint64(data[off+40:])
				leFileSize = bo.Uint64(data[off+48:])
			case "__TEXT":
				textOff = bo.Uint64(data[off+40:])
				textSize = bo.Uint64(data[off+48:])
				nsects := int(bo.Uint32(data[off+64:]))
				for j := 0; j < nsects; j++ {
					if so := bo.Uint32(data[off+72+80*j+48:]); so > 0 && so < firstSection {
						firstSection = so
					}
				}
			}
		}
		off += size
	}
	if leCmd < 0 {
		return nil, nil, errors.New("no __LINKEDIT segment")
	}

	if csCmd >= 0 {
		if uint64(csDataOff)+uint64(csDataSize) != uint64(len(data)) {
			return nil, nil, errors.New("existing code signature is not at the end of the file")
		}
		data = append([]byte(nil), data[:csDataOff]...)
	} else {
		end := hdrSize + sizeofcmds
		if uint32(end+16) > firstSection {
			return nil, nil, errors.New("no room in the Mach-O header for LC_CODE_SIGNATURE")
		}
		data = append([]byte(nil), data[:leFileOff+leFileSize]...)
		csDataOff = uint32((len(data) + 15) &^ 15)
		csCmd = end
		bo.PutUint32(data[csCmd:], lcCodeSignature)
		bo.PutUint32(data[csCmd+4:], 16)
		bo.PutUint32(data[16:], f.Ncmd+1)
		bo.PutUint32(data[20:], uint32(sizeofcmds+16))
	}
	for uint32(len(data)) < csDataOff {
		data = append(data, 0)
	}

	reserved := uint32(len(data)) / 4096 * 40
	reserved = (reserved + 16384 + 15) &^ 15
	if s != nil && s.Key != nil {
		for _, c := range s.Chain {
			reserved += uint32(len(c.Raw))
		}
		reserved += 8192
	}
	bo.PutUint32(data[csCmd+8:], csDataOff)
	bo.PutUint32(data[csCmd+12:], reserved)
	fileSize := uint64(csDataOff) + uint64(reserved) - leFileOff
	bo.PutUint64(data[leCmd+48:], fileSize)
	if vm := bo.Uint64(data[leCmd+32:]); vm < fileSize {
		bo.PutUint64(data[leCmd+32:], (fileSize+linkeditPageSize-1)&^(linkeditPageSize-1))
	}

	blob, cdBytes, err := buildSignature(codeInfo{code: data, textOff: textOff, textSize: textSize, exec: f.Type == macho.TypeExec, sdk: sdk}, s, o)
	if err != nil {
		return nil, nil, err
	}
	if uint32(len(blob)) > reserved {
		return nil, nil, fmt.Errorf("code signature needs %d bytes, %d reserved", len(blob), reserved)
	}
	out := append(data, blob...)
	out = append(out, make([]byte, int(reserved)-len(blob))...)
	sum := sha256.Sum256(cdBytes)
	return out, sum[:], nil
}

func Requirement(identifier string, s *Signer, cdhashes [][]byte) string {
	if s == nil || s.Key == nil {
		var parts []string
		for _, h := range cdhashes {
			parts = append(parts, fmt.Sprintf("cdhash H\"%x\"", h[:20]))
		}
		return strings.Join(parts, " or ")
	}
	leaf := s.Chain[0]
	for _, ext := range leaf.Extensions {
		if ext.Id.String() == "1.2.840.113635.100.6.1.13" {
			return fmt.Sprintf("identifier %q and anchor apple generic and certificate 1[field.1.2.840.113635.100.6.2.6] /* exists */ and certificate leaf[field.1.2.840.113635.100.6.1.13] /* exists */ and certificate leaf[subject.OU] = %q", identifier, TeamID(leaf))
		}
	}
	return fmt.Sprintf("identifier %q and anchor apple generic and certificate leaf[subject.CN] = %q and certificate 1[field.1.2.840.113635.100.6.2.1] /* exists */", identifier, leaf.Subject.CommonName)
}

func IsMachO(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	switch string(data[:4]) {
	case "\xca\xfe\xba\xbe", "\xcf\xfa\xed\xfe", "\xce\xfa\xed\xfe":
		return true
	}
	return false
}

func CDHash(data []byte) ([]byte, error) {
	if len(data) >= 4 && data[0] == 0xca && data[1] == 0xfe && data[2] == 0xba && data[3] == 0xbe {
		ff, err := macho.NewFatFile(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		a := ff.Arches[0]
		data = data[a.Offset : a.Offset+a.Size]
	}
	f, err := macho.NewFile(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	bo := f.ByteOrder
	off := 32
	for i := 0; i < int(f.Ncmd); i++ {
		cmd := bo.Uint32(data[off:])
		size := int(bo.Uint32(data[off+4:]))
		if cmd == lcCodeSignature {
			sb := data[bo.Uint32(data[off+8:]):]
			count := int(binary.BigEndian.Uint32(sb[8:]))
			for j := 0; j < count; j++ {
				if binary.BigEndian.Uint32(sb[12+8*j:]) == 0 {
					o := binary.BigEndian.Uint32(sb[16+8*j:])
					l := binary.BigEndian.Uint32(sb[o+4:])
					sum := sha256.Sum256(sb[o : o+l])
					return sum[:20], nil
				}
			}
		}
		off += size
	}
	return nil, errors.New("no code signature")
}
