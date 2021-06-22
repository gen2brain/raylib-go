package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// cptr returns C pointer
func (m *Mesh) cptr() *C.Mesh {
	return (*C.Mesh)(unsafe.Pointer(m))
}

// cptr returns C pointer
func (m *Material) cptr() *C.Material {
	return (*C.Material)(unsafe.Pointer(m))
}

// cptr returns C pointer
func (m *Model) cptr() *C.Model {
	return (*C.Model)(unsafe.Pointer(m))
}

// cptr returns C pointer
func (r *Ray) cptr() *C.Ray {
	return (*C.Ray)(unsafe.Pointer(r))
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
	C.UnloadMesh(*cmesh)
}

// ExportMesh - Export mesh as an OBJ file
func ExportMesh(mesh Mesh, fileName string) {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cmesh := mesh.cptr()
	C.ExportMesh(*cmesh, cfileName)
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

// LoadMaterials - Load material data (.MTL)
func LoadMaterials(fileName string) Material {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ccount := C.int(0)
	ret := C.LoadMaterials(cfileName, &ccount)
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
func DrawBillboardRec(camera Camera, texture Texture2D, sourceRec Rectangle, center Vector3, size Vector2, tint Color) {
	ccamera := camera.cptr()
	ctexture := texture.cptr()
	csourceRec := sourceRec.cptr()
	ccenter := center.cptr()
	csize := size.cptr()
	ctint := tint.cptr()
	C.DrawBillboardRec(*ccamera, *ctexture, *csourceRec, *ccenter, *csize, *ctint)
}

// DrawMesh - Draw a single mesh
func DrawMesh(mesh Mesh, material Material, transform Matrix) {
	C.DrawMesh(*mesh.cptr(), *material.cptr(), *transform.cptr())
}

// DrawMeshInstanced - Draw mesh with instanced rendering
func DrawMeshInstanced(mesh Mesh, material Material, transforms []Matrix, instances int) {
	C.DrawMeshInstanced(*mesh.cptr(), *material.cptr(), transforms[0].cptr(), C.int(instances))
}

// GetMeshBoundingBox - Compute mesh bounding box limits
func GetMeshBoundingBox(mesh Mesh) BoundingBox {
	cmesh := mesh.cptr()
	ret := C.GetMeshBoundingBox(*cmesh)
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

// CheckCollisionRaySphere - Detect collision between ray and sphere
func CheckCollisionRaySphere(ray Ray, spherePosition Vector3, sphereRadius float32) bool {
	cray := ray.cptr()
	cspherePosition := spherePosition.cptr()
	csphereRadius := (C.float)(sphereRadius)
	ret := C.CheckCollisionRaySphere(*cray, *cspherePosition, csphereRadius)
	v := bool(ret)
	return v
}

// CheckCollisionRaySphereEx - Detect collision between ray and sphere with extended parameters and collision point detection
func CheckCollisionRaySphereEx(ray Ray, spherePosition Vector3, sphereRadius float32, collisionPoint Vector3) bool {
	cray := ray.cptr()
	cspherePosition := spherePosition.cptr()
	csphereRadius := (C.float)(sphereRadius)
	ccollisionPoint := collisionPoint.cptr()
	ret := C.CheckCollisionRaySphereEx(*cray, *cspherePosition, csphereRadius, ccollisionPoint)
	v := bool(ret)
	return v
}

// CheckCollisionRayBox - Detect collision between ray and box
func CheckCollisionRayBox(ray Ray, box BoundingBox) bool {
	cray := ray.cptr()
	cbox := box.cptr()
	ret := C.CheckCollisionRayBox(*cray, *cbox)
	v := bool(ret)
	return v
}
