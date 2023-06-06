package rres

/*
#define RRES_IMPLEMENTATION
#include <rres.h>
#include <stdlib.h>
#include <string.h>
rresResourceChunkInfo GetResourceChunkInfoFromArray(rresResourceChunkInfo *infos, int index)
{
	return infos[index];
}
*/
import "C"
import (
	"hash/crc32"
	"unsafe"
)

const MaxFilenameSize = 1024

// FileHeader - (16 bytes)
type FileHeader struct {
	Id         [4]byte // File identifier: rres
	Version    uint16  // File version: 100 for version 1.0
	ChunkCount uint16  // Number of resource chunks in the file (MAX: 65535)
	CdOffset   uint32  // Central Directory offset in file (0 if not available)
	Reserved   uint32  // <reserved>
}

// ResourceChunkInfo - header (32 bytes)
type ResourceChunkInfo struct {
	Type       [4]byte // Resource chunk type (FourCC)
	Id         uint32  // Resource chunk identifier (generated from filename CRC32 hash)
	CompType   byte    // Data compression algorithm
	CipherType byte    // Data encryption algorithm
	Flags      uint16  // Data flags (if required)
	PackedSize uint32  // Data chunk size (compressed/encrypted + custom data appended)
	BaseSize   uint32  // Data base size (uncompressed/unencrypted)
	NextOffset uint32  // Next resource chunk global offset (if resource has multiple chunks)
	Reserved   uint32  // <reserved>
	Crc32      uint32  // Data chunk CRC32 (propCount + props[] + data)
}

// ResourceChunkData
type ResourceChunkData struct {
	PropCount uint32         // Resource chunk properties count
	Props     *uint32        // Resource chunk properties
	Raw       unsafe.Pointer // Resource chunk raw data
}

// ResourceChunk
type ResourceChunk struct {
	Info ResourceChunkInfo // Resource chunk info
	Data ResourceChunkData // Resource chunk packed data, contains propCount, props[] and raw data
}

// ResourceMulti
//
// NOTE: It supports multiple resource chunks
type ResourceMulti struct {
	Count  uint32         // Resource chunks count
	chunks *ResourceChunk // Resource chunks
}

// Chunks - Resource chunks
func (r *ResourceMulti) Chunks() []ResourceChunk {
	return unsafe.Slice(r.chunks, r.Count)
}

// DirEntry - CDIR: rres central directory entry
type DirEntry struct {
	Id           uint32                // Resource id
	Offset       uint32                // Resource global offset in file
	Reserved     uint32                // reserved
	FileNameSize uint32                // Resource fileName size (NULL terminator and 4-byte alignment padding considered)
	fileName     [MaxFilenameSize]int8 // Resource original fileName (NULL terminated and padded to 4-byte alignment)
}

// FileName - Resource original fileName
func (d *DirEntry) FileName() string {
	cpointer := (*C.char)(unsafe.Pointer(&d.fileName[0]))
	clength := C.int(d.FileNameSize)
	fileName := C.GoStringN(cpointer, clength)
	return fileName
}

// CentralDir - CDIR: rres central directory
//
// NOTE: This data conforms the ResourceChunkData
type CentralDir struct {
	Count   uint32    // Central directory entries count
	entries *DirEntry // Central directory entries
}

// Entries - Central directory entries
func (c *CentralDir) Entries() []DirEntry {
	return unsafe.Slice(c.entries, c.Count)
}

// FontGlyphInfo - FNTG: rres font glyphs info (32 bytes)
//
// NOTE: Array of this type conforms the ResourceChunkData
type FontGlyphInfo struct {
	X, Y, Width, Height int32 // Glyph rectangle in the atlas image
	Value               int32 // Glyph codepoint value
	OffsetX, OffsetY    int32 // Glyph drawing offset (from base line)
	AdvanceX            int32 // Glyph advance X for next character
}

