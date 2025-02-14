package rl

const (
	// Texture parameters (equivalent to OpenGL defines)
	TextureWrapS     = 0x2802 // GL_TEXTURE_WRAP_S
	TextureWrapT     = 0x2803 // GL_TEXTURE_WRAP_T
	TextureMagFilter = 0x2800 // GL_TEXTURE_MAG_FILTER
	TextureMinFilter = 0x2801 // GL_TEXTURE_MIN_FILTER

	TextureFilterNearest          = 0x2600 // GL_NEAREST
	TextureFilterLinear           = 0x2601 // GL_LINEAR
	TextureFilterMipNearest       = 0x2700 // GL_NEAREST_MIPMAP_NEAREST
	TextureFilterNearestMipLinear = 0x2702 // GL_NEAREST_MIPMAP_LINEAR
	TextureFilterLinearMipNearest = 0x2701 // GL_LINEAR_MIPMAP_NEAREST
	TextureFilterMipLinear        = 0x2703 // GL_LINEAR_MIPMAP_LINEAR
	TextureFilterAnisotropic      = 0x3000 // Anisotropic filter (custom identifier)
	TextureMipmapBiasRatio        = 0x4000 // Texture mipmap bias, percentage ratio (custom identifier)

	TextureWrapRepeat       = 0x2901 // GL_REPEAT
	TextureWrapClamp        = 0x812F // GL_CLAMP_TO_EDGE
	TextureWrapMirrorRepeat = 0x8370 // GL_MIRRORED_REPEAT
	TextureWrapMirrorClamp  = 0x8742 // GL_MIRROR_CLAMP_EXT

	// Matrix modes (equivalent to OpenGL)
	Modelview  = 0x1700 // GL_MODELVIEW
	Projection = 0x1701 // GL_PROJECTION
	Texture    = 0x1702 // GL_TEXTURE

	// Primitive assembly draw modes
	Lines     = 0x0001 // GL_LINES
	Triangles = 0x0004 // GL_TRIANGLES
	Quads     = 0x0007 // GL_QUADS

	// GL equivalent data types
	UnsignedByte = 0x1401 // GL_UNSIGNED_BYTE
	Float        = 0x1406 // GL_FLOAT

	// Buffer usage hint
	StreamDraw  = 0x88E0 // GL_STREAM_DRAW
	StreamRead  = 0x88E1 // GL_STREAM_READ
	StreamCopy  = 0x88E2 // GL_STREAM_COPY
	StaticDraw  = 0x88E4 // GL_STATIC_DRAW
	StaticRead  = 0x88E5 // GL_STATIC_READ
	StaticCopy  = 0x88E6 // GL_STATIC_COPY
	DynamicDraw = 0x88E8 // GL_DYNAMIC_DRAW
	DynamicRead = 0x88E9 // GL_DYNAMIC_READ
	DynamicCopy = 0x88EA // GL_DYNAMIC_COPY

	// GL Shader type
	FragmentShader = 0x8B30 // GL_FRAGMENT_SHADER
	VertexShader   = 0x8B31 // GL_VERTEX_SHADER
	ComputeShader  = 0x91B9 // GL_COMPUTE_SHADER

	// GL blending factors
	Zero                  = 0      // GL_ZERO
	One                   = 1      // GL_ONE
	SrcColor              = 0x0300 // GL_SRC_COLOR
	OneMinusSrcColor      = 0x0301 // GL_ONE_MINUS_SRC_COLOR
	SrcAlpha              = 0x0302 // GL_SRC_ALPHA
	OneMinusSrcAlpha      = 0x0303 // GL_ONE_MINUS_SRC_ALPHA
	DstAlpha              = 0x0304 // GL_DST_ALPHA
	OneMinusDstAlpha      = 0x0305 // GL_ONE_MINUS_DST_ALPHA
	DstColor              = 0x0306 // GL_DST_COLOR
	OneMinusDstColor      = 0x0307 // GL_ONE_MINUS_DST_COLOR
	SrcAlphaSaturate      = 0x0308 // GL_SRC_ALPHA_SATURATE
	ConstantColor         = 0x8001 // GL_CONSTANT_COLOR
	OneMinusConstantColor = 0x8002 // GL_ONE_MINUS_CONSTANT_COLOR
	ConstantAlpha         = 0x8003 // GL_CONSTANT_ALPHA
	OneMinusConstantAlpha = 0x8004 // GL_ONE_MINUS_CONSTANT_ALPHA

	// GL blending functions/equations
	FuncAdd             = 0x8006        // GL_FUNC_ADD
	Min                 = 0x8007        // GL_MIN
	Max                 = 0x8008        // GL_MAX
	FuncSubtract        = 0x800A        // GL_FUNC_SUBTRACT
	FuncReverseSubtract = 0x800B        // GL_FUNC_REVERSE_SUBTRACT
	BlendEquation       = 0x8009        // GL_BLEND_EQUATION
	BlendEquationRgb    = BlendEquation // GL_BLEND_EQUATION_RGB (Same as BLEND_EQUATION)
	BlendEquationAlpha  = 0x883D        // GL_BLEND_EQUATION_ALPHA
	BlendDstRgb         = 0x80C8        // GL_BLEND_DST_RGB
	BlendSrcRgb         = 0x80C9        // GL_BLEND_SRC_RGB
	BlendDstAlpha       = 0x80CA        // GL_BLEND_DST_ALPHA
	BlendSrcAlpha       = 0x80CB        // GL_BLEND_SRC_ALPHA
	BlendColor          = 0x8005        // GL_BLEND_COLOR

	ReadFramebuffer = 0x8CA8 // GL_READ_FRAMEBUFFER
	DrawFramebuffer = 0x8CA9 // GL_DRAW_FRAMEBUFFER
)

