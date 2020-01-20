package main

import (
	"fmt"
	"math/rand"

	//"math/rand"
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

type SequentialState struct {
	greedy        bool
	parent        *SequentialState
	currentX      int
	currentY      int
	nSize         int
	puzzleState   [][]int
	goalState     *[][]int
	cost          int
	lastMove      rune
	possibleMoves *[]rune
}

//nPuzzle functions

func (s *SequentialState) shuffle(shuffleAmount int) {
	for x := 0; x < shuffleAmount; {
		thisMove := (*s.possibleMoves)[ rand.Intn(len(*s.possibleMoves))]
		if s.makeMove(thisMove) {
			x++
			//s.printCurrentGoalState()
		}
	}
}

func (s *SequentialState) isValidMove(thisMove rune) bool {
	if thisMove == 'u' {
		return s.currentY > 0
	} else if thisMove == 'd' {
		return s.currentY < s.nSize-1
	} else if thisMove == 'l' {
		return s.currentX > 0
	} else if thisMove == 'r' {
		return s.currentX < s.nSize-1
	}
	return false
}

func (s *SequentialState) makeSwap(x1 int, y1 int, x2 int, y2 int) {
	temp := s.puzzleState[y1][x1]
	s.puzzleState[y1][x1] = s.puzzleState[y2][x2]
	s.puzzleState[y2][x2] = temp
}

func (s *SequentialState) makeMove(thisMove rune) bool {
	if !s.isValidMove(thisMove) {
		return false
	}
	if thisMove == 'u' {
		s.makeSwap(s.currentX, s.currentY, s.currentX, s.currentY-1)
		s.currentY -= 1
		s.lastMove = thisMove
		return true
	} else if thisMove == 'd' {
		s.makeSwap(s.currentX, s.currentY, s.currentX, s.currentY+1)
		s.currentY += 1
		s.lastMove = thisMove
		return true

	} else if thisMove == 'l' {
		s.makeSwap(s.currentX, s.currentY, s.currentX-1, s.currentY)
		s.currentX -= 1
		s.lastMove = thisMove
		return true

	} else if thisMove == 'r' {
		s.makeSwap(s.currentX, s.currentY, s.currentX+1, s.currentY)
		s.currentX += 1
		s.lastMove = thisMove
		return true
	}
	return false
}

func (s *SequentialState) lt(other *SequentialInterface) bool {
	return s.getH() < (*other).getH()
}

func (s *SequentialState) getChildren() []*SequentialInterface {
	panic("implement me: getChildren")
}

func (s *SequentialState) getH() int {
	panic("implement me: getH")
}

func (s *SequentialState) getExpectedCost() int {
	panic("implement me: getExpectedCost")
}

func (s SequentialState) getStateIdentifier() string {
	panic("implement me: getStateIdentifier")
}

func (s SequentialState) isGoal() bool {
	panic("implement me: isGoal")
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
	possibleMoves := make([]rune, 0)
	possibleMoves = append(possibleMoves, 'u', 'l', 'd', 'r')
	//possibleMoves := [4]rune{'u', 'l', 'd', 'r'}
	startState := &SequentialState{
		goalState:     goalState,
		currentX:      nSize - 1,
		currentY:      nSize - 1,
		possibleMoves: &possibleMoves,
		nSize:         nSize,
	}
	startState.puzzleState = *goalState
	startState.shuffle(initShuffleAmount)
	return startState
}

func (s *SequentialState) printCurrentGoalState() {
	for y := 0; y < s.nSize; y++ {
		fmt.Printf("%v\n", s.puzzleState[y])
	}
	fmt.Printf("\n")
}

func main() {

	nSize := 3
	var startState SequentialInterface

	startState = createStartState(nSize, 10)
	describe(startState)
	//fmt.Printf("(%v, %T)\n", goalState, goalState)

	//	Start
	//	Create sample as goal
	//  Shuffle to make "start point"
	//
}
