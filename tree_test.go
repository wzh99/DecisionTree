package main

import (
	"fmt"
	"testing"
)

func TestTreePrint(t *testing.T) {
	train, _ := DataFromFile("data/monks-1.train")
	id3 := ID3Tree{useC4p5: false}
	id3.Train(train)
	id3.Print()
}

func TestID3Tree(t *testing.T) {
	for i := 1; i <= 3; i++ {
		fmt.Printf("MONKS %d: \n", i)
		train, _ := DataFromFile(fmt.Sprintf("data/monks-%d.train", i))
		test, _ := DataFromFile(fmt.Sprintf("data/monks-%d.test", i))
		id3 := &ID3Tree{useC4p5: false}
		fmt.Printf("ID3: %5.2f%%\n", 100 * testTree(id3, train, test))
		c4p5 := &ID3Tree{useC4p5: true}
		fmt.Printf("C4.5: %5.2f%%\n", 100 * testTree(c4p5, train, test))
	}
}

func TestRandomForest(t *testing.T) {
	for i := 1; i <= 3; i++ {
		fmt.Printf("MONKS %d: \n", i)
		train, _ := DataFromFile(fmt.Sprintf("data/monks-%d.train", i))
		test, _ := DataFromFile(fmt.Sprintf("data/monks-%d.test", i))
		for nTree := 1; nTree <= 10; nTree++ {
			// fmt.Printf("| %d | ", nTree)
			for ratio := 0.2; ratio <= 1; ratio += 0.1 {
				forest := NewRandomForest(nTree, ratio, false)
				fmt.Printf("%5.2f%% ", 100 * testTree(forest, train, test))
				// fmt.Printf("%6.4f ", testTree(forest, train, test))
			}
			fmt.Println()
		}
	}
}

func testTree(tree DecisionTree, train []Instance, test []Instance) (ratio float64){
	tree.Train(train)
	success := 0
	for _, item := range test {
		if tree.Decide(item) == item.decision {
			success++
		}
	}
	return float64(success) / float64(len(test))
}

