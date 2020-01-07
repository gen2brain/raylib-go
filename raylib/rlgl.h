/**********************************************************************************************
*
*   rlgl - raylib OpenGL abstraction layer
*
*   rlgl is a wrapper for multiple OpenGL versions (1.1, 2.1, 3.3 Core, ES 2.0) to
*   pseudo-OpenGL 1.1 style functions (rlVertex, rlTranslate, rlRotate...).
*
*   When chosing an OpenGL version greater than OpenGL 1.1, rlgl stores vertex data on internal
*   VBO buffers (and VAOs if available). It requires calling 3 functions:
*       rlglInit()  - Initialize internal buffers and auxiliar resources
*       rlglDraw()  - Process internal buffers and send required draw calls
*       rlglClose() - De-initialize internal buffers data and other auxiliar resources
*
*   CONFIGURATION:
*
*   #define GRAPHICS_API_OPENGL_11
*   #define GRAPHICS_API_OPENGL_21
*   #define GRAPHICS_API_OPENGL_33
*   #define GRAPHICS_API_OPENGL_ES2
*       Use selected OpenGL graphics backend, should be supported by platform
*       Those preprocessor defines are only used on rlgl module, if OpenGL version is
*       required by any other module, use rlGetVersion() tocheck it
*
*   #define RLGL_IMPLEMENTATION
*       Generates the implementation of the library into the included file.
*       If not defined, the library is in header only mode and can be included in other headers
*       or source files without problems. But only ONE file should hold the implementation.
*
*   #define RLGL_STANDALONE
*       Use rlgl as standalone library (no raylib dependency)
*
*   #define SUPPORT_VR_SIMULATOR
*       Support VR simulation functionality (stereo rendering)
*
*   DEPENDENCIES:
*       raymath     - 3D math functionality (Vector3, Matrix, Quaternion)
*       GLAD        - OpenGL extensions loading (OpenGL 3.3 Core only)
*
*
*   LICENSE: zlib/libpng
*
*   Copyright (c) 2014-2019 Ramon Santamaria (@raysan5)
*
*   This software is provided "as-is", without any express or implied warranty. In no event
*   will the authors be held liable for any damages arising from the use of this software.
*
*   Permission is granted to anyone to use this software for any purpose, including commercial
*   applications, and to alter it and redistribute it freely, subject to the following restrictions:
*
*     1. The origin of this software must not be misrepresented; you must not claim that you
*     wrote the original software. If you use this software in a product, an acknowledgment
*     in the product documentation would be appreciated but is not required.
*
*     2. Altered source versions must be plainly marked as such, and must not be misrepresented
*     as being the original software.
*
*     3. This notice may not be removed or altered from any source distribution.
*
**********************************************************************************************/

#ifndef RLGL_H
#define RLGL_H

#if defined(RLGL_STANDALONE)
    #define RAYMATH_STANDALONE
    #define RAYMATH_HEADER_ONLY

    #if defined(_WIN32) && defined(BUILD_LIBTYPE_SHARED)
        #define RLAPI __declspec(dllexport)         // We are building raylib as a Win32 shared library (.dll)
    #elif defined(_WIN32) && defined(USE_LIBTYPE_SHARED)
        #define RLAPI __declspec(dllimport)         // We are using raylib as a Win32 shared library (.dll)
    #else
        #define RLAPI   // We are building or using raylib as a static library (or Linux shared library)
    #endif

    // Allow custom memory allocators
    #ifndef RL_MALLOC
        #define RL_MALLOC(sz)       malloc(sz)
    #endif
    #ifndef RL_CALLOC
        #define RL_CALLOC(n,sz)     calloc(n,sz)
    #endif
    #ifndef RL_FREE
        #define RL_FREE(p)          free(p)
    #endif
#else
    #include "raylib.h"         // Required for: Model, Shader, Texture2D, TraceLog()
#endif

#include "raymath.h"            // Required for: Vector3, Matrix

// Security check in case no GRAPHICS_API_OPENGL_* defined
#if !defined(GRAPHICS_API_OPENGL_11) && \
    !defined(GRAPHICS_API_OPENGL_21) && \
    !defined(GRAPHICS_API_OPENGL_33) && \
    !defined(GRAPHICS_API_OPENGL_ES2)
        #define GRAPHICS_API_OPENGL_33
#endif

// Security check in case multiple GRAPHICS_API_OPENGL_* defined
#if defined(GRAPHICS_API_OPENGL_11)
    #if defined(GRAPHICS_API_OPENGL_21)
        #undef GRAPHICS_API_OPENGL_21
    #endif
    #if defined(GRAPHICS_API_OPENGL_33)
        #undef GRAPHICS_API_OPENGL_33
    #endif
    #if defined(GRAPHICS_API_OPENGL_ES2)
        #undef GRAPHICS_API_OPENGL_ES2
    #endif
#endif

#if defined(GRAPHICS_API_OPENGL_21)
    #define GRAPHICS_API_OPENGL_33
#endif

//----------------------------------------------------------------------------------
// Defines and Macros
//----------------------------------------------------------------------------------
#if defined(GRAPHICS_API_OPENGL_11) || defined(GRAPHICS_API_OPENGL_33)
    // This is the maximum amount of elements (quads) per batch
    // NOTE: Be careful with text, every letter maps to a quad
    #define MAX_BATCH_ELEMENTS            8192
#elif defined(GRAPHICS_API_OPENGL_ES2)
    // We reduce memory sizes for embedded systems (RPI and HTML5)
    // NOTE: On HTML5 (emscripten) this is allocated on heap, by default it's only 16MB!...just take care...
    #define MAX_BATCH_ELEMENTS            2048
#endif

#define MAX_BATCH_BUFFERING                  1      // Max number of buffers for batching (multi-buffering)
#define MAX_MATRIX_STACK_SIZE               32      // Max size of Matrix stack
#define MAX_DRAWCALL_REGISTERED            256      // Max draws by state changes (mode, texture)

// Shader and material limits
#define MAX_SHADER_LOCATIONS                32      // Maximum number of predefined locations stored in shader struct
#define MAX_MATERIAL_MAPS                   12      // Maximum number of texture maps stored in shader struct

// Texture parameters (equivalent to OpenGL defines)
#define RL_TEXTURE_WRAP_S               0x2802      // GL_TEXTURE_WRAP_S
#define RL_TEXTURE_WRAP_T               0x2803      // GL_TEXTURE_WRAP_T
#define RL_TEXTURE_MAG_FILTER           0x2800      // GL_TEXTURE_MAG_FILTER
#define RL_TEXTURE_MIN_FILTER           0x2801      // GL_TEXTURE_MIN_FILTER
#define RL_TEXTURE_ANISOTROPIC_FILTER   0x3000      // Anisotropic filter (custom identifier)

#define RL_FILTER_NEAREST               0x2600      // GL_NEAREST
#define RL_FILTER_LINEAR                0x2601      // GL_LINEAR
#define RL_FILTER_MIP_NEAREST           0x2700      // GL_NEAREST_MIPMAP_NEAREST
#define RL_FILTER_NEAREST_MIP_LINEAR    0x2702      // GL_NEAREST_MIPMAP_LINEAR
#define RL_FILTER_LINEAR_MIP_NEAREST    0x2701      // GL_LINEAR_MIPMAP_NEAREST
#define RL_FILTER_MIP_LINEAR            0x2703      // GL_LINEAR_MIPMAP_LINEAR

#define RL_WRAP_REPEAT                  0x2901      // GL_REPEAT
#define RL_WRAP_CLAMP                   0x812F      // GL_CLAMP_TO_EDGE
#define RL_WRAP_MIRROR_REPEAT           0x8370      // GL_MIRRORED_REPEAT
#define RL_WRAP_MIRROR_CLAMP            0x8742      // GL_MIRROR_CLAMP_EXT

// Matrix modes (equivalent to OpenGL)
#define RL_MODELVIEW                    0x1700      // GL_MODELVIEW
#define RL_PROJECTION                   0x1701      // GL_PROJECTION
#define RL_TEXTURE                      0x1702      // GL_TEXTURE

// Primitive assembly draw modes
#define RL_LINES                        0x0001      // GL_LINES
#define RL_TRIANGLES                    0x0004      // GL_TRIANGLES
#define RL_QUADS                        0x0007      // GL_QUADS

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------
typedef enum { OPENGL_11 = 1, OPENGL_21, OPENGL_33, OPENGL_ES_20 } GlVersion;

typedef unsigned char byte;

#if defined(RLGL_STANDALONE)
    #ifndef __cplusplus
    // Boolean type
    typedef enum { false, true } bool;
    #endif

    // Color type, RGBA (32bit)
    typedef struct Color {
        unsigned char r;
        unsigned char g;
        unsigned char b;
        unsigned char a;
    } Color;

    // Texture2D type
    // NOTE: Data stored in GPU memory
    typedef struct Texture2D {
        unsigned int id;        // OpenGL texture id
        int width;              // Texture base width
        int height;             // Texture base height
        int mipmaps;            // Mipmap levels, 1 by default
        int format;             // Data format (PixelFormat)
    } Texture2D;

    // Texture type, same as Texture2D
    typedef Texture2D Texture;

    // TextureCubemap type, actually, same as Texture2D
    typedef Texture2D TextureCubemap;

    // RenderTexture2D type, for texture rendering
    typedef struct RenderTexture2D {
        unsigned int id;        // OpenGL framebuffer (fbo) id
        Texture2D texture;      // Color buffer attachment texture
        Texture2D depth;        // Depth buffer attachment texture
        bool depthTexture;      // Track if depth attachment is a texture or renderbuffer
    } RenderTexture2D;

    // RenderTexture type, same as RenderTexture2D
    typedef RenderTexture2D RenderTexture;

    // Vertex data definning a mesh
    typedef struct Mesh {
        int vertexCount;        // number of vertices stored in arrays
        int triangleCount;      // number of triangles stored (indexed or not)
        float *vertices;        // vertex position (XYZ - 3 components per vertex) (shader-location = 0)
        float *texcoords;       // vertex texture coordinates (UV - 2 components per vertex) (shader-location = 1)
        float *texcoords2;      // vertex second texture coordinates (useful for lightmaps) (shader-location = 5)
        float *normals;         // vertex normals (XYZ - 3 components per vertex) (shader-location = 2)
        float *tangents;        // vertex tangents (XYZW - 4 components per vertex) (shader-location = 4)
        unsigned char *colors;  // vertex colors (RGBA - 4 components per vertex) (shader-location = 3)
        unsigned short *indices;// vertex indices (in case vertex data comes indexed)

        // Animation vertex data
        float *animVertices;    // Animated vertex positions (after bones transformations)
        float *animNormals;     // Animated normals (after bones transformations)
        int *boneIds;           // Vertex bone ids, up to 4 bones influence by vertex (skinning)
        float *boneWeights;     // Vertex bone weight, up to 4 bones influence by vertex (skinning)

        // OpenGL identifiers
        unsigned int vaoId;     // OpenGL Vertex Array Object id
        unsigned int *vboId;    // OpenGL Vertex Buffer Objects id (7 types of vertex data)
    } Mesh;

    // Shader and material limits
    #define MAX_SHADER_LOCATIONS    32
    #define MAX_MATERIAL_MAPS       12

    // Shader type (generic)
    typedef struct Shader {
        unsigned int id;        // Shader program id
        int *locs;              // Shader locations array (MAX_SHADER_LOCATIONS)
    } Shader;

    // Material texture map
    typedef struct MaterialMap {
        Texture2D texture;      // Material map texture
        Color color;            // Material map color
        float value;            // Material map value
    } MaterialMap;

    // Material type (generic)
    typedef struct Material {
        Shader shader;          // Material shader
        MaterialMap *maps;      // Material maps (MAX_MATERIAL_MAPS)
        float *params;          // Material generic parameters (if required)
    } Material;

    // Camera type, defines a camera position/orientation in 3d space
    typedef struct Camera {
        Vector3 position;       // Camera position
        Vector3 target;         // Camera target it looks-at
        Vector3 up;             // Camera up vector (rotation over its axis)
        float fovy;             // Camera field-of-view apperture in Y (degrees)
    } Camera;

    // Head-Mounted-Display device parameters
    typedef struct VrDeviceInfo {
        int hResolution;                // HMD horizontal resolution in pixels
        int vResolution;                // HMD vertical resolution in pixels
        float hScreenSize;              // HMD horizontal size in meters
        float vScreenSize;              // HMD vertical size in meters
        float vScreenCenter;            // HMD screen center in meters
        float eyeToScreenDistance;      // HMD distance between eye and display in meters
        float lensSeparationDistance;   // HMD lens separation distance in meters
        float interpupillaryDistance;   // HMD IPD (distance between pupils) in meters
        float lensDistortionValues[4];  // HMD lens distortion constant parameters
        float chromaAbCorrection[4];    // HMD chromatic aberration correction parameters
    } VrDeviceInfo;

    // VR Stereo rendering configuration for simulator
    typedef struct VrStereoConfig {
        Shader distortionShader;        // VR stereo rendering distortion shader
        Matrix eyesProjection[2];       // VR stereo rendering eyes projection matrices
        Matrix eyesViewOffset[2];       // VR stereo rendering eyes view offset matrices
        int eyeViewportRight[4];        // VR stereo rendering right eye viewport [x, y, w, h]
        int eyeViewportLeft[4];         // VR stereo rendering left eye viewport [x, y, w, h]
    } VrStereoConfig;


    // TraceLog message types
    typedef enum {
        LOG_ALL,
        LOG_TRACE,
        LOG_DEBUG,
        LOG_INFO,
        LOG_WARNING,
        LOG_ERROR,
        LOG_FATAL,
        LOG_NONE
    } TraceLogType;

    // Texture formats (support depends on OpenGL version)
    typedef enum {
        UNCOMPRESSED_GRAYSCALE = 1,     // 8 bit per pixel (no alpha)
        UNCOMPRESSED_GRAY_ALPHA,
        UNCOMPRESSED_R5G6B5,            // 16 bpp
        UNCOMPRESSED_R8G8B8,            // 24 bpp
        UNCOMPRESSED_R5G5B5A1,          // 16 bpp (1 bit alpha)
        UNCOMPRESSED_R4G4B4A4,          // 16 bpp (4 bit alpha)
        UNCOMPRESSED_R8G8B8A8,          // 32 bpp
        UNCOMPRESSED_R32,               // 32 bpp (1 channel - float)
        UNCOMPRESSED_R32G32B32,         // 32*3 bpp (3 channels - float)
        UNCOMPRESSED_R32G32B32A32,      // 32*4 bpp (4 channels - float)
        COMPRESSED_DXT1_RGB,            // 4 bpp (no alpha)
        COMPRESSED_DXT1_RGBA,           // 4 bpp (1 bit alpha)
        COMPRESSED_DXT3_RGBA,           // 8 bpp
        COMPRESSED_DXT5_RGBA,           // 8 bpp
        COMPRESSED_ETC1_RGB,            // 4 bpp
        COMPRESSED_ETC2_RGB,            // 4 bpp
        COMPRESSED_ETC2_EAC_RGBA,       // 8 bpp
        COMPRESSED_PVRT_RGB,            // 4 bpp
        COMPRESSED_PVRT_RGBA,           // 4 bpp
        COMPRESSED_ASTC_4x4_RGBA,       // 8 bpp
        COMPRESSED_ASTC_8x8_RGBA        // 2 bpp
    } PixelFormat;

    // Texture parameters: filter mode
    // NOTE 1: Filtering considers mipmaps if available in the texture
    // NOTE 2: Filter is accordingly set for minification and magnification
    typedef enum {
        FILTER_POINT = 0,               // No filter, just pixel aproximation
        FILTER_BILINEAR,                // Linear filtering
        FILTER_TRILINEAR,               // Trilinear filtering (linear with mipmaps)
        FILTER_ANISOTROPIC_4X,          // Anisotropic filtering 4x
        FILTER_ANISOTROPIC_8X,          // Anisotropic filtering 8x
        FILTER_ANISOTROPIC_16X,         // Anisotropic filtering 16x
    } TextureFilterMode;

    // Color blending modes (pre-defined)
    typedef enum {
        BLEND_ALPHA = 0,
        BLEND_ADDITIVE,
        BLEND_MULTIPLIED
    } BlendMode;

    // Shader location point type
    typedef enum {
        LOC_VERTEX_POSITION = 0,
        LOC_VERTEX_TEXCOORD01,
        LOC_VERTEX_TEXCOORD02,
        LOC_VERTEX_NORMAL,
        LOC_VERTEX_TANGENT,
        LOC_VERTEX_COLOR,
        LOC_MATRIX_MVP,
        LOC_MATRIX_MODEL,
        LOC_MATRIX_VIEW,
        LOC_MATRIX_PROJECTION,
        LOC_VECTOR_VIEW,
        LOC_COLOR_DIFFUSE,
        LOC_COLOR_SPECULAR,
        LOC_COLOR_AMBIENT,
        LOC_MAP_ALBEDO,          // LOC_MAP_DIFFUSE
        LOC_MAP_METALNESS,       // LOC_MAP_SPECULAR
        LOC_MAP_NORMAL,
        LOC_MAP_ROUGHNESS,
        LOC_MAP_OCCLUSION,
        LOC_MAP_EMISSION,
        LOC_MAP_HEIGHT,
        LOC_MAP_CUBEMAP,
        LOC_MAP_IRRADIANCE,
        LOC_MAP_PREFILTER,
        LOC_MAP_BRDF
    } ShaderLocationIndex;

    // Shader uniform data types
    typedef enum {
        UNIFORM_FLOAT = 0,
        UNIFORM_VEC2,
        UNIFORM_VEC3,
        UNIFORM_VEC4,
        UNIFORM_INT,
        UNIFORM_IVEC2,
        UNIFORM_IVEC3,
        UNIFORM_IVEC4,
        UNIFORM_SAMPLER2D
    } ShaderUniformDataType;

    #define LOC_MAP_DIFFUSE      LOC_MAP_ALBEDO
    #define LOC_MAP_SPECULAR     LOC_MAP_METALNESS

    // Material map type
    typedef enum {
        MAP_ALBEDO    = 0,       // MAP_DIFFUSE
        MAP_METALNESS = 1,       // MAP_SPECULAR
        MAP_NORMAL    = 2,
        MAP_ROUGHNESS = 3,
        MAP_OCCLUSION,
        MAP_EMISSION,
        MAP_HEIGHT,
        MAP_CUBEMAP,             // NOTE: Uses GL_TEXTURE_CUBE_MAP
        MAP_IRRADIANCE,          // NOTE: Uses GL_TEXTURE_CUBE_MAP
        MAP_PREFILTER,           // NOTE: Uses GL_TEXTURE_CUBE_MAP
        MAP_BRDF
    } MaterialMapType;

    #define MAP_DIFFUSE      MAP_ALBEDO
    #define MAP_SPECULAR     MAP_METALNESS
#endif

#if defined(__cplusplus)
extern "C" {            // Prevents name mangling of functions
#endif

//------------------------------------------------------------------------------------
// Functions Declaration - Matrix operations
//------------------------------------------------------------------------------------
RLAPI void rlMatrixMode(int mode);                    // Choose the current matrix to be transformed
RLAPI void rlPushMatrix(void);                        // Push the current matrix to stack
RLAPI void rlPopMatrix(void);                         // Pop lattest inserted matrix from stack
RLAPI void rlLoadIdentity(void);                      // Reset current matrix to identity matrix
RLAPI void rlTranslatef(float x, float y, float z);   // Multiply the current matrix by a translation matrix
RLAPI void rlRotatef(float angleDeg, float x, float y, float z);  // Multiply the current matrix by a rotation matrix
RLAPI void rlScalef(float x, float y, float z);       // Multiply the current matrix by a scaling matrix
RLAPI void rlMultMatrixf(float *matf);                // Multiply the current matrix by another matrix
RLAPI void rlFrustum(double left, double right, double bottom, double top, double znear, double zfar);
RLAPI void rlOrtho(double left, double right, double bottom, double top, double znear, double zfar);
RLAPI void rlViewport(int x, int y, int width, int height); // Set the viewport area

//------------------------------------------------------------------------------------
// Functions Declaration - Vertex level operations
//------------------------------------------------------------------------------------
RLAPI void rlBegin(int mode);                         // Initialize drawing mode (how to organize vertex)
RLAPI void rlEnd(void);                               // Finish vertex providing
RLAPI void rlVertex2i(int x, int y);                  // Define one vertex (position) - 2 int
RLAPI void rlVertex2f(float x, float y);              // Define one vertex (position) - 2 float
RLAPI void rlVertex3f(float x, float y, float z);     // Define one vertex (position) - 3 float
RLAPI void rlTexCoord2f(float x, float y);            // Define one vertex (texture coordinate) - 2 float
RLAPI void rlNormal3f(float x, float y, float z);     // Define one vertex (normal) - 3 float
RLAPI void rlColor4ub(byte r, byte g, byte b, byte a);    // Define one vertex (color) - 4 byte
RLAPI void rlColor3f(float x, float y, float z);          // Define one vertex (color) - 3 float
RLAPI void rlColor4f(float x, float y, float z, float w); // Define one vertex (color) - 4 float

//------------------------------------------------------------------------------------
// Functions Declaration - OpenGL equivalent functions (common to 1.1, 3.3+, ES2)
// NOTE: This functions are used to completely abstract raylib code from OpenGL layer
//------------------------------------------------------------------------------------
void rlEnableTexture(unsigned int id);                  // Enable texture usage
void rlDisableTexture(void);                            // Disable texture usage
void rlTextureParameters(unsigned int id, int param, int value); // Set texture parameters (filter, wrap)
void rlEnableRenderTexture(unsigned int id);            // Enable render texture (fbo)
void rlDisableRenderTexture(void);                      // Disable render texture (fbo), return to default framebuffer
void rlEnableDepthTest(void);                           // Enable depth test
void rlDisableDepthTest(void);                          // Disable depth test
void rlEnableScissorTest(void);                         // Enable scissor test
void rlDisableScissorTest(void);                        // Disable scissor test
void rlScissor(int x, int y, int width, int height);    // Scissor test
void rlEnableWireMode(void);                            // Enable wire mode
void rlDisableWireMode(void);                           // Disable wire mode
void rlDeleteTextures(unsigned int id);                 // Delete OpenGL texture from GPU
void rlDeleteRenderTextures(RenderTexture2D target);    // Delete render textures (fbo) from GPU
void rlDeleteShader(unsigned int id);                   // Delete OpenGL shader program from GPU
void rlDeleteVertexArrays(unsigned int id);             // Unload vertex data (VAO) from GPU memory
void rlDeleteBuffers(unsigned int id);                  // Unload vertex data (VBO) from GPU memory
void rlClearColor(byte r, byte g, byte b, byte a);      // Clear color buffer with color
void rlClearScreenBuffers(void);                        // Clear used screen buffers (color and depth)

//------------------------------------------------------------------------------------
// Functions Declaration - rlgl functionality
//------------------------------------------------------------------------------------
RLAPI void rlglInit(int width, int height);           // Initialize rlgl (buffers, shaders, textures, states)
RLAPI void rlglClose(void);                           // De-inititialize rlgl (buffers, shaders, textures)
RLAPI void rlglDraw(void);                            // Update and draw default internal buffers

RLAPI int rlGetVersion(void);                         // Returns current OpenGL version
RLAPI bool rlCheckBufferLimit(int vCount);            // Check internal buffer overflow for a given number of vertex
RLAPI void rlSetDebugMarker(const char *text);        // Set debug marker for analysis
RLAPI void rlLoadExtensions(void *loader);            // Load OpenGL extensions
RLAPI Vector3 rlUnproject(Vector3 source, Matrix proj, Matrix view);  // Get world coordinates from screen coordinates

// Textures data management
RLAPI unsigned int rlLoadTexture(void *data, int width, int height, int format, int mipmapCount); // Load texture in GPU
RLAPI unsigned int rlLoadTextureDepth(int width, int height, int bits, bool useRenderBuffer);     // Load depth texture/renderbuffer (to be attached to fbo)
RLAPI unsigned int rlLoadTextureCubemap(void *data, int size, int format);                        // Load texture cubemap
RLAPI void rlUpdateTexture(unsigned int id, int width, int height, int format, const void *data); // Update GPU texture with new data
RLAPI void rlGetGlTextureFormats(int format, unsigned int *glInternalFormat, unsigned int *glFormat, unsigned int *glType);  // Get OpenGL internal formats
RLAPI void rlUnloadTexture(unsigned int id);                              // Unload texture from GPU memory

RLAPI void rlGenerateMipmaps(Texture2D *texture);                         // Generate mipmap data for selected texture
RLAPI void *rlReadTexturePixels(Texture2D texture);                       // Read texture pixel data
RLAPI unsigned char *rlReadScreenPixels(int width, int height);           // Read screen pixel data (color buffer)

// Render texture management (fbo)
RLAPI RenderTexture2D rlLoadRenderTexture(int width, int height, int format, int depthBits, bool useDepthTexture);    // Load a render texture (with color and depth attachments)
RLAPI void rlRenderTextureAttach(RenderTexture target, unsigned int id, int attachType);  // Attach texture/renderbuffer to an fbo
RLAPI bool rlRenderTextureComplete(RenderTexture target);                 // Verify render texture is complete

// Vertex data management
RLAPI void rlLoadMesh(Mesh *mesh, bool dynamic);                          // Upload vertex data into GPU and provided VAO/VBO ids
RLAPI void rlUpdateMesh(Mesh mesh, int buffer, int num);                  // Update vertex or index data on GPU (upload new data to one buffer)
RLAPI void rlUpdateMeshAt(Mesh mesh, int buffer, int num, int index);     // Update vertex or index data on GPU, at index
RLAPI void rlDrawMesh(Mesh mesh, Material material, Matrix transform);    // Draw a 3d mesh with material and transform
RLAPI void rlUnloadMesh(Mesh mesh);                                       // Unload mesh data from CPU and GPU

// NOTE: There is a set of shader related functions that are available to end user,
// to avoid creating function wrappers through core module, they have been directly declared in raylib.h

#if defined(RLGL_STANDALONE)
//------------------------------------------------------------------------------------
// Shaders System Functions (Module: rlgl)
// NOTE: This functions are useless when using OpenGL 1.1
//------------------------------------------------------------------------------------
// Shader loading/unloading functions
RLAPI char *LoadText(const char *fileName);                               // Load chars array from text file
RLAPI Shader LoadShader(const char *vsFileName, const char *fsFileName);  // Load shader from files and bind default locations
RLAPI Shader LoadShaderCode(const char *vsCode, const char *fsCode);                  // Load shader from code strings and bind default locations
RLAPI void UnloadShader(Shader shader);                                   // Unload shader from GPU memory (VRAM)

RLAPI Shader GetShaderDefault(void);                                      // Get default shader
RLAPI Texture2D GetTextureDefault(void);                                  // Get default texture

// Shader configuration functions
RLAPI int GetShaderLocation(Shader shader, const char *uniformName);              // Get shader uniform location
RLAPI void SetShaderValue(Shader shader, int uniformLoc, const void *value, int uniformType);               // Set shader uniform value
RLAPI void SetShaderValueV(Shader shader, int uniformLoc, const void *value, int uniformType, int count);   // Set shader uniform value vector
RLAPI void SetShaderValueMatrix(Shader shader, int uniformLoc, Matrix mat);       // Set shader uniform value (matrix 4x4)
RLAPI void SetMatrixProjection(Matrix proj);                              // Set a custom projection matrix (replaces internal projection matrix)
RLAPI void SetMatrixModelview(Matrix view);                               // Set a custom modelview matrix (replaces internal modelview matrix)
RLAPI Matrix GetMatrixModelview(void);                                    // Get internal modelview matrix

// Texture maps generation (PBR)
// NOTE: Required shaders should be provided
RLAPI Texture2D GenTextureCubemap(Shader shader, Texture2D skyHDR, int size);       // Generate cubemap texture from HDR texture
RLAPI Texture2D GenTextureIrradiance(Shader shader, Texture2D cubemap, int size);   // Generate irradiance texture using cubemap data
RLAPI Texture2D GenTexturePrefilter(Shader shader, Texture2D cubemap, int size);    // Generate prefilter texture using cubemap data
RLAPI Texture2D GenTextureBRDF(Shader shader, int size);                  // Generate BRDF texture using cubemap data

// Shading begin/end functions
RLAPI void BeginShaderMode(Shader shader);              // Begin custom shader drawing
RLAPI void EndShaderMode(void);                         // End custom shader drawing (use default shader)
RLAPI void BeginBlendMode(int mode);                    // Begin blending mode (alpha, additive, multiplied)
RLAPI void EndBlendMode(void);                          // End blending mode (reset to default: alpha blending)

// VR control functions
RLAPI void InitVrSimulator(void);                       // Init VR simulator for selected device parameters
RLAPI void CloseVrSimulator(void);                      // Close VR simulator for current device
RLAPI void UpdateVrTracking(Camera *camera);            // Update VR tracking (position and orientation) and camera
RLAPI void SetVrConfiguration(VrDeviceInfo info, Shader distortion);      // Set stereo rendering configuration parameters
RLAPI bool IsVrSimulatorReady(void);                    // Detect if VR simulator is ready
RLAPI void ToggleVrMode(void);                          // Enable/Disable VR experience
RLAPI void BeginVrDrawing(void);                        // Begin VR simulator stereo rendering
RLAPI void EndVrDrawing(void);                          // End VR simulator stereo rendering

RLAPI void TraceLog(int msgType, const char *text, ...);      // Show trace log messages (LOG_INFO, LOG_WARNING, LOG_ERROR, LOG_DEBUG)
RLAPI int GetPixelDataSize(int width, int height, int format);// Get pixel data size in bytes (image or texture)
#endif

#if defined(__cplusplus)
}
#endif

