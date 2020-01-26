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
	getParent() *SequentialInterface
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func main() {

	nSize := 3
	var startState SequentialInterface

	startState = createStartState(nSize, 1000)

	//describe(startState)

	mySolver := createSolver()

	mySolver.solve(&startState, false)

	solvedList:=mySolver.solveAStar(&startState)

	for _,singleNode := range(*solvedList){
		thisState:=singleNode.(*NPuzzleState)
		thisState.printCurrentPuzzleState()
	}

}
