package raylib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/dsnet/compress/bzip2"
	"github.com/klauspost/compress/flate"
	"github.com/pierrec/lz4"
	"github.com/rootlch/encrypt"
	"github.com/ulikunitz/xz"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/xtea"

	"github.com/gen2brain/raylib-go/rres"
)

// LoadResource - Load resource from file by id
// NOTE: Returns uncompressed data with parameters, search resource by id
func LoadResource(reader io.ReadSeeker, rresID int, key []byte) (data rres.Data) {
	var fileHeader rres.FileHeader
	var infoHeader rres.InfoHeader

	reader.Seek(0, 0)

	// Read rres file header
	err := binary.Read(reader, binary.LittleEndian, &fileHeader)
	if err != nil {
		TraceLog(LogWarning, err.Error())
		return
	}

	// Verify "rRES" identifier
	id := fmt.Sprintf("%c", fileHeader.ID)
	if id != "[r R E S]" {
		TraceLog(LogWarning, "not a valid raylib resource file")
		return
	}

	reader.Seek(int64(unsafe.Sizeof(fileHeader)), os.SEEK_CUR)

	for i := 0; i < int(fileHeader.Count); i++ {
		// Read resource info and parameters
		err = binary.Read(reader, binary.LittleEndian, &infoHeader)
		if err != nil {
			TraceLog(LogWarning, err.Error())
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

			// Uncompress data
			data.Data, err = uncompress(b, int(infoHeader.CompType), int(infoHeader.UncompSize))
			if err != nil {
				TraceLog(LogWarning, "[ID %d] %v", infoHeader.ID, err)
			}

			// Decrypt data
			data.Data, err = decrypt(key, data.Data, int(infoHeader.CryptoType))
			if err != nil {
				TraceLog(LogWarning, "[ID %d] %v", infoHeader.ID, err)
			}

			if data.Data != nil && len(data.Data) == int(infoHeader.UncompSize) {
				TraceLog(LogInfo, "[ID %d] Resource data loaded successfully", infoHeader.ID)
			}
		} else {
			// Skip required data to read next resource infoHeader
			reader.Seek(int64(infoHeader.DataSize), os.SEEK_CUR)
		}
	}

	if data.Data == nil {
		TraceLog(LogInfo, "[ID %d] Requested resource could not be found", rresID)
	}

	return
}

// decrypt data
func decrypt(key, data []byte, cryptoType int) ([]byte, error) {
	switch cryptoType {
	case rres.CryptoXOR:
		c, err := encrypt.NewXor(string(key))
		if err != nil {
			return nil, err
		}

		b := c.Encode(data)
		return b, nil
	case rres.CryptoAES:
		b, err := decryptAES(key, data)
		if err != nil {
			return nil, err
		}
		return b, nil
	case rres.Crypto3DES:
		b, err := decrypt3DES(key, data)
		if err != nil {
			return nil, err
		}
		return b, nil
	case rres.CryptoBlowfish:
		b, err := decryptBlowfish(key, data)
		if err != nil {
			return nil, err
		}
		return b, nil
	case rres.CryptoXTEA:
		b, err := decryptXTEA(key, data)
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		return data, nil
	}
}

// uncompress data
func uncompress(data []byte, compType, uncompSize int) ([]byte, error) {
	switch compType {
	case rres.CompNone:
		return data, nil
	case rres.CompDeflate:
		r := flate.NewReader(bytes.NewReader(data))

		u := make([]byte, uncompSize)
		r.Read(u)

		r.Close()

		return u, nil
	case rres.CompLZ4:
		r := lz4.NewReader(bytes.NewReader(data))

		u := make([]byte, uncompSize)
		r.Read(u)

		return u, nil
	case rres.CompLZMA2:
		r, err := xz.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}

		u := make([]byte, uncompSize)
		r.Read(u)

		return u, nil
	case rres.CompBZIP2:
		r, err := bzip2.NewReader(bytes.NewReader(data), &bzip2.ReaderConfig{})
		if err != nil {
			return nil, err
		}

		u := make([]byte, uncompSize)
		r.Read(u)

		return u, nil
	default:
		return data, nil
	}
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
