//go:build js && wasm

package rl

import (
	"image/color"
	"syscall/js"
	"unsafe"
)

//go:wasmimport rlgo copyIn
func copyIn(dst uint32, src unsafe.Pointer, n uint32)

//go:wasmimport rlgo copyOut
func copyOut(dst unsafe.Pointer, src uint32, n uint32)

//go:wasmimport rlgo cString
func copyString(dst uint32, src unsafe.Pointer, n uint32)

//go:wasmimport raylib malloc
func cMalloc(n uint32) uint32

//go:wasmimport raylib free
func cFree(p uint32)

//go:wasmimport raylib strlen
func cStrlen(p uint32) uint32

//go:wasmimport raylib rlwRet
func cRet() uint32

const (
	scratchSize           = 64 << 10
	maxMeshVertexBuffers  = 9
	defaultBatchDrawCalls = 256
)

type scratchMark struct {
	off   uint32
	extra int
}

type cScratch struct {
	base  uint32
	off   uint32
	extra []uint32
}

var scratch cScratch

func (s *cScratch) mark() scratchMark {
	return scratchMark{s.off, len(s.extra)}
}

func (s *cScratch) release(m scratchMark) {
	for _, p := range s.extra[m.extra:] {
		cFree(p)
	}
	s.extra = s.extra[:m.extra]
	s.off = m.off
}

func (s *cScratch) alloc(n uint32) uint32 {
	if s.base == 0 {
		s.base = cMalloc(scratchSize)
	}
	n = (n + 15) &^ 15
	if s.off+n > scratchSize {
		p := cMalloc(n)
		s.extra = append(s.extra, p)
		return p
	}
	p := s.base + s.off
	s.off += n
	return p
}

func (s *cScratch) str(v string) uint32 {
	n := uint32(len(v))
	p := s.alloc(n + 1)
	copyString(p, unsafe.Pointer(unsafe.StringData(v)), n)
	return p
}

func (s *cScratch) strNullable(v string) uint32 {
	if v == "" {
		return 0
	}
	return s.str(v)
}

func (s *cScratch) bytes(src unsafe.Pointer, n uint32) uint32 {
	if src == nil || n == 0 {
		return 0
	}
	p := s.alloc(n)
	copyIn(p, src, n)
	return p
}

func scratchSlice[T any](v []T) uint32 {
	if len(v) == 0 {
		return 0
	}
	var t T
	return scratch.bytes(unsafe.Pointer(unsafe.SliceData(v)), uint32(len(v))*uint32(unsafe.Sizeof(t)))
}

func scratchValue[T any](v *T) uint32 {
	return scratch.bytes(unsafe.Pointer(v), uint32(unsafe.Sizeof(*v)))
}

func readC[T any](p uint32) T {
	var v T
	if p != 0 {
		copyOut(unsafe.Pointer(&v), p, uint32(unsafe.Sizeof(v)))
	}
	return v
}

func writeC[T any](p uint32, v *T) {
	copyIn(p, unsafe.Pointer(v), uint32(unsafe.Sizeof(*v)))
}

func readRet[T any]() T {
	var v T
	copyOut(unsafe.Pointer(&v), retAddr, uint32(unsafe.Sizeof(v)))
	return v
}

func scratchPtr[T any](v *T) uint32 {
	if v == nil {
		return 0
	}
	return scratchValue(v)
}

func readBack[T any](v *T, p uint32) {
	if v != nil && p != 0 {
		copyOut(unsafe.Pointer(v), p, uint32(unsafe.Sizeof(*v)))
	}
}

func persist[T any](v *T) uint32 {
	p := cMalloc(uint32(unsafe.Sizeof(*v)))
	writeC(p, v)
	return p
}

func goSlice[T any](p uint32, n int) []T {
	if p == 0 || n <= 0 {
		return nil
	}
	v := make([]T, n)
	var t T
	copyOut(unsafe.Pointer(&v[0]), p, uint32(n)*uint32(unsafe.Sizeof(t)))
	return v
}

func cOwned[T any](v []T) uint32 {
	if len(v) == 0 {
		return 0
	}
	var t T
	n := uint32(len(v)) * uint32(unsafe.Sizeof(t))
	p := cMalloc(n)
	copyIn(p, unsafe.Pointer(unsafe.SliceData(v)), n)
	return p
}

func ptrSlice[T any](p *T, n int) []T {
	if p == nil || n <= 0 {
		return nil
	}
	return unsafe.Slice(p, n)
}

