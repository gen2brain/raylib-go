package keys

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadPEM(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "Leaf"}, NotBefore: time.Now(), NotAfter: time.Now().Add(time.Hour)}
	leaf, _ := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	tmpl.SerialNumber = big.NewInt(2)
	tmpl.Subject.CommonName = "CA"
	ca, _ := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	keyDER, _ := x509.MarshalPKCS8PrivateKey(key)

	path := filepath.Join(t.TempDir(), "id.pem")
	var data []byte
	data = append(data, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: leaf})...)
	data = append(data, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})...)
	data = append(data, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca})...)
	os.WriteFile(path, data, 0o600)

	id, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if id.Cert.Subject.CommonName != "Leaf" || len(id.Extra) != 1 || id.Extra[0].Subject.CommonName != "CA" {
		t.Errorf("chain order wrong: %s, %d extra", id.Cert.Subject.CommonName, len(id.Extra))
	}
	if _, ok := id.Key.(*ecdsa.PrivateKey); !ok {
		t.Errorf("key type %T", id.Key)
	}
	if len(id.Certificates()) != 2 {
		t.Error("Certificates() must be leaf + extras")
	}
	if _, err := Load(filepath.Join(t.TempDir(), "missing.pem")); err == nil {
		t.Error("missing file must fail")
	}
}
