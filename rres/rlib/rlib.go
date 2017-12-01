// Package rlib provides Encrypt/Decrypt and Compress/Uncompress functions
package rlib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/xtea"

	"github.com/dsnet/compress/bzip2"
	"github.com/golang/snappy"
	"github.com/klauspost/compress/flate"
	"github.com/pierrec/lz4"
	xor "github.com/rootlch/encrypt"
	"github.com/ulikunitz/xz"

	"github.com/gen2brain/raylib-go/rres"
)

// Encrypt data
func Encrypt(key, data []byte, cryptoType int) ([]byte, error) {
	switch cryptoType {
	case rres.CryptoXOR:
		c, err := xor.NewXor(string(key))
		if err != nil {
			return nil, err
		}

		return c.Encode(data), nil
	case rres.CryptoAES:
		b, err := encryptAES(key, data)
		if err != nil {
			return nil, err
		}

		return b, nil
	case rres.Crypto3DES:
		b, err := encrypt3DES(key, data)
		if err != nil {
			return nil, err
		}

		return b, nil
	case rres.CryptoBlowfish:
		b, err := encryptBlowfish(key, data)
		if err != nil {
			return nil, err
		}

		return b, nil
	case rres.CryptoXTEA:
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
	case rres.CryptoXOR:
		c, err := xor.NewXor(string(key))
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

// Compress data
func Compress(data []byte, compType int) ([]byte, error) {
	switch compType {
	case rres.CompNone:
		return data, nil
	case rres.CompDeflate:
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
	case rres.CompLZ4:
		buf := new(bytes.Buffer)

		w := lz4.NewWriter(buf)

		_, err := w.Write(data)
		if err != nil {
			return nil, err
		}

		w.Close()

		return buf.Bytes(), nil
	case rres.CompLZMA2:
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
	case rres.CompBZIP2:
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
	case rres.CompSnappy:
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
	case rres.CompNone:
		return data, nil
	case rres.CompDeflate:
		r := flate.NewReader(bytes.NewReader(data))

		u, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		r.Close()

		return u, nil
	case rres.CompLZ4:
		r := lz4.NewReader(bytes.NewReader(data))

		u, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return u, nil
	case rres.CompLZMA2:
		r, err := xz.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}

		u, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return u, nil
	case rres.CompBZIP2:
		r, err := bzip2.NewReader(bytes.NewReader(data), &bzip2.ReaderConfig{})
		if err != nil {
			return nil, err
		}

		u, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return u, nil
	case rres.CompSnappy:
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