// ResourceDataType
//
// NOTE 1: Data type determines the properties and the data included in every chunk
//
// NOTE 2: This enum defines the basic resource data types,
// some input files could generate multiple resource chunks:
//
//	Fonts processed could generate (2) resource chunks:
//	- [FNTG] rres[0]: RRES_DATA_FONT_GLYPHS
//	- [IMGE] rres[1]: RRES_DATA_IMAGE
//
//	Mesh processed could generate (n) resource chunks:
//	- [VRTX] rres[0]: RRES_DATA_VERTEX
//	...
//	- [VRTX] rres[n]: RRES_DATA_VERTEX
type ResourceDataType int32

const (
	// FourCC: NULL - Reserved for empty chunks, no props/data
	DataNull ResourceDataType = iota
	// FourCC: RAWD - Raw file data, 4 properties
	//    props[0]:size (bytes)
	//    props[1]:extension01 (big-endian: ".png" = 0x2e706e67)
	//    props[2]:extension02 (additional part, extensions with +3 letters)
	//    props[3]:reserved
	//    data: raw bytes
	DataRaw
	// FourCC: TEXT - Text file data, 4 properties
	//    props[0]:size (bytes)
	//    props[1]:rresTextEncoding
	//    props[2]:rresCodeLang
	//    props[3]:cultureCode
	//    data: text
	DataText
	// FourCC: IMGE - Image file data, 4 properties
	//    props[0]:width
	//    props[1]:height
	//    props[2]:rresPixelFormat
	//    props[3]:mipmaps
	//    data: pixels
	DataImage
	// FourCC: WAVE - Audio file data, 4 properties
	//    props[0]:frameCount
	//    props[1]:sampleRate
	//    props[2]:sampleSize
	//    props[3]:channels
	//    data: samples
	DataWave
	// FourCC: VRTX - Vertex file data, 4 properties
	//    props[0]:vertexCount
	//    props[1]:rresVertexAttribute
	//    props[2]:componentCount
	//    props[3]:rresVertexFormat
	//    data: vertex
	DataVertex
	// FourCC: FNTG - Font glyphs info data, 4 properties
	//    props[0]:baseSize
	//    props[1]:glyphCount
	//    props[2]:glyphPadding
	//    props[3]:rresFontStyle
	//    data: rresFontGlyphInfo[0..glyphCount]
	DataFontGlyphs
	// FourCC: LINK - External linked file, 1 property
	//    props[0]:size (bytes)
	//    data: filepath (as provided on input)
	DataLink ResourceDataType = 99
	// FourCC: CDIR - Central directory for input files
	//    props[0]:entryCount, 1 property
	//    data: rresDirEntry[0..entryCount]
	DataDirectory ResourceDataType = 100
)

// CompressionType - Compression algorithms
//
// value required by ResourceChunkInfo.CompType
//
// NOTE 1: This enum just lists some common data compression algorithms for convenience,
// The rres packer tool and the engine-specific library are responsible to implement the desired ones,
//
// NOTE 2: ResourceChunkInfo.CompType is a byte-size value, limited to [0..255]
type CompressionType int32

const (
	CompNone    CompressionType = 0  // No data compression
	CompRle     CompressionType = 1  // RLE compression
	CompDeflate CompressionType = 10 // DEFLATE compression
	CompLz4     CompressionType = 20 // LZ4 compression
	CompLzma2   CompressionType = 30 // LZMA2 compression
	CompQoi     CompressionType = 40 // QOI compression, useful for RGB(A) image data
)

// EncryptionType - Encryption algorithms
//
// value required by ResourceChunkInfo.CipherType
//
// NOTE 1: This enum just lists some common data encryption algorithms for convenience,
// The rres packer tool and the engine-specific library are responsible to implement the desired ones,
//
// NOTE 2: Some encryption algorithms could require/generate additional data (seed, salt, nonce, MAC...)
// in those cases, that extra data must be appended to the original encrypted message and added to the resource data chunk
//
// NOTE 3: ResourceChunkInfo.CipherType is a byte-size value, limited to [0..255]
type EncryptionType int32