#endif // RLGL_H

/***********************************************************************************
*
*   RLGL IMPLEMENTATION
*
************************************************************************************/

#if defined(RLGL_IMPLEMENTATION)

#if !defined(RLGL_STANDALONE)
    // Check if config flags have been externally provided on compilation line
    #if !defined(EXTERNAL_CONFIG_FLAGS)
        #include "config.h"         // Defines module configuration flags
    #endif
#endif

#include <stdio.h>                  // Required for: fopen(), fclose(), fread()... [Used only on LoadText()]
#include <stdlib.h>                 // Required for: malloc(), free(), rand()
#include <string.h>                 // Required for: strcmp(), strlen(), strtok() [Used only in extensions loading]
#include <math.h>                   // Required for: atan2()

#if !defined(RLGL_STANDALONE)
    #include "raymath.h"            // Required for: Vector3 and Matrix functions
#endif

#if defined(GRAPHICS_API_OPENGL_11)
    #if defined(__APPLE__)
        #include <OpenGL/gl.h>      // OpenGL 1.1 library for OSX
        #include <OpenGL/glext.h>
    #else
        // APIENTRY for OpenGL function pointer declarations is required
        #ifndef APIENTRY
            #if defined(_WIN32)
                #define APIENTRY __stdcall
            #else
                #define APIENTRY
            #endif
        #endif
        // WINGDIAPI definition. Some Windows OpenGL headers need it
        #if !defined(WINGDIAPI) && defined(_WIN32)
            #define WINGDIAPI __declspec(dllimport)
        #endif

        #include <GL/gl.h>              // OpenGL 1.1 library
    #endif
#endif

#if defined(GRAPHICS_API_OPENGL_21)
    #define GRAPHICS_API_OPENGL_33      // OpenGL 2.1 uses mostly OpenGL 3.3 Core functionality
#endif

#if defined(GRAPHICS_API_OPENGL_33)
    #if defined(__APPLE__)
        #include <OpenGL/gl3.h>         // OpenGL 3 library for OSX
        #include <OpenGL/gl3ext.h>      // OpenGL 3 extensions library for OSX
    #else
        #define GLAD_IMPLEMENTATION
        #if defined(RLGL_STANDALONE)
            #include "glad.h"           // GLAD extensions loading library, includes OpenGL headers
        #else
            #include "external/glad.h"  // GLAD extensions loading library, includes OpenGL headers
        #endif
    #endif
#endif

#if defined(GRAPHICS_API_OPENGL_ES2)
    #include <EGL/egl.h>                // EGL library
    #include <GLES2/gl2.h>              // OpenGL ES 2.0 library
    #include <GLES2/gl2ext.h>           // OpenGL ES 2.0 extensions library
#endif

#if defined(RLGL_STANDALONE)
    #include <stdarg.h>                 // Required for: va_list, va_start(), vfprintf(), va_end() [Used only on TraceLog()]
#endif

//----------------------------------------------------------------------------------
// Defines and Macros
//----------------------------------------------------------------------------------
#ifndef GL_SHADING_LANGUAGE_VERSION
    #define GL_SHADING_LANGUAGE_VERSION         0x8B8C
#endif

#ifndef GL_COMPRESSED_RGB_S3TC_DXT1_EXT
    #define GL_COMPRESSED_RGB_S3TC_DXT1_EXT     0x83F0
#endif
#ifndef GL_COMPRESSED_RGBA_S3TC_DXT1_EXT
    #define GL_COMPRESSED_RGBA_S3TC_DXT1_EXT    0x83F1
#endif
#ifndef GL_COMPRESSED_RGBA_S3TC_DXT3_EXT
    #define GL_COMPRESSED_RGBA_S3TC_DXT3_EXT    0x83F2
#endif
#ifndef GL_COMPRESSED_RGBA_S3TC_DXT5_EXT
    #define GL_COMPRESSED_RGBA_S3TC_DXT5_EXT    0x83F3
#endif
#ifndef GL_ETC1_RGB8_OES
    #define GL_ETC1_RGB8_OES                    0x8D64
#endif
#ifndef GL_COMPRESSED_RGB8_ETC2
    #define GL_COMPRESSED_RGB8_ETC2             0x9274
#endif
#ifndef GL_COMPRESSED_RGBA8_ETC2_EAC
    #define GL_COMPRESSED_RGBA8_ETC2_EAC        0x9278
#endif
#ifndef GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG
    #define GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG  0x8C00
#endif
#ifndef GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG
    #define GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG 0x8C02
#endif
#ifndef GL_COMPRESSED_RGBA_ASTC_4x4_KHR
    #define GL_COMPRESSED_RGBA_ASTC_4x4_KHR     0x93b0
#endif
#ifndef GL_COMPRESSED_RGBA_ASTC_8x8_KHR
    #define GL_COMPRESSED_RGBA_ASTC_8x8_KHR     0x93b7
#endif

#ifndef GL_MAX_TEXTURE_MAX_ANISOTROPY_EXT
    #define GL_MAX_TEXTURE_MAX_ANISOTROPY_EXT   0x84FF
#endif

#ifndef GL_TEXTURE_MAX_ANISOTROPY_EXT
    #define GL_TEXTURE_MAX_ANISOTROPY_EXT       0x84FE
#endif

#if defined(GRAPHICS_API_OPENGL_11)
    #define GL_UNSIGNED_SHORT_5_6_5             0x8363
    #define GL_UNSIGNED_SHORT_5_5_5_1           0x8034
    #define GL_UNSIGNED_SHORT_4_4_4_4           0x8033
#endif

#if defined(GRAPHICS_API_OPENGL_21)
    #define GL_LUMINANCE                        0x1909
    #define GL_LUMINANCE_ALPHA                  0x190A
#endif

#if defined(GRAPHICS_API_OPENGL_ES2)
    #define glClearDepth                glClearDepthf
    #define GL_READ_FRAMEBUFFER         GL_FRAMEBUFFER
    #define GL_DRAW_FRAMEBUFFER         GL_FRAMEBUFFER
#endif

// Default vertex attribute names on shader to set location points
#define DEFAULT_ATTRIB_POSITION_NAME    "vertexPosition"    // shader-location = 0
#define DEFAULT_ATTRIB_TEXCOORD_NAME    "vertexTexCoord"    // shader-location = 1
#define DEFAULT_ATTRIB_NORMAL_NAME      "vertexNormal"      // shader-location = 2
#define DEFAULT_ATTRIB_COLOR_NAME       "vertexColor"       // shader-location = 3
#define DEFAULT_ATTRIB_TANGENT_NAME     "vertexTangent"     // shader-location = 4
#define DEFAULT_ATTRIB_TEXCOORD2_NAME   "vertexTexCoord2"   // shader-location = 5

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------

// Dynamic vertex buffers (position + texcoords + colors + indices arrays)
typedef struct DynamicBuffer {
    int vCounter;               // vertex position counter to process (and draw) from full buffer
    int tcCounter;              // vertex texcoord counter to process (and draw) from full buffer
    int cCounter;               // vertex color counter to process (and draw) from full buffer
    float *vertices;            // vertex position (XYZ - 3 components per vertex) (shader-location = 0)
    float *texcoords;           // vertex texture coordinates (UV - 2 components per vertex) (shader-location = 1)
    unsigned char *colors;      // vertex colors (RGBA - 4 components per vertex) (shader-location = 3)
#if defined(GRAPHICS_API_OPENGL_11) || defined(GRAPHICS_API_OPENGL_33)
    unsigned int *indices;      // vertex indices (in case vertex data comes indexed) (6 indices per quad)
#elif defined(GRAPHICS_API_OPENGL_ES2)
    unsigned short *indices;    // vertex indices (in case vertex data comes indexed) (6 indices per quad)
                                // NOTE: 6*2 byte = 12 byte, not alignment problem!
#endif
    unsigned int vaoId;         // OpenGL Vertex Array Object id
    unsigned int vboId[4];      // OpenGL Vertex Buffer Objects id (4 types of vertex data)
} DynamicBuffer;

// Draw call type
typedef struct DrawCall {
    int mode;                   // Drawing mode: LINES, TRIANGLES, QUADS
    int vertexCount;            // Number of vertex of the draw
    int vertexAlignment;        // Number of vertex required for index alignment (LINES, TRIANGLES)
    //unsigned int vaoId;         // Vertex array id to be used on the draw
    //unsigned int shaderId;      // Shader id to be used on the draw
    unsigned int textureId;     // Texture id to be used on the draw
                                // TODO: Support additional texture units?

    //Matrix projection;        // Projection matrix for this draw
    //Matrix modelview;         // Modelview matrix for this draw
} DrawCall;

#if defined(SUPPORT_VR_SIMULATOR)
// VR Stereo rendering configuration for simulator
typedef struct VrStereoConfig {
    Shader distortionShader;        // VR stereo rendering distortion shader
    Matrix eyesProjection[2];       // VR stereo rendering eyes projection matrices
    Matrix eyesViewOffset[2];       // VR stereo rendering eyes view offset matrices
    int eyeViewportRight[4];        // VR stereo rendering right eye viewport [x, y, w, h]
    int eyeViewportLeft[4];         // VR stereo rendering left eye viewport [x, y, w, h]
} VrStereoConfig;
#endif

//----------------------------------------------------------------------------------
// Global Variables Definition
//----------------------------------------------------------------------------------
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
static Matrix stack[MAX_MATRIX_STACK_SIZE] = { 0 };
static int stackCounter = 0;

static Matrix modelview = { 0 };
static Matrix projection = { 0 };
static Matrix *currentMatrix = NULL;
static int currentMatrixMode = -1;
static float currentDepth = -1.0f;

// Default dynamic buffer for elements data
// NOTE: A multi-buffering system is supported
static DynamicBuffer vertexData[MAX_BATCH_BUFFERING] = { 0 };
static int currentBuffer = 0;

// Transform matrix to be used with rlTranslate, rlRotate, rlScale
static Matrix transformMatrix = { 0 };
static bool useTransformMatrix = false;

// Default buffers draw calls
static DrawCall *draws = NULL;
static int drawsCounter = 0;

// Default texture (1px white) useful for plain color polys (required by shader)
static unsigned int defaultTextureId = 0;

// Default shaders
static unsigned int defaultVShaderId = 0;   // Default vertex shader id (used by default shader program)
static unsigned int defaultFShaderId = 0;   // Default fragment shader Id (used by default shader program)

static Shader defaultShader = { 0 };        // Basic shader, support vertex color and diffuse texture
static Shader currentShader = { 0 };        // Shader to be used on rendering (by default, defaultShader)

// Extension supported flag: VAO
static bool vaoSupported = false;           // VAO support (OpenGL ES2 could not support VAO extension)

// Extension supported flag: Compressed textures
static bool texCompDXTSupported = false;    // DDS texture compression support
static bool texCompETC1Supported = false;   // ETC1 texture compression support
static bool texCompETC2Supported = false;   // ETC2/EAC texture compression support
static bool texCompPVRTSupported = false;   // PVR texture compression support
static bool texCompASTCSupported = false;   // ASTC texture compression support

// Extension supported flag: Textures format
static bool texNPOTSupported = false;       // NPOT textures full support
static bool texFloatSupported = false;      // float textures support (32 bit per channel)
static bool texDepthSupported = false;      // Depth textures supported
static int maxDepthBits = 16;               // Maximum bits for depth component

// Extension supported flag: Clamp mirror wrap mode
static bool texMirrorClampSupported = false;        // Clamp mirror wrap mode supported

// Extension supported flag: Anisotropic filtering
static bool texAnisotropicFilterSupported = false;  // Anisotropic texture filtering support
static float maxAnisotropicLevel = 0.0f;            // Maximum anisotropy level supported (minimum is 2.0f)

static bool debugMarkerSupported = false;   // Debug marker support

#if defined(GRAPHICS_API_OPENGL_ES2)
// NOTE: VAO functionality is exposed through extensions (OES)
static PFNGLGENVERTEXARRAYSOESPROC glGenVertexArrays;
static PFNGLBINDVERTEXARRAYOESPROC glBindVertexArray;
static PFNGLDELETEVERTEXARRAYSOESPROC glDeleteVertexArrays;
//static PFNGLISVERTEXARRAYOESPROC glIsVertexArray;   // NOTE: Fails in WebGL, omitted
#endif

#if defined(SUPPORT_VR_SIMULATOR)
// VR global variables
static VrStereoConfig vrConfig = { 0 };     // VR stereo configuration for simulator
static RenderTexture2D stereoFbo = { 0 };   // VR stereo rendering framebuffer
static bool vrSimulatorReady = false;       // VR simulator ready flag
static bool vrStereoRender = false;         // VR stereo rendering enabled/disabled flag
                                            // NOTE: This flag is useful to render data over stereo image (i.e. FPS)
#endif  // SUPPORT_VR_SIMULATOR

#endif  // GRAPHICS_API_OPENGL_33 || GRAPHICS_API_OPENGL_ES2

static int blendMode = 0;                   // Track current blending mode

// Default framebuffer size
static int framebufferWidth = 0;            // Default framebuffer width
static int framebufferHeight = 0;           // Default framebuffer height

//----------------------------------------------------------------------------------
// Module specific Functions Declaration
//----------------------------------------------------------------------------------
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
static unsigned int CompileShader(const char *shaderStr, int type);     // Compile custom shader and return shader id
static unsigned int LoadShaderProgram(unsigned int vShaderId, unsigned int fShaderId);  // Load custom shader program

static Shader LoadShaderDefault(void);      // Load default shader (just vertex positioning and texture coloring)
static void SetShaderDefaultLocations(Shader *shader); // Bind default shader locations (attributes and uniforms)
static void UnloadShaderDefault(void);      // Unload default shader

static void LoadBuffersDefault(void);       // Load default internal buffers
static void UpdateBuffersDefault(void);     // Update default internal buffers (VAOs/VBOs) with vertex data
static void DrawBuffersDefault(void);       // Draw default internal buffers vertex data
static void UnloadBuffersDefault(void);     // Unload default internal buffers vertex data from CPU and GPU

static void GenDrawCube(void);              // Generate and draw cube
static void GenDrawQuad(void);              // Generate and draw quad

#if defined(SUPPORT_VR_SIMULATOR)
static void SetStereoView(int eye, Matrix matProjection, Matrix matModelView);  // Set internal projection and modelview matrix depending on eye
#endif

#endif  // GRAPHICS_API_OPENGL_33 || GRAPHICS_API_OPENGL_ES2

#if defined(GRAPHICS_API_OPENGL_11)
static int GenerateMipmaps(unsigned char *data, int baseWidth, int baseHeight);
static Color *GenNextMipmap(Color *srcData, int srcWidth, int srcHeight);
#endif

//----------------------------------------------------------------------------------
// Module Functions Definition - Matrix operations
//----------------------------------------------------------------------------------

#if defined(GRAPHICS_API_OPENGL_11)

// Fallback to OpenGL 1.1 function calls
//---------------------------------------
void rlMatrixMode(int mode)
{
    switch (mode)
    {
        case RL_PROJECTION: glMatrixMode(GL_PROJECTION); break;
        case RL_MODELVIEW: glMatrixMode(GL_MODELVIEW); break;
        case RL_TEXTURE: glMatrixMode(GL_TEXTURE); break;
        default: break;
    }
}

void rlFrustum(double left, double right, double bottom, double top, double znear, double zfar)
{
    glFrustum(left, right, bottom, top, znear, zfar);
}

void rlOrtho(double left, double right, double bottom, double top, double znear, double zfar)
{
    glOrtho(left, right, bottom, top, znear, zfar);
}

void rlPushMatrix(void) { glPushMatrix(); }
void rlPopMatrix(void) { glPopMatrix(); }
void rlLoadIdentity(void) { glLoadIdentity(); }
void rlTranslatef(float x, float y, float z) { glTranslatef(x, y, z); }
void rlRotatef(float angleDeg, float x, float y, float z) { glRotatef(angleDeg, x, y, z); }
void rlScalef(float x, float y, float z) { glScalef(x, y, z); }
void rlMultMatrixf(float *matf) { glMultMatrixf(matf); }

#elif defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)

// Choose the current matrix to be transformed
void rlMatrixMode(int mode)
{
    if (mode == RL_PROJECTION) currentMatrix = &projection;
    else if (mode == RL_MODELVIEW) currentMatrix = &modelview;
    //else if (mode == RL_TEXTURE) // Not supported

    currentMatrixMode = mode;
}

// Push the current matrix into stack
void rlPushMatrix(void)
{
    if (stackCounter >= MAX_MATRIX_STACK_SIZE) TraceLog(LOG_ERROR, "Matrix stack overflow");

    if (currentMatrixMode == RL_MODELVIEW)
    {
        useTransformMatrix = true;
        currentMatrix = &transformMatrix;
    }

    stack[stackCounter] = *currentMatrix;
    stackCounter++;
}

// Pop lattest inserted matrix from stack
void rlPopMatrix(void)
{
    if (stackCounter > 0)
    {
        Matrix mat = stack[stackCounter - 1];
        *currentMatrix = mat;
        stackCounter--;
    }

    if ((stackCounter == 0) && (currentMatrixMode == RL_MODELVIEW))
    {
        currentMatrix = &modelview;
        useTransformMatrix = false;
    }
}

// Reset current matrix to identity matrix
void rlLoadIdentity(void)
{
    *currentMatrix = MatrixIdentity();
}

// Multiply the current matrix by a translation matrix
void rlTranslatef(float x, float y, float z)
{
    Matrix matTranslation = MatrixTranslate(x, y, z);

    // NOTE: We transpose matrix with multiplication order
    *currentMatrix = MatrixMultiply(matTranslation, *currentMatrix);
}

// Multiply the current matrix by a rotation matrix
void rlRotatef(float angleDeg, float x, float y, float z)
{
    Matrix matRotation = MatrixIdentity();

    Vector3 axis = (Vector3){ x, y, z };
    matRotation = MatrixRotate(Vector3Normalize(axis), angleDeg*DEG2RAD);

    // NOTE: We transpose matrix with multiplication order
    *currentMatrix = MatrixMultiply(matRotation, *currentMatrix);
}

// Multiply the current matrix by a scaling matrix
void rlScalef(float x, float y, float z)
{
    Matrix matScale = MatrixScale(x, y, z);

    // NOTE: We transpose matrix with multiplication order
    *currentMatrix = MatrixMultiply(matScale, *currentMatrix);
}

// Multiply the current matrix by another matrix
void rlMultMatrixf(float *matf)
{
    // Matrix creation from array
    Matrix mat = { matf[0], matf[4], matf[8], matf[12],
                   matf[1], matf[5], matf[9], matf[13],
                   matf[2], matf[6], matf[10], matf[14],
                   matf[3], matf[7], matf[11], matf[15] };

    *currentMatrix = MatrixMultiply(*currentMatrix, mat);
}

// Multiply the current matrix by a perspective matrix generated by parameters
void rlFrustum(double left, double right, double bottom, double top, double znear, double zfar)
{
    Matrix matPerps = MatrixFrustum(left, right, bottom, top, znear, zfar);

    *currentMatrix = MatrixMultiply(*currentMatrix, matPerps);
}

// Multiply the current matrix by an orthographic matrix generated by parameters
void rlOrtho(double left, double right, double bottom, double top, double znear, double zfar)
{
    Matrix matOrtho = MatrixOrtho(left, right, bottom, top, znear, zfar);

    *currentMatrix = MatrixMultiply(*currentMatrix, matOrtho);
}

#endif

// Set the viewport area (transformation from normalized device coordinates to window coordinates)
// NOTE: Updates global variables: framebufferWidth, framebufferHeight
void rlViewport(int x, int y, int width, int height)
{
    glViewport(x, y, width, height);
}

//----------------------------------------------------------------------------------
// Module Functions Definition - Vertex level operations
//----------------------------------------------------------------------------------
#if defined(GRAPHICS_API_OPENGL_11)

// Fallback to OpenGL 1.1 function calls
//---------------------------------------
void rlBegin(int mode)
{
    switch (mode)
    {
        case RL_LINES: glBegin(GL_LINES); break;
        case RL_TRIANGLES: glBegin(GL_TRIANGLES); break;
        case RL_QUADS: glBegin(GL_QUADS); break;
        default: break;
    }
}

void rlEnd() { glEnd(); }
void rlVertex2i(int x, int y) { glVertex2i(x, y); }
void rlVertex2f(float x, float y) { glVertex2f(x, y); }
void rlVertex3f(float x, float y, float z) { glVertex3f(x, y, z); }
void rlTexCoord2f(float x, float y) { glTexCoord2f(x, y); }
void rlNormal3f(float x, float y, float z) { glNormal3f(x, y, z); }
void rlColor4ub(byte r, byte g, byte b, byte a) { glColor4ub(r, g, b, a); }
void rlColor3f(float x, float y, float z) { glColor3f(x, y, z); }
void rlColor4f(float x, float y, float z, float w) { glColor4f(x, y, z, w); }

#elif defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)

// Initialize drawing mode (how to organize vertex)
void rlBegin(int mode)
{
    // Draw mode can be RL_LINES, RL_TRIANGLES and RL_QUADS
    // NOTE: In all three cases, vertex are accumulated over default internal vertex buffer
    if (draws[drawsCounter - 1].mode != mode)
    {
        if (draws[drawsCounter - 1].vertexCount > 0)
        {
            // Make sure current draws[i].vertexCount is aligned a multiple of 4,
            // that way, following QUADS drawing will keep aligned with index processing
            // It implies adding some extra alignment vertex at the end of the draw,
            // those vertex are not processed but they are considered as an additional offset
            // for the next set of vertex to be drawn
            if (draws[drawsCounter - 1].mode == RL_LINES) draws[drawsCounter - 1].vertexAlignment = ((draws[drawsCounter - 1].vertexCount < 4)? draws[drawsCounter - 1].vertexCount : draws[drawsCounter - 1].vertexCount%4);
            else if (draws[drawsCounter - 1].mode == RL_TRIANGLES) draws[drawsCounter - 1].vertexAlignment = ((draws[drawsCounter - 1].vertexCount < 4)? 1 : (4 - (draws[drawsCounter - 1].vertexCount%4)));

            else draws[drawsCounter - 1].vertexAlignment = 0;

            if (rlCheckBufferLimit(draws[drawsCounter - 1].vertexAlignment)) rlglDraw();
            else
            {
                vertexData[currentBuffer].vCounter += draws[drawsCounter - 1].vertexAlignment;
                vertexData[currentBuffer].cCounter += draws[drawsCounter - 1].vertexAlignment;
                vertexData[currentBuffer].tcCounter += draws[drawsCounter - 1].vertexAlignment;

                drawsCounter++;
            }
        }

        if (drawsCounter >= MAX_DRAWCALL_REGISTERED) rlglDraw();

        draws[drawsCounter - 1].mode = mode;
        draws[drawsCounter - 1].vertexCount = 0;
        draws[drawsCounter - 1].textureId = defaultTextureId;
    }
}

// Finish vertex providing
void rlEnd(void)
{
    // Make sure vertexCount is the same for vertices, texcoords, colors and normals
    // NOTE: In OpenGL 1.1, one glColor call can be made for all the subsequent glVertex calls

    // Make sure colors count match vertex count
    if (vertexData[currentBuffer].vCounter != vertexData[currentBuffer].cCounter)
    {
        int addColors = vertexData[currentBuffer].vCounter - vertexData[currentBuffer].cCounter;

        for (int i = 0; i < addColors; i++)
        {
            vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter] = vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter - 4];
            vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter + 1] = vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter - 3];
            vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter + 2] = vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter - 2];
            vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter + 3] = vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter - 1];
            vertexData[currentBuffer].cCounter++;
        }
    }

    // Make sure texcoords count match vertex count
    if (vertexData[currentBuffer].vCounter != vertexData[currentBuffer].tcCounter)
    {
        int addTexCoords = vertexData[currentBuffer].vCounter - vertexData[currentBuffer].tcCounter;

        for (int i = 0; i < addTexCoords; i++)
        {
            vertexData[currentBuffer].texcoords[2*vertexData[currentBuffer].tcCounter] = 0.0f;
            vertexData[currentBuffer].texcoords[2*vertexData[currentBuffer].tcCounter + 1] = 0.0f;
            vertexData[currentBuffer].tcCounter++;
        }
    }

    // TODO: Make sure normals count match vertex count... if normals support is added in a future... :P

    // NOTE: Depth increment is dependant on rlOrtho(): z-near and z-far values,
    // as well as depth buffer bit-depth (16bit or 24bit or 32bit)
    // Correct increment formula would be: depthInc = (zfar - znear)/pow(2, bits)
    currentDepth += (1.0f/20000.0f);

    // Verify internal buffers limits
    // NOTE: This check is combined with usage of rlCheckBufferLimit()
    if ((vertexData[currentBuffer].vCounter) >= (MAX_BATCH_ELEMENTS*4 - 4))
    {
        // WARNING: If we are between rlPushMatrix() and rlPopMatrix() and we need to force a rlglDraw(),
        // we need to call rlPopMatrix() before to recover *currentMatrix (modelview) for the next forced draw call!
        // If we have multiple matrix pushed, it will require "stackCounter" pops before launching the draw
        for (int i = stackCounter; i >= 0; i--) rlPopMatrix();
        rlglDraw();
    }
}

