package authenticode

import (
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/binary"
	"errors"
	"slices"
	"unicode/utf16"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/cms"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/tsa"
)

var (
	oidSpcIndirectData = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 4}
	oidSpcStatement    = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 11}
	oidSpcOpusInfo     = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 12}
	oidSpcPEImageData  = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 15}
	oidSpcIndividual   = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 21}
	oidSpcSipInfo      = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 30}
	oidMSCounterSign   = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 3, 3, 1}
)

const DefaultTimestampURL = "http://timestamp.digicert.com"

type Signer struct {
	Key          crypto.Signer
	Chain        []*x509.Certificate
	TimestampURL string
}

type spcAttributeTypeAndValue struct {
	Type  asn1.ObjectIdentifier
	Value asn1.RawValue
}

type spcPEImageData struct {
	Flags asn1.BitString
	File  asn1.RawValue
}

type digestInfo struct {
	Algorithm cms.AlgorithmIdentifier
	Digest    []byte
}

type spcIndirectDataContent struct {
	Data          spcAttributeTypeAndValue
	MessageDigest digestInfo
}

type spcSipInfo struct {
	Version            int
	UUID               []byte
	R1, R2, R3, R4, R5 int
}

var appxSipUUID = []byte{0x4B, 0xDF, 0xC5, 0x0A, 0x07, 0xCE, 0xE2, 0x4D, 0xB7, 0x6E, 0x23, 0xC8, 0x39, 0xA0, 0x9F, 0xD1}

type peLayout struct {
	checksumOff int
	certDirOff  int
	certOff     int
	certSize    int
	headersSize int
	sections    [][2]int
}

func layout(pe []byte) (*peLayout, error) {
	if len(pe) < 0x40 || pe[0] != 'M' || pe[1] != 'Z' {
		return nil, errors.New("not a PE file")
	}
	le := binary.LittleEndian
	off := int(le.Uint32(pe[0x3c:]))
	if off+24 > len(pe) || string(pe[off:off+4]) != "PE\x00\x00" {
		return nil, errors.New("bad PE header")
	}
	nsect := int(le.Uint16(pe[off+6:]))
	optSize := int(le.Uint16(pe[off+20:]))
	opt := off + 24
	if opt+optSize > len(pe) {
		return nil, errors.New("truncated optional header")
	}
	dirs := opt + 96
	if le.Uint16(pe[opt:]) == 0x20b {
		dirs = opt + 112
	}
	l := &peLayout{
		checksumOff: opt + 64,
		certDirOff:  dirs + 4*8,
		headersSize: int(le.Uint32(pe[opt+60:])),
	}
	l.certOff = int(le.Uint32(pe[l.certDirOff:]))
	l.certSize = int(le.Uint32(pe[l.certDirOff+4:]))
	sect := opt + optSize
	for i := 0; i < nsect; i++ {
		s := sect + 40*i
		size := int(le.Uint32(pe[s+16:]))
		ptr := int(le.Uint32(pe[s+20:]))
		if size > 0 {
			l.sections = append(l.sections, [2]int{ptr, size})
		}
	}
	slices.SortFunc(l.sections, func(a, b [2]int) int { return a[0] - b[0] })
	return l, nil
}

func imageHash(pe []byte, l *peLayout) ([]byte, error) {
	h := sha256.New()
	h.Write(pe[:l.checksumOff])
	h.Write(pe[l.checksumOff+4 : l.certDirOff])
	h.Write(pe[l.certDirOff+8 : l.headersSize])
	sum := l.headersSize
	for _, s := range l.sections {
		if s[0]+s[1] > len(pe) {
			return nil, errors.New("section beyond end of file")
		}
		h.Write(pe[s[0] : s[0]+s[1]])
		sum += s[1]
	}
	end := len(pe)
	if l.certSize > 0 {
		end = l.certOff
	}
	if end > sum {
		h.Write(pe[sum:end])
	}
	return h.Sum(nil), nil
}

func indirectData(hash []byte) ([]byte, error) {
	name := utf16.Encode([]rune("<<<Obsolete>>>"))
	bmp := make([]byte, 0, 2*len(name))
	for _, c := range name {
		bmp = binary.BigEndian.AppendUint16(bmp, c)
	}
	unicode, err := asn1.Marshal(asn1.RawValue{Class: asn1.ClassContextSpecific, Tag: 0, Bytes: bmp})
	if err != nil {
		return nil, err
	}
	image, err := asn1.Marshal(spcPEImageData{
		Flags: asn1.BitString{},
		File:  cms.Explicit(0, cms.Marshal(cms.Explicit(2, unicode))),
	})
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(spcIndirectDataContent{
		Data:          spcAttributeTypeAndValue{oidSpcPEImageData, asn1.RawValue{FullBytes: image}},
		MessageDigest: digestInfo{cms.SHA256, hash},
	})
}

