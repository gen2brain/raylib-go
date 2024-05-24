//go:build go1.21
// +build go1.21

package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"slices"
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

	pinner := runtime.Pinner{}
	// Mesh pointer fields must be pinned to allow a Mesh pointer to be passed to C.UploadMesh() below
	// nil checks are required because Pin() will panic if passed nil
	if mesh.Vertices != nil {
		pinner.Pin(mesh.Vertices)
	}
	if mesh.Texcoords != nil {
		pinner.Pin(mesh.Texcoords)
	}
	if mesh.Texcoords2 != nil {
		pinner.Pin(mesh.Texcoords2)
	}
	if mesh.Normals != nil {
		pinner.Pin(mesh.Normals)
	}
	if mesh.Tangents != nil {
		pinner.Pin(mesh.Tangents)
	}
	if mesh.Colors != nil {
		pinner.Pin(mesh.Colors)
	}
	if mesh.Indices != nil {
		pinner.Pin(mesh.Indices)
	}
	if mesh.AnimVertices != nil {
		pinner.Pin(mesh.AnimVertices)
	}
	if mesh.AnimNormals != nil {
		pinner.Pin(mesh.AnimNormals)
	}
	if mesh.BoneIds != nil {
		pinner.Pin(mesh.BoneIds)
	}
	if mesh.BoneWeights != nil {
		pinner.Pin(mesh.BoneWeights)
	}
	// VboID of a new mesh should always be nil before uploading, but including this in case a mesh happens to have it set.
	if mesh.VboID != nil {
		pinner.Pin(mesh.VboID)
	}

	cMesh := mesh.cptr()
	C.UploadMesh(cMesh, C.bool(dynamic))

	// Add new mesh VaoID to list
	goManagedMeshIDs = append(goManagedMeshIDs, mesh.VaoID)

	pinner.Unpin()
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