func first[T any](v []T) *T {
	if len(v) == 0 {
		return nil
	}
	return &v[0]
}

func goString(p uint32) string {
	if p == 0 {
		return ""
	}
	n := cStrlen(p)
	if n == 0 {
		return ""
	}
	b := make([]byte, n)
	copyOut(unsafe.Pointer(&b[0]), p, n)
	return unsafe.String(&b[0], n)
}

func colorToUint32(c color.RGBA) uint32 {
	return *(*uint32)(unsafe.Pointer(&c))
}

func uint32ToColor(v uint32) color.RGBA {
	return *(*color.RGBA)(unsafe.Pointer(&v))
}

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func imageDataSize(width, height, mipmaps, format int32) uint32 {
	size := int32(0)
	for i := int32(0); i < max(mipmaps, 1); i++ {
		size += cGetPixelDataSize(width, height, format)
		width = max(width/2, 1)
		height = max(height/2, 1)
	}
	return uint32(size)
}

func imageToC(img *Image) cImage {
	c := cImage{Width: img.Width, Height: img.Height, Mipmaps: img.Mipmaps, Format: int32(img.Format)}
	if img.Data != nil {
		n := imageDataSize(img.Width, img.Height, img.Mipmaps, int32(img.Format))
		c.Data = cMalloc(n)
		copyIn(c.Data, img.Data, n)
	}
	return c
}

func imageFromC(c cImage) Image {
	img := Image{Width: c.Width, Height: c.Height, Mipmaps: c.Mipmaps, Format: PixelFormat(c.Format)}
	if c.Data != 0 {
		n := imageDataSize(c.Width, c.Height, c.Mipmaps, c.Format)
		data := make([]byte, n)
		copyOut(unsafe.Pointer(&data[0]), c.Data, n)
		img.Data = unsafe.Pointer(&data[0])
		cFree(c.Data)
	}
	return img
}

func imagePixel(img *Image, x, y int32) (unsafe.Pointer, uint32) {
	if img.Data == nil || x < 0 || x >= img.Width || y < 0 || y >= img.Height || img.Format >= CompressedDxt1Rgb {
		return nil, 0
	}
	n := cGetPixelDataSize(1, 1, int32(img.Format))
	return unsafe.Add(img.Data, (y*img.Width+x)*n), uint32(n)
}

func imageIn(img *Image) uint32 {
	c := imageToC(img)
	return scratchValue(&c)
}

func imageDone(p uint32) {
	cFree(readC[cImage](p).Data)
}

func imageBack(img *Image, p uint32) {
	*img = imageFromC(readC[cImage](p))
}

func imageRet() *Image {
	img := imageFromC(readRet[cImage]())
	return &img
}

func waveDataSize(c *cWave) uint32 {
	return c.FrameCount * c.Channels * c.SampleSize / 8
}

func waveToC(w *Wave) cWave {
	c := cWave{FrameCount: w.FrameCount, SampleRate: w.SampleRate, SampleSize: w.SampleSize, Channels: w.Channels}
	if w.Data != nil {
		n := waveDataSize(&c)
		c.Data = cMalloc(n)
		copyIn(c.Data, w.Data, n)
	}
	return c
}

func waveFromC(c cWave) Wave {
	w := Wave{FrameCount: c.FrameCount, SampleRate: c.SampleRate, SampleSize: c.SampleSize, Channels: c.Channels}
	if c.Data != 0 {
		n := waveDataSize(&c)
		data := make([]byte, n)
		copyOut(unsafe.Pointer(&data[0]), c.Data, n)
		w.Data = unsafe.Pointer(&data[0])
		cFree(c.Data)
	}
	return w
}

func waveIn(w *Wave) uint32 {
	c := waveToC(w)
	return scratchValue(&c)
}

func waveDone(p uint32) {
	cFree(readC[cWave](p).Data)
}

func waveBack(w *Wave, p uint32) {
	*w = waveFromC(readC[cWave](p))
}

type shaderEntry struct {
	at   uint32
	locs [MaxShaderLocations]int32
}

var shaderLocs = map[uint32]*shaderEntry{}

func shaderFromC(c cShader) Shader {
	s := Shader{ID: c.Id}
	if c.Locs != 0 {
		e := &shaderEntry{at: c.Locs}
		copyOut(unsafe.Pointer(&e.locs[0]), c.Locs, MaxShaderLocations*4)
		shaderLocs[c.Id] = e
		locs := e.locs
		s.Locs = &locs[0]
	}
	return s
}