// Define one vertex (position)
// NOTE: Vertex position data is the basic information required for drawing
void rlVertex3f(float x, float y, float z)
{
    Vector3 vec = { x, y, z };

    // Transform provided vector if required
    if (useTransformMatrix) vec = Vector3Transform(vec, transformMatrix);

    // Verify that MAX_BATCH_ELEMENTS limit not reached
    if (vertexData[currentBuffer].vCounter < (MAX_BATCH_ELEMENTS*4))
    {
        vertexData[currentBuffer].vertices[3*vertexData[currentBuffer].vCounter] = vec.x;
        vertexData[currentBuffer].vertices[3*vertexData[currentBuffer].vCounter + 1] = vec.y;
        vertexData[currentBuffer].vertices[3*vertexData[currentBuffer].vCounter + 2] = vec.z;
        vertexData[currentBuffer].vCounter++;

        draws[drawsCounter - 1].vertexCount++;
    }
    else TraceLog(LOG_ERROR, "MAX_BATCH_ELEMENTS overflow");
}

// Define one vertex (position)
void rlVertex2f(float x, float y)
{
    rlVertex3f(x, y, currentDepth);
}

// Define one vertex (position)
void rlVertex2i(int x, int y)
{
    rlVertex3f((float)x, (float)y, currentDepth);
}

// Define one vertex (texture coordinate)
// NOTE: Texture coordinates are limited to QUADS only
void rlTexCoord2f(float x, float y)
{
    vertexData[currentBuffer].texcoords[2*vertexData[currentBuffer].tcCounter] = x;
    vertexData[currentBuffer].texcoords[2*vertexData[currentBuffer].tcCounter + 1] = y;
    vertexData[currentBuffer].tcCounter++;
}

// Define one vertex (normal)
// NOTE: Normals limited to TRIANGLES only?
void rlNormal3f(float x, float y, float z)
{
    // TODO: Normals usage...
}

// Define one vertex (color)
void rlColor4ub(byte x, byte y, byte z, byte w)
{
    vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter] = x;
    vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter + 1] = y;
    vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter + 2] = z;
    vertexData[currentBuffer].colors[4*vertexData[currentBuffer].cCounter + 3] = w;
    vertexData[currentBuffer].cCounter++;
}

// Define one vertex (color)
void rlColor4f(float r, float g, float b, float a)
{
    rlColor4ub((byte)(r*255), (byte)(g*255), (byte)(b*255), (byte)(a*255));
}

// Define one vertex (color)
void rlColor3f(float x, float y, float z)
{
    rlColor4ub((byte)(x*255), (byte)(y*255), (byte)(z*255), 255);
}

#endif

//----------------------------------------------------------------------------------
// Module Functions Definition - OpenGL equivalent functions (common to 1.1, 3.3+, ES2)
//----------------------------------------------------------------------------------

// Enable texture usage
void rlEnableTexture(unsigned int id)
{
#if defined(GRAPHICS_API_OPENGL_11)
    glEnable(GL_TEXTURE_2D);
    glBindTexture(GL_TEXTURE_2D, id);
#endif

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if (draws[drawsCounter - 1].textureId != id)
    {
        if (draws[drawsCounter - 1].vertexCount > 0)
        {
            // Make sure current draws[i].vertexCount is aligned a multiple of 4,
            // that way, following QUADS drawing will keep aligned with index processing
            // It implies adding some extra alignment vertex at the end of the draw,
            // those vertex are not processed but they are considered as an additional offset
            // for the next set of vertex to be drawn
            if (draws[drawsCounter - 1].mode == RL_LINES) draws[drawsCounter - 1].vertexAlignment = ((draws[drawsCounter - 1].vertexCount < 4)? draws[drawsCounter - 1].vertexCount : draws[drawsCounter - 1].vertexCount%4);
            else if (draws[drawsCounter - 1].mode == RL_TRIANGLES) draws[drawsCounter - 1].vertexAlignment = ((draws[drawsCounter - 1].vertexCount < 4)? 1 : (4 - (draws[drawsCounter - 1].vertexCount%4)));

            else draws[drawsCounter - 1].vertexAlignment = 0;

            if (rlCheckBufferLimit(draws[drawsCounter - 1].vertexAlignment)) rlglDraw();
            else
            {
                vertexData[currentBuffer].vCounter += draws[drawsCounter - 1].vertexAlignment;
                vertexData[currentBuffer].cCounter += draws[drawsCounter - 1].vertexAlignment;
                vertexData[currentBuffer].tcCounter += draws[drawsCounter - 1].vertexAlignment;

                drawsCounter++;
            }
        }

        if (drawsCounter >= MAX_DRAWCALL_REGISTERED) rlglDraw();

        draws[drawsCounter - 1].textureId = id;
        draws[drawsCounter - 1].vertexCount = 0;
    }
#endif
}

// Disable texture usage
void rlDisableTexture(void)
{
#if defined(GRAPHICS_API_OPENGL_11)
    glDisable(GL_TEXTURE_2D);
    glBindTexture(GL_TEXTURE_2D, 0);
#else
    // NOTE: If quads batch limit is reached,
    // we force a draw call and next batch starts
    if (vertexData[currentBuffer].vCounter >= (MAX_BATCH_ELEMENTS*4)) rlglDraw();
#endif
}

// Set texture parameters (wrap mode/filter mode)
void rlTextureParameters(unsigned int id, int param, int value)
{
    glBindTexture(GL_TEXTURE_2D, id);

    switch (param)
    {
        case RL_TEXTURE_WRAP_S:
        case RL_TEXTURE_WRAP_T:
        {
            if (value == RL_WRAP_MIRROR_CLAMP)
            {
                if (texMirrorClampSupported) glTexParameteri(GL_TEXTURE_2D, param, value);
                else TraceLog(LOG_WARNING, "Clamp mirror wrap mode not supported");
            }
            else glTexParameteri(GL_TEXTURE_2D, param, value);

        } break;
        case RL_TEXTURE_MAG_FILTER:
        case RL_TEXTURE_MIN_FILTER: glTexParameteri(GL_TEXTURE_2D, param, value); break;
        case RL_TEXTURE_ANISOTROPIC_FILTER:
        {
#if !defined(GRAPHICS_API_OPENGL_11)
            if (value <= maxAnisotropicLevel) glTexParameterf(GL_TEXTURE_2D, GL_TEXTURE_MAX_ANISOTROPY_EXT, (float)value);
            else if (maxAnisotropicLevel > 0.0f)
            {
                TraceLog(LOG_WARNING, "[TEX ID %i] Maximum anisotropic filter level supported is %iX", id, maxAnisotropicLevel);
                glTexParameterf(GL_TEXTURE_2D, GL_TEXTURE_MAX_ANISOTROPY_EXT, (float)value);
            }
            else TraceLog(LOG_WARNING, "Anisotropic filtering not supported");
#endif
        } break;
        default: break;
    }

    glBindTexture(GL_TEXTURE_2D, 0);
}

// Enable rendering to texture (fbo)
void rlEnableRenderTexture(unsigned int id)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    glBindFramebuffer(GL_FRAMEBUFFER, id);

    //glDisable(GL_CULL_FACE);    // Allow double side drawing for texture flipping
    //glCullFace(GL_FRONT);
#endif
}

// Disable rendering to texture
void rlDisableRenderTexture(void)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    glBindFramebuffer(GL_FRAMEBUFFER, 0);

    //glEnable(GL_CULL_FACE);
    //glCullFace(GL_BACK);
#endif
}

// Enable depth test
void rlEnableDepthTest(void) { glEnable(GL_DEPTH_TEST); }

// Disable depth test
void rlDisableDepthTest(void) { glDisable(GL_DEPTH_TEST); }

// Enable backface culling
void rlEnableBackfaceCulling(void) { glEnable(GL_CULL_FACE); }

// Disable backface culling
void rlDisableBackfaceCulling(void) { glDisable(GL_CULL_FACE); }

// Enable scissor test
RLAPI void rlEnableScissorTest(void) { glEnable(GL_SCISSOR_TEST); }

// Disable scissor test
RLAPI void rlDisableScissorTest(void) { glDisable(GL_SCISSOR_TEST); }

// Scissor test
RLAPI void rlScissor(int x, int y, int width, int height) { glScissor(x, y, width, height); }

// Enable scissor test
void rlEnableScissorTest(void) {
    glEnable(GL_SCISSOR_TEST);
}

// Disable scissor test
void rlDisableScissorTest(void) {
    glDisable(GL_SCISSOR_TEST);
}

// Scissor test
void rlScissor(int x, int y, int width, int height) {
    glScissor(x, y, width, height);
}

// Enable wire mode
void rlEnableWireMode(void)
{
#if defined (GRAPHICS_API_OPENGL_11) || defined(GRAPHICS_API_OPENGL_33)
    // NOTE: glPolygonMode() not available on OpenGL ES
    glPolygonMode(GL_FRONT_AND_BACK, GL_LINE);
#endif
}

// Disable wire mode
void rlDisableWireMode(void)
{
#if defined (GRAPHICS_API_OPENGL_11) || defined(GRAPHICS_API_OPENGL_33)
    // NOTE: glPolygonMode() not available on OpenGL ES
    glPolygonMode(GL_FRONT_AND_BACK, GL_FILL);
#endif
}

// Unload texture from GPU memory
void rlDeleteTextures(unsigned int id)
{
    if (id > 0) glDeleteTextures(1, &id);
}

// Unload render texture from GPU memory
void rlDeleteRenderTextures(RenderTexture2D target)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if (target.texture.id > 0) glDeleteTextures(1, &target.texture.id);
    if (target.depth.id > 0)
    {
        if (target.depthTexture) glDeleteTextures(1, &target.depth.id);
        else glDeleteRenderbuffers(1, &target.depth.id);
    }

    if (target.id > 0) glDeleteFramebuffers(1, &target.id);

    TraceLog(LOG_INFO, "[FBO ID %i] Unloaded render texture data from VRAM (GPU)", target.id);
#endif
}

// Unload shader from GPU memory
void rlDeleteShader(unsigned int id)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if (id != 0) glDeleteProgram(id);
#endif
}

// Unload vertex data (VAO) from GPU memory
void rlDeleteVertexArrays(unsigned int id)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if (vaoSupported)
    {
        if (id != 0) glDeleteVertexArrays(1, &id);
        TraceLog(LOG_INFO, "[VAO ID %i] Unloaded model data from VRAM (GPU)", id);
    }
#endif
}

// Unload vertex data (VBO) from GPU memory
void rlDeleteBuffers(unsigned int id)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if (id != 0)
    {
        glDeleteBuffers(1, &id);
        if (!vaoSupported) TraceLog(LOG_INFO, "[VBO ID %i] Unloaded model vertex data from VRAM (GPU)", id);
    }
#endif
}

// Clear color buffer with color
void rlClearColor(byte r, byte g, byte b, byte a)
{
    // Color values clamp to 0.0f(0) and 1.0f(255)
    float cr = (float)r/255;
    float cg = (float)g/255;
    float cb = (float)b/255;
    float ca = (float)a/255;

    glClearColor(cr, cg, cb, ca);
}

// Clear used screen buffers (color and depth)
void rlClearScreenBuffers(void)
{
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);     // Clear used buffers: Color and Depth (Depth is used for 3D)
    //glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT | GL_STENCIL_BUFFER_BIT);     // Stencil buffer not used...
}

// Update GPU buffer with new data
void rlUpdateBuffer(int bufferId, void *data, int dataSize)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    glBindBuffer(GL_ARRAY_BUFFER, bufferId);
    glBufferSubData(GL_ARRAY_BUFFER, 0, dataSize, data);
#endif
}

//----------------------------------------------------------------------------------
// Module Functions Definition - rlgl Functions
//----------------------------------------------------------------------------------

// Initialize rlgl: OpenGL extensions, default buffers/shaders/textures, OpenGL states
void rlglInit(int width, int height)
{
    // Check OpenGL information and capabilities
    //------------------------------------------------------------------------------

    // Print current OpenGL and GLSL version
    TraceLog(LOG_INFO, "GPU: Vendor:   %s", glGetString(GL_VENDOR));
    TraceLog(LOG_INFO, "GPU: Renderer: %s", glGetString(GL_RENDERER));
    TraceLog(LOG_INFO, "GPU: Version:  %s", glGetString(GL_VERSION));
    TraceLog(LOG_INFO, "GPU: GLSL:     %s", glGetString(GL_SHADING_LANGUAGE_VERSION));

    // NOTE: We can get a bunch of extra information about GPU capabilities (glGet*)
    //int maxTexSize;
    //glGetIntegerv(GL_MAX_TEXTURE_SIZE, &maxTexSize);
    //TraceLog(LOG_INFO, "GL_MAX_TEXTURE_SIZE: %i", maxTexSize);

    //GL_MAX_TEXTURE_IMAGE_UNITS
    //GL_MAX_VIEWPORT_DIMS

    //int numAuxBuffers;
    //glGetIntegerv(GL_AUX_BUFFERS, &numAuxBuffers);
    //TraceLog(LOG_INFO, "GL_AUX_BUFFERS: %i", numAuxBuffers);

    //GLint numComp = 0;
    //GLint format[32] = { 0 };
    //glGetIntegerv(GL_NUM_COMPRESSED_TEXTURE_FORMATS, &numComp);
    //glGetIntegerv(GL_COMPRESSED_TEXTURE_FORMATS, format);
    //for (int i = 0; i < numComp; i++) TraceLog(LOG_INFO, "Supported compressed format: 0x%x", format[i]);

    // NOTE: We don't need that much data on screen... right now...

    // TODO: Automatize extensions loading using rlLoadExtensions() and GLAD
    // Actually, when rlglInit() is called in InitWindow() in core.c,
    // OpenGL required extensions have already been loaded (PLATFORM_DESKTOP)

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    // Get supported extensions list
    GLint numExt = 0;

#if defined(GRAPHICS_API_OPENGL_33)
    // NOTE: On OpenGL 3.3 VAO and NPOT are supported by default
    vaoSupported = true;

    // Multiple texture extensions supported by default
    texNPOTSupported = true;
    texFloatSupported = true;
    texDepthSupported = true;

    // We get a list of available extensions and we check for some of them (compressed textures)
    // NOTE: We don't need to check again supported extensions but we do (GLAD already dealt with that)
    glGetIntegerv(GL_NUM_EXTENSIONS, &numExt);

    // Allocate numExt strings pointers
    const char **extList = RL_MALLOC(sizeof(const char *)*numExt);

    // Get extensions strings
    for (int i = 0; i < numExt; i++) extList[i] = (const char *)glGetStringi(GL_EXTENSIONS, i);

#elif defined(GRAPHICS_API_OPENGL_ES2)
    // Allocate 512 strings pointers (2 KB)
    const char **extList = RL_MALLOC(sizeof(const char *)*512);

    const char *extensions = (const char *)glGetString(GL_EXTENSIONS);  // One big const string

    // NOTE: We have to duplicate string because glGetString() returns a const string
    int len = strlen(extensions) + 1;
    char *extensionsDup = (char *)RL_CALLOC(len, sizeof(char));
    strcpy(extensionsDup, extensions);

    extList[numExt] = extensionsDup;

    for (int i = 0; i < len; i++)
    {
        if (extensionsDup[i] == ' ')
        {
            extensionsDup[i] = '\0';

            numExt++;
            extList[numExt] = &extensionsDup[i + 1];
        }
    }

    // NOTE: Duplicated string (extensionsDup) must be deallocated
#endif

    TraceLog(LOG_INFO, "Number of supported extensions: %i", numExt);

    // Show supported extensions
    //for (int i = 0; i < numExt; i++)  TraceLog(LOG_INFO, "Supported extension: %s", extList[i]);

    // Check required extensions
    for (int i = 0; i < numExt; i++)
    {
#if defined(GRAPHICS_API_OPENGL_ES2)
        // Check VAO support
        // NOTE: Only check on OpenGL ES, OpenGL 3.3 has VAO support as core feature
        if (strcmp(extList[i], (const char *)"GL_OES_vertex_array_object") == 0)
        {
            // The extension is supported by our hardware and driver, try to get related functions pointers
            // NOTE: emscripten does not support VAOs natively, it uses emulation and it reduces overall performance...
            glGenVertexArrays = (PFNGLGENVERTEXARRAYSOESPROC)eglGetProcAddress("glGenVertexArraysOES");
            glBindVertexArray = (PFNGLBINDVERTEXARRAYOESPROC)eglGetProcAddress("glBindVertexArrayOES");
            glDeleteVertexArrays = (PFNGLDELETEVERTEXARRAYSOESPROC)eglGetProcAddress("glDeleteVertexArraysOES");
            //glIsVertexArray = (PFNGLISVERTEXARRAYOESPROC)eglGetProcAddress("glIsVertexArrayOES");     // NOTE: Fails in WebGL, omitted

            if ((glGenVertexArrays != NULL) && (glBindVertexArray != NULL) && (glDeleteVertexArrays != NULL)) vaoSupported = true;
        }

        // Check NPOT textures support
        // NOTE: Only check on OpenGL ES, OpenGL 3.3 has NPOT textures full support as core feature
        if (strcmp(extList[i], (const char *)"GL_OES_texture_npot") == 0) texNPOTSupported = true;

        // Check texture float support
        if (strcmp(extList[i], (const char *)"GL_OES_texture_float") == 0) texFloatSupported = true;

        // Check depth texture support
        if ((strcmp(extList[i], (const char *)"GL_OES_depth_texture") == 0) ||
            (strcmp(extList[i], (const char *)"GL_WEBGL_depth_texture") == 0)) texDepthSupported = true;

        if (strcmp(extList[i], (const char *)"GL_OES_depth24") == 0) maxDepthBits = 24;
        if (strcmp(extList[i], (const char *)"GL_OES_depth32") == 0) maxDepthBits = 32;
#endif
        // DDS texture compression support
        if ((strcmp(extList[i], (const char *)"GL_EXT_texture_compression_s3tc") == 0) ||
            (strcmp(extList[i], (const char *)"GL_WEBGL_compressed_texture_s3tc") == 0) ||
            (strcmp(extList[i], (const char *)"GL_WEBKIT_WEBGL_compressed_texture_s3tc") == 0)) texCompDXTSupported = true;

        // ETC1 texture compression support
        if ((strcmp(extList[i], (const char *)"GL_OES_compressed_ETC1_RGB8_texture") == 0) ||
            (strcmp(extList[i], (const char *)"GL_WEBGL_compressed_texture_etc1") == 0)) texCompETC1Supported = true;

        // ETC2/EAC texture compression support
        if (strcmp(extList[i], (const char *)"GL_ARB_ES3_compatibility") == 0) texCompETC2Supported = true;

        // PVR texture compression support
        if (strcmp(extList[i], (const char *)"GL_IMG_texture_compression_pvrtc") == 0) texCompPVRTSupported = true;

        // ASTC texture compression support
        if (strcmp(extList[i], (const char *)"GL_KHR_texture_compression_astc_hdr") == 0) texCompASTCSupported = true;

        // Anisotropic texture filter support
        if (strcmp(extList[i], (const char *)"GL_EXT_texture_filter_anisotropic") == 0)
        {
            texAnisotropicFilterSupported = true;
            glGetFloatv(0x84FF, &maxAnisotropicLevel);   // GL_MAX_TEXTURE_MAX_ANISOTROPY_EXT
        }

        // Clamp mirror wrap mode supported
        if (strcmp(extList[i], (const char *)"GL_EXT_texture_mirror_clamp") == 0) texMirrorClampSupported = true;

        // Debug marker support
        if (strcmp(extList[i], (const char *)"GL_EXT_debug_marker") == 0) debugMarkerSupported = true;
    }

    // Free extensions pointers
    RL_FREE(extList);

#if defined(GRAPHICS_API_OPENGL_ES2)
    RL_FREE(extensionsDup);    // Duplicated string must be deallocated

    if (vaoSupported) TraceLog(LOG_INFO, "[EXTENSION] VAO extension detected, VAO functions initialized successfully");
    else TraceLog(LOG_WARNING, "[EXTENSION] VAO extension not found, VAO usage not supported");

    if (texNPOTSupported) TraceLog(LOG_INFO, "[EXTENSION] NPOT textures extension detected, full NPOT textures supported");
    else TraceLog(LOG_WARNING, "[EXTENSION] NPOT textures extension not found, limited NPOT support (no-mipmaps, no-repeat)");
#endif

    if (texCompDXTSupported) TraceLog(LOG_INFO, "[EXTENSION] DXT compressed textures supported");
    if (texCompETC1Supported) TraceLog(LOG_INFO, "[EXTENSION] ETC1 compressed textures supported");
    if (texCompETC2Supported) TraceLog(LOG_INFO, "[EXTENSION] ETC2/EAC compressed textures supported");
    if (texCompPVRTSupported) TraceLog(LOG_INFO, "[EXTENSION] PVRT compressed textures supported");
    if (texCompASTCSupported) TraceLog(LOG_INFO, "[EXTENSION] ASTC compressed textures supported");

    if (texAnisotropicFilterSupported) TraceLog(LOG_INFO, "[EXTENSION] Anisotropic textures filtering supported (max: %.0fX)", maxAnisotropicLevel);
    if (texMirrorClampSupported) TraceLog(LOG_INFO, "[EXTENSION] Mirror clamp wrap texture mode supported");

    if (debugMarkerSupported) TraceLog(LOG_INFO, "[EXTENSION] Debug Marker supported");

    // Initialize buffers, default shaders and default textures
    //----------------------------------------------------------
    // Init default white texture
    unsigned char pixels[4] = { 255, 255, 255, 255 };   // 1 pixel RGBA (4 bytes)
    defaultTextureId = rlLoadTexture(pixels, 1, 1, UNCOMPRESSED_R8G8B8A8, 1);

    if (defaultTextureId != 0) TraceLog(LOG_INFO, "[TEX ID %i] Base white texture loaded successfully", defaultTextureId);
    else TraceLog(LOG_WARNING, "Base white texture could not be loaded");

    // Init default Shader (customized for GL 3.3 and ES2)
    defaultShader = LoadShaderDefault();
    currentShader = defaultShader;

    // Init default vertex arrays buffers
    LoadBuffersDefault();

    // Init transformations matrix accumulator
    transformMatrix = MatrixIdentity();

    // Init draw calls tracking system
    draws = (DrawCall *)RL_MALLOC(sizeof(DrawCall)*MAX_DRAWCALL_REGISTERED);

    for (int i = 0; i < MAX_DRAWCALL_REGISTERED; i++)
    {
        draws[i].mode = RL_QUADS;
        draws[i].vertexCount = 0;
        draws[i].vertexAlignment = 0;
        //draws[i].vaoId = 0;
        //draws[i].shaderId = 0;
        draws[i].textureId = defaultTextureId;
        //draws[i].projection = MatrixIdentity();
        //draws[i].modelview = MatrixIdentity();
    }

    drawsCounter = 1;

    // Init internal matrix stack (emulating OpenGL 1.1)
    for (int i = 0; i < MAX_MATRIX_STACK_SIZE; i++) stack[i] = MatrixIdentity();

    // Init internal projection and modelview matrices
    projection = MatrixIdentity();
    modelview = MatrixIdentity();
    currentMatrix = &modelview;
#endif      // GRAPHICS_API_OPENGL_33 || GRAPHICS_API_OPENGL_ES2

    // Initialize OpenGL default states
    //----------------------------------------------------------
    // Init state: Depth test
    glDepthFunc(GL_LEQUAL);                                 // Type of depth testing to apply
    glDisable(GL_DEPTH_TEST);                               // Disable depth testing for 2D (only used for 3D)

    // Init state: Blending mode
    glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);      // Color blending function (how colors are mixed)
    glEnable(GL_BLEND);                                     // Enable color blending (required to work with transparencies)

    // Init state: Culling
    // NOTE: All shapes/models triangles are drawn CCW
    glCullFace(GL_BACK);                                    // Cull the back face (default)
    glFrontFace(GL_CCW);                                    // Front face are defined counter clockwise (default)
    glEnable(GL_CULL_FACE);                                 // Enable backface culling

#if defined(GRAPHICS_API_OPENGL_11)
    // Init state: Color hints (deprecated in OpenGL 3.0+)
    glHint(GL_PERSPECTIVE_CORRECTION_HINT, GL_NICEST);      // Improve quality of color and texture coordinate interpolation
    glShadeModel(GL_SMOOTH);                                // Smooth shading between vertex (vertex colors interpolation)
#endif

    // Init state: Color/Depth buffers clear
    glClearColor(0.0f, 0.0f, 0.0f, 1.0f);                   // Set clear color (black)
    glClearDepth(1.0f);                                     // Set clear depth value (default)
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);     // Clear color and depth buffers (depth buffer required for 3D)

    // Store screen size into global variables
    framebufferWidth = width;
    framebufferHeight = height;

    TraceLog(LOG_INFO, "OpenGL default states initialized successfully");
}

// Vertex Buffer Object deinitialization (memory free)
void rlglClose(void)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    UnloadShaderDefault();              // Unload default shader
    UnloadBuffersDefault();             // Unload default buffers
    glDeleteTextures(1, &defaultTextureId); // Unload default texture

    TraceLog(LOG_INFO, "[TEX ID %i] Unloaded texture data (base white texture) from VRAM", defaultTextureId);

    RL_FREE(draws);
#endif
}

// Update and draw internal buffers
void rlglDraw(void)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    // Only process data if we have data to process
    if (vertexData[currentBuffer].vCounter > 0)
    {
        UpdateBuffersDefault();
        DrawBuffersDefault();       // NOTE: Stereo rendering is checked inside
    }
#endif
}

// Returns current OpenGL version
int rlGetVersion(void)
{
#if defined(GRAPHICS_API_OPENGL_11)
    return OPENGL_11;
#elif defined(GRAPHICS_API_OPENGL_21)
    #if defined(__APPLE__)
        return OPENGL_33;       // NOTE: Force OpenGL 3.3 on OSX
    #else
        return OPENGL_21;
    #endif
#elif defined(GRAPHICS_API_OPENGL_33)
    return OPENGL_33;
#elif defined(GRAPHICS_API_OPENGL_ES2)
    return OPENGL_ES_20;
#endif
}

// Check internal buffer overflow for a given number of vertex
bool rlCheckBufferLimit(int vCount)
{
    bool overflow = false;
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if ((vertexData[currentBuffer].vCounter + vCount) >= (MAX_BATCH_ELEMENTS*4)) overflow = true;
#endif
    return overflow;
}

// Set debug marker
void rlSetDebugMarker(const char *text)
{
#if defined(GRAPHICS_API_OPENGL_33)
    if (debugMarkerSupported) glInsertEventMarkerEXT(0, text);
#endif
}

