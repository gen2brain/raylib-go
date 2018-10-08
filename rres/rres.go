package rres

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unsafe"

	"github.com/dsnet/compress/bzip2"
	"github.com/golang/snappy"
	"github.com/klauspost/compress/flate"
	"github.com/pierrec/lz4"
	xor "github.com/rootlch/encrypt"
	"github.com/ulikunitz/xz"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/xtea"

	"github.com/gen2brain/raylib-go/raylib"
)

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
	// Snappy compression
	CompSnappy
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

// LoadResource - Load resource from file by id
// NOTE: Returns uncompressed data with parameters, search resource by id
func LoadResource(reader io.ReadSeeker, rresID int, key []byte) (data Data) {
	var fileHeader FileHeader
	var infoHeader InfoHeader

	reader.Seek(0, 0)

	// Read rres file header
	err := binary.Read(reader, binary.LittleEndian, &fileHeader)
	if err != nil {
		rl.TraceLog(rl.LogWarning, err.Error())
		return
	}

	// Verify "rRES" identifier
	if string(fileHeader.ID[:]) != "rRES" {
		rl.TraceLog(rl.LogWarning, "not a valid raylib resource file")
		return
	}

	reader.Seek(int64(unsafe.Sizeof(fileHeader)), os.SEEK_CUR)

	for i := 0; i < int(fileHeader.Count); i++ {
		// Read resource info and parameters
		err = binary.Read(reader, binary.LittleEndian, &infoHeader)
		if err != nil {
			rl.TraceLog(rl.LogWarning, err.Error())
			return
		}

		reader.Seek(int64(unsafe.Sizeof(infoHeader)), os.SEEK_CUR)

		if int(infoHeader.ID) == rresID {
			data.Type = uint32(infoHeader.DataType)
			data.Param1 = infoHeader.Param1
			data.Param2 = infoHeader.Param2
			data.Param3 = infoHeader.Param3
			data.Param4 = infoHeader.Param4

			// Read resource data block
			b := make([]byte, infoHeader.DataSize)
			reader.Read(b)

			// Decompress data
			data.Data, err = Decompress(b, int(infoHeader.CompType))
			if err != nil {
				rl.TraceLog(rl.LogWarning, "[ID %d] %v", infoHeader.ID, err)
			}

			// Decrypt data
			data.Data, err = Decrypt(key, data.Data, int(infoHeader.CryptoType))
			if err != nil {
				rl.TraceLog(rl.LogWarning, "[ID %d] %v", infoHeader.ID, err)
			}

			if data.Data != nil && len(data.Data) == int(infoHeader.UncompSize) {
				rl.TraceLog(rl.LogInfo, "[ID %d] Resource data loaded successfully", infoHeader.ID)
			}
		} else {
			// Skip required data to read next resource infoHeader
			reader.Seek(int64(infoHeader.DataSize), os.SEEK_CUR)
		}
	}

	if data.Data == nil {
		rl.TraceLog(rl.LogWarning, "[ID %d] Requested resource could not be found", rresID)
	}

	return
}

// Encrypt data
func Encrypt(key, data []byte, cryptoType int) ([]byte, error) {
	switch cryptoType {
	case CryptoXOR:
		c, err := xor.NewXor(string(key))
		if err != nil {
			return nil, err
		}

		return c.Encode(data), nil
	case CryptoAES:
		b, err := encryptAES(key, data)
		if err != nil {
			return nil, err
		}

		return b, nil
	case Crypto3DES:
		b, err := encrypt3DES(key, data)
		if err != nil {
			return nil, err
		}

		return b, nil
	case CryptoBlowfish:
		b, err := encryptBlowfish(key, data)
		if err != nil {
			return nil, err
		}

		return b, nil
	case CryptoXTEA:
		b, err := encryptXTEA(key, data)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		return b, nil
	default:
		return data, nil
	}
}

// Decrypt data
func Decrypt(key, data []byte, cryptoType int) ([]byte, error) {
	switch cryptoType {
	case CryptoXOR:
		c, err := xor.NewXor(string(key))
		if err != nil {
			return nil, err
		}

		b := c.Encode(data)
		return b, nil
	case CryptoAES:
		b, err := decryptAES(key, data)
		if err != nil {
			return nil, err
		}
		return b, nil
	case Crypto3DES:
		b, err := decrypt3DES(key, data)
		if err != nil {
			return nil, err
		}
		return b, nil
	case CryptoBlowfish:
		b, err := decryptBlowfish(key, data)
		if err != nil {
			return nil, err
		}
		return b, nil
	case CryptoXTEA:
		b, err := decryptXTEA(key, data)
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		return data, nil
	}
}

