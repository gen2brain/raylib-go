package apple

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

func base64Std(b []byte) string { return base64.StdEncoding.EncodeToString(b) }

func ParsePlist(data []byte) (any, error) {
	dec := xml.NewDecoder(bytes.NewReader(data))
	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		if se, ok := tok.(xml.StartElement); ok && se.Name.Local == "plist" {
			return plistValue(dec)
		}
	}
}

func plistValue(dec *xml.Decoder) (any, error) {
	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			return plistElement(dec, t)
		case xml.EndElement:
			return nil, io.ErrUnexpectedEOF
		}
	}
}

func plistElement(dec *xml.Decoder, se xml.StartElement) (any, error) {
	switch se.Name.Local {
	case "dict":
		m := map[string]any{}
		for {
			tok, err := dec.Token()
			if err != nil {
				return nil, err
			}
			switch t := tok.(type) {
			case xml.EndElement:
				return m, nil
			case xml.StartElement:
				if t.Name.Local != "key" {
					return nil, fmt.Errorf("plist: expected key, got %s", t.Name.Local)
				}
				var key string
				if err := dec.DecodeElement(&key, &t); err != nil {
					return nil, err
				}
				v, err := plistValue(dec)
				if err != nil {
					return nil, err
				}
				m[key] = v
			}
		}
	case "array":
		var list []any
		for {
			tok, err := dec.Token()
			if err != nil {
				return nil, err
			}
			switch t := tok.(type) {
			case xml.EndElement:
				if list == nil {
					list = []any{}
				}
				return list, nil
			case xml.StartElement:
				v, err := plistElement(dec, t)
				if err != nil {
					return nil, err
				}
				list = append(list, v)
			}
		}
	case "true", "false":
		dec.Skip()
		return se.Name.Local == "true", nil
	}
	var text string
	if err := dec.DecodeElement(&text, &se); err != nil {
		return nil, err
	}
	switch se.Name.Local {
	case "string", "date":
		return text, nil
	case "integer":
		return strconv.ParseInt(strings.TrimSpace(text), 10, 64)
	case "real":
		return strconv.ParseFloat(strings.TrimSpace(text), 64)
	case "data":
		return base64.StdEncoding.DecodeString(strings.Join(strings.Fields(text), ""))
	}
	return nil, fmt.Errorf("plist: unsupported element %s", se.Name.Local)
}

func WritePlist(v any) []byte {
	var b bytes.Buffer
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n")
	writePlistValue(&b, v, 0)
	b.WriteString("</plist>\n")
	return b.Bytes()
}

func writePlistValue(b *bytes.Buffer, v any, depth int) {
	indent := strings.Repeat("\t", depth)
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		b.WriteString(indent + "<dict>\n")
		for _, k := range keys {
			b.WriteString(indent + "\t<key>")
			xml.EscapeText(b, []byte(k))
			b.WriteString("</key>\n")
			writePlistValue(b, t[k], depth+1)
		}
		b.WriteString(indent + "</dict>\n")
	case []any:
		b.WriteString(indent + "<array>\n")
		for _, e := range t {
			writePlistValue(b, e, depth+1)
		}
		b.WriteString(indent + "</array>\n")
	case string:
		b.WriteString(indent + "<string>")
		xml.EscapeText(b, []byte(t))
		b.WriteString("</string>\n")
	case bool:
		if t {
			b.WriteString(indent + "<true/>\n")
		} else {
			b.WriteString(indent + "<false/>\n")
		}
	case int64:
		fmt.Fprintf(b, "%s<integer>%d</integer>\n", indent, t)
	case float64:
		fmt.Fprintf(b, "%s<real>%g</real>\n", indent, t)
	case []byte:
		b.WriteString(indent + "<data>\n" + indent + base64Std(t) + "\n" + indent + "</data>\n")
	}
}

func PlistDER(v any) ([]byte, error) {
	body, err := derValue(v)
	if err != nil {
		return nil, err
	}
	inner := append(derTag(0x02, []byte{1}), body...)
	return derTag(0x70, inner), nil
}

func derValue(v any) ([]byte, error) {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var body []byte
		for _, k := range keys {
			val, err := derValue(t[k])
			if err != nil {
				return nil, err
			}
			body = append(body, derTag(0x30, append(derTag(0x0c, []byte(k)), val...))...)
		}
		return derTag(0xb0, body), nil
	case []any:
		var body []byte
		for _, e := range t {
			val, err := derValue(e)
			if err != nil {
				return nil, err
			}
			body = append(body, val...)
		}
		return derTag(0x30, body), nil
	case string:
		return derTag(0x0c, []byte(t)), nil
	case bool:
		if t {
			return derTag(0x01, []byte{0xff}), nil
		}
		return derTag(0x01, []byte{0}), nil
	case int64:
		var b []byte
		n := t
		for {
			b = append([]byte{byte(n)}, b...)
			if n >= -128 && n < 128 {
				break
			}
			n >>= 8
		}
		return derTag(0x02, b), nil
	case []byte:
		return derTag(0x04, t), nil
	}
	return nil, errors.New("entitlements: unsupported value type")
}

func derTag(tag byte, body []byte) []byte {
	out := []byte{tag}
	n := len(body)
	switch {
	case n < 0x80:
		out = append(out, byte(n))
	case n < 0x100:
		out = append(out, 0x81, byte(n))
	case n < 0x10000:
		out = append(out, 0x82, byte(n>>8), byte(n))
	default:
		out = append(out, 0x83, byte(n>>16), byte(n>>8), byte(n))
	}
	return append(out, body...)
}