// Load OpenGL extensions
// NOTE: External loader function could be passed as a pointer
void rlLoadExtensions(void *loader)
{
#if defined(GRAPHICS_API_OPENGL_33)
    // NOTE: glad is generated and contains only required OpenGL 3.3 Core extensions (and lower versions)
    #if !defined(__APPLE__)
        if (!gladLoadGLLoader((GLADloadproc)loader)) TraceLog(LOG_WARNING, "GLAD: Cannot load OpenGL extensions");
        else TraceLog(LOG_INFO, "GLAD: OpenGL extensions loaded successfully");

        #if defined(GRAPHICS_API_OPENGL_21)
        if (GLAD_GL_VERSION_2_1) TraceLog(LOG_INFO, "OpenGL 2.1 profile supported");
        #elif defined(GRAPHICS_API_OPENGL_33)
        if (GLAD_GL_VERSION_3_3) TraceLog(LOG_INFO, "OpenGL 3.3 Core profile supported");
        else TraceLog(LOG_ERROR, "OpenGL 3.3 Core profile not supported");
        #endif
    #endif

    // With GLAD, we can check if an extension is supported using the GLAD_GL_xxx booleans
    //if (GLAD_GL_ARB_vertex_array_object) // Use GL_ARB_vertex_array_object
#endif
}

// Get world coordinates from screen coordinates
Vector3 rlUnproject(Vector3 source, Matrix proj, Matrix view)
{
    Vector3 result = { 0.0f, 0.0f, 0.0f };

    // Calculate unproject matrix (multiply view patrix by projection matrix) and invert it
    Matrix matViewProj = MatrixMultiply(view, proj);
    matViewProj = MatrixInvert(matViewProj);

    // Create quaternion from source point
    Quaternion quat = { source.x, source.y, source.z, 1.0f };

    // Multiply quat point by unproject matrix
    quat = QuaternionTransform(quat, matViewProj);

    // Normalized world points in vectors
    result.x = quat.x/quat.w;
    result.y = quat.y/quat.w;
    result.z = quat.z/quat.w;

    return result;
}

// Convert image data to OpenGL texture (returns OpenGL valid Id)
unsigned int rlLoadTexture(void *data, int width, int height, int format, int mipmapCount)
{
    glBindTexture(GL_TEXTURE_2D, 0);    // Free any old binding

    unsigned int id = 0;

    // Check texture format support by OpenGL 1.1 (compressed textures not supported)
#if defined(GRAPHICS_API_OPENGL_11)
    if (format >= COMPRESSED_DXT1_RGB)
    {
        TraceLog(LOG_WARNING, "OpenGL 1.1 does not support GPU compressed texture formats");
        return id;
    }
#else
    if ((!texCompDXTSupported) && ((format == COMPRESSED_DXT1_RGB) || (format == COMPRESSED_DXT1_RGBA) ||
        (format == COMPRESSED_DXT3_RGBA) || (format == COMPRESSED_DXT5_RGBA)))
    {
        TraceLog(LOG_WARNING, "DXT compressed texture format not supported");
        return id;
    }
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if ((!texCompETC1Supported) && (format == COMPRESSED_ETC1_RGB))
    {
        TraceLog(LOG_WARNING, "ETC1 compressed texture format not supported");
        return id;
    }

    if ((!texCompETC2Supported) && ((format == COMPRESSED_ETC2_RGB) || (format == COMPRESSED_ETC2_EAC_RGBA)))
    {
        TraceLog(LOG_WARNING, "ETC2 compressed texture format not supported");
        return id;
    }

    if ((!texCompPVRTSupported) && ((format == COMPRESSED_PVRT_RGB) || (format == COMPRESSED_PVRT_RGBA)))
    {
        TraceLog(LOG_WARNING, "PVRT compressed texture format not supported");
        return id;
    }

    if ((!texCompASTCSupported) && ((format == COMPRESSED_ASTC_4x4_RGBA) || (format == COMPRESSED_ASTC_8x8_RGBA)))
    {
        TraceLog(LOG_WARNING, "ASTC compressed texture format not supported");
        return id;
    }
#endif
#endif      // GRAPHICS_API_OPENGL_11

    glPixelStorei(GL_UNPACK_ALIGNMENT, 1);

    glGenTextures(1, &id);              // Generate texture id

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    //glActiveTexture(GL_TEXTURE0);     // If not defined, using GL_TEXTURE0 by default (shader texture)
#endif

    glBindTexture(GL_TEXTURE_2D, id);

    int mipWidth = width;
    int mipHeight = height;
    int mipOffset = 0;          // Mipmap data offset

    TraceLog(LOG_DEBUG, "Load texture from data memory address: 0x%x", data);

    // Load the different mipmap levels
    for (int i = 0; i < mipmapCount; i++)
    {
        unsigned int mipSize = GetPixelDataSize(mipWidth, mipHeight, format);

        unsigned int glInternalFormat, glFormat, glType;
        rlGetGlTextureFormats(format, &glInternalFormat, &glFormat, &glType);

        TraceLog(LOG_DEBUG, "Load mipmap level %i (%i x %i), size: %i, offset: %i", i, mipWidth, mipHeight, mipSize, mipOffset);

        if (glInternalFormat != -1)
        {
            if (format < COMPRESSED_DXT1_RGB) glTexImage2D(GL_TEXTURE_2D, i, glInternalFormat, mipWidth, mipHeight, 0, glFormat, glType, (unsigned char *)data + mipOffset);
        #if !defined(GRAPHICS_API_OPENGL_11)
            else glCompressedTexImage2D(GL_TEXTURE_2D, i, glInternalFormat, mipWidth, mipHeight, 0, mipSize, (unsigned char *)data + mipOffset);
        #endif

        #if defined(GRAPHICS_API_OPENGL_33)
            if (format == UNCOMPRESSED_GRAYSCALE)
            {
                GLint swizzleMask[] = { GL_RED, GL_RED, GL_RED, GL_ONE };
                glTexParameteriv(GL_TEXTURE_2D, GL_TEXTURE_SWIZZLE_RGBA, swizzleMask);
            }
            else if (format == UNCOMPRESSED_GRAY_ALPHA)
            {
            #if defined(GRAPHICS_API_OPENGL_21)
                GLint swizzleMask[] = { GL_RED, GL_RED, GL_RED, GL_ALPHA };
            #elif defined(GRAPHICS_API_OPENGL_33)
                GLint swizzleMask[] = { GL_RED, GL_RED, GL_RED, GL_GREEN };
            #endif
                glTexParameteriv(GL_TEXTURE_2D, GL_TEXTURE_SWIZZLE_RGBA, swizzleMask);
            }
        #endif
        }

        mipWidth /= 2;
        mipHeight /= 2;
        mipOffset += mipSize;

        // Security check for NPOT textures
        if (mipWidth < 1) mipWidth = 1;
        if (mipHeight < 1) mipHeight = 1;
    }

    // Texture parameters configuration
    // NOTE: glTexParameteri does NOT affect texture uploading, just the way it's used
#if defined(GRAPHICS_API_OPENGL_ES2)
    // NOTE: OpenGL ES 2.0 with no GL_OES_texture_npot support (i.e. WebGL) has limited NPOT support, so CLAMP_TO_EDGE must be used
    if (texNPOTSupported)
    {
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_REPEAT);       // Set texture to repeat on x-axis
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_REPEAT);       // Set texture to repeat on y-axis
    }
    else
    {
        // NOTE: If using negative texture coordinates (LoadOBJ()), it does not work!
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE);       // Set texture to clamp on x-axis
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE);       // Set texture to clamp on y-axis
    }
#else
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_REPEAT);       // Set texture to repeat on x-axis
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_REPEAT);       // Set texture to repeat on y-axis
#endif

    // Magnification and minification filters
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);  // Alternative: GL_LINEAR
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);  // Alternative: GL_LINEAR

#if defined(GRAPHICS_API_OPENGL_33)
    if (mipmapCount > 1)
    {
        // Activate Trilinear filtering if mipmaps are available
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR_MIPMAP_LINEAR);
    }
#endif

    // At this point we have the texture loaded in GPU and texture parameters configured

    // NOTE: If mipmaps were not in data, they are not generated automatically

    // Unbind current texture
    glBindTexture(GL_TEXTURE_2D, 0);

    if (id > 0) TraceLog(LOG_INFO, "[TEX ID %i] Texture created successfully (%ix%i - %i mipmaps)", id, width, height, mipmapCount);
    else TraceLog(LOG_WARNING, "Texture could not be created");

    return id;
}

// Load depth texture/renderbuffer (to be attached to fbo)
// WARNING: OpenGL ES 2.0 requires GL_OES_depth_texture/WEBGL_depth_texture extensions
unsigned int rlLoadTextureDepth(int width, int height, int bits, bool useRenderBuffer)
{
    unsigned int id = 0;

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    unsigned int glInternalFormat = GL_DEPTH_COMPONENT16;

    if ((bits != 16) && (bits != 24) && (bits != 32)) bits = 16;

    if (bits == 24)
    {
#if defined(GRAPHICS_API_OPENGL_33)
        glInternalFormat = GL_DEPTH_COMPONENT24;
#elif defined(GRAPHICS_API_OPENGL_ES2)
        if (maxDepthBits >= 24) glInternalFormat = GL_DEPTH_COMPONENT24_OES;
#endif
    }

    if (bits == 32)
    {
#if defined(GRAPHICS_API_OPENGL_33)
        glInternalFormat = GL_DEPTH_COMPONENT32;
#elif defined(GRAPHICS_API_OPENGL_ES2)
        if (maxDepthBits == 32) glInternalFormat = GL_DEPTH_COMPONENT32_OES;
#endif
    }

    if (!useRenderBuffer && texDepthSupported)
    {
        glGenTextures(1, &id);
        glBindTexture(GL_TEXTURE_2D, id);
        glTexImage2D(GL_TEXTURE_2D, 0, glInternalFormat, width, height, 0, GL_DEPTH_COMPONENT, GL_UNSIGNED_INT, NULL);

        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE);
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE);

        glBindTexture(GL_TEXTURE_2D, 0);
    }
    else
    {
        // Create the renderbuffer that will serve as the depth attachment for the framebuffer
        // NOTE: A renderbuffer is simpler than a texture and could offer better performance on embedded devices
        glGenRenderbuffers(1, &id);
        glBindRenderbuffer(GL_RENDERBUFFER, id);
        glRenderbufferStorage(GL_RENDERBUFFER, glInternalFormat, width, height);

        glBindRenderbuffer(GL_RENDERBUFFER, 0);
    }
#endif

    return id;
}

// Load texture cubemap
// NOTE: Cubemap data is expected to be 6 images in a single column,
// expected the following convention: +X, -X, +Y, -Y, +Z, -Z
unsigned int rlLoadTextureCubemap(void *data, int size, int format)
{
    unsigned int cubemapId = 0;
    unsigned int dataSize = GetPixelDataSize(size, size, format);

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    glGenTextures(1, &cubemapId);
    glBindTexture(GL_TEXTURE_CUBE_MAP, cubemapId);

    unsigned int glInternalFormat, glFormat, glType;
    rlGetGlTextureFormats(format, &glInternalFormat, &glFormat, &glType);

    if (glInternalFormat != -1)
    {
        // Load cubemap faces
        for (unsigned int i = 0; i < 6; i++)
        {
            if (format < COMPRESSED_DXT1_RGB) glTexImage2D(GL_TEXTURE_CUBE_MAP_POSITIVE_X + i, 0, glInternalFormat, size, size, 0, glFormat, glType, (unsigned char *)data + i*dataSize);
#if !defined(GRAPHICS_API_OPENGL_11)
            else glCompressedTexImage2D(GL_TEXTURE_CUBE_MAP_POSITIVE_X + i, 0, glInternalFormat, size, size, 0, dataSize, (unsigned char *)data + i*dataSize);
#endif
#if defined(GRAPHICS_API_OPENGL_33)
            if (format == UNCOMPRESSED_GRAYSCALE)
            {
                GLint swizzleMask[] = { GL_RED, GL_RED, GL_RED, GL_ONE };
                glTexParameteriv(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_SWIZZLE_RGBA, swizzleMask);
            }
            else if (format == UNCOMPRESSED_GRAY_ALPHA)
            {
#if defined(GRAPHICS_API_OPENGL_21)
                GLint swizzleMask[] = { GL_RED, GL_RED, GL_RED, GL_ALPHA };
#elif defined(GRAPHICS_API_OPENGL_33)
                GLint swizzleMask[] = { GL_RED, GL_RED, GL_RED, GL_GREEN };
#endif
                glTexParameteriv(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_SWIZZLE_RGBA, swizzleMask);
            }
#endif
        }
    }

    // Set cubemap texture sampling parameters
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE);
#if defined(GRAPHICS_API_OPENGL_33)
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_R, GL_CLAMP_TO_EDGE);  // Flag not supported on OpenGL ES 2.0
#endif

    glBindTexture(GL_TEXTURE_CUBE_MAP, 0);
#endif

    return cubemapId;
}

// Update already loaded texture in GPU with new data
// NOTE: We don't know safely if internal texture format is the expected one...
void rlUpdateTexture(unsigned int id, int width, int height, int format, const void *data)
{
    glBindTexture(GL_TEXTURE_2D, id);

    unsigned int glInternalFormat, glFormat, glType;
    rlGetGlTextureFormats(format, &glInternalFormat, &glFormat, &glType);

    if ((glInternalFormat != -1) && (format < COMPRESSED_DXT1_RGB))
    {
        glTexSubImage2D(GL_TEXTURE_2D, 0, 0, 0, width, height, glFormat, glType, (unsigned char *)data);
    }
    else TraceLog(LOG_WARNING, "Texture format updating not supported");
}

// Get OpenGL internal formats and data type from raylib PixelFormat
void rlGetGlTextureFormats(int format, unsigned int *glInternalFormat, unsigned int *glFormat, unsigned int *glType)
{
    *glInternalFormat = -1;
    *glFormat = -1;
    *glType = -1;

    switch (format)
    {
    #if defined(GRAPHICS_API_OPENGL_11) || defined(GRAPHICS_API_OPENGL_21) || defined(GRAPHICS_API_OPENGL_ES2)
        // NOTE: on OpenGL ES 2.0 (WebGL), internalFormat must match format and options allowed are: GL_LUMINANCE, GL_RGB, GL_RGBA
        case UNCOMPRESSED_GRAYSCALE: *glInternalFormat = GL_LUMINANCE; *glFormat = GL_LUMINANCE; *glType = GL_UNSIGNED_BYTE; break;
        case UNCOMPRESSED_GRAY_ALPHA: *glInternalFormat = GL_LUMINANCE_ALPHA; *glFormat = GL_LUMINANCE_ALPHA; *glType = GL_UNSIGNED_BYTE; break;
        case UNCOMPRESSED_R5G6B5: *glInternalFormat = GL_RGB; *glFormat = GL_RGB; *glType = GL_UNSIGNED_SHORT_5_6_5; break;
        case UNCOMPRESSED_R8G8B8: *glInternalFormat = GL_RGB; *glFormat = GL_RGB; *glType = GL_UNSIGNED_BYTE; break;
        case UNCOMPRESSED_R5G5B5A1: *glInternalFormat = GL_RGBA; *glFormat = GL_RGBA; *glType = GL_UNSIGNED_SHORT_5_5_5_1; break;
        case UNCOMPRESSED_R4G4B4A4: *glInternalFormat = GL_RGBA; *glFormat = GL_RGBA; *glType = GL_UNSIGNED_SHORT_4_4_4_4; break;
        case UNCOMPRESSED_R8G8B8A8: *glInternalFormat = GL_RGBA; *glFormat = GL_RGBA; *glType = GL_UNSIGNED_BYTE; break;
        #if !defined(GRAPHICS_API_OPENGL_11)
        case UNCOMPRESSED_R32: if (texFloatSupported) *glInternalFormat = GL_LUMINANCE; *glFormat = GL_LUMINANCE; *glType = GL_FLOAT; break;   // NOTE: Requires extension OES_texture_float
        case UNCOMPRESSED_R32G32B32: if (texFloatSupported) *glInternalFormat = GL_RGB; *glFormat = GL_RGB; *glType = GL_FLOAT; break;         // NOTE: Requires extension OES_texture_float
        case UNCOMPRESSED_R32G32B32A32: if (texFloatSupported) *glInternalFormat = GL_RGBA; *glFormat = GL_RGBA; *glType = GL_FLOAT; break;    // NOTE: Requires extension OES_texture_float
        #endif
    #elif defined(GRAPHICS_API_OPENGL_33)
        case UNCOMPRESSED_GRAYSCALE: *glInternalFormat = GL_R8; *glFormat = GL_RED; *glType = GL_UNSIGNED_BYTE; break;
        case UNCOMPRESSED_GRAY_ALPHA: *glInternalFormat = GL_RG8; *glFormat = GL_RG; *glType = GL_UNSIGNED_BYTE; break;
        case UNCOMPRESSED_R5G6B5: *glInternalFormat = GL_RGB565; *glFormat = GL_RGB; *glType = GL_UNSIGNED_SHORT_5_6_5; break;
        case UNCOMPRESSED_R8G8B8: *glInternalFormat = GL_RGB8; *glFormat = GL_RGB; *glType = GL_UNSIGNED_BYTE; break;
        case UNCOMPRESSED_R5G5B5A1: *glInternalFormat = GL_RGB5_A1; *glFormat = GL_RGBA; *glType = GL_UNSIGNED_SHORT_5_5_5_1; break;
        case UNCOMPRESSED_R4G4B4A4: *glInternalFormat = GL_RGBA4; *glFormat = GL_RGBA; *glType = GL_UNSIGNED_SHORT_4_4_4_4; break;
        case UNCOMPRESSED_R8G8B8A8: *glInternalFormat = GL_RGBA8; *glFormat = GL_RGBA; *glType = GL_UNSIGNED_BYTE; break;
        case UNCOMPRESSED_R32: if (texFloatSupported) *glInternalFormat = GL_R32F; *glFormat = GL_RED; *glType = GL_FLOAT; break;
        case UNCOMPRESSED_R32G32B32: if (texFloatSupported) *glInternalFormat = GL_RGB32F; *glFormat = GL_RGB; *glType = GL_FLOAT; break;
        case UNCOMPRESSED_R32G32B32A32: if (texFloatSupported) *glInternalFormat = GL_RGBA32F; *glFormat = GL_RGBA; *glType = GL_FLOAT; break;
    #endif
        #if !defined(GRAPHICS_API_OPENGL_11)
        case COMPRESSED_DXT1_RGB: if (texCompDXTSupported) *glInternalFormat = GL_COMPRESSED_RGB_S3TC_DXT1_EXT; break;
        case COMPRESSED_DXT1_RGBA: if (texCompDXTSupported) *glInternalFormat = GL_COMPRESSED_RGBA_S3TC_DXT1_EXT; break;
        case COMPRESSED_DXT3_RGBA: if (texCompDXTSupported) *glInternalFormat = GL_COMPRESSED_RGBA_S3TC_DXT3_EXT; break;
        case COMPRESSED_DXT5_RGBA: if (texCompDXTSupported) *glInternalFormat = GL_COMPRESSED_RGBA_S3TC_DXT5_EXT; break;
        case COMPRESSED_ETC1_RGB: if (texCompETC1Supported) *glInternalFormat = GL_ETC1_RGB8_OES; break;                      // NOTE: Requires OpenGL ES 2.0 or OpenGL 4.3
        case COMPRESSED_ETC2_RGB: if (texCompETC2Supported) *glInternalFormat = GL_COMPRESSED_RGB8_ETC2; break;               // NOTE: Requires OpenGL ES 3.0 or OpenGL 4.3
        case COMPRESSED_ETC2_EAC_RGBA: if (texCompETC2Supported) *glInternalFormat = GL_COMPRESSED_RGBA8_ETC2_EAC; break;     // NOTE: Requires OpenGL ES 3.0 or OpenGL 4.3
        case COMPRESSED_PVRT_RGB: if (texCompPVRTSupported) *glInternalFormat = GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG; break;    // NOTE: Requires PowerVR GPU
        case COMPRESSED_PVRT_RGBA: if (texCompPVRTSupported) *glInternalFormat = GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG; break;  // NOTE: Requires PowerVR GPU
        case COMPRESSED_ASTC_4x4_RGBA: if (texCompASTCSupported) *glInternalFormat = GL_COMPRESSED_RGBA_ASTC_4x4_KHR; break;  // NOTE: Requires OpenGL ES 3.1 or OpenGL 4.3
        case COMPRESSED_ASTC_8x8_RGBA: if (texCompASTCSupported) *glInternalFormat = GL_COMPRESSED_RGBA_ASTC_8x8_KHR; break;  // NOTE: Requires OpenGL ES 3.1 or OpenGL 4.3
        #endif
        default: TraceLog(LOG_WARNING, "Texture format not supported"); break;
    }
}

// Unload texture from GPU memory
void rlUnloadTexture(unsigned int id)
{
    if (id > 0) glDeleteTextures(1, &id);
}

// Load a texture to be used for rendering (fbo with default color and depth attachments)
// NOTE: If colorFormat or depthBits are no supported, no attachment is done
RenderTexture2D rlLoadRenderTexture(int width, int height, int format, int depthBits, bool useDepthTexture)
{
    RenderTexture2D target = { 0 };

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if (useDepthTexture && texDepthSupported) target.depthTexture = true;

    // Create the framebuffer object
    glGenFramebuffers(1, &target.id);
    glBindFramebuffer(GL_FRAMEBUFFER, target.id);

    // Create fbo color texture attachment
    //-----------------------------------------------------------------------------------------------------
    if ((format != -1) && (format < COMPRESSED_DXT1_RGB))
    {
        // WARNING: Some texture formats are not supported for fbo color attachment
        target.texture.id = rlLoadTexture(NULL, width, height, format, 1);
        target.texture.width = width;
        target.texture.height = height;
        target.texture.format = format;
        target.texture.mipmaps = 1;
    }
    //-----------------------------------------------------------------------------------------------------

    // Create fbo depth renderbuffer/texture
    //-----------------------------------------------------------------------------------------------------
    if (depthBits > 0)
    {
        target.depth.id = rlLoadTextureDepth(width, height, depthBits, !useDepthTexture);
        target.depth.width = width;
        target.depth.height = height;
        target.depth.format = 19;       //DEPTH_COMPONENT_24BIT?
        target.depth.mipmaps = 1;
    }
    //-----------------------------------------------------------------------------------------------------

    // Attach color texture and depth renderbuffer to FBO
    //-----------------------------------------------------------------------------------------------------
    rlRenderTextureAttach(target, target.texture.id, 0);    // COLOR attachment
    rlRenderTextureAttach(target, target.depth.id, 1);      // DEPTH attachment
    //-----------------------------------------------------------------------------------------------------

    // Check if fbo is complete with attachments (valid)
    //-----------------------------------------------------------------------------------------------------
    if (rlRenderTextureComplete(target)) TraceLog(LOG_INFO, "[FBO ID %i] Framebuffer object created successfully", target.id);
    //-----------------------------------------------------------------------------------------------------

    glBindFramebuffer(GL_FRAMEBUFFER, 0);
#endif

    return target;
}

// Attach color buffer texture to an fbo (unloads previous attachment)
// NOTE: Attach type: 0-Color, 1-Depth renderbuffer, 2-Depth texture
void rlRenderTextureAttach(RenderTexture2D target, unsigned int id, int attachType)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    glBindFramebuffer(GL_FRAMEBUFFER, target.id);

    if (attachType == 0) glFramebufferTexture2D(GL_FRAMEBUFFER, GL_COLOR_ATTACHMENT0, GL_TEXTURE_2D, id, 0);
    else if (attachType == 1)
    {
        if (target.depthTexture) glFramebufferTexture2D(GL_FRAMEBUFFER, GL_DEPTH_ATTACHMENT, GL_TEXTURE_2D, id, 0);
        else glFramebufferRenderbuffer(GL_FRAMEBUFFER, GL_DEPTH_ATTACHMENT, GL_RENDERBUFFER, id);
    }

    glBindFramebuffer(GL_FRAMEBUFFER, 0);
#endif
}

// Verify render texture is complete
bool rlRenderTextureComplete(RenderTexture target)
{
    bool result = false;

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    glBindFramebuffer(GL_FRAMEBUFFER, target.id);

    GLenum status = glCheckFramebufferStatus(GL_FRAMEBUFFER);

    if (status != GL_FRAMEBUFFER_COMPLETE)
    {
        switch (status)
        {
            case GL_FRAMEBUFFER_UNSUPPORTED: TraceLog(LOG_WARNING, "Framebuffer is unsupported"); break;
            case GL_FRAMEBUFFER_INCOMPLETE_ATTACHMENT: TraceLog(LOG_WARNING, "Framebuffer has incomplete attachment"); break;
#if defined(GRAPHICS_API_OPENGL_ES2)
            case GL_FRAMEBUFFER_INCOMPLETE_DIMENSIONS: TraceLog(LOG_WARNING, "Framebuffer has incomplete dimensions"); break;
#endif
            case GL_FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT: TraceLog(LOG_WARNING, "Framebuffer has a missing attachment"); break;
            default: break;
        }
    }

    glBindFramebuffer(GL_FRAMEBUFFER, 0);

    result = (status == GL_FRAMEBUFFER_COMPLETE);
#endif

    return result;
}

// Generate mipmap data for selected texture
void rlGenerateMipmaps(Texture2D *texture)
{
    glBindTexture(GL_TEXTURE_2D, texture->id);

    // Check if texture is power-of-two (POT)
    bool texIsPOT = false;

    if (((texture->width > 0) && ((texture->width & (texture->width - 1)) == 0)) &&
        ((texture->height > 0) && ((texture->height & (texture->height - 1)) == 0))) texIsPOT = true;

#if defined(GRAPHICS_API_OPENGL_11)
    if (texIsPOT)
    {
        // WARNING: Manual mipmap generation only works for RGBA 32bit textures!
        if (texture->format == UNCOMPRESSED_R8G8B8A8)
        {
            // Retrieve texture data from VRAM
            void *data = rlReadTexturePixels(*texture);

            // NOTE: data size is reallocated to fit mipmaps data
            // NOTE: CPU mipmap generation only supports RGBA 32bit data
            int mipmapCount = GenerateMipmaps(data, texture->width, texture->height);

            int size = texture->width*texture->height*4;
            int offset = size;

            int mipWidth = texture->width/2;
            int mipHeight = texture->height/2;

            // Load the mipmaps
            for (int level = 1; level < mipmapCount; level++)
            {
                glTexImage2D(GL_TEXTURE_2D, level, GL_RGBA8, mipWidth, mipHeight, 0, GL_RGBA, GL_UNSIGNED_BYTE, (unsigned char *)data + offset);

                size = mipWidth*mipHeight*4;
                offset += size;

                mipWidth /= 2;
                mipHeight /= 2;
            }

            texture->mipmaps = mipmapCount + 1;
            RL_FREE(data); // Once mipmaps have been generated and data has been uploaded to GPU VRAM, we can discard RAM data

            TraceLog(LOG_WARNING, "[TEX ID %i] Mipmaps [%i] generated manually on CPU side", texture->id, texture->mipmaps);
        }
        else TraceLog(LOG_WARNING, "[TEX ID %i] Mipmaps could not be generated for texture format", texture->id);
    }
#elif defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if ((texIsPOT) || (texNPOTSupported))
    {
        //glHint(GL_GENERATE_MIPMAP_HINT, GL_DONT_CARE);   // Hint for mipmaps generation algorythm: GL_FASTEST, GL_NICEST, GL_DONT_CARE
        glGenerateMipmap(GL_TEXTURE_2D);    // Generate mipmaps automatically
        TraceLog(LOG_INFO, "[TEX ID %i] Mipmaps generated automatically", texture->id);

        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR_MIPMAP_LINEAR);   // Activate Trilinear filtering for mipmaps

        #define MIN(a,b) (((a)<(b))?(a):(b))
        #define MAX(a,b) (((a)>(b))?(a):(b))

        texture->mipmaps =  1 + (int)floor(log(MAX(texture->width, texture->height))/log(2));
    }