// Compress data
func Compress(data []byte, compType int) ([]byte, error) {
	switch compType {
	case CompNone:
		return data, nil
	case CompDeflate:
		buf := new(bytes.Buffer)

		w, err := flate.NewWriter(buf, flate.DefaultCompression)
		if err != nil {
			return nil, err
		}

		_, err = w.Write(data)
		if err != nil {
			return nil, err
		}

		w.Close()

		return buf.Bytes(), nil
	case CompLZ4:
		buf := new(bytes.Buffer)

		w := lz4.NewWriter(buf)

		_, err := w.Write(data)
		if err != nil {
			return nil, err
		}

		w.Close()

		return buf.Bytes(), nil
	case CompLZMA2:
		buf := new(bytes.Buffer)

		w, err := xz.NewWriter(buf)
		if err != nil {
			return nil, err
		}

		_, err = w.Write(data)
		if err != nil {
			return nil, err
		}

		w.Close()

		return buf.Bytes(), nil
	case CompBZIP2:
		buf := new(bytes.Buffer)

		w, err := bzip2.NewWriter(buf, &bzip2.WriterConfig{Level: bzip2.BestCompression})
		if err != nil {
			return nil, err
		}

		_, err = w.Write(data)
		if err != nil {
			return nil, err
		}

		w.Close()

		return buf.Bytes(), nil
	case CompSnappy:
		buf := new(bytes.Buffer)

		w := snappy.NewWriter(buf)

		_, err := w.Write(data)
		if err != nil {
			return nil, err
		}

		w.Close()

		return buf.Bytes(), nil
	default:
		return data, nil
	}
}

// Decompress data
func Decompress(data []byte, compType int) ([]byte, error) {
	switch compType {
	case CompNone:
		return data, nil
	case CompDeflate:
		r := flate.NewReader(bytes.NewReader(data))

		u, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		r.Close()

		return u, nil
	case CompLZ4:
		r := lz4.NewReader(bytes.NewReader(data))

		u, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return u, nil
	case CompLZMA2:
		r, err := xz.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}

		u, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return u, nil
	case CompBZIP2:
		r, err := bzip2.NewReader(bytes.NewReader(data), &bzip2.ReaderConfig{})
		if err != nil {
			return nil, err
		}

		u, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return u, nil
	case CompSnappy:
		r := snappy.NewReader(bytes.NewReader(data))

		u, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return u, nil
	default:
		return data, nil
	}
}

// pad to block size
func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// unpad
func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, fmt.Errorf("unpad error. This can happen when incorrect encryption key is used.")
	}

	return src[:(length - unpadding)], nil
}

// encryptAES
func encryptAES(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	msg := pad(text, aes.BlockSize)
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], msg)

	return ciphertext, nil
}

// decryptAES
func decryptAES(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if (len(text) % aes.BlockSize) != 0 {
		return nil, fmt.Errorf("blocksize must be multiple of decoded message length")
	}

	iv := text[:aes.BlockSize]
	msg := text[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := unpad(msg)
	if err != nil {
		return nil, err
	}

	return unpadMsg, nil
}

// encrypt3DES
func encrypt3DES(key, text []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	msg := pad(text, des.BlockSize)
	ciphertext := make([]byte, des.BlockSize+len(msg))
	iv := ciphertext[:des.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[des.BlockSize:], msg)

	return ciphertext, nil
}

// decrypt3DES
func decrypt3DES(key, text []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if (len(text) % des.BlockSize) != 0 {
		return nil, fmt.Errorf("blocksize must be multiple of decoded message length")
	}

	iv := text[:des.BlockSize]
	msg := text[des.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(msg, msg)

	unpadMsg, err := unpad(msg)
	if err != nil {
		return nil, err
	}

	return unpadMsg, nil
}

// encryptBlowfish
func encryptBlowfish(key, text []byte) ([]byte, error) {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	msg := pad(text, blowfish.BlockSize)
	ciphertext := make([]byte, blowfish.BlockSize+len(msg))
	iv := ciphertext[:blowfish.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[blowfish.BlockSize:], msg)

	return ciphertext, nil
}

// decryptBlowfish
func decryptBlowfish(key, text []byte) ([]byte, error) {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if (len(text) % blowfish.BlockSize) != 0 {
		return nil, fmt.Errorf("blocksize must be multiple of decoded message length")
	}

	iv := text[:blowfish.BlockSize]
	msg := text[blowfish.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(msg, msg)

	unpadMsg, err := unpad(msg)
	if err != nil {
		return nil, err
	}

	return unpadMsg, nil
}

// encryptXTEA
func encryptXTEA(key, text []byte) ([]byte, error) {
	block, err := xtea.NewCipher(key)
	if err != nil {
		return nil, err
	}

	msg := pad(text, xtea.BlockSize)
	ciphertext := make([]byte, xtea.BlockSize+len(msg))
	iv := ciphertext[:xtea.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[xtea.BlockSize:], msg)

	return ciphertext, nil
}

// decryptXTEA
func decryptXTEA(key, text []byte) ([]byte, error) {
	block, err := xtea.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if (len(text) % xtea.BlockSize) != 0 {
		return nil, fmt.Errorf("blocksize must be multiple of decoded message length")
	}

	iv := text[:xtea.BlockSize]
	msg := text[xtea.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(msg, msg)

	unpadMsg, err := unpad(msg)
	if err != nil {
		return nil, err
	}

	return unpadMsg, nil
}
