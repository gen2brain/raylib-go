## rrem

rREM - raylib Resource EMbedder.

### Usage

```
Usage of ./rrem:
  -base string
    	Resources file basename (default "data")
  -bin
    	Generate Go bindata (.go file)
  -comp int
    	Compression type, 0=NONE, 1=DEFLATE, 2=LZ4, 5=LZMA2 (XZ), 6=BZIP2, 7=Snappy (default 1)
  -enc int
    	Encryption type, 0=NONE, 1=XOR, 2=AES, 3=3DES, 4=Blowfish, 5=XTEA
  -header
    	Generate C header (.h file)
  -key string
    	Encryption key
  -source
    	Generate C source (.c file)
```

### Example

[Example](https://github.com/gen2brain/raylib-go/tree/master/examples/others/resources).