const (
	CipherNone              EncryptionType = 0  // No data encryption
	CipherXor               EncryptionType = 1  // XOR encryption, generic using 128bit key in blocks
	CipherDes               EncryptionType = 10 // DES encryption
	CipherTdes              EncryptionType = 11 // Triple DES encryption
	CipherIdea              EncryptionType = 20 // IDEA encryption
	CipherAes               EncryptionType = 30 // AES (128bit or 256bit) encryption
	CipherAesGCM            EncryptionType = 31 // AES Galois/Counter Mode (Galois Message Authentication Code - GMAC)
	CipherXtea              EncryptionType = 40 // XTEA encryption
	CipherBlowfish          EncryptionType = 50 // BLOWFISH encryption
	CipherRsa               EncryptionType = 60 // RSA asymmetric encryption
	CipherSalsa20           EncryptionType = 70 // SALSA20 encryption
	CipherChacha20          EncryptionType = 71 // CHACHA20 encryption
	CipherXchacha20         EncryptionType = 72 // XCHACHA20 encryption
	CipherXchacha20Poly1305 EncryptionType = 73 // XCHACHA20 with POLY1305 for message authentication (MAC)
)

// ErrorType - error codes
//
// NOTE: Error codes when processing rres files
type ErrorType int32

const (
	Success           ErrorType = iota // rres file loaded/saved successfully
	ErrorFileNotFound                  // rres file can not be opened (spelling issues, file actually does not exist...)
	ErrorFileFormat                    // rres file format not a supported (wrong header, wrong identifier)
	ErrorMemoryAlloc                   // Memory could not be allocated for operation.
)

// TextEncoding - TEXT: Text encoding property values
type TextEncoding int32

const (
	TextEncodingUndefined TextEncoding = 0  // Not defined, usually UTF-8
	TextEncodingUtf8      TextEncoding = 1  // UTF-8 text encoding
	TextEncodingUtf8Bom   TextEncoding = 2  // UTF-8 text encoding with Byte-Order-Mark
	TextEncodingUtf16Le   TextEncoding = 10 // UTF-16 Little Endian text encoding
	TextEncodingUtf16Be   TextEncoding = 11 // UTF-16 Big Endian text encoding
)

// CodeLang - TEXT: Text code language
//
// NOTE: It could be useful for code script resources
type CodeLang int32

const (
	CodeLangUndefined CodeLang = iota // Undefined code language, text is plain text
	CodeLangC                         // Text contains C code
	CodeLangCpp                       // Text contains C++ code
	CodeLangCs                        // Text contains C# code
	CodeLangLua                       // Text contains Lua code
	CodeLangJs                        // Text contains JavaScript code
	CodeLangPython                    // Text contains Python code
	CodeLangRust                      // Text contains Rust code
	CodeLangZig                       // Text contains Zig code
	CodeLangOdin                      // Text contains Odin code
	CodeLangJai                       // Text contains Jai code
	CodeLangGdscript                  // Text contains GDScript (Godot) code
	CodeLangGlsl                      // Text contains GLSL shader code
)

// PixelFormat - IMGE: Image/Texture pixel formats
type PixelFormat int32

const (
	PixelFormatUndefined          PixelFormat = iota // Undefined pixel format
	PixelFormatUncompGrayscale                       // 8 bit per pixel (no alpha)
	PixelFormatUncompGrayAlpha                       // 16 bpp (2 channels)
	PixelFormatUncompR5g6b5                          // 16 bpp
	PixelFormatUncompR8g8b8                          // 24 bpp
	PixelFormatUncompR5g5b5a1                        // 16 bpp (1 bit alpha)
	PixelFormatUncompR4g4b4a4                        // 16 bpp (4 bit alpha)
	PixelFormatUncompR8g8b8a8                        // 32 bpp
	PixelFormatUncompR32                             // 32 bpp (1 channel - float)
	PixelFormatUncompR32g32b32                       // 32*3 bpp (3 channels - float)
	PixelFormatUncompR32g32b32a32                    // 32*4 bpp (4 channels - float)
	PixelFormatCompDxt1Rgb                           // 4 bpp (no alpha)
	PixelFormatCompDxt1Rgba                          // 4 bpp (1 bit alpha)
	PixelFormatCompDxt3Rgba                          // 8 bpp
	PixelFormatCompDxt5Rgba                          // 8 bpp
	PixelFormatCompEtc1Rgb                           // 4 bpp
	PixelFormatCompEtc2Rgb                           // 4 bpp
	PixelFormatCompEtc2EacRgba                       // 8 bpp
	PixelFormatCompPvrtRgb                           // 4 bpp
	PixelFormatCompPvrtRgba                          // 4 bpp
	PixelFormatCompAtsc4x4Rgba                       // 8 bpp
	PixelFormatCompAtsc8x8Rgba                       // 2 bpp
)

