package raylib

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Shader location point type
const (
	LocVertexPosition = iota
	LocVertexTexcoord01
	LocVertexTexcoord02
	LocVertexNormal
	LocVertexTangent
	LocVertexColor
	LocMatrixMvp
	LocMatrixModel
	LocMatrixView
	LocMatrixProjection
	LocVectorView
	LocColorDiffuse
	LocColorSpecular
	LocColorAmbient
	LocMapAlbedo
	LocMapMetalness
	LocMapNormal
	LocMapRoughness
	LocMapOccusion
	LocMapEmission
	LocMapHeight
	LocMapCubemap
	LocMapIrradiance
	LocMapPrefilter
	LocMapBrdf
)

// Material map type
const (
	// MapDiffuse
	MapAlbedo = iota
	MapMetalness
	MapNormal
	MapRoughness
	MapOcclusion
	MapEmission
	MapHeight
	// NOTE: Uses GL_TEXTURE_CUBE_MAP
	MapCubemap
	// NOTE: Uses GL_TEXTURE_CUBE_MAP
	MapIrradiance
	// NOTE: Uses GL_TEXTURE_CUBE_MAP
	MapPrefilter
	MapBrdf
)

// Material map type
const (
	MapDiffuse     = MapAlbedo
	MapSpecular    = MapMetalness
	LocMapDiffuse  = LocMapAlbedo
	LocMapSpecular = LocMapMetalness
)

// Shader and material limits
const (
	// Maximum number of predefined locations stored in shader struct
	MaxShaderLocations = 32
	// Maximum number of texture maps stored in shader struct
	MaxMaterialMaps = 12
)

// Mesh - Vertex data definning a mesh
type Mesh struct {
	// Number of vertices stored in arrays
	VertexCount int32
	// Number of triangles stored (indexed or not)
	TriangleCount int32
	// Vertex position (XYZ - 3 components per vertex) (shader-location = 0)
	Vertices *[]float32
	// Vertex texture coordinates (UV - 2 components per vertex) (shader-location = 1)
	Texcoords *[]float32
	// Vertex second texture coordinates (useful for lightmaps) (shader-location = 5)
	Texcoords2 *[]float32
	// Vertex normals (XYZ - 3 components per vertex) (shader-location = 2)
	Normals *[]float32
	// Vertex tangents (XYZ - 3 components per vertex) (shader-location = 4)
	Tangents *[]float32
	// Vertex colors (RGBA - 4 components per vertex) (shader-location = 3)
	Colors *[]uint8
	// Vertex indices (in case vertex data comes indexed)
	Indices *[]uint16
	// OpenGL Vertex Array Object id
	VaoID uint32
	// OpenGL Vertex Buffer Objects id (7 types of vertex data)
	VboID [7]uint32
}

func (m *Mesh) cptr() *C.Mesh {
	return (*C.Mesh)(unsafe.Pointer(m))
}

// NewMesh - Returns new Mesh
func NewMesh(vertexCount, triangleCount int32, vertices, texcoords, texcoords2, normals, tangents *[]float32, colors *[]uint8, indices *[]uint16, vaoID uint32, vboID [7]uint32) Mesh {
	return Mesh{vertexCount, triangleCount, vertices, texcoords, texcoords2, normals, tangents, colors, indices, vaoID, vboID}
}

// newMeshFromPointer - Returns new Mesh from pointer
func newMeshFromPointer(ptr unsafe.Pointer) Mesh {
	return *(*Mesh)(ptr)
}

// Material type
type Material struct {
	// Shader
	Shader Shader
	// Maps
	Maps [MaxMaterialMaps]MaterialMap
	// Padding
	_ [4]byte
	// Generic parameters (if required)
	Params *[]float32
}

func (m *Material) cptr() *C.Material {
	return (*C.Material)(unsafe.Pointer(m))
}

// NewMaterial - Returns new Material
func NewMaterial(shader Shader, maps [MaxMaterialMaps]MaterialMap, params *[]float32) Material {
	return Material{shader, maps, [4]byte{}, params}
}

// newMaterialFromPointer - Returns new Material from pointer
func newMaterialFromPointer(ptr unsafe.Pointer) Material {
	return *(*Material)(ptr)
}

