package main

import (
	"container/heap"
	"fmt"
	"sync"
)

type SequentialInterface interface {
	lt(other *SequentialInterface) bool
	getChildren() []*SequentialInterface
	getH() int
	getExpectedCost() int
	getStateIdentifier() string
	getGoalIdentifier() string
	isGoal() bool
	isValidMove(singleMove rune) bool
	makeMove(singleMove rune) bool
	shuffle(shuffleAmount int)
	getParent() *SequentialInterface
	createSequentialState(goalState interface{}, startState interface{}) *SequentialInterface
	exportCurrentState() interface{}
	setParent(node *SequentialInterface)
	strandDeepCopy() *SequentialInterface
}

type Solver struct {
	solutionMemo map[string]*SequentialInterface
	memoQueue    []*SequentialInterface
	memoLock     sync.Mutex
	memoSuccess  int
	greedy       bool
	debugLog     bool
}

//createSolver will initialize a solver state and return a pointer to it.
func createSolver() *Solver {

	memoQueue := make([]*SequentialInterface, 0)
	solutionMemo := make(map[string]*SequentialInterface)

	returnedSolver := Solver{
		solutionMemo: solutionMemo,
		memoQueue:    memoQueue,
		memoSuccess:  0,
		debugLog:     false,
	}

	return &returnedSolver
}

//makeTrackbackArray will take a given SequentialInterface or nil, and
//return a slice of the sequential states leading to the goal.
func makeTrackbackArray(tailNode *SequentialInterface) *[]*SequentialInterface {
	returnArray := make([]*SequentialInterface, 0)

	thisNode := tailNode
	for thisNode != nil {
		returnArray = append(returnArray, thisNode)
		thisNode = (*thisNode).getParent()
	}

	return &returnArray
}

//solveAStar will attempt to find a solution using the solve algorithm
//below, continuing until a state returns isGoal. The priorityQueue will be
//arranged by getExpectedCost. solveAStar will return a slice of the Sequential
//states in order of the solution.
func (solver *Solver) solveAStar(startState *SequentialInterface) *[]*SequentialInterface {
	tailNode := solver.solve(startState, false)
	return makeTrackbackArray(tailNode)
}

//solveGreedy will attempt to find a solution using the solve algorithm
//below, continuing until a state returns isGoal. The priorityQueue will be
//arranged by getH. solveGreedy will return a slice of the Sequential
//states in order of the solution.
func (solver *Solver) solveGreedy(startState *SequentialInterface) *[]*SequentialInterface {
	tailNode := solver.solve(startState, true)
	return makeTrackbackArray(tailNode)
}

//solve will attempt to find a series of SequentialInterface steps to reach a
//state that satisfies isGoal. The function will return a pointer to the goalState
//which can be used to track backwards, or nil.
func (solver *Solver) solve(startState *SequentialInterface, greedy bool) *SequentialInterface {
	var repeatedStates int64
	frontierQueue := make(PriorityQueue, 0)

	exploredStateCache := make(map[string]*SequentialInterface, 0)
	var memoId string
	heap.Init(&frontierQueue)
	var priority int
	if greedy {
		priority = (*startState).getH()
	} else {
		priority = (*startState).getExpectedCost()
		memoId = (*startState).getStateIdentifier() + "sep" + (*startState).getGoalIdentifier()

		solver.memoLock.Lock()
		val, ok := solver.solutionMemo[memoId]
		solver.memoLock.Unlock()
		if ok {
			return (*val).strandDeepCopy()
		}
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
				//fmt.Printf(memoId+"\n")

				solver.memoLock.Lock()
				solver.solutionMemo[memoId] = exploringNode.strandDeepCopy()
				solver.memoLock.Unlock()
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
				frontierQueue.PushSequentialInterface(singleChild, childPriority)
			}
		}
	}
	//(*startState).(*NPuzzleState).printCurrentPuzzleState()
	//(*startState).(*NPuzzleState).printCurrentGoalState()
	panic("unable to find solution")
}

type solutionEnvelope struct {
	position     int
	solutionTail *SequentialInterface
}

func (solver *Solver) findSolutionPart(solutionList *[]*SequentialInterface, envelopePosition int, startIndex int, endIndex int, envelopeChannel chan solutionEnvelope) {
	if endIndex >= len(*solutionList) {
		endIndex = len(*solutionList) - 1
	}
	if startIndex < 0 {
		startIndex = 0
	}

	if solver.debugLog {
		fmt.Printf("Finding solution for:\n pos %v\n start: %v\n end: %v\n", envelopePosition, startIndex, endIndex)
	}
	partialSolutionDNA := (*((*solutionList)[startIndex])).createSequentialState((*((*solutionList)[endIndex])).exportCurrentState(), (*((*solutionList)[startIndex])).exportCurrentState())

	solutionTail := solver.solve(partialSolutionDNA, false)
	envelopeChannel <- solutionEnvelope{
		position:     envelopePosition,
		solutionTail: solutionTail,
	}
}

