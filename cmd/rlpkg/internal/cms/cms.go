package cms

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"math/big"
	"slices"
)

var (
	OIDData          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 1}
	OIDSignedData    = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 2}
	OIDContentType   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 3}
	OIDMessageDigest = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 4}
	OIDSigningTime   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 5}
	OIDSHA256        = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}

	oidRSA             = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	oidECDSAWithSHA256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 2}
)

type AlgorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
}

var SHA256 = AlgorithmIdentifier{Algorithm: OIDSHA256, Parameters: asn1.RawValue{Tag: asn1.TagNull}}

type Attribute struct {
	Type   asn1.ObjectIdentifier
	Values []asn1.RawValue `asn1:"set"`
}

func NewAttribute(oid asn1.ObjectIdentifier, v any) Attribute {
	return RawAttribute(oid, Marshal(v))
}

func RawAttribute(oid asn1.ObjectIdentifier, der []byte) Attribute {
	return Attribute{Type: oid, Values: []asn1.RawValue{{FullBytes: der}}}
}

func Marshal(v any) []byte {
	b, _ := asn1.Marshal(v)
	return b
}

func Explicit(tag int, der []byte) asn1.RawValue {
	return asn1.RawValue{Class: asn1.ClassContextSpecific, Tag: tag, IsCompound: true, Bytes: der}
}

type issuerAndSerial struct {
	Issuer       asn1.RawValue
	SerialNumber *big.Int
}

type signerInfo struct {
	Version            int
	IssuerAndSerial    issuerAndSerial
	DigestAlgorithm    AlgorithmIdentifier
	SignedAttrs        []Attribute `asn1:"tag:0,optional"`
	SignatureAlgorithm AlgorithmIdentifier
	Signature          []byte
	UnsignedAttrs      []Attribute `asn1:"tag:1,optional"`
}

type contentInfo struct {
	ContentType asn1.ObjectIdentifier
	Content     asn1.RawValue `asn1:"optional"`
}

type signedData struct {
	Version          int
	DigestAlgorithms []AlgorithmIdentifier `asn1:"set"`
	ContentInfo      contentInfo
	Certificates     asn1.RawValue `asn1:"optional,tag:0"`
	SignerInfos      []signerInfo  `asn1:"set"`
}

type Options struct {
	ContentType asn1.ObjectIdentifier
	Content     []byte
	Signed      []byte
	SignedAttrs []Attribute
	Unsigned    func(signature []byte) ([]Attribute, error)
	Key         crypto.Signer
	Chain       []*x509.Certificate
}

func Sign(o Options) ([]byte, error) {
	if o.Key == nil || len(o.Chain) == 0 {
		return nil, errors.New("cms: key and certificate required")
	}
	var sigAlg AlgorithmIdentifier
	switch o.Key.(type) {
	case *rsa.PrivateKey:
		sigAlg = AlgorithmIdentifier{Algorithm: oidRSA, Parameters: asn1.RawValue{Tag: asn1.TagNull}}
	case *ecdsa.PrivateKey:
		sigAlg = AlgorithmIdentifier{Algorithm: oidECDSAWithSHA256}
	default:
		return nil, errors.New("cms: signing key must be RSA or ECDSA")
	}

	signed := o.Signed
	attrs := o.SignedAttrs
	if len(attrs) > 0 {
		slices.SortFunc(attrs, func(a, b Attribute) int { return bytes.Compare(Marshal(a), Marshal(b)) })
		signed = Marshal(attrs)
		signed[0] = 0x31
	}
	sum := sha256.Sum256(signed)
	signature, err := o.Key.Sign(rand.Reader, sum[:], crypto.SHA256)
	if err != nil {
		return nil, err
	}
	var unsigned []Attribute
	if o.Unsigned != nil {
		if unsigned, err = o.Unsigned(signature); err != nil {
			return nil, err
		}
	}

	leaf := o.Chain[0]
	var certs []byte
	for _, c := range o.Chain {
		certs = append(certs, c.Raw...)
	}
	contentType := o.ContentType
	if contentType == nil {
		contentType = OIDData
	}
	var content asn1.RawValue
	if o.Content != nil {
		content = Explicit(0, o.Content)
	}
	sd, err := asn1.Marshal(signedData{
		Version:          1,
		DigestAlgorithms: []AlgorithmIdentifier{SHA256},
		ContentInfo:      contentInfo{ContentType: contentType, Content: content},
		Certificates:     Explicit(0, certs),
		SignerInfos: []signerInfo{{
			Version:            1,
			IssuerAndSerial:    issuerAndSerial{Issuer: asn1.RawValue{FullBytes: leaf.RawIssuer}, SerialNumber: leaf.SerialNumber},
			DigestAlgorithm:    SHA256,
			SignedAttrs:        attrs,
			SignatureAlgorithm: sigAlg,
			Signature:          signature,
			UnsignedAttrs:      unsigned,
		}},
	})
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(contentInfo{ContentType: OIDSignedData, Content: Explicit(0, sd)})
}
