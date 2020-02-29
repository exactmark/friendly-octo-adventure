package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	//for x := 0; x < 1000; x++ {
	//	if x == 27 {
	//
	//	} else {
	//		//mainBasicRun(0)
	//		mainLargeRun(x)
	//		//mainProfileRun(x)
	//	}
	//}
	for targetLength := 2; targetLength < 200; targetLength++ {
		mainCreateKnownLengthRun(0, targetLength)
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

	//seed n4,seed2,shuf250 is solving via a* in 20 seconds, 53 steps.
	rand.Seed(int64(seed))

	startState = createNPuzzleStartState(nSize, 250)

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

	if validateSolution(solvedList) {
		fmt.Printf("Found solution is valid.\n")
	} else {
		fmt.Printf("Found solution is NOT valid.\n")
	}

}

func mainLargeRun(seed int) {
	// this basic setup will solve A* in between 10 and 20 seconds and 51 steps.
	// the greedy version is non-deterministic. Tends to solve in 20ms, and around 150 steps.
	// And golang is going to be non-deterministic in both versions. I'm using
	// GoRoutines for child creation, which means that the order of moves is
	// potentially different. This explains the difference in solve times on
	// A* but does not really explain why different seeds will crash.

	fmt.Printf("Starting solver for seed %v.\n", seed)

	nSize := 4
	var startState SequentialInterface

	rand.Seed(int64(seed))

	startState = createNPuzzleStartState(nSize, 300)

	//startState = createNPuzzleStartState(nSize, 250)

	startState.(*NPuzzleState).printCurrentPuzzleState()
	startState.(*NPuzzleState).printCurrentGoalState()
	startTime := time.Now()
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

	if validateSolution(solvedList) {
		fmt.Printf("Found solution is valid.\n")
	} else {
		fmt.Printf("Found solution is NOT valid.\n")
	}
	PrintMemUsage()
	fmt.Printf("\n")

}

func mainProfileRun(seed int) {
	// this is the current standard for my profiling runs on the desktop, using
	//loop values of
	//for x := 0; x < 12; x++ {
	//	mainProfileRun(x)
	//}

	fmt.Printf("Starting solver.\n")

	nSize := 4
	var startState SequentialInterface

	rand.Seed(int64(seed))

	startState = createNPuzzleStartState(nSize, 10000)

	//startState = createNPuzzleStartState(nSize, 250)

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

	if validateSolution(solvedList) {
		fmt.Printf("Found solution is valid.\n")
	} else {
		fmt.Printf("Found solution is NOT valid.\n")
	}

	fmt.Printf("\n")

}

func mainCreateKnownLengthRun(seed int, targetSolLen int) {

	fmt.Printf("Starting solver for seed %v, target len %v.\n", seed, targetSolLen)

	nSize := 4
	var startState SequentialInterface

	rand.Seed(int64(seed))

	startState = createNPuzzleStartStateWithSolLen(nSize, targetSolLen)

	//startState = createNPuzzleStartState(nSize, 250)

	startState.(*NPuzzleState).printCurrentPuzzleState()
	startState.(*NPuzzleState).printCurrentGoalState()
	startTime := time.Now()
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

	if validateSolution(solvedList) {
		fmt.Printf("Found solution is valid.\n")
	} else {
		fmt.Printf("Found solution is NOT valid.\n")
	}
	PrintMemUsage()
	fmt.Printf("\n")

}

func validateSolution(solutionList *[]*SequentialInterface) bool {
	switch (*(*solutionList)[0]).(type) {
	case *NPuzzleState:
		// here v has type T
		//return (*(*solutionList)[len(*solutionList)-1]).(*NPuzzleState).testSolution()
		return (*(*solutionList)[0]).(*NPuzzleState).testSolution()
	default:
		// no match; here v has the same type as i
	}
	return false
}

//PrintMemUsage taken from:
//https://golangcode.com/print-the-current-memory-usage/
// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
