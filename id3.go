package main

import (
	"fmt"
	"math"
	"math/rand"
)

type ID3Tree struct {
	root *TreeNode
	useC4p5	bool // use C4.5 algorithm to evaluate attribute
	nDecisions int
}

func (t *ID3Tree) Train(data []Instance) {
	if len(data) == 0 {
		panic("No instance in data.")
	}
	t.root = new(TreeNode)
	available := [6]bool{true, true, true, true, true, true}
	dataPtrs := make([]*Instance, len(data))
	for i := range data {
		dataPtrs[i] = &data[i]
	}
	t.buildNode(t.root, dataPtrs, available)
}

func (t *ID3Tree) buildNode(node *TreeNode, data []*Instance, available [6]bool) {
	// Try building leaf node
	nAttrib := 0 // count the number of available attributes
	for _, ok := range available {
		if ok {
			nAttrib++
		}
	}
	decisionCount := make(map[int]int) // lookup table for the number of decisions
	for _, inst := range data {
		decisionCount[inst.decision]++
	}
	if node == t.root {
		t.nDecisions = len(decisionCount)
	}

	// Build lead node if all instance have the same decision
	if len(decisionCount) == 1 {
		node.flag = LeafNode
		node.decision = data[0].decision
		return
	}
	// Build leaf node if no attributes are available
	if nAttrib == 0 {
		maxDec, maxCount := 0, decisionCount[0] // choose the decision of the largest count
		for dec, cnt := range decisionCount {
			if cnt > maxCount {
				maxDec, maxCount = dec, cnt
			}
		}
		node.flag = LeafNode
		node.decision = maxDec
		return
	}

	// Compute entropy of decisions
	var entropy float64 // H(D)
	for _, count := range decisionCount {
		prob := float64(count) / float64(len(data))
		entropy -= prob * math.Log(prob)
	}

	// Choose the attribute that best divides the instances
	var bestIdx int
	var bestMap map[int][]int
	bestMetric := math.Inf(-1)
	for ia, ok := range available {
		var relative float64
		var ratio float64
		if !ok {
			continue
		}

		// Get distribution of current attribute
		attribMap := make(map[int][]int) // attribute value -> instance index
		for id, inst := range data {
			if attribMap[inst.attrib[ia]] == nil {
				attribMap[inst.attrib[ia]] = make([]int, 0)
			}
			attribMap[inst.attrib[ia]] = append(attribMap[inst.attrib[ia]], id)
		}

		// Compute cross entropy
		var cross float64 // H(S|A)
		var attribEntropy float64 // H_A(D)
		for _, list := range attribMap {
			probAttrib := float64(len(list)) / float64(len(data))
			attribEntropy -= probAttrib * math.Log(probAttrib)
			attribDecCount := make(map[int]int) // decision of instances of specific attribute -> count
			for _, iInst := range list {
				attribDecCount[data[iInst].decision]++
			}
			var entropySum float64
			for _, count := range attribDecCount {
				prob := float64(count) / float64(len(list))
				entropySum -= prob * math.Log(prob)
			}
			cross += probAttrib * entropySum
		}

		// Compute relative entropy g(D,A)
		relative = entropy - cross

		// Compute entropy ratio
		ratio = relative / attribEntropy

		// Choose metric to evaluate attributes
		metric := relative
		if t.useC4p5 {
			metric = ratio
		}

		if metric > bestMetric {
			bestIdx, bestMetric, bestMap = ia, metric, attribMap
		}
	}

	// Create interior node
	node.flag = InteriorNode
	node.iAttrib = bestIdx
	node.children = make(map[int]*TreeNode)
	available[bestIdx] = false // this attribute will no longer be considered
	for attrib, list := range bestMap {
		node.children[attrib] = new(TreeNode)
		nodeData := make([]*Instance, 0)
		for _, idx := range list {
			nodeData = append(nodeData, data[idx])
		}
		t.buildNode(node.children[attrib], nodeData, available)
	}
}

func (t *ID3Tree) Decide(datum Instance) int {
	node := t.root
	for {
		switch node.flag {
		case LeafNode:
			return node.decision
		case InteriorNode:
			child := node.children[datum.attrib[node.iAttrib]]
			if child == nil { // can't decide which node to go
				return rand.Int() % t.nDecisions // randomly choose one
			}
			node = child
		}
	}
}

func (t *ID3Tree) Print() {
	t.printNode(t.root, 0)
}

func (t *ID3Tree) printNode(node *TreeNode, depth int) {
	for it := 0; it < depth; it++ {
		fmt.Print("\t")
	}
	switch node.flag {
	case LeafNode:
		fmt.Print("LEAF ")
		fmt.Print("decision: ", node.decision)
		fmt.Println()
	case InteriorNode:
		fmt.Print("INTERIOR ")
		fmt.Print("iAttrib: ", node.iAttrib, " ")
		list := make([]int, 0)
		for attrib := range node.children {
			list = append(list, attrib)
		}
		fmt.Print("children: ", list)
		fmt.Println()
		for _, child := range node.children {
			t.printNode(child, depth + 1)
		}
	}
}
