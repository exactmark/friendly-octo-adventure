package main

import "fmt"

type SequentialInterface interface {
	lt() bool
	getChildren() []*SequentialInterface
	getH() int
	getExpectedCost() int
	getStateIdentifier() string
	isGoal() bool
	shuffle(int)

}

type SequentialState struct {
	greedy      bool
	parent      *SequentialState
	currentX    int
	currentY    int
	puzzleState [][]int
	goalState   *[][]int
	cost        int
	lastMove    rune
}

func (s SequentialState) shuffle(int) {
	panic("implement me")
}

//nPuzzle functions

func (s SequentialState) lt() bool {
	panic("implement me")
}

func (s SequentialState) getChildren() []*SequentialInterface {
	panic("implement me")
}

func (s SequentialState) getH() int {
	panic("implement me")
}

func (s SequentialState) getExpectedCost() int {
	panic("implement me")
}

func (s SequentialState) getStateIdentifier() string {
	panic("implement me")
}

func (s SequentialState) isGoal() bool {
	panic("implement me")
}

func describe(i SequentialInterface) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func getGoalState(nSize int) *[][]int {

	goalState := make([][]int, nSize)
	for x := 0; x < nSize; x++ {
		goalState[x] = make([]int, nSize)
		for y := 0; y < nSize; y++ {
			goalState[x][y] = (x * nSize) + (y) + 1
		}
	}

	goalState[nSize-1][nSize-1] = 0

	return &goalState
}

//func createStartState(nSize int)*SequentialState{
//
//
//}

func main() {

	nSize := 4
	goalState := getGoalState(nSize)
	var startState SequentialInterface

	startState = &SequentialState{goalState: goalState}
	describe(startState)
	fmt.Printf("(%v, %T)\n", goalState, goalState)

	//	Start
	//	Create sample as goal
	//  Shuffle to make "start point"
	//
}
