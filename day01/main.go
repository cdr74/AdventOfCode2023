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

func SolvePuzzle() string {
	return "tbd"
}

func main() {
	var result string
	var input []string

	stopwatch := utils.NewStopwatch()
	stopwatch.Start()

	if runTest {
		input = utils.ReadDataFile(TEST_FILE)
	} else {
		input = utils.ReadDataFile(DATA_FILE)
	}

	result = SolvePuzzle()
	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result:\t\t\t", result)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