#endif
    else TraceLog(LOG_WARNING, "[TEX ID %i] Mipmaps can not be generated", texture->id);

    glBindTexture(GL_TEXTURE_2D, 0);
}

// Upload vertex data into a VAO (if supported) and VBO
void rlLoadMesh(Mesh *mesh, bool dynamic)
{
    if (mesh->vaoId > 0)
    {
        // Check if mesh has already been loaded in GPU
        TraceLog(LOG_WARNING, "Trying to re-load an already loaded mesh");
        return;
    }

    mesh->vaoId = 0;        // Vertex Array Object
    mesh->vboId[0] = 0;     // Vertex positions VBO
    mesh->vboId[1] = 0;     // Vertex texcoords VBO
    mesh->vboId[2] = 0;     // Vertex normals VBO
    mesh->vboId[3] = 0;     // Vertex colors VBO
    mesh->vboId[4] = 0;     // Vertex tangents VBO
    mesh->vboId[5] = 0;     // Vertex texcoords2 VBO
    mesh->vboId[6] = 0;     // Vertex indices VBO

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    int drawHint = GL_STATIC_DRAW;
    if (dynamic) drawHint = GL_DYNAMIC_DRAW;

    if (vaoSupported)
    {
        // Initialize Quads VAO (Buffer A)
        glGenVertexArrays(1, &mesh->vaoId);
        glBindVertexArray(mesh->vaoId);
    }

    // NOTE: Attributes must be uploaded considering default locations points

    // Enable vertex attributes: position (shader-location = 0)
    glGenBuffers(1, &mesh->vboId[0]);
    glBindBuffer(GL_ARRAY_BUFFER, mesh->vboId[0]);
    glBufferData(GL_ARRAY_BUFFER, sizeof(float)*3*mesh->vertexCount, mesh->vertices, drawHint);
    glVertexAttribPointer(0, 3, GL_FLOAT, 0, 0, 0);
    glEnableVertexAttribArray(0);

    // Enable vertex attributes: texcoords (shader-location = 1)
    glGenBuffers(1, &mesh->vboId[1]);
    glBindBuffer(GL_ARRAY_BUFFER, mesh->vboId[1]);
    glBufferData(GL_ARRAY_BUFFER, sizeof(float)*2*mesh->vertexCount, mesh->texcoords, drawHint);
    glVertexAttribPointer(1, 2, GL_FLOAT, 0, 0, 0);
    glEnableVertexAttribArray(1);

    // Enable vertex attributes: normals (shader-location = 2)
    if (mesh->normals != NULL)
    {
        glGenBuffers(1, &mesh->vboId[2]);
        glBindBuffer(GL_ARRAY_BUFFER, mesh->vboId[2]);
        glBufferData(GL_ARRAY_BUFFER, sizeof(float)*3*mesh->vertexCount, mesh->normals, drawHint);
        glVertexAttribPointer(2, 3, GL_FLOAT, 0, 0, 0);
        glEnableVertexAttribArray(2);
    }
    else
    {
        // Default color vertex attribute set to WHITE
        glVertexAttrib3f(2, 1.0f, 1.0f, 1.0f);
        glDisableVertexAttribArray(2);
    }

    // Default color vertex attribute (shader-location = 3)
    if (mesh->colors != NULL)
    {
        glGenBuffers(1, &mesh->vboId[3]);
        glBindBuffer(GL_ARRAY_BUFFER, mesh->vboId[3]);
        glBufferData(GL_ARRAY_BUFFER, sizeof(unsigned char)*4*mesh->vertexCount, mesh->colors, drawHint);
        glVertexAttribPointer(3, 4, GL_UNSIGNED_BYTE, GL_TRUE, 0, 0);
        glEnableVertexAttribArray(3);
    }
    else
    {
        // Default color vertex attribute set to WHITE
        glVertexAttrib4f(3, 1.0f, 1.0f, 1.0f, 1.0f);
        glDisableVertexAttribArray(3);
    }

    // Default tangent vertex attribute (shader-location = 4)
    if (mesh->tangents != NULL)
    {
        glGenBuffers(1, &mesh->vboId[4]);
        glBindBuffer(GL_ARRAY_BUFFER, mesh->vboId[4]);
        glBufferData(GL_ARRAY_BUFFER, sizeof(float)*4*mesh->vertexCount, mesh->tangents, drawHint);
        glVertexAttribPointer(4, 4, GL_FLOAT, 0, 0, 0);
        glEnableVertexAttribArray(4);
    }
    else
    {
        // Default tangents vertex attribute
        glVertexAttrib4f(4, 0.0f, 0.0f, 0.0f, 0.0f);
        glDisableVertexAttribArray(4);
    }

    // Default texcoord2 vertex attribute (shader-location = 5)
    if (mesh->texcoords2 != NULL)
    {
        glGenBuffers(1, &mesh->vboId[5]);
        glBindBuffer(GL_ARRAY_BUFFER, mesh->vboId[5]);
        glBufferData(GL_ARRAY_BUFFER, sizeof(float)*2*mesh->vertexCount, mesh->texcoords2, drawHint);
        glVertexAttribPointer(5, 2, GL_FLOAT, 0, 0, 0);
        glEnableVertexAttribArray(5);
    }
    else
    {
        // Default texcoord2 vertex attribute
        glVertexAttrib2f(5, 0.0f, 0.0f);
        glDisableVertexAttribArray(5);
    }

    if (mesh->indices != NULL)
    {
        glGenBuffers(1, &mesh->vboId[6]);
        glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, mesh->vboId[6]);
        glBufferData(GL_ELEMENT_ARRAY_BUFFER, sizeof(unsigned short)*mesh->triangleCount*3, mesh->indices, drawHint);
    }

    if (vaoSupported)
    {
        if (mesh->vaoId > 0) TraceLog(LOG_INFO, "[VAO ID %i] Mesh uploaded successfully to VRAM (GPU)", mesh->vaoId);
        else TraceLog(LOG_WARNING, "Mesh could not be uploaded to VRAM (GPU)");
    }
    else
    {
        TraceLog(LOG_INFO, "[VBOs] Mesh uploaded successfully to VRAM (GPU)");
    }
#endif
}

// Load a new attributes buffer
unsigned int rlLoadAttribBuffer(unsigned int vaoId, int shaderLoc, void *buffer, int size, bool dynamic)
{
    unsigned int id = 0;

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    int drawHint = GL_STATIC_DRAW;
    if (dynamic) drawHint = GL_DYNAMIC_DRAW;

    if (vaoSupported) glBindVertexArray(vaoId);

    glGenBuffers(1, &id);
    glBindBuffer(GL_ARRAY_BUFFER, id);
    glBufferData(GL_ARRAY_BUFFER, size, buffer, drawHint);
    glVertexAttribPointer(shaderLoc, 2, GL_FLOAT, 0, 0, 0);
    glEnableVertexAttribArray(shaderLoc);

    if (vaoSupported) glBindVertexArray(0);
#endif

    return id;
}

// Update vertex or index data on GPU (upload new data to one buffer)
void rlUpdateMesh(Mesh mesh, int buffer, int num)
{
    rlUpdateMeshAt(mesh, buffer, num, 0);
}

// Update vertex or index data on GPU, at index
// WARNING: error checking is in place that will cause the data to not be
//          updated if offset + size exceeds what the buffer can hold
void rlUpdateMeshAt(Mesh mesh, int buffer, int num, int index)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    // Activate mesh VAO
    if (vaoSupported) glBindVertexArray(mesh.vaoId);

    switch (buffer)
    {
        case 0:     // Update vertices (vertex position)
        {
            glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[0]);
            if (index == 0 && num >= mesh.vertexCount) glBufferData(GL_ARRAY_BUFFER, sizeof(float)*3*num, mesh.vertices, GL_DYNAMIC_DRAW);
            else if (index + num >= mesh.vertexCount) break;
            else glBufferSubData(GL_ARRAY_BUFFER, sizeof(float)*3*index, sizeof(float)*3*num, mesh.vertices);

        } break;
        case 1:     // Update texcoords (vertex texture coordinates)
        {
            glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[1]);
            if (index == 0 && num >= mesh.vertexCount) glBufferData(GL_ARRAY_BUFFER, sizeof(float)*2*num, mesh.texcoords, GL_DYNAMIC_DRAW);
            else if (index + num >= mesh.vertexCount) break;
            else glBufferSubData(GL_ARRAY_BUFFER, sizeof(float)*2*index, sizeof(float)*2*num, mesh.texcoords);

        } break;
        case 2:     // Update normals (vertex normals)
        {
            glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[2]);
            if (index == 0 && num >= mesh.vertexCount) glBufferData(GL_ARRAY_BUFFER, sizeof(float)*3*num, mesh.normals, GL_DYNAMIC_DRAW);
            else if (index + num >= mesh.vertexCount) break;
            else glBufferSubData(GL_ARRAY_BUFFER, sizeof(float)*3*index, sizeof(float)*3*num, mesh.normals);

        } break;
        case 3:     // Update colors (vertex colors)
        {
            glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[3]);
            if (index == 0 && num >= mesh.vertexCount) glBufferData(GL_ARRAY_BUFFER, sizeof(float)*4*num, mesh.colors, GL_DYNAMIC_DRAW);
            else if (index + num >= mesh.vertexCount) break;
            else glBufferSubData(GL_ARRAY_BUFFER, sizeof(unsigned char)*4*index, sizeof(unsigned char)*4*num, mesh.colors);

        } break;
        case 4:     // Update tangents (vertex tangents)
        {
            glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[4]);
            if (index == 0 && num >= mesh.vertexCount) glBufferData(GL_ARRAY_BUFFER, sizeof(float)*4*num, mesh.tangents, GL_DYNAMIC_DRAW);
            else if (index + num >= mesh.vertexCount) break;
            else glBufferSubData(GL_ARRAY_BUFFER, sizeof(float)*4*index, sizeof(float)*4*num, mesh.tangents);
        } break;
        case 5:     // Update texcoords2 (vertex second texture coordinates)
        {
            glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[5]);
            if (index == 0 && num >= mesh.vertexCount) glBufferData(GL_ARRAY_BUFFER, sizeof(float)*2*num, mesh.texcoords2, GL_DYNAMIC_DRAW);
            else if (index + num >= mesh.vertexCount) break;
            else glBufferSubData(GL_ARRAY_BUFFER, sizeof(float)*2*index, sizeof(float)*2*num, mesh.texcoords2);
        } break;
        case 6:     // Update indices (triangle index buffer)
        {
            // the * 3 is because each triangle has 3 indices
            unsigned short *indices = mesh.indices;
            glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, mesh.vboId[6]);
            if (index == 0 && num >= mesh.triangleCount)
                glBufferData(GL_ELEMENT_ARRAY_BUFFER, sizeof(*indices)*num*3, indices, GL_DYNAMIC_DRAW);
            else if (index + num >= mesh.triangleCount)
                break;
            else
                glBufferSubData(GL_ELEMENT_ARRAY_BUFFER, sizeof(*indices)*index*3, sizeof(*indices)*num*3, indices);
        } break;
        default: break;
    }

    // Unbind the current VAO
    if (vaoSupported) glBindVertexArray(0);

    // Another option would be using buffer mapping...
    //mesh.vertices = glMapBuffer(GL_ARRAY_BUFFER, GL_READ_WRITE);
    // Now we can modify vertices
    //glUnmapBuffer(GL_ARRAY_BUFFER);
#endif
}

// Draw a 3d mesh with material and transform
void rlDrawMesh(Mesh mesh, Material material, Matrix transform)
{
#if defined(GRAPHICS_API_OPENGL_11)
    glEnable(GL_TEXTURE_2D);
    glBindTexture(GL_TEXTURE_2D, material.maps[MAP_DIFFUSE].texture.id);

    // NOTE: On OpenGL 1.1 we use Vertex Arrays to draw model
    glEnableClientState(GL_VERTEX_ARRAY);                   // Enable vertex array
    glEnableClientState(GL_TEXTURE_COORD_ARRAY);            // Enable texture coords array
    if (mesh.normals != NULL) glEnableClientState(GL_NORMAL_ARRAY);     // Enable normals array
    if (mesh.colors != NULL) glEnableClientState(GL_COLOR_ARRAY);       // Enable colors array

    glVertexPointer(3, GL_FLOAT, 0, mesh.vertices);         // Pointer to vertex coords array
    glTexCoordPointer(2, GL_FLOAT, 0, mesh.texcoords);      // Pointer to texture coords array
    if (mesh.normals != NULL) glNormalPointer(GL_FLOAT, 0, mesh.normals);           // Pointer to normals array
    if (mesh.colors != NULL) glColorPointer(4, GL_UNSIGNED_BYTE, 0, mesh.colors);   // Pointer to colors array

    rlPushMatrix();
        rlMultMatrixf(MatrixToFloat(transform));
        rlColor4ub(material.maps[MAP_DIFFUSE].color.r, material.maps[MAP_DIFFUSE].color.g, material.maps[MAP_DIFFUSE].color.b, material.maps[MAP_DIFFUSE].color.a);

        if (mesh.indices != NULL) glDrawElements(GL_TRIANGLES, mesh.triangleCount*3, GL_UNSIGNED_SHORT, mesh.indices);
        else glDrawArrays(GL_TRIANGLES, 0, mesh.vertexCount);
    rlPopMatrix();

    glDisableClientState(GL_VERTEX_ARRAY);                  // Disable vertex array
    glDisableClientState(GL_TEXTURE_COORD_ARRAY);           // Disable texture coords array
    if (mesh.normals != NULL) glDisableClientState(GL_NORMAL_ARRAY);    // Disable normals array
    if (mesh.colors != NULL) glDisableClientState(GL_NORMAL_ARRAY);     // Disable colors array

    glDisable(GL_TEXTURE_2D);
    glBindTexture(GL_TEXTURE_2D, 0);
#endif

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    // Bind shader program
    glUseProgram(material.shader.id);

    // Matrices and other values required by shader
    //-----------------------------------------------------
    // Calculate and send to shader model matrix (used by PBR shader)
    if (material.shader.locs[LOC_MATRIX_MODEL] != -1) SetShaderValueMatrix(material.shader, material.shader.locs[LOC_MATRIX_MODEL], transform);

    // Upload to shader material.colDiffuse
    if (material.shader.locs[LOC_COLOR_DIFFUSE] != -1)
        glUniform4f(material.shader.locs[LOC_COLOR_DIFFUSE], (float)material.maps[MAP_DIFFUSE].color.r/255.0f,
                                                           (float)material.maps[MAP_DIFFUSE].color.g/255.0f,
                                                           (float)material.maps[MAP_DIFFUSE].color.b/255.0f,
                                                           (float)material.maps[MAP_DIFFUSE].color.a/255.0f);

    // Upload to shader material.colSpecular (if available)
    if (material.shader.locs[LOC_COLOR_SPECULAR] != -1)
        glUniform4f(material.shader.locs[LOC_COLOR_SPECULAR], (float)material.maps[MAP_SPECULAR].color.r/255.0f,
                                                               (float)material.maps[MAP_SPECULAR].color.g/255.0f,
                                                               (float)material.maps[MAP_SPECULAR].color.b/255.0f,
                                                               (float)material.maps[MAP_SPECULAR].color.a/255.0f);

    if (material.shader.locs[LOC_MATRIX_VIEW] != -1) SetShaderValueMatrix(material.shader, material.shader.locs[LOC_MATRIX_VIEW], modelview);
    if (material.shader.locs[LOC_MATRIX_PROJECTION] != -1) SetShaderValueMatrix(material.shader, material.shader.locs[LOC_MATRIX_PROJECTION], projection);

    // At this point the modelview matrix just contains the view matrix (camera)
    // That's because BeginMode3D() sets it an no model-drawing function modifies it, all use rlPushMatrix() and rlPopMatrix()
    Matrix matView = modelview;         // View matrix (camera)
    Matrix matProjection = projection;  // Projection matrix (perspective)

    // TODO: Matrix nightmare! Trying to combine stack matrices with view matrix and local model transform matrix..
    // There is some problem in the order matrices are multiplied... it requires some time to figure out...
    Matrix matStackTransform = MatrixIdentity();

    // TODO: Consider possible transform matrices in the stack
    // Is this the right order? or should we start with the first stored matrix instead of the last one?
    //for (int i = stackCounter; i > 0; i--) matStackTransform = MatrixMultiply(stack[i], matStackTransform);

    Matrix matModel = MatrixMultiply(transform, matStackTransform); // Apply local model transformation
    Matrix matModelView = MatrixMultiply(matModel, matView);        // Transform to camera-space coordinates
    //-----------------------------------------------------

    // Bind active texture maps (if available)
    for (int i = 0; i < MAX_MATERIAL_MAPS; i++)
    {
        if (material.maps[i].texture.id > 0)
        {
            glActiveTexture(GL_TEXTURE0 + i);
            if ((i == MAP_IRRADIANCE) || (i == MAP_PREFILTER) || (i == MAP_CUBEMAP)) glBindTexture(GL_TEXTURE_CUBE_MAP, material.maps[i].texture.id);
            else glBindTexture(GL_TEXTURE_2D, material.maps[i].texture.id);

            glUniform1i(material.shader.locs[LOC_MAP_DIFFUSE + i], i);
        }
    }

    // Bind vertex array objects (or VBOs)
    if (vaoSupported) glBindVertexArray(mesh.vaoId);
    else
    {
        // Bind mesh VBO data: vertex position (shader-location = 0)
        glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[0]);
        glVertexAttribPointer(material.shader.locs[LOC_VERTEX_POSITION], 3, GL_FLOAT, 0, 0, 0);
        glEnableVertexAttribArray(material.shader.locs[LOC_VERTEX_POSITION]);

        // Bind mesh VBO data: vertex texcoords (shader-location = 1)
        glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[1]);
        glVertexAttribPointer(material.shader.locs[LOC_VERTEX_TEXCOORD01], 2, GL_FLOAT, 0, 0, 0);
        glEnableVertexAttribArray(material.shader.locs[LOC_VERTEX_TEXCOORD01]);

        // Bind mesh VBO data: vertex normals (shader-location = 2, if available)
        if (material.shader.locs[LOC_VERTEX_NORMAL] != -1)
        {
            glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[2]);
            glVertexAttribPointer(material.shader.locs[LOC_VERTEX_NORMAL], 3, GL_FLOAT, 0, 0, 0);
            glEnableVertexAttribArray(material.shader.locs[LOC_VERTEX_NORMAL]);
        }

        // Bind mesh VBO data: vertex colors (shader-location = 3, if available)
        if (material.shader.locs[LOC_VERTEX_COLOR] != -1)
        {
            if (mesh.vboId[3] != 0)
            {
                glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[3]);
                glVertexAttribPointer(material.shader.locs[LOC_VERTEX_COLOR], 4, GL_UNSIGNED_BYTE, GL_TRUE, 0, 0);
                glEnableVertexAttribArray(material.shader.locs[LOC_VERTEX_COLOR]);
            }
            else
            {
                // Set default value for unused attribute
                // NOTE: Required when using default shader and no VAO support
                glVertexAttrib4f(material.shader.locs[LOC_VERTEX_COLOR], 1.0f, 1.0f, 1.0f, 1.0f);
                glDisableVertexAttribArray(material.shader.locs[LOC_VERTEX_COLOR]);
            }
        }

        // Bind mesh VBO data: vertex tangents (shader-location = 4, if available)
        if (material.shader.locs[LOC_VERTEX_TANGENT] != -1)
        {
            glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[4]);
            glVertexAttribPointer(material.shader.locs[LOC_VERTEX_TANGENT], 4, GL_FLOAT, 0, 0, 0);
            glEnableVertexAttribArray(material.shader.locs[LOC_VERTEX_TANGENT]);
        }

        // Bind mesh VBO data: vertex texcoords2 (shader-location = 5, if available)
        if (material.shader.locs[LOC_VERTEX_TEXCOORD02] != -1)
        {
            glBindBuffer(GL_ARRAY_BUFFER, mesh.vboId[5]);
            glVertexAttribPointer(material.shader.locs[LOC_VERTEX_TEXCOORD02], 2, GL_FLOAT, 0, 0, 0);
            glEnableVertexAttribArray(material.shader.locs[LOC_VERTEX_TEXCOORD02]);
        }

        if (mesh.indices != NULL) glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, mesh.vboId[6]);
    }

    int eyesCount = 1;
#if defined(SUPPORT_VR_SIMULATOR)
    if (vrStereoRender) eyesCount = 2;
#endif

    for (int eye = 0; eye < eyesCount; eye++)
    {
        if (eyesCount == 1) modelview = matModelView;
        #if defined(SUPPORT_VR_SIMULATOR)
        else SetStereoView(eye, matProjection, matModelView);
        #endif

        // Calculate model-view-projection matrix (MVP)
        Matrix matMVP = MatrixMultiply(modelview, projection);        // Transform to screen-space coordinates

        // Send combined model-view-projection matrix to shader
        glUniformMatrix4fv(material.shader.locs[LOC_MATRIX_MVP], 1, false, MatrixToFloat(matMVP));

        // Draw call!
        if (mesh.indices != NULL) glDrawElements(GL_TRIANGLES, mesh.triangleCount*3, GL_UNSIGNED_SHORT, 0); // Indexed vertices draw
        else glDrawArrays(GL_TRIANGLES, 0, mesh.vertexCount);
    }

    // Unbind all binded texture maps
    for (int i = 0; i < MAX_MATERIAL_MAPS; i++)
    {
        glActiveTexture(GL_TEXTURE0 + i);       // Set shader active texture
        if ((i == MAP_IRRADIANCE) || (i == MAP_PREFILTER) || (i == MAP_CUBEMAP)) glBindTexture(GL_TEXTURE_CUBE_MAP, 0);
        else glBindTexture(GL_TEXTURE_2D, 0);   // Unbind current active texture
    }

    // Unind vertex array objects (or VBOs)
    if (vaoSupported) glBindVertexArray(0);
    else
    {
        glBindBuffer(GL_ARRAY_BUFFER, 0);
        if (mesh.indices != NULL) glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, 0);
    }

    // Unbind shader program
    glUseProgram(0);

    // Restore projection/modelview matrices
    // NOTE: In stereo rendering matrices are being modified to fit every eye
    projection = matProjection;
    modelview = matView;
#endif
}

// Unload mesh data from CPU and GPU
void rlUnloadMesh(Mesh mesh)
{
    RL_FREE(mesh.vertices);
    RL_FREE(mesh.texcoords);
    RL_FREE(mesh.normals);
    RL_FREE(mesh.colors);
    RL_FREE(mesh.tangents);
    RL_FREE(mesh.texcoords2);
    RL_FREE(mesh.indices);

    RL_FREE(mesh.animVertices);
    RL_FREE(mesh.animNormals);
    RL_FREE(mesh.boneWeights);
    RL_FREE(mesh.boneIds);

    rlDeleteBuffers(mesh.vboId[0]);   // vertex
    rlDeleteBuffers(mesh.vboId[1]);   // texcoords
    rlDeleteBuffers(mesh.vboId[2]);   // normals
    rlDeleteBuffers(mesh.vboId[3]);   // colors
    rlDeleteBuffers(mesh.vboId[4]);   // tangents
    rlDeleteBuffers(mesh.vboId[5]);   // texcoords2
    rlDeleteBuffers(mesh.vboId[6]);   // indices

    rlDeleteVertexArrays(mesh.vaoId);
}

// Read screen pixel data (color buffer)
unsigned char *rlReadScreenPixels(int width, int height)
{
    unsigned char *screenData = (unsigned char *)RL_CALLOC(width*height*4, sizeof(unsigned char));

    // NOTE 1: glReadPixels returns image flipped vertically -> (0,0) is the bottom left corner of the framebuffer
    // NOTE 2: We are getting alpha channel! Be careful, it can be transparent if not cleared properly!
    glReadPixels(0, 0, width, height, GL_RGBA, GL_UNSIGNED_BYTE, screenData);

    // Flip image vertically!
    unsigned char *imgData = (unsigned char *)RL_MALLOC(width*height*sizeof(unsigned char)*4);

    for (int y = height - 1; y >= 0; y--)
    {
        for (int x = 0; x < (width*4); x++)
        {
            imgData[((height - 1) - y)*width*4 + x] = screenData[(y*width*4) + x];  // Flip line

            // Set alpha component value to 255 (no trasparent image retrieval)
            // NOTE: Alpha value has already been applied to RGB in framebuffer, we don't need it!
            if (((x + 1)%4) == 0) imgData[((height - 1) - y)*width*4 + x] = 255;
        }
    }

    RL_FREE(screenData);

    return imgData;     // NOTE: image data should be freed
}