func shaderToC(s Shader) cShader {
	e := shaderLocs[s.ID]
	if s.Locs != nil {
		locs := (*[MaxShaderLocations]int32)(unsafe.Pointer(s.Locs))
		fresh := e == nil
		if fresh {
			e = &shaderEntry{at: cMalloc(MaxShaderLocations * 4)}
			shaderLocs[s.ID] = e
		}
		if fresh || e.locs != *locs {
			e.locs = *locs
			copyIn(e.at, unsafe.Pointer(&e.locs[0]), MaxShaderLocations*4)
		}
	}
	if e == nil {
		return cShader{Id: s.ID}
	}
	return cShader{Id: s.ID, Locs: e.at}
}

func shaderIn(s Shader) uint32 {
	c := shaderToC(s)
	return scratchValue(&c)
}

func materialFromC(c cMaterial) Material {
	return Material{Shader: shaderFromC(c.Shader), Maps: first(goSlice[MaterialMap](c.Maps, MaxMaterialMaps)), Params: c.Params}
}

type materialMaps = [MaxMaterialMaps]MaterialMap

func materialToC(m *Material, maps uint32) cMaterial {
	if m.Maps != nil {
		if maps == 0 {
			maps = scratch.alloc(uint32(unsafe.Sizeof(materialMaps{})))
		}
		copyIn(maps, unsafe.Pointer(m.Maps), uint32(unsafe.Sizeof(materialMaps{})))
	}
	return cMaterial{Shader: shaderToC(m.Shader), Maps: maps, Params: m.Params}
}

func materialRet() Material {
	c := readRet[cMaterial]()
	m := materialFromC(c)
	cFree(c.Maps)
	return m
}

func materialIn(m *Material) uint32 {
	c := materialToC(m, 0)
	return scratchValue(&c)
}

type fontEntry struct {
	at   uint32
	c    cFont
	font Font
}

var fonts = map[uint32]*fontEntry{}

func glyphFromC(c cGlyphInfo, copyImage bool) GlyphInfo {
	g := GlyphInfo{Value: c.Value, OffsetX: c.OffsetX, OffsetY: c.OffsetY, AdvanceX: c.AdvanceX}
	if copyImage {
		img := c.Image
		if img.Data != 0 {
			n := imageDataSize(img.Width, img.Height, img.Mipmaps, img.Format)
			data := make([]byte, n)
			copyOut(unsafe.Pointer(&data[0]), img.Data, n)
			g.Image = Image{Data: unsafe.Pointer(&data[0]), Width: img.Width, Height: img.Height, Mipmaps: img.Mipmaps, Format: PixelFormat(img.Format)}
		}
	}
	return g
}

func glyphsFromC(p uint32, n int) []GlyphInfo {
	cs := goSlice[cGlyphInfo](p, n)
	gs := make([]GlyphInfo, len(cs))
	for i := range cs {
		gs[i] = glyphFromC(cs[i], true)
	}
	return gs
}

func glyphsToC(gs []GlyphInfo) uint32 {
	if len(gs) == 0 {
		return 0
	}
	cs := make([]cGlyphInfo, len(gs))
	for i, g := range gs {
		img := &gs[i].Image
		cs[i] = cGlyphInfo{Value: g.Value, OffsetX: g.OffsetX, OffsetY: g.OffsetY, AdvanceX: g.AdvanceX,
			Image: cImage{Width: img.Width, Height: img.Height, Mipmaps: img.Mipmaps, Format: int32(img.Format)}}
		if img.Data != nil {
			cs[i].Image.Data = scratch.bytes(img.Data, imageDataSize(img.Width, img.Height, img.Mipmaps, int32(img.Format)))
		}
	}
	return scratchSlice(cs)
}

func fontFromC(c cFont) Font {
	e := fonts[c.Texture.ID]
	if e != nil && e.c == c {
		return e.font
	}
	n := int(c.GlyphCount)
	f := Font{BaseSize: c.BaseSize, CharsCount: c.GlyphCount, CharsPadding: c.GlyphPadding, Texture: c.Texture,
		Recs: first(goSlice[Rectangle](c.Recs, n)), Chars: first(glyphsFromC(c.Glyphs, n))}
	if c.Texture.ID != 0 {
		if e != nil {
			cFree(e.at)
		}
		fonts[c.Texture.ID] = &fontEntry{at: persist(&c), c: c, font: f}
	}
	return f
}

