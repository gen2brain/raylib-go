//go:build !go1.21
// +build !go1.21

package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"runtime"

	"golang.org/x/exp/slices"
)

// List of VaoIDs of meshes created by calling UploadMesh()
// Used by UnloadMesh() to determine if mesh is go-managed or C-allocated
var goManagedMeshIDs = make([]uint32, 0)

// UploadMesh - Upload vertex data into a VAO (if supported) and VBO
func UploadMesh(mesh *Mesh, dynamic bool) {
	// check if mesh has already been uploaded to prevent duplication
	if mesh.VaoID != 0 {
		TraceLog(LogWarning, "VAO: [ID %d] Trying to re-load an already loaded mesh", mesh.VaoID)
		return
	}

	defer runtime.KeepAlive(mesh.Vertices)
	defer runtime.KeepAlive(mesh.Texcoords)
	defer runtime.KeepAlive(mesh.Texcoords2)
	defer runtime.KeepAlive(mesh.Normals)
	defer runtime.KeepAlive(mesh.Tangents)
	defer runtime.KeepAlive(mesh.Colors)
	defer runtime.KeepAlive(mesh.Indices)
	defer runtime.KeepAlive(mesh.AnimVertices)
	defer runtime.KeepAlive(mesh.AnimNormals)
	defer runtime.KeepAlive(mesh.BoneIds)
	defer runtime.KeepAlive(mesh.BoneWeights)
	defer runtime.KeepAlive(mesh.VboID)

	cMesh := mesh.cptr()
	C.UploadMesh(cMesh, C.bool(dynamic))

	// Add new mesh VaoID to list
	goManagedMeshIDs = append(goManagedMeshIDs, mesh.VaoID)
}

// UnloadMesh - Unload mesh from memory (RAM and/or VRAM)
func UnloadMesh(mesh *Mesh) {
	// Check list of go-managed mesh IDs
	if slices.Contains(goManagedMeshIDs, mesh.VaoID) {
		// C.UnloadMesh() only needs to read the VaoID & VboID
		// passing a temporary struct with all other fields nil makes it safe for the C code to call free()
		tempMesh := Mesh{
			VaoID: mesh.VaoID,
			VboID: mesh.VboID,
		}
		cmesh := tempMesh.cptr()
		C.UnloadMesh(*cmesh)

		// remove mesh VaoID from list
		goManagedMeshIDs = slices.DeleteFunc(goManagedMeshIDs, func(id uint32) bool { return id == mesh.VaoID })
	} else {
		cmesh := mesh.cptr()
		C.UnloadMesh(*cmesh)
	}
}
