package main

import (
	"fmt"
	"math"
	"math/rand"
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

type NPuzzleState struct {
	greedy        bool
	parent        *NPuzzleState
	currentX      int
	currentY      int
	currentH      int
	nSize         int
	puzzleState   [][]int
	goalState     *[][]int
	goalDict      *map[int]coord
	cost          int
	lastMove      rune
	possibleMoves *[]rune
}

type coord struct {
	x int
	y int
}

//nPuzzle functions

func (s *NPuzzleState) shuffle(shuffleAmount int) {
	for x := 0; x < shuffleAmount; {
		thisMove := (*s.possibleMoves)[rand.Intn(len(*s.possibleMoves))]
		if s.makeMove(thisMove) {
			x++
			//s.printCurrentGoalState()
		}
	}
}

func (s *NPuzzleState) isValidMove(thisMove rune) bool {
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

func (s *NPuzzleState) makeSwap(x1 int, y1 int, x2 int, y2 int) {
	temp := s.puzzleState[y1][x1]
	s.puzzleState[y1][x1] = s.puzzleState[y2][x2]
	s.puzzleState[y2][x2] = temp
}

func (s *NPuzzleState) makeMove(thisMove rune) bool {
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

func (s *NPuzzleState) lt(other *SequentialInterface) bool {
	return s.getH() < (*other).getH()
}

func (s *NPuzzleState) getChildren() []*SequentialInterface {
	panic("implement me: getChildren")
}

func (s *NPuzzleState) getSingleManhattanScore(x int, y int, intChan chan int) {

}

func (s *NPuzzleState) getManhattanDistanceScore() int {
	//intChan := make(chan int)
	var returnSum int
	for x := 0; x < s.nSize; x++ {
		for y := 0; y < s.nSize; y++ {
			if s.puzzleState[y][x] != 0 {
				goal_coord := (*s.goalDict)[s.puzzleState[y][x]]
				returnSum += int(math.Abs(float64( goal_coord.x-x)) + math.Abs(float64(goal_coord.y-y)))
			}
		}
	}

	//for x := 0; x < s.nSize*s.nSize; x++ {
	//	returnSum += <-intChan
	//}

	return returnSum
}

func (s *NPuzzleState) getH() int {
	if s.currentH != -1 {
		return s.currentH
	} else {

		//	calculate h
		// store h
		s.currentH = s.getManhattanDistanceScore()
		return s.currentH
		// return h
	}
	//return -1
}

func (s *NPuzzleState) getExpectedCost() int {
	panic("implement me: getExpectedCost")
}

func (s *NPuzzleState) getStateIdentifier() string {
	panic("implement me: getStateIdentifier")
}

func (s *NPuzzleState) isGoal() bool {
	return s.getH() == 0
}

func describe(i SequentialInterface) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func getGoalState(nSize int) *[][]int {

	goalState := make([][]int, nSize)
	for y := 0; y < nSize; y++ {
		goalState[y] = make([]int, nSize)
		for x := 0; x < nSize; x++ {
			goalState[y][x] = (y * nSize) + (x) + 1
		}
	}

	goalState[nSize-1][nSize-1] = 0

	return &goalState
}

func (s *NPuzzleState) populateGoalDict() {

	for y:=0;y<s.nSize;y++{
		for x:=0;x<s.nSize ;x++  {
			new_coord:=coord{
				x: x,
				y: y,
			}
			(*s.goalDict)[(*s.goalState)[y][x]]=new_coord
			//fmt.Printf("%v,%v\n",(*s.goalState)[y][x],(*s.goalDict)[(*s.goalState)[y][x]])
		}
	}
}

func createStartState(nSize int, initShuffleAmount int) *NPuzzleState {
	goalState := getGoalState(nSize)
	possibleMoves := make([]rune, 0)
	possibleMoves = append(possibleMoves, 'u', 'l', 'd', 'r')

	goalDict := make(map[int]coord)

	startState := &NPuzzleState{
		goalState:     goalState,
		currentX:      nSize - 1,
		currentH:      -1,
		currentY:      nSize - 1,
		possibleMoves: &possibleMoves,
		nSize:         nSize,
		goalDict:      &goalDict,
	}

	startState.populateGoalDict()

	startState.puzzleState = *goalState
	startState.shuffle(initShuffleAmount)
	startState.getH()
	return startState
}

func (s *NPuzzleState) printCurrentGoalState() {
	for y := 0; y < s.nSize; y++ {
		fmt.Printf("%v\n", s.puzzleState[y])
	}
	fmt.Printf("\n")
}

func main() {

	nSize := 3
	var startState SequentialInterface

	startState = createStartState(nSize, 10000)
	describe(startState)
	//fmt.Printf("(%v, %T)\n", goalState, goalState)

	//	Start
	//	Create sample as goal
	//  Shuffle to make "start point"
	//
}
