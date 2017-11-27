// rREM - raylib Resource EMbedder
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/blezek/tga"
	_ "github.com/jbuchbinder/gopnm"
	_ "golang.org/x/image/bmp"

	"github.com/dsnet/compress/bzip2"
	"github.com/jfreymuth/oggvorbis"
	"github.com/jteeuwen/go-bindata"
	"github.com/klauspost/compress/flate"
	"github.com/moutend/go-wav"
	"github.com/pierrec/lz4"
	"github.com/rootlch/encrypt"
	"github.com/ulikunitz/xz"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/xtea"

	"github.com/gen2brain/raylib-go/rres"
)

func init() {
	tga.RegisterFormat()
}

func main() {
	base := flag.String("base", "data", "Resources file basename")
	comp := flag.Int("comp", rres.CompLZMA2, "Compression type, 0=NONE, 1=DEFLATE, 2=LZ4, 5=LZMA2 (XZ), 6=BZIP2")
	enc := flag.Int("enc", rres.CryptoNone, "Encryption type, 0=NONE, 1=XOR, 2=AES, 3=3DES, 4=Blowfish, 5=XTEA")
	key := flag.String("key", "", "Encryption key")
	header := flag.Bool("header", false, "Generate C header (.h file)")
	source := flag.Bool("source", false, "Generate C source (.c file)")
	bin := flag.Bool("bin", false, "Generate Go bindata (.go file)")
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	switch *comp {
	case rres.CompNone:
	case rres.CompDeflate:
	case rres.CompLZ4:
	case rres.CompLZMA2:
	case rres.CompBZIP2:
	default:
		fmt.Printf("compression type %d not implemented\n", *comp)
		os.Exit(1)
	}

	switch *enc {
	case rres.CryptoNone:
	case rres.CryptoXOR:
	case rres.CryptoAES:
	case rres.Crypto3DES:
	case rres.CryptoBlowfish:
	case rres.CryptoXTEA:
	default:
		fmt.Printf("encryption type %d not implemented\n", *enc)
		os.Exit(1)
	}

	if *enc != 0 {
		if *key == "" {
			fmt.Printf("encryption requires key (-k)\n")
			os.Exit(1)
		}
		if len(*key) != 16 && len(*key) != 24 {
			fmt.Printf("wrong key length, it should be 16 or 24\n")
			os.Exit(1)
		}
	}

	rresFile, err := os.Create(fmt.Sprintf("%s.rres", *base))
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	defer rresFile.Close()

	var headerFile *os.File
	if *header {
		headerFile, err = os.Create(fmt.Sprintf("%s.h", *base))
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		defer headerFile.Close()
	}

	var sourceFile *os.File
	if *source {
		sourceFile, err = os.Create(fmt.Sprintf("%s.c", *base))
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		defer sourceFile.Close()
	}

	var fileHeader rres.FileHeader

	// "rRES" identifier
	copy(fileHeader.ID[:], "rRES")
	fileHeader.Count = uint16(len(flag.Args()))
	fileHeader.Version = 1

	// Write file header
	err = binary.Write(rresFile, binary.LittleEndian, &fileHeader)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	rresFile.Seek(int64(unsafe.Sizeof(fileHeader)), os.SEEK_CUR)

	if *header {
		// Write C header file
		_, err = headerFile.Write([]byte(fmt.Sprintf("#define NUM_RESOURCES %d\n\n", flag.NArg())))
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}

	for id, filename := range flag.Args() {
		var data []byte
		var infoHeader rres.InfoHeader

		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("%s: %v\n", filename, err)
			continue
		}

		data, err = ioutil.ReadAll(file)
		if err != nil {
			fmt.Printf("%s: %v\n", filename, err)
		}

		file.Close()

		infoHeader.ID = uint32(id)
		infoHeader.CompType = uint8(*comp)
		infoHeader.CryptoType = uint8(*enc)
		infoHeader.DataType = uint8(fileType(filename))
		infoHeader.PartsCount = uint8(1)

		// Params
		switch infoHeader.DataType {
		case rres.TypeImage:
			img, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				fmt.Printf("%s: %v\n", filename, err)
				continue
			}

			rect := img.Bounds()
			width, height := rect.Dx(), rect.Dy()

			infoHeader.Param1 = uint32(width)
			infoHeader.Param2 = uint32(height)

			switch img.ColorModel() {
			case color.GrayModel:
				infoHeader.Param3 = rres.ImUncompGrayscale

				i := image.NewGray(rect)
				draw.Draw(i, rect, img, rect.Min, draw.Src)
				data = i.Pix
			case color.Gray16Model:
				infoHeader.Param3 = rres.ImUncompGrayAlpha

				i := image.NewGray16(rect)
				draw.Draw(i, rect, img, rect.Min, draw.Src)
				data = i.Pix
			default:
				infoHeader.Param3 = rres.ImUncompR8g8b8a8

				i := image.NewNRGBA(rect)
				draw.Draw(i, rect, img, rect.Min, draw.Src)
				data = i.Pix
			}

		case rres.TypeWave:
			a := &wav.File{}
			err := wav.Unmarshal(data, a)
			if err != nil {
				fmt.Printf("%s: %v\n", filename, err)
			}

			data, err = ioutil.ReadAll(a)
			if err != nil {
				fmt.Printf("%s: %v\n", filename, err)
			}

			infoHeader.Param1 = uint32(a.Samples())
			infoHeader.Param2 = uint32(a.SamplesPerSec())
			infoHeader.Param3 = uint32(a.BitsPerSample())
			infoHeader.Param4 = uint32(a.Channels())
		case rres.TypeVorbis:
			r, err := oggvorbis.NewReader(bytes.NewReader(data))
			if err != nil {
				fmt.Printf("%s: %v\n", filename, err)
			}

			d, _, err := oggvorbis.ReadAll(bytes.NewReader(data))
			if err != nil {
				fmt.Printf("%s: %v\n", filename, err)
			}

			// Convert []float32 to []byte
			header := *(*reflect.SliceHeader)(unsafe.Pointer(&d))
			header.Len *= 4
			header.Cap *= 4
			data = *(*[]byte)(unsafe.Pointer(&header))

			infoHeader.Param1 = uint32(r.SampleRate())
			infoHeader.Param2 = uint32(r.Bitrate().Nominal)
			infoHeader.Param3 = uint32(r.Channels())
		case rres.TypeVertex:
			// TODO https://github.com/sheenobu/go-obj
		case rres.TypeText, rres.TypeRaw:
		}

		// Encryption
		switch infoHeader.CryptoType {
		case rres.CryptoXOR:
			c, err := encrypt.NewXor(*key)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			b := c.Encode(data)
			data = b
		case rres.CryptoAES:
			b, err := encryptAES([]byte(*key), data)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			data = b
		case rres.Crypto3DES:
			b, err := encrypt3DES([]byte(*key), data)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			data = b
		case rres.CryptoBlowfish:
			b, err := encryptBlowfish([]byte(*key), data)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			data = b
		case rres.CryptoXTEA:
			b, err := encryptXTEA([]byte(*key), data)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			data = b
		}

		infoHeader.UncompSize = uint32(len(data))

		// Compression
		switch infoHeader.CompType {
		case rres.CompNone:
			infoHeader.DataSize = uint32(len(data))
		case rres.CompDeflate:
			buf := new(bytes.Buffer)

			w, err := flate.NewWriter(buf, flate.BestCompression)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			_, err = w.Write(data)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			w.Close()

			infoHeader.DataSize = uint32(len(buf.Bytes()))
			data = buf.Bytes()
		case rres.CompLZ4:
			buf := new(bytes.Buffer)

			w := lz4.NewWriter(buf)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			_, err = w.Write(data)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			w.Close()

			infoHeader.DataSize = uint32(len(buf.Bytes()))
			data = buf.Bytes()
		case rres.CompLZMA2:
			buf := new(bytes.Buffer)

			w, err := xz.NewWriter(buf)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			_, err = w.Write(data)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			w.Close()

			infoHeader.DataSize = uint32(len(buf.Bytes()))
			data = buf.Bytes()
		case rres.CompBZIP2:
			buf := new(bytes.Buffer)

			w, err := bzip2.NewWriter(buf, &bzip2.WriterConfig{Level: bzip2.BestCompression})
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			_, err = w.Write(data)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			w.Close()

			infoHeader.DataSize = uint32(len(buf.Bytes()))
			data = buf.Bytes()
		}

		// Write resource info and parameters
		err = binary.Write(rresFile, binary.LittleEndian, &infoHeader)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		rresFile.Seek(int64(unsafe.Sizeof(infoHeader)), os.SEEK_CUR)

		// Write resource data
		_, err = rresFile.Write(data)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		var typeName string
		switch infoHeader.DataType {
		case rres.TypeImage:
			typeName = "IMAGE"
		case rres.TypeWave:
			typeName = "WAVE"
		case rres.TypeVorbis:
			typeName = "VORBIS"
		case rres.TypeText:
			typeName = "TEXT"
		default:
			typeName = "RAW"
		}

		fmt.Printf("%s %d // Embedded as %s\n", filepath.Base(filename), id, typeName)

		if *header {
			headerFile.Write([]byte(fmt.Sprintf("#define RES_%s 0x%08x\t\t// Embedded as %s\n", filepath.Base(filename), id, typeName)))
		}
	}

	// Generate C source
	if *source {
		rresFile.Seek(0, 0)
		d, err := ioutil.ReadAll(rresFile)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		_, err = sourceFile.Write([]byte("// This file has been automatically generated by rREM - raylib Resource Embedder\n\n"))
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		_, err = sourceFile.Write([]byte(fmt.Sprintf("const unsigned char data[%d] = {\n    ", len(d))))
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		blCounter := 0 // break line counter

		for i := 0; i < len(d)-1; i++ {
			blCounter++

			_, err = sourceFile.Write([]byte(fmt.Sprintf("0x%.2x, ", d[i])))
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			if blCounter >= 24 {
				_, err = sourceFile.Write([]byte("\n    "))
				if err != nil {
					fmt.Printf("%v\n", err)
				}

				blCounter = 0
			}
		}

		_, err = sourceFile.Write([]byte(fmt.Sprintf("0x%.2x };\n", d[len(d)-1])))
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	// Generate bindata
	if *bin {
		cfg := bindata.NewConfig()
		cfg.NoCompress = true
		cfg.Output = fmt.Sprintf("%s.go", *base)
		cfg.Input = make([]bindata.InputConfig, 1)
		cfg.Input[0] = bindata.InputConfig{Path: fmt.Sprintf("%s.rres", *base), Recursive: false}

		err := bindata.Translate(cfg)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}

// fileType returns resource file type
func fileType(f string) int {
	switch strings.ToLower(filepath.Ext(f)) {
	case ".jpg", ".jpeg", ".png", ".bmp", ".tga", ".gif":
		return rres.TypeImage
	case ".txt", ".csv", ".info", ".md":
		return rres.TypeText
	case ".wav":
		return rres.TypeWave
	case ".ogg":
		return rres.TypeVorbis
	case ".obj":
		return rres.TypeVertex
	default:
		return rres.TypeRaw
	}
}

// pad to block size
func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
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
