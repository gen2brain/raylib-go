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
	xor "github.com/rootlch/encrypt"
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

	if !validComp(*comp) {
		fmt.Printf("compression type %d not implemented\n", *comp)
		os.Exit(1)
	}

	if !validEnc(*enc) {
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
		data, infoHeader.Param1, infoHeader.Param2, infoHeader.Param3, infoHeader.Param4, err = params(data, int(infoHeader.DataType))
		if err != nil {
			fmt.Printf("%s: %v\n", filename, err)
		}

		// Encryption
		data, err = encrypt([]byte(*key), data, int(infoHeader.CryptoType))
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		infoHeader.UncompSize = uint32(len(data))

		// Compression
		data, err = compress(data, int(infoHeader.CompType))
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		infoHeader.DataSize = uint32(len(data))

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

		fmt.Printf("%s %d // Embedded as %s\n", filepath.Base(filename), id, typeName(int(infoHeader.DataType)))

		if *header {
			headerFile.Write([]byte(fmt.Sprintf("#define RES_%s 0x%08x\t\t// Embedded as %s\n", filepath.Base(filename), id, typeName(int(infoHeader.DataType)))))
		}
	}

	err = rresFile.Sync()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	err = rresFile.Close()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// Generate C source
	if *source {
		fname := fmt.Sprintf("%s.rres", *base)
		file, err := os.Open(fname)
		if err != nil {
			fmt.Printf("%s: %v\n", fname, err)
		}

		d, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Printf("%s: %v\n", fname, err)
		}

		file.Close()

		err = genSource(sourceFile, d)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	// Generate Go bindata
	if *bin {
		err = genBin(*base)
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

// typeName returns data type name
func typeName(dataType int) string {
	switch dataType {
	case rres.TypeImage:
		return "IMAGE"
	case rres.TypeWave:
		return "WAVE"
	case rres.TypeVorbis:
		return "VORBIS"
	case rres.TypeText:
		return "TEXT"
	default:
		return "RAW"
	}
}

// validEnc checks if encryption type is valid
func validEnc(encType int) bool {
	switch encType {
	case rres.CryptoNone, rres.CryptoXOR:
		return true
	case rres.CryptoAES, rres.Crypto3DES:
		return true
	case rres.CryptoBlowfish, rres.CryptoXTEA:
		return true
	}
	return false
}

// validComp checks if compression type is valid
func validComp(compType int) bool {
	switch compType {
	case rres.CompNone, rres.CompDeflate:
		return true
	case rres.CompLZ4, rres.CompLZMA2:
		return true
	case rres.CompBZIP2:
		return true
	}
	return false
}

// params returns data params
func params(data []byte, dataType int) (d []byte, p1, p2, p3, p4 uint32, err error) {
	switch dataType {
	case rres.TypeImage:
		var img image.Image

		img, _, err = image.Decode(bytes.NewReader(data))
		if err != nil {
			return
		}

		rect := img.Bounds()
		width, height := rect.Dx(), rect.Dy()

		p1 = uint32(width)
		p2 = uint32(height)

		switch img.ColorModel() {
		case color.GrayModel:
			p3 = rres.ImUncompGrayscale

			i := image.NewGray(rect)
			draw.Draw(i, rect, img, rect.Min, draw.Src)
			d = i.Pix
			return
		case color.Gray16Model:
			p3 = rres.ImUncompGrayAlpha

			i := image.NewGray16(rect)
			draw.Draw(i, rect, img, rect.Min, draw.Src)
			d = i.Pix
			return
		default:
			p3 = rres.ImUncompR8g8b8a8

			i := image.NewNRGBA(rect)
			draw.Draw(i, rect, img, rect.Min, draw.Src)
			d = i.Pix
			return
		}

	case rres.TypeWave:
		a := &wav.File{}
		err = wav.Unmarshal(data, a)
		if err != nil {
			return
		}

		d, err = ioutil.ReadAll(a)
		if err != nil {
			return
		}

		p1 = uint32(a.Samples())
		p2 = uint32(a.SamplesPerSec())
		p3 = uint32(a.BitsPerSample())
		p4 = uint32(a.Channels())
		return
	case rres.TypeVorbis:
		r, e := oggvorbis.NewReader(bytes.NewReader(data))
		if e != nil {
			err = e
			return
		}

		o, _, e := oggvorbis.ReadAll(bytes.NewReader(data))
		if e != nil {
			err = e
			return
		}

		// Convert []float32 to []byte
		header := *(*reflect.SliceHeader)(unsafe.Pointer(&o))
		header.Len *= 4
		header.Cap *= 4
		d = *(*[]byte)(unsafe.Pointer(&header))

		p1 = uint32(r.SampleRate())
		p2 = uint32(r.Bitrate().Nominal)
		p3 = uint32(r.Channels())
		return
	case rres.TypeVertex:
		// TODO https://github.com/sheenobu/go-obj
	case rres.TypeText, rres.TypeRaw:
	}

	return
}

// encrypt data
func encrypt(key, data []byte, cryptoType int) ([]byte, error) {
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

// compress data
func compress(data []byte, compType int) ([]byte, error) {
	switch compType {
	case rres.CompNone:
		return data, nil
	case rres.CompDeflate:
		buf := new(bytes.Buffer)

		w, err := flate.NewWriter(buf, flate.BestCompression)
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
	default:
		return data, nil
	}
}

// genSource generates C source file
func genSource(w io.Writer, data []byte) error {
	length := len(data)

	_, err := w.Write([]byte("// This file has been automatically generated by rREM - raylib Resource Embedder\n\n"))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(fmt.Sprintf("const unsigned char data[%d] = {\n    ", length)))
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	blCounter := 0 // break line counter

	for i := 0; i < len(data)-1; i++ {
		blCounter++

		_, err = w.Write([]byte(fmt.Sprintf("0x%.2x, ", data[i])))
		if err != nil {
			return err
		}

		if blCounter >= 24 {
			_, err = w.Write([]byte("\n    "))
			if err != nil {
				return err
			}

			blCounter = 0
		}
	}

	_, err = w.Write([]byte(fmt.Sprintf("0x%.2x };\n", data[length-1])))
	if err != nil {
		return err
	}

	return nil
}

//genBin generates go-bindata file
func genBin(base string) error {
	cfg := bindata.NewConfig()
	cfg.NoCompress = true
	cfg.Output = fmt.Sprintf("%s.go", base)
	cfg.Input = make([]bindata.InputConfig, 1)
	cfg.Input[0] = bindata.InputConfig{Path: fmt.Sprintf("%s.rres", base), Recursive: false}

	return bindata.Translate(cfg)
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
