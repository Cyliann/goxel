package voxel_data

import (
	"math"
)

const WORLD_SIZE = 128

type OctreeNode struct {
	Children [8]*OctreeNode
	IsLeaf   int32
}

type FlatNode struct {
	ChildIndices [8]int32 // -1 if no child; 32 bytes
	IsLeaf       int32    //4 bytes
	_            [3]int32 // 12 bytes padding
}

func GetVoxels() []FlatNode {
	tree := buildOctreeFromGrid()
	nodes, _ := flattenTree(tree)
	return nodes
}

func flattenTree(root *OctreeNode) ([]FlatNode, int32) {
	var flat []FlatNode
	var helper func(node *OctreeNode) int32

	helper = func(node *OctreeNode) int32 {
		index := int32(len(flat))
		flat = append(flat, FlatNode{}) // placeholder

		var children [8]int32
		for i := range 8 {
			if node.Children[i] != nil {
				children[i] = helper(node.Children[i])
			} else {
				children[i] = -1
			}
		}

		flat[index] = FlatNode{
			ChildIndices: children,
			IsLeaf:       node.IsLeaf,
		}

		return index
	}

	rootIndex := helper(root)
	return flat, rootIndex
}

// Entry point
func buildOctreeFromGrid() *OctreeNode {
	return buildRecursive(0, 0, 0, WORLD_SIZE)
}

// Recursively builds the octree
func buildRecursive(ox, oy, oz, size int) *OctreeNode {
	// Check if this region is empty or full
	full := false
	empty := true
	for z := oz; z < oz+size; z++ {
		for y := oy; y < oy+size; y++ {
			for x := ox; x < ox+size; x++ {
				if voxelAt(x, y, z) {
					empty = false
				} else {
					full = true
				}
				if !empty && full {
					break
				}
			}
		}
	}

	// Base cases
	if empty {
		return nil // Skip empty space
	}
	if size == 1 || !full {
		return &OctreeNode{IsLeaf: 1}
	}

	// Otherwise, subdivide
	node := &OctreeNode{}
	half := size / 2
	index := 0
	for dz := range 2 {
		for dy := range 2 {
			for dx := range 2 {
				cx := ox + dx*half
				cy := oy + dy*half
				cz := oz + dz*half
				child := buildRecursive(cx, cy, cz, half)
				node.Children[index] = child
				index++
			}
		}
	}
	return node
}

// Helper to access voxel value
func voxelAt(x, y, z int) bool {
	radius := 10
	return math.Pow(float64(x%(3*radius)-radius), 2.)+math.Pow(float64(y%(3*radius)-radius), 2)+math.Pow(float64(z%(3*radius)-radius), 2) < float64(radius*radius)
}
