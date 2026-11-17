package apple

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const notaryAPI = "https://appstoreconnect.apple.com/notary/v2"

type NotaryKey struct {
	KeyID  string
	Issuer string
	Key    *ecdsa.PrivateKey
}

func LoadNotaryKey(path, keyID, issuer string) (*NotaryKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("%s: not a PEM key", path)
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	ec, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%s: App Store Connect keys are ECDSA P-256", path)
	}
	if keyID == "" {
		keyID = strings.TrimSuffix(strings.TrimPrefix(filepath.Base(path), "AuthKey_"), ".p8")
	}
	return &NotaryKey{KeyID: keyID, Issuer: strings.TrimSpace(issuer), Key: ec}, nil
}

func (k *NotaryKey) token() (string, error) {
	enc := func(v any) string {
		b, _ := json.Marshal(v)
		return base64.RawURLEncoding.EncodeToString(b)
	}
	now := time.Now()
	head := enc(map[string]string{"alg": "ES256", "kid": k.KeyID, "typ": "JWT"})
	claims := enc(map[string]any{"iss": k.Issuer, "iat": now.Unix(), "exp": now.Add(15 * time.Minute).Unix(), "aud": "appstoreconnect-v1", "scope": []string{"/notary/v2"}})
	sum := sha256.Sum256([]byte(head + "." + claims))
	r, s, err := ecdsa.Sign(rand.Reader, k.Key, sum[:])
	if err != nil {
		return "", err
	}
	sig := make([]byte, 64)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:])
	return head + "." + claims + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

func (k *NotaryKey) call(method, path string, body any, out any) error {
	tok, err := k.token()
	if err != nil {
		return err
	}
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, notaryAPI+path, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("notary API %s %s: HTTP %d: %s", method, path, resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if out != nil {
		return json.Unmarshal(data, out)
	}
	return nil
}

type submission struct {
	Data struct {
		ID         string `json:"id"`
		Attributes struct {
			AWSAccessKeyID     string `json:"awsAccessKeyId"`
			AWSSecretAccessKey string `json:"awsSecretAccessKey"`
			AWSSessionToken    string `json:"awsSessionToken"`
			Bucket             string `json:"bucket"`
			Object             string `json:"object"`
			Status             string `json:"status"`
		} `json:"attributes"`
	} `json:"data"`
}

func Notarize(k *NotaryKey, name string, data []byte, log io.Writer) (string, error) {
	sum := sha256.Sum256(data)
	var sub submission
	if err := k.call("POST", "/submissions", map[string]any{"submissionName": name, "sha256": hex.EncodeToString(sum[:])}, &sub); err != nil {
		return "", err
	}
	a := sub.Data.Attributes
	if err := s3Put(a.Bucket, a.Object, a.AWSAccessKeyID, a.AWSSecretAccessKey, a.AWSSessionToken, data); err != nil {
		return "", err
	}
	fmt.Fprintf(log, "notary: submission %s uploaded, waiting\n", sub.Data.ID)

	for {
		time.Sleep(15 * time.Second)
		var st submission
		if err := k.call("GET", "/submissions/"+sub.Data.ID, nil, &st); err != nil {
			return sub.Data.ID, err
		}
		switch st.Data.Attributes.Status {
		case "In Progress":
			continue
		case "Accepted":
			return sub.Data.ID, nil
		}
		var logs struct {
			Data struct {
				Attributes struct {
					URL string `json:"developerLogUrl"`
				} `json:"attributes"`
			} `json:"data"`
		}
		msg := st.Data.Attributes.Status
		if err := k.call("GET", "/submissions/"+sub.Data.ID+"/logs", nil, &logs); err == nil && logs.Data.Attributes.URL != "" {
			if resp, err := http.Get(logs.Data.Attributes.URL); err == nil {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				msg += "\n" + string(body)
			}
		}
		return sub.Data.ID, errors.New("notarization " + msg)
	}
}

func s3Put(bucket, object, accessKey, secretKey, sessionToken string, data []byte) error {
	const region, service = "us-west-2", "s3"
	host := bucket + ".s3." + region + ".amazonaws.com"
	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	payloadHash := hex.EncodeToString(sha256Sum(data))

	path := "/" + object
	headers := map[string]string{
		"host":                 host,
		"x-amz-content-sha256": payloadHash,
		"x-amz-date":           amzDate,
		"x-amz-security-token": sessionToken,
	}
	signedHeaders := "host;x-amz-content-sha256;x-amz-date;x-amz-security-token"
	canonical := "PUT\n" + uriEncode(path) + "\n\n"
	for _, h := range strings.Split(signedHeaders, ";") {
		canonical += h + ":" + headers[h] + "\n"
	}
	canonical += "\n" + signedHeaders + "\n" + payloadHash
	scope := dateStamp + "/" + region + "/" + service + "/aws4_request"
	toSign := "AWS4-HMAC-SHA256\n" + amzDate + "\n" + scope + "\n" + hex.EncodeToString(sha256Sum([]byte(canonical)))
	kDate := hmacSHA256([]byte("AWS4"+secretKey), dateStamp)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "aws4_request")
	signature := hex.EncodeToString(hmacSHA256(kSigning, toSign))

	req, err := http.NewRequest("PUT", "https://"+host+path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.ContentLength = int64(len(data))
	for k, v := range headers {
		if k != "host" {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Authorization", fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s", accessKey, scope, signedHeaders, signature))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload: HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func uriEncode(path string) string {
	var b strings.Builder
	for _, c := range []byte(path) {
		if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' || c == '_' || c == '.' || c == '~' || c == '/' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

func sha256Sum(b []byte) []byte {
	h := sha256.Sum256(b)
	return h[:]
}

func hmacSHA256(key []byte, msg string) []byte {
	m := hmac.New(sha256.New, key)
	m.Write([]byte(msg))
	return m.Sum(nil)
}
