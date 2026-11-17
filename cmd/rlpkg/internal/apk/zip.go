package apk

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strings"
)

type Entry struct {
	Name   string
	Method uint16
	time   uint16
	date   uint16
	crc    uint32
	csize  uint32
	usize  uint32
	Data   []byte
	Align  int
}

const (
	Align    = 4
	AlignLib = 16384
)

func ReadTemplate(path string) ([]*Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ReadArchive(data)
}

func ReadArchive(data []byte) ([]*Entry, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	var entries []*Entry
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "META-INF/") || strings.HasSuffix(f.Name, "/") {
			continue
		}
		r, err := f.OpenRaw()
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		align := 0
		if f.Method == zip.Store {
			align = Align
			if strings.HasPrefix(f.Name, "lib/") && strings.HasSuffix(f.Name, ".so") {
				align = AlignLib
			}
		}
		entries = append(entries, &Entry{
			Name: f.Name, Method: f.Method, time: f.ModifiedTime, date: f.ModifiedDate,
			crc: f.CRC32, csize: uint32(f.CompressedSize64), usize: uint32(f.UncompressedSize64),
			Data: data, Align: align,
		})
	}
	return entries, nil
}

func Find(entries []*Entry, name string) *Entry {
	for _, e := range entries {
		if e.Name == name {
			return e
		}
	}
	return nil
}

func (e *Entry) Content() ([]byte, error) {
	if e.Method == zip.Store {
		return e.Data, nil
	}
	return io.ReadAll(flate.NewReader(bytes.NewReader(e.Data)))
}

func (e *Entry) SetContent(data []byte) error {
	e.crc = crc32.ChecksumIEEE(data)
	e.usize = uint32(len(data))
	if e.Method == zip.Store {
		e.Data = data
		e.csize = e.usize
		return nil
	}
	var buf bytes.Buffer
	fw, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		return err
	}
	fw.Write(data)
	if err := fw.Close(); err != nil {
		return err
	}
	e.Data = buf.Bytes()
	e.csize = uint32(buf.Len())
	return nil
}

func NewEntry(name string, data []byte, method uint16, align int) (*Entry, error) {
	e := &Entry{Name: name, Method: method, Align: align, date: 0x0221, time: 0x0800}
	return e, e.SetContent(data)
}

func WriteZip(entries []*Entry) ([]byte, error) {
	var buf bytes.Buffer
	var central bytes.Buffer

	for _, e := range entries {
		if len(e.Name) > 0xffff {
			return nil, fmt.Errorf("%s: name too long", e.Name)
		}
		var extra []byte
		if e.Align > 0 {
			dataOff := buf.Len() + 30 + len(e.Name)
			pad := (e.Align - dataOff%e.Align) % e.Align
			for pad > 0 && pad < 6 {
				pad += e.Align
			}
			if pad > 0 {
				extra = make([]byte, pad)
				binary.LittleEndian.PutUint16(extra, 0xd935)
				binary.LittleEndian.PutUint16(extra[2:], uint16(pad-4))
				binary.LittleEndian.PutUint16(extra[4:], uint16(e.Align))
			}
		}

		offset := uint32(buf.Len())
		hdr := make([]byte, 30)
		binary.LittleEndian.PutUint32(hdr, 0x04034b50)
		binary.LittleEndian.PutUint16(hdr[4:], 20)
		binary.LittleEndian.PutUint16(hdr[8:], e.Method)
		binary.LittleEndian.PutUint16(hdr[10:], e.time)
		binary.LittleEndian.PutUint16(hdr[12:], e.date)
		binary.LittleEndian.PutUint32(hdr[14:], e.crc)
		binary.LittleEndian.PutUint32(hdr[18:], e.csize)
		binary.LittleEndian.PutUint32(hdr[22:], e.usize)
		binary.LittleEndian.PutUint16(hdr[26:], uint16(len(e.Name)))
		binary.LittleEndian.PutUint16(hdr[28:], uint16(len(extra)))
		buf.Write(hdr)
		buf.WriteString(e.Name)
		buf.Write(extra)
		buf.Write(e.Data)

		cd := make([]byte, 46)
		binary.LittleEndian.PutUint32(cd, 0x02014b50)
		binary.LittleEndian.PutUint16(cd[4:], 20)
		binary.LittleEndian.PutUint16(cd[6:], 20)
		binary.LittleEndian.PutUint16(cd[10:], e.Method)
		binary.LittleEndian.PutUint16(cd[12:], e.time)
		binary.LittleEndian.PutUint16(cd[14:], e.date)
		binary.LittleEndian.PutUint32(cd[16:], e.crc)
		binary.LittleEndian.PutUint32(cd[20:], e.csize)
		binary.LittleEndian.PutUint32(cd[24:], e.usize)
		binary.LittleEndian.PutUint16(cd[28:], uint16(len(e.Name)))
		binary.LittleEndian.PutUint32(cd[42:], offset)
		central.Write(cd)
		central.WriteString(e.Name)
	}

	cdOffset := uint32(buf.Len())
	buf.Write(central.Bytes())

	eocd := make([]byte, 22)
	binary.LittleEndian.PutUint32(eocd, 0x06054b50)
	binary.LittleEndian.PutUint16(eocd[8:], uint16(len(entries)))
	binary.LittleEndian.PutUint16(eocd[10:], uint16(len(entries)))
	binary.LittleEndian.PutUint32(eocd[12:], uint32(central.Len()))
	binary.LittleEndian.PutUint32(eocd[16:], cdOffset)
	buf.Write(eocd)
	return buf.Bytes(), nil
}