func pkcs7(content []byte, s *Signer) ([]byte, error) {
	var inner asn1.RawValue
	rest, err := asn1.Unmarshal(content, &inner)
	if err != nil || len(rest) != 0 {
		return nil, errors.New("bad indirect data encoding")
	}
	digest := sha256.Sum256(inner.Bytes)

	o := cms.Options{
		ContentType: oidSpcIndirectData,
		Content:     content,
		Key:         s.Key,
		Chain:       s.Chain,
		SignedAttrs: []cms.Attribute{
			cms.NewAttribute(cms.OIDContentType, oidSpcIndirectData),
			cms.NewAttribute(oidSpcStatement, []asn1.ObjectIdentifier{oidSpcIndividual}),
			cms.RawAttribute(oidSpcOpusInfo, []byte{0x30, 0x00}),
			cms.NewAttribute(cms.OIDMessageDigest, digest[:]),
		},
	}
	if s.TimestampURL != "" {
		o.Unsigned = func(sig []byte) ([]cms.Attribute, error) {
			token, err := tsa.Token(s.TimestampURL, sig)
			if err != nil {
				return nil, err
			}
			return []cms.Attribute{cms.RawAttribute(oidMSCounterSign, token)}, nil
		}
	}
	return cms.Sign(o)
}

func checksum(pe []byte, checksumOff int) uint32 {
	var sum uint64
	for i := 0; i < len(pe); i += 4 {
		if i == checksumOff {
			continue
		}
		var word uint32
		for j := 0; j < 4 && i+j < len(pe); j++ {
			word |= uint32(pe[i+j]) << (8 * j)
		}
		sum += uint64(word)
		sum = (sum & 0xffffffff) + (sum >> 32)
	}
	sum = (sum & 0xffff) + (sum >> 16)
	sum += sum >> 16
	sum &= 0xffff
	return uint32(sum) + uint32(len(pe))
}

func SignAppx(digests []byte, s *Signer) ([]byte, error) {
	if s == nil || s.Key == nil || len(s.Chain) == 0 {
		return nil, errors.New("authenticode: a certificate is required")
	}
	sip, err := asn1.Marshal(spcSipInfo{Version: 0x01010000, UUID: appxSipUUID})
	if err != nil {
		return nil, err
	}
	content, err := asn1.Marshal(spcIndirectDataContent{
		Data:          spcAttributeTypeAndValue{oidSpcSipInfo, asn1.RawValue{FullBytes: sip}},
		MessageDigest: digestInfo{cms.SHA256, digests},
	})
	if err != nil {
		return nil, err
	}
	return pkcs7(content, s)
}

func Sign(pe []byte, s *Signer) ([]byte, error) {
	if s == nil || s.Key == nil || len(s.Chain) == 0 {
		return nil, errors.New("authenticode: a certificate is required")
	}
	l, err := layout(pe)
	if err != nil {
		return nil, err
	}
	if l.certSize > 0 && l.certOff+l.certSize == len(pe) {
		pe = pe[:l.certOff]
	} else if l.certSize > 0 {
		return nil, errors.New("existing certificate table is not at the end of the file")
	}
	out := append([]byte(nil), pe...)
	for len(out)%8 != 0 {
		out = append(out, 0)
	}
	le := binary.LittleEndian
	le.PutUint32(out[l.certDirOff:], 0)
	le.PutUint32(out[l.certDirOff+4:], 0)
	l.certSize = 0

	hash, err := imageHash(out, l)
	if err != nil {
		return nil, err
	}
	content, err := indirectData(hash)
	if err != nil {
		return nil, err
	}
	p7, err := pkcs7(content, s)
	if err != nil {
		return nil, err
	}

	certLen := (8 + len(p7) + 7) &^ 7
	cert := make([]byte, certLen)
	le.PutUint32(cert, uint32(certLen))
	le.PutUint16(cert[4:], 0x0200)
	le.PutUint16(cert[6:], 0x0002)
	copy(cert[8:], p7)

	le.PutUint32(out[l.certDirOff:], uint32(len(out)))
	le.PutUint32(out[l.certDirOff+4:], uint32(certLen))
	out = append(out, cert...)
	le.PutUint32(out[l.checksumOff:], checksum(out, l.checksumOff))
	return out, nil
}
