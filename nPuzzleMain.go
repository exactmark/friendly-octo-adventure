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
	createSequentialState(goalState interface{}, startState interface{}) *SequentialInterface
	exportCurrentState() interface{}
	setParent(node *SequentialInterface)
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func main() {
	for x := 1; x < 8; x++ {
		//mainBasicRun(2)
		mainLargeRun(2)
	}
	//mainBasicRun(4)
}

func mainBasicRun(seed int) {
	// this basic setup will solve A* in between 10 and 20 seconds and 51 steps.
	// the greedy version is non-deterministic. Tends to solve in 20ms, and around 150 steps.
	// And golang is going to be non-deterministic in both versions. I'm using
	// GoRoutines for child creation, which means that the order of moves is
	// potentially different. This explains the difference in solve times on
	// A* but does not really explain why different seeds will crash.

	nSize := 4
	var startState SequentialInterface

	rand.Seed(int64(seed))

	startState = createStartState(nSize, 250)

	startTime := time.Now()
	//startState.(*NPuzzleState).printCurrentPuzzleState()
	//startState.(*NPuzzleState).printCurrentGoalState()
	//describe(startState)

	mySolver := createSolver()

	//mySolver.solve(&startState, false)

	solvedList := mySolver.solveAStar(&startState)
	//solvedList:=mySolver.solveGreedy(&startState)
	//solvedList := mySolver.greedyGuidedAStar(&startState)


	//for _, singleNode := range *solvedList {
	//	var thisState *NPuzzleState
	//	thisState = (*singleNode).(*NPuzzleState)
	//	thisState.printCurrentPuzzleState()
	//}

	fmt.Printf("Found solution in %v time.\n", time.Since(startTime))

	fmt.Printf("Found solution in %v steps\n", len(*solvedList))
	//mySolver=nil
	//startState=nil
}

func mainLargeRun(seed int) {
	// this basic setup will solve A* in between 10 and 20 seconds and 51 steps.
	// the greedy version is non-deterministic. Tends to solve in 20ms, and around 150 steps.
	// And golang is going to be non-deterministic in both versions. I'm using
	// GoRoutines for child creation, which means that the order of moves is
	// potentially different. This explains the difference in solve times on
	// A* but does not really explain why different seeds will crash.

	nSize := 5
	var startState SequentialInterface

	rand.Seed(int64(seed))

	startState = createStartState(nSize, 1000)

	startTime := time.Now()
	//startState.(*NPuzzleState).printCurrentPuzzleState()
	//startState.(*NPuzzleState).printCurrentGoalState()
	//describe(startState)

	mySolver := createSolver()

	//mySolver.solve(&startState, false)

	//solvedList := mySolver.solveAStar(&startState)
	//solvedList:=mySolver.solveGreedy(&startState)
	solvedList := mySolver.greedyGuidedAStar(&startState)


	//for _, singleNode := range *solvedList {
	//	var thisState *NPuzzleState
	//	thisState = (*singleNode).(*NPuzzleState)
	//	thisState.printCurrentPuzzleState()
	//}

	fmt.Printf("Found solution in %v time.\n", time.Since(startTime))

	fmt.Printf("Found solution in %v steps\n", len(*solvedList))
	//mySolver=nil
	//startState=nil
}

