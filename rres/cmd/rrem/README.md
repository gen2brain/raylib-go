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
    	Compression type, 0=None, 1=Deflate, 2=LZ4, 5=LZMA2 (XZ), 6=BZIP2 (default 5)
  -enc int
    	Encryption type, 0=None, 1=XOR, 2=AES, 3=3DES
  -header
    	Generate C header (.h file)
  -key string
    	Encryption key
```

### Example

[Example](https://github.com/gen2brain/raylib-go/tree/master/examples/others/resources).