func fontIn(f *Font) uint32 {
	if e := fonts[f.Texture.ID]; e != nil {
		c := e.c
		c.BaseSize, c.GlyphCount, c.GlyphPadding, c.Texture = f.BaseSize, f.CharsCount, f.CharsPadding, f.Texture
		if c != e.c {
			e.c = c
			writeC(e.at, &c)
		}
		return e.at
	}
	n := int(f.CharsCount)
	c := cFont{BaseSize: f.BaseSize, GlyphCount: f.CharsCount, GlyphPadding: f.CharsPadding, Texture: f.Texture,
		Recs: scratchSlice(ptrSlice(f.Recs, n)), Glyphs: glyphsToC(ptrSlice(f.Chars, n))}
	return scratchValue(&c)
}

type meshEntry struct {
	at    uint32
	owned bool
}

var meshes = map[*uint32]meshEntry{}

func meshFromC(c cMesh, at uint32, owned bool) Mesh {
	v, t := int(c.VertexCount), int(c.TriangleCount)
	m := Mesh{
		VertexCount:   c.VertexCount,
		TriangleCount: c.TriangleCount,
		Vertices:      first(goSlice[float32](c.Vertices, v*3)),
		Texcoords:     first(goSlice[float32](c.Texcoords, v*2)),
		Texcoords2:    first(goSlice[float32](c.Texcoords2, v*2)),
		Normals:       first(goSlice[float32](c.Normals, v*3)),
		Tangents:      first(goSlice[float32](c.Tangents, v*4)),
		Colors:        first(goSlice[uint8](c.Colors, v*4)),
		Indices:       first(goSlice[uint16](c.Indices, t*3)),
		BoneCount:     uint32(c.BoneCount),
		BoneIndices:   first(goSlice[uint8](c.BoneIndices, v*4)),
		BoneWeights:   first(goSlice[float32](c.BoneWeights, v*4)),
		AnimVertices:  first(goSlice[float32](c.AnimVertices, v*3)),
		AnimNormals:   first(goSlice[float32](c.AnimNormals, v*3)),
		VaoID:         c.VaoId,
		VboID:         first(goSlice[uint32](c.VboId, maxMeshVertexBuffers)),
	}
	if m.VboID != nil {
		meshes[m.VboID] = meshEntry{at, owned}
	} else if owned {
		cFree(at)
	}
	return m
}

func meshToC(m *Mesh, owned bool) cMesh {
	v, t := int(m.VertexCount), int(m.TriangleCount)
	data := func(p unsafe.Pointer, n uint32) uint32 {
		if p == nil {
			return 0
		}
		if owned {
			q := cMalloc(n)
			copyIn(q, p, n)
			return q
		}
		return scratch.bytes(p, n)
	}
	f := uint32(v) * 4
	c := cMesh{
		VertexCount:   m.VertexCount,
		TriangleCount: m.TriangleCount,
		Vertices:      data(unsafe.Pointer(m.Vertices), f*3),
		Texcoords:     data(unsafe.Pointer(m.Texcoords), f*2),
		Texcoords2:    data(unsafe.Pointer(m.Texcoords2), f*2),
		Normals:       data(unsafe.Pointer(m.Normals), f*3),
		Tangents:      data(unsafe.Pointer(m.Tangents), f*4),
		Colors:        data(unsafe.Pointer(m.Colors), uint32(v)*4),
		Indices:       data(unsafe.Pointer(m.Indices), uint32(t)*6),
		BoneCount:     int32(m.BoneCount),
		BoneIndices:   data(unsafe.Pointer(m.BoneIndices), uint32(v)*4),
		BoneWeights:   data(unsafe.Pointer(m.BoneWeights), f*4),
		AnimVertices:  data(unsafe.Pointer(m.AnimVertices), f*3),
		AnimNormals:   data(unsafe.Pointer(m.AnimNormals), f*3),
		VaoId:         m.VaoID,
		VboId:         data(unsafe.Pointer(m.VboID), maxMeshVertexBuffers*4),
	}
	return c
}

func meshRet() Mesh {
	c := readRet[cMesh]()
	return meshFromC(c, persist(&c), true)
}

func meshIn(m *Mesh) uint32 {
	if e, ok := meshes[m.VboID]; ok && m.VboID != nil {
		return e.at
	}
	c := meshToC(m, false)
	return scratchValue(&c)
}

