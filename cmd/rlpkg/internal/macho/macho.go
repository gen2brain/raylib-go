package macho

import (
	"bytes"
	"debug/macho"
	"encoding/binary"
	"errors"
	"fmt"
)

func Universal(machos [][]byte) ([]byte, error) {
	if len(machos) == 1 {
		return machos[0], nil
	}

	type slice struct {
		data   []byte
		cpu    macho.Cpu
		subCpu uint32
		align  uint32
		offset uint32
	}
	var list []slice
	for _, data := range machos {
		f, err := macho.NewFile(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		align := uint32(12)
		if f.Cpu == macho.CpuArm64 {
			align = 14
		}
		list = append(list, slice{data: data, cpu: f.Cpu, subCpu: f.SubCpu, align: align})
	}

	offset := uint32(8 + 20*len(list))
	for i := range list {
		a := uint32(1) << list[i].align
		offset = (offset + a - 1) &^ (a - 1)
		list[i].offset = offset
		offset += uint32(len(list[i].data))
	}

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, []uint32{macho.MagicFat, uint32(len(list))})
	for _, s := range list {
		binary.Write(&buf, binary.BigEndian, []uint32{uint32(s.cpu), s.subCpu, s.offset, uint32(len(s.data)), s.align})
	}
	for _, s := range list {
		buf.Write(make([]byte, int(s.offset)-buf.Len()))
		buf.Write(s.data)
	}
	return buf.Bytes(), nil
}

func MinVersion(data []byte) (string, error) {
	f, err := macho.NewFile(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	for _, l := range f.Loads {
		raw := l.Raw()
		if len(raw) < 16 {
			continue
		}
		var v uint32
		switch f.ByteOrder.Uint32(raw) {
		case 0x32:
			v = f.ByteOrder.Uint32(raw[12:])
		case 0x24:
			v = f.ByteOrder.Uint32(raw[8:])
		default:
			continue
		}
		return fmt.Sprintf("%d.%d", v>>16, v>>8&0xff), nil
	}
	return "", errors.New("no minimum macOS version in the executable")
}