// Read texture pixel data
void *rlReadTexturePixels(Texture2D texture)
{
    void *pixels = NULL;

#if defined(GRAPHICS_API_OPENGL_11) || defined(GRAPHICS_API_OPENGL_33)
    glBindTexture(GL_TEXTURE_2D, texture.id);

    // NOTE: Using texture.id, we can retrieve some texture info (but not on OpenGL ES 2.0)
    // Possible texture info: GL_TEXTURE_RED_SIZE, GL_TEXTURE_GREEN_SIZE, GL_TEXTURE_BLUE_SIZE, GL_TEXTURE_ALPHA_SIZE
    //int width, height, format;
    //glGetTexLevelParameteriv(GL_TEXTURE_2D, 0, GL_TEXTURE_WIDTH, &width);
    //glGetTexLevelParameteriv(GL_TEXTURE_2D, 0, GL_TEXTURE_HEIGHT, &height);
    //glGetTexLevelParameteriv(GL_TEXTURE_2D, 0, GL_TEXTURE_INTERNAL_FORMAT, &format);

    // NOTE: Each row written to or read from by OpenGL pixel operations like glGetTexImage are aligned to a 4 byte boundary by default, which may add some padding.
    // Use glPixelStorei to modify padding with the GL_[UN]PACK_ALIGNMENT setting.
    // GL_PACK_ALIGNMENT affects operations that read from OpenGL memory (glReadPixels, glGetTexImage, etc.)
    // GL_UNPACK_ALIGNMENT affects operations that write to OpenGL memory (glTexImage, etc.)
    glPixelStorei(GL_PACK_ALIGNMENT, 1);

    unsigned int glInternalFormat, glFormat, glType;
    rlGetGlTextureFormats(texture.format, &glInternalFormat, &glFormat, &glType);
    unsigned int size = GetPixelDataSize(texture.width, texture.height, texture.format);

    if ((glInternalFormat != -1) && (texture.format < COMPRESSED_DXT1_RGB))
    {
        pixels = (unsigned char *)RL_MALLOC(size);
        glGetTexImage(GL_TEXTURE_2D, 0, glFormat, glType, pixels);
    }
    else TraceLog(LOG_WARNING, "Texture data retrieval not suported for pixel format");

    glBindTexture(GL_TEXTURE_2D, 0);
#endif

#if defined(GRAPHICS_API_OPENGL_ES2)
    // glGetTexImage() is not available on OpenGL ES 2.0
    // Texture2D width and height are required on OpenGL ES 2.0. There is no way to get it from texture id.
    // Two possible Options:
    // 1 - Bind texture to color fbo attachment and glReadPixels()
    // 2 - Create an fbo, activate it, render quad with texture, glReadPixels()
    // We are using Option 1, just need to care for texture format on retrieval
    // NOTE: This behaviour could be conditioned by graphic driver...
    RenderTexture2D fbo = rlLoadRenderTexture(texture.width, texture.height, UNCOMPRESSED_R8G8B8A8, 16, false);

    glBindFramebuffer(GL_FRAMEBUFFER, fbo.id);
    glBindTexture(GL_TEXTURE_2D, 0);

    // Attach our texture to FBO
    // NOTE: Previoust attached texture is automatically detached
    glFramebufferTexture2D(GL_FRAMEBUFFER, GL_COLOR_ATTACHMENT0, GL_TEXTURE_2D, texture.id, 0);

    // Allocate enough memory to read back our texture data
    pixels = (unsigned char *)RL_MALLOC(GetPixelDataSize(texture.width, texture.height, texture.format));

    // Get OpenGL internal formats and data type from our texture format
    unsigned int glInternalFormat, glFormat, glType;
    rlGetGlTextureFormats(texture.format, &glInternalFormat, &glFormat, &glType);

    // NOTE: We read data as RGBA because FBO texture is configured as RGBA, despite binding a RGB texture...
    glReadPixels(0, 0, texture.width, texture.height, glFormat, glType, pixels);

    // Re-attach internal FBO color texture before deleting it
    glFramebufferTexture2D(GL_FRAMEBUFFER, GL_COLOR_ATTACHMENT0, GL_TEXTURE_2D, fbo.texture.id, 0);

    glBindFramebuffer(GL_FRAMEBUFFER, 0);

    // Clean up temporal fbo
    rlDeleteRenderTextures(fbo);
#endif

    return pixels;
}

//----------------------------------------------------------------------------------
// Module Functions Definition - Shaders Functions
// NOTE: Those functions are exposed directly to the user in raylib.h
//----------------------------------------------------------------------------------

// Get default internal texture (white texture)
Texture2D GetTextureDefault(void)
{
    Texture2D texture = { 0 };
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    texture.id = defaultTextureId;
    texture.width = 1;
    texture.height = 1;
    texture.mipmaps = 1;
    texture.format = UNCOMPRESSED_R8G8B8A8;
#endif
    return texture;
}

// Get default shader
Shader GetShaderDefault(void)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    return defaultShader;
#else
    Shader shader = { 0 };
    return shader;
#endif
}

// Load text data from file
// NOTE: text chars array should be freed manually
char *LoadText(const char *fileName)
{
    FILE *textFile = NULL;
    char *text = NULL;

    if (fileName != NULL)
    {
        textFile = fopen(fileName,"rt");

        if (textFile != NULL)
        {
            fseek(textFile, 0, SEEK_END);
            int size = ftell(textFile);
            fseek(textFile, 0, SEEK_SET);

            if (size > 0)
            {
                text = (char *)RL_MALLOC(sizeof(char)*(size + 1));
                int count = fread(text, sizeof(char), size, textFile);
                text[count] = '\0';
            }

            fclose(textFile);
        }
        else TraceLog(LOG_WARNING, "[%s] Text file could not be opened", fileName);
    }

    return text;
}

// Load shader from files and bind default locations
// NOTE: If shader string is NULL, using default vertex/fragment shaders
Shader LoadShader(const char *vsFileName, const char *fsFileName)
{
    Shader shader = { 0 };

    // NOTE: Shader.locs is allocated by LoadShaderCode()

    char *vShaderStr = NULL;
    char *fShaderStr = NULL;

    if (vsFileName != NULL) vShaderStr = LoadText(vsFileName);
    if (fsFileName != NULL) fShaderStr = LoadText(fsFileName);

    shader = LoadShaderCode(vShaderStr, fShaderStr);

    if (vShaderStr != NULL) RL_FREE(vShaderStr);
    if (fShaderStr != NULL) RL_FREE(fShaderStr);

    return shader;
}

// Load shader from code strings
// NOTE: If shader string is NULL, using default vertex/fragment shaders
Shader LoadShaderCode(const char *vsCode, const char *fsCode)
{
    Shader shader = { 0 };
    shader.locs = (int *)RL_CALLOC(MAX_SHADER_LOCATIONS, sizeof(int));

    // NOTE: All locations must be reseted to -1 (no location)
    for (int i = 0; i < MAX_SHADER_LOCATIONS; i++) shader.locs[i] = -1;

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    unsigned int vertexShaderId = defaultVShaderId;
    unsigned int fragmentShaderId = defaultFShaderId;

    if (vsCode != NULL) vertexShaderId = CompileShader(vsCode, GL_VERTEX_SHADER);
    if (fsCode != NULL) fragmentShaderId = CompileShader(fsCode, GL_FRAGMENT_SHADER);

    if ((vertexShaderId == defaultVShaderId) && (fragmentShaderId == defaultFShaderId)) shader = defaultShader;
    else
    {
        shader.id = LoadShaderProgram(vertexShaderId, fragmentShaderId);

        if (vertexShaderId != defaultVShaderId) glDeleteShader(vertexShaderId);
        if (fragmentShaderId != defaultFShaderId) glDeleteShader(fragmentShaderId);

        if (shader.id == 0)
        {
            TraceLog(LOG_WARNING, "Custom shader could not be loaded");
            shader = defaultShader;
        }

        // After shader loading, we TRY to set default location names
        if (shader.id > 0) SetShaderDefaultLocations(&shader);
    }

    // Get available shader uniforms
    // NOTE: This information is useful for debug...
    int uniformCount = -1;

    glGetProgramiv(shader.id, GL_ACTIVE_UNIFORMS, &uniformCount);

    for (int i = 0; i < uniformCount; i++)
    {
        int namelen = -1;
        int num = -1;
        char name[256]; // Assume no variable names longer than 256
        GLenum type = GL_ZERO;

        // Get the name of the uniforms
        glGetActiveUniform(shader.id, i,sizeof(name) - 1, &namelen, &num, &type, name);

        name[namelen] = 0;

        // Get the location of the named uniform
        unsigned int location = glGetUniformLocation(shader.id, name);

        TraceLog(LOG_DEBUG, "[SHDR ID %i] Active uniform [%s] set at location: %i", shader.id, name, location);
    }
#endif

    return shader;
}

// Unload shader from GPU memory (VRAM)
void UnloadShader(Shader shader)
{
    if (shader.id > 0)
    {
        rlDeleteShader(shader.id);
        TraceLog(LOG_INFO, "[SHDR ID %i] Unloaded shader program data", shader.id);
    }

    RL_FREE(shader.locs);
}

// Begin custom shader mode
void BeginShaderMode(Shader shader)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if (currentShader.id != shader.id)
    {
        rlglDraw();
        currentShader = shader;
    }
#endif
}

// End custom shader mode (returns to default shader)
void EndShaderMode(void)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    BeginShaderMode(defaultShader);
#endif
}

// Get shader uniform location
int GetShaderLocation(Shader shader, const char *uniformName)
{
    int location = -1;
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    location = glGetUniformLocation(shader.id, uniformName);

    if (location == -1) TraceLog(LOG_WARNING, "[SHDR ID %i][%s] Shader uniform could not be found", shader.id, uniformName);
    else TraceLog(LOG_INFO, "[SHDR ID %i][%s] Shader uniform set at location: %i", shader.id, uniformName, location);
#endif
    return location;
}

// Set shader uniform value
void SetShaderValue(Shader shader, int uniformLoc, const void *value, int uniformType)
{
    SetShaderValueV(shader, uniformLoc, value, uniformType, 1);
}

// Set shader uniform value vector
void SetShaderValueV(Shader shader, int uniformLoc, const void *value, int uniformType, int count)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    glUseProgram(shader.id);

    switch (uniformType)
    {
        case UNIFORM_FLOAT: glUniform1fv(uniformLoc, count, (float *)value); break;
        case UNIFORM_VEC2: glUniform2fv(uniformLoc, count, (float *)value); break;
        case UNIFORM_VEC3: glUniform3fv(uniformLoc, count, (float *)value); break;
        case UNIFORM_VEC4: glUniform4fv(uniformLoc, count, (float *)value); break;
        case UNIFORM_INT: glUniform1iv(uniformLoc, count, (int *)value); break;
        case UNIFORM_IVEC2: glUniform2iv(uniformLoc, count, (int *)value); break;
        case UNIFORM_IVEC3: glUniform3iv(uniformLoc, count, (int *)value); break;
        case UNIFORM_IVEC4: glUniform4iv(uniformLoc, count, (int *)value); break;
        case UNIFORM_SAMPLER2D: glUniform1iv(uniformLoc, count, (int *)value); break;
        default: TraceLog(LOG_WARNING, "Shader uniform could not be set data type not recognized");
    }

    //glUseProgram(0);      // Avoid reseting current shader program, in case other uniforms are set
#endif
}


// Set shader uniform value (matrix 4x4)
void SetShaderValueMatrix(Shader shader, int uniformLoc, Matrix mat)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    glUseProgram(shader.id);

    glUniformMatrix4fv(uniformLoc, 1, false, MatrixToFloat(mat));

    //glUseProgram(0);
#endif
}

// Set shader uniform value for texture
void SetShaderValueTexture(Shader shader, int uniformLoc, Texture2D texture)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    glUseProgram(shader.id);

    glUniform1i(uniformLoc, texture.id);

    //glUseProgram(0);
#endif
}

// Set a custom projection matrix (replaces internal projection matrix)
void SetMatrixProjection(Matrix proj)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    projection = proj;
#endif
}

// Return internal projection matrix
Matrix GetMatrixProjection(void) {
#if defined(GRAPHICS_API_OPENGL_11)
    float mat[16];
    glGetFloatv(GL_PROJECTION_MATRIX,mat);
    Matrix m;
    m.m0  = mat[0];     m.m1  = mat[1];     m.m2  = mat[2];     m.m3  = mat[3];
    m.m4  = mat[4];     m.m5  = mat[5];     m.m6  = mat[6];     m.m7  = mat[7];
    m.m8  = mat[8];     m.m9  = mat[9];     m.m10 = mat[10];    m.m11 = mat[11];
    m.m12 = mat[12];    m.m13 = mat[13];    m.m14 = mat[14];    m.m15 = mat[15];
    return m;
#else
    return projection;
#endif
#
}

// Set a custom modelview matrix (replaces internal modelview matrix)
void SetMatrixModelview(Matrix view)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    modelview = view;
#endif
}

// Return internal modelview matrix
Matrix GetMatrixModelview(void)
{
    Matrix matrix = MatrixIdentity();
#if defined(GRAPHICS_API_OPENGL_11)
    float mat[16];
    glGetFloatv(GL_MODELVIEW_MATRIX, mat);
    matrix.m0  = mat[0];     matrix.m1  = mat[1];     matrix.m2  = mat[2];     matrix.m3  = mat[3];
    matrix.m4  = mat[4];     matrix.m5  = mat[5];     matrix.m6  = mat[6];     matrix.m7  = mat[7];
    matrix.m8  = mat[8];     matrix.m9  = mat[9];     matrix.m10 = mat[10];    matrix.m11 = mat[11];
    matrix.m12 = mat[12];    matrix.m13 = mat[13];    matrix.m14 = mat[14];    matrix.m15 = mat[15];
#else
    matrix = modelview;
#endif
    return matrix;
}

// Generate cubemap texture from HDR texture
// TODO: OpenGL ES 2.0 does not support GL_RGB16F texture format, neither GL_DEPTH_COMPONENT24
Texture2D GenTextureCubemap(Shader shader, Texture2D skyHDR, int size)
{
    Texture2D cubemap = { 0 };
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    // NOTE: SetShaderDefaultLocations() already setups locations for projection and view Matrix in shader
    // Other locations should be setup externally in shader before calling the function

    // Set up depth face culling and cubemap seamless
    glDisable(GL_CULL_FACE);
#if defined(GRAPHICS_API_OPENGL_33)
    glEnable(GL_TEXTURE_CUBE_MAP_SEAMLESS);     // Flag not supported on OpenGL ES 2.0
#endif

    // Setup framebuffer
    unsigned int fbo, rbo;
    glGenFramebuffers(1, &fbo);
    glGenRenderbuffers(1, &rbo);
    glBindFramebuffer(GL_FRAMEBUFFER, fbo);
    glBindRenderbuffer(GL_RENDERBUFFER, rbo);
#if defined(GRAPHICS_API_OPENGL_33)
    glRenderbufferStorage(GL_RENDERBUFFER, GL_DEPTH_COMPONENT24, size, size);
#elif defined(GRAPHICS_API_OPENGL_ES2)
    glRenderbufferStorage(GL_RENDERBUFFER, GL_DEPTH_COMPONENT16, size, size);
#endif
    glFramebufferRenderbuffer(GL_FRAMEBUFFER, GL_DEPTH_ATTACHMENT, GL_RENDERBUFFER, rbo);

    // Set up cubemap to render and attach to framebuffer
    // NOTE: Faces are stored as 32 bit floating point values
    glGenTextures(1, &cubemap.id);
    glBindTexture(GL_TEXTURE_CUBE_MAP, cubemap.id);
    for (unsigned int i = 0; i < 6; i++)
    {
#if defined(GRAPHICS_API_OPENGL_33)
        glTexImage2D(GL_TEXTURE_CUBE_MAP_POSITIVE_X + i, 0, GL_RGB32F, size, size, 0, GL_RGB, GL_FLOAT, NULL);
#elif defined(GRAPHICS_API_OPENGL_ES2)
        if (texFloatSupported) glTexImage2D(GL_TEXTURE_CUBE_MAP_POSITIVE_X + i, 0, GL_RGB, size, size, 0, GL_RGB, GL_FLOAT, NULL);
#endif
    }

    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE);
#if defined(GRAPHICS_API_OPENGL_33)
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_R, GL_CLAMP_TO_EDGE);  // Flag not supported on OpenGL ES 2.0
#endif
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_MAG_FILTER, GL_LINEAR);

    // Create projection and different views for each face
    Matrix fboProjection = MatrixPerspective(90.0*DEG2RAD, 1.0, 0.01, 1000.0);
    Matrix fboViews[6] = {
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 1.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ -1.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, 1.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, 1.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, -1.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, 1.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, -1.0f }, (Vector3){ 0.0f, -1.0f, 0.0f })
    };

    // Convert HDR equirectangular environment map to cubemap equivalent
    glUseProgram(shader.id);
    glActiveTexture(GL_TEXTURE0);
    glBindTexture(GL_TEXTURE_2D, skyHDR.id);
    SetShaderValueMatrix(shader, shader.locs[LOC_MATRIX_PROJECTION], fboProjection);

    // Note: don't forget to configure the viewport to the capture dimensions
    glViewport(0, 0, size, size);
    glBindFramebuffer(GL_FRAMEBUFFER, fbo);

    for (int i = 0; i < 6; i++)
    {
        SetShaderValueMatrix(shader, shader.locs[LOC_MATRIX_VIEW], fboViews[i]);
        glFramebufferTexture2D(GL_FRAMEBUFFER, GL_COLOR_ATTACHMENT0, GL_TEXTURE_CUBE_MAP_POSITIVE_X + i, cubemap.id, 0);
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
        GenDrawCube();
    }

    // Unbind framebuffer and textures
    glBindFramebuffer(GL_FRAMEBUFFER, 0);

    // Reset viewport dimensions to default
    glViewport(0, 0, framebufferWidth, framebufferHeight);
    //glEnable(GL_CULL_FACE);

    // NOTE: Texture2D is a GL_TEXTURE_CUBE_MAP, not a GL_TEXTURE_2D!
    cubemap.width = size;
    cubemap.height = size;
    cubemap.mipmaps = 1;
    cubemap.format = UNCOMPRESSED_R32G32B32;
#endif
    return cubemap;
}

// Generate irradiance texture using cubemap data
// TODO: OpenGL ES 2.0 does not support GL_RGB16F texture format, neither GL_DEPTH_COMPONENT24
Texture2D GenTextureIrradiance(Shader shader, Texture2D cubemap, int size)
{
    Texture2D irradiance = { 0 };

#if defined(GRAPHICS_API_OPENGL_33) // || defined(GRAPHICS_API_OPENGL_ES2)
    // NOTE: SetShaderDefaultLocations() already setups locations for projection and view Matrix in shader
    // Other locations should be setup externally in shader before calling the function

    // Setup framebuffer
    unsigned int fbo, rbo;
    glGenFramebuffers(1, &fbo);
    glGenRenderbuffers(1, &rbo);
    glBindFramebuffer(GL_FRAMEBUFFER, fbo);
    glBindRenderbuffer(GL_RENDERBUFFER, rbo);
    glRenderbufferStorage(GL_RENDERBUFFER, GL_DEPTH_COMPONENT24, size, size);
    glFramebufferRenderbuffer(GL_FRAMEBUFFER, GL_DEPTH_ATTACHMENT, GL_RENDERBUFFER, rbo);

    // Create an irradiance cubemap, and re-scale capture FBO to irradiance scale
    glGenTextures(1, &irradiance.id);
    glBindTexture(GL_TEXTURE_CUBE_MAP, irradiance.id);
    for (unsigned int i = 0; i < 6; i++)
    {
        glTexImage2D(GL_TEXTURE_CUBE_MAP_POSITIVE_X + i, 0, GL_RGB16F, size, size, 0, GL_RGB, GL_FLOAT, NULL);
    }

    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_R, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_MAG_FILTER, GL_LINEAR);

    // Create projection (transposed) and different views for each face
    Matrix fboProjection = MatrixPerspective(90.0*DEG2RAD, 1.0, 0.01, 1000.0);
    Matrix fboViews[6] = {
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 1.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ -1.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, 1.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, 1.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, -1.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, 1.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, -1.0f }, (Vector3){ 0.0f, -1.0f, 0.0f })
    };

    // Solve diffuse integral by convolution to create an irradiance cubemap
    glUseProgram(shader.id);
    glActiveTexture(GL_TEXTURE0);
    glBindTexture(GL_TEXTURE_CUBE_MAP, cubemap.id);
    SetShaderValueMatrix(shader, shader.locs[LOC_MATRIX_PROJECTION], fboProjection);

    // Note: don't forget to configure the viewport to the capture dimensions
    glViewport(0, 0, size, size);
    glBindFramebuffer(GL_FRAMEBUFFER, fbo);

    for (int i = 0; i < 6; i++)
    {
        SetShaderValueMatrix(shader, shader.locs[LOC_MATRIX_VIEW], fboViews[i]);
        glFramebufferTexture2D(GL_FRAMEBUFFER, GL_COLOR_ATTACHMENT0, GL_TEXTURE_CUBE_MAP_POSITIVE_X + i, irradiance.id, 0);
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
        GenDrawCube();
    }

    // Unbind framebuffer and textures
    glBindFramebuffer(GL_FRAMEBUFFER, 0);

    // Reset viewport dimensions to default
    glViewport(0, 0, framebufferWidth, framebufferHeight);

    irradiance.width = size;
    irradiance.height = size;
    irradiance.mipmaps = 1;
    //irradiance.format = UNCOMPRESSED_R16G16B16;
#endif
    return irradiance;
}

// Generate prefilter texture using cubemap data
// TODO: OpenGL ES 2.0 does not support GL_RGB16F texture format, neither GL_DEPTH_COMPONENT24
Texture2D GenTexturePrefilter(Shader shader, Texture2D cubemap, int size)
{
    Texture2D prefilter = { 0 };

#if defined(GRAPHICS_API_OPENGL_33) // || defined(GRAPHICS_API_OPENGL_ES2)
    // NOTE: SetShaderDefaultLocations() already setups locations for projection and view Matrix in shader
    // Other locations should be setup externally in shader before calling the function
    // TODO: Locations should be taken out of this function... too shader dependant...
    int roughnessLoc = GetShaderLocation(shader, "roughness");

    // Setup framebuffer
    unsigned int fbo, rbo;
    glGenFramebuffers(1, &fbo);
    glGenRenderbuffers(1, &rbo);
    glBindFramebuffer(GL_FRAMEBUFFER, fbo);
    glBindRenderbuffer(GL_RENDERBUFFER, rbo);
    glRenderbufferStorage(GL_RENDERBUFFER, GL_DEPTH_COMPONENT24, size, size);
    glFramebufferRenderbuffer(GL_FRAMEBUFFER, GL_DEPTH_ATTACHMENT, GL_RENDERBUFFER, rbo);

    // Create a prefiltered HDR environment map
    glGenTextures(1, &prefilter.id);
    glBindTexture(GL_TEXTURE_CUBE_MAP, prefilter.id);
    for (unsigned int i = 0; i < 6; i++)
    {
        glTexImage2D(GL_TEXTURE_CUBE_MAP_POSITIVE_X + i, 0, GL_RGB16F, size, size, 0, GL_RGB, GL_FLOAT, NULL);
    }

    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_WRAP_R, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_MIN_FILTER, GL_LINEAR_MIPMAP_LINEAR);
    glTexParameteri(GL_TEXTURE_CUBE_MAP, GL_TEXTURE_MAG_FILTER, GL_LINEAR);

    // Generate mipmaps for the prefiltered HDR texture
    glGenerateMipmap(GL_TEXTURE_CUBE_MAP);

    // Create projection (transposed) and different views for each face
    Matrix fboProjection = MatrixPerspective(90.0*DEG2RAD, 1.0, 0.01, 1000.0);
    Matrix fboViews[6] = {
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 1.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ -1.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, 1.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, 1.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, -1.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, 1.0f }, (Vector3){ 0.0f, -1.0f, 0.0f }),
        MatrixLookAt((Vector3){ 0.0f, 0.0f, 0.0f }, (Vector3){ 0.0f, 0.0f, -1.0f }, (Vector3){ 0.0f, -1.0f, 0.0f })
    };

    // Prefilter HDR and store data into mipmap levels
    glUseProgram(shader.id);
    glActiveTexture(GL_TEXTURE0);
    glBindTexture(GL_TEXTURE_CUBE_MAP, cubemap.id);
    SetShaderValueMatrix(shader, shader.locs[LOC_MATRIX_PROJECTION], fboProjection);

    glBindFramebuffer(GL_FRAMEBUFFER, fbo);

    #define MAX_MIPMAP_LEVELS   5   // Max number of prefilter texture mipmaps

    for (int mip = 0; mip < MAX_MIPMAP_LEVELS; mip++)
    {
        // Resize framebuffer according to mip-level size.
        unsigned int mipWidth  = size*(int)powf(0.5f, (float)mip);
        unsigned int mipHeight = size*(int)powf(0.5f, (float)mip);

        glBindRenderbuffer(GL_RENDERBUFFER, rbo);
        glRenderbufferStorage(GL_RENDERBUFFER, GL_DEPTH_COMPONENT24, mipWidth, mipHeight);
        glViewport(0, 0, mipWidth, mipHeight);

        float roughness = (float)mip/(float)(MAX_MIPMAP_LEVELS - 1);
        glUniform1f(roughnessLoc, roughness);

        for (int i = 0; i < 6; i++)
        {
            SetShaderValueMatrix(shader, shader.locs[LOC_MATRIX_VIEW], fboViews[i]);
            glFramebufferTexture2D(GL_FRAMEBUFFER, GL_COLOR_ATTACHMENT0, GL_TEXTURE_CUBE_MAP_POSITIVE_X + i, prefilter.id, mip);
            glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
            GenDrawCube();
        }
    }

    // Unbind framebuffer and textures
    glBindFramebuffer(GL_FRAMEBUFFER, 0);

    // Reset viewport dimensions to default
    glViewport(0, 0, framebufferWidth, framebufferHeight);

    prefilter.width = size;
    prefilter.height = size;
    //prefilter.mipmaps = 1 + (int)floor(log(size)/log(2));
    //prefilter.format = UNCOMPRESSED_R16G16B16;
#endif
    return prefilter;
}

// Generate BRDF texture using cubemap data
// NOTE: OpenGL ES 2.0 does not support GL_RGB16F texture format, neither GL_DEPTH_COMPONENT24
// TODO: Review implementation: https://github.com/HectorMF/BRDFGenerator
Texture2D GenTextureBRDF(Shader shader, int size)
{
    Texture2D brdf = { 0 };
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    // Generate BRDF convolution texture
    glGenTextures(1, &brdf.id);
    glBindTexture(GL_TEXTURE_2D, brdf.id);
#if defined(GRAPHICS_API_OPENGL_33)
    glTexImage2D(GL_TEXTURE_2D, 0, GL_RGB32F, size, size, 0, GL_RGB, GL_FLOAT, NULL);
#elif defined(GRAPHICS_API_OPENGL_ES2)
    if (texFloatSupported) glTexImage2D(GL_TEXTURE_2D, 0, GL_RGB, size, size, 0, GL_RGB, GL_FLOAT, NULL);
#endif

    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);

    // Render BRDF LUT into a quad using FBO
    unsigned int fbo, rbo;
    glGenFramebuffers(1, &fbo);
    glGenRenderbuffers(1, &rbo);
    glBindFramebuffer(GL_FRAMEBUFFER, fbo);
    glBindRenderbuffer(GL_RENDERBUFFER, rbo);
#if defined(GRAPHICS_API_OPENGL_33)
    glRenderbufferStorage(GL_RENDERBUFFER, GL_DEPTH_COMPONENT24, size, size);
#elif defined(GRAPHICS_API_OPENGL_ES2)
    glRenderbufferStorage(GL_RENDERBUFFER, GL_DEPTH_COMPONENT16, size, size);
#endif
    glFramebufferTexture2D(GL_FRAMEBUFFER, GL_COLOR_ATTACHMENT0, GL_TEXTURE_2D, brdf.id, 0);

    glViewport(0, 0, size, size);
    glUseProgram(shader.id);
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
    GenDrawQuad();

    // Unbind framebuffer and textures
    glBindFramebuffer(GL_FRAMEBUFFER, 0);

    // Unload framebuffer but keep color texture
    glDeleteRenderbuffers(1, &rbo);
    glDeleteFramebuffers(1, &fbo);

    // Reset viewport dimensions to default
    glViewport(0, 0, framebufferWidth, framebufferHeight);

    brdf.width = size;
    brdf.height = size;
    brdf.mipmaps = 1;
    brdf.format = UNCOMPRESSED_R32G32B32;
#endif
    return brdf;
}

