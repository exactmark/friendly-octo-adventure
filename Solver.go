package main

import (
	"container/heap"
	"fmt"
	"strings"
	"sync"
)

type SequentialInterface interface {
	lt(other *SequentialInterface) bool
	getChildren() []*SequentialInterface
	getH() int
	getExpectedCost() int
	getCurrentCost() int
	getStateIdentifier() string
	getGoalIdentifier() string
	isGoal() bool
	isValidMove(singleMove rune) bool
	makeMove(singleMove rune) bool
	shuffle(shuffleAmount int)
	getParent() *SequentialInterface
	createSequentialState(goalState interface{}, startState interface{}) *SequentialInterface
	exportCurrentState() interface{}
	exportGoalState() interface{}
	setParent(node *SequentialInterface)
	strandDeepCopy() *SequentialInterface
}

type Solver struct {
	useMemo      bool
	solutionMemo map[string]int
	memoLock     sync.Mutex
	memoSuccess  int
	cacheHit     int
	cacheMiss    int
	greedy       bool
	debugLog     bool
}

//createSolver will initialize a solver state and return a pointer to it.
func createSolver() *Solver {

	//memoQueue := make([]*SequentialInterface, 0)
	solutionMemo := make(map[string]int)

	returnedSolver := Solver{
		solutionMemo: solutionMemo,
		useMemo:      true,
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

	for i, j := 0, len(returnArray)-1; i < j; i, j = i+1, j-1 {
		returnArray[i], returnArray[j] = returnArray[j], returnArray[i]
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
	//var memoId string
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

	memoId := (*((*solutionList)[startIndex])).getStateIdentifier() + "sep" + (*((*solutionList)[endIndex])).getStateIdentifier()

	if solver.useMemo {
		solver.memoLock.Lock()
		_, ok := solver.solutionMemo[memoId]
		//solver.memoLock.Unlock()
		if ok {
			(*((*solutionList)[startIndex])).setParent(nil)
			//fmt.Printf("Sending cached\n")
			//fmt.Printf((*(*solutionList)[startIndex]).getStateIdentifier() + "\n")
			//
			//fmt.Printf((*(*solutionList)[endIndex]).getStateIdentifier() + "\n\n")

			envelopeChannel <- solutionEnvelope{
				position:     envelopePosition,
				solutionTail: ((*solutionList)[endIndex]),
			}
			solver.memoLock.Unlock()
			solver.cacheHit++
			return
		}
		solver.cacheMiss++
		solver.memoLock.Unlock()
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
	if solver.useMemo {
		trackback := makeTrackbackArray(solutionTail)
		solver.memoLock.Lock()
		solver.solutionMemo[memoId] = len(*trackback)
		solver.memoLock.Unlock()
	}
}

func (solver *Solver) greedyGuidedAStar(s *SequentialInterface) *[]*SequentialInterface {
	return solver.greedyGuidedAStarWithArgs(s, 4, 25)
}

func (solver *Solver) greedyGuidedAStarWithArgs(s *SequentialInterface, startInc int, stopInc int) *[]*SequentialInterface {

	initialSolution := solver.solve(s, true)
	currentSolution := makeTrackbackArray(initialSolution)
	//if solver.debugLog {
	fmt.Printf("initial solution is length %v\n", len(*currentSolution))
	//}
	for solver.spliceOutRepeatedLoops(initialSolution) {

	}
	currentSolution = makeTrackbackArray(initialSolution)
	if solver.debugLog {
		fmt.Printf("loopless initial solution is length %v\n", len(*currentSolution))
	}
	var endIndex int
	endIndex = len(*currentSolution)
	currentInc := startInc
	stepper := 0

	//highestCurrentInc := currentInc

	goal_state := (*(*currentSolution)[0]).getStateIdentifier()

	start_state := (*(*currentSolution)[endIndex-1]).getStateIdentifier()

	//This function is way too long.
	for currentInc < stopInc {
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

		currentSolution = makeTrackbackArray(lastNode)
		newSolutionLength := len(*currentSolution)
		if solver.debugLog {
			if !validateSolution(currentSolution) {
				for _, node := range *currentSolution {
					fmt.Printf((*node).getStateIdentifier() + "\n")
				}
				fmt.Printf("bad solution\n")
			}
			if goal_state != (*(*currentSolution)[0]).getStateIdentifier() ||
				start_state != (*(*currentSolution)[newSolutionLength-1]).getStateIdentifier() {
				for _, node := range *currentSolution {
					fmt.Printf((*node).getStateIdentifier() + "\n")
				}
				fmt.Printf("StartState: %v\n", start_state)
				fmt.Printf("GoalState: %v\n", goal_state)

				fmt.Printf("goal_state_actual: %v\n", (*(*currentSolution)[newSolutionLength-1]).getStateIdentifier())
				fmt.Printf("0state_actual: %v\n", (*(*currentSolution)[0]).getStateIdentifier())

				fmt.Printf("goalstate compare: %v\n", goal_state != (*(*currentSolution)[0]).getStateIdentifier())
				fmt.Printf("startstate compare: %v\n",
					start_state != (*(*currentSolution)[newSolutionLength-1]).getStateIdentifier())

				fmt.Printf("GoalState COMPARE: %v\n", strings.Compare(goal_state, (*(*currentSolution)[0]).getStateIdentifier()))
				fmt.Printf("GoalState COMPARE: %v\n", strings.Compare(goal_state, (*(*currentSolution)[newSolutionLength-1]).getStateIdentifier()))

				fmt.Printf("GoalState COMPARE: %v\n", strings.Compare(start_state, (*(*currentSolution)[newSolutionLength-1]).getStateIdentifier()))
				fmt.Printf("GoalState COMPARE: %v\n", strings.Compare(start_state, (*(*currentSolution)[0]).getStateIdentifier()))
				fmt.Printf("solution constants have changed\n")

			}
		}
		if newSolutionLength > oldSolutionLen {
			for _, node := range *currentSolution {
				fmt.Printf((*node).getStateIdentifier() + "\n")
			}
			fmt.Printf("Solution length has grown.\n")
		} else if newSolutionLength < oldSolutionLen {
			//fmt.Printf("Found better: sol len %v, stepper %v,currentInc %v\n", newSolutionLength, stepper, currentInc)
			currentInc = startInc
			stepper = 0
			for solver.spliceOutRepeatedLoops(lastNode) {
				currentSolution = makeTrackbackArray(lastNode)
				newSolutionLength = len(*currentSolution)
				//fmt.Printf("After splice: sol len %v, stepper %v,currentInc %v\n", newSolutionLength, stepper, currentInc)
			}
			currentSolution = makeTrackbackArray(lastNode)
			newSolutionLength = len(*currentSolution)
			stopInc = len(*currentSolution) / 2
			//stopInc = 40
			//fmt.Printf("CacheMisses/Hits: %v/ %v\n", solver.cacheMiss, solver.cacheHit)
		} else {
			//fmt.Printf("Found: sol len %v, stepper %v,currentInc %v\n", newSolutionLength, stepper, currentInc)

			if stepper == currentInc-1 {
				currentInc++
				stepper = 0
				//if currentInc > highestCurrentInc {
				//	highestCurrentInc = currentInc
				//	fmt.Printf("new highestCurrentInc %v\n", highestCurrentInc)
				//}
				//fmt.Printf("New currentInc: %v\n",currentInc)
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

//spliceOutRepeatedLoops will currently go through each state, find the first repeat
// and point the first node to the child of the second node.
// This lends itself to removing very simple useless moves where a state is tried then
// immediately backed out, i.e. **ABA***.
// However, consider the case of **ABAC******BDEF*** . In this case, the FirstA will
// now point to FirstC, however a better splice would be to point FirstB to FirstD.
func (solver *Solver) spliceOutRepeatedLoops(node *SequentialInterface) bool {

	uniqueMap := make(map[string]*SequentialInterface, 0)
	//doubles := make([][]*SequentialInterface, 0)
	inspectedNode := (*node).getParent()
	for inspectedNode != nil {
		earlierNode, ok := uniqueMap[(*inspectedNode).getStateIdentifier()]
		if ok {
			//fmt.Printf("found repeat\n")
			(*earlierNode).setParent((*inspectedNode).getParent())
			return true
		} else {
			uniqueMap[(*inspectedNode).getStateIdentifier()] = inspectedNode
		}
		inspectedNode = (*inspectedNode).getParent()
	}
	return false

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
