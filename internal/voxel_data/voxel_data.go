package voxel_data

import "math"

const WORLD_SIZE = 32

type OctreeNode struct {
	Children [8]*OctreeNode
	IsLeaf   bool
}

type FlatNode struct {
	ChildIndices [8]int32 // -1 if no child
	IsLeaf       bool
}

func GetVoxels() []FlatNode {
	tree := buildOctreeFromGrid(createWorldMap())
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
func buildOctreeFromGrid(data [WORLD_SIZE * WORLD_SIZE * WORLD_SIZE]float32) *OctreeNode {
	if len(data) != WORLD_SIZE*WORLD_SIZE*WORLD_SIZE {
		panic("input data must be WORLD_SIZE^3 elements")
	}
	return buildRecursive(data, 0, 0, 0, WORLD_SIZE)
}

// Recursively builds the octree
func buildRecursive(data [WORLD_SIZE * WORLD_SIZE * WORLD_SIZE]float32, ox, oy, oz, size int) *OctreeNode {
	// Check if this region is empty or full
	full := false
	empty := true
	for z := oz; z < oz+size; z++ {
		for y := oy; y < oy+size; y++ {
			for x := ox; x < ox+size; x++ {
				if voxelAt(data, x, y, z) == 1.0 {
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
		return &OctreeNode{IsLeaf: true}
	}

	// Otherwise, subdivide
	node := &OctreeNode{}
	half := size / 2
	index := 0
	for dz := 0; dz < 2; dz++ {
		for dy := 0; dy < 2; dy++ {
			for dx := 0; dx < 2; dx++ {
				cx := ox + dx*half
				cy := oy + dy*half
				cz := oz + dz*half
				child := buildRecursive(data, cx, cy, cz, half)
				node.Children[index] = child
				index++
			}
		}
	}
	return node
}

// Helper to access voxel value
func voxelAt(data [WORLD_SIZE * WORLD_SIZE * WORLD_SIZE]float32, x, y, z int) float32 {
	if x < 0 || y < 0 || z < 0 || x >= WORLD_SIZE || y >= WORLD_SIZE || z >= WORLD_SIZE {
		return 0.0
	}
	return data[x+y*WORLD_SIZE+z*WORLD_SIZE*WORLD_SIZE]
}

func createWorldMap() [WORLD_SIZE * WORLD_SIZE * WORLD_SIZE]float32 {
	var data [WORLD_SIZE * WORLD_SIZE * WORLD_SIZE]float32
	radius := 15
	for x := range WORLD_SIZE {
		for y := range WORLD_SIZE {
			for z := range WORLD_SIZE {
				i := x + WORLD_SIZE*y + WORLD_SIZE*WORLD_SIZE*z
				if math.Pow(float64(x-radius), 2.)+math.Pow(float64(y-radius), 2)+math.Pow(float64(z-radius), 2) < float64(radius*radius) {
					data[i] = 1
				} else {
					data[i] = 0
				}
			}
		}
	}

	return data
}
