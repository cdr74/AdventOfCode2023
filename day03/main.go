package main

import (
	"fmt"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// the flag runTest defines which data file to read
const runTest bool = true
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------

func SolvePuzzle1() int {
	var result int = 0

	return result
}

// ---------------------------------------------------------------------------

func SolvePuzzle2(games []Game) int {
	var result int = 0

	return result
}

// ---------------------------------------------------------------------------

func main() {
	var input []string

	stopwatch := utils.NewStopwatch()
	stopwatch.Start()

	if runTest {
		input = utils.ReadDataFile(TEST_FILE)
	} else {
		input = utils.ReadDataFile(DATA_FILE)
	}

	result1 := SolvePuzzle1()
	result2 := SolvePuzzle2()
	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t\t", result1)
	fmt.Println("Result 2:\t\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