// MaterialMap type
type MaterialMap struct {
	// Texture
	Texture Texture2D
	// Color
	Color Color
	// Value
	Value float32
}

// Model type
type Model struct {
	// Vertex data buffers (RAM and VRAM)
	Mesh Mesh
	// Local transform matrix
	Transform Matrix
	// Shader and textures data
	Material Material
	// Padding
	_ [4]byte
}

func (m *Model) cptr() *C.Model {
	return (*C.Model)(unsafe.Pointer(m))
}

// NewModel - Returns new Model
func NewModel(mesh Mesh, transform Matrix, material Material) Model {
	return Model{mesh, transform, material, [4]byte{}}
}

// newModelFromPointer - Returns new Model from pointer
func newModelFromPointer(ptr unsafe.Pointer) Model {
	return *(*Model)(ptr)
}

// Ray type (useful for raycast)
type Ray struct {
	// Ray position (origin)
	Position Vector3
	// Ray direction
	Direction Vector3
}

func (r *Ray) cptr() *C.Ray {
	return (*C.Ray)(unsafe.Pointer(r))
}

// NewRay - Returns new Ray
func NewRay(position, direction Vector3) Ray {
	return Ray{position, direction}
}

// newRayFromPointer - Returns new Ray from pointer
func newRayFromPointer(ptr unsafe.Pointer) Ray {
	return *(*Ray)(ptr)
}

// DrawLine3D - Draw a line in 3D world space
func DrawLine3D(startPos Vector3, endPos Vector3, color Color) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	ccolor := color.cptr()
	C.DrawLine3D(*cstartPos, *cendPos, *ccolor)
}

// DrawCircle3D - Draw a circle in 3D world space
func DrawCircle3D(center Vector3, radius float32, rotationAxis Vector3, rotationAngle float32, color Color) {
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	crotationAxis := rotationAxis.cptr()
	crotationAngle := (C.float)(rotationAngle)
	ccolor := color.cptr()
	C.DrawCircle3D(*ccenter, cradius, *crotationAxis, crotationAngle, *ccolor)
}

// DrawCube - Draw cube
func DrawCube(position Vector3, width float32, height float32, length float32, color Color) {
	cposition := position.cptr()
	cwidth := (C.float)(width)
	cheight := (C.float)(height)
	clength := (C.float)(length)
	ccolor := color.cptr()
	C.DrawCube(*cposition, cwidth, cheight, clength, *ccolor)
}

// DrawCubeV - Draw cube (Vector version)
func DrawCubeV(position Vector3, size Vector3, color Color) {
	cposition := position.cptr()
	csize := size.cptr()
	ccolor := color.cptr()
	C.DrawCubeV(*cposition, *csize, *ccolor)
}

// DrawCubeWires - Draw cube wires
func DrawCubeWires(position Vector3, width float32, height float32, length float32, color Color) {
	cposition := position.cptr()
	cwidth := (C.float)(width)
	cheight := (C.float)(height)
	clength := (C.float)(length)
	ccolor := color.cptr()
	C.DrawCubeWires(*cposition, cwidth, cheight, clength, *ccolor)
}

// DrawCubeTexture - Draw cube textured
func DrawCubeTexture(texture Texture2D, position Vector3, width float32, height float32, length float32, color Color) {
	ctexture := texture.cptr()
	cposition := position.cptr()
	cwidth := (C.float)(width)
	cheight := (C.float)(height)
	clength := (C.float)(length)
	ccolor := color.cptr()
	C.DrawCubeTexture(*ctexture, *cposition, cwidth, cheight, clength, *ccolor)
}

// DrawSphere - Draw sphere
func DrawSphere(centerPos Vector3, radius float32, color Color) {
	ccenterPos := centerPos.cptr()
	cradius := (C.float)(radius)
	ccolor := color.cptr()
	C.DrawSphere(*ccenterPos, cradius, *ccolor)
}

// DrawSphereEx - Draw sphere with extended parameters
func DrawSphereEx(centerPos Vector3, radius float32, rings int32, slices int32, color Color) {
	ccenterPos := centerPos.cptr()
	cradius := (C.float)(radius)
	crings := (C.int)(rings)
	cslices := (C.int)(slices)
	ccolor := color.cptr()
	C.DrawSphereEx(*ccenterPos, cradius, crings, cslices, *ccolor)
}

