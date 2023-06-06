/**********************************************************************************************
*
*   rres-raylib v1.2 - rres loaders specific for raylib data structures
*
*   CONFIGURATION:
*
*   #define RRES_RAYLIB_IMPLEMENTATION
*       Generates the implementation of the library into the included file.
*       If not defined, the library is in header only mode and can be included in other headers
*       or source files without problems. But only ONE file should hold the implementation.
*
*   #define RRES_SUPPORT_COMPRESSION_LZ4
*       Support data compression algorithm LZ4, provided by lz4.h/lz4.c library
*
*   #define RRES_SUPPORT_ENCRYPTION_AES
*       Support data encryption algorithm AES, provided by aes.h/aes.c library
*
*   #define RRES_SUPPORT_ENCRYPTION_XCHACHA20
*       Support data encryption algorithm XChaCha20-Poly1305,
*       provided by monocypher.h/monocypher.c library
*
*   DEPENDENCIES:
*
*     - raylib.h: Data types definition and data loading from memory functions
*                 WARNING: raylib.h MUST be included before including rres-raylib.h
*     - rres.h:   Base implementation of rres specs, required to read rres files and resource chunks
*     - lz4.h:    LZ4 compression support (optional)
*     - aes.h:    AES-256 CTR encryption support (optional)
*     - monocypher.h: for XChaCha20-Poly1305 encryption support (optional) 
*
*   VERSION HISTORY:
*
*     - 1.2 (15-Apr-2023): Updated to monocypher 4.0.1
*     - 1.0 (11-May-2022): Initial implementation release
*
*
*   LICENSE: MIT
*
*   Copyright (c) 2020-2023 Ramon Santamaria (@raysan5)
*
*   Permission is hereby granted, free of charge, to any person obtaining a copy
*   of this software and associated documentation files (the "Software"), to deal
*   in the Software without restriction, including without limitation the rights
*   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
*   copies of the Software, and to permit persons to whom the Software is
*   furnished to do so, subject to the following conditions:
*
*   The above copyright notice and this permission notice shall be included in all
*   copies or substantial portions of the Software.
*
*   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
*   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
*   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
*   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
*   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
*   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
*   SOFTWARE.
*
**********************************************************************************************/

#ifndef RRES_RAYLIB_H
#define RRES_RAYLIB_H

#ifndef RRES_H
    #include "rres.h"
#endif

//----------------------------------------------------------------------------------
// Defines and Macros
//----------------------------------------------------------------------------------
//...

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------
//...

//----------------------------------------------------------------------------------
// Global variables
//----------------------------------------------------------------------------------
//...

//----------------------------------------------------------------------------------
// Module Functions Declaration
//----------------------------------------------------------------------------------
#if defined(__cplusplus)
extern "C" {            // Prevents name mangling of functions
#endif

// rres data loading to raylib data structures
// NOTE: Chunk data must be provided uncompressed/unencrypted
RLAPI void *LoadDataFromResource(rresResourceChunk chunk, unsigned int *size); // Load raw data from rres resource chunk
RLAPI char *LoadTextFromResource(rresResourceChunk chunk);      // Load text data from rres resource chunk
RLAPI Image LoadImageFromResource(rresResourceChunk chunk);     // Load Image data from rres resource chunk
RLAPI Wave LoadWaveFromResource(rresResourceChunk chunk);       // Load Wave data from rres resource chunk
RLAPI Font LoadFontFromResource(rresResourceMulti multi);       // Load Font data from rres resource multiple chunks
RLAPI Mesh LoadMeshFromResource(rresResourceMulti multi);       // Load Mesh data from rres resource multiple chunks

// Unpack resource chunk data (decompres/decrypt data)
// NOTE: Function return 0 on success or other value on failure
RLAPI int UnpackResourceChunk(rresResourceChunk *chunk);        // Unpack resource chunk data (decompress/decrypt)
                                                            
// Set base directory for externally linked data
// NOTE: When resource chunk contains an external link (FourCC: LINK, Type: RRES_DATA_LINK),
// a base directory is required to be prepended to link path
// If not provided, the application path is prepended to link by default 
RLAPI void SetBaseDirectory(const char *baseDir);               // Set base directory for externally linked data

#if defined(__cplusplus)
}
#endif

#endif // RRES_RAYLIB_H

/***********************************************************************************
*
*   RRES RAYLIB IMPLEMENTATION
*
************************************************************************************/

#if defined(RRES_RAYLIB_IMPLEMENTATION)

// Compression/Encryption algorithms supported
// NOTE: They should be the same supported by the rres packaging tool (rrespacker)
// https://github.com/phoboslab/qoi
#include "external/qoi.h"                   // Compression algorithm: QOI (implementation in raylib)

#if defined(RRES_SUPPORT_COMPRESSION_LZ4)
    // https://github.com/lz4/lz4
    #include "external/lz4.h"               // Compression algorithm: LZ4
    #include "external/lz4.c"               // Compression algorithm implementation: LZ4
#endif
#if defined(RRES_SUPPORT_ENCRYPTION_AES)
    // https://github.com/kokke/tiny-AES-c
    #include "external/aes.h"               // Encryption algorithm: AES
    #include "external/aes.c"               // Encryption algorithm implementation: AES
#endif
#if defined(RRES_SUPPORT_ENCRYPTION_XCHACHA20)
    // https://github.com/LoupVaillant/Monocypher
    #include "external/monocypher.h"        // Encryption algorithm: XChaCha20-Poly1305
    #include "external/monocypher.c"        // Encryption algorithm implementation: XChaCha20-Poly1305
#endif

//----------------------------------------------------------------------------------
// Defines and Macros
//----------------------------------------------------------------------------------
//...

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------
//...

//----------------------------------------------------------------------------------
// Global Variables Definition
//----------------------------------------------------------------------------------
static const char *baseDir = NULL;      // Base directory pointer, used on external linked data loading

//----------------------------------------------------------------------------------
// Module specific Functions Declaration
//----------------------------------------------------------------------------------

// Load simple data chunks that are later required by multi-chunk resources
// NOTE: Chunk data must be provided uncompressed/unencrypted
static void *LoadDataFromResourceLink(rresResourceChunk chunk, unsigned int *size);      // Load chunk: RRES_DATA_LINK
static void *LoadDataFromResourceChunk(rresResourceChunk chunk, unsigned int *size);     // Load chunk: RRES_DATA_RAW
static char *LoadTextFromResourceChunk(rresResourceChunk chunk, unsigned int *codeLang); // Load chunk: RRES_DATA_TEXT
static Image LoadImageFromResourceChunk(rresResourceChunk chunk);                        // Load chunk: RRES_DATA_IMAGE