type modelEntry struct {
	at   uint32
	c    cModel
	mats []cMaterial
	maps []materialMaps
	mesh []int32
}

var models = map[*Mesh]*modelEntry{}

func modelFromC(c cModel) Model {
	m := Model{Transform: c.Transform, MeshCount: c.MeshCount, MaterialCount: c.MaterialCount}
	ms := make([]Mesh, c.MeshCount)
	meshSize := uint32(unsafe.Sizeof(cMesh{}))
	for i := range ms {
		at := c.Meshes + uint32(i)*meshSize
		ms[i] = meshFromC(readC[cMesh](at), at, false)
	}
	cms := goSlice[cMaterial](c.Materials, int(c.MaterialCount))
	mats := make([]Material, len(cms))
	for i := range cms {
		mats[i] = materialFromC(cms[i])
	}
	bones := int(c.Skeleton.BoneCount)
	m.Meshes = first(ms)
	m.Materials = first(mats)
	m.MeshMaterial = first(goSlice[int32](c.MeshMaterial, int(c.MeshCount)))
	m.Skeleton = ModelSkeleton{BoneCount: c.Skeleton.BoneCount, Bones: first(goSlice[BoneInfo](c.Skeleton.Bones, bones)),
		BindPose: first(goSlice[Transform](c.Skeleton.BindPose, bones))}
	m.CurrentPose = first(goSlice[Transform](c.CurrentPose, bones))
	m.BoneMatrices = first(goSlice[Matrix](c.BoneMatrices, bones))
	if m.Meshes != nil {
		e := &modelEntry{at: persist(&c), c: c, mats: cms, maps: make([]materialMaps, len(mats)), mesh: goSlice[int32](c.MeshMaterial, int(c.MeshCount))}
		for i := range mats {
			if mats[i].Maps != nil {
				e.maps[i] = *(*materialMaps)(unsafe.Pointer(mats[i].Maps))
			}
		}
		models[m.Meshes] = e
	}
	return m
}

func modelIn(m *Model) uint32 {
	e := models[m.Meshes]
	if e == nil || m.Meshes == nil {
		return modelTemp(m)
	}
	c := e.c
	c.Transform = m.Transform
	matSize := uint32(unsafe.Sizeof(cMaterial{}))
	for i, mat := range ptrSlice(m.Materials, min(int(m.MaterialCount), len(e.mats))) {
		if mat.Maps != nil && e.mats[i].Maps != 0 {
			if maps := (*materialMaps)(unsafe.Pointer(mat.Maps)); *maps != e.maps[i] {
				e.maps[i] = *maps
				copyIn(e.mats[i].Maps, unsafe.Pointer(&e.maps[i]), uint32(unsafe.Sizeof(materialMaps{})))
			}
		}
		cm := cMaterial{Shader: shaderToC(mat.Shader), Maps: e.mats[i].Maps, Params: mat.Params}
		if cm != e.mats[i] {
			e.mats[i] = cm
			writeC(c.Materials+uint32(i)*matSize, &cm)
		}
	}
	if mm := ptrSlice(m.MeshMaterial, len(e.mesh)); mm != nil && c.MeshMaterial != 0 {
		for i := range mm {
			if mm[i] != e.mesh[i] {
				copy(e.mesh, mm)
				copyIn(c.MeshMaterial, unsafe.Pointer(&e.mesh[0]), uint32(len(e.mesh))*4)
				break
			}
		}
	}
	if c != e.c {
		e.c = c
		writeC(e.at, &c)
	}
	return e.at
}

func modelTemp(m *Model) uint32 {
	c := cModel{Transform: m.Transform, MeshCount: m.MeshCount, MaterialCount: m.MaterialCount}
	ms := ptrSlice(m.Meshes, int(m.MeshCount))
	cms := make([]cMesh, len(ms))
	for i := range ms {
		if e, ok := meshes[ms[i].VboID]; ok && ms[i].VboID != nil {
			cms[i] = readC[cMesh](e.at)
		} else {
			cms[i] = meshToC(&ms[i], false)
		}
	}
	c.Meshes = scratchSlice(cms)
	mats := ptrSlice(m.Materials, int(m.MaterialCount))
	cmats := make([]cMaterial, len(mats))
	for i := range mats {
		cmats[i] = materialToC(&mats[i], 0)
	}
	c.Materials = scratchSlice(cmats)
	c.MeshMaterial = scratchSlice(ptrSlice(m.MeshMaterial, int(m.MeshCount)))
	bones := int(m.Skeleton.BoneCount)
	c.Skeleton = cModelSkeleton{BoneCount: m.Skeleton.BoneCount, Bones: scratchSlice(ptrSlice(m.Skeleton.Bones, bones)),
		BindPose: scratchSlice(ptrSlice(m.Skeleton.BindPose, bones))}
	c.CurrentPose = scratchSlice(ptrSlice(m.CurrentPose, bones))
	c.BoneMatrices = scratchSlice(ptrSlice(m.BoneMatrices, bones))
	return scratchValue(&c)
}

