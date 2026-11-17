package apple

import (
	"bytes"
	"crypto/x509"
	"embed"
	"encoding/pem"
	"errors"
)

//go:embed certs/*.pem
var certFiles embed.FS

func appleCAs() []*x509.Certificate {
	var cas []*x509.Certificate
	entries, _ := certFiles.ReadDir("certs")
	for _, e := range entries {
		data, _ := certFiles.ReadFile("certs/" + e.Name())
		for {
			var block *pem.Block
			block, data = pem.Decode(data)
			if block == nil {
				break
			}
			if c, err := x509.ParseCertificate(block.Bytes); err == nil {
				cas = append(cas, c)
			}
		}
	}
	return cas
}

func Chain(leaf *x509.Certificate, extra []*x509.Certificate) ([]*x509.Certificate, error) {
	pool := append(appleCAs(), extra...)
	chain := []*x509.Certificate{leaf}
	cur := leaf
	for i := 0; i < 8; i++ {
		if bytes.Equal(cur.RawSubject, cur.RawIssuer) {
			return chain, nil
		}
		var next *x509.Certificate
		for _, c := range pool {
			if bytes.Equal(c.RawSubject, cur.RawIssuer) && signedBy(cur, c) {
				next = c
				break
			}
		}
		if next == nil {
			return nil, errors.New("issuer of " + cur.Subject.CommonName + " not found: not an Apple-issued certificate")
		}
		chain = append(chain, next)
		cur = next
	}
	return nil, errors.New("certificate chain too long")
}

func signedBy(cert, issuer *x509.Certificate) bool {
	err := cert.CheckSignatureFrom(issuer)
	var insecure x509.InsecureAlgorithmError
	return err == nil || errors.As(err, &insecure)
}

func TeamID(cert *x509.Certificate) string {
	if len(cert.Subject.OrganizationalUnit) > 0 {
		return cert.Subject.OrganizationalUnit[0]
	}
	return ""
}
