package main

import (
	"container/heap"
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

	frontierQueue := make(PriorityQueue, 0)

	heap.Init(&frontierQueue)
	var priority int
	if greedy{
		priority=(*startState).getH()
	}	else{
		priority=(*startState).getExpectedCost()
	}
	frontierQueue.PushSequentialInterface(startState, priority)

	childList:= (*startState).getChildren()
	for _,val := range(childList){
		frontierQueue.PushSequentialInterface(val,(*val).getH())
	}

	for frontierQueue.Len() > 0 {
		item := frontierQueue.PopSequentialInterface()
		fmt.Printf("%v\n",(*item).getH())
		describe(*item)
		//fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
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