func modelPoseBack(m *Model) {
	e := models[m.Meshes]
	if e == nil || m.Meshes == nil {
		return
	}
	bones := uint32(e.c.Skeleton.BoneCount)
	if m.CurrentPose != nil && e.c.CurrentPose != 0 {
		copyOut(unsafe.Pointer(m.CurrentPose), e.c.CurrentPose, bones*uint32(unsafe.Sizeof(Transform{})))
	}
	if m.BoneMatrices != nil && e.c.BoneMatrices != 0 {
		copyOut(unsafe.Pointer(m.BoneMatrices), e.c.BoneMatrices, bones*uint32(unsafe.Sizeof(Matrix{})))
	}
}

func forgetModel(m *Model) {
	e := models[m.Meshes]
	if e == nil {
		return
	}
	for _, mesh := range ptrSlice(m.Meshes, int(m.MeshCount)) {
		delete(meshes, mesh.VboID)
	}
	cFree(e.at)
	delete(models, m.Meshes)
}

type animEntry struct {
	at    uint32
	base  uint32
	count int32
}

var anims = map[*ModelAnimPose]animEntry{}

func animsFromC(p uint32, count int32) []ModelAnimation {
	animSize := uint32(unsafe.Sizeof(cModelAnimation{}))
	out := make([]ModelAnimation, count)
	for i := range out {
		at := p + uint32(i)*animSize
		c := readC[cModelAnimation](at)
		frames := goSlice[uint32](c.KeyframePoses, int(c.KeyframeCount))
		poses := make([]ModelAnimPose, len(frames))
		for k, f := range frames {
			poses[k] = first(goSlice[Transform](f, int(c.BoneCount)))
		}
		out[i] = ModelAnimation{Name: c.Name, BoneCount: c.BoneCount, KeyframeCount: c.KeyframeCount, KeyframePoses: first(poses)}
		if out[i].KeyframePoses != nil {
			anims[out[i].KeyframePoses] = animEntry{at, p, count}
		}
	}
	return out
}

func animIn(a *ModelAnimation) uint32 {
	if e, ok := anims[a.KeyframePoses]; ok && a.KeyframePoses != nil {
		return e.at
	}
	poses := ptrSlice(a.KeyframePoses, int(a.KeyframeCount))
	frames := make([]uint32, len(poses))
	for i, p := range poses {
		frames[i] = scratchSlice(ptrSlice(p, int(a.BoneCount)))
	}
	c := cModelAnimation{Name: a.Name, BoneCount: a.BoneCount, KeyframeCount: a.KeyframeCount, KeyframePoses: scratchSlice(frames)}
	return scratchValue(&c)
}

var (
	cHandles  = map[uint32]unsafe.Pointer{}
	goHandles = map[unsafe.Pointer]uint32{}
)

func handleFor[T any](c uint32) *T {
	if c == 0 {
		return nil
	}
	if p, ok := cHandles[c]; ok {
		return (*T)(p)
	}
	p := unsafe.Pointer(new(T))
	cHandles[c] = p
	goHandles[p] = c
	return (*T)(p)
}

func handleC[T any](p *T) uint32 {
	return goHandles[unsafe.Pointer(p)]
}

func dropHandle[T any](p *T) {
	if c, ok := goHandles[unsafe.Pointer(p)]; ok {
		delete(cHandles, c)
		delete(goHandles, unsafe.Pointer(p))
	}
}

func audioStreamFromC(c cAudioStream) AudioStream {
	return AudioStream{Buffer: handleFor[AudioBuffer](c.Buffer), Processor: handleFor[AudioProcessor](c.Processor),
		SampleRate: c.SampleRate, SampleSize: c.SampleSize, Channels: c.Channels}
}