static const char *GetExtensionFromProps(unsigned int ext01, unsigned int ext02);        // Get file extension from RRES_DATA_RAW properties (unsigned int) 
static unsigned int *ComputeMD5(unsigned char *data, int size);                          // Compute MD5 hash code, returns 4 integers array (static)

//----------------------------------------------------------------------------------
// Module Functions Definition
//----------------------------------------------------------------------------------

// Load raw data from rres resource
void *LoadDataFromResource(rresResourceChunk chunk, unsigned int *size)
{
    void *rawData = NULL;

    // Data can be provided in the resource or linked to an external file
    if (rresGetDataType(chunk.info.type) == RRES_DATA_RAW)       // Raw data
    {
        rawData = LoadDataFromResourceChunk(chunk, size);
    }
    else if (rresGetDataType(chunk.info.type) == RRES_DATA_LINK) // Link to external file
    {
        // Get raw data from external linked file
        unsigned int dataSize = 0;
        void *data = LoadDataFromResourceLink(chunk, &dataSize);

        rawData = data;
        *size = dataSize;
    }

    return rawData;
}

// Load text data from rres resource
// NOTE: Text must be NULL terminated
char *LoadTextFromResource(rresResourceChunk chunk)
{
    char *text = NULL;
    int codeLang = 0;

    if (rresGetDataType(chunk.info.type) == RRES_DATA_TEXT)       // Text data
    {
        text = LoadTextFromResourceChunk(chunk, &codeLang);

        // TODO: Consider text code language to load shader or code scripts
    }
    else if (rresGetDataType(chunk.info.type) == RRES_DATA_RAW)   // Raw text file
    {
        unsigned int size = 0;
        text = LoadDataFromResourceChunk(chunk, &size);
    }
    else if (rresGetDataType(chunk.info.type) == RRES_DATA_LINK)  // Link to external file
    {
        // Get raw data from external linked file
        unsigned int dataSize = 0;
        void *data = LoadDataFromResourceLink(chunk, &dataSize);
        text = data;
    }

    return text;
}

// Load Image data from rres resource
Image LoadImageFromResource(rresResourceChunk chunk)
{
    Image image = { 0 };

    if (rresGetDataType(chunk.info.type) == RRES_DATA_IMAGE)          // Image data
    {
        image = LoadImageFromResourceChunk(chunk);
    }
    else if (rresGetDataType(chunk.info.type) == RRES_DATA_RAW)       // Raw image file
    {
        unsigned int dataSize = 0;
        unsigned char *data = LoadDataFromResourceChunk(chunk, &dataSize);

        image = LoadImageFromMemory(GetExtensionFromProps(chunk.data.props[1], chunk.data.props[2]), data, dataSize);

        RL_FREE(data);
    }
    else if (rresGetDataType(chunk.info.type) == RRES_DATA_LINK)      // Link to external file
    {
        // Get raw data from external linked file
        unsigned int dataSize = 0;
        void *data = LoadDataFromResourceLink(chunk, &dataSize);

        // Load image from linked file data
        // NOTE: Function checks internally if the file extension is supported to
        // properly load the data, if it fails it logs the result and image.data = NULL
        image = LoadImageFromMemory(GetFileExtension(chunk.data.raw), data, dataSize);
    }

    return image;
}

// Load Wave data from rres resource
Wave LoadWaveFromResource(rresResourceChunk chunk)
{
    Wave wave = { 0 };

    if (rresGetDataType(chunk.info.type) == RRES_DATA_WAVE)       // Wave data
    {
        if ((chunk.info.compType == RRES_COMP_NONE) && (chunk.info.cipherType == RRES_CIPHER_NONE))
        {
            wave.frameCount = chunk.data.props[0];
            wave.sampleRate = chunk.data.props[1];
            wave.sampleSize = chunk.data.props[2];
            wave.channels = chunk.data.props[3];

            unsigned int size = wave.frameCount*wave.sampleSize/8;
            wave.data = RL_CALLOC(size, 1);
            memcpy(wave.data, chunk.data.raw, size);
        }
        RRES_LOG("RRES: %c%c%c%c: WARNING: Data must be decompressed/decrypted\n", chunk.info.type[0], chunk.info.type[1], chunk.info.type[2], chunk.info.type[3]);
    }
    else if (rresGetDataType(chunk.info.type) == RRES_DATA_RAW)   // Raw wave file
    {
        unsigned int dataSize = 0;
        unsigned char *data = LoadDataFromResourceChunk(chunk, &dataSize);

        wave = LoadWaveFromMemory(GetExtensionFromProps(chunk.data.props[1], chunk.data.props[2]), data, dataSize);

        RL_FREE(data);
    }
    else if (rresGetDataType(chunk.info.type) == RRES_DATA_LINK)  // Link to external file
    {
        // Get raw data from external linked file
        unsigned int dataSize = 0;
        void *data = LoadDataFromResourceLink(chunk, &dataSize);

        // Load wave from linked file data
        // NOTE: Function checks internally if the file extension is supported to
        // properly load the data, if it fails it logs the result and wave.data = NULL
        wave = LoadWaveFromMemory(GetFileExtension(chunk.data.raw), data, dataSize);
    }

    return wave;
}

