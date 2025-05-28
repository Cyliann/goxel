package voxel_data

import (
	"math"
	"sync"
)

const WORLD_SIZE = 256

type OctreeNode struct {
	Children [8]*OctreeNode
	IsLeaf   bool
}

type FlatNode struct {
	ChildIndices [8]int32 // -1 if no child; 32 bytes
	IsLeaf       bool     // 1 byte
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
	if size == 0 {
		return nil
	}

	firstVoxel := voxelAt(ox, oy, oz)
	isUniform := true

check:
	for z := oz; z < oz+size; z++ {
		for y := oy; y < oy+size; y++ {
			for x := ox; x < ox+size; x++ {
				if voxelAt(x, y, z) != firstVoxel {
					isUniform = false
					break check
				}
			}
		}
	}

	if !firstVoxel && isUniform {
		return nil
	}
	if size == 1 || isUniform {
		return &OctreeNode{IsLeaf: true}
	}

	node := &OctreeNode{}
	half := size / 2
	var wg sync.WaitGroup
	wg.Add(8)

	for dz := range 2 {
		for dy := range 2 {
			for dx := range 2 {
				index := dz*4 + dy*2 + dx
				cx := ox + dx*half
				cy := oy + dy*half
				cz := oz + dz*half

				go func(i, x, y, z int) {
					defer wg.Done()
					child := buildRecursive(x, y, z, half)
					node.Children[i] = child
				}(index, cx, cy, cz)
			}
		}
	}

	wg.Wait()
	return node
}

// Helper to access voxel value
func voxelAt(x, y, z int) bool {
	radius := 10
	return math.Pow(float64(x%(3*radius)-radius), 2.)+math.Pow(float64(y%(3*radius)-radius), 2)+math.Pow(float64(z%(3*radius)-radius), 2) < float64(radius*radius)
}
