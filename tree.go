package main

type DecisionTree interface {
	Train(data []Instance)
	Decide(datum Instance) int // original decision will be ignored in this method
	Print()
}

type NodeType int

const (
	LeafNode NodeType = iota
	InteriorNode
)

// Shared tree node struct by all implementations of decision tree
type TreeNode struct {
	flag NodeType
	// Leaf node
	decision int // the decision of current node
	// Interior node
	iAttrib	int // the index of attribute that this node partitions
	children map[int]*TreeNode // look up table for all children according to their related attribute
}