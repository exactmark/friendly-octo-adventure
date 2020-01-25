package main

import (
	"fmt"
	//"math"
	//"math/rand"
	//"strings"
)

type SequentialInterface interface {
	lt(other *SequentialInterface) bool
	getChildren() []*SequentialInterface
	getH() int
	getExpectedCost() int
	getStateIdentifier() string
	isGoal() bool
	isValidMove(singleMove rune) bool
	makeMove(singleMove rune) bool
	shuffle(shuffleAmount int)
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}


type Solver struct {
	solutionMemo map[string]*SequentialInterface
	memoQueue    []*SequentialInterface
	memoSuccess  int
	greedy       bool
}

func createSolver() *Solver {

	memoQueue := make([]*SequentialInterface, 0)
	solutionMemo := make(map[string]*SequentialInterface)

	returnedSolver := Solver{
		solutionMemo: solutionMemo,
		memoQueue:    memoQueue,
		memoSuccess:  0,
		greedy:       true,
	}

	return &returnedSolver
}

func insertAt(myArray *[]*SequentialInterface, item *SequentialInterface, index int) *[]*SequentialInterface {

	actingArray := *myArray
	actingArray = append(actingArray, nil)
	copy(actingArray[index+1:], actingArray[index:])
	actingArray[index] = item
	return &actingArray
}

func (solver *Solver) solve(startState *SequentialInterface, greedy bool) *SequentialInterface {

	//frontierQueue := make(, 0)

	//heap.Init(&frontierQueue)

	//exploredStateCache := make(map[string]*SequentialInterface, 0)
	//
	//for ; len(frontierQueue) > 0; {
	//	exploringNode :=
	//}

	panic("implement me")
}

func main() {

	nSize := 3
	var startState SequentialInterface

	startState = createStartState(nSize, 10)

	describe(startState)

	mySolver := createSolver()

	mySolver.solve(&startState, true)

	//fmt.Printf("(%v, %T)\n", goalState, goalState)

}

