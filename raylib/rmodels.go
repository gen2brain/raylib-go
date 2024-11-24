package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"image/color"
	"unsafe"
)

// newMeshFromPointer - Returns new Mesh from pointer
func newMeshFromPointer(ptr unsafe.Pointer) Mesh {
	return *(*Mesh)(ptr)
}

// cptr returns C pointer
func (m *Mesh) cptr() *C.Mesh {
	return (*C.Mesh)(unsafe.Pointer(m))
}

// newMaterialFromPointer - Returns new Material from pointer
func newMaterialFromPointer(ptr unsafe.Pointer) Material {
	return *(*Material)(ptr)
}

// cptr returns C pointer
func (m *Material) cptr() *C.Material {
	return (*C.Material)(unsafe.Pointer(m))
}

// newModelFromPointer - Returns new Model from pointer
func newModelFromPointer(ptr unsafe.Pointer) Model {
	return *(*Model)(ptr)
}

// cptr returns C pointer
func (m *Model) cptr() *C.Model {
	return (*C.Model)(unsafe.Pointer(m))
}

// newRayFromPointer - Returns new Ray from pointer
func newRayFromPointer(ptr unsafe.Pointer) Ray {
	return *(*Ray)(ptr)
}

// cptr returns C pointer
func (r *Ray) cptr() *C.Ray {
	return (*C.Ray)(unsafe.Pointer(r))
}

// newModelAnimationFromPointer - Returns new ModelAnimation from pointer
func newModelAnimationFromPointer(ptr unsafe.Pointer) ModelAnimation {
	return *(*ModelAnimation)(ptr)
}

// cptr returns C pointer
func (r *ModelAnimation) cptr() *C.ModelAnimation {
	return (*C.ModelAnimation)(unsafe.Pointer(r))
}

// newRayCollisionFromPointer - Returns new RayCollision from pointer
func newRayCollisionFromPointer(ptr unsafe.Pointer) RayCollision {
	return *(*RayCollision)(ptr)
}

// DrawLine3D - Draw a line in 3D world space
func DrawLine3D(startPos Vector3, endPos Vector3, col color.RGBA) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	ccolor := colorCptr(col)
	C.DrawLine3D(*cstartPos, *cendPos, *ccolor)
}

// DrawPoint3D - Draw a point in 3D space, actually a small line
func DrawPoint3D(position Vector3, col color.RGBA) {
	cposition := position.cptr()
	ccolor := colorCptr(col)
	C.DrawPoint3D(*cposition, *ccolor)
}

// DrawCircle3D - Draw a circle in 3D world space
func DrawCircle3D(center Vector3, radius float32, rotationAxis Vector3, rotationAngle float32, col color.RGBA) {
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	crotationAxis := rotationAxis.cptr()
	crotationAngle := (C.float)(rotationAngle)
	ccolor := colorCptr(col)
	C.DrawCircle3D(*ccenter, cradius, *crotationAxis, crotationAngle, *ccolor)
}

// DrawTriangle3D - Draw a color-filled triangle (vertex in counter-clockwise order!)
func DrawTriangle3D(v1 Vector3, v2 Vector3, v3 Vector3, col color.RGBA) {
	cv1 := v1.cptr()
	cv2 := v2.cptr()
	cv3 := v3.cptr()
	ccolor := colorCptr(col)
	C.DrawTriangle3D(*cv1, *cv2, *cv3, *ccolor)
}

// DrawCube - Draw cube
func DrawCube(position Vector3, width float32, height float32, length float32, col color.RGBA) {
	cposition := position.cptr()
	cwidth := (C.float)(width)
	cheight := (C.float)(height)
	clength := (C.float)(length)
	ccolor := colorCptr(col)
	C.DrawCube(*cposition, cwidth, cheight, clength, *ccolor)
}

// DrawCubeV - Draw cube (Vector version)
func DrawCubeV(position Vector3, size Vector3, col color.RGBA) {
	cposition := position.cptr()
	csize := size.cptr()
	ccolor := colorCptr(col)
	C.DrawCubeV(*cposition, *csize, *ccolor)
}