// DrawSphereWires - Draw sphere wires
func DrawSphereWires(centerPos Vector3, radius float32, rings int32, slices int32, color Color) {
	ccenterPos := centerPos.cptr()
	cradius := (C.float)(radius)
	crings := (C.int)(rings)
	cslices := (C.int)(slices)
	ccolor := color.cptr()
	C.DrawSphereWires(*ccenterPos, cradius, crings, cslices, *ccolor)
}

// DrawCylinder - Draw a cylinder/cone
func DrawCylinder(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, color Color) {
	cposition := position.cptr()
	cradiusTop := (C.float)(radiusTop)
	cradiusBottom := (C.float)(radiusBottom)
	cheight := (C.float)(height)
	cslices := (C.int)(slices)
	ccolor := color.cptr()
	C.DrawCylinder(*cposition, cradiusTop, cradiusBottom, cheight, cslices, *ccolor)
}

// DrawCylinderWires - Draw a cylinder/cone wires
func DrawCylinderWires(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, color Color) {
	cposition := position.cptr()
	cradiusTop := (C.float)(radiusTop)
	cradiusBottom := (C.float)(radiusBottom)
	cheight := (C.float)(height)
	cslices := (C.int)(slices)
	ccolor := color.cptr()
	C.DrawCylinderWires(*cposition, cradiusTop, cradiusBottom, cheight, cslices, *ccolor)
}

// DrawPlane - Draw a plane XZ
func DrawPlane(centerPos Vector3, size Vector2, color Color) {
	ccenterPos := centerPos.cptr()
	csize := size.cptr()
	ccolor := color.cptr()
	C.DrawPlane(*ccenterPos, *csize, *ccolor)
}

// DrawRay - Draw a ray line
func DrawRay(ray Ray, color Color) {
	cray := ray.cptr()
	ccolor := color.cptr()
	C.DrawRay(*cray, *ccolor)
}

// DrawGrid - Draw a grid (centered at (0, 0, 0))
func DrawGrid(slices int32, spacing float32) {
	cslices := (C.int)(slices)
	cspacing := (C.float)(spacing)
	C.DrawGrid(cslices, cspacing)
}

// DrawGizmo - Draw simple gizmo
func DrawGizmo(position Vector3) {
	cposition := position.cptr()
	C.DrawGizmo(*cposition)
}

