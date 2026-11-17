package apple

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const ticketLookupURL = "https://api.apple-cloudkit.com/database/1/com.apple.gk.ticket-delivery/production/public/records/lookup"

func FetchTicket(cdhash []byte) ([]byte, error) {
	body, _ := json.Marshal(map[string]any{"records": []map[string]string{{"recordName": "2/2/" + hex.EncodeToString(cdhash[:20])}}})
	var last error
	for attempt := 0; attempt < 12; attempt++ {
		if attempt > 0 {
			time.Sleep(10 * time.Second)
		}
		resp, err := http.Post(ticketLookupURL, "application/json", bytes.NewReader(body))
		if err != nil {
			last = err
			continue
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			last = err
			continue
		}
		var result struct {
			Records []struct {
				Fields struct {
					SignedTicket struct {
						Value string `json:"value"`
					} `json:"signedTicket"`
				} `json:"fields"`
				ServerErrorCode string `json:"serverErrorCode"`
			} `json:"records"`
		}
		if err := json.Unmarshal(data, &result); err != nil || len(result.Records) == 0 {
			last = fmt.Errorf("ticket lookup: %s", bytes.TrimSpace(data))
			continue
		}
		if v := result.Records[0].Fields.SignedTicket.Value; v != "" {
			return base64.StdEncoding.DecodeString(v)
		}
		last = errors.New("ticket lookup: " + result.Records[0].ServerErrorCode)
	}
	return nil, last
}

func Staple(app string, cdhash []byte) error {
	if cdhash == nil {
		infoData, err := os.ReadFile(filepath.Join(app, "Contents", "Info.plist"))
		if err != nil {
			return err
		}
		info, err := ParsePlist(infoData)
		if err != nil {
			return err
		}
		exe, _ := info.(map[string]any)["CFBundleExecutable"].(string)
		data, err := os.ReadFile(filepath.Join(app, "Contents", "MacOS", exe))
		if err != nil {
			return err
		}
		if cdhash, err = CDHash(data); err != nil {
			return err
		}
	}
	ticket, err := FetchTicket(cdhash)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(app, "Contents", "CodeResources"), ticket, 0o644)
}
