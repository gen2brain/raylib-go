package raylib

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"os"
	"unsafe"
)

// rRES file header (8 byte)
type RRESFileHeader struct {
	// File identifier: rRES (4 byte)
	ID [4]int8
	// File version and subversion (2 byte)
	Version uint16
	// Number of resources in this file (2 byte)
	Count uint16
}

// rRES info header, every resource includes this header (16 byte + 16 byte)
type RRESInfoHeader struct {
	// Resource unique identifier (4 byte)
	ID uint32
	// Resource data type (1 byte)
	DataType uint8
	// Resource data compression type (1 byte)
	CompType uint8
	// Resource data encryption type (1 byte)
	CryptoType uint8
	// Resource data parts count, used for splitted data (1 byte)
	PartsCount uint8
	// Resource data size (compressed or not, only DATA) (4 byte)
	DataSize uint32
	// Resource data size (uncompressed, only DATA) (4 byte)
	UncompSize uint32
	// Resouce parameter 1 (4 byte)
	Param1 uint32
	// Resouce parameter 2 (4 byte)
	Param2 uint32
	// Resouce parameter 3 (4 byte)
	Param3 uint32
	// Resouce parameter 4 (4 byte)
	Param4 uint32
}

// rRES data types
const (
	RRESTypeRaw = iota
	RRESTypeImage
	RRESTypeWave
	RRESTypeVertex
	RRESTypeText
	RRESTypeFontImage
	RRESTypeFontData
	RRESTypeDirectory
)

// Compression types
const (
	// No data compression
	RRESCompNone = iota
	// DEFLATE compression
	RRESCompDeflate
	// LZ4 compression
	RRESCompLz4
	// LZMA compression
	RRESCompLzma
	// BROTLI compression
	RRESCompBrotli
)

// Image formats
const (
	// 8 bit per pixel (no alpha)
	RRESImUncompGrayscale = iota + 1
	// 16 bpp (2 channels)
	RRESImUncompGrayAlpha
	// 16 bpp
	RRESImUncompR5g6b5
	// 24 bpp
	RRESImUncompR8g8b8
	// 16 bpp (1 bit alpha)
	RRESImUncompR5g5b5a1
	// 16 bpp (4 bit alpha)
	RRESImUncompR4g4b4a4
	// 32 bpp
	RRESImUncompR8g8b8a8
	// 4 bpp (no alpha)
	RRESImCompDxt1Rgb
	// 4 bpp (1 bit alpha)
	RRESImCompDxt1Rgba
	// 8 bpp
	RRESImCompDxt3Rgba
	// 8 bpp
	RRESImCompDxt5Rgba
	// 4 bpp
	RRESImCompEtc1Rgb
	// 4 bpp
	RRESImCompEtc2Rgb
	// 8 bpp
	RRESImCompEtc2EacRgba
	// 4 bpp
	RRESImCompPvrtRgb
	// 4 bpp
	RRESImCompPvrtRgba
	// 8 bpp
	RRESImCompAstc4x4Rgba
	// 2 bpp
	RRESImCompAstc8x8Rgba
)

// RRESVert
const (
	RRESVertPosition = iota
	RRESVertTexcoord1
	RRESVertTexcoord2
	RRESVertTexcoord3
	RRESVertTexcoord4
	RRESVertNormal
	RRESVertTangent
	RRESVertColor
	RRESVertIndex
)

// RRESVert
const (
	RRESVertByte = iota
	RRESVertShort
	RRESVertInt
	RRESVertHfloat
	RRESVertFloat
)

// LoadResource - Load resource from file (only one)
// NOTE: Returns uncompressed data with parameters, only first resource found
func LoadResource(fileName string) []byte {
	return LoadResourceByID(fileName, 0)
}

// LoadResourceByID - Load resource from file by id
// NOTE: Returns uncompressed data with parameters, search resource by id
func LoadResourceByID(fileName string, rresID int) (data []byte) {
	file, err := OpenAsset(fileName)
	if err != nil {
		TraceLog(LogWarning, "[%s] rRES raylib resource file could not be opened", fileName)
		return
	}
	defer file.Close()

	fileHeader := RRESFileHeader{}
	infoHeader := RRESInfoHeader{}

	// Read rres file header
	err = binary.Read(file, binary.LittleEndian, &fileHeader)
	if err != nil {
		TraceLog(LogWarning, err.Error())
		return
	}

	//fmt.Printf("%+v\n", fileHeader)

	// Verify "rRES" identifier
	id := fmt.Sprintf("%c", fileHeader.ID)
	if id != "[r R E S]" {
		TraceLog(LogWarning, "[%s] is not a valid raylib resource file", fileName)
		return
	}

	file.Seek(int64(unsafe.Sizeof(fileHeader)), os.SEEK_CUR)

	for i := 0; i < int(fileHeader.Count); i++ {
		// Read resource info and parameters
		err = binary.Read(file, binary.LittleEndian, &infoHeader)
		if err != nil {
			TraceLog(LogWarning, err.Error())
			return
		}

		//fmt.Printf("%+v\n", infoHeader)

		file.Seek(int64(unsafe.Sizeof(infoHeader)), os.SEEK_CUR)

		if int(infoHeader.ID) == rresID {
			// Read resource data block
			data = make([]byte, infoHeader.DataSize)
			file.Read(data)

			if infoHeader.CompType == RRESCompDeflate {
				// Uncompress data
				b := bytes.NewReader(data)
				r := flate.NewReader(b)

				data = make([]byte, infoHeader.UncompSize)
				r.Read(data)
			}

			if len(data) > 0 {
				TraceLog(LogInfo, "[%s][ID %d] Resource data loaded successfully", fileName, infoHeader.ID)
			}
		} else {
			// Skip required data to read next resource infoHeader
			file.Seek(int64(infoHeader.DataSize), os.SEEK_CUR)
		}
	}

	if len(data) == 0 {
		TraceLog(LogInfo, "[%s][ID %d] Requested resource could not be found", fileName, rresID)
	}

	return
}
