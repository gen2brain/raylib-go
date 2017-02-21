package raylib

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "reflect"

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

// NewMeshFromPointer - Returns new Mesh from pointer
func NewMeshFromPointer(ptr unsafe.Pointer) Mesh {
	return *(*Mesh)(ptr)
}

// Material type
type Material struct {
	// Standard shader (supports 3 map textures)
	Shader Shader
	// Diffuse texture  (binded to shader mapTexture0Loc)
	TexDiffuse Texture2D
	// Normal texture   (binded to shader mapTexture1Loc)
	TexNormal Texture2D
	// Specular texture (binded to shader mapTexture2Loc)
	TexSpecular Texture2D
	// Diffuse color
	ColDiffuse Color
	// Ambient color
	ColAmbient Color
	// Specular color
	ColSpecular Color
	// Glossiness level (Ranges from 0 to 1000)
	Glossiness float32
}

func (m *Material) cptr() *C.Material {
	return (*C.Material)(unsafe.Pointer(m))
}

// NewMaterial - Returns new Material
func NewMaterial(shader Shader, texDiffuse, texNormal, texSpecular Texture2D, colDiffuse, colAmbient, colSpecular Color, glossiness float32) Material {
	return Material{shader, texDiffuse, texNormal, texSpecular, colDiffuse, colAmbient, colSpecular, glossiness}
}

// NewMaterialFromPointer - Returns new Material from pointer
func NewMaterialFromPointer(ptr unsafe.Pointer) Material {
	return *(*Material)(ptr)
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

// NewModelFromPointer - Returns new Model from pointer
func NewModelFromPointer(ptr unsafe.Pointer) Model {
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

// NewRayFromPointer - Returns new Ray from pointer
func NewRayFromPointer(ptr unsafe.Pointer) Ray {
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
	v := NewMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadMeshEx - Load mesh from vertex data
func LoadMeshEx(numVertex int32, vData []float32, vtData []float32, vnData []float32, cData []Color) Mesh {
	cnumVertex := (C.int)(numVertex)
	cvData := (*C.float)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vData)).Data))
	cvtData := (*C.float)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vtData)).Data))
	cvnData := (*C.float)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vnData)).Data))
	ccData := (*C.Color)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&cData)).Data))
	ret := C.LoadMeshEx(cnumVertex, cvData, cvtData, cvnData, ccData)
	v := NewMeshFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadModel - Load model from file
func LoadModel(fileName string) Model {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadModel(cfileName)
	v := NewModelFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadModelFromMesh - Load model from mesh data
func LoadModelFromMesh(data Mesh, dynamic bool) Model {
	cdata := data.cptr()
	cdynamic := 0
	if dynamic {
		cdynamic = 1
	}
	ret := C.LoadModelFromMesh(*cdata, C.bool(cdynamic))
	v := NewModelFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadHeightmap - Load heightmap model from image data
func LoadHeightmap(heightmap *Image, size Vector3) Model {
	cheightmap := heightmap.cptr()
	csize := size.cptr()
	ret := C.LoadHeightmap(*cheightmap, *csize)
	v := NewModelFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadCubicmap - Load cubes-based map model from image data
func LoadCubicmap(cubicmap *Image) Model {
	ccubicmap := cubicmap.cptr()
	ret := C.LoadCubicmap(*ccubicmap)
	v := NewModelFromPointer(unsafe.Pointer(&ret))
	return v
}

// UnloadMesh - Unload mesh from memory (RAM and/or VRAM)
func UnloadMesh(mesh *Mesh) {
	cmesh := mesh.cptr()
	C.UnloadMesh(cmesh)
}

// UnloadModel - Unload model from memory (RAM and/or VRAM)
func UnloadModel(model Model) {
	cmodel := model.cptr()
	C.UnloadModel(*cmodel)
}

// LoadMaterial - Load material data (.MTL)
func LoadMaterial(fileName string) Material {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadMaterial(cfileName)
	v := NewMaterialFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadDefaultMaterial - Load default material (uses default models shader)
func LoadDefaultMaterial() Material {
	ret := C.LoadDefaultMaterial()
	v := NewMaterialFromPointer(unsafe.Pointer(&ret))
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
	v := NewBoundingBoxFromPointer(unsafe.Pointer(&ret))
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