// DrawCubeWires - Draw cube wires
func DrawCubeWires(position Vector3, width float32, height float32, length float32, col color.RGBA) {
	cposition := position.cptr()
	cwidth := (C.float)(width)
	cheight := (C.float)(height)
	clength := (C.float)(length)
	ccolor := colorCptr(col)
	C.DrawCubeWires(*cposition, cwidth, cheight, clength, *ccolor)
}

// DrawCubeWiresV - Draw cube wires (Vector version)
func DrawCubeWiresV(position Vector3, size Vector3, col color.RGBA) {
	cposition := position.cptr()
	csize := size.cptr()
	ccolor := colorCptr(col)
	C.DrawCubeWiresV(*cposition, *csize, *ccolor)
}

// DrawSphere - Draw sphere
func DrawSphere(centerPos Vector3, radius float32, col color.RGBA) {
	ccenterPos := centerPos.cptr()
	cradius := (C.float)(radius)
	ccolor := colorCptr(col)
	C.DrawSphere(*ccenterPos, cradius, *ccolor)
}

// DrawSphereEx - Draw sphere with extended parameters
func DrawSphereEx(centerPos Vector3, radius float32, rings int32, slices int32, col color.RGBA) {
	ccenterPos := centerPos.cptr()
	cradius := (C.float)(radius)
	crings := (C.int)(rings)
	cslices := (C.int)(slices)
	ccolor := colorCptr(col)
	C.DrawSphereEx(*ccenterPos, cradius, crings, cslices, *ccolor)
}

// DrawSphereWires - Draw sphere wires
func DrawSphereWires(centerPos Vector3, radius float32, rings int32, slices int32, col color.RGBA) {
	ccenterPos := centerPos.cptr()
	cradius := (C.float)(radius)
	crings := (C.int)(rings)
	cslices := (C.int)(slices)
	ccolor := colorCptr(col)
	C.DrawSphereWires(*ccenterPos, cradius, crings, cslices, *ccolor)
}

// DrawCylinder - Draw a cylinder/cone
func DrawCylinder(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, col color.RGBA) {
	cposition := position.cptr()
	cradiusTop := (C.float)(radiusTop)
	cradiusBottom := (C.float)(radiusBottom)
	cheight := (C.float)(height)
	cslices := (C.int)(slices)
	ccolor := colorCptr(col)
	C.DrawCylinder(*cposition, cradiusTop, cradiusBottom, cheight, cslices, *ccolor)
}

// DrawCylinderEx - Draw a cylinder with base at startPos and top at endPos
func DrawCylinderEx(startPos Vector3, endPos Vector3, startRadius float32, endRadius float32, sides int32, col color.RGBA) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	cstartRadius := (C.float)(startRadius)
	cendRadius := (C.float)(endRadius)
	csides := (C.int)(sides)
	ccolor := colorCptr(col)
	C.DrawCylinderEx(*cstartPos, *cendPos, cstartRadius, cendRadius, csides, *ccolor)
}

// DrawCylinderWires - Draw a cylinder/cone wires
func DrawCylinderWires(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, col color.RGBA) {
	cposition := position.cptr()
	cradiusTop := (C.float)(radiusTop)
	cradiusBottom := (C.float)(radiusBottom)
	cheight := (C.float)(height)
	cslices := (C.int)(slices)
	ccolor := colorCptr(col)
	C.DrawCylinderWires(*cposition, cradiusTop, cradiusBottom, cheight, cslices, *ccolor)
}

// DrawCylinderWiresEx - Draw a cylinder wires with base at startPos and top at endPos
func DrawCylinderWiresEx(startPos Vector3, endPos Vector3, startRadius float32, endRadius float32, sides int32, col color.RGBA) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	cstartRadius := (C.float)(startRadius)
	cendRadius := (C.float)(endRadius)
	csides := (C.int)(sides)
	ccolor := colorCptr(col)
	C.DrawCylinderWiresEx(*cstartPos, *cendPos, cstartRadius, cendRadius, csides, *ccolor)
}