// VertexAttribute - VRTX: Vertex data attribute
//
// NOTE: The expected number of components for every vertex attribute is provided as a property to data,
// the listed components count are the expected/default ones
type VertexAttribute int32

const (
	VertexAttributePosition  VertexAttribute = 0   // Vertex position attribute: [x, y, z]
	VertexAttributeTexcoord1 VertexAttribute = 10  // Vertex texture coordinates attribute: [u, v]
	VertexAttributeTexcoord2 VertexAttribute = 11  // Vertex texture coordinates attribute: [u, v]
	VertexAttributeTexcoord3 VertexAttribute = 12  // Vertex texture coordinates attribute: [u, v]
	VertexAttributeTexcoord4 VertexAttribute = 13  // Vertex texture coordinates attribute: [u, v]
	VertexAttributeNormal    VertexAttribute = 20  // Vertex normal attribute: [x, y, z]
	VertexAttributeTangent   VertexAttribute = 30  // Vertex tangent attribute: [x, y, z, w]
	VertexAttributeColor     VertexAttribute = 40  // Vertex color attribute: [r, g, b, a]
	VertexAttributeIndex     VertexAttribute = 100 // Vertex index attribute: [i]
)

// VertexFormat - VRTX: Vertex data format type
type VertexFormat int32

const (
	VertexFormatUbyte  VertexFormat = iota // 8 bit unsigned integer data
	VertexFormatByte                       // 8 bit signed integer data
	VertexFormatUshort                     // 16 bit unsigned integer data
	VertexFormatShort                      // 16 bit signed integer data
	VertexFormatUint                       // 32 bit unsigned integer data
	VertexFormatInt                        // 32 bit integer data
	VertexFormatHfloat                     // 16 bit float data
	VertexFormatFloat                      // 32 bit float data
)

// FontStyle - FNTG: Font style
type FontStyle int32

const (
	FontStyleUndefined FontStyle = iota // Undefined font style
	FontStyleRegular                    // Regular font style
	FontStyleBold                       // Bold font style
	FontStyleItalic                     // Italic font style
)

// LoadResourceChunk - Load one resource chunk for provided id
func LoadResourceChunk(fileName string, rresId int32) ResourceChunk {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.rresLoadResourceChunk(cfileName, C.int(rresId))
	v := *(*ResourceChunk)(unsafe.Pointer(&ret))
	return v
}

// UnloadResourceChunk - Unload resource chunk from memory
func UnloadResourceChunk(chunk *ResourceChunk) {
	cchunk := *(*C.rresResourceChunk)(unsafe.Pointer(chunk))
	C.rresUnloadResourceChunk(cchunk)
}

// LoadResourceMulti - Load resource for provided id (multiple resource chunks)
func LoadResourceMulti(fileName string, rresId int32) ResourceMulti {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.rresLoadResourceMulti(cfileName, C.int(rresId))
	v := *(*ResourceMulti)(unsafe.Pointer(&ret))
	return v
}

// UnloadResourceMulti - Unload resource from memory (multiple resource chunks)
func UnloadResourceMulti(multi *ResourceMulti) {
	cmulti := *(*C.rresResourceMulti)(unsafe.Pointer(multi))
	C.rresUnloadResourceMulti(cmulti)
}

