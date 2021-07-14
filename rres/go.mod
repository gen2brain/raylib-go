module github.com/gen2brain/raylib-go/rres

go 1.16

replace github.com/gen2brain/raylib-go/raylib => ../raylib

require (
	github.com/dsnet/compress v0.0.1
	github.com/frankban/quicktest v1.13.0 // indirect
	github.com/gen2brain/raylib-go/raylib v0.0.0-00010101000000-000000000000
	github.com/golang/snappy v0.0.4
	github.com/klauspost/compress v1.13.1
	github.com/pierrec/lz4 v2.6.1+incompatible
	github.com/rootlch/encrypt v0.0.0-20120717064446-58e7f47a860b
	github.com/ulikunitz/xz v0.5.10
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
)
