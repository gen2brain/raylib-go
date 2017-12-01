package raylib

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/gen2brain/raylib-go/rres"
	"github.com/gen2brain/raylib-go/rres/rlib"
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

			// Decompress data
			data.Data, err = rlib.Decompress(b, int(infoHeader.CompType))
			if err != nil {
				TraceLog(LogWarning, "[ID %d] %v", infoHeader.ID, err)
			}

			// Decrypt data
			data.Data, err = rlib.Decrypt(key, data.Data, int(infoHeader.CryptoType))
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