// LoadMesh - Load mesh from file
func LoadMesh(fileName string) Mesh {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadMesh(cfileName)
	v := newMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadModel - Load model from file
func LoadModel(fileName string) Model {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadModel(cfileName)
	v := newModelFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadModelFromMesh - Load model from mesh data
func LoadModelFromMesh(data Mesh) Model {
	cdata := data.cptr()
	ret := C.LoadModelFromMesh(*cdata)
	v := newModelFromPointer(unsafe.Pointer(&ret))
	return v
}

// UnloadModel - Unload model from memory (RAM and/or VRAM)
func UnloadModel(model Model) {
	cmodel := model.cptr()
	C.UnloadModel(*cmodel)
}

// UnloadMesh - Unload mesh from memory (RAM and/or VRAM)
func UnloadMesh(mesh *Mesh) {
	cmesh := mesh.cptr()
	C.UnloadMesh(cmesh)
}

// GenMeshPlane - Generate plane mesh (with subdivisions)
func GenMeshPlane(width, length float32, resX, resZ int) Mesh {
	cwidth := (C.float)(width)
	clength := (C.float)(length)
	cresX := (C.int)(resX)
	cresZ := (C.int)(resZ)

	ret := C.GenMeshPlane(cwidth, clength, cresX, cresZ)
	v := newMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenMeshCube - Generate cuboid mesh
func GenMeshCube(width, height, length float32) Mesh {
	cwidth := (C.float)(width)
	cheight := (C.float)(height)
	clength := (C.float)(length)

	ret := C.GenMeshCube(cwidth, cheight, clength)
	v := newMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenMeshSphere - Generate sphere mesh (standard sphere)
func GenMeshSphere(radius float32, rings, slices int) Mesh {
	cradius := (C.float)(radius)
	crings := (C.int)(rings)
	cslices := (C.int)(slices)

	ret := C.GenMeshSphere(cradius, crings, cslices)
	v := newMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenMeshHemiSphere - Generate half-sphere mesh (no bottom cap)
func GenMeshHemiSphere(radius float32, rings, slices int) Mesh {
	cradius := (C.float)(radius)
	crings := (C.int)(rings)
	cslices := (C.int)(slices)

	ret := C.GenMeshHemiSphere(cradius, crings, cslices)
	v := newMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenMeshCylinder - Generate cylinder mesh
func GenMeshCylinder(radius, height float32, slices int) Mesh {
	cradius := (C.float)(radius)
	cheight := (C.float)(height)
	cslices := (C.int)(slices)

	ret := C.GenMeshCylinder(cradius, cheight, cslices)
	v := newMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenMeshTorus - Generate torus mesh
func GenMeshTorus(radius, size float32, radSeg, sides int) Mesh {
	cradius := (C.float)(radius)
	csize := (C.float)(size)
	cradSeg := (C.int)(radSeg)
	csides := (C.int)(sides)

	ret := C.GenMeshTorus(cradius, csize, cradSeg, csides)
	v := newMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenMeshKnot - Generate trefoil knot mesh
func GenMeshKnot(radius, size float32, radSeg, sides int) Mesh {
	cradius := (C.float)(radius)
	csize := (C.float)(size)
	cradSeg := (C.int)(radSeg)
	csides := (C.int)(sides)

	ret := C.GenMeshKnot(cradius, csize, cradSeg, csides)
	v := newMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenMeshHeightmap - Generate heightmap mesh from image data
func GenMeshHeightmap(heightmap Image, size Vector3) Mesh {
	cheightmap := heightmap.cptr()
	csize := size.cptr()

	ret := C.GenMeshHeightmap(*cheightmap, *csize)
	v := newMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenMeshCubicmap - Generate cubes-based map mesh from image data
func GenMeshCubicmap(cubicmap Image, size Vector3) Mesh {
	ccubicmap := cubicmap.cptr()
	csize := size.cptr()

	ret := C.GenMeshCubicmap(*ccubicmap, *csize)
	v := newMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadMaterial - Load material data (.MTL)
func LoadMaterial(fileName string) Material {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadMaterial(cfileName)
	v := newMaterialFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadMaterialDefault - Load default material (Supports: DIFFUSE, SPECULAR, NORMAL maps)
func LoadMaterialDefault() Material {
	ret := C.LoadMaterialDefault()
	v := newMaterialFromPointer(unsafe.Pointer(&ret))
	return v
}

// UnloadMaterial - Unload material textures from VRAM
func UnloadMaterial(material Material) {
	cmaterial := material.cptr()
	C.UnloadMaterial(*cmaterial)
}

// DrawModel - Draw a model (with texture if set)
func DrawModel(model Model, position Vector3, scale float32, tint Color) {
	cmodel := model.cptr()
	cposition := position.cptr()
	cscale := (C.float)(scale)
	ctint := tint.cptr()
	C.DrawModel(*cmodel, *cposition, cscale, *ctint)
}

// DrawModelEx - Draw a model with extended parameters
func DrawModelEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint Color) {
	cmodel := model.cptr()
	cposition := position.cptr()
	crotationAxis := rotationAxis.cptr()
	crotationAngle := (C.float)(rotationAngle)
	cscale := scale.cptr()
	ctint := tint.cptr()
	C.DrawModelEx(*cmodel, *cposition, *crotationAxis, crotationAngle, *cscale, *ctint)
}

// DrawModelWires - Draw a model wires (with texture if set)
func DrawModelWires(model Model, position Vector3, scale float32, tint Color) {
	cmodel := model.cptr()
	cposition := position.cptr()
	cscale := (C.float)(scale)
	ctint := tint.cptr()
	C.DrawModelWires(*cmodel, *cposition, cscale, *ctint)
}

// DrawModelWiresEx - Draw a model wires (with texture if set) with extended parameters
func DrawModelWiresEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint Color) {
	cmodel := model.cptr()
	cposition := position.cptr()
	crotationAxis := rotationAxis.cptr()
	crotationAngle := (C.float)(rotationAngle)
	cscale := scale.cptr()
	ctint := tint.cptr()
	C.DrawModelWiresEx(*cmodel, *cposition, *crotationAxis, crotationAngle, *cscale, *ctint)
}

// DrawBoundingBox - Draw bounding box (wires)
func DrawBoundingBox(box BoundingBox, color Color) {
	cbox := box.cptr()
	ccolor := color.cptr()
	C.DrawBoundingBox(*cbox, *ccolor)
}

// DrawBillboard - Draw a billboard texture
func DrawBillboard(camera Camera, texture Texture2D, center Vector3, size float32, tint Color) {
	ccamera := camera.cptr()
	ctexture := texture.cptr()
	ccenter := center.cptr()
	csize := (C.float)(size)
	ctint := tint.cptr()
	C.DrawBillboard(*ccamera, *ctexture, *ccenter, csize, *ctint)
}

// DrawBillboardRec - Draw a billboard texture defined by sourceRec
func DrawBillboardRec(camera Camera, texture Texture2D, sourceRec Rectangle, center Vector3, size float32, tint Color) {
	ccamera := camera.cptr()
	ctexture := texture.cptr()
	csourceRec := sourceRec.cptr()
	ccenter := center.cptr()
	csize := (C.float)(size)
	ctint := tint.cptr()
	C.DrawBillboardRec(*ccamera, *ctexture, *csourceRec, *ccenter, csize, *ctint)
}

// CalculateBoundingBox - Calculate mesh bounding box limits
func CalculateBoundingBox(mesh Mesh) BoundingBox {
	cmesh := mesh.cptr()
	ret := C.CalculateBoundingBox(*cmesh)
	v := newBoundingBoxFromPointer(unsafe.Pointer(&ret))
	return v
}

// CheckCollisionSpheres - Detect collision between two spheres
func CheckCollisionSpheres(centerA Vector3, radiusA float32, centerB Vector3, radiusB float32) bool {
	ccenterA := centerA.cptr()
	cradiusA := (C.float)(radiusA)
	ccenterB := centerB.cptr()
	cradiusB := (C.float)(radiusB)
	ret := C.CheckCollisionSpheres(*ccenterA, cradiusA, *ccenterB, cradiusB)
	v := bool(int(ret) == 1)
	return v
}

// CheckCollisionBoxes - Detect collision between two bounding boxes
func CheckCollisionBoxes(box1 BoundingBox, box2 BoundingBox) bool {
	cbox1 := box1.cptr()
	cbox2 := box2.cptr()
	ret := C.CheckCollisionBoxes(*cbox1, *cbox2)
	v := bool(int(ret) == 1)
	return v
}

// CheckCollisionBoxSphere - Detect collision between box and sphere
func CheckCollisionBoxSphere(box BoundingBox, centerSphere Vector3, radiusSphere float32) bool {
	cbox := box.cptr()
	ccenterSphere := centerSphere.cptr()
	cradiusSphere := (C.float)(radiusSphere)
	ret := C.CheckCollisionBoxSphere(*cbox, *ccenterSphere, cradiusSphere)
	v := bool(int(ret) == 1)
	return v
}

// CheckCollisionRaySphere - Detect collision between ray and sphere
func CheckCollisionRaySphere(ray Ray, spherePosition Vector3, sphereRadius float32) bool {
	cray := ray.cptr()
	cspherePosition := spherePosition.cptr()
	csphereRadius := (C.float)(sphereRadius)
	ret := C.CheckCollisionRaySphere(*cray, *cspherePosition, csphereRadius)
	v := bool(int(ret) == 1)
	return v
}

// CheckCollisionRaySphereEx - Detect collision between ray and sphere with extended parameters and collision point detection
func CheckCollisionRaySphereEx(ray Ray, spherePosition Vector3, sphereRadius float32, collisionPoint Vector3) bool {
	cray := ray.cptr()
	cspherePosition := spherePosition.cptr()
	csphereRadius := (C.float)(sphereRadius)
	ccollisionPoint := collisionPoint.cptr()
	ret := C.CheckCollisionRaySphereEx(*cray, *cspherePosition, csphereRadius, ccollisionPoint)
	v := bool(int(ret) == 1)
	return v
}

// CheckCollisionRayBox - Detect collision between ray and box
func CheckCollisionRayBox(ray Ray, box BoundingBox) bool {
	cray := ray.cptr()
	cbox := box.cptr()
	ret := C.CheckCollisionRayBox(*cray, *cbox)
	v := bool(int(ret) == 1)
	return v
}
