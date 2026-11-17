package tsa

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"fmt"
	"io"
	"math/big"
	"net/http"
)

var oidSHA256 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}

type algorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
}

type messageImprint struct {
	HashAlgorithm algorithmIdentifier
	HashedMessage []byte
}

type timeStampReq struct {
	Version        int
	MessageImprint messageImprint
	Nonce          *big.Int
	CertReq        bool
}

type pkiStatusInfo struct {
	Status int
	Rest   asn1.RawValue `asn1:"optional"`
}

type timeStampResp struct {
	Status         pkiStatusInfo
	TimeStampToken asn1.RawValue `asn1:"optional"`
}

func Token(url string, signature []byte) ([]byte, error) {
	sum := sha256.Sum256(signature)
	nonce, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 64))
	if err != nil {
		return nil, err
	}
	req, err := asn1.Marshal(timeStampReq{
		Version:        1,
		MessageImprint: messageImprint{algorithmIdentifier{Algorithm: oidSHA256, Parameters: asn1.RawValue{Tag: asn1.TagNull}}, sum[:]},
		Nonce:          nonce,
		CertReq:        true,
	})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/timestamp-query", bytes.NewReader(req))
	if err != nil {
		return nil, fmt.Errorf("timestamp: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("timestamp: HTTP %d", resp.StatusCode)
	}
	var tsr timeStampResp
	if _, err := asn1.Unmarshal(body, &tsr); err != nil {
		return nil, fmt.Errorf("timestamp: %w", err)
	}
	if tsr.Status.Status > 1 || len(tsr.TimeStampToken.FullBytes) == 0 {
		return nil, fmt.Errorf("timestamp: status %d", tsr.Status.Status)
	}
	return tsr.TimeStampToken.FullBytes, nil
}
