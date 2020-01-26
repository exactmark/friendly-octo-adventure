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

func main() {

	nSize := 3
	var startState SequentialInterface

	startState = createStartState(nSize, 1000)

	describe(startState)

	mySolver := createSolver()

	mySolver.solve(&startState, false)

	//fmt.Printf("(%v, %T)\n", goalState, goalState)

}
