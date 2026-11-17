package apk

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	wireVarint  = 0
	wireFixed64 = 1
	wireBytes   = 2
	wireFixed32 = 5
)

type protoField struct {
	num  int
	wire int
	v    uint64
	b    []byte
}

func protoParse(data []byte) ([]protoField, error) {
	var fields []protoField
	for len(data) > 0 {
		key, n := binary.Uvarint(data)
		if n <= 0 {
			return nil, errors.New("protobuf: bad key")
		}
		data = data[n:]
		f := protoField{num: int(key >> 3), wire: int(key & 7)}
		switch f.wire {
		case wireVarint:
			v, n := binary.Uvarint(data)
			if n <= 0 {
				return nil, errors.New("protobuf: bad varint")
			}
			f.v, data = v, data[n:]
		case wireFixed64:
			if len(data) < 8 {
				return nil, errors.New("protobuf: short fixed64")
			}
			f.b, data = data[:8], data[8:]
		case wireFixed32:
			if len(data) < 4 {
				return nil, errors.New("protobuf: short fixed32")
			}
			f.b, data = data[:4], data[4:]
		case wireBytes:
			l, n := binary.Uvarint(data)
			if n <= 0 || uint64(len(data)-n) < l {
				return nil, errors.New("protobuf: bad length")
			}
			f.b, data = data[n:n+int(l)], data[n+int(l):]
		default:
			return nil, fmt.Errorf("protobuf: unsupported wire type %d", f.wire)
		}
		fields = append(fields, f)
	}
	return fields, nil
}

func protoEncode(fields []protoField) []byte {
	var out []byte
	for _, f := range fields {
		out = binary.AppendUvarint(out, uint64(f.num)<<3|uint64(f.wire))
		switch f.wire {
		case wireVarint:
			out = binary.AppendUvarint(out, f.v)
		case wireBytes:
			out = binary.AppendUvarint(out, uint64(len(f.b)))
			out = append(out, f.b...)
		default:
			out = append(out, f.b...)
		}
	}
	return out
}

func protoString(num int, s string) protoField {
	return protoField{num: num, wire: wireBytes, b: []byte(s)}
}

func protoMessage(num int, fields []protoField) protoField {
	return protoField{num: num, wire: wireBytes, b: protoEncode(fields)}
}

func protoFind(fields []protoField, num int) int {
	for i, f := range fields {
		if f.num == num {
			return i
		}
	}
	return -1
}

func protoGetString(fields []protoField, num int) string {
	if i := protoFind(fields, num); i >= 0 {
		return string(fields[i].b)
	}
	return ""
}

func protoSetString(fields []protoField, num int, s string) []protoField {
	if i := protoFind(fields, num); i >= 0 {
		fields[i].b = []byte(s)
		return fields
	}
	return append(fields, protoString(num, s))
}
