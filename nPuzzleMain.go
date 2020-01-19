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
	greedy        bool
	parent        *SequentialState
	currentX      int
	currentY      int
	puzzleState   [][]int
	goalState     *[][]int
	cost          int
	lastMove      rune
	possibleMoves *[]rune
}

func (s SequentialState) shuffle(shuffleAmount int) {
	//panic("implement me")
	fmt.Printf("%v\n",*s.possibleMoves)
	numPossibleMoves:= len(*s.possibleMoves)
	for x:=0; x<shuffleAmount;x++  {
		
	}

	fmt.Printf("implement shuffle\n")
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

func createStartState(nSize int, initShuffleAmount int) *SequentialState {
	goalState := getGoalState(nSize)
	possibleMoves := make([]rune,0)
	possibleMoves = append(possibleMoves, 'u', 'l', 'd', 'r')
	//possibleMoves := [4]rune{'u', 'l', 'd', 'r'}
	startState := &SequentialState{
		goalState:     goalState,
		currentX:      nSize - 1,
		currentY:      nSize - 1,
		possibleMoves: &possibleMoves,
	}
	startState.puzzleState = *goalState
	startState.shuffle(initShuffleAmount)
	return startState
}

func main() {

	nSize := 3
	var startState SequentialInterface

	startState = createStartState(nSize, 1)
	describe(startState)
	//fmt.Printf("(%v, %T)\n", goalState, goalState)

	//	Start
	//	Create sample as goal
	//  Shuffle to make "start point"
	//
}