// VertexBuffer - Dynamic vertex buffers (position + texcoords + colors + indices arrays)
type VertexBuffer struct {
	ElementCount int32     // Number of elements in the buffer (QUADS)
	Vertices     *float32  // Vertex position (XYZ - 3 components per vertex) (shader-location = 0)
	Texcoords    *float32  // Vertex texture coordinates (UV - 2 components per vertex) (shader-location = 1)
	Normals      *float32  // Vertex normal (XYZ - 3 components per vertex) (shader-location = 2)
	Colors       *uint8    // Vertex colors (RGBA - 4 components per vertex) (shader-location = 3)
	Indices      *uint32   // Vertex indices (in case vertex data comes indexed) (6 indices per quad)
	VaoId        uint32    // OpenGL Vertex Array Object id
	VboId        [5]uint32 // OpenGL Vertex Buffer Objects id (5 types of vertex data)
}

// DrawCall - Draw call type
// NOTE: Only texture changes register a new draw, other state-change-related elements are not
// used at this moment (vaoId, shaderId, matrices), raylib just forces a batch draw call if any
// of those state-change happens (this is done in core module)
type DrawCall struct {
	Mode            int32
	VertexCount     int32
	VertexAlignment int32
	TextureId       uint32
}

// RenderBatch type
type RenderBatch struct {
	BufferCount   int32
	CurrentBuffer int32
	VertexBuffer  *VertexBuffer
	Draws         *DrawCall
	DrawCounter   int32
	DurrentDepth  float32
}

// OpenGL version
const (
	Opengl11   int32 = 1
	Opengl21   int32 = 2
	Opengl33   int32 = 3
	Opengl43   int32 = 4
	OpenglEs20 int32 = 5
)

// GlVersion type
type GlVersion = int32

// Shader attribute data types
const (
	ShaderAttribFloat int32 = 0
	ShaderAttribVec2  int32 = 1
	ShaderAttribVec3  int32 = 2
	ShaderAttribVec4  int32 = 3
)

// ShaderAttributeDataType type
type ShaderAttributeDataType = int32

// Framebuffer attachment type
// NOTE: By default up to 8 color channels defined but it can be more
const (
	AttachmentColorChannel0 int32 = 0
	AttachmentColorChannel1 int32 = 1
	AttachmentColorChannel2 int32 = 2
	AttachmentColorChannel3 int32 = 3
	AttachmentColorChannel4 int32 = 4
	AttachmentColorChannel5 int32 = 5
	AttachmentColorChannel6 int32 = 6
	AttachmentColorChannel7 int32 = 7
	AttachmentDepth         int32 = 100
	AttachmentStencil       int32 = 200
)

// FramebufferAttachType type
type FramebufferAttachType = int32

// Framebuffer texture attachment type
const (
	AttachmentCubemapPositiveX int32 = 0
	AttachmentCubemapNegativeX int32 = 1
	AttachmentCubemapPositiveY int32 = 2
	AttachmentCubemapNegativeY int32 = 3
	AttachmentCubemapPositiveZ int32 = 4
	AttachmentCubemapNegativeZ int32 = 5
	AttachmentTexture2d        int32 = 100
	AttachmentRenderbuffer     int32 = 200
)

// FramebufferAttachTextureType type
type FramebufferAttachTextureType = int32