// Begin blending mode (alpha, additive, multiplied)
// NOTE: Only 3 blending modes supported, default blend mode is alpha
void BeginBlendMode(int mode)
{
    if ((blendMode != mode) && (mode < 3))
    {
        rlglDraw();

        switch (mode)
        {
            case BLEND_ALPHA: glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA); break;
            case BLEND_ADDITIVE: glBlendFunc(GL_SRC_ALPHA, GL_ONE); break; // Alternative: glBlendFunc(GL_ONE, GL_ONE);
            case BLEND_MULTIPLIED: glBlendFunc(GL_DST_COLOR, GL_ONE_MINUS_SRC_ALPHA); break;
            default: break;
        }

        blendMode = mode;
    }
}

// End blending mode (reset to default: alpha blending)
void EndBlendMode(void)
{
    BeginBlendMode(BLEND_ALPHA);
}

#if defined(SUPPORT_VR_SIMULATOR)
// Init VR simulator for selected device parameters
// NOTE: It modifies the global variable: stereoFbo
void InitVrSimulator(void)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    // Initialize framebuffer and textures for stereo rendering
    // NOTE: Screen size should match HMD aspect ratio
    stereoFbo = rlLoadRenderTexture(framebufferWidth, framebufferHeight, UNCOMPRESSED_R8G8B8A8, 24, false);

    vrSimulatorReady = true;
#else
    TraceLog(LOG_WARNING, "VR Simulator not supported on OpenGL 1.1");
#endif
}

// Update VR tracking (position and orientation) and camera
// NOTE: Camera (position, target, up) gets update with head tracking information
void UpdateVrTracking(Camera *camera)
{
    // TODO: Simulate 1st person camera system
}

// Close VR simulator for current device
void CloseVrSimulator(void)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if (vrSimulatorReady) rlDeleteRenderTextures(stereoFbo);        // Unload stereo framebuffer and texture
#endif
}

// Set stereo rendering configuration parameters
void SetVrConfiguration(VrDeviceInfo hmd, Shader distortion)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    // Reset vrConfig for a new values assignment
    memset(&vrConfig, 0, sizeof(vrConfig));

    // Assign distortion shader
    vrConfig.distortionShader = distortion;

    // Compute aspect ratio
    float aspect = ((float)hmd.hResolution*0.5f)/(float)hmd.vResolution;

    // Compute lens parameters
    float lensShift = (hmd.hScreenSize*0.25f - hmd.lensSeparationDistance*0.5f)/hmd.hScreenSize;
    float leftLensCenter[2] = { 0.25f + lensShift, 0.5f };
    float rightLensCenter[2] = { 0.75f - lensShift, 0.5f };
    float leftScreenCenter[2] = { 0.25f, 0.5f };
    float rightScreenCenter[2] = { 0.75f, 0.5f };

    // Compute distortion scale parameters
    // NOTE: To get lens max radius, lensShift must be normalized to [-1..1]
    float lensRadius = (float)fabs(-1.0f - 4.0f*lensShift);
    float lensRadiusSq = lensRadius*lensRadius;
    float distortionScale = hmd.lensDistortionValues[0] +
                            hmd.lensDistortionValues[1]*lensRadiusSq +
                            hmd.lensDistortionValues[2]*lensRadiusSq*lensRadiusSq +
                            hmd.lensDistortionValues[3]*lensRadiusSq*lensRadiusSq*lensRadiusSq;

    TraceLog(LOG_DEBUG, "VR: Distortion Scale: %f", distortionScale);

    float normScreenWidth = 0.5f;
    float normScreenHeight = 1.0f;
    float scaleIn[2] = { 2.0f/normScreenWidth, 2.0f/normScreenHeight/aspect };
    float scale[2] = { normScreenWidth*0.5f/distortionScale, normScreenHeight*0.5f*aspect/distortionScale };

    TraceLog(LOG_DEBUG, "VR: Distortion Shader: LeftLensCenter = { %f, %f }", leftLensCenter[0], leftLensCenter[1]);
    TraceLog(LOG_DEBUG, "VR: Distortion Shader: RightLensCenter = { %f, %f }", rightLensCenter[0], rightLensCenter[1]);
    TraceLog(LOG_DEBUG, "VR: Distortion Shader: Scale = { %f, %f }", scale[0], scale[1]);
    TraceLog(LOG_DEBUG, "VR: Distortion Shader: ScaleIn = { %f, %f }", scaleIn[0], scaleIn[1]);

    // Fovy is normally computed with: 2*atan2(hmd.vScreenSize, 2*hmd.eyeToScreenDistance)
    // ...but with lens distortion it is increased (see Oculus SDK Documentation)
    //float fovy = 2.0f*atan2(hmd.vScreenSize*0.5f*distortionScale, hmd.eyeToScreenDistance);     // Really need distortionScale?
    float fovy = 2.0f*(float)atan2(hmd.vScreenSize*0.5f, hmd.eyeToScreenDistance);

    // Compute camera projection matrices
    float projOffset = 4.0f*lensShift;      // Scaled to projection space coordinates [-1..1]
    Matrix proj = MatrixPerspective(fovy, aspect, 0.01, 1000.0);
    vrConfig.eyesProjection[0] = MatrixMultiply(proj, MatrixTranslate(projOffset, 0.0f, 0.0f));
    vrConfig.eyesProjection[1] = MatrixMultiply(proj, MatrixTranslate(-projOffset, 0.0f, 0.0f));

    // Compute camera transformation matrices
    // NOTE: Camera movement might seem more natural if we model the head.
    // Our axis of rotation is the base of our head, so we might want to add
    // some y (base of head to eye level) and -z (center of head to eye protrusion) to the camera positions.
    vrConfig.eyesViewOffset[0] = MatrixTranslate(-hmd.interpupillaryDistance*0.5f, 0.075f, 0.045f);
    vrConfig.eyesViewOffset[1] = MatrixTranslate(hmd.interpupillaryDistance*0.5f, 0.075f, 0.045f);

    // Compute eyes Viewports
    vrConfig.eyeViewportRight[2] = hmd.hResolution/2;
    vrConfig.eyeViewportRight[3] = hmd.vResolution;

    vrConfig.eyeViewportLeft[0] = hmd.hResolution/2;
    vrConfig.eyeViewportLeft[1] = 0;
    vrConfig.eyeViewportLeft[2] = hmd.hResolution/2;
    vrConfig.eyeViewportLeft[3] = hmd.vResolution;

    if (vrConfig.distortionShader.id > 0)
    {
        // Update distortion shader with lens and distortion-scale parameters
        SetShaderValue(vrConfig.distortionShader, GetShaderLocation(vrConfig.distortionShader, "leftLensCenter"), leftLensCenter, UNIFORM_VEC2);
        SetShaderValue(vrConfig.distortionShader, GetShaderLocation(vrConfig.distortionShader, "rightLensCenter"), rightLensCenter, UNIFORM_VEC2);
        SetShaderValue(vrConfig.distortionShader, GetShaderLocation(vrConfig.distortionShader, "leftScreenCenter"), leftScreenCenter, UNIFORM_VEC2);
        SetShaderValue(vrConfig.distortionShader, GetShaderLocation(vrConfig.distortionShader, "rightScreenCenter"), rightScreenCenter, UNIFORM_VEC2);

        SetShaderValue(vrConfig.distortionShader, GetShaderLocation(vrConfig.distortionShader, "scale"), scale, UNIFORM_VEC2);
        SetShaderValue(vrConfig.distortionShader, GetShaderLocation(vrConfig.distortionShader, "scaleIn"), scaleIn, UNIFORM_VEC2);
        SetShaderValue(vrConfig.distortionShader, GetShaderLocation(vrConfig.distortionShader, "hmdWarpParam"), hmd.lensDistortionValues, UNIFORM_VEC4);
        SetShaderValue(vrConfig.distortionShader, GetShaderLocation(vrConfig.distortionShader, "chromaAbParam"), hmd.chromaAbCorrection, UNIFORM_VEC4);
    }
#endif
}

// Detect if VR simulator is running
bool IsVrSimulatorReady(void)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    return vrSimulatorReady;
#else
    return false;
#endif
}

// Enable/Disable VR experience (device or simulator)
void ToggleVrMode(void)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    vrSimulatorReady = !vrSimulatorReady;

    if (!vrSimulatorReady)
    {
        vrStereoRender = false;

        // Reset viewport and default projection-modelview matrices
        rlViewport(0, 0, framebufferWidth, framebufferHeight);
        projection = MatrixOrtho(0.0, framebufferWidth, framebufferHeight, 0.0, 0.0, 1.0);
        modelview = MatrixIdentity();
    }
    else vrStereoRender = true;
#endif
}

// Begin VR drawing configuration
void BeginVrDrawing(void)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if (vrSimulatorReady)
    {

        rlEnableRenderTexture(stereoFbo.id);    // Setup framebuffer for stereo rendering
        //glEnable(GL_FRAMEBUFFER_SRGB);        // Enable SRGB framebuffer (only if required)

        //glViewport(0, 0, buffer.width, buffer.height);    // Useful if rendering to separate framebuffers (every eye)
        rlClearScreenBuffers();                 // Clear current framebuffer

        vrStereoRender = true;
    }
#endif
}

// End VR drawing process (and desktop mirror)
void EndVrDrawing(void)
{
#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
    if (vrSimulatorReady)
    {
        vrStereoRender = false;         // Disable stereo render

        rlDisableRenderTexture();       // Unbind current framebuffer

        rlClearScreenBuffers();         // Clear current framebuffer

        // Set viewport to default framebuffer size (screen size)
        rlViewport(0, 0, framebufferWidth, framebufferHeight);

        // Let rlgl reconfigure internal matrices
        rlMatrixMode(RL_PROJECTION);                            // Enable internal projection matrix
        rlLoadIdentity();                                       // Reset internal projection matrix
        rlOrtho(0.0, framebufferWidth, framebufferHeight, 0.0, 0.0, 1.0); // Recalculate internal projection matrix
        rlMatrixMode(RL_MODELVIEW);                             // Enable internal modelview matrix
        rlLoadIdentity();                                       // Reset internal modelview matrix

        // Draw RenderTexture (stereoFbo) using distortion shader if available
        if (vrConfig.distortionShader.id > 0) currentShader = vrConfig.distortionShader;
        else currentShader = GetShaderDefault();

        rlEnableTexture(stereoFbo.texture.id);

        rlPushMatrix();
            rlBegin(RL_QUADS);
                rlColor4ub(255, 255, 255, 255);
                rlNormal3f(0.0f, 0.0f, 1.0f);

                // Bottom-left corner for texture and quad
                rlTexCoord2f(0.0f, 1.0f);
                rlVertex2f(0.0f, 0.0f);

                // Bottom-right corner for texture and quad
                rlTexCoord2f(0.0f, 0.0f);
                rlVertex2f(0.0f, (float)stereoFbo.texture.height);

                // Top-right corner for texture and quad
                rlTexCoord2f(1.0f, 0.0f);
                rlVertex2f( (float)stereoFbo.texture.width, (float)stereoFbo.texture.height);

                // Top-left corner for texture and quad
                rlTexCoord2f(1.0f, 1.0f);
                rlVertex2f( (float)stereoFbo.texture.width, 0.0f);
            rlEnd();
        rlPopMatrix();

        rlDisableTexture();

        // Update and draw render texture fbo with distortion to backbuffer
        UpdateBuffersDefault();
        DrawBuffersDefault();

        // Restore defaultShader
        currentShader = defaultShader;

        // Reset viewport and default projection-modelview matrices
        rlViewport(0, 0, framebufferWidth, framebufferHeight);
        projection = MatrixOrtho(0.0, framebufferWidth, framebufferHeight, 0.0, 0.0, 1.0);
        modelview = MatrixIdentity();

        rlDisableDepthTest();
    }
#endif
}
#endif          // SUPPORT_VR_SIMULATOR

//----------------------------------------------------------------------------------
// Module specific Functions Definition
//----------------------------------------------------------------------------------

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)
// Compile custom shader and return shader id
static unsigned int CompileShader(const char *shaderStr, int type)
{
    unsigned int shader = glCreateShader(type);
    glShaderSource(shader, 1, &shaderStr, NULL);

    GLint success = 0;
    glCompileShader(shader);
    glGetShaderiv(shader, GL_COMPILE_STATUS, &success);

    if (success != GL_TRUE)
    {
        TraceLog(LOG_WARNING, "[SHDR ID %i] Failed to compile shader...", shader);
        int maxLength = 0;
        int length;
        glGetShaderiv(shader, GL_INFO_LOG_LENGTH, &maxLength);

#if defined(_MSC_VER)
        char *log = RL_MALLOC(maxLength);
#else
        char log[maxLength];
#endif
        glGetShaderInfoLog(shader, maxLength, &length, log);

        TraceLog(LOG_INFO, "%s", log);

#if defined(_MSC_VER)
        RL_FREE(log);
#endif
    }
    else TraceLog(LOG_INFO, "[SHDR ID %i] Shader compiled successfully", shader);

    return shader;
}

// Load custom shader strings and return program id
static unsigned int LoadShaderProgram(unsigned int vShaderId, unsigned int fShaderId)
{
    unsigned int program = 0;

#if defined(GRAPHICS_API_OPENGL_33) || defined(GRAPHICS_API_OPENGL_ES2)

    GLint success = 0;
    program = glCreateProgram();

    glAttachShader(program, vShaderId);
    glAttachShader(program, fShaderId);

    // NOTE: Default attribute shader locations must be binded before linking
    glBindAttribLocation(program, 0, DEFAULT_ATTRIB_POSITION_NAME);
    glBindAttribLocation(program, 1, DEFAULT_ATTRIB_TEXCOORD_NAME);
    glBindAttribLocation(program, 2, DEFAULT_ATTRIB_NORMAL_NAME);
    glBindAttribLocation(program, 3, DEFAULT_ATTRIB_COLOR_NAME);
    glBindAttribLocation(program, 4, DEFAULT_ATTRIB_TANGENT_NAME);
    glBindAttribLocation(program, 5, DEFAULT_ATTRIB_TEXCOORD2_NAME);

    // NOTE: If some attrib name is no found on the shader, it locations becomes -1

    glLinkProgram(program);

    // NOTE: All uniform variables are intitialised to 0 when a program links

    glGetProgramiv(program, GL_LINK_STATUS, &success);

    if (success == GL_FALSE)
    {
        TraceLog(LOG_WARNING, "[SHDR ID %i] Failed to link shader program...", program);

        int maxLength = 0;
        int length;

        glGetProgramiv(program, GL_INFO_LOG_LENGTH, &maxLength);

#if defined(_MSC_VER)
        char *log = RL_MALLOC(maxLength);
#else
        char log[maxLength];
#endif
        glGetProgramInfoLog(program, maxLength, &length, log);

        TraceLog(LOG_INFO, "%s", log);

#if defined(_MSC_VER)
        RL_FREE(log);
#endif
        glDeleteProgram(program);

        program = 0;
    }
    else TraceLog(LOG_INFO, "[SHDR ID %i] Shader program loaded successfully", program);
#endif
    return program;
}


// Load default shader (just vertex positioning and texture coloring)
// NOTE: This shader program is used for internal buffers
static Shader LoadShaderDefault(void)
{
    Shader shader = { 0 };
    shader.locs = (int *)RL_CALLOC(MAX_SHADER_LOCATIONS, sizeof(int));

    // NOTE: All locations must be reseted to -1 (no location)
    for (int i = 0; i < MAX_SHADER_LOCATIONS; i++) shader.locs[i] = -1;

    // Vertex shader directly defined, no external file required
    const char *defaultVShaderStr =
#if defined(GRAPHICS_API_OPENGL_21)
    "#version 120                       \n"
#elif defined(GRAPHICS_API_OPENGL_ES2)
    "#version 100                       \n"
#endif
#if defined(GRAPHICS_API_OPENGL_ES2) || defined(GRAPHICS_API_OPENGL_21)
    "attribute vec3 vertexPosition;     \n"
    "attribute vec2 vertexTexCoord;     \n"
    "attribute vec4 vertexColor;        \n"
    "varying vec2 fragTexCoord;         \n"
    "varying vec4 fragColor;            \n"
#elif defined(GRAPHICS_API_OPENGL_33)
    "#version 330                       \n"
    "in vec3 vertexPosition;            \n"
    "in vec2 vertexTexCoord;            \n"
    "in vec4 vertexColor;               \n"
    "out vec2 fragTexCoord;             \n"
    "out vec4 fragColor;                \n"
#endif
    "uniform mat4 mvp;                  \n"
    "void main()                        \n"
    "{                                  \n"
    "    fragTexCoord = vertexTexCoord; \n"
    "    fragColor = vertexColor;       \n"
    "    gl_Position = mvp*vec4(vertexPosition, 1.0); \n"
    "}                                  \n";

    // Fragment shader directly defined, no external file required
    const char *defaultFShaderStr =
#if defined(GRAPHICS_API_OPENGL_21)
    "#version 120                       \n"
#elif defined(GRAPHICS_API_OPENGL_ES2)
    "#version 100                       \n"
    "precision mediump float;           \n"     // precision required for OpenGL ES2 (WebGL)
#endif
#if defined(GRAPHICS_API_OPENGL_ES2) || defined(GRAPHICS_API_OPENGL_21)
    "varying vec2 fragTexCoord;         \n"
    "varying vec4 fragColor;            \n"
#elif defined(GRAPHICS_API_OPENGL_33)
    "#version 330       \n"
    "in vec2 fragTexCoord;              \n"
    "in vec4 fragColor;                 \n"
    "out vec4 finalColor;               \n"
#endif
    "uniform sampler2D texture0;        \n"
    "uniform vec4 colDiffuse;           \n"
    "void main()                        \n"
    "{                                  \n"
#if defined(GRAPHICS_API_OPENGL_ES2) || defined(GRAPHICS_API_OPENGL_21)
    "    vec4 texelColor = texture2D(texture0, fragTexCoord); \n" // NOTE: texture2D() is deprecated on OpenGL 3.3 and ES 3.0
    "    gl_FragColor = texelColor*colDiffuse*fragColor;      \n"
#elif defined(GRAPHICS_API_OPENGL_33)
    "    vec4 texelColor = texture(texture0, fragTexCoord);   \n"
    "    finalColor = texelColor*colDiffuse*fragColor;        \n"
#endif
    "}                                  \n";

    // NOTE: Compiled vertex/fragment shaders are kept for re-use
    defaultVShaderId = CompileShader(defaultVShaderStr, GL_VERTEX_SHADER);     // Compile default vertex shader
    defaultFShaderId = CompileShader(defaultFShaderStr, GL_FRAGMENT_SHADER);   // Compile default fragment shader

    shader.id = LoadShaderProgram(defaultVShaderId, defaultFShaderId);

    if (shader.id > 0)
    {
        TraceLog(LOG_INFO, "[SHDR ID %i] Default shader loaded successfully", shader.id);

        // Set default shader locations: attributes locations
        shader.locs[LOC_VERTEX_POSITION] = glGetAttribLocation(shader.id, "vertexPosition");
        shader.locs[LOC_VERTEX_TEXCOORD01] = glGetAttribLocation(shader.id, "vertexTexCoord");
        shader.locs[LOC_VERTEX_COLOR] = glGetAttribLocation(shader.id, "vertexColor");

        // Set default shader locations: uniform locations
        shader.locs[LOC_MATRIX_MVP]  = glGetUniformLocation(shader.id, "mvp");
        shader.locs[LOC_COLOR_DIFFUSE] = glGetUniformLocation(shader.id, "colDiffuse");
        shader.locs[LOC_MAP_DIFFUSE] = glGetUniformLocation(shader.id, "texture0");

        // NOTE: We could also use below function but in case DEFAULT_ATTRIB_* points are
        // changed for external custom shaders, we just use direct bindings above
        //SetShaderDefaultLocations(&shader);
    }
    else TraceLog(LOG_WARNING, "[SHDR ID %i] Default shader could not be loaded", shader.id);

    return shader;
}

// Get location handlers to for shader attributes and uniforms
// NOTE: If any location is not found, loc point becomes -1
static void SetShaderDefaultLocations(Shader *shader)
{
    // NOTE: Default shader attrib locations have been fixed before linking:
    //          vertex position location    = 0
    //          vertex texcoord location    = 1
    //          vertex normal location      = 2
    //          vertex color location       = 3
    //          vertex tangent location     = 4
    //          vertex texcoord2 location   = 5

    // Get handles to GLSL input attibute locations
    shader->locs[LOC_VERTEX_POSITION] = glGetAttribLocation(shader->id, DEFAULT_ATTRIB_POSITION_NAME);
    shader->locs[LOC_VERTEX_TEXCOORD01] = glGetAttribLocation(shader->id, DEFAULT_ATTRIB_TEXCOORD_NAME);
    shader->locs[LOC_VERTEX_TEXCOORD02] = glGetAttribLocation(shader->id, DEFAULT_ATTRIB_TEXCOORD2_NAME);
    shader->locs[LOC_VERTEX_NORMAL] = glGetAttribLocation(shader->id, DEFAULT_ATTRIB_NORMAL_NAME);
    shader->locs[LOC_VERTEX_TANGENT] = glGetAttribLocation(shader->id, DEFAULT_ATTRIB_TANGENT_NAME);
    shader->locs[LOC_VERTEX_COLOR] = glGetAttribLocation(shader->id, DEFAULT_ATTRIB_COLOR_NAME);

    // Get handles to GLSL uniform locations (vertex shader)
    shader->locs[LOC_MATRIX_MVP]  = glGetUniformLocation(shader->id, "mvp");
    shader->locs[LOC_MATRIX_PROJECTION]  = glGetUniformLocation(shader->id, "projection");
    shader->locs[LOC_MATRIX_VIEW]  = glGetUniformLocation(shader->id, "view");

    // Get handles to GLSL uniform locations (fragment shader)
    shader->locs[LOC_COLOR_DIFFUSE] = glGetUniformLocation(shader->id, "colDiffuse");
    shader->locs[LOC_MAP_DIFFUSE] = glGetUniformLocation(shader->id, "texture0");
    shader->locs[LOC_MAP_SPECULAR] = glGetUniformLocation(shader->id, "texture1");
    shader->locs[LOC_MAP_NORMAL] = glGetUniformLocation(shader->id, "texture2");
}

// Unload default shader
static void UnloadShaderDefault(void)
{
    glUseProgram(0);

    glDetachShader(defaultShader.id, defaultVShaderId);
    glDetachShader(defaultShader.id, defaultFShaderId);
    glDeleteShader(defaultVShaderId);
    glDeleteShader(defaultFShaderId);

    glDeleteProgram(defaultShader.id);
}

// Load default internal buffers
static void LoadBuffersDefault(void)
{
    // Initialize CPU (RAM) arrays (vertex position, texcoord, color data and indexes)
    //--------------------------------------------------------------------------------------------
    for (int i = 0; i < MAX_BATCH_BUFFERING; i++)
    {
        vertexData[i].vertices = (float *)RL_MALLOC(sizeof(float)*3*4*MAX_BATCH_ELEMENTS);        // 3 float by vertex, 4 vertex by quad
        vertexData[i].texcoords = (float *)RL_MALLOC(sizeof(float)*2*4*MAX_BATCH_ELEMENTS);       // 2 float by texcoord, 4 texcoord by quad
        vertexData[i].colors = (unsigned char *)RL_MALLOC(sizeof(unsigned char)*4*4*MAX_BATCH_ELEMENTS);  // 4 float by color, 4 colors by quad
#if defined(GRAPHICS_API_OPENGL_33)
        vertexData[i].indices = (unsigned int *)RL_MALLOC(sizeof(unsigned int)*6*MAX_BATCH_ELEMENTS);      // 6 int by quad (indices)
#elif defined(GRAPHICS_API_OPENGL_ES2)
        vertexData[i].indices = (unsigned short *)RL_MALLOC(sizeof(unsigned short)*6*MAX_BATCH_ELEMENTS);  // 6 int by quad (indices)
#endif

        for (int j = 0; j < (3*4*MAX_BATCH_ELEMENTS); j++) vertexData[i].vertices[j] = 0.0f;
        for (int j = 0; j < (2*4*MAX_BATCH_ELEMENTS); j++) vertexData[i].texcoords[j] = 0.0f;
        for (int j = 0; j < (4*4*MAX_BATCH_ELEMENTS); j++) vertexData[i].colors[j] = 0;

        int k = 0;

        // Indices can be initialized right now
        for (int j = 0; j < (6*MAX_BATCH_ELEMENTS); j += 6)
        {
            vertexData[i].indices[j] = 4*k;
            vertexData[i].indices[j + 1] = 4*k + 1;
            vertexData[i].indices[j + 2] = 4*k + 2;
            vertexData[i].indices[j + 3] = 4*k;
            vertexData[i].indices[j + 4] = 4*k + 2;
            vertexData[i].indices[j + 5] = 4*k + 3;

            k++;
        }

        vertexData[i].vCounter = 0;
        vertexData[i].tcCounter = 0;
        vertexData[i].cCounter = 0;
    }

    TraceLog(LOG_INFO, "Internal buffers initialized successfully (CPU)");
    //--------------------------------------------------------------------------------------------

    // Upload to GPU (VRAM) vertex data and initialize VAOs/VBOs
    //--------------------------------------------------------------------------------------------
    for (int i = 0; i < MAX_BATCH_BUFFERING; i++)
    {
        if (vaoSupported)
        {
            // Initialize Quads VAO
            glGenVertexArrays(1, &vertexData[i].vaoId);
            glBindVertexArray(vertexData[i].vaoId);
        }

        // Quads - Vertex buffers binding and attributes enable
        // Vertex position buffer (shader-location = 0)
        glGenBuffers(1, &vertexData[i].vboId[0]);
        glBindBuffer(GL_ARRAY_BUFFER, vertexData[i].vboId[0]);
        glBufferData(GL_ARRAY_BUFFER, sizeof(float)*3*4*MAX_BATCH_ELEMENTS, vertexData[i].vertices, GL_DYNAMIC_DRAW);
        glEnableVertexAttribArray(currentShader.locs[LOC_VERTEX_POSITION]);
        glVertexAttribPointer(currentShader.locs[LOC_VERTEX_POSITION], 3, GL_FLOAT, 0, 0, 0);

        // Vertex texcoord buffer (shader-location = 1)
        glGenBuffers(1, &vertexData[i].vboId[1]);
        glBindBuffer(GL_ARRAY_BUFFER, vertexData[i].vboId[1]);
        glBufferData(GL_ARRAY_BUFFER, sizeof(float)*2*4*MAX_BATCH_ELEMENTS, vertexData[i].texcoords, GL_DYNAMIC_DRAW);
        glEnableVertexAttribArray(currentShader.locs[LOC_VERTEX_TEXCOORD01]);
        glVertexAttribPointer(currentShader.locs[LOC_VERTEX_TEXCOORD01], 2, GL_FLOAT, 0, 0, 0);

        // Vertex color buffer (shader-location = 3)
        glGenBuffers(1, &vertexData[i].vboId[2]);
        glBindBuffer(GL_ARRAY_BUFFER, vertexData[i].vboId[2]);
        glBufferData(GL_ARRAY_BUFFER, sizeof(unsigned char)*4*4*MAX_BATCH_ELEMENTS, vertexData[i].colors, GL_DYNAMIC_DRAW);
        glEnableVertexAttribArray(currentShader.locs[LOC_VERTEX_COLOR]);
        glVertexAttribPointer(currentShader.locs[LOC_VERTEX_COLOR], 4, GL_UNSIGNED_BYTE, GL_TRUE, 0, 0);

        // Fill index buffer
        glGenBuffers(1, &vertexData[i].vboId[3]);
        glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, vertexData[i].vboId[3]);
#if defined(GRAPHICS_API_OPENGL_33)
        glBufferData(GL_ELEMENT_ARRAY_BUFFER, sizeof(int)*6*MAX_BATCH_ELEMENTS, vertexData[i].indices, GL_STATIC_DRAW);
#elif defined(GRAPHICS_API_OPENGL_ES2)
        glBufferData(GL_ELEMENT_ARRAY_BUFFER, sizeof(short)*6*MAX_BATCH_ELEMENTS, vertexData[i].indices, GL_STATIC_DRAW);
#endif
    }

    TraceLog(LOG_INFO, "Internal buffers uploaded successfully (GPU)");

    // Unbind the current VAO
    if (vaoSupported) glBindVertexArray(0);
    //--------------------------------------------------------------------------------------------
}