func audioStreamToC(s *AudioStream) cAudioStream {
	return cAudioStream{Buffer: handleC(s.Buffer), Processor: handleC(s.Processor), SampleRate: s.SampleRate, SampleSize: s.SampleSize, Channels: s.Channels}
}

func forgetAudioStream(s *AudioStream) {
	forgetStreamCallbacks(s.Buffer)
	dropHandle(s.Buffer)
	dropHandle(s.Processor)
}

func audioStreamIn(s *AudioStream) uint32 {
	c := audioStreamToC(s)
	return scratchValue(&c)
}

func soundFromC(c cSound) Sound {
	return Sound{Stream: audioStreamFromC(c.Stream), FrameCount: c.FrameCount}
}

func soundIn(s *Sound) uint32 {
	c := cSound{Stream: audioStreamToC(&s.Stream), FrameCount: s.FrameCount}
	return scratchValue(&c)
}

var musicData = map[unsafe.Pointer]uint32{}

func musicFromC(c cMusic) Music {
	return Music{Stream: audioStreamFromC(c.Stream), FrameCount: c.FrameCount, Looping: c.Looping, CtxType: c.CtxType,
		CtxData: unsafe.Pointer(handleFor[byte](c.CtxData))}
}

func musicIn(m *Music) uint32 {
	c := cMusic{Stream: audioStreamToC(&m.Stream), FrameCount: m.FrameCount, Looping: m.Looping, CtxType: m.CtxType,
		CtxData: handleC((*byte)(m.CtxData))}
	return scratchValue(&c)
}

type automationEntry struct {
	at uint32
	c  cAutomationEventList
}

var (
	automationLists   = map[*AutomationEvent]*automationEntry{}
	currentAutomation *AutomationEventList
)

func automationListFromC(c cAutomationEventList) AutomationEventList {
	events := goSlice[AutomationEvent](c.Events, int(c.Capacity))
	l := AutomationEventList{Capacity: c.Capacity, Count: c.Count, Events: first(events)}
	if l.Events != nil {
		automationLists[l.Events] = &automationEntry{at: persist(&c), c: c}
	}
	return l
}

func automationListIn(l *AutomationEventList) uint32 {
	if e := automationLists[l.Events]; e != nil && l.Events != nil {
		if l.Count != e.c.Count && l.Count <= e.c.Capacity {
			copyIn(e.c.Events, unsafe.Pointer(l.Events), l.Count*uint32(unsafe.Sizeof(AutomationEvent{})))
			e.c.Count = l.Count
			writeC(e.at, &e.c)
		}
		return e.at
	}
	c := cAutomationEventList{Capacity: l.Capacity, Count: l.Count, Events: scratchSlice(ptrSlice(l.Events, int(l.Capacity)))}
	return scratchValue(&c)
}

func syncAutomationList() {
	l := currentAutomation
	if l == nil {
		return
	}
	e := automationLists[l.Events]
	if e == nil {
		return
	}
	e.c = readC[cAutomationEventList](e.at)
	l.Count = e.c.Count
	if l.Count > 0 {
		copyOut(unsafe.Pointer(l.Events), e.c.Events, l.Count*uint32(unsafe.Sizeof(AutomationEvent{})))
	}
}

//go:wasmimport raylib rlwTraceLog
func cTraceLogText(logLevel int32, text uint32)

//go:wasmimport raylib rlwSetTraceLogCallbackEnabled
func cSetTraceLogCallbackEnabled(enabled int32)

//go:wasmimport raylib rlwAttachAudioMixedProcessorSlot
func cAttachAudioMixedProcessorSlot()

//go:wasmimport raylib rlwDetachAudioMixedProcessorSlot
func cDetachAudioMixedProcessorSlot()

//go:wasmimport raylib rlwAttachAudioStreamProcessorSlot
func cAttachAudioStreamProcessorSlot(stream uint32, slot int32)

//go:wasmimport raylib rlwDetachAudioStreamProcessorSlot
func cDetachAudioStreamProcessorSlot(stream uint32, slot int32)

//go:wasmimport raylib rlwSetAudioStreamCallbackSlot
func cSetAudioStreamCallbackSlot(stream uint32, slot int32)

const streamCallbackBase = 1000

type streamCallback struct {
	buffer   *AudioBuffer
	channels int
	fn       AudioCallback
}