// Load Font data from rres resource
Font LoadFontFromResource(rresResourceMulti multi)
{
    Font font = { 0 };

    // Font resource consist of (2) chunks:
    //  - RRES_DATA_FONT_GLYPHS: Basic font and glyphs properties/data
    //  - RRES_DATA_IMAGE: Image atlas for the font characters
    if (multi.count >= 2)
    {
        if (rresGetDataType(multi.chunks[0].info.type) == RRES_DATA_FONT_GLYPHS)
        {
            if ((multi.chunks[0].info.compType == RRES_COMP_NONE) && (multi.chunks[0].info.cipherType == RRES_CIPHER_NONE))
            {
                // Load font basic properties from chunk[0]
                font.baseSize = multi.chunks[0].data.props[0];           // Base size (default chars height)
                font.glyphCount = multi.chunks[0].data.props[1];         // Number of characters (glyphs)
                font.glyphPadding = multi.chunks[0].data.props[2];      // Padding around the chars

                font.recs = (Rectangle *)RL_CALLOC(font.glyphCount, sizeof(Rectangle));
                font.glyphs = (GlyphInfo *)RL_CALLOC(font.glyphCount, sizeof(GlyphInfo));

                for (int i = 0; i < font.glyphCount; i++)
                {
                    // Font glyphs info comes as a data blob
                    font.recs[i].x = (float)((rresFontGlyphInfo *)multi.chunks[0].data.raw)[i].x;
                    font.recs[i].y = (float)((rresFontGlyphInfo *)multi.chunks[0].data.raw)[i].y;
                    font.recs[i].width = (float)((rresFontGlyphInfo *)multi.chunks[0].data.raw)[i].width;
                    font.recs[i].height = (float)((rresFontGlyphInfo *)multi.chunks[0].data.raw)[i].height;

                    font.glyphs[i].value = ((rresFontGlyphInfo *)multi.chunks[0].data.raw)[i].value;
                    font.glyphs[i].offsetX = ((rresFontGlyphInfo *)multi.chunks[0].data.raw)[i].offsetX;
                    font.glyphs[i].offsetY = ((rresFontGlyphInfo *)multi.chunks[0].data.raw)[i].offsetY;
                    font.glyphs[i].advanceX = ((rresFontGlyphInfo *)multi.chunks[0].data.raw)[i].advanceX;

                    // NOTE: font.glyphs[i].image is not loaded
                }
            }
            else RRES_LOG("RRES: %s: WARNING: Data must be decompressed/decrypted\n", multi.chunks[0].info.type);
        }

        // Load font image chunk
        if (rresGetDataType(multi.chunks[1].info.type) == RRES_DATA_IMAGE)
        {
            if ((multi.chunks[0].info.compType == RRES_COMP_NONE) && (multi.chunks[0].info.cipherType == RRES_CIPHER_NONE))
            {
                Image image = LoadImageFromResourceChunk(multi.chunks[1]);
                font.texture = LoadTextureFromImage(image);
                UnloadImage(image);
            }
            else RRES_LOG("RRES: %s: WARNING: Data must be decompressed/decrypted\n", multi.chunks[1].info.type);
        }
    }
    else    // One chunk of data: RRES_DATA_RAW or RRES_DATA_LINK?
    {
        if (rresGetDataType(multi.chunks[0].info.type) == RRES_DATA_RAW)      // Raw font file
        {
            unsigned int dataSize = 0;
            unsigned char *rawData = LoadDataFromResourceChunk(multi.chunks[0], &dataSize);

            font = LoadFontFromMemory(GetExtensionFromProps(multi.chunks[0].data.props[1], multi.chunks[0].data.props[2]), rawData, dataSize, 32, NULL, 0);

            RL_FREE(rawData);
        }
        if (rresGetDataType(multi.chunks[0].info.type) == RRES_DATA_LINK)     // Link to external font file
        {
            // Get raw data from external linked file
            unsigned int dataSize = 0;
            void *rawData = LoadDataFromResourceLink(multi.chunks[0], &dataSize);

            // Load image from linked file data
            // NOTE 1: Loading font at 32px base size and default charset (95 glyphs)
            // NOTE 2: Function checks internally if the file extension is supported to
            // properly load the data, if it fails it logs the result and font.texture.id = 0
            font = LoadFontFromMemory(GetFileExtension(multi.chunks[0].data.raw), rawData, dataSize, 32, NULL, 0);

            RRES_FREE(rawData);
        }
    }

    return font;
}