// LoadResourceChunkInfo - Load resource chunk info for provided id
func LoadResourceChunkInfo(fileName string, rresId int32) ResourceChunkInfo {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.rresLoadResourceChunkInfo(cfileName, C.int(rresId))
	v := *(*ResourceChunkInfo)(unsafe.Pointer(&ret))
	return v
}

// LoadResourceChunkInfoAll - Load all resource chunks info
func LoadResourceChunkInfoAll(fileName string) []ResourceChunkInfo {
	// Convert the fileName into a CString and releases the memory afterwards
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))

	// The length of the resulted array is saved in the chunkCount variable
	var chunkCount C.uint
	cinfos := C.rresLoadResourceChunkInfoAll(cfileName, &chunkCount)

	// The C array can be released afterwards, because the values are stored in a golang slice
	defer C.free(unsafe.Pointer(cinfos))

	// Iterate over the C array and store the values in a golang slice
	infos := make([]ResourceChunkInfo, chunkCount)
	for i := 0; i < int(chunkCount); i++ {
		// Get the C value from the C array
		ret := C.GetResourceChunkInfoFromArray(cinfos, C.int(i))
		// Convert the C value into a golang value
		v := *(*ResourceChunkInfo)(unsafe.Pointer(&ret))
		// Save the golang value in the golang slice
		infos[i] = v
	}

	return infos
}

// LoadCentralDirectory - Load central directory resource chunk from file
func LoadCentralDirectory(fileName string) CentralDir {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.rresLoadCentralDirectory(cfileName)
	v := *(*CentralDir)(unsafe.Pointer(&ret))
	return v
}

// UnloadCentralDirectory - Unload central directory resource chunk
func UnloadCentralDirectory(dir *CentralDir) {
	cdir := *(*C.rresCentralDir)(unsafe.Pointer(dir))
	C.rresUnloadCentralDirectory(cdir)
}

// GetDataType - Get ResourceDataType from FourCC code
func GetDataType(fourCC [4]byte) ResourceDataType {
	value := string(fourCC[:])
	switch value {
	case "NULL":
		return DataNull
	case "RAWD":
		return DataRaw
	case "TEXT":
		return DataText
	case "IMGE":
		return DataImage
	case "WAVE":
		return DataWave
	case "VRTX":
		return DataVertex
	case "FNTG":
		return DataFontGlyphs
	case "LINK":
		return DataLink
	case "CDIR":
		return DataDirectory
	default:
		return 0
	}
}

// GetResourceId - Get resource id for a provided filename
//
// NOTE: It requires CDIR available in the file (it's optinal by design)
func GetResourceId(dir CentralDir, fileName string) int32 {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cdir := *(*C.rresCentralDir)(unsafe.Pointer(&dir))
	ret := C.rresGetResourceId(cdir, cfileName)
	v := int32(ret)
	return v
}

// ComputeCRC32 - Compute CRC32 hash
//
// NOTE: CRC32 is used as rres id, generated from original filename
func ComputeCRC32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

// SetCipherPassword - Set password to be used on data decryption
//
// NOTE: The cipher password is kept as an internal pointer to provided string, it's up to the user to manage that sensible data properly
//
// Password should be to allocate and set before loading an encrypted resource and it should be cleaned/wiped after the encrypted resource has been loaded
//
// You can use the WipeCipherPassword function to clear the password
func SetCipherPassword(pass string) {
	cpass := C.CString(pass)
	C.rresSetCipherPassword(cpass)
}

// GetCipherPassword - Get password to be used on data decryption
func GetCipherPassword() string {
	cpass := C.rresGetCipherPassword()
	return C.GoString(cpass)
}

// WipeCipherPassword - Clears the password from the C memory using explicit_bzero
//
// This is an approach but no guarantee
func WipeCipherPassword() {
	cpass := C.rresGetCipherPassword()
	C.explicit_bzero(unsafe.Pointer(cpass), C.strlen(cpass))
	C.free(unsafe.Pointer(cpass))
}
