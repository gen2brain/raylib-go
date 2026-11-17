package msix

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"slices"
	"strings"
)

var attributeNames = map[string]string{
	"2.5.4.3":                    "CN",
	"2.5.4.4":                    "SN",
	"2.5.4.5":                    "SERIALNUMBER",
	"2.5.4.6":                    "C",
	"2.5.4.7":                    "L",
	"2.5.4.8":                    "S",
	"2.5.4.9":                    "STREET",
	"2.5.4.10":                   "O",
	"2.5.4.11":                   "OU",
	"2.5.4.12":                   "T",
	"2.5.4.42":                   "G",
	"2.5.4.43":                   "I",
	"0.9.2342.19200300.100.1.25": "DC",
	"1.2.840.113549.1.9.1":       "E",
}

func Publisher(cert *x509.Certificate) (string, error) {
	var rdns pkix.RDNSequence
	if rest, err := asn1.Unmarshal(cert.RawSubject, &rdns); err != nil || len(rest) != 0 {
		return "", fmt.Errorf("certificate subject: %w", err)
	}
	var parts []string
	for _, rdn := range slices.Backward(rdns) {
		var pairs []string
		for _, atv := range rdn {
			name, ok := attributeNames[atv.Type.String()]
			if !ok {
				name = "OID." + atv.Type.String()
			}
			pairs = append(pairs, name+"="+quote(fmt.Sprint(atv.Value)))
		}
		parts = append(parts, strings.Join(pairs, " + "))
	}
	return strings.Join(parts, ", "), nil
}

func quote(value string) string {
	if strings.ContainsAny(value, ",+=\"\n<>#;") || value != strings.TrimSpace(value) {
		return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
	}
	return value
}

func CommonName(name string) string { return "CN=" + quote(name) }

func ValidName(name string) bool {
	if len(name) < 3 || len(name) > 50 {
		return false
	}
	for _, c := range name {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '.' || c == '-') {
			return false
		}
	}
	return true
}

func canonical(dn string) string {
	var b strings.Builder
	quoted, separator := false, false
	for _, c := range dn {
		switch {
		case c == '"':
			quoted = !quoted
		case !quoted && separator && c == ' ':
			continue
		}
		separator = !quoted && (c == ',' || c == '+')
		b.WriteRune(c)
	}
	return b.String()
}
