package main

import "container/heap"

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
	}

	return &returnedSolver
}

func makeTrackbackArray(tailNode *SequentialInterface)*[]*SequentialInterface{
	returnArray:=make([]*SequentialInterface,0)

	thisNode:=tailNode
	for thisNode!=nil{
		returnArray= append(returnArray, thisNode)
		thisNode= (*thisNode).getParent()
	}

	return &returnArray
}

func (solver *Solver) solveAStar(startState *SequentialInterface) *[]*SequentialInterface {
	tailNode:=solver.solve(startState,false)
	return makeTrackbackArray(tailNode)
}

func (solver *Solver) solveGreedy(startState *SequentialInterface) *[]*SequentialInterface {
	tailNode:=solver.solve(startState,true)
	return makeTrackbackArray(tailNode)
}

func (solver *Solver) solve(startState *SequentialInterface, greedy bool) *SequentialInterface {
	var repeatedStates int64
	frontierQueue := make(PriorityQueue, 0)

	exploredStateCache := make(map[string]*SequentialInterface, 0)

	heap.Init(&frontierQueue)
	var priority int
	if greedy {
		priority = (*startState).getH()
	} else {
		priority = (*startState).getExpectedCost()
	}
	frontierQueue.PushSequentialInterface(startState, priority)

	childList := (*startState).getChildren()
	for _, val := range childList {
		frontierQueue.PushSequentialInterface(val, (*val).getH())
	}

	for frontierQueue.Len() > 0 {
		exploringNode := *frontierQueue.PopSequentialInterface()
		if exploringNode.isGoal() {
			if !greedy {
				//	store to solutionMemo cache
			}
			return &exploringNode
		}
		exploredStateCache[exploringNode.getStateIdentifier()] = &exploringNode
		childList := exploringNode.getChildren()
		for _, singleChild := range childList {
			_, existingState := exploredStateCache[(*singleChild).getStateIdentifier()]
			if existingState {
				repeatedStates += 1
			} else {
				var childPriority int
				if greedy {
					childPriority = (*singleChild).getH()
				} else {
					childPriority = (*singleChild).getExpectedCost()
				}
				frontierQueue.PushSequentialInterface(singleChild,childPriority)
			}
		}
	}
	panic("unable to find solution")
}

