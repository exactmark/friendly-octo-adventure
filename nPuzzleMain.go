package main

import (
	"fmt"
	"math/rand"
	"time"

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

	nSize := 4
	var startState SequentialInterface

	rand.Seed((1))

	startState = createStartState(nSize, 1000)


	startTime:=time.Now()

	//describe(startState)

	mySolver := createSolver()

	//mySolver.solve(&startState, false)

	solvedList:=mySolver.solveAStar(&startState)
	//solvedList:=mySolver.solveGreedy(&startState)

	for _,singleNode := range(*solvedList){
		var thisState *NPuzzleState
		thisState = (*singleNode).(*NPuzzleState)
		thisState.printCurrentPuzzleState()
	}

	fmt.Printf("Found solution in %v time.\n",time.Since(startTime))

	fmt.Printf("Found solution in %v steps\n",len(*solvedList))

}
