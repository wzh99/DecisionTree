package main

import (
	"fmt"
	"math/rand"
)

type RandomForest struct {
	trees []DecisionTree
	ratio float64 // ratio of training data used in every tree
}

func NewRandomForest(nTree int, ratio float64, useC4p5 bool) (f *RandomForest) {
	f = &RandomForest{ratio: ratio}
	f.trees = make([]DecisionTree, nTree)
	for i := range f.trees {
		f.trees[i] = &ID3Tree{useC4p5: useC4p5}
	}
	return
}

func (f *RandomForest) Train(data []Instance) {
	if len(data) == 0 {
		panic("No instance in data.")
	}
	for _, tree := range f.trees {
		// Choose data for the training of current tree
		treeData := make([]Instance, 0)
		for i := 0; i < int(f.ratio * float64(len(data))); i++ {
			treeData = append(treeData, data[rand.Int()%len(data)])
		}
		tree.Train(treeData)
	}
}

func (f *RandomForest) Decide(datum Instance) int {
	votes := make(map[int]int)
	for _, tree := range f.trees {
		decision := tree.Decide(datum)
		votes[decision]++
	}
	bestDecision, maxVotes := 0, 0
	for dec, vote := range votes {
		if vote > maxVotes {
			bestDecision, maxVotes = dec, vote
		}
	}
	return bestDecision
}

func (f *RandomForest) Print() {
	for i, tree := range f.trees {
		fmt.Println("Tree ", i)
		tree.Print()
	}
}
