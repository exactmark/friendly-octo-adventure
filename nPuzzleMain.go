package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
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
	parent                 *NPuzzleState
	currentX               int
	currentY               int
	currentH               int
	nSize                  int
	puzzleState            [][]int
	goalState              *[][]int
	goalDict               *map[int]coord
	cost                   int
	lastMove               rune
	possibleMoves          *[]rune
	stateIdentifier        string
	stateIdentifierCreated bool
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

func makeChild(s *NPuzzleState, direction rune, returnChan chan *NPuzzleState) {
	childOne := *s
	childOne.currentH = -1
	childOne.puzzleState = make([][]int, s.nSize)
	for y := 0; y < s.nSize; y++ {
		childOne.puzzleState[y] = make([]int, s.nSize)
		copy(childOne.puzzleState[y], s.puzzleState[y])
	}

	childOne.stateIdentifierCreated = false

	childOne.makeMove(direction)
	childOne.getH()
	returnChan <- &childOne

}

func (s *NPuzzleState) getChildren() []*SequentialInterface {

	returnList := make([]*SequentialInterface, 0)
	var counter int
	returnChannel := make(chan *NPuzzleState)
	for _, direction := range *(s.possibleMoves) {
		//describe(s)
		if s.isValidMove(direction) {
			counter++
			go makeChild(s, direction, returnChannel)
		}
	}

	for x := 0; x < counter; x++ {
		singleChild := SequentialInterface(<-returnChannel)
		returnList = append(returnList, &singleChild)
	}

	return returnList
}

func (s *NPuzzleState) getManhattanDistanceScore() int {
	//intChan := make(chan int)
	var returnSum int
	for x := 0; x < s.nSize; x++ {
		for y := 0; y < s.nSize; y++ {
			if s.puzzleState[y][x] != 0 {
				goalCoord := (*s.goalDict)[s.puzzleState[y][x]]
				returnSum += int(math.Abs(float64(goalCoord.x-x)) + math.Abs(float64(goalCoord.y-y)))
			}
		}
	}

	return returnSum
}

func (s *NPuzzleState) getH() int {
	if s.currentH != -1 {
		return s.currentH
	} else {
		s.currentH = s.getManhattanDistanceScore()
		return s.currentH
	}
}

func (s *NPuzzleState) getExpectedCost() int {
	return s.cost + s.getH()
}

func (s *NPuzzleState) getStateIdentifier() string {
	if s.stateIdentifierCreated {

	} else {
		stateId := ""
		for y := 0; y < s.nSize; y++ {
			subStrings := make([]string, 0)
			for x := 0; x < s.nSize; x++ {
				subStrings = append(subStrings, string(s.puzzleState[y][x]))
			}
			stateId += strings.Join(subStrings, "")
		}
		s.stateIdentifier = stateId
		s.stateIdentifierCreated = true
	}
	return s.stateIdentifier
}

func (s *NPuzzleState) isGoal() bool {
	return s.getH() == 0
}

func describe(i interface{}) {
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

	for y := 0; y < s.nSize; y++ {
		for x := 0; x < s.nSize; x++ {
			new_coord := coord{
				x: x,
				y: y,
			}
			(*s.goalDict)[(*s.goalState)[y][x]] = new_coord
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

func insertAt(myArray *[]*interface{}, item *interface{}, index int)*[]*interface{} {

	actingArray:=*myArray
	actingArray = append(actingArray, nil)
	copy(actingArray[index+1:], actingArray[index:])
	actingArray[index] = item
	return &actingArray
}

func testInsert() {
	var coord0, coord1, coord2, coord3 interface{}
	coord0 = coord{
		x: 0,
		y: 0,
	}
	coord1 = coord{
		x: 1,
		y: 1,
	}
	coord2 = coord{
		x: 2,
		y: 2,
	}
	coord3 = coord{
		x: 3,
		y: 3,
	}

	coordArray := make([]*interface{}, 0)
	coordArray = append(coordArray, &coord0, &coord1, &coord3)
	//coordArray = append(coordArray, nil)
	//copy(coordArray[2+1:], coordArray[2:])
	//coordArray[2] = &coord2

	describe(coordArray)

	for _, val := range coordArray {
		describe(val)
	}

	coordAddr := &coordArray
	coordArray = *insertAt(coordAddr, &coord2, 2)

	describe(coordArray)
	for _, val := range coordArray {
		describe(val)
	}
}

func (solver *Solver) solve(startState *SequentialInterface, greedy bool) *SequentialInterface {

	//frontier_queue:=make([]*SequentialInterface,0)

	panic("implement me")
}

func main() {

	testInsert()
	nSize := 3
	var startState SequentialInterface

	startState = createStartState(nSize, 10)

	describe(startState)

	mySolver := createSolver()

	mySolver.solve(&startState, true)

	//fmt.Printf("(%v, %T)\n", goalState, goalState)

}