// DrawCapsule - Draw a capsule with the center of its sphere caps at startPos and endPos
func DrawCapsule(startPos, endPos Vector3, radius float32, slices, rings int32, col color.RGBA) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	cradius := (C.float)(radius)
	cslices := (C.int)(slices)
	crings := (C.int)(rings)
	ccolor := colorCptr(col)
	C.DrawCapsule(*cstartPos, *cendPos, cradius, cslices, crings, *ccolor)
}

// DrawCapsuleWires - Draw capsule wireframe with the center of its sphere caps at startPos and endPos
func DrawCapsuleWires(startPos, endPos Vector3, radius float32, slices, rings int32, col color.RGBA) {
	cstartPos := startPos.cptr()
	cendPos := endPos.cptr()
	cradius := (C.float)(radius)
	cslices := (C.int)(slices)
	crings := (C.int)(rings)
	ccolor := colorCptr(col)
	C.DrawCapsuleWires(*cstartPos, *cendPos, cradius, cslices, crings, *ccolor)
}

// DrawPlane - Draw a plane XZ
func DrawPlane(centerPos Vector3, size Vector2, col color.RGBA) {
	ccenterPos := centerPos.cptr()
	csize := size.cptr()
	ccolor := colorCptr(col)
	C.DrawPlane(*ccenterPos, *csize, *ccolor)
}

// DrawRay - Draw a ray line
func DrawRay(ray Ray, col color.RGBA) {
	cray := ray.cptr()
	ccolor := colorCptr(col)
	C.DrawRay(*cray, *ccolor)
}

// DrawGrid - Draw a grid (centered at (0, 0, 0))
func DrawGrid(slices int32, spacing float32) {
	cslices := (C.int)(slices)
	cspacing := (C.float)(spacing)
	C.DrawGrid(cslices, cspacing)
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

// IsModelValid - Check if a model is valid (loaded in GPU, VAO/VBOs)
func IsModelValid(model Model) bool {
	cmodel := model.cptr()
	ret := C.IsModelValid(*cmodel)
	v := bool(ret)
	return v
}

// UnloadModel - Unload model from memory (RAM and/or VRAM)
func UnloadModel(model Model) {
	cmodel := model.cptr()
	C.UnloadModel(*cmodel)
}

// GetModelBoundingBox - Compute model bounding box limits (considers all meshes
func GetModelBoundingBox(model Model) BoundingBox {
	cmodel := model.cptr()
	ret := C.GetModelBoundingBox(*cmodel)
	v := newBoundingBoxFromPointer(unsafe.Pointer(&ret))
	return v
}

// DrawModel - Draw a model (with texture if set)
func DrawModel(model Model, position Vector3, scale float32, tint color.RGBA) {
	cmodel := model.cptr()
	cposition := position.cptr()
	cscale := (C.float)(scale)
	ctint := colorCptr(tint)
	C.DrawModel(*cmodel, *cposition, cscale, *ctint)
}

// DrawModelEx - Draw a model with extended parameters
func DrawModelEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
	cmodel := model.cptr()
	cposition := position.cptr()
	crotationAxis := rotationAxis.cptr()
	crotationAngle := (C.float)(rotationAngle)
	cscale := scale.cptr()
	ctint := colorCptr(tint)
	C.DrawModelEx(*cmodel, *cposition, *crotationAxis, crotationAngle, *cscale, *ctint)
}

// DrawModelWires - Draw a model wires (with texture if set)
func DrawModelWires(model Model, position Vector3, scale float32, tint color.RGBA) {
	cmodel := model.cptr()
	cposition := position.cptr()
	cscale := (C.float)(scale)
	ctint := colorCptr(tint)
	C.DrawModelWires(*cmodel, *cposition, cscale, *ctint)
}

