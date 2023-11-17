package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// the flag runTest defines which data file to read
const runTest bool = true
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

func ReadDataFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a scanner.
	scanner := bufio.NewScanner(file)

	// Create an array to store the lines.
	var lines []string

	// Scan the file line by line.
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Return the array of lines.
	return lines
}

func SolvePuzzle() string {
	return "tbd"
}

func main() {
	var result string
	stopwatch := utils.NewStopwatch()
	stopwatch.Start()

	if runTest {
		ReadDataFile(TEST_FILE)
	} else {
		ReadDataFile(DATA_FILE)
	}

	result = SolvePuzzle()
	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result:\t\t\t", result)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
