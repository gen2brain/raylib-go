package rres

type Data struct {
	// Resource type (4 byte)
	Type uint32

	// Resource parameter 1 (4 byte)
	Param1 uint32
	// Resource parameter 2 (4 byte)
	Param2 uint32
	// Resource parameter 3 (4 byte)
	Param3 uint32
	// Resource parameter 4 (4 byte)
	Param4 uint32

	// Resource data
	Data []byte
}

// FileHeader - rRES file header (8 byte)
type FileHeader struct {
	// File identifier: rRES (4 byte)
	ID [4]byte
	// File version and subversion (2 byte)
	Version uint16
	// Number of resources in this file (2 byte)
	Count uint16
}

// InfoHeader - rRES info header, every resource includes this header (16 byte + 16 byte)
type InfoHeader struct {
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
	// Resource parameter 1 (4 byte)
	Param1 uint32
	// Resource parameter 2 (4 byte)
	Param2 uint32
	// Resource parameter 3 (4 byte)
	Param3 uint32
	// Resource parameter 4 (4 byte)
	Param4 uint32
}

// rRES data types
const (
	TypeRaw = iota
	TypeImage
	TypeWave
	TypeVertex
	TypeText
	TypeFontImage
	TypeFontCharData
	TypeDirectory
	TypeVorbis
)

// Compression types
const (
	// No data compression
	CompNone = iota
	// DEFLATE compression
	CompDeflate
	// LZ4 compression
	CompLZ4
	// LZMA compression
	CompLZMA
	// BROTLI compression
	CompBrotli
	// LZMA2 (XZ) compression
	CompLZMA2
	// BZIP2 compression
	CompBZIP2
)

// Encryption types
const (
	// No data encryption
	CryptoNone = iota
	// XOR (128 bit) encryption
	CryptoXOR
	// RIJNDAEL (128 bit) encryption (AES)
	CryptoAES
	// Triple DES encryption
	Crypto3DES
	// Blowfish encryption
	CryptoBlowfish
	// Extended TEA encryption
	CryptoXTEA
)

// Image formats
const (
	// 8 bit per pixel (no alpha)
	ImUncompGrayscale = iota + 1
	// 16 bpp (2 channels)
	ImUncompGrayAlpha
	// 16 bpp
	ImUncompR5g6b5
	// 24 bpp
	ImUncompR8g8b8
	// 16 bpp (1 bit alpha)
	ImUncompR5g5b5a1
	// 16 bpp (4 bit alpha)
	ImUncompR4g4b4a4
	// 32 bpp
	ImUncompR8g8b8a8
	// 4 bpp (no alpha)
	ImCompDxt1Rgb
	// 4 bpp (1 bit alpha)
	ImCompDxt1Rgba
	// 8 bpp
	ImCompDxt3Rgba
	// 8 bpp
	ImCompDxt5Rgba
	// 4 bpp
	ImCompEtc1Rgb
	// 4 bpp
	ImCompEtc2Rgb
	// 8 bpp
	ImCompEtc2EacRgba
	// 4 bpp
	ImCompPvrtRgb
	// 4 bpp
	ImCompPvrtRgba
	// 8 bpp
	ImCompAstc4x4Rgba
	// 2 bpp
	ImCompAstc8x8Rgba
)

// Vert
const (
	VertPosition = iota
	VertTexcoord1
	VertTexcoord2
	VertTexcoord3
	VertTexcoord4
	VertNormal
	VertTangent
	VertColor
	VertIndex
)

// Vert
const (
	VertByte = iota
	VertShort
	VertInt
	VertHfloat
	VertFloat
)