// DrawModelWiresEx - Draw a model wires (with texture if set) with extended parameters
func DrawModelWiresEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
	cmodel := model.cptr()
	cposition := position.cptr()
	crotationAxis := rotationAxis.cptr()
	crotationAngle := (C.float)(rotationAngle)
	cscale := scale.cptr()
	ctint := colorCptr(tint)
	C.DrawModelWiresEx(*cmodel, *cposition, *crotationAxis, crotationAngle, *cscale, *ctint)
}

// DrawModelPoints - Draw a model as points
func DrawModelPoints(model Model, position Vector3, scale float32, tint color.RGBA) {
	cmodel := model.cptr()
	cposition := position.cptr()
	cscale := (C.float)(scale)
	ctint := colorCptr(tint)
	C.DrawModelPoints(*cmodel, *cposition, cscale, *ctint)
}

// DrawModelPointsEx - Draw a model as points with extended parameters
func DrawModelPointsEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
	cmodel := model.cptr()
	cposition := position.cptr()
	crotationAxis := rotationAxis.cptr()
	crotationAngle := (C.float)(rotationAngle)
	cscale := scale.cptr()
	ctint := colorCptr(tint)
	C.DrawModelPointsEx(*cmodel, *cposition, *crotationAxis, crotationAngle, *cscale, *ctint)
}

// DrawBoundingBox - Draw bounding box (wires)
func DrawBoundingBox(box BoundingBox, col color.RGBA) {
	cbox := box.cptr()
	ccolor := colorCptr(col)
	C.DrawBoundingBox(*cbox, *ccolor)
}

// DrawBillboard - Draw a billboard texture
func DrawBillboard(camera Camera, texture Texture2D, center Vector3, scale float32, tint color.RGBA) {
	ccamera := camera.cptr()
	ctexture := texture.cptr()
	ccenter := center.cptr()
	cscale := (C.float)(scale)
	ctint := colorCptr(tint)
	C.DrawBillboard(*ccamera, *ctexture, *ccenter, cscale, *ctint)
}

// DrawBillboardRec - Draw a billboard texture defined by sourceRec
func DrawBillboardRec(camera Camera, texture Texture2D, sourceRec Rectangle, center Vector3, size Vector2, tint color.RGBA) {
	ccamera := camera.cptr()
	ctexture := texture.cptr()
	csourceRec := sourceRec.cptr()
	ccenter := center.cptr()
	csize := size.cptr()
	ctint := colorCptr(tint)
	C.DrawBillboardRec(*ccamera, *ctexture, *csourceRec, *ccenter, *csize, *ctint)
}

// DrawBillboardPro - Draw a billboard texture with pro parameters
func DrawBillboardPro(camera Camera, texture Texture2D, sourceRec Rectangle, position Vector3, up Vector3, size Vector2, origin Vector2, rotation float32, tint Color) {
	ccamera := camera.cptr()
	ctexture := texture.cptr()
	csourceRec := sourceRec.cptr()
	cposition := position.cptr()
	cup := up.cptr()
	csize := size.cptr()
	corigin := origin.cptr()
	crotation := (C.float)(rotation)
	ctint := colorCptr(tint)
	C.DrawBillboardPro(*ccamera, *ctexture, *csourceRec, *cposition, *cup, *csize, *corigin, crotation, *ctint)
}

// UpdateMeshBuffer - Update mesh vertex data in GPU for a specific buffer index
func UpdateMeshBuffer(mesh Mesh, index int, data []byte, offset int) {
	cindex := (C.int)(index)
	coffset := (C.int)(offset)
	cdataSize := (C.int)(len(data))
	C.UpdateMeshBuffer(*mesh.cptr(), cindex, unsafe.Pointer(&data[0]), cdataSize, coffset)
}

// DrawMesh - Draw a single mesh
func DrawMesh(mesh Mesh, material Material, transform Matrix) {
	C.DrawMesh(*mesh.cptr(), *material.cptr(), *transform.cptr())
}

// DrawMeshInstanced - Draw mesh with instanced rendering
func DrawMeshInstanced(mesh Mesh, material Material, transforms []Matrix, instances int) {
	C.DrawMeshInstanced(*mesh.cptr(), *material.cptr(), transforms[0].cptr(), C.int(instances))
}