// Update default internal buffers (VAOs/VBOs) with vertex array data
// NOTE: If there is not vertex data, buffers doesn't need to be updated (vertexCount > 0)
// TODO: If no data changed on the CPU arrays --> No need to re-update GPU arrays (change flag required)
static void UpdateBuffersDefault(void)
{
    // Update vertex buffers data
    if (vertexData[currentBuffer].vCounter > 0)
    {
        // Activate elements VAO
        if (vaoSupported) glBindVertexArray(vertexData[currentBuffer].vaoId);

        // Vertex positions buffer
        glBindBuffer(GL_ARRAY_BUFFER, vertexData[currentBuffer].vboId[0]);
        glBufferSubData(GL_ARRAY_BUFFER, 0, sizeof(float)*3*vertexData[currentBuffer].vCounter, vertexData[currentBuffer].vertices);
        //glBufferData(GL_ARRAY_BUFFER, sizeof(float)*3*4*MAX_BATCH_ELEMENTS, vertexData[currentBuffer].vertices, GL_DYNAMIC_DRAW);  // Update all buffer

        // Texture coordinates buffer
        glBindBuffer(GL_ARRAY_BUFFER, vertexData[currentBuffer].vboId[1]);
        glBufferSubData(GL_ARRAY_BUFFER, 0, sizeof(float)*2*vertexData[currentBuffer].vCounter, vertexData[currentBuffer].texcoords);
        //glBufferData(GL_ARRAY_BUFFER, sizeof(float)*2*4*MAX_BATCH_ELEMENTS, vertexData[currentBuffer].texcoords, GL_DYNAMIC_DRAW); // Update all buffer

        // Colors buffer
        glBindBuffer(GL_ARRAY_BUFFER, vertexData[currentBuffer].vboId[2]);
        glBufferSubData(GL_ARRAY_BUFFER, 0, sizeof(unsigned char)*4*vertexData[currentBuffer].vCounter, vertexData[currentBuffer].colors);
        //glBufferData(GL_ARRAY_BUFFER, sizeof(float)*4*4*MAX_BATCH_ELEMENTS, vertexData[currentBuffer].colors, GL_DYNAMIC_DRAW);    // Update all buffer

        // NOTE: glMapBuffer() causes sync issue.
        // If GPU is working with this buffer, glMapBuffer() will wait(stall) until GPU to finish its job.
        // To avoid waiting (idle), you can call first glBufferData() with NULL pointer before glMapBuffer().
        // If you do that, the previous data in PBO will be discarded and glMapBuffer() returns a new
        // allocated pointer immediately even if GPU is still working with the previous data.

        // Another option: map the buffer object into client's memory
        // Probably this code could be moved somewhere else...
        // vertexData[currentBuffer].vertices = (float *)glMapBuffer(GL_ARRAY_BUFFER, GL_READ_WRITE);
        // if (vertexData[currentBuffer].vertices)
        // {
            // Update vertex data
        // }
        // glUnmapBuffer(GL_ARRAY_BUFFER);

        // Unbind the current VAO
        if (vaoSupported) glBindVertexArray(0);
    }
}

// Draw default internal buffers vertex data
static void DrawBuffersDefault(void)
{
    Matrix matProjection = projection;
    Matrix matModelView = modelview;

    int eyesCount = 1;
#if defined(SUPPORT_VR_SIMULATOR)
    if (vrStereoRender) eyesCount = 2;
#endif

    for (int eye = 0; eye < eyesCount; eye++)
    {
#if defined(SUPPORT_VR_SIMULATOR)
        if (eyesCount == 2) SetStereoView(eye, matProjection, matModelView);
#endif

        // Draw buffers
        if (vertexData[currentBuffer].vCounter > 0)
        {
            // Set current shader and upload current MVP matrix
            glUseProgram(currentShader.id);

            // Create modelview-projection matrix
            Matrix matMVP = MatrixMultiply(modelview, projection);

            glUniformMatrix4fv(currentShader.locs[LOC_MATRIX_MVP], 1, false, MatrixToFloat(matMVP));
            glUniform4f(currentShader.locs[LOC_COLOR_DIFFUSE], 1.0f, 1.0f, 1.0f, 1.0f);
            glUniform1i(currentShader.locs[LOC_MAP_DIFFUSE], 0);    // Provided value refers to the texture unit (active)

            // TODO: Support additional texture units on custom shader
            //if (currentShader->locs[LOC_MAP_SPECULAR] > 0) glUniform1i(currentShader.locs[LOC_MAP_SPECULAR], 1);
            //if (currentShader->locs[LOC_MAP_NORMAL] > 0) glUniform1i(currentShader.locs[LOC_MAP_NORMAL], 2);

            // NOTE: Right now additional map textures not considered for default buffers drawing

            int vertexOffset = 0;

            if (vaoSupported) glBindVertexArray(vertexData[currentBuffer].vaoId);
            else
            {
                // Bind vertex attrib: position (shader-location = 0)
                glBindBuffer(GL_ARRAY_BUFFER, vertexData[currentBuffer].vboId[0]);
                glVertexAttribPointer(currentShader.locs[LOC_VERTEX_POSITION], 3, GL_FLOAT, 0, 0, 0);
                glEnableVertexAttribArray(currentShader.locs[LOC_VERTEX_POSITION]);

                // Bind vertex attrib: texcoord (shader-location = 1)
                glBindBuffer(GL_ARRAY_BUFFER, vertexData[currentBuffer].vboId[1]);
                glVertexAttribPointer(currentShader.locs[LOC_VERTEX_TEXCOORD01], 2, GL_FLOAT, 0, 0, 0);
                glEnableVertexAttribArray(currentShader.locs[LOC_VERTEX_TEXCOORD01]);

                // Bind vertex attrib: color (shader-location = 3)
                glBindBuffer(GL_ARRAY_BUFFER, vertexData[currentBuffer].vboId[2]);
                glVertexAttribPointer(currentShader.locs[LOC_VERTEX_COLOR], 4, GL_UNSIGNED_BYTE, GL_TRUE, 0, 0);
                glEnableVertexAttribArray(currentShader.locs[LOC_VERTEX_COLOR]);

                glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, vertexData[currentBuffer].vboId[3]);
            }

            glActiveTexture(GL_TEXTURE0);

            for (int i = 0; i < drawsCounter; i++)
            {
                glBindTexture(GL_TEXTURE_2D, draws[i].textureId);

                // TODO: Find some way to bind additional textures --> Use global texture IDs? Register them on draw[i]?
                //if (currentShader->locs[LOC_MAP_SPECULAR] > 0) { glActiveTexture(GL_TEXTURE1); glBindTexture(GL_TEXTURE_2D, textureUnit1_id); }
                //if (currentShader->locs[LOC_MAP_SPECULAR] > 0) { glActiveTexture(GL_TEXTURE2); glBindTexture(GL_TEXTURE_2D, textureUnit2_id); }

                if ((draws[i].mode == RL_LINES) || (draws[i].mode == RL_TRIANGLES)) glDrawArrays(draws[i].mode, vertexOffset, draws[i].vertexCount);
                else
                {
#if defined(GRAPHICS_API_OPENGL_33)
                    // We need to define the number of indices to be processed: quadsCount*6
                    // NOTE: The final parameter tells the GPU the offset in bytes from the
                    // start of the index buffer to the location of the first index to process
                    glDrawElements(GL_TRIANGLES, draws[i].vertexCount/4*6, GL_UNSIGNED_INT, (GLvoid *)(sizeof(GLuint)*vertexOffset/4*6));
#elif defined(GRAPHICS_API_OPENGL_ES2)
                    glDrawElements(GL_TRIANGLES, draws[i].vertexCount/4*6, GL_UNSIGNED_SHORT, (GLvoid *)(sizeof(GLushort)*vertexOffset/4*6));
#endif
                }

                vertexOffset += (draws[i].vertexCount + draws[i].vertexAlignment);
            }

            if (!vaoSupported)
            {
                glBindBuffer(GL_ARRAY_BUFFER, 0);
                glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, 0);
            }

            glBindTexture(GL_TEXTURE_2D, 0);    // Unbind textures
        }

        if (vaoSupported) glBindVertexArray(0); // Unbind VAO

        glUseProgram(0);    // Unbind shader program
    }

    // Reset vertex counters for next frame
    vertexData[currentBuffer].vCounter = 0;
    vertexData[currentBuffer].tcCounter = 0;
    vertexData[currentBuffer].cCounter = 0;

    // Reset depth for next draw
    currentDepth = -1.0f;

    // Restore projection/modelview matrices
    projection = matProjection;
    modelview = matModelView;

    // Reset draws array
    for (int i = 0; i < MAX_DRAWCALL_REGISTERED; i++)
    {
        draws[i].mode = RL_QUADS;
        draws[i].vertexCount = 0;
        draws[i].textureId = defaultTextureId;
    }

    drawsCounter = 1;

    // Change to next buffer in the list
    currentBuffer++;
    if (currentBuffer >= MAX_BATCH_BUFFERING) currentBuffer = 0;
}

// Unload default internal buffers vertex data from CPU and GPU
static void UnloadBuffersDefault(void)
{
    // Unbind everything
    if (vaoSupported) glBindVertexArray(0);
    glDisableVertexAttribArray(0);
    glDisableVertexAttribArray(1);
    glDisableVertexAttribArray(2);
    glDisableVertexAttribArray(3);
    glBindBuffer(GL_ARRAY_BUFFER, 0);
    glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, 0);

    for (int i = 0; i < MAX_BATCH_BUFFERING; i++)
    {
        // Delete VBOs from GPU (VRAM)
        glDeleteBuffers(1, &vertexData[i].vboId[0]);
        glDeleteBuffers(1, &vertexData[i].vboId[1]);
        glDeleteBuffers(1, &vertexData[i].vboId[2]);
        glDeleteBuffers(1, &vertexData[i].vboId[3]);

        // Delete VAOs from GPU (VRAM)
        if (vaoSupported) glDeleteVertexArrays(1, &vertexData[i].vaoId);

        // Free vertex arrays memory from CPU (RAM)
        RL_FREE(vertexData[i].vertices);
        RL_FREE(vertexData[i].texcoords);
        RL_FREE(vertexData[i].colors);
        RL_FREE(vertexData[i].indices);
    }
}

// Renders a 1x1 XY quad in NDC
static void GenDrawQuad(void)
{
    unsigned int quadVAO = 0;
    unsigned int quadVBO = 0;

    float vertices[] = {
        // Positions        // Texture Coords
        -1.0f, 1.0f, 0.0f, 0.0f, 1.0f,
        -1.0f, -1.0f, 0.0f, 0.0f, 0.0f,
        1.0f, 1.0f, 0.0f, 1.0f, 1.0f,
        1.0f, -1.0f, 0.0f, 1.0f, 0.0f,
    };

    // Set up plane VAO
    glGenVertexArrays(1, &quadVAO);
    glGenBuffers(1, &quadVBO);
    glBindVertexArray(quadVAO);

    // Fill buffer
    glBindBuffer(GL_ARRAY_BUFFER, quadVBO);
    glBufferData(GL_ARRAY_BUFFER, sizeof(vertices), &vertices, GL_STATIC_DRAW);

    // Link vertex attributes
    glEnableVertexAttribArray(0);
    glVertexAttribPointer(0, 3, GL_FLOAT, GL_FALSE, 5*sizeof(float), (void *)0);
    glEnableVertexAttribArray(1);
    glVertexAttribPointer(1, 2, GL_FLOAT, GL_FALSE, 5*sizeof(float), (void *)(3*sizeof(float)));

    // Draw quad
    glBindVertexArray(quadVAO);
    glDrawArrays(GL_TRIANGLE_STRIP, 0, 4);
    glBindVertexArray(0);

    glDeleteBuffers(1, &quadVBO);
    glDeleteVertexArrays(1, &quadVAO);
}

// Renders a 1x1 3D cube in NDC
static void GenDrawCube(void)
{
    unsigned int cubeVAO = 0;
    unsigned int cubeVBO = 0;

    float vertices[] = {
        -1.0f, -1.0f, -1.0f,  0.0f, 0.0f, -1.0f, 0.0f, 0.0f,
        1.0f, 1.0f, -1.0f, 0.0f, 0.0f, -1.0f, 1.0f, 1.0f,
        1.0f, -1.0f, -1.0f, 0.0f, 0.0f, -1.0f, 1.0f, 0.0f,
        1.0f, 1.0f, -1.0f, 0.0f, 0.0f, -1.0f, 1.0f, 1.0f,
        -1.0f, -1.0f, -1.0f, 0.0f, 0.0f, -1.0f, 0.0f, 0.0f,
        -1.0f, 1.0f, -1.0f, 0.0f, 0.0f, -1.0f, 0.0f, 1.0f,
        -1.0f, -1.0f, 1.0f, 0.0f, 0.0f, 1.0f, 0.0f, 0.0f,
        1.0f, -1.0f, 1.0f, 0.0f, 0.0f, 1.0f, 1.0f, 0.0f,
        1.0f, 1.0f, 1.0f, 0.0f, 0.0f, 1.0f, 1.0f, 1.0f,
        1.0f, 1.0f, 1.0f, 0.0f, 0.0f, 1.0f, 1.0f, 1.0f,
        -1.0f, 1.0f, 1.0f, 0.0f, 0.0f, 1.0f, 0.0f, 1.0f,
        -1.0f, -1.0f, 1.0f, 0.0f, 0.0f, 1.0f, 0.0f, 0.0f,
        -1.0f, 1.0f, 1.0f, -1.0f, 0.0f, 0.0f, 1.0f, 0.0f,
        -1.0f, 1.0f, -1.0f, -1.0f, 0.0f, 0.0f, 1.0f, 1.0f,
        -1.0f, -1.0f, -1.0f, -1.0f, 0.0f, 0.0f, 0.0f, 1.0f,
        -1.0f, -1.0f, -1.0f, -1.0f, 0.0f, 0.0f, 0.0f, 1.0f,
        -1.0f, -1.0f, 1.0f, -1.0f, 0.0f, 0.0f, 0.0f, 0.0f,
        -1.0f, 1.0f, 1.0f, -1.0f, 0.0f, 0.0f, 1.0f, 0.0f,
        1.0f, 1.0f, 1.0f, 1.0f, 0.0f, 0.0f, 1.0f, 0.0f,
        1.0f, -1.0f, -1.0f, 1.0f, 0.0f, 0.0f, 0.0f, 1.0f,
        1.0f, 1.0f, -1.0f, 1.0f, 0.0f, 0.0f, 1.0f, 1.0f,
        1.0f, -1.0f, -1.0f, 1.0f, 0.0f, 0.0f, 0.0f, 1.0f,
        1.0f, 1.0f, 1.0f, 1.0f, 0.0f, 0.0f, 1.0f, 0.0f,
        1.0f, -1.0f, 1.0f, 1.0f, 0.0f, 0.0f, 0.0f, 0.0f,
        -1.0f, -1.0f, -1.0f, 0.0f, -1.0f, 0.0f, 0.0f, 1.0f,
        1.0f, -1.0f, -1.0f, 0.0f, -1.0f, 0.0f, 1.0f, 1.0f,
        1.0f, -1.0f, 1.0f, 0.0f, -1.0f, 0.0f, 1.0f, 0.0f,
        1.0f, -1.0f, 1.0f, 0.0f, -1.0f, 0.0f, 1.0f, 0.0f,
        -1.0f, -1.0f, 1.0f, 0.0f, -1.0f, 0.0f, 0.0f, 0.0f,
        -1.0f, -1.0f, -1.0f, 0.0f, -1.0f, 0.0f, 0.0f, 1.0f,
        -1.0f, 1.0f, -1.0f, 0.0f, 1.0f, 0.0f, 0.0f, 1.0f,
        1.0f, 1.0f , 1.0f, 0.0f, 1.0f, 0.0f, 1.0f, 0.0f,
        1.0f, 1.0f, -1.0f, 0.0f, 1.0f, 0.0f, 1.0f, 1.0f,
        1.0f, 1.0f, 1.0f, 0.0f, 1.0f, 0.0f, 1.0f, 0.0f,
        -1.0f, 1.0f, -1.0f, 0.0f, 1.0f, 0.0f, 0.0f, 1.0f,
        -1.0f, 1.0f, 1.0f, 0.0f, 1.0f, 0.0f, 0.0f, 0.0f
    };

    // Set up cube VAO
    glGenVertexArrays(1, &cubeVAO);
    glGenBuffers(1, &cubeVBO);

    // Fill buffer
    glBindBuffer(GL_ARRAY_BUFFER, cubeVBO);
    glBufferData(GL_ARRAY_BUFFER, sizeof(vertices), vertices, GL_STATIC_DRAW);

    // Link vertex attributes
    glBindVertexArray(cubeVAO);
    glEnableVertexAttribArray(0);
    glVertexAttribPointer(0, 3, GL_FLOAT, GL_FALSE, 8*sizeof(float), (void *)0);
    glEnableVertexAttribArray(1);
    glVertexAttribPointer(1, 3, GL_FLOAT, GL_FALSE, 8*sizeof(float), (void *)(3*sizeof(float)));
    glEnableVertexAttribArray(2);
    glVertexAttribPointer(2, 2, GL_FLOAT, GL_FALSE, 8*sizeof(float), (void *)(6*sizeof(float)));
    glBindBuffer(GL_ARRAY_BUFFER, 0);
    glBindVertexArray(0);

    // Draw cube
    glBindVertexArray(cubeVAO);
    glDrawArrays(GL_TRIANGLES, 0, 36);
    glBindVertexArray(0);

    glDeleteBuffers(1, &cubeVBO);
    glDeleteVertexArrays(1, &cubeVAO);
}

#if defined(SUPPORT_VR_SIMULATOR)
// Set internal projection and modelview matrix depending on eyes tracking data
static void SetStereoView(int eye, Matrix matProjection, Matrix matModelView)
{
    Matrix eyeProjection = matProjection;
    Matrix eyeModelView = matModelView;

    // Setup viewport and projection/modelview matrices using tracking data
    rlViewport(eye*framebufferWidth/2, 0, framebufferWidth/2, framebufferHeight);

    // Apply view offset to modelview matrix
    eyeModelView = MatrixMultiply(matModelView, vrConfig.eyesViewOffset[eye]);

    // Set current eye projection matrix
    eyeProjection = vrConfig.eyesProjection[eye];

    SetMatrixModelview(eyeModelView);
    SetMatrixProjection(eyeProjection);
}
#endif  // SUPPORT_VR_SIMULATOR

#endif  // GRAPHICS_API_OPENGL_33 || GRAPHICS_API_OPENGL_ES2

#if defined(GRAPHICS_API_OPENGL_11)
// Mipmaps data is generated after image data
// NOTE: Only works with RGBA (4 bytes) data!
static int GenerateMipmaps(unsigned char *data, int baseWidth, int baseHeight)
{
    int mipmapCount = 1;                // Required mipmap levels count (including base level)
    int width = baseWidth;
    int height = baseHeight;
    int size = baseWidth*baseHeight*4;  // Size in bytes (will include mipmaps...), RGBA only

    // Count mipmap levels required
    while ((width != 1) && (height != 1))
    {
        if (width != 1) width /= 2;
        if (height != 1) height /= 2;

        TraceLog(LOG_DEBUG, "Next mipmap size: %i x %i", width, height);

        mipmapCount++;

        size += (width*height*4);       // Add mipmap size (in bytes)
    }

    TraceLog(LOG_DEBUG, "Total mipmaps required: %i", mipmapCount);
    TraceLog(LOG_DEBUG, "Total size of data required: %i", size);

    unsigned char *temp = realloc(data, size);

    if (temp != NULL) data = temp;
    else TraceLog(LOG_WARNING, "Mipmaps required memory could not be allocated");

    width = baseWidth;
    height = baseHeight;
    size = (width*height*4);

    // Generate mipmaps
    // NOTE: Every mipmap data is stored after data
    Color *image = (Color *)RL_MALLOC(width*height*sizeof(Color));
    Color *mipmap = NULL;
    int offset = 0;
    int j = 0;

    for (int i = 0; i < size; i += 4)
    {
        image[j].r = data[i];
        image[j].g = data[i + 1];
        image[j].b = data[i + 2];
        image[j].a = data[i + 3];
        j++;
    }

    TraceLog(LOG_DEBUG, "Mipmap base (%ix%i)", width, height);

    for (int mip = 1; mip < mipmapCount; mip++)
    {
        mipmap = GenNextMipmap(image, width, height);

        offset += (width*height*4); // Size of last mipmap
        j = 0;

        width /= 2;
        height /= 2;
        size = (width*height*4);    // Mipmap size to store after offset

        // Add mipmap to data
        for (int i = 0; i < size; i += 4)
        {
            data[offset + i] = mipmap[j].r;
            data[offset + i + 1] = mipmap[j].g;
            data[offset + i + 2] = mipmap[j].b;
            data[offset + i + 3] = mipmap[j].a;
            j++;
        }

        RL_FREE(image);

        image = mipmap;
        mipmap = NULL;
    }

    RL_FREE(mipmap);       // free mipmap data

    return mipmapCount;
}

// Manual mipmap generation (basic scaling algorithm)
static Color *GenNextMipmap(Color *srcData, int srcWidth, int srcHeight)
{
    int x2, y2;
    Color prow, pcol;

    int width = srcWidth/2;
    int height = srcHeight/2;

    Color *mipmap = (Color *)RL_MALLOC(width*height*sizeof(Color));

    // Scaling algorithm works perfectly (box-filter)
    for (int y = 0; y < height; y++)
    {
        y2 = 2*y;

        for (int x = 0; x < width; x++)
        {
            x2 = 2*x;

            prow.r = (srcData[y2*srcWidth + x2].r + srcData[y2*srcWidth + x2 + 1].r)/2;
            prow.g = (srcData[y2*srcWidth + x2].g + srcData[y2*srcWidth + x2 + 1].g)/2;
            prow.b = (srcData[y2*srcWidth + x2].b + srcData[y2*srcWidth + x2 + 1].b)/2;
            prow.a = (srcData[y2*srcWidth + x2].a + srcData[y2*srcWidth + x2 + 1].a)/2;

            pcol.r = (srcData[(y2+1)*srcWidth + x2].r + srcData[(y2+1)*srcWidth + x2 + 1].r)/2;
            pcol.g = (srcData[(y2+1)*srcWidth + x2].g + srcData[(y2+1)*srcWidth + x2 + 1].g)/2;
            pcol.b = (srcData[(y2+1)*srcWidth + x2].b + srcData[(y2+1)*srcWidth + x2 + 1].b)/2;
            pcol.a = (srcData[(y2+1)*srcWidth + x2].a + srcData[(y2+1)*srcWidth + x2 + 1].a)/2;

            mipmap[y*width + x].r = (prow.r + pcol.r)/2;
            mipmap[y*width + x].g = (prow.g + pcol.g)/2;
            mipmap[y*width + x].b = (prow.b + pcol.b)/2;
            mipmap[y*width + x].a = (prow.a + pcol.a)/2;
        }
    }

    TraceLog(LOG_DEBUG, "Mipmap generated successfully (%ix%i)", width, height);

    return mipmap;
}
#endif

#if defined(RLGL_STANDALONE)
// Show trace log messages (LOG_INFO, LOG_WARNING, LOG_ERROR, LOG_DEBUG)
void TraceLog(int msgType, const char *text, ...)
{
    va_list args;
    va_start(args, text);

    switch (msgType)
    {
        case LOG_INFO: fprintf(stdout, "INFO: "); break;
        case LOG_ERROR: fprintf(stdout, "ERROR: "); break;
        case LOG_WARNING: fprintf(stdout, "WARNING: "); break;
        case LOG_DEBUG: fprintf(stdout, "DEBUG: "); break;
        default: break;
    }

    vfprintf(stdout, text, args);
    fprintf(stdout, "\n");

    va_end(args);

    if (msgType == LOG_ERROR) exit(1);
}

// Get pixel data size in bytes (image or texture)
// NOTE: Size depends on pixel format
int GetPixelDataSize(int width, int height, int format)
{
    int dataSize = 0;       // Size in bytes
    int bpp = 0;            // Bits per pixel

    switch (format)
    {
        case UNCOMPRESSED_GRAYSCALE: bpp = 8; break;
        case UNCOMPRESSED_GRAY_ALPHA:
        case UNCOMPRESSED_R5G6B5:
        case UNCOMPRESSED_R5G5B5A1:
        case UNCOMPRESSED_R4G4B4A4: bpp = 16; break;
        case UNCOMPRESSED_R8G8B8A8: bpp = 32; break;
        case UNCOMPRESSED_R8G8B8: bpp = 24; break;
        case UNCOMPRESSED_R32: bpp = 32; break;
        case UNCOMPRESSED_R32G32B32: bpp = 32*3; break;
        case UNCOMPRESSED_R32G32B32A32: bpp = 32*4; break;
        case COMPRESSED_DXT1_RGB:
        case COMPRESSED_DXT1_RGBA:
        case COMPRESSED_ETC1_RGB:
        case COMPRESSED_ETC2_RGB:
        case COMPRESSED_PVRT_RGB:
        case COMPRESSED_PVRT_RGBA: bpp = 4; break;
        case COMPRESSED_DXT3_RGBA:
        case COMPRESSED_DXT5_RGBA:
        case COMPRESSED_ETC2_EAC_RGBA:
        case COMPRESSED_ASTC_4x4_RGBA: bpp = 8; break;
        case COMPRESSED_ASTC_8x8_RGBA: bpp = 2; break;
        default: break;
    }

    dataSize = width*height*bpp/8;  // Total data size in bytes

    // Most compressed formats works on 4x4 blocks,
    // if texture is smaller, minimum dataSize is 8 or 16
    if ((width < 4) && (height < 4))
    {
        if ((format >= COMPRESSED_DXT1_RGB) && (format < COMPRESSED_DXT3_RGBA)) dataSize = 8;
        else if ((format >= COMPRESSED_DXT3_RGBA) && (format < COMPRESSED_ASTC_8x8_RGBA)) dataSize = 16;
    }

    return dataSize;
}
#endif  // RLGL_STANDALONE

#endif  // RLGL_IMPLEMENTATION
