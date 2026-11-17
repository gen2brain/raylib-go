package apple

import (
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"time"

	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/cms"
	"github.com/gen2brain/raylib-go/cmd/rlpkg/internal/tsa"
)

var (
	oidTimestampToken   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 14}
	oidAppleCDHashPlist = asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 9, 1}
	oidAppleCDHash      = asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 9, 2}
)

const TimestampURL = "http://timestamp.apple.com/ts01"

type cdHashAlgo struct {
	Algorithm asn1.ObjectIdentifier
	Hash      []byte
}

func SignCMS(cd []byte, key crypto.Signer, chain []*x509.Certificate, timestamp bool) ([]byte, error) {
	digest := sha256.Sum256(cd)
	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>cdhashes</key>
	<array>
		<data>
		%s
		</data>
	</array>
</dict>
</plist>
`, base64Std(digest[:20]))

	o := cms.Options{
		Key:   key,
		Chain: chain,
		SignedAttrs: []cms.Attribute{
			cms.NewAttribute(cms.OIDContentType, cms.OIDData),
			cms.NewAttribute(cms.OIDSigningTime, time.Now().UTC().Truncate(time.Second)),
			cms.NewAttribute(oidAppleCDHashPlist, []byte(plist)),
			cms.NewAttribute(cms.OIDMessageDigest, digest[:]),
			cms.NewAttribute(oidAppleCDHash, cdHashAlgo{cms.OIDSHA256, digest[:]}),
		},
	}
	if timestamp {
		o.Unsigned = func(sig []byte) ([]cms.Attribute, error) {
			token, err := tsa.Token(TimestampURL, sig)
			if err != nil {
				return nil, err
			}
			return []cms.Attribute{cms.RawAttribute(oidTimestampToken, token)}, nil
		}
	}
	return cms.Sign(o)
}