// ExportMesh - Export mesh as an OBJ file
func ExportMesh(mesh Mesh, fileName string) {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cmesh := mesh.cptr()
	C.ExportMesh(*cmesh, cfileName)
}

// GetMeshBoundingBox - Compute mesh bounding box limits
func GetMeshBoundingBox(mesh Mesh) BoundingBox {
	cmesh := mesh.cptr()
	ret := C.GetMeshBoundingBox(*cmesh)
	v := newBoundingBoxFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenMeshPoly - Generate polygonal mesh
func GenMeshPoly(sides int, radius float32) Mesh {
	csides := (C.int)(sides)
	cradius := (C.float)(radius)

	ret := C.GenMeshPoly(csides, cradius)
	v := newMeshFromPointer(unsafe.Pointer(&ret))
	return v
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

// GenMeshCone - Generate cone/pyramid mesh
func GenMeshCone(radius, height float32, slices int) Mesh {
	cradius := (C.float)(radius)
	cheight := (C.float)(height)
	cslices := (C.int)(slices)

	ret := C.GenMeshCone(cradius, cheight, cslices)
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

// LoadMaterials - Load material data (.MTL)
func LoadMaterials(fileName string) []Material {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ccount := C.int(0)
	ret := C.LoadMaterials(cfileName, &ccount)
	v := (*[1 << 24]Material)(unsafe.Pointer(ret))[:int(ccount)]
	return v
}

// LoadMaterialDefault - Load default material (Supports: DIFFUSE, SPECULAR, NORMAL maps)
func LoadMaterialDefault() Material {
	ret := C.LoadMaterialDefault()
	v := newMaterialFromPointer(unsafe.Pointer(&ret))
	return v
}

// IsMaterialValid - Check if a material is valid (shader assigned, map textures loaded in GPU)
func IsMaterialValid(material Material) bool {
	cmaterial := material.cptr()
	ret := C.IsMaterialValid(*cmaterial)
	v := bool(ret)
	return v
}

// UnloadMaterial - Unload material textures from VRAM
func UnloadMaterial(material Material) {
	cmaterial := material.cptr()
	C.UnloadMaterial(*cmaterial)
}

// SetMaterialTexture - Set texture for a material map type (MATERIAL_MAP_DIFFUSE, MATERIAL_MAP_SPECULAR...)
func SetMaterialTexture(material *Material, mapType int32, texture Texture2D) {
	cmaterial := material.cptr()
	cmapType := (C.int)(mapType)
	ctexture := texture.cptr()
	C.SetMaterialTexture(cmaterial, cmapType, *ctexture)
}

// SetModelMeshMaterial - Set material for a mesh
func SetModelMeshMaterial(model *Model, meshId int32, materialId int32) {
	cmodel := model.cptr()
	cmeshId := (C.int)(meshId)
	cmaterialId := (C.int)(materialId)
	C.SetModelMeshMaterial(cmodel, cmeshId, cmaterialId)
}

// LoadModelAnimations - Load model animations from file
func LoadModelAnimations(fileName string) []ModelAnimation {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ccount := C.int(0)
	ret := C.LoadModelAnimations(cfileName, &ccount)
	v := (*[1 << 24]ModelAnimation)(unsafe.Pointer(ret))[:int(ccount)]
	return v
}

// UpdateModelAnimation - Update model animation pose (CPU)
func UpdateModelAnimation(model Model, anim ModelAnimation, frame int32) {
	cmodel := model.cptr()
	canim := anim.cptr()
	cframe := (C.int)(frame)
	C.UpdateModelAnimation(*cmodel, *canim, cframe)
}

// UpdateModelAnimationBones - Update model animation mesh bone matrices (GPU skinning)
func UpdateModelAnimationBones(model Model, anim ModelAnimation, frame int32) {
	cmodel := model.cptr()
	canim := anim.cptr()
	cframe := (C.int)(frame)
	C.UpdateModelAnimationBones(*cmodel, *canim, cframe)
}

// UnloadModelAnimation - Unload animation data
func UnloadModelAnimation(anim ModelAnimation) {
	canim := anim.cptr()
	C.UnloadModelAnimation(*canim)
}

// UnloadModelAnimations - Unload animation array data
func UnloadModelAnimations(animations []ModelAnimation) {
	C.UnloadModelAnimations((*C.ModelAnimation)(unsafe.Pointer(&animations[0])), (C.int)(len(animations)))
}

// IsModelAnimationValid - Check model animation skeleton match
func IsModelAnimationValid(model Model, anim ModelAnimation) bool {
	cmodel := model.cptr()
	canim := anim.cptr()
	ret := C.IsModelAnimationValid(*cmodel, *canim)
	v := bool(ret)
	return v
}

// CheckCollisionSpheres - Detect collision between two spheres
func CheckCollisionSpheres(centerA Vector3, radiusA float32, centerB Vector3, radiusB float32) bool {
	ccenterA := centerA.cptr()
	cradiusA := (C.float)(radiusA)
	ccenterB := centerB.cptr()
	cradiusB := (C.float)(radiusB)
	ret := C.CheckCollisionSpheres(*ccenterA, cradiusA, *ccenterB, cradiusB)
	v := bool(ret)
	return v
}

// CheckCollisionBoxes - Detect collision between two bounding boxes
func CheckCollisionBoxes(box1 BoundingBox, box2 BoundingBox) bool {
	cbox1 := box1.cptr()
	cbox2 := box2.cptr()
	ret := C.CheckCollisionBoxes(*cbox1, *cbox2)
	v := bool(ret)
	return v
}

// CheckCollisionBoxSphere - Detect collision between box and sphere
func CheckCollisionBoxSphere(box BoundingBox, centerSphere Vector3, radiusSphere float32) bool {
	cbox := box.cptr()
	ccenterSphere := centerSphere.cptr()
	cradiusSphere := (C.float)(radiusSphere)
	ret := C.CheckCollisionBoxSphere(*cbox, *ccenterSphere, cradiusSphere)
	v := bool(ret)
	return v
}

// GetRayCollisionSphere - Get collision info between ray and sphere
func GetRayCollisionSphere(ray Ray, center Vector3, radius float32) RayCollision {
	cray := ray.cptr()
	ccenter := center.cptr()
	cradius := (C.float)(radius)
	ret := C.GetRayCollisionSphere(*cray, *ccenter, cradius)
	v := newRayCollisionFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetRayCollisionBox - Get collision info between ray and box
func GetRayCollisionBox(ray Ray, box BoundingBox) RayCollision {
	cray := ray.cptr()
	cbox := box.cptr()
	ret := C.GetRayCollisionBox(*cray, *cbox)
	v := newRayCollisionFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetRayCollisionMesh - Get collision info between ray and mesh
func GetRayCollisionMesh(ray Ray, mesh Mesh, transform Matrix) RayCollision {
	cray := ray.cptr()
	cmesh := mesh.cptr()
	ctransform := transform.cptr()
	ret := C.GetRayCollisionMesh(*cray, *cmesh, *ctransform)
	v := newRayCollisionFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetRayCollisionTriangle - Get collision info between ray and triangle
func GetRayCollisionTriangle(ray Ray, p1, p2, p3 Vector3) RayCollision {
	cray := ray.cptr()
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	ret := C.GetRayCollisionTriangle(*cray, *cp1, *cp2, *cp3)
	v := newRayCollisionFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetRayCollisionQuad - Get collision info between ray and quad
func GetRayCollisionQuad(ray Ray, p1, p2, p3, p4 Vector3) RayCollision {
	cray := ray.cptr()
	cp1 := p1.cptr()
	cp2 := p2.cptr()
	cp3 := p3.cptr()
	cp4 := p4.cptr()
	ret := C.GetRayCollisionQuad(*cray, *cp1, *cp2, *cp3, *cp4)
	v := newRayCollisionFromPointer(unsafe.Pointer(&ret))
	return v
}