// Load Mesh data from rres resource
// NOTE: We try to load vertex data following raylib structure constraints,
// in case data does not fit raylib Mesh structure, it is not loaded
Mesh LoadMeshFromResource(rresResourceMulti multi)
{
    Mesh mesh = { 0 };

    // TODO: Support externally linked mesh resource?

    // Mesh resource consist of (n) chunks:
    for (unsigned int i = 0; i < multi.count; i++)
    {
        if ((multi.chunks[0].info.compType == RRES_COMP_NONE) && (multi.chunks[0].info.cipherType == RRES_CIPHER_NONE))
        {
            // NOTE: raylib only supports vertex arrays with same vertex count,
            // rres.chunks[0] defined vertexCount will be the reference for the following chunks
            // The only exception to vertexCount is the mesh.indices array
            if (mesh.vertexCount == 0) mesh.vertexCount = multi.chunks[0].data.props[0];

            // Verify chunk type and vertex count
            if (rresGetDataType(multi.chunks[i].info.type) == RRES_DATA_VERTEX)
            {
                // In case vertex count do not match we skip that resource chunk
                if ((multi.chunks[i].data.props[1] != RRES_VERTEX_ATTRIBUTE_INDEX) && (multi.chunks[i].data.props[0] != mesh.vertexCount)) continue;

                // NOTE: We are only loading raylib supported rresVertexFormat and raylib expected components count
                switch (multi.chunks[i].data.props[1])    // Check rresVertexAttribute value
                {
                    case RRES_VERTEX_ATTRIBUTE_POSITION:
                    {
                        // raylib expects 3 components per vertex and float vertex format
                        if ((multi.chunks[i].data.props[2] == 3) && (multi.chunks[i].data.props[3] == RRES_VERTEX_FORMAT_FLOAT))
                        {
                            mesh.vertices = (float *)RL_CALLOC(mesh.vertexCount*3, sizeof(float));
                            memcpy(mesh.vertices, multi.chunks[i].data.raw, mesh.vertexCount*3*sizeof(float));
                        }
                        else RRES_LOG("RRES: WARNING: MESH: Vertex attribute position not valid, componentCount/vertexFormat do not fit\n");

                    } break;
                    case RRES_VERTEX_ATTRIBUTE_TEXCOORD1:
                    {
                        // raylib expects 2 components per vertex and float vertex format
                        if ((multi.chunks[i].data.props[2] == 2) && (multi.chunks[i].data.props[3] == RRES_VERTEX_FORMAT_FLOAT))
                        {
                            mesh.texcoords = (float *)RL_CALLOC(mesh.vertexCount*2, sizeof(float));
                            memcpy(mesh.texcoords, multi.chunks[i].data.raw, mesh.vertexCount*2*sizeof(float));
                        }
                        else RRES_LOG("RRES: WARNING: MESH: Vertex attribute texcoord1 not valid, componentCount/vertexFormat do not fit\n");

                    } break;
                    case RRES_VERTEX_ATTRIBUTE_TEXCOORD2:
                    {
                        // raylib expects 2 components per vertex and float vertex format
                        if ((multi.chunks[i].data.props[2] == 2) && (multi.chunks[i].data.props[3] == RRES_VERTEX_FORMAT_FLOAT))
                        {
                            mesh.texcoords2 = (float *)RL_CALLOC(mesh.vertexCount*2, sizeof(float));
                            memcpy(mesh.texcoords2, multi.chunks[i].data.raw, mesh.vertexCount*2*sizeof(float));
                        }
                        else RRES_LOG("RRES: WARNING: MESH: Vertex attribute texcoord2 not valid, componentCount/vertexFormat do not fit\n");

                    } break;
                    case RRES_VERTEX_ATTRIBUTE_TEXCOORD3:
                    {
                        RRES_LOG("RRES: WARNING: MESH: Vertex attribute texcoord3 not supported\n");

                    } break;
                    case RRES_VERTEX_ATTRIBUTE_TEXCOORD4:
                    {
                        RRES_LOG("RRES: WARNING: MESH: Vertex attribute texcoord4 not supported\n");

                    } break;
                    case RRES_VERTEX_ATTRIBUTE_NORMAL:
                    {
                        // raylib expects 3 components per vertex and float vertex format
                        if ((multi.chunks[i].data.props[2] == 3) && (multi.chunks[i].data.props[3] == RRES_VERTEX_FORMAT_FLOAT))
                        {
                            mesh.normals = (float *)RL_CALLOC(mesh.vertexCount*3, sizeof(float));
                            memcpy(mesh.normals, multi.chunks[i].data.raw, mesh.vertexCount*3*sizeof(float));
                        }
                        else RRES_LOG("RRES: WARNING: MESH: Vertex attribute normal not valid, componentCount/vertexFormat do not fit\n");

                    } break;
                    case RRES_VERTEX_ATTRIBUTE_TANGENT:
                    {
                        // raylib expects 4 components per vertex and float vertex format
                        if ((multi.chunks[i].data.props[2] == 4) && (multi.chunks[i].data.props[3] == RRES_VERTEX_FORMAT_FLOAT))
                        {
                            mesh.tangents = (float *)RL_CALLOC(mesh.vertexCount*4, sizeof(float));
                            memcpy(mesh.tangents, multi.chunks[i].data.raw, mesh.vertexCount*4*sizeof(float));
                        }
                        else RRES_LOG("RRES: WARNING: MESH: Vertex attribute tangent not valid, componentCount/vertexFormat do not fit\n");

                    } break;
                    case RRES_VERTEX_ATTRIBUTE_COLOR:
                    {
                        // raylib expects 4 components per vertex and unsigned char vertex format
                        if ((multi.chunks[i].data.props[2] == 4) && (multi.chunks[i].data.props[3] == RRES_VERTEX_FORMAT_UBYTE))
                        {
                            mesh.colors = (unsigned char *)RL_CALLOC(mesh.vertexCount*4, sizeof(unsigned char));
                            memcpy(mesh.colors, multi.chunks[i].data.raw, mesh.vertexCount*4*sizeof(unsigned char));
                        }
                        else RRES_LOG("RRES: WARNING: MESH: Vertex attribute color not valid, componentCount/vertexFormat do not fit\n");

                    } break;
                    case RRES_VERTEX_ATTRIBUTE_INDEX:
                    {
                        // raylib expects 1 components per index and unsigned short vertex format
                        if ((multi.chunks[i].data.props[2] == 1) && (multi.chunks[i].data.props[3] == RRES_VERTEX_FORMAT_USHORT))
                        {
                            mesh.indices = (unsigned short *)RL_CALLOC(multi.chunks[i].data.props[0], sizeof(unsigned short));
                            memcpy(mesh.indices, multi.chunks[i].data.raw, multi.chunks[i].data.props[0]*sizeof(unsigned short));
                        }
                        else RRES_LOG("RRES: WARNING: MESH: Vertex attribute index not valid, componentCount/vertexFormat do not fit\n");

                    } break;
                    default: break;
                }
            }
        }
        else RRES_LOG("RRES: WARNING: Vertex provided data must be decompressed/decrypted\n");
    }

    return mesh;
}