func (solver *Solver) greedyGuidedAStar(s *SequentialInterface) *[]*SequentialInterface {
	return solver.greedyGuidedAStarWithArgs(s, 9, 25)
}

func (solver *Solver) greedyGuidedAStarWithArgs(s *SequentialInterface, startInc int, lastInc int) *[]*SequentialInterface {
	currentSolution := makeTrackbackArray(solver.solve(s, true))
	if solver.debugLog {
		fmt.Printf("initial solution is length %v\n", len(*currentSolution))
	}
	var endIndex int
	endIndex = len(*currentSolution)
	currentInc := startInc
	stepper := 0

	//This function is way too long.
	for currentInc < lastInc {
		//lastInc=len(*currentSolution)
		oldSolutionLen := len(*currentSolution)
		solutionPart := 0
		envelopeChannel := make(chan solutionEnvelope, 0)
		for singleIndex := 0 - stepper; singleIndex < endIndex; singleIndex += currentInc {
			//fmt.Printf("Solve %v, %v\n",singleIndex, singleIndex+currentInc-1)
			go solver.findSolutionPart(currentSolution, solutionPart, singleIndex, singleIndex+currentInc-1, envelopeChannel)
			solutionPart++
		}

		envelopeHeap := make(PriorityQueue, 0)

		heap.Init(&envelopeHeap)

		for returnedEnvelopes := 0; returnedEnvelopes < solutionPart; returnedEnvelopes++ {
			thisEnvelope := <-envelopeChannel
			envelopeHeap.PushSequentialInterface(thisEnvelope.solutionTail, thisEnvelope.position)
			if solver.debugLog {
				fmt.Printf("got part %v\n", thisEnvelope.position)
			}
		}

		var lastNode *SequentialInterface

		lastNode = nil

		for len(envelopeHeap) > 0 {
			thisNode := envelopeHeap.PopSequentialInterface()
			pointHeadAtLastNode(thisNode, lastNode)
			if solver.debugLog {
				fmt.Printf("Last node is \n")
				if lastNode == nil {
					fmt.Printf("nil\n")
				} else {
					(*lastNode).(*NPuzzleState).printCurrentPuzzleState()
				}
			}
			lastNode = &(*thisNode)
		}

		if solver.debugLog {
			fmt.Printf("found parts...\n")
		}

		//solver.spliceOutRepeatedLoops(lastNode)

		currentSolution = makeTrackbackArray(lastNode)
		newSolutionLength := len(*currentSolution)
		if newSolutionLength < oldSolutionLen {
			fmt.Printf("Found better: sol len %v, stepper %v,currentInc %v\n", newSolutionLength, stepper, currentInc)
			currentInc = startInc
			stepper = 0
		} else {
			fmt.Printf("Found: sol len %v, stepper %v,currentInc %v\n", newSolutionLength, stepper, currentInc)

			if stepper == currentInc-1 {
				currentInc++
				stepper = 0
			} else {
				stepper++
			}
		}

		//fmt.Printf("new solution len %v\n", len(*currentSolution))
		if solver.debugLog {
			for _, val := range *currentSolution {
				(*val).(*NPuzzleState).printCurrentPuzzleState()
			}
		}

		endIndex = len(*currentSolution)

	}
	if solver.debugLog {
		fmt.Printf("last solution len %v\n", len(*currentSolution))
	}
	return currentSolution
}

func (solver *Solver) spliceOutRepeatedLoops(node *SequentialInterface) {

	uniqueMap := make(map[string]*SequentialInterface, 0)
	doubles := make([][]*SequentialInterface, 0)
	nextNode := (*node).getParent()
	for nextNode != nil {
		laterNode, ok := uniqueMap[(*nextNode).getStateIdentifier()]
		if ok {
			fmt.Printf("found repeat\n")
			thisPair := make([]*SequentialInterface, 0)
			thisPair = append(thisPair, laterNode)
			thisPair = append(thisPair, nextNode)
			doubles = append(doubles, thisPair)
			//	todo do pointer repoint. Maybe stick in loop until no repeats found
		} else {
			uniqueMap[(*nextNode).getStateIdentifier()] = nextNode
		}
		nextNode = (*nextNode).getParent()
	}

}

func pointHeadAtLastNode(thisNode *SequentialInterface, lastNode *SequentialInterface) {
	if lastNode == nil {

	} else {
		inspectedNode := thisNode
		for (*inspectedNode).getParent() != nil {
			inspectedNode = (*inspectedNode).getParent()
		}
		(*inspectedNode).setParent(lastNode)
	}
}