var (
	traceLogCallback TraceLogCallbackFun
	streamCallbacks  [audioProcessorSlots]streamCallback
	audioScratch     []float32
)

func init() {
	hooks := js.Global().Get("rlgo")
	hooks.Set("traceLog", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if traceLogCallback != nil {
			traceLogCallback(args[0].Int(), goString(uint32(args[1].Int())))
		}
		return nil
	}))
	hooks.Set("audio", js.FuncOf(func(_ js.Value, args []js.Value) any {
		audioCallback(args[0].Int(), uint32(args[1].Int()), args[2].Int())
		return nil
	}))
}

func audioCallback(slot int, data uint32, frames int) {
	channels := audioDeviceChannels
	var cb *streamCallback
	if slot >= streamCallbackBase {
		cb = &streamCallbacks[slot-streamCallbackBase]
		if cb.fn == nil {
			return
		}
		channels = cb.channels
	}
	n := frames * channels
	if n <= 0 {
		return
	}
	if len(audioScratch) < n {
		audioScratch = make([]float32, n)
	}
	buf := audioScratch[:n]
	size := uint32(n) * 4
	copyOut(unsafe.Pointer(&buf[0]), data, size)
	switch {
	case cb != nil:
		cb.fn(buf, frames)
	case slot < 0:
		callAudioMixedProcessors(unsafe.Pointer(&buf[0]), frames)
	default:
		callAudioStreamProcessor(slot, unsafe.Pointer(&buf[0]), frames)
	}
	copyIn(data, unsafe.Pointer(&buf[0]), size)
}

func streamCallbackSlot(stream *AudioStream, fn AudioCallback) int32 {
	free := -1
	for i := range streamCallbacks {
		if streamCallbacks[i].buffer == stream.Buffer && streamCallbacks[i].fn != nil {
			free = i
			break
		}
		if free < 0 && streamCallbacks[i].fn == nil {
			free = i
		}
	}
	if fn == nil {
		if free >= 0 && streamCallbacks[free].buffer == stream.Buffer {
			streamCallbacks[free] = streamCallback{}
		}
		return -1
	}
	if free < 0 {
		TraceLog(LogWarning, "AUDIO: All %d stream callback slots are in use", audioProcessorSlots)
		return -2
	}
	streamCallbacks[free] = streamCallback{stream.Buffer, max(int(stream.Channels), 1), fn}
	return int32(free)
}

func forgetStreamCallbacks(buffer *AudioBuffer) {
	for i := range streamCallbacks {
		if streamCallbacks[i].buffer == buffer {
			streamCallbacks[i] = streamCallback{}
		}
	}
}

var renderBatches = map[*VertexBuffer]uint32{}

func renderBatchFromC(c cRlRenderBatch) RenderBatch {
	cvbs := goSlice[cRlVertexBuffer](c.VertexBuffer, int(c.BufferCount))
	vbs := make([]VertexBuffer, len(cvbs))
	for i, v := range cvbs {
		vbs[i] = VertexBuffer{ElementCount: v.ElementCount, VaoId: v.VaoId, VboId: v.VboId}
	}
	b := RenderBatch{BufferCount: c.BufferCount, CurrentBuffer: c.CurrentBuffer, VertexBuffer: first(vbs),
		Draws: first(goSlice[DrawCall](c.Draws, defaultBatchDrawCalls)), DrawCounter: c.DrawCounter, CurrentDepth: c.CurrentDepth}
	if b.VertexBuffer != nil {
		renderBatches[b.VertexBuffer] = persist(&c)
	}
	return b
}

func renderBatchC(b *RenderBatch) uint32 {
	if b == nil {
		return 0
	}
	return renderBatches[b.VertexBuffer]
}

var (
	retAddr = cRet()

	frameChan = make(chan float64, 1)
	frameFunc = js.FuncOf(func(_ js.Value, args []js.Value) any {
		select {
		case frameChan <- args[0].Float():
		default:
		}
		return nil
	})
	requestAnimationFrame = js.Global().Get("requestAnimationFrame")

	targetFPS int32
	lastFrame float64
)

func waitFrame() {
	for {
		requestAnimationFrame.Invoke(frameFunc)
		t := <-frameChan
		if targetFPS <= 0 {
			lastFrame = t
			return
		}
		interval := 1000 / float64(targetFPS)
		if t >= lastFrame+interval-1 {
			lastFrame += interval
			if t-lastFrame > interval {
				lastFrame = t
			}
			return
		}
	}
}