// Unpack compressed/encrypted data from resource chunk
// In case data could not be processed by rres.h, it is just copied in chunk.data.raw for processing here
// NOTE 1: Function return 0 on success or an error code on failure
// NOTE 2: Data corruption CRC32 check has already been performed by rresLoadResourceMulti() on rres.h
int UnpackResourceChunk(rresResourceChunk *chunk)
{
    int result = 0;
    bool updateProps = false;

    // Result error codes:
    //  0 - No error, decompression/decryption successful
    //  1 - Encryption algorithm not supported
    //  2 - Invalid password on decryption
    //  3 - Compression algorithm not supported
    //  4 - Error on data decompression

    // NOTE 1: If data is compressed/encrypted the properties are not loaded by rres.h because
    // it's up to the user to process the data; *chunk must be properly updated by this function
    // NOTE 2: rres-raylib should support the same algorithms and libraries used by rrespacker tool
    void *unpackedData = NULL;    

    // STEP 1. Data decryption
    //-------------------------------------------------------------------------------------
    unsigned char *decryptedData = NULL;

    switch (chunk->info.cipherType)
    {
        case RRES_CIPHER_NONE: decryptedData = chunk->data.raw; break;
#if defined(RRES_SUPPORT_ENCRYPTION_AES)
        case RRES_CIPHER_AES:
        {
            // WARNING: Implementation dependant!
            // rrespacker tool appends (salt[16] + MD5[16]) to encrypted data for convenience,
            // Actually, chunk->info.packedSize considers those additional elements

            // Get some memory for the possible message output
            decryptedData = (unsigned char *)RL_CALLOC(chunk->info.packedSize - 16 - 16, 1);
            if (decryptedData != NULL) memcpy(decryptedData, chunk->data.raw, chunk->info.packedSize - 16 - 16);

            // Required variables for key stretching
            uint8_t key[32] = { 0 };                    // Encryption key
            uint8_t salt[16] = { 0 };                   // Key stretching salt

            // Retrieve salt from chunk packed data
            // salt is stored at the end of packed data, before nonce and MAC: salt[16] + MD5[16]
            memcpy(salt, ((unsigned char *)chunk->data.raw) + (chunk->info.packedSize - 16 - 16), 16);
            
            // Key stretching configuration
            crypto_argon2_config config = {
                .algorithm = CRYPTO_ARGON2_I,           // Algorithm: Argon2i
                .nb_blocks = 16384,                     // Blocks: 16 MB
                .nb_passes = 3,                         // Iterations
                .nb_lanes  = 1                          // Single-threaded
            };
            crypto_argon2_inputs inputs = {
                .pass = (const uint8_t *)rresGetCipherPassword(),     // User password
                .pass_size = 16,                        // Password length
                .salt = salt,                           // Salt for the password
                .salt_size = 16
            };
            crypto_argon2_extras extras = { 0 };        // Extra parameters unused

            void *workArea = RL_MALLOC(config.nb_blocks*1024);    // Key stretching work area

            // Generate strong encryption key, generated from user password using Argon2i algorithm (256 bit)
            crypto_argon2(key, 32, workArea, config, inputs, extras);

            // Wipe key generation secrets, they are no longer needed
            crypto_wipe(salt, 16);
            RL_FREE(workArea);

            // Required variables for decryption and message authentication
            unsigned int md5[4] = { 0 };                // Message Authentication Code generated on encryption

            // Retrieve MD5 from chunk packed data
            // NOTE: MD5 is stored at the end of packed data, after salt: salt[16] + MD5[16]
            memcpy(md5, ((unsigned char *)chunk->data.raw) + (chunk->info.packedSize - 16), 4*sizeof(unsigned int));

            // Message decryption, requires key
            struct AES_ctx ctx = { 0 };
            AES_init_ctx(&ctx, key);
            AES_CTR_xcrypt_buffer(&ctx, (uint8_t *)decryptedData, chunk->info.packedSize - 16 - 16);   // AES Counter mode, stream cipher

            // Verify MD5 to check if data decryption worked
            unsigned int decryptMD5[4] = { 0 };
            unsigned int *md5Ptr = ComputeMD5(decryptedData, chunk->info.packedSize - 16 - 16);
            for (int i = 0; i < 4; i++) decryptMD5[i] = md5Ptr[i];

            // Wipe secrets if they are no longer needed
            crypto_wipe(key, 32);

            if (memcmp(decryptMD5, md5, 4*sizeof(unsigned int)) == 0)    // Decrypted successfully!
            {
                chunk->info.packedSize -= (16 + 16);    // We remove additional data size from packed size (salt[16] + MD5[16])
                RRES_LOG("RRES: %c%c%c%c: Data decrypted successfully (AES)\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);
            }
            else
            {
                result = 2;    // Data was not decrypted as expected, wrong password or message corrupted
                RRES_LOG("RRES: WARNING: %c%c%c%c: Data decryption failed, wrong password or corrupted data\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);
            }

        } break;
#endif
#if defined(RRES_SUPPORT_ENCRYPTION_XCHACHA20)
        case RRES_CIPHER_XCHACHA20_POLY1305:
        {
            // WARNING: Implementation dependant!
            // rrespacker tool appends (salt[16] + nonce[24] + MAC[16]) to encrypted data for convenience,
            // Actually, chunk->info.packedSize considers those additional elements

            // Get some memory for the possible message output
            decryptedData = (unsigned char *)RL_CALLOC(chunk->info.packedSize - 16 - 24 - 16, 1);

            // Required variables for key stretching
            uint8_t key[32] = { 0 };                    // Encryption key
            uint8_t salt[16] = { 0 };                   // Key stretching salt

            // Retrieve salt from chunk packed data
            // salt is stored at the end of packed data, before nonce and MAC: salt[16] + nonce[24] + MAC[16]
            memcpy(salt, ((unsigned char *)chunk->data.raw) + (chunk->info.packedSize - 16 - 24 - 16), 16);
            
            // Key stretching configuration
            crypto_argon2_config config = {
                .algorithm = CRYPTO_ARGON2_I,           // Algorithm: Argon2i
                .nb_blocks = 16384,                     // Blocks: 16 MB
                .nb_passes = 3,                         // Iterations
                .nb_lanes  = 1                          // Single-threaded
            };
            crypto_argon2_inputs inputs = {
                .pass = (const uint8_t *)rresGetCipherPassword(),     // User password
                .pass_size = 16,                        // Password length
                .salt = salt,                           // Salt for the password
                .salt_size = 16
            };
            crypto_argon2_extras extras = { 0 };        // Extra parameters unused

            void *workArea = RL_MALLOC(config.nb_blocks*1024);    // Key stretching work area

            // Generate strong encryption key, generated from user password using Argon2i algorithm (256 bit)
            crypto_argon2(key, 32, workArea, config, inputs, extras);

            // Wipe key generation secrets, they are no longer needed
            crypto_wipe(salt, 16);
            RL_FREE(workArea);

            // Required variables for decryption and message authentication
            uint8_t nonce[24] = { 0 };                  // nonce used on encryption, unique to processed file
            uint8_t mac[16] = { 0 };                    // Message Authentication Code generated on encryption

            // Retrieve nonce and MAC from chunk packed data
            // nonce and MAC are stored at the end of packed data, after salt: salt[16] + nonce[24] + MAC[16]
            memcpy(nonce, ((unsigned char *)chunk->data.raw) + (chunk->info.packedSize - 16 - 24), 24);
            memcpy(mac, ((unsigned char *)chunk->data.raw) + (chunk->info.packedSize - 16), 16);

            // Message decryption requires key, nonce and MAC
            int decryptResult = crypto_aead_unlock(decryptedData, mac, key, nonce, NULL, 0, chunk->data.raw, (chunk->info.packedSize - 16 - 24 - 16));

            // Wipe secrets if they are no longer needed
            crypto_wipe(nonce, 24);
            crypto_wipe(key, 32);

            if (decryptResult == 0)    // Decrypted successfully!
            {
                chunk->info.packedSize -= (16 + 24 + 16);    // We remove additional data size from packed size
                RRES_LOG("RRES: %c%c%c%c: Data decrypted successfully (XChaCha20)\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);
            }
            else if (decryptResult == -1)
            {
                result = 2;   // Wrong password or message corrupted
                RRES_LOG("RRES: WARNING: %c%c%c%c: Data decryption failed, wrong password or corrupted data\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);
            }
        } break;
#endif
        default: 
        {
            result = 1;    // Decryption algorithm not supported
            RRES_LOG("RRES: WARNING: %c%c%c%c: Chunk data encryption algorithm not supported\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);

        } break;
    }

    if ((result == 0) && (chunk->info.cipherType != RRES_CIPHER_NONE))
    {
        // Data is not encrypted any more, register it
        chunk->info.cipherType = RRES_CIPHER_NONE;
        updateProps = true;
    }

    // STEP 2: Data decompression (if decryption was successful)
    //-------------------------------------------------------------------------------------
    unsigned char *uncompData = NULL;

    if (result == 0)
    {
        switch (chunk->info.compType)
        {
            case RRES_COMP_NONE: unpackedData = decryptedData; break;
            case RRES_COMP_DEFLATE:
            {
                int uncompDataSize = 0;

                // TODO: WARNING: Possible issue with allocators: RL_CALLOC() vs RRES_CALLOC()
                uncompData = DecompressData(decryptedData, chunk->info.packedSize, &uncompDataSize);

                if ((uncompData != NULL) && (uncompDataSize > 0))     // Decompression successful
                {
                    unpackedData = uncompData;
                    chunk->info.packedSize = uncompDataSize;
                    RRES_LOG("RRES: %c%c%c%c: Data decompressed successfully (DEFLATE)\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);
                }
                else
                {
                    result = 4;    // Decompression process failed
                    RRES_LOG("RRES: WARNING: %c%c%c%c: Chunk data decompression failed\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);
                }

                // Security check, uncompDataSize must match the provided chunk->baseSize
                if (uncompDataSize != chunk->info.baseSize) RRES_LOG("RRES: WARNING: Decompressed data could be corrupted, unexpected size\n");
            } break;
#if defined(RRES_SUPPORT_COMPRESSION_LZ4)
            case RRES_COMP_LZ4:
            {
                int uncompDataSize = 0;
                uncompData = (unsigned char *)RRES_CALLOC(chunk->info.baseSize, 1);
                uncompDataSize = LZ4_decompress_safe(decryptedData, uncompData, chunk->info.packedSize, chunk->info.baseSize);

                if ((uncompData != NULL) && (uncompDataSize > 0))     // Decompression successful
                {
                    unpackedData = uncompData;
                    chunk->info.packedSize = uncompDataSize;
                    RRES_LOG("RRES: %c%c%c%c: Data decompressed successfully (LZ4)\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);
                }
                else
                {
                    result = 4;    // Decompression process failed
                    RRES_LOG("RRES: WARNING: %c%c%c%c: Chunk data decompression failed\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);
                }

                // WARNING: Decompression could be successful but not the original message size returned
                if (uncompDataSize != chunk->info.baseSize) RRES_LOG("RRES: WARNING: Decompressed data could be corrupted, unexpected size\n");
            } break;
#endif
            case RRES_COMP_QOI:
            {
                int uncompDataSize = 0;
                qoi_desc desc = { 0 };

                // TODO: WARNING: Possible issue with allocators: QOI_MALLOC() vs RRES_MALLOC()
                uncompData = qoi_decode(decryptedData, chunk->info.packedSize, &desc, 0);
                uncompDataSize = (desc.width*desc.height*desc.channels) + 20;   // Add the 20 bytes of (propCount + props[4])

                if ((uncompData != NULL) && (uncompDataSize > 0))     // Decompression successful
                {
                    unpackedData = uncompData;
                    chunk->info.packedSize = uncompDataSize;
                    RRES_LOG("RRES: %c%c%c%c: Data decompressed successfully (QOI)\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);
                }
                else
                {
                    result = 4;    // Decompression process failed
                    RRES_LOG("RRES: WARNING: %c%c%c%c: Chunk data decompression failed\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);
                }

                if (uncompDataSize != chunk->info.baseSize) RRES_LOG("RRES: WARNING: Decompressed data could be corrupted, unexpected size\n");
            } break;
            default:
            {
                result = 3;
                RRES_LOG("RRES: WARNING: %c%c%c%c: Chunk data compression algorithm not supported\n", chunk->info.type[0], chunk->info.type[1], chunk->info.type[2], chunk->info.type[3]);
            } break;
        }
    }

    if ((result == 0) && (chunk->info.compType != RRES_COMP_NONE))
    {
        // Data is not encrypted any more, register it
        chunk->info.compType = RRES_COMP_NONE;
        updateProps = true;
    }

    // Update chunk->data.propCount and chunk->data.props if required
    if (updateProps && (unpackedData != NULL))
    {
        // Data is decompressed/decrypted into chunk->data.raw but data.propCount and data.props[] are still empty, 
        // they must be filled with the just updated chunk->data.raw (that contains everything)
        chunk->data.propCount = ((int *)unpackedData)[0];

        if (chunk->data.propCount > 0)
        {
            chunk->data.props = (unsigned int *)RRES_CALLOC(chunk->data.propCount, sizeof(int));
            for (unsigned int i = 0; i < chunk->data.propCount; i++) chunk->data.props[i] = ((int *)unpackedData)[1 + i];
        }

        // Move chunk->data.raw pointer (chunk->data.propCount*sizeof(int)) positions
        void *raw = RRES_CALLOC(chunk->info.baseSize - 20, 1);
        if (raw != NULL) memcpy(raw, ((unsigned char *)unpackedData) + 20, chunk->info.baseSize - 20);
        RRES_FREE(chunk->data.raw);
        chunk->data.raw = raw;
        RL_FREE(unpackedData);
    }

    return result;
}

//----------------------------------------------------------------------------------
// Module specific Functions Definition
//----------------------------------------------------------------------------------

// Load data chunk: RRES_DATA_LINK
static void *LoadDataFromResourceLink(rresResourceChunk chunk, unsigned int *size)
{
    unsigned char fullFilePath[2048] = { 0 };
    void *data = NULL;
    *size = 0;

    // Get external link filepath
    unsigned char *linkFilePath = RL_CALLOC(chunk.data.props[0], 1);
    if (linkFilePath != NULL) memcpy(linkFilePath, chunk.data.raw, chunk.data.props[0]);

    // Get base directory to append filepath if not provided by user
    if (baseDir == NULL) baseDir = GetApplicationDirectory();
    
    strcpy(fullFilePath, baseDir);
    strcat(fullFilePath, linkFilePath);

    RRES_LOG("RRES: %c%c%c%c: Data file linked externally: %s\n", chunk.info.type[0], chunk.info.type[1], chunk.info.type[2], chunk.info.type[3], linkFilePath);

    if (FileExists(fullFilePath))
    {
        // Load external file as raw data
        // NOTE: We check if file is a text file to allow automatic line-endings processing
        if (IsFileExtension(linkFilePath, ".txt;.md;.vs;.fs;.info;.c;.h;.json;.xml;.glsl"))     // Text file
        {
            data = LoadFileText(fullFilePath);
            *size = TextLength(data);
        }
        else data = LoadFileData(fullFilePath, size);

        if ((data != NULL) && (*size > 0)) RRES_LOG("RRES: %c%c%c%c: External linked file loaded successfully\n", chunk.info.type[0], chunk.info.type[1], chunk.info.type[2], chunk.info.type[3]);
    }
    else RRES_LOG("RRES: WARNING: [%s] Linked external file could not be found\n", linkFilePath);

    return data;
}

// Load data chunk: RRES_DATA_RAW
// NOTE: This chunk can be used raw files embedding or other binary blobs
static void *LoadDataFromResourceChunk(rresResourceChunk chunk, unsigned int *size)
{
    void *rawData = NULL;

    if ((chunk.info.compType == RRES_COMP_NONE) && (chunk.info.cipherType == RRES_CIPHER_NONE))
    {
        rawData = RL_CALLOC(chunk.data.props[0], 1);
        if (rawData != NULL) memcpy(rawData, chunk.data.raw, chunk.data.props[0]);
        *size = chunk.data.props[0];
    }
    else RRES_LOG("RRES: %c%c%c%c: WARNING: Data must be decompressed/decrypted\n", chunk.info.type[0], chunk.info.type[1], chunk.info.type[2], chunk.info.type[3]);

    return rawData;
}

// Load data chunk: RRES_DATA_TEXT
// NOTE: This chunk can be used for shaders or other text data elements (materials?)
static char *LoadTextFromResourceChunk(rresResourceChunk chunk, unsigned int *codeLang)
{
    void *text = NULL;

    if ((chunk.info.compType == RRES_COMP_NONE) && (chunk.info.cipherType == RRES_CIPHER_NONE))
    {
        text = (char *)RL_CALLOC(chunk.data.props[0] + 1, 1);    // We add NULL terminator, just in case
        if (text != NULL) memcpy(text, chunk.data.raw, chunk.data.props[0]);

        // TODO: We got some extra text properties, in case they could be useful for users:
        // chunk.props[1]:rresTextEncoding, chunk.props[2]:rresCodeLang, chunk. props[3]:cultureCode
        *codeLang = chunk.data.props[2];
        //chunks.props[3]:cultureCode could be useful for localized text
    }
    else RRES_LOG("RRES: %c%c%c%c: WARNING: Data must be decompressed/decrypted\n", chunk.info.type[0], chunk.info.type[1], chunk.info.type[2], chunk.info.type[3]);

    return text;
}

// Load data chunk: RRES_DATA_IMAGE
// NOTE: Many data types use images data in some way (font, material...)
static Image LoadImageFromResourceChunk(rresResourceChunk chunk)
{
    Image image = { 0 };

    if ((chunk.info.compType == RRES_COMP_NONE) && (chunk.info.cipherType == RRES_CIPHER_NONE))
    {
        image.width = chunk.data.props[0];
        image.height = chunk.data.props[1];
        int format = chunk.data.props[2];

        // Assign equivalent pixel formats for our engine
        // NOTE: In this case rresPixelFormat defined values match raylib PixelFormat values
        switch (format)
        {
            case RRES_PIXELFORMAT_UNCOMP_GRAYSCALE: image.format = PIXELFORMAT_UNCOMPRESSED_GRAYSCALE; break;
            case RRES_PIXELFORMAT_UNCOMP_GRAY_ALPHA: image.format = PIXELFORMAT_UNCOMPRESSED_GRAY_ALPHA; break;
            case RRES_PIXELFORMAT_UNCOMP_R5G6B5: image.format = PIXELFORMAT_UNCOMPRESSED_R5G6B5; break;
            case RRES_PIXELFORMAT_UNCOMP_R8G8B8: image.format = PIXELFORMAT_UNCOMPRESSED_R8G8B8; break;
            case RRES_PIXELFORMAT_UNCOMP_R5G5B5A1: image.format = PIXELFORMAT_UNCOMPRESSED_R5G5B5A1; break;
            case RRES_PIXELFORMAT_UNCOMP_R4G4B4A4: image.format = PIXELFORMAT_UNCOMPRESSED_R4G4B4A4; break;
            case RRES_PIXELFORMAT_UNCOMP_R8G8B8A8: image.format = PIXELFORMAT_UNCOMPRESSED_R8G8B8A8; break;
            case RRES_PIXELFORMAT_UNCOMP_R32: image.format = PIXELFORMAT_UNCOMPRESSED_R32; break;
            case RRES_PIXELFORMAT_UNCOMP_R32G32B32: image.format = PIXELFORMAT_UNCOMPRESSED_R32G32B32; break;
            case RRES_PIXELFORMAT_UNCOMP_R32G32B32A32: image.format = PIXELFORMAT_UNCOMPRESSED_R32G32B32A32; break;
            case RRES_PIXELFORMAT_COMP_DXT1_RGB: image.format = PIXELFORMAT_COMPRESSED_DXT1_RGB; break;
            case RRES_PIXELFORMAT_COMP_DXT1_RGBA: image.format = PIXELFORMAT_COMPRESSED_DXT1_RGBA; break;
            case RRES_PIXELFORMAT_COMP_DXT3_RGBA: image.format = PIXELFORMAT_COMPRESSED_DXT3_RGBA; break;
            case RRES_PIXELFORMAT_COMP_DXT5_RGBA: image.format = PIXELFORMAT_COMPRESSED_DXT5_RGBA; break;
            case RRES_PIXELFORMAT_COMP_ETC1_RGB: image.format = PIXELFORMAT_COMPRESSED_ETC1_RGB; break;
            case RRES_PIXELFORMAT_COMP_ETC2_RGB: image.format = PIXELFORMAT_COMPRESSED_ETC2_RGB; break;
            case RRES_PIXELFORMAT_COMP_ETC2_EAC_RGBA: image.format = PIXELFORMAT_COMPRESSED_ETC2_EAC_RGBA; break;
            case RRES_PIXELFORMAT_COMP_PVRT_RGB: image.format = PIXELFORMAT_COMPRESSED_PVRT_RGB; break;
            case RRES_PIXELFORMAT_COMP_PVRT_RGBA: image.format = PIXELFORMAT_COMPRESSED_PVRT_RGBA; break;
            case RRES_PIXELFORMAT_COMP_ASTC_4x4_RGBA: image.format = PIXELFORMAT_COMPRESSED_ASTC_4x4_RGBA; break;
            case RRES_PIXELFORMAT_COMP_ASTC_8x8_RGBA: image.format = PIXELFORMAT_COMPRESSED_ASTC_8x8_RGBA; break;
            default: break;
        }

        image.mipmaps = chunk.data.props[3];

        // Image data size can be computed from image properties
        unsigned int size = GetPixelDataSize(image.width, image.height, image.format);

        // NOTE: Computed image data must match the data size of the chunk processed (minus propCount + props[4] size)
        if (size == (chunk.info.baseSize - 20))
        {
            image.data = RL_CALLOC(size, 1);
            if (image.data != NULL) memcpy(image.data, chunk.data.raw, size);
        }
        else RRES_LOG("RRES: WARNING: IMGE: Chunk data size do not match expected image data size\n");
    }
    else RRES_LOG("RRES: %c%c%c%c: WARNING: Data must be decompressed/decrypted\n", chunk.info.type[0], chunk.info.type[1], chunk.info.type[2], chunk.info.type[3]);

    return image;
}

// Get file extension from RRES_DATA_RAW properties (unsigned int) 
static const char *GetExtensionFromProps(unsigned int ext01, unsigned int ext02)
{
    static char extension[8] = { 0 };
    memset(extension, 0, 8);

    // Convert file extension provided as 2 unsigned int properties, to a char[] array 
    // NOTE: Extension is defined as 2 unsigned int big-endian values (4 bytes each), 
    // starting with a dot, i.e 0x2e706e67 => ".png"
    extension[0] = (unsigned char)((ext01 & 0xff000000) >> 24);
    extension[1] = (unsigned char)((ext01 & 0x00ff0000) >> 16);
    extension[2] = (unsigned char)((ext01 & 0x0000ff00) >> 8);
    extension[3] = (unsigned char)(ext01 & 0x000000ff);

    extension[4] = (unsigned char)((ext02 & 0xff000000) >> 24);
    extension[5] = (unsigned char)((ext02 & 0x00ff0000) >> 16);
    extension[6] = (unsigned char)((ext02 & 0x0000ff00) >> 8);
    extension[7] = (unsigned char)(ext02 & 0x000000ff);

    return extension;
}

// Compute MD5 hash code, returns 4 integers array (static)
static unsigned int *ComputeMD5(unsigned char *data, int size)
{
#define LEFTROTATE(x, c) (((x) << (c)) | ((x) >> (32 - (c))))

    static unsigned int hash[4] = { 0 };

    // NOTE: All variables are unsigned 32 bit and wrap modulo 2^32 when calculating

    // r specifies the per-round shift amounts
    unsigned int r[] = {
        7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
        5,  9, 14, 20, 5,  9, 14, 20, 5,  9, 14, 20, 5,  9, 14, 20,
        4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
        6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21
    };

    // Use binary integer part of the sines of integers (in radians) as constants// Initialize variables:
    unsigned int k[] = {
        0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee,
        0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501,
        0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be,
        0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821,
        0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa,
        0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
        0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed,
        0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a,
        0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c,
        0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70,
        0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x04881d05,
        0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
        0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039,
        0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1,
        0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1,
        0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391
    };

    hash[0] = 0x67452301;
    hash[1] = 0xefcdab89;
    hash[2] = 0x98badcfe;
    hash[3] = 0x10325476;

    // Pre-processing: adding a single 1 bit
    // Append '1' bit to message
    // NOTE: The input bytes are considered as bits strings,
    // where the first bit is the most significant bit of the byte

    // Pre-processing: padding with zeros
    // Append '0' bit until message length in bit 448 (mod 512)
    // Append length mod (2 pow 64) to message

    int newDataSize = ((((size + 8)/64) + 1)*64) - 8;

    unsigned char *msg = RL_CALLOC(newDataSize + 64, 1);   // Also appends "0" bits (we alloc also 64 extra bytes...)
    memcpy(msg, data, size);
    msg[size] = 128;                 // Write the "1" bit

    unsigned int bitsLen = 8*size;
    memcpy(msg + newDataSize, &bitsLen, 4);  // We append the len in bits at the end of the buffer

    // Process the message in successive 512-bit chunks for each 512-bit chunk of message
    for (int offset = 0; offset < newDataSize; offset += (512/8))
    {
        // Break chunk into sixteen 32-bit words w[j], 0 <= j <= 15
        unsigned int *w = (unsigned int *)(msg + offset);

        // Initialize hash value for this chunk
        unsigned int a = hash[0];
        unsigned int b = hash[1];
        unsigned int c = hash[2];
        unsigned int d = hash[3];

        for (int i = 0; i < 64; i++)
        {
            unsigned int f, g;

            if (i < 16)
            {
                f = (b & c) | ((~b) & d);
                g = i;
            }
            else if (i < 32)
            {
                f = (d & b) | ((~d) & c);
                g = (5*i + 1)%16;
            }
            else if (i < 48)
            {
                f = b ^ c ^ d;
                g = (3*i + 5)%16;
            }
            else
            {
                f = c ^ (b | (~d));
                g = (7*i)%16;
            }

            unsigned int temp = d;
            d = c;
            c = b;
            b = b + LEFTROTATE((a + f + k[i] + w[g]), r[i]);
            a = temp;
        }

        // Add chunk's hash to result so far
        hash[0] += a;
        hash[1] += b;
        hash[2] += c;
        hash[3] += d;
    }

    RL_FREE(msg);

    return hash;
}

#endif // RRES_RAYLIB_IMPLEMENTATION
