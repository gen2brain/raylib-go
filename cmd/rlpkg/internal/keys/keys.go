package keys

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"software.sslmate.com/src/go-pkcs12"
)

const PasswordEnv = "RLPKG_P12_PASSWORD"

type Identity struct {
	Key   crypto.Signer
	Cert  *x509.Certificate
	Extra []*x509.Certificate
}

func (id *Identity) Certificates() []*x509.Certificate {
	return append([]*x509.Certificate{id.Cert}, id.Extra...)
}

func Load(path string) (*Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lower := strings.ToLower(path)
	if strings.HasSuffix(lower, ".p12") || strings.HasSuffix(lower, ".pfx") {
		key, cert, cas, err := pkcs12.DecodeChain(data, os.Getenv(PasswordEnv))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		signer, ok := key.(crypto.Signer)
		if !ok {
			return nil, fmt.Errorf("%s: unsupported key type", path)
		}
		return &Identity{Key: signer, Cert: cert, Extra: cas}, nil
	}
	id, err := ParsePEM(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return id, nil
}

func ParsePEM(data []byte) (*Identity, error) {
	id := &Identity{}
	for {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		var key any
		var err error
		switch block.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			if id.Cert == nil {
				id.Cert = cert
			} else {
				id.Extra = append(id.Extra, cert)
			}
			continue
		case "PRIVATE KEY":
			key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		case "RSA PRIVATE KEY":
			key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		case "EC PRIVATE KEY":
			key, err = x509.ParseECPrivateKey(block.Bytes)
		default:
			continue
		}
		if err != nil {
			return nil, err
		}
		id.Key, _ = key.(crypto.Signer)
	}
	if id.Key == nil || id.Cert == nil {
		return nil, fmt.Errorf("needs a PEM private key and certificate")
	}
	return id, nil
}
